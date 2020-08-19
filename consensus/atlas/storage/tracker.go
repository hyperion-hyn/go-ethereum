package storage

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
)

type StateTracker struct {
	db      *state.StateDB
	dirties map[common.Hash]common.Hash
}

func NewStateTracker(db *state.StateDB) *StateTracker {
	return &StateTracker{
		db:      db,
		dirties: make(map[common.Hash]common.Hash),
	}
}
func (m *StateTracker) GetState(addr common.Address, hash common.Hash) common.Hash {
	return m.db.GetState(addr, hash)
}

func (m *StateTracker) SetState(addr common.Address, key, value common.Hash) {
	m.dirties[key] = value
	m.db.SetState(addr, key, value)
}

func (m *StateTracker) Dirties() map[common.Hash]common.Hash {
	return m.dirties
}
