package restaking

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// NewDelegation creates a new delegation object
func NewRedelegation(delegatorAddr common.Address, amount *big.Int) Redelegation_ {
	return Redelegation_{
		DelegatorAddress: delegatorAddr,
		Amount:           amount,
		Reward:           big.NewInt(0),
	}
}

func NewRedelegationMap() RedelegationMap_ {
	return RedelegationMap_{
		Keys: []*Address{},
		Map:  map[Address]*RedelegationMapEntry_{},
	}
}

func (r *RedelegationMap_) Contain(delegator Address) bool {
	_, ok := r.Map[delegator]
	return ok
}

func (r *RedelegationMap_) Put(delegator Address, redelegation Redelegation_) {
	entry, ok := r.Map[delegator]
	if ok {
		entry.Entry = redelegation
	} else {
		r.Keys = append(r.Keys, &delegator)
		r.Map[delegator] = &RedelegationMapEntry_{
			Entry: redelegation,
			Index: big.NewInt(int64(len(r.Keys))),
		}
	}
}

func (r *RedelegationMap_) Remove(delegator Address) {
	if valueEntry, ok := r.Map[delegator]; ok {
		index := valueEntry.Index.Uint64()
		if int(index) != len(r.Keys) { // the last one
			lastDelegator := r.Keys[len(r.Keys)-1]
			r.Keys[index-1] = lastDelegator
			r.Map[*lastDelegator].Index = big.NewInt(int64(index))
		}
		r.Keys = r.Keys[:len(r.Keys)-1]
		delete(r.Map, delegator)
	}
}

func (r *RedelegationMap_) Get(delegator Address) (Redelegation_, bool) {
	if entry, ok := r.Map[delegator]; ok {
		return entry.Entry, true
	}
	return Redelegation_{}, false
}

// Storage_Redelegation_
func (s *Storage_Redelegation_) AddReward(reward *big.Int) {
	rewardTemp := s.Reward().Value()
	rewardTemp = rewardTemp.Add(rewardTemp, reward)
	s.Reward().SetValue(rewardTemp)
}

func (s *Storage_Redelegation_) AddAmount(amount *big.Int) {
	amountTemp := s.Amount().Value()
	amountTemp = amountTemp.Add(amountTemp, amount)
	s.Amount().SetValue(amountTemp)
}

func (s *Storage_Redelegation_) CanReleaseAt(epoch *big.Int) bool {
	return s.Undelegation().Amount().Value().Cmp(common.Big0) > 0 && s.Undelegation().Epoch().Value().Cmp(epoch) >= 0
}

// Storage_RedelegationMap_
func (s *Storage_RedelegationMap_) AllKeys() []common.Address {
	result := make([]common.Address, 0)
	length := s.Keys().Length()
	for i := 0; i < length; i++ {
		result = append(result, s.Keys().Get(i).Value())
	}
	return result
}

func (s *Storage_RedelegationMap_) Put(key common.Address, redelegation *Redelegation_) {
	if s.Contain(key) {
		s.Map().Get(key).Entry().Clear()
		s.Map().Get(key).Entry().Save(redelegation)
	} else {
		length := s.Keys().Length()
		//set keys
		s.Keys().Get(length).SetValue(key)

		//set map
		entry := s.Map().Get(key)
		entry.Index().SetValue(big.NewInt(int64(length + 1))) // because index start with 1
		entry.Entry().Save(redelegation)
	}
}

func (s *Storage_RedelegationMap_) Contain(key common.Address) bool {
	return s.Map().Get(key).Index().Value().Cmp(common.Big0) > 0
}

func (s *Storage_RedelegationMap_) Get(key common.Address) (*Storage_Redelegation_, bool) {
	if s.Contain(key) {
		return s.Map().Get(key).Entry(), true
	} else {
		return nil, false
	}
}

func (s *Storage_RedelegationMap_) Remove(key common.Address) {
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

func (s *Storage_RedelegationMap_) LoadFully() (*RedelegationMap_, error) {
	s.Keys().load()
	length := s.Keys().Length()
	for i := 0; i < length; i++ {
		k := s.Keys().Get(i).Value()
		s.Map().Get(k).load()
	}

	// copy
	src := s.obj
	des := RedelegationMap_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}