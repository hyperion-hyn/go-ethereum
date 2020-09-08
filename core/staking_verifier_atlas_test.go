package core

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/bls"
	"github.com/ethereum/go-ethereum/params"
	common2 "github.com/ethereum/go-ethereum/staking/types/common"
	"github.com/ethereum/go-ethereum/staking/types/microstaking"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"github.com/pkg/errors"
	"math/big"
	"strings"
	"testing"
)

const (
	defNumWrappersInState = 5
	defNumPubPerAddr      = 1

	validatorIndex  = 0
	validatorIndex2 = 7
	delegatorIndex  = 6
)

var (
	blsKeys = makeKeyPairs(20)

	createOperatorAddr  = makeTestAddr("operator")
	createValidatorAddr = crypto.CreateAddress(createOperatorAddr, defaultNonce)
	validatorAddr       = makeTestAddr(fmt.Sprint("val", validatorIndex))
	validatorAddr2      = makeTestAddr(fmt.Sprint("val", validatorIndex2))
	operatorAddr        = makeTestAddr(fmt.Sprint("op", validatorIndex))
	operatorAddr2       = makeTestAddr(fmt.Sprint("op", validatorIndex2))
	delegatorAddr       = makeTestAddr(delegatorIndex)
)

var (
	oneBig          = big.NewInt(1e18)
	fiveKOnes       = new(big.Int).Mul(big.NewInt(5000), oneBig)
	tenKOnes        = new(big.Int).Mul(big.NewInt(10000), oneBig)
	twelveKOnes     = new(big.Int).Mul(big.NewInt(12000), oneBig)
	fifteenKOnes    = new(big.Int).Mul(big.NewInt(15000), oneBig)
	twentyKOnes     = new(big.Int).Mul(big.NewInt(20000), oneBig)
	twentyFiveKOnes = new(big.Int).Mul(big.NewInt(25000), oneBig)
	thirtyKOnes     = new(big.Int).Mul(big.NewInt(30000), oneBig)
	millionOnes     = new(big.Int).Mul(big.NewInt(1000000), oneBig)

	negRate           = common.NewDecWithPrec(-1, 10)
	pointOneDec       = common.NewDecWithPrec(1, 1)
	pointTwoDec       = common.NewDecWithPrec(2, 1)
	pointFiveDec      = common.NewDecWithPrec(5, 1)
	pointSevenDec     = common.NewDecWithPrec(7, 1)
	pointEightFiveDec = common.NewDecWithPrec(85, 2)
	pointNineDec      = common.NewDecWithPrec(9, 1)
	oneDec            = common.OneDec()

	defaultDelAmount = millionOnes
)

const (
	defaultEpoch           = 5
	defaultNextEpoch       = 6
	defaultSnapBlockNumber = 90
	defaultBlockNumber     = 100
	defaultNonce           = 5
)

var (
	defaultDesc = restaking.Description_{
		Name:            "SuperHero",
		Identity:        "YouWouldNotKnow",
		Website:         "Secret Website",
		SecurityContact: "LicenseToKill",
		Details:         "blah blah blah",
	}

	defaultCommissionRates = restaking.CommissionRates_{
		Rate:          pointOneDec,
		MaxRate:       pointNineDec,
		MaxChangeRate: pointFiveDec,
	}
)

func TestCheckValidatorDuplicatedFields(t *testing.T) {
	type args struct {
		state    vm.StateDB
		identity string
		keys     restaking.BLSPublicKeys_
	}
	tests := []struct {
		name   string
		args   args
		expErr error
	}{
		{
			name: "no duplicated fields",
			args: args{
				state:    makeStateDBForRestaking(t),
				identity: makeIdentityStr("new validator"),
				keys:     restaking.NewBLSKeysWithBLSKey(blsKeys[11].pub),
			},
			expErr: nil,
		},
		{
			name: "empty bls keys",
			args: args{
				state:    makeStateDBForRestaking(t),
				identity: makeIdentityStr("new validator"),
				keys:     restaking.NewEmptyBLSKeys(),
			},
			expErr: nil,
		},
		{
			name: "empty identity",
			args: args{
				state:    makeStateDBForRestaking(t),
				identity: "",
				keys:     restaking.NewBLSKeysWithBLSKey(blsKeys[11].pub),
			},
			expErr: nil,
		},
		{
			name: "identity duplication",
			args: args{
				state:    makeStateDBForRestaking(t),
				identity: makeIdentityStr(0),
				keys:     restaking.NewBLSKeysWithBLSKey(blsKeys[11].pub),
			},
			expErr: errDupIdentity,
		},
		{
			name: "bls key duplication",
			args: args{
				state:    makeStateDBForRestaking(t),
				identity: makeIdentityStr("new validator"),
				keys:     restaking.NewBLSKeysWithBLSKey(blsKeys[0].pub),
			},
			expErr: errDuplicateSlotKeys,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkValidatorDuplicatedFields(tt.args.state, tt.args.identity, tt.args.keys)
			if assErr := assertError(err, tt.expErr); assErr != nil {
				t.Errorf("Test - %v: %v", tt.name, assErr)
			}
		})
	}
}

