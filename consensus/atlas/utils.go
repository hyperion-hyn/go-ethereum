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
	"github.com/hyperion-hyn/bls/ffi/go/bls"
	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

func RLPHash(v interface{}) (h common.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, v)
	hw.Sum(h[:0])
	return h
}

// GetSignatureAddress gets the signer address from the signature
func GetSignatureAddress(data []byte, sig []byte) (common.Address, error) {
	// 1. Keccak data
	hashData := crypto.Keccak256([]byte(data))
	// 2. Recover public key
	pubkey, err := crypto.SigToPub(hashData, sig)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*pubkey), nil
}

func CheckValidatorSignature(hashdata []byte, sig []byte, pubKey []byte) error {
	// 1. deserialize signature
	var sign bls.Sign
	if err := sign.Deserialize(sig); err != nil {
		log.Debug("Failed to deserialize bls signature", "err", err)
		return err
	}

	// 2. deserialize publicKey
	var publicKey bls.PublicKey
	if err := publicKey.Deserialize(pubKey); err != nil {
		log.Error("Failed to deserialize publicKey", "err", err)
		return err
	}

	// 3. verify signature
	if !sign.VerifyHash(&publicKey, hashdata) {
		log.Error("Failed to verify data")
		return ErrInvalidSignature
	}

	return nil
}
