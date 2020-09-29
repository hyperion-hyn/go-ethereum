// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package microstaking

import (
	"bytes"
	"encoding/json"
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

type Decimal = common.Dec
type Storage_Decimal struct {
	// Builtin-Type

	// &wrapper.tmplStruct{T:0xf, Name:"Decimal", Fields:[]*wrapper.tmplField(nil), Type:"common.Dec", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xf, numberOfBytes:0x20, stringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc000127580)}, TupleRawNames:[]string{"f"}, TupleType:(*reflect.rtype)(0xc00054a660)}}
	obj           *Decimal
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Decimal) Value() common.Dec {
	// &wrapper.tmplStruct{T:0xf, Name:"Decimal", Fields:[]*wrapper.tmplField(nil), Type:"common.Dec", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xf, numberOfBytes:0x20, stringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc000127580)}, TupleRawNames:[]string{"f"}, TupleType:(*reflect.rtype)(0xc00054a660)}}
	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	*s.obj = common.NewDecFromBigIntWithPrec(hash.Big(), common.Precision)
	return *s.obj
}

func (s *Storage_Decimal) SetValue(value common.Dec) {
	// &wrapper.tmplStruct{T:0xf, Name:"Decimal", Fields:[]*wrapper.tmplField(nil), Type:"common.Dec", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xf, numberOfBytes:0x20, stringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc000127580)}, TupleRawNames:[]string{"f"}, TupleType:(*reflect.rtype)(0xc00054a660)}}
	hash := value.BigInt()
	if hash != nil {
		s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(hash))
	}
	*s.obj = value
}

type String = string
type Storage_String struct {
	// Builtin-Type

	// &wrapper.tmplStruct{T:0x3, Name:"String", Fields:[]*wrapper.tmplField(nil), Type:"string", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x3, numberOfBytes:0x20, stringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           *String
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_String) Value() string {
	// &wrapper.tmplStruct{T:0x3, Name:"String", Fields:[]*wrapper.tmplField(nil), Type:"string", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x3, numberOfBytes:0x20, stringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	rv := GetStateAsBytes(s.db, s.addr, s.slot)
	*s.obj = string(rv)
	return *s.obj
}

