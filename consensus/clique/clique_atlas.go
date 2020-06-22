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

// Package clique implements the proof-of-authority consensus engine.
package clique

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/consensus/clique/reward"
	"github.com/ethereum/go-ethereum/consensus/clique/votepower"
	"github.com/ethereum/go-ethereum/numeric"
	"github.com/ethereum/go-ethereum/staking/availability"
	"github.com/ethereum/go-ethereum/staking/committee"
	"github.com/ethereum/go-ethereum/staking/network"
	staking "github.com/ethereum/go-ethereum/staking/types"
	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
	"math/big"
	"math/rand"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/misc"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	lru "github.com/hashicorp/golang-lru"
)

// Clique is the proof-of-authority consensus engine proposed to support the
// Ethereum testnet following the Ropsten attacks.
type AtlasClique struct {
	config *params.CliqueConfig // Consensus engine configuration parameters
	db     ethdb.Database       // Database to store and retrieve snapshot checkpoints

	recents    *lru.ARCCache // Snapshots for recent block to speed up reorgs
	signatures *lru.ARCCache // Signatures of recent blocks to speed up mining

	proposals map[common.Address]bool // Current list of proposals we are pushing

	signer common.Address // Ethereum address of the signing key
	signFn SignerFn       // Signer function to authorize hashes with
	lock   sync.RWMutex   // Protects the signer fields

	// The fields below are for testing only
	fakeDiff bool // Skip difficulty verifications
}

// New creates a Clique proof-of-authority consensus engine with the initial
// signers set to the ones provided by the user.
// Just for staking testing.
func NewAtlasClique(config *params.CliqueConfig, db ethdb.Database) *AtlasClique {
	// Set any missing consensus parameters to their defaults
	conf := *config
	if conf.Epoch == 0 {
		conf.Epoch = epochLength
	}
	// Allocate the snapshot caches and create the engine
	recents, _ := lru.NewARC(inmemorySnapshots)
	signatures, _ := lru.NewARC(inmemorySignatures)

	return &AtlasClique{
		config:     &conf,
		db:         db,
		recents:    recents,
		signatures: signatures,
		proposals:  make(map[common.Address]bool),
	}
}

// Author implements consensus.Engine, returning the Ethereum address recovered
// from the signature in the header's extra-data section.
func (c *AtlasClique) Author(header *types.Header) (common.Address, error) {
	return ecrecover(header, c.signatures)
}

// VerifyHeader checks whether a header conforms to the consensus rules.
func (c *AtlasClique) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	return c.verifyHeader(chain, header, nil)
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers. The
// method returns a quit channel to abort the operations and a results channel to
// retrieve the async verifications (the order is that of the input slice).
func (c *AtlasClique) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
	results := make(chan error, len(headers))

	go func() {
		for i, header := range headers {
			err := c.verifyHeader(chain, header, headers[:i])

			select {
			case <-abort:
				return
			case results <- err:
			}
		}
	}()
	return abort, results
}

