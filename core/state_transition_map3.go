package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/staking/network"
	"github.com/ethereum/go-ethereum/staking/types/microstaking"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"math/big"
)

func (st *StateTransition) verifyAndApplyCreateMap3NodeTx(verifier StakingVerifier, msg *microstaking.CreateMap3Node, signer common.Address) error {
	epoch, blockNum := st.evm.EpochNumber, st.evm.BlockNumber
	newNode, err := verifier.VerifyCreateMap3NodeMsg(st.state, st.bc, epoch, blockNum, msg, signer)
	if err != nil {
		return err
	}
	saveNewMap3NodeToPool(newNode, st.state.Map3NodePool())
	st.state.SubBalance(signer, msg.Amount)
	return nil
}

func (st *StateTransition) verifyAndApplyEditMap3NodeTx(verifier StakingVerifier, msg *microstaking.EditMap3Node, signer common.Address) error {
	if err := verifier.VerifyEditMap3NodeMsg(st.state, st.bc, st.evm.EpochNumber, st.evm.BlockNumber, msg, signer); err != nil {
		return err
	}
	nodePool := st.state.Map3NodePool()
	wrapper, _ := st.state.Map3NodeByAddress(msg.Map3NodeAddress)
	updateMap3NodeFromPoolByMsg(wrapper, nodePool, msg)
	return nil
}

func (st *StateTransition) verifyAndApplyTerminateMap3NodeTx(verifier StakingVerifier, msg *microstaking.TerminateMap3Node, signer common.Address) error {
	if err := verifier.VerifyTerminateMap3NodeMsg(st.state, st.evm.EpochNumber, msg, signer); err != nil {
		return err
	}
	node, _ := st.state.Map3NodeByAddress(msg.Map3NodeAddress)
	allDelegators := node.Microdelegations().AllKeys()
	if err := releaseMicrodelegationFromMap3Node(st.state, node, allDelegators); err != nil {
		return err
	}
	node.Terminate()
	return nil
}

func (st *StateTransition) verifyAndApplyMicrodelegateTx(verifier StakingVerifier, msg *microstaking.Microdelegate, signer common.Address) error {
	blockNum, epoch := st.evm.BlockNumber, st.evm.EpochNumber
	if err := verifier.VerifyMicrodelegateMsg(st.state, st.bc, blockNum, msg, signer); err != nil {
		return err
	}
	st.state.SubBalance(msg.DelegatorAddress, msg.Amount)
	node, _ := st.state.Map3NodeByAddress(msg.Map3NodeAddress)

	// TODO(ATLAS): only allow in pending status
	isNewDelegator := node.AddMicrodelegation(msg.DelegatorAddress, msg.Amount, true, epoch)
	if isNewDelegator {
		st.state.Map3NodePool().UpdateDelegationIndex(msg.DelegatorAddress, &microstaking.DelegationIndex_{
			Map3Address: msg.Map3NodeAddress,
			IsOperator:  node.IsOperator(msg.DelegatorAddress),
		})
	}
	return nil
}

func (st *StateTransition) verifyAndApplyUnmicrodelegateTx(verifier StakingVerifier, msg *microstaking.Unmicrodelegate, signer common.Address) error {
	blockNum, epoch := st.evm.BlockNumber, st.evm.EpochNumber
	if err := verifier.VerifyUnmicrodelegateMsg(st.state, st.bc, blockNum, epoch, msg, signer); err != nil {
		return err
	}
	node, _ := st.state.Map3NodeByAddress(msg.Map3NodeAddress)
	toBalance, completed := node.Unmicrodelegate(msg.DelegatorAddress, msg.Amount)
	if completed {
		st.state.Map3NodePool().RemoveDelegationIndex(msg.DelegatorAddress, msg.Map3NodeAddress)
	}
	st.state.AddBalance(msg.DelegatorAddress, toBalance)
	return nil
}

func (st *StateTransition) verifyAndApplyCollectMicrodelRewardsTx(verifier StakingVerifier, msg *microstaking.CollectRewards, signer common.Address) (*big.Int, error) {
	if err := verifier.VerifyCollectMicrostakingRewardsMsg(st.state, msg, signer); err != nil {
		return network.NoReward, err
	}
	map3NodePool := st.state.Map3NodePool()
	return payoutMicrodelegationRewards(st.state, map3NodePool, msg.DelegatorAddress)
}

