// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package test

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
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
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = math.MaxInt32
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

	// &wrapper.tmplStruct{T:0x7, Name:"Address", Fields:[]*wrapper.tmplField(nil), Type:"common.Address", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           *Address
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Address) Value() common.Address {
	// &wrapper.tmplStruct{T:0x7, Name:"Address", Fields:[]*wrapper.tmplField(nil), Type:"common.Address", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}

	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	data := hash.Bytes()[32-(s.offset+s.numberOfBytes) : 32-s.offset]

	*s.obj = common.BytesToAddress(data)
	return *s.obj
}

func (s *Storage_Address) SetValue(value common.Address) {
	// &wrapper.tmplStruct{T:0x7, Name:"Address", Fields:[]*wrapper.tmplField(nil), Type:"common.Address", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	val := value

	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	data := hash.Bytes()
	copy(data[32-(s.offset+s.numberOfBytes):32-s.offset], val[len(val)-s.numberOfBytes:])
	hash.SetBytes(data)
	s.db.SetState(s.addr, common.BigToHash(s.slot), hash)

	*s.obj = value
}

type BigInt = *big.Int
type Storage_BigInt struct {
	// Builtin-Type

	// &wrapper.tmplStruct{T:0x1, Name:"BigInt", Fields:[]*wrapper.tmplField(nil), Type:"*big.Int", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           BigInt
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_BigInt) Value() *big.Int {
	// &wrapper.tmplStruct{T:0x1, Name:"BigInt", Fields:[]*wrapper.tmplField(nil), Type:"*big.Int", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}

	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	data := hash.Bytes()[32-(s.offset+s.numberOfBytes) : 32-s.offset]

	*s.obj = *(big.NewInt(0).SetBytes(data))
	retval := big.NewInt(0).SetBytes(data)
	return retval
}

func (s *Storage_BigInt) SetValue(value *big.Int) {
	// &wrapper.tmplStruct{T:0x1, Name:"BigInt", Fields:[]*wrapper.tmplField(nil), Type:"*big.Int", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	val := math.PaddedBigBytes(value, 32)

	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	data := hash.Bytes()
	copy(data[32-(s.offset+s.numberOfBytes):32-s.offset], val[len(val)-s.numberOfBytes:])
	hash.SetBytes(data)
	s.db.SetState(s.addr, common.BigToHash(s.slot), hash)

	*s.obj = *(hash.Big())
}

type Bool = bool
type Storage_Bool struct {
	// Builtin-Type

	// &wrapper.tmplStruct{T:0x2, Name:"Bool", Fields:[]*wrapper.tmplField(nil), Type:"bool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x2, numberOfBytes:0x1, stringKind:"t_bool", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           *Bool
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Bool) Value() bool {
	// &wrapper.tmplStruct{T:0x2, Name:"Bool", Fields:[]*wrapper.tmplField(nil), Type:"bool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x2, numberOfBytes:0x1, stringKind:"t_bool", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}

	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	data := hash.Bytes()[32-(s.offset+s.numberOfBytes) : 32-s.offset]

	*s.obj = bool(big.NewInt(0).SetBytes(data).Int64() != 0)
	return *s.obj
}

func (s *Storage_Bool) SetValue(value bool) {
	// &wrapper.tmplStruct{T:0x2, Name:"Bool", Fields:[]*wrapper.tmplField(nil), Type:"bool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x2, numberOfBytes:0x1, stringKind:"t_bool", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	var flag uint
	if value {
		flag = 1
	} else {
		flag = 0
	}

	val := math.PaddedBigBytes(big.NewInt(0).SetInt64(int64(flag)), 32)

	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	data := hash.Bytes()
	copy(data[32-(s.offset+s.numberOfBytes):32-s.offset], val[len(val)-s.numberOfBytes:])
	hash.SetBytes(data)
	s.db.SetState(s.addr, common.BigToHash(s.slot), hash)

	*s.obj = value
}

type Bytes = []uint8
type Storage_Bytes struct {
	// Builtin-Type

	// &wrapper.tmplStruct{T:0x9, Name:"Bytes", Fields:[]*wrapper.tmplField(nil), Type:"[]uint8", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x9, numberOfBytes:0x20, stringKind:"t_bytes_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           *Bytes
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Bytes) Value() []uint8 {
	// &wrapper.tmplStruct{T:0x9, Name:"Bytes", Fields:[]*wrapper.tmplField(nil), Type:"[]uint8", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x9, numberOfBytes:0x20, stringKind:"t_bytes_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	rv := GetStateAsBytes(s.db, s.addr, s.slot)
	*s.obj = []uint8(rv)
	return *s.obj
}

func (s *Storage_Bytes) SetValue(value []uint8) {
	// &wrapper.tmplStruct{T:0x9, Name:"Bytes", Fields:[]*wrapper.tmplField(nil), Type:"[]uint8", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x9, numberOfBytes:0x20, stringKind:"t_bytes_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	SetStateAsBytes(s.db, s.addr, s.slot, []byte(value))
	*s.obj = value
}

type Bytes5 = [5]uint8
type Storage_Bytes5 struct {
	// Builtin-Type

	// &wrapper.tmplStruct{T:0x8, Name:"Bytes5", Fields:[]*wrapper.tmplField(nil), Type:"[5]uint8", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:5, T:0x8, numberOfBytes:0x5, stringKind:"t_bytes5", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           *Bytes5
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Bytes5) Value() [5]uint8 {
	// &wrapper.tmplStruct{T:0x8, Name:"Bytes5", Fields:[]*wrapper.tmplField(nil), Type:"[5]uint8", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:5, T:0x8, numberOfBytes:0x5, stringKind:"t_bytes5", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}

	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	data := hash.Bytes()[32-(s.offset+s.numberOfBytes) : 32-s.offset]

	copy((*s.obj)[:], data[:])
	return *s.obj
}

func (s *Storage_Bytes5) SetValue(value [5]uint8) {
	// &wrapper.tmplStruct{T:0x8, Name:"Bytes5", Fields:[]*wrapper.tmplField(nil), Type:"[5]uint8", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:5, T:0x8, numberOfBytes:0x5, stringKind:"t_bytes5", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	val := value

	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	data := hash.Bytes()
	copy(data[32-(s.offset+s.numberOfBytes):32-s.offset], val[len(val)-s.numberOfBytes:])
	hash.SetBytes(data)
	s.db.SetState(s.addr, common.BigToHash(s.slot), hash)

	copy((*s.obj)[:], value[:])
}

type Decimal = common.Dec
type Storage_Decimal struct {
	// Builtin-Type

	// &wrapper.tmplStruct{T:0xf, Name:"Decimal", Fields:[]*wrapper.tmplField(nil), Type:"common.Dec", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xf, numberOfBytes:0x20, stringKind:"t_struct(Decimal)5_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048b080)}, TupleRawNames:[]string{"f"}, TupleType:(*reflect.rtype)(0xc000124a20)}}
	obj           *Decimal
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Decimal) Value() common.Dec {
	// &wrapper.tmplStruct{T:0xf, Name:"Decimal", Fields:[]*wrapper.tmplField(nil), Type:"common.Dec", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xf, numberOfBytes:0x20, stringKind:"t_struct(Decimal)5_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048b080)}, TupleRawNames:[]string{"f"}, TupleType:(*reflect.rtype)(0xc000124a20)}}
	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	*s.obj = common.NewDecFromBigIntWithPrec(hash.Big(), common.Precision)
	return *s.obj
}

func (s *Storage_Decimal) SetValue(value common.Dec) {
	// &wrapper.tmplStruct{T:0xf, Name:"Decimal", Fields:[]*wrapper.tmplField(nil), Type:"common.Dec", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xf, numberOfBytes:0x20, stringKind:"t_struct(Decimal)5_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048b080)}, TupleRawNames:[]string{"f"}, TupleType:(*reflect.rtype)(0xc000124a20)}}
	hash := value.BigInt()
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(hash))
	*s.obj = value
}

type String = string
type Storage_String struct {
	// Builtin-Type

	// &wrapper.tmplStruct{T:0x3, Name:"String", Fields:[]*wrapper.tmplField(nil), Type:"string", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x3, numberOfBytes:0x20, stringKind:"t_string_memory_ptr", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           *String
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_String) Value() string {
	// &wrapper.tmplStruct{T:0x3, Name:"String", Fields:[]*wrapper.tmplField(nil), Type:"string", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x3, numberOfBytes:0x20, stringKind:"t_string_memory_ptr", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	rv := GetStateAsBytes(s.db, s.addr, s.slot)
	*s.obj = string(rv)
	return *s.obj
}

func (s *Storage_String) SetValue(value string) {
	// &wrapper.tmplStruct{T:0x3, Name:"String", Fields:[]*wrapper.tmplField(nil), Type:"string", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x3, numberOfBytes:0x20, stringKind:"t_string_memory_ptr", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	SetStateAsBytes(s.db, s.addr, s.slot, []byte(value))
	*s.obj = value
}

type Uint32 = uint32
type Storage_Uint32 struct {
	// Builtin-Type

	// &wrapper.tmplStruct{T:0x1, Name:"Uint32", Fields:[]*wrapper.tmplField(nil), Type:"uint32", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:32, T:0x1, numberOfBytes:0x4, stringKind:"t_uint32", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           *Uint32
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Uint32) Value() uint32 {
	// &wrapper.tmplStruct{T:0x1, Name:"Uint32", Fields:[]*wrapper.tmplField(nil), Type:"uint32", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:32, T:0x1, numberOfBytes:0x4, stringKind:"t_uint32", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}

	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	data := hash.Bytes()[32-(s.offset+s.numberOfBytes) : 32-s.offset]

	*s.obj = uint32(big.NewInt(0).SetBytes(data).Uint64())
	return *s.obj
}

func (s *Storage_Uint32) SetValue(value uint32) {
	// &wrapper.tmplStruct{T:0x1, Name:"Uint32", Fields:[]*wrapper.tmplField(nil), Type:"uint32", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:32, T:0x1, numberOfBytes:0x4, stringKind:"t_uint32", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	val := math.PaddedBigBytes(big.NewInt(0).SetUint64(uint64(value)), 32)

	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	data := hash.Bytes()
	copy(data[32-(s.offset+s.numberOfBytes):32-s.offset], val[len(val)-s.numberOfBytes:])
	hash.SetBytes(data)
	s.db.SetState(s.addr, common.BigToHash(s.slot), hash)

	*s.obj = value
}

type Uint64 = uint64
type Storage_Uint64 struct {
	// Builtin-Type

	// &wrapper.tmplStruct{T:0x1, Name:"Uint64", Fields:[]*wrapper.tmplField(nil), Type:"uint64", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:64, T:0x1, numberOfBytes:0x8, stringKind:"t_uint64", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           *Uint64
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Uint64) Value() uint64 {
	// &wrapper.tmplStruct{T:0x1, Name:"Uint64", Fields:[]*wrapper.tmplField(nil), Type:"uint64", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:64, T:0x1, numberOfBytes:0x8, stringKind:"t_uint64", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}

	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	data := hash.Bytes()[32-(s.offset+s.numberOfBytes) : 32-s.offset]

	*s.obj = uint64(big.NewInt(0).SetBytes(data).Uint64())
	return *s.obj
}

func (s *Storage_Uint64) SetValue(value uint64) {
	// &wrapper.tmplStruct{T:0x1, Name:"Uint64", Fields:[]*wrapper.tmplField(nil), Type:"uint64", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:64, T:0x1, numberOfBytes:0x8, stringKind:"t_uint64", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	val := math.PaddedBigBytes(big.NewInt(0).SetUint64(uint64(value)), 32)

	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	data := hash.Bytes()
	copy(data[32-(s.offset+s.numberOfBytes):32-s.offset], val[len(val)-s.numberOfBytes:])
	hash.SetBytes(data)
	s.db.SetState(s.addr, common.BigToHash(s.slot), hash)

	*s.obj = value
}

type Uint8 = uint8
type Storage_Uint8 struct {
	// Builtin-Type

	// &wrapper.tmplStruct{T:0x1, Name:"Uint8", Fields:[]*wrapper.tmplField(nil), Type:"uint8", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:8, T:0x1, numberOfBytes:0x1, stringKind:"uint8", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           *Uint8
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Uint8) Value() uint8 {
	// &wrapper.tmplStruct{T:0x1, Name:"Uint8", Fields:[]*wrapper.tmplField(nil), Type:"uint8", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:8, T:0x1, numberOfBytes:0x1, stringKind:"uint8", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}

	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	data := hash.Bytes()[32-(s.offset+s.numberOfBytes) : 32-s.offset]

	*s.obj = uint8(big.NewInt(0).SetBytes(data).Uint64())
	return *s.obj
}

func (s *Storage_Uint8) SetValue(value uint8) {
	// &wrapper.tmplStruct{T:0x1, Name:"Uint8", Fields:[]*wrapper.tmplField(nil), Type:"uint8", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:8, T:0x1, numberOfBytes:0x1, stringKind:"uint8", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	val := math.PaddedBigBytes(big.NewInt(0).SetUint64(uint64(value)), 32)

	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	data := hash.Bytes()
	copy(data[32-(s.offset+s.numberOfBytes):32-s.offset], val[len(val)-s.numberOfBytes:])
	hash.SetBytes(data)
	s.db.SetState(s.addr, common.BigToHash(s.slot), hash)

	*s.obj = value
}

// Array_0c8aba4f3757 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x5, Name:"Array_0c8aba4f3757", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0001b01a0)}, Type:"[3]Array_bbc2fab83295", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048a600), Size:3, T:0x5, numberOfBytes:0x180, stringKind:"t_array(t_array(t_struct(PendingDelegation_t)83_storage)2_storage)3_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "[3][2]struct { Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=0,offset=0\\\"\"; Epoch *big.Int \"json:\\\"Epoch\\\" storage:\\\"slot=1,offset=0\\\"\" }"
type Array_0c8aba4f3757 [3]Array_bbc2fab83295

type Storage_Array_0c8aba4f3757 struct {

	// &wrapper.tmplStruct{T:0x5, Name:"Array_0c8aba4f3757", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0001b01a0)}, Type:"[3]Array_bbc2fab83295", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048a600), Size:3, T:0x5, numberOfBytes:0x180, stringKind:"t_array(t_array(t_struct(PendingDelegation_t)83_storage)2_storage)3_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           *Array_0c8aba4f3757
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Array_0c8aba4f3757) Length() int {
	// Value: &wrapper.tmplStruct{T:0x5, Name:"Array_0c8aba4f3757", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0001b01a0)}, Type:"[3]Array_bbc2fab83295", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048a600), Size:3, T:0x5, numberOfBytes:0x180, stringKind:"t_array(t_array(t_struct(PendingDelegation_t)83_storage)2_storage)3_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	return len(s.obj)
}

