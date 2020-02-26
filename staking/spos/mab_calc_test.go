package spos

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/numeric"
	"math/big"
	"reflect"
	"testing"
)

func TestDelegationMABCalc(t *testing.T) {
	type fields struct {
		DelegatorAddress common.Address
		UpdateItems      []struct {
			blockNum *big.Int
			amount   *big.Int
		}
	}
	type args struct {
		blockNum *big.Int
	}
	type table struct {
		args
		want numeric.Dec
	}
	tests := []struct {
		name   string
		fields fields
		tables []table
	}{
		{
			name: "DelegationMABCal",
			fields: fields{
				DelegatorAddress: common.StringToAddress("0xed9d02e382b34818e88b88a309c7fe71e65f419d"),
				UpdateItems: []struct {
					blockNum *big.Int
					amount   *big.Int
				}{
					{
						blockNum: big.NewInt(0),
						amount:   big.NewInt(100),
					},
					{
						blockNum: big.NewInt(3),
						amount:   big.NewInt(200),
					},
					{
						blockNum: big.NewInt(15),
						amount:   big.NewInt(50),
					},
					{
						blockNum: big.NewInt(20),
						amount:   big.NewInt(5),
					},
				},
			},
			tables: []table{
				{
					args: args{blockNum: big.NewInt(0)},
					want: numeric.NewDec(100),
				},
				{
					args: args{blockNum: big.NewInt(3)},
					want: numeric.NewDec(100),
				},
				{
					args: args{blockNum: big.NewInt(10)},
					want: numeric.NewDec(170),
				},
				{
					args: args{blockNum: big.NewInt(15)},
					want: numeric.NewDec(50),
				},
				{
					args: args{blockNum: big.NewInt(18)},
					want: numeric.NewDec(50),
				},
				{
					args: args{blockNum: big.NewInt(20)},
					want: numeric.NewDec(5),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DelegationMAB{
				DelegatorAddress: tt.fields.DelegatorAddress,
			}
			for i := range tt.fields.UpdateItems {
				d.UpdateMAB(tt.fields.UpdateItems[i].blockNum, tt.fields.UpdateItems[i].amount)
			}
			for i := range tt.tables {
				if got := d.Calc(tt.tables[i].blockNum); !reflect.DeepEqual(got, tt.tables[i].want) {
					t.Errorf("Calc() = %v, want %v", got, tt.tables[i].want)
				}
			}
		})
	}
}

func TestValidatorMABCalc(t *testing.T) {
	type fields struct {
		Validator   common.Address
		UpdateItems []struct {
			blockNum *big.Int
			amount   *big.Int
		}
	}
	type args struct {
		blockNum *big.Int
	}
	type table struct {
		args
		want numeric.Dec
	}
	tests := []struct {
		name   string
		fields fields
		tables []table
	}{
		{
			name: "ValidatorMABCal",
			fields: fields{
				Validator: common.StringToAddress("0xed9d02e382b34818e88b88a309c7fe71e65f419d"),
				UpdateItems: []struct {
					blockNum *big.Int
					amount   *big.Int
				}{
					{
						blockNum: big.NewInt(0),
						amount:   big.NewInt(100),
					},
					{
						blockNum: big.NewInt(3),
						amount:   big.NewInt(200),
					},
					{
						blockNum: big.NewInt(15),
						amount:   big.NewInt(50),
					},
					{
						blockNum: big.NewInt(20),
						amount:   big.NewInt(5),
					},
				},
			},
			tables: []table{
				{
					args: args{blockNum: big.NewInt(0)},
					want: numeric.NewDec(100),
				},
				{
					args: args{blockNum: big.NewInt(3)},
					want: numeric.NewDec(100),
				},
				{
					args: args{blockNum: big.NewInt(10)},
					want: numeric.NewDec(170),
				},
				{
					args: args{blockNum: big.NewInt(15)},
					want: numeric.NewDec(50),
				},
				{
					args: args{blockNum: big.NewInt(18)},
					want: numeric.NewDec(50),
				},
				{
					args: args{blockNum: big.NewInt(20)},
					want: numeric.NewDec(5),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := ValidatorMAB{
				ValidatorAddress: tt.fields.Validator,
			}
			for i := range tt.fields.UpdateItems {
				d.UpdateMAB(tt.fields.UpdateItems[i].blockNum, tt.fields.UpdateItems[i].amount)
			}
			for i := range tt.tables {
				if got := d.Calc(tt.tables[i].blockNum); !reflect.DeepEqual(got, tt.tables[i].want) {
					t.Errorf("Calc() = %v, want %v", got, tt.tables[i].want)
				}
			}
		})
	}
}
