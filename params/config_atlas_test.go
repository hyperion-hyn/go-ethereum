package params

import (
	"fmt"
	"testing"
)

func TestAtlasConfig_IsFirstBlock(t *testing.T) {
	tests := []struct {
		BlocksPerEpoch uint64
		blockNum       uint64
		want           bool
	}{
		{
			BlocksPerEpoch: 15,
			blockNum:       0,
			want:           true,
		},
		{
			BlocksPerEpoch: 15,
			blockNum:       1,
			want:           true,
		},
		{
			BlocksPerEpoch: 15,
			blockNum:       2,
			want:           false,
		},
		{
			BlocksPerEpoch: 15,
			blockNum:       15,
			want:           false,
		},
		{
			BlocksPerEpoch: 15,
			blockNum:       16,
			want:           true,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("Case-%v", i), func(t *testing.T) {
			c := &AtlasConfig{
				BlocksPerEpoch: tt.BlocksPerEpoch,
			}
			if got := c.IsFirstBlock(tt.blockNum); got != tt.want {
				t.Errorf("IsLastBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAtlasConfig_IsLastBlock(t *testing.T) {
	tests := []struct {
		BlocksPerEpoch uint64
		blockNum       uint64
		want           bool
	}{
		{
			BlocksPerEpoch: 15,
			blockNum:       0,
			want:           true,
		},
		{
			BlocksPerEpoch: 15,
			blockNum:       1,
			want:           false,
		},
		{
			BlocksPerEpoch: 15,
			blockNum:       15,
			want:           true,
		},
		{
			BlocksPerEpoch: 15,
			blockNum:       16,
			want:           false,
		},
		{
			BlocksPerEpoch: 15,
			blockNum:       30,
			want:           true,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("Case-%v", i), func(t *testing.T) {
			c := &AtlasConfig{
				BlocksPerEpoch: tt.BlocksPerEpoch,
			}
			if got := c.IsLastBlock(tt.blockNum); got != tt.want {
				t.Errorf("IsLastBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAtlasConfig_EpochLastBlock(t *testing.T) {
	tests := []struct {
		BlocksPerEpoch uint64
		epochNum       uint64
		want           uint64
	}{
		{
			BlocksPerEpoch: 15,
			epochNum:       0,
			want:           0,
		},
		{
			BlocksPerEpoch: 15,
			epochNum:       1,
			want:           15,
		},
		{
			BlocksPerEpoch: 15,
			epochNum:       2,
			want:           30,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("Case-%v", i), func(t *testing.T) {
			c := &AtlasConfig{
				BlocksPerEpoch: tt.BlocksPerEpoch,
			}
			if got := c.EpochLastBlock(tt.epochNum); got != tt.want {
				t.Errorf("EpochLastBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAtlasConfig_EpochFirstBlock(t *testing.T) {
	tests := []struct {
		BlocksPerEpoch uint64
		epochNum       uint64
		want           uint64
	}{
		{
			BlocksPerEpoch: 15,
			epochNum:       0,
			want:           0,
		},
		{
			BlocksPerEpoch: 15,
			epochNum:       1,
			want:           1,
		},
		{
			BlocksPerEpoch: 15,
			epochNum:       2,
			want:           16,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("Case-%v", i), func(t *testing.T) {
			c := &AtlasConfig{
				BlocksPerEpoch: tt.BlocksPerEpoch,
			}
			if got := c.EpochFirstBlock(tt.epochNum); got != tt.want {
				t.Errorf("EpochFirstBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAtlasConfig_EpochByBlock(t *testing.T) {
	tests := []struct {
		BlocksPerEpoch uint64
		blockNum       uint64
		want           uint64
	}{
		{
			BlocksPerEpoch: 15,
			blockNum:       0,
			want:           0,
		},
		{
			BlocksPerEpoch: 15,
			blockNum:       1,
			want:           1,
		},
		{
			BlocksPerEpoch: 15,
			blockNum:       15,
			want:           1,
		},
		{
			BlocksPerEpoch: 15,
			blockNum:       16,
			want:           2,
		},
		{
			BlocksPerEpoch: 15,
			blockNum:       31,
			want:           3,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("Case-%v", i), func(t *testing.T) {
			c := &AtlasConfig{
				BlocksPerEpoch: tt.BlocksPerEpoch,
			}
			if got := c.EpochByBlock(tt.blockNum); got != tt.want {
				t.Errorf("EpochByBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}
