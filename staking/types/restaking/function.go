package restaking

import (
	"bytes"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
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
)

var (
	emptyBLSPubKey = BLSPublicKey_{}

	errCommissionRateTooLarge  = errors.New("commission rate and change rate can not be larger than max commission rate")
	errInvalidCommissionRate   = errors.New("commission rate, change rate and max rate should be a value ranging from 0.0 to 1.0")
	errNeedAtLeastOneSlotKey   = errors.New("need at least one slot key")
	ErrExcessiveBLSKeys        = errors.New("more slot keys provided than allowed")
	errDuplicateSlotKeys       = errors.New("slot keys can not have duplicates")
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

// FromLibBLSPublicKey replaces the key contents with the given key,
func (pk *BLSPublicKey_) FromLibBLSPublicKey(key *bls.PublicKey) error {
	bs := key.Serialize()
	if len(bs) != len(pk.Key) {
		return errors.Errorf(
			"key size (BLS) size mismatch, expected %d have %d", len(pk.Key), len(bs),
		)
	}
	copy(pk.Key[:], bs)
	return nil
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
		if addr := slot.EcdsaAddress; slot.EffectiveStake.IsPositive() {
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
