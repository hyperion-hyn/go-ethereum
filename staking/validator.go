package staking

import (
	"github.com/ethereum/go-ethereum/common"
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
)

var (
	StakingAddress = common.StringToAddress("XXXXX")

	errAddressNotMatch           = errors.New("Validator key not match")
)

// ValidatorWrapper contains validator and its delegation information
type ValidatorWrapper struct {
	Validators  map[common.Address]*Validator   `json:"validators"`
	Delegations map[common.Address]*Delegations `json:"delegations"`
}

// Validator - data fields for a validator
type Validator struct {
	// ECDSA address of the validator
	Address common.Address
	// description for the validator
	Description
}

// Description - some possible IRL connections
type Description struct {
	Name            string `json:"name"`             // name
	Identity        string `json:"identity"`         // optional identity signature (ex. UPort or Keybase)
	Website         string `json:"website"`          // optional website link
	SecurityContact string `json:"security_contact"` // optional security contact info
	Details         string `json:"details"`          // optional details
}

// EnsureLength ensures the length of a validator's description.
func (d Description) EnsureLength() (Description, error) {
	if len(d.Name) > MaxNameLength {
		return d, errors.Errorf("[EnsureLength] Exceed Maximum Length, have: %s, maxNameLen: %s", len(d.Name), MaxNameLength)
	}
	if len(d.Identity) > MaxIdentityLength {
		return d, errors.Errorf("[EnsureLength] Exceed Maximum Length, have: %s, maxIdentityLen: %s", len(d.Identity), MaxIdentityLength)
	}
	if len(d.Website) > MaxWebsiteLength {
		return d, errors.Errorf("[EnsureLength] Exceed Maximum Length, have: %s, maxWebsiteLen: %s", len(d.Website), MaxWebsiteLength)
	}
	if len(d.SecurityContact) > MaxSecurityContactLength {
		return d, errors.Errorf("[EnsureLength] Exceed Maximum Length, have: %s, maxSecurityContactLen: %s", len(d.SecurityContact), MaxSecurityContactLength)
	}
	if len(d.Details) > MaxDetailsLength {
		return d, errors.Errorf("[EnsureLength] Exceed Maximum Length, have: %s, maxDetailsLen: %s", len(d.Details), MaxDetailsLength)
	}

	return d, nil
}

func NewValidatorWrapper() *ValidatorWrapper {
	return &ValidatorWrapper{
		Validators: make(map[common.Address]*Validator),
		Delegations: make(map[common.Address]*Delegations),
	}
}

func (wrapper ValidatorWrapper) Amount() map[common.Address]*big.Int {
	amount := make(map[common.Address]*big.Int)
	for address, delegations := range wrapper.Delegations {
		amount[address] = delegations.Amount()
	}
	return amount
}

// UpdateDescription returns a new Description object with d1 as the base and the fields that's not empty in d2 updated
// accordingly. An error is returned if the resulting description fields have invalid length.
func UpdateDescription(d1, d2 Description) (Description, error) {
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

// CreateValidatorFromNewMsg creates validator from NewValidator message
func CreateValidatorFromNewMsg(val *CreateValidator) (*Validator, error) {
	desc, err := val.Description.EnsureLength()
	if err != nil {
		return nil, err
	}
	v := Validator{val.ValidatorAddress, desc}
	return &v, nil
}

// UpdateValidatorFromEditMsg updates validator from EditValidator message
func UpdateValidatorFromEditMsg(validator *Validator, edit *EditValidator) error {
	if validator.Address != edit.ValidatorAddress {
		return errAddressNotMatch
	}
	desc, err := UpdateDescription(validator.Description, edit.Description)
	if err != nil {
		return err
	}

	validator.Description = desc
	return nil
}