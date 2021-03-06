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
	"encoding/hex"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/atlas"
	"github.com/ethereum/go-ethereum/consensus/atlas/validator"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
)

const (
	dbKeySnapshotPrefix = "atlas-snapshot"
)

// Vote represents a single vote that an authorized validator made to modify the
// list of authorizations.
type Vote struct {
	Validator common.Address `json:"validator"` // Authorized validator that cast this vote
	Block     uint64         `json:"block"`     // Block number the vote was cast in (expire old votes)
	Address   common.Address `json:"address"`   // Account being voted on to change its authorization
	Authorize bool           `json:"authorize"` // Whether to authorize or deauthorize the voted account
}

// Tally is a simple vote tally to keep the current score of votes. Votes that
// go against the proposal aren't counted since it's equivalent to not voting.
type Tally struct {
	Authorize bool `json:"authorize"` // Whether the vote it about authorizing or kicking someone
	Votes     int  `json:"votes"`     // Number of votes until now wanting to pass the proposal
}

// Snapshot is the state of the authorization voting at a given point in time.
type Snapshot struct {
	Epoch uint64 // The number of blocks after which to checkpoint and reset the pending votes

	Number uint64                   `json:"number"`     // Block number where the snapshot was created
	Hash   common.Hash              `json:"hash"`       // Block hash where the snapshot was created
	Votes  []*Vote                  `json:"votes"`      // List of votes cast in chronological order
	Tally  map[common.Address]Tally `json:"tally"`      // Current vote tally to avoid recalculating
	ValSet atlas.ValidatorSet       `json:"validators"` // Set of authorized validators at this moment
}

// newSnapshot create a new snapshot with the specified startup parameters. This
// method does not initialize the set of recent validators, so only ever use if for
// the genesis block.
func newSnapshot(epoch uint64, number uint64, hash common.Hash, valSet atlas.ValidatorSet) *Snapshot {
	snap := &Snapshot{
		Epoch:  epoch,
		Number: number,
		Hash:   hash,
		ValSet: valSet,
		Tally:  make(map[common.Address]Tally),
	}
	return snap
}

// loadSnapshot loads an existing snapshot from the database.
func loadSnapshot(epoch uint64, db ethdb.Database, hash common.Hash) (*Snapshot, error) {
	blob, err := db.Get(append([]byte(dbKeySnapshotPrefix), hash[:]...))
	if err != nil {
		return nil, err
	}
	snap := new(Snapshot)
	if err := json.Unmarshal(blob, snap); err != nil {
		return nil, err
	}
	snap.Epoch = epoch

	return snap, nil
}

// store inserts the snapshot into the database.
func (s *Snapshot) store(db ethdb.Database) error {
	blob, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return db.Put(append([]byte(dbKeySnapshotPrefix), s.Hash[:]...), blob)
}

// copy creates a deep copy of the snapshot, though not the individual votes.
func (s *Snapshot) copy() *Snapshot {
	cpy := &Snapshot{
		Epoch:  s.Epoch,
		Number: s.Number,
		Hash:   s.Hash,
		ValSet: s.ValSet.Copy(),
		Votes:  make([]*Vote, len(s.Votes)),
		Tally:  make(map[common.Address]Tally),
	}

	for address, tally := range s.Tally {
		cpy.Tally[address] = tally
	}
	copy(cpy.Votes, s.Votes)

	return cpy
}

// checkVote return whether it's a valid vote
func (s *Snapshot) checkVote(address common.Address, authorize bool) bool {
	_, validator := s.ValSet.GetBySigner(address)
	return (validator != nil && !authorize) || (validator == nil && authorize)
}

// cast adds a new vote into the tally.
func (s *Snapshot) cast(address common.Address, authorize bool) bool {
	// Ensure the vote is meaningful
	if !s.checkVote(address, authorize) {
		return false
	}
	// Cast the vote into an existing or new tally
	if old, ok := s.Tally[address]; ok {
		old.Votes++
		s.Tally[address] = old
	} else {
		s.Tally[address] = Tally{Authorize: authorize, Votes: 1}
	}
	return true
}

// uncast removes a previously cast vote from the tally.
func (s *Snapshot) uncast(address common.Address, authorize bool) bool {
	// If there's no tally, it's a dangling vote, just drop
	tally, ok := s.Tally[address]
	if !ok {
		return false
	}
	// Ensure we only revert counted votes
	if tally.Authorize != authorize {
		return false
	}
	// Otherwise revert the vote
	if tally.Votes > 1 {
		tally.Votes--
		s.Tally[address] = tally
	} else {
		delete(s.Tally, address)
	}
	return true
}