func (s *Storage_Array_0c8aba4f3757) Get(index int) *Storage_Array_bbc2fab83295 {
	// Value: &wrapper.tmplField{Type:"Array_bbc2fab83295", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048a680), Size:2, T:0x5, numberOfBytes:0x80, stringKind:"t_array(t_struct(PendingDelegation_t)83_storage)2_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	base := s.slot

	slotsPerItem := (128 + 31) / 32
	actual := big.NewInt(0).Add(base, big.NewInt(0).SetUint64(uint64(index*slotsPerItem)))
	offset := 0

	return &Storage_Array_bbc2fab83295{
		obj:           &s.obj[index],
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 128,
		dirty:         s.dirty,
	}
}

// Array_19387431c506 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x5, Name:"Array_19387431c506", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000493110)}, Type:"[2]BigInt", SolKind:abi.Type{Elem:(*abi.Type)(0xc0001ac400), Size:2, T:0x5, numberOfBytes:0x40, stringKind:"t_array(t_uint256)2_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "[2]*big.Int"
type Array_19387431c506 [2]BigInt

type Storage_Array_19387431c506 struct {

	// &wrapper.tmplStruct{T:0x5, Name:"Array_19387431c506", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000493110)}, Type:"[2]BigInt", SolKind:abi.Type{Elem:(*abi.Type)(0xc0001ac400), Size:2, T:0x5, numberOfBytes:0x40, stringKind:"t_array(t_uint256)2_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           *Array_19387431c506
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Array_19387431c506) Length() int {
	// Value: &wrapper.tmplStruct{T:0x5, Name:"Array_19387431c506", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000493110)}, Type:"[2]BigInt", SolKind:abi.Type{Elem:(*abi.Type)(0xc0001ac400), Size:2, T:0x5, numberOfBytes:0x40, stringKind:"t_array(t_uint256)2_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	return len(s.obj)
}

func (s *Storage_Array_19387431c506) Get(index int) *Storage_BigInt {
	// Value: &wrapper.tmplField{Type:"BigInt", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	base := s.slot

	itemsPerSlot := 32 / 32
	actual := big.NewInt(0).Add(base, big.NewInt(0).SetUint64(uint64(index/itemsPerSlot)))
	offset := ((index % itemsPerSlot) * 32)

	if s.obj[index] == nil {

		instance := big.NewInt(0)

		s.obj[index] = instance
	}

	return &Storage_BigInt{
		obj:           s.obj[index],
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Array_35cbea3b1b88 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x5, Name:"Array_35cbea3b1b88", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000492a90)}, Type:"[9]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048bc80), Size:9, T:0x5, numberOfBytes:0x20, stringKind:"t_array(t_bytes1)9_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "[9]uint8"
type Array_35cbea3b1b88 [9]Uint8

type Storage_Array_35cbea3b1b88 struct {

	// &wrapper.tmplStruct{T:0x5, Name:"Array_35cbea3b1b88", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000492a90)}, Type:"[9]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048bc80), Size:9, T:0x5, numberOfBytes:0x20, stringKind:"t_array(t_bytes1)9_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           *Array_35cbea3b1b88
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Array_35cbea3b1b88) Value() [9]Uint8 {
	// &wrapper.tmplStruct{T:0x5, Name:"Array_35cbea3b1b88", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000492a90)}, Type:"[9]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048bc80), Size:9, T:0x5, numberOfBytes:0x20, stringKind:"t_array(t_bytes1)9_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
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
			(*s.obj)[offset+j] = data[31-j]
		}
	}
	return *s.obj
}

func (s *Storage_Array_35cbea3b1b88) SetValue(value [9]Uint8) {
	// &wrapper.tmplStruct{T:0x5, Name:"Array_35cbea3b1b88", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000492a90)}, Type:"[9]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048bc80), Size:9, T:0x5, numberOfBytes:0x20, stringKind:"t_array(t_bytes1)9_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
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

// Array_52c9a19dddcb is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x5, Name:"Array_52c9a19dddcb", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000492dd0)}, Type:"[4]Uint32", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048bf80), Size:4, T:0x5, numberOfBytes:0x20, stringKind:"t_array(t_uint32)4_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "[4]uint32"
type Array_52c9a19dddcb [4]Uint32

type Storage_Array_52c9a19dddcb struct {

	// &wrapper.tmplStruct{T:0x5, Name:"Array_52c9a19dddcb", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000492dd0)}, Type:"[4]Uint32", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048bf80), Size:4, T:0x5, numberOfBytes:0x20, stringKind:"t_array(t_uint32)4_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           *Array_52c9a19dddcb
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Array_52c9a19dddcb) Length() int {
	// Value: &wrapper.tmplStruct{T:0x5, Name:"Array_52c9a19dddcb", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000492dd0)}, Type:"[4]Uint32", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048bf80), Size:4, T:0x5, numberOfBytes:0x20, stringKind:"t_array(t_uint32)4_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	return len(s.obj)
}

func (s *Storage_Array_52c9a19dddcb) Get(index int) *Storage_Uint32 {
	// Value: &wrapper.tmplField{Type:"Uint32", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:32, T:0x1, numberOfBytes:0x4, stringKind:"t_uint32", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	base := s.slot

	itemsPerSlot := 32 / 4
	actual := big.NewInt(0).Add(base, big.NewInt(0).SetUint64(uint64(index/itemsPerSlot)))
	offset := ((index % itemsPerSlot) * 4)

	return &Storage_Uint32{
		obj:           &s.obj[index],
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 4,
		dirty:         s.dirty,
	}
}

// Array_61783fec1216 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x5, Name:"Array_61783fec1216", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000492f70)}, Type:"[12]BigInt", SolKind:abi.Type{Elem:(*abi.Type)(0xc0001ac200), Size:12, T:0x5, numberOfBytes:0x60, stringKind:"t_array(t_uint48)12_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "[12]*big.Int"
type Array_61783fec1216 [12]BigInt

type Storage_Array_61783fec1216 struct {

	// &wrapper.tmplStruct{T:0x5, Name:"Array_61783fec1216", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000492f70)}, Type:"[12]BigInt", SolKind:abi.Type{Elem:(*abi.Type)(0xc0001ac200), Size:12, T:0x5, numberOfBytes:0x60, stringKind:"t_array(t_uint48)12_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           *Array_61783fec1216
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Array_61783fec1216) Length() int {
	// Value: &wrapper.tmplStruct{T:0x5, Name:"Array_61783fec1216", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000492f70)}, Type:"[12]BigInt", SolKind:abi.Type{Elem:(*abi.Type)(0xc0001ac200), Size:12, T:0x5, numberOfBytes:0x60, stringKind:"t_array(t_uint48)12_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	return len(s.obj)
}

func (s *Storage_Array_61783fec1216) Get(index int) *Storage_BigInt {
	// Value: &wrapper.tmplField{Type:"BigInt", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:48, T:0x1, numberOfBytes:0x6, stringKind:"t_uint48", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	base := s.slot

	itemsPerSlot := 32 / 6
	actual := big.NewInt(0).Add(base, big.NewInt(0).SetUint64(uint64(index/itemsPerSlot)))
	offset := ((index % itemsPerSlot) * 6)

	if s.obj[index] == nil {

		instance := big.NewInt(0)

		s.obj[index] = instance
	}

	return &Storage_BigInt{
		obj:           s.obj[index],
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 6,
		dirty:         s.dirty,
	}
}

// Array_9d1a033deeae is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x5, Name:"Array_9d1a033deeae", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000492c30)}, Type:"[300]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048be00), Size:300, T:0x5, numberOfBytes:0x140, stringKind:"t_array(t_bytes1)300_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "[300]uint8"
type Array_9d1a033deeae [300]Uint8

type Storage_Array_9d1a033deeae struct {

	// &wrapper.tmplStruct{T:0x5, Name:"Array_9d1a033deeae", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000492c30)}, Type:"[300]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048be00), Size:300, T:0x5, numberOfBytes:0x140, stringKind:"t_array(t_bytes1)300_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           *Array_9d1a033deeae
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Array_9d1a033deeae) Value() [300]Uint8 {
	// &wrapper.tmplStruct{T:0x5, Name:"Array_9d1a033deeae", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000492c30)}, Type:"[300]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048be00), Size:300, T:0x5, numberOfBytes:0x140, stringKind:"t_array(t_bytes1)300_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
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
			(*s.obj)[offset+j] = data[31-j]
		}
	}
	return *s.obj
}

func (s *Storage_Array_9d1a033deeae) SetValue(value [300]Uint8) {
	// &wrapper.tmplStruct{T:0x5, Name:"Array_9d1a033deeae", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000492c30)}, Type:"[300]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048be00), Size:300, T:0x5, numberOfBytes:0x140, stringKind:"t_array(t_bytes1)300_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
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

// Array_b7edf2cad8cd is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x5, Name:"Array_b7edf2cad8cd", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0004928f0)}, Type:"[1]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048ba80), Size:1, T:0x5, numberOfBytes:0x20, stringKind:"t_array(t_bytes1)1_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "[1]uint8"
type Array_b7edf2cad8cd [1]Uint8

type Storage_Array_b7edf2cad8cd struct {

	// &wrapper.tmplStruct{T:0x5, Name:"Array_b7edf2cad8cd", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0004928f0)}, Type:"[1]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048ba80), Size:1, T:0x5, numberOfBytes:0x20, stringKind:"t_array(t_bytes1)1_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           *Array_b7edf2cad8cd
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Array_b7edf2cad8cd) Value() [1]Uint8 {
	// &wrapper.tmplStruct{T:0x5, Name:"Array_b7edf2cad8cd", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0004928f0)}, Type:"[1]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048ba80), Size:1, T:0x5, numberOfBytes:0x20, stringKind:"t_array(t_bytes1)1_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
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
			(*s.obj)[offset+j] = data[31-j]
		}
	}
	return *s.obj
}

