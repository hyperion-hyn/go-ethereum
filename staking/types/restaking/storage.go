package restaking

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// Storage_Description_
func (s *Storage_Description_) UpdateDescription(newDesc *Description_) {
	if newDesc.Name != "" {
		s.Name().SetValue(newDesc.Name)
	}
	if newDesc.Identity != "" {
		s.Identity().SetValue(newDesc.Identity)
	}
	if newDesc.Website != "" {
		s.Website().SetValue(newDesc.Website)
	}
	if newDesc.SecurityContact != "" {
		s.SecurityContact().SetValue(newDesc.SecurityContact)
	}
	if newDesc.Details != "" {
		s.Details().SetValue(newDesc.Details)
	}
}

// Storage_BLSPublicKeys_
func (s *Storage_BLSPublicKeys_) Length() int {
	return int(s.Keys().Length().Int64())
}

func (s *Storage_BLSPublicKeys_) Save(keys *BLSPublicKeys_) {
	length := len(keys.Keys)
	s.Keys().Resize(uint64(length))
	for i := 0; i < length; i++ {
		s.Keys().Get(uint64(i)).Key().SetValue(keys.Keys[i].Key)
	}
}

func (s *Storage_BLSPublicKeys_) Get(index int) *BLSPublicKey_ {
	s.Keys().Get(uint64(index)).Key().Value()
	return s.Keys().Get(uint64(index)).obj
}

func (s *Storage_BLSPublicKeys_) Set(index int, key *BLSPublicKey_) {
	s.Keys().Get(uint64(index)).Key().SetValue(key.Key)
}

func (s *Storage_BLSPublicKeys_) Remove(index int, keepOrder bool) {

	//remove current
	length := s.Length()
	lastOneStorage := s.Keys().Get(uint64(length - 1))
	//remove lastOne
	s.Keys().Get(uint64(length - 1)).Key().SetValue([48]uint8{})
	//replace lastOne to index
	s.Keys().Get(uint64(index)).Key().SetValue(lastOneStorage.Key().Value())
	//resize length
	s.Keys().Resize(uint64(length - 1))
}

func (s *Storage_BLSPublicKeys_) Push(key *BLSPublicKey_) {
	length := s.Length()

	//over length will auto resize , not resize again
	s.Keys().Get(uint64(length)).Key().SetValue(key.Key)
}

func (s *Storage_BLSPublicKeys_) Pop() *BLSPublicKey_ {

	length := s.Length()

	blsPublicKeyTemp :=
		BLSPublicKey_{Key: s.Keys().Get(uint64(length - 1)).Key().Value()}

	s.Keys().Get(uint64(length - 1)).Key().SetValue([48]uint8{})
	s.Keys().Resize(uint64(length - 1))
	return &blsPublicKeyTemp
}

func (s *Storage_BLSPublicKeys_) Load() *BLSPublicKeys_ {
	length := s.Length()

	for i := 0; i < length; i++ {
		s.Keys().Get(uint64(i)).Key().Value()
	}
	return s.obj
}

// Storage_AddressSet_
func (s *Storage_AddressSet_) AllKeys() []common.Address {

	addressSlice := make([]common.Address, 0)
	addressLength := s.Keys().Length().Int64()
	for i := int64(0); i < addressLength; i++ {
		addressSlice = append(addressSlice, s.Keys().Get(uint64(i)).Value())
	}
	return addressSlice
}

// Storage_Validator_
func (s *Storage_Validator_) Load() *Validator_ {

	s.ValidatorAddress().Value()
	s.SlotPubKeys().Load() // need check
	s.LastEpochInCommittee().Value()
	s.Status().Value()
	s.Commission().CommissionRates().Rate().Value()
	s.Commission().CommissionRates().MaxChangeRate().Value()
	s.Commission().CommissionRates().MaxRate().Value()
	s.Commission().UpdateHeight().Value()
	s.Description().Name().Value()
	s.Description().Identity().Value()
	s.Description().Website().Value()
	s.Description().SecurityContact().Value()
	s.Description().Details().Value()
	s.CreationHeight().Value()

	return s.obj
}

