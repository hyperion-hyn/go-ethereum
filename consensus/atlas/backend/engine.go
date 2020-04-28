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
	"math"

	"github.com/docker/docker/daemon/logger"

	"github.com/harmony-one/bls/ffi/go/bls"
	bls_cosi "github.com/ethereum/go-ethereum/crypto/bls"
	"github.com/ethereum/go-ethereum/staking"
	"math/big"
	"math/rand"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/atlas"
	atlasCore "github.com/ethereum/go-ethereum/consensus/atlas/core"
	"github.com/ethereum/go-ethereum/consensus/atlas/validator"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	lru "github.com/hashicorp/golang-lru"
	"golang.org/x/crypto/sha3"
)

const (
	checkpointInterval = 1024 // Number of blocks after which to save the vote snapshot to the database
	inmemorySnapshots  = 128  // Number of recent vote snapshots to keep in memory
	inmemoryPeers      = 40
	inmemoryMessages   = 1024
	MaxValidatorCount  = 88
)

var (
	// errInvalidProposal is returned when a prposal is malformed.
	errInvalidProposal = errors.New("invalid proposal")
	// errInvalidSignature is returned when given signature is not signed by given
	// address.
	errInvalidSignature = errors.New("invalid signature")
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
	// errMismatchTxhashes is returned if the TxHash in header is mismatch.
	errMismatchTxhashes = errors.New("mismatch transcations hashes")
	// errMismatchTxhashes is returned if the TxHash in header is mismatch.
	err = errors.New("mismatch transcations hashes")
)
var (
	defaultDifficulty = big.NewInt(1)
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
	number := header.Number.Uint64()
	snap, err := sb.snapshot(sb.chain, number-1, header.ParentHash, nil)
	if err != nil {
		return common.Address{}, err
	}

	return ecrecover(snap, header)
}

// Signers extracts all the addresses who have signed the given header
// It will extract for each seal who signed it, regardless of if the seal is
// repeated
func (sb *backend) Signers(header *types.Header) ([], error) {
	// ATLAS(yhx): return signers who signed the given header, should return public keys
	extra, err := types.ExtractAtlasExtra(header)
	if err != nil {
		return []common.Address{}, err
	}

	number := header.Number.Uint64()
	snap, err := sb.snapshot(sb.chain, number-1, header.ParentHash, nil)
	if err != nil {
		return []common.Address{}, err
	}

	var sign bls.Sign
	if err := sign.Deserialize(extra.Signature); err != nil {
		return []common.Address{}, err
	}

	var pubKey bls.PublicKey
	if err := pubKey.Deserialize(extra.PublicKey); err != nil {
		return []common.Address{}, err
	}

	// ATLAS(zgx): should we verify signature here?
	// if sign.Verify(&pubKey, header.Number) == false {
	// 	logger.Error("Failed to verify signature with signer's public key")
	// 	return errInvalidSignature
	// }

	bitmap, _ := bls_cosi.NewMask(snap.ValSet.GetPublicKeys(), nil)
	if err := bitmap.SetMask(extra.AggBitmap); err != nil {
		return []common.Address{}, err
	}

	signers := make([]atlas.Validator)
	publicKeys := bitmap.GetPubKeyFromMask(true)
	for i, publicKey := range(publicKeys) {
		idx, validator := snap.ValSet.GetByPublicKey(publicKey)
		signers = append(signers, validator)
	}

	// ATLAS(zgx): match return values.
	return signers, nil
}

// VerifyHeader checks whether a header conforms to the consensus rules of a
// given engine. Verifying the seal may be done optionally here, or explicitly
// via the VerifySeal method.
func (sb *backend) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	return sb.verifyHeader(chain, header, nil)
}

