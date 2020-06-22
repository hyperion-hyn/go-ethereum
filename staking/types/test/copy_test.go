package staketest

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/numeric"
	"github.com/ethereum/go-ethereum/staking/effective"
	"github.com/ethereum/go-ethereum/staking/types"
	"math/big"
	"reflect"
	"testing"
)

var (
	testPub = types.BLSPublicKey{1}
)

func TestCopyValidatorWrapper(t *testing.T) {
	tests := []struct {
		w types.ValidatorWrapper
	}{
		{makeNonZeroValidatorWrapper()},
		{makeZeroValidatorWrapper()},
		{types.ValidatorWrapper{}},
	}
	for i, test := range tests {
		cp := CopyValidatorWrapper(test.w)

		if err := assertValidatorWrapperDeepCopy(cp, test.w); err != nil {
			t.Errorf("Test %v: %v", i, err)
		}
	}
}

func makeNonZeroValidatorWrapper() types.ValidatorWrapper {
	w := types.ValidatorWrapper{
		Validator:     makeNonZeroValidator(),
		Redelegations: types.Redelegations{nonZeroDelegation, zeroDelegation},
		BlockReward:   common.Big1,
	}
	w.Counters.NumBlocksToSign = common.Big1
	w.Counters.NumBlocksSigned = common.Big2
	return w
}

func makeZeroValidatorWrapper() types.ValidatorWrapper {
	w := types.ValidatorWrapper{
		Redelegations: make(types.Redelegations, 0),
		BlockReward:   common.Big0,
	}
	w.Counters.NumBlocksSigned = common.Big0
	w.Counters.NumBlocksToSign = common.Big0
	return w
}

func TestCopyValidator(t *testing.T) {
	tests := []struct {
		v types.Validator
	}{
		{makeNonZeroValidator()},
		{makeZeroValidator()},
		{types.Validator{}},
	}
	for i, test := range tests {
		cp := CopyValidator(test.v)
		if err := assertValidatorDeepCopy(test.v, cp); err != nil {
			t.Errorf("Test %v: %v", i, err)
		}
	}
}

// makeNonZeroValidator makes a valid Validator data structure
func makeNonZeroValidator() types.Validator {
	d := types.Description{
		Name:     "Wayne",
		Identity: "wen",
		Website:  "harmony.one.wen",
		Details:  "best",
	}
	v := types.Validator{
		Address:              common.BigToAddress(common.Big0),
		SlotPubKeys:          []types.BLSPublicKey{testPub},
		LastEpochInCommittee: big.NewInt(20),
		MinSelfDelegation:    common.Big1,
		MaxTotalDelegation:   common.Big1,
		Status:               effective.Active,
		Commission:           nonZeroCommission,
		Description:          d,
		CreationHeight:       big.NewInt(12306),
	}
	return v
}

func makeZeroValidator() types.Validator {
	v := types.Validator{
		SlotPubKeys:          make([]types.BLSPublicKey, 0),
		LastEpochInCommittee: common.Big0,
		MinSelfDelegation:    common.Big0,
		MaxTotalDelegation:   common.Big0,
		Commission:           zeroCommission,
		CreationHeight:       common.Big0,
	}
	return v
}

func TestCopyCommission(t *testing.T) {
	tests := []struct {
		c types.Commission
	}{
		{nonZeroCommission},
		{zeroCommission},
		{types.Commission{}},
	}
	for i, test := range tests {
		cp := CopyCommission(test.c)

		if err := assertCommissionDeepCopy(cp, test.c); err != nil {
			t.Errorf("Test %v: %v", i, err)
		}
	}
}

func TestCopyDelegation(t *testing.T) {
	tests := []struct {
		d types.Redelegation
	}{
		{nonZeroDelegation},
		{zeroDelegation},
		{types.Redelegation{}},
	}
	for i, test := range tests {
		cp := CopyDelegation(test.d)
		if err := assertDelegationDeepCopy(cp, test.d); err != nil {
			t.Errorf("Test %v: %v", i, err)
		}
	}
}

