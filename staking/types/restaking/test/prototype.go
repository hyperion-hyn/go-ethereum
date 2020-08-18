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
	millionOnes  = new(big.Int).Mul(big.NewInt(1000000), oneBig)

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
			Validator:                 validatorPrototype,
			Redelegations:             restaking.NewRedelegationMap(),
			TotalDelegation:           big.NewInt(0),
			TotalDelegationByOperator: big.NewInt(0),
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

// GetDefaultValidatorWrapper return the default staking.ValidatorWrapper for testing
func GetDefaultValidatorWrapper() restaking.ValidatorWrapper_ {
	return CopyValidatorWrapper(vWrapperPrototype)
}

type ValidatorWrapperBuilder struct {
	wrapper restaking.ValidatorWrapper_
}

func (b *ValidatorWrapperBuilder) SetValidatorAddress(validator common.Address) *ValidatorWrapperBuilder {
	b.wrapper.Validator.ValidatorAddress = validator
	return b
}

func (b *ValidatorWrapperBuilder) AddOperatorAddress(operator common.Address) *ValidatorWrapperBuilder {
	b.wrapper.Validator.OperatorAddresses.Put(operator)
	return b
}

func (b *ValidatorWrapperBuilder) AddSlotPubKey(key restaking.BLSPublicKey_) *ValidatorWrapperBuilder {
	b.wrapper.Validator.SlotPubKeys.Keys = append(b.wrapper.Validator.SlotPubKeys.Keys, &key)
	return b
}

func (b *ValidatorWrapperBuilder) AddSlotPubKeys(keys restaking.BLSPublicKeys_) *ValidatorWrapperBuilder {
	for _, k := range keys.Keys {
		b.AddSlotPubKey(*k)
	}
	return b
}

func (b *ValidatorWrapperBuilder) SetLastEpochInCommittee(lastEpochInCommittee *big.Int) *ValidatorWrapperBuilder {
	b.wrapper.Validator.LastEpochInCommittee = big.NewInt(0).Set(lastEpochInCommittee)
	return b
}

func (b *ValidatorWrapperBuilder) SetMaxTotalDelegation(maxTotalDelegation *big.Int) *ValidatorWrapperBuilder {
	b.wrapper.Validator.MaxTotalDelegation = big.NewInt(0).Set(maxTotalDelegation)
	return b
}

func (b *ValidatorWrapperBuilder) SetStatus(status uint8) *ValidatorWrapperBuilder {
	b.wrapper.Validator.Status = status
	return b
}

func (b *ValidatorWrapperBuilder) SetCommission(commission restaking.Commission_) *ValidatorWrapperBuilder {
	b.wrapper.Validator.Commission = CopyCommission(commission)
	return b
}

func (b *ValidatorWrapperBuilder) SetDescription(description restaking.Description_) *ValidatorWrapperBuilder {
	b.wrapper.Validator.Description = description
	return b
}

func (b *ValidatorWrapperBuilder) SetCreationHeight(creationHeight *big.Int) *ValidatorWrapperBuilder {
	b.wrapper.Validator.CreationHeight = big.NewInt(0).Set(creationHeight)
	return b
}

func (b *ValidatorWrapperBuilder) AddRedelegation(redelegation restaking.Redelegation_) *ValidatorWrapperBuilder {
	b.wrapper.Redelegations.Put(redelegation.DelegatorAddress, CopyRedelegation(redelegation))
	b.wrapper.TotalDelegation.Add(b.wrapper.TotalDelegation, redelegation.Amount)
	if b.wrapper.Validator.OperatorAddresses.Contain(redelegation.DelegatorAddress) {
		b.wrapper.TotalDelegationByOperator.Add(b.wrapper.TotalDelegationByOperator, redelegation.Amount)
	}
	return b
}

func (b *ValidatorWrapperBuilder) SetNumBlocksToSign(numBlocksToSign *big.Int) *ValidatorWrapperBuilder {
	b.wrapper.Counters.NumBlocksToSign = big.NewInt(0).Set(numBlocksToSign)
	return b
}

func (b *ValidatorWrapperBuilder) SetNumBlocksSigned(numBlocksSigned *big.Int) *ValidatorWrapperBuilder {
	b.wrapper.Counters.NumBlocksSigned = big.NewInt(0).Set(numBlocksSigned)
	return b
}

func (b *ValidatorWrapperBuilder) SetBlockReward(blockReward *big.Int) *ValidatorWrapperBuilder {
	b.wrapper.BlockReward = big.NewInt(0).Set(blockReward)
	return b
}

func (b *ValidatorWrapperBuilder) Build() restaking.ValidatorWrapper_ {
	return b.wrapper
}

func NewValidatorWrapperBuilder() *ValidatorWrapperBuilder {
	return &ValidatorWrapperBuilder{wrapper: CopyValidatorWrapper(vWrapperPrototype)}
}