func TestVerifyCreateValidatorMsg(t *testing.T) {
	type args struct {
		stateDB      vm.StateDB
		chainContext ChainContext
		blockNum     *big.Int
		msg          restaking.CreateValidator
		signer       common.Address
	}
	tests := []struct {
		name    string
		args    args
		want    restaking.ValidatorWrapper_
		wantErr error
	}{
		{
			name: "valid request",
			args: args{
				stateDB:      makeStateDBForRestaking(t),
				chainContext: makeFakeChainContextForStake(t),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg:          defaultMsgCreateValidator(),
				signer:       createOperatorAddr,
			},
			want: defaultExpCreatedValidator(),
		},
		{
			name: "state db nil",
			args: args{
				stateDB:      nil,
				chainContext: makeFakeChainContextForStake(t),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg:          defaultMsgCreateValidator(),
				signer:       createOperatorAddr,
			},
			wantErr: errStateDBIsMissing,
		},
		{
			name: "block number nil",
			args: args{
				stateDB:      makeStateDBForRestaking(t),
				chainContext: makeFakeChainContextForStake(t),
				blockNum:     nil,
				msg:          defaultMsgCreateValidator(),
				signer:       createOperatorAddr,
			},
			wantErr: errBlockNumMissing,
		},
		{
			name: "bls collision (checkDuplicateFields)",
			args: args{
				stateDB:      makeStateDBForRestaking(t),
				chainContext: makeFakeChainContextForStake(t),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() restaking.CreateValidator {
					m := defaultMsgCreateValidator()
					m.SlotPubKey = blsKeys[0].pub
					return m
				}(),
				signer: createOperatorAddr,
			},
			wantErr: errors.Wrapf(errDuplicateSlotKeys, "duplicate public key %x", blsKeys[0].pub.Hex()),
		},
		{
			name: "incorrect signature",
			args: args{
				stateDB:      makeStateDBForRestaking(t),
				chainContext: makeFakeChainContextForStake(t),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() restaking.CreateValidator {
					m := defaultMsgCreateValidator()
					m.SlotKeySig = blsKeys[12].sig
					return m
				}(),
				signer: createOperatorAddr,
			},
			wantErr: errors.New("bls keys and corresponding signatures could not be verified"),
		},
		{
			name: "maxTotalDelegation less currentTotalDelegation",
			args: args{
				stateDB:      makeStateDBForRestaking(t),
				chainContext: makeFakeChainContextForStake(t),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() restaking.CreateValidator {
					m := defaultMsgCreateValidator()
					m.MaxTotalDelegation = new(big.Int).Sub(defaultStakingAmount, big.NewInt(1))
					return m
				}(),
				signer: createOperatorAddr,
			},
			wantErr: errors.New("total delegation can not be bigger than max_total_delegation"),
		},
		// TODO(ATLAS): restaking test cases
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifier, _ := NewStakingVerifier(tt.args.chainContext)
			got, err := verifier.VerifyCreateValidatorMsg(tt.args.stateDB, tt.args.blockNum, &tt.args.msg, tt.args.signer)
			if assErr := assertError(err, tt.wantErr); assErr != nil {
				t.Errorf("Test - %v: %v", tt.name, err)
			}
			if err != nil || tt.wantErr != nil {
				return
			}
			if err := restaking.CheckValidatorWrapperEqual(*got.NewValidator, tt.want); err != nil {
				t.Errorf("Test - %v: %v", tt.name, err)
			}
		})
	}
}

