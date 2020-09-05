package restaking

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"testing"
)

var (
	testPub = BLSPublicKey_{Key: [48]byte{1}}

	nonZeroCommissionRates = CommissionRates_{
		Rate:          common.NewDecWithPrec(1, 1),
		MaxRate:       common.NewDecWithPrec(2, 1),
		MaxChangeRate: common.NewDecWithPrec(3, 1),
	}

	zeroCommissionRates = CommissionRates_{
		Rate:          common.ZeroDec(),
		MaxRate:       common.ZeroDec(),
		MaxChangeRate: common.ZeroDec(),
	}

	nonZeroCommission = Commission_{
		CommissionRates: nonZeroCommissionRates,
		UpdateHeight:    common.Big1,
	}

	zeroCommission = Commission_{
		CommissionRates: zeroCommissionRates,
		UpdateHeight:    common.Big1,
	}

	nonZeroDelegation = Redelegation_{
		DelegatorAddress: common.BigToAddress(common.Big1),
		Amount:           common.Big1,
		Reward:           common.Big2,
		Undelegation: Undelegation_{
			Amount: common.Big1,
			Epoch:  common.Big2,
		},
	}

	zeroDelegation = Redelegation_{
		Amount: common.Big0,
		Reward: common.Big0,
		Undelegation: Undelegation_{
			Amount: common.Big0,
			Epoch:  common.Big0,
		},
	}
)

func TestValidatorWrapperBuilder(t *testing.T) {
	tests := []struct {
		validatorAddr        common.Address
		operatorAddr         common.Address
		key                  BLSPublicKey_
		lastEpochInCommittee *big.Int
		maxTotalDelegation   *big.Int
		commission           Commission_
		description          Description_
		creationHeight       *big.Int
		redelegation         Redelegation_
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
			description: Description_{
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
		exp.Validator.OperatorAddresses = NewAddressSetWithAddress(test.operatorAddr)
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

		if err := CheckValidatorWrapperEqual(*v, exp); err != nil {
			t.Errorf("Test %v: %v", i, err)
		}
	}
}

// makeNonZeroValidator makes a valid Validator data structure
func makeNonZeroValidator() Validator_ {
	d := Description_{
		Name:            "Wayne",
		Identity:        "wen",
		Website:         "harmony.one.wen",
		Details:         "best",
		SecurityContact: "sc",
	}
	v := Validator_{
		ValidatorAddress:     common.BigToAddress(common.Big1),
		OperatorAddresses:    NewAddressSetWithAddress(common.BigToAddress(common.Big1)),
		SlotPubKeys:          NewBLSKeysWithBLSKey(testPub),
		LastEpochInCommittee: big.NewInt(20),
		MaxTotalDelegation:   common.Big1,
		Status:               uint8(Active),
		Commission:           nonZeroCommission,
		Description:          d,
		CreationHeight:       big.NewInt(12306),
	}
	return v
}

func makeNonZeroValidatorWrapper() ValidatorWrapper_ {
	w := ValidatorWrapper_{
		Validator: makeNonZeroValidator(),
		Redelegations: func() RedelegationMap_ {
			m := NewRedelegationMap()
			m.Put(nonZeroDelegation.DelegatorAddress, nonZeroDelegation)
			m.Put(zeroDelegation.DelegatorAddress, zeroDelegation)
			return m
		}(),
		BlockReward: common.Big1,
	}
	w.Counters.NumBlocksToSign = common.Big1
	w.Counters.NumBlocksSigned = common.Big2
	return w
}
