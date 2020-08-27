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

package validator

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/atlas"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hyperion-hyn/bls/ffi/go/bls"
)

func New(publicKey []byte, coinbase common.Address) (atlas.Validator, error) {
	var blsPublicKey bls.PublicKey
	if err := blsPublicKey.Deserialize(publicKey); err != nil {
		return &defaultValidator{}, err
	}

	signer := crypto.PubkeyToSigner(&blsPublicKey)
	return &defaultValidator{
		signer:   signer,
		coinbase: coinbase,
		pubKey:   &blsPublicKey,
	}, nil
}

func NewSet(addrs []atlas.Validator, policy atlas.ProposerPolicy) atlas.ValidatorSet {
	return newDefaultSet(addrs, policy)
}

// Check whether the extraData is presented in prescribed form
func ValidExtraData(extraData []byte) bool {
	return len(extraData) == types.AtlasExtraSeal
}
