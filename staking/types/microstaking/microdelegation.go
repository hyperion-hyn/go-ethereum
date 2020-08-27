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

func (s *Storage_PendingDelegation_) SetNil() {
	s.Amount().SetValue(common.Big0)
	s.UnlockedEpoch().SetValue(common.ZeroDec())
}

// Storage_Microdelegation_
func (s *Storage_Microdelegation_) Save(microdelegation Microdelegation_) {
	s.DelegatorAddress().SetValue(microdelegation.DelegatorAddress)
	if microdelegation.Amount != nil {
		s.Amount().SetValue(microdelegation.Amount)
	}
	if microdelegation.Reward != nil {
		s.Reward().SetValue(microdelegation.Reward)
	}
	if microdelegation.Undelegation.Amount != nil {
		s.Undelegation().Amount().SetValue(microdelegation.Undelegation.Amount)
	}
	if microdelegation.Undelegation.Epoch != nil {
		s.Undelegation().Epoch().SetValue(microdelegation.Undelegation.Epoch)
	}
	if microdelegation.PendingDelegation.Amount != nil {
		s.Undelegation().Amount().SetValue(microdelegation.PendingDelegation.Amount)
	}
	if !microdelegation.PendingDelegation.UnlockedEpoch.IsNil() {
		s.PendingDelegation().UnlockedEpoch().SetValue(microdelegation.PendingDelegation.UnlockedEpoch)
	}
	panic("implement")
}
func (s *Storage_Microdelegation_) SetNil() {
	s.DelegatorAddress().SetValue(common.BigToAddress(common.Big0))
	s.Amount().SetValue(common.Big0)
	s.Reward().SetValue(common.Big0)
	s.PendingDelegation().SetNil()
	s.Undelegation().Amount().SetValue(common.Big0)
	s.Undelegation().Epoch().SetValue(common.Big0)
	panic("implement")
}

func (s *Storage_Microdelegation_) Load() *Microdelegation_ {
	s.DelegatorAddress().Value()
	s.Amount().Value()
	s.Reward().Value()
	s.Undelegation().Amount().Value()
	s.Undelegation().Epoch().Value()
	panic("implement")
	return s.obj
}

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
	addressSlice := make([]common.Address, 0)
	addressLength := s.Keys().Length()
	for i := 0; i < addressLength; i++ {
		addressSlice = append(addressSlice, s.Keys().Get(i).Value())
	}
	return addressSlice
}

func (s *Storage_MicrodelegationMap_) Put(key common.Address, microdelegation *Microdelegation_) {
	if s.Contain(key) {
		s.Map().Get(key).Entry().Save(*microdelegation)
	} else {
		keysLength := s.Keys().Length()
		//set keys
		s.Keys().Get(keysLength).SetValue(key)

		s.Get(key)
		//set map
		sMicrodelegation := s.Map().Get(key)
		//set map entity
		entry := sMicrodelegation.Entry()
		entry.Save(*microdelegation)
		//set map index
		sMicrodelegation.Index().SetValue(big.NewInt(0).Add(big.NewInt(int64(keysLength)), common.Big1)) //because index start with 1
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
	//remove keys
	keysStorage := s.Keys()
	keysLength := keysStorage.Length()
	lastKey := keysStorage.Get(keysLength - 1).Value()
	keyIndex := s.Map().Get(key).Index().Value()
	keysStorage.Get(int(keyIndex.Uint64() - 1)).SetValue(keysStorage.Get(keysLength - 1).Value())
	keysStorage.Get(keysLength - 1).SetValue(common.BigToAddress(common.Big0))
	s.Keys().Resize(keysLength - 1)

	//remove map entry
	maps := s.Map()
	delegationElem := maps.Get(key)
	lastDelegationElem := maps.Get(lastKey)
	lastDelegationElem.Index().SetValue(keyIndex)

	delegationElem.Entry().SetNil()
	delegationElem.Index().SetValue(common.Big0)
}

func (s *Storage_MicrodelegationMap_) Save(microdelegationMap MicrodelegationMap_) {
	delegationKeys := microdelegationMap.Keys
	s.Keys().Resize(len(delegationKeys))
	for i := 0; i < len(delegationKeys); i++ {
		addressTemp := delegationKeys[i]
		s.Keys().Get(i).SetValue(*addressTemp)
		s.Map().Get(*addressTemp).Entry().Save(microdelegationMap.Map[*addressTemp].Entry)
		s.Map().Get(*addressTemp).Index().SetValue(microdelegationMap.Map[*addressTemp].Index)
	}
}

func (s *Storage_MicrodelegationMap_) Load() *MicrodelegationMap_ {
	length := s.Keys().Length()
	for i := 0; i < length; i++ {
		k := s.Keys().Get(i).Value()
		s.Map().Get(k).Index().Value()
		s.Map().Get(k).Entry().Load()
	}
	return s.obj
}

func (s *Storage_DelegationIndex_) SetNil() {
	s.Map3Address().SetValue(common.BigToAddress(common.Big0))
	s.IsOperator().SetValue(false)
}

func (s *Storage_DelegationIndex_) Save(delegationIndex DelegationIndex_) {
	s.Map3Address().SetValue(delegationIndex.Map3Address)
	s.IsOperator().SetValue(delegationIndex.IsOperator)
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

func (s *Storage_DelegationIndexMap_) Put(key common.Address, delegationIndex DelegationIndex_) {
	if s.Contain(key) {
		s.Map().Get(key).Entry().Save(delegationIndex)
	} else {
		keysLength := s.Keys().Length()
		//set keys
		s.Keys().Get(keysLength).SetValue(key)

		s.Get(key)
		//set map
		sMicrodelegation := s.Map().Get(key)
		//set map entity
		entry := sMicrodelegation.Entry()
		entry.Save(delegationIndex)
		//set map index
		sMicrodelegation.Index().SetValue(big.NewInt(0).Add(big.NewInt(int64(keysLength)), common.Big1)) //because index start with 1
	}
}

func (s *Storage_DelegationIndexMap_) Remove(key common.Address) {
	//remove keys
	keysStorage := s.Keys()
	keysLength := keysStorage.Length()
	lastKey := keysStorage.Get(keysLength - 1).Value()
	keyIndex := s.Map().Get(key).Index().Value()
	keysStorage.Get(int(keyIndex.Uint64() - 1)).SetValue(keysStorage.Get(keysLength - 1).Value())
	keysStorage.Get(keysLength - 1).SetValue(common.BigToAddress(common.Big0))
	s.Keys().Resize(keysLength - 1)

	//remove map entry
	maps := s.Map()
	entry := maps.Get(key)
	lastEntry := maps.Get(lastKey)
	lastEntry.Index().SetValue(keyIndex)

	entry.Entry().SetNil()
	entry.Index().SetValue(common.Big0)
}
