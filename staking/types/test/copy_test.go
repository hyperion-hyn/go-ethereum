package staketest

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/staking/effective"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"math/big"
	"reflect"
	"testing"
)

var (
	testPub = restaking.BLSPublicKey_{Key: [48]byte{1}}
)

func TestCopyValidatorWrapper(t *testing.T) {
	tests := []struct {
		w restaking.ValidatorWrapper_
	}{
		{makeNonZeroValidatorWrapper()},
		{makeZeroValidatorWrapper()},
		{restaking.ValidatorWrapper_{}},
	}
	for i, test := range tests {
		cp := CopyValidatorWrapper(test.w)

		if err := assertValidatorWrapperDeepCopy(cp, test.w); err != nil {
			t.Errorf("Test %v: %v", i, err)
		}
	}
}

func makeNonZeroValidatorWrapper() restaking.ValidatorWrapper_ {
	w := restaking.ValidatorWrapper_{
		Validator: makeNonZeroValidator(),
		Redelegations: func() restaking.RedelegationMap_ {
			m := restaking.NewRelegationMap()
			m.Put(nonZeroDelegation.DelegatorAddress, &nonZeroDelegation)
			m.Put(zeroDelegation.DelegatorAddress, &zeroDelegation)
			return m
		}(),
		BlockReward: common.Big1,
	}
	w.Counters.NumBlocksToSign = common.Big1
	w.Counters.NumBlocksSigned = common.Big2
	return w
}

func makeZeroValidatorWrapper() restaking.ValidatorWrapper_ {
	w := restaking.ValidatorWrapper_{
		Redelegations: restaking.NewRelegationMap(),
		BlockReward:   common.Big0,
	}
	w.Counters.NumBlocksSigned = common.Big0
	w.Counters.NumBlocksToSign = common.Big0
	return w
}

func TestCopyValidator(t *testing.T) {
	tests := []struct {
		v restaking.Validator_
	}{
		{makeNonZeroValidator()},
		{makeZeroValidator()},
		{restaking.Validator_{}},
	}
	for i, test := range tests {
		cp := CopyValidator(test.v)
		if err := assertValidatorDeepCopy(test.v, cp); err != nil {
			t.Errorf("Test %v: %v", i, err)
		}
	}
}

// makeNonZeroValidator makes a valid Validator data structure
func makeNonZeroValidator() restaking.Validator_ {
	d := restaking.Description_{
		Name:     "Wayne",
		Identity: "wen",
		Website:  "harmony.one.wen",
		Details:  "best",
	}
	v := restaking.Validator_{
		ValidatorAddress:     common.BigToAddress(common.Big0),
		OperatorAddresses:    restaking.NewAddressSetWithAddress(common.BigToAddress(common.Big1)),
		SlotPubKeys:          restaking.BLSPublicKeys_{Keys: []*restaking.BLSPublicKey_{&testPub}},
		LastEpochInCommittee: big.NewInt(20),
		MaxTotalDelegation:   common.Big1,
		Status:               big.NewInt(int64(effective.Active)),
		Commission:           nonZeroCommission,
		Description:          d,
		CreationHeight:       big.NewInt(12306),
	}
	return v
}

func makeZeroValidator() restaking.Validator_ {
	v := restaking.Validator_{
		SlotPubKeys:          restaking.BLSPublicKeys_{},
		LastEpochInCommittee: common.Big0,
		MaxTotalDelegation:   common.Big0,
		Commission:           zeroCommission,
		CreationHeight:       common.Big0,
	}
	return v
}

func TestCopyCommission(t *testing.T) {
	tests := []struct {
		c restaking.Commission_
	}{
		{nonZeroCommission},
		{zeroCommission},
		{restaking.Commission_{}},
	}
	for i, test := range tests {
		cp := CopyCommission(test.c)

		if err := assertCommissionDeepCopy(cp, test.c); err != nil {
			t.Errorf("Test %v: %v", i, err)
		}
	}
}

func TestCopyRedelegation(t *testing.T) {
	tests := []struct {
		d restaking.Redelegation_
	}{
		{nonZeroDelegation},
		{zeroDelegation},
		{restaking.Redelegation_{}},
	}
	for i, test := range tests {
		cp := CopyDelegation(test.d)
		if err := assertRedelegationDeepCopy(cp, test.d); err != nil {
			t.Errorf("Test %v: %v", i, err)
		}
	}
}

var (
	nonZeroCommissionRates = restaking.CommissionRates_{
		Rate:          common.NewDecWithPrec(1, 1),
		MaxRate:       common.NewDecWithPrec(2, 1),
		MaxChangeRate: common.NewDecWithPrec(3, 1),
	}

	zeroCommissionRates = restaking.CommissionRates_{
		Rate:          common.ZeroDec(),
		MaxRate:       common.ZeroDec(),
		MaxChangeRate: common.ZeroDec(),
	}

	nonZeroCommission = restaking.Commission_{
		CommissionRates: nonZeroCommissionRates,
		UpdateHeight:    common.Big1,
	}

	zeroCommission = restaking.Commission_{
		CommissionRates: zeroCommissionRates,
		UpdateHeight:    common.Big1,
	}

	nonZeroDelegation = restaking.Redelegation_{
		DelegatorAddress: common.BigToAddress(common.Big1),
		Amount:           common.Big1,
		Reward:           common.Big2,
		Undelegation: restaking.Undelegation_{
			Amount: common.Big1,
			Epoch:  common.Big2,
		},
	}

	zeroDelegation = restaking.Redelegation_{
		Amount: common.Big0,
		Reward: common.Big0,
	}
)