func (s *Storage_Array_b7edf2cad8cd) SetValue(value [1]Uint8) {
	// &wrapper.tmplStruct{T:0x5, Name:"Array_b7edf2cad8cd", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0004928f0)}, Type:"[1]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048ba80), Size:1, T:0x5, numberOfBytes:0x20, stringKind:"t_array(t_bytes1)1_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
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

// Array_bbc2fab83295 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x5, Name:"Array_bbc2fab83295", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000493ee0)}, Type:"[2]PendingDelegation_t", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048a680), Size:2, T:0x5, numberOfBytes:0x80, stringKind:"t_array(t_struct(PendingDelegation_t)83_storage)2_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "[2]struct { Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=0,offset=0\\\"\"; Epoch *big.Int \"json:\\\"Epoch\\\" storage:\\\"slot=1,offset=0\\\"\" }"
type Array_bbc2fab83295 [2]PendingDelegation_t

type Storage_Array_bbc2fab83295 struct {

	// &wrapper.tmplStruct{T:0x5, Name:"Array_bbc2fab83295", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000493ee0)}, Type:"[2]PendingDelegation_t", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048a680), Size:2, T:0x5, numberOfBytes:0x80, stringKind:"t_array(t_struct(PendingDelegation_t)83_storage)2_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           *Array_bbc2fab83295
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Array_bbc2fab83295) Length() int {
	// Value: &wrapper.tmplStruct{T:0x5, Name:"Array_bbc2fab83295", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000493ee0)}, Type:"[2]PendingDelegation_t", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048a680), Size:2, T:0x5, numberOfBytes:0x80, stringKind:"t_array(t_struct(PendingDelegation_t)83_storage)2_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	return len(s.obj)
}

func (s *Storage_Array_bbc2fab83295) Get(index int) *Storage_PendingDelegation_t {
	// Value: &wrapper.tmplField{Type:"PendingDelegation_t", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(PendingDelegation_t)83_storage", TupleRawName:"PendingDelegation_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048a700), (*abi.Type)(0xc00048a780)}, TupleRawNames:[]string{"Amount", "Epoch"}, TupleType:(*reflect.rtype)(0xc000124600)}, Tag:"", Slot:<nil>, Offset:0}
	base := s.slot

	slotsPerItem := (64 + 31) / 32
	actual := big.NewInt(0).Add(base, big.NewInt(0).SetUint64(uint64(index*slotsPerItem)))
	offset := 0

	return &Storage_PendingDelegation_t{
		obj:           &s.obj[index],
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 64,
		dirty:         s.dirty,
	}
}

// Map_3783ff483121 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xe, Name:"Map_3783ff483121", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0001b11e0), (*wrapper.tmplField)(0xc0001b12b0)}, Type:"map[Address]*Bool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0001adc00), (*abi.Type)(0xc0001add00)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "map[common.Address]*bool"
type Map_3783ff483121 map[Address]*Bool

type Storage_Map_3783ff483121 struct {

	// &wrapper.tmplStruct{T:0xe, Name:"Map_3783ff483121", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0001b11e0), (*wrapper.tmplField)(0xc0001b12b0)}, Type:"map[Address]*Bool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0001adc00), (*abi.Type)(0xc0001add00)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           Map_3783ff483121
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Map_3783ff483121) Get(key Address) *Storage_Bool {
	// Key:	&wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"Bool", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc0001adc80), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	keyBytes := key.Hash().Bytes()
	hash := crypto.Keccak256Hash(append(keyBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()
	offset := 0
	if s.obj[key] == nil {

		instance := new(Bool)

		s.obj[key] = instance
	}

	return &Storage_Bool{
		obj:           s.obj[key],
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 1,
		dirty:         s.dirty,
	}
}

// Map_624bcc0df495 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xe, Name:"Map_624bcc0df495", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012add0), (*wrapper.tmplField)(0xc00012aea0)}, Type:"map[String]*Bool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_string_memory_ptr,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0001ad880), (*abi.Type)(0xc0001ada00)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "map[string]*bool"
type Map_624bcc0df495 map[String]*Bool

type Storage_Map_624bcc0df495 struct {

	// &wrapper.tmplStruct{T:0xe, Name:"Map_624bcc0df495", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012add0), (*wrapper.tmplField)(0xc00012aea0)}, Type:"map[String]*Bool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_string_memory_ptr,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0001ad880), (*abi.Type)(0xc0001ada00)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           Map_624bcc0df495
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Map_624bcc0df495) Get(key String) *Storage_Bool {
	// Key:	&wrapper.tmplField{Type:"String", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x3, numberOfBytes:0x20, stringKind:"t_string_memory_ptr", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"Bool", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc0001ad980), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	keyBytes := []byte(key)
	hash := crypto.Keccak256Hash(append(keyBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()
	offset := 0
	if s.obj[key] == nil {

		instance := new(Bool)

		s.obj[key] = instance
	}

	return &Storage_Bool{
		obj:           s.obj[key],
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 1,
		dirty:         s.dirty,
	}
}

// Map_cf7837dafed3 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xe, Name:"Map_cf7837dafed3", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012b110), (*wrapper.tmplField)(0xc0001b1040)}, Type:"map[Uint64]*Map3NodeSnapshot_t", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_uint64,t_struct(Map3NodeSnapshot_t)148_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000c8400), (*abi.Type)(0xc0000c8500)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "map[uint64]*struct { Map3Nodes map[common.Address]*struct { Map3Node struct { Actived bool \"json:\\\"Actived\\\" storage:\\\"slot=0,offset=0\\\"\"; NodeAddress common.Address \"json:\\\"NodeAddress\\\" storage:\\\"slot=0,offset=1\\\"\"; InitiatorAddress common.Address \"json:\\\"InitiatorAddress\\\" storage:\\\"slot=1,offset=0\\\"\"; NodeKeys []uint8 \"json:\\\"NodeKeys\\\" storage:\\\"slot=2,offset=0\\\"\"; Commission struct { CommissionRates struct { Rate *big.Int \"json:\\\"Rate\\\" storage:\\\"slot=0,offset=0\\\"\"; MaxRate struct { F *big.Int \"json:\\\"f\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"MaxRate\\\" storage:\\\"slot=1,offset=0\\\"\"; MaxChangeRate *big.Int \"json:\\\"MaxChangeRate\\\" storage:\\\"slot=2,offset=0\\\"\" } \"json:\\\"CommissionRates\\\" storage:\\\"slot=0,offset=0\\\"\"; UpdateHeight *big.Int \"json:\\\"UpdateHeight\\\" storage:\\\"slot=3,offset=0\\\"\" } \"json:\\\"Commission\\\" storage:\\\"slot=3,offset=0\\\"\"; Description struct { Name string \"json:\\\"Name\\\" storage:\\\"slot=0,offset=0\\\"\"; Identity string \"json:\\\"Identity\\\" storage:\\\"slot=1,offset=0\\\"\"; WebSite string \"json:\\\"WebSite\\\" storage:\\\"slot=2,offset=0\\\"\"; SecurityContract string \"json:\\\"SecurityContract\\\" storage:\\\"slot=3,offset=0\\\"\"; Details string \"json:\\\"Details\\\" storage:\\\"slot=4,offset=0\\\"\"; Frozen uint8 \"json:\\\"Frozen\\\" storage:\\\"slot=5,offset=0\\\"\"; Serial [5]uint8 \"json:\\\"Serial\\\" storage:\\\"slot=5,offset=1\\\"\"; Destroyed bool \"json:\\\"Destroyed\\\" storage:\\\"slot=5,offset=6\\\"\"; Flag [1]uint8 \"json:\\\"Flag\\\" storage:\\\"slot=6,offset=0\\\"\"; Symbol [9]uint8 \"json:\\\"Symbol\\\" storage:\\\"slot=7,offset=0\\\"\"; Signature [300]uint8 \"json:\\\"Signature\\\" storage:\\\"slot=8,offset=0\\\"\"; Feature [4]uint32 \"json:\\\"Feature\\\" storage:\\\"slot=18,offset=0\\\"\"; Mac [12]*big.Int \"json:\\\"Mac\\\" storage:\\\"slot=19,offset=0\\\"\"; Version [2]*big.Int \"json:\\\"Version\\\" storage:\\\"slot=22,offset=0\\\"\" } \"json:\\\"Description\\\" storage:\\\"slot=7,offset=0\\\"\"; SplittedFrom common.Address \"json:\\\"SplittedFrom\\\" storage:\\\"slot=31,offset=0\\\"\" } \"json:\\\"Map3Node\\\" storage:\\\"slot=0,offset=0\\\"\"; Microdelegations map[common.Address]*struct { DelegatorAddress common.Address \"json:\\\"DelegatorAddress\\\" storage:\\\"slot=0,offset=0\\\"\"; Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=1,offset=0\\\"\"; Reward *big.Int \"json:\\\"Reward\\\" storage:\\\"slot=2,offset=0\\\"\"; Undelegations []*struct { Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=0,offset=0\\\"\"; Epoch *big.Int \"json:\\\"Epoch\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"Undelegations\\\" storage:\\\"slot=3,offset=0\\\"\"; PendingDelegations []*struct { Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=0,offset=0\\\"\"; Epoch *big.Int \"json:\\\"Epoch\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"PendingDelegations\\\" storage:\\\"slot=4,offset=0\\\"\"; PendingDelegationsfixed [2]struct { Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=0,offset=0\\\"\"; Epoch *big.Int \"json:\\\"Epoch\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"PendingDelegationsfixed\\\" storage:\\\"slot=5,offset=0\\\"\"; PendingDelegationsfixed2dimension [3][2]struct { Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=0,offset=0\\\"\"; Epoch *big.Int \"json:\\\"Epoch\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"PendingDelegationsfixed2dimension\\\" storage:\\\"slot=9,offset=0\\\"\"; AutoRenew bool \"json:\\\"AutoRenew\\\" storage:\\\"slot=21,offset=0\\\"\" } \"json:\\\"Microdelegations\\\" storage:\\\"slot=32,offset=0\\\"\"; RedelegationReference struct { ValidatorAddress common.Address \"json:\\\"ValidatorAddress\\\" storage:\\\"slot=0,offset=0\\\"\"; ReleasedTotalDelegation *big.Int \"json:\\\"ReleasedTotalDelegation\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"RedelegationReference\\\" storage:\\\"slot=33,offset=0\\\"\"; AccumulatedReward *big.Int \"json:\\\"AccumulatedReward\\\" storage:\\\"slot=35,offset=0\\\"\"; NodeState struct { Status uint8 \"json:\\\"Status\\\" storage:\\\"slot=0,offset=0\\\"\"; NodeAge *big.Int \"json:\\\"NodeAge\\\" storage:\\\"slot=1,offset=0\\\"\"; CreationEpoch *big.Int \"json:\\\"CreationEpoch\\\" storage:\\\"slot=2,offset=0\\\"\"; ActivationEpoch *big.Int \"json:\\\"ActivationEpoch\\\" storage:\\\"slot=3,offset=0\\\"\"; ReleaseEpoch *big.Int \"json:\\\"ReleaseEpoch\\\" storage:\\\"slot=4,offset=0\\\"\" } \"json:\\\"nodeState\\\" storage:\\\"slot=36,offset=0\\\"\"; TotalDelegation *big.Int \"json:\\\"TotalDelegation\\\" storage:\\\"slot=41,offset=0\\\"\"; TotalPendingDelegation *big.Int \"json:\\\"TotalPendingDelegation\\\" storage:\\\"slot=42,offset=0\\\"\" } \"json:\\\"Map3Nodes\\\" storage:\\\"slot=0,offset=0\\\"\"; Epoch *big.Int \"json:\\\"Epoch\\\" storage:\\\"slot=1,offset=0\\\"\" }"
type Map_cf7837dafed3 map[Uint64]*Map3NodeSnapshot_t

type Storage_Map_cf7837dafed3 struct {

	// &wrapper.tmplStruct{T:0xe, Name:"Map_cf7837dafed3", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012b110), (*wrapper.tmplField)(0xc0001b1040)}, Type:"map[Uint64]*Map3NodeSnapshot_t", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_uint64,t_struct(Map3NodeSnapshot_t)148_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000c8400), (*abi.Type)(0xc0000c8500)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           Map_cf7837dafed3
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Map_cf7837dafed3) Get(key Uint64) *Storage_Map3NodeSnapshot_t {
	// Key:	&wrapper.tmplField{Type:"Uint64", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:64, T:0x1, numberOfBytes:0x8, stringKind:"t_uint64", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"Map3NodeSnapshot_t", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000c8480), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	keyBytes := big.NewInt(0).SetUint64(uint64(key)).Bytes()
	hash := crypto.Keccak256Hash(append(keyBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()
	offset := 0
	if s.obj[key] == nil {

		instance := new(Map3NodeSnapshot_t)

		s.obj[key] = instance
	}

	return &Storage_Map3NodeSnapshot_t{
		obj:           s.obj[key],
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 64,
		dirty:         s.dirty,
	}
}

// Map_db3cea57c803 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xe, Name:"Map_db3cea57c803", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012b1e0), (*wrapper.tmplField)(0xc0001b0dd0)}, Type:"map[Address]*Map3NodeWrapper_t", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(Map3NodeWrapper_t)141_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048a980), (*abi.Type)(0xc0001ad780)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "map[common.Address]*struct { Map3Node struct { Actived bool \"json:\\\"Actived\\\" storage:\\\"slot=0,offset=0\\\"\"; NodeAddress common.Address \"json:\\\"NodeAddress\\\" storage:\\\"slot=0,offset=1\\\"\"; InitiatorAddress common.Address \"json:\\\"InitiatorAddress\\\" storage:\\\"slot=1,offset=0\\\"\"; NodeKeys []uint8 \"json:\\\"NodeKeys\\\" storage:\\\"slot=2,offset=0\\\"\"; Commission struct { CommissionRates struct { Rate *big.Int \"json:\\\"Rate\\\" storage:\\\"slot=0,offset=0\\\"\"; MaxRate struct { F *big.Int \"json:\\\"f\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"MaxRate\\\" storage:\\\"slot=1,offset=0\\\"\"; MaxChangeRate *big.Int \"json:\\\"MaxChangeRate\\\" storage:\\\"slot=2,offset=0\\\"\" } \"json:\\\"CommissionRates\\\" storage:\\\"slot=0,offset=0\\\"\"; UpdateHeight *big.Int \"json:\\\"UpdateHeight\\\" storage:\\\"slot=3,offset=0\\\"\" } \"json:\\\"Commission\\\" storage:\\\"slot=3,offset=0\\\"\"; Description struct { Name string \"json:\\\"Name\\\" storage:\\\"slot=0,offset=0\\\"\"; Identity string \"json:\\\"Identity\\\" storage:\\\"slot=1,offset=0\\\"\"; WebSite string \"json:\\\"WebSite\\\" storage:\\\"slot=2,offset=0\\\"\"; SecurityContract string \"json:\\\"SecurityContract\\\" storage:\\\"slot=3,offset=0\\\"\"; Details string \"json:\\\"Details\\\" storage:\\\"slot=4,offset=0\\\"\"; Frozen uint8 \"json:\\\"Frozen\\\" storage:\\\"slot=5,offset=0\\\"\"; Serial [5]uint8 \"json:\\\"Serial\\\" storage:\\\"slot=5,offset=1\\\"\"; Destroyed bool \"json:\\\"Destroyed\\\" storage:\\\"slot=5,offset=6\\\"\"; Flag [1]uint8 \"json:\\\"Flag\\\" storage:\\\"slot=6,offset=0\\\"\"; Symbol [9]uint8 \"json:\\\"Symbol\\\" storage:\\\"slot=7,offset=0\\\"\"; Signature [300]uint8 \"json:\\\"Signature\\\" storage:\\\"slot=8,offset=0\\\"\"; Feature [4]uint32 \"json:\\\"Feature\\\" storage:\\\"slot=18,offset=0\\\"\"; Mac [12]*big.Int \"json:\\\"Mac\\\" storage:\\\"slot=19,offset=0\\\"\"; Version [2]*big.Int \"json:\\\"Version\\\" storage:\\\"slot=22,offset=0\\\"\" } \"json:\\\"Description\\\" storage:\\\"slot=7,offset=0\\\"\"; SplittedFrom common.Address \"json:\\\"SplittedFrom\\\" storage:\\\"slot=31,offset=0\\\"\" } \"json:\\\"Map3Node\\\" storage:\\\"slot=0,offset=0\\\"\"; Microdelegations map[common.Address]*struct { DelegatorAddress common.Address \"json:\\\"DelegatorAddress\\\" storage:\\\"slot=0,offset=0\\\"\"; Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=1,offset=0\\\"\"; Reward *big.Int \"json:\\\"Reward\\\" storage:\\\"slot=2,offset=0\\\"\"; Undelegations []*struct { Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=0,offset=0\\\"\"; Epoch *big.Int \"json:\\\"Epoch\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"Undelegations\\\" storage:\\\"slot=3,offset=0\\\"\"; PendingDelegations []*struct { Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=0,offset=0\\\"\"; Epoch *big.Int \"json:\\\"Epoch\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"PendingDelegations\\\" storage:\\\"slot=4,offset=0\\\"\"; PendingDelegationsfixed [2]struct { Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=0,offset=0\\\"\"; Epoch *big.Int \"json:\\\"Epoch\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"PendingDelegationsfixed\\\" storage:\\\"slot=5,offset=0\\\"\"; PendingDelegationsfixed2dimension [3][2]struct { Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=0,offset=0\\\"\"; Epoch *big.Int \"json:\\\"Epoch\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"PendingDelegationsfixed2dimension\\\" storage:\\\"slot=9,offset=0\\\"\"; AutoRenew bool \"json:\\\"AutoRenew\\\" storage:\\\"slot=21,offset=0\\\"\" } \"json:\\\"Microdelegations\\\" storage:\\\"slot=32,offset=0\\\"\"; RedelegationReference struct { ValidatorAddress common.Address \"json:\\\"ValidatorAddress\\\" storage:\\\"slot=0,offset=0\\\"\"; ReleasedTotalDelegation *big.Int \"json:\\\"ReleasedTotalDelegation\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"RedelegationReference\\\" storage:\\\"slot=33,offset=0\\\"\"; AccumulatedReward *big.Int \"json:\\\"AccumulatedReward\\\" storage:\\\"slot=35,offset=0\\\"\"; NodeState struct { Status uint8 \"json:\\\"Status\\\" storage:\\\"slot=0,offset=0\\\"\"; NodeAge *big.Int \"json:\\\"NodeAge\\\" storage:\\\"slot=1,offset=0\\\"\"; CreationEpoch *big.Int \"json:\\\"CreationEpoch\\\" storage:\\\"slot=2,offset=0\\\"\"; ActivationEpoch *big.Int \"json:\\\"ActivationEpoch\\\" storage:\\\"slot=3,offset=0\\\"\"; ReleaseEpoch *big.Int \"json:\\\"ReleaseEpoch\\\" storage:\\\"slot=4,offset=0\\\"\" } \"json:\\\"nodeState\\\" storage:\\\"slot=36,offset=0\\\"\"; TotalDelegation *big.Int \"json:\\\"TotalDelegation\\\" storage:\\\"slot=41,offset=0\\\"\"; TotalPendingDelegation *big.Int \"json:\\\"TotalPendingDelegation\\\" storage:\\\"slot=42,offset=0\\\"\" }"
type Map_db3cea57c803 map[Address]*Map3NodeWrapper_t

type Storage_Map_db3cea57c803 struct {

	// &wrapper.tmplStruct{T:0xe, Name:"Map_db3cea57c803", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012b1e0), (*wrapper.tmplField)(0xc0001b0dd0)}, Type:"map[Address]*Map3NodeWrapper_t", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(Map3NodeWrapper_t)141_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048a980), (*abi.Type)(0xc0001ad780)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           Map_db3cea57c803
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Map_db3cea57c803) Get(key Address) *Storage_Map3NodeWrapper_t {
	// Key:	&wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"Map3NodeWrapper_t", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048aa80), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	keyBytes := key.Hash().Bytes()
	hash := crypto.Keccak256Hash(append(keyBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()
	offset := 0
	if s.obj[key] == nil {

		instance := new(Map3NodeWrapper_t)

		s.obj[key] = instance
	}

	return &Storage_Map3NodeWrapper_t{
		obj:           s.obj[key],
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 1376,
		dirty:         s.dirty,
	}
}

// Map_dfd3c4f07730 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xe, Name:"Map_dfd3c4f07730", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0001b15f0), (*wrapper.tmplField)(0xc0001b16c0)}, Type:"map[Address]Map_3783ff483121", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_mapping(t_address,t_bool))", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0001adb00), (*abi.Type)(0xc0001adb80)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "map[common.Address]map[common.Address]*bool"
type Map_dfd3c4f07730 map[Address]Map_3783ff483121

type Storage_Map_dfd3c4f07730 struct {

	// &wrapper.tmplStruct{T:0xe, Name:"Map_dfd3c4f07730", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0001b15f0), (*wrapper.tmplField)(0xc0001b16c0)}, Type:"map[Address]Map_3783ff483121", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_mapping(t_address,t_bool))", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0001adb00), (*abi.Type)(0xc0001adb80)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           Map_dfd3c4f07730
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Map_dfd3c4f07730) Get(key Address) *Storage_Map_3783ff483121 {
	// Key:	&wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"Map_3783ff483121", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0001adc00), (*abi.Type)(0xc0001add00)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	keyBytes := key.Hash().Bytes()
	hash := crypto.Keccak256Hash(append(keyBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()
	offset := 0
	if s.obj[key] == nil {

		instance := make(Map_3783ff483121)

		s.obj[key] = instance
	}

	return &Storage_Map_3783ff483121{
		obj:           s.obj[key],
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Map_f6daa253cd33 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xe, Name:"Map_f6daa253cd33", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000493520), (*wrapper.tmplField)(0xc0001b0410)}, Type:"map[Address]*Microdelegation_t", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(Microdelegation_t)108_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0001ac700), (*abi.Type)(0xc0001ad180)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "map[common.Address]*struct { DelegatorAddress common.Address \"json:\\\"DelegatorAddress\\\" storage:\\\"slot=0,offset=0\\\"\"; Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=1,offset=0\\\"\"; Reward *big.Int \"json:\\\"Reward\\\" storage:\\\"slot=2,offset=0\\\"\"; Undelegations []*struct { Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=0,offset=0\\\"\"; Epoch *big.Int \"json:\\\"Epoch\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"Undelegations\\\" storage:\\\"slot=3,offset=0\\\"\"; PendingDelegations []*struct { Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=0,offset=0\\\"\"; Epoch *big.Int \"json:\\\"Epoch\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"PendingDelegations\\\" storage:\\\"slot=4,offset=0\\\"\"; PendingDelegationsfixed [2]struct { Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=0,offset=0\\\"\"; Epoch *big.Int \"json:\\\"Epoch\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"PendingDelegationsfixed\\\" storage:\\\"slot=5,offset=0\\\"\"; PendingDelegationsfixed2dimension [3][2]struct { Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=0,offset=0\\\"\"; Epoch *big.Int \"json:\\\"Epoch\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"PendingDelegationsfixed2dimension\\\" storage:\\\"slot=9,offset=0\\\"\"; AutoRenew bool \"json:\\\"AutoRenew\\\" storage:\\\"slot=21,offset=0\\\"\" }"
type Map_f6daa253cd33 map[Address]*Microdelegation_t

type Storage_Map_f6daa253cd33 struct {

	// &wrapper.tmplStruct{T:0xe, Name:"Map_f6daa253cd33", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000493520), (*wrapper.tmplField)(0xc0001b0410)}, Type:"map[Address]*Microdelegation_t", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(Microdelegation_t)108_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0001ac700), (*abi.Type)(0xc0001ad180)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           Map_f6daa253cd33
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Map_f6daa253cd33) Get(key Address) *Storage_Microdelegation_t {
	// Key:	&wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"Microdelegation_t", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc0001ac780), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	keyBytes := key.Hash().Bytes()
	hash := crypto.Keccak256Hash(append(keyBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()
	offset := 0
	if s.obj[key] == nil {

		instance := new(Microdelegation_t)

		s.obj[key] = instance
	}

	return &Storage_Microdelegation_t{
		obj:           s.obj[key],
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 704,
		dirty:         s.dirty,
	}
}

// Map_f75907254fb8 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xe, Name:"Map_f75907254fb8", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0001b1a00), (*wrapper.tmplField)(0xc0001b1ad0)}, Type:"map[String]BigInt", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_string_memory_ptr,t_uint256)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0001ade80), (*abi.Type)(0xc0001adf00)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "map[string]*big.Int"
type Map_f75907254fb8 map[String]BigInt

type Storage_Map_f75907254fb8 struct {

	// &wrapper.tmplStruct{T:0xe, Name:"Map_f75907254fb8", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0001b1a00), (*wrapper.tmplField)(0xc0001b1ad0)}, Type:"map[String]BigInt", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_string_memory_ptr,t_uint256)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0001ade80), (*abi.Type)(0xc0001adf00)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           Map_f75907254fb8
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Map_f75907254fb8) Get(key String) *Storage_BigInt {
	// Key:	&wrapper.tmplField{Type:"String", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x3, numberOfBytes:0x20, stringKind:"t_string_memory_ptr", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"BigInt", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	keyBytes := []byte(key)
	hash := crypto.Keccak256Hash(append(keyBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()
	offset := 0
	if s.obj[key] == nil {

		instance := big.NewInt(0)

		s.obj[key] = instance
	}

	return &Storage_BigInt{
		obj:           s.obj[key],
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Slice_2f14e711180c is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x4, Name:"Slice_2f14e711180c", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0001b1110)}, Type:"[]BigInt", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000c8280), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_uint48)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "[]*big.Int"
type Slice_2f14e711180c []BigInt

type Storage_Slice_2f14e711180c struct {

	// &wrapper.tmplStruct{T:0x4, Name:"Slice_2f14e711180c", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0001b1110)}, Type:"[]BigInt", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000c8280), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_uint48)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           *Slice_2f14e711180c
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Slice_2f14e711180c) Length() int {
	rv := s.db.GetState(s.addr, common.BigToHash(s.slot))
	return int(rv.Big().Int64())
}

func (s *Storage_Slice_2f14e711180c) Resize(length int) {
	// Value: &wrapper.tmplField{Type:"BigInt", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:48, T:0x1, numberOfBytes:0x6, stringKind:"t_uint48", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(big.NewInt(0).SetUint64(uint64(length))))

	_length := length
	_capacity := length + 50

	slice := make([]BigInt, _length, _capacity)

	copy(slice, *s.obj)
	*s.obj = slice
}

func (s *Storage_Slice_2f14e711180c) Get(index int) *Storage_BigInt {
	// Value: &wrapper.tmplField{Type:"BigInt", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:48, T:0x1, numberOfBytes:0x6, stringKind:"t_uint48", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	length := s.Length()
	if length <= index {
		s.Resize(index + 1)
	}

	base := crypto.Keccak256Hash(common.BigToHash(s.slot).Bytes()).Big()

	itemsPerSlot := 32 / 6
	actual := big.NewInt(0).Add(base, big.NewInt(0).SetUint64(uint64(index/itemsPerSlot)))
	offset := ((index % itemsPerSlot) * 6)

	if (*s.obj)[index] == nil {

		instance := big.NewInt(0)

		(*s.obj)[index] = instance
	}

	return &Storage_BigInt{
		obj:           (*s.obj)[index],
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 6,
		dirty:         s.dirty,
	}
}

// Slice_760032d9d704 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x4, Name:"Slice_760032d9d704", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00038e000)}, Type:"[]*Address", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000c8d00), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_address)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "[]*common.Address"
type Slice_760032d9d704 []*Address

type Storage_Slice_760032d9d704 struct {

	// &wrapper.tmplStruct{T:0x4, Name:"Slice_760032d9d704", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00038e000)}, Type:"[]*Address", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000c8d00), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_address)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           *Slice_760032d9d704
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Slice_760032d9d704) Length() int {
	rv := s.db.GetState(s.addr, common.BigToHash(s.slot))
	return int(rv.Big().Int64())
}

func (s *Storage_Slice_760032d9d704) Resize(length int) {
	// Value: &wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000c8c80), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(big.NewInt(0).SetUint64(uint64(length))))

	_length := length
	_capacity := length + 50

	slice := make([]*Address, _length, _capacity)

	copy(slice, *s.obj)
	*s.obj = slice
}

func (s *Storage_Slice_760032d9d704) Get(index int) *Storage_Address {
	// Value: &wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000c8c80), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	length := s.Length()
	if length <= index {
		s.Resize(index + 1)
	}

	base := crypto.Keccak256Hash(common.BigToHash(s.slot).Bytes()).Big()

	itemsPerSlot := 32 / 20
	actual := big.NewInt(0).Add(base, big.NewInt(0).SetUint64(uint64(index/itemsPerSlot)))
	offset := ((index % itemsPerSlot) * 20)

	if (*s.obj)[index] == nil {

		instance := new(Address)

		(*s.obj)[index] = instance
	}

	return &Storage_Address{
		obj:           (*s.obj)[index],
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 20,
		dirty:         s.dirty,
	}
}

// Slice_91a6287af0aa is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x4, Name:"Slice_91a6287af0aa", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000493a00)}, Type:"[]*Undelegation_t", SolKind:abi.Type{Elem:(*abi.Type)(0xc0001acc00), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_struct(Undelegation_t)78_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "[]*struct { Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=0,offset=0\\\"\"; Epoch *big.Int \"json:\\\"Epoch\\\" storage:\\\"slot=1,offset=0\\\"\" }"
type Slice_91a6287af0aa []*Undelegation_t

type Storage_Slice_91a6287af0aa struct {

	// &wrapper.tmplStruct{T:0x4, Name:"Slice_91a6287af0aa", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000493a00)}, Type:"[]*Undelegation_t", SolKind:abi.Type{Elem:(*abi.Type)(0xc0001acc00), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_struct(Undelegation_t)78_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           *Slice_91a6287af0aa
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Slice_91a6287af0aa) Length() int {
	rv := s.db.GetState(s.addr, common.BigToHash(s.slot))
	return int(rv.Big().Int64())
}

func (s *Storage_Slice_91a6287af0aa) Resize(length int) {
	// Value: &wrapper.tmplField{Type:"Undelegation_t", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc0001aca00), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(big.NewInt(0).SetUint64(uint64(length))))

	_length := length
	_capacity := length + 50

	slice := make([]*Undelegation_t, _length, _capacity)

	copy(slice, *s.obj)
	*s.obj = slice
}

func (s *Storage_Slice_91a6287af0aa) Get(index int) *Storage_Undelegation_t {
	// Value: &wrapper.tmplField{Type:"Undelegation_t", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc0001aca00), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	length := s.Length()
	if length <= index {
		s.Resize(index + 1)
	}

	base := crypto.Keccak256Hash(common.BigToHash(s.slot).Bytes()).Big()

	slotsPerItem := (64 + 31) / 32
	actual := big.NewInt(0).Add(base, big.NewInt(0).SetUint64(uint64(index*slotsPerItem)))
	offset := 0

	if (*s.obj)[index] == nil {

		instance := new(Undelegation_t)

		(*s.obj)[index] = instance
	}

	return &Storage_Undelegation_t{
		obj:           (*s.obj)[index],
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 64,
		dirty:         s.dirty,
	}
}

// Slice_df61efaa55f0 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x4, Name:"Slice_df61efaa55f0", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000493d40)}, Type:"[]*PendingDelegation_t", SolKind:abi.Type{Elem:(*abi.Type)(0xc0001ace00), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_struct(PendingDelegation_t)83_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "[]*struct { Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=0,offset=0\\\"\"; Epoch *big.Int \"json:\\\"Epoch\\\" storage:\\\"slot=1,offset=0\\\"\" }"
type Slice_df61efaa55f0 []*PendingDelegation_t

type Storage_Slice_df61efaa55f0 struct {

	// &wrapper.tmplStruct{T:0x4, Name:"Slice_df61efaa55f0", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000493d40)}, Type:"[]*PendingDelegation_t", SolKind:abi.Type{Elem:(*abi.Type)(0xc0001ace00), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_struct(PendingDelegation_t)83_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           *Slice_df61efaa55f0
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Slice_df61efaa55f0) Length() int {
	rv := s.db.GetState(s.addr, common.BigToHash(s.slot))
	return int(rv.Big().Int64())
}

func (s *Storage_Slice_df61efaa55f0) Resize(length int) {
	// Value: &wrapper.tmplField{Type:"PendingDelegation_t", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc0001acd80), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(big.NewInt(0).SetUint64(uint64(length))))

	_length := length
	_capacity := length + 50

	slice := make([]*PendingDelegation_t, _length, _capacity)

	copy(slice, *s.obj)
	*s.obj = slice
}

func (s *Storage_Slice_df61efaa55f0) Get(index int) *Storage_PendingDelegation_t {
	// Value: &wrapper.tmplField{Type:"PendingDelegation_t", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc0001acd80), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	length := s.Length()
	if length <= index {
		s.Resize(index + 1)
	}

	base := crypto.Keccak256Hash(common.BigToHash(s.slot).Bytes()).Big()

	slotsPerItem := (64 + 31) / 32
	actual := big.NewInt(0).Add(base, big.NewInt(0).SetUint64(uint64(index*slotsPerItem)))
	offset := 0

	if (*s.obj)[index] == nil {

		instance := new(PendingDelegation_t)

		(*s.obj)[index] = instance
	}

	return &Storage_PendingDelegation_t{
		obj:           (*s.obj)[index],
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 64,
		dirty:         s.dirty,
	}
}

// CommissionRates_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"CommissionRates_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012bc70), (*wrapper.tmplField)(0xc00012bd40), (*wrapper.tmplField)(0xc00012be10)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x60, stringKind:"t_struct(CommissionRates_t)12_storage", TupleRawName:"CommissionRates_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048af80), (*abi.Type)(0xc00048b000), (*abi.Type)(0xc00048b180)}, TupleRawNames:[]string{"Rate", "MaxRate", "MaxChangeRate"}, TupleType:(*reflect.rtype)(0xc000124a80)}}
type CommissionRates_t struct {
	Rate          BigInt  `json:"Rate" storage:"slot=0,offset=0"`
	MaxRate       Decimal `json:"MaxRate" storage:"slot=1,offset=0"`
	MaxChangeRate BigInt  `json:"MaxChangeRate" storage:"slot=2,offset=0"`
}

// Commission_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Commission_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012bee0), (*wrapper.tmplField)(0xc000492000)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x80, stringKind:"t_struct(Commission_t)17_storage", TupleRawName:"Commission_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048af00), (*abi.Type)(0xc00048b280)}, TupleRawNames:[]string{"CommissionRates", "UpdateHeight"}, TupleType:(*reflect.rtype)(0xc000124ae0)}}
type Commission_t struct {
	CommissionRates CommissionRates_t `json:"CommissionRates" storage:"slot=0,offset=0"`
	UpdateHeight    BigInt            `json:"UpdateHeight" storage:"slot=3,offset=0"`
}

// Description_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Description_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000492270), (*wrapper.tmplField)(0xc000492340), (*wrapper.tmplField)(0xc000492410), (*wrapper.tmplField)(0xc0004924e0), (*wrapper.tmplField)(0xc0004925b0), (*wrapper.tmplField)(0xc000492680), (*wrapper.tmplField)(0xc000492750), (*wrapper.tmplField)(0xc000492820), (*wrapper.tmplField)(0xc0004929c0), (*wrapper.tmplField)(0xc000492b60), (*wrapper.tmplField)(0xc000492d00), (*wrapper.tmplField)(0xc000492ea0), (*wrapper.tmplField)(0xc000493040), (*wrapper.tmplField)(0xc0004931e0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x300, stringKind:"t_struct(Description_t)58_storage", TupleRawName:"Description_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048b480), (*abi.Type)(0xc00048b580), (*abi.Type)(0xc00048b600), (*abi.Type)(0xc00048b680), (*abi.Type)(0xc00048b700), (*abi.Type)(0xc00048b800), (*abi.Type)(0xc00048b880), (*abi.Type)(0xc00048b980), (*abi.Type)(0xc00048ba00), (*abi.Type)(0xc00048bc00), (*abi.Type)(0xc00048bd80), (*abi.Type)(0xc00048bf00), (*abi.Type)(0xc0001ac180), (*abi.Type)(0xc0001ac380)}, TupleRawNames:[]string{"Name", "Identity", "WebSite", "SecurityContract", "Details", "Frozen", "Serial", "Destroyed", "Flag", "Symbol", "Signature", "Feature", "Mac", "Version"}, TupleType:(*reflect.rtype)(0xc000124f00)}}
type Description_t struct {
	Name             String             `json:"Name" storage:"slot=0,offset=0"`
	Identity         String             `json:"Identity" storage:"slot=1,offset=0"`
	WebSite          String             `json:"WebSite" storage:"slot=2,offset=0"`
	SecurityContract String             `json:"SecurityContract" storage:"slot=3,offset=0"`
	Details          String             `json:"Details" storage:"slot=4,offset=0"`
	Frozen           Uint8              `json:"Frozen" storage:"slot=5,offset=0"`
	Serial           Bytes5             `json:"Serial" storage:"slot=5,offset=1"`
	Destroyed        Bool               `json:"Destroyed" storage:"slot=5,offset=6"`
	Flag             Array_b7edf2cad8cd `json:"Flag" storage:"slot=6,offset=0"`
	Symbol           Array_35cbea3b1b88 `json:"Symbol" storage:"slot=7,offset=0"`
	Signature        Array_9d1a033deeae `json:"Signature" storage:"slot=8,offset=0"`
	Feature          Array_52c9a19dddcb `json:"Feature" storage:"slot=18,offset=0"`
	Mac              Array_61783fec1216 `json:"Mac" storage:"slot=19,offset=0"`
	Version          Array_19387431c506 `json:"Version" storage:"slot=22,offset=0"`
}

// Global_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Global_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0001b1c70), (*wrapper.tmplField)(0xc0001b1d40), (*wrapper.tmplField)(0xc0001b1e10), (*wrapper.tmplField)(0xc0001b1ee0), (*wrapper.tmplField)(0xc00038e0d0), (*wrapper.tmplField)(0xc00038e1a0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x0, stringKind:"t_struct(Global_t)", TupleRawName:"Global_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000c8600), (*abi.Type)(0xc0000c8a80), (*abi.Type)(0xc0000c8b00), (*abi.Type)(0xc0000c8b80), (*abi.Type)(0xc0000c8c00), (*abi.Type)(0xc0000c8e80)}, TupleRawNames:[]string{"pool", "node", "version", "name", "addrs", "signatures"}, TupleType:(*reflect.rtype)(0xc00022c660)}}
type Global_t struct {
	Pool       Map3NodePool_t     `json:"pool" storage:"slot=0,offset=0"`
	Node       Map3Node_t         `json:"node" storage:"slot=6,offset=0"`
	Version    BigInt             `json:"version" storage:"slot=38,offset=0"`
	Name       String             `json:"name" storage:"slot=39,offset=0"`
	Addrs      Slice_760032d9d704 `json:"addrs" storage:"slot=40,offset=0"`
	Signatures Slice_2f14e711180c `json:"signatures" storage:"slot=41,offset=0"`
}

// Map3NodePool_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Map3NodePool_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0001b1450), (*wrapper.tmplField)(0xc0001b1520), (*wrapper.tmplField)(0xc0001b1790), (*wrapper.tmplField)(0xc0001b1860), (*wrapper.tmplField)(0xc0001b1930), (*wrapper.tmplField)(0xc0001b1ba0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xc0, stringKind:"t_struct(Map3NodePool_t)175_storage", TupleRawName:"Map3NodePool_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000c8680), (*abi.Type)(0xc0000c8700), (*abi.Type)(0xc0000c8780), (*abi.Type)(0xc0000c8800), (*abi.Type)(0xc0000c8880), (*abi.Type)(0xc0000c8980)}, TupleRawNames:[]string{"Nodes", "NodeSnapshotByEpoch", "NodeAddressSetByDelegator", "NodeKeySet", "DescriptionIdentitySet", "NodePriority"}, TupleType:(*reflect.rtype)(0xc00022c5a0)}}
type Map3NodePool_t struct {
	Nodes                     Map_db3cea57c803 `json:"Nodes" storage:"slot=0,offset=0"`
	NodeSnapshotByEpoch       Map_cf7837dafed3 `json:"NodeSnapshotByEpoch" storage:"slot=1,offset=0"`
	NodeAddressSetByDelegator Map_dfd3c4f07730 `json:"NodeAddressSetByDelegator" storage:"slot=2,offset=0"`
	NodeKeySet                Map_624bcc0df495 `json:"NodeKeySet" storage:"slot=3,offset=0"`
	DescriptionIdentitySet    Map_624bcc0df495 `json:"DescriptionIdentitySet" storage:"slot=4,offset=0"`
	NodePriority              Map_f75907254fb8 `json:"NodePriority" storage:"slot=5,offset=0"`
}

// Map3NodeSnapshot_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Map3NodeSnapshot_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0001b0ea0), (*wrapper.tmplField)(0xc0001b0f70)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Map3NodeSnapshot_t)148_storage", TupleRawName:"Map3NodeSnapshot_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000c8080), (*abi.Type)(0xc0000c8100)}, TupleRawNames:[]string{"Map3Nodes", "Epoch"}, TupleType:(*reflect.rtype)(0xc000125740)}}
type Map3NodeSnapshot_t struct {
	Map3Nodes Map_db3cea57c803 `json:"Map3Nodes" storage:"slot=0,offset=0"`
	Epoch     BigInt           `json:"Epoch" storage:"slot=1,offset=0"`
}

// Map3NodeWrapper_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Map3NodeWrapper_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000493450), (*wrapper.tmplField)(0xc0001b04e0), (*wrapper.tmplField)(0xc0001b05b0), (*wrapper.tmplField)(0xc0001b0680), (*wrapper.tmplField)(0xc0001b0b60), (*wrapper.tmplField)(0xc0001b0c30), (*wrapper.tmplField)(0xc0001b0d00)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x560, stringKind:"t_struct(Map3NodeWrapper_t)141_storage", TupleRawName:"Map3NodeWrapper_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048ab00), (*abi.Type)(0xc0001ac680), (*abi.Type)(0xc0001ad280), (*abi.Type)(0xc0001ad480), (*abi.Type)(0xc0001ad500), (*abi.Type)(0xc0001ad600), (*abi.Type)(0xc0001ad680)}, TupleRawNames:[]string{"Map3Node", "Microdelegations", "RedelegationReference", "AccumulatedReward", "nodeState", "TotalDelegation", "TotalPendingDelegation"}, TupleType:(*reflect.rtype)(0xc000125380)}}
type Map3NodeWrapper_t struct {
	Map3Node               Map3Node_t              `json:"Map3Node" storage:"slot=0,offset=0"`
	Microdelegations       Map_f6daa253cd33        `json:"Microdelegations" storage:"slot=32,offset=0"`
	RedelegationReference  RedelegationReference_t `json:"RedelegationReference" storage:"slot=33,offset=0"`
	AccumulatedReward      BigInt                  `json:"AccumulatedReward" storage:"slot=35,offset=0"`
	NodeState              NodeState_t             `json:"nodeState" storage:"slot=36,offset=0"`
	TotalDelegation        BigInt                  `json:"TotalDelegation" storage:"slot=41,offset=0"`
	TotalPendingDelegation BigInt                  `json:"TotalPendingDelegation" storage:"slot=42,offset=0"`
}

// Map3Node_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Map3Node_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012b2b0), (*wrapper.tmplField)(0xc00012b450), (*wrapper.tmplField)(0xc00012b520), (*wrapper.tmplField)(0xc00012b860), (*wrapper.tmplField)(0xc0004921a0), (*wrapper.tmplField)(0xc0004932b0), (*wrapper.tmplField)(0xc000493380)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x400, stringKind:"t_struct(Map3Node_t)73_storage", TupleRawName:"Map3Node_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048ab80), (*abi.Type)(0xc00048ac80), (*abi.Type)(0xc00048ad00), (*abi.Type)(0xc00048ad80), (*abi.Type)(0xc00048ae80), (*abi.Type)(0xc00048b400), (*abi.Type)(0xc0001ac580)}, TupleRawNames:[]string{"Actived", "NodeAddress", "InitiatorAddress", "NodeKeys", "Commission", "Description", "SplittedFrom"}, TupleType:(*reflect.rtype)(0xc000124f60)}}
type Map3Node_t struct {
	Actived          Bool          `json:"Actived" storage:"slot=0,offset=0"`
	NodeAddress      Address       `json:"NodeAddress" storage:"slot=0,offset=1"`
	InitiatorAddress Address       `json:"InitiatorAddress" storage:"slot=1,offset=0"`
	NodeKeys         Bytes         `json:"NodeKeys" storage:"slot=2,offset=0"`
	Commission       Commission_t  `json:"Commission" storage:"slot=3,offset=0"`
	Description      Description_t `json:"Description" storage:"slot=7,offset=0"`
	SplittedFrom     Address       `json:"SplittedFrom" storage:"slot=31,offset=0"`
}

// Microdelegation_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Microdelegation_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0004935f0), (*wrapper.tmplField)(0xc0004936c0), (*wrapper.tmplField)(0xc000493790), (*wrapper.tmplField)(0xc000493ad0), (*wrapper.tmplField)(0xc000493e10), (*wrapper.tmplField)(0xc0001b0000), (*wrapper.tmplField)(0xc0001b0270), (*wrapper.tmplField)(0xc0001b0340)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x2c0, stringKind:"t_struct(Microdelegation_t)108_storage", TupleRawName:"Microdelegation_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc0001ac800), (*abi.Type)(0xc0001ac880), (*abi.Type)(0xc0001ac900), (*abi.Type)(0xc0001ac980), (*abi.Type)(0xc0001acd00), (*abi.Type)(0xc0001acf80), (*abi.Type)(0xc0001ad000), (*abi.Type)(0xc0001ad080)}, TupleRawNames:[]string{"DelegatorAddress", "Amount", "Reward", "Undelegations", "PendingDelegations", "PendingDelegationsfixed", "PendingDelegationsfixed2dimension", "AutoRenew"}, TupleType:(*reflect.rtype)(0xc000125200)}}
type Microdelegation_t struct {
	DelegatorAddress                  Address            `json:"DelegatorAddress" storage:"slot=0,offset=0"`
	Amount                            BigInt             `json:"Amount" storage:"slot=1,offset=0"`
	Reward                            BigInt             `json:"Reward" storage:"slot=2,offset=0"`
	Undelegations                     Slice_91a6287af0aa `json:"Undelegations" storage:"slot=3,offset=0"`
	PendingDelegations                Slice_df61efaa55f0 `json:"PendingDelegations" storage:"slot=4,offset=0"`
	PendingDelegationsfixed           Array_bbc2fab83295 `json:"PendingDelegationsfixed" storage:"slot=5,offset=0"`
	PendingDelegationsfixed2dimension Array_0c8aba4f3757 `json:"PendingDelegationsfixed2dimension" storage:"slot=9,offset=0"`
	AutoRenew                         Bool               `json:"AutoRenew" storage:"slot=21,offset=0"`
}

// NodeState_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"NodeState_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0001b0750), (*wrapper.tmplField)(0xc0001b0820), (*wrapper.tmplField)(0xc0001b08f0), (*wrapper.tmplField)(0xc0001b09c0), (*wrapper.tmplField)(0xc0001b0a90)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xa0, stringKind:"t_struct(NodeState_t)124_storage", TupleRawName:"NodeState_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048a000), (*abi.Type)(0xc00048a100), (*abi.Type)(0xc00048a280), (*abi.Type)(0xc00048a300), (*abi.Type)(0xc00048a380)}, TupleRawNames:[]string{"Status", "NodeAge", "CreationEpoch", "ActivationEpoch", "ReleaseEpoch"}, TupleType:(*reflect.rtype)(0xc000124300)}}
type NodeState_t struct {
	Status          Uint8  `json:"Status" storage:"slot=0,offset=0"`
	NodeAge         BigInt `json:"NodeAge" storage:"slot=1,offset=0"`
	CreationEpoch   BigInt `json:"CreationEpoch" storage:"slot=2,offset=0"`
	ActivationEpoch BigInt `json:"ActivationEpoch" storage:"slot=3,offset=0"`
	ReleaseEpoch    BigInt `json:"ReleaseEpoch" storage:"slot=4,offset=0"`
}

// PendingDelegation_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"PendingDelegation_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000493ba0), (*wrapper.tmplField)(0xc000493c70)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(PendingDelegation_t)83_storage", TupleRawName:"PendingDelegation_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048a700), (*abi.Type)(0xc00048a780)}, TupleRawNames:[]string{"Amount", "Epoch"}, TupleType:(*reflect.rtype)(0xc000124600)}}
type PendingDelegation_t struct {
	Amount BigInt `json:"Amount" storage:"slot=0,offset=0"`
	Epoch  BigInt `json:"Epoch" storage:"slot=1,offset=0"`
}

// RedelegationReference_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"RedelegationReference_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012af70), (*wrapper.tmplField)(0xc00012b040)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(RedelegationReference_t)113_storage", TupleRawName:"RedelegationReference_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc0001ad300), (*abi.Type)(0xc0001ad380)}, TupleRawNames:[]string{"ValidatorAddress", "ReleasedTotalDelegation"}, TupleType:(*reflect.rtype)(0xc000125320)}}
type RedelegationReference_t struct {
	ValidatorAddress        Address `json:"ValidatorAddress" storage:"slot=0,offset=0"`
	ReleasedTotalDelegation BigInt  `json:"ReleasedTotalDelegation" storage:"slot=1,offset=0"`
}

// Undelegation_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Undelegation_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000493860), (*wrapper.tmplField)(0xc000493930)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Undelegation_t)78_storage", TupleRawName:"Undelegation_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc0001aca80), (*abi.Type)(0xc0001acb00)}, TupleRawNames:[]string{"Amount", "Epoch"}, TupleType:(*reflect.rtype)(0xc000124600)}}
type Undelegation_t struct {
	Amount BigInt `json:"Amount" storage:"slot=0,offset=0"`
	Epoch  BigInt `json:"Epoch" storage:"slot=1,offset=0"`
}

type Storage_CommissionRates_t struct {

	// &wrapper.tmplStruct{T:0x6, Name:"CommissionRates_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012bc70), (*wrapper.tmplField)(0xc00012bd40), (*wrapper.tmplField)(0xc00012be10)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x60, stringKind:"t_struct(CommissionRates_t)12_storage", TupleRawName:"CommissionRates_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048af80), (*abi.Type)(0xc00048b000), (*abi.Type)(0xc00048b180)}, TupleRawNames:[]string{"Rate", "MaxRate", "MaxChangeRate"}, TupleType:(*reflect.rtype)(0xc000124a80)}}
	obj           *CommissionRates_t
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Commission_t struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Commission_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012bee0), (*wrapper.tmplField)(0xc000492000)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x80, stringKind:"t_struct(Commission_t)17_storage", TupleRawName:"Commission_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048af00), (*abi.Type)(0xc00048b280)}, TupleRawNames:[]string{"CommissionRates", "UpdateHeight"}, TupleType:(*reflect.rtype)(0xc000124ae0)}}
	obj           *Commission_t
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Description_t struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Description_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000492270), (*wrapper.tmplField)(0xc000492340), (*wrapper.tmplField)(0xc000492410), (*wrapper.tmplField)(0xc0004924e0), (*wrapper.tmplField)(0xc0004925b0), (*wrapper.tmplField)(0xc000492680), (*wrapper.tmplField)(0xc000492750), (*wrapper.tmplField)(0xc000492820), (*wrapper.tmplField)(0xc0004929c0), (*wrapper.tmplField)(0xc000492b60), (*wrapper.tmplField)(0xc000492d00), (*wrapper.tmplField)(0xc000492ea0), (*wrapper.tmplField)(0xc000493040), (*wrapper.tmplField)(0xc0004931e0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x300, stringKind:"t_struct(Description_t)58_storage", TupleRawName:"Description_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048b480), (*abi.Type)(0xc00048b580), (*abi.Type)(0xc00048b600), (*abi.Type)(0xc00048b680), (*abi.Type)(0xc00048b700), (*abi.Type)(0xc00048b800), (*abi.Type)(0xc00048b880), (*abi.Type)(0xc00048b980), (*abi.Type)(0xc00048ba00), (*abi.Type)(0xc00048bc00), (*abi.Type)(0xc00048bd80), (*abi.Type)(0xc00048bf00), (*abi.Type)(0xc0001ac180), (*abi.Type)(0xc0001ac380)}, TupleRawNames:[]string{"Name", "Identity", "WebSite", "SecurityContract", "Details", "Frozen", "Serial", "Destroyed", "Flag", "Symbol", "Signature", "Feature", "Mac", "Version"}, TupleType:(*reflect.rtype)(0xc000124f00)}}
	obj           *Description_t
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Global_t struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Global_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0001b1c70), (*wrapper.tmplField)(0xc0001b1d40), (*wrapper.tmplField)(0xc0001b1e10), (*wrapper.tmplField)(0xc0001b1ee0), (*wrapper.tmplField)(0xc00038e0d0), (*wrapper.tmplField)(0xc00038e1a0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x0, stringKind:"t_struct(Global_t)", TupleRawName:"Global_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000c8600), (*abi.Type)(0xc0000c8a80), (*abi.Type)(0xc0000c8b00), (*abi.Type)(0xc0000c8b80), (*abi.Type)(0xc0000c8c00), (*abi.Type)(0xc0000c8e80)}, TupleRawNames:[]string{"pool", "node", "version", "name", "addrs", "signatures"}, TupleType:(*reflect.rtype)(0xc00022c660)}}
	obj           *Global_t
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Map3NodePool_t struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Map3NodePool_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0001b1450), (*wrapper.tmplField)(0xc0001b1520), (*wrapper.tmplField)(0xc0001b1790), (*wrapper.tmplField)(0xc0001b1860), (*wrapper.tmplField)(0xc0001b1930), (*wrapper.tmplField)(0xc0001b1ba0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xc0, stringKind:"t_struct(Map3NodePool_t)175_storage", TupleRawName:"Map3NodePool_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000c8680), (*abi.Type)(0xc0000c8700), (*abi.Type)(0xc0000c8780), (*abi.Type)(0xc0000c8800), (*abi.Type)(0xc0000c8880), (*abi.Type)(0xc0000c8980)}, TupleRawNames:[]string{"Nodes", "NodeSnapshotByEpoch", "NodeAddressSetByDelegator", "NodeKeySet", "DescriptionIdentitySet", "NodePriority"}, TupleType:(*reflect.rtype)(0xc00022c5a0)}}
	obj           *Map3NodePool_t
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Map3NodeSnapshot_t struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Map3NodeSnapshot_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0001b0ea0), (*wrapper.tmplField)(0xc0001b0f70)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Map3NodeSnapshot_t)148_storage", TupleRawName:"Map3NodeSnapshot_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000c8080), (*abi.Type)(0xc0000c8100)}, TupleRawNames:[]string{"Map3Nodes", "Epoch"}, TupleType:(*reflect.rtype)(0xc000125740)}}
	obj           *Map3NodeSnapshot_t
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Map3NodeWrapper_t struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Map3NodeWrapper_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000493450), (*wrapper.tmplField)(0xc0001b04e0), (*wrapper.tmplField)(0xc0001b05b0), (*wrapper.tmplField)(0xc0001b0680), (*wrapper.tmplField)(0xc0001b0b60), (*wrapper.tmplField)(0xc0001b0c30), (*wrapper.tmplField)(0xc0001b0d00)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x560, stringKind:"t_struct(Map3NodeWrapper_t)141_storage", TupleRawName:"Map3NodeWrapper_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048ab00), (*abi.Type)(0xc0001ac680), (*abi.Type)(0xc0001ad280), (*abi.Type)(0xc0001ad480), (*abi.Type)(0xc0001ad500), (*abi.Type)(0xc0001ad600), (*abi.Type)(0xc0001ad680)}, TupleRawNames:[]string{"Map3Node", "Microdelegations", "RedelegationReference", "AccumulatedReward", "nodeState", "TotalDelegation", "TotalPendingDelegation"}, TupleType:(*reflect.rtype)(0xc000125380)}}
	obj           *Map3NodeWrapper_t
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Map3Node_t struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Map3Node_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012b2b0), (*wrapper.tmplField)(0xc00012b450), (*wrapper.tmplField)(0xc00012b520), (*wrapper.tmplField)(0xc00012b860), (*wrapper.tmplField)(0xc0004921a0), (*wrapper.tmplField)(0xc0004932b0), (*wrapper.tmplField)(0xc000493380)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x400, stringKind:"t_struct(Map3Node_t)73_storage", TupleRawName:"Map3Node_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048ab80), (*abi.Type)(0xc00048ac80), (*abi.Type)(0xc00048ad00), (*abi.Type)(0xc00048ad80), (*abi.Type)(0xc00048ae80), (*abi.Type)(0xc00048b400), (*abi.Type)(0xc0001ac580)}, TupleRawNames:[]string{"Actived", "NodeAddress", "InitiatorAddress", "NodeKeys", "Commission", "Description", "SplittedFrom"}, TupleType:(*reflect.rtype)(0xc000124f60)}}
	obj           *Map3Node_t
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Microdelegation_t struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Microdelegation_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0004935f0), (*wrapper.tmplField)(0xc0004936c0), (*wrapper.tmplField)(0xc000493790), (*wrapper.tmplField)(0xc000493ad0), (*wrapper.tmplField)(0xc000493e10), (*wrapper.tmplField)(0xc0001b0000), (*wrapper.tmplField)(0xc0001b0270), (*wrapper.tmplField)(0xc0001b0340)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x2c0, stringKind:"t_struct(Microdelegation_t)108_storage", TupleRawName:"Microdelegation_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc0001ac800), (*abi.Type)(0xc0001ac880), (*abi.Type)(0xc0001ac900), (*abi.Type)(0xc0001ac980), (*abi.Type)(0xc0001acd00), (*abi.Type)(0xc0001acf80), (*abi.Type)(0xc0001ad000), (*abi.Type)(0xc0001ad080)}, TupleRawNames:[]string{"DelegatorAddress", "Amount", "Reward", "Undelegations", "PendingDelegations", "PendingDelegationsfixed", "PendingDelegationsfixed2dimension", "AutoRenew"}, TupleType:(*reflect.rtype)(0xc000125200)}}
	obj           *Microdelegation_t
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_NodeState_t struct {

	// &wrapper.tmplStruct{T:0x6, Name:"NodeState_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0001b0750), (*wrapper.tmplField)(0xc0001b0820), (*wrapper.tmplField)(0xc0001b08f0), (*wrapper.tmplField)(0xc0001b09c0), (*wrapper.tmplField)(0xc0001b0a90)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xa0, stringKind:"t_struct(NodeState_t)124_storage", TupleRawName:"NodeState_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048a000), (*abi.Type)(0xc00048a100), (*abi.Type)(0xc00048a280), (*abi.Type)(0xc00048a300), (*abi.Type)(0xc00048a380)}, TupleRawNames:[]string{"Status", "NodeAge", "CreationEpoch", "ActivationEpoch", "ReleaseEpoch"}, TupleType:(*reflect.rtype)(0xc000124300)}}
	obj           *NodeState_t
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_PendingDelegation_t struct {

	// &wrapper.tmplStruct{T:0x6, Name:"PendingDelegation_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000493ba0), (*wrapper.tmplField)(0xc000493c70)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(PendingDelegation_t)83_storage", TupleRawName:"PendingDelegation_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048a700), (*abi.Type)(0xc00048a780)}, TupleRawNames:[]string{"Amount", "Epoch"}, TupleType:(*reflect.rtype)(0xc000124600)}}
	obj           *PendingDelegation_t
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_RedelegationReference_t struct {

	// &wrapper.tmplStruct{T:0x6, Name:"RedelegationReference_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012af70), (*wrapper.tmplField)(0xc00012b040)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(RedelegationReference_t)113_storage", TupleRawName:"RedelegationReference_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc0001ad300), (*abi.Type)(0xc0001ad380)}, TupleRawNames:[]string{"ValidatorAddress", "ReleasedTotalDelegation"}, TupleType:(*reflect.rtype)(0xc000125320)}}
	obj           *RedelegationReference_t
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Undelegation_t struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Undelegation_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000493860), (*wrapper.tmplField)(0xc000493930)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Undelegation_t)78_storage", TupleRawName:"Undelegation_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc0001aca80), (*abi.Type)(0xc0001acb00)}, TupleRawNames:[]string{"Amount", "Epoch"}, TupleType:(*reflect.rtype)(0xc000124600)}}
	obj           *Undelegation_t
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

