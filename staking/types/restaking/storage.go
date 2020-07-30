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
	return 0
}

func (s *Storage_BLSPublicKeys_) Get(index int) *BLSPublicKey_ {
	return nil
}

func (s *Storage_BLSPublicKeys_) Set(index int, key *BLSPublicKey_) {
}

func (s *Storage_BLSPublicKeys_) Remove(index int, keepOrder bool) {
}

func (s *Storage_BLSPublicKeys_) Push(key *BLSPublicKey_) {
}

func (s *Storage_BLSPublicKeys_) Pop() *BLSPublicKey_ {
	return nil
}

func (s *Storage_BLSPublicKeys_) Load() *BLSPublicKeys_ {
	return nil
}



// Storage_AddressSet_
func (s *Storage_AddressSet_) AllKeys() []common.Address {
	return nil
}




// Storage_Validator_
func (s *Storage_Validator_) Load() *Validator_ {
	return nil
}


// Storage_ValidatorWrapper_
func (s *Storage_ValidatorWrapper_) AddBlockReward(reward *big.Int) {
}

func (s *Storage_ValidatorWrapper_) AddTotalDelegation(amount *big.Int) {
}

func (s *Storage_ValidatorWrapper_) SubTotalDelegation(amount *big.Int) {
}

func (s *Storage_ValidatorWrapper_) AddTotalDelegationByOperator(amount *big.Int) {
}

func (s *Storage_ValidatorWrapper_) SubTotalDelegationByOperator(amount *big.Int) {
}


// Storage_ValidatorWrapperMap_
func (s *Storage_ValidatorWrapperMap_) AllKeys() []common.Address {
	return nil
}

func (s *Storage_ValidatorWrapperMap_) Put(key common.Address, validator *ValidatorWrapper_) {
}

func (s *Storage_ValidatorWrapperMap_) Contain(key common.Address) bool {
	return false
}

func (s *Storage_ValidatorWrapperMap_) Get(key common.Address) (*Storage_ValidatorWrapper_, bool) {
	return nil, false
}

func (s *Storage_ValidatorWrapperMap_) Remove(key common.Address) {
	panic("no implement")
}

// Storage_Redelegation_
func (s *Storage_Redelegation_) AddReward(reward *big.Int) {
}

func (s *Storage_Redelegation_) AddAmount(amount *big.Int) {
}

// Storage_RedelegationMap_
func (s *Storage_RedelegationMap_) AllKeys() []common.Address {
	return nil
}

func (s *Storage_RedelegationMap_) Put(key common.Address, redelegation *Redelegation_) {
}

func (s *Storage_RedelegationMap_) Contain(key common.Address) bool {
	return false
}

func (s *Storage_RedelegationMap_) Get(key common.Address) (*Storage_Redelegation_, bool) {
	return nil, false
}

func (s *Storage_RedelegationMap_) Remove(key common.Address) {
	panic("no implement")
}



// Storage_Slots_
func (s *Storage_Slots_) Length() int {
	return 0
}

func (s *Storage_Slots_) Get(index int) *Storage_Slot_ {
	return nil
}

func (s *Storage_Slots_) Set(index int, key *Slot_) {
}

func (s *Storage_Slots_) Remove(index int, keepOrder bool) {
}

func (s *Storage_Slots_) Push(slot *Slot_) {
}

func (s *Storage_Slots_) Pop() *Storage_Slot_ {
	return nil
}





// Storage_Committee_
func (s *Storage_Committee_) Load() *Committee_ {
	return nil
}



// Storage_ValidatorPool_
func (s *Storage_ValidatorPool_) UpdateCommittee(committee *Committee_) {
	panic("no implement")
}
