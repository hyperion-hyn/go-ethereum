package network

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
)

var (
	minSelfDelegationProportion = common.NewDecWithPrec(10, 2) // 10%
	minDelegationProportion     = common.NewDecWithPrec(1, 3)  // 0.1%
	minimalDelegation           = common.NewDec(100)

	baseStakingRequirement = common.NewDec(550000).MulInt64(params.Ether)
)

type Map3NodeStakingScheduler struct {
	Config *params.ChainConfig
}

func LatestMap3StakingRequirement(blockHeight *big.Int, Config *params.ChainConfig) (*big.Int, *big.Int, *big.Int) {
	// TODO(ATLAS): requirement change by time
	minTotal := baseStakingRequirement.RoundInt()
	minSelf := baseStakingRequirement.Mul(minSelfDelegationProportion).RoundInt()
	minDel := baseStakingRequirement.Mul(minDelegationProportion)
	return minTotal, minSelf, common.MaxDec(minDel, minimalDelegation).RoundInt()
}