// apply creates a new authorization snapshot by applying the given headers to
// the original one.
func (s *Snapshot) apply(config *atlas.Config, chain consensus.ChainReader, headers []*types.Header) (*Snapshot, error) {
	// Allow passing in no headers for cleaner code
	if len(headers) == 0 {
		return s, nil
	}
	// Sanity check that the headers can be applied
	for i := 0; i < len(headers)-1; i++ {
		if headers[i+1].Number.Uint64() != headers[i].Number.Uint64()+1 {
			return nil, errInvalidVotingChain
		}
	}
	if headers[0].Number.Uint64() != s.Number+1 {
		return nil, errInvalidVotingChain
	}
	// Iterate through the headers and create a new snapshot
	snap := s.copy()

	for _, header := range headers[len(headers)-1:] {
		// Remove any votes on checkpoint blocks
		number := header.Number.Uint64()
		if number%s.Epoch == 0 {
			snap.Votes = nil
			snap.Tally = make(map[common.Address]Tally)
		}
		stateDB, err := chain.StateAt(header.Root)
		if err != nil {
			return nil, err
		}
		validators, err := getValidators(stateDB, MaxValidatorCount)
		if err != nil {
			return nil, err
		}
		snap.ValSet = validator.NewSet(validators, config.ProposerPolicy)
	}
	snap.Number += uint64(len(headers))
	snap.Hash = headers[len(headers)-1].Hash()

	return snap, nil
}

// validators retrieves the list of authorized validators in ascending order.
func (s *Snapshot) validators() []atlas.Validator {
	validators := make([]atlas.Validator, 0, s.ValSet.Size())
	for _, validator := range s.ValSet.List() {
		validators = append(validators, validator)
	}
	for i := 0; i < len(validators); i++ {
		for j := i + 1; j < len(validators); j++ {
			addri := validators[i].Signer()
			addrj := validators[j].Signer()
			if bytes.Compare(addri[:], addrj[:]) > 0 {
				validators[i], validators[j] = validators[j], validators[i]
			}
		}
	}
	return validators
}

type snapshotJSON struct {
	Epoch  uint64                   `json:"epoch"`
	Number uint64                   `json:"number"`
	Hash   common.Hash              `json:"hash"`
	Votes  []*Vote                  `json:"votes"`
	Tally  map[common.Address]Tally `json:"tally"`

	// for validator set
	Validators []validatorJSON      `json:"validators"`
	Policy     atlas.ProposerPolicy `json:"policy"`
}

func (s *Snapshot) toJSONStruct() *snapshotJSON {
	retval := &snapshotJSON{
		Epoch:  s.Epoch,
		Number: s.Number,
		Hash:   s.Hash,
		Votes:  s.Votes,
		Tally:  s.Tally,
		Policy: s.ValSet.Policy(),
	}
	vals := s.validators()
	retval.Validators = make([]validatorJSON, 0, len(vals))
	for _, val := range vals {
		retval.Validators = append(retval.Validators, validatorJSON{
			Coinbase: val.Coinbase(),
			PubKey:   val.PublicKey().SerializeToHexStr(),
		})
	}
	return retval
}

// Unmarshal from a json byte array
func (s *Snapshot) UnmarshalJSON(b []byte) error {
	var j snapshotJSON
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}

	s.Epoch = j.Epoch
	s.Number = j.Number
	s.Hash = j.Hash
	s.Votes = j.Votes
	s.Tally = j.Tally
	vals := make([]atlas.Validator, 0, len(j.Validators))
	for _, item := range j.Validators {
		if b, err := hex.DecodeString(item.PubKey); err != nil {
			return err
		} else if val, err := validator.New(b, item.Coinbase); err != nil {
			return err
		} else {
			vals = append(vals, val)
		}
	}
	s.ValSet = validator.NewSet(vals, j.Policy)
	return nil
}

// Marshal to a json byte array
func (s *Snapshot) MarshalJSON() ([]byte, error) {
	j := s.toJSONStruct()
	return json.Marshal(j)
}

type validatorJSON struct {
	Coinbase common.Address `json:"coinbase"`
	PubKey   string         `json:"pub_key"`
}
