package restaking

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// CreateValidator - type for creating a new validator
type CreateValidator struct {
	OperatorAddress    common.Address   `json:"operator-address"`
	Description        Description_     `json:"description"`
	CommissionRates    CommissionRates_ `json:"commission"`
	MaxTotalDelegation *big.Int         `json:"max-total-delegation"`
	SlotPubKey        BLSPublicKey_   `json:"slot-pub-key"`
	SlotKeySig        BLSSignature   `json:"slot-key-sig"`
}

// EditValidator - type for edit existing validator
type EditValidator struct {
	ValidatorAddress   common.Address  `json:"validator-address"`
	OperatorAddress    common.Address  `json:"operator-address"`
	Description        Description_   `json:"description"`
	CommissionRate     *common.Dec     `json:"commission-rate" rlp:"nil"`
	MaxTotalDelegation *big.Int        `json:"max-total-delegation"`
	SlotKeyToRemove    *BLSPublicKey_  `json:"slot-key-to_remove" rlp:"nil"`
	SlotKeyToAdd       *BLSPublicKey_  `json:"slot-key-to_add" rlp:"nil"`
	SlotKeyToAddSig    *BLSSignature   `json:"slot-key-to-add-sig" rlp:"nil"`
	EPOSStatus         ValidatorStatus `json:"epos-eligibility-status" rlp:"nil"`
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

// CollectReward - type for collecting token rewards
type CollectReward struct {
	DelegatorAddress common.Address `json:"delegator_address"`
	ValidatorAddress common.Address `json:"validator_address"`
}
