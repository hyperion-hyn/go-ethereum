package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/staking/effective"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"math/big"
)

// CreateValidator - type for creating a new validator
type CreateValidator struct {
	OperatorAddress    common.Address             `json:"operator-address"`
	Description        restaking.Description_     `json:"description"`
	CommissionRates    restaking.CommissionRates_ `json:"commission"`
	MaxTotalDelegation *big.Int                   `json:"max-total-delegation"`
	SlotPubKeys        restaking.BLSPublicKeys_   `json:"slot-pub-keys"`
	SlotKeySigs        []BLSSignature             `json:"slot-key-sigs"`
}

// EditValidator - type for edit existing validator
type EditValidator struct {
	ValidatorAddress   common.Address           `json:"validator-address"`
	OperatorAddress    common.Address           `json:"operator-address"`
	Description        *restaking.Description_  `json:"description"`
	CommissionRate     *common.Dec              `json:"commission-rate" rlp:"nil"`
	MaxTotalDelegation *big.Int                 `json:"max-total-delegation"`
	SlotKeyToRemove    *restaking.BLSPublicKey_ `json:"slot-key-to_remove" rlp:"nil"`
	SlotKeyToAdd       *restaking.BLSPublicKey_ `json:"slot-key-to_add" rlp:"nil"`
	SlotKeyToAddSig    *BLSSignature            `json:"slot-key-to-add-sig" rlp:"nil"`
	EPOSStatus         effective.Eligibility    `json:"epos-eligibility-status" rlp:"nil"`
}

// Redelegate - type for delegating to a validator
type Redelegate struct {
	DelegatorAddress common.Address `json:"delegator_address"`
	ValidatorAddress common.Address `json:"validator_address"`
}

// Unredelegate - type for removing delegation responsibility
type Unredelegate struct {
	DelegatorAddress common.Address `json:"delegator_address"`
	ValidatorAddress common.Address `json:"validator_address"`
}

// CollectRedelegationRewards - type for collecting token rewards
type CollectRedelegationRewards struct {
	DelegatorAddress common.Address `json:"delegator_address"`
	ValidatorAddress common.Address `json:"validator_address"`
}