func (st *StateTransition) verifyAndApplyRenewMap3NodeTx(verifier StakingVerifier, msg *microstaking.RenewMap3Node, signer common.Address) error {
	blockNum, epoch := st.evm.BlockNumber, st.evm.EpochNumber
	if err := verifier.VerifyRenewMap3NodeMsg(st.state, st.bc, blockNum, epoch, msg, signer); err != nil {
		return err
	}

	node, _ := st.state.Map3NodeByAddress(msg.Map3NodeAddress)
	md, _ := node.Microdelegations().Get(msg.DelegatorAddress)
	status := microstaking.NotRenewed
	if msg.IsRenew {
		status = microstaking.Renewed
	}
	md.Renewal().Save(&microstaking.Renewal_{
		Status:       uint8(status),
		UpdateHeight: st.evm.BlockNumber,
	})

	if !msg.NewCommissionRate.IsNil() {
		node.Map3Node().Commission().RateForNextPeriod().SetValue(msg.NewCommissionRate)
	}
	return nil
}

/**
 * save the new map3 node into pool
 */
func saveNewMap3NodeToPool(wrapper *microstaking.Map3NodeWrapper_, map3NodePool *microstaking.Storage_Map3NodePool_) {
	map3NodePool.Nodes().Put(wrapper.Map3Node.Map3Address, wrapper)
	keySet := map3NodePool.NodeKeySet()
	for _, key := range wrapper.Map3Node.NodeKeys.Keys {
		keySet.Get(key.Hex()).SetValue(true)
	}
	if identity := wrapper.Map3Node.Description.Identity; identity != "" {
		map3NodePool.DescriptionIdentitySet().Get(identity).SetValue(true)
	}
	map3Address, operator := wrapper.Map3Node.Map3Address, wrapper.Map3Node.OperatorAddress
	map3NodePool.UpdateDelegationIndex(operator, &microstaking.DelegationIndex_{
		Map3Address: map3Address,
		IsOperator:  true,
	})
}

func updateMap3NodeFromPoolByMsg(map3Node *microstaking.Storage_Map3NodeWrapper_, pool *microstaking.Storage_Map3NodePool_,
	msg *microstaking.EditMap3Node) {
	// update description
	if msg.Description.Identity != "" {
		i := map3Node.Map3Node().Description().Identity().Value()
		pool.DescriptionIdentitySet().Get(i).SetValue(false)
		pool.DescriptionIdentitySet().Get(msg.Description.Identity).SetValue(true)
	}
	map3Node.Map3Node().Description().IncrementalUpdateFrom(msg.Description)

	if msg.NodeKeyToRemove != nil {
		for i := 0; i < map3Node.Map3Node().NodeKeys().Length(); i++ {
			if map3Node.Map3Node().NodeKeys().Get(i).Equal(msg.NodeKeyToRemove) {
				map3Node.Map3Node().NodeKeys().Remove(i)
				pool.NodeKeySet().Get(msg.NodeKeyToRemove.Hex()).SetValue(false)
				break
			}
		}
	}

	if msg.NodeKeyToAdd != nil {
		map3Node.Map3Node().NodeKeys().Push(msg.NodeKeyToAdd)
		pool.NodeKeySet().Get(msg.NodeKeyToAdd.Hex()).SetValue(true)
	}
}

func payoutMicrodelegationRewards(stateDB vm.StateDB, map3NodePool *microstaking.Storage_Map3NodePool_, delegator common.Address) (*big.Int, error) {
	delegationIndexMap := map3NodePool.DelegationIndexMapByDelegator().Get(delegator)
	totalRewards := big.NewInt(0)
	for i := 0; i < delegationIndexMap.Keys().Length(); i++ {
		nodeAddr := delegationIndexMap.Keys().Get(i).Value()
		node, err := stateDB.Map3NodeByAddress(nodeAddr)
		if err != nil {
			return network.NoReward, err
		}

		if micro, ok := node.Microdelegations().Get(delegator); ok {
			r := micro.Reward().Value()
			if r.Sign() > 0 {
				totalRewards = totalRewards.Add(totalRewards, r)
				micro.Reward().Clear()
			}
		} else {
			return network.NoReward, errMicrodelegationNotExist
		}
	}
	stateDB.AddBalance(delegator, totalRewards)
	return totalRewards, nil
}

