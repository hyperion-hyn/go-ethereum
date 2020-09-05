package core

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"math/big"
	"testing"
)

var (
	oneMill = new(big.Int).Mul(big.NewInt(1000000), big.NewInt(1e18)) //
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
				maxTotalDelegation: oneMill,
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
				maxTotalDelegation: oneMill,
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
				maxTotalDelegation: oneMill,
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
				blsKeys[11].pub.Hex(): false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ctx.makeStateAndValidator(t)
			saveNewValidatorToPool(tt.ctx.newValidator, tt.ctx.validatorPool)
			v, _ := tt.ctx.stateDB.ValidatorByAddress(tt.ctx.validatorAddr)
			got, _ := v.LoadFully()
			exp, _ := tt.ctx.newValidator.Copy()
			if err := restaking.CheckValidatorWrapperEqual(*got, *exp); err != nil {
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
	builder := restaking.NewValidatorWrapperBuilder().
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
	c.newValidator = w
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
			name: "UpdateHeight not update when same rate",
			msg: func() restaking.EditValidator {
				ev := defaultMsgEditValidator()
				rate := common.NewDecWithPrec(5, 1)
				ev.CommissionRate = &rate
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

			exp, _ := validatorSt.Validator().LoadFully()
			_ = restaking.UpdateValidatorFromEditMsg(exp, &tt.msg)
			if tt.blockNum != nil {
				exp.Commission.UpdateHeight = tt.blockNum
			}

			updateValidatorFromPoolByMsg(validatorSt, validatorPool, &tt.msg, tt.blockNum)
			v, _ := stateDB.ValidatorByAddress(tt.msg.ValidatorAddress)
			got, _ := v.Validator().LoadFully()

			if err := restaking.CheckValidatorEqual(*got, *exp); err != nil {
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
	tests := []struct {
		name      string
		validator common.Address
		delegator common.Address
		want      *big.Int
		wantErr   error
	}{
		{
			name:      "collect reward",
			validator: validatorAddr,
			delegator: operatorAddr,
			want:      reward00,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateDB := makeStateForReward(t)
			handler := &RewardToBalance{StateDB: stateDB}
			validator, _ := stateDB.ValidatorByAddress(tt.validator)

			got, gotErr := payoutRedelegationReward(validator, tt.delegator, handler, nil)
			if err := assertError(gotErr, tt.wantErr); err != nil {
				t.Errorf("Test - %v: gotErr = %v, want %v", tt.name, gotErr, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			if tt.want.Cmp(got) != 0 {
				t.Errorf("Test - %v: got = %v, want %v", tt.name, got, tt.want)
			}

			redelgation, _ := validator.Redelegations().Get(tt.delegator)
			if redelgation.Reward().Value().Cmp(common.Big0) != 0 {
				t.Errorf("Test - %v: fail to collect reward", tt.name)
			}
		})
	}
}
