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
	"bytes"
	"errors"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/trie"

	"io"
	"math"

	"math/big"
	"time"

	"github.com/hyperion-hyn/bls/ffi/go/bls"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	bls_cosi "github.com/ethereum/go-ethereum/crypto/bls"

	lru "github.com/hashicorp/golang-lru"
	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/atlas"
	"github.com/ethereum/go-ethereum/consensus/atlas/validator"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	checkpointInterval = 1024 // Number of blocks after which to save the vote snapshot to the database
	inmemorySnapshots  = 128  // Number of recent vote snapshots to keep in memory
	inmemoryPeers      = 40
	inmemoryMessages   = 1024 * 10
	MaxValidatorCount  = 88
)

const (
	CONSORTIUM_BOARD = "0xa40bFc4701562c3fBe246E1da2Ac980c929b7d3e"
)

var (
	// errInvalidProposal is returned when a prposal is malformed.
	errInvalidProposal = errors.New("invalid proposal")
	// errInvalidSignature is returned when given signature is not signed by given
	// address.
	errInvalidSignature = errors.New("invalid signature")
	// errInvalidPublicKey is returned when given public key is valid
	errInvalidPublicKey = errors.New("invalid public key")
	// errUnknownBlock is returned when the list of validators is requested for a block
	// that is not part of the local blockchain.
	errUnknownBlock = errors.New("unknown block")
	// errUnauthorized is returned if a header is signed by a non authorized entity.
	errUnauthorized = errors.New("unauthorized")
	// errInvalidDifficulty is returned if the difficulty of a block is not 1
	errInvalidDifficulty = errors.New("invalid difficulty")
	// errInvalidExtraDataFormat is returned when the extra data format is incorrect
	errInvalidExtraDataFormat = errors.New("invalid extra data format")
	// errInvalidMixDigest is returned if a block's mix digest is not Atlas digest.
	errInvalidMixDigest = errors.New("invalid Atlas mix digest")
	// errInvalidNonce is returned if a block's nonce is invalid
	errInvalidNonce = errors.New("invalid nonce")
	// errInvalidUncleHash is returned if a block contains an non-empty uncle list.
	errInvalidUncleHash = errors.New("non empty uncle hash")
	// errInconsistentValidatorSet is returned if the validator set is inconsistent
	errInconsistentValidatorSet = errors.New("non empty uncle hash")
	// errInvalidTimestamp is returned if the timestamp of a block is lower than the previous block's timestamp + the minimum block period.
	errInvalidTimestamp = errors.New("invalid timestamp")
	// errInvalidVotingChain is returned if an authorization list is attempted to
	// be modified via out-of-range or non-contiguous headers.
	errInvalidVotingChain = errors.New("invalid voting chain")
	// errInvalidVote is returned if a nonce value is something else that the two
	// allowed constants of 0x00..0 or 0xff..f.
	errInvalidVote = errors.New("vote nonce not 0x00..0 or 0xff..f")
	// errInvalidCommittedSeals is returned if the committed seal is not signed by any of parent validators.
	errInvalidCommittedSeals = errors.New("invalid committed seals")
	// errEmptyCommittedSeals is returned if the field of committed seals is zero.
	errEmptyCommittedSeals = errors.New("zero committed seals")
	// errInvalidLastCommits is returned if LastCommits is invalid
	errInvalidLastCommits = errors.New("invalid lastCommits")
	// errInvalidAggregatedSignature is returned if the field of aggregated signature is invalid.
	errInvalidAggregatedSignature = errors.New("invalid aggregated signature")
	// errMismatchTxhashes is returned if the TxHash in header is mismatch.
	errMismatchTxhashes = errors.New("mismatch transcations hashes")
	// errCountBetweenPublicKeyAndSignatureNotMatch is returned if count of public keys and count of signature is mismatch
	errCountBetweenPublicKeyAndSignatureNotMatch = errors.New("Count between public key and signature not match.")
)
var (
	DefaultDifficulty = big.NewInt(1)
	nilUncleHash      = types.CalcUncleHash(nil) // Always Keccak256(RLP([])) as uncles are meaningless outside of PoW.
	emptyNonce        = types.BlockNonce{}
	now               = time.Now

	nonceAuthVote = hexutil.MustDecode("0xffffffffffffffff") // Magic nonce number to vote on adding a new validator
	nonceDropVote = hexutil.MustDecode("0x0000000000000000") // Magic nonce number to vote on removing a validator.

	inmemoryAddresses  = 20 // Number of recent addresses from ecrecover
	recentAddresses, _ = lru.NewARC(inmemoryAddresses)
)

