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
	"io"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/atlas"
	bls_cosi "github.com/ethereum/go-ethereum/crypto/bls"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/hyperion-hyn/bls/ffi/go/bls"
)

// newRoundState creates a new roundState instance with the given view and validatorSet
// lockedHash and preprepare are for round change when lock exists,
// we need to keep a reference of preprepare in order to propose locked proposal when there is a lock and itself is the proposer
func newRoundState(view *atlas.View, validatorSet atlas.ValidatorSet, lockedHash common.Hash, preprepare *atlas.Preprepare, pendingRequest *atlas.Request, hasBadProposal func(hash common.Hash) bool) *roundState {
	pubKeys := validatorSet.GetPublicKeys()
	prepareBitmap, _ := bls_cosi.NewMask(pubKeys, nil)
	confirmBitmap, _ := bls_cosi.NewMask(pubKeys, nil)

	return &roundState{
		round:          view.Round,
		sequence:       view.Sequence,
		Preprepare:     preprepare,
		Prepares:       newMessageSet(validatorSet),
		Commits:        newMessageSet(validatorSet),
		lockedHash:     lockedHash,
		mu:             new(sync.RWMutex),
		pendingRequest: pendingRequest,
		hasBadProposal: hasBadProposal,
		prepareBitmap:  prepareBitmap,
		confirmBitmap:  confirmBitmap,
	}
}

// roundState stores the consensus state
type roundState struct {
	round          *big.Int
	sequence       *big.Int
	Preprepare     *atlas.Preprepare
	Prepares       *messageSet
	Expect         *atlas.Subject
	Commits        *messageSet
	Confirm        *atlas.Subject
	lockedHash     common.Hash
	pendingRequest *atlas.Request

	mu             *sync.RWMutex
	hasBadProposal func(hash common.Hash) bool

	aggregatedPrepareSig       *bls.Sign
	aggregatedPreparePublicKey *bls.PublicKey
	prepareBitmap              *bls_cosi.Mask
	aggregatedConfirmSig       *bls.Sign
	aggregatedConfirmPublicKey *bls.PublicKey
	confirmBitmap              *bls_cosi.Mask
}

func (s *roundState) GetPrepareSize() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.prepareBitmap.CountEnabled()
}

func (s *roundState) GetConfirmSize() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.confirmBitmap.CountEnabled()
}

func (s *roundState) Subject() *atlas.Subject {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.Preprepare == nil {
		return nil
	}

	return &atlas.Subject{
		View: &atlas.View{
			Round:    new(big.Int).Set(s.round),
			Sequence: new(big.Int).Set(s.sequence),
		},
		Digest: s.Preprepare.Proposal.Hash(),
	}
}

func (s *roundState) SetPreprepare(preprepare *atlas.Preprepare) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Preprepare = preprepare
}

func (s *roundState) SetExpect(expect *atlas.Subject) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Expect = expect
}

func (s *roundState) SetCommitted(committed *atlas.Subject) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Confirm = committed
}

func (s *roundState) Proposal() atlas.Proposal {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.Preprepare != nil {
		return s.Preprepare.Proposal
	}

	return nil
}

func (s *roundState) SetRound(r *big.Int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.round = new(big.Int).Set(r)
}

func (s *roundState) Round() *big.Int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.round
}

func (s *roundState) SetSequence(seq *big.Int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sequence = seq
}

func (s *roundState) Sequence() *big.Int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.sequence
}

func (s *roundState) LockHash() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Preprepare != nil {
		s.lockedHash = s.Preprepare.Proposal.Hash()
	}
}

func (s *roundState) UnlockHash() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lockedHash = common.Hash{}
}

func (s *roundState) IsHashLocked() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if common.EmptyHash(s.lockedHash) {
		return false
	}
	return !s.hasBadProposal(s.GetLockedHash())
}

func (s *roundState) GetLockedHash() common.Hash {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.lockedHash
}

// The DecodeRLP method should read one value from the given
// Stream. It is not forbidden to read less or more, but it might
// be confusing.
func (s *roundState) DecodeRLP(stream *rlp.Stream) error {
	var ss struct {
		Round          *big.Int
		Sequence       *big.Int
		Preprepare     *atlas.Preprepare
		Prepares       *messageSet
		Commits        *messageSet
		lockedHash     common.Hash
		pendingRequest *atlas.Request

		aggregatedPrepareSig *bls.Sign
		prepareBitmap        *bls_cosi.Mask
		aggregatedCommitSig  *bls.Sign
		commitBitmap         *bls_cosi.Mask
	}

	if err := stream.Decode(&ss); err != nil {
		return err
	}
	s.round = ss.Round
	s.sequence = ss.Sequence
	s.Preprepare = ss.Preprepare
	s.Prepares = ss.Prepares
	s.Commits = ss.Commits
	s.lockedHash = ss.lockedHash
	s.pendingRequest = ss.pendingRequest

	s.aggregatedPrepareSig = ss.aggregatedPrepareSig
	s.prepareBitmap = ss.prepareBitmap
	s.aggregatedConfirmSig = ss.aggregatedCommitSig
	s.confirmBitmap = ss.commitBitmap

	s.mu = new(sync.RWMutex)

	return nil
}

// EncodeRLP should write the RLP encoding of its receiver to w.
// If the implementation is a pointer method, it may also be
// called for nil pointers.
//
// Implementations should generate valid RLP. The data written is
// not verified at the moment, but a future version might. It is
// recommended to write only a single value but writing multiple
// values or no value at all is also permitted.
func (s *roundState) EncodeRLP(w io.Writer) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return rlp.Encode(w, []interface{}{
		s.round,
		s.sequence,
		s.Preprepare,
		s.Prepares,
		s.Commits,
		s.lockedHash,
		s.pendingRequest,
		s.aggregatedPrepareSig,
		s.prepareBitmap,
		s.aggregatedConfirmSig,
		s.confirmBitmap,
	})
}
