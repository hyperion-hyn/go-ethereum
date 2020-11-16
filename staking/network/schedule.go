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
	numOfScalingCycle := NumOfScalingCycle(blockHeight, config)

	requirement := baseStakingRequirement
	coefficient := scalingCoefficient
	for i := 1; i < numOfScalingCycle; i++ { // cycle num starts form 1
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

func NumOfScalingCycle(blockHeight *big.Int, config *params.ChainConfig) int {
	num := common.NewDecFromBigInt(blockHeight).QuoInt64(int64(config.Atlas.ScalingCycle))
	if !num.TruncateDec().Equal(num) { // check if it is an integer
		num = num.Add(common.OneDec())
	}
	numInt := num.TruncateInt64()

	if numInt > 4 && numInt <= 10 {
		h := common.NewDecFromBigInt(blockHeight).Sub(common.NewDec(int64(config.Atlas.ScalingCycle)*4))
		num := h.QuoInt64(int64(config.Atlas.ScalingCycle)*2)
		if !num.TruncateDec().Equal(num) { // check if it is an integer
			num = num.Add(common.OneDec())
		}
		numInt = num.TruncateInt64() + 4
	} else if numInt > 10 {
		h := common.NewDecFromBigInt(blockHeight).Sub(common.NewDec(int64(config.Atlas.ScalingCycle)*10))
		num := h.QuoInt64(int64(config.Atlas.ScalingCycle)*3)
		if !num.TruncateDec().Equal(num) { // check if it is an integer
			num = num.Add(common.OneDec())
		}
		numInt = num.TruncateInt64() + 7
	}
	return int(numInt)
}