// verifyHeader checks whether a header conforms to the consensus rules.The
// caller may optionally pass in a batch of parents (ascending order) to avoid
// looking those up from the database. This is useful for concurrently verifying
// a batch of new headers.
func (c *AtlasClique) verifyHeader(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	if header.Number == nil {
		return errUnknownBlock
	}
	number := header.Number.Uint64()

	// Don't waste time checking blocks from the future
	if header.Time > uint64(time.Now().Unix()) {
		return consensus.ErrFutureBlock
	}
	// Checkpoint blocks need to enforce zero beneficiary
	checkpoint := (number % c.config.Epoch) == 0
	if checkpoint && header.Coinbase != (common.Address{}) {
		return errInvalidCheckpointBeneficiary
	}
	// Nonces must be 0x00..0 or 0xff..f, zeroes enforced on checkpoints
	if !bytes.Equal(header.Nonce[:], nonceAuthVote) && !bytes.Equal(header.Nonce[:], nonceDropVote) {
		return errInvalidVote
	}
	if checkpoint && !bytes.Equal(header.Nonce[:], nonceDropVote) {
		return errInvalidCheckpointVote
	}
	// Check that the extra-data contains both the vanity and signature
	if len(header.Extra) < extraVanity {
		return errMissingVanity
	}
	if len(header.Extra) < extraVanity+extraSeal {
		return errMissingSignature
	}
	// Ensure that the extra-data contains a signer list on checkpoint, but none otherwise
	signersBytes := len(header.Extra) - extraVanity - extraSeal
	if !checkpoint && signersBytes != 0 {
		return errExtraSigners
	}
	if checkpoint && signersBytes%common.AddressLength != 0 {
		return errInvalidCheckpointSigners
	}
	// Ensure that the mix digest is zero as we don't have fork protection currently
	if header.MixDigest != (common.Hash{}) {
		return errInvalidMixDigest
	}
	// Ensure that the block doesn't contain any uncles which are meaningless in PoA
	if header.UncleHash != uncleHash {
		return errInvalidUncleHash
	}
	// Ensure that the block's difficulty is meaningful (may not be correct at this point)
	if number > 0 {
		if header.Difficulty == nil || (header.Difficulty.Cmp(diffInTurn) != 0 && header.Difficulty.Cmp(diffNoTurn) != 0) {
			return errInvalidDifficulty
		}
	}
	// If all checks passed, validate any special fields for hard forks
	if err := misc.VerifyForkHashes(chain.Config(), header, false); err != nil {
		return err
	}
	// All basic checks passed, verify cascading fields
	return c.verifyCascadingFields(chain, header, parents)
}

// verifyCascadingFields verifies all the header fields that are not standalone,
// rather depend on a batch of previous headers. The caller may optionally pass
// in a batch of parents (ascending order) to avoid looking those up from the
// database. This is useful for concurrently verifying a batch of new headers.
func (c *AtlasClique) verifyCascadingFields(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	// The genesis block is the always valid dead-end
	number := header.Number.Uint64()
	if number == 0 {
		return nil
	}
	// Ensure that the block's timestamp isn't too close to its parent
	var parent *types.Header
	if len(parents) > 0 {
		parent = parents[len(parents)-1]
	} else {
		parent = chain.GetHeader(header.ParentHash, number-1)
	}
	if parent == nil || parent.Number.Uint64() != number-1 || parent.Hash() != header.ParentHash {
		return consensus.ErrUnknownAncestor
	}
	if parent.Time+c.config.Period > header.Time {
		return ErrInvalidTimestamp
	}
	// Retrieve the snapshot needed to verify this header and cache it
	snap, err := c.snapshot(chain, number-1, header.ParentHash, parents)
	if err != nil {
		return err
	}
	// If the block is a checkpoint block, verify the signer list
	if number%c.config.Epoch == 0 {
		signers := make([]byte, len(snap.Signers)*common.AddressLength)
		for i, signer := range snap.signers() {
			copy(signers[i*common.AddressLength:], signer[:])
		}
		extraSuffix := len(header.Extra) - extraSeal
		if !bytes.Equal(header.Extra[extraVanity:extraSuffix], signers) {
			return errMismatchingCheckpointSigners
		}
	}
	// All basic checks passed, verify the seal and return
	return c.verifySeal(chain, header, parents)
}

