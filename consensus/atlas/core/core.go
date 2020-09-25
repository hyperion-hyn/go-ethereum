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
	"bytes"
	"errors"
	"fmt"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/hyperion-hyn/bls/ffi/go/bls"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/atlas"
	"github.com/ethereum/go-ethereum/crypto"
	bls_cosi "github.com/ethereum/go-ethereum/crypto/bls"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/rlp"
)

// New creates an Atlas consensus core
func New(backend atlas.Backend, config *atlas.Config) Engine {
	r := metrics.NewRegistry()
	c := &core{
		config:             config,
		address:            backend.Address(),
		state:              StateAcceptRequest,
		handlerWg:          new(sync.WaitGroup),
		logger:             log.New("annotation", backend.Annotation(), "signer", backend.Signer()),
		backend:            backend,
		backlogs:           make(map[common.Address]*prque.Prque),
		backlogsMu:         new(sync.Mutex),
		pendingRequests:    prque.New(),
		pendingRequestsMu:  new(sync.Mutex),
		consensusTimestamp: time.Time{},
		roundMeter:         metrics.NewMeter(),
		sequenceMeter:      metrics.NewMeter(),
		consensusTimer:     metrics.NewTimer(),
	}

	r.Register("consensus/atlas/core/round", c.roundMeter)
	r.Register("consensus/atlas/core/sequence", c.sequenceMeter)
	r.Register("consensus/atlas/core/consensus", c.consensusTimer)

	c.validateHashFn = c.checkValidatorSignature
	return c
}

// ----------------------------------------------------------------------------

type core struct {
	config  *atlas.Config
	address common.Address // owner's address
	state   State
	logger  log.Logger

	backend               atlas.Backend
	events                *event.TypeMuxSubscription
	finalCommittedSub     *event.TypeMuxSubscription
	timeoutSub            *event.TypeMuxSubscription
	futurePreprepareTimer *time.Timer

	valSet                atlas.ValidatorSet
	waitingForRoundChange bool
	validateHashFn        func(hash common.Hash, sig []byte, pubKey []byte) error

	backlogs   map[common.Address]*prque.Prque
	backlogsMu *sync.Mutex

	current   *roundState
	handlerWg *sync.WaitGroup

	roundChangeSet   *roundChangeSet
	roundChangeTimer *time.Timer

	pendingRequests   *prque.Prque
	pendingRequestsMu *sync.Mutex

	consensusTimestamp time.Time
	// the meter to record the round change rate
	roundMeter metrics.Meter
	// the meter to record the sequence update rate
	sequenceMeter metrics.Meter
	// the timer to record consensus duration (from accepting a preprepare to final committed stage)
	consensusTimer metrics.Timer
}

