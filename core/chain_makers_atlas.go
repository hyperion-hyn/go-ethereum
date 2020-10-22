package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/staking/types/microstaking"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"math/big"
)

func (cr *fakeChainReader) Engine() consensus.Engine {
	panic("implement me")
}

func (cr *fakeChainReader) ReadValidatorSnapshotAtBlock(blockNum *big.Int, validatorAddress common.Address) (*restaking.Storage_ValidatorWrapper_, error) {
	panic("implement me")
}

func (cr *fakeChainReader) ReadCommitteeAtBlock(blockNum *big.Int) (*restaking.Storage_Committee_, error) {
	panic("implement me")
}

func (cr *fakeChainReader) Database() ethdb.Database {
	panic("implement me")
}

func (cr *fakeChainReader) ReadMap3NodeSnapshotAtBlock(blockNum *big.Int, map3Address common.Address) (*microstaking.Storage_Map3NodeWrapper_, error) {
	panic("implement me")
}