func (s *Storage_String) SetValue(value string) {
	// &wrapper.tmplStruct{T:0x3, Name:"String", Fields:[]*wrapper.tmplField(nil), Type:"string", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x3, numberOfBytes:0x20, stringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	SetStateAsBytes(s.db, s.addr, s.slot, []byte(value))
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

// Array_aaa956d410fd is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x5, Name:"Array_aaa956d410fd", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054e820)}, Type:"[48]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc000127100), Size:48, T:0x5, numberOfBytes:0x40, stringKind:"t_array(t_bytes1)48_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "[48]uint8"
type Array_aaa956d410fd [48]Uint8

type Storage_Array_aaa956d410fd struct {

	// &wrapper.tmplStruct{T:0x5, Name:"Array_aaa956d410fd", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054e820)}, Type:"[48]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc000127100), Size:48, T:0x5, numberOfBytes:0x40, stringKind:"t_array(t_bytes1)48_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           *Array_aaa956d410fd
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Array_aaa956d410fd) Value() [48]Uint8 {
	// &wrapper.tmplStruct{T:0x5, Name:"Array_aaa956d410fd", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054e820)}, Type:"[48]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc000127100), Size:48, T:0x5, numberOfBytes:0x40, stringKind:"t_array(t_bytes1)48_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
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

func (s *Storage_Array_aaa956d410fd) SetValue(value [48]Uint8) {
	// &wrapper.tmplStruct{T:0x5, Name:"Array_aaa956d410fd", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054e820)}, Type:"[48]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc000127100), Size:48, T:0x5, numberOfBytes:0x40, stringKind:"t_array(t_bytes1)48_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
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

// Map_2f2f5e901078 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xe, Name:"Map_2f2f5e901078", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054f450), (*wrapper.tmplField)(0xc000422000)}, Type:"map[Address]*MicrodelegationMapEntry_", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(MicrodelegationMapEntry_)87_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc000550380), (*abi.Type)(0xc0000d2180)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "map[common.Address]*struct { Entry struct { DelegatorAddress common.Address \"json:\\\"DelegatorAddress\\\" storage:\\\"slot=0,offset=0\\\"\"; Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=1,offset=0\\\"\"; Reward *big.Int \"json:\\\"Reward\\\" storage:\\\"slot=2,offset=0\\\"\"; PendingDelegation struct { Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=0,offset=0\\\"\"; UnlockedEpoch struct { F *big.Int \"json:\\\"f\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"UnlockedEpoch\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"PendingDelegation\\\" storage:\\\"slot=3,offset=0\\\"\"; Undelegation struct { Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=0,offset=0\\\"\"; Epoch *big.Int \"json:\\\"Epoch\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"Undelegation\\\" storage:\\\"slot=5,offset=0\\\"\"; Renewal struct { Status uint8 \"json:\\\"Status\\\" storage:\\\"slot=0,offset=0\\\"\"; UpdateHeight *big.Int \"json:\\\"UpdateHeight\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"Renewal\\\" storage:\\\"slot=7,offset=0\\\"\" } \"json:\\\"Entry\\\" storage:\\\"slot=0,offset=0\\\"\"; Index *big.Int \"json:\\\"Index\\\" storage:\\\"slot=9,offset=0\\\"\" }"
type Map_2f2f5e901078 map[Address]*MicrodelegationMapEntry_

type Storage_Map_2f2f5e901078 struct {

	// &wrapper.tmplStruct{T:0xe, Name:"Map_2f2f5e901078", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054f450), (*wrapper.tmplField)(0xc000422000)}, Type:"map[Address]*MicrodelegationMapEntry_", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(MicrodelegationMapEntry_)87_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc000550380), (*abi.Type)(0xc0000d2180)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           Map_2f2f5e901078
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Map_2f2f5e901078) Get(key Address) *Storage_MicrodelegationMapEntry_ {
	// Key:	&wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"MicrodelegationMapEntry_", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc000550400), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	keyBytes := key.Hash().Bytes()
	hash := crypto.Keccak256Hash(append(keyBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()
	offset := 0
	if s.obj[key] == nil {

		instance := new(MicrodelegationMapEntry_)

		s.obj[key] = instance
	}

	return &Storage_MicrodelegationMapEntry_{
		obj:           s.obj[key],
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 320,
		dirty:         s.dirty,
	}
}

// Map_619a78bdf676 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xe, Name:"Map_619a78bdf676", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054e5b0), (*wrapper.tmplField)(0xc000422820)}, Type:"map[Address]*Map3NodeWrapperMapEntry_", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(Map3NodeWrapperMapEntry_)116_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc000126c00), (*abi.Type)(0xc0000d2800)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "map[common.Address]*struct { Entry struct { Map3Node struct { Map3Address common.Address \"json:\\\"Map3Address\\\" storage:\\\"slot=0,offset=0\\\"\"; OperatorAddress common.Address \"json:\\\"OperatorAddress\\\" storage:\\\"slot=1,offset=0\\\"\"; NodeKeys struct { Keys []*struct { Key [48]uint8 \"json:\\\"Key\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"Keys\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"NodeKeys\\\" storage:\\\"slot=2,offset=0\\\"\"; Commission struct { Rate struct { F *big.Int \"json:\\\"f\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"Rate\\\" storage:\\\"slot=0,offset=0\\\"\"; RateForNextPeriod struct { F *big.Int \"json:\\\"f\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"RateForNextPeriod\\\" storage:\\\"slot=1,offset=0\\\"\"; UpdateHeight *big.Int \"json:\\\"UpdateHeight\\\" storage:\\\"slot=2,offset=0\\\"\" } \"json:\\\"Commission\\\" storage:\\\"slot=3,offset=0\\\"\"; Description struct { Name string \"json:\\\"Name\\\" storage:\\\"slot=0,offset=0\\\"\"; Identity string \"json:\\\"Identity\\\" storage:\\\"slot=1,offset=0\\\"\"; Website string \"json:\\\"Website\\\" storage:\\\"slot=2,offset=0\\\"\"; SecurityContact string \"json:\\\"SecurityContact\\\" storage:\\\"slot=3,offset=0\\\"\"; Details string \"json:\\\"Details\\\" storage:\\\"slot=4,offset=0\\\"\" } \"json:\\\"Description\\\" storage:\\\"slot=6,offset=0\\\"\"; CreationHeight *big.Int \"json:\\\"CreationHeight\\\" storage:\\\"slot=11,offset=0\\\"\"; Age struct { F *big.Int \"json:\\\"f\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"Age\\\" storage:\\\"slot=12,offset=0\\\"\"; Status uint8 \"json:\\\"Status\\\" storage:\\\"slot=13,offset=0\\\"\"; PendingEpoch *big.Int \"json:\\\"PendingEpoch\\\" storage:\\\"slot=14,offset=0\\\"\"; ActivationEpoch *big.Int \"json:\\\"ActivationEpoch\\\" storage:\\\"slot=15,offset=0\\\"\"; ReleaseEpoch struct { F *big.Int \"json:\\\"f\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"ReleaseEpoch\\\" storage:\\\"slot=16,offset=0\\\"\" } \"json:\\\"Map3Node\\\" storage:\\\"slot=0,offset=0\\\"\"; Microdelegations struct { Keys []*common.Address \"json:\\\"Keys\\\" storage:\\\"slot=0,offset=0\\\"\"; Map map[common.Address]*struct { Entry struct { DelegatorAddress common.Address \"json:\\\"DelegatorAddress\\\" storage:\\\"slot=0,offset=0\\\"\"; Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=1,offset=0\\\"\"; Reward *big.Int \"json:\\\"Reward\\\" storage:\\\"slot=2,offset=0\\\"\"; PendingDelegation struct { Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=0,offset=0\\\"\"; UnlockedEpoch struct { F *big.Int \"json:\\\"f\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"UnlockedEpoch\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"PendingDelegation\\\" storage:\\\"slot=3,offset=0\\\"\"; Undelegation struct { Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=0,offset=0\\\"\"; Epoch *big.Int \"json:\\\"Epoch\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"Undelegation\\\" storage:\\\"slot=5,offset=0\\\"\"; Renewal struct { Status uint8 \"json:\\\"Status\\\" storage:\\\"slot=0,offset=0\\\"\"; UpdateHeight *big.Int \"json:\\\"UpdateHeight\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"Renewal\\\" storage:\\\"slot=7,offset=0\\\"\" } \"json:\\\"Entry\\\" storage:\\\"slot=0,offset=0\\\"\"; Index *big.Int \"json:\\\"Index\\\" storage:\\\"slot=9,offset=0\\\"\" } \"json:\\\"Map\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"Microdelegations\\\" storage:\\\"slot=17,offset=0\\\"\"; RestakingReference struct { ValidatorAddress common.Address \"json:\\\"ValidatorAddress\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"RestakingReference\\\" storage:\\\"slot=19,offset=0\\\"\"; AccumulatedReward *big.Int \"json:\\\"AccumulatedReward\\\" storage:\\\"slot=20,offset=0\\\"\"; TotalDelegation *big.Int \"json:\\\"TotalDelegation\\\" storage:\\\"slot=21,offset=0\\\"\"; TotalPendingDelegation *big.Int \"json:\\\"TotalPendingDelegation\\\" storage:\\\"slot=22,offset=0\\\"\" } \"json:\\\"Entry\\\" storage:\\\"slot=0,offset=0\\\"\"; Index *big.Int \"json:\\\"Index\\\" storage:\\\"slot=23,offset=0\\\"\" }"
type Map_619a78bdf676 map[Address]*Map3NodeWrapperMapEntry_

type Storage_Map_619a78bdf676 struct {

	// &wrapper.tmplStruct{T:0xe, Name:"Map_619a78bdf676", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054e5b0), (*wrapper.tmplField)(0xc000422820)}, Type:"map[Address]*Map3NodeWrapperMapEntry_", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(Map3NodeWrapperMapEntry_)116_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc000126c00), (*abi.Type)(0xc0000d2800)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           Map_619a78bdf676
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Map_619a78bdf676) Get(key Address) *Storage_Map3NodeWrapperMapEntry_ {
	// Key:	&wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"Map3NodeWrapperMapEntry_", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc000126c80), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	keyBytes := key.Hash().Bytes()
	hash := crypto.Keccak256Hash(append(keyBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()
	offset := 0
	if s.obj[key] == nil {

		instance := new(Map3NodeWrapperMapEntry_)

		s.obj[key] = instance
	}

	return &Storage_Map3NodeWrapperMapEntry_{
		obj:           s.obj[key],
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 768,
		dirty:         s.dirty,
	}
}

// Map_624bcc0df495 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xe, Name:"Map_624bcc0df495", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054e340), (*wrapper.tmplField)(0xc00054e410)}, Type:"map[String]*Bool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_string_memory_ptr,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000d3500), (*abi.Type)(0xc0000d3680)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "map[string]*bool"
type Map_624bcc0df495 map[String]*Bool

type Storage_Map_624bcc0df495 struct {

	// &wrapper.tmplStruct{T:0xe, Name:"Map_624bcc0df495", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054e340), (*wrapper.tmplField)(0xc00054e410)}, Type:"map[String]*Bool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_string_memory_ptr,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000d3500), (*abi.Type)(0xc0000d3680)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
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
	// Value: &wrapper.tmplField{Type:"Bool", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000d3600), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
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

// Map_972804fe1df3 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xe, Name:"Map_972804fe1df3", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000422a90), (*wrapper.tmplField)(0xc000422b60)}, Type:"map[Address]*DelegationIndexMap_", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(DelegationIndexMap_)142_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000d2a00), (*abi.Type)(0xc0000d3380)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "map[common.Address]*struct { Keys []*common.Address \"json:\\\"Keys\\\" storage:\\\"slot=0,offset=0\\\"\"; Map map[common.Address]*struct { Entry struct { Map3Address common.Address \"json:\\\"Map3Address\\\" storage:\\\"slot=0,offset=0\\\"\"; IsOperator bool \"json:\\\"IsOperator\\\" storage:\\\"slot=0,offset=20\\\"\" } \"json:\\\"Entry\\\" storage:\\\"slot=0,offset=0\\\"\"; Index *big.Int \"json:\\\"Index\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"Map\\\" storage:\\\"slot=1,offset=0\\\"\" }"
type Map_972804fe1df3 map[Address]*DelegationIndexMap_

type Storage_Map_972804fe1df3 struct {

	// &wrapper.tmplStruct{T:0xe, Name:"Map_972804fe1df3", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000422a90), (*wrapper.tmplField)(0xc000422b60)}, Type:"map[Address]*DelegationIndexMap_", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(DelegationIndexMap_)142_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000d2a00), (*abi.Type)(0xc0000d3380)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           Map_972804fe1df3
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Map_972804fe1df3) Get(key Address) *Storage_DelegationIndexMap_ {
	// Key:	&wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"DelegationIndexMap_", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000d2c00), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	keyBytes := key.Hash().Bytes()
	hash := crypto.Keccak256Hash(append(keyBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()
	offset := 0
	if s.obj[key] == nil {

		instance := new(DelegationIndexMap_)

		s.obj[key] = instance
	}

	return &Storage_DelegationIndexMap_{
		obj:           s.obj[key],
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 64,
		dirty:         s.dirty,
	}
}

// Map_f414efdd5b9c is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xe, Name:"Map_f414efdd5b9c", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000377930), (*wrapper.tmplField)(0xc000377d40)}, Type:"map[Address]*DelegationIndexMapEntry_", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(DelegationIndexMapEntry_)134_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000d2d80), (*abi.Type)(0xc0000d3200)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "map[common.Address]*struct { Entry struct { Map3Address common.Address \"json:\\\"Map3Address\\\" storage:\\\"slot=0,offset=0\\\"\"; IsOperator bool \"json:\\\"IsOperator\\\" storage:\\\"slot=0,offset=20\\\"\" } \"json:\\\"Entry\\\" storage:\\\"slot=0,offset=0\\\"\"; Index *big.Int \"json:\\\"Index\\\" storage:\\\"slot=1,offset=0\\\"\" }"
type Map_f414efdd5b9c map[Address]*DelegationIndexMapEntry_

type Storage_Map_f414efdd5b9c struct {

	// &wrapper.tmplStruct{T:0xe, Name:"Map_f414efdd5b9c", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000377930), (*wrapper.tmplField)(0xc000377d40)}, Type:"map[Address]*DelegationIndexMapEntry_", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(DelegationIndexMapEntry_)134_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000d2d80), (*abi.Type)(0xc0000d3200)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           Map_f414efdd5b9c
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Map_f414efdd5b9c) Get(key Address) *Storage_DelegationIndexMapEntry_ {
	// Key:	&wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"DelegationIndexMapEntry_", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000d2e00), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	keyBytes := key.Hash().Bytes()
	hash := crypto.Keccak256Hash(append(keyBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()
	offset := 0
	if s.obj[key] == nil {

		instance := new(DelegationIndexMapEntry_)

		s.obj[key] = instance
	}

	return &Storage_DelegationIndexMapEntry_{
		obj:           s.obj[key],
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 64,
		dirty:         s.dirty,
	}
}

// Slice_2f1765e04dad is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x4, Name:"Slice_2f1765e04dad", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054e9c0)}, Type:"[]*BLSPublicKey_", SolKind:abi.Type{Elem:(*abi.Type)(0xc000127300), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_struct(BLSPublicKey_)20_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "[]*struct { Key [48]uint8 \"json:\\\"Key\\\" storage:\\\"slot=0,offset=0\\\"\" }"
type Slice_2f1765e04dad []*BLSPublicKey_

type Storage_Slice_2f1765e04dad struct {

	// &wrapper.tmplStruct{T:0x4, Name:"Slice_2f1765e04dad", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054e9c0)}, Type:"[]*BLSPublicKey_", SolKind:abi.Type{Elem:(*abi.Type)(0xc000127300), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_struct(BLSPublicKey_)20_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           *Slice_2f1765e04dad
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Slice_2f1765e04dad) Length() int {
	rv := s.db.GetState(s.addr, common.BigToHash(s.slot))
	return int(rv.Big().Int64())
}

func (s *Storage_Slice_2f1765e04dad) Resize(length int) {
	// Value: &wrapper.tmplField{Type:"BLSPublicKey_", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc000127000), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(big.NewInt(0).SetUint64(uint64(length))))

	_length := length
	_capacity := length + 50

	slice := make([]*BLSPublicKey_, _length, _capacity)

	copy(slice, *s.obj)
	*s.obj = slice
}

func (s *Storage_Slice_2f1765e04dad) Get(index int) *Storage_BLSPublicKey_ {
	// Value: &wrapper.tmplField{Type:"BLSPublicKey_", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc000127000), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	length := s.Length()
	if length <= index {
		s.Resize(index + 1)
	}

	base := crypto.Keccak256Hash(common.BigToHash(s.slot).Bytes()).Big()

	slotsPerItem := (64 + 31) / 32
	actual := big.NewInt(0).Add(base, big.NewInt(0).SetUint64(uint64(index*slotsPerItem)))
	offset := 0

	if (*s.obj)[index] == nil {

		instance := new(BLSPublicKey_)

		(*s.obj)[index] = instance
	}

	return &Storage_BLSPublicKey_{
		obj:           (*s.obj)[index],
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 64,
		dirty:         s.dirty,
	}
}

// Slice_760032d9d704 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x4, Name:"Slice_760032d9d704", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000377790)}, Type:"[]*Address", SolKind:abi.Type{Elem:(*abi.Type)(0xc000126a80), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_address)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "[]*common.Address"
type Slice_760032d9d704 []*Address

type Storage_Slice_760032d9d704 struct {

	// &wrapper.tmplStruct{T:0x4, Name:"Slice_760032d9d704", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000377790)}, Type:"[]*Address", SolKind:abi.Type{Elem:(*abi.Type)(0xc000126a80), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_address)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
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
	// Value: &wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc000126980), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(big.NewInt(0).SetUint64(uint64(length))))

	_length := length
	_capacity := length + 50

	slice := make([]*Address, _length, _capacity)

	copy(slice, *s.obj)
	*s.obj = slice
}

func (s *Storage_Slice_760032d9d704) Get(index int) *Storage_Address {
	// Value: &wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc000126980), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
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

// BLSPublicKey_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"BLSPublicKey_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054e8f0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(BLSPublicKey_)20_storage", TupleRawName:"BLSPublicKey_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000127080)}, TupleRawNames:[]string{"Key"}, TupleType:(*reflect.rtype)(0xc00054a480)}}
type BLSPublicKey_ struct {
	Key Array_aaa956d410fd `json:"Key" storage:"slot=0,offset=0"`
}

// BLSPublicKeys_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"BLSPublicKeys_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054ea90)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(BLSPublicKeys_)24_storage", TupleRawName:"BLSPublicKeys_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000126f80)}, TupleRawNames:[]string{"Keys"}, TupleType:(*reflect.rtype)(0xc00054a4e0)}}
type BLSPublicKeys_ struct {
	Keys Slice_2f1765e04dad `json:"Keys" storage:"slot=0,offset=0"`
}

// Commission_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Commission_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000377110), (*wrapper.tmplField)(0xc0003771e0), (*wrapper.tmplField)(0xc0003772b0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x60, stringKind:"t_struct(Commission_)31_storage", TupleRawName:"Commission_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000127500), (*abi.Type)(0xc000127700), (*abi.Type)(0xc000127780)}, TupleRawNames:[]string{"Rate", "RateForNextPeriod", "UpdateHeight"}, TupleType:(*reflect.rtype)(0xc00054a6c0)}}
type Commission_ struct {
	Rate              Decimal `json:"Rate" storage:"slot=0,offset=0"`
	RateForNextPeriod Decimal `json:"RateForNextPeriod" storage:"slot=1,offset=0"`
	UpdateHeight      BigInt  `json:"UpdateHeight" storage:"slot=2,offset=0"`
}

// DelegationIndexMapEntry_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"DelegationIndexMapEntry_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000377ba0), (*wrapper.tmplField)(0xc000377c70)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(DelegationIndexMapEntry_)134_storage", TupleRawName:"DelegationIndexMapEntry_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000d2e80), (*abi.Type)(0xc0000d3100)}, TupleRawNames:[]string{"Entry", "Index"}, TupleType:(*reflect.rtype)(0xc0003c8ba0)}}
type DelegationIndexMapEntry_ struct {
	Entry DelegationIndex_ `json:"Entry" storage:"slot=0,offset=0"`
	Index BigInt           `json:"Index" storage:"slot=1,offset=0"`
}

// DelegationIndexMap_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"DelegationIndexMap_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000377860), (*wrapper.tmplField)(0xc000377e10)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(DelegationIndexMap_)142_storage", TupleRawName:"DelegationIndexMap_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000d2c80), (*abi.Type)(0xc0000d2d00)}, TupleRawNames:[]string{"Keys", "Map"}, TupleType:(*reflect.rtype)(0xc0003c8c60)}}
type DelegationIndexMap_ struct {
	Keys Slice_760032d9d704 `json:"Keys" storage:"slot=0,offset=0"`
	Map  Map_f414efdd5b9c   `json:"Map" storage:"slot=1,offset=0"`
}

// DelegationIndex_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"DelegationIndex_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000377a00), (*wrapper.tmplField)(0xc000377ad0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(DelegationIndex_)129_storage", TupleRawName:"DelegationIndex_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000d2f00), (*abi.Type)(0xc0000d2f80)}, TupleRawNames:[]string{"Map3Address", "IsOperator"}, TupleType:(*reflect.rtype)(0xc0003c8b40)}}
type DelegationIndex_ struct {
	Map3Address Address `json:"Map3Address" storage:"slot=0,offset=0"`
	IsOperator  Bool    `json:"IsOperator" storage:"slot=0,offset=20"`
}

// Description_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Description_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000377380), (*wrapper.tmplField)(0xc000377450), (*wrapper.tmplField)(0xc000377520), (*wrapper.tmplField)(0xc0003775f0), (*wrapper.tmplField)(0xc0003776c0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xa0, stringKind:"t_struct(Description_)15_storage", TupleRawName:"Description_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000127900), (*abi.Type)(0xc000127a00), (*abi.Type)(0xc000127a80), (*abi.Type)(0xc000127b00), (*abi.Type)(0xc000127b80)}, TupleRawNames:[]string{"Name", "Identity", "Website", "SecurityContact", "Details"}, TupleType:(*reflect.rtype)(0xc00054a7e0)}}
type Description_ struct {
	Name            String `json:"Name" storage:"slot=0,offset=0"`
	Identity        String `json:"Identity" storage:"slot=1,offset=0"`
	Website         String `json:"Website" storage:"slot=2,offset=0"`
	SecurityContact String `json:"SecurityContact" storage:"slot=3,offset=0"`
	Details         String `json:"Details" storage:"slot=4,offset=0"`
}

// Global_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Global_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000422ea0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x0, stringKind:"t_struct(Global_t)", TupleRawName:"Global_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000d3880)}, TupleRawNames:[]string{"Map3NodePool"}, TupleType:(*reflect.rtype)(0xc0003c8f60)}}
type Global_t struct {
	Map3NodePool Map3NodePool_ `json:"Map3NodePool" storage:"slot=0,offset=0"`
}

// Map3NodePool_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Map3NodePool_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0004229c0), (*wrapper.tmplField)(0xc000422c30), (*wrapper.tmplField)(0xc000422d00), (*wrapper.tmplField)(0xc000422dd0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xa0, stringKind:"t_struct(Map3NodePool_)157_storage", TupleRawName:"Map3NodePool_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000126880), (*abi.Type)(0xc0000d2980), (*abi.Type)(0xc0000d3480), (*abi.Type)(0xc0000d3780)}, TupleRawNames:[]string{"Nodes", "DelegationIndexMapByDelegator", "NodeKeySet", "DescriptionIdentitySet"}, TupleType:(*reflect.rtype)(0xc0003c8ea0)}}
type Map3NodePool_ struct {
	Nodes                         Map3NodeWrapperMap_ `json:"Nodes" storage:"slot=0,offset=0"`
	DelegationIndexMapByDelegator Map_972804fe1df3    `json:"DelegationIndexMapByDelegator" storage:"slot=2,offset=0"`
	NodeKeySet                    Map_624bcc0df495    `json:"NodeKeySet" storage:"slot=3,offset=0"`
	DescriptionIdentitySet        Map_624bcc0df495    `json:"DescriptionIdentitySet" storage:"slot=4,offset=0"`
}

// Map3NodeWrapperMapEntry_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Map3NodeWrapperMapEntry_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000422680), (*wrapper.tmplField)(0xc000422750)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x300, stringKind:"t_struct(Map3NodeWrapperMapEntry_)116_storage", TupleRawName:"Map3NodeWrapperMapEntry_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000126d00), (*abi.Type)(0xc0000d2700)}, TupleRawNames:[]string{"Entry", "Index"}, TupleType:(*reflect.rtype)(0xc0003c82a0)}}
type Map3NodeWrapperMapEntry_ struct {
	Entry Map3NodeWrapper_ `json:"Entry" storage:"slot=0,offset=0"`
	Index BigInt           `json:"Index" storage:"slot=23,offset=0"`
}

// Map3NodeWrapperMap_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Map3NodeWrapperMap_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054e4e0), (*wrapper.tmplField)(0xc0004228f0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Map3NodeWrapperMap_)124_storage", TupleRawName:"Map3NodeWrapperMap_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000126900), (*abi.Type)(0xc000126b80)}, TupleRawNames:[]string{"Keys", "Map"}, TupleType:(*reflect.rtype)(0xc0003c8360)}}
type Map3NodeWrapperMap_ struct {
	Keys Slice_760032d9d704 `json:"Keys" storage:"slot=0,offset=0"`
	Map  Map_619a78bdf676   `json:"Map" storage:"slot=1,offset=0"`
}

// Map3NodeWrapper_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Map3NodeWrapper_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054f2b0), (*wrapper.tmplField)(0xc000422270), (*wrapper.tmplField)(0xc000422340), (*wrapper.tmplField)(0xc000422410), (*wrapper.tmplField)(0xc0004224e0), (*wrapper.tmplField)(0xc0004225b0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x2e0, stringKind:"t_struct(Map3NodeWrapper_)111_storage", TupleRawName:"Map3NodeWrapper_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000126d80), (*abi.Type)(0xc000550200), (*abi.Type)(0xc0000d2300), (*abi.Type)(0xc0000d2480), (*abi.Type)(0xc0000d2500), (*abi.Type)(0xc0000d2600)}, TupleRawNames:[]string{"Map3Node", "Microdelegations", "RestakingReference", "AccumulatedReward", "TotalDelegation", "TotalPendingDelegation"}, TupleType:(*reflect.rtype)(0xc0003c8240)}}
type Map3NodeWrapper_ struct {
	Map3Node               Map3Node_           `json:"Map3Node" storage:"slot=0,offset=0"`
	Microdelegations       MicrodelegationMap_ `json:"Microdelegations" storage:"slot=17,offset=0"`
	RestakingReference     RestakingReference_ `json:"RestakingReference" storage:"slot=19,offset=0"`
	AccumulatedReward      BigInt              `json:"AccumulatedReward" storage:"slot=20,offset=0"`
	TotalDelegation        BigInt              `json:"TotalDelegation" storage:"slot=21,offset=0"`
	TotalPendingDelegation BigInt              `json:"TotalPendingDelegation" storage:"slot=22,offset=0"`
}

// Map3Node_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Map3Node_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054e680), (*wrapper.tmplField)(0xc00054e750), (*wrapper.tmplField)(0xc00054eb60), (*wrapper.tmplField)(0xc00054ec30), (*wrapper.tmplField)(0xc00054ed00), (*wrapper.tmplField)(0xc00054edd0), (*wrapper.tmplField)(0xc00054eea0), (*wrapper.tmplField)(0xc00054ef70), (*wrapper.tmplField)(0xc00054f040), (*wrapper.tmplField)(0xc00054f110), (*wrapper.tmplField)(0xc00054f1e0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x220, stringKind:"t_struct(Map3Node_)54_storage", TupleRawName:"Map3Node_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000126e00), (*abi.Type)(0xc000126e80), (*abi.Type)(0xc000126f00), (*abi.Type)(0xc000127480), (*abi.Type)(0xc000127880), (*abi.Type)(0xc000127e00), (*abi.Type)(0xc000127e80), (*abi.Type)(0xc000127f00), (*abi.Type)(0xc000127f80), (*abi.Type)(0xc000550080), (*abi.Type)(0xc000550100)}, TupleRawNames:[]string{"Map3Address", "OperatorAddress", "NodeKeys", "Commission", "Description", "CreationHeight", "Age", "Status", "PendingEpoch", "ActivationEpoch", "ReleaseEpoch"}, TupleType:(*reflect.rtype)(0xc00054a840)}}
type Map3Node_ struct {
	Map3Address     Address        `json:"Map3Address" storage:"slot=0,offset=0"`
	OperatorAddress Address        `json:"OperatorAddress" storage:"slot=1,offset=0"`
	NodeKeys        BLSPublicKeys_ `json:"NodeKeys" storage:"slot=2,offset=0"`
	Commission      Commission_    `json:"Commission" storage:"slot=3,offset=0"`
	Description     Description_   `json:"Description" storage:"slot=6,offset=0"`
	CreationHeight  BigInt         `json:"CreationHeight" storage:"slot=11,offset=0"`
	Age             Decimal        `json:"Age" storage:"slot=12,offset=0"`
	Status          Uint8          `json:"Status" storage:"slot=13,offset=0"`
	PendingEpoch    BigInt         `json:"PendingEpoch" storage:"slot=14,offset=0"`
	ActivationEpoch BigInt         `json:"ActivationEpoch" storage:"slot=15,offset=0"`
	ReleaseEpoch    Decimal        `json:"ReleaseEpoch" storage:"slot=16,offset=0"`
}

// MicrodelegationMapEntry_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"MicrodelegationMapEntry_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054fd40), (*wrapper.tmplField)(0xc00054fe10)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x140, stringKind:"t_struct(MicrodelegationMapEntry_)87_storage", TupleRawName:"MicrodelegationMapEntry_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000550480), (*abi.Type)(0xc0000d2080)}, TupleRawNames:[]string{"Entry", "Index"}, TupleType:(*reflect.rtype)(0xc0003c8000)}}
type MicrodelegationMapEntry_ struct {
	Entry Microdelegation_ `json:"Entry" storage:"slot=0,offset=0"`
	Index BigInt           `json:"Index" storage:"slot=9,offset=0"`
}

// MicrodelegationMap_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"MicrodelegationMap_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054f380), (*wrapper.tmplField)(0xc0004221a0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(MicrodelegationMap_)95_storage", TupleRawName:"MicrodelegationMap_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000550280), (*abi.Type)(0xc000550300)}, TupleRawNames:[]string{"Keys", "Map"}, TupleType:(*reflect.rtype)(0xc0003c80c0)}}
type MicrodelegationMap_ struct {
	Keys Slice_760032d9d704 `json:"Keys" storage:"slot=0,offset=0"`
	Map  Map_2f2f5e901078   `json:"Map" storage:"slot=1,offset=0"`
}

// Microdelegation_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Microdelegation_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054f520), (*wrapper.tmplField)(0xc00054f5f0), (*wrapper.tmplField)(0xc00054f6c0), (*wrapper.tmplField)(0xc00054f930), (*wrapper.tmplField)(0xc00054fa00), (*wrapper.tmplField)(0xc00054fc70)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x120, stringKind:"t_struct(Microdelegation_)77_storage", TupleRawName:"Microdelegation_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000550500), (*abi.Type)(0xc000550580), (*abi.Type)(0xc000550600), (*abi.Type)(0xc000550680), (*abi.Type)(0xc000550880), (*abi.Type)(0xc000550b00)}, TupleRawNames:[]string{"DelegatorAddress", "Amount", "Reward", "PendingDelegation", "Undelegation", "Renewal"}, TupleType:(*reflect.rtype)(0xc000704000)}}
type Microdelegation_ struct {
	DelegatorAddress  Address            `json:"DelegatorAddress" storage:"slot=0,offset=0"`
	Amount            BigInt             `json:"Amount" storage:"slot=1,offset=0"`
	Reward            BigInt             `json:"Reward" storage:"slot=2,offset=0"`
	PendingDelegation PendingDelegation_ `json:"PendingDelegation" storage:"slot=3,offset=0"`
	Undelegation      Undelegation_      `json:"Undelegation" storage:"slot=5,offset=0"`
	Renewal           Renewal_           `json:"Renewal" storage:"slot=7,offset=0"`
}

// PendingDelegation_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"PendingDelegation_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054f790), (*wrapper.tmplField)(0xc00054f860)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(PendingDelegation_)59_storage", TupleRawName:"PendingDelegation_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000550700), (*abi.Type)(0xc000550780)}, TupleRawNames:[]string{"Amount", "UnlockedEpoch"}, TupleType:(*reflect.rtype)(0xc00054aa80)}}
type PendingDelegation_ struct {
	Amount        BigInt  `json:"Amount" storage:"slot=0,offset=0"`
	UnlockedEpoch Decimal `json:"UnlockedEpoch" storage:"slot=1,offset=0"`
}

// Renewal_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Renewal_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054fad0), (*wrapper.tmplField)(0xc00054fba0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Renewal_)82_storage", TupleRawName:"Renewal_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000550b80), (*abi.Type)(0xc000550c00)}, TupleRawNames:[]string{"Status", "UpdateHeight"}, TupleType:(*reflect.rtype)(0xc00054ac00)}}
type Renewal_ struct {
	Status       Uint8  `json:"Status" storage:"slot=0,offset=0"`
	UpdateHeight BigInt `json:"UpdateHeight" storage:"slot=1,offset=0"`
}

// RestakingReference_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"RestakingReference_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054e270)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(RestakingReference_)98_storage", TupleRawName:"RestakingReference_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000d2380)}, TupleRawNames:[]string{"ValidatorAddress"}, TupleType:(*reflect.rtype)(0xc0003c81e0)}}
type RestakingReference_ struct {
	ValidatorAddress Address `json:"ValidatorAddress" storage:"slot=0,offset=0"`
}

// Undelegation_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Undelegation_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054e000), (*wrapper.tmplField)(0xc00054e1a0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Undelegation_)64_storage", TupleRawName:"Undelegation_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000550900), (*abi.Type)(0xc000550980)}, TupleRawNames:[]string{"Amount", "Epoch"}, TupleType:(*reflect.rtype)(0xc00054ab40)}}
type Undelegation_ struct {
	Amount BigInt `json:"Amount" storage:"slot=0,offset=0"`
	Epoch  BigInt `json:"Epoch" storage:"slot=1,offset=0"`
}

type Storage_BLSPublicKey_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"BLSPublicKey_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054e8f0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(BLSPublicKey_)20_storage", TupleRawName:"BLSPublicKey_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000127080)}, TupleRawNames:[]string{"Key"}, TupleType:(*reflect.rtype)(0xc00054a480)}}
	obj           *BLSPublicKey_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_BLSPublicKeys_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"BLSPublicKeys_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054ea90)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(BLSPublicKeys_)24_storage", TupleRawName:"BLSPublicKeys_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000126f80)}, TupleRawNames:[]string{"Keys"}, TupleType:(*reflect.rtype)(0xc00054a4e0)}}
	obj           *BLSPublicKeys_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Commission_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Commission_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000377110), (*wrapper.tmplField)(0xc0003771e0), (*wrapper.tmplField)(0xc0003772b0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x60, stringKind:"t_struct(Commission_)31_storage", TupleRawName:"Commission_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000127500), (*abi.Type)(0xc000127700), (*abi.Type)(0xc000127780)}, TupleRawNames:[]string{"Rate", "RateForNextPeriod", "UpdateHeight"}, TupleType:(*reflect.rtype)(0xc00054a6c0)}}
	obj           *Commission_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_DelegationIndexMapEntry_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"DelegationIndexMapEntry_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000377ba0), (*wrapper.tmplField)(0xc000377c70)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(DelegationIndexMapEntry_)134_storage", TupleRawName:"DelegationIndexMapEntry_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000d2e80), (*abi.Type)(0xc0000d3100)}, TupleRawNames:[]string{"Entry", "Index"}, TupleType:(*reflect.rtype)(0xc0003c8ba0)}}
	obj           *DelegationIndexMapEntry_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_DelegationIndexMap_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"DelegationIndexMap_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000377860), (*wrapper.tmplField)(0xc000377e10)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(DelegationIndexMap_)142_storage", TupleRawName:"DelegationIndexMap_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000d2c80), (*abi.Type)(0xc0000d2d00)}, TupleRawNames:[]string{"Keys", "Map"}, TupleType:(*reflect.rtype)(0xc0003c8c60)}}
	obj           *DelegationIndexMap_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_DelegationIndex_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"DelegationIndex_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000377a00), (*wrapper.tmplField)(0xc000377ad0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(DelegationIndex_)129_storage", TupleRawName:"DelegationIndex_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000d2f00), (*abi.Type)(0xc0000d2f80)}, TupleRawNames:[]string{"Map3Address", "IsOperator"}, TupleType:(*reflect.rtype)(0xc0003c8b40)}}
	obj           *DelegationIndex_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Description_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Description_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000377380), (*wrapper.tmplField)(0xc000377450), (*wrapper.tmplField)(0xc000377520), (*wrapper.tmplField)(0xc0003775f0), (*wrapper.tmplField)(0xc0003776c0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xa0, stringKind:"t_struct(Description_)15_storage", TupleRawName:"Description_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000127900), (*abi.Type)(0xc000127a00), (*abi.Type)(0xc000127a80), (*abi.Type)(0xc000127b00), (*abi.Type)(0xc000127b80)}, TupleRawNames:[]string{"Name", "Identity", "Website", "SecurityContact", "Details"}, TupleType:(*reflect.rtype)(0xc00054a7e0)}}
	obj           *Description_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Global_t struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Global_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000422ea0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x0, stringKind:"t_struct(Global_t)", TupleRawName:"Global_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000d3880)}, TupleRawNames:[]string{"Map3NodePool"}, TupleType:(*reflect.rtype)(0xc0003c8f60)}}
	obj           *Global_t
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Map3NodePool_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Map3NodePool_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0004229c0), (*wrapper.tmplField)(0xc000422c30), (*wrapper.tmplField)(0xc000422d00), (*wrapper.tmplField)(0xc000422dd0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xa0, stringKind:"t_struct(Map3NodePool_)157_storage", TupleRawName:"Map3NodePool_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000126880), (*abi.Type)(0xc0000d2980), (*abi.Type)(0xc0000d3480), (*abi.Type)(0xc0000d3780)}, TupleRawNames:[]string{"Nodes", "DelegationIndexMapByDelegator", "NodeKeySet", "DescriptionIdentitySet"}, TupleType:(*reflect.rtype)(0xc0003c8ea0)}}
	obj           *Map3NodePool_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Map3NodeWrapperMapEntry_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Map3NodeWrapperMapEntry_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000422680), (*wrapper.tmplField)(0xc000422750)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x300, stringKind:"t_struct(Map3NodeWrapperMapEntry_)116_storage", TupleRawName:"Map3NodeWrapperMapEntry_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000126d00), (*abi.Type)(0xc0000d2700)}, TupleRawNames:[]string{"Entry", "Index"}, TupleType:(*reflect.rtype)(0xc0003c82a0)}}
	obj           *Map3NodeWrapperMapEntry_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Map3NodeWrapperMap_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Map3NodeWrapperMap_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054e4e0), (*wrapper.tmplField)(0xc0004228f0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Map3NodeWrapperMap_)124_storage", TupleRawName:"Map3NodeWrapperMap_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000126900), (*abi.Type)(0xc000126b80)}, TupleRawNames:[]string{"Keys", "Map"}, TupleType:(*reflect.rtype)(0xc0003c8360)}}
	obj           *Map3NodeWrapperMap_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Map3NodeWrapper_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Map3NodeWrapper_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054f2b0), (*wrapper.tmplField)(0xc000422270), (*wrapper.tmplField)(0xc000422340), (*wrapper.tmplField)(0xc000422410), (*wrapper.tmplField)(0xc0004224e0), (*wrapper.tmplField)(0xc0004225b0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x2e0, stringKind:"t_struct(Map3NodeWrapper_)111_storage", TupleRawName:"Map3NodeWrapper_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000126d80), (*abi.Type)(0xc000550200), (*abi.Type)(0xc0000d2300), (*abi.Type)(0xc0000d2480), (*abi.Type)(0xc0000d2500), (*abi.Type)(0xc0000d2600)}, TupleRawNames:[]string{"Map3Node", "Microdelegations", "RestakingReference", "AccumulatedReward", "TotalDelegation", "TotalPendingDelegation"}, TupleType:(*reflect.rtype)(0xc0003c8240)}}
	obj           *Map3NodeWrapper_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Map3Node_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Map3Node_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054e680), (*wrapper.tmplField)(0xc00054e750), (*wrapper.tmplField)(0xc00054eb60), (*wrapper.tmplField)(0xc00054ec30), (*wrapper.tmplField)(0xc00054ed00), (*wrapper.tmplField)(0xc00054edd0), (*wrapper.tmplField)(0xc00054eea0), (*wrapper.tmplField)(0xc00054ef70), (*wrapper.tmplField)(0xc00054f040), (*wrapper.tmplField)(0xc00054f110), (*wrapper.tmplField)(0xc00054f1e0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x220, stringKind:"t_struct(Map3Node_)54_storage", TupleRawName:"Map3Node_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000126e00), (*abi.Type)(0xc000126e80), (*abi.Type)(0xc000126f00), (*abi.Type)(0xc000127480), (*abi.Type)(0xc000127880), (*abi.Type)(0xc000127e00), (*abi.Type)(0xc000127e80), (*abi.Type)(0xc000127f00), (*abi.Type)(0xc000127f80), (*abi.Type)(0xc000550080), (*abi.Type)(0xc000550100)}, TupleRawNames:[]string{"Map3Address", "OperatorAddress", "NodeKeys", "Commission", "Description", "CreationHeight", "Age", "Status", "PendingEpoch", "ActivationEpoch", "ReleaseEpoch"}, TupleType:(*reflect.rtype)(0xc00054a840)}}
	obj           *Map3Node_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_MicrodelegationMapEntry_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"MicrodelegationMapEntry_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054fd40), (*wrapper.tmplField)(0xc00054fe10)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x140, stringKind:"t_struct(MicrodelegationMapEntry_)87_storage", TupleRawName:"MicrodelegationMapEntry_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000550480), (*abi.Type)(0xc0000d2080)}, TupleRawNames:[]string{"Entry", "Index"}, TupleType:(*reflect.rtype)(0xc0003c8000)}}
	obj           *MicrodelegationMapEntry_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_MicrodelegationMap_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"MicrodelegationMap_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054f380), (*wrapper.tmplField)(0xc0004221a0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(MicrodelegationMap_)95_storage", TupleRawName:"MicrodelegationMap_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000550280), (*abi.Type)(0xc000550300)}, TupleRawNames:[]string{"Keys", "Map"}, TupleType:(*reflect.rtype)(0xc0003c80c0)}}
	obj           *MicrodelegationMap_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Microdelegation_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Microdelegation_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054f520), (*wrapper.tmplField)(0xc00054f5f0), (*wrapper.tmplField)(0xc00054f6c0), (*wrapper.tmplField)(0xc00054f930), (*wrapper.tmplField)(0xc00054fa00), (*wrapper.tmplField)(0xc00054fc70)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x120, stringKind:"t_struct(Microdelegation_)77_storage", TupleRawName:"Microdelegation_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000550500), (*abi.Type)(0xc000550580), (*abi.Type)(0xc000550600), (*abi.Type)(0xc000550680), (*abi.Type)(0xc000550880), (*abi.Type)(0xc000550b00)}, TupleRawNames:[]string{"DelegatorAddress", "Amount", "Reward", "PendingDelegation", "Undelegation", "Renewal"}, TupleType:(*reflect.rtype)(0xc000704000)}}
	obj           *Microdelegation_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_PendingDelegation_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"PendingDelegation_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054f790), (*wrapper.tmplField)(0xc00054f860)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(PendingDelegation_)59_storage", TupleRawName:"PendingDelegation_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000550700), (*abi.Type)(0xc000550780)}, TupleRawNames:[]string{"Amount", "UnlockedEpoch"}, TupleType:(*reflect.rtype)(0xc00054aa80)}}
	obj           *PendingDelegation_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Renewal_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Renewal_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054fad0), (*wrapper.tmplField)(0xc00054fba0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Renewal_)82_storage", TupleRawName:"Renewal_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000550b80), (*abi.Type)(0xc000550c00)}, TupleRawNames:[]string{"Status", "UpdateHeight"}, TupleType:(*reflect.rtype)(0xc00054ac00)}}
	obj           *Renewal_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_RestakingReference_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"RestakingReference_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054e270)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(RestakingReference_)98_storage", TupleRawName:"RestakingReference_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000d2380)}, TupleRawNames:[]string{"ValidatorAddress"}, TupleType:(*reflect.rtype)(0xc0003c81e0)}}
	obj           *RestakingReference_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Undelegation_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Undelegation_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00054e000), (*wrapper.tmplField)(0xc00054e1a0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Undelegation_)64_storage", TupleRawName:"Undelegation_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000550900), (*abi.Type)(0xc000550980)}, TupleRawNames:[]string{"Amount", "Epoch"}, TupleType:(*reflect.rtype)(0xc00054ab40)}}
	obj           *Undelegation_
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

