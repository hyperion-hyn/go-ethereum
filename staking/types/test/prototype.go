package staketest

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"math/big"
)

var (
	oneBig       = big.NewInt(1e18)
	tenKOnes     = new(big.Int).Mul(big.NewInt(10000), oneBig)
	twentyKOnes  = new(big.Int).Mul(big.NewInt(20000), oneBig)
	hundredKOnes = new(big.Int).Mul(big.NewInt(100000), oneBig)
	millionOnes = new(big.Int).Mul(big.NewInt(1000000), oneBig)

	// DefaultDelAmount is the default delegation amount
	DefaultDelAmount = new(big.Int).Set(twentyKOnes)

	// DefaultMinSelfDel is the default value of MinSelfDelegation
	DefaultMinSelfDel = new(big.Int).Set(tenKOnes)

	// DefaultMaxTotalDel is the default value of MaxTotalDelegation
	DefaultMaxTotalDel = new(big.Int).Set(millionOnes)
)

var (
	vWrapperPrototype = func() restaking.ValidatorWrapper_ {
		w := restaking.ValidatorWrapper_{
			Validator: validatorPrototype,
			Redelegations: restaking.NewRedelegationMap(),
			TotalDelegation:           big.NewInt(0).Set(DefaultDelAmount),
			TotalDelegationByOperator: big.NewInt(0).Set(DefaultDelAmount),
			BlockReward:               big.NewInt(0),
		}
		w.Counters.NumBlocksToSign = big.NewInt(0)
		w.Counters.NumBlocksSigned = big.NewInt(0)
		return w
	}()

	validatorPrototype = restaking.Validator_{
		ValidatorAddress:     common.Address{},
		OperatorAddresses:    restaking.NewEmptyAddressSet(),
		SlotPubKeys:          restaking.NewEmptyBLSKeys(),
		LastEpochInCommittee: big.NewInt(0),
		MaxTotalDelegation:   DefaultMaxTotalDel,
		Status:               uint8(restaking.Active),
		Commission:           commission,
		Description:          description,
		CreationHeight:       big.NewInt(0),
	}

	commissionRates = restaking.CommissionRates_{
		Rate:          common.NewDecWithPrec(5, 1),
		MaxRate:       common.NewDecWithPrec(9, 1),
		MaxChangeRate: common.NewDecWithPrec(3, 1),
	}

	commission = restaking.Commission_{
		CommissionRates: commissionRates,
		UpdateHeight:    big.NewInt(0),
	}

	description = restaking.Description_{
		Name:            "SuperHero",
		Identity:        "YouWouldNotKnow",
		Website:         "Secret Website",
		SecurityContact: "LicenseToKill",
		Details:         "blah blah blah",
	}
)

// GetDefaultValidator return the default staking.Validator for testing
func GetDefaultValidator() restaking.Validator_ {
	return CopyValidator(validatorPrototype)
}

// GetDefaultValidatorWithAddr return the default staking.Validator with the
// given validator address and bls keys
func GetDefaultValidatorWithAddr(validator, operator common.Address, pubs restaking.BLSPublicKeys_) restaking.Validator_ {
	v := CopyValidator(validatorPrototype)
	v.ValidatorAddress = validator
	v.OperatorAddresses.Put(operator)
	v.SlotPubKeys = CopySlotPubKeys(pubs)
	return v
}

// GetDefaultValidatorWrapper return the default staking.ValidatorWrapper for testing
func GetDefaultValidatorWrapper() restaking.ValidatorWrapper_ {
	return CopyValidatorWrapper(vWrapperPrototype)
}

// GetDefaultValidatorWrapperWithAddr return the default staking.ValidatorWrapper
// with the given validator address and bls keys.
func GetDefaultValidatorWrapperWithAddr(validator, operator common.Address, pubs restaking.BLSPublicKeys_) restaking.ValidatorWrapper_ {
	w := CopyValidatorWrapper(vWrapperPrototype)
	w.Validator.ValidatorAddress = validator
	w.Validator.OperatorAddresses.Put(operator)
	w.Validator.SlotPubKeys = CopySlotPubKeys(pubs)
	w.Redelegations.Put(operator, restaking.NewRedelegation(operator, big.NewInt(0).Set(DefaultDelAmount)))
	return w
}
