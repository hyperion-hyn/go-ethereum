package state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/staking/types/microstaking"
	"github.com/pkg/errors"
)

var (
	errMap3NodeNotExist    = errors.New("map3 node does not exist")

	map3StorageAddress = common.BigToAddress(common.Big2) // TODO(ATLAS): what address?
)

func (s *StateDB) Map3NodePool() *microstaking.Storage_Map3NodePool_ {
	// singleton
	s.map3Once.Do(func() {
		var g microstaking.Global_t
		globalSt := microstaking.New(&g, s, map3StorageAddress, common.Big0)
		s.map3NodePool = globalSt.Map3NodePool()
	})
	return s.map3NodePool
}

func (s *StateDB) Map3NodeByAddress(map3Address common.Address) (*microstaking.Storage_Map3NodeWrapper_, error) {
	node, ok := s.Map3NodePool().Nodes().Get(map3Address)
	if !ok {
		return nil, errMap3NodeNotExist
	}
	return node, nil
}