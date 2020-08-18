package staketest

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
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
			m := restaking.NewRedelegationMap()
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

func makeZeroValidatorWrapper() restaking.ValidatorWrapper_ {
	w := restaking.ValidatorWrapper_{
		Validator:                 makeZeroValidator(),
		Redelegations:             restaking.NewRedelegationMap(),
		BlockReward:               common.Big0,
		TotalDelegation:           common.Big0,
		TotalDelegationByOperator: common.Big0,
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
		Name:            "Wayne",
		Identity:        "wen",
		Website:         "harmony.one.wen",
		Details:         "best",
		SecurityContact: "sc",
	}
	v := restaking.Validator_{
		ValidatorAddress:     common.BigToAddress(common.Big1),
		OperatorAddresses:    restaking.NewAddressSetWithAddress(common.BigToAddress(common.Big1)),
		SlotPubKeys:          restaking.NewBLSKeysWithBLSKey(testPub),
		LastEpochInCommittee: big.NewInt(20),
		MaxTotalDelegation:   common.Big1,
		Status:               uint8(restaking.Active),
		Commission:           nonZeroCommission,
		Description:          d,
		CreationHeight:       big.NewInt(12306),
	}
	return v
}

func makeZeroValidator() restaking.Validator_ {
	v := restaking.Validator_{
		OperatorAddresses:    restaking.NewEmptyAddressSet(),
		SlotPubKeys:          restaking.NewEmptyBLSKeys(),
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
		cp := CopyRedelegation(test.d)
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
		Undelegation: restaking.Undelegation_{
			Amount: common.Big0,
			Epoch:  common.Big0,
		},
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
	if err := assertBigIntCopy(w1.BlockReward, w2.BlockReward); err != nil {
		return fmt.Errorf("BlockReward %v", err)
	}
	return nil
}

func assertValidatorDeepCopy(v1, v2 restaking.Validator_) error {
	if !reflect.DeepEqual(v1, v2) {
		return fmt.Errorf("not deep equal")
	}

	if &v1.OperatorAddresses == &v2.OperatorAddresses {
		return fmt.Errorf("OperatorAddresses same pointer")
	}
	if err := assertAddressSetDeepCopy(v1.OperatorAddresses, v2.OperatorAddresses); err != nil {
		return fmt.Errorf("OperatorAddresses not deep copy")
	}

	if &v1.SlotPubKeys == &v2.SlotPubKeys {
		return fmt.Errorf("SlotPubKeys same pointer")
	}
	for i := range v1.SlotPubKeys.Keys {
		if v1.SlotPubKeys.Keys[i] == v2.SlotPubKeys.Keys[i] {
			return fmt.Errorf("SlotPubKeys[%v] same address", i)
		}
	}
	if err := assertBigIntCopy(v1.LastEpochInCommittee, v2.LastEpochInCommittee); err != nil {
		return fmt.Errorf("LastEpochInCommittee %v", err)
	}
	if err := assertBigIntCopy(v1.CreationHeight, v2.CreationHeight); err != nil {
		return fmt.Errorf("CreationHeight %v", err)
	}
	if &v1.Description == &v2.Description {
		return fmt.Errorf("same description")
	}
	if err := assertCommissionDeepCopy(v1.Commission, v2.Commission); err != nil {
		return fmt.Errorf("CommissionRates: %v", err)
	}
	return nil
}

func assertAddressSetDeepCopy(as1, as2 restaking.AddressSet_) error {
	if !reflect.DeepEqual(as1, as2) {
		return fmt.Errorf("not deep equal")
	}
	if len(as1.Keys) != 0 {
		for i := range as1.Keys {
			if as1.Keys[i] == as2.Keys[i] {
				return fmt.Errorf("AddressSet key same address")
			}
			key := as1.Keys[i]
			if as1.Contain(*key) != as2.Contain(*key) {
				return fmt.Errorf("AddressSet [%v]: not equal", key)
			}
		}
	}
	return nil
}

func assertCommissionDeepCopy(c1, c2 restaking.Commission_) error {
	if !reflect.DeepEqual(c1, c2) {
		return fmt.Errorf("not deep equal")
	}
	if &c1.CommissionRates == &c2.CommissionRates {
		return fmt.Errorf("CommissionRates same address")
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
	if len(ds1.Keys) != 0 {
		for i := range ds1.Keys {
			if ds1.Keys[i] == ds2.Keys[i] {
				return fmt.Errorf("RedelegationMap key same address")
			}
			key := ds1.Keys[i]
			r1, _ := ds1.Get(*key)
			r2, _ := ds2.Get(*key)
			if err := assertRedelegationDeepCopy(r1, r2); err != nil {
				return fmt.Errorf("[%v]: %v", key, err)
			}
		}
	}
	return nil
}

func assertRedelegationDeepCopy(d1, d2 restaking.Redelegation_) error {
	if !reflect.DeepEqual(d1, d2) {
		return fmt.Errorf("not deep equal")
	}
	if err := assertBigIntCopy(d1.Amount, d2.Amount); err != nil {
		return fmt.Errorf("amount %v", err)
	}
	if err := assertBigIntCopy(d1.Reward, d2.Reward); err != nil {
		return fmt.Errorf("reward %v", err)
	}
	if &d1.Undelegation == &d2.Undelegation {
		return fmt.Errorf("undelegations same address")
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
	if !reflect.DeepEqual(d1, d2) {
		return fmt.Errorf("not deep equal")
	}
	//if err := assertBigIntCopy(d1.I, d2.I); err != nil {
	//	return fmt.Errorf("int: %v", err)
	//}
	return nil
}

func assertBigIntCopy(i1, i2 *big.Int) error {
	if (i1 == nil) != (i2 == nil) {
		return errors.New("is nil not equal")
	}
	if i1 != nil && i1 == i2 {
		return errors.New("not copy, same address")
	}
	if i1 != nil && i1.Cmp(i2) != 0 {
		return errors.New("big int not equal")
	}
	return nil
}
