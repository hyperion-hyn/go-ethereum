package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"math/big"
)

func (cr *fakeChainReader) Engine() consensus.Engine {
	panic("implement me")
}

func (cr *fakeChainReader) ReadValidatorPoolAtBlock(blockNum *big.Int) (*restaking.Storage_ValidatorPool_, error) {
	panic("implement me")
}

func (cr *fakeChainReader) ReadValidatorAtBlock(blockNum *big.Int, validatorAddress common.Address) (*restaking.Storage_ValidatorWrapper_, error) {
	panic("implement me")
}

func (cr *fakeChainReader) ReadValidatorAtEpoch(epoch *big.Int, validatorAddress common.Address) (*restaking.Storage_ValidatorWrapper_, error) {
	panic("implement me")
}

func (cr *fakeChainReader) ReadCommitteeAtEpoch(epoch *big.Int) (*restaking.Storage_Committee_, error) {
	panic("implement me")
}