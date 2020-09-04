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
)

func (c *core) sendCommit() {
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
			Code: msgCommit,
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

	// Commit the proposal once we have enough COMMIT messages and we are not in the Confirm state.
	//
	// If we already have a proposal, we may have chance to speed up the consensus process
	// by committing the proposal without PREPARE messages.
	if c.current.Confirms.Size() >= c.QuorumSize() && c.state.Cmp(StateCommitted) < 0 {
		// Still need to call LockHash here since state can skip Expect state and jump directly to the Confirm state.
		c.current.LockHash()
		c.commit()

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

	if commit.Digest != c.current.Preprepare.Proposal.SealHash(c.backend) {
		logger.Warn("Inconsistent subjects between EXPECT and proposal", "expected", c.current.Preprepare.Proposal.SealHash(c.backend), "got", commit.Digest)
		return errInconsistentSubject
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

	if bitmap.CountEnabled() < c.QuorumSize() {
		return errNotSatisfyQuorum
	}

	return nil
}
