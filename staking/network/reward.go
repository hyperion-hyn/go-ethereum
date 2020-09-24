package network

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/atlas/backend/reward"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/staking/types/restaking"
	"math"
	"math/big"
)

var (
	// BlockReward is the block reward, to be split evenly among block signers.
	BaseBlockReward = new(big.Int).Mul(big.NewInt(2), big.NewInt(params.Ether))
	// NoReward ..
	NoReward = big.NewInt(0)
	// EmptyPayout ..
	EmptyPayout = noReward{}

	rewardStorageAddress = common.BigToAddress(common.Big3)
	blockRewardHashKey   = common.BigToHash(common.Big0)
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

type stateDB interface {
	GetState(addr common.Address, hash common.Hash) common.Hash
	SetState(addr common.Address, key, value common.Hash)
}

func NewRewardPool(state stateDB) RewardPool {
	return RewardPool{state: state}
}

type RewardPool struct {
	state stateDB
}

func (p *RewardPool) AddTxFeeAsReward(blockHeight, fee *big.Int) {
	total := p.TxFeeAsReward(blockHeight)
	total.Add(total, fee)
	p.state.SetState(rewardStorageAddress, common.BigToHash(blockHeight), common.BigToHash(total))
}

func (p *RewardPool) TakeReward(blockHeight *big.Int, config *params.ChainConfig) *big.Int {
	// total reward = block reward + tx fee
	// calculate block reward
	blockReward, remaining := p.calcBlockReward(blockHeight, config)

	// get tx fee by block height and compute total reward
	txFee := p.TxFeeAsReward(blockHeight)
	totalReward := big.NewInt(0).Add(blockReward, txFee)

	// remove reward from pool
	remaining.Sub(remaining, blockReward)
	p.state.SetState(rewardStorageAddress, blockRewardHashKey, common.BigToHash(remaining))
	p.state.SetState(rewardStorageAddress, common.BigToHash(blockHeight), common.BigToHash(common.Big0))
	return totalReward
}

func  (p *RewardPool) TxFeeAsReward(blockHeight *big.Int) *big.Int {
	return p.state.GetState(rewardStorageAddress, common.BigToHash(blockHeight)).Big()
}

func (p *RewardPool) BlockReward(blockHeight *big.Int, config *params.ChainConfig) *big.Int {
	r, _ := p.calcBlockReward(blockHeight, config)
	return r
}

func (p *RewardPool) calcBlockReward(blockHeight *big.Int, config *params.ChainConfig) (*big.Int, *big.Int) {
	blockReward := big.NewInt(0)
	remaining := p.state.GetState(rewardStorageAddress, blockRewardHashKey).Big()
	if remaining.Sign() > 0 {
		if blockHeight.Cmp(common.Big0) == 0 {
			blockReward.Set(BaseBlockReward)
		} else {
			quo := big.NewInt(0).Quo(blockHeight, big.NewInt(int64(config.Atlas.BlocksPerHalfingCycle))).Uint64()
			quoFloat64 := float64(quo)
			r := big.NewInt(0).Mul(BaseBlockReward, big.NewInt(int64(math.Pow(3, quoFloat64))))
			blockReward = r.Quo(r, big.NewInt(int64(math.Pow(4, quoFloat64))))
		}

		if blockReward.Cmp(remaining) > 0 {
			blockReward = remaining
		}
	}
	return blockReward, remaining
}