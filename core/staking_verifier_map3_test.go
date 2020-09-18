package core

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	common2 "github.com/ethereum/go-ethereum/staking/types/common"
	"github.com/ethereum/go-ethereum/staking/types/microstaking"
	"github.com/pkg/errors"
	"math/big"
	"testing"
)

const (
	defNumNodeWrappersInState = 5
	defNumPubPerNode          = 1

	map3NodeIndex  = 0
	map3NodeIndex2 = 1
	map3NodeIndex3 = 2
)

var (
	twoHundredKOnes = new(big.Int).Mul(big.NewInt(200000), oneBig)

	createMap3NodeAddr = crypto.CreateAddress(createOperatorAddr, defaultNonce)
	map3NodeAddr       = makeTestAddr(fmt.Sprint("map3", map3NodeIndex))
	map3NodeAddr2      = makeTestAddr(fmt.Sprint("map3", map3NodeIndex2))
	map3NodeAddr3      = makeTestAddr(fmt.Sprint("map3", map3NodeIndex3))
	map3OperatorAddr   = makeTestAddr(fmt.Sprint("op", map3NodeIndex))
	map3OperatorAddr2  = makeTestAddr(fmt.Sprint("op", map3NodeIndex2))
	map3OperatorAddr3  = makeTestAddr(fmt.Sprint("op", map3NodeIndex3))

	defaultDesc2 = microstaking.Description_{
		Name:            "SuperHero",
		Identity:        "YouWouldNotKnow",
		Website:         "Secret Website",
		SecurityContact: "LicenseToKill",
		Details:         "blah blah blah",
	}

	defaultCommissionRate = pointOneDec
	defaultCommission     = microstaking.Commission_{
		Rate:              defaultCommissionRate,
		RateForNextPeriod: defaultCommissionRate,
		UpdateHeight:      big.NewInt(10),
	}
)

func TestCheckMap3DuplicatedFields(t *testing.T) {
	type args struct {
		state    vm.StateDB
		identity string
		keys     microstaking.BLSPublicKeys_
	}
	tests := []struct {
		name   string
		args   args
		expErr error
	}{
		{
			name: "no duplicated fields",
			args: args{
				state:    makeStateDBForMicrostaking(t),
				identity: makeIdentityStr("new map3 node"),
				keys:     microstaking.NewBLSKeysWithBLSKey(blsKeys[11].pub2),
			},
			expErr: nil,
		},
		{
			name: "empty bls keys",
			args: args{
				state:    makeStateDBForMicrostaking(t),
				identity: makeIdentityStr("new map3 node"),
				keys:     microstaking.NewEmptyBLSKeys(),
			},
			expErr: nil,
		},
		{
			name: "empty identity",
			args: args{
				state:    makeStateDBForMicrostaking(t),
				identity: "",
				keys:     microstaking.NewBLSKeysWithBLSKey(blsKeys[11].pub2),
			},
			expErr: nil,
		},
		{
			name: "identity duplication",
			args: args{
				state:    makeStateDBForMicrostaking(t),
				identity: makeIdentityStr(0),
				keys:     microstaking.NewBLSKeysWithBLSKey(blsKeys[11].pub2),
			},
			expErr: errDupMap3NodeIdentity,
		},
		{
			name: "bls key duplication",
			args: args{
				state:    makeStateDBForMicrostaking(t),
				identity: makeIdentityStr("new map3 node"),
				keys:     microstaking.NewBLSKeysWithBLSKey(blsKeys[0].pub2),
			},
			expErr: errDupMap3NodePubKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkMap3DuplicatedFields(tt.args.state, tt.args.identity, tt.args.keys)
			if assErr := assertError(err, tt.expErr); assErr != nil {
				t.Errorf("Test - %v: %v", tt.name, assErr)
			}
		})
	}
}