func (s *Storage_Validator_) Save(validator *Validator_) {

	s.ValidatorAddress().SetValue(validator.ValidatorAddress)
	s.SlotPubKeys().Save(&validator.SlotPubKeys)
	s.LastEpochInCommittee().SetValue(validator.LastEpochInCommittee)
	s.Status().SetValue(validator.Status)
	s.Commission().CommissionRates().Rate().SetValue(validator.Commission.CommissionRates.Rate)
	s.Commission().CommissionRates().MaxChangeRate().SetValue(validator.Commission.CommissionRates.MaxChangeRate)
	s.Commission().CommissionRates().MaxRate().SetValue(validator.Commission.CommissionRates.MaxRate)
	s.Commission().UpdateHeight().SetValue(validator.Commission.UpdateHeight)
	s.Description().Name().SetValue(validator.Description.Name)
	s.Description().Identity().SetValue(validator.Description.Identity)
	s.Description().Website().SetValue(validator.Description.Website)
	s.Description().SecurityContact().SetValue(validator.Description.SecurityContact)
	s.Description().Details().SetValue(validator.Description.Details)
	s.CreationHeight().SetValue(validator.CreationHeight)
}

// Storage_ValidatorWrapper_

func (s *Storage_ValidatorWrapper_) Save(validatorWrapper *ValidatorWrapper_) {

	s.Validator().Save(&validatorWrapper.Validator)
	s.Redelegations().Save(validatorWrapper.Redelegations)
	s.Counters().NumBlocksSigned().SetValue(validatorWrapper.Counters.NumBlocksSigned)
	s.Counters().NumBlocksToSign().SetValue(validatorWrapper.Counters.NumBlocksToSign)
	s.BlockReward().SetValue(validatorWrapper.BlockReward)
	s.TotalDelegation().SetValue(validatorWrapper.TotalDelegation)
	s.TotalDelegationByOperator().SetValue(validatorWrapper.TotalDelegationByOperator)

}

func (s *Storage_ValidatorWrapper_) AddBlockReward(reward *big.Int) {
	blockReward := s.BlockReward().Value()
	blockReward = blockReward.Add(blockReward, reward)
	s.BlockReward().SetValue(blockReward)
}

func (s *Storage_ValidatorWrapper_) AddTotalDelegation(amount *big.Int) {
	totalDelegation := s.TotalDelegation().Value()
	totalDelegation = totalDelegation.Add(totalDelegation, amount)
	s.TotalDelegation().SetValue(totalDelegation)
}

func (s *Storage_ValidatorWrapper_) SubTotalDelegation(amount *big.Int) {
	totalDelegation := s.TotalDelegation().Value()
	totalDelegation = totalDelegation.Sub(totalDelegation, amount)
	s.TotalDelegation().SetValue(totalDelegation)
}

func (s *Storage_ValidatorWrapper_) AddTotalDelegationByOperator(amount *big.Int) {
	totalDelegationByOperator := s.TotalDelegationByOperator().Value()
	totalDelegationByOperator = totalDelegationByOperator.Add(totalDelegationByOperator, amount)
	s.TotalDelegationByOperator().SetValue(totalDelegationByOperator)
}

func (s *Storage_ValidatorWrapper_) SubTotalDelegationByOperator(amount *big.Int) {
	totalDelegationByOperator := s.TotalDelegationByOperator().Value()
	totalDelegationByOperator = totalDelegationByOperator.Sub(totalDelegationByOperator, amount)
	s.TotalDelegationByOperator().SetValue(totalDelegationByOperator)
}

// Storage_ValidatorWrapperMap_
func (s *Storage_ValidatorWrapperMap_) AllKeys() []common.Address {

	addressSlice := make([]common.Address, 0)
	addressLength := s.Keys().Length().Int64()
	for i := int64(0); i < addressLength; i++ {
		addressSlice = append(addressSlice, s.Keys().Get(uint64(i)).Value())
	}
	return addressSlice
}

