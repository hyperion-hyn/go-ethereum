package microstaking

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	migrationAddAgeAddress = []common.Address{
		common.ParseAddr(""),
	}
)

//todo need double check code
func MigrateNodeAge(chain consensus.ChainReader, header *types.Header, stateDB *state.StateDB, activeNode *Storage_Map3NodeWrapper_) error {
	atlasConfig := chain.Config().Atlas
	// 判断当前的块是否是在允许的高度内（迁移高度+期限）
	if header.Number.Uint64() > (atlasConfig.MigrationBlock + atlasConfig.MigrationDuration) {
		return nil
	}
	// 判断当前节点的operatorAddress属于需要给予增加节龄的地址
	operator := activeNode.Map3Node().OperatorAddress().Value()
	isContain := false
	for _, addressTemp := range migrationAddAgeAddress {
		if operator == addressTemp {
			isContain = true
		}
	}
	if !isContain {
		return nil
	}
	// 判断当前map3节点是该地址的第一个节点

	// 取回该节点抵押过的所有map3
	delegationIndexMap := stateDB.Map3NodePool().DelegationIndexMapByDelegator().Get(operator)
	var createMap3Addrs = make([]common.Address, 0)
	for i := 0; i < delegationIndexMap.Keys().Length(); i++ {
		nodeAddr := delegationIndexMap.Keys().Get(i).Value()
		node, err := stateDB.Map3NodeByAddress(nodeAddr)
		if err != nil {
			return err
		}
		nodeOperator := node.Map3Node().OperatorAddress().Value()
		if operator == nodeOperator {
			createMap3Addrs = append(createMap3Addrs, nodeAddr)
		}
	}
	// 判断创建的节点数是否>1
	if len(createMap3Addrs) > 1 {
		return nil
	}

	if activeNode.Map3Node().Map3Address().Value() != createMap3Addrs[0] {
		return errors.New("active node address not equal")
	}
	// 给当前节点增加节龄（180）
	nodeAge := activeNode.Map3Node().Age().Value()
	activeNode.Map3Node().Age().SetValue(nodeAge.Add(common.NewDec(180)))
	return nil

}
