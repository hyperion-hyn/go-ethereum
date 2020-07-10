package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/numeric"
	"github.com/ethereum/go-ethereum/staking/effective"
	"math/big"
)

// StakeMsg defines the interface of Stake Message
type StakeMsg interface {
	Copy() StakeMsg
}

type CreateMap3Node struct {
	InitiatorAddress common.Address   `json:"initiator-address"`
	Description      Description     `json:"description"`
	CommissionRates  CommissionRates `json:"commission"`
	NodeKeys         Map3NodeKeys
	Amount           *big.Int `json:"amount"`
}

type EditMap3Node struct {
	Map3NodeAddress    common.Address `json:"validator-address"`
	Description        *Description   `json:"description"`
	CommissionRate     numeric.Dec    `json:"commission-rate" rlp:"nil"`
	MinSelfDelegation  *big.Int       `json:"min-self-delegation" rlp:"nil"`
	MaxTotalDelegation *big.Int       `json:"max-total-delegation" rlp:"nil"`
	NodeKeyToRemove    *Map3NodeKey
	NodeKeyToAdd       *Map3NodeKey
}

// Microdelegate - type for delegating to a map3 node
type Microdelegate struct {
	DelegatorAddress common.Address `json:"delegator_address"`
	Map3NodeAddress  common.Address
	Amount           *big.Int
}

// Unmicrodelegate - type for removing delegation responsibility
type Unmicrodelegate struct {
	DelegatorAddress common.Address
	Map3NodeAddress  common.Address
	Amount           *big.Int
}

// CollectMicrodelegationRewards - type for collecting token rewards
type CollectMicrodelegationRewards struct {
	DelegatorAddress common.Address `json:"delegator_address"`
}

type DivideMap3NodeStake struct {
	Map3NodeAddress common.Address
}

type RenewMap3NodeStake struct {
	Map3NodeAddress common.Address
	DelegatorAddress common.Address
	IsRenew         bool
}

// CreateValidator - type for creating a new validator
type CreateValidator struct {
	InitiatorAddress common.Address  `json:"initiator-address"`
	Description      Description     `json:"description"`
	CommissionRates  CommissionRates `json:"commission"`
	SlotPubKeys      BLSPublicKeys   `json:"slot-pub-keys"`
	SlotKeySigs      []BLSSignature  `json:"slot-key-sigs"`
}

// Copy returns a deep copy of the CreateValidator as a StakeMsg interface
func (v CreateValidator) Copy() StakeMsg {
	cp := CreateValidator{
		InitiatorAddress: v.InitiatorAddress,
		Description:      v.Description,
		CommissionRates:  v.CommissionRates.Copy(),
	}

	if v.SlotPubKeys != nil {
		cp.SlotPubKeys = make([]BLSPublicKey, len(v.SlotPubKeys))
		copy(cp.SlotPubKeys, v.SlotPubKeys)
	}
	if v.SlotKeySigs != nil {
		cp.SlotKeySigs = make([]BLSSignature, len(v.SlotKeySigs))
		copy(cp.SlotKeySigs, v.SlotKeySigs)
	}
	return cp
}

// EditValidator - type for edit existing validator
type EditValidator struct {
	ValidatorAddress   common.Address        `json:"validator-address"`
	Description        *Description          `json:"description"`
	CommissionRate     numeric.Dec           `json:"commission-rate" rlp:"nil"`
	MinSelfDelegation  *big.Int              `json:"min-self-delegation" rlp:"nil"`
	MaxTotalDelegation *big.Int              `json:"max-total-delegation" rlp:"nil"`
	SlotKeyToRemove    *BLSPublicKey         `json:"slot-key-to_remove" rlp:"nil"`
	SlotKeyToAdd       *BLSPublicKey         `json:"slot-key-to_add" rlp:"nil"`
	SlotKeyToAddSig    *BLSSignature         `json:"slot-key-to-add-sig" rlp:"nil"`
	EPOSStatus         effective.Eligibility `json:"epos-eligibility-status" rlp:"nil"`
}

// Copy returns a deep copy of the EditValidator as a StakeMsg interface
func (v EditValidator) Copy() StakeMsg {
	cp := EditValidator{
		ValidatorAddress: v.ValidatorAddress,
		Description:      v.Description,
		EPOSStatus:       v.EPOSStatus,
	}
	if !v.CommissionRate.IsNil() {
		cr := v.CommissionRate.Copy()
		cp.CommissionRate = cr
	}
	if v.MinSelfDelegation != nil {
		cp.MinSelfDelegation = new(big.Int).Set(v.MinSelfDelegation)
	}
	if v.MaxTotalDelegation != nil {
		cp.MaxTotalDelegation = new(big.Int).Set(v.MaxTotalDelegation)
	}
	if v.SlotKeyToRemove != nil {
		keyRem := *v.SlotKeyToRemove
		cp.SlotKeyToRemove = &keyRem
	}
	if v.SlotKeyToAdd != nil {
		keyAdd := *v.SlotKeyToAdd
		cp.SlotKeyToAdd = &keyAdd
	}
	if v.SlotKeyToAddSig != nil {
		sigAdd := *v.SlotKeyToAddSig
		cp.SlotKeyToAddSig = &sigAdd
	}
	return cp
}

// Redelegate - type for delegating to a validator
type Redelegate struct {
	DelegatorAddress common.Address `json:"delegator_address"`
	ValidatorAddress common.Address `json:"validator_address"`
	Amount           *big.Int       `json:"amount"`
}

// Copy returns a deep copy of the Redelegate as a StakeMsg interface
func (v Redelegate) Copy() StakeMsg {
	cp := Redelegate{
		DelegatorAddress: v.DelegatorAddress,
		ValidatorAddress: v.ValidatorAddress,
	}
	if v.Amount != nil {
		cp.Amount = new(big.Int).Set(v.Amount)
	}
	return cp
}

// Unredelegate - type for removing delegation responsibility
type Unredelegate struct {
	DelegatorAddress common.Address `json:"delegator_address"`
	ValidatorAddress common.Address `json:"validator_address"`
}

// Copy returns a deep copy of the Unredelegate as a StakeMsg interface
func (v Unredelegate) Copy() StakeMsg {
	cp := Unredelegate{
		DelegatorAddress: v.DelegatorAddress,
		ValidatorAddress: v.ValidatorAddress,
	}
	return cp
}

// CollectRedelegationRewards - type for collecting token rewards
type CollectRedelegationRewards struct {
	DelegatorAddress common.Address `json:"delegator_address"`
	ValidatorAddress common.Address `json:"validator_address"`
}

// Copy returns a deep copy of the CollectRedelegationRewards as a StakeMsg interface
func (v CollectRedelegationRewards) Copy() StakeMsg {
	return CollectRedelegationRewards{
		DelegatorAddress: v.DelegatorAddress,
		ValidatorAddress: v.ValidatorAddress,
	}
}
