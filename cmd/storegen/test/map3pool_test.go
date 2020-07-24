// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package test

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
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
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// use backtick in text/template: https://github.com/golang/go/issues/18221

type StateValues map[common.Hash]common.Hash

func GetStateAsBytes(db *state.StateDB, addr common.Address, slot *big.Int) []byte {
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
	obj *common.Address

	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func (s *Storage_Address) Value() common.Address {
	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	*s.obj = common.BigToAddress(hash.Big())
	return *s.obj
}

func (s *Storage_Address) SetValue(value common.Address) {
	hash := value.Hash()
	s.db.SetState(s.addr, common.BigToHash(s.slot), hash)
	*s.obj = value
}

type BigInt = *big.Int
type Storage_BigInt struct {
	obj *big.Int

	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func (s *Storage_BigInt) Value() *big.Int {
	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	*s.obj = *hash.Big()
	return s.obj
}

func (s *Storage_BigInt) SetValue(value *big.Int) {
	hash := value
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(hash))
	*s.obj = *value
}

type Bool = bool
type Storage_Bool struct {
	obj *bool

	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func (s *Storage_Bool) Value() bool {
	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	*s.obj = (hash.Big().Int64() != 0)
	return *s.obj
}

func (s *Storage_Bool) SetValue(value bool) {
	var val uint
	if value {
		val = 1
	} else {
		val = 0
	}
	hash := big.NewInt(0).SetUint64(uint64(val))
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(hash))
	*s.obj = value
}

type Bytes = []byte
type Storage_Bytes struct {
	obj []byte

	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func (s *Storage_Bytes) Value() []byte {
	rv := GetStateAsBytes(s.db, s.addr, s.slot)
	s.obj = []byte(rv)
	return s.obj
}

func (s *Storage_Bytes) SetValue(value []byte) {
	SetStateAsBytes(s.db, s.addr, s.slot, []byte(value))
	s.obj = value
}

type String = string
type Storage_String struct {
	obj *string

	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func (s *Storage_String) Value() string {
	rv := GetStateAsBytes(s.db, s.addr, s.slot)
	*s.obj = string(rv)
	return *s.obj
}

func (s *Storage_String) SetValue(value string) {
	SetStateAsBytes(s.db, s.addr, s.slot, []byte(value))
	*s.obj = value
}

type Uint64 = uint64
type Storage_Uint64 struct {
	obj *uint64

	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func (s *Storage_Uint64) Value() uint64 {
	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	*s.obj = uint64(hash.Big().Uint64())
	return *s.obj
}

func (s *Storage_Uint64) SetValue(value uint64) {
	hash := big.NewInt(0).SetUint64(uint64(value))
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(hash))
	*s.obj = value
}

type Uint8 = uint8
type Storage_Uint8 struct {
	obj *uint8

	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func (s *Storage_Uint8) Value() uint8 {
	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	*s.obj = uint8(hash.Big().Uint64())
	return *s.obj
}

func (s *Storage_Uint8) SetValue(value uint8) {
	hash := big.NewInt(0).SetUint64(uint64(value))
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(hash))
	*s.obj = value
}

// Array_15880d9287c7 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x5, Name:"Array_15880d9287c7", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0004ddea0)}, Type:"[2]PendingDelegation_t", SolKind:abi.Type{Elem:(*abi.Type)(0xc00062e750), Kind:0x11, Type:(*reflect.rtype)(0xc00048aaf0), Size:2, T:0x5, NumberOfBytes:0x80, StringKind:"t_array(t_struct(PendingDelegation_t)52_storage)2_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Array_15880d9287c7 [2]PendingDelegation_t

type Storage_Array_15880d9287c7 struct {
	obj   *Array_15880d9287c7
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func (s *Storage_Array_15880d9287c7) Length() int {
	return len(s.obj)
}

func (s *Storage_Array_15880d9287c7) Get(index uint64) *Storage_PendingDelegation_t {
	// Value: &wrapper.tmplField{Type:"PendingDelegation_t", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000243860), Size:0, T:0x6, NumberOfBytes:0x40, StringKind:"t_struct(PendingDelegation_t)52_storage", TupleRawName:"PendingDelegation_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00062e7e0), (*abi.Type)(0xc00062e870)}, TupleRawNames:[]string{"Amount", "Epoch"}}, Tag:"", Slot:<nil>, Offset:0}
	actual := big.NewInt(0).Add(s.slot, big.NewInt(0).SetUint64(index*(64/32)))

	return &Storage_PendingDelegation_t{
		obj:   &s.obj[index],
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Array_19387431c506 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x5, Name:"Array_19387431c506", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0004dcfc0)}, Type:"[2]BigInt", SolKind:abi.Type{Elem:(*abi.Type)(0xc00060bd40), Kind:0x11, Type:(*reflect.rtype)(0xc00048aaf0), Size:2, T:0x5, NumberOfBytes:0x40, StringKind:"t_array(t_uint256)2_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Array_19387431c506 [2]BigInt

type Storage_Array_19387431c506 struct {
	obj   *Array_19387431c506
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func (s *Storage_Array_19387431c506) Length() int {
	return len(s.obj)
}

func (s *Storage_Array_19387431c506) Get(index uint64) *Storage_BigInt {
	// Value: &wrapper.tmplField{Type:"BigInt", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48b9ce0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"", Slot:<nil>, Offset:0}
	actual := big.NewInt(0).Add(s.slot, big.NewInt(0).SetUint64(index*(32/32)))
	if s.obj[index] == nil {
		s.obj[index] = big.NewInt(0)

	}

	return &Storage_BigInt{
		obj:   s.obj[index],
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Array_54580f1bd6f5 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x5, Name:"Array_54580f1bd6f5", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00051c0e0)}, Type:"[3]Array_15880d9287c7", SolKind:abi.Type{Elem:(*abi.Type)(0xc00062e6c0), Kind:0x11, Type:(*reflect.rtype)(0xc00048ab40), Size:3, T:0x5, NumberOfBytes:0x180, StringKind:"t_array(t_array(t_struct(PendingDelegation_t)52_storage)2_storage)3_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Array_54580f1bd6f5 [3]Array_15880d9287c7

type Storage_Array_54580f1bd6f5 struct {
	obj   *Array_54580f1bd6f5
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func (s *Storage_Array_54580f1bd6f5) Length() int {
	return len(s.obj)
}

func (s *Storage_Array_54580f1bd6f5) Get(index uint64) *Storage_Array_15880d9287c7 {
	// Value: &wrapper.tmplField{Type:"Array_15880d9287c7", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc00062e750), Kind:0x11, Type:(*reflect.rtype)(0xc00048aaf0), Size:2, T:0x5, NumberOfBytes:0x80, StringKind:"t_array(t_struct(PendingDelegation_t)52_storage)2_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"", Slot:<nil>, Offset:0}
	actual := big.NewInt(0).Add(s.slot, big.NewInt(0).SetUint64(index*(128/32)))

	return &Storage_Array_15880d9287c7{
		obj:   &s.obj[index],
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Map_3783ff483121 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xd, Name:"Map_3783ff483121", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00051d420), (*wrapper.tmplField)(0xc00051d500)}, Type:"map[Address]*Bool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc0006363c0), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_address,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc000642ab0), (*abi.Type)(0xc000642b40)}, TupleRawNames:[]string(nil)}}
type Map_3783ff483121 map[Address]*Bool

type Storage_Map_3783ff483121 struct {
	obj   Map_3783ff483121
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func (s *Storage_Map_3783ff483121) Get(key Address) *Storage_Bool {
	// Key:	&wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x11, Type:(*reflect.rtype)(0x485c1e0), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"Bool", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x47720c0), Size:0, T:0x2, NumberOfBytes:0x1, StringKind:"t_bool", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"", Slot:<nil>, Offset:0}
	keyBytes := key.Hash().Bytes()
	hash := crypto.Keccak256Hash(append(keyBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()
	if s.obj[key] == nil {
		s.obj[key] = new(Bool)

	}

	return &Storage_Bool{
		obj:   s.obj[key],
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Map_624bcc0df495 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xd, Name:"Map_624bcc0df495", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00051db20), (*wrapper.tmplField)(0xc00051dc00)}, Type:"map[String]*Bool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc0006364e0), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_string_memory_ptr,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc000642d80), (*abi.Type)(0xc000642e10)}, TupleRawNames:[]string(nil)}}
type Map_624bcc0df495 map[String]*Bool

type Storage_Map_624bcc0df495 struct {
	obj   Map_624bcc0df495
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func (s *Storage_Map_624bcc0df495) Get(key String) *Storage_Bool {
	// Key:	&wrapper.tmplField{Type:"String", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x18, Type:(*reflect.rtype)(0x479b560), Size:0, T:0x3, NumberOfBytes:0x20, StringKind:"t_string_memory_ptr", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"Bool", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x47720c0), Size:0, T:0x2, NumberOfBytes:0x1, StringKind:"t_bool", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"", Slot:<nil>, Offset:0}
	keyBytes := []byte(key)
	hash := crypto.Keccak256Hash(append(keyBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()
	if s.obj[key] == nil {
		s.obj[key] = new(Bool)

	}

	return &Storage_Bool{
		obj:   s.obj[key],
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Map_aa079853f0d3 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xd, Name:"Map_aa079853f0d3", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00051cfc0), (*wrapper.tmplField)(0xc00051d0a0)}, Type:"map[Address]*Map3NodeWrapper_t", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc0006361e0), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_address,t_struct(Map3NodeWrapper_t)110_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc00062f560), (*abi.Type)(0xc00062f5f0)}, TupleRawNames:[]string(nil)}}
type Map_aa079853f0d3 map[Address]*Map3NodeWrapper_t

type Storage_Map_aa079853f0d3 struct {
	obj   Map_aa079853f0d3
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func (s *Storage_Map_aa079853f0d3) Get(key Address) *Storage_Map3NodeWrapper_t {
	// Key:	&wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x11, Type:(*reflect.rtype)(0x485c1e0), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"Map3NodeWrapper_t", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0xc00048dc80), Size:0, T:0x6, NumberOfBytes:0x340, StringKind:"t_struct(Map3NodeWrapper_t)110_storage", TupleRawName:"Map3NodeWrapper_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00062f680), (*abi.Type)(0xc00062f710), (*abi.Type)(0xc000642120), (*abi.Type)(0xc0006421b0), (*abi.Type)(0xc000642240), (*abi.Type)(0xc0006422d0), (*abi.Type)(0xc000642360)}, TupleRawNames:[]string{"Map3Node", "Microdelegations", "RedelegationReference", "AccumulatedReward", "nodeState", "TotalDelegation", "TotalPendingDelegation"}}, Tag:"", Slot:<nil>, Offset:0}
	keyBytes := key.Hash().Bytes()
	hash := crypto.Keccak256Hash(append(keyBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()
	if s.obj[key] == nil {
		s.obj[key] = new(Map3NodeWrapper_t)

	}

	return &Storage_Map3NodeWrapper_t{
		obj:   s.obj[key],
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Map_da81b317a4e9 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xd, Name:"Map_da81b317a4e9", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0004dd420), (*wrapper.tmplField)(0xc00051c380)}, Type:"map[Address]*Microdelegation_t", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc000636120), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_address,t_struct(Microdelegation_t)77_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc00062f7a0), (*abi.Type)(0xc00062f830)}, TupleRawNames:[]string(nil)}}
type Map_da81b317a4e9 map[Address]*Microdelegation_t

type Storage_Map_da81b317a4e9 struct {
	obj   Map_da81b317a4e9
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func (s *Storage_Map_da81b317a4e9) Get(key Address) *Storage_Microdelegation_t {
	// Key:	&wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x11, Type:(*reflect.rtype)(0x485c1e0), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"Microdelegation_t", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0xc00048dbc0), Size:0, T:0x6, NumberOfBytes:0x2c0, StringKind:"t_struct(Microdelegation_t)77_storage", TupleRawName:"Microdelegation_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00062f8c0), (*abi.Type)(0xc00062f950), (*abi.Type)(0xc00062f9e0), (*abi.Type)(0xc00062fa70), (*abi.Type)(0xc00062fc20), (*abi.Type)(0xc00062fdd0), (*abi.Type)(0xc00062fe60), (*abi.Type)(0xc00062fef0)}, TupleRawNames:[]string{"DelegatorAddress", "Amount", "Reward", "Undelegations", "PendingDelegations", "PendingDelegationsfixed", "PendingDelegationsfixed2dimension", "AutoRenew"}}, Tag:"", Slot:<nil>, Offset:0}
	keyBytes := key.Hash().Bytes()
	hash := crypto.Keccak256Hash(append(keyBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()
	if s.obj[key] == nil {
		s.obj[key] = new(Microdelegation_t)

	}

	return &Storage_Microdelegation_t{
		obj:   s.obj[key],
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Map_dfd3c4f07730 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xd, Name:"Map_dfd3c4f07730", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00051d340), (*wrapper.tmplField)(0xc00051d5e0)}, Type:"map[Address]Map_3783ff483121", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc000636420), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_address,t_mapping(t_address,t_bool))", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc000642990), (*abi.Type)(0xc000642a20)}, TupleRawNames:[]string(nil)}}
type Map_dfd3c4f07730 map[Address]Map_3783ff483121

type Storage_Map_dfd3c4f07730 struct {
	obj   Map_dfd3c4f07730
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func (s *Storage_Map_dfd3c4f07730) Get(key Address) *Storage_Map_3783ff483121 {
	// Key:	&wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x11, Type:(*reflect.rtype)(0x485c1e0), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"Map_3783ff483121", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc0006363c0), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_address,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc000642ab0), (*abi.Type)(0xc000642b40)}, TupleRawNames:[]string(nil)}, Tag:"", Slot:<nil>, Offset:0}
	keyBytes := key.Hash().Bytes()
	hash := crypto.Keccak256Hash(append(keyBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()
	if s.obj[key] == nil {
		s.obj[key] = make(Map_3783ff483121)

	}

	return &Storage_Map_3783ff483121{
		obj:   s.obj[key],
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Map_eed20e3b37e8 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xd, Name:"Map_eed20e3b37e8", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00051d7a0), (*wrapper.tmplField)(0xc00051d880)}, Type:"map[Uint64]*Map3NodeSnapshot_t", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc0006362a0), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_uint64,t_struct(Map3NodeSnapshot_t)117_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc00062f320), (*abi.Type)(0xc00062f440)}, TupleRawNames:[]string(nil)}}
type Map_eed20e3b37e8 map[Uint64]*Map3NodeSnapshot_t

type Storage_Map_eed20e3b37e8 struct {
	obj   Map_eed20e3b37e8
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func (s *Storage_Map_eed20e3b37e8) Get(key Uint64) *Storage_Map3NodeSnapshot_t {
	// Key:	&wrapper.tmplField{Type:"Uint64", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0xb, Type:(*reflect.rtype)(0x479b7a0), Size:64, T:0x1, NumberOfBytes:0x8, StringKind:"t_uint64", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"Map3NodeSnapshot_t", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0xc00048dcc0), Size:0, T:0x6, NumberOfBytes:0x40, StringKind:"t_struct(Map3NodeSnapshot_t)117_storage", TupleRawName:"Map3NodeSnapshot_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00062f4d0), (*abi.Type)(0xc000642510)}, TupleRawNames:[]string{"Map3Nodes", "Epoch"}}, Tag:"", Slot:<nil>, Offset:0}
	keyBytes := big.NewInt(0).SetUint64(uint64(key)).Bytes()
	hash := crypto.Keccak256Hash(append(keyBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()
	if s.obj[key] == nil {
		s.obj[key] = new(Map3NodeSnapshot_t)

	}

	return &Storage_Map3NodeSnapshot_t{
		obj:   s.obj[key],
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Map_f75907254fb8 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xd, Name:"Map_f75907254fb8", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00051dea0), (*wrapper.tmplField)(0xc000538000)}, Type:"map[String]BigInt", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc000243ce0), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_string_memory_ptr,t_uint256)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc00062f050), (*abi.Type)(0xc00062f170)}, TupleRawNames:[]string(nil)}}
type Map_f75907254fb8 map[String]BigInt

type Storage_Map_f75907254fb8 struct {
	obj   Map_f75907254fb8
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func (s *Storage_Map_f75907254fb8) Get(key String) *Storage_BigInt {
	// Key:	&wrapper.tmplField{Type:"String", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x18, Type:(*reflect.rtype)(0x479b560), Size:0, T:0x3, NumberOfBytes:0x20, StringKind:"t_string_memory_ptr", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"BigInt", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48b9ce0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"", Slot:<nil>, Offset:0}
	keyBytes := []byte(key)
	hash := crypto.Keccak256Hash(append(keyBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()
	if s.obj[key] == nil {
		s.obj[key] = big.NewInt(0)

	}

	return &Storage_BigInt{
		obj:   s.obj[key],
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Slice_7553d099f0a1 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x4, Name:"Slice_7553d099f0a1", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0004ddce0)}, Type:"[]*PendingDelegation_t", SolKind:abi.Type{Elem:(*abi.Type)(0xc00062fcb0), Kind:0x17, Type:(*reflect.rtype)(0x4788540), Size:0, T:0x4, NumberOfBytes:0x20, StringKind:"t_array(t_struct(PendingDelegation_t)52_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Slice_7553d099f0a1 []*PendingDelegation_t

type Storage_Slice_7553d099f0a1 struct {
	obj   Slice_7553d099f0a1
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func (s *Storage_Slice_7553d099f0a1) Length() *big.Int {
	rv := s.db.GetState(s.addr, common.BigToHash(s.slot))
	return rv.Big()
}

func (s *Storage_Slice_7553d099f0a1) Get(index uint64) *Storage_PendingDelegation_t {
	// Value: &wrapper.tmplField{Type:"PendingDelegation_t", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0xc00048db40), Size:0, T:0x6, NumberOfBytes:0x40, StringKind:"t_struct(PendingDelegation_t)52_storage", TupleRawName:"PendingDelegation_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00062e7e0), (*abi.Type)(0xc00062e870)}, TupleRawNames:[]string{"Amount", "Epoch"}}, Tag:"", Slot:<nil>, Offset:0}
	indexBytes := common.BigToHash(big.NewInt(0).SetUint64(index)).Bytes()
	hash := crypto.Keccak256Hash(append(indexBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()

	return &Storage_PendingDelegation_t{
		obj:   s.obj[index],
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Slice_9f684cf5f20e is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x4, Name:"Slice_9f684cf5f20e", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0004dd960)}, Type:"[]*Undelegation_t", SolKind:abi.Type{Elem:(*abi.Type)(0xc00062fb00), Kind:0x17, Type:(*reflect.rtype)(0x4788540), Size:0, T:0x4, NumberOfBytes:0x20, StringKind:"t_array(t_struct(Undelegation_t)47_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Slice_9f684cf5f20e []*Undelegation_t

type Storage_Slice_9f684cf5f20e struct {
	obj   Slice_9f684cf5f20e
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func (s *Storage_Slice_9f684cf5f20e) Length() *big.Int {
	rv := s.db.GetState(s.addr, common.BigToHash(s.slot))
	return rv.Big()
}

func (s *Storage_Slice_9f684cf5f20e) Get(index uint64) *Storage_Undelegation_t {
	// Value: &wrapper.tmplField{Type:"Undelegation_t", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0xc00048db40), Size:0, T:0x6, NumberOfBytes:0x40, StringKind:"t_struct(Undelegation_t)47_storage", TupleRawName:"Undelegation_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00062e240), (*abi.Type)(0xc00062e2d0)}, TupleRawNames:[]string{"Amount", "Epoch"}}, Tag:"", Slot:<nil>, Offset:0}
	indexBytes := common.BigToHash(big.NewInt(0).SetUint64(index)).Bytes()
	hash := crypto.Keccak256Hash(append(indexBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()

	return &Storage_Undelegation_t{
		obj:   s.obj[index],
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// CommissionRates_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"CommissionRates_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0004dc620), (*wrapper.tmplField)(0xc0004dc700), (*wrapper.tmplField)(0xc0004dc7e0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type CommissionRates_t struct {
	Rate          BigInt `json:"Rate" storage:"slot=0,offset=0"`
	MaxRate       BigInt `json:"MaxRate" storage:"slot=1,offset=0"`
	MaxChangeRate BigInt `json:"MaxChangeRate" storage:"slot=2,offset=0"`
}

// Commission_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Commission_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0004dc8c0), (*wrapper.tmplField)(0xc0004dc9a0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Commission_t struct {
	CommissionRates CommissionRates_t `json:"CommissionRates" storage:"slot=0,offset=0"`
	UpdateHeight    BigInt            `json:"UpdateHeight" storage:"slot=3,offset=0"`
}

// Description_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Description_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0004dcb60), (*wrapper.tmplField)(0xc0004dcc40), (*wrapper.tmplField)(0xc0004dcd20), (*wrapper.tmplField)(0xc0004dce00), (*wrapper.tmplField)(0xc0004dcee0), (*wrapper.tmplField)(0xc0004dd0a0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Description_t struct {
	Name             String             `json:"Name" storage:"slot=0,offset=0"`
	Identity         String             `json:"Identity" storage:"slot=1,offset=0"`
	WebSite          String             `json:"WebSite" storage:"slot=2,offset=0"`
	SecurityContract String             `json:"SecurityContract" storage:"slot=3,offset=0"`
	Details          String             `json:"Details" storage:"slot=4,offset=0"`
	Version          Array_19387431c506 `json:"Version" storage:"slot=5,offset=0"`
}

// Global_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Global_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0005381c0), (*wrapper.tmplField)(0xc0005382a0), (*wrapper.tmplField)(0xc000538380), (*wrapper.tmplField)(0xc000538460)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Global_t struct {
	Pool    Map3NodePool_t `json:"pool" storage:"slot=0,offset=0"`
	Node    Map3Node_t     `json:"node" storage:"slot=6,offset=0"`
	Version BigInt         `json:"version" storage:"slot=21,offset=0"`
	Name    String         `json:"name" storage:"slot=22,offset=0"`
}

// Map3NodePool_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Map3NodePool_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00051d6c0), (*wrapper.tmplField)(0xc00051d960), (*wrapper.tmplField)(0xc00051da40), (*wrapper.tmplField)(0xc00051dce0), (*wrapper.tmplField)(0xc00051ddc0), (*wrapper.tmplField)(0xc0005380e0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Map3NodePool_t struct {
	Nodes                     Map_aa079853f0d3 `json:"Nodes" storage:"slot=0,offset=0"`
	NodeSnapshotByEpoch       Map_eed20e3b37e8 `json:"NodeSnapshotByEpoch" storage:"slot=1,offset=0"`
	NodeAddressSetByDelegator Map_dfd3c4f07730 `json:"NodeAddressSetByDelegator" storage:"slot=2,offset=0"`
	NodeKeySet                Map_624bcc0df495 `json:"NodeKeySet" storage:"slot=3,offset=0"`
	DescriptionIdentitySet    Map_624bcc0df495 `json:"DescriptionIdentitySet" storage:"slot=4,offset=0"`
	NodePriority              Map_f75907254fb8 `json:"NodePriority" storage:"slot=5,offset=0"`
}

// Map3NodeSnapshot_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Map3NodeSnapshot_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00051d180), (*wrapper.tmplField)(0xc00051d260)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Map3NodeSnapshot_t struct {
	Map3Nodes Map_aa079853f0d3 `json:"Map3Nodes" storage:"slot=0,offset=0"`
	Epoch     BigInt           `json:"Epoch" storage:"slot=1,offset=0"`
}

// Map3NodeWrapper_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Map3NodeWrapper_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0004dd340), (*wrapper.tmplField)(0xc00051c460), (*wrapper.tmplField)(0xc00051c700), (*wrapper.tmplField)(0xc00051c7e0), (*wrapper.tmplField)(0xc00051cd20), (*wrapper.tmplField)(0xc00051ce00), (*wrapper.tmplField)(0xc00051cee0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Map3NodeWrapper_t struct {
	Map3Node               Map3Node_t              `json:"Map3Node" storage:"slot=0,offset=0"`
	Microdelegations       Map_da81b317a4e9        `json:"Microdelegations" storage:"slot=15,offset=0"`
	RedelegationReference  RedelegationReference_t `json:"RedelegationReference" storage:"slot=16,offset=0"`
	AccumulatedReward      BigInt                  `json:"AccumulatedReward" storage:"slot=18,offset=0"`
	NodeState              NodeState_t             `json:"nodeState" storage:"slot=19,offset=0"`
	TotalDelegation        BigInt                  `json:"TotalDelegation" storage:"slot=24,offset=0"`
	TotalPendingDelegation BigInt                  `json:"TotalPendingDelegation" storage:"slot=25,offset=0"`
}

// Map3Node_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Map3Node_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0004dc0e0), (*wrapper.tmplField)(0xc0004dc460), (*wrapper.tmplField)(0xc0004dc540), (*wrapper.tmplField)(0xc0004dca80), (*wrapper.tmplField)(0xc0004dd180), (*wrapper.tmplField)(0xc0004dd260)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Map3Node_t struct {
	NodeAddress      Address       `json:"NodeAddress" storage:"slot=0,offset=0"`
	InitiatorAddress Address       `json:"InitiatorAddress" storage:"slot=1,offset=0"`
	NodeKeys         Bytes         `json:"NodeKeys" storage:"slot=2,offset=0"`
	Commission       Commission_t  `json:"Commission" storage:"slot=3,offset=0"`
	Description      Description_t `json:"Description" storage:"slot=7,offset=0"`
	SplittedFrom     Address       `json:"SplittedFrom" storage:"slot=14,offset=0"`
}

// Microdelegation_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Microdelegation_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0004dd500), (*wrapper.tmplField)(0xc0004dd5e0), (*wrapper.tmplField)(0xc0004dd6c0), (*wrapper.tmplField)(0xc0004dda40), (*wrapper.tmplField)(0xc0004dddc0), (*wrapper.tmplField)(0xc00051c000), (*wrapper.tmplField)(0xc00051c1c0), (*wrapper.tmplField)(0xc00051c2a0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Microdelegation_t struct {
	DelegatorAddress                  Address            `json:"DelegatorAddress" storage:"slot=0,offset=0"`
	Amount                            BigInt             `json:"Amount" storage:"slot=1,offset=0"`
	Reward                            BigInt             `json:"Reward" storage:"slot=2,offset=0"`
	Undelegations                     Slice_9f684cf5f20e `json:"Undelegations" storage:"slot=3,offset=0"`
	PendingDelegations                Slice_7553d099f0a1 `json:"PendingDelegations" storage:"slot=4,offset=0"`
	PendingDelegationsfixed           Array_15880d9287c7 `json:"PendingDelegationsfixed" storage:"slot=5,offset=0"`
	PendingDelegationsfixed2dimension Array_54580f1bd6f5 `json:"PendingDelegationsfixed2dimension" storage:"slot=9,offset=0"`
	AutoRenew                         Bool               `json:"AutoRenew" storage:"slot=21,offset=0"`
}

// NodeState_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"NodeState_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00051c8c0), (*wrapper.tmplField)(0xc00051c9a0), (*wrapper.tmplField)(0xc00051ca80), (*wrapper.tmplField)(0xc00051cb60), (*wrapper.tmplField)(0xc00051cc40)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type NodeState_t struct {
	Status          Uint8  `json:"Status" storage:"slot=0,offset=0"`
	NodeAge         BigInt `json:"NodeAge" storage:"slot=1,offset=0"`
	CreationEpoch   BigInt `json:"CreationEpoch" storage:"slot=2,offset=0"`
	ActivationEpoch BigInt `json:"ActivationEpoch" storage:"slot=3,offset=0"`
	ReleaseEpoch    BigInt `json:"ReleaseEpoch" storage:"slot=4,offset=0"`
}

// PendingDelegation_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"PendingDelegation_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0004ddb20), (*wrapper.tmplField)(0xc0004ddc00)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type PendingDelegation_t struct {
	Amount BigInt `json:"Amount" storage:"slot=0,offset=0"`
	Epoch  BigInt `json:"Epoch" storage:"slot=1,offset=0"`
}

// RedelegationReference_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"RedelegationReference_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00051c540), (*wrapper.tmplField)(0xc00051c620)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type RedelegationReference_t struct {
	ValidatorAddress        Address `json:"ValidatorAddress" storage:"slot=0,offset=0"`
	ReleasedTotalDelegation BigInt  `json:"ReleasedTotalDelegation" storage:"slot=1,offset=0"`
}

// Undelegation_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Undelegation_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0004dd7a0), (*wrapper.tmplField)(0xc0004dd880)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Undelegation_t struct {
	Amount BigInt `json:"Amount" storage:"slot=0,offset=0"`
	Epoch  BigInt `json:"Epoch" storage:"slot=1,offset=0"`
}

type Storage_CommissionRates_t struct {
	obj   *CommissionRates_t
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

type Storage_Commission_t struct {
	obj   *Commission_t
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

type Storage_Description_t struct {
	obj   *Description_t
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

type Storage_Global_t struct {
	obj   *Global_t
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

type Storage_Map3NodePool_t struct {
	obj   *Map3NodePool_t
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

type Storage_Map3NodeSnapshot_t struct {
	obj   *Map3NodeSnapshot_t
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

type Storage_Map3NodeWrapper_t struct {
	obj   *Map3NodeWrapper_t
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

type Storage_Map3Node_t struct {
	obj   *Map3Node_t
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

type Storage_Microdelegation_t struct {
	obj   *Microdelegation_t
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

type Storage_NodeState_t struct {
	obj   *NodeState_t
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

type Storage_PendingDelegation_t struct {
	obj   *PendingDelegation_t
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

type Storage_RedelegationReference_t struct {
	obj   *RedelegationReference_t
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

type Storage_Undelegation_t struct {
	obj   *Undelegation_t
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func New(g *Global_t, db *state.StateDB, addr common.Address, slot *big.Int) *Storage_Global_t {
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
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Rate", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48b9ce0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Rate\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Rate == nil {
		s.obj.Rate = big.NewInt(0)
	}

	return &Storage_BigInt{
		obj:   s.obj.Rate,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// MaxRate is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_CommissionRates_t) MaxRate() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"MaxRate", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48b9ce0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"MaxRate\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.MaxRate == nil {
		s.obj.MaxRate = big.NewInt(0)
	}

	return &Storage_BigInt{
		obj:   s.obj.MaxRate,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// MaxChangeRate is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_CommissionRates_t) MaxChangeRate() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"MaxChangeRate", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48b9ce0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"MaxChangeRate\" storage:\"slot=2,offset=0\"", Slot:2, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("2", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.MaxChangeRate == nil {
		s.obj.MaxChangeRate = big.NewInt(0)
	}

	return &Storage_BigInt{
		obj:   s.obj.MaxChangeRate,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// CommissionRates is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Commission_t) CommissionRates() *Storage_CommissionRates_t {
	// Field: &wrapper.tmplField{Type:"CommissionRates_t", ElemType:"CommissionRates_t", Name:"CommissionRates", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000243560), Size:0, T:0x6, NumberOfBytes:0x60, StringKind:"t_struct(CommissionRates_t)9_storage", TupleRawName:"CommissionRates_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00060b4d0), (*abi.Type)(0xc00060b5f0), (*abi.Type)(0xc00060b680)}, TupleRawNames:[]string{"Rate", "MaxRate", "MaxChangeRate"}}, Tag:"json:\"CommissionRates\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_CommissionRates_t{
		obj:   &s.obj.CommissionRates,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// UpdateHeight is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Commission_t) UpdateHeight() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"UpdateHeight", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48b9ce0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"UpdateHeight\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("3", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.UpdateHeight == nil {
		s.obj.UpdateHeight = big.NewInt(0)
	}

	return &Storage_BigInt{
		obj:   s.obj.UpdateHeight,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Name is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Description_t) Name() *Storage_String {
	// Field: &wrapper.tmplField{Type:"String", ElemType:"String", Name:"Name", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x18, Type:(*reflect.rtype)(0x479b560), Size:0, T:0x3, NumberOfBytes:0x20, StringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Name\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_String{
		obj:   &s.obj.Name,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Identity is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Description_t) Identity() *Storage_String {
	// Field: &wrapper.tmplField{Type:"String", ElemType:"String", Name:"Identity", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x18, Type:(*reflect.rtype)(0x479b560), Size:0, T:0x3, NumberOfBytes:0x20, StringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Identity\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_String{
		obj:   &s.obj.Identity,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// WebSite is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Description_t) WebSite() *Storage_String {
	// Field: &wrapper.tmplField{Type:"String", ElemType:"String", Name:"WebSite", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x18, Type:(*reflect.rtype)(0x479b560), Size:0, T:0x3, NumberOfBytes:0x20, StringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"WebSite\" storage:\"slot=2,offset=0\"", Slot:2, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("2", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_String{
		obj:   &s.obj.WebSite,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// SecurityContract is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Description_t) SecurityContract() *Storage_String {
	// Field: &wrapper.tmplField{Type:"String", ElemType:"String", Name:"SecurityContract", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x18, Type:(*reflect.rtype)(0x479b560), Size:0, T:0x3, NumberOfBytes:0x20, StringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"SecurityContract\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("3", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_String{
		obj:   &s.obj.SecurityContract,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Details is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Description_t) Details() *Storage_String {
	// Field: &wrapper.tmplField{Type:"String", ElemType:"String", Name:"Details", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x18, Type:(*reflect.rtype)(0x479b560), Size:0, T:0x3, NumberOfBytes:0x20, StringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Details\" storage:\"slot=4,offset=0\"", Slot:4, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("4", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_String{
		obj:   &s.obj.Details,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Version is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Description_t) Version() *Storage_Array_19387431c506 {
	// Field: &wrapper.tmplField{Type:"Array_19387431c506", ElemType:"Array_19387431c506", Name:"Version", SolKind:abi.Type{Elem:(*abi.Type)(0xc00060bd40), Kind:0x11, Type:(*reflect.rtype)(0xc00048aaf0), Size:2, T:0x5, NumberOfBytes:0x40, StringKind:"t_array(t_uint256)2_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Version\" storage:\"slot=5,offset=0\"", Slot:5, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("5", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Array_19387431c506{
		obj:   &s.obj.Version,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Pool is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Global_t) Pool() *Storage_Map3NodePool_t {
	// Field: &wrapper.tmplField{Type:"Map3NodePool_t", ElemType:"Map3NodePool_t", Name:"Pool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc0006365a0), Size:0, T:0x6, NumberOfBytes:0xc0, StringKind:"t_struct(Map3NodePool_t)144_storage", TupleRawName:"Map3NodePool_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc000643050), (*abi.Type)(0xc0006430e0), (*abi.Type)(0xc000643170), (*abi.Type)(0xc000643200), (*abi.Type)(0xc000643290), (*abi.Type)(0xc000643320)}, TupleRawNames:[]string{"Nodes", "NodeSnapshotByEpoch", "NodeAddressSetByDelegator", "NodeKeySet", "DescriptionIdentitySet", "NodePriority"}}, Tag:"json:\"pool\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Map3NodePool_t{
		obj:   &s.obj.Pool,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Node is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Global_t) Node() *Storage_Map3Node_t {
	// Field: &wrapper.tmplField{Type:"Map3Node_t", ElemType:"Map3Node_t", Name:"Node", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc0002437a0), Size:0, T:0x6, NumberOfBytes:0x1e0, StringKind:"t_struct(Map3Node_t)42_storage", TupleRawName:"Map3Node_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00060b0e0), (*abi.Type)(0xc00060b200), (*abi.Type)(0xc00060b290), (*abi.Type)(0xc00060b3b0), (*abi.Type)(0xc00060b8c0), (*abi.Type)(0xc00062e000)}, TupleRawNames:[]string{"NodeAddress", "InitiatorAddress", "NodeKeys", "Commission", "Description", "SplittedFrom"}}, Tag:"json:\"node\" storage:\"slot=6,offset=0\"", Slot:6, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("6", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Map3Node_t{
		obj:   &s.obj.Node,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Version is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Global_t) Version() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Version", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48b9ce0), Size:256, T:0x0, NumberOfBytes:0x20, StringKind:"t_int256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"version\" storage:\"slot=21,offset=0\"", Slot:21, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("21", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Version == nil {
		s.obj.Version = big.NewInt(0)
	}

	return &Storage_BigInt{
		obj:   s.obj.Version,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Name is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Global_t) Name() *Storage_String {
	// Field: &wrapper.tmplField{Type:"String", ElemType:"String", Name:"Name", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x18, Type:(*reflect.rtype)(0x479b560), Size:0, T:0x3, NumberOfBytes:0x20, StringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"name\" storage:\"slot=22,offset=0\"", Slot:22, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("22", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_String{
		obj:   &s.obj.Name,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Nodes is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodePool_t) Nodes() *Storage_Map_aa079853f0d3 {
	// Field: &wrapper.tmplField{Type:"Map_aa079853f0d3", ElemType:"Map_aa079853f0d3", Name:"Nodes", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc0006361e0), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_address,t_struct(Map3NodeWrapper_t)110_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc00062f560), (*abi.Type)(0xc00062f5f0)}, TupleRawNames:[]string(nil)}, Tag:"json:\"Nodes\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Nodes == nil {
		s.obj.Nodes = make(Map_aa079853f0d3)
	}

	return &Storage_Map_aa079853f0d3{
		obj:   s.obj.Nodes,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// NodeSnapshotByEpoch is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodePool_t) NodeSnapshotByEpoch() *Storage_Map_eed20e3b37e8 {
	// Field: &wrapper.tmplField{Type:"Map_eed20e3b37e8", ElemType:"Map_eed20e3b37e8", Name:"NodeSnapshotByEpoch", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc0006362a0), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_uint64,t_struct(Map3NodeSnapshot_t)117_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc00062f320), (*abi.Type)(0xc00062f440)}, TupleRawNames:[]string(nil)}, Tag:"json:\"NodeSnapshotByEpoch\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.NodeSnapshotByEpoch == nil {
		s.obj.NodeSnapshotByEpoch = make(Map_eed20e3b37e8)
	}

	return &Storage_Map_eed20e3b37e8{
		obj:   s.obj.NodeSnapshotByEpoch,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// NodeAddressSetByDelegator is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodePool_t) NodeAddressSetByDelegator() *Storage_Map_dfd3c4f07730 {
	// Field: &wrapper.tmplField{Type:"Map_dfd3c4f07730", ElemType:"Map_dfd3c4f07730", Name:"NodeAddressSetByDelegator", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc000636420), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_address,t_mapping(t_address,t_bool))", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc000642990), (*abi.Type)(0xc000642a20)}, TupleRawNames:[]string(nil)}, Tag:"json:\"NodeAddressSetByDelegator\" storage:\"slot=2,offset=0\"", Slot:2, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("2", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.NodeAddressSetByDelegator == nil {
		s.obj.NodeAddressSetByDelegator = make(Map_dfd3c4f07730)
	}

	return &Storage_Map_dfd3c4f07730{
		obj:   s.obj.NodeAddressSetByDelegator,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// NodeKeySet is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodePool_t) NodeKeySet() *Storage_Map_624bcc0df495 {
	// Field: &wrapper.tmplField{Type:"Map_624bcc0df495", ElemType:"Map_624bcc0df495", Name:"NodeKeySet", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc0006364e0), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_string_memory_ptr,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc000642d80), (*abi.Type)(0xc000642e10)}, TupleRawNames:[]string(nil)}, Tag:"json:\"NodeKeySet\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("3", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.NodeKeySet == nil {
		s.obj.NodeKeySet = make(Map_624bcc0df495)
	}

	return &Storage_Map_624bcc0df495{
		obj:   s.obj.NodeKeySet,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// DescriptionIdentitySet is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodePool_t) DescriptionIdentitySet() *Storage_Map_624bcc0df495 {
	// Field: &wrapper.tmplField{Type:"Map_624bcc0df495", ElemType:"Map_624bcc0df495", Name:"DescriptionIdentitySet", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc0006364e0), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_string_memory_ptr,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc000642d80), (*abi.Type)(0xc000642e10)}, TupleRawNames:[]string(nil)}, Tag:"json:\"DescriptionIdentitySet\" storage:\"slot=4,offset=0\"", Slot:4, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("4", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.DescriptionIdentitySet == nil {
		s.obj.DescriptionIdentitySet = make(Map_624bcc0df495)
	}

	return &Storage_Map_624bcc0df495{
		obj:   s.obj.DescriptionIdentitySet,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// NodePriority is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodePool_t) NodePriority() *Storage_Map_f75907254fb8 {
	// Field: &wrapper.tmplField{Type:"Map_f75907254fb8", ElemType:"Map_f75907254fb8", Name:"NodePriority", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc000243ce0), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_string_memory_ptr,t_uint256)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc00062f050), (*abi.Type)(0xc00062f170)}, TupleRawNames:[]string(nil)}, Tag:"json:\"NodePriority\" storage:\"slot=5,offset=0\"", Slot:5, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("5", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.NodePriority == nil {
		s.obj.NodePriority = make(Map_f75907254fb8)
	}

	return &Storage_Map_f75907254fb8{
		obj:   s.obj.NodePriority,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Map3Nodes is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeSnapshot_t) Map3Nodes() *Storage_Map_aa079853f0d3 {
	// Field: &wrapper.tmplField{Type:"Map_aa079853f0d3", ElemType:"Map_aa079853f0d3", Name:"Map3Nodes", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc0006361e0), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_address,t_struct(Map3NodeWrapper_t)110_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc00062f560), (*abi.Type)(0xc00062f5f0)}, TupleRawNames:[]string(nil)}, Tag:"json:\"Map3Nodes\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Map3Nodes == nil {
		s.obj.Map3Nodes = make(Map_aa079853f0d3)
	}

	return &Storage_Map_aa079853f0d3{
		obj:   s.obj.Map3Nodes,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Epoch is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeSnapshot_t) Epoch() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Epoch", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48b9ce0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Epoch\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Epoch == nil {
		s.obj.Epoch = big.NewInt(0)
	}

	return &Storage_BigInt{
		obj:   s.obj.Epoch,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Map3Node is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeWrapper_t) Map3Node() *Storage_Map3Node_t {
	// Field: &wrapper.tmplField{Type:"Map3Node_t", ElemType:"Map3Node_t", Name:"Map3Node", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc0002437a0), Size:0, T:0x6, NumberOfBytes:0x1e0, StringKind:"t_struct(Map3Node_t)42_storage", TupleRawName:"Map3Node_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00060b0e0), (*abi.Type)(0xc00060b200), (*abi.Type)(0xc00060b290), (*abi.Type)(0xc00060b3b0), (*abi.Type)(0xc00060b8c0), (*abi.Type)(0xc00062e000)}, TupleRawNames:[]string{"NodeAddress", "InitiatorAddress", "NodeKeys", "Commission", "Description", "SplittedFrom"}}, Tag:"json:\"Map3Node\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Map3Node_t{
		obj:   &s.obj.Map3Node,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Microdelegations is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeWrapper_t) Microdelegations() *Storage_Map_da81b317a4e9 {
	// Field: &wrapper.tmplField{Type:"Map_da81b317a4e9", ElemType:"Map_da81b317a4e9", Name:"Microdelegations", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc000636120), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_address,t_struct(Microdelegation_t)77_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc00062f7a0), (*abi.Type)(0xc00062f830)}, TupleRawNames:[]string(nil)}, Tag:"json:\"Microdelegations\" storage:\"slot=15,offset=0\"", Slot:15, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("15", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Microdelegations == nil {
		s.obj.Microdelegations = make(Map_da81b317a4e9)
	}

	return &Storage_Map_da81b317a4e9{
		obj:   s.obj.Microdelegations,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// RedelegationReference is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeWrapper_t) RedelegationReference() *Storage_RedelegationReference_t {
	// Field: &wrapper.tmplField{Type:"RedelegationReference_t", ElemType:"RedelegationReference_t", Name:"RedelegationReference", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000243920), Size:0, T:0x6, NumberOfBytes:0x40, StringKind:"t_struct(RedelegationReference_t)82_storage", TupleRawName:"RedelegationReference_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00062e480), (*abi.Type)(0xc00062e510)}, TupleRawNames:[]string{"ValidatorAddress", "ReleasedTotalDelegation"}}, Tag:"json:\"RedelegationReference\" storage:\"slot=16,offset=0\"", Slot:16, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("16", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_RedelegationReference_t{
		obj:   &s.obj.RedelegationReference,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// AccumulatedReward is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeWrapper_t) AccumulatedReward() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"AccumulatedReward", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48b9ce0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"AccumulatedReward\" storage:\"slot=18,offset=0\"", Slot:18, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("18", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.AccumulatedReward == nil {
		s.obj.AccumulatedReward = big.NewInt(0)
	}

	return &Storage_BigInt{
		obj:   s.obj.AccumulatedReward,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// NodeState is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeWrapper_t) NodeState() *Storage_NodeState_t {
	// Field: &wrapper.tmplField{Type:"NodeState_t", ElemType:"NodeState_t", Name:"NodeState", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000243bc0), Size:0, T:0x6, NumberOfBytes:0xa0, StringKind:"t_struct(NodeState_t)93_storage", TupleRawName:"NodeState_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00062ec60), (*abi.Type)(0xc00062ecf0), (*abi.Type)(0xc00062ed80), (*abi.Type)(0xc00062ee10), (*abi.Type)(0xc00062eea0)}, TupleRawNames:[]string{"Status", "NodeAge", "CreationEpoch", "ActivationEpoch", "ReleaseEpoch"}}, Tag:"json:\"nodeState\" storage:\"slot=19,offset=0\"", Slot:19, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("19", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_NodeState_t{
		obj:   &s.obj.NodeState,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// TotalDelegation is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeWrapper_t) TotalDelegation() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"TotalDelegation", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48b9ce0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"TotalDelegation\" storage:\"slot=24,offset=0\"", Slot:24, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("24", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.TotalDelegation == nil {
		s.obj.TotalDelegation = big.NewInt(0)
	}

	return &Storage_BigInt{
		obj:   s.obj.TotalDelegation,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// TotalPendingDelegation is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeWrapper_t) TotalPendingDelegation() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"TotalPendingDelegation", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48b9ce0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"TotalPendingDelegation\" storage:\"slot=25,offset=0\"", Slot:25, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("25", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.TotalPendingDelegation == nil {
		s.obj.TotalPendingDelegation = big.NewInt(0)
	}

	return &Storage_BigInt{
		obj:   s.obj.TotalPendingDelegation,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// NodeAddress is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3Node_t) NodeAddress() *Storage_Address {
	// Field: &wrapper.tmplField{Type:"Address", ElemType:"Address", Name:"NodeAddress", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x11, Type:(*reflect.rtype)(0x485c1e0), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"NodeAddress\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Address{
		obj:   &s.obj.NodeAddress,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// InitiatorAddress is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3Node_t) InitiatorAddress() *Storage_Address {
	// Field: &wrapper.tmplField{Type:"Address", ElemType:"Address", Name:"InitiatorAddress", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x11, Type:(*reflect.rtype)(0x485c1e0), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"InitiatorAddress\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Address{
		obj:   &s.obj.InitiatorAddress,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// NodeKeys is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3Node_t) NodeKeys() *Storage_Bytes {
	// Field: &wrapper.tmplField{Type:"Bytes", ElemType:"Bytes", Name:"NodeKeys", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x17, Type:(*reflect.rtype)(0x478f1c0), Size:0, T:0x9, NumberOfBytes:0x20, StringKind:"t_bytes_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"NodeKeys\" storage:\"slot=2,offset=0\"", Slot:2, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("2", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.NodeKeys == nil {
		s.obj.NodeKeys = make(Bytes, 0)
	}

	return &Storage_Bytes{
		obj:   s.obj.NodeKeys,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Commission is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3Node_t) Commission() *Storage_Commission_t {
	// Field: &wrapper.tmplField{Type:"Commission_t", ElemType:"Commission_t", Name:"Commission", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc0002435c0), Size:0, T:0x6, NumberOfBytes:0x80, StringKind:"t_struct(Commission_t)14_storage", TupleRawName:"Commission_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00060b440), (*abi.Type)(0xc00060b7a0)}, TupleRawNames:[]string{"CommissionRates", "UpdateHeight"}}, Tag:"json:\"Commission\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("3", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Commission_t{
		obj:   &s.obj.Commission,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Description is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3Node_t) Description() *Storage_Description_t {
	// Field: &wrapper.tmplField{Type:"Description_t", ElemType:"Description_t", Name:"Description", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000243740), Size:0, T:0x6, NumberOfBytes:0xe0, StringKind:"t_struct(Description_t)29_storage", TupleRawName:"Description_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc00060b950), (*abi.Type)(0xc00060ba70), (*abi.Type)(0xc00060bb00), (*abi.Type)(0xc00060bb90), (*abi.Type)(0xc00060bc20), (*abi.Type)(0xc00060bcb0)}, TupleRawNames:[]string{"Name", "Identity", "WebSite", "SecurityContract", "Details", "Version"}}, Tag:"json:\"Description\" storage:\"slot=7,offset=0\"", Slot:7, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("7", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Description_t{
		obj:   &s.obj.Description,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// SplittedFrom is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3Node_t) SplittedFrom() *Storage_Address {
	// Field: &wrapper.tmplField{Type:"Address", ElemType:"Address", Name:"SplittedFrom", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x11, Type:(*reflect.rtype)(0x485c1e0), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"SplittedFrom\" storage:\"slot=14,offset=0\"", Slot:14, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("14", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Address{
		obj:   &s.obj.SplittedFrom,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// DelegatorAddress is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Microdelegation_t) DelegatorAddress() *Storage_Address {
	// Field: &wrapper.tmplField{Type:"Address", ElemType:"Address", Name:"DelegatorAddress", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x11, Type:(*reflect.rtype)(0x485c1e0), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"DelegatorAddress\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Address{
		obj:   &s.obj.DelegatorAddress,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Amount is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Microdelegation_t) Amount() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Amount", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48b9ce0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Amount\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Amount == nil {
		s.obj.Amount = big.NewInt(0)
	}

	return &Storage_BigInt{
		obj:   s.obj.Amount,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Reward is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Microdelegation_t) Reward() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Reward", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48b9ce0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Reward\" storage:\"slot=2,offset=0\"", Slot:2, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("2", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Reward == nil {
		s.obj.Reward = big.NewInt(0)
	}

	return &Storage_BigInt{
		obj:   s.obj.Reward,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Undelegations is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Microdelegation_t) Undelegations() *Storage_Slice_9f684cf5f20e {
	// Field: &wrapper.tmplField{Type:"Slice_9f684cf5f20e", ElemType:"Slice_9f684cf5f20e", Name:"Undelegations", SolKind:abi.Type{Elem:(*abi.Type)(0xc00062fb00), Kind:0x17, Type:(*reflect.rtype)(0x4788540), Size:0, T:0x4, NumberOfBytes:0x20, StringKind:"t_array(t_struct(Undelegation_t)47_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Undelegations\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("3", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Undelegations == nil {
		s.obj.Undelegations = make(Slice_9f684cf5f20e, 0)
	}

	return &Storage_Slice_9f684cf5f20e{
		obj:   s.obj.Undelegations,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// PendingDelegations is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Microdelegation_t) PendingDelegations() *Storage_Slice_7553d099f0a1 {
	// Field: &wrapper.tmplField{Type:"Slice_7553d099f0a1", ElemType:"Slice_7553d099f0a1", Name:"PendingDelegations", SolKind:abi.Type{Elem:(*abi.Type)(0xc00062fcb0), Kind:0x17, Type:(*reflect.rtype)(0x4788540), Size:0, T:0x4, NumberOfBytes:0x20, StringKind:"t_array(t_struct(PendingDelegation_t)52_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"PendingDelegations\" storage:\"slot=4,offset=0\"", Slot:4, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("4", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.PendingDelegations == nil {
		s.obj.PendingDelegations = make(Slice_7553d099f0a1, 0)
	}

	return &Storage_Slice_7553d099f0a1{
		obj:   s.obj.PendingDelegations,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// PendingDelegationsfixed is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Microdelegation_t) PendingDelegationsfixed() *Storage_Array_15880d9287c7 {
	// Field: &wrapper.tmplField{Type:"Array_15880d9287c7", ElemType:"Array_15880d9287c7", Name:"PendingDelegationsfixed", SolKind:abi.Type{Elem:(*abi.Type)(0xc00062e750), Kind:0x11, Type:(*reflect.rtype)(0xc00048aaf0), Size:2, T:0x5, NumberOfBytes:0x80, StringKind:"t_array(t_struct(PendingDelegation_t)52_storage)2_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"PendingDelegationsfixed\" storage:\"slot=5,offset=0\"", Slot:5, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("5", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Array_15880d9287c7{
		obj:   &s.obj.PendingDelegationsfixed,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// PendingDelegationsfixed2dimension is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Microdelegation_t) PendingDelegationsfixed2dimension() *Storage_Array_54580f1bd6f5 {
	// Field: &wrapper.tmplField{Type:"Array_54580f1bd6f5", ElemType:"Array_54580f1bd6f5", Name:"PendingDelegationsfixed2dimension", SolKind:abi.Type{Elem:(*abi.Type)(0xc00062e6c0), Kind:0x11, Type:(*reflect.rtype)(0xc00048ab40), Size:3, T:0x5, NumberOfBytes:0x180, StringKind:"t_array(t_array(t_struct(PendingDelegation_t)52_storage)2_storage)3_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"PendingDelegationsfixed2dimension\" storage:\"slot=9,offset=0\"", Slot:9, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("9", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Array_54580f1bd6f5{
		obj:   &s.obj.PendingDelegationsfixed2dimension,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// AutoRenew is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Microdelegation_t) AutoRenew() *Storage_Bool {
	// Field: &wrapper.tmplField{Type:"Bool", ElemType:"Bool", Name:"AutoRenew", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x1, Type:(*reflect.rtype)(0x47921a0), Size:0, T:0x2, NumberOfBytes:0x1, StringKind:"t_bool", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"AutoRenew\" storage:\"slot=21,offset=0\"", Slot:21, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("21", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Bool{
		obj:   &s.obj.AutoRenew,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Status is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_NodeState_t) Status() *Storage_Uint8 {
	// Field: &wrapper.tmplField{Type:"Uint8", ElemType:"Uint8", Name:"Status", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x8, Type:(*reflect.rtype)(0x479b7e0), Size:0, T:0x1, NumberOfBytes:0x1, StringKind:"uint8", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Status\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Uint8{
		obj:   &s.obj.Status,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// NodeAge is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_NodeState_t) NodeAge() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"NodeAge", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48b9ce0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"NodeAge\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.NodeAge == nil {
		s.obj.NodeAge = big.NewInt(0)
	}

	return &Storage_BigInt{
		obj:   s.obj.NodeAge,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// CreationEpoch is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_NodeState_t) CreationEpoch() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"CreationEpoch", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48b9ce0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"CreationEpoch\" storage:\"slot=2,offset=0\"", Slot:2, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("2", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.CreationEpoch == nil {
		s.obj.CreationEpoch = big.NewInt(0)
	}

	return &Storage_BigInt{
		obj:   s.obj.CreationEpoch,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// ActivationEpoch is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_NodeState_t) ActivationEpoch() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"ActivationEpoch", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48b9ce0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"ActivationEpoch\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("3", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.ActivationEpoch == nil {
		s.obj.ActivationEpoch = big.NewInt(0)
	}

	return &Storage_BigInt{
		obj:   s.obj.ActivationEpoch,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// ReleaseEpoch is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_NodeState_t) ReleaseEpoch() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"ReleaseEpoch", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48b9ce0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"ReleaseEpoch\" storage:\"slot=4,offset=0\"", Slot:4, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("4", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.ReleaseEpoch == nil {
		s.obj.ReleaseEpoch = big.NewInt(0)
	}

	return &Storage_BigInt{
		obj:   s.obj.ReleaseEpoch,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Amount is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_PendingDelegation_t) Amount() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Amount", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48b9ce0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Amount\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Amount == nil {
		s.obj.Amount = big.NewInt(0)
	}

	return &Storage_BigInt{
		obj:   s.obj.Amount,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Epoch is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_PendingDelegation_t) Epoch() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Epoch", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48b9ce0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Epoch\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Epoch == nil {
		s.obj.Epoch = big.NewInt(0)
	}

	return &Storage_BigInt{
		obj:   s.obj.Epoch,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// ValidatorAddress is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_RedelegationReference_t) ValidatorAddress() *Storage_Address {
	// Field: &wrapper.tmplField{Type:"Address", ElemType:"Address", Name:"ValidatorAddress", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x11, Type:(*reflect.rtype)(0x485c1e0), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"ValidatorAddress\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Address{
		obj:   &s.obj.ValidatorAddress,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// ReleasedTotalDelegation is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_RedelegationReference_t) ReleasedTotalDelegation() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"ReleasedTotalDelegation", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48b9ce0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"ReleasedTotalDelegation\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.ReleasedTotalDelegation == nil {
		s.obj.ReleasedTotalDelegation = big.NewInt(0)
	}

	return &Storage_BigInt{
		obj:   s.obj.ReleasedTotalDelegation,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Amount is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Undelegation_t) Amount() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Amount", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48b9ce0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Amount\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Amount == nil {
		s.obj.Amount = big.NewInt(0)
	}

	return &Storage_BigInt{
		obj:   s.obj.Amount,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Epoch is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Undelegation_t) Epoch() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Epoch", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48b9ce0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Epoch\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Epoch == nil {
		s.obj.Epoch = big.NewInt(0)
	}

	return &Storage_BigInt{
		obj:   s.obj.Epoch,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

