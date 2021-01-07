// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/consensus/atlas/backend"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/hyperion-hyn/bls/ffi/go/bls"

	"gopkg.in/karalabe/cookiejar.v2/collections/prque"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/atlas"
	"github.com/ethereum/go-ethereum/consensus/atlas/validator"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	elog "github.com/ethereum/go-ethereum/log"
)

var testLogger = elog.New()

type testSystemBackend struct {
	id  uint64
	sys *testSystem

	engine Engine
	peers  atlas.ValidatorSet
	events *event.TypeMux

	committedMsgs []testCommittedMsgs
	sentMsgs      [][]byte // store the message when Send is called by core

	signers   []common.Address // signer's id (address format)
	signerKey *bls.SecretKey
	address   common.Address // owner's address
	db        ethdb.Database
}

func (self *testSystemBackend) Signer() []common.Address {
	return self.signers
}

func (self *testSystemBackend) Annotation() string {
	return fmt.Sprintf("backend-%d", self.id)
}

type testCommittedMsgs struct {
	commitProposal atlas.Proposal
	committedSeals []byte
}

// ==============================================
//
// define the functions that needs to be provided for Istanbul.

func (self *testSystemBackend) Address() common.Address {
	return self.address
}

// Peers returns all connected peers
func (self *testSystemBackend) Validators(proposal atlas.Proposal) atlas.ValidatorSet {
	return self.peers
}

func (self *testSystemBackend) EventMux() *event.TypeMux {
	return self.events
}

func (self *testSystemBackend) Send(message []byte, target common.Address) error {
	testLogger.Info("enqueuing a message...", "address", self.Signer())
	self.sentMsgs = append(self.sentMsgs, message)
	self.sys.queuedMessage <- atlas.MessageEvent{
		Payload: message,
	}
	return nil
}

func (self *testSystemBackend) Broadcast(valSet atlas.ValidatorSet, message []byte) error {
	testLogger.Info("enqueuing a message...", "address", self.Signer())
	self.sentMsgs = append(self.sentMsgs, message)
	self.sys.queuedMessage <- atlas.MessageEvent{
		Payload: message,
	}
	return nil
}

func (self *testSystemBackend) Gossip(valSet atlas.ValidatorSet, message []byte) error {
	testLogger.Warn("Gossip")
	return nil
}

func (self *testSystemBackend) Commit(proposal atlas.Proposal, signature []byte, bitmap []byte) error {
	// func (self *testSystemBackend) Commit(proposal atlas.Proposal, seals [][]byte) error {
	testLogger.Info("commit message", "address", self.Signer())
	self.committedMsgs = append(self.committedMsgs, testCommittedMsgs{
		commitProposal: proposal,
		// TODO(zgx): how to commit a message?
		committedSeals: signature,
	})

	// fake new head events
	go self.events.Post(atlas.FinalCommittedEvent{})
	return nil
}

func (self *testSystemBackend) Verify(proposal atlas.Proposal) (time.Duration, error) {
	return 0, nil
}

func (self *testSystemBackend) SignHash(singer common.Address, hash common.Hash) ([]byte, []byte, []byte, error) {
	testLogger.Info(fmt.Sprintf("testSystemBackend.Sign: %x...", hash.Hex()))
	sighash := self.signerKey.SignHash(hash.Bytes()).Serialize()
	pubkey := self.signerKey.GetPublicKey().Serialize()
	testLogger.Info(fmt.Sprintf("testSystemBackend.Sign: data: %x... signature: %x, publicKey: %x", hash.Hex(), sighash[:10], pubkey[:10]))
	return sighash, pubkey, nil, nil
}

func (self *testSystemBackend) CheckSignature([]byte, []byte, []byte) error {
	return nil
}

func (self *testSystemBackend) CheckValidatorSignature(hash common.Hash, sig []byte, pubKey []byte) error {
	return atlas.CheckValidatorSignature(hash.Bytes(), sig, pubKey)
}

func (self *testSystemBackend) Hash(b interface{}) common.Hash {
	return common.StringToHash("Test")
}

func (self *testSystemBackend) NewRequest(request atlas.Proposal) {
	go self.events.Post(atlas.RequestEvent{
		Proposal: request,
	})
}

func (self *testSystemBackend) HasBadProposal(hash common.Hash) bool {
	return false
}

func (self *testSystemBackend) LastProposal() (atlas.Proposal, common.Address) {
	l := len(self.committedMsgs)
	if l > 0 {
		return self.committedMsgs[l-1].commitProposal, common.Address{}
	}
	return makeBlock(0), common.Address{}
}

// Only block height 5 will return true
func (self *testSystemBackend) HasPropsal(hash common.Hash, number *big.Int) bool {
	return number.Cmp(big.NewInt(5)) == 0
}

func (self *testSystemBackend) GetProposer(number uint64) common.Address {
	return common.Address{}
}

