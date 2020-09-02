package microstaking

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

const (
	PendingDelegationLockPeriodInEpoch = 7
)

// NewMicrodelegation creates a new microdelegation object
func NewMicrodelegation(delegator common.Address, amount *big.Int, unlockedEpoch common.Dec, pending bool) Microdelegation_ {
	d := Microdelegation_{
		DelegatorAddress: delegator,
		Amount:           big.NewInt(0),
		Reward:           big.NewInt(0),
	}
	if pending {
		d.PendingDelegation = PendingDelegation_{
			Amount:        amount,
			UnlockedEpoch: unlockedEpoch,
		}
	} else {
		d.Amount = amount
	}
	return d
}

func NewMicrodelegationMap() MicrodelegationMap_ {
	return MicrodelegationMap_{
		Keys: []*Address{},
		Map:  map[Address]*MicrodelegationMapEntry_{},
	}
}

func (s *Storage_PendingDelegation_) AddAmount(amount, epoch *big.Int) {
	panic("no implement")
}

func (s *Storage_PendingDelegation_) SubAmount(amount *big.Int) {
	amt := s.Amount().Value()
	amt = amt.Sub(amt, amount)
	s.Amount().SetValue(amt)
}

// Storage_Microdelegation_
func (s *Storage_Microdelegation_) AddReward(reward *big.Int) {
	rewardTemp := s.Reward().Value()
	rewardTemp = rewardTemp.Add(rewardTemp, reward)
	s.Reward().SetValue(rewardTemp)
}

func (s *Storage_Microdelegation_) AddAmount(amount *big.Int) {
	amountTemp := s.Amount().Value()
	amountTemp = amountTemp.Add(amountTemp, amount)
	s.Amount().SetValue(amountTemp)
}

func (m *MicrodelegationMap_) Put(delegator Address, microdelegation Microdelegation_) {
	entry, ok := m.Map[delegator]
	if ok {
		entry.Entry = microdelegation
	} else {
		m.Keys = append(m.Keys, &delegator)
		m.Map[delegator] = &MicrodelegationMapEntry_{
			Entry: microdelegation,
			Index: big.NewInt(int64(len(m.Keys))),
		}
	}
}

func (m *MicrodelegationMap_) Remove(delegator Address) {
	if valueEntry, ok := m.Map[delegator]; ok {
		index := valueEntry.Index.Uint64()
		if int(index) != len(m.Keys) { // the last one
			lastDelegator := m.Keys[len(m.Keys)-1]
			m.Keys[index-1] = lastDelegator
			m.Map[*lastDelegator].Index = big.NewInt(int64(index))
		}
		m.Keys = m.Keys[:len(m.Keys)-1]
		delete(m.Map, delegator)
	}
}

func (m *MicrodelegationMap_) Get(delegator Address) (Microdelegation_, bool) {
	if entry, ok := m.Map[delegator]; ok {
		return entry.Entry, true
	}
	return Microdelegation_{}, false
}

// Storage_MicrodelegationMap_
func (s *Storage_MicrodelegationMap_) AllKeys() []common.Address {
	result := make([]common.Address, 0)
	length := s.Keys().Length()
	for i := 0; i < length; i++ {
		result = append(result, s.Keys().Get(i).Value())
	}
	return result
}

func (s *Storage_MicrodelegationMap_) Put(key common.Address, microdelegation *Microdelegation_) {
	if s.Contain(key) {
		s.Map().Get(key).Entry().Clear()
		s.Map().Get(key).Entry().Save(microdelegation)
	} else {
		length := s.Keys().Length()
		//set keys
		s.Keys().Get(length).SetValue(key)

		//set map
		entry := s.Map().Get(key)
		entry.Index().SetValue(big.NewInt(int64(length + 1))) // because index start with 1
		entry.Entry().Save(microdelegation)
	}
}

func (s *Storage_MicrodelegationMap_) Contain(key common.Address) bool {
	return s.Map().Get(key).Index().Value().Cmp(common.Big0) > 0
}

func (s *Storage_MicrodelegationMap_) Get(key common.Address) (*Storage_Microdelegation_, bool) {
	if s.Contain(key) {
		return s.Map().Get(key).Entry(), true
	} else {
		return nil, false
	}
}

func (s *Storage_MicrodelegationMap_) Remove(key common.Address) {
	if !s.Contain(key) {
		return
	}

	entry := s.Map().Get(key)
	index := int(entry.Index().Value().Int64()) // index + 1

	// Move the last key to the index where the key to delete is
	length := s.Keys().Length()
	lastKey := s.Keys().Get(length - 1).Value()
	s.Keys().Get(index - 1).SetValue(lastKey)

	// Update the index for the moved key
	s.Map().Get(lastKey).Index().SetValue(big.NewInt(int64(index)))

	// Delete the slot where the moved key was stored
	s.Keys().Get(length - 1).Clear()
	s.Keys().Resize(length - 1)
	entry.Clear()
}

func (s *Storage_MicrodelegationMap_) LoadFully() (*MicrodelegationMap_, error) {
	s.Keys().load()
	length := s.Keys().Length()
	for i := 0; i < length; i++ {
		k := s.Keys().Get(i).Value()
		s.Map().Get(k).load()
	}

	// copy
	src := s.obj
	des := MicrodelegationMap_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_DelegationIndexMap_) Contain(key common.Address) bool {
	return s.Map().Get(key).Index().Value().Cmp(common.Big0) > 0
}

func (s *Storage_DelegationIndexMap_) Get(key common.Address) (*Storage_DelegationIndex_, bool) {
	if s.Contain(key) {
		return s.Map().Get(key).Entry(), true
	} else {
		return nil, false
	}
}

func (s *Storage_DelegationIndexMap_) Put(key common.Address, delegationIndex *DelegationIndex_) {
	if s.Contain(key) {
		s.Map().Get(key).Entry().Clear()
		s.Map().Get(key).Entry().Save(delegationIndex)
	} else {
		length := s.Keys().Length()
		//set keys
		s.Keys().Get(length).SetValue(key)

		//set map
		entry := s.Map().Get(key)
		entry.Index().SetValue(big.NewInt(int64(length + 1))) // because index start with 1
		entry.Entry().Save(delegationIndex)
	}
}

func (s *Storage_DelegationIndexMap_) Remove(key common.Address) {
	if !s.Contain(key) {
		return
	}

	entry := s.Map().Get(key)
	index := int(entry.Index().Value().Int64()) // index + 1

	// Move the last key to the index where the key to delete is
	length := s.Keys().Length()
	lastKey := s.Keys().Get(length - 1).Value()
	s.Keys().Get(index - 1).SetValue(lastKey)

	// Update the index for the moved key
	s.Map().Get(lastKey).Index().SetValue(big.NewInt(int64(index)))

	// Delete the slot where the moved key was stored
	s.Keys().Get(length - 1).Clear()
	s.Keys().Resize(length - 1)
	entry.Clear()
}