// snapshot retrieves the authorization snapshot at a given point in time.
func (c *AtlasClique) snapshot(chain consensus.ChainReader, number uint64, hash common.Hash, parents []*types.Header) (*Snapshot, error) {
	// Search for a snapshot in memory or on disk for checkpoints
	var (
		headers []*types.Header
		snap    *Snapshot
	)
	for snap == nil {
		// If an in-memory snapshot was found, use that
		if s, ok := c.recents.Get(hash); ok {
			snap = s.(*Snapshot)
			break
		}
		// If an on-disk checkpoint snapshot can be found, use that
		if number%checkpointInterval == 0 {
			if s, err := loadSnapshot(c.config, c.signatures, c.db, hash); err == nil {
				log.Trace("Loaded voting snapshot from disk", "number", number, "hash", hash)
				snap = s
				break
			}
		}
		// If we're at the genesis, snapshot the initial state. Alternatively if we're
		// at a checkpoint block without a parent (light client CHT), or we have piled
		// up more headers than allowed to be reorged (chain reinit from a freezer),
		// consider the checkpoint trusted and snapshot it.
		if number == 0 || (number%c.config.Epoch == 0 && (len(headers) > params.ImmutabilityThreshold || chain.GetHeaderByNumber(number-1) == nil)) {
			checkpoint := chain.GetHeaderByNumber(number)
			if checkpoint != nil {
				hash := checkpoint.Hash()

				signers := make([]common.Address, (len(checkpoint.Extra)-extraVanity-extraSeal)/common.AddressLength)
				for i := 0; i < len(signers); i++ {
					copy(signers[i][:], checkpoint.Extra[extraVanity+i*common.AddressLength:])
				}
				snap = newSnapshot(c.config, c.signatures, number, hash, signers)
				if err := snap.store(c.db); err != nil {
					return nil, err
				}
				log.Info("Stored checkpoint snapshot to disk", "number", number, "hash", hash)
				break
			}
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
	c.recents.Add(snap.Hash, snap)

	// If we've generated a new checkpoint snapshot, save to disk
	if snap.Number%checkpointInterval == 0 && len(headers) > 0 {
		if err = snap.store(c.db); err != nil {
			return nil, err
		}
		log.Trace("Stored voting snapshot to disk", "number", snap.Number, "hash", snap.Hash)
	}
	return snap, err
}

// VerifyUncles implements consensus.Engine, always returning an error for any
// uncles as this consensus mechanism doesn't permit uncles.
func (c *AtlasClique) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	if len(block.Uncles()) > 0 {
		return errors.New("uncles not allowed")
	}
	return nil
}

// VerifySeal implements consensus.Engine, checking whether the signature contained
// in the header satisfies the consensus protocol requirements.
func (c *AtlasClique) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	return c.verifySeal(chain, header, nil)
}

// verifySeal checks whether the signature contained in the header satisfies the
// consensus protocol requirements. The method accepts an optional list of parent
// headers that aren't yet part of the local blockchain to generate the snapshots
// from.
func (c *AtlasClique) verifySeal(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	// Verifying the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return errUnknownBlock
	}
	// Retrieve the snapshot needed to verify this header and cache it
	snap, err := c.snapshot(chain, number-1, header.ParentHash, parents)
	if err != nil {
		return err
	}

	// Resolve the authorization key and check against signers
	signer, err := ecrecover(header, c.signatures)
	if err != nil {
		return err
	}
	if _, ok := snap.Signers[signer]; !ok {
		return errUnauthorizedSigner
	}
	for seen, recent := range snap.Recents {
		if recent == signer {
			// Signer is among recents, only fail if the current block doesn't shift it out
			if limit := uint64(len(snap.Signers)/2 + 1); seen > number-limit {
				return errRecentlySigned
			}
		}
	}
	// Ensure that the difficulty corresponds to the turn-ness of the signer
	if !c.fakeDiff {
		inturn := snap.inturn(header.Number.Uint64(), signer)
		if inturn && header.Difficulty.Cmp(diffInTurn) != 0 {
			return errWrongDifficulty
		}
		if !inturn && header.Difficulty.Cmp(diffNoTurn) != 0 {
			return errWrongDifficulty
		}
	}
	return nil
}

