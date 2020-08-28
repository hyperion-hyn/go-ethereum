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
	"github.com/ethereum/go-ethereum/crypto"
	bls_cosi "github.com/ethereum/go-ethereum/crypto/bls"
)

func (c *core) sendCommit() {
	logger := c.logger.New("state", c.state)

	// If I'm the proposer and I have the same sequence with the proposal
	if c.IsProposer() {
		sub, err := c.AssembleSignedSubject()
		if err != nil {
			logger.Error("Failed to sign", "view", c.currentView(), "err", err)
		}

		encodedSubject, err := Encode(sub)
		if err != nil {
			logger.Error("Failed to encode", "subject", sub, "err", err)
			return
		}

		c.broadcast(&message{
			Code: msgConfirm,
			Msg:  encodedSubject,
		})
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
	c.broadcast(&message{
		Code: msgCommit,
		Msg:  encodedSubject,
	})
}

func (c *core) handleCommit(msg *message, src atlas.Validator) error {
	// Decode COMMIT message
	var commit *atlas.SignedSubject
	err := msg.Decode(&commit)
	if err != nil {
		return errFailedDecodeCommit
	}

	if err := c.checkMessage(msgCommit, commit.Subject.View); err != nil {
		return err
	}

	if err := c.verifyCommit(commit, src); err != nil {
		return err
	}

	c.acceptCommit(msg, src, c.valSet)

	// Commit the proposal once we have enough COMMIT messages and we are not in the Confirm state.
	//
	// If we already have a proposal, we may have chance to speed up the consensus process
	// by committing the proposal without PREPARE messages.
	if c.current.Commits.Size() >= c.QuorumSize() && c.state.Cmp(StateCommitted) < 0 {
		// Still need to call LockHash here since state can skip Expect state and jump directly to the Confirm state.
		c.current.LockHash()
		c.commit()

	}

	return nil
}

// verifyCommit verifies if the received COMMIT message is equivalent to our subject
func (c *core) verifyCommit(commit *atlas.SignedSubject, src atlas.Validator) error {
	return nil
}

func (c *core) acceptCommit(msg *message, src atlas.Validator, validatorSet atlas.ValidatorSet) error {
	logger := c.logger.New("from", src, "state", c.state)

	// Add the COMMIT message to current round state
	if err := c.current.Commits.Add(msg); err != nil {
		logger.Error("Failed to record commit message", "msg", msg, "err", err)
		return err
	}

	var commit *atlas.SignedSubject
	if err := msg.Decode(&commit); err != nil {
		return errFailedDecodePrepare
	}

	if commit.Subject.Digest != c.current.Preprepare.Proposal.Hash() {
		logger.Warn("Inconsistent subjects between EXPECT and proposal", "expected", c.current.Preprepare.Proposal.Hash(), "got", commit.Subject.Digest)
		return errInconsistentSubject
	}

	var sign bls.Sign
	if err := sign.Deserialize(commit.Signature); err != nil {
		logger.Error("Failed to deserialize signature", "msg", msg, "err", err)
		return err
	}

	var pubKey bls.PublicKey
	if err := pubKey.Deserialize(commit.PublicKey); err != nil {
		logger.Error("Failed to deserialize signer's public key", "msg", msg, "err", err)
		return err
	}

	hash := crypto.Keccak256Hash(commit.Subject.Digest.Bytes())
	if sign.VerifyHash(&pubKey, hash.Bytes()) == false {
		logger.Error("Failed to verify signature with signer's public key commit", "msg", msg)
		return errInvalidSignature
	}

	bitmap, _ := bls_cosi.NewMask(validatorSet.GetPublicKeys(), nil)
	if err := bitmap.SetMask(commit.Mask); err != nil {
		logger.Error("Failed to SetMask", "view", c.currentView(), "err", err)
		return err
	}

	if bitmap.CountEnabled() < c.QuorumSize() {
		return errNotSatisfyQuorum
	}

	return nil
}