// Author retrieves the Ethereum address of the account that minted the given
// block, which may be different from the header's coinbase if a consensus
// engine is based on signatures.
func (sb *backend) Author(header *types.Header) (common.Address, error) {
	return common.Address{}, nil
}

// Signers extracts all the addresses who have signed the given header
// It will extract for each seal who signed it, regardless of if the seal is
// repeated
func (sb *backend) Signers(header *types.Header) ([]atlas.Validator, error) {
	number := header.Number.Uint64()
	snap, err := sb.snapshot(sb.chain, number-1, header.ParentHash, nil)
	if err != nil {
		return []atlas.Validator{}, err
	}

	signers, err := getSigners(snap.ValSet, header.LastCommits[types.AtlasExtraSignature:])
	if err != nil {
		return nil, err
	}

	return signers, nil
}

func getSigners(valSet atlas.ValidatorSet, mask []byte) ([]atlas.Validator, error) {
	bitmap, _ := bls_cosi.NewMask(valSet.GetPublicKeys(), nil)
	if err := bitmap.SetMask(mask); err != nil {
		return nil, err
	}

	signers := make([]atlas.Validator, bitmap.CountEnabled())
	publicKeys := bitmap.GetPubKeyFromMask(true)
	for _, publicKey := range publicKeys {
		_, validator := valSet.GetByPublicKey(publicKey)
		signers = append(signers, validator)
	}
	return signers, nil
}

// VerifyHeader checks whether a header conforms to the consensus rules of a
// given engine. Verifying the seal may be done optionally here, or explicitly
// via the VerifySeal method.
func (sb *backend) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, seal bool) error {
	return sb._VerifyHeader(chain.(consensus.ChainReader), header, seal)
}

func (sb *backend) _VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	return sb.verifyHeader(chain, header, nil, seal)
}

// verifyHeader checks whether a header conforms to the consensus rules.The
// caller may optionally pass in a batch of parents (ascending order) to avoid
// looking those up from the database. This is useful for concurrently verifying
// a batch of new headers.
func (sb *backend) verifyHeader(chain consensus.ChainReader, header *types.Header, parents []*types.Header, seal bool) error {
	if header.Number == nil {
		return errUnknownBlock
	}

	// Don't waste time checking blocks from the future
	if header.Time > uint64(now().Unix()) {
		return consensus.ErrFutureBlock
	}

	// Ensure that the coinbase is valid
	if header.Nonce != (emptyNonce) {
		return errInvalidNonce
	}
	// Ensure that the mix digest is zero as we don't have fork protection currently
	if header.MixDigest != types.AtlasDigest {
		return errInvalidMixDigest
	}
	// Ensure that the block doesn't contain any uncles which are meaningless in Atlas
	if header.UncleHash != nilUncleHash {
		return errInvalidUncleHash
	}
	// Ensure that the block's difficulty is meaningful (may not be correct at this point)
	if header.Difficulty == nil || header.Difficulty.Cmp(DefaultDifficulty) != 0 {
		return errInvalidDifficulty
	}

	return sb.verifyCascadingFields(chain, header, parents, seal)
}

// verifyCascadingFields verifies all the header fields that are not standalone,
// rather depend on a batch of previous headers. The caller may optionally pass
// in a batch of parents (ascending order) to avoid looking those up from the
// database. This is useful for concurrently verifying a batch of new headers.
func (sb *backend) verifyCascadingFields(chain consensus.ChainReader, header *types.Header, parents []*types.Header, seal bool) error {
	// The genesis block is the always valid dead-end
	number := header.Number.Uint64()
	if number == 0 {
		return nil
	}
	// Ensure that the block's timestamp isn't too close to it's parent
	var parent *types.Header
	if len(parents) > 0 {
		parent = parents[len(parents)-1]
	} else {
		parent = chain.GetHeader(header.ParentHash, number-1)
	}
	if parent == nil || parent.Number.Uint64() != number-1 || parent.Hash() != header.ParentHash {
		return consensus.ErrUnknownAncestor
	}
	if parent.Time+sb.config.BlockPeriod > header.Time {
		return errInvalidTimestamp
	}

	if err := sb.verifySigner(chain, header, parents); err != nil {
		return err
	}

	if seal {
		return sb.verifyCommittedSeals(chain, header, parents)
	}
	return nil
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications (the order is that of
// the input slice).
func (sb *backend) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	return sb._VerifyHeaders(chain.(consensus.ChainReader), headers, seals)
}

func (sb *backend) _VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
	results := make(chan error, len(headers))
	go func() {
		for i, header := range headers {
			err := sb.verifyHeader(chain, header, headers[:i], seals[i])

			select {
			case <-abort:
				return
			case results <- err:
			}
		}
	}()
	return abort, results
}

