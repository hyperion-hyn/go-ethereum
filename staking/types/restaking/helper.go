package restaking

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/harmony-one/bls/ffi/go/bls"
	"github.com/pkg/errors"
	"math/big"
)

const (
	BLSVerificationStr = "hyperion-hyn"
	MaxPubKeyAllowed   = 1
)

var (
	errAddressNotMatch         = errors.New("validator key not match")
	errBLSKeysNotMatchSigs     = errors.New("bls keys and corresponding signatures could not be verified")
	errSlotKeyToRemoveNotFound = errors.New("slot key to remove not found")
	errSlotKeyToAddExists      = errors.New("slot key to add already exists")
	errCannotChangeBannedTrait = errors.New("cannot change validator banned status")
)

func NewEmptyBLSKeys() BLSPublicKeys_ {
	return BLSPublicKeys_{Keys: make([]*BLSPublicKey_, 0)}
}

func NewBLSKeysWithBLSKey(key BLSPublicKey_) BLSPublicKeys_ {
	keys := NewEmptyBLSKeys()
	keys.Keys = append(keys.Keys, &key)
	return keys
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
func CreateValidatorFromNewMsg(msg *CreateValidator, valAddr common.Address, blockNum *big.Int) (*Validator_, error) {
	if err := msg.Description.EnsureLength(); err != nil {
		return nil, err
	}
	// TODO(ATLAS): default max?
	commission := Commission_{msg.CommissionRates, blockNum}

	if err := VerifyBLSKey(&msg.SlotPubKey, &msg.SlotKeySig); err != nil {
		return nil, err
	}

	v := Validator_{
		ValidatorAddress:     valAddr,
		OperatorAddresses:    NewAddressSetWithAddress(msg.OperatorAddress),
		SlotPubKeys:          NewBLSKeysWithBLSKey(msg.SlotPubKey),
		LastEpochInCommittee: new(big.Int),
		MaxTotalDelegation:   msg.MaxTotalDelegation,	// TODO(ATLAS): default max?
		Status:               uint8(Active),
		Commission:           commission,
		Description:          msg.Description,
		CreationHeight:       blockNum,
	}
	return &v, nil
}

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

// UpdateValidatorFromEditMsg updates validator from EditValidator message
func UpdateValidatorFromEditMsg(validator *Validator_, edit *EditValidator) error {
	if validator.ValidatorAddress != edit.ValidatorAddress {
		return errAddressNotMatch
	}

	if err := validator.Description.UpdateFrom(edit.Description); err != nil {
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
			if err := VerifyBLSKey(edit.SlotKeyToAdd, edit.SlotKeyToAddSig); err != nil {
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

// NewDelegation creates a new delegation object
func NewRedelegation(delegatorAddr common.Address, amount *big.Int) Redelegation_ {
	return Redelegation_{
		DelegatorAddress: delegatorAddr,
		Amount:           amount,
		Reward:           big.NewInt(0),
	}
}

func NewRedelegationMap() RedelegationMap_ {
	return RedelegationMap_{
		Keys: []*Address{},
		Map:  map[Address]*RedelegationMapEntry_{},
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
