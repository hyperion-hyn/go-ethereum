// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package storage

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// use backtick in text/template: https://github.com/golang/go/issues/18221

type StateDB interface {
	GetState(addr common.Address, hash common.Hash) common.Hash
	SetState(addr common.Address, key, value common.Hash)
}

type StateValues map[common.Hash]common.Hash

func GetStateAsBytes(db StateDB, addr common.Address, slot *big.Int) []byte {
	var retval []byte

	hash := db.GetState(addr, common.BigToHash(slot))
	flag := int(hash[31]) & 0x01
	if flag == 0 {
		// less than 32 bytes
		length := int(hash[31]) / 2
		retval = hash[0:length]
	} else {
		// more than 32 bytes
		length := (hash.Big().Uint64() - 1) / 2
		base := crypto.Keccak256Hash(common.BigToHash(slot).Bytes()).Big()
		for remaining, i := length, uint64(0); remaining <= length; remaining, i = remaining-32, i+1 {
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
		db.SetState(addr, common.BigToHash(slot), common.BigToHash(big.NewInt(int64(length*2+1))))
		base := crypto.Keccak256Hash(common.BigToHash(slot).Bytes()).Big()
		for remaining, i := uint64(length), uint64(0); remaining <= length; remaining, i = remaining-32, i+1 {
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

type Address = common.Address
type Storage_Address struct {
	// Builtin-Type

	// &wrapper.tmplStruct{T:0x7, Name:"Address", Fields:[]*wrapper.tmplField(nil), Type:"common.Address", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x11, Type:(*reflect.rtype)(0x4862e20), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *Address
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Address) Value() common.Address {
	// &wrapper.tmplStruct{T:0x7, Name:"Address", Fields:[]*wrapper.tmplField(nil), Type:"common.Address", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x11, Type:(*reflect.rtype)(0x4862e20), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}

	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	data := hash.Bytes()[32-(s.offset+s.numberOfBytes) : 32-s.offset]

	*s.obj = common.BytesToAddress(data)
	return *s.obj
}

func (s *Storage_Address) SetValue(value common.Address) {
	// &wrapper.tmplStruct{T:0x7, Name:"Address", Fields:[]*wrapper.tmplField(nil), Type:"common.Address", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x11, Type:(*reflect.rtype)(0x4862e20), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	val := value

	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	data := hash.Bytes()
	copy(data[32-(s.offset+s.numberOfBytes):32-s.offset], val[len(val)-s.numberOfBytes:])
	hash.SetBytes(data)
	s.db.SetState(s.addr, common.BigToHash(s.slot), hash)

	*s.obj = value
}

type Bytes = []byte
type Storage_Bytes struct {
	// Builtin-Type

	// &wrapper.tmplStruct{T:0x9, Name:"Bytes", Fields:[]*wrapper.tmplField(nil), Type:"[]byte", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x17, Type:(*reflect.rtype)(0x4795b20), Size:0, T:0x9, NumberOfBytes:0x20, StringKind:"t_bytes_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *Bytes
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Bytes) Value() []byte {
	// &wrapper.tmplStruct{T:0x9, Name:"Bytes", Fields:[]*wrapper.tmplField(nil), Type:"[]byte", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x17, Type:(*reflect.rtype)(0x4795b20), Size:0, T:0x9, NumberOfBytes:0x20, StringKind:"t_bytes_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	rv := GetStateAsBytes(s.db, s.addr, s.slot)
	*s.obj = []byte(rv)
	return *s.obj
}

func (s *Storage_Bytes) SetValue(value []byte) {
	// &wrapper.tmplStruct{T:0x9, Name:"Bytes", Fields:[]*wrapper.tmplField(nil), Type:"[]byte", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x17, Type:(*reflect.rtype)(0x4795b20), Size:0, T:0x9, NumberOfBytes:0x20, StringKind:"t_bytes_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	SetStateAsBytes(s.db, s.addr, s.slot, []byte(value))
	*s.obj = value
}

// Slice_a38c5e05aa47 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x4, Name:"Slice_a38c5e05aa47", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0005307e0)}, Type:"[]*Validator_t", SolKind:abi.Type{Elem:(*abi.Type)(0xc000441170), Kind:0x17, Type:(*reflect.rtype)(0xc000435fc0), Size:0, T:0x4, NumberOfBytes:0x20, StringKind:"t_array(t_struct(Validator_t)6_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
// "[]*struct { PublicKey []uint8 \"json:\\\"PublicKey\\\" storage:\\\"slot=0,offset=0\\\"\"; Coinbase common.Address \"json:\\\"Coinbase\\\" storage:\\\"slot=1,offset=0\\\"\" }"
type Slice_a38c5e05aa47 []*Validator_t

type Storage_Slice_a38c5e05aa47 struct {

	// &wrapper.tmplStruct{T:0x4, Name:"Slice_a38c5e05aa47", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0005307e0)}, Type:"[]*Validator_t", SolKind:abi.Type{Elem:(*abi.Type)(0xc000441170), Kind:0x17, Type:(*reflect.rtype)(0xc000435fc0), Size:0, T:0x4, NumberOfBytes:0x20, StringKind:"t_array(t_struct(Validator_t)6_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *Slice_a38c5e05aa47
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Slice_a38c5e05aa47) Length() int {
	rv := s.db.GetState(s.addr, common.BigToHash(s.slot))
	return int(rv.Big().Int64())
}

func (s *Storage_Slice_a38c5e05aa47) Resize(length int) {
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(big.NewInt(0).SetUint64(uint64(length))))

	slice := make([]*Validator_t, length, length+50)
	copy(slice, *s.obj)
	*s.obj = slice
}

func (s *Storage_Slice_a38c5e05aa47) Get(index int) *Storage_Validator_t {
	// Value: &wrapper.tmplField{Type:"Validator_t", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0xc000435f80), Size:0, T:0x6, NumberOfBytes:0x40, StringKind:"t_struct(Validator_t)6_storage", TupleRawName:"Validator_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc000441200), (*abi.Type)(0xc000441320)}, TupleRawNames:[]string{"PublicKey", "Coinbase"}}, Tag:"", Slot:<nil>, Offset:0}
	length := s.Length()
	if length < index {
		s.Resize(index + 1)
	}

	hash := crypto.Keccak256Hash(common.BigToHash(s.slot).Bytes())
	actual := big.NewInt(0).Add(hash.Big(), big.NewInt(0).SetUint64(uint64(index*(64/32))))
	offset := 0
	if (*s.obj)[index] == nil {

		instance := new(Validator_t)

		(*s.obj)[index] = instance
	}

	return &Storage_Validator_t{
		obj:           (*s.obj)[index],
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 64,
		dirty:         s.dirty,
	}
}

// Committee_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Committee_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0005308c0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Committee_t struct {
	Members Slice_a38c5e05aa47 `json:"members" storage:"slot=0,offset=0"`
}

// Global_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Global_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0005309a0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Global_t struct {
	Committee Committee_t `json:"committee" storage:"slot=0,offset=0"`
}

// Validator_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Validator_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000530620), (*wrapper.tmplField)(0xc000530700)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Validator_t struct {
	PublicKey Bytes   `json:"PublicKey" storage:"slot=0,offset=0"`
	Coinbase  Address `json:"Coinbase" storage:"slot=1,offset=0"`
}

