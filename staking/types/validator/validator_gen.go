// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package validator

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
	obj   *Address
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
	obj   BigInt
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
	obj   Bool
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

type Decimal = common.Dec
type Storage_Decimal struct {
	obj   *Decimal
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func (s *Storage_Decimal) Value() common.Dec {
	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	*s.obj = common.NewDecFromBigIntWithPrec(hash.Big(), common.Precision)
	return *s.obj
}

func (s *Storage_Decimal) SetValue(value common.Dec) {
	hash := value.BigInt()
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(hash))
	*s.obj = value
}

type String = string
type Storage_String struct {
	obj   *String
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

type Uint8 = uint8
type Storage_Uint8 struct {
	obj   *Uint8
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

// Array_aaa956d410fd is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x5, Name:"Array_aaa956d410fd", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0003fa7e0)}, Type:"[48]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc0004c3b00), Kind:0x11, Type:(*reflect.rtype)(0xc000166b40), Size:48, T:0x5, NumberOfBytes:0x40, StringKind:"t_array(t_bytes1)48_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Array_aaa956d410fd [48]Uint8

type Storage_Array_aaa956d410fd struct {
	obj   *Array_aaa956d410fd
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func (s *Storage_Array_aaa956d410fd) Length() int {
	return len(s.obj)
}

func (s *Storage_Array_aaa956d410fd) Get(index uint64) *Storage_Uint8 {
	// Value: &wrapper.tmplField{Type:"Uint8", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x8, Type:(*reflect.rtype)(0x47c71c0), Size:0, T:0x1, NumberOfBytes:0x1, StringKind:"uint8", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"", Slot:<nil>, Offset:0}
	actual := big.NewInt(0).Add(s.slot, big.NewInt(0).SetUint64(index*(1/32)))

	return &Storage_Uint8{
		obj:   &s.obj[index],
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Map_219baf077902 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xd, Name:"Map_219baf077902", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0003faa80), (*wrapper.tmplField)(0xc000016c40)}, Type:"map[Address]*ValidatorWrapper_", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc000166cd0), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_address,t_struct(ValidatorWrapper_)101_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004dda70), (*abi.Type)(0xc0004ddb00)}, TupleRawNames:[]string(nil)}}
type Map_219baf077902 map[Address]*ValidatorWrapper_

