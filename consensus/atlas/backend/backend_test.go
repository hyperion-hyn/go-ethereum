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

package backend

import (
	"crypto/ecdsa"
	"math/rand"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/hyperion-hyn/bls/ffi/go/bls"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/atlas"
	"github.com/ethereum/go-ethereum/consensus/atlas/validator"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	bls_cosi "github.com/ethereum/go-ethereum/crypto/bls"
)

func TestSign(t *testing.T) {
	b := newBackend()
	data := []byte("Here is a string....")
	sig, key, _, err := b.Sign(data)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	//Check signature recover
	hashData := crypto.Keccak256([]byte(data))

	var pubKey bls.PublicKey
	if err := pubKey.Deserialize(key); err != nil {
		t.Errorf("failed to deserialize public key: #{key}")
	}

	var sign bls.Sign
	if err := sign.Deserialize(sig); err != nil {
		t.Errorf("failed to deserialize signature: #{sig}")
	}

	if sign.Verify(&pubKey, string(hashData)) == false {
		t.Errorf("failed to verify signature: #{sig}")
	}
}

func TestCheckSignature(t *testing.T) {
	key, _ := generatePrivateKey()
	data := []byte("Here is a string....")
	hashData := crypto.Keccak256([]byte(data))
	sig, _ := crypto.Sign(hashData, key)
	b := newBackend()
	a := getAddress()
	err := b.CheckSignature(data, a, sig)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	a = getInvalidAddress()
	err = b.CheckSignature(data, a, sig)
	if err != errInvalidSignature {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidSignature)
	}
}

func TestCheckValidatorSignature(t *testing.T) {
	_, keys := newTestValidatorSet(5)

	// 1. Positive test: sign with validator's key should succeed
	data := []byte("dummy data")
	hashData := crypto.Keccak256([]byte(data))
	for _, k := range keys {
		// Sign
		sig := k.SignHash(hashData)
		if sig != nil {
			t.Errorf("failed to sign hash data: have nil")
		}
		// CheckValidatorSignature should succeed
		err := atlas.CheckValidatorSignature(data, sig.Serialize(), k.GetPublicKey().Serialize())
		if err != nil {
			t.Errorf("error mismatch: have %v, want nil", err)
		}
	}

	// 2. Negative test: sign with any key other than validator's key should return error
	key, err := crypto.GenerateBLSKey()
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	// Sign
	sig := key.SignHash(hashData)
	if sig == nil {
		t.Errorf("failed to sign hash data: have nil")
	}

	// CheckValidatorSignature should return ErrUnauthorizedAddress
	err = atlas.CheckValidatorSignature(data, sig.Serialize(), key.GetPublicKey().Serialize())
	if err != atlas.ErrUnauthorizedAddress {
		t.Errorf("error mismatch: have %v, want %v", err, atlas.ErrUnauthorizedAddress)
	}
}

