package reward

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/numeric"
	"github.com/ethereum/go-ethereum/staking/spos"
	"github.com/pkg/errors"
	"math/big"
)

var (
	// BaseStakedReward is the base block reward for staking, 5 HYN.
	BaseStakedReward = numeric.NewDecFromBigInt(new(big.Int).Mul(big.NewInt(5), big.NewInt(18)))

	// ErrPayoutNotEqualBlockReward ..
	ErrPayoutNotEqualBlockReward = errors.New("total payout not equal to blockreward")
	// NoReward ..
	NoReward = common.Big0
)

type Distributor interface {
	GetPercentage(address common.Address) numeric.Dec
}

type DelegationDistributorBasedMAB struct {
	ValidatorMAB *spos.ValidatorMAB
	BlockNum     *big.Int
}

func (d DelegationDistributorBasedMAB) GetPercentage(address common.Address) numeric.Dec {
	for _, delegationMAB := range d.ValidatorMAB.DelegationMABs {
		if delegationMAB.DelegatorAddress == address {
			return delegationMAB.Calc(d.BlockNum).Quo(d.ValidatorMAB.Calc(d.BlockNum))
		}
	}
	return numeric.NewDec(0)
}