// Prepare implements consensus.Engine, preparing all the consensus fields of the
// header for running the transactions on top.
func (c *AtlasClique) Prepare(chain consensus.ChainReader, header *types.Header) error {
	// If the block isn't a checkpoint, cast a random vote (good enough for now)
	header.Coinbase = common.Address{}
	header.Nonce = types.BlockNonce{}

	number := header.Number.Uint64()
	// Assemble the voting snapshot to check which votes make sense
	snap, err := c.snapshot(chain, number-1, header.ParentHash, nil)
	if err != nil {
		return err
	}
	if number%c.config.Epoch != 0 {
		c.lock.RLock()

		// Gather all the proposals that make sense voting on
		addresses := make([]common.Address, 0, len(c.proposals))
		for address, authorize := range c.proposals {
			if snap.validVote(address, authorize) {
				addresses = append(addresses, address)
			}
		}
		// If there's pending proposals, cast a vote on them
		if len(addresses) > 0 {
			header.Coinbase = addresses[rand.Intn(len(addresses))]
			if c.proposals[header.Coinbase] {
				copy(header.Nonce[:], nonceAuthVote)
			} else {
				copy(header.Nonce[:], nonceDropVote)
			}
		}
		c.lock.RUnlock()
	}
	// Set the correct difficulty
	header.Difficulty = CalcDifficulty(snap, c.signer)

	// Ensure the extra data has all its components
	if len(header.Extra) < extraVanity {
		header.Extra = append(header.Extra, bytes.Repeat([]byte{0x00}, extraVanity-len(header.Extra))...)
	}
	header.Extra = header.Extra[:extraVanity]

	if number%c.config.Epoch == 0 {
		for _, signer := range snap.signers() {
			header.Extra = append(header.Extra, signer[:]...)
		}
	}
	header.Extra = append(header.Extra, make([]byte, extraSeal)...)

	// Mix digest is reserved for now, set to empty
	header.MixDigest = common.Hash{}

	// Ensure the timestamp has the correct delay
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	header.Time = parent.Time + c.config.Period
	if header.Time < uint64(time.Now().Unix()) {
		header.Time = uint64(time.Now().Unix())
	}
	return nil
}

// Finalize implements consensus.Engine, ensuring no uncles are set, nor block
// rewards given.
func (c *AtlasClique) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header) {
	// No block rewards in PoA, so the state remains as is and uncles are dropped
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = types.CalcUncleHash(nil)
}

// FinalizeAndAssemble implements consensus.Engine, ensuring no uncles are set,
// nor block rewards given, and returns the final block.
func (c *AtlasClique) FinalizeAndAssemble(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	// ATLAS
	isNewEpoch := chain.Config().Atlas.IsLastBlock(header.Number.Uint64())
	if isNewEpoch {
		// unredelegation
		if err := payoutUnredelegations(chain, header, state); err != nil {
			return nil, err
		}

		// Needs to be before AccumulateRewardsAndCountSigs because
		// ComputeAndMutateEPOSStatus depends on the signing counts that's
		// consistent with the counts when the new shardState was proposed.
		// Refer to committee.IsEligibleForEPoSAuction()
		for _, addr := range curComm.StakedValidators().Addrs {		// current committee
			if err := availability.ComputeAndMutateEPOSStatus(
				chain, state, addr,
			); err != nil {
				return nil, err
			}
		}

		// committee
		newComm := updateCommitteeForNextEpoch(chain, header, state)

		// Needs to be after payoutUndelegations because payoutUndelegations
		// depends on the old LastEpochInCommittee
		if err := setLastEpochInCommittee(newComm, state); err != nil {
			return nil, err
		}
	}

	// TODO: map3 node undelegation by time
	// get all map3 node addresses
	// check if each node is to be released
	// auto renew node

	// reward
	accumulateRewardsAndCountSigs(chain, state ,header)

	// TODO: slash

	// ATLAS - END

	// No block rewards in PoA, so the state remains as is and uncles are dropped
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = types.CalcUncleHash(nil)

	// Assemble and return the final block for sealing
	return types.NewBlock(header, txs, nil, receipts), nil
}


// Authorize injects a private key into the consensus engine to mint new blocks
// with.
func (c *AtlasClique) Authorize(signer common.Address, signFn SignerFn) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.signer = signer
	c.signFn = signFn
}

