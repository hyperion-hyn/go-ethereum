package burning

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"
)

func TestCalculateInternalBurningAmount(t *testing.T) {
	tests := []struct {
		activeNodeCount     int
		scalingCycleNum     int
		requireMicrostaking *big.Int
		want                *big.Int
	}{
		{
			activeNodeCount:     100,
			scalingCycleNum:     1,
			requireMicrostaking: big.NewInt(55000),
			want:                big.NewInt(27500000),
		},
		{
			activeNodeCount:     2500,
			scalingCycleNum:     4,
			requireMicrostaking: big.NewInt(55000),
			want:                big.NewInt(68750000),
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test-%v", i), func(t *testing.T) {
			got, _ := CalculateInternalBurningAmount(tt.activeNodeCount, tt.scalingCycleNum, tt.requireMicrostaking)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculateInternalBurningAmount() got = %v, want %v", got, tt.want)
			}
		})
	}
}
