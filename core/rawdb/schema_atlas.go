package rawdb

import "encoding/binary"

var (
	mutateMap3KeyPrefix  = []byte("MutateMap3Node")
	renewedMap3KeyPrefix = []byte("RenewedMap3Node")
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