type Storage_Map_219baf077902 struct {
	obj   Map_219baf077902
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func (s *Storage_Map_219baf077902) Get(key Address) *Storage_ValidatorWrapper_ {
	// Key:	&wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x11, Type:(*reflect.rtype)(0x488c320), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"ValidatorWrapper_", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0xc0000abcc0), Size:0, T:0x6, NumberOfBytes:0x2e0, StringKind:"t_struct(ValidatorWrapper_)101_storage", TupleRawName:"ValidatorWrapper_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004ddb90), (*abi.Type)(0xc0004ea5a0), (*abi.Type)(0xc0004ea7e0), (*abi.Type)(0xc0004ea870), (*abi.Type)(0xc0004ea900), (*abi.Type)(0xc0004ea990)}, TupleRawNames:[]string{"Validator", "Redelegations", "Counters", "BlockReward", "TotalDelegation", "TotalDelegationByOperator"}}, Tag:"", Slot:<nil>, Offset:0}
	keyBytes := key.Hash().Bytes()
	hash := crypto.Keccak256Hash(append(keyBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()
	if s.obj[key] == nil {

		instance := new(ValidatorWrapper_)

		s.obj[key] = instance
	}

	return &Storage_ValidatorWrapper_{
		obj:   s.obj[key],
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Map_3783ff483121 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xd, Name:"Map_3783ff483121", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0003fae00), (*wrapper.tmplField)(0xc0003faee0)}, Type:"map[Address]*Bool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc000166be0), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_address,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004dc3f0), (*abi.Type)(0xc0004dc480)}, TupleRawNames:[]string(nil)}}
type Map_3783ff483121 map[Address]*Bool

type Storage_Map_3783ff483121 struct {
	obj   Map_3783ff483121
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func (s *Storage_Map_3783ff483121) Get(key Address) *Storage_Bool {
	// Key:	&wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x11, Type:(*reflect.rtype)(0x488c320), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"Bool", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x479ba20), Size:0, T:0x2, NumberOfBytes:0x1, StringKind:"t_bool", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"", Slot:<nil>, Offset:0}
	keyBytes := key.Hash().Bytes()
	hash := crypto.Keccak256Hash(append(keyBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()
	if s.obj[key] == nil {

		instance := new(Bool)

		s.obj[key] = instance
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
// &wrapper.tmplStruct{T:0xd, Name:"Map_624bcc0df495", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000016d20), (*wrapper.tmplField)(0xc000016e00)}, Type:"map[String]*Bool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc000166af0), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_string_memory_ptr,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004c2e10), (*abi.Type)(0xc0004c2f30)}, TupleRawNames:[]string(nil)}}
type Map_624bcc0df495 map[String]*Bool

type Storage_Map_624bcc0df495 struct {
	obj   Map_624bcc0df495
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func (s *Storage_Map_624bcc0df495) Get(key String) *Storage_Bool {
	// Key:	&wrapper.tmplField{Type:"String", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x18, Type:(*reflect.rtype)(0x47c6f00), Size:0, T:0x3, NumberOfBytes:0x20, StringKind:"t_string_memory_ptr", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"Bool", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x479ba20), Size:0, T:0x2, NumberOfBytes:0x1, StringKind:"t_bool", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"", Slot:<nil>, Offset:0}
	keyBytes := []byte(key)
	hash := crypto.Keccak256Hash(append(keyBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()
	if s.obj[key] == nil {

		instance := new(Bool)

		s.obj[key] = instance
	}

	return &Storage_Bool{
		obj:   s.obj[key],
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Map_f644f6b1160a is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xd, Name:"Map_f644f6b1160a", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000016000), (*wrapper.tmplField)(0xc000016460)}, Type:"map[Address]*Redelegation_", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc000166c30), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_address,t_struct(Redelegation_)80_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004dcf30), (*abi.Type)(0xc0004dcfc0)}, TupleRawNames:[]string(nil)}}
type Map_f644f6b1160a map[Address]*Redelegation_

type Storage_Map_f644f6b1160a struct {
	obj   Map_f644f6b1160a
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func (s *Storage_Map_f644f6b1160a) Get(key Address) *Storage_Redelegation_ {
	// Key:	&wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x11, Type:(*reflect.rtype)(0x488c320), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"Redelegation_", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0xc0000abb00), Size:0, T:0x6, NumberOfBytes:0xa0, StringKind:"t_struct(Redelegation_)80_storage", TupleRawName:"Redelegation_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004dc6c0), (*abi.Type)(0xc0004dc750), (*abi.Type)(0xc0004dc7e0), (*abi.Type)(0xc0004dc870)}, TupleRawNames:[]string{"DelegatorAddress", "Amount", "Reward", "Undelegation"}}, Tag:"", Slot:<nil>, Offset:0}
	keyBytes := key.Hash().Bytes()
	hash := crypto.Keccak256Hash(append(keyBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()
	if s.obj[key] == nil {

		instance := new(Redelegation_)

		s.obj[key] = instance
	}

	return &Storage_Redelegation_{
		obj:   s.obj[key],
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Slice_6bb0ade2b43c is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x4, Name:"Slice_6bb0ade2b43c", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0003fa9a0)}, Type:"[]*BLSPublicKey_", SolKind:abi.Type{Elem:(*abi.Type)(0xc0004ea090), Kind:0x17, Type:(*reflect.rtype)(0x47b3820), Size:0, T:0x4, NumberOfBytes:0x20, StringKind:"t_array(t_struct(BLSPublicKey_)32_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Slice_6bb0ade2b43c []*BLSPublicKey_

type Storage_Slice_6bb0ade2b43c struct {
	obj   *Slice_6bb0ade2b43c
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func (s *Storage_Slice_6bb0ade2b43c) Length() *big.Int {
	rv := s.db.GetState(s.addr, common.BigToHash(s.slot))
	return rv.Big()
}

func (s *Storage_Slice_6bb0ade2b43c) Resize(length uint64) {
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(big.NewInt(0).SetUint64(length)))

	slice := make([]*BLSPublicKey_, length, length+50)
	copy(slice, *s.obj)
	*s.obj = slice
}

func (s *Storage_Slice_6bb0ade2b43c) Get(index uint64) *Storage_BLSPublicKey_ {
	// Value: &wrapper.tmplField{Type:"BLSPublicKey_", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0xc0000abb80), Size:0, T:0x6, NumberOfBytes:0x40, StringKind:"t_struct(BLSPublicKey_)32_storage", TupleRawName:"BLSPublicKey_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004c3a70)}, TupleRawNames:[]string{"Key"}}, Tag:"", Slot:<nil>, Offset:0}
	length := s.Length().Uint64()
	if index >= length {
		s.Resize(index + 1)
	}

	hash := crypto.Keccak256Hash(common.BigToHash(s.slot).Bytes())
	actual := big.NewInt(0).Add(hash.Big(), big.NewInt(0).SetUint64(index*(64/32)))
	if (*s.obj)[index] == nil {

		instance := new(BLSPublicKey_)

		(*s.obj)[index] = instance
	}

	return &Storage_BLSPublicKey_{
		obj:   (*s.obj)[index],
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Slice_760032d9d704 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x4, Name:"Slice_760032d9d704", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0003fac40)}, Type:"[]*Address", SolKind:abi.Type{Elem:(*abi.Type)(0xc0004dd8c0), Kind:0x17, Type:(*reflect.rtype)(0x47b3820), Size:0, T:0x4, NumberOfBytes:0x20, StringKind:"t_array(t_address)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Slice_760032d9d704 []*Address

type Storage_Slice_760032d9d704 struct {
	obj   *Slice_760032d9d704
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func (s *Storage_Slice_760032d9d704) Length() *big.Int {
	rv := s.db.GetState(s.addr, common.BigToHash(s.slot))
	return rv.Big()
}

func (s *Storage_Slice_760032d9d704) Resize(length uint64) {
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(big.NewInt(0).SetUint64(length)))

	slice := make([]*Address, length, length+50)
	copy(slice, *s.obj)
	*s.obj = slice
}

func (s *Storage_Slice_760032d9d704) Get(index uint64) *Storage_Address {
	// Value: &wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48a4640), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"", Slot:<nil>, Offset:0}
	length := s.Length().Uint64()
	if index >= length {
		s.Resize(index + 1)
	}

	hash := crypto.Keccak256Hash(common.BigToHash(s.slot).Bytes())
	actual := big.NewInt(0).Add(hash.Big(), big.NewInt(0).SetUint64(index*(20/32)))
	if (*s.obj)[index] == nil {

		instance := new(Address)

		(*s.obj)[index] = instance
	}

	return &Storage_Address{
		obj:   (*s.obj)[index],
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Slice_d5f808851da6 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x4, Name:"Slice_d5f808851da6", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000176c0)}, Type:"[]*Slot_", SolKind:abi.Type{Elem:(*abi.Type)(0xc0004c3830), Kind:0x17, Type:(*reflect.rtype)(0x47b3820), Size:0, T:0x4, NumberOfBytes:0x20, StringKind:"t_array(t_struct(Slot_)116_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Slice_d5f808851da6 []*Slot_

type Storage_Slice_d5f808851da6 struct {
	obj   *Slice_d5f808851da6
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

func (s *Storage_Slice_d5f808851da6) Length() *big.Int {
	rv := s.db.GetState(s.addr, common.BigToHash(s.slot))
	return rv.Big()
}

func (s *Storage_Slice_d5f808851da6) Resize(length uint64) {
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(big.NewInt(0).SetUint64(length)))

	slice := make([]*Slot_, length, length+50)
	copy(slice, *s.obj)
	*s.obj = slice
}

func (s *Storage_Slice_d5f808851da6) Get(index uint64) *Storage_Slot_ {
	// Value: &wrapper.tmplField{Type:"Slot_", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0xc0000aba80), Size:0, T:0x6, NumberOfBytes:0x80, StringKind:"t_struct(Slot_)116_storage", TupleRawName:"Slot_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004c38c0), (*abi.Type)(0xc0004c39e0), (*abi.Type)(0xc0004c3d40)}, TupleRawNames:[]string{"EcdsaAddress", "BLSPublicKey", "EffectiveStake"}}, Tag:"", Slot:<nil>, Offset:0}
	length := s.Length().Uint64()
	if index >= length {
		s.Resize(index + 1)
	}

	hash := crypto.Keccak256Hash(common.BigToHash(s.slot).Bytes())
	actual := big.NewInt(0).Add(hash.Big(), big.NewInt(0).SetUint64(index*(128/32)))
	if (*s.obj)[index] == nil {

		instance := new(Slot_)

		(*s.obj)[index] = instance
	}

	return &Storage_Slot_{
		obj:   (*s.obj)[index],
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// AddressSet_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"AddressSet_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0003fad20), (*wrapper.tmplField)(0xc0003fa000)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type AddressSet_ struct {
	Keys Slice_760032d9d704 `json:"Keys" storage:"slot=0,offset=0"`
	Set  Map_3783ff483121   `json:"Set" storage:"slot=1,offset=0"`
}

// BLSPublicKey_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"BLSPublicKey_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0003fa8c0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type BLSPublicKey_ struct {
	Key Array_aaa956d410fd `json:"Key" storage:"slot=0,offset=0"`
}

// BLSPublicKeys_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"BLSPublicKeys_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0003fa460)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type BLSPublicKeys_ struct {
	Keys Slice_6bb0ade2b43c `json:"Keys" storage:"slot=0,offset=0"`
}

// CommissionRates_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"CommissionRates_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0003fb260), (*wrapper.tmplField)(0xc0003fb340), (*wrapper.tmplField)(0xc0003fb420)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type CommissionRates_ struct {
	Rate          Decimal `json:"Rate" storage:"slot=0,offset=0"`
	MaxRate       Decimal `json:"MaxRate" storage:"slot=1,offset=0"`
	MaxChangeRate Decimal `json:"MaxChangeRate" storage:"slot=2,offset=0"`
}

// Commission_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Commission_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0003fb500), (*wrapper.tmplField)(0xc0003fb5e0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Commission_ struct {
	CommissionRates CommissionRates_ `json:"CommissionRates" storage:"slot=0,offset=0"`
	UpdateHeight    BigInt           `json:"UpdateHeight" storage:"slot=3,offset=0"`
}

// Committee_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Committee_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000017340), (*wrapper.tmplField)(0xc000017880)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Committee_ struct {
	Epoch BigInt `json:"Epoch" storage:"slot=0,offset=0"`
	Slots Slots_ `json:"Slots" storage:"slot=1,offset=0"`
}

// Counters_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Counters_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000016700), (*wrapper.tmplField)(0xc0000167e0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Counters_ struct {
	NumBlocksToSign BigInt `json:"NumBlocksToSign" storage:"slot=0,offset=0"`
	NumBlocksSigned BigInt `json:"NumBlocksSigned" storage:"slot=1,offset=0"`
}

// Description_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Description_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0003fb7a0), (*wrapper.tmplField)(0xc0003fb880), (*wrapper.tmplField)(0xc0003fb960), (*wrapper.tmplField)(0xc0003fba40), (*wrapper.tmplField)(0xc0003fbb20)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Description_ struct {
	Name             String `json:"Name" storage:"slot=0,offset=0"`
	Identity         String `json:"Identity" storage:"slot=1,offset=0"`
	WebSite          String `json:"WebSite" storage:"slot=2,offset=0"`
	SecurityContract String `json:"SecurityContract" storage:"slot=3,offset=0"`
	Details          String `json:"Details" storage:"slot=4,offset=0"`
}

// Global_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Global_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000017a40)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Global_t struct {
	ValidatorPool ValidatorPool_ `json:"ValidatorPool" storage:"slot=0,offset=0"`
}

// RedelegationMap_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"RedelegationMap_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0003fbea0), (*wrapper.tmplField)(0xc000016540)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type RedelegationMap_ struct {
	Keys Slice_760032d9d704 `json:"Keys" storage:"slot=0,offset=0"`
	Map  Map_f644f6b1160a   `json:"Map" storage:"slot=1,offset=0"`
}

// Redelegation_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Redelegation_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000160e0), (*wrapper.tmplField)(0xc0000161c0), (*wrapper.tmplField)(0xc0000162a0), (*wrapper.tmplField)(0xc000016380)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Redelegation_ struct {
	DelegatorAddress Address       `json:"DelegatorAddress" storage:"slot=0,offset=0"`
	Amount           BigInt        `json:"Amount" storage:"slot=1,offset=0"`
	Reward           BigInt        `json:"Reward" storage:"slot=2,offset=0"`
	Undelegation     Undelegation_ `json:"Undelegation" storage:"slot=3,offset=0"`
}

// Slot_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Slot_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000017420), (*wrapper.tmplField)(0xc000017500), (*wrapper.tmplField)(0xc0000175e0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Slot_ struct {
	EcdsaAddress   Address       `json:"EcdsaAddress" storage:"slot=0,offset=0"`
	BLSPublicKey   BLSPublicKey_ `json:"BLSPublicKey" storage:"slot=1,offset=0"`
	EffectiveStake Decimal       `json:"EffectiveStake" storage:"slot=3,offset=0"`
}

// Slots_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Slots_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000177a0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Slots_ struct {
	Entrys Slice_d5f808851da6 `json:"Entrys" storage:"slot=0,offset=0"`
}

// Undelegation_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Undelegation_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0003fa620), (*wrapper.tmplField)(0xc0003fa700)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Undelegation_ struct {
	Amount BigInt `json:"Amount" storage:"slot=0,offset=0"`
	Epoch  BigInt `json:"Epoch" storage:"slot=1,offset=0"`
}

