package availability

import (
	"github.com/ethereum/go-ethereum/common"
	staking "github.com/ethereum/go-ethereum/staking/types"
	"math/big"
)

// Reader ..
type Reader interface {
	GetValidatorAtEpoch(epoch *big.Int, validatorAddress common.Address) (*staking.ValidatorWrapperStorage, error)
}

// RoundHeader is the interface of block.Header for calculating the BallotResult.
type RoundHeader interface {
	Number() *big.Int
	ShardID() uint32
	LastCommitBitmap() []byte
}

// ValidatorState is the interface of state.DB
type ValidatorState interface {
	ValidatorByAddress(validatorAddress common.Address) (*staking.ValidatorWrapperStorage, error)
}
