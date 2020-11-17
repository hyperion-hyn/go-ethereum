package burning

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/atlas"
	"github.com/pkg/errors"
	"math/big"
)

type Receipt struct {
	Hash           *common.Hash `json:"hash" rlp:"-"`
	InternalAmount *big.Int     `json:"internal-amount"`
	ExternalAmount *big.Int     `json:"external-amount"`
	BlockNum       *big.Int     `json:"block"`
}

func (r *Receipt) DoHash() common.Hash {
	if r.Hash != nil {
		return *r.Hash
	}
	h := atlas.RLPHash(r)
	r.Hash = &h
	return h
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