// verifyHeader checks whether a header conforms to the consensus rules.The
// caller may optionally pass in a batch of parents (ascending order) to avoid
// looking those up from the database. This is useful for concurrently verifying
// a batch of new headers.
func (sb *backend) verifyHeader(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	if header.Number == nil {
		return errUnknownBlock
	}

	// Don't waste time checking blocks from the future
	if header.Time > uint64(now().Unix()) {
		return consensus.ErrFutureBlock
	}

	// Ensure that the extra data format is satisfied
	if _, err := types.ExtractAtlasExtra(header); err != nil {
		return errInvalidExtraDataFormat
	}

	// Ensure that the coinbase is valid
	if header.Nonce != (emptyNonce) && !bytes.Equal(header.Nonce[:], nonceAuthVote) && !bytes.Equal(header.Nonce[:], nonceDropVote) {
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
	if header.Difficulty == nil || header.Difficulty.Cmp(defaultDifficulty) != 0 {
		return errInvalidDifficulty
	}

	return sb.verifyCascadingFields(chain, header, parents)
}

// verifyCascadingFields verifies all the header fields that are not standalone,
// rather depend on a batch of previous headers. The caller may optionally pass
// in a batch of parents (ascending order) to avoid looking those up from the
// database. This is useful for concurrently verifying a batch of new headers.
func (sb *backend) verifyCascadingFields(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
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
	// Verify validators in extraData. Validators in snapshot and extraData should be the same.
	snap, err := sb.snapshot(chain, number-1, header.ParentHash, parents)
	if err != nil {
		return err
	}
	validators := make([]byte, len(snap.validators())*common.AddressLength)
	for i, validator := range snap.validators() {
		copy(validators[i*common.AddressLength:], validator[:])
	}
	if err := sb.verifySigner(chain, header, parents); err != nil {
		return err
	}

	return sb.verifyCommittedSeals(chain, header, parents)
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications (the order is that of
// the input slice).
func (sb *backend) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
	results := make(chan error, len(headers))
	go func() {
		for i, header := range headers {
			err := sb.verifyHeader(chain, header, headers[:i])

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

	// Retrieve the snapshot needed to verify this header and cache it
	snap, err := sb.snapshot(chain, number-1, header.ParentHash, parents)
	if err != nil {
		return err
	}

	// resolve the authorization key and check against signers
	signer, err := ecrecover(snap, header)
	if err != nil {
		return err
	}

	// Signer should be in the validator set of previous block's extraData.
	if _, v := snap.ValSet.GetByAddress(signer); v == nil {
		return errUnauthorized
	}
	return nil
}

// verifyCommittedSeals checks whether every committed seal is signed by one of the parent's validators
func (sb *backend) verifyCommittedSeals(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	number := header.Number.Uint64()
	// We don't need to verify committed seals in the genesis block
	if number == 0 {
		return nil
	}

	// Retrieve the snapshot needed to verify this header and cache it
	snap, err := sb.snapshot(chain, number-1, header.ParentHash, parents)
	if err != nil {
		return err
	}

	extra, err := types.ExtractAtlasExtra(header)
	if err != nil {
		return err
	}
	// The length of Confirm seals should be larger than 0
	if len(extra.AggSignature) == 0 || len(extra.AggBitmap) == 0{
		return errEmptyCommittedSeals
	}

	validators := snap.ValSet.Copy()

	mask, err := bls_cosi.NewMask(snap.ValSet.GetPublicKeys(), nil)
	if err != nil {
		return err
	}

	quorumSize := int(math.Ceil(float64(2*validators.Size()) / 3))
	if mask.CountEnabled() < quorumSize {
		return errInvalidCommittedSeals
	}

	// ATLAS(zgx): maybe we should make sure signer satisify a low limit.
	return nil
}

// VerifySeal checks whether the crypto seal on a header is valid according to
// the consensus rules of the given engine.
func (sb *backend) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	// get parent header and ensure the signer is in parent's validator set
	number := header.Number.Uint64()
	if number == 0 {
		return errUnknownBlock
	}

	// ensure that the difficulty equals to defaultDifficulty
	if header.Difficulty.Cmp(defaultDifficulty) != 0 {
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
func (sb *backend) Prepare(chain consensus.ChainReader, header *types.Header) error {
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
	header.Difficulty = defaultDifficulty

	// set header's timestamp
	header.Time = parent.Time + sb.config.BlockPeriod
	if header.Time < uint64(time.Now().Unix()) {
		header.Time = uint64(time.Now().Unix())
	}

	signature, mask, err := sb.ReadLastCommits()
	if err != nil {
		return err
	}

	if err = writeCommittedSeals(header, signature, mask); err != nil {
		return err
	}

	return nil
}

// Finalize runs any post-transaction state modifications (e.g. block rewards)
// and assembles the final block.
//
// Note, the block header and state database might be updated to reflect any
// consensus rules that happen at finalization (e.g. block rewards).
func (sb *backend) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header) {
	// ATLAS(yhx): reward
	// No block rewards in Atlas, so the state remains as is and uncles are dropped
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = nilUncleHash
}

// Finalize runs any post-transaction state modifications (e.g. block rewards)
// and assembles the final block.
//
// Note, the block header and state database might be updated to reflect any
// consensus rules that happen at finalization (e.g. block rewards).
func (sb *backend) FinalizeAndAssemble(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	// No block rewards in Atlas, so the state remains as is and uncles are dropped
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = nilUncleHash

	// Assemble and return the final block for sealing
	return types.NewBlock(header, txs, nil, receipts), nil
}

// Seal generates a new block for the given input block with the local miner's
// seal place on top.
func (sb *backend) Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {

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

	if _, v := snap.ValSet.GetByAddress(sb.address); v == nil {
		return errUnauthorized
	}

	block, err = sb.updateBlock(parent, block)
	if err != nil {
		return err
	}

	go func() {
		// get the proposed block hash and clear it if the seal() is completed.
		sb.sealMu.Lock()
		// ATLAS(zgx): what is the purpose of proposedBlockHash, can we ignore it in FBFT?
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
				if result != nil && block.Hash() == result.Hash() {
					// wait for the timestamp of header, use this to adjust the block period
					delay := time.Unix(int64(header.Time), 0).Sub(now())
					// ATLAS(zgx): what if delay is negative?
					select {
					case <-time.After(delay):
					}
					results <- result
					return
				}
			case <-stop:
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
	// sign the hash
	seal, pubkey, err := sb.Sign(sigHash(header).Bytes())
	if err != nil {
		return nil, err
	}

	err = writeSeal(header, seal, pubkey)
	if err != nil {
		return nil, err
	}

	return block.WithSeal(header), nil
}

// APIs returns the RPC APIs this consensus engine provides.
func (sb *backend) APIs(chain consensus.ChainReader) []rpc.API {
	return []rpc.API{{
		Namespace: "atlas",
		Version:   "1.0",
		Service:   &API{chain: chain, atlas: sb},
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
	// Search for a snapshot in memory or on disk for checkpoints
	var (
		headers []*types.Header
		snap    *Snapshot
	)
	for snap == nil {
		// If an in-memory snapshot was found, use that
		if s, ok := sb.recents.Get(hash); ok {
			snap = s.(*Snapshot)
			break
		}
		// If an on-disk checkpoint snapshot can be found, use that
		if number%checkpointInterval == 0 {
			if s, err := loadSnapshot(sb.config.Epoch, sb.db, hash); err == nil {
				log.Trace("Loaded voting snapshot form disk", "number", number, "hash", hash)
				snap = s
				break
			}
		}
		// If we're at block zero, make a snapshot
		if number == 0 {
			genesis := chain.GetHeaderByNumber(0)
			if err := sb.VerifyHeader(chain, genesis, false); err != nil {
				return nil, err
			}
			// ATLAS(yhx): get validators from storage
			stateDB, err := chain.StateAt(chain.GetBlock(hash, number).Root())
			if err != nil {
				return nil, err
			}
			validators, err := getLargestAmountStakingValidators(stateDB, MaxValidatorCount)
			if err != nil {
				return nil, err
			}
			snap = newSnapshot(sb.config.Epoch, 0, genesis.Hash(), validator.NewSet(validators, sb.config.ProposerPolicy))
			if err := snap.store(sb.db); err != nil {
				return nil, err
			}
			log.Trace("Stored genesis voting snapshot to disk")
			break
		}
		// No snapshot for this header, gather the header and move backward
		var header *types.Header
		if len(parents) > 0 {
			// If we have explicit parents, pick from there (enforced)
			header = parents[len(parents)-1]
			if header.Hash() != hash || header.Number.Uint64() != number {
				return nil, consensus.ErrUnknownAncestor
			}
			parents = parents[:len(parents)-1]
		} else {
			// No explicit parents (or no more left), reach out to the database
			header = chain.GetHeader(hash, number)
			if header == nil {
				return nil, consensus.ErrUnknownAncestor
			}
		}
		headers = append(headers, header)
		number, hash = number-1, header.ParentHash
	}
	// Previous snapshot found, apply any pending headers on top of it
	for i := 0; i < len(headers)/2; i++ {
		headers[i], headers[len(headers)-1-i] = headers[len(headers)-1-i], headers[i]
	}
	snap, err := snap.apply(headers)
	if err != nil {
		return nil, err
	}
	sb.recents.Add(snap.Hash, snap)

	// If we've generated a new checkpoint snapshot, save to disk
	if snap.Number%checkpointInterval == 0 && len(headers) > 0 {
		if err = snap.store(sb.db); err != nil {
			return nil, err
		}
		log.Trace("Stored voting snapshot to disk", "number", snap.Number, "hash", snap.Hash)
	}
	return snap, err
}

// FIXME: Need to update this for Atlas
// sigHash returns the hash which is used as input for the Atlas
// signing. It is the hash of the entire header apart from the 65 byte signature
// contained at the end of the extra data.
//
// Note, the method requires the extra data to be at least 65 bytes, otherwise it
// panics. This is done to avoid accidentally using both forms (signature present
// or not), which could be abused to produce different hashes for the same header.
func sigHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewLegacyKeccak256()

	// Clean seal is required for calculating proposer seal.
	rlp.Encode(hasher, types.AtlasFilteredHeader(header, false))
	hasher.Sum(hash[:0])
	return hash
}

// SealHash returns the hash of a block prior to it being sealed.
func (sb *backend) SealHash(header *types.Header) common.Hash {
	return sigHash(header)
}

// ecrecover extracts the Ethereum account address from a signed header.
func ecrecover(snap * Snapshot, header *types.Header) (common.Address, error) {
	hash := header.Hash()
	if addr, ok := recentAddresses.Get(hash); ok {
		return addr.(common.Address), nil
	}

	// Retrieve the signature from the header extra-data
	atlasExtra, err := types.ExtractAtlasExtra(header)
	if err != nil {
		return common.Address{}, err
	}

	var sign bls.Sign
	if err := sign.Deserialize(atlasExtra.Signature); err != nil {
		return common.Address{}, err
	}

	var pubKey bls.PublicKey
	if err := pubKey.Deserialize(atlasExtra.PublicKey); err != nil {
		return common.Address{}, err
	}

	// ATLAS(zgx): should verify signature here?


	_, validator := snap.ValSet.GetByPublicKey(pubKey)
	if validator == nil {
		return common.Address{}, errValidatorNotExist
	}

	addr := validator.Coinbase()
	recentAddresses.Add(hash, addr)
	return addr, nil
}

// prepareExtra returns a extra-data of the given header and validators
func prepareExtra(header *types.Header, vals []common.Address) ([]byte, error) {
	var buf bytes.Buffer

	// compensate the lack bytes if header.Extra is not enough AtlasExtraVanity bytes.
	if len(header.Extra) < types.AtlasExtraVanity {
		header.Extra = append(header.Extra, bytes.Repeat([]byte{0x00}, types.AtlasExtraVanity-len(header.Extra))...)
	}
	buf.Write(header.Extra[:types.AtlasExtraVanity])

	ist := &types.AtlasExtra{
		Signature:     vals,
		Bitmap:        []byte{},
		CommittedSeal: [][]byte{},
	}

	payload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		return nil, err
	}

	return append(buf.Bytes(), payload...), nil
}

// writeSeal writes the extra-data field of the given header with the given seals.
// suggest to rename to writeSeal.
func writeSeal(h *types.Header, seal []byte, pubkey []byte) error {
	if len(seal) != types.AtlasExtraSignature || len(pubkey) != types.AtlasExtraPublicKey {
		return errInvalidSignature
	}

	atlasExtra, err := types.ExtractAtlasExtra(h)
	if err != nil {
		return err
	}

	atlasExtra.Signature = seal
	atlasExtra.PublicKey = pubkey

	payload, err := rlp.EncodeToBytes(&atlasExtra)
	if err != nil {
		return err
	}

	h.Extra = append(h.Extra[:types.AtlasExtraVanity], payload...)
	return nil
}

// writeCommittedSeals writes the extra-data field of a block header with given committed seals.
func writeCommittedSeals(h *types.Header, signature []byte, bitmap []byte) error {
	if len(signature) != types.AtlasExtraSignature || len(bitmap) != types.AtlasExtraMask {
		return errInvalidCommittedSeals
	}

	atlasExtra, err := types.ExtractAtlasExtra(h)
	if err != nil {
		return err
	}

	atlasExtra.AggSignature = make([]byte, len(signature))
	copy(atlasExtra.AggSignature, signature)
	atlasExtra.AggBitmap = make([]byte, len(bitmap))
	copy(atlasExtra.AggBitmap, bitmap)

	payload, err := rlp.EncodeToBytes(&atlasExtra)
	if err != nil {
		return err
	}

	h.Extra = append(h.Extra[:types.AtlasExtraVanity], payload...)
	return nil
}

// ATLAS(yhx): getLargestAmountStakingValidators
func getLargestAmountStakingValidators(state *state.StateDB, numVal int) ([]staking.Validator, error) {
	container := state.GetStakingInfo(staking.StakingInfoAddress)

	// ATLAS(yhx): sort validators based on staking amount, select top numVal validators
	if container == nil {
		return nil, errValidatorNotExist
	}
	amount := make(map[common.Address]*big.Int)
	for _, val := range container.Validators {
		amount[val.Validator.Address] = val.Amount()
	}

	// sort by amount
	type pair struct {
		key   common.Address
		value *big.Int
	}
	var pairs []pair
	for k, v := range amount {
		pairs = append(pairs, pair{k, v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].value.Cmp(pairs[j].value) > 0
	})

	if numVal > len(pairs) {
		numVal = len(pairs)
	}
	addresses := make([]common.Address, numVal)
	for i := 0; i < numVal; i++ {
		addresses[i] = pairs[i].key
	}
	return addresses, nil
}

func convertValidators(vals []staking.Validator) []atlas.Validator{
	var retval []atlas.Validator = make([]atlas.Validator, 0, len(vals))

	for i, v := range vals {
		retval = append(retval, {

		})
	}
}
