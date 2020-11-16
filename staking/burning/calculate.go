package burning

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"math/big"
)

const (
	preBurningAmount = 1 // TODO(ATLAS): how much ?
)

type Record struct {
	InternalAmount *big.Int
	ExternalAmount *big.Int
	BlockNum       *big.Int
}

func CalculateInternalBurningAmount(activeNodeCount int, scalingCycleNum int, requireMicrostaking *big.Int) (*big.Int, error) {
	activeNodeCountSqrt, err := common.NewDec(int64(activeNodeCount)).ApproxSqrt()
	if err != nil {
		return nil, errors.Wrap(err, "failed to calculate internal burning amount")
	}
	scalingCycleNumSqrt, err := common.NewDec(int64(scalingCycleNum)).ApproxSqrt()
	if err != nil {
		return nil, errors.Wrap(err, "failed to calculate internal burning amount")
	}
	burningAmount := common.NewDec(50).MulInt(requireMicrostaking).Mul(activeNodeCountSqrt).Quo(scalingCycleNumSqrt).RoundInt()
	return burningAmount, nil
}