// Rate is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_CommissionRates_t) Rate() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Rate", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Rate\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Rate == nil {

		instance := big.NewInt(0)

		s.obj.Rate = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.Rate,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// MaxRate is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_CommissionRates_t) MaxRate() *Storage_Decimal {
	// Field: &wrapper.tmplField{Type:"Decimal", ElemType:"Decimal", Name:"MaxRate", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(Decimal)5_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048b080)}, TupleRawNames:[]string{"f"}, TupleType:(*reflect.rtype)(0xc000124a20)}, Tag:"json:\"MaxRate\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Decimal{
		obj:           &s.obj.MaxRate,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// MaxChangeRate is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_CommissionRates_t) MaxChangeRate() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"MaxChangeRate", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"MaxChangeRate\" storage:\"slot=2,offset=0\"", Slot:2, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("2", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.MaxChangeRate == nil {

		instance := big.NewInt(0)

		s.obj.MaxChangeRate = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.MaxChangeRate,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// CommissionRates is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Commission_t) CommissionRates() *Storage_CommissionRates_t {
	// Field: &wrapper.tmplField{Type:"CommissionRates_t", ElemType:"CommissionRates_t", Name:"CommissionRates", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x60, stringKind:"t_struct(CommissionRates_t)12_storage", TupleRawName:"CommissionRates_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048af80), (*abi.Type)(0xc00048b000), (*abi.Type)(0xc00048b180)}, TupleRawNames:[]string{"Rate", "MaxRate", "MaxChangeRate"}, TupleType:(*reflect.rtype)(0xc000124a80)}, Tag:"json:\"CommissionRates\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_CommissionRates_t{
		obj:           &s.obj.CommissionRates,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 96,
		dirty:         s.dirty,
	}
}

