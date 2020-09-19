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

	"github.com/hyperion-hyn/bls/ffi/go/bls"

	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/atlas"
)

func (c *core) sendPreprepare(request *atlas.Request) {
	logger := c.logger.New("state", c.state)

	// If I'm the proposer and I have the same sequence with the proposal
	if c.current.Sequence().Cmp(request.Proposal.Number()) == 0 && c.IsProposer() {
		curView := c.currentView()

		hash := request.Proposal.SealHash(c.backend)
		signature, _, _, err := c.backend.SignHash(hash)
		if err != nil {
			logger.Error("Failed to SignHash", "err", err)
			return
		}

		preprepare, err := Encode(&atlas.Preprepare{
			View:      curView,
			Proposal:  request.Proposal,
			Signature: signature,
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
		return err
	}

	if err := c.verifyPreprepare(&preprepare, src); err != nil {
		return err
	}

	// Verify the proposal we received
	if duration, err := c.backend.Verify(preprepare.Proposal); err != nil {
		// if it's a future block, we will handle it again after the duration
		if err == consensus.ErrFutureBlock {
			logger.Info("Proposed block will be handled in the future", "err", err, "duration", duration)
			c.stopFuturePreprepareTimer()
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
		//   1. the locked proposal and the received proposal match
		//   2. we have no locked proposal
		if err := c.acceptPreprepare(&preprepare); err != nil {
			return err
		}
		c.setState(StatePreprepared)
		c.current.LockHash()
		c.sendPrepare()
	}

	return nil
}

// verifyPrepare verifies if the received PREPARE message is equivalent to our subject
func (c *core) verifyPreprepare(preprepare *atlas.Preprepare, src atlas.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	var sign bls.Sign
	if err := sign.Deserialize(preprepare.Signature); err != nil {
		logger.Error("Failed to deserialize signature", "err", err)
		return err
	}

	hash := preprepare.Proposal.SealHash(c.backend)
	if ok := sign.VerifyHash(src.PublicKey(), hash.Bytes()); !ok {
		return errInvalidSignature
	}

	return nil
}

func (c *core) acceptPreprepare(preprepare *atlas.Preprepare) error {
	// ATLAS(zgx): please refer to accpetPrepare
	c.consensusTimestamp = time.Now()
	c.current.SetPreprepare(preprepare)
	return nil
}
