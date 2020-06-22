package network

import (
	"errors"
	"github.com/ethereum/go-ethereum/consensus/clique/reward"
	"github.com/ethereum/go-ethereum/numeric"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/staking/committee"
	"math/big"

)

var (
	// BlockReward is the block reward, to be split evenly among block signers.
	BlockReward = numeric.NewDecFromBigInt(new(big.Int).Mul(big.NewInt(2), big.NewInt(params.Ether)))
	// BlockRewardStakedCase is the baseline block reward in staked case -
	totalTokens = numeric.NewDecFromBigInt(
		new(big.Int).Mul(big.NewInt(12600000000), big.NewInt(params.Ether)),
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

func (ignoreMissing) MissingSigners() committee.SlotList {
	return committee.SlotList{}
}

type noReward struct{ ignoreMissing }

func (noReward) ReadRoundResult() *reward.CompletedRound {
	return &reward.CompletedRound{
		Total:           big.NewInt(0),
		Award:           []reward.Payout{},
	}
}

type stakingEra struct {
	reward.CompletedRound
	missingSigners committee.SlotList
}

// NewStakingEraRewardForRound ..
func NewStakingEraRewardForRound(
	totalPayout *big.Int,
	mia committee.SlotList,
	payouts []reward.Payout,
) reward.Reader {
	return &stakingEra{
		CompletedRound: reward.CompletedRound{
			Total: totalPayout,
			Award: payouts,
		},
		missingSigners: mia,
	}
}

// MissingSigners ..
func (r *stakingEra) MissingSigners() committee.SlotList {
	return r.missingSigners
}

// ReadRoundResult ..
func (r *stakingEra) ReadRoundResult() *reward.CompletedRound {
	return &r.CompletedRound
}

