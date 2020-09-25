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
	"math"
	"reflect"
	"sort"
	"sync"

	"github.com/hyperion-hyn/bls/ffi/go/bls"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/atlas"
)

type defaultValidator struct {
	signer   common.Address // validator's id (address format)
	coinbase common.Address // validator's account address
	pubKey   *bls.PublicKey
}

func (val *defaultValidator) Signer() common.Address {
	return val.signer
}

func (val *defaultValidator) Coinbase() common.Address {
	return val.coinbase
}

func (val *defaultValidator) PublicKey() *bls.PublicKey {
	return val.pubKey
}

func (val *defaultValidator) String() string {
	return "signer: " + val.Signer().String() + ", account: " + val.Coinbase().String()
}

// ----------------------------------------------------------------------------

type defaultSet struct {
	validators atlas.Validators
	policy     atlas.ProposerPolicy

	proposer    atlas.Validator
	validatorMu sync.RWMutex
	selector    atlas.ProposalSelector
}

func newDefaultSet(validators []atlas.Validator, policy atlas.ProposerPolicy) *defaultSet {
	valSet := &defaultSet{}

	valSet.policy = policy
	// init validators
	valSet.validators = make([]atlas.Validator, len(validators))
	for i, validator := range validators {
		valSet.validators[i] = validator
	}
	// sort validator
	//sort.Sort(valSet.validators)
	// init proposer
	if valSet.Size() > 0 {
		valSet.proposer = valSet.GetByIndex(0)
	}
	switch policy {
	case atlas.RoundRobin:
		valSet.selector = roundRobinProposer
	case atlas.Sticky:
		valSet.selector = stickyProposer
	default:
		valSet.selector = roundRobinProposer
	}

	return valSet
}

func (valSet *defaultSet) Size() int {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()
	return len(valSet.validators)
}

func (valSet *defaultSet) List() []atlas.Validator {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()
	return valSet.validators
}

func (valSet *defaultSet) GetByIndex(i uint64) atlas.Validator {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()
	if i < uint64(valSet.Size()) {
		return valSet.validators[i]
	}
	return nil
}

func (valSet *defaultSet) GetBySigner(addr common.Address) (int, atlas.Validator) {
	for i, val := range valSet.List() {
		if addr == val.Signer() {
			return i, val
		}
	}
	return -1, nil
}

func (valSet *defaultSet) GetByCoinbase(addr common.Address) (int, atlas.Validator) {
	for i, val := range valSet.List() {
		if addr == val.Coinbase() {
			return i, val
		}
	}
	return -1, nil
}

func (valSet *defaultSet) GetByPublicKey(pubKey *bls.PublicKey) (int, atlas.Validator) {
	for i, val := range valSet.List() {
		if reflect.DeepEqual(pubKey, val.PublicKey()) == true {
			return i, val
		}
	}
	return -1, nil
}

func (valSet *defaultSet) GetProposer() atlas.Validator {
	return valSet.proposer
}

func (valSet *defaultSet) IsProposer(address common.Address) bool {
	_, val := valSet.GetBySigner(address)
	return reflect.DeepEqual(valSet.GetProposer(), val)
}

func (valSet *defaultSet) CalcProposer(lastProposer common.Address, round uint64) {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()
	valSet.proposer = valSet.selector(valSet, lastProposer, round)
}

func calcSeed(valSet atlas.ValidatorSet, proposer common.Address, round uint64) uint64 {
	offset := 0
	if idx, val := valSet.GetBySigner(proposer); val != nil {
		offset = idx
	}
	return uint64(offset) + round
}

func emptyAddress(addr common.Address) bool {
	return addr == common.Address{}
}

func roundRobinProposer(valSet atlas.ValidatorSet, proposer common.Address, round uint64) atlas.Validator {
	if valSet.Size() == 0 {
		return nil
	}
	seed := uint64(0)
	if emptyAddress(proposer) {
		seed = round
	} else {
		seed = calcSeed(valSet, proposer, round) + 1
	}
	pick := seed % uint64(valSet.Size())
	return valSet.GetByIndex(pick)
}

func stickyProposer(valSet atlas.ValidatorSet, proposer common.Address, round uint64) atlas.Validator {
	if valSet.Size() == 0 {
		return nil
	}
	seed := uint64(0)
	if emptyAddress(proposer) {
		seed = round
	} else {
		seed = calcSeed(valSet, proposer, round)
	}
	pick := seed % uint64(valSet.Size())
	return valSet.GetByIndex(pick)
}

func (valSet *defaultSet) AddValidator(validator atlas.Validator) bool {
	valSet.validatorMu.Lock()
	defer valSet.validatorMu.Unlock()
	for _, v := range valSet.validators {
		if v.Signer() == validator.Signer() {
			return false
		}
	}
	valSet.validators = append(valSet.validators, validator)
	// sort validator
	sort.Sort(valSet.validators)
	return true
}

func (valSet *defaultSet) RemoveValidator(address common.Address) bool {
	valSet.validatorMu.Lock()
	defer valSet.validatorMu.Unlock()

	for i, v := range valSet.validators {
		if v.Signer() == address {
			valSet.validators = append(valSet.validators[:i], valSet.validators[i+1:]...)
			return true
		}
	}
	return false
}

func (valSet *defaultSet) Copy() atlas.ValidatorSet {
	valSet.validatorMu.RLock()
	defer valSet.validatorMu.RUnlock()

	addresses := make([]atlas.Validator, 0, len(valSet.validators))
	for _, v := range valSet.validators {
		addresses = append(addresses, v)
	}
	return NewSet(addresses, valSet.policy)
}

func (valSet *defaultSet) F() int { return int(math.Ceil(float64(valSet.Size())/3)) - 1 }

func (valSet *defaultSet) Policy() atlas.ProposerPolicy { return valSet.policy }

func (valSet *defaultSet) GetPublicKeys() []*bls.PublicKey {
	publicKeys := make([]*bls.PublicKey, 0, len(valSet.validators))
	for _, v := range valSet.validators {
		publicKeys = append(publicKeys, v.PublicKey())
	}
	return publicKeys
}
