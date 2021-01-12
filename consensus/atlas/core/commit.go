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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/atlas"
	bls_cosi "github.com/ethereum/go-ethereum/crypto/bls"
	"github.com/hyperion-hyn/bls/ffi/go/bls"
	"time"
)

func (c *core) sendCommit() {
	logger := c.logger.New("state", c.state)

	signers := c.Signer()
	for _, signer := range signers {
		// If I'm the proposer and I have the same sequence with the proposal
		if c.IsProposer(signer) {
			sub, err := c.AssembleSignedSubject(signer, c.current.Subject())
			if err != nil {
				logger.Error("Failed to sign", "view", c.currentView(), "err", err)
				break
			}

			encodedSubject, err := Encode(sub)
			if err != nil {
				logger.Error("Failed to encode", "subject", sub, "err", err)
				break
			}
			c.consensusConfirmGauge.Update(time.Since(c.confirmTimestamp).Milliseconds())
			c.broadcast(signer, &message{
				Code: msgCommit,
				Msg:  encodedSubject,
			})
			break
		}
	}
}

func (c *core) sendCommitForOldBlock(view *atlas.View, digest common.Hash) {
	sub := &atlas.Subject{
		View:   view,
		Digest: digest,
	}
	c.broadcastCommit(sub)
}

func (c *core) broadcastCommit(sub *atlas.Subject) {
	logger := c.logger.New("state", c.state)

	encodedSubject, err := Encode(sub)
	if err != nil {
		logger.Error("Failed to encode", "subject", sub)
		return
	}
	signers := c.Signer()
	for _, signer := range signers { // TODO(Z): every signer?
		c.broadcast(signer, &message{
			Code: msgCommit,
			Msg:  encodedSubject,
		})
	}
}

func (c *core) handleCommit(msg *message, src atlas.Validator) error {
	// Decode COMMIT message
	var commit atlas.Subject
	err := msg.Decode(&commit)
	if err != nil {
		return errFailedDecodeCommit
	}

	if err := c.checkMessage(msgCommit, commit.View); err != nil {
		return err
	}

	if err := c.verifyCommit(&commit, src); err != nil {
		return err
	}

	if err := c.acceptCommit(msg, src); err != nil {
		return err
	}

	// Commit the proposal once we have enough CONFIRM signature and we are not in the Confirm state.
	//
	// If we already have a proposal in PREPREPARED state, we may have chance to speed up the consensus process
	// by committing the proposal without PREPARE messages.
	if c.current.confirmBitmap.CountEnabled() >= c.QuorumSize() && (c.state.Cmp(StatePreprepared) >= 0 && c.state.Cmp(StateCommitted) < 0) {
		// commit need proposal which was set in the PREPREPARED state, in other state can jump directly to the Confirm state.
		c.setState(StateCommitted)
		c.commit()
	} else {
		c.sendNextRoundChange()
	}

	return nil
}

// verifyCommit verifies if the received COMMIT message is equivalent to our subject
func (c *core) verifyCommit(commit *atlas.Subject, src atlas.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	sub := c.current.Subject()
	if !atlas.IsConsistentSubject(commit, sub) {
		logger.Warn("Inconsistent subjects between commit and proposal", "expected", sub, "got", commit)
		return errInconsistentSubject
	}

	return nil
}

func (c *core) acceptCommit(msg *message, src atlas.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	var commit atlas.Subject
	if err := msg.Decode(&commit); err != nil {
		return errFailedDecodeConfirm
	}

	if err := c.backend.CheckSignature(commit.Payload, c.valSet.GetProposer().PublicKey().Serialize(), commit.Signature); err == errInvalidSignature {
		logger.Error("Leader give a commit with invalid signature")
		c.sendNextRoundChange()
		return errInvalidSignature
	}

	signPayload, err := c.verifySignPayload(&commit, c.valSet)
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

	c.current.aggregatedConfirmSig = &sign
	c.current.confirmBitmap = bitmap

	return nil
}