// UpdateHeight is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Commission_t) UpdateHeight() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"UpdateHeight", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"UpdateHeight\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("3", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.UpdateHeight == nil {

		instance := big.NewInt(0)

		s.obj.UpdateHeight = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.UpdateHeight,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Name is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Description_t) Name() *Storage_String {
	// Field: &wrapper.tmplField{Type:"String", ElemType:"String", Name:"Name", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x3, numberOfBytes:0x20, stringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Name\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_String{
		obj:           &s.obj.Name,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Identity is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Description_t) Identity() *Storage_String {
	// Field: &wrapper.tmplField{Type:"String", ElemType:"String", Name:"Identity", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x3, numberOfBytes:0x20, stringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Identity\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_String{
		obj:           &s.obj.Identity,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// WebSite is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Description_t) WebSite() *Storage_String {
	// Field: &wrapper.tmplField{Type:"String", ElemType:"String", Name:"WebSite", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x3, numberOfBytes:0x20, stringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"WebSite\" storage:\"slot=2,offset=0\"", Slot:2, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("2", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_String{
		obj:           &s.obj.WebSite,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// SecurityContract is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Description_t) SecurityContract() *Storage_String {
	// Field: &wrapper.tmplField{Type:"String", ElemType:"String", Name:"SecurityContract", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x3, numberOfBytes:0x20, stringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"SecurityContract\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("3", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_String{
		obj:           &s.obj.SecurityContract,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Details is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Description_t) Details() *Storage_String {
	// Field: &wrapper.tmplField{Type:"String", ElemType:"String", Name:"Details", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x3, numberOfBytes:0x20, stringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Details\" storage:\"slot=4,offset=0\"", Slot:4, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("4", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_String{
		obj:           &s.obj.Details,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Frozen is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Description_t) Frozen() *Storage_Uint8 {
	// Field: &wrapper.tmplField{Type:"Uint8", ElemType:"Uint8", Name:"Frozen", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:8, T:0x1, numberOfBytes:0x1, stringKind:"uint8", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Frozen\" storage:\"slot=5,offset=0\"", Slot:5, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("5", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Uint8{
		obj:           &s.obj.Frozen,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 1,
		dirty:         s.dirty,
	}
}

