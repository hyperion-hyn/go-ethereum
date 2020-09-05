package restaking

import (
	"testing"
)

func TestCheckValidatorEqual(t *testing.T) {
	tests := []struct {
		v1, v2 Validator_
	}{
		{GetDefaultValidator(), GetDefaultValidator()},
		{Validator_{}, Validator_{}},
	}
	for i, test := range tests {
		if err := CheckValidatorEqual(test.v1, test.v2); err != nil {
			t.Errorf("Test %v: %v", i, err)
		}
	}
}

func TestCheckValidatorWrapperEqual(t *testing.T) {
	tests := []struct {
		w1, w2 ValidatorWrapper_
	}{
		{GetDefaultValidatorWrapper(), GetDefaultValidatorWrapper()},
		{ValidatorWrapper_{}, ValidatorWrapper_{}},
	}
	for i, test := range tests {
		if err := CheckValidatorWrapperEqual(test.w1, test.w2); err != nil {
			t.Errorf("Test %v: %v", i, err)
		}
	}
}

// GetDefaultValidator return the default staking.Validator for testing
func GetDefaultValidator() Validator_ {
	v := GetDefaultValidatorWrapper().Validator
	return v
}

// GetDefaultValidatorWrapper return the default staking.ValidatorWrapper for testing
func GetDefaultValidatorWrapper() ValidatorWrapper_ {
	v := NewValidatorWrapperBuilder().Build()
	return *v
}
