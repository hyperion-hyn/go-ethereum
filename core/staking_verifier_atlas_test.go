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
	"github.com/ethereum/go-ethereum/staking/effective"
	staking "github.com/ethereum/go-ethereum/staking/types"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	staketest "github.com/ethereum/go-ethereum/staking/types/test"
	"github.com/pkg/errors"
	"math/big"
	"strings"
	"testing"
)

const (
	defNumWrappersInState = 5
	defNumPubPerAddr      = 1

	validatorIndex  = 0
	validator2Index = 1
	delegatorIndex  = 6
)

var (
	blsKeys = makeKeyPairs(20)

	createOperatorAddr  = makeTestAddr("operator")
	createValidatorAddr = crypto.CreateAddress(createOperatorAddr, defaultNonce)
	validatorAddr       = makeTestAddr(fmt.Sprint("val", validatorIndex))
	operatorAddr        = makeTestAddr(fmt.Sprint("op", validatorIndex))
	validatorAddr2      = makeTestAddr(validator2Index)
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
	hundredKOnes    = new(big.Int).Mul(big.NewInt(100000), oneBig)

	negRate           = common.NewDecWithPrec(-1, 10)
	pointOneDec       = common.NewDecWithPrec(1, 1)
	pointTwoDec       = common.NewDecWithPrec(2, 1)
	pointFiveDec      = common.NewDecWithPrec(5, 1)
	pointSevenDec     = common.NewDecWithPrec(7, 1)
	pointEightFiveDec = common.NewDecWithPrec(85, 2)
	pointNineDec      = common.NewDecWithPrec(9, 1)
	oneDec            = common.OneDec()
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
		keys     []*restaking.BLSPublicKey_
	}
	tests := []struct {
		name   string
		args   args
		expErr error
	}{
		{
			name: "no duplicated fields",
			args: args{
				state:    makeStateDBForStake(t),
				identity: makeIdentityStr("new validator"),
				keys:     []*restaking.BLSPublicKey_{&blsKeys[11].pub},
			},
			expErr: nil,
		},
		{
			name: "empty bls keys",
			args: args{
				state:    makeStateDBForStake(t),
				identity: makeIdentityStr("new validator"),
				keys:     []*restaking.BLSPublicKey_{},
			},
			expErr: nil,
		},
		{
			name: "empty identity",
			args: args{
				state:    makeStateDBForStake(t),
				identity: "",
				keys:     []*restaking.BLSPublicKey_{&blsKeys[11].pub},
			},
			expErr: nil,
		},
		{
			name: "identity duplication",
			args: args{
				state:    makeStateDBForStake(t),
				identity: makeIdentityStr(0),
				keys:     []*restaking.BLSPublicKey_{&blsKeys[11].pub},
			},
			expErr: errDupIdentity,
		},
		{
			name: "bls key duplication",
			args: args{
				state:    makeStateDBForStake(t),
				identity: makeIdentityStr("new validator"),
				keys:     []*restaking.BLSPublicKey_{&blsKeys[0].pub},
			},
			expErr: errDuplicateSlotKeys,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkValidatorDuplicatedFields(tt.args.state, tt.args.identity, restaking.BLSPublicKeys_{Keys: tt.args.keys})
			if assErr := assertError(err, tt.expErr); assErr != nil {
				t.Errorf("Test - %v: %v", tt.name, assErr)
			}
		})
	}
}

