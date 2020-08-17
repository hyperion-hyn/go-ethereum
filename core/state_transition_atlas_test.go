package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"math/big"
	"reflect"
	"testing"
)

func TestSaveNewValidatorToPool(t *testing.T) {
	type args struct {
		wrapper       *restaking.ValidatorWrapper_
		validatorPool *restaking.Storage_ValidatorPool_
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			saveNewValidatorToPool(tt.args.wrapper, tt.args.validatorPool)
		})
	}
}

func TestUpdateValidatorFromPoolByMsg(t *testing.T) {
	type args struct {
		validator *restaking.Storage_ValidatorWrapper_
		pool      *restaking.Storage_ValidatorPool_
		msg       *restaking.EditValidator
		blockNum  *big.Int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestPayoutRedelegationReward(t *testing.T) {
	type args struct {
		s         *restaking.Storage_ValidatorWrapper_
		delegator common.Address
		handler   RestakingRewardHandler
		epoch     *big.Int
	}
	tests := []struct {
		name    string
		args    args
		want    *big.Int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := payoutRedelegationReward(tt.args.s, tt.args.delegator, tt.args.handler, tt.args.epoch)
			if (err != nil) != tt.wantErr {
				t.Errorf("payoutRedelegationReward() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("payoutRedelegationReward() got = %v, want %v", got, tt.want)
			}
		})
	}
}