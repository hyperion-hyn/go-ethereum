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
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/atlas"
	"github.com/ethereum/go-ethereum/consensus/atlas/validator"
)

func TestHandleCommit(t *testing.T) {
	N := uint64(4)
	F := uint64(1)

	proposal := newTestProposal()
	expectedSubject := &atlas.Subject{
		View: &atlas.View{
			Round:    big.NewInt(0),
			Sequence: proposal.Number(),
		},
		Digest: atlas.SealHash(proposal.Header()),
	}

	testCases := []struct {
		system      *testSystem
		expectedErr error
	}{
		{
			// normal case
			func() *testSystem {
				sys := NewTestSystemWithBackend(N, F)

				for i, backend := range sys.backends {
					c := backend.engine.(*core)
					c.valSet = backend.peers
					c.current = newTestRoundState(
						&atlas.View{
							Round:    big.NewInt(0),
							Sequence: big.NewInt(1),
						},
						c.valSet,
					)

					if i == 0 {
						// replica 0 is the proposer
						c.state = StateConfirmed
					} else {
						c.state = StateExpected
					}
				}
				return sys
			}(),
			nil,
		},
		{
			// future message
			func() *testSystem {
				sys := NewTestSystemWithBackend(N, F)

				for i, backend := range sys.backends {
					c := backend.engine.(*core)
					c.valSet = backend.peers
					if i == 0 {
						// replica 0 is the proposer
						c.current = newTestRoundState(
							expectedSubject.View,
							c.valSet,
						)
						c.state = StateConfirmed
					} else {
						c.state = StateExpected
						c.current = newTestRoundState(
							&atlas.View{
								Round:    big.NewInt(2),
								Sequence: big.NewInt(3),
							},
							c.valSet,
						)
					}
				}
				return sys
			}(),
			errOldMessage,
		},
		{
			// subject not match
			func() *testSystem {
				sys := NewTestSystemWithBackend(N, F)

				for i, backend := range sys.backends {
					c := backend.engine.(*core)
					c.valSet = backend.peers
					if i == 0 {
						// replica 0 is the proposer
						c.current = newTestRoundState(
							expectedSubject.View,
							c.valSet,
						)
						c.state = StateConfirmed
					} else {
						c.state = StateExpected
						c.current = newTestRoundState(
							&atlas.View{
								Round:    big.NewInt(0),
								Sequence: big.NewInt(0),
							},
							c.valSet,
						)
					}
				}
				return sys
			}(),
			errFutureMessage,
		},
		{
			// jump state
			func() *testSystem {
				sys := NewTestSystemWithBackend(N, F)

				for i, backend := range sys.backends {
					c := backend.engine.(*core)
					c.valSet = backend.peers
					c.current = newTestRoundState(
						&atlas.View{
							Round:    big.NewInt(0),
							Sequence: proposal.Number(),
						},
						c.valSet,
					)

					// only replica0 stays at StatePreprepared
					// other replicas are at StatePrepared
					if i != 0 {
						c.state = StatePreprepared
					} else {
						c.state = StateConfirmed
					}
				}
				return sys
			}(),
			nil,
		},
		// TODO: double send message
	}

OUTER:
	for _, test := range testCases {
		test.system.Run(false)

		v0 := test.system.backends[0]
		r0 := v0.engine.(*core)

		for _, v := range test.system.backends {
			c := v.engine.(*core)
			_, validator := r0.valSet.GetBySigner(c.Signer())
			signedSubject, err := c.SignSubject(c.current.Subject())
			if err != nil {
				t.Errorf("failed to sing subject")
			}
			m, _ := Encode(signedSubject)
			s := r0.state
			r0.state = StateExpected
			if err := r0.acceptConfirm(&message{
				Code:          msgConfirm,
				Msg:           m,
				Signer:        validator.Signer(),
				Signature:     []byte{},
				CommittedSeal: validator.Signer().Bytes(), // small hack
			}, validator); err != nil {
				t.Errorf("failed to acceptConfirm message: %v", err)
			}
			r0.state = s
		}

		for _, v := range test.system.backends {
			validator := r0.valSet.GetProposer()
			c := v.engine.(*core)

			s := r0.state
			r0.state = StateConfirmed
			signedSubject, err := r0.AssembleSignedSubject(c.current.Subject())
			if err != nil {
				t.Errorf("failed to assemble subject: %v", err)
			}
			r0.state = s

			m, _ := Encode(signedSubject)

			if err := c.handleCommit(&message{
				Code:          msgCommit,
				Msg:           m,
				Signer:        validator.Signer(),
				Signature:     []byte{},
				CommittedSeal: validator.Signer().Bytes(), // small hack
			}, validator); err != nil {
				if err != test.expectedErr {
					t.Errorf("error mismatch: have %v, want %v", err, test.expectedErr)
				}
				if c.current.IsHashLocked() {
					t.Errorf("block should not be locked")
				}
				continue OUTER
			}
		}

		// prepared is normal case
		if r0.state != StateCommitted {
			// There are not enough commit messages in core
			if r0.state != StateConfirmed {
				t.Errorf("state mismatch: have %v, want %v", r0.state, StateConfirmed)
			}
			if r0.current.confirmBitmap.CountEnabled() >= r0.QuorumSize() {
				t.Errorf("the size of commit messages should be less than %v", r0.QuorumSize())
			}
			if r0.current.IsHashLocked() {
				t.Errorf("block should not be locked")
			}
			continue
		}

		// core should have 2F+1 before Ceil2Nby3Block or Ceil(2N/3) prepare messages
		if r0.current.confirmBitmap.CountEnabled() < r0.QuorumSize() {
			t.Errorf("the size of commit messages should be larger than 2F+1 or Ceil(2N/3): size %v", r0.QuorumSize())
		}

		// check signatures large than F
		signedCount := r0.current.confirmBitmap.CountEnabled()
		if signedCount <= r0.valSet.F() {
			t.Errorf("the expected signed count should be larger than %v, but got %v", r0.valSet.F(), signedCount)
		}
		if !r0.current.IsHashLocked() {
			t.Errorf("block should be locked")
		}
	}
}

