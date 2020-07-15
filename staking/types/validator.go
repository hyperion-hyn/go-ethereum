package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/clique/votepower"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/numeric"
	"github.com/ethereum/go-ethereum/staking/committee"
	"github.com/ethereum/go-ethereum/staking/effective"
	"github.com/harmony-one/bls/ffi/go/bls"
	"github.com/pkg/errors"
	"math/big"
)

// Define validator staking related const
const (
	MaxNameLength            = 140
	MaxIdentityLength        = 140
	MaxWebsiteLength         = 140
	MaxSecurityContactLength = 140
	MaxDetailsLength         = 280
	BLSVerificationStr       = "hyperion-hyn"
	TenThousand              = 10000
	APRHistoryLength         = 30
)

var (
	errAddressNotMatch = errors.New("validator key not match")
	// ErrInvalidSelfDelegation ..
	ErrInvalidSelfDelegation = errors.New(
		"self delegation can not be less than min_self_delegation",
	)
	errInvalidTotalDelegation = errors.New(
		"total delegation can not be bigger than max_total_delegation",
	)
	errMinSelfDelegationTooSmall = errors.New(
		"min_self_delegation must be greater than or equal to 10,000 ONE",
	)
	errCommissionRateTooLarge = errors.New(
		"commission rate and change rate can not be larger than max commission rate",
	)
	errInvalidCommissionRate = errors.New(
		"commission rate, change rate and max rate should be a value ranging from 0.0 to 1.0",
	)
	errNeedAtLeastOneSlotKey = errors.New("need at least one slot key")
	errBLSKeysNotMatchSigs   = errors.New(
		"bls keys and corresponding signatures could not be verified",
	)
	errNilMinSelfDelegation    = errors.New("MinSelfDelegation can not be nil")
	errNilMaxTotalDelegation   = errors.New("MaxTotalDelegation can not be nil")
	errSlotKeyToRemoveNotFound = errors.New("slot key to remove not found")
	errSlotKeyToAddExists      = errors.New("slot key to add already exists")
	errDuplicateSlotKeys       = errors.New("slot keys can not have duplicates")
	// ErrExcessiveBLSKeys ..
	ErrExcessiveBLSKeys        = errors.New("more slot keys provided than allowed")
	errCannotChangeBannedTrait = errors.New("cannot change validator banned status")
)

const (
	DoNotEnforceMaxBLS      = -1
	PublicKeySizeInBytes    = 48
	BLSSignatureSizeInBytes = 96
)

var (
	emptyBLSPubKey = BLSPublicKey{}
)

// BLSPublicKey defines the bls public key
type BLSPublicKey [PublicKeySizeInBytes]byte

type BLSPublicKeys []BLSPublicKey

// BLSSignature defines the bls signature
type BLSSignature [BLSSignatureSizeInBytes]byte

// Big ..
func (pk BLSPublicKey) Big() *big.Int {
	return new(big.Int).SetBytes(pk[:])
}

// IsEmpty returns whether the bls public key is empty 0 bytes
func (pk BLSPublicKey) IsEmpty() bool {
	return bytes.Equal(pk[:], emptyBLSPubKey[:])
}

// Hex returns the hex string of bls public key
func (pk BLSPublicKey) Hex() string {
	return hex.EncodeToString(pk[:])
}

// ToLibBLSPublicKey copies the key contents into the given key.
func (pk *BLSPublicKey) ToLibBLSPublicKey(key *bls.PublicKey) error {
	return key.Deserialize(pk[:])
}

type Counters struct {
	// The number of blocks the validator
	// should've signed when in active mode (selected in committee)
	NumBlocksToSign *big.Int `json:"to-sign",rlp:"nil"`
	// The number of blocks the validator actually signed
	NumBlocksSigned *big.Int `json:"signed",rlp:"nil"`
}

