package rawdb

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

//Write active map3 node address at epoch to rawdb
func WriteMutateMap3Nodes(db DatabaseWriter, epoch uint64, mutateAddrs []common.Address) {
	// Convert the receipts into their mutateAddrs form and serialize them
	bytes, err := rlp.EncodeToBytes(mutateAddrs)
	if err != nil {
		log.Error("Fail to encode mutateMap3Nodes")
	}
	// Store the flattened receipt slice
	if err := db.Put(mutateMap3Key(epoch), bytes); err != nil {
		log.Error("Fail to store mutateMap3Nodes")
	}
}

func ReadMutateMap3Nodes(db DatabaseReader, epoch uint64) []common.Address {
	data, _ := db.Get(mutateMap3Key(epoch))
	if len(data) == 0 {
		return nil
	}
	var mutateMap3Nodes []common.Address
	if err := rlp.DecodeBytes(data, &mutateMap3Nodes); err != nil {
		log.Error("Invalid mutateMap3Nodes RLP")
		return nil
	}
	return mutateMap3Nodes
}

func WriteRenewedMap3Nodes(db DatabaseWriter, epoch uint64, renewActiveMap3Nodes []common.Address) {
	// Convert the receipts into their mutateAddrs form and serialize them
	bytes, err := rlp.EncodeToBytes(renewActiveMap3Nodes)
	if err != nil {
		log.Error("Fail to encode renewedMap3Nodes")
	}
	// Store the flattened receipt slice
	if err := db.Put(renewedMap3Key(epoch), bytes); err != nil {
		log.Error("Fail to store renewedMap3Nodes")
	}
}

func ReadRenewedMap3Nodes(db DatabaseReader, epoch uint64) []common.Address {
	data, _ := db.Get(renewedMap3Key(epoch))
	if len(data) == 0 {
		return nil
	}
	var renewActiveMap3Nodes []common.Address
	if err := rlp.DecodeBytes(data, &renewActiveMap3Nodes); err != nil {
		log.Error("Invalid renewedMap3Nodes RLP")
		return nil
	}
	return renewActiveMap3Nodes
}
