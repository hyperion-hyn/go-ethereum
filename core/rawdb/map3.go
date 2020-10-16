package rawdb

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

//Write active map3 node address at epoch to rawdb
func WriteActiveMap3Nodes(db DatabaseWriter, epoch uint64, activeMap3Nodes []common.Address) {
	// Convert the receipts into their storage form and serialize them
	bytes, err := rlp.EncodeToBytes(activeMap3Nodes)
	if err != nil {
		log.Error("Fail to encode activeMap3Node")
	}
	// Store the flattened receipt slice
	if err := db.Put(activeMap3Key(epoch), bytes); err != nil {
		log.Error("Fail to store activeMap3Node")
	}
}

func ReadActiveMap3Nodes(db DatabaseReader, epoch uint64) []common.Address {
	data, _ := db.Get(activeMap3Key(epoch))
	if len(data) == 0 {
		return nil
	}
	var activeMap3Nodes []common.Address
	if err := rlp.DecodeBytes(data, &activeMap3Nodes); err != nil {
		log.Error("Invalid activeMap3Node RLP")
		return nil
	}
	return activeMap3Nodes
}

//Write terminated map3 node address at epoch to rawdb
func WriteTerminatedMap3Nodes(db DatabaseWriter, epoch uint64, terminatedMap3Nodes []common.Address) {
	// Convert the receipts into their storage form and serialize them
	bytes, err := rlp.EncodeToBytes(terminatedMap3Nodes)
	if err != nil {
		log.Error("Fail to encode terminatedMap3Nodes")
	}
	// Store the flattened receipt slice
	if err := db.Put(terminatedMap3Key(epoch), bytes); err != nil {
		log.Error("Fail to store terminatedMap3Nodes")
	}
}

func ReadTerminatedMap3Nodes(db DatabaseReader, epoch uint64) []common.Address {
	data, _ := db.Get(terminatedMap3Key(epoch))
	if len(data) == 0 {
		return nil
	}
	var terminatedMap3Nodes []common.Address
	if err := rlp.DecodeBytes(data, &terminatedMap3Nodes); err != nil {
		log.Error("Invalid terminatedMap3Nodes RLP")
		return nil
	}
	return terminatedMap3Nodes
}