func defaultMsgCreateValidator() restaking.CreateValidator {
	pub, sig := blsKeys[11].pub, blsKeys[11].sig
	cv := restaking.CreateValidator{
		OperatorAddress:    createOperatorAddr,
		Description:        defaultDesc,
		CommissionRates:    defaultCommissionRates,
		MaxTotalDelegation: millionOnes,
		SlotPubKey:         pub,
		SlotKeySig:         sig,
	}
	return cv
}

func defaultExpCreatedValidator() restaking.ValidatorWrapper_ {
	pub := blsKeys[11].pub
	v := restaking.NewValidatorWrapperBuilder().
		SetValidatorAddress(createValidatorAddr).
		AddOperatorAddress(createOperatorAddr).
		AddSlotPubKey(pub).
		SetMaxTotalDelegation(millionOnes).
		SetCommission(restaking.Commission_{
			CommissionRates: defaultCommissionRates,
			UpdateHeight:    big.NewInt(defaultBlockNumber),
		}).
		SetDescription(defaultDesc).
		SetCreationHeight(big.NewInt(defaultBlockNumber)).
		AddRedelegation(restaking.Redelegation_{
			DelegatorAddress: createOperatorAddr,
			Amount:           defaultDelAmount,
		}).Build()
	return *v
}

func TestVerifyEditValidatorMsg(t *testing.T) {
	type args struct {
		stateDB      vm.StateDB
		chainContext ChainContext
		epoch        *big.Int
		blockNum     *big.Int
		msg          restaking.EditValidator
		signer       common.Address
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "valid request",
			args: args{
				stateDB:      makeStateDBForRestaking(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg:          defaultMsgEditValidator(),
				signer:       operatorAddr,
			},
			wantErr: nil,
		},
		{
			name: "state db nil",
			args: args{
				stateDB:      nil,
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg:          defaultMsgEditValidator(),
				signer:       operatorAddr,
			},
			wantErr: errStateDBIsMissing,
		},
		{
			name: "chain context nil",
			args: args{
				stateDB:      makeStateDBForRestaking(t),
				chainContext: nil,
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg:          defaultMsgEditValidator(),
				signer:       operatorAddr,
			},
			wantErr: errChainContextMissing,
		},
		{
			name: "epoch nil",
			args: args{
				stateDB:      makeStateDBForRestaking(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        nil,
				blockNum:     big.NewInt(defaultBlockNumber),
				msg:          defaultMsgEditValidator(),
				signer:       operatorAddr,
			},
			wantErr: errEpochMissing,
		},
		{
			name: "block number nil",
			args: args{
				stateDB:      makeStateDBForRestaking(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     nil,
				msg:          defaultMsgEditValidator(),
				signer:       operatorAddr,
			},
			wantErr: errBlockNumMissing,
		},
		{
			name: "bls key collision",
			args: args{
				stateDB:      makeStateDBForRestaking(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() restaking.EditValidator {
					msg := defaultMsgEditValidator()
					msg.SlotKeyToAdd = &blsKeys[3].pub
					msg.SlotKeyToAddSig = &blsKeys[3].sig
					return msg
				}(),
				signer: operatorAddr,
			},
			wantErr: errDuplicateSlotKeys,
		},
		{
			name: "identity collision",
			args: args{
				stateDB:      makeStateDBForRestaking(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() restaking.EditValidator {
					msg := defaultMsgEditValidator()
					msg.Description.Identity = makeIdentityStr(0)
					return msg
				}(),
				signer: createOperatorAddr,
			},
			wantErr: errDupIdentity,
		},
		{
			name: "validator not exist",
			args: args{
				stateDB:      makeStateDBForRestaking(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() restaking.EditValidator {
					msg := defaultMsgEditValidator()
					msg.ValidatorAddress = makeTestAddr("addr not in chain")
					return msg
				}(),
				signer: operatorAddr,
			},
			wantErr: errValidatorNotExist,
		},
		{
			name: "invalid operator",
			args: args{
				stateDB:      makeStateDBForRestaking(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() restaking.EditValidator {
					msg := defaultMsgEditValidator()
					msg.OperatorAddress = makeTestAddr("invalid operator")
					return msg
				}(),
				signer: makeTestAddr("invalid operator"),
			},
			wantErr: errInvalidValidatorOperator,
		},
		{
			name: "invalid signer",
			args: args{
				stateDB:      makeStateDBForRestaking(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() restaking.EditValidator {
					msg := defaultMsgEditValidator()
					return msg
				}(),
				signer: makeTestAddr("invalid operator"),
			},
			wantErr: errInvalidSigner,
		},
		{
			name: "signature cannot be verified",
			args: args{
				stateDB:      makeStateDBForRestaking(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() restaking.EditValidator {
					msg := defaultMsgEditValidator()
					msg.SlotKeyToAddSig = &blsKeys[13].sig
					return msg
				}(),
				signer: operatorAddr,
			},
			wantErr: errors.New("bls keys and corresponding signatures could not be verified"),
		},
		{
			name: "rate exceed maxRate",
			args: args{
				stateDB:      makeStateDBForRestaking(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() restaking.EditValidator {
					msg := defaultMsgEditValidator()
					msg.CommissionRate = &oneDec
					return msg
				}(),
				signer: operatorAddr,
			},
			wantErr: errors.New("rate:1.000000000000000000 max rate:0.900000000000000000"),
		},
		{
			name: "rate exceed maxChangeRate",
			args: args{
				stateDB:      makeStateDBForRestaking(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() restaking.EditValidator {
					msg := defaultMsgEditValidator()
					msg.CommissionRate = &pointEightFiveDec
					return msg
				}(),
				signer: operatorAddr,
			},
			wantErr: errCommissionRateChangeTooFast,
		},
		{
			name: "max total delegation too small",
			args: args{
				stateDB:      makeStateDBForRestaking(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() restaking.EditValidator {
					msg := defaultMsgEditValidator()
					msg.MaxTotalDelegation = oneBig
					return msg
				}(),
				signer: operatorAddr,
			},
			wantErr: errors.New("total delegation can not be bigger than max_total_delegation"),
		},
		{
			name: "banned validator",
			args: args{
				stateDB: func(t *testing.T) *state.StateDB {
					sdb := makeStateDBForRestaking(t)
					vw, err := sdb.ValidatorByAddress(validatorAddr)
					if err != nil {
						t.Fatal(err)
					}
					vw.Validator().Status().SetValue(uint8(restaking.Banned))
					return sdb
				}(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg:          defaultMsgEditValidator(),
				signer:       operatorAddr,
			},
			wantErr: errors.New("cannot change validator banned status"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifier, _ := NewStakingVerifier(tt.args.chainContext)
			_, err := verifier.VerifyEditValidatorMsg(tt.args.stateDB, tt.args.chainContext, tt.args.epoch, tt.args.blockNum, &tt.args.msg, tt.args.signer)
			if assErr := assertError(err, tt.wantErr); assErr != nil {
				t.Errorf("Test - %v: %v", tt.name, err)
			}
		})
	}
}

var (
	editDesc = restaking.Description_{
		Name:            "batman",
		Identity:        "batman",
		Website:         "",
		SecurityContact: "",
		Details:         "",
	}
)

func defaultMsgEditValidator() restaking.EditValidator {
	var (
		pub0Copy  restaking.BLSPublicKey_
		pub12Copy restaking.BLSPublicKey_
		sig12Copy common2.BLSSignature
	)
	copy(pub0Copy.Key[:], blsKeys[0].pub.Key[:])
	copy(pub12Copy.Key[:], blsKeys[12].pub.Key[:])
	copy(sig12Copy[:], blsKeys[12].sig[:])

	return restaking.EditValidator{
		ValidatorAddress: validatorAddr,
		OperatorAddress:  operatorAddr,
		Description:      editDesc,
		CommissionRate:   &pointTwoDec,
		SlotKeyToRemove:  &pub0Copy,
		SlotKeyToAdd:     &pub12Copy,
		SlotKeyToAddSig:  &sig12Copy,
		EPOSStatus:       restaking.Inactive,
	}
}

func TestVerifyRedelegateMsg(t *testing.T) {
	type args struct {
		stateDB      vm.StateDB
		chainContext ChainContext
		msg          restaking.Redelegate
		signer       common.Address
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "delegate successfully",
			args: args{
				stateDB:      makeStateDBForRestaking(t),
				chainContext: makeFakeChainContextForStake(t),
				msg:          defaultMsgDelegate(),
				signer:       delegatorAddr,
			},
			wantErr: nil,
		},
		{
			name: "state db nil",
			args: args{
				stateDB:      nil,
				chainContext: makeFakeChainContextForStake(t),
				msg:          defaultMsgDelegate(),
				signer:       delegatorAddr,
			},
			wantErr: errStateDBIsMissing,
		},
		{
			name: "validator not exist",
			args: args{
				stateDB:      makeStateDBForRestaking(t),
				chainContext: makeFakeChainContextForStake(t),
				msg: func() restaking.Redelegate {
					msg := defaultMsgDelegate()
					msg.ValidatorAddress = makeTestAddr("addr not in chain")
					return msg
				}(),
				signer: delegatorAddr,
			},
			wantErr: errValidatorNotExist,
		},
		{
			name: "invalid signer",
			args: args{
				stateDB:      makeStateDBForRestaking(t),
				chainContext: makeFakeChainContextForStake(t),
				msg:          defaultMsgDelegate(),
				signer:       makeTestAddr("invalid operator"),
			},
			wantErr: errInvalidSigner,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifier, _ := NewStakingVerifier(tt.args.chainContext)
			_, err := verifier.VerifyRedelegateMsg(tt.args.stateDB, &tt.args.msg, tt.args.signer)
			if assErr := assertError(err, tt.wantErr); assErr != nil {
				t.Errorf("Test - %v: %v", tt.name, err)
			}
		})
	}
}

func defaultMsgDelegate() restaking.Redelegate {
	return restaking.Redelegate{
		DelegatorAddress: delegatorAddr,
		ValidatorAddress: validatorAddr,
	}
}

func TestVerifyUnredelegateMsg(t *testing.T) {
	type args struct {
		stateDB      vm.StateDB
		chainContext ChainContext
		epoch        *big.Int
		msg          restaking.Unredelegate
		signer       common.Address
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "unredelegate successfully",
			args: args{
				stateDB:      makeDefaultStateForUndelegate(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				msg:          defaultMsgUndelegate(),
				signer:       operatorAddr2,
			},
			wantErr: nil,
		},
		{
			name: "state db nil",
			args: args{
				stateDB:      nil,
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				msg:          defaultMsgUndelegate(),
				signer:       operatorAddr2,
			},
			wantErr: errStateDBIsMissing,
		},
		{
			name: "epoch nil",
			args: args{
				stateDB:      makeDefaultStateForUndelegate(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        nil,
				msg:          defaultMsgUndelegate(),
				signer:       operatorAddr2,
			},
			wantErr: errEpochMissing,
		},
		{
			name: "invalid signer",
			args: args{
				stateDB:      makeDefaultStateForUndelegate(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				msg:          defaultMsgUndelegate(),
				signer:       makeTestAddr("invalid operator"),
			},
			wantErr: errInvalidSigner,
		},
		{
			name: "validator not exist",
			args: args{
				stateDB:      makeDefaultStateForUndelegate(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				msg: func() restaking.Unredelegate {
					msg := defaultMsgUndelegate()
					msg.ValidatorAddress = makeTestAddr("addr not in chain")
					return msg
				}(),
				signer: operatorAddr2,
			},
			wantErr: errValidatorNotExist,
		},
		{
			name: "redelegation not exist",
			args: args{
				stateDB:      makeDefaultStateForUndelegate(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				msg: func() restaking.Unredelegate {
					msg := defaultMsgUndelegate()
					msg.DelegatorAddress = makeTestAddr("addr not in chain")
					return msg
				}(),
				signer: makeTestAddr("addr not in chain"),
			},
			wantErr: errRedelegationNotExist,
		},
		{
			name: "insufficient balance to undelegate",
			args: args{
				stateDB:      makeDefaultStateForUndelegate(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				msg: func() restaking.Unredelegate {
					msg := defaultMsgUndelegate()
					msg.DelegatorAddress = delegatorAddr
					return msg
				}(),
				signer: delegatorAddr,
			},
			wantErr: errInsufficientBalanceToUndelegate,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifier, _ := NewStakingVerifier(tt.args.chainContext)
			_, err := verifier.VerifyUnredelegateMsg(tt.args.stateDB, tt.args.epoch, &tt.args.msg, tt.args.signer)
			if assErr := assertError(err, tt.wantErr); assErr != nil {
				t.Errorf("Test - %v: %v", tt.name, err)
			}
		})
	}
}

func makeDefaultSnapVWrapperForUndelegate() *restaking.ValidatorWrapper_ {
	w := makeVWrapperByIndex(validatorIndex2)
	newRedelegation := restaking.Redelegation_{
		DelegatorAddress: delegatorAddr,
		Amount:           common.Big0,
		Undelegation: restaking.Undelegation_{
			Amount: new(big.Int).Set(twentyKOnes),
			Epoch:  defaultStakingAmount,
		},
	}
	w.Redelegations.Put(delegatorAddr, newRedelegation)
	return w
}

func makeDefaultStateForUndelegate(t *testing.T) *state.StateDB {
	sdb := makeStateDBForRestaking(t)
	w := makeDefaultSnapVWrapperForUndelegate()
	if err := updateStateValidators(sdb, []*restaking.ValidatorWrapper_{w}); err != nil {
		t.Fatal(err)
	}
	sdb.IntermediateRoot(false)
	return sdb
}

// undelegate from delegator which has already go one entry for undelegation
func defaultMsgUndelegate() restaking.Unredelegate {
	return restaking.Unredelegate{
		DelegatorAddress: operatorAddr2,
		ValidatorAddress: validatorAddr2,
	}
}

var (
	reward00 = twentyKOnes
	reward01 = tenKOnes
	reward10 = thirtyKOnes
	reward11 = twentyFiveKOnes
)

func TestVerifyCollectRedelRewardsMsg(t *testing.T) {
	type args struct {
		stateDB      vm.StateDB
		chainContext ChainContext
		msg          restaking.CollectReward
		signer       common.Address
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "collect successfully",
			args: args{
				stateDB:      makeStateForReward(t),
				chainContext: makeFakeChainContextForStake(t),
				msg:          defaultMsgCollectReward(),
				signer:       operatorAddr,
			},
			wantErr: nil,
		},
		{
			name: "state db nil",
			args: args{
				stateDB:      nil,
				chainContext: makeFakeChainContextForStake(t),
				msg:          defaultMsgCollectReward(),
				signer:       operatorAddr,
			},
			wantErr: errStateDBIsMissing,
		},
		{
			name: "invalid signer",
			args: args{
				stateDB:      makeStateForReward(t),
				chainContext: makeFakeChainContextForStake(t),
				msg:          defaultMsgCollectReward(),
				signer:       makeTestAddr("invalid operator"),
			},
			wantErr: errInvalidSigner,
		},
		{
			name: "validator not exist",
			args: args{
				stateDB:      makeStateForReward(t),
				chainContext: makeFakeChainContextForStake(t),
				msg: func() restaking.CollectReward {
					msg := defaultMsgCollectReward()
					msg.ValidatorAddress = makeTestAddr("addr not in chain")
					return msg
				}(),
				signer: operatorAddr,
			},
			wantErr: errValidatorNotExist,
		},
		{
			name: "redelegation not exist",
			args: args{
				stateDB:      makeStateForReward(t),
				chainContext: makeFakeChainContextForStake(t),
				msg: func() restaking.CollectReward {
					msg := defaultMsgCollectReward()
					msg.DelegatorAddress = makeTestAddr("addr not in chain")
					return msg
				}(),
				signer: makeTestAddr("addr not in chain"),
			},
			wantErr: errRedelegationNotExist,
		},
		{
			name: "no reward",
			args: args{
				stateDB:      makeStateDBForRestaking(t),
				chainContext: makeFakeChainContextForStake(t),
				msg:          defaultMsgCollectReward(),
				signer:       operatorAddr,
			},
			wantErr: errNoRewardsToCollect,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifier, _ := NewStakingVerifier(tt.args.chainContext)
			_, err := verifier.VerifyCollectRedelRewardMsg(tt.args.stateDB, &tt.args.msg, tt.args.signer)
			if assErr := assertError(err, tt.wantErr); assErr != nil {
				t.Errorf("Test - %v: %v", tt.name, err)
			}
		})
	}
}

// undelegate from delegator which has already go one entry for undelegation
func defaultMsgCollectReward() restaking.CollectReward {
	return restaking.CollectReward{
		DelegatorAddress: operatorAddr,
		ValidatorAddress: validatorAddr,
	}
}

func makeStateForReward(t *testing.T) *state.StateDB {
	sdb := makeStateDBForRestaking(t)
	if err := addStateRewardForAddr(sdb, validatorAddr, reward00); err != nil {
		t.Fatal(err)
	}
	sdb.IntermediateRoot(false)
	return sdb
}

func addStateRewardForAddr(sdb *state.StateDB, validator common.Address, reward *big.Int) error {
	w, err := sdb.ValidatorByAddress(validator)
	if err != nil {
		return err
	}
	redelegation, _ := w.Redelegations().Get(w.Redelegations().AllKeys()[0])
	redelegation.AddReward(reward)
	return nil
}

// makeFakeChainContextForStake makes the default fakeChainContext for staking test
func makeFakeChainContextForStake(t *testing.T) *fakeChainContext {
	stateDB := makeStateDBForRestaking(t)
	return &fakeChainContext{stateDBs: map[uint64]*state.StateDB{
		defaultEpoch: stateDB,
	}}
}

// makeStateDBForRestaking make the default state db for restaking test
func makeStateDBForRestaking(t *testing.T) *state.StateDB {
	sdb, err := newTestStateDB()
	if err != nil {
		t.Fatal(err)
	}
	ws := makeVWrappersForRestaking(defNumWrappersInState, defNumPubPerAddr)
	if err := updateStateValidators(sdb, ws); err != nil {
		t.Fatalf("make default state: %v", err)
	}
	sdb.SetNonce(createOperatorAddr, defaultNonce)
	sdb.AddBalance(createOperatorAddr, millionOnes)
	sdb.AddBalance(delegatorAddr, millionOnes)
	sdb.Commit(true)
	return sdb
}

func updateStateValidators(sdb *state.StateDB, ws []*restaking.ValidatorWrapper_) error {
	for _, w := range ws {
		sdb.ValidatorPool().Validators().Put(w.Validator.ValidatorAddress, w)
		sdb.IncrementValidatorNonce()
		for _, k := range w.Validator.SlotPubKeys.Keys {
			sdb.ValidatorPool().SlotKeySet().Get(k.Hex()).SetValue(true)
		}
		sdb.ValidatorPool().DescriptionIdentitySet().Get(w.Validator.Description.Identity).SetValue(true)
	}
	return nil
}

func makeVWrapperByIndex(index int) *restaking.ValidatorWrapper_ {
	pubGetter := newBLSPubGetter(blsKeys[index*defNumPubPerAddr:])

	return makeStateVWrapperFromGetter(index, defNumPubPerAddr, pubGetter)
}

func newTestStateDB() (*state.StateDB, error) {
	return state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
}

// makeVWrappersForRestaking makes the default restaking.ValidatorWrappers for
// initialization of default state db for restaking test
func makeVWrappersForRestaking(num, numPubsPerVal int) []*restaking.ValidatorWrapper_ {
	ws := make([]*restaking.ValidatorWrapper_, 0, num)
	pubGetter := newBLSPubGetter(blsKeys)
	for i := 0; i != num; i++ {
		w := makeStateVWrapperFromGetter(i, numPubsPerVal, pubGetter)
		ws = append(ws, w)
	}
	return ws
}

func makeStateVWrapperFromGetter(index int, numPubs int, pubGetter *BLSPubGetter) *restaking.ValidatorWrapper_ {
	validatorAddr := makeTestAddr(fmt.Sprint("val", index))
	operator := makeTestAddr(fmt.Sprint("op", index))
	pubs := restaking.NewEmptyBLSKeys()
	for i := 0; i != numPubs; i++ {
		pub := pubGetter.getPub()
		pubs.Keys = append(pubs.Keys, &pub)
	}
	w := restaking.NewValidatorWrapperBuilder().
		SetValidatorAddress(validatorAddr).
		AddOperatorAddress(operator).
		AddSlotPubKeys(pubs).
		SetDescription(defaultDesc).
		SetCommission(restaking.Commission_{
			CommissionRates: defaultCommissionRates,
			UpdateHeight:    big.NewInt(defaultSnapBlockNumber),
		}).
		SetMaxTotalDelegation(defaultDelAmount).
		AddRedelegation(restaking.NewRedelegation(operator, defaultDelAmount)).
		Build()
	w.Validator.Description.Identity = makeIdentityStr(index)
	return w
}

type BLSPubGetter struct {
	keys  []blsPubSigPair
	index int
}

func newBLSPubGetter(keys []blsPubSigPair) *BLSPubGetter {
	return &BLSPubGetter{
		keys:  keys,
		index: 0,
	}
}

func (g *BLSPubGetter) getPub() restaking.BLSPublicKey_ {
	key := g.keys[g.index]
	g.index++
	return key.pub
}

func (g *BLSPubGetter) getPub2() microstaking.BLSPublicKey_ {
	key := g.keys[g.index]
	g.index++
	return key.pub2
}

// fakeChainContext is the fake structure of ChainContext for testing
type fakeChainContext struct {
	stateDBs map[uint64]*state.StateDB
}

func (chain *fakeChainContext) Engine() consensus.Engine {
	panic("no implement")
}

func (chain *fakeChainContext) GetHeader(common.Hash, uint64) *types.Header {
	panic("no implement")
}

func (chain *fakeChainContext) ReadValidatorAtEpoch(epoch *big.Int, validator common.Address) (*restaking.Storage_ValidatorWrapper_, error) {
	stateDB := chain.stateDBs[epoch.Uint64()]
	return stateDB.ValidatorByAddress(validatorAddr)
}

func (chain *fakeChainContext) ReadValidatorAtEpochOrCurrentBlock(epoch *big.Int, validator common.Address) (*restaking.Storage_ValidatorWrapper_, error) {
	stateDB := chain.stateDBs[epoch.Uint64()]
	return stateDB.ValidatorByAddress(validatorAddr)
}

func (chain *fakeChainContext) Config() *params.ChainConfig {
	// TODO(ATLAS): restaking enable
	return &params.ChainConfig{Atlas: &params.AtlasConfig{RestakingEnable: false}}
}

func makeIdentityStr(item interface{}) string {
	return fmt.Sprintf("hyperion-hyn-%v", item)
}

func makeTestAddr(item interface{}) common.Address {
	s := fmt.Sprintf("hyperion-hyn-%v", item)
	return common.BytesToAddress([]byte(s))
}

func makeKeyPairs(size int) []blsPubSigPair {
	pairs := make([]blsPubSigPair, 0, size)
	for i := 0; i != size; i++ {
		pairs = append(pairs, makeBLSKeyPair())
	}
	return pairs
}

type blsPubSigPair struct {
	pub  restaking.BLSPublicKey_
	pub2 microstaking.BLSPublicKey_
	sig  common2.BLSSignature
}

func makeBLSKeyPair() blsPubSigPair {
	blsPriv := bls.RandPrivateKey()
	blsPub := blsPriv.GetPublicKey()
	msgHash := crypto.Keccak256([]byte(common2.BLSVerificationStr))
	sig := blsPriv.SignHash(msgHash)

	var pub restaking.BLSPublicKey_
	copy(pub.Key[:], blsPub.Serialize())

	var pub2 microstaking.BLSPublicKey_
	copy(pub2.Key[:], blsPub.Serialize())

	var signature common2.BLSSignature
	copy(signature[:], sig.Serialize())

	return blsPubSigPair{pub: pub, pub2: pub2, sig: signature}
}

func assertError(got, expect error) error {
	if (got == nil) != (expect == nil) {
		return fmt.Errorf("unexpected error [%v] / [%v]", got, expect)
	}
	if (got == nil) || (expect == nil) {
		return nil
	}
	if !strings.Contains(got.Error(), expect.Error()) {
		return fmt.Errorf("unexpected error [%v] / [%v]", got, expect)
	}
	return nil
}