// Key is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_BLSPublicKey_) Key() *Storage_Array_aaa956d410fd {
	// Field: &wrapper.tmplField{Type:"Array_aaa956d410fd", ElemType:"Array_aaa956d410fd", Name:"Key", SolKind:abi.Type{Elem:(*abi.Type)(0xc000127100), Size:48, T:0x5, numberOfBytes:0x40, stringKind:"t_array(t_bytes1)48_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Key\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Array_aaa956d410fd{
		obj:           &s.obj.Key,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 64,
		dirty:         s.dirty,
	}
}

// Keys is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_BLSPublicKeys_) Keys() *Storage_Slice_2f1765e04dad {
	// Field: &wrapper.tmplField{Type:"Slice_2f1765e04dad", ElemType:"Slice_2f1765e04dad", Name:"Keys", SolKind:abi.Type{Elem:(*abi.Type)(0xc000127300), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_struct(BLSPublicKey_)20_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Keys\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Keys == nil {

		hash := s.db.GetState(s.addr, common.BigToHash(actual))
		instance := make(Slice_2f1765e04dad, hash.Big().Int64())

		s.obj.Keys = instance
	}

	return &Storage_Slice_2f1765e04dad{
		obj:           &s.obj.Keys,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Rate is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Commission_) Rate() *Storage_Decimal {
	// Field: &wrapper.tmplField{Type:"Decimal", ElemType:"Decimal", Name:"Rate", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc000127580)}, TupleRawNames:[]string{"f"}, TupleType:(*reflect.rtype)(0xc00054a660)}, Tag:"json:\"Rate\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Decimal{
		obj:           &s.obj.Rate,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// RateForNextPeriod is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Commission_) RateForNextPeriod() *Storage_Decimal {
	// Field: &wrapper.tmplField{Type:"Decimal", ElemType:"Decimal", Name:"RateForNextPeriod", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc000127580)}, TupleRawNames:[]string{"f"}, TupleType:(*reflect.rtype)(0xc00054a660)}, Tag:"json:\"RateForNextPeriod\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Decimal{
		obj:           &s.obj.RateForNextPeriod,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// UpdateHeight is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Commission_) UpdateHeight() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"UpdateHeight", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"UpdateHeight\" storage:\"slot=2,offset=0\"", Slot:2, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("2", 10)
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

// Entry is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_DelegationIndexMapEntry_) Entry() *Storage_DelegationIndex_ {
	// Field: &wrapper.tmplField{Type:"DelegationIndex_", ElemType:"DelegationIndex_", Name:"Entry", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(DelegationIndex_)129_storage", TupleRawName:"DelegationIndex_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000d2f00), (*abi.Type)(0xc0000d2f80)}, TupleRawNames:[]string{"Map3Address", "IsOperator"}, TupleType:(*reflect.rtype)(0xc0003c8b40)}, Tag:"json:\"Entry\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_DelegationIndex_{
		obj:           &s.obj.Entry,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Index is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_DelegationIndexMapEntry_) Index() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Index", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Index\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Index == nil {

		instance := big.NewInt(0)

		s.obj.Index = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.Index,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Keys is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_DelegationIndexMap_) Keys() *Storage_Slice_760032d9d704 {
	// Field: &wrapper.tmplField{Type:"Slice_760032d9d704", ElemType:"Slice_760032d9d704", Name:"Keys", SolKind:abi.Type{Elem:(*abi.Type)(0xc000126a80), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_address)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Keys\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Keys == nil {

		hash := s.db.GetState(s.addr, common.BigToHash(actual))
		instance := make(Slice_760032d9d704, hash.Big().Int64())

		s.obj.Keys = instance
	}

	return &Storage_Slice_760032d9d704{
		obj:           &s.obj.Keys,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Map is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_DelegationIndexMap_) Map() *Storage_Map_f414efdd5b9c {
	// Field: &wrapper.tmplField{Type:"Map_f414efdd5b9c", ElemType:"Map_f414efdd5b9c", Name:"Map", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(DelegationIndexMapEntry_)134_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000d2d80), (*abi.Type)(0xc0000d3200)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Map\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Map == nil {

		instance := make(Map_f414efdd5b9c)

		s.obj.Map = instance
	}

	return &Storage_Map_f414efdd5b9c{
		obj:           s.obj.Map,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Map3Address is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_DelegationIndex_) Map3Address() *Storage_Address {
	// Field: &wrapper.tmplField{Type:"Address", ElemType:"Address", Name:"Map3Address", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Map3Address\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Address{
		obj:           &s.obj.Map3Address,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 20,
		dirty:         s.dirty,
	}
}

