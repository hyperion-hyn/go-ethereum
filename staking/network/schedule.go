package network

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
)

var (
	minSelfDelegationProportion = common.NewDecWithPrec(10, 2) // 10%
	minDelegationProportion     = common.NewDecWithPrec(1, 3)  // 0.1%
	minimalDelegation           = common.NewDec(100).MulInt64(params.Ether)

	baseStakingRequirement        = common.NewDec(550000).MulInt64(params.Ether)
	scalingCoefficient            = common.NewDecWithPrec(55, 2)
	scalingCoefficientIncremental = common.NewDecWithPrec(1, 2)
)

// according to hyperion economic model v2.1
func LatestMicrostakingRequirement(blockHeight *big.Int, config *params.ChainConfig) (*big.Int, *big.Int, *big.Int) {
	times := new(big.Int).Quo(blockHeight, big.NewInt(int64(config.Atlas.ScalingCycle))).Int64()
	if times > 4 && times <= 10 {
		h := new(big.Int).Sub(blockHeight, big.NewInt(int64(config.Atlas.ScalingCycle)*4))
		times = h.Quo(h, big.NewInt(int64(config.Atlas.ScalingCycle)*2)).Int64()
		times += 4
	} else if times > 10 {
		h := new(big.Int).Sub(blockHeight, big.NewInt(int64(config.Atlas.ScalingCycle)*10))
		times = h.Quo(h, big.NewInt(int64(config.Atlas.ScalingCycle)*3)).Int64()
		times += 7
	}

	requirement := baseStakingRequirement
	coefficient := scalingCoefficient
	for i := int64(0); i < times; i++ {
		if coefficient.GTE(common.OneDec()) {
			break
		}

		coefficient = coefficient.Add(scalingCoefficientIncremental)
		requirement = requirement.Mul(coefficient)
	}
	requireTotal := requirement.RoundInt()
	requireSelf := requirement.Mul(minSelfDelegationProportion).RoundInt()
	requireDel := requirement.Mul(minDelegationProportion)
	return requireTotal, requireSelf, common.MaxDec(requireDel, minimalDelegation).RoundInt()
}
