package types

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/numeric"
	"github.com/ethereum/go-ethereum/staking/reward"
	"github.com/pkg/errors"
	"math/big"
)

//go:generate gencodec -type Validator -field-override validatorMarshaling -out gen_validator.go

// Define validator staking related const
const (
	MaxNameLength            = 140
	MaxIdentityLength        = 140
	MaxWebsiteLength         = 140
	MaxSecurityContactLength = 140
	MaxDetailsLength         = 280
	MinSelfDelegation        = 1000000		// TODO to be confirmed
)

var (
	StakingInfoAddress = common.StringToAddress("0x0000000000000000000000000000000123456789") // used to save staking state in state db

	errAddressNotMatch = errors.New("Validator key not match")
)

type ValidatorContainer struct {
	Validators []ValidatorWrapper `json:"validators"`
}

// ValidatorWrapper contains one validator and its delegation information
type ValidatorWrapper struct {
	Validator   *Validator  `json:"validator"`
	Delegations Delegations `json:"delegations"`
}

// Validator - data fields for a validator
type Validator struct {
	// ECDSA address of the validator
	Address common.Address `json:"validator_address" gencodec:"required"`
	// description for the validator
	Description Description `json:"description"`
	// validator's self declared minimum self delegation
	MinSelfDelegation *big.Int `json:"min_self_delegation"`
	// maximum total delegation allowed
	MaxTotalDelegation *big.Int `json:"max_total_delegation"`
	// Is the validator active in participating committee selection process or not
	Active bool `json:"-"`
	// commission parameters
	CommissionRate numeric.Dec `json:"commission_rate" gencodec:"required"`
	// CreationHeight is the height of creation
	CreationHeight *big.Int `json:"-"`
	// Banned records whether this validator is banned from the network because they double-signed
	Banned bool `json:"-"`
}

// field type overrides for gencodec
type validatorMarshaling struct {
	Address            common.UnprefixedAddress
	MinSelfDelegation  *math.HexOrDecimal256
	MaxTotalDelegation *math.HexOrDecimal256
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

// String returns a human readable string representation of a validator.
func (v *Validator) String() string {
	return fmt.Sprintf(`Validator
  Address:                    %s
  Minimum Self Delegation:     %v
  Maximum Total Delegation:     %v
  Description:                %v
  Commission:                 %v`,
		v.Address.String(), v.MinSelfDelegation, v.MaxTotalDelegation, v.Description, v.CommissionRate,
	)
}

func (wrapper ValidatorWrapper) TotalDelegation() *big.Int {
	amount := big.NewInt(0)
	for _, delegation := range wrapper.Delegations {
		amount = amount.Add(amount, delegation.Amount)
	}
	return amount
}

func (container ValidatorContainer) IsValidator(addr common.Address) bool {
	return container.Validator(addr) != nil
}

func (container ValidatorContainer) Validator(addr common.Address) *ValidatorWrapper {
	// TODO search algo?
	for _, val := range container.Validators {
		if addr == val.Validator.Address {
			return &val
		}
	}
	return nil
}

// Is the validator's self delegation more than MinSelfDelegation or not
func (wrapper ValidatorWrapper) IsStakingValid() bool {
	return wrapper.Delegations[0].Amount.Cmp(big.NewInt(MinSelfDelegation)) >= 0
}

func (wrapper ValidatorWrapper) AddReward(reward *big.Int, distributor reward.Distributor)  {
	rewardPool := big.NewInt(0).Set(reward)

	// Payout commission
	commissionInt := wrapper.Validator.CommissionRate.MulInt(reward).RoundInt()
	wrapper.Delegations[0].Reward.Add(wrapper.Delegations[0].Reward, commissionInt)
	rewardPool.Sub(rewardPool, commissionInt)

	// Payout each delegator's reward
	totalRewardForDelegators := big.NewInt(0).Set(rewardPool)
	for i := range wrapper.Delegations {
		delegation := wrapper.Delegations[i]
		percentage := distributor.GetPercentage(delegation.DelegatorAddress)
		rewardInt := percentage.MulInt(totalRewardForDelegators).RoundInt()
		delegation.Reward.Add(delegation.Reward, rewardInt)
		rewardPool.Sub(rewardPool, rewardInt)
	}

	// The last remaining bit belongs to the validator (remember the validator's self delegation is always at index 0)
	if rewardPool.Cmp(big.NewInt(0)) > 0 {
		wrapper.Delegations[0].Reward.Add(wrapper.Delegations[0].Reward, rewardPool)
	}
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
func CreateValidatorFromNewMsg(val *CreateValidator, blockNum *big.Int) (*Validator, error) {
	desc, err := val.Description.EnsureLength()
	if err != nil {
		return nil, err
	}
	v := Validator{
		Address:            val.ValidatorAddress,
		MinSelfDelegation:  val.MinSelfDelegation,
		MaxTotalDelegation: val.MaxTotalDelegation,
		Active:             true,
		CommissionRate:     val.CommissionRate,
		Description:        desc,
		CreationHeight:     blockNum,
		Banned:             false,
	}
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

	if !edit.CommissionRate.IsNil() {
		validator.CommissionRate = edit.CommissionRate
	}

	if edit.MinSelfDelegation != nil && edit.MinSelfDelegation.Sign() != 0 {
		validator.MinSelfDelegation = edit.MinSelfDelegation
	}

	if edit.MaxTotalDelegation != nil && edit.MaxTotalDelegation.Sign() != 0 {
		validator.MaxTotalDelegation = edit.MaxTotalDelegation
	}

	if edit.Active != nil {
		validator.Active = *edit.Active
	}

	return nil
}
