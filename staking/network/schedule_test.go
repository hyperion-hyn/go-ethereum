package network

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
	"reflect"
	"testing"
)

func TestLatestMicrostakingRequirement(t *testing.T) {
	tests := []struct {
		ScalingCycle uint64
		blockHeight  *big.Int
		want         common.Dec
		want1        common.Dec
		want2        common.Dec
	}{
		{
			ScalingCycle: 10,
			blockHeight:  common.Big0,
			want:         common.NewDec(550000),
			want1:        common.NewDec(55000),
			want2:        common.NewDec(550),
		},
		{
			ScalingCycle: 10,
			blockHeight:  common.Big1,
			want:         common.NewDec(550000),
			want1:        common.NewDec(55000),
			want2:        common.NewDec(550),
		},
		{
			ScalingCycle: 10,
			blockHeight:  big.NewInt(10),
			want:         common.NewDec(550000),
			want1:        common.NewDec(55000),
			want2:        common.NewDec(550),
		},
		{
			ScalingCycle: 10,
			blockHeight:  big.NewInt(11),
			want:         common.NewDec(308000),
			want1:        common.NewDec(30800),
			want2:        common.NewDec(308),
		},
		{
			ScalingCycle: 10,
			blockHeight:  big.NewInt(41),
			want:         common.NewDecWithPrec(60076632, 3),
			want1:        common.NewDecWithPrec(60076632, 4),
			want2:        common.NewDec(100),
		},
		{
			ScalingCycle: 10,
			blockHeight:  big.NewInt(100),
			want:         common.NewDecWithPrec(21988047312, 6),
			want1:        common.NewDecWithPrec(21988047312, 7),
			want2:        common.NewDec(100),
		},
		{
			ScalingCycle: 10,
			blockHeight:  big.NewInt(101),
			want:         common.NewDecWithPrec(1363258933344, 8),
			want1:        common.NewDecWithPrec(1363258933344, 9),
			want2:        common.NewDec(100),
		},
		{
			ScalingCycle: 10,
			blockHeight:  big.NewInt(131),
			want:         common.NewDecWithPrec(85885312800672, 10),
			want1:        common.NewDecWithPrec(85885312800672, 11),
			want2:        common.NewDec(100),
		},
		{
			ScalingCycle: 1296000,
			blockHeight:  common.Big0,
			want:         common.NewDec(550000),
			want1:        common.NewDec(55000),
			want2:        common.NewDec(550),
		},
		{
			ScalingCycle: 1296000,
			blockHeight:  common.Big1,
			want:         common.NewDec(550000),
			want1:        common.NewDec(55000),
			want2:        common.NewDec(550),
		},
		{
			ScalingCycle: 1296000,
			blockHeight:  big.NewInt(293087),
			want:         common.NewDec(550000),
			want1:        common.NewDec(55000),
			want2:        common.NewDec(550),
		},
		{
			ScalingCycle: 1296000,
			blockHeight:  big.NewInt(1296000),
			want:         common.NewDec(550000),
			want1:        common.NewDec(55000),
			want2:        common.NewDec(550),
		},
		{
			ScalingCycle: 1296000,
			blockHeight:  big.NewInt(1296001),
			want:         common.NewDec(308000),
			want1:        common.NewDec(30800),
			want2:        common.NewDec(308),
		},
		{
			ScalingCycle: 1296000,
			blockHeight:  big.NewInt(2592000),
			want:         common.NewDec(308000),
			want1:        common.NewDec(30800),
			want2:        common.NewDec(308),
		},
		{
			ScalingCycle: 1296000,
			blockHeight:  big.NewInt(5184001),
			want:         common.NewDecWithPrec(60076632, 3),
			want1:        common.NewDecWithPrec(60076632, 4),
			want2:        common.NewDec(100),
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test Case-%v", i), func(t *testing.T) {
			config := params.ChainConfig{
				Atlas: &params.AtlasConfig{
					ScalingCycle: tt.ScalingCycle,
				},
			}
			got, got1, got2 := LatestMicrostakingRequirement(tt.blockHeight, &config)
			want := tt.want.MulInt(big.NewInt(params.Ether)).RoundInt()
			if !reflect.DeepEqual(got, want) {
				t.Errorf("LatestMicrostakingRequirement() got = %v, want %v", got, want)
			}
			want1 := tt.want1.MulInt(big.NewInt(params.Ether)).RoundInt()
			if !reflect.DeepEqual(got1, want1) {
				t.Errorf("LatestMicrostakingRequirement() got1 = %v, want %v", got1, want1)
			}
			want2 := tt.want2.MulInt(big.NewInt(params.Ether)).RoundInt()
			if !reflect.DeepEqual(got2, want2) {
				t.Errorf("LatestMicrostakingRequirement() got2 = %v, want %v", got2, want2)
			}
		})
	}
}

func TestNumOfScalingCycle(t *testing.T) {
	tests := []struct {
		ScalingCycle uint64
		blockHeight  *big.Int
		want         int
	}{
		{
			ScalingCycle: 1296000,
			blockHeight:  common.Big0,
			want:         0,
		},
		{
			ScalingCycle: 1296000,
			blockHeight:  big.NewInt(432000),
			want:         1,
		},
		{
			ScalingCycle: 1296000,
			blockHeight:  big.NewInt(1296000),
			want:         1,
		},
		{
			ScalingCycle: 1296000,
			blockHeight:  big.NewInt(1296001),
			want:         2,
		},
		{
			ScalingCycle: 1296000,
			blockHeight:  big.NewInt(5184000),
			want:         4,
		},
		{
			ScalingCycle: 1296000,
			blockHeight:  big.NewInt(5184001),
			want:         5,
		},
		{
			ScalingCycle: 1296000,
			blockHeight:  big.NewInt(7776000),
			want:         5,
		},
		{
			ScalingCycle: 1296000,
			blockHeight:  big.NewInt(7776001),
			want:         6,
		},
		{
			ScalingCycle: 1296000,
			blockHeight:  big.NewInt(12960000),
			want:         7,
		},
		{
			ScalingCycle: 1296000,
			blockHeight:  big.NewInt(12960001),
			want:         8,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test Case-%v", i), func(t *testing.T) {
			config := params.ChainConfig{
				Atlas: &params.AtlasConfig{
					ScalingCycle: tt.ScalingCycle,
				},
			}
			if got := NumOfScalingCycle(tt.blockHeight, &config); got != tt.want {
				t.Errorf("NumOfScalingCycle() = %v, want %v", got, tt.want)
			}
		})
	}
}
