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
	bls_cosi "github.com/ethereum/go-ethereum/crypto/bls"
	"github.com/ethereum/go-ethereum/rlp"
)

type CommitPayload struct {
	Digest    common.Hash
	Signature []byte
	PubKey    []byte
	Bitmap    []byte
}

func (c *core) sendCommit() {
	logger := c.logger.New("state", c.state)

	// If I'm the proposer and I have the same sequence with the proposal
	if c.IsProposer() {
		curView := c.currentView()

		hash := c.current.Preprepare.Proposal.Hash()
		sign := c.current.aggregatedConfirmSig.Serialize()
		pubKey := c.current.confirmBitmap.AggregatePublic.Serialize()
		bitmap := c.current.confirmBitmap.Bitmap

		payload := ExpectPayload{
			Digest:    hash,
			Signature: sign,
			PubKey:    pubKey,
			Bitmap:    bitmap,
		}

		proposal, err := rlp.EncodeToBytes(payload);
		if err != nil {
			logger.Error("failed to encode payload", "view", curView)
			return
		}

		confirmed, err := Encode(&atlas.Confirm{
			View:     curView,
			Proposal: proposal,
		})
		if err != nil {
			logger.Error("Failed to encode", "view", curView)
			return
		}

		c.broadcast(&message{
			Code: msgConfirm,
			Msg:  confirmed,
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
	var commit *atlas.Commit
	err := msg.Decode(&commit)
	if err != nil {
		return errFailedDecodeCommit
	}

	if err := c.checkMessage(msgCommit, commit.View); err != nil {
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
func (c *core) verifyCommit(commit *atlas.Commit, src atlas.Validator) error {
	return nil
}

func (c *core) acceptCommit(msg *message, src atlas.Validator, validatorSet atlas.ValidatorSet) error {
	logger := c.logger.New("from", src, "state", c.state)

	// Add the COMMIT message to current round state
	if err := c.current.Commits.Add(msg); err != nil {
		logger.Error("Failed to record commit message", "msg", msg, "err", err)
		return err
	}

	var commit *atlas.Commit
	if err := msg.Decode(&commit); err != nil {
		return errFailedDecodePrepare
	}

	var proposal CommitPayload
	if err := rlp.DecodeBytes(commit.Proposal, proposal); err != nil {
		logger.Error("Failed to decode payload", "view", c.currentView(), "err", err)
		return err
	}

	if proposal.Digest != c.current.Preprepare.Proposal.Hash() {
		logger.Warn("Inconsistent subjects between EXPECT and proposal", "expected", c.current.Preprepare.Proposal.Hash(), "got", proposal.Digest)
		return errInconsistentSubject
	}

	var sign bls.Sign
	if err := sign.Deserialize(proposal.Signature); err != nil {
		logger.Error("Failed to deserialize signature", "msg", msg, "err", err)
		return err
	}

	var pubKey bls.PublicKey
	if err := pubKey.Deserialize(proposal.PubKey); err != nil {
		logger.Error("Failed to deserialize signer's public key", "msg", msg, "err", err)
		return err
	}

	if sign.Verify(&pubKey, proposal.Digest.String()) == false {
		logger.Error("Failed to verify signature with signer's public key", "msg", msg)
		return errInvalidSignature
	}

	bitmap, _ := bls_cosi.NewMask(validatorSet.GetPublicKeys(), nil)
	if err := bitmap.SetMask(proposal.Bitmap); err != nil {
		logger.Error("Failed to SetMask", "view", c.currentView(), "err", err)
		return err
	}

	if bitmap.CountEnabled()  < c.QuorumSize() {
		return errNotSatisfyQuorum
	}

	return nil
}
