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

	"github.com/harmony-one/bls/ffi/go/bls"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/atlas"
	bls_cosi "github.com/ethereum/go-ethereum/crypto/bls"
	"github.com/ethereum/go-ethereum/rlp"
)

type ExpectPayload struct {
	Digest    common.Hash
	Signature []byte
	PubKey    []byte
	Bitmap    []byte
}

func (c *core) sendExpect() {
	logger := c.logger.New("state", c.state)

	// If I'm the proposer and I have the same sequence with the proposal
	if c.IsProposer() {
		curView := c.currentView()

		hash := c.current.Preprepare.Proposal.Hash()
		sign := c.current.aggregatedPrepareSig.Serialize()
		pubKey := c.current.prepareBitmap.AggregatePublic.Serialize()
		bitmap := c.current.prepareBitmap.Bitmap

		payload := ExpectPayload{
			Digest: hash,
			Signature: sign,
			PubKey: pubKey,
			Bitmap: bitmap,
		}

		proposal, err := rlp.EncodeToBytes(payload);
		if err != nil {
			logger.Error("failed to encode payload", "view", curView)
			return
		}

		prepared, err := Encode(&atlas.Expect{
			View:     curView,
			Proposal: proposal,
		})
		if err != nil {
			logger.Error("Failed to encode", "view", curView)
			return
		}

		c.broadcast(&message{
			Code: msgExpect,
			Msg:  prepared,
		})
	}
}

func (c *core) handleExpect(msg *message, src atlas.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	// Decode EXPECT
	var expect *atlas.Expect
	err := msg.Decode(&expect)
	if err != nil {
		return errFailedDecodePrepared
	}

	// Ensure we have the same view with the PREPARED message
	// If it is old message, see if we need to broadcast COMMIT
	if err := c.checkMessage(msgExpect, expect.View); err != nil {
		if err == errOldMessage {
			// ATLAS(zgx): what if old message is different from preprepare.proposal?
			// Get validator set for the given proposal
			valSet := c.backend.ParentValidators(c.current.Preprepare.Proposal).Copy()
			previousProposer := c.backend.GetProposer(c.current.Preprepare.Proposal.Number().Uint64() - 1)
			valSet.CalcProposer(previousProposer, c.current.Preprepare.View.Round.Uint64())

			if err := c.verifyExpect(msg, src, valSet); err != nil {
				return err
			}

			// Broadcast COMMIT if it is an existing block
			// 1. The proposer needs to be a proposer matches the given (Sequence + Round)
			// 2. The given block must exist
			if valSet.IsProposer(src.Address()) && c.backend.HasPropsal(c.current.Preprepare.Proposal.Hash(), c.current.Preprepare.Proposal.Number()) {
				// ATLAS(zgx): maybe nothing can be done for old block for lacking multiple-signature.
				// c.sendCommitForOldBlock(c.current.Preprepare.View, c.current.Preprepare.Proposal.Hash())
				return nil
			}
		}
		return err
	}

	// Check if the message comes from current proposer
	if !c.valSet.IsProposer(src.Address()) {
		logger.Warn("Ignore expect messages from non-proposer")
		return errNotFromProposer
	}

	// Here is about to accept the PREPARED
	if c.state == StatePreprepared || c.state == StatePrepared {
		if err := c.verifyExpect(msg, src, c.valSet); err != nil {
			c.sendNextRoundChange()
			return err
		}

		var proposal ExpectPayload
		if err := rlp.DecodeBytes(expect.Proposal, proposal); err != nil {
			logger.Error("Failed to decode payload", "view", c.currentView(), "err", err)
			return err
		}

		if proposal.Digest != c.current.Preprepare.Proposal.Hash() {
			return errInconsistentSubject
		}

		// Send ROUND CHANGE if the locked proposal and the received proposal are different
		if c.IsProposer() && c.current.IsHashLocked() {
			if proposal.Digest == c.current.GetLockedHash() {
				// Broadcast COMMIT and enters Expect state directly
				c.acceptExpect(expect)
				// ATLAS(zgx): LockHash in handlePrepare, so set state to StatePrepared directly
				c.setState(StateExpected)
				c.sendConfirm()
			} else {
				// Send round change
				c.sendNextRoundChange()
			}
		} else {
			// Either
			//   1. the locked proposal and the received proposal match
			//   2. we have no locked proposal
			c.acceptExpect(expect)
			c.setState(StateExpected)
			c.sendConfirm()
		}
	}

	return nil
}


func (c *core) acceptExpect(prepare *atlas.Expect) {
	c.consensusTimestamp = time.Now()
	c.current.SetExpect(prepare)
}

func (c *core) verifyExpect(msg *message, src atlas.Validator, validatorSet atlas.ValidatorSet) error {
	logger := c.logger.New("from", src, "state", c.state)

	var expect *atlas.Expect
	if err := msg.Decode(&expect); err != nil {
		return errFailedDecodePrepare
	}

	var proposal ExpectPayload
	if err := rlp.DecodeBytes(expect.Proposal, proposal); err != nil {
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