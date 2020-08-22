package wrapper

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// tmplData is the data structure required to fill the binding template.
type tmplData struct {
	Package   string // Name of the package to place the generated file in
	Defines   map[string]*tmplStruct
	Structs   map[string]*tmplStruct // Contract struct type definitions
	Basics    map[string]*tmplStruct // Contract struct type definitions
	Variables map[string]abi.Type    // Contract struct type definitions
}

// tmplField is a wrapper around a struct field with binding language
// struct type definition and relative filed name.
type tmplField struct {
	Type     string // Field type representation depends on target binding language
	ElemType string
	Name     string   // Field name converted from the raw user-defined field name
	SolKind  abi.Type // Raw abi type information
	Tag      string
	Slot     *big.Int
	Offset   int
}

// tmplStruct is a wrapper around an abi.tuple contains a auto-generated
// struct name.
type tmplStruct struct {
	T       byte
	Name    string       // Auto-generated struct name(before solidity v0.5.11) or raw name.
	Fields  []*tmplField // Struct fields definition depends on the binding language.
	Type    string
	SolKind abi.Type
}

// tmplSource is language to template mapping containing all the supported
// programming languages the package can generate to.
var tmplSource = map[bind.Lang]string{
	bind.LangGo: tmplSourceGo,
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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/common/math"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = math.MaxInt32
)