// ValidatorPool_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"ValidatorPool_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000170a0), (*wrapper.tmplField)(0xc000017180), (*wrapper.tmplField)(0xc000017260), (*wrapper.tmplField)(0xc000017960)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type ValidatorPool_ struct {
	Validators             ValidatorWrapperMap_ `json:"Validators" storage:"slot=0,offset=0"`
	PublicKeySet           Map_624bcc0df495     `json:"PublicKeySet" storage:"slot=2,offset=0"`
	DescriptionIdentitySet Map_624bcc0df495     `json:"DescriptionIdentitySet" storage:"slot=3,offset=0"`
	Committee              Committee_           `json:"Committee" storage:"slot=4,offset=0"`
}

// ValidatorWrapperMap_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"ValidatorWrapperMap_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000016ee0), (*wrapper.tmplField)(0xc000016fc0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type ValidatorWrapperMap_ struct {
	Keys Slice_760032d9d704 `json:"Keys" storage:"slot=0,offset=0"`
	Map  Map_219baf077902   `json:"Map" storage:"slot=1,offset=0"`
}

// ValidatorWrapper_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"ValidatorWrapper_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0003fbdc0), (*wrapper.tmplField)(0xc000016620), (*wrapper.tmplField)(0xc0000168c0), (*wrapper.tmplField)(0xc0000169a0), (*wrapper.tmplField)(0xc000016a80), (*wrapper.tmplField)(0xc000016b60)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type ValidatorWrapper_ struct {
	Validator                 Validator_       `json:"Validator" storage:"slot=0,offset=0"`
	Redelegations             RedelegationMap_ `json:"Redelegations" storage:"slot=16,offset=0"`
	Counters                  Counters_        `json:"Counters" storage:"slot=18,offset=0"`
	BlockReward               BigInt           `json:"BlockReward" storage:"slot=20,offset=0"`
	TotalDelegation           BigInt           `json:"TotalDelegation" storage:"slot=21,offset=0"`
	TotalDelegationByOperator BigInt           `json:"TotalDelegationByOperator" storage:"slot=22,offset=0"`
}