// IsOperator is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_DelegationIndex_) IsOperator() *Storage_Bool {
	// Field: &wrapper.tmplField{Type:"Bool", ElemType:"Bool", Name:"IsOperator", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x2, numberOfBytes:0x1, stringKind:"t_bool", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"IsOperator\" storage:\"slot=0,offset=20\"", Slot:0, Offset:20}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 20

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Bool{
		obj:           &s.obj.IsOperator,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 1,
		dirty:         s.dirty,
	}
}

// Name is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Description_) Name() *Storage_String {
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
func (s *Storage_Description_) Identity() *Storage_String {
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

// Website is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Description_) Website() *Storage_String {
	// Field: &wrapper.tmplField{Type:"String", ElemType:"String", Name:"Website", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x3, numberOfBytes:0x20, stringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Website\" storage:\"slot=2,offset=0\"", Slot:2, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("2", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_String{
		obj:           &s.obj.Website,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// SecurityContact is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Description_) SecurityContact() *Storage_String {
	// Field: &wrapper.tmplField{Type:"String", ElemType:"String", Name:"SecurityContact", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x3, numberOfBytes:0x20, stringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"SecurityContact\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("3", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_String{
		obj:           &s.obj.SecurityContact,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Details is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Description_) Details() *Storage_String {
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

// Map3NodePool is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Global_t) Map3NodePool() *Storage_Map3NodePool_ {
	// Field: &wrapper.tmplField{Type:"Map3NodePool_", ElemType:"Map3NodePool_", Name:"Map3NodePool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xa0, stringKind:"t_struct(Map3NodePool_)157_storage", TupleRawName:"Map3NodePool_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000126880), (*abi.Type)(0xc0000d2980), (*abi.Type)(0xc0000d3480), (*abi.Type)(0xc0000d3780)}, TupleRawNames:[]string{"Nodes", "DelegationIndexMapByDelegator", "NodeKeySet", "DescriptionIdentitySet"}, TupleType:(*reflect.rtype)(0xc0003c8ea0)}, Tag:"json:\"Map3NodePool\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Map3NodePool_{
		obj:           &s.obj.Map3NodePool,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 160,
		dirty:         s.dirty,
	}
}