func TestVerifyCreateMap3NodeMsg(t *testing.T) {
	type args struct {
		stateDB      vm.StateDB
		chainContext ChainContext
		epoch        *big.Int
		blockNum     *big.Int
		msg          microstaking.CreateMap3Node
		signer       common.Address
	}
	tests := []struct {
		name    string
		args    args
		want    microstaking.Map3NodeWrapper_
		wantErr error
	}{
		{
			name: "valid request",
			args: args{
				stateDB:      makeStateDBForMicrostaking(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg:          defaultMsgCreateMap3Node(),
				signer:       createOperatorAddr,
			},
			want: defaultExpCreatedMap3Node(),
		},
		{
			name: "state db nil",
			args: args{
				stateDB:      nil,
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg:          defaultMsgCreateMap3Node(),
				signer:       createOperatorAddr,
			},
			wantErr: errStateDBIsMissing,
		},
		{
			name: "chain context nil",
			args: args{
				stateDB:      makeStateDBForMicrostaking(t),
				chainContext: nil,
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg:          defaultMsgCreateMap3Node(),
				signer:       createOperatorAddr,
			},
			wantErr: errChainContextMissing,
		},
		{
			name: "epoch nil",
			args: args{
				stateDB:      makeStateDBForMicrostaking(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        nil,
				blockNum:     big.NewInt(defaultBlockNumber),
				msg:          defaultMsgCreateMap3Node(),
				signer:       createOperatorAddr,
			},
			wantErr: errEpochMissing,
		},
		{
			name: "block number nil",
			args: args{
				stateDB:      makeStateDBForMicrostaking(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     nil,
				msg:          defaultMsgCreateMap3Node(),
				signer:       createOperatorAddr,
			},
			wantErr: errBlockNumMissing,
		},
		{
			name: "negative amount",
			args: args{
				stateDB:      makeStateDBForMicrostaking(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() microstaking.CreateMap3Node {
					m := defaultMsgCreateMap3Node()
					m.Amount = big.NewInt(-1)
					return m
				}(),
				signer: createOperatorAddr,
			},
			wantErr: errNegativeAmount,
		},
		{
			name: "invalid signer",
			args: args{
				stateDB:      makeStateDBForMicrostaking(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg:          defaultMsgCreateMap3Node(),
				signer:       makeTestAddr("invalid operator"),
			},
			wantErr: errInvalidSigner,
		},
		{
			name: "bls collision (checkDuplicateFields)",
			args: args{
				stateDB:      makeStateDBForMicrostaking(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() microstaking.CreateMap3Node {
					m := defaultMsgCreateMap3Node()
					m.NodePubKey = blsKeys[0].pub2
					return m
				}(),
				signer: createOperatorAddr,
			},
			wantErr: errors.Wrapf(errDupMap3NodePubKey, "duplicate public key %x", blsKeys[0].pub2.Hex()),
		},
		{
			name: "insufficient balance",
			args: args{
				stateDB:      makeStateDBForMicrostaking(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() microstaking.CreateMap3Node {
					m := defaultMsgCreateMap3Node()
					m.Amount = oneMill
					return m
				}(),
				signer: createOperatorAddr,
			},
			wantErr: errInsufficientBalanceForStake,
		},
		{
			name: "self delegation too small",
			args: args{
				stateDB:      makeStateDBForMicrostaking(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() microstaking.CreateMap3Node {
					m := defaultMsgCreateMap3Node()
					m.Amount = oneBig
					return m
				}(),
				signer: createOperatorAddr,
			},
			wantErr: errSelfDelegationTooSmall,
		},
		{
			name: "incorrect signature",
			args: args{
				stateDB:      makeStateDBForMicrostaking(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() microstaking.CreateMap3Node {
					m := defaultMsgCreateMap3Node()
					m.NodeKeySig = blsKeys[12].sig
					return m
				}(),
				signer: createOperatorAddr,
			},
			wantErr: errors.New("bls keys and corresponding signatures could not be verified"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifier, _ := NewStakingVerifier(tt.args.chainContext)
			got, err := verifier.VerifyCreateMap3NodeMsg(tt.args.stateDB, tt.args.chainContext, tt.args.epoch, tt.args.blockNum, &tt.args.msg, tt.args.signer)
			if assErr := assertError(err, tt.wantErr); assErr != nil {
				t.Errorf("Test - %v: %v", tt.name, err)
			}
			if err != nil || tt.wantErr != nil {
				return
			}
			if err := microstaking.CheckMap3NodeWrapperEqual(*got, tt.want); err != nil {
				t.Errorf("Test - %v: %v", tt.name, err)
			}
		})
	}
}

func defaultMsgCreateMap3Node() microstaking.CreateMap3Node {
	pub, sig := blsKeys[11].pub2, blsKeys[11].sig
	cn := microstaking.CreateMap3Node{
		OperatorAddress: createOperatorAddr,
		Description:     defaultDesc2,
		Commission:      defaultCommissionRate,
		NodePubKey:      pub,
		NodeKeySig:      sig,
		Amount:          twoHundredKOnes,
	}
	return cn
}

func defaultExpCreatedMap3Node() microstaking.Map3NodeWrapper_ {
	pub := blsKeys[11].pub2
	v := microstaking.NewMap3NodeWrapperBuilder().
		SetMap3Address(createValidatorAddr).
		SetOperatorAddress(createOperatorAddr).
		AddNodeKey(pub).
		SetCommission(microstaking.Commission_{
			Rate:              defaultCommissionRate,
			RateForNextPeriod: defaultCommissionRate,
			UpdateHeight:      big.NewInt(defaultBlockNumber),
		}).
		SetDescription(defaultDesc2).
		SetCreationHeight(big.NewInt(defaultBlockNumber)).
		SetPendingEpoch(big.NewInt(defaultEpoch)).
		AddMicrodelegation(microstaking.NewMicrodelegation(createOperatorAddr, twoHundredKOnes,
			common.NewDec(defaultEpoch).Add(microstaking.PendingLockInEpoch), true)).
		Build()
	return *v
}

func TestVerifyEditMap3NodeMsg(t *testing.T) {
	type args struct {
		stateDB  vm.StateDB
		epoch    *big.Int
		blockNum *big.Int
		msg      microstaking.EditMap3Node
		signer   common.Address
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "valid request",
			args: args{
				stateDB:  makeStateDBForMicrostaking(t),
				epoch:    big.NewInt(defaultEpoch),
				blockNum: big.NewInt(defaultBlockNumber),
				msg:      defaultMsgEditMap3Node(),
				signer:   map3OperatorAddr,
			},
			wantErr: nil,
		},
		{
			name: "state db nil",
			args: args{
				stateDB:  nil,
				epoch:    big.NewInt(defaultEpoch),
				blockNum: big.NewInt(defaultBlockNumber),
				msg:      defaultMsgEditMap3Node(),
				signer:   map3OperatorAddr,
			},
			wantErr: errStateDBIsMissing,
		},
		{
			name: "epoch nil",
			args: args{
				stateDB:  makeStateDBForMicrostaking(t),
				epoch:    nil,
				blockNum: big.NewInt(defaultBlockNumber),
				msg:      defaultMsgEditMap3Node(),
				signer:   map3OperatorAddr,
			},
			wantErr: errEpochMissing,
		},
		{
			name: "block number nil",
			args: args{
				stateDB:  makeStateDBForMicrostaking(t),
				epoch:    big.NewInt(defaultEpoch),
				blockNum: nil,
				msg:      defaultMsgEditMap3Node(),
				signer:   map3OperatorAddr,
			},
			wantErr: errBlockNumMissing,
		},
		{
			name: "invalid signer",
			args: args{
				stateDB:  makeStateDBForMicrostaking(t),
				epoch:    big.NewInt(defaultEpoch),
				blockNum: big.NewInt(defaultBlockNumber),
				msg: func() microstaking.EditMap3Node {
					msg := defaultMsgEditMap3Node()
					return msg
				}(),
				signer: makeTestAddr("invalid operator"),
			},
			wantErr: errInvalidSigner,
		},
		{
			name: "bls key collision",
			args: args{
				stateDB:  makeStateDBForMicrostaking(t),
				epoch:    big.NewInt(defaultEpoch),
				blockNum: big.NewInt(defaultBlockNumber),
				msg: func() microstaking.EditMap3Node {
					msg := defaultMsgEditMap3Node()
					msg.NodeKeyToAdd = &blsKeys[3].pub2
					msg.NodeKeyToAddSig = &blsKeys[3].sig
					return msg
				}(),
				signer: map3OperatorAddr,
			},
			wantErr: errDupMap3NodePubKey,
		},
		{
			name: "identity collision",
			args: args{
				stateDB:  makeStateDBForMicrostaking(t),
				epoch:    big.NewInt(defaultEpoch),
				blockNum: big.NewInt(defaultBlockNumber),
				msg: func() microstaking.EditMap3Node {
					msg := defaultMsgEditMap3Node()
					msg.Description.Identity = makeIdentityStr(0)
					return msg
				}(),
				signer: map3OperatorAddr,
			},
			wantErr: errDupMap3NodeIdentity,
		},
		{
			name: "map3 node not exist",
			args: args{
				stateDB:  makeStateDBForMicrostaking(t),
				epoch:    big.NewInt(defaultEpoch),
				blockNum: big.NewInt(defaultBlockNumber),
				msg: func() microstaking.EditMap3Node {
					msg := defaultMsgEditMap3Node()
					msg.Map3NodeAddress = makeTestAddr("addr not in chain")
					return msg
				}(),
				signer: map3OperatorAddr,
			},
			wantErr: errMap3NodeNotExist,
		},
		{
			name: "invalid operator",
			args: args{
				stateDB:  makeStateDBForMicrostaking(t),
				epoch:    big.NewInt(defaultEpoch),
				blockNum: big.NewInt(defaultBlockNumber),
				msg: func() microstaking.EditMap3Node {
					msg := defaultMsgEditMap3Node()
					msg.OperatorAddress = makeTestAddr("invalid operator")
					return msg
				}(),
				signer: makeTestAddr("invalid operator"),
			},
			wantErr: errInvalidMap3NodeOperator,
		},
		{
			name: "signature cannot be verified",
			args: args{
				stateDB:  makeStateDBForMicrostaking(t),
				epoch:    big.NewInt(defaultEpoch),
				blockNum: big.NewInt(defaultBlockNumber),
				msg: func() microstaking.EditMap3Node {
					msg := defaultMsgEditMap3Node()
					msg.NodeKeyToAddSig = &blsKeys[13].sig
					return msg
				}(),
				signer: map3OperatorAddr,
			},
			wantErr: errors.New("bls keys and corresponding signatures could not be verified"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifier, _ := NewStakingVerifier(makeFakeChainContextForStake(t))
			err := verifier.VerifyEditMap3NodeMsg(tt.args.stateDB, tt.args.epoch, tt.args.blockNum, &tt.args.msg, tt.args.signer)
			if assErr := assertError(err, tt.wantErr); assErr != nil {
				t.Errorf("Test - %v: %v", tt.name, err)
			}
		})
	}
}

var (
	editDesc2 = microstaking.Description_{
		Name:            "batman",
		Identity:        "batman",
		Website:         "",
		SecurityContact: "",
		Details:         "",
	}
)

func defaultMsgEditMap3Node() microstaking.EditMap3Node {
	var (
		pub0Copy  microstaking.BLSPublicKey_
		pub12Copy microstaking.BLSPublicKey_
		sig12Copy common2.BLSSignature
	)
	copy(pub0Copy.Key[:], blsKeys[0].pub2.Key[:])
	copy(pub12Copy.Key[:], blsKeys[12].pub2.Key[:])
	copy(sig12Copy[:], blsKeys[12].sig[:])

	return microstaking.EditMap3Node{
		Map3NodeAddress: map3NodeAddr,
		OperatorAddress: map3OperatorAddr,
		Description:     editDesc2,
		NodeKeyToRemove: &pub0Copy,
		NodeKeyToAdd:    &pub12Copy,
		NodeKeyToAddSig: &sig12Copy,
	}
}

func TestVerifyTerminateMap3NodeMsg(t *testing.T) {
	type args struct {
		stateDB vm.StateDB
		epoch   *big.Int
		msg     microstaking.TerminateMap3Node
		signer  common.Address
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "terminate successfully",
			args: args{
				stateDB: makeStateForTerminating(t),
				epoch:   big.NewInt(defaultEpoch),
				msg:     defaultMsgTerminateMap3Node(),
				signer:  map3OperatorAddr,
			},
			wantErr: nil,
		},
		{
			name: "state db nil",
			args: args{
				stateDB: nil,
				epoch:   big.NewInt(defaultEpoch),
				msg:     defaultMsgTerminateMap3Node(),
				signer:  map3OperatorAddr,
			},
			wantErr: errStateDBIsMissing,
		},
		{
			name: "epoch nil",
			args: args{
				stateDB: makeStateForTerminating(t),
				epoch:   nil,
				msg:     defaultMsgTerminateMap3Node(),
				signer:  map3OperatorAddr,
			},
			wantErr: errEpochMissing,
		},
		{
			name: "invalid signer",
			args: args{
				stateDB: makeStateForTerminating(t),
				epoch:   big.NewInt(defaultEpoch),
				msg:     defaultMsgTerminateMap3Node(),
				signer:  makeTestAddr("invalid operator"),
			},
			wantErr: errInvalidSigner,
		},
		{
			name: "map3 node not exist",
			args: args{
				stateDB: makeStateForTerminating(t),
				epoch:   big.NewInt(defaultEpoch),
				msg: func() microstaking.TerminateMap3Node {
					msg := defaultMsgTerminateMap3Node()
					msg.Map3NodeAddress = makeTestAddr("addr not in chain")
					return msg
				}(),
				signer: map3OperatorAddr,
			},
			wantErr: errMap3NodeNotExist,
		},
		{
			name: "invalid operator",
			args: args{
				stateDB: makeStateForTerminating(t),
				epoch:   big.NewInt(defaultEpoch),
				msg: func() microstaking.TerminateMap3Node {
					msg := defaultMsgTerminateMap3Node()
					msg.OperatorAddress = makeTestAddr("invalid operator")
					return msg
				}(),
				signer: makeTestAddr("invalid operator"),
			},
			wantErr: errInvalidMap3NodeOperator,
		},
		{
			name: "invalid status",
			args: args{
				stateDB: makeStateForTerminating(t),
				epoch:   big.NewInt(defaultEpoch),
				msg: func() microstaking.TerminateMap3Node {
					msg := defaultMsgTerminateMap3Node()
					msg.Map3NodeAddress = map3NodeAddr2
					msg.OperatorAddress = map3OperatorAddr2
					return msg
				}(),
				signer: map3OperatorAddr2,
			},
			wantErr: errTerminateMap3NodeNotAllowed,
		},
		{
			name: "microdelegation still locked",
			args: args{
				stateDB: makeStateForTerminating(t),
				epoch:   big.NewInt(defaultEpoch - 1),
				msg: func() microstaking.TerminateMap3Node {
					msg := defaultMsgTerminateMap3Node()
					msg.Map3NodeAddress = map3NodeAddr3
					msg.OperatorAddress = map3OperatorAddr3
					return msg
				}(),
				signer: map3OperatorAddr3,
			},
			wantErr: errMicrodelegationStillLocked,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifier, _ := NewStakingVerifier(makeFakeChainContextForStake(t))
			err := verifier.VerifyTerminateMap3NodeMsg(tt.args.stateDB, tt.args.epoch, &tt.args.msg, tt.args.signer)
			if assErr := assertError(err, tt.wantErr); assErr != nil {
				t.Errorf("Test - %v: %v", tt.name, err)
			}
		})
	}
}

func defaultMsgTerminateMap3Node() microstaking.TerminateMap3Node {
	return microstaking.TerminateMap3Node{
		Map3NodeAddress: map3NodeAddr,
		OperatorAddress: map3OperatorAddr,
	}
}

func makeStateForTerminating(t *testing.T) *state.StateDB {
	sdb := makeStateDBForMicrostaking(t)
	if err := changeMap3StatusForAddr(sdb, map3NodeAddr2, microstaking.Active); err != nil {
		t.Fatal(err)
	}
	sdb.IntermediateRoot(true)
	return sdb
}

func changeMap3StatusForAddr(sdb *state.StateDB, map3Addr common.Address, status microstaking.Map3Status) error {
	n, err := sdb.Map3NodeByAddress(map3Addr)
	if err != nil {
		return err
	}
	n.Map3Node().Status().SetValue(uint8(status))
	return nil
}

func TestVerifyMicrodelegateMsg(t *testing.T) {
	type args struct {
		stateDB      vm.StateDB
		chainContext ChainContext
		blockNum     *big.Int
		msg          microstaking.Microdelegate
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
				stateDB:      makeStateForMicrodelegating(t),
				chainContext: makeFakeChainContextForStake(t),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg:          defaultMsgMicrodelegate(),
				signer:       delegatorAddr,
			},
			wantErr: nil,
		},
		{
			name: "state db nil",
			args: args{
				stateDB:      nil,
				chainContext: makeFakeChainContextForStake(t),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg:          defaultMsgMicrodelegate(),
				signer:       delegatorAddr,
			},
			wantErr: errStateDBIsMissing,
		},
		{
			name: "chain context nil",
			args: args{
				stateDB:      makeStateForMicrodelegating(t),
				chainContext: nil,
				blockNum:     big.NewInt(defaultBlockNumber),
				msg:          defaultMsgMicrodelegate(),
				signer:       createOperatorAddr,
			},
			wantErr: errChainContextMissing,
		},
		{
			name: "block number nil",
			args: args{
				stateDB:      makeStateForMicrodelegating(t),
				chainContext: makeFakeChainContextForStake(t),
				blockNum:     nil,
				msg:          defaultMsgMicrodelegate(),
				signer:       createOperatorAddr,
			},
			wantErr: errBlockNumMissing,
		},
		{
			name: "invalid signer",
			args: args{
				stateDB:      makeStateForMicrodelegating(t),
				chainContext: makeFakeChainContextForStake(t),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg:          defaultMsgMicrodelegate(),
				signer:       makeTestAddr("invalid operator"),
			},
			wantErr: errInvalidSigner,
		},
		{
			name: "negative amount",
			args: args{
				stateDB:      makeStateForMicrodelegating(t),
				chainContext: makeFakeChainContextForStake(t),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() microstaking.Microdelegate {
					m := defaultMsgMicrodelegate()
					m.Amount = big.NewInt(-1)
					return m
				}(),
				signer: delegatorAddr,
			},
			wantErr: errNegativeAmount,
		},
		{
			name: "map3 node not exist",
			args: args{
				stateDB:      makeStateForMicrodelegating(t),
				chainContext: makeFakeChainContextForStake(t),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() microstaking.Microdelegate {
					msg := defaultMsgMicrodelegate()
					msg.Map3NodeAddress = makeTestAddr("addr not in chain")
					return msg
				}(),
				signer: delegatorAddr,
			},
			wantErr: errMap3NodeNotExist,
		},
		{
			name: "invalid status",
			args: args{
				stateDB:      makeStateForMicrodelegating(t),
				chainContext: makeFakeChainContextForStake(t),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() microstaking.Microdelegate {
					msg := defaultMsgMicrodelegate()
					msg.Map3NodeAddress = map3NodeAddr2
					return msg
				}(),
				signer: delegatorAddr,
			},
			wantErr: errInvalidNodeStatusForDelegation,
		},
		{
			name: "insufficient balance",
			args: args{
				stateDB:      makeStateForMicrodelegating(t),
				chainContext: makeFakeChainContextForStake(t),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() microstaking.Microdelegate {
					m := defaultMsgMicrodelegate()
					m.Amount = oneMill
					return m
				}(),
				signer: delegatorAddr,
			},
			wantErr: errInsufficientBalanceForStake,
		},
		{
			name: "delegation too small",
			args: args{
				stateDB:      makeStateForMicrodelegating(t),
				chainContext: makeFakeChainContextForStake(t),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() microstaking.Microdelegate {
					m := defaultMsgMicrodelegate()
					m.Amount = oneBig
					return m
				}(),
				signer: delegatorAddr,
			},
			wantErr: errDelegationTooSmall,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifier, _ := NewStakingVerifier(tt.args.chainContext)
			err := verifier.VerifyMicrodelegateMsg(tt.args.stateDB, tt.args.chainContext, tt.args.blockNum, &tt.args.msg, tt.args.signer)
			if assErr := assertError(err, tt.wantErr); assErr != nil {
				t.Errorf("Test - %v: %v", tt.name, err)
			}
		})
	}
}

func defaultMsgMicrodelegate() microstaking.Microdelegate {
	return microstaking.Microdelegate{
		Map3NodeAddress:  map3NodeAddr,
		DelegatorAddress: delegatorAddr,
		Amount:           tenKOnes,
	}
}

func makeStateForMicrodelegating(t *testing.T) *state.StateDB {
	sdb := makeStateDBForMicrostaking(t)
	if err := changeMap3StatusForAddr(sdb, map3NodeAddr2, microstaking.Active); err != nil {
		t.Fatal(err)
	}
	sdb.IntermediateRoot(true)
	return sdb
}

func TestVerifyUnmicrodelegateMsg(t *testing.T) {
	type args struct {
		stateDB      vm.StateDB
		chainContext ChainContext
		blockNum     *big.Int
		epoch        *big.Int
		msg          microstaking.Unmicrodelegate
		signer       common.Address
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "unmicrodelegate successfully",
			args: args{
				stateDB:      makeDefaultStateForUnmicrodelegate(t),
				chainContext: makeFakeChainContextForStake(t),
				blockNum:     big.NewInt(defaultBlockNumber),
				epoch:        big.NewInt(defaultEpoch + 1),
				msg:          defaultMsgUnmicrodelegate(),
				signer:       map3OperatorAddr,
			},
			wantErr: nil,
		},
		{
			name: "state db nil",
			args: args{
				stateDB:      nil,
				chainContext: makeFakeChainContextForStake(t),
				blockNum:     big.NewInt(defaultBlockNumber),
				epoch:        big.NewInt(defaultEpoch + 1),
				msg:          defaultMsgUnmicrodelegate(),
				signer:       map3OperatorAddr,
			},
			wantErr: errStateDBIsMissing,
		},
		{
			name: "chain context nil",
			args: args{
				stateDB:      makeDefaultStateForUnmicrodelegate(t),
				chainContext: nil,
				blockNum:     big.NewInt(defaultBlockNumber),
				epoch:        big.NewInt(defaultEpoch + 1),
				msg:          defaultMsgUnmicrodelegate(),
				signer:       map3OperatorAddr,
			},
			wantErr: errChainContextMissing,
		},
		{
			name: "block number nil",
			args: args{
				stateDB:      makeDefaultStateForUnmicrodelegate(t),
				chainContext: makeFakeChainContextForStake(t),
				blockNum:     nil,
				epoch:        big.NewInt(defaultEpoch + 1),
				msg:          defaultMsgUnmicrodelegate(),
				signer:       map3OperatorAddr,
			},
			wantErr: errBlockNumMissing,
		},
		{
			name: "epoch nil",
			args: args{
				stateDB:      makeDefaultStateForUnmicrodelegate(t),
				chainContext: makeFakeChainContextForStake(t),
				blockNum:     big.NewInt(defaultBlockNumber),
				epoch:        nil,
				msg:          defaultMsgUnmicrodelegate(),
				signer:       map3OperatorAddr,
			},
			wantErr: errEpochMissing,
		},
		{
			name: "negative amount",
			args: args{
				stateDB:      makeDefaultStateForUnmicrodelegate(t),
				chainContext: makeFakeChainContextForStake(t),
				blockNum:     big.NewInt(defaultBlockNumber),
				epoch:        big.NewInt(defaultEpoch + 1),
				msg: func() microstaking.Unmicrodelegate {
					m := defaultMsgUnmicrodelegate()
					m.Amount = big.NewInt(-1)
					return m
				}(),
				signer: map3OperatorAddr,
			},
			wantErr: errNegativeAmount,
		},
		{
			name: "invalid signer",
			args: args{
				stateDB:      makeDefaultStateForUnmicrodelegate(t),
				chainContext: makeFakeChainContextForStake(t),
				blockNum:     big.NewInt(defaultBlockNumber),
				epoch:        big.NewInt(defaultEpoch + 1),
				msg:          defaultMsgUnmicrodelegate(),
				signer:       makeTestAddr("invalid delegator"),
			},
			wantErr: errInvalidSigner,
		},
		{
			name: "map3 node not exist",
			args: args{
				stateDB:      makeDefaultStateForUnmicrodelegate(t),
				chainContext: makeFakeChainContextForStake(t),
				blockNum:     big.NewInt(defaultBlockNumber),
				epoch:        big.NewInt(defaultEpoch + 1),
				msg: func() microstaking.Unmicrodelegate {
					msg := defaultMsgUnmicrodelegate()
					msg.Map3NodeAddress = makeTestAddr("addr not in chain")
					return msg
				}(),
				signer: map3OperatorAddr,
			},
			wantErr: errMap3NodeNotExist,
		},
		{
			name: "invalid status",
			args: args{
				stateDB:      makeDefaultStateForUnmicrodelegate(t),
				chainContext: makeFakeChainContextForStake(t),
				blockNum:     big.NewInt(defaultBlockNumber),
				epoch:        big.NewInt(defaultEpoch + 1),
				msg: func() microstaking.Unmicrodelegate {
					msg := defaultMsgUnmicrodelegate()
					msg.Map3NodeAddress = map3NodeAddr2
					return msg
				}(),
				signer: map3OperatorAddr,
			},
			wantErr: errUnmicrodelegateNotAllowed,
		},
		{
			name: "microdelegation not exist",
			args: args{
				stateDB:      makeDefaultStateForUnmicrodelegate(t),
				chainContext: makeFakeChainContextForStake(t),
				blockNum:     big.NewInt(defaultBlockNumber),
				epoch:        big.NewInt(defaultEpoch + 1),
				msg: func() microstaking.Unmicrodelegate {
					msg := defaultMsgUnmicrodelegate()
					msg.DelegatorAddress = makeTestAddr("addr not in chain")
					return msg
				}(),
				signer: makeTestAddr("addr not in chain"),
			},
			wantErr: errMicrodelegationNotExist,
		},
		{
			name: "insufficient balance to unmicrodelegate",
			args: args{
				stateDB:      makeDefaultStateForUnmicrodelegate(t),
				chainContext: makeFakeChainContextForStake(t),
				blockNum:     big.NewInt(defaultBlockNumber),
				epoch:        big.NewInt(defaultEpoch + 1),
				msg: func() microstaking.Unmicrodelegate {
					msg := defaultMsgUnmicrodelegate()
					msg.Amount = big.NewInt(0).Add(oneMill, common.Big1)
					return msg
				}(),
				signer: map3OperatorAddr,
			},
			wantErr: errInsufficientBalanceToUnmicrodelegate,
		},
		{
			name: "microdelegation still locked",
			args: args{
				stateDB:      makeDefaultStateForUnmicrodelegate(t),
				chainContext: makeFakeChainContextForStake(t),
				blockNum:     big.NewInt(defaultBlockNumber),
				epoch:        big.NewInt(defaultEpoch),
				msg:          defaultMsgUnmicrodelegate(),
				signer:       map3OperatorAddr,
			},
			wantErr: errMicrodelegationStillLocked,
		},
		{
			name: "self delegation too small",
			args: args{
				stateDB:      makeDefaultStateForUnmicrodelegate(t),
				chainContext: makeFakeChainContextForStake(t),
				epoch:        big.NewInt(defaultEpoch + 1),
				blockNum:     big.NewInt(defaultBlockNumber),
				msg: func() microstaking.Unmicrodelegate {
					m := defaultMsgUnmicrodelegate()
					m.Amount = big.NewInt(0).Sub(oneMill, common.Big1)
					return m
				}(),
				signer: map3OperatorAddr,
			},
			wantErr: errSelfDelegationTooSmall,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifier, _ := NewStakingVerifier(tt.args.chainContext)
			err := verifier.VerifyUnmicrodelegateMsg(tt.args.stateDB, tt.args.chainContext, tt.args.blockNum, tt.args.epoch, &tt.args.msg, tt.args.signer)
			if assErr := assertError(err, tt.wantErr); assErr != nil {
				t.Errorf("Test - %v: %v", tt.name, err)
			}
		})
	}
}

func makeDefaultStateForUnmicrodelegate(t *testing.T) *state.StateDB {
	sdb := makeStateDBForMicrostaking(t)
	if err := changeMap3StatusForAddr(sdb, map3NodeAddr2, microstaking.Active); err != nil {
		t.Fatal(err)
	}
	sdb.IntermediateRoot(true)
	return sdb
}

// undelegate from delegator which has already go one entry for undelegation
func defaultMsgUnmicrodelegate() microstaking.Unmicrodelegate {
	return microstaking.Unmicrodelegate{
		Map3NodeAddress:  map3NodeAddr,
		DelegatorAddress: map3OperatorAddr,
		Amount:           tenKOnes,
	}
}

func TestVerifyCollectMicordelRewardsMsg(t *testing.T) {
	type args struct {
		stateDB vm.StateDB
		msg     microstaking.CollectRewards
		signer  common.Address
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "collect successfully",
			args: args{
				stateDB: makeStateForMicrostakingReward(t),
				msg:     defaultMsgCollectMicrodelRewards(),
				signer:  map3OperatorAddr,
			},
			wantErr: nil,
		},
		{
			name: "state db nil",
			args: args{
				stateDB: nil,
				msg:     defaultMsgCollectMicrodelRewards(),
				signer:  map3OperatorAddr,
			},
			wantErr: errStateDBIsMissing,
		},
		{
			name: "invalid signer",
			args: args{
				stateDB: makeStateForMicrostakingReward(t),
				msg:     defaultMsgCollectMicrodelRewards(),
				signer:  makeTestAddr("invalid operator"),
			},
			wantErr: errInvalidSigner,
		},
		{
			name: "no reward",
			args: args{
				stateDB: makeStateForMicrostakingReward(t),
				msg: func() microstaking.CollectRewards {
					m := defaultMsgCollectMicrodelRewards()
					m.DelegatorAddress = operatorAddr2
					return m
				}(),
				signer: operatorAddr2,
			},
			wantErr: errNoRewardsToCollect,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifier, _ := NewStakingVerifier(makeFakeChainContextForStake(t))
			err := verifier.VerifyCollectMicrostakingRewardsMsg(tt.args.stateDB, &tt.args.msg, tt.args.signer)
			if assErr := assertError(err, tt.wantErr); assErr != nil {
				t.Errorf("Test - %v: %v", tt.name, err)
			}
		})
	}
}

