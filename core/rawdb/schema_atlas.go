package rawdb

import "encoding/binary"

var (
	mutateMap3KeyPrefix      = []byte("MutateMap3Node")
	renewActiveMap3KeyPrefix = []byte("RenewActiveMap3Node")
)

func mutateMap3Key(epoch uint64) []byte {
	epochKey := make([]byte, 8)
	binary.BigEndian.PutUint64(epochKey, epoch)
	return append(mutateMap3KeyPrefix, epochKey...)
}

func renewActiveMap3Key(epoch uint64) []byte {
	epochKey := make([]byte, 8)
	binary.BigEndian.PutUint64(epochKey, epoch)
	return append(renewActiveMap3KeyPrefix, epochKey...)
}
