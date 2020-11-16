package rawdb

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/staking/burning"
	"github.com/ethereum/go-ethereum/staking/types/microstaking"
	"math/big"
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

func WriteUnmicrodelegationReturnRecords(db DatabaseWriter, block *big.Int, unmicrodelegationReturnRecords []microstaking.MicrostakingReturnRecord) {
	bytes, err := rlp.EncodeToBytes(unmicrodelegationReturnRecords)
	if err != nil {
		log.Error("Fail to encode unmicrodelegationReturnRecords")
	}
	// Store the flattened receipt slice
	if err := db.Put(unmicrodelegationReturnRecordKey(block), bytes); err != nil {
		log.Error("Fail to store unmicrodelegationReturnRecords")
	}
}

func ReadUnmicrodelegationReturnRecords(db DatabaseReader, blockNum *big.Int) []microstaking.MicrostakingReturnRecord {
	data, _ := db.Get(unmicrodelegationReturnRecordKey(blockNum))
	if len(data) == 0 {
		return nil
	}
	var unmicrodelegationReturnRecords []microstaking.MicrostakingReturnRecord
	if err := rlp.DecodeBytes(data, &unmicrodelegationReturnRecords); err != nil {
		log.Error("Invalid unmicrodelegationReturnRecords RLP")
		return nil
	}
	return unmicrodelegationReturnRecords
}

func WriteTerminateMap3ReturnRecords(db DatabaseWriter, blockNum *big.Int, map3Node common.Address, terminateMap3ReturnRecords []microstaking.MicrostakingReturnRecord) {
	bytes, err := rlp.EncodeToBytes(terminateMap3ReturnRecords)
	if err != nil {
		log.Error("Fail to encode terminateMap3ReturnRecords")
	}
	// Store the flattened receipt slice
	if err := db.Put(terminateMap3ReturnRecordKey(blockNum, map3Node), bytes); err != nil {
		log.Error("Fail to store terminateMap3ReturnRecords")
	}
}

func ReadTerminateMap3ReturnRecords(db DatabaseReader, blockNum *big.Int, map3Node common.Address) []microstaking.MicrostakingReturnRecord {
	data, _ := db.Get(terminateMap3ReturnRecordKey(blockNum, map3Node))
	if len(data) == 0 {
		return nil
	}
	var terminateMap3ReturnRecords []microstaking.MicrostakingReturnRecord
	if err := rlp.DecodeBytes(data, &terminateMap3ReturnRecords); err != nil {
		log.Error("Invalid terminateMap3ReturnRecords RLP")
		return nil
	}
	return terminateMap3ReturnRecords
}

func WriteTokenBurningRecord(writer DatabaseWriter, blockNum *big.Int, record burning.Record) {
	bytes, err := rlp.EncodeToBytes(record)
	if err != nil {
		log.Error("Fail to encode burningRecord")
	}
	if err := writer.Put(tokenBurningRecordKey(blockNum), bytes); err != nil {
		log.Error("Fail to store burningRecord")
	}
}

func ReadTokenBurningRecords(db DatabaseReader, blockNum *big.Int) *burning.Record {
	data, _ := db.Get(tokenBurningRecordKey(blockNum))
	if len(data) == 0 {
		return nil
	}
	var record burning.Record
	if err := rlp.DecodeBytes(data, &record); err != nil {
		log.Error("Invalid burningRecord RLP")
		return nil
	}
	return &record
}