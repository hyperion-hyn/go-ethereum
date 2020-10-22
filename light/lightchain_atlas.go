package light

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/staking/types/microstaking"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"math/big"
)

func (lc *LightChain) ReadValidatorSnapshotAtBlock(blockNum *big.Int, validatorAddress common.Address) (*restaking.Storage_ValidatorWrapper_, error) {
	panic("implement me")
}

func (lc *LightChain) ReadCommitteeAtBlock(blockNum *big.Int) (*restaking.Storage_Committee_, error) {
	panic("implement me")
}

func (lc *LightChain) ReadMap3NodeSnapshotAtBlock(blockNum *big.Int, map3Address common.Address) (*microstaking.Storage_Map3NodeWrapper_, error) {
	panic("implement me")
}