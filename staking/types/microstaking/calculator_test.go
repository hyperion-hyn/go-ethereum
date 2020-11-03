package microstaking

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"reflect"
	"testing"
)

const (
	blocksPerEpoch = 100
)

func TestCalculatorForActivationAtOnce_Calculate(t *testing.T) {
	type args struct {
		epoch    *big.Int
		blockNum *big.Int
	}
	tests := []struct {
		name                string
		args                args
		wantActivationEpoch *big.Int
		wantReleaseEpoch    common.Dec
	}{
		{
			name: "not the last block",
			args: args{
				epoch:    big.NewInt(5),
				blockNum: big.NewInt(450),
			},
			wantActivationEpoch: big.NewInt(5),
			wantReleaseEpoch:    common.NewDec(184),
		},
		{
			name: "the last block",
			args: args{
				epoch:    big.NewInt(5),
				blockNum: big.NewInt(500),
			},
			wantActivationEpoch: big.NewInt(6),
			wantReleaseEpoch:    common.NewDec(185),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CalculatorForActivationAtOnce{
				config: fakeConfig{blocksPerEpoch: blocksPerEpoch},
			}
			gotActivationEpoch, gotReleaseEpoch := c.Calculate(tt.args.epoch, tt.args.blockNum)
			if !reflect.DeepEqual(gotActivationEpoch, tt.wantActivationEpoch) {
				t.Errorf("Calculate() gotActivationEpoch = %v, want %v", gotActivationEpoch, tt.wantActivationEpoch)
			}
			if !reflect.DeepEqual(gotReleaseEpoch, tt.wantReleaseEpoch) {
				t.Errorf("Calculate() gotReleaseEpoch = %v, want %v", gotReleaseEpoch, tt.wantReleaseEpoch)
			}
		})
	}
}

type fakeConfig struct {
	blocksPerEpoch uint64
}

func (c fakeConfig) IsMicrostakingImprove(blockNum *big.Int) bool {
	panic("no implement")
}

func (c fakeConfig) IsLastBlock(blockNum uint64) bool {
	if blockNum == 0 {
		return true
	}
	blocks := c.blocksPerEpoch
	return blockNum%blocks == 0
}

func TestCalculatorForActivationAtEndOfEpoch_Calculate(t *testing.T) {
	type args struct {
		epoch    *big.Int
		blockNum *big.Int
	}
	tests := []struct {
		name                string
		args                args
		wantActivationEpoch *big.Int
		wantReleaseEpoch    common.Dec
	}{
		{
			name: "not the last block",
			args: args{
				epoch:    big.NewInt(5),
				blockNum: big.NewInt(450),
			},
			wantActivationEpoch: big.NewInt(5),
			wantReleaseEpoch:    common.NewDec(185),
		},
		{
			name: "the last block",
			args: args{
				epoch:    big.NewInt(5),
				blockNum: big.NewInt(500),
			},
			wantActivationEpoch: big.NewInt(6),
			wantReleaseEpoch:    common.NewDec(186),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CalculatorForActivationAtEndOfEpoch{}
			gotActivationEpoch, gotReleaseEpoch := c.Calculate(tt.args.epoch, tt.args.blockNum)
			if !reflect.DeepEqual(gotActivationEpoch, tt.wantActivationEpoch) {
				t.Errorf("Calculate() gotActivationEpoch = %v, want %v", gotActivationEpoch, tt.wantActivationEpoch)
			}
			if !reflect.DeepEqual(gotReleaseEpoch, tt.wantReleaseEpoch) {
				t.Errorf("Calculate() gotReleaseEpoch = %v, want %v", gotReleaseEpoch, tt.wantReleaseEpoch)
			}
		})
	}
}