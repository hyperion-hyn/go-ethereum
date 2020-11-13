package rawdb

import (
	"encoding/binary"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

var (
	mutateMap3KeyPrefix                    = []byte("MutateMap3Node")
	renewedMap3KeyPrefix                   = []byte("RenewedMap3Node")
	unmicrodelegationReturnRecordKeyPrefix = []byte("UnmicrodelegationReturnRecord")
	terminateMap3ReturnRecordKeyPrefix     = []byte("TerminateMap3ReturnRecord")
)

func mutateMap3Key(epoch uint64) []byte {
	epochKey := make([]byte, 8)
	binary.BigEndian.PutUint64(epochKey, epoch)
	return append(mutateMap3KeyPrefix, epochKey...)
}

func renewedMap3Key(epoch uint64) []byte {
	epochKey := make([]byte, 8)
	binary.BigEndian.PutUint64(epochKey, epoch)
	return append(renewedMap3KeyPrefix, epochKey...)
}

func unmicrodelegationReturnRecordKey(block *big.Int) []byte {
	blockKey := block.Bytes()
	return append(unmicrodelegationReturnRecordKeyPrefix, blockKey...)
}

func terminateMap3ReturnRecordKey(block *big.Int, map3Address common.Address) []byte {
	blockKey := block.Bytes()
	key := append(terminateMap3ReturnRecordKeyPrefix, blockKey...)
	key = append(key, map3Address.Bytes()...)
	return key
}
