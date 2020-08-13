package staketest

import (
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"testing"
)

func TestCheckValidatorEqual(t *testing.T) {
	tests := []struct {
		v1, v2 restaking.Validator_
	}{
		{validatorPrototype, validatorPrototype},
		{makeZeroValidator(), makeZeroValidator()},
		{restaking.Validator_{}, restaking.Validator_{}},
	}
	for i, test := range tests {
		if err := CheckValidatorEqual(&test.v1, &test.v2); err != nil {
			t.Errorf("Test %v: %v", i, err)
		}
	}
}
