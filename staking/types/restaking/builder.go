package restaking

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type ValidatorWrapperBuilder struct {
	wrapper ValidatorWrapper_
}

func (b *ValidatorWrapperBuilder) SetValidatorAddress(validator common.Address) *ValidatorWrapperBuilder {
	b.wrapper.Validator.ValidatorAddress = validator
	return b
}

func (b *ValidatorWrapperBuilder) AddOperatorAddress(operator common.Address) *ValidatorWrapperBuilder {
	b.wrapper.Validator.OperatorAddresses.Put(operator)
	return b
}

func (b *ValidatorWrapperBuilder) AddSlotPubKey(key BLSPublicKey_) *ValidatorWrapperBuilder {
	b.wrapper.Validator.SlotPubKeys.Keys = append(b.wrapper.Validator.SlotPubKeys.Keys, &key)
	return b
}

func (b *ValidatorWrapperBuilder) AddSlotPubKeys(keys BLSPublicKeys_) *ValidatorWrapperBuilder {
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

func (b *ValidatorWrapperBuilder) SetStatus(status ValidatorStatus) *ValidatorWrapperBuilder {
	b.wrapper.Validator.Status = uint8(status)
	return b
}

func (b *ValidatorWrapperBuilder) SetCommission(commission Commission_) *ValidatorWrapperBuilder {
	c, _ := commission.Copy()
	b.wrapper.Validator.Commission = *c
	return b
}

func (b *ValidatorWrapperBuilder) SetDescription(description Description_) *ValidatorWrapperBuilder {
	b.wrapper.Validator.Description = description
	return b
}

func (b *ValidatorWrapperBuilder) SetCreationHeight(creationHeight *big.Int) *ValidatorWrapperBuilder {
	b.wrapper.Validator.CreationHeight = big.NewInt(0).Set(creationHeight)
	return b
}

func (b *ValidatorWrapperBuilder) AddRedelegation(redelegation Redelegation_) *ValidatorWrapperBuilder {
	b.wrapper.Redelegations.Put(redelegation.DelegatorAddress, redelegation)
	b.wrapper.TotalDelegation.Add(b.wrapper.TotalDelegation, redelegation.Amount)
	if b.wrapper.Validator.OperatorAddresses.Contain(redelegation.DelegatorAddress) {
		b.wrapper.TotalDelegationFromOperators.Add(b.wrapper.TotalDelegationFromOperators, redelegation.Amount)
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

func (b *ValidatorWrapperBuilder) Build() *ValidatorWrapper_ {
	return &b.wrapper
}

func NewValidatorWrapperBuilder() *ValidatorWrapperBuilder {
	return &ValidatorWrapperBuilder{
		wrapper: ValidatorWrapper_{
			Validator: Validator_{
				ValidatorAddress:     common.Address{},
				OperatorAddresses:    NewEmptyAddressSet(),
				SlotPubKeys:          NewEmptyBLSKeys(),
				LastEpochInCommittee: big.NewInt(0),
				MaxTotalDelegation:   big.NewInt(0),
				Status:               uint8(Active),
				Commission: Commission_{
					CommissionRates: CommissionRates_{
						Rate:          common.NewDec(0),
						MaxRate:       common.NewDec(0),
						MaxChangeRate: common.NewDec(0),
					},
					UpdateHeight: big.NewInt(0),
				},
				Description:    Description_{},
				CreationHeight: big.NewInt(0),
			},
			Redelegations: NewRedelegationMap(),
			Counters: Counters_{
				NumBlocksToSign: big.NewInt(0),
				NumBlocksSigned: big.NewInt(0),
			},
			BlockReward:                  big.NewInt(0),
			TotalDelegation:              big.NewInt(0),
			TotalDelegationFromOperators: big.NewInt(0),
		},
	}
}
