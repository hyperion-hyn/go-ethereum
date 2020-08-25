package light

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"math/big"
)

func (lc *LightChain) ReadValidatorAtEpoch(*big.Int, common.Address) (*restaking.Storage_ValidatorWrapper_, error) {
	panic("implement me")
}

func (lc *LightChain) ReadValidatorAtEpochOrCurrentBlock(*big.Int, common.Address) (*restaking.Storage_ValidatorWrapper_, error) {
	panic("implement me")
}
