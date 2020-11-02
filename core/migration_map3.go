package core

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/staking/types/microstaking"
	"github.com/pkg/errors"
	"math/big"
)

var (
	foundationAddress     = common.HexToAddress("0xB9F8d03176f040c6f1A739DDF73ae8Dbde9657F6")
	operatorsToAddNodeAge = map[common.Address]bool{
		common.HexToAddress("0x42B66C58c1D1bf304c37398Ccb89963f7d1e3E38"): true,
		common.HexToAddress("0xE72fcA6d8945f5805F537835f850f1945DCfb72a"): true,
		common.HexToAddress("0x697097A6fD21c5F254dadb0d5aC8bC3d24F483aD"): true,
	}
)

func increaseMap3NodeAgeFromEthereum(node *microstaking.Map3NodeWrapper_, blockNum *big.Int, stateDB vm.StateDB, chain ChainContext) error {
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
	// migrate map3 nodes from ethereum
	if blockNum.Cmp(big.NewInt(int64(config.Map3MigrationBlock))) == 0 {

		// parse map3 nodes from string
		var ns []microstaking.PlainMap3NodeWrapper
		if err := json.Unmarshal([]byte(getMap3NodesToBeMigrated()), &ns); err != nil {
			return errors.Wrap(err, "failed to parse map3 nodes to be migrated")
		}

		// save map3 nodes into pool
		pool := stateDB.Map3NodePool()
		for _, n := range ns {
			saveNewMap3NodeToPool(n.ToMap3NodeWrapper(), pool)
		}
		log.Info("migrate eth map3 nodes", "nodes", len(ns), "height", blockNum)

		// set operator nonce to be 1
		for _, n := range ns {
			stateDB.SetNonce(n.Map3Node.OperatorAddress, 1)
		}

		// pay out from foundation account
		total := calculateTotalDelegationAndReward(ns)
		stateDB.SubBalance(foundationAddress, total)
		log.Info("take out from foundation account for eth map3 nodes", "total", total)
	}

	// repair microdelegation index for map3 nodes from ethereum
	if blockNum.Cmp(big.NewInt(int64(config.MicrodelegationIndexRepairBlock))) == 0 {
		// parse map3 nodes from string
		var ns []microstaking.PlainMap3NodeWrapper
		if err := json.Unmarshal([]byte(getMap3NodesToBeMigrated()), &ns); err != nil {
			return errors.Wrap(err, "failed to parse map3 nodes to repair index")
		}

		// repair index
		pool := stateDB.Map3NodePool()
		for _, node := range ns {
			map3Address := node.Map3Node.Map3Address
			operatorAddress := node.Map3Node.OperatorAddress
			for _, delegation := range node.Microdelegations {
				delegator := delegation.DelegatorAddress
				pool.UpdateDelegationIndex(delegator, &microstaking.DelegationIndex_{
					Map3Address: map3Address,
					IsOperator:  delegator == operatorAddress,
				})
			}
		}
		log.Info("repaired microdelegation index for map3 nodes from ethereum", "nodes", len(ns))
	}

	return nil
}

func calculateTotalDelegationAndReward(nodes []microstaking.PlainMap3NodeWrapper) *big.Int {
	total := big.NewInt(0)
	for _, node := range nodes {
		for _, microdelegation := range node.Microdelegations {
			total = total.Add(total, microdelegation.Reward)
		}
		total = total.Add(total, node.TotalDelegation)
	}
	return total
}
