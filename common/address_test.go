// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package common

import (
	"encoding/json"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/log"
)

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(log.LvlDebug), log.StreamHandler(os.Stdout, log.TerminalFormat(true))))
}
func TestIsBech32Address(t *testing.T) {
	tests := []struct {
		str string
		exp bool
	}{
		{"hyn1t2htvpfl862vnwdqnuekd9p4ulh3h6hdldamnd", true},
		{"thyn1t2htvpfl862vnwdqnuekd9p4ulh3h6hd3c5lnu", true},
		{"HYN1T2HTVPFL862VNWDQNUEKD9P4ULH3H6HDLDAMND", true},
		{"THYN1T2HTVPFL862VNWDQNUEKD9P4ULH3H6HD3C5LNU", true},
		{"hym1t2htvpfl862vnwdqnuekd9p4ulh3h6hdldamnd", false},
		{"hyn2t2htvpfl862vnwdqnuekd9p4ulh3h6hdldamnd", false},
		{"hyn1t2htvpfl862vnwdqnuekd9p4ulh3h6hdldamnf", false},
		{"hyn1t2htvpfl862vnwdqnuekd9p4ulh3h6hdldamn", false},
		{"yn1t2htvpfl862vnwdqnuekd9p4ulh3h6hdldamnd", false},
	}

	for _, test := range tests {
		if result := IsBech32Address(test.str); result != test.exp {
			t.Errorf("IsBech32Address(%s) == %v; expected %v",
				test.str, result, test.exp)
		}
	}
}

func TestBech32AddressUnmarshalJSON(t *testing.T) {
	var tests = []struct {
		Input     string
		ShouldErr bool
		Output    *big.Int
	}{
		{"", true, nil},
		{`""`, true, nil},
		{`"0x"`, true, nil},
		{`"0x00"`, true, nil},
		{`"0xG000000000000000000000000000000000000000"`, true, nil},
		{`"0x0000000000000000000000000000000000000000"`, false, big.NewInt(0)},
		{`"0x0000000000000000000000000000000000000010"`, false, big.NewInt(16)},
		{`"hyn1t2htvpfl862vnwdqnuekd9p4ulh3h6hdldamnd"`, false, HexToAddress("0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed").Big()},
		{`"thyn1t2htvpfl862vnwdqnuekd9p4ulh3h6hd3c5lnu"`, false, HexToAddress("0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed").Big()},
	}
	for i, test := range tests {
		var v Address
		err := json.Unmarshal([]byte(test.Input), &v)
		if err != nil && !test.ShouldErr {
			t.Errorf("test #%d: unexpected error: %v", i, err)
		}
		if err == nil {
			if test.ShouldErr {
				t.Errorf("test #%d: expected error, got none", i)
			}
			if got := new(big.Int).SetBytes(v.Bytes()); got.Cmp(test.Output) != 0 {
				t.Errorf("test #%d: address mismatch: have %v, want %v", i, got, test.Output)
			}
		}
	}
}

func TestAddressToBech32(t *testing.T) {
	var tests = []struct {
		Input  string
		Output string
	}{
		// Test cases from https://github.com/ethereum/EIPs/blob/master/EIPS/eip-55.md#specification
		{"0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed", "hyn1t2htvpfl862vnwdqnuekd9p4ulh3h6hdldamnd"},
		{"0xfb6916095ca1df60bb79ce92ce3ea74c37c5d359", "hyn1ld53vz2u580kpwmee6fvu048fsmut56ehjtunc"},
		{"0xdbf03b407c01e7cd3cbea99509d93f8dddc8c6fb", "hyn1m0crksruq8nu60974x2snkfl3hwu33hm9s2gua"},
		{"0xd1220a0cf47c7b9be7a2e6ba89f429762e7b9adb", "hyn16y3q5r8503aeheazu6agnapfwch8hxkmzcz9df"},
		{"0x70247395aFFd13C2347aA8c748225f1bFeD2C32A", "hyn1wqj889d0l5fuydr64rr5sgjlr0ld9se20c3ckj"},
		// Ensure that non-standard length input values are handled correctly
		{"0xa", "hyn1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq2kgxvca"},
		{"0x0a", "hyn1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq2kgxvca"},
		{"0x00a", "hyn1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq2kgxvca"},
		{"0x000000000000000000000000000000000000000a", "hyn1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq2kgxvca"},
	}
	for i, test := range tests {
		output := HexToAddress(test.Input).Bech32()
		if output != test.Output {
			t.Errorf("test #%d: failed to match when it should (%s != %s)", i, output, test.Output)
		}
	}
}

func BenchmarkAddressToBech32(b *testing.B) {
	testAddr := HexToAddress("0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed")
	for n := 0; n < b.N; n++ {
		testAddr.Bech32()
	}
}

func TestParseAddr(t *testing.T) {
	var tests = []struct {
		input string
		want  string		
	}{
		{"0x15a128e599b74842bccba860311efa92991bffb5", "0x15a128e599b74842BCcBa860311Efa92991bffb5"},
		{"hyn1zksj3evekayy90xt4psrz8h6j2v3hla4840g9v", "0x15a128e599b74842BCcBa860311Efa92991bffb5"},
		{"hyn1wqj889d0l5fuydr64rr5sgjlr0ld9se20c3ckj", "0x70247395aFFd13C2347aA8c748225f1bFeD2C32A"},
		{"deadbeef", "0x00000000000000000000000000000000DeaDBeef"},
		{"helloworld", "0x0000000000000000000000000000000000000000"},
	}

	for i, test := range tests {

		addr := ParseAddr(test.input)
		expected := HexToAddress(test.want)

		if addr.Big().Cmp(expected.Big()) != 0 {
			t.Errorf("test #%d: parse %v, want %v, have %v", i, test.input, expected.Hex(), addr.Hex())
		}
	}
}
