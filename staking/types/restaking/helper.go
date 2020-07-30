package restaking

import (
	"bytes"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/staking/effective"
	"github.com/harmony-one/bls/ffi/go/bls"
	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"
	"golang.org/x/sync/singleflight"
	"math/big"
	"time"
)

const (
	DoNotEnforceMaxBLS       = -1
	BLSSignatureSizeInBytes  = 96
	MaxNameLength            = 140
	MaxIdentityLength        = 140
	MaxWebsiteLength         = 140
	MaxSecurityContactLength = 140
	MaxDetailsLength         = 280
	BLSVerificationStr       = "hyperion-hyn"
	MaxPubKeyAllowed         = 1
)

var (
	emptyBLSPubKey = BLSPublicKey_{}

	errAddressNotMatch         = errors.New("validator key not match")
	errBLSKeysNotMatchSigs     = errors.New("bls keys and corresponding signatures could not be verified")
	errCommissionRateTooLarge  = errors.New("commission rate and change rate can not be larger than max commission rate")
	errInvalidCommissionRate   = errors.New("commission rate, change rate and max rate should be a value ranging from 0.0 to 1.0")
	errNeedAtLeastOneSlotKey   = errors.New("need at least one slot key")
	ErrExcessiveBLSKeys        = errors.New("more slot keys provided than allowed")
	errDuplicateSlotKeys       = errors.New("slot keys can not have duplicates")
	errSlotKeyToRemoveNotFound = errors.New("slot key to remove not found")
	errSlotKeyToAddExists      = errors.New("slot key to add already exists")
	errCannotChangeBannedTrait = errors.New("cannot change validator banned status")
	ErrCommitteeNil            = errors.New("subcommittee is nil pointer")
)

// BLSSignature defines the bls signature
type BLSSignature [BLSSignatureSizeInBytes]byte

// Big ..
func (pk BLSPublicKey_) Big() *big.Int {
	return new(big.Int).SetBytes(pk.Key[:])
}

// IsEmpty returns whether the bls public key is empty 0 bytes
func (pk BLSPublicKey_) IsEmpty() bool {
	return bytes.Equal(pk.Key[:], emptyBLSPubKey.Key[:])
}

// Hex returns the hex string of bls public key
func (pk BLSPublicKey_) Hex() string {
	return hex.EncodeToString(pk.Key[:])
}

// ToLibBLSPublicKey copies the key contents into the given key.
func (pk *BLSPublicKey_) ToLibBLSPublicKey(key *bls.PublicKey) error {
	return key.Deserialize(pk.Key[:])
}

// EnsureLength ensures the length of a validator's description.
func (d Description_) EnsureLength() (Description_, error) {
	if len(d.Name) > MaxNameLength {
		return d, errors.Errorf(
			"exceed maximum name length %d %d", len(d.Name), MaxNameLength,
		)
	}
	if len(d.Identity) > MaxIdentityLength {
		return d, errors.Errorf(
			"exceed Maximum Length identity %d %d", len(d.Identity), MaxIdentityLength,
		)
	}
	if len(d.Website) > MaxWebsiteLength {
		return d, errors.Errorf(
			"exceed Maximum Length website %d %d", len(d.Website), MaxWebsiteLength,
		)
	}
	if len(d.SecurityContact) > MaxSecurityContactLength {
		return d, errors.Errorf(
			"exceed Maximum Length %d %d", len(d.SecurityContact), MaxSecurityContactLength,
		)
	}
	if len(d.Details) > MaxDetailsLength {
		return d, errors.Errorf(
			"exceed Maximum Length for details %d %d", len(d.Details), MaxDetailsLength,
		)
	}

	return d, nil
}

// VerifyBLSKeys checks if the public BLS key at index i of pubKeys matches the
// BLS key signature at index i of pubKeysSigs.
func VerifyBLSKeys(pubKeys BLSPublicKeys_, pubKeySigs []BLSSignature) error {
	if len(pubKeys.Keys) != len(pubKeySigs) {
		return errBLSKeysNotMatchSigs
	}

	for i := 0; i < len(pubKeys.Keys); i++ {
		if err := VerifyBLSKey(pubKeys.Keys[i], &pubKeySigs[i]); err != nil {
			return err
		}
	}

	return nil
}

// VerifyBLSKey checks if the public BLS key matches the BLS signature
func VerifyBLSKey(pubKey *BLSPublicKey_, pubKeySig *BLSSignature) error {
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
func CreateValidatorFromNewMsg(msg *staking.CreateValidator, valAddr common.Address, blockNum *big.Int) (*Validator_, error) {
	desc, err := msg.Description.EnsureLength()
	if err != nil {
		return nil, err
	}
	commission := Commission_{msg.CommissionRates, blockNum}

	if err = VerifyBLSKeys(msg.SlotPubKeys, msg.SlotKeySigs); err != nil {
		return nil, err
	}

	exist := true
	v := Validator_{
		ValidatorAddress: valAddr,
		OperatorAddresses: AddressSet_{
			Keys: []*Address{&msg.OperatorAddress},
			Set:  map[Address]*Bool{msg.OperatorAddress: &exist},
		},
		SlotPubKeys:          msg.SlotPubKeys,
		LastEpochInCommittee: new(big.Int),
		Status:               big.NewInt(int64(effective.Active)),
		Commission:           commission,
		Description:          desc,
		CreationHeight:       blockNum,
	}
	return &v, nil
}

// UpdateValidatorFromEditMsg updates validator from EditValidator message
func UpdateValidatorFromEditMsg(validator *Validator_, edit *staking.EditValidator) error {
	if validator.ValidatorAddress != edit.ValidatorAddress {
		return errAddressNotMatch
	}
	desc, err := UpdateDescription(validator.Description, *edit.Description)
	if err != nil {
		return err
	}
	validator.Description = desc

	if edit.CommissionRate != nil {
		validator.Commission.CommissionRates.Rate = *edit.CommissionRate
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
			if err := VerifyBLSKey(edit.SlotKeyToAdd, edit.SlotKeyToAddSig); err != nil {
				return err
			}
			validator.SlotPubKeys.Keys = append(validator.SlotPubKeys.Keys, edit.SlotKeyToAdd)
		} else {
			return errSlotKeyToAddExists
		}
	}

	switch validator.Status.Uint64() {
	case uint64(effective.Banned):
		return errCannotChangeBannedTrait
	default:
		switch edit.EPOSStatus {
		case effective.Active, effective.Inactive:
			validator.Status = big.NewInt(int64(edit.EPOSStatus))
		default:
		}
	}
	return nil
}

