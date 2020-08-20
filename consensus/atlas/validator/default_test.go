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
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/hyperion-hyn/bls/ffi/go/bls"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/atlas"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	testAddress  = "70524d664ffe731100208a0154e556f9bb679ae6"
	testAddress2 = "b37866a925bccd69cfa98d43b510f1d23d78a851"
)

func TestValidatorSet(t *testing.T) {
	testNewValidatorSet(t)
	testNormalValSet(t)
	testEmptyValSet(t)
	testStickyProposer(t)
	testAddAndRemoveValidator(t)
}

func testNewValidatorSet(t *testing.T) {
	var validators []atlas.Validator
	const ValCnt = 100

	// Create 100 validators with random addresses
	b := []byte{}
	for i := 0; i < ValCnt; i++ {
		key, _ := crypto.GenerateKey()
		coinbase := crypto.PubkeyToAddress(key.PublicKey)
		blsKey, _ := crypto.GenerateBLSKey()
		val, err := New(coinbase, blsKey.GetPublicKey().Serialize())
		if err != nil {
			t.Errorf("failed to new a validator: %v", err)
		}
		validators = append(validators, val)
		b = append(b, val.Address().Bytes()...)
	}

	// Create ValidatorSet
	valSet := NewSet(validators, atlas.RoundRobin)
	if valSet == nil {
		t.Errorf("the validator byte array cannot be parsed")
		t.FailNow()
	}

	// Check validators sorting: should be in ascending order
	for i := 0; i < ValCnt-1; i++ {
		val := valSet.GetByIndex(uint64(i))
		nextVal := valSet.GetByIndex(uint64(i + 1))
		if strings.Compare(val.Address().String(), nextVal.Address().String()) >= 0 {
			t.Errorf("validator set is not sorted in ascending order")
		}
	}
}

func testNormalValSet(t *testing.T) {
	b1 := common.Hex2Bytes(testAddress)
	b2 := common.Hex2Bytes(testAddress2)
	addr1 := common.BytesToAddress(b1)
	addr2 := common.BytesToAddress(b2)

	secretKey1, _ := crypto.GenerateBLSKey()
	secretKey2, _ := crypto.GenerateBLSKey()
	if strings.Compare(crypto.PubkeyToSigner(secretKey1.GetPublicKey()).String(), crypto.PubkeyToSigner(secretKey2.GetPublicKey()).String()) == 1 {
		secretKey1, secretKey2 = secretKey2, secretKey1
	}

	publicKey1 := secretKey1.GetPublicKey()
	publicKey2 := secretKey2.GetPublicKey()
	// s1 := crypto.PubkeyToSigner(publicKey1)
	s2 := crypto.PubkeyToSigner(publicKey2)
	val1, _ := New(addr1, publicKey1.Serialize())
	val2, _ := New(addr2, publicKey2.Serialize())

	valSet := newDefaultSet([]atlas.Validator{val1, val2}, atlas.RoundRobin)
	if valSet == nil {
		t.Errorf("the format of validator set is invalid")
		t.FailNow()
	}

	// check size
	if size := valSet.Size(); size != 2 {
		t.Errorf("the size of validator set is wrong: have %v, want 2", size)
	}
	// test get by index
	if val := valSet.GetByIndex(uint64(0)); !reflect.DeepEqual(val, val1) {
		t.Errorf("validator mismatch: have %v, want %v", val, val1)
	}
	// test get by invalid index
	if val := valSet.GetByIndex(uint64(2)); val != nil {
		t.Errorf("validator mismatch: have %v, want nil", val)
	}
	// test get by address
	if _, val := valSet.GetByAddress(s2); !reflect.DeepEqual(val, val2) {
		t.Errorf("validator mismatch: have %v, want %v", val, val2)
	}
	// test get by invalid address
	invalidAddr := common.HexToAddress("0x9535b2e7faaba5288511d89341d94a38063a349b")
	if _, val := valSet.GetByAddress(invalidAddr); val != nil {
		t.Errorf("validator mismatch: have %v, want nil", val)
	}
	// test get proposer
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val1) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val1)
	}
	// test calculate proposer
	lastProposer := addr1
	valSet.CalcProposer(lastProposer, uint64(0))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val2) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val2)
	}
	valSet.CalcProposer(lastProposer, uint64(3))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val1) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val1)
	}
	// test empty last proposer
	lastProposer = common.Address{}
	valSet.CalcProposer(lastProposer, uint64(3))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val2) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val2)
	}
}

