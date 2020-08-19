package core

import (
	"fmt"
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
		name        string
		ctx         saveNewValidatorCtx
		slotKeySet  map[string]bool
		identitySet map[string]bool
	}{
		{
			name: "save new validator",
			ctx: saveNewValidatorCtx{
				validatorAddr:      createValidatorAddr,
				operatorAddr:       createOperatorAddr,
				key:                &blsKeys[11].pub,
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
			identitySet: map[string]bool{
				defaultDesc.Identity: true,
			},
			slotKeySet: map[string]bool{
				blsKeys[11].pub.Hex(): true,
			},
		},
		{
			name: "no new identity",
			ctx: saveNewValidatorCtx{
				validatorAddr:      createValidatorAddr,
				operatorAddr:       createOperatorAddr,
				key:                &blsKeys[11].pub,
				maxTotalDelegation: staketest.DefaultMaxTotalDel,
				commission: restaking.Commission_{
					CommissionRates: defaultCommissionRates,
					UpdateHeight:    big.NewInt(defaultBlockNumber),
				},
				description: func(d restaking.Description_) restaking.Description_ {
					d.Identity = ""
					return d
				}(defaultDesc),
				creationHeight: big.NewInt(defaultBlockNumber),
				redelegation: restaking.Redelegation_{
					DelegatorAddress: createOperatorAddr,
					Amount:           defaultDelAmount,
				},
			},
			identitySet: map[string]bool{
				defaultDesc.Identity: false,
			},
			slotKeySet: map[string]bool{
				blsKeys[11].pub.Hex(): true,
			},
		},
		{
			name: "no new slot key",
			ctx: saveNewValidatorCtx{
				validatorAddr:      createValidatorAddr,
				operatorAddr:       createOperatorAddr,
				key:                nil,
				maxTotalDelegation: staketest.DefaultMaxTotalDel,
				commission: restaking.Commission_{
					CommissionRates: defaultCommissionRates,
					UpdateHeight:    big.NewInt(defaultBlockNumber),
				},
				description: defaultDesc,
				creationHeight: big.NewInt(defaultBlockNumber),
				redelegation: restaking.Redelegation_{
					DelegatorAddress: createOperatorAddr,
					Amount:           defaultDelAmount,
				},
			},
			identitySet: map[string]bool{
				defaultDesc.Identity: true,
			},
			slotKeySet: map[string]bool{
				blsKeys[11].pub.Hex(): false,
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
			if err := assertIdentityAndSlotKeySet(tt.ctx.validatorPool, tt.identitySet, tt.slotKeySet); err != nil {
				t.Errorf("Test - %v: %v", tt.name, err)
			}
		})
	}
}

type saveNewValidatorCtx struct {
	// input args
	validatorAddr      common.Address
	operatorAddr       common.Address
	key                *restaking.BLSPublicKey_
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

func (c *saveNewValidatorCtx) makeStateAndValidator(t *testing.T) {
	builder := staketest.NewValidatorWrapperBuilder().
		SetValidatorAddress(c.validatorAddr).
		AddOperatorAddress(c.operatorAddr).
		SetMaxTotalDelegation(c.maxTotalDelegation).
		SetCommission(c.commission).
		SetDescription(c.description).
		SetCreationHeight(c.creationHeight).
		AddRedelegation(c.redelegation)
	if c.key != nil {
		builder.AddSlotPubKey(*c.key)
	}
	w := builder.Build()
	c.newValidator = &w
	c.stateDB = makeStateDBForStake(t)
	c.validatorPool = c.stateDB.ValidatorPool()
}

func TestUpdateValidatorFromPoolByMsg(t *testing.T) {
	tests := []struct {
		name        string
		msg         restaking.EditValidator
		blockNum    *big.Int
		slotKeySet  map[string]bool
		identitySet map[string]bool
	}{
		{
			name:     "edit completely",
			msg:      defaultMsgEditValidator(),
			blockNum: big.NewInt(111),
			identitySet: map[string]bool{
				editDesc.Identity:  true,
				makeIdentityStr(0): false,
			},
			slotKeySet: map[string]bool{
				blsKeys[12].pub.Hex(): true,
				blsKeys[0].pub.Hex():  false,
			},
		},
		{
			name: "UpdateHeight not update",
			msg: func() restaking.EditValidator {
				ev := defaultMsgEditValidator()
				ev.CommissionRate = nil
				return ev
			}(),
			identitySet: map[string]bool{
				editDesc.Identity:  true,
				makeIdentityStr(0): false,
			},
			slotKeySet: map[string]bool{
				blsKeys[12].pub.Hex(): true,
				blsKeys[0].pub.Hex():  false,
			},
		},
		{
			name: "add new key, not remove old key",
			msg: func() restaking.EditValidator {
				ev := defaultMsgEditValidator()
				ev.SlotKeyToRemove = nil
				return ev
			}(),
			blockNum: big.NewInt(111),
			identitySet: map[string]bool{
				editDesc.Identity:  true,
				makeIdentityStr(0): false,
			},
			slotKeySet: map[string]bool{
				blsKeys[12].pub.Hex(): true,
				blsKeys[0].pub.Hex():  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateDB := makeStateDBForStake(t)
			validatorPool := stateDB.ValidatorPool()
			validatorSt, _ := stateDB.ValidatorByAddress(tt.msg.ValidatorAddress)

			exp := staketest.CopyValidator(*validatorSt.Validator().Load())
			_ = restaking.UpdateValidatorFromEditMsg(&exp, &tt.msg)
			if tt.blockNum != nil {
				exp.Commission.UpdateHeight = tt.blockNum
			}

			updateValidatorFromPoolByMsg(validatorSt, validatorPool, &tt.msg, tt.blockNum)
			got, _ := stateDB.ValidatorByAddress(tt.msg.ValidatorAddress)

			if err := staketest.CheckValidatorEqual(*got.Validator().Load(), exp); err != nil {
				t.Errorf("Test - %v: %v", tt.name, err)
			}
			if err := assertIdentityAndSlotKeySet(validatorPool, tt.identitySet, tt.slotKeySet); err != nil {
				t.Errorf("Test - %v: %v", tt.name, err)
			}
		})
	}
}

func assertIdentityAndSlotKeySet(validatorPool *restaking.Storage_ValidatorPool_, expIdentitySet, expSlotKeySet map[string]bool) error {
	identitySet := validatorPool.DescriptionIdentitySet()
	for i, b := range expIdentitySet {
		got := identitySet.Get(i).Value()
		if got != b {
			return fmt.Errorf("identity %v: %v, exp: %v", i, got, b)
		}
	}
	keySet := validatorPool.SlotKeySet()
	for i, b := range expSlotKeySet {
		got := keySet.Get(i).Value()
		if got != b {
			return fmt.Errorf("slot key %v: %v, exp: %v", i, got, b)
		}
	}
	return nil
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
