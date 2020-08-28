package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"math/big"
)

func (bc *BlockChain) ReadValidatorPoolAtBlock(blockNum *big.Int) (*restaking.Storage_ValidatorPool_, error) {
	header := bc.GetHeaderByNumber(blockNum.Uint64())
	stateDB, err := bc.StateAt(header.Root)
	if err != nil {
		return nil, err
	}
	return stateDB.ValidatorPool(), nil
}

func (bc *BlockChain) ReadValidatorAtBlock(blockNum *big.Int, validatorAddress common.Address) (*restaking.Storage_ValidatorWrapper_, error) {
	header := bc.GetHeaderByNumber(blockNum.Uint64())
	stateDB, err := bc.StateAt(header.Root)
	if err != nil {
		return nil, err
	}
	return stateDB.ValidatorByAddress(validatorAddress)
}

func (bc *BlockChain) ReadValidatorAtEpoch(epoch *big.Int, validatorAddress common.Address) (*restaking.Storage_ValidatorWrapper_, error) {
	blockNum := bc.Config().Atlas.EpochFirstBlock(epoch.Uint64())
	if blockNum > 0 {
		blockNum--
	}
	return bc.ReadValidatorAtBlock(big.NewInt(int64(blockNum)), validatorAddress)
}

func (bc *BlockChain) ReadValidatorAtEpochOrCurrentBlock(epoch *big.Int, validatorAddress common.Address) (*restaking.Storage_ValidatorWrapper_, error) {
	validatorWrapper, err := bc.ReadValidatorAtEpoch(epoch, validatorAddress)
	// if cannot read validator at epoch begin ,read validator at current block stateDB
	if err != nil && err == state.ErrValidatorNotExist {
		return bc.ReadValidatorAtBlock(bc.CurrentBlock().Number(), validatorAddress)
	} else {
		return validatorWrapper, err
	}
}

func (bc *BlockChain) ReadCommitteeAtEpoch(epoch *big.Int) (*restaking.Storage_Committee_, error) {
	blockNum := bc.Config().Atlas.EpochFirstBlock(epoch.Uint64())
	if blockNum > 0 {
		blockNum--
	}
	pool, err := bc.ReadValidatorPoolAtBlock(big.NewInt(int64(blockNum)))
	if err != nil {
		return nil, err
	}
	return pool.Committee(), nil
}
