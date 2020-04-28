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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

// API is a user facing RPC API to dump Atlas state
type API struct {
	chain consensus.ChainReader
	atlas *backend
}

// BlockSigners is contains who created and who signed a particular block, denoted by its number and hash
type BlockSigners struct {
	Number     uint64
	Hash       common.Hash
	Author     common.Address
	Committers []common.Address
}

// SignerAddress returns the public address that is used to sign block headers in IBFT
func (api *API) SignerAddress() common.Address {
	return api.atlas.Address()
}

// GetSignersFromBlock returns the signers and minter for a given block number, or the
// latest block available if none is specified
func (api *API) GetSignersFromBlock(number *rpc.BlockNumber) (*BlockSigners, error) {
	// Retrieve the requested block number (or current if none requested)
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}

	if header == nil {
		return nil, errUnknownBlock
	}

	return api.signers(header)
}

// GetSignersFromBlockByHash returns the signers and minter for a given block hash
func (api *API) GetSignersFromBlockByHash(hash common.Hash) (*BlockSigners, error) {
	header := api.chain.GetHeaderByHash(hash)
	if header == nil {
		return nil, errUnknownBlock
	}

	return api.signers(header)
}

func (api *API) signers(header *types.Header) (*BlockSigners, error) {
	author, err := api.atlas.Author(header)
	if err != nil {
		return nil, err
	}

	committers, err := api.atlas.Signers(header)
	if err != nil {
		return nil, err
	}

	return &BlockSigners{
		Number:     header.Number.Uint64(),
		Hash:       header.Hash(),
		Author:     author,
		Committers: committers,
	}, nil
}

// GetSnapshot retrieves the state snapshot at a given block.
func (api *API) GetSnapshot(number *rpc.BlockNumber) (*Snapshot, error) {
	// Retrieve the requested block number (or current if none requested)
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	// Ensure we have an actually valid block and return its snapshot
	if header == nil {
		return nil, errUnknownBlock
	}
	return api.atlas.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil)
}

// GetSnapshotAtHash retrieves the state snapshot at a given block.
func (api *API) GetSnapshotAtHash(hash common.Hash) (*Snapshot, error) {
	header := api.chain.GetHeaderByHash(hash)
	if header == nil {
		return nil, errUnknownBlock
	}
	return api.atlas.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil)
}

// GetValidators retrieves the list of authorized validators at the specified block.
func (api *API) GetValidators(number *rpc.BlockNumber) ([]common.Address, error) {
	// Retrieve the requested block number (or current if none requested)
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	// Ensure we have an actually valid block and return the validators from its snapshot
	if header == nil {
		return nil, errUnknownBlock
	}
	snap, err := api.atlas.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil)
	if err != nil {
		return nil, err
	}
	return snap.validators(), nil
}

// GetValidatorsAtHash retrieves the state snapshot at a given block.
func (api *API) GetValidatorsAtHash(hash common.Hash) ([]common.Address, error) {
	header := api.chain.GetHeaderByHash(hash)
	if header == nil {
		return nil, errUnknownBlock
	}
	snap, err := api.atlas.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil)
	if err != nil {
		return nil, err
	}
	return snap.validators(), nil
}

// Candidates returns the current candidates the node tries to uphold and vote on.
func (api *API) Candidates() map[common.Address]ValidatorProposal {
	api.atlas.candidatesLock.RLock()
	defer api.atlas.candidatesLock.RUnlock()

	proposals := make(map[common.Address]ValidatorProposal)
	for address, proposal := range api.atlas.candidates {
		proposals[address] = proposal
	}
	return proposals
}

// Propose injects a new authorization candidate that the validator will attempt to
// push through.
func (api *API) Propose(signer common.Address, coinbase common.Address, auth bool) {
	api.atlas.candidatesLock.Lock()
	defer api.atlas.candidatesLock.Unlock()

	api.atlas.candidates[signer] = ValidatorProposal{Signer: signer, Coinbase:coinbase, Vote:auth}
}

// Discard drops a currently running candidate, stopping the validator from casting
// further votes (either for or against).
func (api *API) Discard(signer common.Address) {
	api.atlas.candidatesLock.Lock()
	defer api.atlas.candidatesLock.Unlock()

	delete(api.atlas.candidates, signer)
}