// Serial is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Description_t) Serial() *Storage_Bytes5 {
	// Field: &wrapper.tmplField{Type:"Bytes5", ElemType:"Bytes5", Name:"Serial", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:5, T:0x8, numberOfBytes:0x5, stringKind:"t_bytes5", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Serial\" storage:\"slot=5,offset=1\"", Slot:5, Offset:1}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("5", 10)
	offset := 1

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Bytes5{
		obj:           &s.obj.Serial,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 5,
		dirty:         s.dirty,
	}
}

// Destroyed is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Description_t) Destroyed() *Storage_Bool {
	// Field: &wrapper.tmplField{Type:"Bool", ElemType:"Bool", Name:"Destroyed", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x2, numberOfBytes:0x1, stringKind:"t_bool", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Destroyed\" storage:\"slot=5,offset=6\"", Slot:5, Offset:6}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("5", 10)
	offset := 6

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Bool{
		obj:           &s.obj.Destroyed,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 1,
		dirty:         s.dirty,
	}
}

// Flag is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Description_t) Flag() *Storage_Array_b7edf2cad8cd {
	// Field: &wrapper.tmplField{Type:"Array_b7edf2cad8cd", ElemType:"Array_b7edf2cad8cd", Name:"Flag", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048ba80), Size:1, T:0x5, numberOfBytes:0x20, stringKind:"t_array(t_bytes1)1_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Flag\" storage:\"slot=6,offset=0\"", Slot:6, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("6", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Array_b7edf2cad8cd{
		obj:           &s.obj.Flag,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Symbol is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Description_t) Symbol() *Storage_Array_35cbea3b1b88 {
	// Field: &wrapper.tmplField{Type:"Array_35cbea3b1b88", ElemType:"Array_35cbea3b1b88", Name:"Symbol", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048bc80), Size:9, T:0x5, numberOfBytes:0x20, stringKind:"t_array(t_bytes1)9_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Symbol\" storage:\"slot=7,offset=0\"", Slot:7, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("7", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Array_35cbea3b1b88{
		obj:           &s.obj.Symbol,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Signature is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Description_t) Signature() *Storage_Array_9d1a033deeae {
	// Field: &wrapper.tmplField{Type:"Array_9d1a033deeae", ElemType:"Array_9d1a033deeae", Name:"Signature", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048be00), Size:300, T:0x5, numberOfBytes:0x140, stringKind:"t_array(t_bytes1)300_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Signature\" storage:\"slot=8,offset=0\"", Slot:8, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("8", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Array_9d1a033deeae{
		obj:           &s.obj.Signature,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 320,
		dirty:         s.dirty,
	}
}