// undelegate from delegator which has already go one entry for undelegation
func defaultMsgCollectMicrodelRewards() microstaking.CollectRewards {
	return microstaking.CollectRewards{
		DelegatorAddress: map3OperatorAddr,
	}
}

func makeStateForMicrostakingReward(t *testing.T) *state.StateDB {
	sdb := makeStateDBForMicrostaking(t)
	if err := addStateMicrostakingRewardForAddr(sdb, map3NodeAddr, map3OperatorAddr, reward00); err != nil {
		t.Fatal(err)
	}
	sdb.IntermediateRoot(true)
	return sdb
}

func addStateMicrostakingRewardForAddr(sdb *state.StateDB, map3Addr, delegator common.Address, reward *big.Int) error {
	w, err := sdb.Map3NodeByAddress(map3Addr)
	if err != nil {
		return err
	}
	redelegation, ok := w.Microdelegations().Get(delegator)
	if !ok {
		return errMicrodelegationNotExist
	}
	redelegation.AddReward(reward)
	return nil
}

// makeStateDBForMicrostaking make the default state db for restaking test
func makeStateDBForMicrostaking(t *testing.T) *state.StateDB {
	sdb, err := newTestStateDB()
	if err != nil {
		t.Fatal(err)
	}
	ws := makeNodeWrappersForMicrostaking(defNumNodeWrappersInState, defNumPubPerNode)
	if err := updateStateMap3Nodes(sdb, ws); err != nil {
		t.Fatalf("make default state: %v", err)
	}
	sdb.SetNonce(createOperatorAddr, defaultNonce)
	sdb.AddBalance(createOperatorAddr, twoHundredKOnes)
	sdb.AddBalance(delegatorAddr, tenKOnes)
	sdb.Commit(true)
	return sdb
}

