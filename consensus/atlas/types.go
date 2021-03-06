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

package atlas

import (
	"fmt"
	"io"
	"math/big"
	"reflect"

	"github.com/hyperion-hyn/bls/ffi/go/bls"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bls_cosi "github.com/ethereum/go-ethereum/crypto/bls"
	"github.com/ethereum/go-ethereum/rlp"
)

// Proposal supports retrieving height and serialized block to be used during Atlas consensus.
type Proposal interface {
	Header() *types.Header

	// Number retrieves the sequence number of this proposal.
	Number() *big.Int

	// Hash retrieves the hash of this proposal.
	Hash() common.Hash

	// SealHash retrieves the sealhash of this proposal.
	SealHash(sealer types.Sealer) common.Hash

	EncodeRLP(w io.Writer) error

	DecodeRLP(s *rlp.Stream) error
}

type Request struct {
	Proposal Proposal
}

// View includes a round number and a sequence number.
// Sequence is the block number we'd like to commit.
// Each round has a number and is composed by 3 steps: preprepare, prepare and commit.
//
// If the given block is not accepted by validators, a round change will occur
// and the validators start a new round with round+1.
type View struct {
	Round    *big.Int
	Sequence *big.Int
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (v *View) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{v.Round, v.Sequence})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (v *View) DecodeRLP(s *rlp.Stream) error {
	var view struct {
		Round    *big.Int
		Sequence *big.Int
	}

	if err := s.Decode(&view); err != nil {
		return err
	}
	v.Round, v.Sequence = view.Round, view.Sequence
	return nil
}

func (v *View) String() string {
	return fmt.Sprintf("{Round: %d, Sequence: %d}", v.Round.Uint64(), v.Sequence.Uint64())
}

// Cmp compares v and y and returns:
//   -1 if v <  y
//    0 if v == y
//   +1 if v >  y
func (v *View) Cmp(y *View) int {
	if v.Sequence.Cmp(y.Sequence) != 0 {
		return v.Sequence.Cmp(y.Sequence)
	}
	if v.Round.Cmp(y.Round) != 0 {
		return v.Round.Cmp(y.Round)
	}
	return 0
}

type Preprepare struct {
	View      *View
	Proposal  Proposal
	Signature []byte
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (b *Preprepare) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{b.View, b.Proposal, b.Signature})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (b *Preprepare) DecodeRLP(s *rlp.Stream) error {
	var preprepare struct {
		View      *View
		Proposal  *types.Block
		Signature []byte
	}

	if err := s.Decode(&preprepare); err != nil {
		return err
	}
	b.View, b.Proposal, b.Signature = preprepare.View, preprepare.Proposal, preprepare.Signature

	return nil
}

type Subject struct {
	View      *View
	Digest    common.Hash
	Payload   []byte
	Signature []byte // Signature of Payload
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (b *Subject) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{b.View, b.Digest, b.Payload, b.Signature})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (b *Subject) DecodeRLP(s *rlp.Stream) error {
	var subject struct {
		View      *View
		Digest    common.Hash
		Payload   []byte
		Signature []byte // Signature of Payload
	}

	if err := s.Decode(&subject); err != nil {
		return err
	}
	b.View, b.Digest, b.Payload, b.Signature = subject.View, subject.Digest, subject.Payload, subject.Signature
	return nil
}

func (b Subject) String() string {
	return fmt.Sprintf("{View: %v, Digest: %v}", b.View, b.Digest.String())
}

type SignPayload struct {
	Signature []byte
	Mask      []byte
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (b *SignPayload) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{b.Signature, b.Mask})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (b *SignPayload) DecodeRLP(s *rlp.Stream) error {
	var obj struct {
		Signature []byte
		Mask      []byte
	}

	if err := s.Decode(&obj); err != nil {
		return err
	}
	b.Signature, b.Mask = obj.Signature, obj.Mask
	return nil
}

func (b SignPayload) String() string {
	return fmt.Sprintf("{ Signature: %x, PublicKey: %x}", b.Signature[:10], b.Mask)
}

type SignHashFn func(hash common.Hash) (signature []byte, publicKey []byte, mask []byte, err error)

func SignSubject(subject *Subject, valset ValidatorSet, signFn SignHashFn) (*Subject, error) {
	hash := subject.Digest
	signature, publicKey, _, err := signFn(hash)
	if err != nil {
		return nil, err
	}

	var pubKey bls.PublicKey
	if err := pubKey.Deserialize(publicKey); err != nil {
		return nil, err
	}

	bitmap, err := bls_cosi.NewMask(valset.GetPublicKeys(), nil)
	if err := bitmap.SetKey(&pubKey, true); err != nil {
		return nil, err
	}
	val := SignPayload{
		Signature: signature,
		Mask:      bitmap.Mask(),
	}

	payload, err := rlp.EncodeToBytes(&val)
	if err != nil {
		return nil, err
	}

	var retval Subject
	{
		var data []byte
		var err error
		if data, err = rlp.EncodeToBytes(subject); err != nil {
			return nil, err
		}
		if err = rlp.DecodeBytes(data, &retval); err != nil {
			return nil, err
		}
	}
	retval.Payload = payload

	return &retval, nil
}

func IsConsistentSubject(a *Subject, b *Subject) bool {
	return reflect.DeepEqual(a.View, b.View) && reflect.DeepEqual(a.Digest, b.Digest)
}