// Feature is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Description_t) Feature() *Storage_Array_52c9a19dddcb {
	// Field: &wrapper.tmplField{Type:"Array_52c9a19dddcb", ElemType:"Array_52c9a19dddcb", Name:"Feature", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048bf80), Size:4, T:0x5, numberOfBytes:0x20, stringKind:"t_array(t_uint32)4_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Feature\" storage:\"slot=18,offset=0\"", Slot:18, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("18", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Array_52c9a19dddcb{
		obj:           &s.obj.Feature,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Mac is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Description_t) Mac() *Storage_Array_61783fec1216 {
	// Field: &wrapper.tmplField{Type:"Array_61783fec1216", ElemType:"Array_61783fec1216", Name:"Mac", SolKind:abi.Type{Elem:(*abi.Type)(0xc0001ac200), Size:12, T:0x5, numberOfBytes:0x60, stringKind:"t_array(t_uint48)12_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Mac\" storage:\"slot=19,offset=0\"", Slot:19, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("19", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Array_61783fec1216{
		obj:           &s.obj.Mac,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 96,
		dirty:         s.dirty,
	}
}

// Version is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Description_t) Version() *Storage_Array_19387431c506 {
	// Field: &wrapper.tmplField{Type:"Array_19387431c506", ElemType:"Array_19387431c506", Name:"Version", SolKind:abi.Type{Elem:(*abi.Type)(0xc0001ac400), Size:2, T:0x5, numberOfBytes:0x40, stringKind:"t_array(t_uint256)2_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Version\" storage:\"slot=22,offset=0\"", Slot:22, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("22", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Array_19387431c506{
		obj:           &s.obj.Version,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 64,
		dirty:         s.dirty,
	}
}

// Pool is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Global_t) Pool() *Storage_Map3NodePool_t {
	// Field: &wrapper.tmplField{Type:"Map3NodePool_t", ElemType:"Map3NodePool_t", Name:"Pool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xc0, stringKind:"t_struct(Map3NodePool_t)175_storage", TupleRawName:"Map3NodePool_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000c8680), (*abi.Type)(0xc0000c8700), (*abi.Type)(0xc0000c8780), (*abi.Type)(0xc0000c8800), (*abi.Type)(0xc0000c8880), (*abi.Type)(0xc0000c8980)}, TupleRawNames:[]string{"Nodes", "NodeSnapshotByEpoch", "NodeAddressSetByDelegator", "NodeKeySet", "DescriptionIdentitySet", "NodePriority"}, TupleType:(*reflect.rtype)(0xc00022c5a0)}, Tag:"json:\"pool\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Map3NodePool_t{
		obj:           &s.obj.Pool,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 192,
		dirty:         s.dirty,
	}
}

// Node is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Global_t) Node() *Storage_Map3Node_t {
	// Field: &wrapper.tmplField{Type:"Map3Node_t", ElemType:"Map3Node_t", Name:"Node", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x400, stringKind:"t_struct(Map3Node_t)73_storage", TupleRawName:"Map3Node_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048ab80), (*abi.Type)(0xc00048ac80), (*abi.Type)(0xc00048ad00), (*abi.Type)(0xc00048ad80), (*abi.Type)(0xc00048ae80), (*abi.Type)(0xc00048b400), (*abi.Type)(0xc0001ac580)}, TupleRawNames:[]string{"Actived", "NodeAddress", "InitiatorAddress", "NodeKeys", "Commission", "Description", "SplittedFrom"}, TupleType:(*reflect.rtype)(0xc000124f60)}, Tag:"json:\"node\" storage:\"slot=6,offset=0\"", Slot:6, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("6", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Map3Node_t{
		obj:           &s.obj.Node,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 1024,
		dirty:         s.dirty,
	}
}

// Version is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Global_t) Version() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Version", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x0, numberOfBytes:0x20, stringKind:"t_int256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"version\" storage:\"slot=38,offset=0\"", Slot:38, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("38", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Version == nil {

		instance := big.NewInt(0)

		s.obj.Version = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.Version,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Name is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Global_t) Name() *Storage_String {
	// Field: &wrapper.tmplField{Type:"String", ElemType:"String", Name:"Name", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x3, numberOfBytes:0x20, stringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"name\" storage:\"slot=39,offset=0\"", Slot:39, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("39", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_String{
		obj:           &s.obj.Name,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Addrs is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Global_t) Addrs() *Storage_Slice_760032d9d704 {
	// Field: &wrapper.tmplField{Type:"Slice_760032d9d704", ElemType:"Slice_760032d9d704", Name:"Addrs", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000c8d00), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_address)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"addrs\" storage:\"slot=40,offset=0\"", Slot:40, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("40", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Addrs == nil {

		hash := s.db.GetState(s.addr, common.BigToHash(actual))
		instance := make(Slice_760032d9d704, hash.Big().Int64())

		s.obj.Addrs = instance
	}

	return &Storage_Slice_760032d9d704{
		obj:           &s.obj.Addrs,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Signatures is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Global_t) Signatures() *Storage_Slice_2f14e711180c {
	// Field: &wrapper.tmplField{Type:"Slice_2f14e711180c", ElemType:"Slice_2f14e711180c", Name:"Signatures", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000c8280), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_uint48)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"signatures\" storage:\"slot=41,offset=0\"", Slot:41, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("41", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Signatures == nil {

		hash := s.db.GetState(s.addr, common.BigToHash(actual))
		instance := make(Slice_2f14e711180c, hash.Big().Int64())

		s.obj.Signatures = instance
	}

	return &Storage_Slice_2f14e711180c{
		obj:           &s.obj.Signatures,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Nodes is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodePool_t) Nodes() *Storage_Map_db3cea57c803 {
	// Field: &wrapper.tmplField{Type:"Map_db3cea57c803", ElemType:"Map_db3cea57c803", Name:"Nodes", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(Map3NodeWrapper_t)141_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048a980), (*abi.Type)(0xc0001ad780)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Nodes\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Nodes == nil {

		instance := make(Map_db3cea57c803)

		s.obj.Nodes = instance
	}

	return &Storage_Map_db3cea57c803{
		obj:           s.obj.Nodes,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// NodeSnapshotByEpoch is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodePool_t) NodeSnapshotByEpoch() *Storage_Map_cf7837dafed3 {
	// Field: &wrapper.tmplField{Type:"Map_cf7837dafed3", ElemType:"Map_cf7837dafed3", Name:"NodeSnapshotByEpoch", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_uint64,t_struct(Map3NodeSnapshot_t)148_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000c8400), (*abi.Type)(0xc0000c8500)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"NodeSnapshotByEpoch\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.NodeSnapshotByEpoch == nil {

		instance := make(Map_cf7837dafed3)

		s.obj.NodeSnapshotByEpoch = instance
	}

	return &Storage_Map_cf7837dafed3{
		obj:           s.obj.NodeSnapshotByEpoch,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// NodeAddressSetByDelegator is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodePool_t) NodeAddressSetByDelegator() *Storage_Map_dfd3c4f07730 {
	// Field: &wrapper.tmplField{Type:"Map_dfd3c4f07730", ElemType:"Map_dfd3c4f07730", Name:"NodeAddressSetByDelegator", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_mapping(t_address,t_bool))", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0001adb00), (*abi.Type)(0xc0001adb80)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"NodeAddressSetByDelegator\" storage:\"slot=2,offset=0\"", Slot:2, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("2", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.NodeAddressSetByDelegator == nil {

		instance := make(Map_dfd3c4f07730)

		s.obj.NodeAddressSetByDelegator = instance
	}

	return &Storage_Map_dfd3c4f07730{
		obj:           s.obj.NodeAddressSetByDelegator,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// NodeKeySet is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodePool_t) NodeKeySet() *Storage_Map_624bcc0df495 {
	// Field: &wrapper.tmplField{Type:"Map_624bcc0df495", ElemType:"Map_624bcc0df495", Name:"NodeKeySet", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_string_memory_ptr,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0001ad880), (*abi.Type)(0xc0001ada00)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"NodeKeySet\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("3", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.NodeKeySet == nil {

		instance := make(Map_624bcc0df495)

		s.obj.NodeKeySet = instance
	}

	return &Storage_Map_624bcc0df495{
		obj:           s.obj.NodeKeySet,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// DescriptionIdentitySet is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodePool_t) DescriptionIdentitySet() *Storage_Map_624bcc0df495 {
	// Field: &wrapper.tmplField{Type:"Map_624bcc0df495", ElemType:"Map_624bcc0df495", Name:"DescriptionIdentitySet", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_string_memory_ptr,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0001ad880), (*abi.Type)(0xc0001ada00)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"DescriptionIdentitySet\" storage:\"slot=4,offset=0\"", Slot:4, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("4", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.DescriptionIdentitySet == nil {

		instance := make(Map_624bcc0df495)

		s.obj.DescriptionIdentitySet = instance
	}

	return &Storage_Map_624bcc0df495{
		obj:           s.obj.DescriptionIdentitySet,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// NodePriority is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodePool_t) NodePriority() *Storage_Map_f75907254fb8 {
	// Field: &wrapper.tmplField{Type:"Map_f75907254fb8", ElemType:"Map_f75907254fb8", Name:"NodePriority", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_string_memory_ptr,t_uint256)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0001ade80), (*abi.Type)(0xc0001adf00)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"NodePriority\" storage:\"slot=5,offset=0\"", Slot:5, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("5", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.NodePriority == nil {

		instance := make(Map_f75907254fb8)

		s.obj.NodePriority = instance
	}

	return &Storage_Map_f75907254fb8{
		obj:           s.obj.NodePriority,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Map3Nodes is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeSnapshot_t) Map3Nodes() *Storage_Map_db3cea57c803 {
	// Field: &wrapper.tmplField{Type:"Map_db3cea57c803", ElemType:"Map_db3cea57c803", Name:"Map3Nodes", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(Map3NodeWrapper_t)141_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048a980), (*abi.Type)(0xc0001ad780)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Map3Nodes\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Map3Nodes == nil {

		instance := make(Map_db3cea57c803)

		s.obj.Map3Nodes = instance
	}

	return &Storage_Map_db3cea57c803{
		obj:           s.obj.Map3Nodes,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Epoch is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeSnapshot_t) Epoch() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Epoch", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Epoch\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Epoch == nil {

		instance := big.NewInt(0)

		s.obj.Epoch = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.Epoch,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Map3Node is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeWrapper_t) Map3Node() *Storage_Map3Node_t {
	// Field: &wrapper.tmplField{Type:"Map3Node_t", ElemType:"Map3Node_t", Name:"Map3Node", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x400, stringKind:"t_struct(Map3Node_t)73_storage", TupleRawName:"Map3Node_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048ab80), (*abi.Type)(0xc00048ac80), (*abi.Type)(0xc00048ad00), (*abi.Type)(0xc00048ad80), (*abi.Type)(0xc00048ae80), (*abi.Type)(0xc00048b400), (*abi.Type)(0xc0001ac580)}, TupleRawNames:[]string{"Actived", "NodeAddress", "InitiatorAddress", "NodeKeys", "Commission", "Description", "SplittedFrom"}, TupleType:(*reflect.rtype)(0xc000124f60)}, Tag:"json:\"Map3Node\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Map3Node_t{
		obj:           &s.obj.Map3Node,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 1024,
		dirty:         s.dirty,
	}
}

