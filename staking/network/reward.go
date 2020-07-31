package network

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/clique/reward"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"math/big"
)

var (
	// BlockReward is the block reward, to be split evenly among block signers.
	BaseBlockReward = new(big.Int).Mul(big.NewInt(2), big.NewInt(params.Ether))
	// BlockRewardStakedCase is the baseline block reward in staked case -
	totalTokens = common.NewDecFromBigInt(
		new(big.Int).Mul(big.NewInt(10000000000), big.NewInt(params.Ether)),
	)
	// ErrPayoutNotEqualBlockReward ..
	ErrPayoutNotEqualBlockReward = errors.New(
		"total payout not equal to blockreward",
	)
	// NoReward ..
	NoReward = big.NewInt(0)
	// EmptyPayout ..
	EmptyPayout = noReward{}
)

type ignoreMissing struct{}

func (ignoreMissing) MissingSigners() restaking.Slots_ {
	return restaking.Slots_{}
}

type noReward struct{ ignoreMissing }

func (noReward) ReadRoundResult() *reward.CompletedRound {
	return &reward.CompletedRound{
		Total: big.NewInt(0),
		Award: []reward.Payout{},
	}
}

type stakingEra struct {
	reward.CompletedRound
	missingSigners restaking.Slots_
}

// NewStakingEraRewardForRound ..
func NewStakingEraRewardForRound(
	totalPayout *big.Int,
	mia *restaking.Slots_,
	payouts []reward.Payout,
) reward.Reader {
	return &stakingEra{
		CompletedRound: reward.CompletedRound{
			Total: totalPayout,
			Award: payouts,
		},
		missingSigners: *mia,
	}
}

// MissingSigners ..
func (r *stakingEra) MissingSigners() restaking.Slots_ {
	return r.missingSigners
}

// ReadRoundResult ..
func (r *stakingEra) ReadRoundResult() *reward.CompletedRound {
	return &r.CompletedRound
}

func CalcBlockReward(blockHeight *big.Int, config *params.ChainConfig) common.Dec {
	// TODO(ATLAS): change block reward by time
	return common.NewDecFromBigInt(BaseBlockReward)
}