// VerifyUncles verifies that the given block's uncles conform to the consensus
// rules of a given engine.
func (sb *backend) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	if len(block.Uncles()) > 0 {
		return errInvalidUncleHash
	}
	return nil
}

// verifySigner checks whether the signer is in parent's validator set
func (sb *backend) verifySigner(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	// Verifying the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return errUnknownBlock
	}

	return nil
}

func (sb *backend) VerifyCommittedSeals(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	return sb.verifyCommittedSeals(chain, header, parents)
}

// verifyCommittedSeals checks whether every committed seal is signed by one of the parent's validators
func (sb *backend) verifyCommittedSeals(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	number := header.Number.Uint64()
	// We don't need to verify committed seals in the genesis block
	if number == 0 {
		return nil
	}

	// Retrieve the snapshot needed to verify this header and cache it
	var (
		snap *Snapshot
		err  error
	)
	if number == 1 {
		snap, err = sb.snapshot(chain, 0, header.ParentHash, parents)
	} else {
		parentHeader := chain.GetHeaderByHash(header.ParentHash)
		snap, err = sb.snapshot(chain, number-2, parentHeader.ParentHash, parents)
	}
	if err != nil {
		return err
	}

	if len(header.LastCommits) != types.AtlasExtraSignature+types.GetMaskByteCount(snap.ValSet.Size()) {
		return errInvalidAggregatedSignature
	}

	var parent *types.Header
	if len(parents) > 0 {
		parent = parents[len(parents)-1]
	} else {
		parent = chain.GetHeader(header.ParentHash, number-1)
	}

	err = verifySignature(snap.ValSet, parent.Hash().Bytes()[:], header.LastCommits[:types.AtlasExtraSignature], header.LastCommits[types.AtlasExtraSignature:])
	if err != nil {
		return err
	}

	return nil
}

func verifySignature(valSet atlas.ValidatorSet, hash []byte, signature []byte, bitmap []byte) error {
	mask, err := bls_cosi.NewMask(valSet.GetPublicKeys(), nil)
	if err != nil {
		return err
	}
	if err := mask.SetMask(bitmap); err != nil {
		return err
	}

	quorumSize := int(math.Ceil(float64(2*valSet.Size()) / 3))
	if mask.CountEnabled() < quorumSize {
		return errInvalidCommittedSeals
	}

	aggregatePublicKey := mask.AggregatePublic

	var sign bls.Sign
	if err := sign.Deserialize(signature); err != nil {
		return err
	}

	if ok := sign.VerifyHash(aggregatePublicKey, hash); !ok {
		log.Error("verify hash error", "hash", common.Bytes2Hex(hash))
		return errInvalidAggregatedSignature
	}

	return nil
}

// VerifySeal checks whether the crypto seal on a header is valid according to
// the consensus rules of the given engine.
func (sb *backend) VerifySeal(chain consensus.ChainHeaderReader, header *types.Header) error {
	return sb._VerifySeal(chain.(consensus.ChainReader), header)
}

func (sb *backend) _VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	// get parent header and ensure the signer is in parent's validator set
	number := header.Number.Uint64()
	if number == 0 {
		return errUnknownBlock
	}

	// ensure that the difficulty equals to DefaultDifficulty
	if header.Difficulty.Cmp(DefaultDifficulty) != 0 {
		return errInvalidDifficulty
	}
	if err := sb.verifySigner(chain, header, nil); err != nil {
		return err
	}

	if err := sb.verifyCommittedSeals(chain, header, nil); err != nil {
		return err
	}
	return nil
}

// Prepare initializes the consensus fields of a block header according to the
// rules of a particular engine. The changes are executed inline.
func (sb *backend) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	return sb._Prepare(chain.(consensus.ChainReader), header)
}