// Seal implements consensus.Engine, attempting to create a sealed block using
// the local signing credentials.
func (c *AtlasClique) Seal(chain consensus.ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	header := block.Header()

	// Sealing the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return errUnknownBlock
	}
	// For 0-period chains, refuse to seal empty blocks (no reward but would spin sealing)
	if c.config.Period == 0 && len(block.Transactions()) == 0 {
		log.Info("Sealing paused, waiting for transactions")
		return nil
	}
	// Don't hold the signer fields for the entire sealing procedure
	c.lock.RLock()
	signer, signFn := c.signer, c.signFn
	c.lock.RUnlock()

	// Bail out if we're unauthorized to sign a block
	snap, err := c.snapshot(chain, number-1, header.ParentHash, nil)
	if err != nil {
		return err
	}
	if _, authorized := snap.Signers[signer]; !authorized {
		return errUnauthorizedSigner
	}
	// If we're amongst the recent signers, wait for the next block
	for seen, recent := range snap.Recents {
		if recent == signer {
			// Signer is among recents, only wait if the current block doesn't shift it out
			if limit := uint64(len(snap.Signers)/2 + 1); number < limit || seen > number-limit {
				log.Info("Signed recently, must wait for others")
				return nil
			}
		}
	}
	// Sweet, the protocol permits us to sign the block, wait for our time
	delay := time.Unix(int64(header.Time), 0).Sub(time.Now()) // nolint: gosimple
	if header.Difficulty.Cmp(diffNoTurn) == 0 {
		// It's not our turn explicitly to sign, delay it a bit
		wiggle := time.Duration(len(snap.Signers)/2+1) * wiggleTime
		delay += time.Duration(rand.Int63n(int64(wiggle)))

		log.Trace("Out-of-turn signing requested", "wiggle", common.PrettyDuration(wiggle))
	}
	// Sign all the things!
	sighash, err := signFn(accounts.Account{Address: signer}, accounts.MimetypeClique, CliqueRLP(header))
	if err != nil {
		return err
	}
	copy(header.Extra[len(header.Extra)-extraSeal:], sighash)
	// Wait until sealing is terminated or delay timeout.
	log.Trace("Waiting for slot to sign and propagate", "delay", common.PrettyDuration(delay))
	go func() {
		select {
		case <-stop:
			return
		case <-time.After(delay):
		}

		select {
		case results <- block.WithSeal(header):
		default:
			log.Warn("Sealing result is not read by miner", "sealhash", SealHash(header))
		}
	}()

	return nil
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
// that a new block should have based on the previous blocks in the chain and the
// current signer.
func (c *AtlasClique) CalcDifficulty(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	snap, err := c.snapshot(chain, parent.Number.Uint64(), parent.Hash(), nil)
	if err != nil {
		return nil
	}
	return CalcDifficulty(snap, c.signer)
}

// SealHash returns the hash of a block prior to it being sealed.
func (c *AtlasClique) SealHash(header *types.Header) common.Hash {
	return SealHash(header)
}

// Close implements consensus.Engine. It's a noop for clique as there are no background threads.
func (c *AtlasClique) Close() error {
	return nil
}

// APIs implements consensus.Engine, returning the user facing RPC API to allow
// controlling the signer voting.
func (c *AtlasClique) APIs(chain consensus.ChainReader) []rpc.API {
	return []rpc.API{{
		Namespace: "clique",
		Version:   "1.0",
		Service:   &API{chain: chain},
		Public:    false,
	}}
}

// ATLAS
func setLastEpochInCommittee(newComm *committee.Committee, state *state.StateDB) error {
	for _, addr := range newComm.StakedValidators().Addrs {
		wrapper, err := state.ValidatorWrapper(addr)
		if err != nil {
			return errors.New(
				"[Finalize] failed to get validator from state to finalize",
			)
		}
		wrapper.LastEpochInCommittee = newShardState.Epoch
	}
	return nil
}