type Storage_Committee_t struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Committee_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0005308c0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *Committee_t
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Global_t struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Global_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0005309a0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *Global_t
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Validator_t struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Validator_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000530620), (*wrapper.tmplField)(0xc000530700)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *Validator_t
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func New(g *Global_t, db StateDB, addr common.Address, slot *big.Int) *Storage_Global_t {
	return &Storage_Global_t{
		obj:   g,
		db:    db,
		addr:  addr,
		slot:  slot,
		dirty: make(StateValues),
	}
}

// Members is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Committee_t) Members() *Storage_Slice_a38c5e05aa47 {
	// Field: &wrapper.tmplField{Type:"Slice_a38c5e05aa47", ElemType:"Slice_a38c5e05aa47", Name:"Members", SolKind:abi.Type{Elem:(*abi.Type)(0xc000441170), Kind:0x17, Type:(*reflect.rtype)(0xc000435fc0), Size:0, T:0x4, NumberOfBytes:0x20, StringKind:"t_array(t_struct(Validator_t)6_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"members\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Members == nil {

		hash := s.db.GetState(s.addr, common.BigToHash(actual))
		instance := make(Slice_a38c5e05aa47, hash.Big().Int64())

		s.obj.Members = instance
	}

	return &Storage_Slice_a38c5e05aa47{
		obj:           &s.obj.Members,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Committee is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Global_t) Committee() *Storage_Committee_t {
	// Field: &wrapper.tmplField{Type:"Committee_t", ElemType:"Committee_t", Name:"Committee", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000425740), Size:0, T:0x6, NumberOfBytes:0x20, StringKind:"t_struct(Committee_t)10_storage", TupleRawName:"Committee_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004415f0)}, TupleRawNames:[]string{"members"}}, Tag:"json:\"committee\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Committee_t{
		obj:           &s.obj.Committee,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// PublicKey is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Validator_t) PublicKey() *Storage_Bytes {
	// Field: &wrapper.tmplField{Type:"Bytes", ElemType:"Bytes", Name:"PublicKey", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x17, Type:(*reflect.rtype)(0x4795b20), Size:0, T:0x9, NumberOfBytes:0x20, StringKind:"t_bytes_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"PublicKey\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.PublicKey == nil {

		hash := s.db.GetState(s.addr, common.BigToHash(actual))
		instance := make(Bytes, hash.Big().Int64())

		s.obj.PublicKey = instance
	}

	return &Storage_Bytes{
		obj:           &s.obj.PublicKey,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Coinbase is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Validator_t) Coinbase() *Storage_Address {
	// Field: &wrapper.tmplField{Type:"Address", ElemType:"Address", Name:"Coinbase", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x11, Type:(*reflect.rtype)(0x4862e20), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Coinbase\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Address{
		obj:           &s.obj.Coinbase,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 20,
		dirty:         s.dirty,
	}
}