// TODO(ATLAS): terminate and Release delegation?
func releaseMicrodelegationFromMap3Node(stateDB vm.StateDB, node *microstaking.Storage_Map3NodeWrapper_,
	delegatorsToBeReleased []common.Address) error {
	totalToReduce, totalPendingToReduce := big.NewInt(0), big.NewInt(0)
	for _, delegator := range delegatorsToBeReleased {
		md, ok := node.Microdelegations().Get(delegator)
		if !ok {
			return errMicrodelegationNotExist
		}

		balance := big.NewInt(0).Add(md.Amount().Value(), md.Reward().Value())
		balance.Add(balance, md.PendingDelegation().Amount().Value())
		stateDB.AddBalance(delegator, balance)

		totalToReduce = totalToReduce.Add(totalToReduce, md.Amount().Value())
		totalPendingToReduce = totalPendingToReduce.Add(totalPendingToReduce, md.PendingDelegation().Amount().Value())
	}

	for _, delegator := range delegatorsToBeReleased {
		node.Microdelegations().Remove(delegator)
		delegationIndexMap := stateDB.Map3NodePool().DelegationIndexMapByDelegator().Get(delegator)
		delegationIndexMap.Remove(node.Map3Node().Map3Address().Value())
	}

	node.SubTotalDelegation(totalToReduce)
	node.SubTotalPendingDelegation(totalPendingToReduce)
	return nil
}

func LookupMicrodelegationShares(node *microstaking.Storage_Map3NodeWrapper_) (map[common.Address]common.Dec, error) {
	shares := map[common.Address]common.Dec{}
	totalDelegationDec := common.NewDecFromBigInt(node.TotalDelegation().Value())
	if totalDelegationDec.IsZero() {
		log.Info("zero total delegation during AddReward delegation payout",
			"validator-snapshot", node.Map3Node().Map3Address().Value().Hex())
		return shares, nil
	}

	for _, key := range node.Microdelegations().AllKeys() {
		delegation, ok := node.Microdelegations().Get(key)
		if !ok {
			return nil, errMicrodelegationNotExist
		}
		percentage := common.NewDecFromBigInt(delegation.Amount().Value()).Quo(totalDelegationDec)
		shares[key] = percentage
	}

	// TODO(ATLAS): cache shares

	return shares, nil
}

type map3NodeAsParticipant struct {
	stateDB vm.StateDB
	node    *microstaking.Storage_Map3NodeWrapper_
}

func (p map3NodeAsParticipant) restakingAmount() *big.Int {
	return p.node.TotalDelegation().Value()
}

func (p map3NodeAsParticipant) postCreateValidator(validator common.Address, amount *big.Int) error {
	p.node.RestakingReference().ValidatorAddress().SetValue(validator)
	return nil
}

func (p map3NodeAsParticipant) postRedelegate(validator common.Address, amount *big.Int) error {
	p.node.RestakingReference().ValidatorAddress().SetValue(validator)
	return nil
}

func (p map3NodeAsParticipant) rewardHandler() RestakingRewardHandler {
	return &RewardToMap3Node{}
}

type RewardToMap3Node struct {
	StateDB vm.StateDB
}

func (handler RewardToMap3Node) HandleReward(redelegation *restaking.Storage_Redelegation_, epoch *big.Int) (*big.Int, error) {
	reward := redelegation.Reward().Value()
	if reward.Sign() == 0 {
		return common.Big0, nil
	}
	map3Address := redelegation.DelegatorAddress().Value()
	// TODO(ATLAS): can not continue to delegate after activating the map3 node
	// calculate shares based on the latest delegation state
	node, err := handler.StateDB.Map3NodeByAddress(map3Address)
	if err != nil {
		return nil, err
	}
	shares, err := LookupMicrodelegationShares(node)
	if err != nil {
		return nil, err
	}
	if err := handler.StateDB.AddMicrodelegationReward(node, reward, shares); err != nil {
		return nil, err
	}
	redelegation.Reward().Clear()
	return reward, nil
}