func (s *Storage_ValidatorWrapperMap_) Put(key common.Address, validator *ValidatorWrapper_) {

	keysLength := s.Keys().Length()
	//set keys
	s.Keys().Get(keysLength.Uint64()).SetValue(key)
	s.Get(key)
	//set map
	sValidatorWrapper := s.Map().Get(key)
	//set map entity
	sValidatorWrapperEntity := sValidatorWrapper.Entry()
	sValidatorWrapperEntity.Save(validator)
	//set map index
	sValidatorWrapper.Index().SetValue(keysLength.Add(keysLength, big.NewInt(1))) //because index start with 1

}

func (s *Storage_ValidatorWrapperMap_) Contain(key common.Address) bool {
	return s.Map().Get(key).Index().Value().Cmp(common.Big0) > 0
}

func (s *Storage_ValidatorWrapperMap_) Get(key common.Address) (*Storage_ValidatorWrapper_, bool) {

	if s.Contain(key) {
		return s.Map().Get(Address{}).Entry(), true
	} else {
		return nil, false
	}

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
	s.Amount().SetValue(redelegation.Amount)
	s.Reward().SetValue(redelegation.Reward)
	s.Undelegation().Amount().SetValue(redelegation.Undelegation.Amount)
	s.Undelegation().Epoch().SetValue(redelegation.Undelegation.Epoch)
}
func (s *Storage_Redelegation_) SetNil() {

	s.DelegatorAddress().SetValue(common.BigToAddress(common.Big0))
	s.Amount().SetValue(common.Big0)
	s.Reward().SetValue(common.Big0)
	s.Undelegation().Amount().SetValue(common.Big0)
	s.Undelegation().Epoch().SetValue(common.Big0)
}

// Storage_RedelegationMap_
func (s *Storage_RedelegationMap_) AllKeys() []common.Address {

	addressSlice := make([]common.Address, 0)
	addressLength := s.Keys().Length().Int64()
	for i := int64(0); i < addressLength; i++ {
		addressSlice = append(addressSlice, s.Keys().Get(uint64(i)).Value())
	}
	return addressSlice
}

func (s *Storage_RedelegationMap_) Put(key common.Address, redelegation *Redelegation_) {

	keysLength := s.Keys().Length()
	//set keys
	s.Keys().Get(keysLength.Uint64()).SetValue(key)

	s.Get(key)
	//set map
	sRedelegation := s.Map().Get(key)
	//set map entity
	sRedelegationEntity := sRedelegation.Entry()
	sRedelegationEntity.Save(*redelegation)
	//set map index
	sRedelegation.Index().SetValue(keysLength.Add(keysLength, big.NewInt(1))) //because index start with 1

}

func (s *Storage_RedelegationMap_) Contain(key common.Address) bool {
	return s.Map().Get(key).Index().Value().Cmp(common.Big0) > 0
}

func (s *Storage_RedelegationMap_) Get(key common.Address) (*Storage_Redelegation_, bool) {
	if s.Contain(key) {
		return s.Map().Get(Address{}).Entry(), true
	} else {
		return nil, false
	}
}

func (s *Storage_RedelegationMap_) Remove(key common.Address) {

	//remove keys
	keysStorage := s.Keys()
	keysLength := keysStorage.Length()
	lastKey := keysStorage.Get(keysLength.Uint64() - 1).Value()
	keyIndex := s.Map().Get(key).Index().Value()
	keysStorage.Get(keyIndex.Uint64() - 1).SetValue(keysStorage.Get(keysLength.Uint64() - 1).Value())
	keysStorage.Get(keysLength.Uint64() - 1).SetValue(common.BigToAddress(common.Big0))
	s.Keys().Resize(keysLength.Uint64() - 1)

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
	s.Keys().Resize(uint64(len(relegationKeys)))
	for i := 0; i < len(relegationKeys); i++ {
		addressTemp := relegationKeys[i]
		s.Keys().Get(uint64(i)).SetValue(*addressTemp)
		s.Map().Get(*addressTemp).Entry().Save(relegationMap.Map[*addressTemp].Entry)
	}
}