func makeNodeWrappersForMicrostaking(num, numPubsPerNode int) []*microstaking.Map3NodeWrapper_ {
	ws := make([]*microstaking.Map3NodeWrapper_, 0, num)
	pubGetter := newBLSPubGetter(blsKeys)
	for i := 0; i != num; i++ {
		w := makeStateNodeWrapperFromGetter(i, numPubsPerNode, pubGetter)
		ws = append(ws, w)
	}
	return ws
}

func makeStateNodeWrapperFromGetter(index int, numPubs int, pubGetter *BLSPubGetter) *microstaking.Map3NodeWrapper_ {
	map3Addr := makeTestAddr(fmt.Sprint("map3", index))
	operator := makeTestAddr(fmt.Sprint("op", index))
	pubs := microstaking.NewEmptyBLSKeys()
	for i := 0; i != numPubs; i++ {
		pub := pubGetter.getPub2()
		pubs.Keys = append(pubs.Keys, &pub)
	}
	w := microstaking.NewMap3NodeWrapperBuilder().
		SetMap3Address(map3Addr).
		SetOperatorAddress(operator).
		AddNodeKeys(pubs).
		SetDescription(defaultDesc2).
		SetCommission(defaultCommission).
		AddMicrodelegation(microstaking.NewMicrodelegation(operator, defaultDelAmount, common.NewDec(defaultEpoch), true)).
		Build()
	w.Map3Node.Description.Identity = makeIdentityStr(index)
	return w
}

func updateStateMap3Nodes(sdb *state.StateDB, ws []*microstaking.Map3NodeWrapper_) error {
	pool := sdb.Map3NodePool()
	for _, w := range ws {
		pool.Nodes().Put(w.Map3Node.Map3Address, w)
		sdb.IncrementMap3NodeNonce()
		for _, k := range w.Map3Node.NodeKeys.Keys {
			pool.NodeKeySet().Get(k.Hex()).SetValue(true)
		}
		pool.DescriptionIdentitySet().Get(w.Map3Node.Description.Identity).SetValue(true)

		for _, key := range w.Microdelegations.Keys {
			delegator := *key
			index := microstaking.DelegationIndex_{
				Map3Address: w.Map3Node.Map3Address,
				IsOperator:  delegator == w.Map3Node.OperatorAddress,
			}
			pool.UpdateDelegationIndex(delegator, &index)
		}
	}
	return nil
}