func sealWithKeys(privateKeys []*bls.SecretKey, hash common.Hash) (*bls.Sign, *bls.PublicKey, *bls_cosi.Mask, error) {
	var signature bls.Sign
	var publicKey bls.PublicKey
	publicKeys := make([]*bls.PublicKey, len(privateKeys))
	for i, privateKey := range privateKeys {
		publicKeys[i] = privateKey.GetPublicKey()
	}

	bitmap, err := bls_cosi.NewMask(publicKeys, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	for _, privateKey := range privateKeys {
		sign := privateKey.SignHash(hash.Bytes())
		signature.Add(sign)
		publicKey.Add(privateKey.GetPublicKey())
	}

	return &signature, &publicKey, bitmap, nil
}

func generateSecretKeys(n int) ([]*bls.SecretKey, []*bls.PublicKey, error) {
	privateKeys := make([]*bls.SecretKey, n)
	publicKeys := make([]*bls.PublicKey, n)
	for i := 0; i < n; i++ {
		privateKey, err := crypto.GenerateBLSKey()
		if err != nil {
			return nil, nil, err
		}
		privateKeys[i] = privateKey
		publicKeys[i] = privateKey.GetPublicKey()
	}
	return privateKeys, publicKeys, nil
}

func randSetBit(mask *bls_cosi.Mask, n int, v bool) {
	set := make(map[int]bool)
	count := 0
	for {
		if count >= n {
			break
		}

		i := rand.Intn(n)
		_, ok := set[i]
		if ok {
			continue
		}

		mask.SetBit(i, v)
	}
}

func TestCommit(t *testing.T) {
	privateKeys, _, err := generateSecretKeys(88)
	if err != nil {
		t.Errorf("failed to generate %d SecretKeys", 88)
	}

	backend := newBackend()

	commitCh := make(chan *types.Block)
	// Case: it's a proposer, so the backend.commit will receive channel result from backend.Commit function
	testCases := []struct {
		expectedErr   error
		expectedSign  func(block *types.Block) (signature []byte, publicKey []byte, bitmap []byte, err error)
		expectedBlock func() *types.Block
	}{
		{
			// normal case
			nil,
			func(block *types.Block) ([]byte, []byte, []byte, error) {
				hashdata := SealHash(block.Header())
				sign, aggregatedPublicKey, bitmap, err := sealWithKeys(privateKeys, hashdata)
				if err != nil {
					return nil, nil, nil, err
				}

				return sign.Serialize(), aggregatedPublicKey.Serialize(), bitmap.Mask(), nil
			},
			func() *types.Block {
				chain, engine, _ := newBlockChain(1)
				block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
				expectedBlock, _ := engine.updateBlock(engine.chain.GetHeader(block.ParentHash(), block.NumberU64()-1), block)
				return expectedBlock
			},
		},
		{
			// invalid signature
			errInvalidCommittedSeals,
			func(block *types.Block) ([]byte, []byte, []byte, error) {
				return nil, nil, nil, nil
			},
			func() *types.Block {
				chain, engine, _ := newBlockChain(1)
				block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
				expectedBlock, _ := engine.updateBlock(engine.chain.GetHeader(block.ParentHash(), block.NumberU64()-1), block)
				return expectedBlock
			},
		},
	}

	for _, test := range testCases {
		expBlock := test.expectedBlock()
		go func() {
			select {
			case result := <-backend.commitCh:
				commitCh <- result
				return
			}
		}()

		backend.proposedBlockHash = expBlock.Hash()
		signature, publicKey, bitmap, err := test.expectedSign(expBlock)
		if err != nil {
			t.Errorf("failed to sign block: %v", err)
		}
		if err := backend.Commit(expBlock, signature, publicKey, bitmap); err != nil {
			if err != test.expectedErr {
				t.Errorf("error mismatch: have %v, want %v", err, test.expectedErr)
			}
		}

		if test.expectedErr == nil {
			// to avoid race condition is occurred by goroutine
			select {
			case result := <-commitCh:
				if result.Hash() != expBlock.Hash() {
					t.Errorf("hash mismatch: have %v, want %v", result.Hash(), expBlock.Hash())
				}
			case <-time.After(10 * time.Second):
				t.Fatal("timeout")
			}
		}
	}
}

func TestGetProposer(t *testing.T) {
	chain, engine, _ := newBlockChain(1)
	block := makeBlock(chain, engine, chain.Genesis())
	chain.InsertChain(types.Blocks{block})
	expected := engine.GetProposer(1)
	actual := engine.Address()
	if actual != expected {
		t.Errorf("proposer mismatch: have %v, want %v", actual.Hex(), expected.Hex())
	}
}

/**
 * SimpleBackend
 * Private key: bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1
 * Public key: 04a2bfb0f7da9e1b9c0c64e14f87e8fb82eb0144e97c25fe3a977a921041a50976984d18257d2495e7bfd3d4b280220217f429287d25ecdf2b0d7c0f7aae9aa624
 * Address: 0x70524d664ffe731100208a0154e556f9bb679ae6
 */
func getAddress() common.Address {
	return common.HexToAddress("0x70524d664ffe731100208a0154e556f9bb679ae6")
}

func getInvalidAddress() common.Address {
	return common.HexToAddress("0x9535b2e7faaba5288511d89341d94a38063a349b")
}

func generatePrivateKey() (*ecdsa.PrivateKey, error) {
	key := "bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1"
	return crypto.HexToECDSA(key)
}

func newTestValidatorSet(n int) (atlas.ValidatorSet, []*bls.SecretKey) {
	// generate validators
	keys := make(Keys, n)
	addrs := make([]atlas.Validator, n)
	for i := 0; i < n; i++ {
		privateKey, _ := crypto.GenerateBLSKey()
		accountKey, _ := crypto.GenerateKey()
		keys[i] = privateKey
		val, _ := validator.New(privateKey.GetPublicKey().Serialize(), crypto.PubkeyToAddress(accountKey.PublicKey))
		addrs[i] = val
	}
	vset := validator.NewSet(addrs, atlas.RoundRobin)
	sort.Sort(keys) //Keys need to be sorted by its public key address
	return vset, keys
}

type Keys []*bls.SecretKey

func (slice Keys) Len() int {
	return len(slice)
}

func (slice Keys) Less(i, j int) bool {
	return strings.Compare(crypto.PubkeyToSigner(slice[i].GetPublicKey()).String(), crypto.PubkeyToSigner(slice[j].GetPublicKey()).String()) < 0
}

func (slice Keys) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func newBackend() (b *backend) {
	_, b, _ = newBlockChain(4)
	return
}