// use backtick in text/template: https://github.com/golang/go/issues/18221
{{ $tick := "` + "`" + `" }}

type StateDB interface {
	GetState(addr common.Address, hash common.Hash) common.Hash
	SetState(addr common.Address, key, value common.Hash)
}


{{$basics := .Basics}}
{{$defines := .Defines}}
{{$structs := .Structs}}
{{$variables := .Variables}}

type StateValues map[common.Hash]common.Hash

func GetStateAsBytes(db StateDB, addr common.Address, slot *big.Int) []byte{
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

func SetStateAsBytes(db StateDB, addr common.Address, slot *big.Int, value []byte) {
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

{{define "storage_fields"}}
// {{ printf "%#v" . }}

{{- if eq .Name "BigInt"}}
	obj    {{.Name}}
{{else if or (isptr .) (ismap .)}}
	obj    {{.Name}}
{{else}}
	obj    *{{.Name}}
{{end -}}
	db     StateDB
	addr   common.Address
	slot   *big.Int
	offset int
	numberOfBytes int
	dirty  StateValues	
{{end}}{{/* storage_fields */}}


{{define "new_instance"}}
{{- if eq .Type "BigInt" }}
		instance := big.NewInt(0)
{{else if or (isptr .)}}
		instance := new({{.Type}})
{{else if or (isslice .)}}
		hash := s.db.GetState(s.addr, common.BigToHash(actual))
		instance := make({{.Type}}, hash.Big().Int64())
{{else if or (ismap .)}}
		instance := make({{.Type}})
{{end -}}
{{end}}{{/* new_instance */}}

{{define "new_slice"}}
{{- if eq .Type "BigInt" }}
		slice := make([]{{.Type}}, _length, _capacity)
{{else}}
		slice := make([]*{{.Type}}, _length, _capacity)
{{end -}}
{{end}}{{/* new_slice */}}

{{define "getBytes"}}
	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	data := hash.Bytes()[32 - (s.offset + s.numberOfBytes) : 32 - s.offset]
{{end}}{{/* getBytes */}}

{{define "setBytes"}}
	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	data := hash.Bytes()
	copy(data[32 - (s.offset + s.numberOfBytes) : 32 - s.offset], val[len(val) - s.numberOfBytes:])
	hash.SetBytes(data)
	s.db.SetState(s.addr, common.BigToHash(s.slot), hash)
{{end}}{{/* setBytes */}}

{{range $basics}}
type {{.Name}}={{.Type}}
type Storage_{{.Name}} struct {
// Builtin-Type
{{template "storage_fields" .}}
}

func (s *Storage_{{.Name}}) Value() {{.Type}} {
	// {{ printf "%#v" . }}
{{- if eq .Name "String"}}
	rv := GetStateAsBytes(s.db, s.addr, s.slot)
	*s.obj = {{.Type}}(rv)
	return *s.obj
{{else if eq .Name "Bytes"}}
	rv := GetStateAsBytes(s.db, s.addr, s.slot)
	*s.obj = {{.Type}}(rv)
	return *s.obj
{{else if eq .Name "BigInt"}}
	{{template "getBytes" .}}
	*s.obj = *(big.NewInt(0).SetBytes(data))
	retval := big.NewInt(0).SetBytes(data)
	return retval
{{else if eq .Name "Uint8" "Uint16" "Uint32" "Uint64"}}
	{{template "getBytes" .}}
	*s.obj = {{.Type}}(big.NewInt(0).SetBytes(data).Uint64())
	return *s.obj
{{else if eq .Name "Int8" "Int16" "Int32" "Int64"}}
	{{template "getBytes" .}}
	*s.obj = {{.Type}}(big.NewInt(0).SetBytes(data).Int64())
	return *s.obj
{{else if eq .Name "Bool"}}
	{{template "getBytes" .}}
	*s.obj = {{.Type}}(big.NewInt(0).SetBytes(data).Int64() != 0)
	return *s.obj
{{else if eq .Name "Address"}}
	{{template "getBytes" .}}
	*s.obj = common.BytesToAddress(data)
	return *s.obj
{{else if eq .Name "Decimal"}}
	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	*s.obj = common.NewDecFromBigIntWithPrec(hash.Big(), common.Precision)
	return *s.obj
{{else if match .Name "Bytes([1-9]|[12][0-9]|3[0-2])" }}
	{{template "getBytes" .}}
	copy((*s.obj)[:], data[:])
	return *s.obj
{{else}}
	UNSUPPORTED {{.Name}} {{.Type}}
{{end -}}
}

func (s *Storage_{{.Name}}) SetValue(value {{.Type}}) {
// {{ printf "%#v" . }}
{{- if eq .Name "String"}}
	SetStateAsBytes(s.db, s.addr, s.slot, []byte(value))
	*s.obj = value
{{else if eq .Name "Bytes" }}
	SetStateAsBytes(s.db, s.addr, s.slot, []byte(value))
	*s.obj = value
{{else if eq .Name "BigInt"}}
	val := math.PaddedBigBytes(value, 32)
	{{template "setBytes" .}}
	*s.obj = *(hash.Big())
{{else if eq .Name "Uint8" "Uint16" "Uint32" "Uint64"}}
	val := math.PaddedBigBytes(big.NewInt(0).SetUint64(uint64(value)), 32)
	{{template "setBytes" .}}
	*s.obj = value
{{else if eq .Name "Int8" "Int16" "Int32" "Int64"}}
	val := math.PaddedBigBytes(big.NewInt(0).SetInt64(int64(value)), 32)
	{{template "setBytes" .}}
	*s.obj = value
{{else if eq .Name "Bool"}}
	var flag uint
	if value {
		flag = 1
	} else {
		flag = 0
	}

	val := math.PaddedBigBytes(big.NewInt(0).SetInt64(int64(flag)), 32)
	{{template "setBytes" .}}
	*s.obj = value
{{else if eq .Name "Address"}}
	val := value
	{{template "setBytes" .}}
	*s.obj = value
{{else if eq .Name "Decimal"}}
	hash := value.BigInt()
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(hash))
	*s.obj = value
{{else if match .Name "Bytes([1-9]|[12][0-9]|3[0-2])" }}
	val := value
	{{template "setBytes" .}}
	copy((*s.obj)[:], value[:])
{{else}}
	UNSUPPORTED {{.Name}} {{.Type}}
{{end -}}
}
{{end}}{{/* basics */}}


{{range $defines}}
// {{.Name}} is an auto generated low-level Go binding around an user-defined struct.
// {{ printf "%#v" . }}
// {{ printf "%#v" (GetReflectType .SolKind).String }}
type {{.Name}} {{.Type}}

type Storage_{{.Name}} struct {
{{template "storage_fields" .}}
}


{{- if isarray .}}

{{$elem := index .Fields 0}}
{{ if isFixedSizeByteArray (GetReflectType .SolKind) }}

func (s* Storage_{{.Name}}) Value() {{.Type}} {
	// {{ printf "%#v" . }}
	length := len(*s.obj)
	base := s.slot
	for offset, i := 0, uint64(0); offset < length; offset, i = offset+32, i+1 {
		var available int
		h := s.db.GetState(s.addr, common.BigToHash(big.NewInt(0).Add(base, big.NewInt(0).SetUint64(i))))
		remaining := length - offset
		if remaining >= 32 {
			available = 32
		} else {
			available = remaining
		}

		data := h.Bytes()
		for j := 0; j < available; j++ {
			(*s.obj)[offset + j] = data[31-j]
		}
	}
	return *s.obj
}

func (s *Storage_{{.Name}}) SetValue(value {{.Type}}) {
	// {{ printf "%#v" . }}
	length := len(*s.obj)
	base := s.slot
	for offset, i := 0, uint64(0); offset < length; offset, i = offset+32, i+1 {
		var available int
		var val [32]byte

		remaining := length - offset
		if remaining >= 32 {
			available = 32
		} else {
			available = remaining
		}

		for j := 0; j < available; j++ {
			val[31-j] = value[offset+j]
		}

		s.db.SetState(s.addr, common.BigToHash(big.NewInt(0).Add(base, big.NewInt(0).SetUint64(i))), common.BytesToHash(val[:]))
	}
}

{{ else }}
func (s* Storage_{{.Name}}) Length() (int) {
	// Value: {{ printf "%#v" . }}
	return len(s.obj)
}

func (s* Storage_{{.Name}}) Get(index int) ( *Storage_{{$elem.Type}} ) {
	// Value: {{ printf "%#v" $elem }}
	base := s.slot
{{ if le $elem.SolKind.NumberOfBytes 32 }}
	itemsPerSlot := 32/{{$elem.SolKind.NumberOfBytes}}
	actual := big.NewInt(0).Add(base, big.NewInt(0).SetUint64(uint64(index/itemsPerSlot)))
	offset := ((index % itemsPerSlot) * {{$elem.SolKind.NumberOfBytes}})
{{else}}
	slotsPerItem := ({{$elem.SolKind.NumberOfBytes}} + 31)/32
	actual := big.NewInt(0).Add(base, big.NewInt(0).SetUint64(uint64(index * slotsPerItem)))
	offset := 0
{{end}}
{{- if or (isptr $elem) (isslice $elem) (ismap $elem) (eq $elem.Type "BigInt") }}
	if s.obj[index] == nil {
		{{template "new_instance" $elem}}
		s.obj[index] = instance
	}
{{end}}

	return &Storage_{{ $elem.Type }} {
{{- if or (isptr $elem) (ismap $elem) (eq $elem.Type "BigInt") }}
		obj: s.obj[index],
{{else}}
		obj: &s.obj[index],
{{end -}}
		db: s.db,
		addr: s.addr,
		slot: actual,
		offset: offset,
		numberOfBytes: {{$elem.SolKind.NumberOfBytes}},
		dirty: s.dirty,
	}
}

{{end}} {{/* fixed-size byte array */}}
{{end}} {{/* array */}}


{{- if isslice .}}
{{$elem := index .Fields 0}}

func (s* Storage_{{.Name}}) Length() (int) {
	rv := s.db.GetState(s.addr, common.BigToHash(s.slot))
	return int(rv.Big().Int64())
}

func (s* Storage_{{.Name}}) Resize(length int) {
	// Value: {{ printf "%#v" $elem }}
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(big.NewInt(0).SetUint64(uint64(length))))

	_length := length
	_capacity := length + 50
	{{template "new_slice" $elem}}
	copy(slice, *s.obj)
	*s.obj = slice
}

func (s* Storage_{{.Name}}) Get(index int) ( *Storage_{{$elem.Type}} ) {
	// Value: {{ printf "%#v" $elem }}
	length := s.Length()
	if length <= index {
		s.Resize(index+1)
	}

	base := crypto.Keccak256Hash(common.BigToHash(s.slot).Bytes()).Big()
{{ if le $elem.SolKind.NumberOfBytes 32 }}
	itemsPerSlot := 32/{{$elem.SolKind.NumberOfBytes}}
	actual := big.NewInt(0).Add(base, big.NewInt(0).SetUint64(uint64(index/itemsPerSlot)))
	offset := ((index % itemsPerSlot) * {{$elem.SolKind.NumberOfBytes}})
{{else}}
	slotsPerItem := ({{$elem.SolKind.NumberOfBytes}} + 31)/32
	actual := big.NewInt(0).Add(base, big.NewInt(0).SetUint64(uint64(index * slotsPerItem)))
	offset := 0
{{end}}

{{- if or (isptr $elem) (isslice $elem) (ismap $elem) (eq $elem.Type "BigInt") }}
	if (*s.obj)[index] == nil {
		{{template "new_instance" $elem}}
		(*s.obj)[index] = instance
	}
{{end}}

	return &Storage_{{ $elem.Type }} {
		obj: (*s.obj)[index],
		db: s.db,
		addr: s.addr,
		slot: actual,
		offset: offset,
		numberOfBytes: {{$elem.SolKind.NumberOfBytes}},
		dirty: s.dirty,
	}
}

{{end}} {{/* slice */}}


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
{{else if eq $elemKey.Type "Decimal"}}
	keyBytes := key.BigInt().Bytes()
{{else}}
	UNSUPPORTED {{$elemKey.Type}}
{{end -}}

	hash := crypto.Keccak256Hash(append(keyBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()
	offset := 0

{{- if or (isptr $elemValue) (isslice $elemValue) (ismap $elemValue) (eq $elemValue.Type "BigInt") }}
	if s.obj[key] == nil {
		{{template "new_instance" $elemValue}}
		s.obj[key] = instance
	}
{{end}}

	return &Storage_{{ $elemValue.Type }} {
{{- if or (isptr $elemValue) (ismap $elemValue) (eq $elemValue.Type "BigInt") }}
		obj: s.obj[key],
{{else}}
		obj: &s.obj[key],
{{end -}}
		db: s.db,
		addr: s.addr,
		slot: actual,
		offset: offset,
		numberOfBytes: {{$elemValue.SolKind.NumberOfBytes}},	
		dirty: s.dirty,
	}
}

{{end}} {{/* map */}}

{{end}}{{/* defines */}}


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
{{template "storage_fields" .}}
}
{{end}}


func New(g *Global_t, db StateDB, addr common.Address, slot *big.Int) *Storage_Global_t {
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
	offset := {{$field.Offset}}
	
	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
{{- if or (isptr $field) (isslice $field) (ismap $field) (eq $field.Type "BigInt") }}
	if s.obj.{{$field.Name}} == nil {
		{{template "new_instance" $field}}
		s.obj.{{$field.Name}} = instance
	}
{{end}}
	return &Storage_{{ $field.Type }} {
{{- if or (isptr $field) (ismap $field) (eq $field.Type "BigInt") }}
		obj: s.obj.{{$field.Name}},
{{else}}
		obj: &s.obj.{{$field.Name}},
{{end -}}
		db: s.db,
		addr: s.addr,
		slot: actual,
		offset: offset,
		numberOfBytes: {{$field.SolKind.NumberOfBytes}},
		dirty: s.dirty,
	}
}
	{{end}}
{{end}}


`