// Microdelegations is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeWrapper_t) Microdelegations() *Storage_Map_f6daa253cd33 {
	// Field: &wrapper.tmplField{Type:"Map_f6daa253cd33", ElemType:"Map_f6daa253cd33", Name:"Microdelegations", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(Microdelegation_t)108_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0001ac700), (*abi.Type)(0xc0001ad180)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Microdelegations\" storage:\"slot=32,offset=0\"", Slot:32, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("32", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Microdelegations == nil {

		instance := make(Map_f6daa253cd33)

		s.obj.Microdelegations = instance
	}

	return &Storage_Map_f6daa253cd33{
		obj:           s.obj.Microdelegations,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// RedelegationReference is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeWrapper_t) RedelegationReference() *Storage_RedelegationReference_t {
	// Field: &wrapper.tmplField{Type:"RedelegationReference_t", ElemType:"RedelegationReference_t", Name:"RedelegationReference", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(RedelegationReference_t)113_storage", TupleRawName:"RedelegationReference_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc0001ad300), (*abi.Type)(0xc0001ad380)}, TupleRawNames:[]string{"ValidatorAddress", "ReleasedTotalDelegation"}, TupleType:(*reflect.rtype)(0xc000125320)}, Tag:"json:\"RedelegationReference\" storage:\"slot=33,offset=0\"", Slot:33, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("33", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_RedelegationReference_t{
		obj:           &s.obj.RedelegationReference,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 64,
		dirty:         s.dirty,
	}
}

// AccumulatedReward is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeWrapper_t) AccumulatedReward() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"AccumulatedReward", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"AccumulatedReward\" storage:\"slot=35,offset=0\"", Slot:35, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("35", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.AccumulatedReward == nil {

		instance := big.NewInt(0)

		s.obj.AccumulatedReward = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.AccumulatedReward,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// NodeState is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeWrapper_t) NodeState() *Storage_NodeState_t {
	// Field: &wrapper.tmplField{Type:"NodeState_t", ElemType:"NodeState_t", Name:"NodeState", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xa0, stringKind:"t_struct(NodeState_t)124_storage", TupleRawName:"NodeState_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048a000), (*abi.Type)(0xc00048a100), (*abi.Type)(0xc00048a280), (*abi.Type)(0xc00048a300), (*abi.Type)(0xc00048a380)}, TupleRawNames:[]string{"Status", "NodeAge", "CreationEpoch", "ActivationEpoch", "ReleaseEpoch"}, TupleType:(*reflect.rtype)(0xc000124300)}, Tag:"json:\"nodeState\" storage:\"slot=36,offset=0\"", Slot:36, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("36", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_NodeState_t{
		obj:           &s.obj.NodeState,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 160,
		dirty:         s.dirty,
	}
}

// TotalDelegation is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeWrapper_t) TotalDelegation() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"TotalDelegation", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"TotalDelegation\" storage:\"slot=41,offset=0\"", Slot:41, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("41", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.TotalDelegation == nil {

		instance := big.NewInt(0)

		s.obj.TotalDelegation = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.TotalDelegation,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// TotalPendingDelegation is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeWrapper_t) TotalPendingDelegation() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"TotalPendingDelegation", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"TotalPendingDelegation\" storage:\"slot=42,offset=0\"", Slot:42, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("42", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.TotalPendingDelegation == nil {

		instance := big.NewInt(0)

		s.obj.TotalPendingDelegation = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.TotalPendingDelegation,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Actived is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3Node_t) Actived() *Storage_Bool {
	// Field: &wrapper.tmplField{Type:"Bool", ElemType:"Bool", Name:"Actived", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x2, numberOfBytes:0x1, stringKind:"t_bool", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Actived\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Bool{
		obj:           &s.obj.Actived,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 1,
		dirty:         s.dirty,
	}
}

// NodeAddress is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3Node_t) NodeAddress() *Storage_Address {
	// Field: &wrapper.tmplField{Type:"Address", ElemType:"Address", Name:"NodeAddress", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"NodeAddress\" storage:\"slot=0,offset=1\"", Slot:0, Offset:1}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 1

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Address{
		obj:           &s.obj.NodeAddress,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 20,
		dirty:         s.dirty,
	}
}

// InitiatorAddress is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3Node_t) InitiatorAddress() *Storage_Address {
	// Field: &wrapper.tmplField{Type:"Address", ElemType:"Address", Name:"InitiatorAddress", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"InitiatorAddress\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Address{
		obj:           &s.obj.InitiatorAddress,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 20,
		dirty:         s.dirty,
	}
}

// NodeKeys is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3Node_t) NodeKeys() *Storage_Bytes {
	// Field: &wrapper.tmplField{Type:"Bytes", ElemType:"Bytes", Name:"NodeKeys", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x9, numberOfBytes:0x20, stringKind:"t_bytes_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"NodeKeys\" storage:\"slot=2,offset=0\"", Slot:2, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("2", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.NodeKeys == nil {

		hash := s.db.GetState(s.addr, common.BigToHash(actual))
		instance := make(Bytes, hash.Big().Int64())

		s.obj.NodeKeys = instance
	}

	return &Storage_Bytes{
		obj:           &s.obj.NodeKeys,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Commission is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3Node_t) Commission() *Storage_Commission_t {
	// Field: &wrapper.tmplField{Type:"Commission_t", ElemType:"Commission_t", Name:"Commission", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x80, stringKind:"t_struct(Commission_t)17_storage", TupleRawName:"Commission_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048af00), (*abi.Type)(0xc00048b280)}, TupleRawNames:[]string{"CommissionRates", "UpdateHeight"}, TupleType:(*reflect.rtype)(0xc000124ae0)}, Tag:"json:\"Commission\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("3", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Commission_t{
		obj:           &s.obj.Commission,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 128,
		dirty:         s.dirty,
	}
}

// Description is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3Node_t) Description() *Storage_Description_t {
	// Field: &wrapper.tmplField{Type:"Description_t", ElemType:"Description_t", Name:"Description", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x300, stringKind:"t_struct(Description_t)58_storage", TupleRawName:"Description_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048b480), (*abi.Type)(0xc00048b580), (*abi.Type)(0xc00048b600), (*abi.Type)(0xc00048b680), (*abi.Type)(0xc00048b700), (*abi.Type)(0xc00048b800), (*abi.Type)(0xc00048b880), (*abi.Type)(0xc00048b980), (*abi.Type)(0xc00048ba00), (*abi.Type)(0xc00048bc00), (*abi.Type)(0xc00048bd80), (*abi.Type)(0xc00048bf00), (*abi.Type)(0xc0001ac180), (*abi.Type)(0xc0001ac380)}, TupleRawNames:[]string{"Name", "Identity", "WebSite", "SecurityContract", "Details", "Frozen", "Serial", "Destroyed", "Flag", "Symbol", "Signature", "Feature", "Mac", "Version"}, TupleType:(*reflect.rtype)(0xc000124f00)}, Tag:"json:\"Description\" storage:\"slot=7,offset=0\"", Slot:7, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("7", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Description_t{
		obj:           &s.obj.Description,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 768,
		dirty:         s.dirty,
	}
}

// SplittedFrom is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3Node_t) SplittedFrom() *Storage_Address {
	// Field: &wrapper.tmplField{Type:"Address", ElemType:"Address", Name:"SplittedFrom", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"SplittedFrom\" storage:\"slot=31,offset=0\"", Slot:31, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("31", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Address{
		obj:           &s.obj.SplittedFrom,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 20,
		dirty:         s.dirty,
	}
}

// DelegatorAddress is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Microdelegation_t) DelegatorAddress() *Storage_Address {
	// Field: &wrapper.tmplField{Type:"Address", ElemType:"Address", Name:"DelegatorAddress", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"DelegatorAddress\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Address{
		obj:           &s.obj.DelegatorAddress,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 20,
		dirty:         s.dirty,
	}
}

// Amount is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Microdelegation_t) Amount() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Amount", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Amount\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Amount == nil {

		instance := big.NewInt(0)

		s.obj.Amount = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.Amount,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Reward is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Microdelegation_t) Reward() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Reward", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Reward\" storage:\"slot=2,offset=0\"", Slot:2, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("2", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Reward == nil {

		instance := big.NewInt(0)

		s.obj.Reward = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.Reward,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Undelegations is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Microdelegation_t) Undelegations() *Storage_Slice_91a6287af0aa {
	// Field: &wrapper.tmplField{Type:"Slice_91a6287af0aa", ElemType:"Slice_91a6287af0aa", Name:"Undelegations", SolKind:abi.Type{Elem:(*abi.Type)(0xc0001acc00), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_struct(Undelegation_t)78_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Undelegations\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("3", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Undelegations == nil {

		hash := s.db.GetState(s.addr, common.BigToHash(actual))
		instance := make(Slice_91a6287af0aa, hash.Big().Int64())

		s.obj.Undelegations = instance
	}

	return &Storage_Slice_91a6287af0aa{
		obj:           &s.obj.Undelegations,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// PendingDelegations is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Microdelegation_t) PendingDelegations() *Storage_Slice_df61efaa55f0 {
	// Field: &wrapper.tmplField{Type:"Slice_df61efaa55f0", ElemType:"Slice_df61efaa55f0", Name:"PendingDelegations", SolKind:abi.Type{Elem:(*abi.Type)(0xc0001ace00), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_struct(PendingDelegation_t)83_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"PendingDelegations\" storage:\"slot=4,offset=0\"", Slot:4, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("4", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.PendingDelegations == nil {

		hash := s.db.GetState(s.addr, common.BigToHash(actual))
		instance := make(Slice_df61efaa55f0, hash.Big().Int64())

		s.obj.PendingDelegations = instance
	}

	return &Storage_Slice_df61efaa55f0{
		obj:           &s.obj.PendingDelegations,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// PendingDelegationsfixed is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Microdelegation_t) PendingDelegationsfixed() *Storage_Array_bbc2fab83295 {
	// Field: &wrapper.tmplField{Type:"Array_bbc2fab83295", ElemType:"Array_bbc2fab83295", Name:"PendingDelegationsfixed", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048a680), Size:2, T:0x5, numberOfBytes:0x80, stringKind:"t_array(t_struct(PendingDelegation_t)83_storage)2_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"PendingDelegationsfixed\" storage:\"slot=5,offset=0\"", Slot:5, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("5", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Array_bbc2fab83295{
		obj:           &s.obj.PendingDelegationsfixed,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 128,
		dirty:         s.dirty,
	}
}

// PendingDelegationsfixed2dimension is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Microdelegation_t) PendingDelegationsfixed2dimension() *Storage_Array_0c8aba4f3757 {
	// Field: &wrapper.tmplField{Type:"Array_0c8aba4f3757", ElemType:"Array_0c8aba4f3757", Name:"PendingDelegationsfixed2dimension", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048a600), Size:3, T:0x5, numberOfBytes:0x180, stringKind:"t_array(t_array(t_struct(PendingDelegation_t)83_storage)2_storage)3_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"PendingDelegationsfixed2dimension\" storage:\"slot=9,offset=0\"", Slot:9, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("9", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Array_0c8aba4f3757{
		obj:           &s.obj.PendingDelegationsfixed2dimension,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 384,
		dirty:         s.dirty,
	}
}

// AutoRenew is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Microdelegation_t) AutoRenew() *Storage_Bool {
	// Field: &wrapper.tmplField{Type:"Bool", ElemType:"Bool", Name:"AutoRenew", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x2, numberOfBytes:0x1, stringKind:"t_bool", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"AutoRenew\" storage:\"slot=21,offset=0\"", Slot:21, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("21", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Bool{
		obj:           &s.obj.AutoRenew,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 1,
		dirty:         s.dirty,
	}
}

// Status is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_NodeState_t) Status() *Storage_Uint8 {
	// Field: &wrapper.tmplField{Type:"Uint8", ElemType:"Uint8", Name:"Status", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:8, T:0x1, numberOfBytes:0x1, stringKind:"uint8", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Status\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Uint8{
		obj:           &s.obj.Status,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 1,
		dirty:         s.dirty,
	}
}

// NodeAge is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_NodeState_t) NodeAge() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"NodeAge", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"NodeAge\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.NodeAge == nil {

		instance := big.NewInt(0)

		s.obj.NodeAge = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.NodeAge,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// CreationEpoch is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_NodeState_t) CreationEpoch() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"CreationEpoch", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"CreationEpoch\" storage:\"slot=2,offset=0\"", Slot:2, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("2", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.CreationEpoch == nil {

		instance := big.NewInt(0)

		s.obj.CreationEpoch = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.CreationEpoch,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// ActivationEpoch is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_NodeState_t) ActivationEpoch() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"ActivationEpoch", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"ActivationEpoch\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("3", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.ActivationEpoch == nil {

		instance := big.NewInt(0)

		s.obj.ActivationEpoch = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.ActivationEpoch,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// ReleaseEpoch is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_NodeState_t) ReleaseEpoch() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"ReleaseEpoch", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"ReleaseEpoch\" storage:\"slot=4,offset=0\"", Slot:4, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("4", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.ReleaseEpoch == nil {

		instance := big.NewInt(0)

		s.obj.ReleaseEpoch = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.ReleaseEpoch,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Amount is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_PendingDelegation_t) Amount() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Amount", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Amount\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Amount == nil {

		instance := big.NewInt(0)

		s.obj.Amount = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.Amount,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Epoch is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_PendingDelegation_t) Epoch() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Epoch", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Epoch\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Epoch == nil {

		instance := big.NewInt(0)

		s.obj.Epoch = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.Epoch,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// ValidatorAddress is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_RedelegationReference_t) ValidatorAddress() *Storage_Address {
	// Field: &wrapper.tmplField{Type:"Address", ElemType:"Address", Name:"ValidatorAddress", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"ValidatorAddress\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Address{
		obj:           &s.obj.ValidatorAddress,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 20,
		dirty:         s.dirty,
	}
}

// ReleasedTotalDelegation is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_RedelegationReference_t) ReleasedTotalDelegation() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"ReleasedTotalDelegation", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"ReleasedTotalDelegation\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.ReleasedTotalDelegation == nil {

		instance := big.NewInt(0)

		s.obj.ReleasedTotalDelegation = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.ReleasedTotalDelegation,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Amount is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Undelegation_t) Amount() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Amount", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Amount\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Amount == nil {

		instance := big.NewInt(0)

		s.obj.Amount = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.Amount,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Epoch is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Undelegation_t) Epoch() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Epoch", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Epoch\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Epoch == nil {

		instance := big.NewInt(0)

		s.obj.Epoch = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.Epoch,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}
