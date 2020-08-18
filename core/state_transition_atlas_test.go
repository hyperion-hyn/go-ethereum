package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	staketest "github.com/ethereum/go-ethereum/staking/types/restaking/test"
	"math/big"
	"reflect"
	"testing"
)

func TestSaveNewValidatorToPool(t *testing.T) {
	tests := []struct {
		name string
		ctx  saveNewValidatorCtx
	}{
		{
			name: "save new validator",
			ctx: saveNewValidatorCtx{
				validatorAddr:      createValidatorAddr,
				operatorAddr:       createOperatorAddr,
				key:                blsKeys[11].pub,
				maxTotalDelegation: staketest.DefaultMaxTotalDel,
				commission: restaking.Commission_{
					CommissionRates: defaultCommissionRates,
					UpdateHeight:    big.NewInt(defaultBlockNumber),
				},
				description:    defaultDesc,
				creationHeight: big.NewInt(defaultBlockNumber),
				redelegation: restaking.Redelegation_{
					DelegatorAddress: createOperatorAddr,
					Amount:           defaultDelAmount,
				},
			},
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ctx.makeStateAndValidator(t)
			saveNewValidatorToPool(tt.ctx.newValidator, tt.ctx.validatorPool)
			got, _ := tt.ctx.stateDB.ValidatorByAddress(tt.ctx.validatorAddr)
			exp := staketest.CopyValidatorWrapper(*tt.ctx.newValidator)
			if err := staketest.CheckValidatorWrapperEqual(*got.Load(), exp); err != nil {
				t.Errorf("Test - %v: %v", tt.name, err)
			}
			identitySet := tt.ctx.validatorPool.DescriptionIdentitySet()
			if !identitySet.Get(tt.ctx.description.Identity).Value() {
				t.Errorf("Test - %v: identity not contain", tt.name)
			}
			keySet := tt.ctx.validatorPool.SlotKeySet()
			if !keySet.Get(tt.ctx.key.Hex()).Value() {
				t.Errorf("Test - %v: identity not contain", tt.name)
			}

		})
	}
}

type saveNewValidatorCtx struct {
	// input args
	validatorAddr      common.Address
	operatorAddr       common.Address
	key                restaking.BLSPublicKey_
	maxTotalDelegation *big.Int
	commission         restaking.Commission_
	description        restaking.Description_
	creationHeight     *big.Int
	redelegation       restaking.Redelegation_

	// computed fields
	newValidator  *restaking.ValidatorWrapper_
	validatorPool *restaking.Storage_ValidatorPool_
	stateDB       *state.StateDB
}

func (s *saveNewValidatorCtx) makeStateAndValidator(t *testing.T) {
	w := staketest.NewValidatorWrapperBuilder().
		SetValidatorAddress(s.validatorAddr).
		AddOperatorAddress(s.operatorAddr).
		AddSlotPubKey(s.key).
		SetMaxTotalDelegation(s.maxTotalDelegation).
		SetCommission(s.commission).
		SetDescription(s.description).
		SetCreationHeight(s.creationHeight).
		AddRedelegation(s.redelegation).
		Build()
	s.newValidator = &w
	s.stateDB = makeStateDBForStake(t)
	s.validatorPool = s.stateDB.ValidatorPool()
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