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

package backend

import (
	"math/big"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/hyperion-hyn/bls/ffi/go/bls"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/atlas"
	kernel "github.com/ethereum/go-ethereum/consensus/atlas/core"
	"github.com/ethereum/go-ethereum/consensus/atlas/validator"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
)

const (
	// fetcherID is the ID indicates the block is from Atlas engine
	fetcherID = "atlas"
)

var (
	lastCommitsKey = []byte("LastCommits")
)

// New creates an Ethereum backend for Atlas core engine.
func New(config *atlas.Config, db ethdb.Database) consensus.Atlas {
	// Allocate the snapshot caches and create the engine
	recents, _ := lru.NewARC(inmemorySnapshots)
	recentMessages, _ := lru.NewARC(inmemoryPeers)
	knownMessages, _ := lru.NewARC(inmemoryMessages)

	backend := &backend{
		config:         config,
		atlasEventMux:  new(event.TypeMux),
		logger:         log.New(),
		db:             db,
		commitCh:       make(chan *types.Block, 1),
		recents:        recents,
		candidates:     make(map[common.Address]ValidatorProposal),
		coreStarted:    false,
		recentMessages: recentMessages,
		knownMessages:  knownMessages,
	}
	backend.core = kernel.New(backend, backend.config)
	return backend
}

// ----------------------------------------------------------------------------
type ValidatorProposal struct {
	Signer   common.Address
	Coinbase common.Address
	Vote     bool
}

type backend struct {
	config        *atlas.Config
	atlasEventMux *event.TypeMux

	address    common.Address       // Signer's account address
	signer     common.Address       // Signer's id (address format)
	signHashFn consensus.SignHashFn // Sign function to authorize hashes with
	lock       sync.RWMutex         // Protects the signer fields

	core         kernel.Engine
	logger       log.Logger
	db           ethdb.Database
	chain        consensus.ChainReader
	currentBlock func() *types.Block
	hasBadBlock  func(hash common.Hash) bool

	// the channels for atlas engine notifications
	commitCh          chan *types.Block
	proposedBlockHash common.Hash
	sealMu            sync.Mutex
	coreStarted       bool
	coreMu            sync.RWMutex

	// Current list of candidates we are pushing
	candidates map[common.Address]ValidatorProposal
	// Protects the signer fields
	candidatesLock sync.RWMutex
	// Snapshots for recent block to speed up reorgs
	recents *lru.ARCCache

	// event subscription for ChainHeadEvent event
	broadcaster consensus.Broadcaster

	recentMessages *lru.ARCCache // the cache of peer's messages
	knownMessages  *lru.ARCCache // the cache of self messages
}

// zekun: HACK
func (sb *backend) CalcDifficulty(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	return new(big.Int)
}

// Address implements atlas.Backend.Address
func (sb *backend) Address() common.Address {
	return sb.address
}

// Signer implements atlas.Backend.Signer
func (sb *backend) Signer() common.Address {
	return sb.signer
}

// Validators implements atlas.Backend.Validators
func (sb *backend) Validators(proposal atlas.Proposal) atlas.ValidatorSet {
	return sb.getValidators(proposal.Number().Uint64(), proposal.Hash())
}

// Broadcast implements atlas.Backend.Broadcast
func (sb *backend) Broadcast(valSet atlas.ValidatorSet, payload []byte) error {
	// send to others
	sb.Gossip(valSet, payload)
	// send to self
	msg := atlas.MessageEvent{
		Payload: payload,
	}
	go sb.atlasEventMux.Post(msg)
	return nil
}

// Broadcast implements atlas.Backend.Gossip
func (sb *backend) Gossip(valSet atlas.ValidatorSet, payload []byte) error {
	hash := atlas.RLPHash(payload)
	sb.knownMessages.Add(hash, true)

	if sb.broadcaster != nil {
		ps := sb.broadcaster.FindPeers(nil)
		for addr, p := range ps {
			ms, ok := sb.recentMessages.Get(addr)
			var m *lru.ARCCache
			if ok {
				m, _ = ms.(*lru.ARCCache)
				if _, k := m.Get(hash); k {
					// This peer had this event, skip it
					continue
				}
			} else {
				m, _ = lru.NewARC(inmemoryMessages)
			}

			m.Add(hash, true)
			sb.recentMessages.Add(addr, m)

			go p.Send(AtlasMsg, payload)
		}
	}
	return nil
}

// Commit implements atlas.Backend.Commit
func (sb *backend) Commit(proposal atlas.Proposal, signature []byte, publicKey []byte, bitmap []byte) error {
	// ATLAS(zgx): should save signature and bitmap into db, proposal is a sealed block.
	// Check if the proposal is a valid block
	block := &types.Block{}
	block, ok := proposal.(*types.Block)
	if !ok {
		sb.logger.Error("Invalid proposal, %v", proposal)
		return errInvalidProposal
	}

	h := block.Header()
	// Append seals into extra-data
	err := WriteCommittedSeals(h, signature, publicKey, bitmap)
	if err != nil {
		return err
	}
	// update block's header
	block = block.WithSeal(h)

	sb.logger.Info("Confirm", "address", sb.Address(), "hash", proposal.Hash(), "number", proposal.Number().Uint64())
	if err := sb.WriteLastCommits(signature, bitmap); err != nil {
		return err
	}

	// - if the proposed and committed blocks are the same, send the proposed hash
	//   to commit channel, which is being watched inside the engine.Seal() function.
	// - otherwise, we try to insert the block.
	// -- if success, the ChainHeadEvent event will be broadcasted, try to build
	//    the next block and the previous Seal() will be stopped.
	// -- otherwise, a error will be returned and a round change event will be fired.
	if sb.proposedBlockHash == block.Hash() {
		// feed block hash to Seal() and wait the Seal() result
		sb.commitCh <- block
		return nil
	}

	if sb.broadcaster != nil {
		sb.broadcaster.Enqueue(fetcherID, block)
	}
	return nil
}

