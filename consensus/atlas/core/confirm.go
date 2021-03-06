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
	"fmt"
	"github.com/ethereum/go-ethereum/consensus/atlas"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/hyperion-hyn/bls/ffi/go/bls"
	"time"
)

func (c *core) sendConfirm() {
	logger := c.logger.New("state", c.state)

	signers := c.Signer()
	for _, signer := range signers {
		sub, err := c.SignSubject(signer, c.current.Subject())
		if err != nil {
			logger.Error("Failed to sign", "view", c.currentView(), "err", err)
			continue
		}

		encodedSubject, err := Encode(sub)
		if err != nil {
			logger.Error("Failed to encode", "subject", sub, "err", err)
			continue
		}
		c.broadcast(signer, &message{
			Code: msgConfirm,
			Msg:  encodedSubject,
		})
	}
}

func (c *core) handleConfirm(msg *message, src atlas.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	// only proposer can/need handle confirm
	if !c.ContainProposer() {
		return nil
	}

	// Decode PREPARE message
	var confirm atlas.Subject
	err := msg.Decode(&confirm)
	if err != nil {
		return errFailedDecodeConfirm
	}

	if err := c.checkMessage(msgConfirm, confirm.View); err != nil {
		return err
	}

	// If it is locked, it can only process on the locked block.
	// Passing verifyPrepare and checkMessage implies it is processing on the locked block since it was verified in the Preprepared state.
	if err := c.verifyConfirm(&confirm, src); err != nil {
		return err
	}

	// Signer should be in the validator set
	if _, v := c.valSet.GetBySigner(src.Signer()); v == nil {
		return errNotFromCommittee
	}

	if err := c.acceptConfirm(msg, src); err != nil {
		logger.Error("acceptConfirm", "err", err)
		return err
	}

	// Change to Confirm state if we've received enough Expect messages
	// and we are in Expected state to provent validator send CONFIRM messages
	// to pass by Prepared and Expected state
	if c.current.GetConfirmSize() >= c.QuorumSize() && c.state.Cmp(StateExpected) == 0 {
		duration := 2 * time.Second
		estimateCommitElapsed := time.Duration(500) * time.Millisecond
		blockPeriod := time.Duration(int64(c.config.BlockPeriod)) * time.Second
		if !c.current.WaitConfirm() {
			c.current.SetWaitConfirm(true)
			//if blockPeriod-time.Now().Sub(c.prePrepareTimestamp) < duration+commitElapsed {
			proposalTime := time.Unix(int64(c.current.Proposal().Header().Time), 0)
			if blockPeriod-time.Now().Sub(proposalTime) < duration+estimateCommitElapsed {
				duration = blockPeriod - time.Now().Sub(proposalTime) - estimateCommitElapsed
			}
			c.logger.Debug(fmt.Sprintf("sleep %v second to accept more message", duration))
			time.AfterFunc(duration, func() {
				c.logger.Debug(fmt.Sprintf("%v second passed, send commit message", duration))
				c.setState(StateConfirmed)
				c.sendCommit()
			})
		} else {
			c.logger.Debug("accept more message")
		}
	}
	return nil
}

// verifyPrepare verifies if the received CONFIRM message is equivalent to our subject
func (c *core) verifyConfirm(confirm *atlas.Subject, src atlas.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	sub := c.current.Subject()
	if !atlas.IsConsistentSubject(confirm, sub) {
		logger.Warn("Inconsistent subjects between CONFIRM and proposal", "expected", sub, "got", confirm)
		return errInconsistentSubject
	}

	return nil
}

func (c *core) acceptConfirm(msg *message, src atlas.Validator) error {
	// TODO(zgx): should reset if error occure
	logger := c.logger.New("from", src, "state", c.state)

	// only in the Expect state can accept CONFIRM signature
	if c.state != StateExpected {
		return nil
	}

	var confirm atlas.Subject
	if err := msg.Decode(&confirm); err != nil {
		return errFailedDecodeConfirm
	}

	var signPayload atlas.SignPayload
	if err := rlp.DecodeBytes(confirm.Payload, &signPayload); err != nil {
		return errFailedDecodeConfirm
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

	err = c.checkValidatorSignature(confirm.Digest, signPayload.Signature, pubKey.Serialize())
	if err != nil {
		logger.Error("Failed to verify signature with signer's public key confirm", "signature", signPayload.Signature[:10], "publicKey", pubKey.Serialize()[:10])
		return err
	}

	enabled, err := c.current.confirmBitmap.KeyEnabled(pubKey)
	if err != nil {
		return err
	} else if enabled == true {
		return errDuplicateMessage
	}

	if err := c.current.confirmBitmap.SetKey(pubKey, true); err == nil {
		c.current.aggregatedConfirmSig.Add(&sign)
	}

	return nil
}
