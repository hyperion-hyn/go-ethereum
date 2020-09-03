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
	"github.com/hyperion-hyn/bls/ffi/go/bls"

	"github.com/ethereum/go-ethereum/consensus/atlas"
	"github.com/ethereum/go-ethereum/rlp"
)

func (c *core) sendPrepare() {
	logger := c.logger.New("state", c.state)

	sub, err := c.SignSubject(c.current.Subject())
	if err != nil {
		logger.Error("Failed to sign", "view", c.currentView(), "err", err)
	}

	encodedSubject, err := Encode(sub)
	if err != nil {
		logger.Error("Failed to encode", "subject", sub, "err", err)
		return
	}

	c.broadcast(&message{
		Code: msgPrepare,
		Msg:  encodedSubject,
	})
}

func (c *core) handlePrepare(msg *message, src atlas.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	// Decode PREPARE message
	var prepare atlas.Subject
	err := msg.Decode(&prepare)
	if err != nil {
		return errFailedDecodePrepare
	}

	if err := c.checkMessage(msgPrepare, prepare.View); err != nil {
		return err
	}

	// If it is locked, it can only process on the locked block.
	// Passing verifyPrepare and checkMessage implies it is processing on the locked block since it was verified in the Preprepared state.
	if err := c.verifyPrepare(&prepare, src); err != nil {
		return err
	}

	// Signer should be in the validator set
	if _, v := c.valSet.GetBySigner(src.Signer()); v == nil {
		return errNotFromCommittee
	}

	// ATLAS(zgx): after send out expect, should stop accept prepare.
	if err := c.acceptPrepare(msg, src); err != nil {
		return err
	}

	if !c.IsProposer() {
		logger.Error("message come from no-proposer", "msg", msg)
		return nil
	}

	// Change to Expect state if we've received enough PREPARE messages or it is locked
	// and we are in earlier state before Expect state.
	if ((c.current.IsHashLocked() && prepare.Digest == c.current.GetLockedHash()) || c.current.GetPrepareSize() >= c.QuorumSize()) &&
		c.state.Cmp(StatePrepared) < 0 {
		c.current.LockHash()
		c.setState(StatePrepared)
		c.sendExpect()
	}

	return nil
}

// verifyPrepare verifies if the received PREPARE message is equivalent to our subject
func (c *core) verifyPrepare(prepare *atlas.Subject, src atlas.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	sub := c.current.Subject()
	if !atlas.IsConsistentSubject(prepare, sub) {
		logger.Warn("Inconsistent subjects between PREPARE and proposal", "expected", sub, "got", prepare)
		return errInconsistentSubject
	}

	// ATLAS(zgx): should verifySignPayload here?
	return nil
}

func (c *core) acceptPrepare(msg *message, src atlas.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	var prepare atlas.Subject
	if err := msg.Decode(&prepare); err != nil {
		return err
	}

	if prepare.Digest != c.current.Preprepare.Proposal.Hash() {
		logger.Warn("Inconsistent subjects between PREPARE and proposal", "expected", c.current.Preprepare.Proposal.Hash(), "got", prepare.Digest)
		return errInconsistentSubject
	}

	var signPayload *atlas.SignPayload
	if err := rlp.DecodeBytes(prepare.Payload, &signPayload); err != nil {
		return err
	}

	var sign bls.Sign
	if err := sign.Deserialize(signPayload.Signature); err != nil {
		logger.Error("Failed to deserialize signature", "signature", signPayload.Signature, "err", err)
		return err
	}

	pubKey, err := c.getValidatorPublicKey(src.Signer(), c.valSet)
	if err != nil {
		return err
	}

	err = c.checkValidatorSignature(prepare.Digest.Bytes(), signPayload.Signature, pubKey.Serialize())
	if err != nil {
		logger.Error("Failed to verify signature with signer's public key prepare", "signature", signPayload.Signature[:10], "publicKey", pubKey.Serialize()[:10])
		return err
	}

	enabled, err := c.current.prepareBitmap.KeyEnabled(pubKey)
	if err != nil {
		return err
	} else if enabled == true {
		return errDuplicateMessage
	}

	if c.state == StatePrepared {
		return nil
	}

	if err := c.current.prepareBitmap.SetKey(pubKey, true); err == nil {
		c.current.aggregatedPrepareSig.Add(&sign)
		c.current.aggregatedPreparePublicKey.Add(pubKey)
	}

	// Add the PREPARE message to current round state
	if err := c.current.Prepares.Add(msg); err != nil {
		logger.Error("Failed to add PREPARE message to round state", "msg", msg, "err", err)
		return err
	}

	return nil
}
