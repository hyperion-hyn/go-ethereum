package restaking

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

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
			lastDelegator := r.Keys[len(r.Keys) - 1]
			r.Keys[index - 1] = lastDelegator
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

func (s *Storage_Redelegation_) Save(redelegation Redelegation_) {
	s.DelegatorAddress().SetValue(redelegation.DelegatorAddress)
	if redelegation.Amount != nil {
		s.Amount().SetValue(redelegation.Amount)
	}
	if redelegation.Reward != nil {
		s.Reward().SetValue(redelegation.Reward)
	}
	if redelegation.Undelegation.Amount != nil {
		s.Undelegation().Amount().SetValue(redelegation.Undelegation.Amount)
	}
	if redelegation.Undelegation.Epoch != nil {
		s.Undelegation().Epoch().SetValue(redelegation.Undelegation.Epoch)
	}
}
func (s *Storage_Redelegation_) SetNil() {
	s.DelegatorAddress().SetValue(common.BigToAddress(common.Big0))
	s.Amount().SetValue(common.Big0)
	s.Reward().SetValue(common.Big0)
	s.Undelegation().Amount().SetValue(common.Big0)
	s.Undelegation().Epoch().SetValue(common.Big0)
}

func (s *Storage_Redelegation_) CanReleaseAt(epoch *big.Int) bool {
	return s.Undelegation().Amount().Value().Cmp(common.Big0) > 0 && s.Undelegation().Epoch().Value().Cmp(epoch) >= 0
}


// Storage_RedelegationMap_
func (s *Storage_RedelegationMap_) AllKeys() []common.Address {
	addressSlice := make([]common.Address, 0)
	addressLength := s.Keys().Length()
	for i := 0; i < addressLength; i++ {
		addressSlice = append(addressSlice, s.Keys().Get(i).Value())
	}
	return addressSlice
}

func (s *Storage_RedelegationMap_) Put(key common.Address, redelegation *Redelegation_) {
	if s.Contain(key) {
		s.Map().Get(key).Entry().Save(*redelegation)
	} else {
		keysLength := s.Keys().Length()
		//set keys
		s.Keys().Get(keysLength).SetValue(key)

		s.Get(key)
		//set map
		sRedelegation := s.Map().Get(key)
		//set map entity
		sRedelegationEntity := sRedelegation.Entry()
		sRedelegationEntity.Save(*redelegation)
		//set map index
		sRedelegation.Index().SetValue(big.NewInt(0).Add(big.NewInt(int64(keysLength)), common.Big1)) //because index start with 1
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

func (s *Storage_RedelegationMap_) Save(relegationMap RedelegationMap_) {
	relegationKeys := relegationMap.Keys
	s.Keys().Resize(len(relegationKeys))
	for i := 0; i < len(relegationKeys); i++ {
		addressTemp := relegationKeys[i]
		s.Keys().Get(i).SetValue(*addressTemp)
		s.Map().Get(*addressTemp).Entry().Save(relegationMap.Map[*addressTemp].Entry)
		s.Map().Get(*addressTemp).Index().SetValue(relegationMap.Map[*addressTemp].Index)
	}
}
