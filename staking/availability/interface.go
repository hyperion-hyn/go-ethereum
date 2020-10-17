package availability

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"math/big"
)

// RoundHeader is the interface of block.Header for calculating the BallotResult.
type RoundHeader interface {
	Number() *big.Int
	LastCommitBitmap() ([]byte, error)
}

// ValidatorState is the interface of state.DB
type ValidatorState interface {
	ValidatorByAddress(validatorAddress common.Address) (*restaking.Storage_ValidatorWrapper_, error)
	ValidatorSnapshotByAddress(validatorAddress common.Address) (*restaking.Storage_ValidatorWrapper_, error)
}

type CommitBitmapReader struct {
	Header *types.Header
}

func (c CommitBitmapReader) Number() *big.Int {
	return c.Header.Number
}

func (c CommitBitmapReader) LastCommitBitmap() ([]byte, error) {
	return c.Header.Bitmap, nil
}