// UpdateDescription returns a new Description object with d1 as the base and the fields that's not empty in d2 updated
// accordingly. An error is returned if the resulting description fields have invalid length.
func UpdateDescription(d1, d2 Description_) (Description_, error) {
	newDesc := d1
	if d2.Name != "" {
		newDesc.Name = d2.Name
	}
	if d2.Identity != "" {
		newDesc.Identity = d2.Identity
	}
	if d2.Website != "" {
		newDesc.Website = d2.Website
	}
	if d2.SecurityContact != "" {
		newDesc.SecurityContact = d2.SecurityContact
	}
	if d2.Details != "" {
		newDesc.Details = d2.Details
	}
	return newDesc.EnsureLength()
}

var (
	hundredPercent = common.OneDec()
	zeroPercent    = common.ZeroDec()
)

// SanityCheck checks basic requirements of a validator
func (v *Validator_) SanityCheck(maxSlotKeyAllowed int) error {
	if _, err := v.Description.EnsureLength(); err != nil {
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

// NewDelegation creates a new delegation object
func NewRedelegation(delegatorAddr common.Address, amount *big.Int) *Redelegation_ {
	return &Redelegation_{
		DelegatorAddress: delegatorAddr,
		Amount:           amount,
		Reward:           big.NewInt(0),
	}
}

// Computed represents current epoch
// availability measures, mostly for RPC
type Computed struct {
	Signed            *big.Int   `json:"current-epoch-signed"`
	ToSign            *big.Int   `json:"current-epoch-to-sign"`
	BlocksLeftInEpoch uint64     `json:"-"`
	Percentage        common.Dec `json:"current-epoch-signing-percentage"`
	IsBelowThreshold  bool       `json:"-"`
}

// NewComputed ..
func NewComputed(
	signed, toSign *big.Int,
	blocksLeft uint64,
	percent common.Dec,
	isBelowNow bool) *Computed {
	return &Computed{signed, toSign, blocksLeft, percent, isBelowNow}
}

// BLSPublicKeys ..
func (c *Committee_) BLSPublicKeys() ([]*bls.PublicKey, error) {
	if c == nil {
		return nil, ErrCommitteeNil
	}
	return lookupBLSPublicKeys(c)
}

func (c *Committee_) Hash() (h common.Hash) {
	hw := sha3.New256()
	rlp.Encode(hw, c)
	hw.Sum(h[:0])
	return h
}

// StakedSlots gives overview of members
// in a subcommittee (aka a shard)
type StakedSlots struct {
	CountStakedValidator int
	CountStakedBLSKey    int
	Addrs                []common.Address
	LookupSet            map[common.Address]struct{}
	TotalEffectiveStaked common.Dec
}

// StakedValidators ..
func (c *Committee_) StakedValidators() *StakedSlots {
	countStakedValidator, countStakedBLSKey := 0, 0
	networkWideSlice, networkWideSet :=
		[]common.Address{}, map[common.Address]struct{}{}
	totalEffectiveStake := common.ZeroDec()

	for _, slot := range c.Slots.Entrys {
		// an external validator,
		// non-nil EffectiveStake is how we known
		if addr := slot.EcdsaAddress; slot.EffectiveStake != nil {
			totalEffectiveStake = totalEffectiveStake.Add(slot.EffectiveStake)
			countStakedBLSKey++
			if _, seen := networkWideSet[addr]; !seen {
				countStakedValidator++
				networkWideSet[addr] = struct{}{}
				networkWideSlice = append(networkWideSlice, addr)
			}
		}
	}

	return &StakedSlots{
		CountStakedValidator: countStakedValidator,
		CountStakedBLSKey:    countStakedBLSKey,
		Addrs:                networkWideSlice,
		LookupSet:            networkWideSet,
		TotalEffectiveStaked: totalEffectiveStake,
	}
}

var (
	blsKeyCache singleflight.Group
)

func lookupBLSPublicKeys(
	c *Committee_,
) ([]*bls.PublicKey, error) {
	key := c.Hash().Hex()
	results, err, _ := blsKeyCache.Do(
		key, func() (interface{}, error) {
			slice := make([]*bls.PublicKey, len(c.Slots.Entrys))
			for j := range c.Slots.Entrys {
				committerKey := &bls.PublicKey{}
				if err := c.Slots.Entrys[j].BLSPublicKey.ToLibBLSPublicKey(
					committerKey,
				); err != nil {
					return nil, err
				}
				slice[j] = committerKey
			}
			// Only made once
			go func() {
				time.Sleep(25 * time.Minute)
				blsKeyCache.Forget(key)
			}()
			return slice, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return results.([]*bls.PublicKey), nil
}
