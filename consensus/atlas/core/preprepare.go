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
)

func (c *core) sendPreprepare(request *atlas.Request) {
	logger := c.logger.New("state", c.state)

	// If I'm the proposer and I have the same sequence with the proposal
	if c.current.Sequence().Cmp(request.Proposal.Number()) == 0 && c.IsProposer() {
		curView := c.currentView()

		preprepare, err := Encode(&atlas.Preprepare{
			View:     curView,
			Proposal: request.Proposal,
		})
		if err != nil {
			logger.Error("Failed to encode", "view", curView)
			return
		}

		c.broadcast(&message{
			Code: msgPreprepare,
			Msg:  preprepare,
		})
	}
}

func (c *core) handlePreprepare(msg *message, src atlas.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	// Decode PRE-PREPARE
	var preprepare atlas.Preprepare
	err := msg.Decode(&preprepare)
	if err != nil {
		return errFailedDecodePreprepare
	}

	// Check if the message comes from current proposer
	if !c.valSet.IsProposer(src.Signer()) {
		logger.Warn("Ignore preprepare messages from non-proposer")
		return errNotFromProposer
	}

	// Ensure we have the same view with the PRE-PREPARE message
	// If it is old message, see if we need to broadcast COMMIT
	if err := c.checkMessage(msgPreprepare, preprepare.View); err != nil {
		if err == errOldMessage {
			// ATLAS(zgx): When/Why will node receive the old message? if this node receive old message, whether these messages have been broadcast to peer?
			// Get validator set for the given proposal
			valSet := c.backend.ParentValidators(preprepare.Proposal).Copy()
			previousProposer := c.backend.GetProposer(preprepare.Proposal.Number().Uint64() - 1)
			valSet.CalcProposer(previousProposer, preprepare.View.Round.Uint64())
			// Broadcast COMMIT if it is an existing block
			// 1. The proposer needs to be a proposer matches the given (Sequence + Round)
			// 2. The given block must exist
			if valSet.IsProposer(src.Signer()) && c.backend.HasPropsal(preprepare.Proposal.Hash(), preprepare.Proposal.Number()) {
				// ATLAS(zgx): maybe nothing can be done for old block for lacking multiple-signature.
				// c.sendCommitForOldBlock(preprepare.View, preprepare.Proposal.Hash())
				return nil
			}
		}
		return err
	}

	// Verify the proposal we received
	if duration, err := c.backend.Verify(preprepare.Proposal); err != nil {
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

	// Here is about to accept the PRE-PREPARE
	if c.state == StateAcceptRequest {
		// Send ROUND CHANGE if the locked proposal and the received proposal are different
		// ATLAS(zgx): have checked message come from proposer, Why I am a proposer and have a different hash? duplicated messages?
		if c.current.IsHashLocked() {
			if preprepare.Proposal.Hash() == c.current.GetLockedHash() {
				// Broadcast COMMIT and enters Expect state directly
				if err := c.acceptPreprepare(&preprepare); err != nil {
					return err
				}
				// ATLAS(zgx): LockHash is invoked in handlePrepare, so set state to StatePrepared directly
				c.setState(StatePrepared)
				c.sendExpect()
			} else {
				// Send round change
				c.sendNextRoundChange()
			}
		} else {
			// Either
			//   1. the locked proposal and the received proposal match
			//   2. we have no locked proposal
			if err := c.acceptPreprepare(&preprepare); err != nil {
				return err
			}
			c.setState(StatePreprepared)
			c.sendPrepare()
		}
	}

	return nil
}

func (c *core) acceptPreprepare(preprepare *atlas.Preprepare) error {
	// ATLAS(zgx): please refer to accpetPrepare
	c.consensusTimestamp = time.Now()
	c.current.SetPreprepare(preprepare)
	return nil
}