func TestVerifyCreateValidatorMsg(t *testing.T) {
	type args struct {
		stateDB  vm.StateDB
		blockNum *big.Int
		msg      staking.CreateValidator
		signer   common.Address
	}
	tests := []struct {
		name    string
		args    args
		want    restaking.Validator_
		wantErr error
	}{
		{
			name: "valid request",
			args: args{
				stateDB:  makeStateDBForStake(t),
				blockNum: big.NewInt(defaultBlockNumber),
				msg:      defaultMsgCreateValidator(),
				signer:   createOperatorAddr,
			},
			want:    defaultExpCreatedValidator(),
			wantErr: nil,
		},
		{
			name: "state db nil",
			args: args{
				stateDB:  nil,
				blockNum: big.NewInt(defaultBlockNumber),
				msg:      defaultMsgCreateValidator(),
				signer:   createOperatorAddr,
			},
			want:    nil,
			wantErr: errStateDBIsMissing,
		},
		{
			name: "block number nil",
			args: args{
				stateDB:  makeStateDBForStake(t),
				blockNum: nil,
				msg:      defaultMsgCreateValidator(),
				signer:   createOperatorAddr,
			},
			want:    nil,
			wantErr: errBlockNumMissing,
		},
		{
			name: "bls collision (checkDuplicateFields)",
			args: args{
				stateDB:  makeStateDBForStake(t),
				blockNum: big.NewInt(defaultBlockNumber),
				msg: func() staking.CreateValidator {
					m := defaultMsgCreateValidator()
					m.SlotPubKeys = restaking.BLSPublicKeys_{Keys: []*restaking.BLSPublicKey_{&blsKeys[0].pub}}
					return m
				}(),
				signer: createOperatorAddr,
			},
			want:    nil,
			wantErr: errors.Wrapf(errDuplicateSlotKeys, "duplicate public key %x", blsKeys[0].pub.Hex()),
		},
		{
			name: "incorrect signature",
			args: args{
				stateDB:  makeStateDBForStake(t),
				blockNum: big.NewInt(defaultBlockNumber),
				msg: func() staking.CreateValidator {
					m := defaultMsgCreateValidator()
					m.SlotKeySigs = []restaking.BLSSignature{blsKeys[12].sig}
					return m
				}(),
				signer: createOperatorAddr,
			},
			want:    nil,
			wantErr: errors.New("bls keys and corresponding signatures could not be verified"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := VerifyCreateValidatorMsg(tt.args.stateDB, tt.args.blockNum, &tt.args.msg, tt.args.signer)
			if assErr := assertError(err, tt.wantErr); assErr != nil {
				t.Errorf("Test - %v: %v", tt.name, err)
			}
			if err != nil || tt.wantErr != nil {
				return
			}
			if err := staketest.CheckValidatorEqual(got, &tt.want); err != nil {
				t.Errorf("Test - %v: %v", tt.name, err)
			}
		})
	}
}

func defaultMsgCreateValidator() staking.CreateValidator {
	pub, sig := blsKeys[11].pub, blsKeys[11].sig
	cv := staking.CreateValidator{
		OperatorAddress:    createOperatorAddr,
		Description:        defaultDesc,
		CommissionRates:    defaultCommissionRates,
		MaxTotalDelegation: staketest.DefaultMaxTotalDel,
		SlotPubKeys:        restaking.BLSPublicKeys_{Keys: []*restaking.BLSPublicKey_{&pub}},
		SlotKeySigs:        []restaking.BLSSignature{sig},
	}
	return cv
}

func defaultExpCreatedValidator() restaking.Validator_ {
	pub := blsKeys[11].pub
	v := restaking.Validator_{
		ValidatorAddress:     validatorAddr,
		OperatorAddresses:    restaking.NewAddressSetWithAddress(createOperatorAddr),
		SlotPubKeys:          restaking.BLSPublicKeys_{Keys: []*restaking.BLSPublicKey_{&pub}},
		LastEpochInCommittee: new(big.Int),
		MaxTotalDelegation:   staketest.DefaultMaxTotalDel,
		Status:               big.NewInt(int64(effective.Active)),
		Commission: restaking.Commission_{
			CommissionRates: defaultCommissionRates,
			UpdateHeight:    big.NewInt(defaultBlockNumber),
		},
		Description:    defaultDesc,
		CreationHeight: big.NewInt(defaultBlockNumber),
	}
	return v
}

