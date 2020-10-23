package abi

import (
	"fmt"
	"testing"
)

func TestPackRevert(t *testing.T) {
	var cases = []struct {
		input     string
	}{
		{"validator does not exist"},
	}
	for index, c := range cases {
		t.Run(fmt.Sprintf("case %d", index), func(t *testing.T) {
			reason, _ := PackRevert(c.input)
			got, _ := UnpackRevert(reason)
			if c.input != got {
				t.Fatalf("Output mismatch, want %v, got %v", c.input, got)
			}
		})
	}
}