func (sb *backend) _Prepare(chain consensus.ChainReader, header *types.Header) error {
	// unused fields, force to set to empty
	header.Coinbase = common.Address{}
	header.Nonce = emptyNonce
	header.MixDigest = types.AtlasDigest

	number := header.Number.Uint64()
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	// use the same difficulty for all blocks
	header.Difficulty = DefaultDifficulty

	var (
		snap *Snapshot
		err  error
	)
	if number == 1 {
		snap, err = sb.snapshot(chain, number-1, header.ParentHash, nil)
	} else {
		// Assemble the voting snapshot
		snap, err = sb.snapshot(chain, number-2, parent.ParentHash, nil)
	}
	if err != nil {
		return err
	}

	lastCommits, err := rawdb.ReadLastCommits(chain.Database(), number-1)
	if err != nil {
		sb.logger.Error("last commit not found. ", "number", number-1)
		return errInvalidLastCommits
	}
	if len(lastCommits) != types.AtlasExtraSignature+types.GetMaskByteCount(snap.ValSet.Size()) {
		sb.logger.Error("last commit length error.", "signature", types.AtlasExtraSignature, "maskCount", types.GetMaskByteCount(snap.ValSet.Size()))
		return errInvalidLastCommits
	}

	// set header's signature and bitmap
	header.LastCommits = make([]byte, len(lastCommits))
	copy(header.LastCommits[:], lastCommits[:])

	// set header's slashes
	header.Slashes = []byte{}

	extra, err := prepareExtra(header, snap.validators())
	if err != nil {
		return err
	}
	header.Extra = extra

	// set header's timestamp
	header.Time = parent.Time + sb.config.BlockPeriod
	if header.Time < uint64(time.Now().Unix()) {
		header.Time = uint64(time.Now().Unix())
	}
	return nil
}

// Finalize runs any post-transaction state modifications (e.g. block rewards)
// and assembles the final block.
//
// Note, the block header and state database might be updated to reflect any
// consensus rules that happen at finalization (e.g. block rewards).
func (sb *backend) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header) {
	sb._Finalize(chain.(consensus.ChainReader), header, state, txs, uncles)
}

func (sb *backend) _Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header) {
	chainReader := chain.(consensus.ChainReader)                    // ATLAS
	_, err := handleMap3AndAtlasStaking(chainReader, header, state) // ATLAS
	if err != nil {
		sb.logger.Error("staking err", "err", err)
	}

	// No block rewards in Atlas, so the state remains as is and uncles are dropped
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = nilUncleHash
}

// Finalize runs any post-transaction state modifications (e.g. block rewards)
// and assembles the final block.
//
// Note, the block header and state database might be updated to reflect any
// consensus rules that happen at finalization (e.g. block rewards).
func (sb *backend) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	return sb._FinalizeAndAssemble(chain.(consensus.ChainReader), header, state, txs, uncles, receipts)
}

func (sb *backend) _FinalizeAndAssemble(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	// ATLAS
	chainReader := chain.(consensus.ChainReader)
	_, err := handleMap3AndAtlasStaking(chainReader, header, state)
	if err != nil {
		sb.logger.Error("handleMap3AndAtlasStaking", "err", err)
		return nil, err
	}
	// ATLAS - END

	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = nilUncleHash

	// Assemble and return the final block for sealing
	return types.NewBlock(header, txs, nil, receipts, new(trie.Trie)), nil
}

// Seal generates a new block for the given input block with the local miner's
// seal place on top.
func (sb *backend) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	return sb._Seal(chain.(consensus.ChainReader), block, results, stop)
}

