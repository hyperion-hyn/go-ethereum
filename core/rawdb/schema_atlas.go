package rawdb

import "encoding/binary"

var (
	activeMap3KeyPrefix     = []byte("ActiveMap3Node")
	terminatedMap3KeyPrefix = []byte("TerminatedMap3Node")
)

func activeMap3Key(epoch uint64) []byte {
	epochKey := make([]byte, 8)
	binary.BigEndian.PutUint64(epochKey, epoch)
	return append(activeMap3KeyPrefix, epochKey...)
}

func terminatedMap3Key(epoch uint64) []byte {
	epochKey := make([]byte, 8)
	binary.BigEndian.PutUint64(epochKey, epoch)
	return append(terminatedMap3KeyPrefix, epochKey...)
}
