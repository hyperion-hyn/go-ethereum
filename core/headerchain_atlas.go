package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"math/big"
)

func (hc *HeaderChain) ReadValidatorSnapshotAtBlock(blockNum *big.Int, validatorAddress common.Address) (*restaking.Storage_ValidatorWrapper_, error) {
	panic("implement me")
}

func (hc *HeaderChain) ReadCommitteeAtBlock(blockNum *big.Int) (*restaking.Storage_Committee_, error) {
	panic("implement me")
}