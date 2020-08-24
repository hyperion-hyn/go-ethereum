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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/atlas"
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
	var prepare *atlas.SignedSubject
	err := msg.Decode(&prepare)
	if err != nil {
		return errFailedDecodePrepare
	}

	if err := c.checkMessage(msgPrepare, prepare.View); err != nil {
		return err
	}

	// If it is locked, it can only process on the locked block.
	// Passing verifyPrepare and checkMessage implies it is processing on the locked block since it was verified in the Preprepared state.
	if err := c.verifyPrepare(prepare, src); err != nil {
		return err
	}

	// Signer should be in the validator set
	if _, v := c.valSet.GetByAddress(src.Address()); v == nil {
		return errNotFromCommittee
	}

	c.acceptPrepare(msg, src)

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
func (c *core) verifyPrepare(prepare *atlas.SignedSubject, src atlas.Validator) error {
	return nil
}

func (c *core) acceptPrepare(msg *message, src atlas.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	// Add the PREPARE message to current round state
	if err := c.current.Prepares.Add(msg); err != nil {
		logger.Error("Failed to add PREPARE message to round state", "msg", msg, "err", err)
		return err
	}

	var prepare *atlas.SignedSubject
	if err := msg.Decode(&prepare); err != nil {
		return errFailedDecodePrepare
	}

	if prepare.Digest != c.current.Preprepare.Proposal.Hash() {
		logger.Warn("Inconsistent subjects between PREPARE and proposal", "expected", c.current.Preprepare.Proposal.Hash(), "got", prepare.Digest)
		return errInconsistentSubject
	}

	var sign bls.Sign
	if err := sign.Deserialize(prepare.Signature); err != nil {
		logger.Error("Failed to deserialize signature", "msg", msg, "err", err)
		return err
	}

	var signer common.Address
	signer = src.Address()

	_, validator := c.valSet.GetByAddress(signer)
	if validator == nil {
		return errInvalidSigner
	}

	var pubKey *bls.PublicKey = validator.PublicKey()
	if validator == nil {
		return errInvalidSigner
	}

	if sign.Verify(pubKey, prepare.Digest.String()) == false {
		logger.Error("Failed to verify signature with signer's public key", "msg", msg)
		return errInvalidSignature
	}

	if err := c.current.prepareBitmap.SetKey(pubKey, true); err != nil {
		c.current.aggregatedPrepareSig.Add(&sign)
	}

	return nil
}