func assertValidatorWrapperDeepCopy(w1, w2 restaking.ValidatorWrapper_) error {
	if err := assertValidatorDeepCopy(w1.Validator, w2.Validator); err != nil {
		return fmt.Errorf("validator %v", err)
	}
	if err := assertRedelegationMapDeepCopy(w1.Redelegations, w2.Redelegations); err != nil {
		return fmt.Errorf("delegations %v", err)
	}
	if err := assertBigIntCopy(w1.Counters.NumBlocksToSign, w2.Counters.NumBlocksToSign); err != nil {
		return fmt.Errorf("counters %v", err)
	}
	if err := assertBigIntCopy(w1.Counters.NumBlocksSigned, w2.Counters.NumBlocksSigned); err != nil {
		return fmt.Errorf("counters %v", err)
	}
	if w1.BlockReward != nil && w1.BlockReward == w2.BlockReward {
		return fmt.Errorf("BlockReward same address")
	}
	return nil
}

func assertValidatorDeepCopy(v1, v2 restaking.Validator_) error {
	if !reflect.DeepEqual(v1, v2) {
		return fmt.Errorf("not deep equal")
	}
	if &v1.SlotPubKeys == &v2.SlotPubKeys {
		return fmt.Errorf("SlotPubKeys same pointer")
	}
	for i := range v1.SlotPubKeys.Keys {
		if v1.SlotPubKeys.Keys[i].Hex() == v2.SlotPubKeys.Keys[i].Hex() {
			return fmt.Errorf("SlotPubKeys[%v] same address", i)
		}
	}
	if v1.LastEpochInCommittee != nil && v1.LastEpochInCommittee == v2.LastEpochInCommittee {
		return fmt.Errorf("LastEpochInCommittee same address")
	}
	if v1.CreationHeight != nil && v1.CreationHeight == v2.CreationHeight {
		return fmt.Errorf("CreationHeight same address")
	}
	if &v1.Description == &v2.Description {
		return fmt.Errorf("same description")
	}
	if err := assertCommissionDeepCopy(v1.Commission, v2.Commission); err != nil {
		return fmt.Errorf("CommissionRates: %v", err)
	}
	return nil
}

func assertCommissionDeepCopy(c1, c2 restaking.Commission_) error {
	if !reflect.DeepEqual(c1, c2) {
		return fmt.Errorf("not deep equal")
	}
	if err := assertCommissionRatesCopy(c1.CommissionRates, c2.CommissionRates); err != nil {
		return fmt.Errorf("CommissionRates: %v", err)
	}
	if err := assertBigIntCopy(c1.UpdateHeight, c2.UpdateHeight); err != nil {
		return fmt.Errorf("UpdateHeight: %v", err)
	}
	return nil
}

func assertCommissionRatesCopy(cr1, cr2 restaking.CommissionRates_) error {
	if err := assertDecCopy(cr1.Rate, cr2.Rate); err != nil {
		return fmt.Errorf("rate: %v", err)
	}
	if err := assertDecCopy(cr1.MaxRate, cr2.MaxRate); err != nil {
		return fmt.Errorf("maxRate: %v", err)
	}
	if err := assertDecCopy(cr1.MaxChangeRate, cr2.MaxChangeRate); err != nil {
		return fmt.Errorf("maxChangeRate: %v", err)
	}
	return nil
}

func assertRedelegationMapDeepCopy(ds1, ds2 restaking.RedelegationMap_) error {
	if !reflect.DeepEqual(ds1, ds2) {
		return fmt.Errorf("not deep equal")
	}
	for _, key := range ds1.Keys {
		if err := assertRedelegationDeepCopy(*ds1.Get(*key), *ds2.Get(*key)); err != nil {
			return fmt.Errorf("[%v]: %v", key, err)
		}
	}
	return nil
}

func assertRedelegationDeepCopy(d1, d2 restaking.Redelegation_) error {
	if !reflect.DeepEqual(d1, d2) {
		return fmt.Errorf("not deep equal")
	}
	if d1.Amount != nil && d1.Amount == d2.Amount {
		return fmt.Errorf("amount same address")
	}
	if d1.Reward != nil && d1.Reward == d2.Reward {
		return fmt.Errorf("reward same address")
	}
	if err := assertUndelegationDeepCopy(d1.Undelegation, d2.Undelegation); err != nil {
		return fmt.Errorf("undelegations %v", err)
	}
	return nil
}

func assertUndelegationDeepCopy(ud1, ud2 restaking.Undelegation_) error {
	if !reflect.DeepEqual(ud1, ud2) {
		return fmt.Errorf("not deep equal")
	}
	if err := assertBigIntCopy(ud1.Amount, ud2.Amount); err != nil {
		return fmt.Errorf("amount %v", err)
	}
	if err := assertBigIntCopy(ud1.Epoch, ud2.Epoch); err != nil {
		return fmt.Errorf("epoch: %v", err)
	}
	return nil
}

func assertDecCopy(d1, d2 common.Dec) error {
	if d1.IsNil() != d2.IsNil() {
		return errors.New("IsNil not equal")
	}
	if d1.IsNil() {
		return nil
	}
	if d1 == d2 {
		return errors.New("same address")
	}
	return nil
}

func assertBigIntCopy(i1, i2 *big.Int) error {
	if (i1 == nil) != (i2 == nil) {
		return errors.New("is nil not equal")
	}
	if i1 != nil && i1 == i2 {
		return errors.New("not copy")
	}
	return nil
}