// Validator_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Validator_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0003fab60), (*wrapper.tmplField)(0xc0003fa0e0), (*wrapper.tmplField)(0xc0003fafc0), (*wrapper.tmplField)(0xc0003fb0a0), (*wrapper.tmplField)(0xc0003fb180), (*wrapper.tmplField)(0xc0003fb6c0), (*wrapper.tmplField)(0xc0003fbc00), (*wrapper.tmplField)(0xc0003fbce0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Validator_ struct {
	ValidatorAddress     Address        `json:"ValidatorAddress" storage:"slot=0,offset=0"`
	OperatorAddresses    AddressSet_    `json:"OperatorAddresses" storage:"slot=1,offset=0"`
	SlotPubKeys          BLSPublicKeys_ `json:"SlotPubKeys" storage:"slot=3,offset=0"`
	LastEpochInCommittee BigInt         `json:"LastEpochInCommittee" storage:"slot=4,offset=0"`
	Status               BigInt         `json:"Status" storage:"slot=5,offset=0"`
	Commission           Commission_    `json:"Commission" storage:"slot=6,offset=0"`
	Description          Description_   `json:"Description" storage:"slot=10,offset=0"`
	CreationHeight       BigInt         `json:"CreationHeight" storage:"slot=15,offset=0"`
}

