package wrapper

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// tmplData is the data structure required to fill the binding template.
type tmplData struct {
	Package   string                   // Name of the package to place the generated file in
	Defines map[string]*tmplStruct
	Structs   map[string]*tmplStruct   // Contract struct type definitions
	Basics   map[string]*tmplStruct   // Contract struct type definitions
	Variables   map[string]abi.Type   // Contract struct type definitions
}


// tmplField is a wrapper around a struct field with binding language
// struct type definition and relative filed name.
type tmplField struct {
	Type    string   // Field type representation depends on target binding language
	ElemType string
	Name    string   // Field name converted from the raw user-defined field name
	SolKind abi.Type // Raw abi type information
	Tag     string
	Slot    *big.Int
	Offset  int
}

// tmplStruct is a wrapper around an abi.tuple contains a auto-generated
// struct name.
type tmplStruct struct {
	T   byte
	Name   string       // Auto-generated struct name(before solidity v0.5.11) or raw name.
	Fields []*tmplField // Struct fields definition depends on the binding language.
	Type  string
	SolKind abi.Type
}

// tmplSource is language to template mapping containing all the supported
// programming languages the package can generate to.
var tmplSource = map[bind.Lang]string{
	bind.LangGo:   tmplSourceGo,
}

// tmplSourceGo is the Go source template use to generate the contract binding
// based on.
const tmplSourceGo = `
// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package {{.Package}}

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// use backtick in text/template: https://github.com/golang/go/issues/18221
{{ $tick := "`+"`"+`" }}


{{$basics := .Basics}}
{{$defines := .Defines}}
{{$structs := .Structs}}
{{$variables := .Variables}}

type StateValues map[common.Hash]common.Hash

func GetStateAsBytes(db *state.StateDB, addr common.Address, slot *big.Int) []byte{
	var retval []byte

	hash := db.GetState(addr, common.BigToHash(slot))
	flag := int(hash[31]) & 0x01
	if flag == 0 { 
		// less than 32 bytes
		length := int(hash[31])/2
		retval = hash[0:length]
	} else {
		// more than 32 bytes
		length := (hash.Big().Uint64() - 1)/2
		base := crypto.Keccak256Hash(common.BigToHash(slot).Bytes()).Big()
		for remaining, i := length, uint64(0); remaining <= length; remaining, i = remaining - 32, i+1 {
			var available uint64
			h := db.GetState(addr, common.BigToHash(big.NewInt(0).Add(base, big.NewInt(0).SetUint64(i))))
			if remaining >= 32 {
				available = 32
			} else {
				available = remaining
			}
			retval = append(retval, h.Bytes()[:available]...)
		}
	}
	
	return retval
}

func SetStateAsBytes(db *state.StateDB, addr common.Address, slot *big.Int, value []byte) {
	length := uint64(len(value))
	if length < 32 {
		// less than 32 bytes
		var hash common.Hash
		copy(hash[:length], value)
		hash[31] = byte(length * 2)
		db.SetState(addr, common.BigToHash(slot), hash)
	} else {
		// more than 32 bytes
		db.SetState(addr, common.BigToHash(slot), common.BigToHash(big.NewInt(int64(length * 2 + 1))))
		base := crypto.Keccak256Hash(common.BigToHash(slot).Bytes()).Big()
		for remaining, i := uint64(length), uint64(0); remaining <= length; remaining, i = remaining - 32, i+1 {
			var available uint64
			if remaining >= 32 {
				available = 32
			} else {
				available = remaining
			}
			var hash common.Hash
			copy(hash[:available], value[length-remaining:length-remaining+available])
			db.SetState(addr, common.BigToHash(big.NewInt(0).Add(base, big.NewInt(0).SetUint64(i))), hash)
		}
	}
}


{{range $basics}}
type {{.Name}}={{.Type}}
type Storage_{{.Name}} struct {
{{- if eq .Name "BigInt"}}
	obj    {{.Type}}
{{else}}
	obj    *{{.Type}}
{{end -}}
	db     *state.StateDB
	addr   common.Address
	slot   *big.Int
	dirty  StateValues	
}

func (s *Storage_{{.Name}}) Value() {{.Type}} {
{{- if eq .Name "String"}}
	rv := GetStateAsBytes(s.db, s.addr, s.slot)
	*s.obj = {{.Type}}(rv)
	return *s.obj
{{else if eq .Name "Bytes"}}
	rv := GetStateAsBytes(s.db, s.addr, s.slot)
	*s.obj = {{.Type}}(rv)
	return *s.obj
{{else if eq .Name "BigInt"}}
	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	*s.obj = *hash.Big()
	return s.obj
{{else}}
	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
{{- if eq .Name "Uint8" "Uint16" "Uint32" "Uint64"}}
	*s.obj = {{.Type}}(hash.Big().Uint64())	
{{else if eq .Name "Int8" "Int16" "Int32" "Int64"}}
	*s.obj = {{.Type}}(hash.Big().Int64())
{{else if eq .Name "Bool"}}
	*s.obj = (hash.Big().Int64() != 0)
{{else if eq .Name "BigInt"}}
{{else if eq .Name "Address"}}
	*s.obj = common.BigToAddress(hash.Big())
{{else}}
	UNSUPPORTED {{.Name}} {{.Type}}
{{end -}}
	return *s.obj
{{end -}}
}

func (s *Storage_{{.Name}}) SetValue(value {{.Type}}) {
{{- if eq .Name "String"}}
	SetStateAsBytes(s.db, s.addr, s.slot, []byte(value))
	*s.obj = value
{{else if eq .Name "Bytes" }}
	SetStateAsBytes(s.db, s.addr, s.slot, []byte(value))
	*s.obj = value
{{else if eq .Name "Uint8" "Uint16" "Uint32" "Uint64"}}
	hash := big.NewInt(0).SetUint64(uint64(value))
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(hash))
	*s.obj = value
{{else if eq .Name "Int8" "Int16" "Int32" "Int64"}}
	hash := big.NewInt(0).SetInt64(int64(value))
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(hash))
	*s.obj = value
{{else if eq .Name "Bool"}}
	var val uint
	if value {
		val = 1
	} else {
		val = 0
	}
	hash := big.NewInt(0).SetUint64(uint64(val))
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(hash))
	*s.obj = value
{{else if eq .Name "BigInt"}}
	hash := value
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(hash))
	*s.obj = *value
{{else if eq .Name "Address"}}
	hash := value.Hash()
	s.db.SetState(s.addr, common.BigToHash(s.slot), hash)
	*s.obj = value
{{else}}
	UNSUPPORTED {{.Name}} {{.Type}}
{{end -}}
}
{{end}}


{{range $defines}}
// {{.Name}} is an auto generated low-level Go binding around an user-defined struct.
// {{ printf "%#v" . }}
type {{.Name}} {{.Type}}

type Storage_{{.Name}} struct {
{{- if or (isptr .) (ismap .)}}
	obj    {{.Name}}
{{else}}
	obj    *{{.Name}}
{{end -}}
	db     *state.StateDB
	addr   common.Address
	slot   *big.Int
	dirty  StateValues	
}


{{- if isarray .}}
{{$elem := index .Fields 0}}
func (s* Storage_{{.Name}}) Length() int {
	return len(s.obj)
}

func (s* Storage_{{.Name}}) Get(index uint64) ( *Storage_{{$elem.Type}} ) {
	// Value: {{ printf "%#v" $elem }}
	actual := big.NewInt(0).Add(s.slot, big.NewInt(0).SetUint64(index*({{$elem.SolKind.NumberOfBytes}}/32)))

{{- if or (isptr $elem) (isslice $elem) (ismap $elem) }}
	if s.obj[index] == nil {
{{- if eq $elem.Type "BigInt" }}
		s.obj[index] = big.NewInt(0)
{{else if or (isptr $elem)}}
		s.obj[index] = new({{$elem.Type}})
{{else if or (isslice $elem)}}
		s.obj[index] = make({{$elem.Type}}, 0)
{{else if or (ismap $elem)}}
		s.obj[index] = make({{$elem.Type}})
{{else if or (isarray $elem)}}
		
{{end}}
	}
{{end}}

	return &Storage_{{ $elem.Type }} {
{{- if or (isptr $elem) (ismap $elem)}}
		obj: s.obj[index],
{{else}}
		obj: &s.obj[index],
{{end -}}
		db: s.db,
		addr: s.addr,
		slot: actual,
		dirty: s.dirty,
	}
}

{{end}}


{{- if isslice .}}
{{$elem := index .Fields 0}}

func (s* Storage_{{.Name}}) Length() (*big.Int) {
	rv := s.db.GetState(s.addr, common.BigToHash(s.slot))
	return rv.Big()
}

func (s* Storage_{{.Name}}) Get(index uint64) ( *Storage_{{$elem.Type}} ) {
	// Value: {{ printf "%#v" $elem }}
	indexBytes := common.BigToHash(big.NewInt(0).SetUint64(index)).Bytes()
	hash := crypto.Keccak256Hash(append(indexBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()

	return &Storage_{{ $elem.Type }} {
		obj: (*s.obj)[index],
		db: s.db,
		addr: s.addr,
		slot: actual,
		dirty: s.dirty,
	}
}

{{end}}


{{- if ismap . }}
{{$elemKey := index .Fields 0}}
{{$elemValue := index .Fields 1}}
func (s* Storage_{{.Name}}) Get(key {{$elemKey.Type}}) ( *Storage_{{$elemValue.Type}} ) {
	// Key:	{{ printf "%#v" $elemKey }}
	// Value: {{ printf "%#v" $elemValue }}
{{- if eq $elemKey.Type "String"}}
	keyBytes := []byte(key)
{{else if eq $elemKey.Type "Uint8" "Uint16" "Uint32" "Uint64"}}
	keyBytes := big.NewInt(0).SetUint64(uint64(key)).Bytes()
{{else if eq $elemKey.Type "Int8" "Int16" "Int32" "Int64"}}
	keyBytes := big.NewInt(0).SetInt64(int64(key)).Bytes()
{{else if eq $elemKey.Type "Bool"}}
	keyBytes := big.NewInt(0).SetInt64(int64(key)).Bytes()
{{else if eq $elemKey.Type "BigInt"}}
	keyBytes := key.Bytes()
{{else if eq $elemKey.Type "Address"}}
	keyBytes := key.Hash().Bytes()
{{else}}
	UNSUPPORTED {{$elemKey.Type}}
{{end -}}

	hash := crypto.Keccak256Hash(append(keyBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()
	
{{- if or (isptr $elemValue) (isslice $elemValue) (ismap $elemValue) }}
	if s.obj[key] == nil {
{{- if eq $elemValue.Type "BigInt"}}
		s.obj[key] = big.NewInt(0)
{{else if isptr $elemValue}}
		s.obj[key] = new({{$elemValue.Type}})
{{else if or (isslice $elemValue)}}
		s.obj[key] = make({{$elemValue.Type}}, 0)
{{else if or (ismap $elemValue)}}
		s.obj[key] = make({{$elemValue.Type}})
{{else if or (isarray $elemValue)}}
{{else}}
	here {{$elemValue.Type}}
{{end}}
	}
{{end}}

	return &Storage_{{ $elemValue.Type }} {
{{- if or (isptr $elemValue) (ismap $elemValue)}}
		obj: s.obj[key],
{{else}}
		obj: &s.obj[key],
{{end -}}
		db: s.db,
		addr: s.addr,
		slot: actual,
		dirty: s.dirty,
	}
}

{{end}}

{{end}}


{{range $structs}}
	// {{.Name}} is an auto generated low-level Go binding around an user-defined struct.
	// {{ printf "%#v" . }}
	type {{.Name}} struct {
	{{range $field := .Fields}}
	{{$field.Name}} {{$field.Type}} {{$tick}}{{$field.Tag}}{{$tick}}{{end}}
	}
{{end}}


{{range $structs}}
type Storage_{{.Name}} struct {
	obj *{{.Name}}
	db     *state.StateDB
	addr common.Address
	slot *big.Int
	dirty StateValues
}
{{end}}


func New(g *Global_t, db *state.StateDB, addr common.Address, slot *big.Int) *Storage_Global_t {
	return &Storage_Global_t {
		obj: g,
		db: db,
		addr: addr,
		slot: slot,
		dirty: make(StateValues),
	}
}


{{range $structs}}
{{$typeName := .Name}}
{{range $field := .Fields}}
	// {{.Name}} is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_{{ $typeName }}) {{$field.Name}}() (*Storage_{{$field.Type}}) {
	// Field: {{ printf "%#v" $field }}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("{{$field.Slot}}", 10)
	
	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
{{- if or (isptr $field) (isslice $field) (ismap $field) }}
	if s.obj.{{$field.Name}} == nil {
{{- if eq $field.Type "BigInt" }}
		s.obj.{{$field.Name}} = big.NewInt(0)
{{else if or (isptr $field)}}
		s.obj.{{$field.Name}} = new({{$field.Type}})
{{else if or (isslice $field)}}
		s.obj.{{$field.Name}} = make({{$field.Type}}, 0)
{{else if or (ismap $field)}}
		s.obj.{{$field.Name}} = make({{$field.Type}})
{{else if or (isarray $field)}}
		
{{end -}}
	}
{{end}}
	return &Storage_{{ $field.Type }} {
{{- if or (isptr $field) (ismap $field)}}
		obj: s.obj.{{$field.Name}},
{{else}}
		obj: &s.obj.{{$field.Name}},
{{end -}}
		db: s.db,
		addr: s.addr,
		slot: actual,
		dirty: s.dirty,
	}
}
	{{end}}
{{end}}


`

