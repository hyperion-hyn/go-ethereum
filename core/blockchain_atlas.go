package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/staking/types/microstaking"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"github.com/pkg/errors"
	"math/big"
)

var (
	errCommitteeNotExist = errors.New("committee not exist")
)

func (bc *BlockChain) ReadValidatorSnapshotAtBlock(blockNum *big.Int, validatorAddress common.Address) (*restaking.Storage_ValidatorWrapper_, error) {
	header := bc.GetHeaderByNumber(blockNum.Uint64())
	stateDB, err := bc.StateAt(header.Root)
	if err != nil {
		return nil, err
	}
	return stateDB.ValidatorSnapshotByAddress(validatorAddress)
}

func (bc *BlockChain) ReadCommitteeAtBlock(blockNum *big.Int) (*restaking.Storage_Committee_, error) {
	header := bc.GetHeaderByNumber(blockNum.Uint64())
	stateDB, err := bc.StateAt(header.Root)
	if err != nil {
		return nil, err
	}
	cmm := stateDB.ValidatorPool().Committee()
	if cmm.Slots().Length() == 0 {
		return nil, errCommitteeNotExist
	}
	return cmm, nil
}

func (bc *BlockChain) ReadMap3NodeSnapshotAtBlock(blockNum *big.Int, map3Address common.Address) (*microstaking.Storage_Map3NodeWrapper_, error) {
	return nil, nil
}