package core

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/staking/types/microstaking"
	"math/big"
)

var (
	operatorsToAddNodeAge = map[common.Address]bool{
		common.HexToAddress("0x42B66C58c1D1bf304c37398Ccb89963f7d1e3E38"): true,
		common.HexToAddress("0xE72fcA6d8945f5805F537835f850f1945DCfb72a"): true,
		common.HexToAddress("0x697097A6fD21c5F254dadb0d5aC8bC3d24F483aD"): true,
	}
)

func increaseMap3NodeAgeOnDemand(node *microstaking.Map3NodeWrapper_, blockNum *big.Int, stateDB vm.StateDB, chain ChainContext) error {
	config := chain.Config().Atlas
	// over deadline
	if blockNum.Cmp(big.NewInt(int64(config.Map3NodeAgeDeadlineBlock))) > 0 {
		return nil
	}

	// check if the first node to be created
	operator := node.Map3Node.OperatorAddress
	if _, ok := operatorsToAddNodeAge[operator]; !ok {
		return nil
	}
	delegationIndexMap := stateDB.Map3NodePool().DelegationIndexMapByDelegator().Get(operator)
	for i := 0; i < delegationIndexMap.Keys().Length(); i++ {
		nodeAddr := delegationIndexMap.Keys().Get(i).Value()
		index, ok := delegationIndexMap.Get(nodeAddr)
		if !ok {
			return errors.New("delegation index not found")
		}
		if index.IsOperator().Value() {
			return nil
		}
	}
	node.Map3Node.Age = common.NewDec(180)
	return nil
}