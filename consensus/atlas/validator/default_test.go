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
	"crypto/ecdsa"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/hyperion-hyn/bls/ffi/go/bls"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/atlas"
	"github.com/ethereum/go-ethereum/crypto"
)

var ()

func TestValidatorSet(t *testing.T) {
	testNewValidatorSet(t)
	testNormalValSet(t)
	testEmptyValSet(t)
	testStickyProposer(t)
	testAddAndRemoveValidator(t)
}

func newValidator() (atlas.Validator, *ecdsa.PrivateKey, *bls.SecretKey, error) {
	privateKey, _ := crypto.GenerateKey()
	secretKey, _ := crypto.GenerateBLSKey()

	peer, err := New(crypto.PubkeyToAddress(privateKey.PublicKey), secretKey.GetPublicKey().Serialize())
	if err != nil {
		return nil, nil, nil, err
	}

	return peer, privateKey, secretKey, nil
}

func generateValidators(count int) ([]atlas.Validator, error) {
	var validators []atlas.Validator

	// Create 100 validators with random addresses
	for i := 0; i < count; i++ {
		val, _, _, err := newValidator()
		if err != nil {
			return nil, err
		}
		validators = append(validators, val)
	}
	return validators, nil
}

func testNewValidatorSet(t *testing.T) {
	const ValCnt = 100
	validators, err := generateValidators(ValCnt)
	if err != nil {
		t.Errorf("failed to new a validator: %v", err)
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
	validators, err := generateValidators(2)
	if err != nil {
		t.Errorf("failed to generate validators: %v", err)
	}
	sort.Slice(validators, func(i, j int) bool {
		return bytes.Compare(validators[i].Address().Bytes(), validators[j].Address().Bytes()) == 1
	})

	if bytes.Compare(validators[0].Address().Bytes(), validators[1].Address().Bytes()) < 1 {
		t.Errorf("validators should be in descending order")
	}
	if err != nil {
		t.Errorf("failed to new a validator: %v", err)
	}
	valSet := newDefaultSet(validators, atlas.RoundRobin)
	if valSet == nil {
		t.Errorf("the format of validator set is invalid")
		t.FailNow()
	}

	// check size
	if size := valSet.Size(); size != 2 {
		t.Errorf("the size of validator set is wrong: have %v, want 2", size)
	}
	// test get by index
	if val := valSet.GetByIndex(uint64(0)); !reflect.DeepEqual(val, validators[1]) {
		t.Errorf("validator mismatch: have %v, want %v", val, validators[1])
	}
	// test get by invalid index
	if val := valSet.GetByIndex(uint64(2)); val != nil {
		t.Errorf("validator mismatch: have %v, want nil", val)
	}
	// test get by address
	if _, val := valSet.GetByAddress(validators[0].Address()); !reflect.DeepEqual(val, validators[0]) {
		t.Errorf("validator mismatch: have %v, want %v", val, validators[0])
	}
	// test get by invalid address
	invalidAddr := common.HexToAddress("0x9535b2e7faaba5288511d89341d94a38063a349b")
	if _, val := valSet.GetByAddress(invalidAddr); val != nil {
		t.Errorf("validator mismatch: have %v, want nil", val)
	}

	val1 := validators[1]
	val2 := validators[0]

	// test get proposer
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val1) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val1)
	}
	// test calculate proposer
	lastProposer := val1.Address()
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

func testAddAndRemoveValidator(t *testing.T) {
	valSet := NewSet([]atlas.Validator{}, atlas.RoundRobin)

	validators, err := generateValidators(3)
	if err != nil {
		t.Errorf("failed to generate validators: %v", err)
	}
	sort.Slice(validators, func(i, j int) bool {
		return bytes.Compare(validators[i].Address().Bytes(), validators[j].Address().Bytes()) == 1
	})

	{
		validator := validators[0]
		if ok := valSet.AddValidator(validator); !ok {
			t.Error("the validator should be added")
		}
		if ok := valSet.AddValidator(validator); ok {
			t.Error("the existing validator should not be added")
		}
	}
	{
		validator := validators[1]
		valSet.AddValidator(validator)
	}
	{
		validator := validators[2]
		valSet.AddValidator(validator)
	}
	if len(valSet.List()) != 3 {
		t.Error("the size of validator set should be 3")
	}

	for i, v := range valSet.List() {
		expected := validators[len(validators)-1-i]
		if v.Address() != expected.Address() {
			t.Errorf("the order of validators is wrong: have %v, want %v", v.Address().Hex(), expected.Address().Hex())
		}
	}

	if ok := valSet.RemoveValidator(validators[2].Address()); !ok {
		t.Error("the validator should be removed")
	}
	if ok := valSet.RemoveValidator(validators[2].Address()); ok {
		t.Error("the non-existing validator should not be removed")
	}
	if len(valSet.List()) != 2 {
		t.Error("the size of validator set should be 2")
	}
	valSet.RemoveValidator(validators[1].Address())
	if len(valSet.List()) != 1 {
		t.Error("the size of validator set should be 1")
	}
	valSet.RemoveValidator(validators[0].Address())
	if len(valSet.List()) != 0 {
		t.Error("the size of validator set should be 0")
	}
}

func testStickyProposer(t *testing.T) {
	validators, err := generateValidators(2)
	if err != nil {
		t.Errorf("failed to generate validators: %v", err)
	}
	sort.Slice(validators, func(i, j int) bool {
		return bytes.Compare(validators[i].Address().Bytes(), validators[j].Address().Bytes()) == 1
	})

	val1 := validators[1]
	val2 := validators[0]

	valSet := newDefaultSet(validators, atlas.Sticky)

	// test get proposer
	if val := valSet.GetProposer(); !reflect.DeepEqual(val, val1) {
		t.Errorf("proposer mismatch: have %v, want %v", val, val1)
	}
	// test calculate proposer
	lastProposer := val1.Address()
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
