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

	if err := c.checkMessage(msgExpect, expect.View); err != nil {
		return err
	}

	if err := c.verifyExpect(&expect, src); err != nil {
		return err
	}

	// Here is about to accept the EXPECT message
	// validator in Preprepared state, leader in Prepared state
	if c.state == StatePreprepared || c.state == StatePrepared {
		// Send ROUND CHANGE if the locked proposal and the received proposal are different
		if expect.Digest == c.current.GetLockedHash() {
			if err := c.acceptExpect(&expect, src); err != nil {
				return err
			}
			c.setState(StateExpected)
			c.sendConfirm()
		} else {
			// Send round change
			logger.Warn("state error", "send next round", "currentView", c.currentView())
			c.sendNextRoundChange()
		}
	}

	return nil
}

func (c *core) verifyExpect(expect *atlas.Subject, src atlas.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	sub := c.current.Subject()
	if !atlas.IsConsistentSubject(expect, sub) {
		logger.Warn("Inconsistent subjects between expect and proposal", "expected", sub, "got", expect)
		return errInconsistentSubject
	}

	return nil
}

func (c *core) acceptExpect(expect *atlas.Subject, src atlas.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	if err := c.backend.CheckSignature(expect.Payload, c.valSet.GetProposer().PublicKey().Serialize(), expect.Signature); err == errInvalidSignature {
		logger.Error("Leader give a expect with invalid signature")
		c.sendNextRoundChange()
		return errInvalidSignature
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

	var sign bls.Sign
	if err := sign.Deserialize(signPayload.Signature); err != nil {
		logger.Error("Failed to deserialize signature", "err", err)
		return err
	}

	c.consensusTimestamp = time.Now()
	c.current.SetExpect(expect)
	return nil
}