func testEmptyValSet(t *testing.T) {
	valSet := NewSet([]atlas.Validator{}, atlas.RoundRobin)
	if valSet == nil {
		t.Errorf("validator set should not be nil")
	}
}

func generateSecretKeys(n int) []*bls.SecretKey {
	secretKeys := make([]*bls.SecretKey, n)
	for i := 0; i < n; i++ {
		key, _ := crypto.GenerateBLSKey()
		secretKeys[i] = key
	}
	return secretKeys
}

func sortSecretKeys(secretKeys []*bls.SecretKey) {
	for i := 0; i < len(secretKeys); i++ {
		for j := i; j < len(secretKeys); j++ {
			if bytes.Compare(crypto.PubkeyToSigner(secretKeys[i].GetPublicKey()).Bytes()[:], crypto.PubkeyToSigner(secretKeys[j].GetPublicKey()).Bytes()[:]) > 0 {
				secretKeys[i], secretKeys[j] = secretKeys[j], secretKeys[i]
			}
		}
	}
}

func testAddAndRemoveValidator(t *testing.T) {
	valSet := NewSet([]atlas.Validator{}, atlas.RoundRobin)

	secretKeys := generateSecretKeys(3)
	sortSecretKeys(secretKeys)

	{
		validator, _ := New(common.StringToAddress(string(2)), secretKeys[2].GetPublicKey().Serialize())
		if !valSet.AddValidator(validator) {
			t.Error("the validator should be added")
		}
		if valSet.AddValidator(validator) {
			t.Error("the existing validator should not be added")
		}
	}
	{
		validator, _ := New(common.StringToAddress(string(1)), secretKeys[1].GetPublicKey().Serialize())
		valSet.AddValidator(validator)
	}
	{
		validator, _ := New(common.StringToAddress(string(0)), secretKeys[0].GetPublicKey().Serialize())
		valSet.AddValidator(validator)
	}
	if len(valSet.List()) != 3 {
		t.Error("the size of validator set should be 3")
	}

	for i, v := range valSet.List() {
		expected := crypto.PubkeyToSigner(secretKeys[i].GetPublicKey())
		if v.Address() != expected {
			t.Errorf("the order of validators is wrong: have %v, want %v", v.Address().Hex(), expected.Hex())
		}
	}

	if !valSet.RemoveValidatorBySigner(crypto.PubkeyToSigner(secretKeys[2].GetPublicKey())) {
		t.Error("the validator should be removed")
	}
	if valSet.RemoveValidatorBySigner(crypto.PubkeyToSigner(secretKeys[2].GetPublicKey())) {
		t.Error("the non-existing validator should not be removed")
	}
	if len(valSet.List()) != 2 {
		t.Error("the size of validator set should be 2")
	}
	valSet.RemoveValidatorBySigner(crypto.PubkeyToSigner(secretKeys[1].GetPublicKey()))
	if len(valSet.List()) != 1 {
		t.Error("the size of validator set should be 1")
	}
	valSet.RemoveValidatorBySigner(crypto.PubkeyToSigner(secretKeys[0].GetPublicKey()))
	if len(valSet.List()) != 0 {
		t.Error("the size of validator set should be 0")
	}
}

func testStickyProposer(t *testing.T) {
	b1 := common.Hex2Bytes(testAddress)
	b2 := common.Hex2Bytes(testAddress2)
	addr1 := common.BytesToAddress(b1)
	addr2 := common.BytesToAddress(b2)

	blsKey1, _ := crypto.GenerateBLSKey()
	blsKey2, _ := crypto.GenerateBLSKey()
	if strings.Compare(crypto.PubkeyToSigner(blsKey1.GetPublicKey()).String(), crypto.PubkeyToSigner(blsKey2.GetPublicKey()).String()) == 1 {
		blsKey1, blsKey2 = blsKey2, blsKey1
	}

	s1 := blsKey1.GetPublicKey().Serialize()
	s2 := blsKey2.GetPublicKey().Serialize()
	val1, _ := New(addr1, s1)
	val2, _ := New(addr2, s2)

	valSet := newDefaultSet([]atlas.Validator{val1, val2}, atlas.Sticky)

	// test get proposer
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val1) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val1)
	}
	// test calculate proposer
	lastProposer := addr1
	valSet.CalcProposer(lastProposer, uint64(0))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val1) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val1)
	}

	valSet.CalcProposer(lastProposer, uint64(1))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val2) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val2)
	}
	// test empty last proposer
	lastProposer = common.Address{}
	valSet.CalcProposer(lastProposer, uint64(3))
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val2) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val2)
	}
}
