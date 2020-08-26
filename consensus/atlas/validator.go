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
	"strings"

	"github.com/hyperion-hyn/bls/ffi/go/bls"

	"github.com/ethereum/go-ethereum/common"
)

type Validator interface {
	// Address returns validator's id
	Address() common.Address

	// Coinbase returns validator's coinbase address
	Coinbase() common.Address

	PublicKey() *bls.PublicKey

	// String representation of Validator
	String() string
}

// ----------------------------------------------------------------------------

type Validators []Validator

func (slice Validators) Len() int {
	return len(slice)
}

func (slice Validators) Less(i, j int) bool {
	return strings.Compare(slice[i].Address().String(), slice[j].Address().String()) < 0
}

func (slice Validators) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// ----------------------------------------------------------------------------

type ValidatorSet interface {
	// Calculate the proposer
	CalcProposer(lastProposer common.Address, round uint64)
	// Return the validator size
	Size() int
	// Return the validator array
	List() []Validator
	// Get validator by index
	GetByIndex(i uint64) Validator
	// Get validator by given address
	GetByAddress(addr common.Address) (int, Validator)
	// Get validator by public key
	GetByPublicKey(pubKey *bls.PublicKey) (int, Validator)
	// Get current proposer
	GetProposer() Validator
	// Check whether the validator with given address is a proposer
	IsProposer(address common.Address) bool
	// Add validator
	AddValidator(validator Validator) bool
	// Remove validator
	RemoveValidator(address common.Address) bool
	// Copy validator set
	Copy() ValidatorSet
	// Get the maximum number of faulty nodes
	F() int
	// Get proposer policy
	Policy() ProposerPolicy
	// Get public keys
	GetPublicKeys() []*bls.PublicKey
}

// ----------------------------------------------------------------------------

type ProposalSelector func(ValidatorSet, common.Address, uint64) Validator
