package restaking

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"math/big"
)

const (
	DoNotEnforceMaxBLS       = -1
	BLSSignatureSizeInBytes  = 96
	MaxNameLength            = 140
	MaxIdentityLength        = 140
	MaxWebsiteLength         = 140
	MaxSecurityContactLength = 140
	MaxDetailsLength         = 280
)

var (
	errCommissionRateTooLarge = errors.New("commission rate and change rate can not be larger than max commission rate")
	errInvalidCommissionRate  = errors.New("commission rate, change rate and max rate should be a value ranging from 0.0 to 1.0")
	errNeedAtLeastOneSlotKey  = errors.New("need at least one slot key")
	ErrExcessiveBLSKeys       = errors.New("more slot keys provided than allowed")
	errDuplicateSlotKeys      = errors.New("slot keys can not have duplicates")
	ErrCommitteeNil           = errors.New("subcommittee is nil pointer")
	errNilMaxTotalDelegation  = errors.New("MaxTotalDelegation can not be nil")
)

func (a *AddressSet_) Contain(address common.Address) bool {
	_, ok := a.Set[address]
	return ok
}

func (a *AddressSet_) Put(address common.Address) {
	if a.Contain(address) {
		return
	}
	a.Keys = append(a.Keys, &address)
	a.Set[address] = func() *bool { t := true; return &t }()
}


// Storage_AddressSet_
func (s *Storage_AddressSet_) AllKeys() []common.Address {
	addressSlice := make([]common.Address, 0)
	addressLength := s.Keys().Length()
	for i := 0; i < addressLength; i++ {
		addressSlice = append(addressSlice, s.Keys().Get(i).Value())
	}
	return addressSlice
}

func (s *Storage_AddressSet_) Save(addressSet AddressSet_) {
	keysLength := len(addressSet.Keys)
	if keysLength != 0 {
		s.Keys().Resize(keysLength)
		for i := 0; i < keysLength; i++ {
			keyTemp := addressSet.Keys[i]
			s.Keys().Get(i).SetValue(*keyTemp)
			s.Set().Get(*keyTemp).SetValue(*addressSet.Set[*keyTemp])
		}
	}
}

func (s *Storage_AddressSet_) Load() *AddressSet_ {
	length := s.Keys().Length()
	for i := 0; i < length; i++ {
		k := s.Keys().Get(i).Value()
		s.Set().Get(k).Value()
	}
	return s.obj
}


var (
	hundredPercent = common.OneDec()
	zeroPercent    = common.ZeroDec()
)

// ValidatorStatus represents ability to participate in EPoS auction
// that occurs just once an epoch
type ValidatorStatus byte

const (
	// Nil is a default state that represents a no-op
	Nil ValidatorStatus = iota
	// Active means allowed in epos auction
	Active
	// Inactive means validator did not sign enough over 66%
	// of the time in an epoch and so they are removed from
	// the possibility of being in the epos auction, which happens
	// only once an epoch and only
	// by beaconchain, aka shard.BeaconChainShardID
	Inactive
	// Banned records whether this validator is banned
	// from the network because they double-signed
	// it can never be undone
	Banned
)

func (e ValidatorStatus) String() string {
	switch e {
	case Active:
		return "active"
	case Inactive:
		return "inactive"
	case Banned:
		return "doubleSigningBanned"
	default:
		return "unknown"
	}
}

