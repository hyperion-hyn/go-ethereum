package availability

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
	LastCommitBitmap() []byte
}

// ValidatorState is the interface of state.DB
type ValidatorState interface {
	ValidatorByAddress(validatorAddress common.Address) (*staking.ValidatorWrapperStorage, error)
}

type CommitBitmapReader struct {
	Header *types.Header
}

func (c CommitBitmapReader) Number() *big.Int {
	return c.Header.Number
}

func (c CommitBitmapReader) LastCommitBitmap() []byte {
	// TODO(ATLAS): get LastCommitBitmap from parent header extra
	panic("implement me")
}