func (sb *backend) _Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {

	// update the block header timestamp and signature and propose the block to core engine
	header := block.Header()
	number := header.Number.Uint64()

	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}

	// Bail out if we're unauthorized to sign a block
	snap, err := sb.snapshot(chain, number-1, header.ParentHash, nil)
	if err != nil {
		return err
	}

	containSigner := false
	for _, signer := range sb.signers {
		if _, v := snap.ValSet.GetBySigner(signer); v != nil {
			containSigner = true
			break
		}
	}

	if !containSigner {
		return errUnauthorized
	}

	block, err = sb.updateBlock(parent, block)
	if err != nil {
		return err
	}

	go func() {
		// get the proposed block hash and clear it if the seal() is completed.
		sb.sealMu.Lock()
		sb.proposedBlockHash = block.Hash()

		defer func() {
			sb.proposedBlockHash = common.Hash{}
			sb.sealMu.Unlock()
		}()
		// post block into Atlas engine
		go sb.EventMux().Post(atlas.RequestEvent{
			Proposal: block,
		})
		for {
			select {
			case result := <-sb.commitCh:
				// if the block hash and the hash from channel are the same,
				// return the result. Otherwise, keep waiting the next hash.
				// ATLAS:
				//  *WARNING*: when a new block was received, worker will commit new work,
				//  which will interrupt previous Seal and begin a new Seal.
				//  if only consume ONE block here and quit, COMMIT messages will pile up and
				//  wait on sending block to commit channel, and hang in handleEvents. then
				//  no more events will be handled.
				//  to make it work and prevent hang in handleEvents, keeping consumming blocks
				//  from commit channel here is necessary.
				if result != nil && block.Hash() == result.Hash() {
					// wait for the timestamp of header, use this to adjust the block period
					delay := time.Unix(int64(header.Time+sb.config.BlockPeriod), 0).Sub(now())
					sb.logger.Debug("mine new block in future", "delay", delay)
					select {
					case <-time.After(delay):
					}
					results <- result
					return
				} else if result != nil {
					sb.logger.Debug("drop block", "number", result.NumberU64(), "blockHash", sb.SealHash(block.Header()), "resultHash", sb.SealHash(result.Header()))
					// ATLAS: keeping consumming blocks from commit channel is necessary
					//  to prevent hang in handleEvents.
				} else {
					sb.logger.Debug("result is null")
				}
			case <-stop:
				sb.logger.Debug("stop seal", "number", block.NumberU64())
				results <- nil
				return
			}
		}
	}()
	return nil
}

// update timestamp and signature of the block based on its number of transactions
func (sb *backend) updateBlock(parent *types.Header, block *types.Block) (*types.Block, error) {
	header := block.Header()

	return block.WithSeal(header), nil
}

// APIs returns the RPC APIs this consensus engine provides.
func (sb *backend) APIs(chain consensus.ChainHeaderReader) []rpc.API {
	return []rpc.API{{
		Namespace: "atlas",
		Version:   "1.0",
		Service:   &API{chain: chain.(consensus.ChainReader), atlas: sb},
		Public:    true,
	}}
}

// Start implements consensus.Atlas.Start
func (sb *backend) Start(chain consensus.ChainReader, currentBlock func() *types.Block, hasBadBlock func(hash common.Hash) bool) error {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()
	if sb.coreStarted {
		return atlas.ErrStartedEngine
	}

	// clear previous data
	sb.proposedBlockHash = common.Hash{}
	if sb.commitCh != nil {
		close(sb.commitCh)
	}
	sb.commitCh = make(chan *types.Block, 1)

	sb.chain = chain
	sb.currentBlock = currentBlock
	sb.hasBadBlock = hasBadBlock

	if err := sb.core.Start(); err != nil {
		return err
	}

	sb.coreStarted = true
	return nil
}

// Stop implements consensus.Atlas.Stop
func (sb *backend) Stop() error {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()
	if !sb.coreStarted {
		return atlas.ErrStoppedEngine
	}
	if err := sb.core.Stop(); err != nil {
		return err
	}
	sb.coreStarted = false
	return nil
}

// snapshot retrieves the authorization snapshot at a given point in time.
func (sb *backend) snapshot(chain consensus.ChainReader, number uint64, hash common.Hash, parents []*types.Header) (*Snapshot, error) {
	if number == 0 || number == 1 {
		hash = chain.GetHeaderByNumber(0).Hash()
	}

	// If an in-memory snapshot was found, use that
	if s, ok := sb.recents.Get(hash); ok {
		snap := s.(*Snapshot)
		if number > 1 && snap.Number != number {
			return nil, consensus.ErrUnknownAncestor
		}
		return snap, nil
	}
	header := chain.GetHeaderByNumber(number)
	if header == nil {
		return nil, consensus.ErrUnknownAncestor
	}
	stateDB, err := chain.StateAt(header.Root)
	if err != nil {
		return nil, err
	}
	validators, err := getValidators(stateDB, MaxValidatorCount)
	log.Debug("get validator", "number", number, "size", len(validators))
	if err != nil {
		return nil, err
	}
	snap := newSnapshot(sb.config.Epoch, number, hash, validator.NewSet(validators, sb.config.ProposerPolicy))
	sb.recents.Add(hash, snap)
	return snap, nil
}

