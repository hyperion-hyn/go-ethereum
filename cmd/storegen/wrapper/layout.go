package wrapper

import (
	"encoding/json"
	"io"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

var (
	bigT      = reflect.TypeOf(&big.Int{})
	derefbigT = reflect.TypeOf(big.Int{})
	uint8T    = reflect.TypeOf(uint8(0))
	uint16T   = reflect.TypeOf(uint16(0))
	uint32T   = reflect.TypeOf(uint32(0))
	uint64T   = reflect.TypeOf(uint64(0))
	int8T     = reflect.TypeOf(int8(0))
	int16T    = reflect.TypeOf(int16(0))
	int32T    = reflect.TypeOf(int32(0))
	int64T    = reflect.TypeOf(int64(0))
	addressT  = reflect.TypeOf(common.Address{})
	stringT   = reflect.TypeOf("")
	decimalT  = reflect.TypeOf(common.Dec{})
)

// U256 converts a big Int into a 256bit EVM number.
func U256(n *big.Int) []byte {
	return math.PaddedBigBytes(math.U256(n), 32)
}

type BigInt struct {
	big.Int
}

type Member struct {
	AstId    int
	Contract string
	Label    string
	Offset   int
	Slot     BigInt
	Type     abi.Type
}

type Layout struct {
	Storage []Member
	Types   map[string]abi.Type
}

type fMember struct {
	AstId    int
	Contract string
	Label    string
	Offset   int
	Slot     BigInt
	Type     string
}

type Type struct {
	Encoding      string
	Label         string
	NumberOfBytes BigInt
	Base          string
	Key           string
	Value         string
	Members       []fMember
}

// JSON returns a parsed ABI interface and error if it failed.
func JSON(reader io.Reader) (Layout, error) {
	dec := json.NewDecoder(reader)

	var layout Layout
	if err := dec.Decode(&layout); err != nil {
		return Layout{}, err
	}
	return layout, nil
}

func (l *Layout) UnmarshalJSON(b []byte) error {
	type Layout struct {
		Storage []fMember
		Types   map[string]Type
	}

	l.Storage = make([]Member, 0)
	l.Types = make(map[string]abi.Type)

	var layout Layout
	err := json.Unmarshal(b, &layout)
	if err != nil {
		return err
	}

	layout.Types["t_struct(Global_t)"] = Type{
		Encoding: "inplace",
		Label:    "struct Global_t",
		Members:  layout.Storage,
	}

	for name, definition := range layout.Types {
		typ, err := NewType(name, definition, layout.Types)
		if err != nil {
			return err
		}
		l.Types[name] = typ
	}

	for _, variable := range layout.Storage {
		definition := layout.Types[variable.Type]
		typ, err := NewType(variable.Type, definition, layout.Types)
		if err != nil {
			return err
		}

		l.Storage = append(l.Storage, Member{
			AstId:    variable.AstId,
			Contract: variable.Contract,
			Label:    capitalise(variable.Label),
			Offset:   variable.Offset,
			Slot:     variable.Slot,
			Type:     typ,
		})
	}
	return nil
}

func (i *BigInt) UnmarshalJSON(b []byte) error {
	var val string
	err := json.Unmarshal(b, &val)
	if err != nil {
		return err
	}

	i.SetString(val, 10)

	return nil
}