// SanityCheck checks basic requirements of a validator
func (v *Validator_) SanityCheck(maxSlotKeyAllowed int) error {
	if err := v.Description.EnsureLength(); err != nil {
		return err
	}

	if len(v.SlotPubKeys.Keys) == 0 {
		return errNeedAtLeastOneSlotKey
	}

	if c := len(v.SlotPubKeys.Keys); maxSlotKeyAllowed != DoNotEnforceMaxBLS &&
		c > maxSlotKeyAllowed {
		return errors.Wrapf(
			ErrExcessiveBLSKeys, "have: %d allowed: %d",
			c, maxSlotKeyAllowed,
		)
	}

	if v.MaxTotalDelegation == nil {
		return errNilMaxTotalDelegation
	}
	// TODO(ATLAS): minimal delegation?

	if err := v.Commission.SanityCheck(); err != nil {
		return err
	}

	allKeys := map[string]struct{}{}
	for i := range v.SlotPubKeys.Keys {
		key := v.SlotPubKeys.Keys[i].Hex()
		if _, ok := allKeys[key]; !ok {
			allKeys[key] = struct{}{}
		} else {
			return errDuplicateSlotKeys
		}
	}
	return nil
}


// Storage_Validator_
func (s *Storage_Validator_) Load() *Validator_ {
	s.ValidatorAddress().Value()
	s.OperatorAddresses().Load()
	s.SlotPubKeys().Load() // need check
	s.LastEpochInCommittee().Value()
	s.MaxTotalDelegation().Value()
	s.Status().Value()
	s.Commission().Load()
	s.Description().Load()
	s.CreationHeight().Value()
	return s.obj
}

func (s *Storage_Validator_) Save(validator *Validator_) {
	if validator.ValidatorAddress != common.BigToAddress(common.Big0) {
		s.ValidatorAddress().SetValue(validator.ValidatorAddress)
	}
	if len(validator.OperatorAddresses.Keys) != 0 {
		s.OperatorAddresses().Save(validator.OperatorAddresses)
	}
	if len(validator.SlotPubKeys.Keys) != 0 {
		s.SlotPubKeys().Save(&validator.SlotPubKeys)
	}

	if validator.LastEpochInCommittee != nil {
		s.LastEpochInCommittee().SetValue(validator.LastEpochInCommittee)
	}

	if validator.MaxTotalDelegation != nil {
		s.MaxTotalDelegation().SetValue(validator.MaxTotalDelegation)
	}

	if validator.Status != uint8(Nil) {
		s.Status().SetValue(validator.Status)
	}

	if !validator.Commission.CommissionRates.Rate.IsNil() {
		s.Commission().CommissionRates().Rate().SetValue(validator.Commission.CommissionRates.Rate)
	}
	if !validator.Commission.CommissionRates.MaxChangeRate.IsNil() {
		s.Commission().CommissionRates().MaxChangeRate().SetValue(validator.Commission.CommissionRates.MaxChangeRate)
	}
	if !validator.Commission.CommissionRates.MaxRate.IsNil() {
		s.Commission().CommissionRates().MaxRate().SetValue(validator.Commission.CommissionRates.MaxRate)
	}
	if validator.Commission.UpdateHeight != nil {
		s.Commission().UpdateHeight().SetValue(validator.Commission.UpdateHeight)
	}

	if validator.Description.Name != "" {
		s.Description().Name().SetValue(validator.Description.Name)
	}
	if validator.Description.Identity != "" {
		s.Description().Identity().SetValue(validator.Description.Identity)
	}
	if validator.Description.Website != "" {
		s.Description().Website().SetValue(validator.Description.Website)
	}
	if validator.Description.SecurityContact != "" {
		s.Description().SecurityContact().SetValue(validator.Description.SecurityContact)
	}
	if validator.Description.Details != "" {
		s.Description().Details().SetValue(validator.Description.Details)
	}

	if validator.CreationHeight != nil {
		s.CreationHeight().SetValue(validator.CreationHeight)
	}
}