func TestVerifyEditValidatorMsg(t *testing.T) {
	type args struct {
		stateDB      vm.StateDB
		chainContext ChainContext
		epoch        *big.Int
		blockNum     *big.Int
		msg          staking.EditValidator
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
				stateDB:      makeStateDBForStake(t),
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
				stateDB:      makeStateDBForStake(t),
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
				stateDB:      makeStateDBForStake(t),
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
				stateDB:      makeStateDBForStake(t),
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
				stateDB:      makeStateDBForStake(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() staking.EditValidator {
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
				stateDB:      makeStateDBForStake(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() staking.EditValidator {
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
				stateDB:      makeStateDBForStake(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() staking.EditValidator {
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
				stateDB:      makeStateDBForStake(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() staking.EditValidator {
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
				stateDB:      makeStateDBForStake(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() staking.EditValidator {
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
				stateDB:      makeStateDBForStake(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() staking.EditValidator {
					msg := defaultMsgEditValidator()
					msg.SlotKeyToAddSig = &blsKeys[13].sig
					return msg
				}(),
				signer: operatorAddr,
			},
			wantErr: errDupIdentity,
		},
		{
			name: "rate exceed maxRate",
			args: args{
				stateDB:      makeStateDBForStake(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() staking.EditValidator {
					msg := defaultMsgEditValidator()
					msg.CommissionRate = &oneDec
					return msg
				}(),
				signer: operatorAddr,
			},
			wantErr: errCommissionRateChangeTooHigh,
		},
		{
			name: "rate exceed maxChangeRate",
			args: args{
				stateDB:      makeStateDBForStake(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() staking.EditValidator {
					msg := defaultMsgEditValidator()
					msg.CommissionRate = &pointEightFiveDec
					return msg
				}(),
				signer: operatorAddr,
			},
			wantErr: errCommissionRateChangeTooFast,
		},
		{
			name: "banned validator",
			args: args{
				//stateDB:      makeStateDBForStake(t),
				stateDB: func(t *testing.T) *state.StateDB {
					sdb := makeStateDBForStake(t)
					vw, err := sdb.ValidatorByAddress(validatorAddr)
					if err != nil {
						t.Fatal(err)
					}
					vw.Validator().Status().SetValue(big.NewInt(int64(effective.Banned)))
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
			err := VerifyEditValidatorMsg(tt.args.stateDB, tt.args.chainContext, tt.args.epoch, tt.args.blockNum, &tt.args.msg, tt.args.signer)
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

func defaultMsgEditValidator() staking.EditValidator {
	var (
		pub0Copy  restaking.BLSPublicKey_
		pub12Copy restaking.BLSPublicKey_
		sig12Copy restaking.BLSSignature
	)
	copy(pub0Copy.Key[:], blsKeys[0].pub.Key[:])
	copy(pub12Copy.Key[:], blsKeys[12].pub.Key[:])
	copy(sig12Copy[:], blsKeys[12].sig[:])

	return staking.EditValidator{
		ValidatorAddress: validatorAddr,
		OperatorAddress:  operatorAddr,
		Description:      &editDesc,
		CommissionRate:   &pointTwoDec,
		SlotKeyToRemove:  &pub0Copy,
		SlotKeyToAdd:     &pub12Copy,
		SlotKeyToAddSig:  &sig12Copy,
		EPOSStatus:       effective.Inactive,
	}
}

func TestVerifyRedelegateMsg(t *testing.T) {
	type args struct {
		stateDB vm.StateDB
		msg     staking.Redelegate
		signer  common.Address
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "new delegate",
			args: args{
				stateDB: makeStateDBForStake(t),
				msg:     defaultMsgDelegate(),
				signer:  delegatorAddr,
			},
			wantErr: nil,
		},
		{
			name: "state db nil",
			args: args{
				stateDB: nil,
				msg:     defaultMsgDelegate(),
				signer:  delegatorAddr,
			},
			wantErr: errStateDBIsMissing,
		},
		{
			name: "validator not exist",
			args: args{
				stateDB: makeStateDBForStake(t),
				msg: func() staking.Redelegate {
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
				stateDB: makeStateDBForStake(t),
				msg:     defaultMsgDelegate(),
				signer:  makeTestAddr("invalid operator"),
			},
			wantErr: errInvalidSigner,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyRedelegateMsg(tt.args.stateDB, &tt.args.msg, tt.args.signer)
			if assErr := assertError(err, tt.wantErr); assErr != nil {
				t.Errorf("Test - %v: %v", tt.name, err)
			}
		})
	}
}

func defaultMsgDelegate() staking.Redelegate {
	return staking.Redelegate{
		DelegatorAddress: delegatorAddr,
		ValidatorAddress: validatorAddr,
	}
}

//
//func TestVerifyAndUndelegateFromMsg(t *testing.T) {
//	tests := []struct {
//		sdb   vm.StateDB
//		epoch *big.Int
//		msg   staking.Undelegate
//
//		expVWrapper staking.ValidatorWrapper
//		expErr      error
//	}{
//		{
//			// 0: Unredelegate at delegation with an entry already exist at the same epoch.
//			// Will increase the amount in undelegate entry
//			sdb:   makeDefaultStateForUndelegate(t),
//			epoch: big.NewInt(defaultEpoch),
//			msg:   defaultMsgUndelegate(),
//
//			expVWrapper: defaultExpVWrapperUndelegateSameEpoch(t),
//		},
//		{
//			// 1: Unredelegate with undelegation entry exist but not in same epoch.
//			// Will create a new undelegate entry
//			sdb:   makeDefaultStateForUndelegate(t),
//			epoch: big.NewInt(defaultNextEpoch),
//			msg:   defaultMsgUndelegate(),
//
//			expVWrapper: defaultExpVWrapperUndelegateNextEpoch(t),
//		},
//		{
//			// 2: Unredelegate from a delegation record with no undelegation entry.
//			// Will create a new undelegate entry
//			sdb:   makeDefaultStateForUndelegate(t),
//			epoch: big.NewInt(defaultEpoch),
//			msg:   defaultMsgSelfUndelegate(),
//
//			expVWrapper: defaultVWrapperSelfUndelegate(t),
//		},
//		{
//			// 3: Self delegation below min self delegation, change status to Inactive
//			sdb:   makeDefaultStateForUndelegate(t),
//			epoch: big.NewInt(defaultEpoch),
//			msg: func() staking.Undelegate {
//				msg := defaultMsgSelfUndelegate()
//				msg.Amount = new(big.Int).Set(fifteenKOnes)
//				return msg
//			}(),
//
//			expVWrapper: func(t *testing.T) staking.ValidatorWrapper {
//				w := defaultVWrapperSelfUndelegate(t)
//
//				w.Delegations[0].Amount = new(big.Int).Set(fiveKOnes)
//				w.Delegations[0].Undelegations[0].Amount = new(big.Int).Set(fifteenKOnes)
//				w.Status = effective.Inactive
//
//				return w
//			}(t),
//		},
//		{
//			// 4: Extract tokens from banned validator
//			sdb: func(t *testing.T) *state.DB {
//				sdb := makeDefaultStateForUndelegate(t)
//				w, err := sdb.ValidatorWrapper(validatorAddr)
//				if err != nil {
//					t.Fatal(err)
//				}
//				w.Status = effective.Banned
//				if err := sdb.UpdateValidatorWrapper(validatorAddr, w); err != nil {
//					t.Fatal(err)
//				}
//				return sdb
//			}(t),
//			epoch: big.NewInt(defaultEpoch),
//			msg: func() staking.Undelegate {
//				msg := defaultMsgSelfUndelegate()
//				msg.Amount = new(big.Int).Set(fifteenKOnes)
//				return msg
//			}(),
//
//			expVWrapper: func(t *testing.T) staking.ValidatorWrapper {
//				w := defaultVWrapperSelfUndelegate(t)
//
//				w.Delegations[0].Amount = new(big.Int).Set(fiveKOnes)
//				w.Delegations[0].Undelegations[0].Amount = new(big.Int).Set(fifteenKOnes)
//				w.Status = effective.Banned
//
//				return w
//			}(t),
//		},
//		{
//			// 5: nil state db
//			sdb:   nil,
//			epoch: big.NewInt(defaultEpoch),
//			msg:   defaultMsgUndelegate(),
//
//			expErr: errStateDBIsMissing,
//		},
//		{
//			// 6: nil epoch
//			sdb:   makeDefaultStateForUndelegate(t),
//			epoch: nil,
//			msg:   defaultMsgUndelegate(),
//
//			expErr: errEpochMissing,
//		},
//		{
//			// 7: negative amount
//			sdb:   makeDefaultStateForUndelegate(t),
//			epoch: big.NewInt(defaultEpoch),
//			msg: func() staking.Undelegate {
//				msg := defaultMsgUndelegate()
//				msg.Amount = big.NewInt(-1)
//				return msg
//			}(),
//
//			expErr: errNegativeAmount,
//		},
//		{
//			// 8: validator flag not set
//			sdb: func() *state.DB {
//				sdb := makeStateDBForStake(t)
//				w := makeVWrapperByIndex(6)
//				if err := sdb.UpdateValidatorWrapper(makeTestAddr(6), &w); err != nil {
//					t.Fatal(err)
//				}
//				return sdb
//			}(),
//			epoch: big.NewInt(defaultEpoch),
//			msg: func() staking.Undelegate {
//				msg := defaultMsgUndelegate()
//				msg.ValidatorAddress = makeTestAddr(6)
//				return msg
//			}(),
//
//			expErr: errValidatorNotExist,
//		},
//		{
//			// 9: vWrapper not in state
//			sdb: func() *state.DB {
//				sdb := makeStateDBForStake(t)
//				sdb.SetValidatorFlag(makeTestAddr(6))
//				return sdb
//			}(),
//			epoch: big.NewInt(defaultEpoch),
//			msg: func() staking.Undelegate {
//				msg := defaultMsgUndelegate()
//				msg.ValidatorAddress = makeTestAddr(6)
//				return msg
//			}(),
//
//			expErr: errors.New("address not present in state"),
//		},
//		{
//			// 10: Insufficient balance to undelegate
//			sdb:   makeDefaultStateForUndelegate(t),
//			epoch: big.NewInt(defaultEpoch),
//			msg: func() staking.Undelegate {
//				msg := defaultMsgUndelegate()
//				msg.Amount = new(big.Int).Set(hundredKOnes)
//				return msg
//			}(),
//
//			expErr: errors.New("insufficient balance to undelegate"),
//		},
//		{
//			// 11: No delegation record
//			sdb:   makeDefaultStateForUndelegate(t),
//			epoch: big.NewInt(defaultEpoch),
//			msg: func() staking.Undelegate {
//				msg := defaultMsgUndelegate()
//				msg.DelegatorAddress = makeTestAddr("not exist")
//				return msg
//			}(),
//
//			expErr: errNoDelegationToUndelegate,
//		},
//	}
//	for i, test := range tests {
//		w, err := VerifyUnredelegateMsg(test.sdb, test.epoch, &test.msg)
//
//		if assErr := assertError(err, test.expErr); assErr != nil {
//			t.Errorf("Test %v: %v", i, assErr)
//		}
//		if err != nil || test.expErr != nil {
//			continue
//		}
//
//		if err := staketest.CheckValidatorWrapperEqual(*w, test.expVWrapper); err != nil {
//			t.Errorf("Test %v: %v", i, err)
//		}
//	}
//}
//
//func makeDefaultSnapVWrapperForUndelegate(t *testing.T) staking.ValidatorWrapper {
//	w := makeVWrapperByIndex(validatorIndex)
//
//	newDelegation := staking.NewDelegation(delegatorAddr, new(big.Int).Set(twentyKOnes))
//	if err := newDelegation.Undelegate(big.NewInt(defaultEpoch), fiveKOnes); err != nil {
//		t.Fatal(err)
//	}
//	w.Delegations = append(w.Delegations, newDelegation)
//
//	return w
//}
//
//func makeDefaultStateForUndelegate(t *testing.T) *state.DB {
//	sdb := makeStateDBForStake(t)
//	w := makeDefaultSnapVWrapperForUndelegate(t)
//
//	if err := sdb.UpdateValidatorWrapper(validatorAddr, &w); err != nil {
//		t.Fatal(err)
//	}
//	sdb.IntermediateRoot(true)
//	return sdb
//}
//
//// undelegate from delegator which has already go one entry for undelegation
//func defaultMsgUndelegate() staking.Undelegate {
//	return staking.Undelegate{
//		DelegatorAddress: delegatorAddr,
//		ValidatorAddress: validatorAddr,
//		Amount:           fiveKOnes,
//	}
//}
//
//func defaultExpVWrapperUndelegateSameEpoch(t *testing.T) staking.ValidatorWrapper {
//	w := makeDefaultSnapVWrapperForUndelegate(t)
//
//	amt := w.Delegations[1].Undelegations[0].Amount
//	w.Delegations[1].Undelegations[0].Amount = new(big.Int).
//		Add(w.Delegations[1].Undelegations[0].Amount, amt)
//	w.Delegations[1].Amount = new(big.Int).Sub(w.Delegations[1].Amount, fiveKOnes)
//
//	return w
//}
//
//func defaultExpVWrapperUndelegateNextEpoch(t *testing.T) staking.ValidatorWrapper {
//	w := makeDefaultSnapVWrapperForUndelegate(t)
//
//	w.Delegations[1].Undelegations = append(w.Delegations[1].Undelegations,
//		staking.Undelegation{Amount: fiveKOnes, Epoch: big.NewInt(defaultNextEpoch)})
//	w.Delegations[1].Amount = new(big.Int).Sub(w.Delegations[1].Amount, fiveKOnes)
//
//	return w
//}
//
//// undelegate from self undelegation (new undelegates)
//func defaultMsgSelfUndelegate() staking.Undelegate {
//	return staking.Undelegate{
//		DelegatorAddress: validatorAddr,
//		ValidatorAddress: validatorAddr,
//		Amount:           fiveKOnes,
//	}
//}
//
//func defaultVWrapperSelfUndelegate(t *testing.T) staking.ValidatorWrapper {
//	w := makeDefaultSnapVWrapperForUndelegate(t)
//
//	w.Delegations[0].Undelegations = staking.Undelegations{
//		staking.Undelegation{Amount: fiveKOnes, Epoch: big.NewInt(defaultEpoch)},
//	}
//	w.Delegations[0].Amount = new(big.Int).Sub(w.Delegations[0].Amount, fiveKOnes)
//
//	return w
//}
//
//var (
//	reward00 = twentyKOnes
//	reward01 = tenKOnes
//	reward10 = thirtyKOnes
//	reward11 = twentyFiveKOnes
//)
//
//func TestVerifyAndCollectRewardsFromDelegation(t *testing.T) {
//	tests := []struct {
//		sdb vm.StateDB
//		ds  []staking.DelegationIndex
//
//		expVWrappers    []*staking.ValidatorWrapper
//		expTotalRewards *big.Int
//		expErr          error
//	}{
//		{
//			// 0: Positive test case
//			sdb: makeStateForReward(t),
//			ds:  makeMsgCollectRewards(),
//
//			expVWrappers:    expVWrappersForReward(),
//			expTotalRewards: new(big.Int).Add(reward01, reward11),
//		},
//		{
//			// 1: No rewards to collect
//			sdb: makeStateDBForStake(t),
//			ds:  []staking.DelegationIndex{{ValidatorAddress: validatorAddr2, Index: 0}},
//
//			expErr: errNoRewardsToCollect,
//		},
//		{
//			// 2: nil state db
//			sdb: nil,
//			ds:  makeMsgCollectRewards(),
//
//			expErr: errStateDBIsMissing,
//		},
//		{
//			// 3: ValidatorWrapper not in state
//			sdb: makeStateForReward(t),
//			ds: func() []staking.DelegationIndex {
//				msg := makeMsgCollectRewards()
//				msg[1].ValidatorAddress = makeTestAddr("addr not exist")
//				return msg
//			}(),
//
//			expErr: errors.New("address not present in state"),
//		},
//		{
//			// 4: Wrong input message - index out of range
//			sdb: makeStateForReward(t),
//			ds: func() []staking.DelegationIndex {
//				dis := makeMsgCollectRewards()
//				dis[1].Index = 2
//				return dis
//			}(),
//
//			expErr: errors.New("index out of bound"),
//		},
//	}
//	for i, test := range tests {
//		ws, tReward, err := VerifyCollectRedelRewardsMsg(test.sdb, nil, )
//
//		if assErr := assertError(err, test.expErr); assErr != nil {
//			t.Fatalf("Test %v: %v", i, err)
//		}
//		if err != nil || test.expErr != nil {
//			continue
//		}
//
//		if len(ws) != len(test.expVWrappers) {
//			t.Fatalf("vwrapper size unexpected: %v / %v", len(ws), len(test.expVWrappers))
//		}
//		for wi := range ws {
//			if err := staketest.CheckValidatorWrapperEqual(*ws[wi], *test.expVWrappers[wi]); err != nil {
//				t.Errorf("%v wrapper: %v", wi, err)
//			}
//		}
//		if tReward.Cmp(test.expTotalRewards) != 0 {
//			t.Errorf("Test %v: total Rewards unexpected: %v / %v", i, tReward, test.expTotalRewards)
//		}
//	}
//}
//
//func makeMsgCollectRewards() []staking.DelegationIndex {
//	dis := []staking.DelegationIndex{
//		{
//			ValidatorAddress: validatorAddr,
//			Index:            1,
//			BlockNum:         big.NewInt(defaultBlockNumber),
//		}, {
//			ValidatorAddress: validatorAddr2,
//			Index:            1,
//			BlockNum:         big.NewInt(defaultBlockNumber),
//		},
//	}
//	return dis
//}
//
//func makeStateForReward(t *testing.T) *state.StateDB {
//	sdb := makeStateDBForStake(t)
//
//	rewards0 := []*big.Int{reward00, reward01}
//	if err := addStateRewardForAddr(sdb, validatorAddr, rewards0); err != nil {
//		t.Fatal(err)
//	}
//	rewards1 := []*big.Int{reward10, reward11}
//	if err := addStateRewardForAddr(sdb, validatorAddr2, rewards1); err != nil {
//		t.Fatal(err)
//	}
//
//	sdb.IntermediateRoot(true)
//	return sdb
//}
//
//func addStateRewardForAddr(sdb *state.DB, addr common.Address, rewards []*big.Int) error {
//	w, err := sdb.ValidatorWrapper(addr)
//	if err != nil {
//		return err
//	}
//	w.Delegations = append(w.Delegations,
//		staking.NewDelegation(delegatorAddr, new(big.Int).Set(twentyKOnes)),
//	)
//	w.Delegations[1].Undelegations = staking.Undelegations{}
//	w.Delegations[0].Reward = new(big.Int).Set(rewards[0])
//	w.Delegations[1].Reward = new(big.Int).Set(rewards[1])
//
//	return sdb.UpdateValidatorWrapper(addr, w)
//}
//
//func expVWrappersForReward() []*staking.ValidatorWrapper {
//	w1 := makeVWrapperByIndex(validatorIndex)
//	w1.Delegations = append(w1.Delegations,
//		staking.NewDelegation(delegatorAddr, new(big.Int).Set(twentyKOnes)),
//	)
//	w1.Delegations[1].Undelegations = staking.Undelegations{}
//	w1.Delegations[0].Reward = new(big.Int).Set(reward00)
//	w1.Delegations[1].Reward = new(big.Int).SetUint64(0)
//
//	w2 := makeVWrapperByIndex(validator2Index)
//	w2.Delegations = append(w2.Delegations,
//		staking.NewDelegation(delegatorAddr, new(big.Int).Set(twentyKOnes)),
//	)
//	w2.Delegations[1].Undelegations = staking.Undelegations{}
//	w2.Delegations[0].Reward = new(big.Int).Set(reward10)
//	w2.Delegations[1].Reward = new(big.Int).SetUint64(0)
//	return []*staking.ValidatorWrapper{&w1, &w2}
//}

// makeFakeChainContextForStake makes the default fakeChainContext for staking test
func makeFakeChainContextForStake(t *testing.T) *fakeChainContext {
	stateDB := makeStateDBForStake(t)
	return &fakeChainContext{stateDBs: map[uint64]*state.StateDB{
		defaultEpoch: stateDB,
	}}
}

// makeStateDBForStake make the default state db for staking test
func makeStateDBForStake(t *testing.T) *state.StateDB {
	sdb, err := newTestStateDB()
	if err != nil {
		t.Fatal(err)
	}
	ws := makeVWrappersForStake(defNumWrappersInState, defNumPubPerAddr)
	if err := updateStateValidators(sdb, ws); err != nil {
		t.Fatalf("make default state: %v", err)
	}
	sdb.AddBalance(createOperatorAddr, hundredKOnes)
	sdb.AddBalance(delegatorAddr, hundredKOnes)

	sdb.IntermediateRoot(true)

	return sdb
}

func updateStateValidators(sdb *state.StateDB, ws []*restaking.ValidatorWrapper_) error {
	for _, w := range ws {
		sdb.ValidatorPool().Validators().Put(w.Validator.ValidatorAddress, w)
		for _, k := range w.Validator.SlotPubKeys.Keys {
			sdb.ValidatorPool().SlotKeySet().Get(k.Hex()).SetValue(true)
		}
		sdb.ValidatorPool().DescriptionIdentitySet().Get(w.Validator.Description.Identity).SetValue(true)
	}
	return nil
}

func makeVWrapperByIndex(index int) restaking.ValidatorWrapper_ {
	pubGetter := newBLSPubGetter(blsKeys[index*defNumPubPerAddr:])

	return makeStateVWrapperFromGetter(index, defNumPubPerAddr, pubGetter)
}

func newTestStateDB() (*state.StateDB, error) {
	return state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()))
}

// makeVWrappersForStake makes the default staking.ValidatorWrappers for
// initialization of default state db for staking test
func makeVWrappersForStake(num, numPubsPerVal int) []*restaking.ValidatorWrapper_ {
	ws := make([]*restaking.ValidatorWrapper_, 0, num)
	pubGetter := newBLSPubGetter(blsKeys)
	for i := 0; i != num; i++ {
		w := makeStateVWrapperFromGetter(i, numPubsPerVal, pubGetter)
		ws = append(ws, &w)
	}
	return ws
}

func makeStateVWrapperFromGetter(index int, numPubs int, pubGetter *BLSPubGetter) restaking.ValidatorWrapper_ {
	validatorAddr := makeTestAddr(fmt.Sprint("val", index))
	operator := makeTestAddr(fmt.Sprint("op", index))
	pubs := make([]*restaking.BLSPublicKey_, 0, numPubs)
	for i := 0; i != numPubs; i++ {
		pub := pubGetter.getPub()
		pubs = append(pubs, &pub)
	}
	w := staketest.GetDefaultValidatorWrapperWithAddr(validatorAddr, operator, restaking.BLSPublicKeys_{Keys: pubs})
	w.Validator.Description.Identity = makeIdentityStr(index)
	w.Validator.Commission.UpdateHeight = big.NewInt(defaultSnapBlockNumber)
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

// fakeChainContext is the fake structure of ChainContext for testing
type fakeChainContext struct {
	stateDBs map[uint64]*state.StateDB
}

func (chain *fakeChainContext) Engine() consensus.Engine {
	panic("no implement")
}

func (chain *fakeChainContext) GetHeader(common.Hash, uint64) *types.Header {
	return nil
}

func (chain *fakeChainContext) ReadValidatorAtEpoch(*big.Int, common.Address) (*restaking.Storage_ValidatorWrapper_, error) {
	panic("implement me")
}

func makeIdentityStr(item interface{}) string {
	return fmt.Sprintf("harmony-one-%v", item)
}

func makeTestAddr(item interface{}) common.Address {
	s := fmt.Sprintf("harmony-one-%v", item)
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
	pub restaking.BLSPublicKey_
	sig restaking.BLSSignature
}

func makeBLSKeyPair() blsPubSigPair {
	blsPriv := bls.RandPrivateKey()
	blsPub := blsPriv.GetPublicKey()
	msgHash := crypto.Keccak256([]byte(restaking.BLSVerificationStr))
	sig := blsPriv.SignHash(msgHash)

	var pub restaking.BLSPublicKey_
	copy(pub.Key[:], blsPub.Serialize())

	var signature restaking.BLSSignature
	copy(signature[:], sig.Serialize())

	return blsPubSigPair{pub, signature}
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