// Nodes is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodePool_) Nodes() *Storage_Map3NodeWrapperMap_ {
	// Field: &wrapper.tmplField{Type:"Map3NodeWrapperMap_", ElemType:"Map3NodeWrapperMap_", Name:"Nodes", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Map3NodeWrapperMap_)124_storage", TupleRawName:"Map3NodeWrapperMap_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000126900), (*abi.Type)(0xc000126b80)}, TupleRawNames:[]string{"Keys", "Map"}, TupleType:(*reflect.rtype)(0xc0003c8360)}, Tag:"json:\"Nodes\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Map3NodeWrapperMap_{
		obj:           &s.obj.Nodes,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 64,
		dirty:         s.dirty,
	}
}

// DelegationIndexMapByDelegator is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodePool_) DelegationIndexMapByDelegator() *Storage_Map_972804fe1df3 {
	// Field: &wrapper.tmplField{Type:"Map_972804fe1df3", ElemType:"Map_972804fe1df3", Name:"DelegationIndexMapByDelegator", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(DelegationIndexMap_)142_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000d2a00), (*abi.Type)(0xc0000d3380)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"DelegationIndexMapByDelegator\" storage:\"slot=2,offset=0\"", Slot:2, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("2", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.DelegationIndexMapByDelegator == nil {

		instance := make(Map_972804fe1df3)

		s.obj.DelegationIndexMapByDelegator = instance
	}

	return &Storage_Map_972804fe1df3{
		obj:           s.obj.DelegationIndexMapByDelegator,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// NodeKeySet is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodePool_) NodeKeySet() *Storage_Map_624bcc0df495 {
	// Field: &wrapper.tmplField{Type:"Map_624bcc0df495", ElemType:"Map_624bcc0df495", Name:"NodeKeySet", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_string_memory_ptr,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000d3500), (*abi.Type)(0xc0000d3680)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"NodeKeySet\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
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
func (s *Storage_Map3NodePool_) DescriptionIdentitySet() *Storage_Map_624bcc0df495 {
	// Field: &wrapper.tmplField{Type:"Map_624bcc0df495", ElemType:"Map_624bcc0df495", Name:"DescriptionIdentitySet", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_string_memory_ptr,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000d3500), (*abi.Type)(0xc0000d3680)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"DescriptionIdentitySet\" storage:\"slot=4,offset=0\"", Slot:4, Offset:0}
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

// Entry is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeWrapperMapEntry_) Entry() *Storage_Map3NodeWrapper_ {
	// Field: &wrapper.tmplField{Type:"Map3NodeWrapper_", ElemType:"Map3NodeWrapper_", Name:"Entry", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x2e0, stringKind:"t_struct(Map3NodeWrapper_)111_storage", TupleRawName:"Map3NodeWrapper_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000126d80), (*abi.Type)(0xc000550200), (*abi.Type)(0xc0000d2300), (*abi.Type)(0xc0000d2480), (*abi.Type)(0xc0000d2500), (*abi.Type)(0xc0000d2600)}, TupleRawNames:[]string{"Map3Node", "Microdelegations", "RestakingReference", "AccumulatedReward", "TotalDelegation", "TotalPendingDelegation"}, TupleType:(*reflect.rtype)(0xc0003c8240)}, Tag:"json:\"Entry\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Map3NodeWrapper_{
		obj:           &s.obj.Entry,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 736,
		dirty:         s.dirty,
	}
}

// Index is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeWrapperMapEntry_) Index() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Index", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Index\" storage:\"slot=23,offset=0\"", Slot:23, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("23", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Index == nil {

		instance := big.NewInt(0)

		s.obj.Index = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.Index,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Keys is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeWrapperMap_) Keys() *Storage_Slice_760032d9d704 {
	// Field: &wrapper.tmplField{Type:"Slice_760032d9d704", ElemType:"Slice_760032d9d704", Name:"Keys", SolKind:abi.Type{Elem:(*abi.Type)(0xc000126a80), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_address)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Keys\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Keys == nil {

		hash := s.db.GetState(s.addr, common.BigToHash(actual))
		instance := make(Slice_760032d9d704, hash.Big().Int64())

		s.obj.Keys = instance
	}

	return &Storage_Slice_760032d9d704{
		obj:           &s.obj.Keys,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Map is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeWrapperMap_) Map() *Storage_Map_619a78bdf676 {
	// Field: &wrapper.tmplField{Type:"Map_619a78bdf676", ElemType:"Map_619a78bdf676", Name:"Map", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(Map3NodeWrapperMapEntry_)116_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc000126c00), (*abi.Type)(0xc0000d2800)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Map\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Map == nil {

		instance := make(Map_619a78bdf676)

		s.obj.Map = instance
	}

	return &Storage_Map_619a78bdf676{
		obj:           s.obj.Map,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Map3Node is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeWrapper_) Map3Node() *Storage_Map3Node_ {
	// Field: &wrapper.tmplField{Type:"Map3Node_", ElemType:"Map3Node_", Name:"Map3Node", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x220, stringKind:"t_struct(Map3Node_)54_storage", TupleRawName:"Map3Node_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000126e00), (*abi.Type)(0xc000126e80), (*abi.Type)(0xc000126f00), (*abi.Type)(0xc000127480), (*abi.Type)(0xc000127880), (*abi.Type)(0xc000127e00), (*abi.Type)(0xc000127e80), (*abi.Type)(0xc000127f00), (*abi.Type)(0xc000127f80), (*abi.Type)(0xc000550080), (*abi.Type)(0xc000550100)}, TupleRawNames:[]string{"Map3Address", "OperatorAddress", "NodeKeys", "Commission", "Description", "CreationHeight", "Age", "Status", "PendingEpoch", "ActivationEpoch", "ReleaseEpoch"}, TupleType:(*reflect.rtype)(0xc00054a840)}, Tag:"json:\"Map3Node\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Map3Node_{
		obj:           &s.obj.Map3Node,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 544,
		dirty:         s.dirty,
	}
}