func (self *testSystemBackend) ParentValidators(proposal atlas.Proposal) atlas.ValidatorSet {
	return self.peers
}

func (self *testSystemBackend) SealHash(header *types.Header) common.Hash {
	return backend.SealHash(header)
}

func (sb *testSystemBackend) Close() error {
	return nil
}

// ==============================================
//
// define the struct that need to be provided for integration tests.

type testSystem struct {
	backends []*testSystemBackend

	queuedMessage chan atlas.MessageEvent
	quit          chan struct{}
}

func newTestSystem(n uint64) *testSystem {
	testLogger.SetHandler(elog.StdoutHandler)
	return &testSystem{
		backends: make([]*testSystemBackend, n),

		queuedMessage: make(chan atlas.MessageEvent),
		quit:          make(chan struct{}),
	}
}

func generateValidators(n int) ([]atlas.Validator, []*bls.SecretKey) {
	vals := make([]atlas.Validator, n)
	keys := make([]*bls.SecretKey, n)
	for i := 0; i < n; i++ {
		privateKey, _ := crypto.GenerateKey()
		secretKey, _ := crypto.GenerateBLSKey()
		val, _ := validator.New(secretKey.GetPublicKey().Serialize(), crypto.PubkeyToAddress(privateKey.PublicKey))
		keys[i] = secretKey
		vals[i] = val
	}
	return vals, keys
}

func newTestValidatorSet(n int) atlas.ValidatorSet {
	vals, _ := generateValidators(n)
	return validator.NewSet(vals, atlas.RoundRobin)
}

func findSecretKeyBySigner(secretKeys []*bls.SecretKey, address common.Address) *bls.SecretKey {
	for _, v := range secretKeys {
		if bytes.Compare(crypto.PubkeyToSigner(v.GetPublicKey()).Bytes(), address.Bytes()) == 0 {
			return v
		}
	}
	return nil
}

// FIXME: int64 is needed for N and F
func NewTestSystemWithBackend(n, f uint64) *testSystem {
	testLogger.SetHandler(elog.StdoutHandler)
	vals, keys := generateValidators(int(n))
	sys := newTestSystem(n)
	config := atlas.DefaultConfig

	for i := uint64(0); i < n; i++ {
		vset := validator.NewSet(vals, atlas.RoundRobin)
		backend := sys.NewBackend(i)
		backend.peers = vset
		backend.address = vset.GetByIndex(i).Coinbase()
		backend.signers = []common.Address{vset.GetByIndex(i).Signer()}
		backend.signerKey = findSecretKeyBySigner(keys, vset.GetByIndex(i).Signer())

		core := New(backend, config).(*core)
		core.state = StateAcceptRequest
		core.current = newRoundState(&atlas.View{
			Round:    big.NewInt(0),
			Sequence: big.NewInt(1),
		}, vset, common.Hash{}, nil, nil, func(hash common.Hash) bool {
			return false
		})
		core.valSet = vset
		core.logger = testLogger
		core.validateHashFn = backend.CheckValidatorSignature
		core.backlogs = make(map[common.Address]*prque.Prque)
		core.backlogsMu = new(sync.Mutex)
		core.roundChangeSet = newRoundChangeSet(core.valSet)
		backend.engine = core
	}

	return sys
}

// listen will consume messages from queue and deliver a message to core
func (t *testSystem) listen() {
	for {
		select {
		case <-t.quit:
			return
		case queuedMessage := <-t.queuedMessage:
			testLogger.Info("consuming a queue message...")
			for _, backend := range t.backends {
				go backend.EventMux().Post(queuedMessage)
			}
		}
	}
}

// Run will start system components based on given flag, and returns a closer
// function that caller can control lifecycle
//
// Given a true for core if you want to initialize core engine.
func (t *testSystem) Run(core bool) func() {
	for _, b := range t.backends {
		if core {
			b.engine.Start() // start Atlas core
		}
	}

	go t.listen()
	closer := func() { t.stop(core) }
	return closer
}

func (t *testSystem) stop(core bool) {
	close(t.quit)

	for _, b := range t.backends {
		if core {
			b.engine.Stop()
		}
	}
}

func (t *testSystem) NewBackend(id uint64) *testSystemBackend {
	// assume always success
	ethDB := rawdb.NewMemoryDatabase()
	backend := &testSystemBackend{
		id:     id,
		sys:    t,
		events: new(event.TypeMux),
		db:     ethDB,
	}

	t.backends[id] = backend
	return backend
}

// ==============================================
//
// helper functions.

func getPublicKeyAddress(privateKey *ecdsa.PrivateKey) common.Address {
	return crypto.PubkeyToAddress(privateKey.PublicKey)
}

func init() {
	elog.Root().SetHandler(elog.LvlFilterHandler(elog.Lvl(elog.LvlDebug), elog.StreamHandler(os.Stdout, elog.TerminalFormat(true))))
}