// ValidatorWrapper contains validator,
// its delegation information
type ValidatorWrapper struct {
	Validator       Validator
	Redelegations   Redelegations
	Counters        Counters
	BlockReward     *big.Int // All the rewarded accumulated so far
	TotalDelegation *big.Int
}

type ValidatorWrappers map[common.Address]ValidatorWrapper

// Computed represents current epoch
// availability measures, mostly for RPC
type Computed struct {
	Signed            *big.Int    `json:"current-epoch-signed"`
	ToSign            *big.Int    `json:"current-epoch-to-sign"`
	BlocksLeftInEpoch uint64      `json:"-"`
	Percentage        numeric.Dec `json:"current-epoch-signing-percentage"`
	IsBelowThreshold  bool        `json:"-"`
}

func (c Computed) String() string {
	s, _ := json.Marshal(c)
	return string(s)
}

// NewComputed ..
func NewComputed(
	signed, toSign *big.Int,
	blocksLeft uint64,
	percent numeric.Dec,
	isBelowNow bool) *Computed {
	return &Computed{signed, toSign, blocksLeft, percent, isBelowNow}
}

// NewEmptyStats ..
func NewEmptyStats() *ValidatorStats {
	return &ValidatorStats{
		[]APREntry{},
		numeric.ZeroDec(),
		[]VoteWithCurrentEpochEarning{},
		effective.Booted,
	}
}

// CurrentEpochPerformance represents validator performance in the context of
// whatever current epoch is
type CurrentEpochPerformance struct {
	CurrentSigningPercentage Computed `json:"current-epoch-signing-percent"`
}

// ValidatorRPCEnhanced contains extra information for RPC consumer
type ValidatorRPCEnhanced struct {
	Wrapper              ValidatorWrapper         `json:"validator"`
	Performance          *CurrentEpochPerformance `json:"current-epoch-performance"`
	ComputedMetrics      *ValidatorStats          `json:"metrics"`
	TotalDelegated       *big.Int                 `json:"total-delegation"`
	CurrentlyInCommittee bool                     `json:"currently-in-committee"`
	EPoSStatus           string                   `json:"epos-status"`
	EPoSWinningStake     *numeric.Dec             `json:"epos-winning-stake"`
	BootedStatus         *string                  `json:"booted-status"`
	ActiveStatus         string                   `json:"active-status"`
	Lifetime             *AccumulatedOverLifetime `json:"lifetime"`
}

// AccumulatedOverLifetime ..
type AccumulatedOverLifetime struct {
	BlockReward *big.Int    `json:"reward-accumulated"`
	Signing     Counters    `json:"blocks"`
	APR         numeric.Dec `json:"apr"`
	EpochAPRs   []APREntry  `json:"epoch-apr"`
}

func (w ValidatorWrapper) String() string {
	s, _ := json.Marshal(w)
	return string(s)
}

// VoteWithCurrentEpochEarning ..
type VoteWithCurrentEpochEarning struct {
	Vote   votepower.PureStakedVote `json:"key"`
	Earned *big.Int                 `json:"earned-reward"`
}

// APREntry ..
type APREntry struct {
	Epoch *big.Int
	Value numeric.Dec
}

// ValidatorStats to record validator's performance and history records
type ValidatorStats struct {
	// APRs is the APR history containing APR's of epochs
	APRs []APREntry `json:"-"`
	// TotalEffectiveStake is the total effective stake this validator has
	TotalEffectiveStake numeric.Dec `json:"-"`
	// MetricsPerShard ..
	MetricsPerShard []VoteWithCurrentEpochEarning `json:"by-bls-key"`
	// BootedStatus
	BootedStatus effective.BootedStatus `json:"-"`
}

func (s ValidatorStats) String() string {
	str, _ := json.Marshal(s)
	return string(str)
}