// EventMux implements atlas.Backend.EventMux
func (sb *backend) EventMux() *event.TypeMux {
	return sb.atlasEventMux
}

// Verify implements atlas.Backend.Verify
func (sb *backend) Verify(proposal atlas.Proposal) (time.Duration, error) {
	// Check if the proposal is a valid block
	block := &types.Block{}
	block, ok := proposal.(*types.Block)
	if !ok {
		sb.logger.Error("Invalid proposal, %v", proposal)
		return 0, errInvalidProposal
	}

	// check bad block
	if sb.HasBadProposal(block.Hash()) {
		return 0, core.ErrBlacklistedHash
	}

	// check block body
	txnHash := types.DeriveSha(block.Transactions())
	uncleHash := types.CalcUncleHash(block.Uncles())
	if txnHash != block.Header().TxHash {
		return 0, errMismatchTxhashes
	}
	if uncleHash != nilUncleHash {
		return 0, errInvalidUncleHash
	}

	// verify the header of proposed block
	err := sb.VerifyHeader(sb.chain, block.Header(), false)
	// ignore errEmptyCommittedSeals error because we don't have the committed seals yet
	if err == nil || err == errEmptyCommittedSeals {
		return 0, nil
	} else if err == consensus.ErrFutureBlock {
		return time.Unix(int64(block.Header().Time), 0).Sub(now()), consensus.ErrFutureBlock
	}
	return 0, err
}

// Sign implements atlas.Backend.Sign
func (sb *backend) SignHash(hash common.Hash) (signature []byte, publicKey []byte, mask []byte, err error) {
	// ATLAS(zgx): Sign is called by finalizeMessage and updateBlock, the former sign message, the latter sign block
	signature, publicKey, mask, err = sb.signHashFn(accounts.Account{Address: sb.signer}, hash)
	if err != nil {
		return nil, nil, nil, err
	}
	return signature, publicKey, mask, nil
}

// CheckSignature implements atlas.Backend.CheckSignature
func (sb *backend) CheckSignature(data []byte, pubKey []byte, sig []byte) error {
	var publicKey bls.PublicKey
	if err := publicKey.Deserialize(pubKey); err != nil {
		return err
	}

	var sign bls.Sign
	if err := sign.Deserialize(sig); err != nil {
		return err
	}

	hashData := crypto.Keccak256([]byte(data))

	if ok := sign.VerifyHash(&publicKey, hashData); !ok {
		return errInvalidSignature
	}
	return nil
}

// HasPropsal implements atlas.Backend.HashBlock
func (sb *backend) HasPropsal(hash common.Hash, number *big.Int) bool {
	return sb.chain.GetHeader(hash, number.Uint64()) != nil
}

// GetProposer implements atlas.Backend.GetProposer
func (sb *backend) GetProposer(number uint64) common.Address {
	if h := sb.chain.GetHeaderByNumber(number); h != nil {
		a, _ := sb.Author(h)
		return a
	}
	return common.Address{}
}

// ParentValidators implements atlas.Backend.GetParentValidators
func (sb *backend) ParentValidators(proposal atlas.Proposal) atlas.ValidatorSet {
	if block, ok := proposal.(*types.Block); ok {
		return sb.getValidators(block.Number().Uint64()-1, block.ParentHash())
	}
	return validator.NewSet(nil, sb.config.ProposerPolicy)
}

func (sb *backend) getValidators(number uint64, hash common.Hash) atlas.ValidatorSet {
	snap, err := sb.snapshot(sb.chain, number, hash, nil)
	if err != nil {
		return validator.NewSet(nil, sb.config.ProposerPolicy)
	}
	return snap.ValSet
}

func (sb *backend) LastProposal() (atlas.Proposal, common.Address) {
	block := sb.currentBlock()

	var proposer common.Address
	if block.Number().Cmp(common.Big0) > 0 {
		var err error
		proposer, err = sb.Author(block.Header())
		if err != nil {
			sb.logger.Error("Failed to get block proposer", "err", err)
			return nil, common.Address{}
		}
	}

	// Return header only block here since we don't need block body
	return block, proposer
}

func (sb *backend) HasBadProposal(hash common.Hash) bool {
	if sb.hasBadBlock == nil {
		return false
	}
	return sb.hasBadBlock(hash)
}

func (sb *backend) Close() error {
	return nil
}

// Authorize injects a private key into the consensus engine to mint new blocks
// with.
func (c *backend) Authorize(signer common.Address, signHashFn consensus.SignHashFn) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.signer = signer
	c.signHashFn = signHashFn
}

func (sb *backend) WriteLastCommits(signature []byte, mask []byte) error {
	if len(signature) != types.AtlasExtraSignature || len(mask) != types.AtlasExtraMask {
		return types.ErrInvalidAtlasHeaderExtra
	}
	data := make([]byte, len(signature)+len(mask))
	copy(data, signature)
	copy(data[len(signature):], mask)
	if err := sb.db.Put(lastCommitsKey, data); err != nil {
		return err
	}
	return nil
}

func (sb *backend) ReadLastCommits() (signature []byte, mask []byte, err error) {
	var data []byte
	data, err = sb.db.Get(lastCommitsKey)
	if err != nil {
		return nil, nil, err
	}

	if len(data) != types.AtlasExtraSignature+types.AtlasExtraMask {
		return nil, nil, types.ErrInvalidAtlasHeaderExtra
	}

	return data[:types.AtlasExtraSignature], data[types.AtlasExtraSignature:], nil
}
