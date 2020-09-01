package staketest

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"math/big"
	"testing"
)

func TestGetDefaultValidator(t *testing.T) {
	v := GetDefaultValidator()
	if err := assertValidatorDeepCopy(v, validatorPrototype); err != nil {
		t.Error(err)
	}
}

func TestGetDefaultValidatorWrapper(t *testing.T) {
	w := GetDefaultValidatorWrapper()
	if err := assertValidatorWrapperDeepCopy(w, vWrapperPrototype); err != nil {
		t.Error(err)
	}
}

func TestValidatorWrapperBuilder(t *testing.T) {
	tests := []struct {
		validatorAddr        common.Address
		operatorAddr         common.Address
		key                  restaking.BLSPublicKey_
		lastEpochInCommittee *big.Int
		maxTotalDelegation   *big.Int
		commission           restaking.Commission_
		description          restaking.Description_
		creationHeight       *big.Int
		redelegation         restaking.Redelegation_
		blockReward          *big.Int
		numBlocksToSign      *big.Int
		numBlocksSigned      *big.Int
	}{
		{
			validatorAddr:        common.BigToAddress(common.Big1),
			operatorAddr:         common.BigToAddress(common.Big1),
			key:                  testPub,
			lastEpochInCommittee: big.NewInt(1),
			maxTotalDelegation:   big.NewInt(1),
			commission:           nonZeroCommission,
			description: restaking.Description_{
				Name:            "A",
				Identity:        "B",
				Website:         "C",
				SecurityContact: "D",
				Details:         "E",
			},
			creationHeight:  big.NewInt(1),
			redelegation:    nonZeroDelegation,
			blockReward:     big.NewInt(1),
			numBlocksToSign: big.NewInt(1),
			numBlocksSigned: big.NewInt(1),
		},
	}

	for i, test := range tests {
		v := NewValidatorWrapperBuilder().
			SetValidatorAddress(test.validatorAddr).
			AddOperatorAddress(test.operatorAddr).
			AddSlotPubKey(test.key).
			SetLastEpochInCommittee(test.lastEpochInCommittee).
			SetMaxTotalDelegation(test.maxTotalDelegation).
			SetCommission(test.commission).
			SetDescription(test.description).
			SetCreationHeight(test.creationHeight).
			AddRedelegation(test.redelegation).
			SetBlockReward(test.blockReward).
			SetNumBlocksToSign(test.numBlocksToSign).
			SetNumBlocksSigned(test.numBlocksSigned).
			Build()

		exp := GetDefaultValidatorWrapper()
		exp.Validator.ValidatorAddress = test.validatorAddr
		exp.Validator.OperatorAddresses = restaking.NewAddressSetWithAddress(test.operatorAddr)
		exp.Validator.SlotPubKeys.Keys = append(exp.Validator.SlotPubKeys.Keys, &test.key)
		exp.Validator.LastEpochInCommittee = test.lastEpochInCommittee
		exp.Validator.MaxTotalDelegation = test.maxTotalDelegation
		exp.Validator.Commission = test.commission
		exp.Validator.Description = test.description
		exp.Validator.CreationHeight = test.creationHeight
		exp.Redelegations.Put(test.redelegation.DelegatorAddress, test.redelegation)
		exp.BlockReward = test.blockReward
		exp.Counters.NumBlocksToSign = test.numBlocksToSign
		exp.Counters.NumBlocksSigned = test.numBlocksSigned
		exp.TotalDelegation = common.Big1
		exp.TotalDelegationByOperator = common.Big1

		if err := assertValidatorWrapperDeepCopy(v, exp); err != nil {
			t.Errorf("Test %v: %v", i, err)
		}
	}
}
