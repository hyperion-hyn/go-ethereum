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
	"math/big"
	"sync"
	"testing"

	"github.com/hyperion-hyn/bls/ffi/go/bls"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/atlas"
	bls_cosi "github.com/ethereum/go-ethereum/crypto/bls"
)

func newTestRoundState(view *atlas.View, validatorSet atlas.ValidatorSet) *roundState {
	pubKeys := validatorSet.GetPublicKeys()
	prepareBitmap, _ := bls_cosi.NewMask(pubKeys, nil)
	confirmBitmap, _ := bls_cosi.NewMask(pubKeys, nil)

	return &roundState{
		round:      view.Round,
		sequence:   view.Sequence,
		Preprepare: newTestPreprepare(view),
		mu:         new(sync.RWMutex),
		hasBadProposal: func(hash common.Hash) bool {
			return false
		},
		aggregatedPrepareSig: &bls.Sign{},
		prepareBitmap:        prepareBitmap,
		aggregatedConfirmSig: &bls.Sign{},
		confirmBitmap:        confirmBitmap,
	}
}

func TestLockHash(t *testing.T) {
	sys := NewTestSystemWithBackend(1, 0)
	rs := newTestRoundState(
		&atlas.View{
			Round:    big.NewInt(0),
			Sequence: big.NewInt(0),
		},
		sys.backends[0].peers,
	)
	if !common.EmptyHash(rs.GetLockedHash()) {
		t.Errorf("error mismatch: have %v, want empty", rs.GetLockedHash())
	}
	if rs.IsHashLocked() {
		t.Error("IsHashLocked should return false")
	}

	// Lock
	expected := rs.Proposal().Hash()
	rs.LockHash()
	if expected != rs.GetLockedHash() {
		t.Errorf("error mismatch: have %v, want %v", rs.GetLockedHash(), expected)
	}
	if !rs.IsHashLocked() {
		t.Error("IsHashLocked should return true")
	}

	// Unlock
	rs.UnlockHash()
	if !common.EmptyHash(rs.GetLockedHash()) {
		t.Errorf("error mismatch: have %v, want empty", rs.GetLockedHash())
	}
	if rs.IsHashLocked() {
		t.Error("IsHashLocked should return false")
	}
}
