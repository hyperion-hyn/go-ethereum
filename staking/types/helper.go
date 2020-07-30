package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/staking/effective"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"github.com/harmony-one/bls/ffi/go/bls"
	"github.com/pkg/errors"
	"math/big"
)

const (
	BLSSignatureSizeInBytes  = 96
	BLSVerificationStr       = "hyperion-hyn"
	MaxPubKeyAllowed         = 1
)

var (
	errAddressNotMatch         = errors.New("validator key not match")
	errBLSKeysNotMatchSigs     = errors.New("bls keys and corresponding signatures could not be verified")
	errSlotKeyToRemoveNotFound = errors.New("slot key to remove not found")
	errSlotKeyToAddExists      = errors.New("slot key to add already exists")
	errCannotChangeBannedTrait = errors.New("cannot change validator banned status")
)

// BLSSignature defines the bls signature
type BLSSignature [BLSSignatureSizeInBytes]byte

// VerifyBLSKeys checks if the public BLS key at index i of pubKeys matches the
// BLS key signature at index i of pubKeysSigs.
func VerifyBLSKeys(pubKeys restaking.BLSPublicKeys_, pubKeySigs []BLSSignature) error {
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
func VerifyBLSKey(pubKey *restaking.BLSPublicKey_, pubKeySig *BLSSignature) error {
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
func CreateValidatorFromNewMsg(msg *CreateValidator, valAddr common.Address, blockNum *big.Int) (*restaking.Validator_, error) {
	desc, err := msg.Description.EnsureLength()
	if err != nil {
		return nil, err
	}
	commission := restaking.Commission_{msg.CommissionRates, blockNum}

	if err = VerifyBLSKeys(msg.SlotPubKeys, msg.SlotKeySigs); err != nil {
		return nil, err
	}

	exist := true
	v := restaking.Validator_{
		ValidatorAddress: valAddr,
		OperatorAddresses: restaking.AddressSet_{
			Keys: []*restaking.Address{&msg.OperatorAddress},
			Set:  map[restaking.Address]*restaking.Bool{msg.OperatorAddress: &exist},
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
func UpdateValidatorFromEditMsg(validator *restaking.Validator_, edit *EditValidator) error {
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
func UpdateDescription(d1, d2 restaking.Description_) (restaking.Description_, error) {
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

// NewDelegation creates a new delegation object
func NewRedelegation(delegatorAddr common.Address, amount *big.Int) *restaking.Redelegation_ {
	return &restaking.Redelegation_{
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
