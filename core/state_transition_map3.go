package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/staking/network"
	"github.com/ethereum/go-ethereum/staking/types/microstaking"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"github.com/pkg/errors"
	"math/big"
)

func (st *StateTransition) verifyAndApplyCreateMap3NodeTx(verifier StakingVerifier, msg *microstaking.CreateMap3Node, signer common.Address) error {
	epoch, blockNum := st.evm.EpochNumber, st.evm.BlockNumber
	newNode, err := verifier.VerifyCreateMap3NodeMsg(st.state, st.bc, epoch, blockNum, msg, signer)
	if err != nil {
		return err
	}
	if err := increaseMap3NodeAgeFromEthereum(newNode, blockNum, st.state, st.bc); err != nil {
		return errors.Wrap(err, "failed to increase node age")
	}
	nodeSt := saveNewMap3NodeToPool(newNode, st.state.Map3NodePool())
	st.state.SubBalance(signer, msg.Amount)

	if st.bc.Config().Atlas.IsMicrostakingImprove(blockNum) {
		requireTotal, requireSelf, _ := network.LatestMicrostakingRequirement(blockNum, st.bc.Config())
		if nodeSt.CanActivate(requireTotal, requireSelf) {
			calculator := microstaking.NewLockDurationCalculator(st.bc.Config().Atlas, blockNum)
			if err := nodeSt.Activate(epoch, blockNum, calculator); err != nil {
				return errors.Wrap(err, "failed to activate new map3 node")
			}

			// update snapshot
			newSnapshot, err := nodeSt.Load()
			if err != nil {
				return errors.Wrap(err, "failed to update snapshot")
			}
			st.state.Map3NodePool().Map3NodeSnapshots().Put(newSnapshot.Map3Node.Map3Address, newSnapshot)
		}
	}
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
	if err := releaseMicrodelegationFromMap3Node(st.bc, st.evm.BlockNumber, st.state, node, allDelegators); err != nil {
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

	if st.bc.Config().Atlas.IsMicrostakingImprove(blockNum) {
		requireTotal, requireSelf, _ := network.LatestMicrostakingRequirement(blockNum, st.bc.Config())
		if node.CanActivate(requireTotal, requireSelf) {
			calculator := microstaking.NewLockDurationCalculator(st.bc.Config().Atlas, blockNum)
			if err := node.Activate(epoch, blockNum, calculator); err != nil {
				return errors.Wrap(err, "failed to activate map3 node")
			}

			// update snapshot
			newSnapshot, err := node.Load()
			if err != nil {
				return errors.Wrap(err, "failed to update snapshot")
			}
			st.state.Map3NodePool().Map3NodeSnapshots().Put(newSnapshot.Map3Node.Map3Address, newSnapshot)
		}
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

type RewardFromMap3 struct {
	Map3NodeAddress common.Address
	RewardAmount    *big.Int
}

func (st *StateTransition) verifyAndApplyCollectMicrodelRewardsTx(verifier StakingVerifier, msg *microstaking.CollectRewards, signer common.Address) (*big.Int, error) {
	if err := verifier.VerifyCollectMicrostakingRewardsMsg(st.state, msg, signer); err != nil {
		return network.NoReward, err
	}
	totalReward, rewardNodes, err := payoutMicrodelegationRewards(st.state, st.state.Map3NodePool(), msg.DelegatorAddress)
	if err != nil {
		return network.NoReward, err
	}

	// Add log if everything is good
	rewardNodesBytes, err := rlp.EncodeToBytes(rewardNodes)
	if err != nil {
		return network.NoReward, err
	}
	st.state.AddLog(&types.Log{
		Address:     msg.DelegatorAddress,
		Topics:      []common.Hash{microstaking.CollectRewardsTopic},
		Data:        rewardNodesBytes,
		BlockNumber: st.evm.BlockNumber.Uint64(),
	})
	return totalReward, nil
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

	if msg.NewCommissionRate != nil && !msg.NewCommissionRate.IsNil() {
		node.Map3Node().Commission().RateForNextPeriod().SetValue(*msg.NewCommissionRate)
	}
	return nil
}

/**
 * save the new map3 node into pool
 */
func saveNewMap3NodeToPool(wrapper *microstaking.Map3NodeWrapper_, map3NodePool *microstaking.Storage_Map3NodePool_) *microstaking.Storage_Map3NodeWrapper_ {
	map3Address, operator := wrapper.Map3Node.Map3Address, wrapper.Map3Node.OperatorAddress
	map3NodePool.Map3Nodes().Put(map3Address, wrapper)
	map3NodePool.Map3NodeSnapshots().Put(map3Address, wrapper)
	keySet := map3NodePool.NodeKeySet()
	for _, key := range wrapper.Map3Node.NodeKeys.Keys {
		keySet.Get(key.Hex()).SetValue(true)
	}
	if identity := wrapper.Map3Node.Description.Identity; identity != "" {
		map3NodePool.DescriptionIdentitySet().Get(identity).SetValue(true)
	}
	map3NodePool.UpdateDelegationIndex(operator, &microstaking.DelegationIndex_{
		Map3Address: map3Address,
		IsOperator:  true,
	})
	nSt, ok := map3NodePool.Map3Nodes().Get(map3Address)
	if !ok {
		log.Error("new map3 node not found in pool", "node", map3Address)
	}
	return nSt
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

func payoutMicrodelegationRewards(stateDB vm.StateDB, map3NodePool *microstaking.Storage_Map3NodePool_, delegator common.Address) (*big.Int, []RewardFromMap3, error) {
	delegationIndexMap := map3NodePool.DelegationIndexMapByDelegator().Get(delegator)
	totalRewards := big.NewInt(0)
	rewardNodes := make([]RewardFromMap3, 0)
	for i := 0; i < delegationIndexMap.Keys().Length(); i++ {
		nodeAddr := delegationIndexMap.Keys().Get(i).Value()
		node, err := stateDB.Map3NodeByAddress(nodeAddr)
		if err != nil {
			return network.NoReward, nil, err
		}

		if micro, ok := node.Microdelegations().Get(delegator); ok {
			r := micro.Reward().Value()
			if r.Sign() > 0 {
				totalRewards = totalRewards.Add(totalRewards, r)
				micro.Reward().Clear()
				rewardNodes = append(rewardNodes, RewardFromMap3{
					Map3NodeAddress: nodeAddr,
					RewardAmount:    r,
				})
			}
		} else {
			return network.NoReward, nil, errMicrodelegationNotExist
		}
	}
	stateDB.AddBalance(delegator, totalRewards)
	return totalRewards, rewardNodes, nil
}

// TODO(ATLAS): terminate and Release delegation?
func releaseMicrodelegationFromMap3Node(chain ChainContext, blockNum *big.Int, stateDB vm.StateDB, node *microstaking.Storage_Map3NodeWrapper_,
	delegatorsToBeReleased []common.Address) error {
	totalToReduce, totalPendingToReduce := big.NewInt(0), big.NewInt(0)
	returnRecords := make([]microstaking.MicrostakingReturnRecord, 0)
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

		returnAmount := big.NewInt(0).Add(md.Amount().Value(), md.PendingDelegation().Amount().Value())
		returnRecord := microstaking.MicrostakingReturnRecord{
			Delegator: delegator,
			Map3Node:  node.Map3Node().Map3Address().Value(),
			Amount:    returnAmount,
			Reward:    md.Reward().Value(),
		}
		returnRecords = append(returnRecords, returnRecord)
	}

	for _, delegator := range delegatorsToBeReleased {
		node.Microdelegations().Remove(delegator)
		delegationIndexMap := stateDB.Map3NodePool().DelegationIndexMapByDelegator().Get(delegator)
		delegationIndexMap.Remove(node.Map3Node().Map3Address().Value())
	}

	node.SubTotalDelegation(totalToReduce)
	node.SubTotalPendingDelegation(totalPendingToReduce)

	// write rawdb
	batch := chain.Database().NewBatch()
	rawdb.WriteTerminateMap3ReturnRecords(batch, blockNum, node.Map3Node().Map3Address().Value(), returnRecords)
	if err := batch.Write(); err != nil {
		return err
	}

	return nil
}

func lookupMicrodelegationShares(node *microstaking.Storage_Map3NodeWrapper_) (map[common.Address]common.Dec, error) {
	shares := map[common.Address]common.Dec{}
	totalDelegationDec := common.NewDecFromBigInt(node.TotalDelegation().Value())
	if totalDelegationDec.IsZero() {
		log.Info("zero total delegation during AddReward delegation payout",
			"map3node-snapshot", node.Map3Node().Map3Address().Value().Hex())
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
	chain   ChainContext
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
	return &RewardToMap3Node{
		StateDB: p.stateDB,
		Chain:   p.chain,
	}
}

type RewardToMap3Node struct {
	StateDB vm.StateDB
	Chain   ChainContext
}

func (handler RewardToMap3Node) HandleReward(redelegation *restaking.Storage_Redelegation_, blockNum *big.Int) (*big.Int, error) {
	reward := redelegation.Reward().Value()
	if reward.Sign() == 0 {
		return common.Big0, nil
	}
	map3Address := redelegation.DelegatorAddress().Value()
	// calculate shares based on the map3 node snapshot
	lastButOneBlockNum := new(big.Int).Sub(blockNum, common.Big2)
	node, err := handler.Chain.ReadMap3NodeSnapshotAtBlock(lastButOneBlockNum, map3Address)
	if err != nil {
		return nil, err
	}
	shares, err := lookupMicrodelegationShares(node)
	if err != nil {
		return nil, err
	}
	if err := handler.StateDB.AddMicrostakingReward(node, reward, shares); err != nil {
		return nil, err
	}
	redelegation.Reward().Clear()
	return reward, nil
}
