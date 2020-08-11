package staketest

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
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

func TestGetDefaultValidatorWithAddr(t *testing.T) {
	tests := []struct {
		validatorAddr common.Address
		operatorAddr  common.Address
		keys          restaking.BLSPublicKeys_
	}{
		{
			validatorAddr: common.BigToAddress(common.Big1),
			operatorAddr:  common.BigToAddress(common.Big2),
			keys:          restaking.BLSPublicKeys_{Keys: []*restaking.BLSPublicKey_{&testPub}},
		},
		{
			validatorAddr: common.Address{},
			operatorAddr:  common.Address{},
			keys:          restaking.BLSPublicKeys_{},
		},
		{},
	}
	for i, test := range tests {
		v := GetDefaultValidatorWithAddr(test.validatorAddr, test.operatorAddr, test.keys)

		exp := CopyValidator(validatorPrototype)
		exp.ValidatorAddress = test.validatorAddr
		exp.OperatorAddresses = restaking.NewAddressSetWithAddress(test.operatorAddr)
		exp.SlotPubKeys = test.keys

		if err := assertValidatorDeepCopy(v, exp); err != nil {
			t.Errorf("Test %v: %v", i, err)
		}
	}
}

func TestGetDefaultValidatorWrapperWithAddr(t *testing.T) {
	tests := []struct {
		validatorAddr common.Address
		operatorAddr  common.Address
		keys          restaking.BLSPublicKeys_
	}{
		{
			validatorAddr: common.BigToAddress(common.Big1),
			operatorAddr:  common.BigToAddress(common.Big2),
			keys:          restaking.BLSPublicKeys_{Keys: []*restaking.BLSPublicKey_{&testPub}},
		},
		{
			validatorAddr: common.Address{},
			operatorAddr:  common.Address{},
			keys:          restaking.BLSPublicKeys_{},
		},
		{},
	}
	for i, test := range tests {
		v := GetDefaultValidatorWrapperWithAddr(test.validatorAddr, test.operatorAddr, test.keys)

		exp := CopyValidatorWrapper(vWrapperPrototype)
		exp.Validator.ValidatorAddress = test.validatorAddr
		exp.Validator.OperatorAddresses = restaking.NewAddressSetWithAddress(test.operatorAddr)
		exp.Validator.SlotPubKeys = test.keys
		exp.Redelegations = func() restaking.RedelegationMap_ {
			m := restaking.NewRelegationMap()
			m.Put(test.operatorAddr, &restaking.Redelegation_{
				DelegatorAddress: test.operatorAddr,
			})
			return m
		}()

		if err := assertValidatorWrapperDeepCopy(v, exp); err != nil {
			t.Errorf("Test %v: %v", i, err)
		}
	}
}
