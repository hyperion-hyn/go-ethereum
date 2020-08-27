package network

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
)

const (
	Million = 1000000
)

var (
	MinSelfDelegation = common.NewDecWithPrec(20, 2) // 20%
	MinDelegation     = common.NewDecWithPrec(1, 2)  // 1%

	// TODO(ATLAS): how many?
	baseStakingRequirement = common.NewDecFromBigInt(new(big.Int).Mul(big.NewInt(params.Ether), big.NewInt(Million)))
)

type Map3NodeStakingScheduler struct {
	Config *params.ChainConfig
}

func LatestMap3StakingRequirement(blockHeight *big.Int, Config *params.ChainConfig) (*big.Int, *big.Int, *big.Int) {
	// TODO(ATLAS): total node state change by time
	return baseStakingRequirement.RoundInt(),
		baseStakingRequirement.Mul(MinSelfDelegation).RoundInt(),
		baseStakingRequirement.Mul(MinDelegation).RoundInt()
}