// round is not checked for now
func TestVerifyCommit(t *testing.T) {
	// for log purpose
	peer, _, _, err := newValidator()
	if err != nil {
		t.Errorf("failed to new a validator: %v", err)
	}
	valSet := validator.NewSet([]atlas.Validator{peer}, atlas.RoundRobin)

	sys := NewTestSystemWithBackend(uint64(1), uint64(0))

	testCases := []struct {
		expected   error
		commit     *atlas.Subject
		roundState *roundState
	}{
		{
			// normal case
			expected: nil,
			commit: &atlas.Subject{
				View:   &atlas.View{Round: big.NewInt(0), Sequence: big.NewInt(0)},
				Digest: atlas.SealHash(newTestProposal().Header()),
			},
			roundState: newTestRoundState(
				&atlas.View{Round: big.NewInt(0), Sequence: big.NewInt(0)},
				valSet,
			),
		},
		{
			// old message
			expected: errInconsistentSubject,
			commit: &atlas.Subject{
				View:   &atlas.View{Round: big.NewInt(0), Sequence: big.NewInt(0)},
				Digest: atlas.SealHash(newTestProposal().Header()),
			},
			roundState: newTestRoundState(
				&atlas.View{Round: big.NewInt(1), Sequence: big.NewInt(1)},
				valSet,
			),
		},
		{
			// different digest
			expected: errInconsistentSubject,
			commit: &atlas.Subject{
				View:   &atlas.View{Round: big.NewInt(0), Sequence: big.NewInt(0)},
				Digest: common.StringToHash("1234567890"),
			},
			roundState: newTestRoundState(
				&atlas.View{Round: big.NewInt(1), Sequence: big.NewInt(1)},
				valSet,
			),
		},
		{
			// malicious package(lack of sequence)
			expected: errInconsistentSubject,
			commit: &atlas.Subject{
				View:   &atlas.View{Round: big.NewInt(0), Sequence: nil},
				Digest: atlas.SealHash(newTestProposal().Header()),
			},
			roundState: newTestRoundState(
				&atlas.View{Round: big.NewInt(1), Sequence: big.NewInt(1)},
				valSet,
			),
		},
		{
			// wrong prepare message with same sequence but different round
			expected: errInconsistentSubject,
			commit: &atlas.Subject{
				View:   &atlas.View{Round: big.NewInt(1), Sequence: big.NewInt(0)},
				Digest: atlas.SealHash(newTestProposal().Header()),
			},
			roundState: newTestRoundState(
				&atlas.View{Round: big.NewInt(0), Sequence: big.NewInt(0)},
				valSet,
			),
		},
		{
			// wrong prepare message with same round but different sequence
			expected: errInconsistentSubject,
			commit: &atlas.Subject{
				View:   &atlas.View{Round: big.NewInt(0), Sequence: big.NewInt(1)},
				Digest: atlas.SealHash(newTestProposal().Header()),
			},
			roundState: newTestRoundState(
				&atlas.View{Round: big.NewInt(0), Sequence: big.NewInt(0)},
				valSet,
			),
		},
	}
	for i, test := range testCases {
		c := sys.backends[0].engine.(*core)
		c.current = test.roundState

		signedSubject, err := c.SignSubject(test.commit)
		if err != nil {
			t.Errorf("failed to sign subject: %v", err)
		}

		err = c.verifyCommit(signedSubject, peer)
		if err != test.expected {
			t.Errorf("result %d: error mismatch: have %v, want %v", i, err, test.expected)
		}
	}
}