// Microdelegations is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeWrapper_) Microdelegations() *Storage_MicrodelegationMap_ {
	// Field: &wrapper.tmplField{Type:"MicrodelegationMap_", ElemType:"MicrodelegationMap_", Name:"Microdelegations", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(MicrodelegationMap_)95_storage", TupleRawName:"MicrodelegationMap_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000550280), (*abi.Type)(0xc000550300)}, TupleRawNames:[]string{"Keys", "Map"}, TupleType:(*reflect.rtype)(0xc0003c80c0)}, Tag:"json:\"Microdelegations\" storage:\"slot=17,offset=0\"", Slot:17, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("17", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_MicrodelegationMap_{
		obj:           &s.obj.Microdelegations,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 64,
		dirty:         s.dirty,
	}
}

// RestakingReference is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeWrapper_) RestakingReference() *Storage_RestakingReference_ {
	// Field: &wrapper.tmplField{Type:"RestakingReference_", ElemType:"RestakingReference_", Name:"RestakingReference", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(RestakingReference_)98_storage", TupleRawName:"RestakingReference_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000d2380)}, TupleRawNames:[]string{"ValidatorAddress"}, TupleType:(*reflect.rtype)(0xc0003c81e0)}, Tag:"json:\"RestakingReference\" storage:\"slot=19,offset=0\"", Slot:19, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("19", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_RestakingReference_{
		obj:           &s.obj.RestakingReference,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// AccumulatedReward is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeWrapper_) AccumulatedReward() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"AccumulatedReward", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"AccumulatedReward\" storage:\"slot=20,offset=0\"", Slot:20, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("20", 10)
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

// TotalDelegation is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeWrapper_) TotalDelegation() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"TotalDelegation", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"TotalDelegation\" storage:\"slot=21,offset=0\"", Slot:21, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("21", 10)
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
func (s *Storage_Map3NodeWrapper_) TotalPendingDelegation() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"TotalPendingDelegation", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"TotalPendingDelegation\" storage:\"slot=22,offset=0\"", Slot:22, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("22", 10)
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

// Map3Address is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3Node_) Map3Address() *Storage_Address {
	// Field: &wrapper.tmplField{Type:"Address", ElemType:"Address", Name:"Map3Address", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Map3Address\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Address{
		obj:           &s.obj.Map3Address,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 20,
		dirty:         s.dirty,
	}
}

// OperatorAddress is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3Node_) OperatorAddress() *Storage_Address {
	// Field: &wrapper.tmplField{Type:"Address", ElemType:"Address", Name:"OperatorAddress", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"OperatorAddress\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Address{
		obj:           &s.obj.OperatorAddress,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 20,
		dirty:         s.dirty,
	}
}

// NodeKeys is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3Node_) NodeKeys() *Storage_BLSPublicKeys_ {
	// Field: &wrapper.tmplField{Type:"BLSPublicKeys_", ElemType:"BLSPublicKeys_", Name:"NodeKeys", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(BLSPublicKeys_)24_storage", TupleRawName:"BLSPublicKeys_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000126f80)}, TupleRawNames:[]string{"Keys"}, TupleType:(*reflect.rtype)(0xc00054a4e0)}, Tag:"json:\"NodeKeys\" storage:\"slot=2,offset=0\"", Slot:2, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("2", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_BLSPublicKeys_{
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
func (s *Storage_Map3Node_) Commission() *Storage_Commission_ {
	// Field: &wrapper.tmplField{Type:"Commission_", ElemType:"Commission_", Name:"Commission", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x60, stringKind:"t_struct(Commission_)31_storage", TupleRawName:"Commission_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000127500), (*abi.Type)(0xc000127700), (*abi.Type)(0xc000127780)}, TupleRawNames:[]string{"Rate", "RateForNextPeriod", "UpdateHeight"}, TupleType:(*reflect.rtype)(0xc00054a6c0)}, Tag:"json:\"Commission\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("3", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Commission_{
		obj:           &s.obj.Commission,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 96,
		dirty:         s.dirty,
	}
}

// Description is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3Node_) Description() *Storage_Description_ {
	// Field: &wrapper.tmplField{Type:"Description_", ElemType:"Description_", Name:"Description", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xa0, stringKind:"t_struct(Description_)15_storage", TupleRawName:"Description_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000127900), (*abi.Type)(0xc000127a00), (*abi.Type)(0xc000127a80), (*abi.Type)(0xc000127b00), (*abi.Type)(0xc000127b80)}, TupleRawNames:[]string{"Name", "Identity", "Website", "SecurityContact", "Details"}, TupleType:(*reflect.rtype)(0xc00054a7e0)}, Tag:"json:\"Description\" storage:\"slot=6,offset=0\"", Slot:6, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("6", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Description_{
		obj:           &s.obj.Description,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 160,
		dirty:         s.dirty,
	}
}

// CreationHeight is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3Node_) CreationHeight() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"CreationHeight", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"CreationHeight\" storage:\"slot=11,offset=0\"", Slot:11, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("11", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.CreationHeight == nil {

		instance := big.NewInt(0)

		s.obj.CreationHeight = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.CreationHeight,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Age is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3Node_) Age() *Storage_Decimal {
	// Field: &wrapper.tmplField{Type:"Decimal", ElemType:"Decimal", Name:"Age", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc000127580)}, TupleRawNames:[]string{"f"}, TupleType:(*reflect.rtype)(0xc00054a660)}, Tag:"json:\"Age\" storage:\"slot=12,offset=0\"", Slot:12, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("12", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Decimal{
		obj:           &s.obj.Age,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Status is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3Node_) Status() *Storage_Uint8 {
	// Field: &wrapper.tmplField{Type:"Uint8", ElemType:"Uint8", Name:"Status", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:8, T:0x1, numberOfBytes:0x1, stringKind:"uint8", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Status\" storage:\"slot=13,offset=0\"", Slot:13, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("13", 10)
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

// PendingEpoch is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3Node_) PendingEpoch() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"PendingEpoch", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"PendingEpoch\" storage:\"slot=14,offset=0\"", Slot:14, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("14", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.PendingEpoch == nil {

		instance := big.NewInt(0)

		s.obj.PendingEpoch = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.PendingEpoch,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// ActivationEpoch is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3Node_) ActivationEpoch() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"ActivationEpoch", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"ActivationEpoch\" storage:\"slot=15,offset=0\"", Slot:15, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("15", 10)
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
func (s *Storage_Map3Node_) ReleaseEpoch() *Storage_Decimal {
	// Field: &wrapper.tmplField{Type:"Decimal", ElemType:"Decimal", Name:"ReleaseEpoch", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc000127580)}, TupleRawNames:[]string{"f"}, TupleType:(*reflect.rtype)(0xc00054a660)}, Tag:"json:\"ReleaseEpoch\" storage:\"slot=16,offset=0\"", Slot:16, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("16", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Decimal{
		obj:           &s.obj.ReleaseEpoch,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Entry is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_MicrodelegationMapEntry_) Entry() *Storage_Microdelegation_ {
	// Field: &wrapper.tmplField{Type:"Microdelegation_", ElemType:"Microdelegation_", Name:"Entry", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x120, stringKind:"t_struct(Microdelegation_)77_storage", TupleRawName:"Microdelegation_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000550500), (*abi.Type)(0xc000550580), (*abi.Type)(0xc000550600), (*abi.Type)(0xc000550680), (*abi.Type)(0xc000550880), (*abi.Type)(0xc000550b00)}, TupleRawNames:[]string{"DelegatorAddress", "Amount", "Reward", "PendingDelegation", "Undelegation", "Renewal"}, TupleType:(*reflect.rtype)(0xc000704000)}, Tag:"json:\"Entry\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Microdelegation_{
		obj:           &s.obj.Entry,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 288,
		dirty:         s.dirty,
	}
}

// Index is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_MicrodelegationMapEntry_) Index() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Index", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Index\" storage:\"slot=9,offset=0\"", Slot:9, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("9", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Index == nil {

		instance := big.NewInt(0)

		s.obj.Index = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.Index,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Keys is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_MicrodelegationMap_) Keys() *Storage_Slice_760032d9d704 {
	// Field: &wrapper.tmplField{Type:"Slice_760032d9d704", ElemType:"Slice_760032d9d704", Name:"Keys", SolKind:abi.Type{Elem:(*abi.Type)(0xc000126a80), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_address)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Keys\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Keys == nil {

		hash := s.db.GetState(s.addr, common.BigToHash(actual))
		instance := make(Slice_760032d9d704, hash.Big().Int64())

		s.obj.Keys = instance
	}

	return &Storage_Slice_760032d9d704{
		obj:           &s.obj.Keys,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Map is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_MicrodelegationMap_) Map() *Storage_Map_2f2f5e901078 {
	// Field: &wrapper.tmplField{Type:"Map_2f2f5e901078", ElemType:"Map_2f2f5e901078", Name:"Map", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(MicrodelegationMapEntry_)87_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc000550380), (*abi.Type)(0xc0000d2180)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Map\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Map == nil {

		instance := make(Map_2f2f5e901078)

		s.obj.Map = instance
	}

	return &Storage_Map_2f2f5e901078{
		obj:           s.obj.Map,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// DelegatorAddress is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Microdelegation_) DelegatorAddress() *Storage_Address {
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
func (s *Storage_Microdelegation_) Amount() *Storage_BigInt {
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
func (s *Storage_Microdelegation_) Reward() *Storage_BigInt {
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

// PendingDelegation is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Microdelegation_) PendingDelegation() *Storage_PendingDelegation_ {
	// Field: &wrapper.tmplField{Type:"PendingDelegation_", ElemType:"PendingDelegation_", Name:"PendingDelegation", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(PendingDelegation_)59_storage", TupleRawName:"PendingDelegation_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000550700), (*abi.Type)(0xc000550780)}, TupleRawNames:[]string{"Amount", "UnlockedEpoch"}, TupleType:(*reflect.rtype)(0xc00054aa80)}, Tag:"json:\"PendingDelegation\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("3", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_PendingDelegation_{
		obj:           &s.obj.PendingDelegation,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 64,
		dirty:         s.dirty,
	}
}

// Undelegation is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Microdelegation_) Undelegation() *Storage_Undelegation_ {
	// Field: &wrapper.tmplField{Type:"Undelegation_", ElemType:"Undelegation_", Name:"Undelegation", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Undelegation_)64_storage", TupleRawName:"Undelegation_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000550900), (*abi.Type)(0xc000550980)}, TupleRawNames:[]string{"Amount", "Epoch"}, TupleType:(*reflect.rtype)(0xc00054ab40)}, Tag:"json:\"Undelegation\" storage:\"slot=5,offset=0\"", Slot:5, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("5", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Undelegation_{
		obj:           &s.obj.Undelegation,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 64,
		dirty:         s.dirty,
	}
}

// Renewal is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Microdelegation_) Renewal() *Storage_Renewal_ {
	// Field: &wrapper.tmplField{Type:"Renewal_", ElemType:"Renewal_", Name:"Renewal", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Renewal_)82_storage", TupleRawName:"Renewal_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000550b80), (*abi.Type)(0xc000550c00)}, TupleRawNames:[]string{"Status", "UpdateHeight"}, TupleType:(*reflect.rtype)(0xc00054ac00)}, Tag:"json:\"Renewal\" storage:\"slot=7,offset=0\"", Slot:7, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("7", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Renewal_{
		obj:           &s.obj.Renewal,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 64,
		dirty:         s.dirty,
	}
}

// Amount is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_PendingDelegation_) Amount() *Storage_BigInt {
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

