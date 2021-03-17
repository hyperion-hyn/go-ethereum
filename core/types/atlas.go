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
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

const (
	AtlasMaxValidator = 128 // Maximum number of validators, type of AtlasExtra.Proposer should large enough to hold max index

	AtlasExtraPublicKey = 48 // Fixed number of extra-data bytes reverved for BLS public key
	AtlasExtraSignature = 96 // Fixed number of extra-data bytes reverved for BLS signature
)

var (
	// AtlasDigest represents a Keccak-256 hash of "Atlas"
	// to identify whether the block is from Atlas consensus engine
	AtlasDigest = common.HexToHash("0x3c2f3117cb7ce8fb13f9b40ae62eb87f02e1c4810729d073fc8f6520ebf84a25")

	Map3Account      = common.HexToAddress("0x6a7ad21ff076440e39020e289debdcb309e12c23")
	ValidatorAccount = common.HexToAddress("0x69270f88069d56dc62bd62b0b9f2b302a2b820a8")
)

func GetMaskByteCount(valSetSize int) int {
	return (valSetSize + 7) >> 3
}

type ActiveMap3Info struct {
	Address    string     `json:"address"`
	StartEpoch uint64     `json:"start_epoch"`
	EndEpoch   common.Dec `json:"end_epoch"`
	Commission common.Dec `json:"commission"`
}

type Map3Requirement struct {
	RequireTotal    *big.Int   `json:"requireTotal"`
	RequireSelf     *big.Int   `json:"requireSelf"`
	RequireDelegate *big.Int   `json:"requireDelegate"`
	MinCommission   common.Dec `json:"minCommission"`
	MaxCommission   common.Dec `json:"maxCommission"`
	Map3LockEpoch   common.Dec `json:"map3LockEpoch"`
}
