package network

import (
	"github.com/ethereum/go-ethereum/numeric"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
)

const (
	Million = 1000000
)

var (
	MinSelfDelegation = numeric.NewDecWithPrec(20, 2) // 20%
	MinDelegation     = numeric.NewDecWithPrec(1, 2)  // 1%

	baseStakingRequirement = numeric.NewDecFromBigInt(new(big.Int).Mul(big.NewInt(params.Ether), big.NewInt(Million)))

)

type Map3NodeStakingScheduler struct {
	Config *params.ChainConfig
}

func LatestStakingRequirement(blockHeight *big.Int, Config *params.ChainConfig) (*big.Int, *big.Int, *big.Int) {
	// TODO: total node state change by time
	return baseStakingRequirement.RoundInt(),
		baseStakingRequirement.Mul(MinSelfDelegation).RoundInt(),
		baseStakingRequirement.Mul(MinDelegation).RoundInt()
}

func HaveRequirementChangeInEpoch(epoch *big.Int, Config *params.ChainConfig) bool {
	return false
}