package microstaking

import (
	"testing"
)

func TestCheckMap3NodeEqual(t *testing.T) {
	tests := []struct {
		v1, v2 Map3Node_
	}{
		{GetDefaultMap3Node(), GetDefaultMap3Node()},
		{Map3Node_{}, Map3Node_{}},
	}
	for i, test := range tests {
		if err := CheckMap3NodeEqual(test.v1, test.v2); err != nil {
			t.Errorf("Test %v: %v", i, err)
		}
	}
}

func TestCheckMap3NodeWrapperEqual(t *testing.T) {
	tests := []struct {
		w1, w2 Map3NodeWrapper_
	}{
		{GetDefaultMap3NodeWrapper(), GetDefaultMap3NodeWrapper()},
		{Map3NodeWrapper_{}, Map3NodeWrapper_{}},
	}
	for i, test := range tests {
		if err := CheckMap3NodeWrapperEqual(test.w1, test.w2); err != nil {
			t.Errorf("Test %v: %v", i, err)
		}
	}
}

// GetDefaultMap3Node return the default microstaking.Map3Node for testing
func GetDefaultMap3Node() Map3Node_ {
	v := GetDefaultMap3NodeWrapper().Map3Node
	return v
}

// GetDefaultMap3NodeWrapper return the default microstaking.Map3NodeWrapper for testing
func GetDefaultMap3NodeWrapper() Map3NodeWrapper_ {
	v := NewMap3NodeWrapperBuilder().Build()
	return *v
}
