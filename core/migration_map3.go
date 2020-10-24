package core

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/staking/types/microstaking"
	"github.com/pkg/errors"
	"math/big"
)

var (
	operatorsToAddNodeAge = map[common.Address]bool{
		common.HexToAddress("0x42B66C58c1D1bf304c37398Ccb89963f7d1e3E38"): true,
		common.HexToAddress("0xE72fcA6d8945f5805F537835f850f1945DCfb72a"): true,
		common.HexToAddress("0x697097A6fD21c5F254dadb0d5aC8bC3d24F483aD"): true,
	}

	map3NodesToBeMigrated = "" // TODO(ATLAS): get from ethereum staking contract
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

func MigrateMap3NodesFromEthereum(chain ChainContext, stateDB vm.StateDB, blockNum *big.Int) error {
	config := chain.Config().Atlas
	if blockNum.Cmp(big.NewInt(int64(config.Map3MigrationBlock))) != 0 {
		return nil
	}

	// parse map3 nodes from string
	var ns []microstaking.PlainMap3NodeWrapper
	if err := json.Unmarshal([]byte(map3NodesToBeMigrated), &ns); err != nil {
		return errors.Wrap(err, "failed to parse map3 nodes to be migrated")
	}

	pool := stateDB.Map3NodePool()
	for _, n := range ns {
		saveNewMap3NodeToPool(n.ToMap3NodeWrapper(), pool)
	}
	return nil
}