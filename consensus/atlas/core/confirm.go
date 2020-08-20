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
	"github.com/ethereum/go-ethereum/rlp"
)

type ConfirmPayload struct {
	Digest    common.Hash
	Signature []byte
	Signer    []byte
}

func (c *core) sendConfirm() {
	logger := c.logger.New("state", c.state)

	hash := c.current.Preprepare.Proposal.Hash()
	sign, _, err := c.backend.Sign(hash.Bytes())
	signer := c.backend.Signer().Bytes()
	if err != nil {
		logger.Error("Failed to sign", "view", c.currentView())
		return
	}

	proposal, err := rlp.EncodeToBytes(&ConfirmPayload{
		Digest:    hash,
		Signature: sign,
		Signer:    signer,
	})
	if err != nil {
		logger.Error("Failed to encode proposal", "view", c.currentView())
	}

	encoded, err := Encode(&atlas.Prepare{
		View:     c.currentView(),
		Proposal: proposal,
	})
	if err != nil {
		logger.Error("Failed to encode payload", "view", c.currentView())
		return
	}

	c.broadcast(&message{
		Code: msgPrepare,
		Msg:  encoded,
	})
}

func (c *core) handleConfirm(msg *message, src atlas.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	// Decode PREPARE message
	var confirm *atlas.Confirm
	err := msg.Decode(&confirm)
	if err != nil {
		return errFailedDecodePrepare
	}

	if err := c.checkMessage(msgPrepare, confirm.View); err != nil {
		return err
	}

	// If it is locked, it can only process on the locked block.
	// Passing verifyPrepare and checkMessage implies it is processing on the locked block since it was verified in the Preprepared state.
	if err := c.verifyConfirm(confirm, src); err != nil {
		return err
	}

	// Signer should be in the validator set
	if _, v := c.valSet.GetByAddress(src.Address()); v == nil {
		return errNotFromCommittee
	}

	c.acceptConfirm(msg, src)

	if !c.IsProposer() {
		return nil
	}

	var proposal PreparePayload
	if err := rlp.DecodeBytes(confirm.Proposal, proposal); err != nil {
		logger.Error("Failed to decode payload", "view", c.currentView(), "err", err)
		return err
	}

	// Change to Expect state if we've received enough CONFIRM messages or it is locked
	// and we are in earlier state before Expect state.
	if ((c.current.IsHashLocked() && proposal.Digest == c.current.GetLockedHash()) || c.current.GetConfirmSize() >= c.QuorumSize()) &&
		c.state.Cmp(StateExpected) < 0 {
		c.current.LockHash()
		c.setState(StateExpected)
		c.sendCommit()
	}

	return nil
}

// verifyPrepare verifies if the received CONFIRM message is equivalent to our subject
func (c *core) verifyConfirm(prepare *atlas.Confirm, src atlas.Validator) error {
	return nil
}

func (c *core) acceptConfirm(msg *message, src atlas.Validator) error {
	// TODO(zgx): should reset if error occure
	logger := c.logger.New("from", src, "state", c.state)

	// Add the PREPARE message to current round state
	if err := c.current.Prepares.Add(msg); err != nil {
		logger.Error("Failed to add PREPARE message to round state", "msg", msg, "err", err)
		return err
	}

	var prepare *atlas.Confirm
	if err := msg.Decode(&prepare); err != nil {
		return errFailedDecodePrepare
	}

	var proposal PreparePayload
	if err := rlp.DecodeBytes(prepare.Proposal, proposal); err != nil {
		logger.Error("Failed to decode payload", "view", c.currentView(), "err", err)
		return err
	}

	if proposal.Digest != c.current.Preprepare.Proposal.Hash() {
		logger.Warn("Inconsistent subjects between PREPARE and proposal", "expected", c.current.Preprepare.Proposal.Hash(), "got", proposal.Digest)
		return errInconsistentSubject
	}

	var sign bls.Sign
	if err := sign.Deserialize(proposal.Signature); err != nil {
		logger.Error("Failed to deserialize signature", "msg", msg, "err", err)
		return err
	}

	var signer common.Address
	signer.SetBytes(proposal.Signer)

	_, validator := c.valSet.GetByAddress(signer)
	if validator == nil {
		return errInvalidSigner
	}

	var pubKey *bls.PublicKey = validator.PublicKey()
	if validator == nil {
		return errInvalidSigner
	}

	if sign.Verify(pubKey, proposal.Digest.String()) == false {
		logger.Error("Failed to verify signature with signer's public key", "msg", msg)
		return errInvalidSignature
	}

	if err := c.current.confirmBitmap.SetKey(pubKey, true); err != nil {
		c.current.aggregatedConfirmSig.Add(&sign)
	}

	c.current.aggregatedConfirmPublicKey.Add(pubKey)

	c.current.confirmBitmap.Mask()

	return nil
}
