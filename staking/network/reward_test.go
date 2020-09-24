package network

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
	"reflect"
	"testing"
)

var (
	totalReward                  = new(big.Int).Mul(big.NewInt(100000000), big.NewInt(params.Ether))
	blocksPerHalfingCycle uint64 = 100
)

func TestRewardPool_calcBlockReward(t *testing.T) {
	tests := []struct {
		name        string
		ctx         *rewardPoolCtx
		blockHeight *big.Int
		want        *big.Int
		want1       *big.Int
	}{
		{
			name:        "genesis block",
			ctx:         makeRewardPoolCtx(t),
			blockHeight: common.Big0,
			want:        BaseBlockReward,
			want1:       totalReward,
		},
		{
			name:        "block number 1",
			ctx:         makeRewardPoolCtx(t),
			blockHeight: common.Big1,
			want:        BaseBlockReward,
			want1:       totalReward,
		},
		{
			name:        "block number 99",
			ctx:         makeRewardPoolCtx(t),
			blockHeight: big.NewInt(99),
			want:        BaseBlockReward,
			want1:       totalReward,
		},
		{
			name:        "block number 100",
			ctx:         makeRewardPoolCtx(t),
			blockHeight: big.NewInt(100),
			want: func() *big.Int {
				return common.NewDecFromInt(BaseBlockReward).MulInt64(3).QuoInt64(4).RoundInt()
			}(),
			want1: totalReward,
		},
		{
			name:        "block number 300",
			ctx:         makeRewardPoolCtx(t),
			blockHeight: big.NewInt(300),
			want: func() *big.Int {
				return common.NewDecFromInt(BaseBlockReward).
					MulInt64(3).MulInt64(3).MulInt64(3).
					QuoInt64(4).QuoInt64(4).QuoInt64(4).RoundInt()
			}(),
			want1: totalReward,
		},
		{
			name: "remaining reward less than block reward",
			ctx: func() *rewardPoolCtx {
				ctx := makeRewardPoolCtx(t)
				ctx.stateDB.SetState(rewardStorageAddress, blockRewardHashKey, common.BigToHash(big.NewInt(params.Ether)))
				return ctx
			}(),
			blockHeight: common.Big1,
			want:        big.NewInt(params.Ether),
			want1:       big.NewInt(params.Ether),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &RewardPool{
				state: tt.ctx.stateDB,
			}
			got, got1 := p.calcBlockReward(tt.blockHeight, &tt.ctx.config)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("block reward got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("remaining reward of pool got1 = %v, want %v", got1, tt.want1)
			}
		})

	}
}

func TestRewardPool_AddTxFeeAsReward(t *testing.T) {
	tests := []struct {
		name        string
		ctx         *rewardPoolCtx
		blockHeight *big.Int
		fees        []*big.Int
		want        *big.Int
	}{
		{
			name:        "total fee",
			ctx:         makeRewardPoolCtx(t),
			blockHeight: common.Big1,
			fees:        []*big.Int{common.Big1, common.Big2, common.Big3},
			want:        big.NewInt(6),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &RewardPool{
				state: tt.ctx.stateDB,
			}
			for _, fee := range tt.fees {
				p.AddTxFeeAsReward(tt.blockHeight, fee)
			}
			got := p.TxFeeAsReward(tt.blockHeight)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("block reward got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRewardPool_TakeReward(t *testing.T) {
	tests := []struct {
		name            string
		ctx             *rewardPoolCtx
		blockHeight     *big.Int
		fee             *big.Int
		wantBlockReward *big.Int
		wantRemaining   *big.Int
	}{
		{
			name:            "total fee",
			ctx:             makeRewardPoolCtx(t),
			blockHeight:     common.Big1,
			fee:             new(big.Int).Mul(big.NewInt(100), big.NewInt(params.Ether)),
			wantBlockReward: new(big.Int).Mul(big.NewInt(102), big.NewInt(params.Ether)),
			wantRemaining:   new(big.Int).Mul(big.NewInt(100000000 - 2), big.NewInt(params.Ether)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &RewardPool{
				state: tt.ctx.stateDB,
			}
			p.AddTxFeeAsReward(tt.blockHeight, tt.fee)

			// check result
			if got := p.TakeReward(tt.blockHeight, &tt.ctx.config); !reflect.DeepEqual(got, tt.wantBlockReward) {
				t.Errorf("TakeReward() = %v, want %v", got, tt.wantBlockReward)
			}
			remaining := p.state.GetState(rewardStorageAddress, blockRewardHashKey).Big()
			if !reflect.DeepEqual(remaining, tt.wantRemaining) {
				t.Errorf("Remaining reward = %v, want %v", remaining, tt.wantRemaining)
			}
		})
	}
}

type rewardPoolCtx struct {
	stateDB stateDB
	config  params.ChainConfig
}

func makeRewardPoolCtx(t *testing.T) *rewardPoolCtx {
	// stateDB
	stateDB, err := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	if err != nil {
		t.Fatal(err)
	}
	stateDB.SetState(rewardStorageAddress, blockRewardHashKey, common.BigToHash(totalReward))
	stateDB.SetNonce(rewardStorageAddress, 1)
	stateDB.Commit(true)

	// config
	config := params.ChainConfig{
		Atlas: &params.AtlasConfig{
			BlocksPerHalfingCycle: blocksPerHalfingCycle,
		},
	}
	return &rewardPoolCtx{
		stateDB: stateDB,
		config:  config,
	}
}