// Validator - data fields for a validator
type Validator struct {
	// ECDSA address of the validator
	ValidatorAddress common.Address `json:"validator-address"`
	// validator's initiators (node address)
	InitiatorAddresses AddressSet `json:"initiator-addresses"`
	// The BLS public key of the validator for consensus
	SlotPubKeys BLSPublicKeys `json:"bls-public-keys"`
	// The number of the last epoch this validator is
	// selected in committee (0 means never selected)
	LastEpochInCommittee *big.Int `json:"last-epoch-in-committee"`
	// Is the validator active in participating
	// committee selection process or not
	Status effective.Eligibility `json:"-"`
	// commission parameters
	Commission Commission
	// description for the validator
	Description Description
	// CreationHeight is the height of creation
	CreationHeight *big.Int `json:"creation-height"`
}

type ValidatorPool struct {
	Validators               ValidatorWrappers
	SlotKeySet               PubKeySet
	DescriptionIdentitySet   DescriptionIdentitySet
	CommitteeByEpoch         committee.CommitteeByEpoch
}

// SanityCheck checks basic requirements of a validator
func (v *Validator) SanityCheck(maxSlotKeyAllowed int) error {
	if _, err := v.Description.EnsureLength(); err != nil {
		return err
	}

	if len(v.SlotPubKeys) == 0 {
		return errNeedAtLeastOneSlotKey
	}

	if c := len(v.SlotPubKeys); maxSlotKeyAllowed != DoNotEnforceMaxBLS &&
		c > maxSlotKeyAllowed {
		return errors.Wrapf(
			ErrExcessiveBLSKeys, "have: %d allowed: %d",
			c, maxSlotKeyAllowed,
		)
	}

	if v.Commission.CommissionRates.Rate.LT(zeroPercent) || v.Commission.CommissionRates.Rate.GT(hundredPercent) {
		return errors.Wrapf(
			errInvalidCommissionRate, "rate:%s", v.Commission.CommissionRates.Rate.String(),
		)
	}

	if v.Commission.CommissionRates.MaxRate.LT(zeroPercent) || v.Commission.CommissionRates.MaxRate.GT(hundredPercent) {
		return errors.Wrapf(
			errInvalidCommissionRate, "max rate:%s", v.Commission.CommissionRates.MaxRate.String(),
		)
	}

	if v.Commission.CommissionRates.MaxChangeRate.LT(zeroPercent) ||
		v.Commission.CommissionRates.MaxChangeRate.GT(hundredPercent) {
		return errors.Wrapf(
			errInvalidCommissionRate, "max change rate:%s", v.Commission.CommissionRates.MaxChangeRate.String(),
		)
	}

	if v.Commission.CommissionRates.Rate.GT(v.Commission.CommissionRates.MaxRate) {
		return errors.Wrapf(
			errCommissionRateTooLarge,
			"rate:%s max rate:%s", v.Commission.CommissionRates.Rate.String(),
			v.Commission.CommissionRates.MaxRate.String(),
		)
	}

	if v.Commission.CommissionRates.MaxChangeRate.GT(v.Commission.CommissionRates.MaxRate) {
		return errors.Wrapf(
			errCommissionRateTooLarge,
			"rate:%s max change rate:%s", v.Commission.CommissionRates.Rate.String(),
			v.Commission.CommissionRates.MaxChangeRate.String(),
		)
	}

	allKeys := map[BLSPublicKey]struct{}{}
	for i := range v.SlotPubKeys {
		if _, ok := allKeys[v.SlotPubKeys[i]]; !ok {
			allKeys[v.SlotPubKeys[i]] = struct{}{}
		} else {
			return errDuplicateSlotKeys
		}
	}
	return nil
}

var (
	hundredPercent = numeric.NewDec(1)
	zeroPercent    = numeric.NewDec(0)
)

// VerifyBLSKeys checks if the public BLS key at index i of pubKeys matches the
// BLS key signature at index i of pubKeysSigs.
func VerifyBLSKeys(pubKeys []BLSPublicKey, pubKeySigs []BLSSignature) error {
	if len(pubKeys) != len(pubKeySigs) {
		return errBLSKeysNotMatchSigs
	}

	for i := 0; i < len(pubKeys); i++ {
		if err := VerifyBLSKey(&pubKeys[i], &pubKeySigs[i]); err != nil {
			return err
		}
	}

	return nil
}

