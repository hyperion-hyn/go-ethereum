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
	"time"

	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/atlas"
	bls_cosi "github.com/ethereum/go-ethereum/crypto/bls"
)

func (c *core) sendExpect() {
	logger := c.logger.New("state", c.state)

	// If I'm the proposer and I have the same sequence with the proposal
	if c.IsProposer() {
		sub, err := c.AssembleSignedSubject(c.current.Subject())
		if err != nil {
			logger.Error("Failed to sign", "view", c.currentView(), "err", err)
		}

		encodedSubject, err := Encode(sub)
		if err != nil {
			logger.Error("Failed to encode", "subject", sub, "err", err)
			return
		}

		c.broadcast(&message{
			Code: msgExpect,
			Msg:  encodedSubject,
		})
	}
}

func (c *core) handleExpect(msg *message, src atlas.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	// Decode EXPECT
	var expect atlas.Subject
	err := msg.Decode(&expect)
	if err != nil {
		return errFailedDecodeExpect
	}

	// Check if the message comes from current proposer
	if !c.valSet.IsProposer(src.Signer()) {
		logger.Warn("Ignore expect messages from non-proposer")
		return errNotFromProposer
	}

	// Ensure we have the same view with the PREPARED message
	// If it is old message, see if we need to broadcast COMMIT
	if err := c.checkMessage(msgExpect, expect.View); err != nil {
		if err == errOldMessage {
			// ATLAS(zgx): what if old message is different from preprepare.proposal?
			// Get validator set for the given proposal
			valSet := c.backend.ParentValidators(c.current.Preprepare.Proposal).Copy()
			previousProposer := c.backend.GetProposer(c.current.Preprepare.Proposal.Number().Uint64() - 1)
			valSet.CalcProposer(previousProposer, c.current.Preprepare.View.Round.Uint64())

			// Broadcast COMMIT if it is an existing block
			// 1. The proposer needs to be a proposer matches the given (Sequence + Round)
			// 2. The given block must exist
			if valSet.IsProposer(src.Signer()) && c.backend.HasPropsal(c.current.Preprepare.Proposal.Hash(), c.current.Preprepare.Proposal.Number()) {
				// ATLAS(zgx): maybe nothing can be done for old block for lacking multiple-signature.
				// c.sendCommitForOldBlock(c.current.Preprepare.View, c.current.Preprepare.Proposal.Hash())
				return nil
			}
		}
		return err
	}

	// Verify the proposal we received
	if duration, err := c.backend.Verify(c.current.Preprepare.Proposal); err != nil {
		// if it's a future block, we will handle it again after the duration
		if err == consensus.ErrFutureBlock {
			logger.Info("Proposed block will be handled in the future", "err", err, "duration", duration)
			c.stopFuturePreprepareTimer()
			// ATLAS(zgx): futurePreprepareTimer hold one timer, how to process multiple future block?
			c.futurePreprepareTimer = time.AfterFunc(duration, func() {
				c.sendEvent(backlogEvent{
					src: src,
					msg: msg,
				})
			})
		} else {
			logger.Warn("Failed to verify proposal", "err", err, "duration", duration)
			c.sendNextRoundChange()
		}
		return err
	}

	// Here is about to accept the PREPARED
	if c.state == StatePreprepared || c.state == StatePrepared {
		// Send ROUND CHANGE if the locked proposal and the received proposal are different
		if c.current.IsHashLocked() {
			if expect.Digest == c.current.GetLockedHash() {
				// Broadcast COMMIT and enters Expect state directly
				if err := c.acceptExpect(&expect); err != nil {
					return err
				}
				// ATLAS(zgx): LockHash in handlePrepare, so set state to StatePrepared directly
				c.setState(StateExpected)
				c.sendConfirm()
			} else {
				// Send round change
				c.sendNextRoundChange()
			}
		} else {
			// Either
			//   1. the locked proposal and the received proposal match
			//   2. we have no locked proposal
			if err := c.acceptExpect(&expect); err != nil {
				return err
			}
			c.setState(StateExpected)
			c.sendConfirm()
		}
	}

	return nil
}

func (c *core) acceptExpect(prepare *atlas.Subject) error {
	// ATLAS(zgx): please refer to acceptPrepare
	c.consensusTimestamp = time.Now()
	c.current.SetExpect(prepare)
	return nil
}

func (c *core) verifyExpect(expect *atlas.Subject, src atlas.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	sub := c.current.Subject()
	if !atlas.IsConsistentSubject(sub, expect) {
		logger.Warn("Inconsistent subjects between expect and proposal", "expected", sub, "got", expect)
		return errInconsistentSubject
	}

	signPayload, err := c.verifySignPayload(expect, c.valSet)
	if err != nil {
		return err
	}

	bitmap, _ := bls_cosi.NewMask(c.valSet.GetPublicKeys(), nil)
	if err := bitmap.SetMask(signPayload.Mask); err != nil {
		logger.Error("Failed to SetMask", "err", err)
		return err
	}

	if bitmap.CountEnabled() < c.QuorumSize() {
		return errNotSatisfyQuorum
	}

	return nil
}