// Storage_Slots_

func (s *Storage_Slots_) Length() int {
	return int(s.Entrys().Length().Int64())
}

func (s *Storage_Slots_) Load() []*Slot_ {

	slotsLength := s.Length()
	for i := 0; i < slotsLength; i++ {
		s.Get(i).Load()
	}

	return s.obj.Entrys
}

func (s *Storage_Slots_) Get(index int) *Storage_Slot_ {
	return s.Entrys().Get(uint64(index))
}

func (s *Storage_Slots_) Set(index int, key *Slot_) {
	s.Entrys().Get(uint64(index)).Save(key)
}

func (s *Storage_Slots_) Remove(index int, keepOrder bool) {

	// remove from index

	oldEntriesLength := s.Entrys().Length()

	//set lastEntity to index
	lastEntry := s.Entrys().Get(oldEntriesLength.Uint64() - 1)
	s.Entrys().Get(uint64(index)).Save(&Slot_{
		EcdsaAddress: lastEntry.EcdsaAddress().Value(),
		BLSPublicKey: BLSPublicKey_{
			Key: lastEntry.BLSPublicKey().Key().Value(),
		},
		EffectiveStake: lastEntry.EffectiveStake().Value(),
	})

	//set lastEntity to zero

	lastEntry.SetNil()

	//resize slice
	s.Entrys().Resize(oldEntriesLength.Uint64() - 1)

}

func (s *Storage_Slots_) Push(slot *Slot_) {
	entityLength := s.Entrys().Length()
	s.Entrys().Get(entityLength.Uint64()).Save(slot)
}

func (s *Storage_Slots_) Pop() *Storage_Slot_ {
	entityLength := s.Entrys().Length()
	storageSlot := s.Entrys().Get(entityLength.Uint64() - 1)
	s.Remove(int(entityLength.Int64()-1), false)
	return storageSlot
}

func (s *Storage_Slots_) UpdateSlots(slots Slots_) {
	// remove old
	length := s.Length()
	for i := 0; i < length; i++ {
		s.Get(i).SetNil()
	}
	//set new
	newSlotsLength := len(slots.Entrys)
	s.Entrys().Resize(uint64(newSlotsLength))
	for i := 0; i < newSlotsLength; i++ {
		s.Get(i).Save(slots.Entrys[i])
	}
}

func (s *Storage_Slot_) SetNil() {

	s.EffectiveStake().SetValue(common.NewDec(int64(0)))
	s.EcdsaAddress().SetValue(common.BigToAddress(big.NewInt(0)))
	s.BLSPublicKey().Key().SetValue([48]uint8{})
}

func (s *Storage_Slot_) Save(key *Slot_) {

	s.BLSPublicKey().Key().SetValue(key.BLSPublicKey.Key)
	s.EcdsaAddress().SetValue(key.EcdsaAddress)
	s.EffectiveStake().SetValue(key.EffectiveStake)
}

func (s *Storage_Slot_) Load() *Slot_ {

	s.BLSPublicKey().Key().Value()
	s.EcdsaAddress().Value()
	s.EffectiveStake().Value()
	return s.obj
}

// Storage_Committee_
func (s *Storage_Committee_) Load() *Committee_ {

	s.Epoch().Value()
	s.Slots().Load()
	return s.obj
}

// Storage_ValidatorPool_
func (s *Storage_ValidatorPool_) UpdateCommittee(committee *Committee_) {

	if committee.Epoch != nil {
		s.Committee().Epoch().SetValue(committee.Epoch)
	}
	if committee.Slots.Entrys != nil {
		s.Committee().Slots().UpdateSlots(committee.Slots)
	}
}
