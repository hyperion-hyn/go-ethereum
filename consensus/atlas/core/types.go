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
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

type Engine interface {
	Start() error
	Stop() error

	IsProposer(signer common.Address) bool

	ContainProposer() bool

	// verify if a hash is the same as the proposed block in the current pending request
	//
	// this is useful when the engine is currently the proposer
	//
	// pending request is populated right at the preprepare stage so this would give us the earliest verification
	// to avoid any race condition of coming propagated blocks
	IsCurrentProposal(blockHash common.Hash) bool

	Authorize()

	GetLockedHash() common.Hash
}

type State uint64

const (
	StateAcceptRequest State = iota
	StatePreprepared
	StatePrepared
	StateExpected
	StateConfirmed
	StateCommitted
)

func (s State) String() string {
	switch s {
	case StateAcceptRequest:
		return "Accept request"
	case StatePreprepared:
		return "Preprepared"
	case StatePrepared:
		return "Prepared"
	case StateExpected:
		return "Expected"
	case StateConfirmed:
		return "Confirmed"
	case StateCommitted:
		return "Committed"
	default:
		return "Unknown"
	}
}

// Cmp compares s and y and returns:
//   -1 if s is the previous state of y
//    0 if s and y are the same state
//   +1 if s is the next state of y
func (s State) Cmp(y State) int {
	if uint64(s) < uint64(y) {
		return -1
	}
	if uint64(s) > uint64(y) {
		return 1
	}
	return 0
}

const (
	msgPreprepare uint64 = iota
	msgPrepare
	msgExpect
	msgConfirm
	msgCommit
	msgRoundChange
	msgAll
)

type message struct {
	Code          uint64
	Msg           []byte
	Signer        common.Address
	Signature     []byte
	SignerPubKey  []byte
	CommittedSeal []byte
}

// ==============================================
//
// define the functions that needs to be provided for rlp Encoder/Decoder.

// EncodeRLP serializes m into the Ethereum RLP format.
func (m *message) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Code, m.Msg, m.Signer, m.Signature, m.CommittedSeal})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (m *message) DecodeRLP(s *rlp.Stream) error {
	var msg struct {
		Code          uint64
		Msg           []byte
		Signer        common.Address
		Signature     []byte
		CommittedSeal []byte
	}

	if err := s.Decode(&msg); err != nil {
		return err
	}
	m.Code, m.Msg, m.Signer, m.Signature, m.CommittedSeal = msg.Code, msg.Msg, msg.Signer, msg.Signature, msg.CommittedSeal
	return nil
}

// ==============================================
//
// define the functions that needs to be provided for core.

func (m *message) FromPayload(b []byte, preprocessorFn func(*message) error, validateFn func(hash common.Hash, signature []byte, publicKey []byte) error) error {
	// Decode message
	err := rlp.DecodeBytes(b, &m)
	if err != nil {
		return err
	}

	if preprocessorFn != nil {
		if err = preprocessorFn(m); err != nil {
			return err
		}
	}

	// Validate message (on a message without Signature)
	if validateFn != nil {
		var payload []byte
		payload, err = m.PayloadNoSig()
		if err != nil {
			return err
		}

		hash := crypto.Keccak256Hash(payload)
		err := validateFn(hash, m.Signature, m.SignerPubKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *message) Payload() ([]byte, error) {
	return rlp.EncodeToBytes(m)
}

func (m *message) PayloadNoSig() ([]byte, error) {
	// ATLAS(zgx): should remove PayloadNoSig?
	return rlp.EncodeToBytes(&message{
		Code:          m.Code,
		Msg:           m.Msg,
		Signer:        m.Signer,
		SignerPubKey:  []byte{},
		Signature:     []byte{},
		CommittedSeal: []byte{},
	})
}

func (m *message) Decode(val interface{}) error {
	return rlp.DecodeBytes(m.Msg, val)
}

func (m *message) String() string {
	return fmt.Sprintf("{Code: %v, Signer: %v, Signature: %x, PublicKey: %x}", m.Code, m.Signer.String(),
		m.Signature[:math.Min(10, len(m.Signature))],
		m.SignerPubKey[:math.Min(10, len(m.Signature))])
}

// ==============================================
//
// helper functions

func Encode(val interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(val)
}
