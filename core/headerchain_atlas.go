package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"math/big"
)

func (hc *HeaderChain) ReadValidatorPoolAtBlock(blockNum *big.Int) (*restaking.Storage_ValidatorPool_, error) {
	panic("implement me")
}

func (hc *HeaderChain) ReadValidatorAtBlock(blockNum *big.Int, validatorAddress common.Address) (*restaking.Storage_ValidatorWrapper_, error) {
	panic("implement me")
}

func (hc *HeaderChain) ReadValidatorAtEpoch(epoch *big.Int, validatorAddress common.Address) (*restaking.Storage_ValidatorWrapper_, error) {
	panic("implement me")
}

func (hc *HeaderChain) ReadCommitteeAtEpoch(epoch *big.Int) (*restaking.Storage_Committee_, error) {
	panic("implement me")
}