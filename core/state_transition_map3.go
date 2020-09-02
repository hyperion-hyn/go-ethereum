package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/staking/network"
	"github.com/ethereum/go-ethereum/staking/types/microstaking"
	"math/big"
)

func (st *StateTransition) verifyAndApplyCreateMap3NodeTx(msg *microstaking.CreateMap3Node, signer common.Address) error {
	epoch, blockNum := st.evm.EpochNumber, st.evm.BlockNumber
	newNode, err := VerifyCreateMap3NodeMsg(st.state, st.bc, epoch, blockNum, msg, signer)
	if err != nil {
		return err
	}
	saveNewMap3NodeToPool(newNode, st.state.Map3NodePool())
	st.state.SubBalance(signer, msg.Amount)
	return nil
}

func (st *StateTransition) verifyAndApplyEditMap3NodeTx(msg *microstaking.EditMap3Node, signer common.Address) error {
	if err := VerifyEditMap3NodeMsg(st.state, st.evm.EpochNumber, st.evm.BlockNumber, msg, signer); err != nil {
		return err
	}
	nodePool := st.state.Map3NodePool()
	wrapper, _ := st.state.Map3NodeByAddress(msg.Map3NodeAddress)
	updateMap3NodeFromPoolByMsg(wrapper, nodePool, msg)
	return nil
}

func (st *StateTransition) verifyAndApplyTerminateMap3NodeTx(msg *microstaking.TerminateMap3Node, signer common.Address) error {
	if err := VerifyTerminateMap3NodeMsg(st.state, st.evm.EpochNumber, msg, signer); err != nil {
		return err
	}
	node, _ := st.state.Map3NodeByAddress(msg.Map3NodeAddress)
	allDelegators := node.Microdelegations().AllKeys()
	if err := ReleaseMicrodelegationFromMap3Node(st.state, node, allDelegators); err != nil {
		return err
	}
	// update node state
	node.Map3Node().Status().SetValue(uint8(microstaking.Terminated))
	node.Map3Node().ActivationEpoch().SetValue(common.Big0)
	node.Map3Node().ReleaseEpoch().SetValue(common.ZeroDec())
	return nil
}

func (st *StateTransition) verifyAndApplyMicrodelegateTx(msg *microstaking.Microdelegate, signer common.Address) error {
	blockNum, epoch := st.evm.BlockNumber, st.evm.EpochNumber
	if err := VerifyMicrodelegateMsg(st.state, st.bc, blockNum, msg, signer); err != nil {
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

func (st *StateTransition) verifyAndApplyUnmicrodelegateTx(msg *microstaking.Unmicrodelegate, signer common.Address) error {
	blockNum, epoch := st.evm.BlockNumber, st.evm.EpochNumber
	if err := VerifyUnmicrodelegateMsg(st.state, st.bc, blockNum, epoch, msg, signer); err != nil {
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

func (st *StateTransition) verifyAndApplyCollectMicrodelRewardsTx(msg *microstaking.CollectRewards,
	signer common.Address) (*big.Int, error) {
	if err := VerifyCollectMicrodelRewardsMsg(st.state, msg, signer); err != nil {
		return network.NoReward, err
	}
	map3NodePool := st.state.Map3NodePool()
	return payoutMicrodelegationRewards(st.state, map3NodePool, msg.DelegatorAddress)
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
	map3Node.Map3Node().Description().UpdateDescription(msg.Description)

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
			if r.Cmp(common.Big0) > 0 {
				totalRewards = totalRewards.Add(totalRewards, r)
				micro.Reward().SetValue(common.Big0)
			}
		} else {
			return network.NoReward, errMicrodelegationNotExist
		}
	}
	stateDB.AddBalance(delegator, totalRewards)
	return totalRewards, nil
}

func ReleaseMicrodelegationFromMap3Node(stateDB vm.StateDB, node *microstaking.Storage_Map3NodeWrapper_,
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
