package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/staking/committee"
	"github.com/ethereum/go-ethereum/staking/types"
	"math/big"
)

func (bc *BlockChain) GetMap3NodePoolAtBlock(blockNum *big.Int) *types.Map3NodePoolStorage {
	return nil
}

func (bc *BlockChain) GetValidatorPoolAtBlock(blockNum *big.Int) *types.ValidatorPoolStorage {
	return nil
}

func (bc *BlockChain) GetMap3NodeAtBlock(blockNum *big.Int, nodeAddress common.Address) (*types.Map3NodeWrapperStorage, error) {
	return nil, nil
}

func (bc *BlockChain) GetValidatorAtBlock(blockNum *big.Int, validatorAddress common.Address) (*types.ValidatorWrapperStorage, error) {
	return nil, nil
}

func (bc *BlockChain) GetMap3NodeAtEpoch(epoch *big.Int, nodeAddress common.Address) (*types.Map3NodeWrapperStorage, error) {
	return nil, nil
}

func (bc *BlockChain) GetValidatorAtEpoch(epoch *big.Int, validatorAddress common.Address) (*types.ValidatorWrapperStorage, error) {
	return nil, nil
}

func (bc *BlockChain) GetCommitteeAtEpoch(epoch *big.Int) (*committee.CommitteeStorage, error) {
	return nil, nil
}