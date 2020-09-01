package restaking

import (
	"github.com/ethereum/go-ethereum/common"
	common2 "github.com/ethereum/go-ethereum/staking/types/common"
	"github.com/pkg/errors"
	"math/big"
)

const (
	MaxPubKeyAllowed         = 1
	DoNotEnforceMaxBLS       = -1
	MaxNameLength            = 140
	MaxIdentityLength        = 140
	MaxWebsiteLength         = 140
	MaxSecurityContactLength = 140
	MaxDetailsLength         = 280
)

var (
	errCommissionRateTooLarge  = errors.New("commission rate and change rate can not be larger than max commission rate")
	errInvalidCommissionRate   = errors.New("commission rate, change rate and max rate should be a value ranging from 0.0 to 1.0")
	errNeedAtLeastOneSlotKey   = errors.New("need at least one slot key")
	ErrExcessiveBLSKeys        = errors.New("more slot keys provided than allowed")
	errDuplicateSlotKeys       = errors.New("slot keys can not have duplicates")
	ErrCommitteeNil            = errors.New("subcommittee is nil pointer")
	errNilMaxTotalDelegation   = errors.New("MaxTotalDelegation can not be nil")
	errAddressNotMatch         = errors.New("validator key not match")
	errSlotKeyToRemoveNotFound = errors.New("slot key to remove not found")
	errSlotKeyToAddExists      = errors.New("slot key to add already exists")
	errCannotChangeBannedTrait = errors.New("cannot change validator banned status")
)

func NewEmptyAddressSet() AddressSet_ {
	return AddressSet_{
		Keys: []*Address{},
		Set:  make(map[Address]*Bool),
	}
}

func NewAddressSetWithAddress(address common.Address) AddressSet_ {
	set := NewEmptyAddressSet()
	set.Put(address)
	return set
}

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


func (s *Storage_AddressSet_) LoadFully() (*AddressSet_, error) {
	s.Keys().load()
	length := s.Keys().Length()
	for i := 0; i < length; i++ {
		k := s.Keys().Get(i).Value()
		s.Set().Get(k).Value()
	}

	// copy
	src := s.obj
	des := AddressSet_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

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
func (s *Storage_Validator_) LoadFully() (*Validator_, error) {
	s.ValidatorAddress().Value()
	if _, err := s.OperatorAddresses().LoadFully(); err != nil {
		return nil, err
	}
	s.SlotPubKeys().load() // need check
	s.LastEpochInCommittee().Value()
	s.MaxTotalDelegation().Value()
	s.Status().Value()
	s.Commission().load()
	s.Description().load()
	s.CreationHeight().Value()

	// copy
	src := s.obj
	des := Validator_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
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

func (s *Storage_ValidatorWrapper_) LoadFully() (*ValidatorWrapper_, error) {
	if _, err := s.Validator().LoadFully(); err != nil {
		return nil, err
	}
	if _, err := s.Redelegations().LoadFully(); err != nil {
		return nil, err
	}
	s.Counters().load()
	s.BlockReward().Value()
	s.TotalDelegation().Value()
	s.TotalDelegationByOperator().Value()

	// copy
	src := s.obj
	des := ValidatorWrapper_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
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
	s.Committee().Clear()
	s.Committee().Save(committee)
}

// CreateValidatorFromNewMsg creates validator from NewValidator message
func CreateValidatorFromNewMsg(msg *CreateValidator, valAddr common.Address, blockNum *big.Int) (*Validator_, error) {
	if err := msg.Description.EnsureLength(); err != nil {
		return nil, err
	}
	// TODO(ATLAS): default max?
	commission := Commission_{msg.CommissionRates, blockNum}

	if err := common2.VerifyBLSKey(&msg.SlotPubKey, &msg.SlotKeySig); err != nil {
		return nil, err
	}

	v := Validator_{
		ValidatorAddress:     valAddr,
		OperatorAddresses:    NewAddressSetWithAddress(msg.OperatorAddress),
		SlotPubKeys:          NewBLSKeysWithBLSKey(msg.SlotPubKey),
		LastEpochInCommittee: new(big.Int),
		MaxTotalDelegation:   msg.MaxTotalDelegation, // TODO(ATLAS): default max?
		Status:               uint8(Active),
		Commission:           commission,
		Description:          msg.Description,
		CreationHeight:       blockNum,
	}
	return &v, nil
}

// UpdateValidatorFromEditMsg updates validator from EditValidator message
func UpdateValidatorFromEditMsg(validator *Validator_, edit *EditValidator) error {
	if validator.ValidatorAddress != edit.ValidatorAddress {
		return errAddressNotMatch
	}

	if err := validator.Description.IncrementalUpdateFrom(edit.Description); err != nil {
		return err
	}

	if edit.CommissionRate != nil {
		validator.Commission.CommissionRates.Rate = *edit.CommissionRate
	}

	if edit.MaxTotalDelegation != nil && edit.MaxTotalDelegation.Sign() != 0 {
		validator.MaxTotalDelegation = edit.MaxTotalDelegation
	}

	if edit.SlotKeyToRemove != nil {
		index := -1
		for i, key := range validator.SlotPubKeys.Keys {
			if *key == *edit.SlotKeyToRemove {
				index = i
				break
			}
		}
		// we found key to be removed
		if index >= 0 {
			validator.SlotPubKeys.Keys = append(
				validator.SlotPubKeys.Keys[:index], validator.SlotPubKeys.Keys[index+1:]...,
			)
		} else {
			return errSlotKeyToRemoveNotFound
		}
	}

	if edit.SlotKeyToAdd != nil {
		found := false
		for _, key := range validator.SlotPubKeys.Keys {
			if *key == *edit.SlotKeyToAdd {
				found = true
				break
			}
		}
		if !found {
			if err := common2.VerifyBLSKey(edit.SlotKeyToAdd, edit.SlotKeyToAddSig); err != nil {
				return err
			}
			validator.SlotPubKeys.Keys = append(validator.SlotPubKeys.Keys, edit.SlotKeyToAdd)
		} else {
			return errSlotKeyToAddExists
		}
	}

	switch validator.Status {
	case Uint8(Banned):
		return errCannotChangeBannedTrait
	default:
		switch edit.EPOSStatus {
		case Active, Inactive:
			validator.Status = Uint8(edit.EPOSStatus)
		default:
		}
	}
	return nil
}