func (c *core) finalizeMessage(msg *message) ([]byte, error) {
	var err error
	// Add sender address
	msg.Signer = c.backend.Signer()

	// Sign message
	data, err := msg.PayloadNoSig()
	if err != nil {
		return nil, err
	}
	hash := crypto.Keccak256Hash(data)
	msg.Signature, msg.SignerPubKey, _, err = c.backend.SignHash(hash)
	if err != nil {
		return nil, err
	}

	if err := c.validateHashFn(hash, msg.Signature, msg.SignerPubKey); err != nil {
		c.logger.Error("Validate failed after Sign")
		return nil, err
	}

	// Convert to payload
	payload, err := msg.Payload()
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (c *core) broadcast(msg *message) {
	logger := c.logger.New("state", c.state)

	payload, err := c.finalizeMessage(msg)
	if err != nil {
		logger.Error("Failed to finalize message", "msg", msg, "err", err)
		return
	}

	// Broadcast payload
	if err = c.backend.Broadcast(c.valSet, payload); err != nil {
		logger.Error("Failed to broadcast message", "msg", msg, "err", err)
		return
	}
}

func (c *core) currentView() *atlas.View {
	return &atlas.View{
		Sequence: new(big.Int).Set(c.current.Sequence()),
		Round:    new(big.Int).Set(c.current.Round()),
	}
}

func (c *core) IsProposer() bool {
	v := c.valSet
	if v == nil {
		return false
	}
	return v.IsProposer(c.backend.Signer())
}

func (c *core) IsCurrentProposal(blockHash common.Hash) bool {
	return c.current != nil && c.current.pendingRequest != nil && c.current.pendingRequest.Proposal.SealHash(c.backend) == blockHash
}

func (c *core) commit() {
	c.setState(StateCommitted)

	proposal := c.current.Proposal()
	if proposal != nil {
		committedSignature := c.current.aggregatedConfirmSig.Serialize()
		committedBitmap := c.current.confirmBitmap.Mask()

		if err := c.backend.Commit(proposal, committedSignature, committedBitmap); err != nil {
			c.current.UnlockHash() //Unlock block when insertion fails
			c.sendNextRoundChange()
			return
		}
	}
}

// startNewRound starts a new round. if round equals to 0, it means to starts a new sequence
func (c *core) startNewRound(round *big.Int) {
	var logger log.Logger
	if c.current == nil {
		logger = c.logger.New("old_round", -1, "old_seq", 0)
	} else {
		logger = c.logger.New("old_round", c.current.Round(), "old_seq", c.current.Sequence())
	}

	roundChange := false
	// Try to get last proposal
	lastProposal, lastProposer := c.backend.LastProposal()
	if c.current == nil {
		logger.Trace("Start to the initial round")
	} else if lastProposal.Number().Cmp(c.current.Sequence()) >= 0 {
		diff := new(big.Int).Sub(lastProposal.Number(), c.current.Sequence())
		c.sequenceMeter.Mark(new(big.Int).Add(diff, common.Big1).Int64())

		if !c.consensusTimestamp.IsZero() {
			c.consensusTimer.UpdateSince(c.consensusTimestamp)
			c.consensusTimestamp = time.Time{}
		}
		logger.Trace("Catch up latest proposal", "number", lastProposal.Number().Uint64(), "hash", lastProposal.Hash())
	} else if lastProposal.Number().Cmp(big.NewInt(c.current.Sequence().Int64()-1)) == 0 {
		if round.Cmp(common.Big0) == 0 {
			// same seq and round, don't need to start new round
			return
		} else if round.Cmp(c.current.Round()) < 0 {
			logger.Warn("New round should not be smaller than current round", "seq", lastProposal.Number().Int64(), "new_round", round, "old_round", c.current.Round())
			return
		}
		roundChange = true
	} else {
		logger.Warn("New sequence should be larger than current sequence", "new_seq", lastProposal.Number().Int64())
		return
	}

	var newView *atlas.View
	if roundChange {
		newView = &atlas.View{
			Sequence: new(big.Int).Set(c.current.Sequence()),
			Round:    new(big.Int).Set(round),
		}
	} else {
		newView = &atlas.View{
			Sequence: new(big.Int).Add(lastProposal.Number(), common.Big1),
			Round:    new(big.Int),
		}
		c.valSet = c.backend.Validators(lastProposal)
	}

	// Update logger
	logger = logger.New("old_proposer", c.valSet.GetProposer())
	// Clear invalid ROUND CHANGE messages
	c.roundChangeSet = newRoundChangeSet(c.valSet)
	// New snapshot for new round
	c.updateRoundState(newView, c.valSet, roundChange)
	// Calculate new proposer
	c.valSet.CalcProposer(lastProposer, newView.Round.Uint64())
	c.waitingForRoundChange = false
	c.setState(StateAcceptRequest)
	if roundChange && c.IsProposer() && c.current != nil {
		// If it is locked, propose the old proposal
		// If we have pending request, propose pending request
		if c.current.IsHashLocked() {
			r := &atlas.Request{
				Proposal: c.current.Proposal(), //c.current.Proposal would be the locked proposal by previous proposer, see updateRoundState
			}
			c.sendPreprepare(r)
		} else if c.current.pendingRequest != nil {
			c.sendPreprepare(c.current.pendingRequest)
		}
	}
	c.newRoundChangeTimer()

	logger.Debug("New round", "new_round", newView.Round, "new_seq", newView.Sequence, "new_proposer", c.valSet.GetProposer(), "valSet", c.valSet.List(), "size", c.valSet.Size(), "IsProposer", c.IsProposer())
}

func (c *core) catchUpRound(view *atlas.View) {
	logger := c.logger.New("old_round", c.current.Round(), "old_seq", c.current.Sequence(), "old_proposer", c.valSet.GetProposer())

	if view.Round.Cmp(c.current.Round()) > 0 {
		c.roundMeter.Mark(new(big.Int).Sub(view.Round, c.current.Round()).Int64())
	}
	c.waitingForRoundChange = true

	// Need to keep block locked for round catching up
	c.updateRoundState(view, c.valSet, true)
	c.roundChangeSet.Clear(view.Round)
	c.newRoundChangeTimer()

	logger.Trace("Catch up round", "new_round", view.Round, "new_seq", view.Sequence, "new_proposer", c.valSet)
}

// updateRoundState updates round state by checking if locking block is necessary
func (c *core) updateRoundState(view *atlas.View, validatorSet atlas.ValidatorSet, roundChange bool) {
	// Lock only if both roundChange is true and it is locked
	if roundChange && c.current != nil {
		if c.current.IsHashLocked() {
			c.current = newRoundState(view, validatorSet, c.current.GetLockedHash(), c.current.Preprepare, c.current.pendingRequest, c.backend.HasBadProposal)
		} else {
			c.current = newRoundState(view, validatorSet, common.Hash{}, nil, c.current.pendingRequest, c.backend.HasBadProposal)
		}
	} else {
		c.current = newRoundState(view, validatorSet, common.Hash{}, nil, nil, c.backend.HasBadProposal)
	}
}

func (c *core) setState(state State) {
	if c.state != state {
		c.state = state
	}
	if state == StateAcceptRequest {
		c.processPendingRequests()
	}
	c.processBacklog()
}

func (c *core) Signer() common.Address {
	return c.backend.Signer()
}

func (c *core) stopFuturePreprepareTimer() {
	if c.futurePreprepareTimer != nil {
		c.futurePreprepareTimer.Stop()
	}
}

func (c *core) stopTimer() {
	c.stopFuturePreprepareTimer()
	if c.roundChangeTimer != nil {
		c.roundChangeTimer.Stop()
	}
}

func (c *core) newRoundChangeTimer() {
	c.stopTimer()

	// set timeout based on the round number
	timeout := time.Duration(c.config.RequestTimeout) * time.Millisecond
	round := c.current.Round().Uint64()
	if round > 0 {
		timeout += time.Duration(math.Pow(2, float64(round))) * time.Second
	}

	c.roundChangeTimer = time.AfterFunc(timeout, func() {
		c.logger.Debug("timeout", "timeout", timeout, "round", c.current.round.Uint64())
		c.sendEvent(timeoutEvent{})
	})
}

func (c *core) checkValidatorSignature(hash common.Hash, sig []byte, pubKey []byte) error {
	return atlas.CheckValidatorSignature(hash.Bytes(), sig, pubKey)
}

func (c *core) QuorumSize() int {
	c.logger.Trace("Confirmation Formula used ceil(2N/3)")
	return int(math.Floor(float64(2*c.valSet.Size())/3)) + 1
}

// PrepareCommittedSeal returns a committed seal for the given hash
func PrepareCommittedSeal(hash common.Hash) []byte {
	var buf bytes.Buffer
	buf.Write(hash.Bytes())
	buf.Write([]byte{byte(msgCommit)})
	return buf.Bytes()
}

func (c *core) SignSubject(subject *atlas.Subject) (*atlas.Subject, error) {
	signedSubject, err := atlas.SignSubject(subject, c.valSet, func(hash common.Hash) (signature []byte, publicKey []byte, mask []byte, err error) {
		signature, publicKey, mask, err = c.backend.SignHash(hash)
		if err != nil {
			return nil, nil, nil, err
		}
		return signature, publicKey, mask, nil
	})
	return signedSubject, err
}

func (c *core) AssembleSignedSubject(subject *atlas.Subject) (*atlas.Subject, error) {
	switch c.state {
	case StatePrepared, StateConfirmed:
		var val *atlas.SignPayload
		switch c.state {
		case StatePrepared:
			val = &atlas.SignPayload{
				Signature: c.current.aggregatedPrepareSig.Serialize(),
				Mask:      c.current.prepareBitmap.Mask(),
			}
		case StateConfirmed:
			val = &atlas.SignPayload{
				Signature: c.current.aggregatedConfirmSig.Serialize(),
				Mask:      c.current.confirmBitmap.Mask(),
			}
		}
		payload, err := rlp.EncodeToBytes(val)
		if err != nil {
			return nil, err
		}
		hash := crypto.Keccak256Hash(payload)
		signature, _, _, err := c.backend.SignHash(hash)
		if err != nil {
			return nil, errFailedSignData
		}
		subject.Payload = payload
		subject.Signature = signature
		return subject, nil
	default:
		return nil, errors.New(fmt.Sprintf("invalid state: %v", c.current))
	}
}

func (c *core) verifySignPayload(subject *atlas.Subject, validatorSet atlas.ValidatorSet) (*atlas.SignPayload, error) {
	var signPayload atlas.SignPayload
	if err := rlp.DecodeBytes(subject.Payload, &signPayload); err != nil {
		return nil, errFailedDecodeSignPayload
	}

	bitmap, _ := bls_cosi.NewMask(c.valSet.GetPublicKeys(), nil)
	if err := bitmap.SetMask(signPayload.Mask); err != nil {
		c.logger.Error("Failed to SetMask", "err", err)
		return nil, err
	}

	publicKey := bitmap.AggregatePublic

	if err := c.checkValidatorSignature(subject.Digest, signPayload.Signature, publicKey.Serialize()); err != nil {
		return nil, err
	}

	return &signPayload, nil
}

func (c *core) getValidatorPublicKey(signer common.Address, valSet atlas.ValidatorSet) (*bls.PublicKey, error) {
	_, validator := valSet.GetBySigner(signer)
	if validator == nil {
		return nil, errInvalidSigner
	}

	var pubKey *bls.PublicKey = validator.PublicKey()
	if pubKey == nil {
		return nil, errInvalidSigner
	}
	return pubKey, nil
}

func (c *core) Authorize() {
	c.logger = log.New("annotation", c.backend.Annotation(), "signer", c.backend.Signer())
}

func (c *core) GetLockedHash() common.Hash {
	if c.current != nil {
		return c.current.GetLockedHash()
	} else {
		return common.Hash{}
	}
}