type Storage_AddressSet_ struct {
	obj   *AddressSet_
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

type Storage_BLSPublicKey_ struct {
	obj   *BLSPublicKey_
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

type Storage_BLSPublicKeys_ struct {
	obj   *BLSPublicKeys_
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

type Storage_CommissionRates_ struct {
	obj   *CommissionRates_
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

type Storage_Commission_ struct {
	obj   *Commission_
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

type Storage_Committee_ struct {
	obj   *Committee_
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

type Storage_Counters_ struct {
	obj   *Counters_
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

type Storage_Description_ struct {
	obj   *Description_
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

type Storage_RedelegationMap_ struct {
	obj   *RedelegationMap_
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

type Storage_Redelegation_ struct {
	obj   *Redelegation_
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

type Storage_Slot_ struct {
	obj   *Slot_
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

type Storage_Slots_ struct {
	obj   *Slots_
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

type Storage_Undelegation_ struct {
	obj   *Undelegation_
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

type Storage_ValidatorPool_ struct {
	obj   *ValidatorPool_
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

type Storage_ValidatorWrapperMap_ struct {
	obj   *ValidatorWrapperMap_
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

type Storage_ValidatorWrapper_ struct {
	obj   *ValidatorWrapper_
	db    *state.StateDB
	addr  common.Address
	slot  *big.Int
	dirty StateValues
}

type Storage_Validator_ struct {
	obj   *Validator_
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

// Keys is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_AddressSet_) Keys() *Storage_Slice_760032d9d704 {
	// Field: &wrapper.tmplField{Type:"Slice_760032d9d704", ElemType:"Slice_760032d9d704", Name:"Keys", SolKind:abi.Type{Elem:(*abi.Type)(0xc0004dd8c0), Kind:0x17, Type:(*reflect.rtype)(0x47b3820), Size:0, T:0x4, NumberOfBytes:0x20, StringKind:"t_array(t_address)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Keys\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Keys == nil {

		hash := s.db.GetState(s.addr, common.BigToHash(actual))
		instance := make(Slice_760032d9d704, hash.Big().Int64())

		s.obj.Keys = instance
	}

	return &Storage_Slice_760032d9d704{
		obj:   &s.obj.Keys,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Set is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_AddressSet_) Set() *Storage_Map_3783ff483121 {
	// Field: &wrapper.tmplField{Type:"Map_3783ff483121", ElemType:"Map_3783ff483121", Name:"Set", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc000166be0), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_address,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004dc3f0), (*abi.Type)(0xc0004dc480)}, TupleRawNames:[]string(nil)}, Tag:"json:\"Set\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Set == nil {

		instance := make(Map_3783ff483121)

		s.obj.Set = instance
	}

	return &Storage_Map_3783ff483121{
		obj:   s.obj.Set,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Key is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_BLSPublicKey_) Key() *Storage_Array_aaa956d410fd {
	// Field: &wrapper.tmplField{Type:"Array_aaa956d410fd", ElemType:"Array_aaa956d410fd", Name:"Key", SolKind:abi.Type{Elem:(*abi.Type)(0xc0004c3b00), Kind:0x11, Type:(*reflect.rtype)(0xc000166b40), Size:48, T:0x5, NumberOfBytes:0x40, StringKind:"t_array(t_bytes1)48_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Key\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Array_aaa956d410fd{
		obj:   &s.obj.Key,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Keys is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_BLSPublicKeys_) Keys() *Storage_Slice_6bb0ade2b43c {
	// Field: &wrapper.tmplField{Type:"Slice_6bb0ade2b43c", ElemType:"Slice_6bb0ade2b43c", Name:"Keys", SolKind:abi.Type{Elem:(*abi.Type)(0xc0004ea090), Kind:0x17, Type:(*reflect.rtype)(0x47b3820), Size:0, T:0x4, NumberOfBytes:0x20, StringKind:"t_array(t_struct(BLSPublicKey_)32_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Keys\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Keys == nil {

		hash := s.db.GetState(s.addr, common.BigToHash(actual))
		instance := make(Slice_6bb0ade2b43c, hash.Big().Int64())

		s.obj.Keys = instance
	}

	return &Storage_Slice_6bb0ade2b43c{
		obj:   &s.obj.Keys,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Rate is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_CommissionRates_) Rate() *Storage_Decimal {
	// Field: &wrapper.tmplField{Type:"Decimal", ElemType:"Decimal", Name:"Rate", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000163bc0), Size:0, T:0x6, NumberOfBytes:0x20, StringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004c3dd0)}, TupleRawNames:[]string{"f"}}, Tag:"json:\"Rate\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Decimal{
		obj:   &s.obj.Rate,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// MaxRate is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_CommissionRates_) MaxRate() *Storage_Decimal {
	// Field: &wrapper.tmplField{Type:"Decimal", ElemType:"Decimal", Name:"MaxRate", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000163bc0), Size:0, T:0x6, NumberOfBytes:0x20, StringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004c3dd0)}, TupleRawNames:[]string{"f"}}, Tag:"json:\"MaxRate\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Decimal{
		obj:   &s.obj.MaxRate,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// MaxChangeRate is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_CommissionRates_) MaxChangeRate() *Storage_Decimal {
	// Field: &wrapper.tmplField{Type:"Decimal", ElemType:"Decimal", Name:"MaxChangeRate", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000163bc0), Size:0, T:0x6, NumberOfBytes:0x20, StringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004c3dd0)}, TupleRawNames:[]string{"f"}}, Tag:"json:\"MaxChangeRate\" storage:\"slot=2,offset=0\"", Slot:2, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("2", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Decimal{
		obj:   &s.obj.MaxChangeRate,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// CommissionRates is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Commission_) CommissionRates() *Storage_CommissionRates_ {
	// Field: &wrapper.tmplField{Type:"CommissionRates_", ElemType:"CommissionRates_", Name:"CommissionRates", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc0004e2180), Size:0, T:0x6, NumberOfBytes:0x60, StringKind:"t_struct(CommissionRates_)22_storage", TupleRawName:"CommissionRates_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004dd290), (*abi.Type)(0xc0004dd320), (*abi.Type)(0xc0004dd3b0)}, TupleRawNames:[]string{"Rate", "MaxRate", "MaxChangeRate"}}, Tag:"json:\"CommissionRates\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_CommissionRates_{
		obj:   &s.obj.CommissionRates,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// UpdateHeight is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Commission_) UpdateHeight() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"UpdateHeight", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48f48e0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"UpdateHeight\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("3", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.UpdateHeight == nil {

		instance := big.NewInt(0)

		s.obj.UpdateHeight = instance
	}

	return &Storage_BigInt{
		obj:   s.obj.UpdateHeight,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Epoch is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Committee_) Epoch() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Epoch", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48f48e0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Epoch\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Epoch == nil {

		instance := big.NewInt(0)

		s.obj.Epoch = instance
	}

	return &Storage_BigInt{
		obj:   s.obj.Epoch,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Slots is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Committee_) Slots() *Storage_Slots_ {
	// Field: &wrapper.tmplField{Type:"Slots_", ElemType:"Slots_", Name:"Slots", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000163c80), Size:0, T:0x6, NumberOfBytes:0x20, StringKind:"t_struct(Slots_)120_storage", TupleRawName:"Slots_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004c37a0)}, TupleRawNames:[]string{"Entrys"}}, Tag:"json:\"Slots\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Slots_{
		obj:   &s.obj.Slots,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// NumBlocksToSign is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Counters_) NumBlocksToSign() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"NumBlocksToSign", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48f48e0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"NumBlocksToSign\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.NumBlocksToSign == nil {

		instance := big.NewInt(0)

		s.obj.NumBlocksToSign = instance
	}

	return &Storage_BigInt{
		obj:   s.obj.NumBlocksToSign,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// NumBlocksSigned is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Counters_) NumBlocksSigned() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"NumBlocksSigned", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48f48e0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"NumBlocksSigned\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.NumBlocksSigned == nil {

		instance := big.NewInt(0)

		s.obj.NumBlocksSigned = instance
	}

	return &Storage_BigInt{
		obj:   s.obj.NumBlocksSigned,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Name is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Description_) Name() *Storage_String {
	// Field: &wrapper.tmplField{Type:"String", ElemType:"String", Name:"Name", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x18, Type:(*reflect.rtype)(0x47c6f00), Size:0, T:0x3, NumberOfBytes:0x20, StringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Name\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
func (s *Storage_Description_) Identity() *Storage_String {
	// Field: &wrapper.tmplField{Type:"String", ElemType:"String", Name:"Identity", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x18, Type:(*reflect.rtype)(0x47c6f00), Size:0, T:0x3, NumberOfBytes:0x20, StringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Identity\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
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
func (s *Storage_Description_) WebSite() *Storage_String {
	// Field: &wrapper.tmplField{Type:"String", ElemType:"String", Name:"WebSite", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x18, Type:(*reflect.rtype)(0x47c6f00), Size:0, T:0x3, NumberOfBytes:0x20, StringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"WebSite\" storage:\"slot=2,offset=0\"", Slot:2, Offset:0}
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
func (s *Storage_Description_) SecurityContract() *Storage_String {
	// Field: &wrapper.tmplField{Type:"String", ElemType:"String", Name:"SecurityContract", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x18, Type:(*reflect.rtype)(0x47c6f00), Size:0, T:0x3, NumberOfBytes:0x20, StringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"SecurityContract\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
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
func (s *Storage_Description_) Details() *Storage_String {
	// Field: &wrapper.tmplField{Type:"String", ElemType:"String", Name:"Details", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x18, Type:(*reflect.rtype)(0x47c6f00), Size:0, T:0x3, NumberOfBytes:0x20, StringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Details\" storage:\"slot=4,offset=0\"", Slot:4, Offset:0}
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

// ValidatorPool is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Global_t) ValidatorPool() *Storage_ValidatorPool_ {
	// Field: &wrapper.tmplField{Type:"ValidatorPool_", ElemType:"ValidatorPool_", Name:"ValidatorPool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc0004e2900), Size:0, T:0x6, NumberOfBytes:0xc0, StringKind:"t_struct(ValidatorPool_)138_storage", TupleRawName:"ValidatorPool_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004dd7a0), (*abi.Type)(0xc0004eac60), (*abi.Type)(0xc0004eacf0), (*abi.Type)(0xc0004ead80)}, TupleRawNames:[]string{"Validators", "PublicKeySet", "DescriptionIdentitySet", "Committee"}}, Tag:"json:\"ValidatorPool\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_ValidatorPool_{
		obj:   &s.obj.ValidatorPool,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Keys is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_RedelegationMap_) Keys() *Storage_Slice_760032d9d704 {
	// Field: &wrapper.tmplField{Type:"Slice_760032d9d704", ElemType:"Slice_760032d9d704", Name:"Keys", SolKind:abi.Type{Elem:(*abi.Type)(0xc0004dd8c0), Kind:0x17, Type:(*reflect.rtype)(0x47b3820), Size:0, T:0x4, NumberOfBytes:0x20, StringKind:"t_array(t_address)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Keys\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Keys == nil {

		hash := s.db.GetState(s.addr, common.BigToHash(actual))
		instance := make(Slice_760032d9d704, hash.Big().Int64())

		s.obj.Keys = instance
	}

	return &Storage_Slice_760032d9d704{
		obj:   &s.obj.Keys,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Map is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_RedelegationMap_) Map() *Storage_Map_f644f6b1160a {
	// Field: &wrapper.tmplField{Type:"Map_f644f6b1160a", ElemType:"Map_f644f6b1160a", Name:"Map", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc000166c30), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_address,t_struct(Redelegation_)80_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004dcf30), (*abi.Type)(0xc0004dcfc0)}, TupleRawNames:[]string(nil)}, Tag:"json:\"Map\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Map == nil {

		instance := make(Map_f644f6b1160a)

		s.obj.Map = instance
	}

	return &Storage_Map_f644f6b1160a{
		obj:   s.obj.Map,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// DelegatorAddress is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Redelegation_) DelegatorAddress() *Storage_Address {
	// Field: &wrapper.tmplField{Type:"Address", ElemType:"Address", Name:"DelegatorAddress", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x11, Type:(*reflect.rtype)(0x488c320), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"DelegatorAddress\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
func (s *Storage_Redelegation_) Amount() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Amount", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48f48e0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Amount\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Amount == nil {

		instance := big.NewInt(0)

		s.obj.Amount = instance
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
func (s *Storage_Redelegation_) Reward() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Reward", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48f48e0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Reward\" storage:\"slot=2,offset=0\"", Slot:2, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("2", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Reward == nil {

		instance := big.NewInt(0)

		s.obj.Reward = instance
	}

	return &Storage_BigInt{
		obj:   s.obj.Reward,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Undelegation is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Redelegation_) Undelegation() *Storage_Undelegation_ {
	// Field: &wrapper.tmplField{Type:"Undelegation_", ElemType:"Undelegation_", Name:"Undelegation", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000163e60), Size:0, T:0x6, NumberOfBytes:0x40, StringKind:"t_struct(Undelegation_)71_storage", TupleRawName:"Undelegation_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004dc900), (*abi.Type)(0xc0004dc990)}, TupleRawNames:[]string{"Amount", "Epoch"}}, Tag:"json:\"Undelegation\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("3", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Undelegation_{
		obj:   &s.obj.Undelegation,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// EcdsaAddress is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Slot_) EcdsaAddress() *Storage_Address {
	// Field: &wrapper.tmplField{Type:"Address", ElemType:"Address", Name:"EcdsaAddress", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x11, Type:(*reflect.rtype)(0x488c320), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"EcdsaAddress\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Address{
		obj:   &s.obj.EcdsaAddress,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// BLSPublicKey is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Slot_) BLSPublicKey() *Storage_BLSPublicKey_ {
	// Field: &wrapper.tmplField{Type:"BLSPublicKey_", ElemType:"BLSPublicKey_", Name:"BLSPublicKey", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000163b00), Size:0, T:0x6, NumberOfBytes:0x40, StringKind:"t_struct(BLSPublicKey_)32_storage", TupleRawName:"BLSPublicKey_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004c3a70)}, TupleRawNames:[]string{"Key"}}, Tag:"json:\"BLSPublicKey\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_BLSPublicKey_{
		obj:   &s.obj.BLSPublicKey,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// EffectiveStake is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Slot_) EffectiveStake() *Storage_Decimal {
	// Field: &wrapper.tmplField{Type:"Decimal", ElemType:"Decimal", Name:"EffectiveStake", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000163bc0), Size:0, T:0x6, NumberOfBytes:0x20, StringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004c3dd0)}, TupleRawNames:[]string{"f"}}, Tag:"json:\"EffectiveStake\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("3", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Decimal{
		obj:   &s.obj.EffectiveStake,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Entrys is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Slots_) Entrys() *Storage_Slice_d5f808851da6 {
	// Field: &wrapper.tmplField{Type:"Slice_d5f808851da6", ElemType:"Slice_d5f808851da6", Name:"Entrys", SolKind:abi.Type{Elem:(*abi.Type)(0xc0004c3830), Kind:0x17, Type:(*reflect.rtype)(0x47b3820), Size:0, T:0x4, NumberOfBytes:0x20, StringKind:"t_array(t_struct(Slot_)116_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Entrys\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Entrys == nil {

		hash := s.db.GetState(s.addr, common.BigToHash(actual))
		instance := make(Slice_d5f808851da6, hash.Big().Int64())

		s.obj.Entrys = instance
	}

	return &Storage_Slice_d5f808851da6{
		obj:   &s.obj.Entrys,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Amount is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Undelegation_) Amount() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Amount", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48f48e0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Amount\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Amount == nil {

		instance := big.NewInt(0)

		s.obj.Amount = instance
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
func (s *Storage_Undelegation_) Epoch() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Epoch", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48f48e0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Epoch\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Epoch == nil {

		instance := big.NewInt(0)

		s.obj.Epoch = instance
	}

	return &Storage_BigInt{
		obj:   s.obj.Epoch,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Validators is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_ValidatorPool_) Validators() *Storage_ValidatorWrapperMap_ {
	// Field: &wrapper.tmplField{Type:"ValidatorWrapperMap_", ElemType:"ValidatorWrapperMap_", Name:"Validators", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc0004e2840), Size:0, T:0x6, NumberOfBytes:0x40, StringKind:"t_struct(ValidatorWrapperMap_)109_storage", TupleRawName:"ValidatorWrapperMap_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004dd830), (*abi.Type)(0xc0004dd9e0)}, TupleRawNames:[]string{"Keys", "Map"}}, Tag:"json:\"Validators\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_ValidatorWrapperMap_{
		obj:   &s.obj.Validators,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// PublicKeySet is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_ValidatorPool_) PublicKeySet() *Storage_Map_624bcc0df495 {
	// Field: &wrapper.tmplField{Type:"Map_624bcc0df495", ElemType:"Map_624bcc0df495", Name:"PublicKeySet", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc000166af0), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_string_memory_ptr,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004c2e10), (*abi.Type)(0xc0004c2f30)}, TupleRawNames:[]string(nil)}, Tag:"json:\"PublicKeySet\" storage:\"slot=2,offset=0\"", Slot:2, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("2", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.PublicKeySet == nil {

		instance := make(Map_624bcc0df495)

		s.obj.PublicKeySet = instance
	}

	return &Storage_Map_624bcc0df495{
		obj:   s.obj.PublicKeySet,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// DescriptionIdentitySet is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_ValidatorPool_) DescriptionIdentitySet() *Storage_Map_624bcc0df495 {
	// Field: &wrapper.tmplField{Type:"Map_624bcc0df495", ElemType:"Map_624bcc0df495", Name:"DescriptionIdentitySet", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc000166af0), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_string_memory_ptr,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004c2e10), (*abi.Type)(0xc0004c2f30)}, TupleRawNames:[]string(nil)}, Tag:"json:\"DescriptionIdentitySet\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("3", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.DescriptionIdentitySet == nil {

		instance := make(Map_624bcc0df495)

		s.obj.DescriptionIdentitySet = instance
	}

	return &Storage_Map_624bcc0df495{
		obj:   s.obj.DescriptionIdentitySet,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Committee is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_ValidatorPool_) Committee() *Storage_Committee_ {
	// Field: &wrapper.tmplField{Type:"Committee_", ElemType:"Committee_", Name:"Committee", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000163ce0), Size:0, T:0x6, NumberOfBytes:0x40, StringKind:"t_struct(Committee_)125_storage", TupleRawName:"Committee_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004c35f0), (*abi.Type)(0xc0004c3710)}, TupleRawNames:[]string{"Epoch", "Slots"}}, Tag:"json:\"Committee\" storage:\"slot=4,offset=0\"", Slot:4, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("4", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Committee_{
		obj:   &s.obj.Committee,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Keys is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_ValidatorWrapperMap_) Keys() *Storage_Slice_760032d9d704 {
	// Field: &wrapper.tmplField{Type:"Slice_760032d9d704", ElemType:"Slice_760032d9d704", Name:"Keys", SolKind:abi.Type{Elem:(*abi.Type)(0xc0004dd8c0), Kind:0x17, Type:(*reflect.rtype)(0x47b3820), Size:0, T:0x4, NumberOfBytes:0x20, StringKind:"t_array(t_address)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Keys\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Keys == nil {

		hash := s.db.GetState(s.addr, common.BigToHash(actual))
		instance := make(Slice_760032d9d704, hash.Big().Int64())

		s.obj.Keys = instance
	}

	return &Storage_Slice_760032d9d704{
		obj:   &s.obj.Keys,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Map is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_ValidatorWrapperMap_) Map() *Storage_Map_219baf077902 {
	// Field: &wrapper.tmplField{Type:"Map_219baf077902", ElemType:"Map_219baf077902", Name:"Map", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc000166cd0), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_address,t_struct(ValidatorWrapper_)101_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004dda70), (*abi.Type)(0xc0004ddb00)}, TupleRawNames:[]string(nil)}, Tag:"json:\"Map\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Map == nil {

		instance := make(Map_219baf077902)

		s.obj.Map = instance
	}

	return &Storage_Map_219baf077902{
		obj:   s.obj.Map,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Validator is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_ValidatorWrapper_) Validator() *Storage_Validator_ {
	// Field: &wrapper.tmplField{Type:"Validator_", ElemType:"Validator_", Name:"Validator", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc0004e26c0), Size:0, T:0x6, NumberOfBytes:0x200, StringKind:"t_struct(Validator_)66_storage", TupleRawName:"Validator_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004ddc20), (*abi.Type)(0xc0004ddcb0), (*abi.Type)(0xc0004ddef0), (*abi.Type)(0xc0004ea240), (*abi.Type)(0xc0004ea2d0), (*abi.Type)(0xc0004ea360), (*abi.Type)(0xc0004ea3f0), (*abi.Type)(0xc0004ea480)}, TupleRawNames:[]string{"ValidatorAddress", "OperatorAddresses", "SlotPubKeys", "LastEpochInCommittee", "Status", "Commission", "Description", "CreationHeight"}}, Tag:"json:\"Validator\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Validator_{
		obj:   &s.obj.Validator,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Redelegations is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_ValidatorWrapper_) Redelegations() *Storage_RedelegationMap_ {
	// Field: &wrapper.tmplField{Type:"RedelegationMap_", ElemType:"RedelegationMap_", Name:"Redelegations", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc0004e2780), Size:0, T:0x6, NumberOfBytes:0x40, StringKind:"t_struct(RedelegationMap_)88_storage", TupleRawName:"RedelegationMap_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004ea630), (*abi.Type)(0xc0004ea6c0)}, TupleRawNames:[]string{"Keys", "Map"}}, Tag:"json:\"Redelegations\" storage:\"slot=16,offset=0\"", Slot:16, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("16", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_RedelegationMap_{
		obj:   &s.obj.Redelegations,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Counters is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_ValidatorWrapper_) Counters() *Storage_Counters_ {
	// Field: &wrapper.tmplField{Type:"Counters_", ElemType:"Counters_", Name:"Counters", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc0004e2000), Size:0, T:0x6, NumberOfBytes:0x40, StringKind:"t_struct(Counters_)41_storage", TupleRawName:"Counters_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004dcc60), (*abi.Type)(0xc0004dccf0)}, TupleRawNames:[]string{"NumBlocksToSign", "NumBlocksSigned"}}, Tag:"json:\"Counters\" storage:\"slot=18,offset=0\"", Slot:18, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("18", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Counters_{
		obj:   &s.obj.Counters,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// BlockReward is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_ValidatorWrapper_) BlockReward() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"BlockReward", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48f48e0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"BlockReward\" storage:\"slot=20,offset=0\"", Slot:20, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("20", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.BlockReward == nil {

		instance := big.NewInt(0)

		s.obj.BlockReward = instance
	}

	return &Storage_BigInt{
		obj:   s.obj.BlockReward,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// TotalDelegation is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_ValidatorWrapper_) TotalDelegation() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"TotalDelegation", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48f48e0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"TotalDelegation\" storage:\"slot=21,offset=0\"", Slot:21, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("21", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.TotalDelegation == nil {

		instance := big.NewInt(0)

		s.obj.TotalDelegation = instance
	}

	return &Storage_BigInt{
		obj:   s.obj.TotalDelegation,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// TotalDelegationByOperator is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_ValidatorWrapper_) TotalDelegationByOperator() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"TotalDelegationByOperator", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48f48e0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"TotalDelegationByOperator\" storage:\"slot=22,offset=0\"", Slot:22, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("22", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.TotalDelegationByOperator == nil {

		instance := big.NewInt(0)

		s.obj.TotalDelegationByOperator = instance
	}

	return &Storage_BigInt{
		obj:   s.obj.TotalDelegationByOperator,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// ValidatorAddress is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Validator_) ValidatorAddress() *Storage_Address {
	// Field: &wrapper.tmplField{Type:"Address", ElemType:"Address", Name:"ValidatorAddress", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x11, Type:(*reflect.rtype)(0x488c320), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"ValidatorAddress\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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

// OperatorAddresses is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Validator_) OperatorAddresses() *Storage_AddressSet_ {
	// Field: &wrapper.tmplField{Type:"AddressSet_", ElemType:"AddressSet_", Name:"OperatorAddresses", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc0004e2540), Size:0, T:0x6, NumberOfBytes:0x40, StringKind:"t_struct(AddressSet_)49_storage", TupleRawName:"AddressSet_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004ddd40), (*abi.Type)(0xc0004dddd0)}, TupleRawNames:[]string{"Keys", "Set"}}, Tag:"json:\"OperatorAddresses\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_AddressSet_{
		obj:   &s.obj.OperatorAddresses,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// SlotPubKeys is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Validator_) SlotPubKeys() *Storage_BLSPublicKeys_ {
	// Field: &wrapper.tmplField{Type:"BLSPublicKeys_", ElemType:"BLSPublicKeys_", Name:"SlotPubKeys", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc0004e2660), Size:0, T:0x6, NumberOfBytes:0x20, StringKind:"t_struct(BLSPublicKeys_)36_storage", TupleRawName:"BLSPublicKeys_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004ea000)}, TupleRawNames:[]string{"Keys"}}, Tag:"json:\"SlotPubKeys\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("3", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_BLSPublicKeys_{
		obj:   &s.obj.SlotPubKeys,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// LastEpochInCommittee is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Validator_) LastEpochInCommittee() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"LastEpochInCommittee", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48f48e0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"LastEpochInCommittee\" storage:\"slot=4,offset=0\"", Slot:4, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("4", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.LastEpochInCommittee == nil {

		instance := big.NewInt(0)

		s.obj.LastEpochInCommittee = instance
	}

	return &Storage_BigInt{
		obj:   s.obj.LastEpochInCommittee,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Status is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Validator_) Status() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Status", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48f48e0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Status\" storage:\"slot=5,offset=0\"", Slot:5, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("5", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Status == nil {

		instance := big.NewInt(0)

		s.obj.Status = instance
	}

	return &Storage_BigInt{
		obj:   s.obj.Status,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Commission is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Validator_) Commission() *Storage_Commission_ {
	// Field: &wrapper.tmplField{Type:"Commission_", ElemType:"Commission_", Name:"Commission", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc0004e21e0), Size:0, T:0x6, NumberOfBytes:0x80, StringKind:"t_struct(Commission_)27_storage", TupleRawName:"Commission_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004dd200), (*abi.Type)(0xc0004dd4d0)}, TupleRawNames:[]string{"CommissionRates", "UpdateHeight"}}, Tag:"json:\"Commission\" storage:\"slot=6,offset=0\"", Slot:6, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("6", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Commission_{
		obj:   &s.obj.Commission,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// Description is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Validator_) Description() *Storage_Description_ {
	// Field: &wrapper.tmplField{Type:"Description_", ElemType:"Description_", Name:"Description", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000163740), Size:0, T:0x6, NumberOfBytes:0xa0, StringKind:"t_struct(Description_)15_storage", TupleRawName:"Description_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0004c3170), (*abi.Type)(0xc0004c3290), (*abi.Type)(0xc0004c3320), (*abi.Type)(0xc0004c33b0), (*abi.Type)(0xc0004c3440)}, TupleRawNames:[]string{"Name", "Identity", "WebSite", "SecurityContract", "Details"}}, Tag:"json:\"Description\" storage:\"slot=10,offset=0\"", Slot:10, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("10", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Description_{
		obj:   &s.obj.Description,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

// CreationHeight is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Validator_) CreationHeight() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"CreationHeight", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48f48e0), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"CreationHeight\" storage:\"slot=15,offset=0\"", Slot:15, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("15", 10)

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.CreationHeight == nil {

		instance := big.NewInt(0)

		s.obj.CreationHeight = instance
	}

	return &Storage_BigInt{
		obj:   s.obj.CreationHeight,
		db:    s.db,
		addr:  s.addr,
		slot:  actual,
		dirty: s.dirty,
	}
}