// SealHash returns the hash of a block prior to it being sealed.
func (sb *backend) SealHash(header *types.Header) common.Hash {
	return SealHash(header)
}

// prepareExtra returns a extra-data of the given header and validators
func prepareExtra(header *types.Header, vals []atlas.Validator) ([]byte, error) {
	var buf bytes.Buffer

	buf.Write(header.Extra[:])

	return append(buf.Bytes(), []byte{}...), nil
}

// writeSeal writes the extra-data field of the given header with the given seals.
// suggest to rename to writeSeal.
func writeSeal(h *types.Header, seal []byte) error {
	if len(seal) != types.AtlasExtraSignature {
		return errInvalidSignature
	}

	copy(h.LastCommits[:], seal)

	return nil
}

// WriteCommittedSeals writes the extra-data field of a block header with given committed seals.
func WriteCommittedSeals(h *types.Header, signature []byte, bitmap []byte, valSetSize int) error {
	if len(signature) != types.AtlasExtraSignature || len(bitmap) != types.GetMaskByteCount(valSetSize) {
		return errInvalidCommittedSeals
	}

	copy(h.LastCommits[:types.AtlasExtraSignature], signature[:])
	copy(h.LastCommits[types.AtlasExtraSignature:], bitmap[:])

	return nil
}

func WriteCommittedSealInGenesis(genesis *core.Genesis, header *types.Header, signatures []*bls.Sign, publicKeys []*bls.PublicKey) error {
	mask, err := bls_cosi.NewMask(publicKeys, nil)
	if err != nil {
		return err
	}

	if len(publicKeys) != len(signatures) {
		return errCountBetweenPublicKeyAndSignatureNotMatch
	}

	var sign bls.Sign
	var publicKey bls.PublicKey

	for i := 0; i < len(publicKeys); i++ {
		mask.SetKey(publicKeys[i], true)
		publicKey.Add(publicKeys[i])
		sign.Add(signatures[i])
	}

	genesis.LastCommits = make([]byte, types.AtlasExtraSignature+types.GetMaskByteCount(len(publicKeys)))
	copy(genesis.LastCommits[:types.AtlasExtraSignature], sign.Serialize())
	copy(genesis.LastCommits[types.AtlasExtraSignature:], mask.Mask())

	return nil
}

// ATLAS(yhx): getValidators
func getValidators(state *state.StateDB, numVal int) ([]atlas.Validator, error) {
	committee, err := state.ValidatorPool().Committee().Load()
	if err != nil {
		return nil, err
	}

	length := len(committee.Slots.Entrys)
	validators := make([]atlas.Validator, length)
	for i := 0; i < length; i++ {
		member := committee.Slots.Entrys[i]
		v, err := validator.New(member.BLSPublicKey.Key[:], member.EcdsaAddress)
		if err != nil {
			return nil, err
		}
		validators[i] = v
	}
	return validators, nil
}

// SealHash returns the hash of a block prior to it being sealed.
func SealHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewLegacyKeccak256()
	encodeSigHeader(hasher, header)
	hasher.Sum(hash[:0])
	return hash
}

// AtlasRLP returns the rlp bytes which needs to be signed for
// sealing. The RLP to sign consists of the entire header apart from the signature
// contained at the end of the extra data.
//
// Note, the method requires the extra data to be at least 65 bytes, otherwise it
// panics. This is done to avoid accidentally using both forms (signature present
// or not), which could be abused to produce different hashes for the same header.
func AtlasRLP(header *types.Header) []byte {
	b := new(bytes.Buffer)
	encodeSigHeader(b, header)
	return b.Bytes()
}

func encodeSigHeader(w io.Writer, header *types.Header) {
	err := rlp.Encode(w, []interface{}{
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra,
		header.MixDigest,
		header.Nonce,
		header.Epoch,
		header.LastCommits,
	})
	if err != nil {
		panic("can't encode: " + err.Error())
	}
}

func retrieveAggregatedPublicKey(valSet atlas.ValidatorSet, bitmap []byte) (*bls.PublicKey, error) {
	mask, err := bls_cosi.NewMask(valSet.GetPublicKeys(), nil)
	if err != nil {
		return nil, err
	}
	if err := mask.SetMask(bitmap); err != nil {
		return nil, err
	}
	return mask.AggregatePublic, nil
}
