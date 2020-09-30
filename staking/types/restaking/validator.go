package restaking

import (
	"github.com/ethereum/go-ethereum/common"
	common2 "github.com/ethereum/go-ethereum/staking/types/common"
	"github.com/pkg/errors"
	"math/big"
)

const (
	MaxPubKeyAllowed   = 1
	DoNotEnforceMaxBLS = -1
)

var (
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
	result := make([]common.Address, 0)
	length := s.Keys().Length()
	for i := 0; i < length; i++ {
		result = append(result, s.Keys().Get(i).Value())
	}
	return result
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

func (v *Validator_) ToPlainValidator() *PlainValidator {
	return &PlainValidator{
		ValidatorAddress: v.ValidatorAddress,
		OperatorAddresses: func() []common.Address {
			var operators []common.Address
			for _, key := range v.OperatorAddresses.Keys {
				operators = append(operators, *key)
			}
			return operators
		}(),
		SlotPubKeys: func() []BLSPublicKey_ {
			var pubKeys []BLSPublicKey_
			for _, pubKey := range v.SlotPubKeys.Keys {
				pubKeys = append(pubKeys, *pubKey)
			}
			return pubKeys
		}(),
		LastEpochInCommittee: v.LastEpochInCommittee,
		MaxTotalDelegation:   v.MaxTotalDelegation,
		Status:               v.Status,
		Commission:           v.Commission,
		Description:          v.Description,
		CreationHeight:       v.CreationHeight,
	}
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

func (v *ValidatorWrapper_) ToPlainValidatorWrapper() *PlainValidatorWrapper {
	return &PlainValidatorWrapper{
		Validator: *v.Validator.ToPlainValidator(),
		Redelegations: func() []Redelegation_ {
			var redelegations []Redelegation_
			for _, key := range v.Redelegations.Keys {
				redegation, ok := v.Redelegations.Get(*key)
				if ok {
					redelegations = append(redelegations, redegation)
				}
			}
			return redelegations
		}(),
		Counters:                  v.Counters,
		BlockReward:               v.BlockReward,
		TotalDelegation:           v.TotalDelegation,
		TotalDelegationByOperator: v.TotalDelegationByOperator,
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
	return s.Validator().OperatorAddresses().Set().Get(delegator).Value()
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

func (s *Storage_ValidatorWrapper_) Undelegate(delegator common.Address, epoch, amount *big.Int) {
	if redelegation, ok := s.Redelegations().Get(delegator); ok {
		amt := redelegation.Amount().Value()
		if amount != nil && amount.Cmp(amt) < 0 {
			amt = amount
		}
		redelegation.Undelegation().Amount().SetValue(amt)
		redelegation.Undelegation().Epoch().SetValue(epoch)
		redelegation.Amount().Clear()
		s.SubTotalDelegation(amt)
		if s.IsOperator(delegator) {
			s.SubTotalDelegationByOperator(amt)
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
		s.Map().Get(key).Entry().Clear() // TODO(ATLAS): not supported
		s.Map().Get(key).Entry().Save(validator)
	} else {
		length := s.Keys().Length()
		//set keys
		s.Keys().Get(length).SetValue(key)
		//set map
		entry := s.Map().Get(key)
		entry.Index().SetValue(big.NewInt(int64(length + 1))) // because index start with 1
		entry.Entry().Save(validator)
	}
}

func (s *Storage_ValidatorWrapperMap_) Contain(key common.Address) bool {
	return s.Map().Get(key).Index().Value().Sign() > 0
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
func CreateValidatorFromNewMsg(msg *CreateValidator, valAddr common.Address, amount, blockNum *big.Int) (*ValidatorWrapper_, error) {
	if err := common2.VerifyBLSKey(&msg.SlotPubKey, &msg.SlotKeySig); err != nil {
		return nil, err
	}

	builder := NewValidatorWrapperBuilder()
	return builder.SetValidatorAddress(valAddr).
		AddOperatorAddress(msg.OperatorAddress).
		AddSlotPubKey(msg.SlotPubKey).
		SetMaxTotalDelegation(msg.MaxTotalDelegation).
		SetStatus(Active).
		SetCommission(Commission_{CommissionRates: msg.CommissionRates, UpdateHeight: blockNum}).
		SetDescription(msg.Description).
		SetCreationHeight(blockNum).
		AddRedelegation(NewRedelegation(msg.OperatorAddress, amount)).
		Build(), nil
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

type PlainValidator struct {
	ValidatorAddress     common.Address   `json:"ValidatorAddress"`
	OperatorAddresses    []common.Address `json:"OperatorAddresses"`
	SlotPubKeys          []BLSPublicKey_  `json:"SlotPubKeys"`
	LastEpochInCommittee *big.Int         `json:"LastEpochInCommittee"`
	MaxTotalDelegation   *big.Int         `json:"MaxTotalDelegation"`
	Status               uint8            `json:"Status"`
	Commission           Commission_      `json:"Commission"`
	Description          Description_     `json:"Description"`
	CreationHeight       *big.Int         `json:"CreationHeight"`
}

func (v *PlainValidator) ToValidator() *Validator_ {
	return &Validator_{
		ValidatorAddress: v.ValidatorAddress,
		OperatorAddresses: func() AddressSet_ {
			set := NewEmptyAddressSet()
			for _, key := range v.OperatorAddresses {
				set.Put(key)
			}
			return set
		}(),
		SlotPubKeys: func() BLSPublicKeys_ {
			pubKeys := NewEmptyBLSKeys()
			for _, pubKey := range v.SlotPubKeys {
				pubKeys.Keys = append(pubKeys.Keys, &pubKey)
			}
			return pubKeys
		}(),
		LastEpochInCommittee: v.LastEpochInCommittee,
		MaxTotalDelegation:   v.MaxTotalDelegation,
		Status:               v.Status,
		Commission:           v.Commission,
		Description:          v.Description,
		CreationHeight:       v.CreationHeight,
	}
}

type PlainValidatorWrapper struct {
	Validator                 PlainValidator  `json:"Validator"`
	Redelegations             []Redelegation_ `json:"Redelegations"`
	Counters                  Counters_       `json:"Counters"`
	BlockReward               BigInt          `json:"BlockReward"`
	TotalDelegation           BigInt          `json:"TotalDelegation"`
	TotalDelegationByOperator BigInt          `json:"TotalDelegationByOperator"`
}

func (v *PlainValidatorWrapper) ToValidatorWrapper() *ValidatorWrapper_ {
	return &ValidatorWrapper_{
		Validator: *v.Validator.ToValidator(),
		Redelegations: func() RedelegationMap_ {
			delegations := NewRedelegationMap()
			for _, delegation := range v.Redelegations {
				delegations.Put(delegation.DelegatorAddress, delegation)
			}
			return delegations
		}(),
		Counters:                  v.Counters,
		BlockReward:               v.BlockReward,
		TotalDelegation:           v.TotalDelegation,
		TotalDelegationByOperator: v.TotalDelegationByOperator,
	}
}
