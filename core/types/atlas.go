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
	"github.com/ethereum/go-ethereum/common"
)

const (
	AtlasMaxValidator = 128 // Maximum number of validators, type of AtlasExtra.Proposer should large enough to hold max index

	AtlasExtraVanity    = 32 // Fixed number of extra-data bytes reserved for validator vanity
	AtlasExtraPublicKey = 48 // Fixed number of extra-data bytes reverved for BLS public key
	AtlasExtraSignature = 96 // Fixed number of extra-data bytes reverved for BLS signature
)

var (
	// AtlasDigest represents a Keccak-256 hash of "Atlas"
	// to identify whether the block is from Atlas consensus engine
	AtlasDigest = common.HexToHash("0x3c2f3117cb7ce8fb13f9b40ae62eb87f02e1c4810729d073fc8f6520ebf84a25")
)

func GetMaskByteCount(valSetSize int) int {
	return (valSetSize + 7) >> 3
}