// Withdraw unlocked tokens to the delegators' accounts
func payoutUnredelegations(chain consensus.ChainReader, header *types.Header, state *state.StateDB) error {
	nowEpoch, blockNow := header.Epoch, header.Number

	// TODO(storage): read all validator addresses from state db
	validators, err := chain.ReadValidatorList()
	if err != nil {
		const msg = "[Finalize] failed to read all validators"
		return errors.New(msg)
	}
	// Payout undelegated/unlocked tokens
	for _, validator := range validators {
		wrapper, err := state.ValidatorWrapper(validator)
		if err != nil {
			return errors.New(
				"[Finalize] failed to get validator from state to finalize",
			)
		}
		lockPeriod := staking.LockPeriodInEpoch
		for i := range wrapper.Redelegations {
			redelegation := &wrapper.Redelegations[i]
			totalWithdraw := redelegation.RemoveUnlockedUndelegations(
				header.Epoch, wrapper.LastEpochInCommittee, lockPeriod,
			)
			//state.AddBalance(delegation.DelegatorAddress, totalWithdraw)
			// TODO: update map3 node balance
		}
	}
	log.Info("paid out delegations", "epoch", nowEpoch.Uint64(), "block-number", blockNow.Uint64())
	return nil
}

func updateCommitteeForNextEpoch(chain consensus.ChainReader, header *types.Header, db *state.StateDB) error {
	if chain.Config().Atlas.IsLastBlock(header.Number.Uint64()) {
		nextComm, err := committee.WithStakingEnabled.Compute(header.Epoch, chain)
		if err != nil {
			return err
		}
		// TODO(storage): save next committee to statedb
	}
	return nil
}

// BlockCommitSig returns the byte array of aggregated
// commit signature and bitmap signed on the block
func blockCommitSig(chain consensus.ChainReader, blockNum uint64) ([]byte, []byte, error) {
	if blockNum < 1 {
		return nil, nil, nil
	}

	extra := chain.GetHeaderByNumber(blockNum).Extra
	// TODO: read aggSig and bitmap from extra
	return aggSig, bitmap, nil
}

// accumulateRewardsAndCountSigs credits the coinbase of the given block with the mining
// reward. The total reward consists of the static block reward
// This func also do IncrementValidatorSigningCounts for validators
func accumulateRewardsAndCountSigs(
	bc consensus.ChainReader, state *state.StateDB, header *types.Header,
) (reward.Reader, error) {
	blockNum := header.Number.Uint64()
	nowEpoch := header.Epoch

	if blockNum == 0 {
		// genesis block has no parent to reward.
		return network.EmptyPayout, nil
	}

	defaultReward := network.BlockReward

	// TODO(ATLAS): adjust block reward by block num

	// If too much is staked, then possible to have negative reward,
	// not an error, just a possible economic situation, hence we return
	if defaultReward.Sign() == -1 {		// negative
		return network.EmptyPayout, nil
	}

	newRewards, payouts := big.NewInt(0), []reward.Payout{}
	members, payable, missing, err := ballotResult(bc, header)
	if err != nil {
		return network.EmptyPayout, err
	}
	subComm := committee.Committee{members}		// TODO: from state?

	if err := availability.IncrementValidatorSigningCounts(
		bc,
		subComm.StakedValidators(),
		state,
		payable,
		missing,
	); err != nil {
		return network.EmptyPayout, err
	}
	votingPower, err := lookupVotingPower(
		nowEpoch, &subComm,
	)
	if err != nil {
		return network.EmptyPayout, err
	}

	allSignersShare := numeric.ZeroDec()
	for j := range payable {
		voter := votingPower.Voters[payable[j].BLSPublicKey]
		voterShare := voter.OverallPercent
		allSignersShare = allSignersShare.Add(voterShare)
	}
	for member := range payable {
		// TODO Give out whatever leftover to the last voter/handle
		// what to do about share of those that didn't sign
		blsKey := payable[member].BLSPublicKey
		voter := votingPower.Voters[blsKey]

		// TODO(storage): read validator detail by epoch
		snapshot, err := bc.ReadValidatorSnapshot(voter.EarningAccount)
		if err != nil {
			return network.EmptyPayout, err
		}
		due := defaultReward.Mul(
			voter.OverallPercent.Quo(allSignersShare),
		).RoundInt()
		newRewards.Add(newRewards, due)

		shares, err := lookupDelegatorShares(snapshot)
		if err != nil {
			return network.EmptyPayout, err
		}
		// TODO(storage): update delegation reward
		if err := state.AddReward(snapshot.Validator, due, shares); err != nil {
			return network.EmptyPayout, err
		}
		payouts = append(payouts, reward.Payout{
			Addr:        voter.EarningAccount,
			NewlyEarned: due,
			EarningKey:  voter.Identity,
		})
	}
	return network.NewStakingEraRewardForRound(newRewards, missing, payouts), nil
}