// VerifyBLSKey checks if the public BLS key matches the BLS signature
func VerifyBLSKey(pubKey *BLSPublicKey, pubKeySig *BLSSignature) error {
	if len(pubKeySig) == 0 {
		return errBLSKeysNotMatchSigs
	}

	blsPubKey := new(bls.PublicKey)
	if err := pubKey.ToLibBLSPublicKey(blsPubKey); err != nil {
		return errBLSKeysNotMatchSigs
	}

	msgSig := bls.Sign{}
	if err := msgSig.Deserialize(pubKeySig[:]); err != nil {
		return err
	}

	messageBytes := []byte(BLSVerificationStr)
	msgHash := crypto.Keccak256(messageBytes)
	if !msgSig.VerifyHash(blsPubKey, msgHash[:]) {
		return errBLSKeysNotMatchSigs
	}

	return nil
}

// CreateValidatorFromNewMsg creates validator from NewValidator message
func CreateValidatorFromNewMsg(val *CreateValidator, valAddr common.Address, blockNum *big.Int) (*Validator, error) {
	desc, err := val.Description.EnsureLength()
	if err != nil {
		return nil, err
	}
	commission := Commission{val.CommissionRates, blockNum}
	pubKeys := append(val.SlotPubKeys[0:0], val.SlotPubKeys...)

	if err = VerifyBLSKeys(pubKeys, val.SlotKeySigs); err != nil {
		return nil, err
	}

	v := Validator{
		ValidatorAddress:     valAddr,
		InitiatorAddresses:   AddressSet{val.InitiatorAddress: struct{}{}},
		SlotPubKeys:          pubKeys,
		LastEpochInCommittee: new(big.Int),
		Status:               effective.Active,
		Commission:           commission,
		Description:          desc,
		CreationHeight:       blockNum,
	}
	return &v, nil
}

// UpdateValidatorFromEditMsg updates validator from EditValidator message
func UpdateValidatorFromEditMsg(validator *Validator, edit *EditValidator) error {
	if validator.ValidatorAddress != edit.ValidatorAddress {
		return errAddressNotMatch
	}
	desc, err := UpdateDescription(validator.Description, *edit.Description)
	if err != nil {
		return err
	}
	validator.Description = desc

	if !edit.CommissionRate.IsNil() {
		validator.Commission.CommissionRates.Rate = edit.CommissionRate
	}

	if edit.SlotKeyToRemove != nil {
		index := -1
		for i, key := range validator.SlotPubKeys {
			if key == *edit.SlotKeyToRemove {
				index = i
				break
			}
		}
		// we found key to be removed
		if index >= 0 {
			validator.SlotPubKeys = append(
				validator.SlotPubKeys[:index], validator.SlotPubKeys[index+1:]...,
			)
		} else {
			return errSlotKeyToRemoveNotFound
		}
	}

	if edit.SlotKeyToAdd != nil {
		found := false
		for _, key := range validator.SlotPubKeys {
			if key == *edit.SlotKeyToAdd {
				found = true
				break
			}
		}
		if !found {
			if err := VerifyBLSKey(edit.SlotKeyToAdd, edit.SlotKeyToAddSig); err != nil {
				return err
			}
			validator.SlotPubKeys = append(validator.SlotPubKeys, *edit.SlotKeyToAdd)
		} else {
			return errSlotKeyToAddExists
		}
	}

	switch validator.Status {
	case effective.Banned:
		return errCannotChangeBannedTrait
	default:
		switch edit.EPOSStatus {
		case effective.Active, effective.Inactive:
			validator.Status = edit.EPOSStatus
		default:
		}
	}
	return nil
}

// String returns a human readable string representation of a validator.
func (v Validator) String() string {
	s, _ := json.Marshal(v)
	return string(s)
}