// Storage_ValidatorWrapper_
func (s *Storage_ValidatorWrapper_) Save(validatorWrapper *ValidatorWrapper_) {
	s.Validator().Save(&validatorWrapper.Validator)
	s.Redelegations().Save(validatorWrapper.Redelegations)
	if validatorWrapper.Counters.NumBlocksSigned != nil {
		s.Counters().NumBlocksSigned().SetValue(validatorWrapper.Counters.NumBlocksSigned)
	}
	if validatorWrapper.Counters.NumBlocksToSign != nil {
		s.Counters().NumBlocksToSign().SetValue(validatorWrapper.Counters.NumBlocksToSign)
	}
	if validatorWrapper.BlockReward != nil {
		s.BlockReward().SetValue(validatorWrapper.BlockReward)
	}
	if validatorWrapper.TotalDelegation != nil {
		s.TotalDelegation().SetValue(validatorWrapper.TotalDelegation)
	}
	if validatorWrapper.TotalDelegationByOperator != nil {
		s.TotalDelegationByOperator().SetValue(validatorWrapper.TotalDelegationByOperator)
	}
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

func (s *Storage_ValidatorWrapper_) IsOperator(delegator common.Address) bool {
	if s.Validator().OperatorAddresses().Set().Get(delegator).Value() {
		return true
	}
	return false
}

func (s *Storage_ValidatorWrapper_) AddRedelegation(delegator common.Address, amount *big.Int) {
	if redelegation, ok := s.Redelegations().Get(delegator); ok {
		redelegation.AddAmount(amount)
	} else {
		m := NewRedelegation(delegator, amount)
		s.Redelegations().Put(delegator, &m)
	}
	s.AddTotalDelegation(amount)
	if s.IsOperator(delegator) {
		s.AddTotalDelegationByOperator(amount)
	}
}

func (s *Storage_ValidatorWrapper_) Undelegate(delegator common.Address, epoch *big.Int) {
	if redelegation, ok := s.Redelegations().Get(delegator); ok {
		amount := redelegation.Amount().Value()
		redelegation.Undelegation().Amount().SetValue(amount)
		redelegation.Undelegation().Epoch().SetValue(epoch)
		redelegation.Amount().SetValue(common.Big0)
		s.SubTotalDelegation(amount)
		if s.IsOperator(delegator) {
			s.SubTotalDelegationByOperator(amount)
		}
	}
}

func (s *Storage_ValidatorWrapper_) Load() *ValidatorWrapper_ {
	s.Validator().Load()
	s.Redelegations().Load()
	s.Counters().NumBlocksSigned().Value()
	s.Counters().NumBlocksToSign().Value()
	s.BlockReward().Value()
	s.TotalDelegation().Value()
	s.TotalDelegationByOperator().Value()
	return s.obj
}


// Storage_ValidatorWrapperMap_
func (s *Storage_ValidatorWrapperMap_) AllKeys() []common.Address {
	addressSlice := make([]common.Address, 0)
	addressLength := s.Keys().Length()
	for i := 0; i < addressLength; i++ {
		addressSlice = append(addressSlice, s.Keys().Get(i).Value())
	}
	return addressSlice
}

func (s *Storage_ValidatorWrapperMap_) Put(key common.Address, validator *ValidatorWrapper_) {
	if s.Contain(key) {
		s.Map().Get(key).Entry().Save(validator)
	} else {
		keysLength := s.Keys().Length()
		//set keys
		s.Keys().Get(keysLength).SetValue(key)
		//set map
		sValidatorWrapper := s.Map().Get(key)
		//set map entity
		sValidatorWrapperEntity := sValidatorWrapper.Entry()
		sValidatorWrapperEntity.Save(validator)
		//set map index
		sValidatorWrapper.Index().SetValue(big.NewInt(0).Add(big.NewInt(int64(keysLength)), common.Big1)) //because index start with 1
	}
}

func (s *Storage_ValidatorWrapperMap_) Contain(key common.Address) bool {
	return s.Map().Get(key).Index().Value().Cmp(common.Big0) > 0
}

func (s *Storage_ValidatorWrapperMap_) Get(key common.Address) (*Storage_ValidatorWrapper_, bool) {
	if s.Contain(key) {
		return s.Map().Get(key).Entry(), true
	}
	return nil, false
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
