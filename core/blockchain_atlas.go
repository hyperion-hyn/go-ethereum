package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"math/big"
)

func (bc *BlockChain) ReadValidatorPoolAtBlock(blockNum *big.Int) *restaking.Storage_ValidatorPool_ {
	// TODO(ATLAS): implement
	return nil
}

func (bc *BlockChain) ReadValidatorAtBlock(blockNum *big.Int, validatorAddress common.Address) (*restaking.Storage_ValidatorWrapper_, error) {
	return nil, nil
}

func (bc *BlockChain) ReadValidatorAtEpoch(epoch *big.Int, validatorAddress common.Address) (*restaking.Storage_ValidatorWrapper_, error) {
	return nil, nil
}

func (bc *BlockChain) ReadCommitteeAtEpoch(epoch *big.Int) (*restaking.Storage_Committee_, error) {
	return nil, nil
}