// UnlockedEpoch is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_PendingDelegation_) UnlockedEpoch() *Storage_Decimal {
	// Field: &wrapper.tmplField{Type:"Decimal", ElemType:"Decimal", Name:"UnlockedEpoch", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc000127580)}, TupleRawNames:[]string{"f"}, TupleType:(*reflect.rtype)(0xc00054a660)}, Tag:"json:\"UnlockedEpoch\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Decimal{
		obj:           &s.obj.UnlockedEpoch,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Status is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Renewal_) Status() *Storage_Uint8 {
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

// UpdateHeight is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Renewal_) UpdateHeight() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"UpdateHeight", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"UpdateHeight\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
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

// ValidatorAddress is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_RestakingReference_) ValidatorAddress() *Storage_Address {
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

// Amount is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Undelegation_) Amount() *Storage_BigInt {
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
func (s *Storage_Undelegation_) Epoch() *Storage_BigInt {
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

// -------------------------------- enhance blew -------------------------------------

func (s *Storage_Address) Clear() {
	// &wrapper.tmplStruct{T:0x7, Name:"Address", Fields:[]*wrapper.tmplField(nil), Type:"common.Address", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	s.SetValue(common.Address{})
}

func (s *Storage_BigInt) Clear() {
	// &wrapper.tmplStruct{T:0x1, Name:"BigInt", Fields:[]*wrapper.tmplField(nil), Type:"*big.Int", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	s.SetValue(big.NewInt(0))
}

func (s *Storage_Bool) Clear() {
	// &wrapper.tmplStruct{T:0x2, Name:"Bool", Fields:[]*wrapper.tmplField(nil), Type:"bool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x2, numberOfBytes:0x1, stringKind:"t_bool", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	s.SetValue(false)
}

func (s *Storage_Decimal) Clear() {
	// &wrapper.tmplStruct{T:0xf, Name:"Decimal", Fields:[]*wrapper.tmplField(nil), Type:"common.Dec", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xf, numberOfBytes:0x20, stringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc000127580)}, TupleRawNames:[]string{"f"}, TupleType:(*reflect.rtype)(0xc00054a660)}}
	s.SetValue(common.NewDecFromBigIntWithPrec(big.NewInt(0), common.Precision))
}

func (s *Storage_String) Clear() {
	// &wrapper.tmplStruct{T:0x3, Name:"String", Fields:[]*wrapper.tmplField(nil), Type:"string", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x3, numberOfBytes:0x20, stringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	rv := s.Value()
	s.SetValue(string(bytes.Repeat([]byte{0x00}, len(rv))))
	s.SetValue("")
}

func (s *Storage_Uint8) Clear() {
	// &wrapper.tmplStruct{T:0x1, Name:"Uint8", Fields:[]*wrapper.tmplField(nil), Type:"uint8", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:8, T:0x1, numberOfBytes:0x1, stringKind:"uint8", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	s.SetValue(0)
}

func (s *Storage_BLSPublicKey_) Save(obj *BLSPublicKey_) {
	s.Key().Save(obj.Key)
}

func (s *Storage_BLSPublicKey_) Clear() {
	s.Key().Clear()
}

func (s *Storage_BLSPublicKey_) load() *BLSPublicKey_ {
	s.Key().load()
	return s.obj
}

func (s *Storage_BLSPublicKey_) Load() (*BLSPublicKey_, error) {
	src := s.load()
	des := BLSPublicKey_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *BLSPublicKey_) Copy() (*BLSPublicKey_, error) {
	des := BLSPublicKey_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_BLSPublicKeys_) Save(obj *BLSPublicKeys_) {
	if obj.Keys != nil {
		s.Keys().Save(obj.Keys)
	}
}

func (s *Storage_BLSPublicKeys_) Clear() {
	s.Keys().Clear()
}

func (s *Storage_BLSPublicKeys_) load() *BLSPublicKeys_ {
	s.Keys().load()
	return s.obj
}

func (s *Storage_BLSPublicKeys_) Load() (*BLSPublicKeys_, error) {
	src := s.load()
	des := BLSPublicKeys_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *BLSPublicKeys_) Copy() (*BLSPublicKeys_, error) {
	des := BLSPublicKeys_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_Commission_) Save(obj *Commission_) {
	s.Rate().SetValue(obj.Rate)
	s.RateForNextPeriod().SetValue(obj.RateForNextPeriod)
	if obj.UpdateHeight != nil {
		s.UpdateHeight().SetValue(obj.UpdateHeight)
	}
}

func (s *Storage_Commission_) Clear() {
	s.Rate().Clear()
	s.RateForNextPeriod().Clear()
	s.UpdateHeight().Clear()
}

func (s *Storage_Commission_) load() *Commission_ {
	s.Rate().Value()
	s.RateForNextPeriod().Value()
	s.UpdateHeight().Value()
	return s.obj
}

func (s *Storage_Commission_) Load() (*Commission_, error) {
	src := s.load()
	des := Commission_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Commission_) Copy() (*Commission_, error) {
	des := Commission_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_DelegationIndexMapEntry_) Save(obj *DelegationIndexMapEntry_) {
	s.Entry().Save(&obj.Entry)
	if obj.Index != nil {
		s.Index().SetValue(obj.Index)
	}
}

func (s *Storage_DelegationIndexMapEntry_) Clear() {
	s.Entry().Clear()
	s.Index().Clear()
}

func (s *Storage_DelegationIndexMapEntry_) load() *DelegationIndexMapEntry_ {
	s.Entry().load()
	s.Index().Value()
	return s.obj
}

func (s *Storage_DelegationIndexMapEntry_) Load() (*DelegationIndexMapEntry_, error) {
	src := s.load()
	des := DelegationIndexMapEntry_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *DelegationIndexMapEntry_) Copy() (*DelegationIndexMapEntry_, error) {
	des := DelegationIndexMapEntry_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_DelegationIndexMap_) Save(obj *DelegationIndexMap_) {
	if obj.Keys != nil {
		s.Keys().Save(obj.Keys)
	}
	if obj.Map != nil {
		s.Map().Save(obj.Map)
	}
}

func (s *Storage_DelegationIndexMap_) Clear() {
	s.Keys().Clear()
	s.Map().Clear()
}

func (s *Storage_DelegationIndexMap_) load() *DelegationIndexMap_ {
	s.Keys().load()
	s.Map().load()
	return s.obj
}

func (s *Storage_DelegationIndexMap_) Load() (*DelegationIndexMap_, error) {
	src := s.load()
	des := DelegationIndexMap_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *DelegationIndexMap_) Copy() (*DelegationIndexMap_, error) {
	des := DelegationIndexMap_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_DelegationIndex_) Save(obj *DelegationIndex_) {
	s.Map3Address().SetValue(obj.Map3Address)
	s.IsOperator().SetValue(obj.IsOperator)
}

func (s *Storage_DelegationIndex_) Clear() {
	s.Map3Address().Clear()
	s.IsOperator().Clear()
}

func (s *Storage_DelegationIndex_) load() *DelegationIndex_ {
	s.Map3Address().Value()
	s.IsOperator().Value()
	return s.obj
}

func (s *Storage_DelegationIndex_) Load() (*DelegationIndex_, error) {
	src := s.load()
	des := DelegationIndex_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *DelegationIndex_) Copy() (*DelegationIndex_, error) {
	des := DelegationIndex_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_Description_) Save(obj *Description_) {
	s.Name().SetValue(obj.Name)
	s.Identity().SetValue(obj.Identity)
	s.Website().SetValue(obj.Website)
	s.SecurityContact().SetValue(obj.SecurityContact)
	s.Details().SetValue(obj.Details)
}

func (s *Storage_Description_) Clear() {
	s.Name().Clear()
	s.Identity().Clear()
	s.Website().Clear()
	s.SecurityContact().Clear()
	s.Details().Clear()
}

func (s *Storage_Description_) load() *Description_ {
	s.Name().Value()
	s.Identity().Value()
	s.Website().Value()
	s.SecurityContact().Value()
	s.Details().Value()
	return s.obj
}

func (s *Storage_Description_) Load() (*Description_, error) {
	src := s.load()
	des := Description_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Description_) Copy() (*Description_, error) {
	des := Description_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_Global_t) Save(obj *Global_t) {
	s.Map3NodePool().Save(&obj.Map3NodePool)
}

func (s *Storage_Global_t) Clear() {
	s.Map3NodePool().Clear()
}

func (s *Storage_Global_t) load() *Global_t {
	s.Map3NodePool().load()
	return s.obj
}

func (s *Storage_Global_t) Load() (*Global_t, error) {
	src := s.load()
	des := Global_t{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Global_t) Copy() (*Global_t, error) {
	des := Global_t{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_Map3NodePool_) Save(obj *Map3NodePool_) {
	s.Nodes().Save(&obj.Nodes)
	if obj.DelegationIndexMapByDelegator != nil {
		s.DelegationIndexMapByDelegator().Save(obj.DelegationIndexMapByDelegator)
	}
	if obj.NodeKeySet != nil {
		s.NodeKeySet().Save(obj.NodeKeySet)
	}
	if obj.DescriptionIdentitySet != nil {
		s.DescriptionIdentitySet().Save(obj.DescriptionIdentitySet)
	}
}

func (s *Storage_Map3NodePool_) Clear() {
	s.Nodes().Clear()
	s.DelegationIndexMapByDelegator().Clear()
	s.NodeKeySet().Clear()
	s.DescriptionIdentitySet().Clear()
}

func (s *Storage_Map3NodePool_) load() *Map3NodePool_ {
	s.Nodes().load()
	s.DelegationIndexMapByDelegator().load()
	s.NodeKeySet().load()
	s.DescriptionIdentitySet().load()
	return s.obj
}

func (s *Storage_Map3NodePool_) Load() (*Map3NodePool_, error) {
	src := s.load()
	des := Map3NodePool_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Map3NodePool_) Copy() (*Map3NodePool_, error) {
	des := Map3NodePool_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_Map3NodeWrapperMapEntry_) Save(obj *Map3NodeWrapperMapEntry_) {
	s.Entry().Save(&obj.Entry)
	if obj.Index != nil {
		s.Index().SetValue(obj.Index)
	}
}

func (s *Storage_Map3NodeWrapperMapEntry_) Clear() {
	s.Entry().Clear()
	s.Index().Clear()
}

func (s *Storage_Map3NodeWrapperMapEntry_) load() *Map3NodeWrapperMapEntry_ {
	s.Entry().load()
	s.Index().Value()
	return s.obj
}

func (s *Storage_Map3NodeWrapperMapEntry_) Load() (*Map3NodeWrapperMapEntry_, error) {
	src := s.load()
	des := Map3NodeWrapperMapEntry_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Map3NodeWrapperMapEntry_) Copy() (*Map3NodeWrapperMapEntry_, error) {
	des := Map3NodeWrapperMapEntry_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_Map3NodeWrapperMap_) Save(obj *Map3NodeWrapperMap_) {
	if obj.Keys != nil {
		s.Keys().Save(obj.Keys)
	}
	if obj.Map != nil {
		s.Map().Save(obj.Map)
	}
}

func (s *Storage_Map3NodeWrapperMap_) Clear() {
	s.Keys().Clear()
	s.Map().Clear()
}

func (s *Storage_Map3NodeWrapperMap_) load() *Map3NodeWrapperMap_ {
	s.Keys().load()
	s.Map().load()
	return s.obj
}

func (s *Storage_Map3NodeWrapperMap_) Load() (*Map3NodeWrapperMap_, error) {
	src := s.load()
	des := Map3NodeWrapperMap_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Map3NodeWrapperMap_) Copy() (*Map3NodeWrapperMap_, error) {
	des := Map3NodeWrapperMap_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_Map3NodeWrapper_) Save(obj *Map3NodeWrapper_) {
	s.Map3Node().Save(&obj.Map3Node)
	s.Microdelegations().Save(&obj.Microdelegations)
	s.RestakingReference().Save(&obj.RestakingReference)
	if obj.AccumulatedReward != nil {
		s.AccumulatedReward().SetValue(obj.AccumulatedReward)
	}
	if obj.TotalDelegation != nil {
		s.TotalDelegation().SetValue(obj.TotalDelegation)
	}
	if obj.TotalPendingDelegation != nil {
		s.TotalPendingDelegation().SetValue(obj.TotalPendingDelegation)
	}
}

func (s *Storage_Map3NodeWrapper_) Clear() {
	s.Map3Node().Clear()
	s.Microdelegations().Clear()
	s.RestakingReference().Clear()
	s.AccumulatedReward().Clear()
	s.TotalDelegation().Clear()
	s.TotalPendingDelegation().Clear()
}

func (s *Storage_Map3NodeWrapper_) load() *Map3NodeWrapper_ {
	s.Map3Node().load()
	s.Microdelegations().load()
	s.RestakingReference().load()
	s.AccumulatedReward().Value()
	s.TotalDelegation().Value()
	s.TotalPendingDelegation().Value()
	return s.obj
}

func (s *Storage_Map3NodeWrapper_) Load() (*Map3NodeWrapper_, error) {
	src := s.load()
	des := Map3NodeWrapper_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Map3NodeWrapper_) Copy() (*Map3NodeWrapper_, error) {
	des := Map3NodeWrapper_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_Map3Node_) Save(obj *Map3Node_) {
	s.Map3Address().SetValue(obj.Map3Address)
	s.OperatorAddress().SetValue(obj.OperatorAddress)
	s.NodeKeys().Save(&obj.NodeKeys)
	s.Commission().Save(&obj.Commission)
	s.Description().Save(&obj.Description)
	if obj.CreationHeight != nil {
		s.CreationHeight().SetValue(obj.CreationHeight)
	}
	s.Age().SetValue(obj.Age)
	s.Status().SetValue(obj.Status)
	if obj.PendingEpoch != nil {
		s.PendingEpoch().SetValue(obj.PendingEpoch)
	}
	if obj.ActivationEpoch != nil {
		s.ActivationEpoch().SetValue(obj.ActivationEpoch)
	}
	s.ReleaseEpoch().SetValue(obj.ReleaseEpoch)
}

func (s *Storage_Map3Node_) Clear() {
	s.Map3Address().Clear()
	s.OperatorAddress().Clear()
	s.NodeKeys().Clear()
	s.Commission().Clear()
	s.Description().Clear()
	s.CreationHeight().Clear()
	s.Age().Clear()
	s.Status().Clear()
	s.PendingEpoch().Clear()
	s.ActivationEpoch().Clear()
	s.ReleaseEpoch().Clear()
}

func (s *Storage_Map3Node_) load() *Map3Node_ {
	s.Map3Address().Value()
	s.OperatorAddress().Value()
	s.NodeKeys().load()
	s.Commission().load()
	s.Description().load()
	s.CreationHeight().Value()
	s.Age().Value()
	s.Status().Value()
	s.PendingEpoch().Value()
	s.ActivationEpoch().Value()
	s.ReleaseEpoch().Value()
	return s.obj
}

func (s *Storage_Map3Node_) Load() (*Map3Node_, error) {
	src := s.load()
	des := Map3Node_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Map3Node_) Copy() (*Map3Node_, error) {
	des := Map3Node_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_MicrodelegationMapEntry_) Save(obj *MicrodelegationMapEntry_) {
	s.Entry().Save(&obj.Entry)
	if obj.Index != nil {
		s.Index().SetValue(obj.Index)
	}
}

func (s *Storage_MicrodelegationMapEntry_) Clear() {
	s.Entry().Clear()
	s.Index().Clear()
}

func (s *Storage_MicrodelegationMapEntry_) load() *MicrodelegationMapEntry_ {
	s.Entry().load()
	s.Index().Value()
	return s.obj
}

func (s *Storage_MicrodelegationMapEntry_) Load() (*MicrodelegationMapEntry_, error) {
	src := s.load()
	des := MicrodelegationMapEntry_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *MicrodelegationMapEntry_) Copy() (*MicrodelegationMapEntry_, error) {
	des := MicrodelegationMapEntry_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_MicrodelegationMap_) Save(obj *MicrodelegationMap_) {
	if obj.Keys != nil {
		s.Keys().Save(obj.Keys)
	}
	if obj.Map != nil {
		s.Map().Save(obj.Map)
	}
}

func (s *Storage_MicrodelegationMap_) Clear() {
	s.Keys().Clear()
	s.Map().Clear()
}

func (s *Storage_MicrodelegationMap_) load() *MicrodelegationMap_ {
	s.Keys().load()
	s.Map().load()
	return s.obj
}

func (s *Storage_MicrodelegationMap_) Load() (*MicrodelegationMap_, error) {
	src := s.load()
	des := MicrodelegationMap_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *MicrodelegationMap_) Copy() (*MicrodelegationMap_, error) {
	des := MicrodelegationMap_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_Microdelegation_) Save(obj *Microdelegation_) {
	s.DelegatorAddress().SetValue(obj.DelegatorAddress)
	if obj.Amount != nil {
		s.Amount().SetValue(obj.Amount)
	}
	if obj.Reward != nil {
		s.Reward().SetValue(obj.Reward)
	}
	s.PendingDelegation().Save(&obj.PendingDelegation)
	s.Undelegation().Save(&obj.Undelegation)
	s.Renewal().Save(&obj.Renewal)
}

func (s *Storage_Microdelegation_) Clear() {
	s.DelegatorAddress().Clear()
	s.Amount().Clear()
	s.Reward().Clear()
	s.PendingDelegation().Clear()
	s.Undelegation().Clear()
	s.Renewal().Clear()
}

func (s *Storage_Microdelegation_) load() *Microdelegation_ {
	s.DelegatorAddress().Value()
	s.Amount().Value()
	s.Reward().Value()
	s.PendingDelegation().load()
	s.Undelegation().load()
	s.Renewal().load()
	return s.obj
}

func (s *Storage_Microdelegation_) Load() (*Microdelegation_, error) {
	src := s.load()
	des := Microdelegation_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Microdelegation_) Copy() (*Microdelegation_, error) {
	des := Microdelegation_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_PendingDelegation_) Save(obj *PendingDelegation_) {
	if obj.Amount != nil {
		s.Amount().SetValue(obj.Amount)
	}
	s.UnlockedEpoch().SetValue(obj.UnlockedEpoch)
}

func (s *Storage_PendingDelegation_) Clear() {
	s.Amount().Clear()
	s.UnlockedEpoch().Clear()
}

func (s *Storage_PendingDelegation_) load() *PendingDelegation_ {
	s.Amount().Value()
	s.UnlockedEpoch().Value()
	return s.obj
}

func (s *Storage_PendingDelegation_) Load() (*PendingDelegation_, error) {
	src := s.load()
	des := PendingDelegation_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *PendingDelegation_) Copy() (*PendingDelegation_, error) {
	des := PendingDelegation_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_Renewal_) Save(obj *Renewal_) {
	s.Status().SetValue(obj.Status)
	if obj.UpdateHeight != nil {
		s.UpdateHeight().SetValue(obj.UpdateHeight)
	}
}

func (s *Storage_Renewal_) Clear() {
	s.Status().Clear()
	s.UpdateHeight().Clear()
}

func (s *Storage_Renewal_) load() *Renewal_ {
	s.Status().Value()
	s.UpdateHeight().Value()
	return s.obj
}

func (s *Storage_Renewal_) Load() (*Renewal_, error) {
	src := s.load()
	des := Renewal_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Renewal_) Copy() (*Renewal_, error) {
	des := Renewal_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_RestakingReference_) Save(obj *RestakingReference_) {
	s.ValidatorAddress().SetValue(obj.ValidatorAddress)
}

func (s *Storage_RestakingReference_) Clear() {
	s.ValidatorAddress().Clear()
}

func (s *Storage_RestakingReference_) load() *RestakingReference_ {
	s.ValidatorAddress().Value()
	return s.obj
}

func (s *Storage_RestakingReference_) Load() (*RestakingReference_, error) {
	src := s.load()
	des := RestakingReference_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *RestakingReference_) Copy() (*RestakingReference_, error) {
	des := RestakingReference_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_Undelegation_) Save(obj *Undelegation_) {
	if obj.Amount != nil {
		s.Amount().SetValue(obj.Amount)
	}
	if obj.Epoch != nil {
		s.Epoch().SetValue(obj.Epoch)
	}
}

func (s *Storage_Undelegation_) Clear() {
	s.Amount().Clear()
	s.Epoch().Clear()
}

func (s *Storage_Undelegation_) load() *Undelegation_ {
	s.Amount().Value()
	s.Epoch().Value()
	return s.obj
}

func (s *Storage_Undelegation_) Load() (*Undelegation_, error) {
	src := s.load()
	des := Undelegation_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Undelegation_) Copy() (*Undelegation_, error) {
	des := Undelegation_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_Array_aaa956d410fd) Save(obj Array_aaa956d410fd) {
	s.SetValue(obj)
}

func (s *Storage_Array_aaa956d410fd) Clear() {
	s.SetValue(Array_aaa956d410fd{})
}

func (s *Storage_Array_aaa956d410fd) load() Array_aaa956d410fd {
	s.Value()
	return *s.obj
}

func (s *Storage_Map_2f2f5e901078) Save(obj Map_2f2f5e901078) {
	for k, v := range obj {
		s.Get(k).Save(v)
	}
}

func (s *Storage_Map_2f2f5e901078) Clear() {
	panic("not support to clear map")
}

func (s *Storage_Map_2f2f5e901078) load() Map_2f2f5e901078 {
	panic("not support to clear map")
}

func (s *Storage_Map_619a78bdf676) Save(obj Map_619a78bdf676) {
	for k, v := range obj {
		s.Get(k).Save(v)
	}
}

func (s *Storage_Map_619a78bdf676) Clear() {
	panic("not support to clear map")
}

func (s *Storage_Map_619a78bdf676) load() Map_619a78bdf676 {
	panic("not support to clear map")
}

func (s *Storage_Map_624bcc0df495) Save(obj Map_624bcc0df495) {
	for k, v := range obj {
		s.Get(k).SetValue(*v)
	}
}

func (s *Storage_Map_624bcc0df495) Clear() {
	panic("not support to clear map")
}

func (s *Storage_Map_624bcc0df495) load() Map_624bcc0df495 {
	panic("not support to clear map")
}

func (s *Storage_Map_972804fe1df3) Save(obj Map_972804fe1df3) {
	for k, v := range obj {
		s.Get(k).Save(v)
	}
}

func (s *Storage_Map_972804fe1df3) Clear() {
	panic("not support to clear map")
}

func (s *Storage_Map_972804fe1df3) load() Map_972804fe1df3 {
	panic("not support to clear map")
}

func (s *Storage_Map_f414efdd5b9c) Save(obj Map_f414efdd5b9c) {
	for k, v := range obj {
		s.Get(k).Save(v)
	}
}

func (s *Storage_Map_f414efdd5b9c) Clear() {
	panic("not support to clear map")
}

func (s *Storage_Map_f414efdd5b9c) load() Map_f414efdd5b9c {
	panic("not support to clear map")
}

func (s *Storage_Slice_2f1765e04dad) Save(obj Slice_2f1765e04dad) {
	for i := 0; i < len(obj); i++ {
		s.Get(i).Save(obj[i])
	}
	s.Resize(len(obj))
}

func (s *Storage_Slice_2f1765e04dad) Clear() {
	for i := 0; i < s.Length(); i++ {
		s.Get(i).Clear()
	}
	s.Resize(0)
}

func (s *Storage_Slice_2f1765e04dad) load() Slice_2f1765e04dad {
	for i := 0; i < s.Length(); i++ {
		s.Get(i).load()
	}
	return *s.obj
}

func (s *Storage_Slice_760032d9d704) Save(obj Slice_760032d9d704) {
	for i := 0; i < len(obj); i++ {
		s.Get(i).SetValue(*obj[i])
	}
	s.Resize(len(obj))
}

func (s *Storage_Slice_760032d9d704) Clear() {
	for i := 0; i < s.Length(); i++ {
		s.Get(i).Clear()
	}
	s.Resize(0)
}

func (s *Storage_Slice_760032d9d704) load() Slice_760032d9d704 {
	for i := 0; i < s.Length(); i++ {
		s.Get(i).Value()
	}
	return *s.obj
}

func deepCopy(src, des interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	if err := json.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(des); err != nil {
		return err
	}
	return nil
}