var (
	nonZeroCommissionRates = types.CommissionRates{
		Rate:          numeric.NewDecWithPrec(1, 1),
		MaxRate:       numeric.NewDecWithPrec(2, 1),
		MaxChangeRate: numeric.NewDecWithPrec(3, 1),
	}

	zeroCommissionRates = types.CommissionRates{
		Rate:          numeric.ZeroDec(),
		MaxRate:       numeric.ZeroDec(),
		MaxChangeRate: numeric.ZeroDec(),
	}

	nonZeroCommission = types.Commission{
		CommissionRates: nonZeroCommissionRates,
		UpdateHeight:    common.Big1,
	}

	zeroCommission = types.Commission{
		CommissionRates: zeroCommissionRates,
		UpdateHeight:    common.Big1,
	}

	nonZeroDelegation = types.Redelegation{
		DelegatorAddress: common.BigToAddress(common.Big1),
		Amount:           common.Big1,
		Reward:           common.Big2,
		Undelegations: types.Undelegations{
			types.Undelegation{
				Amount: common.Big1,
				Epoch:  common.Big2,
			},
			types.Undelegation{
				Amount: common.Big3,
				Epoch:  common.Big1,
			},
		},
	}

	zeroDelegation = types.Redelegation{
		Amount:        common.Big0,
		Reward:        common.Big0,
		Undelegations: make(types.Undelegations, 0),
	}
)

func assertValidatorWrapperDeepCopy(w1, w2 types.ValidatorWrapper) error {
	if err := assertValidatorDeepCopy(w1.Validator, w2.Validator); err != nil {
		return fmt.Errorf("validator %v", err)
	}
	if err := assertDelegationsDeepCopy(w1.Redelegations, w2.Redelegations); err != nil {
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

func assertValidatorDeepCopy(v1, v2 types.Validator) error {
	if !reflect.DeepEqual(v1, v2) {
		return fmt.Errorf("not deep equal")
	}
	if &v1.SlotPubKeys == &v2.SlotPubKeys {
		return fmt.Errorf("SlotPubKeys same pointer")
	}
	for i := range v1.SlotPubKeys {
		if &v1.SlotPubKeys[i] == &v2.SlotPubKeys[i] {
			return fmt.Errorf("SlotPubKeys[%v] same address", i)
		}
	}
	if v1.LastEpochInCommittee != nil && v1.LastEpochInCommittee == v2.LastEpochInCommittee {
		return fmt.Errorf("LastEpochInCommittee same address")
	}
	if v1.MinSelfDelegation != nil && v1.MinSelfDelegation == v2.MinSelfDelegation {
		return fmt.Errorf("MinSelfDelegation same address")
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

func assertCommissionDeepCopy(c1, c2 types.Commission) error {
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

func assertCommissionRatesCopy(cr1, cr2 types.CommissionRates) error {
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

func assertDelegationsDeepCopy(ds1, ds2 types.Redelegations) error {
	if !reflect.DeepEqual(ds1, ds2) {
		return fmt.Errorf("not deep equal")
	}
	for i := range ds1 {
		if err := assertDelegationDeepCopy(ds1[i], ds2[i]); err != nil {
			return fmt.Errorf("[%v]: %v", i, err)
		}
	}
	return nil
}

func assertDelegationDeepCopy(d1, d2 types.Redelegation) error {
	if !reflect.DeepEqual(d1, d2) {
		return fmt.Errorf("not deep equal")
	}
	if d1.Amount != nil && d1.Amount == d2.Amount {
		return fmt.Errorf("amount same address")
	}
	if d1.Reward != nil && d1.Reward == d2.Reward {
		return fmt.Errorf("reward same address")
	}
	if err := assertUndelegationsDeepCopy(d1.Undelegations, d2.Undelegations); err != nil {
		return fmt.Errorf("undelegations %v", err)
	}
	return nil
}

func assertUndelegationsDeepCopy(uds1, uds2 types.Undelegations) error {
	if !reflect.DeepEqual(uds1, uds2) {
		return fmt.Errorf("not deep equal")
	}
	for i := range uds1 {
		if err := assertUndelegationDeepCopy(uds1[i], uds2[i]); err != nil {
			return fmt.Errorf("[%v]: %v", i, err)
		}
	}
	return nil
}

func assertUndelegationDeepCopy(ud1, ud2 types.Undelegation) error {
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

func assertDecCopy(d1, d2 numeric.Dec) error {
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
