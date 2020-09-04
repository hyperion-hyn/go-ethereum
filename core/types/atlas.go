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

package types

import (
	"errors"
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	AtlasMaxValidator = 128 // Maximum number of validators, type of AtlasExtra.Proposer should large enough to hold max index

	AtlasExtraVanity    = 32                           // Fixed number of extra-data bytes reserved for validator vanity
	AtlasExtraPublicKey = 48                           // Fixed number of extra-data bytes reverved for BLS public key
	AtlasExtraSignature = 96                           // Fixed number of extra-data bytes reverved for BLS signature
	AtlasExtraMask      = (AtlasMaxValidator + 7) >> 3 // Fixed number of extra-data bytes reserved for BLS signature bitmap
	AtlasExtraProposer  = ((AtlasMaxValidator + ((1 << 8) - 1)) >> 8) << 1
	AtlasExtraSeal      = AtlasExtraSignature + AtlasExtraMask
)

var (
	// AtlasDigest represents a Keccak-256 hash of "Atlas"
	// to identify whether the block is from Atlas consensus engine
	AtlasDigest = common.HexToHash("0x3c2f3117cb7ce8fb13f9b40ae62eb87f02e1c4810729d073fc8f6520ebf84a25")

	// ErrInvalidAtlasHeaderExtra is returned if the length of extra-data is less than 32 bytes
	ErrInvalidAtlasHeaderExtra = errors.New("invalid atlas header extra-data")
)

// Note: If change type of field from fixed-length array to slice, should check all the usage of this struct.
// or it will fail to initialize instance of this struct.
type AtlasExtra struct {
	AggSignature [AtlasExtraSignature]byte // aggregated signature
	AggBitmap    [AtlasExtraMask]byte      // aggregated bitmap
}

// EncodeRLP serializes ist into the Ethereum RLP format.
func (ist *AtlasExtra) EncodeRLP(w io.Writer) error {

	return rlp.Encode(w, []interface{}{
		ist.AggSignature,
		ist.AggBitmap,
	})
}

// DecodeRLP implements rlp.Decoder, and load the atlas fields from a RLP stream.
func (ist *AtlasExtra) DecodeRLP(s *rlp.Stream) error {
	var atlasExtra struct {
		AggSignature [AtlasExtraSignature]byte // aggregated signature
		AggBitmap    [AtlasExtraMask]byte      // aggregated bitmap
	}
	if err := s.Decode(&atlasExtra); err != nil {
		return err
	}
	ist.AggSignature, ist.AggBitmap = atlasExtra.AggSignature, atlasExtra.AggBitmap
	return nil
}

// ExtractAtlasExtra extracts all values of the AtlasExtra from the header. It returns an
// error if the length of the given extra-data is less than 32 bytes or the extra-data can not
// be decoded.
func ExtractAtlasExtra(h *Header) (*AtlasExtra, error) {
	return ExtractAtlasExtraField(h.Extra)
}

func ExtractAtlasExtraField(extra []byte) (*AtlasExtra, error) {
	if len(extra) < AtlasExtraVanity {
		return nil, ErrInvalidAtlasHeaderExtra
	}

	var atlasExtra AtlasExtra
	err := rlp.DecodeBytes(extra[AtlasExtraVanity:], &atlasExtra)
	if err != nil {
		return nil, err
	}
	return &atlasExtra, nil

}

func reset(data []byte) {
	copy(data, make([]byte, len(data)))
}

// AtlasFilteredHeader returns a filtered header which some information (like seal, committed seals)
// are clean to fulfill the Atlas hash rules. It returns nil if the extra-data cannot be
// decoded/encoded by rlp.
func AtlasFilteredHeader(h *Header, keepSeal bool) *Header {
	// Ethash ignore .MixDigest and .Nonce because these two field wiil change after seal.
	// Clique ignore .Extra[crypto.SignatureLength:] because signature in extra will change after seal.

	newHeader := CopyHeader(h)
	atlasExtra, err := ExtractAtlasExtra(newHeader)
	if err != nil {
		return nil
	}

	if !keepSeal {
		reset(atlasExtra.AggSignature[:])
		reset(atlasExtra.AggBitmap[:])
	}

	payload, err := rlp.EncodeToBytes(&atlasExtra)
	if err != nil {
		return nil
	}

	newHeader.Extra = append(newHeader.Extra[:AtlasExtraVanity], payload...)

	return newHeader
}