func ballotResult(
	bc consensus.ChainReader, header *types.Header,
) (committee.SlotList, committee.SlotList, committee.SlotList, error) {
	//parentHeader := bc.GetHeaderByHash(header.ParentHash)
	//if parentHeader == nil {
	//	return nil, nil, nil, errors.Errorf(
	//		"cannot find parent block header in DB %s",
	//		header.ParentHash.Hex(),
	//	)
	//}
	//parentShardState, err := bc.ReadShardState(parentHeader.Epoch())
	//if err != nil {
	//	return nil, nil, nil, errors.Errorf(
	//		"cannot read shard state %v", parentHeader.Epoch(),
	//	)
	//}
	// TODO(ATLAS): get parent committee by epoch
	return availability.BallotResult(header, parentCommittee)
}

var (
	votingPowerCache   singleflight.Group
	delegateShareCache singleflight.Group
)

func lookupVotingPower(epoch *big.Int, comm *committee.Committee) (*votepower.Roster, error) {
	key := epoch.String()
	results, err, _ := votingPowerCache.Do(
		key, func() (interface{}, error) {
			votingPower, err := votepower.Compute(comm, epoch)
			if err != nil {
				return nil, err
			}

			// For new calc, remove old data from 3 epochs ago
			deleteEpoch := big.NewInt(0).Sub(epoch, big.NewInt(3))
			deleteKey := deleteEpoch.String()
			votingPowerCache.Forget(deleteKey)

			return votingPower, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return results.(*votepower.Roster), nil
}

// Lookup or compute the shares of stake for all delegators in a validator
func lookupDelegatorShares(
	snapshot *staking.ValidatorSnapshot,
) (map[common.Address]numeric.Dec, error) {
	epoch := snapshot.Epoch
	validatorSnapshot := snapshot.Validator
	key := fmt.Sprintf("%s-%s", epoch.String(), validatorSnapshot.ValidatorAddress.Hex())

	shares, err, _ := delegateShareCache.Do(
		key, func() (interface{}, error) {
			result := map[common.Address]numeric.Dec{}

			totalDelegationDec := numeric.NewDecFromBigInt(validatorSnapshot.TotalDelegation())
			if totalDelegationDec.IsZero() {
				log.Info("zero total delegation during AddReward delegation payout",
					"validator-snapshot", validatorSnapshot.ValidatorAddress.Hex())
				return result, nil
			}

			for i := range validatorSnapshot.Redelegations {
				delegation := validatorSnapshot.Redelegations[i]
				// NOTE percentage = <this_delegator_amount>/<total_delegation>
				percentage := numeric.NewDecFromBigInt(delegation.Amount).Quo(totalDelegationDec)
				result[delegation.DelegatorAddress] = percentage
			}

			// For new calc, remove old data from 3 epochs ago
			deleteEpoch := big.NewInt(0).Sub(epoch, big.NewInt(3))
			deleteKey := fmt.Sprintf("%s-%s", deleteEpoch.String(), validatorSnapshot.ValidatorAddress.Hex())
			votingPowerCache.Forget(deleteKey)

			return result, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return shares.(map[common.Address]numeric.Dec), nil
}

// ATLAS - END