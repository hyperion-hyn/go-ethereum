// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package restaking

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

	// &wrapper.tmplStruct{T:0xf, Name:"Decimal", Fields:[]*wrapper.tmplField(nil), Type:"common.Dec", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xf, numberOfBytes:0x20, stringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fd100)}, TupleRawNames:[]string{"f"}, TupleType:(*reflect.rtype)(0xc0004387e0)}}
	obj           *Decimal
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Decimal) Value() common.Dec {
	// &wrapper.tmplStruct{T:0xf, Name:"Decimal", Fields:[]*wrapper.tmplField(nil), Type:"common.Dec", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xf, numberOfBytes:0x20, stringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fd100)}, TupleRawNames:[]string{"f"}, TupleType:(*reflect.rtype)(0xc0004387e0)}}
	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	*s.obj = common.NewDecFromBigIntWithPrec(hash.Big(), common.Precision)
	return *s.obj
}

func (s *Storage_Decimal) SetValue(value common.Dec) {
	// &wrapper.tmplStruct{T:0xf, Name:"Decimal", Fields:[]*wrapper.tmplField(nil), Type:"common.Dec", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xf, numberOfBytes:0x20, stringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fd100)}, TupleRawNames:[]string{"f"}, TupleType:(*reflect.rtype)(0xc0004387e0)}}
	hash := value.BigInt()
	if hash != nil {
		s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(hash))
	}
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
// &wrapper.tmplStruct{T:0x5, Name:"Array_aaa956d410fd", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000b81a0)}, Type:"[48]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000fcc80), Size:48, T:0x5, numberOfBytes:0x40, stringKind:"t_array(t_bytes1)48_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "[48]uint8"
type Array_aaa956d410fd [48]Uint8

type Storage_Array_aaa956d410fd struct {

	// &wrapper.tmplStruct{T:0x5, Name:"Array_aaa956d410fd", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000b81a0)}, Type:"[48]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000fcc80), Size:48, T:0x5, numberOfBytes:0x40, stringKind:"t_array(t_bytes1)48_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           *Array_aaa956d410fd
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Array_aaa956d410fd) Value() [48]Uint8 {
	// &wrapper.tmplStruct{T:0x5, Name:"Array_aaa956d410fd", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000b81a0)}, Type:"[48]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000fcc80), Size:48, T:0x5, numberOfBytes:0x40, stringKind:"t_array(t_bytes1)48_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
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
	// &wrapper.tmplStruct{T:0x5, Name:"Array_aaa956d410fd", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000b81a0)}, Type:"[48]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000fcc80), Size:48, T:0x5, numberOfBytes:0x40, stringKind:"t_array(t_bytes1)48_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
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

// Map_11675a09c025 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xe, Name:"Map_11675a09c025", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0003885b0), (*wrapper.tmplField)(0xc000388820)}, Type:"map[Address]*ValidatorWrapperMapEntry_", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(ValidatorWrapperMapEntry_)113_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048c000), (*abi.Type)(0xc000156680)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "map[common.Address]*struct { Entry struct { Validator struct { ValidatorAddress common.Address \"json:\\\"ValidatorAddress\\\" storage:\\\"slot=0,offset=0\\\"\"; OperatorAddresses struct { Keys []*common.Address \"json:\\\"Keys\\\" storage:\\\"slot=0,offset=0\\\"\"; Set map[common.Address]*bool \"json:\\\"Set\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"OperatorAddresses\\\" storage:\\\"slot=1,offset=0\\\"\"; SlotPubKeys struct { Keys []*struct { Key [48]uint8 \"json:\\\"Key\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"Keys\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"SlotPubKeys\\\" storage:\\\"slot=3,offset=0\\\"\"; LastEpochInCommittee *big.Int \"json:\\\"LastEpochInCommittee\\\" storage:\\\"slot=4,offset=0\\\"\"; MaxTotalDelegation *big.Int \"json:\\\"MaxTotalDelegation\\\" storage:\\\"slot=5,offset=0\\\"\"; Status uint8 \"json:\\\"Status\\\" storage:\\\"slot=6,offset=0\\\"\"; Commission struct { CommissionRates struct { Rate struct { F *big.Int \"json:\\\"f\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"Rate\\\" storage:\\\"slot=0,offset=0\\\"\"; MaxRate struct { F *big.Int \"json:\\\"f\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"MaxRate\\\" storage:\\\"slot=1,offset=0\\\"\"; MaxChangeRate struct { F *big.Int \"json:\\\"f\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"MaxChangeRate\\\" storage:\\\"slot=2,offset=0\\\"\" } \"json:\\\"CommissionRates\\\" storage:\\\"slot=0,offset=0\\\"\"; UpdateHeight *big.Int \"json:\\\"UpdateHeight\\\" storage:\\\"slot=3,offset=0\\\"\" } \"json:\\\"Commission\\\" storage:\\\"slot=7,offset=0\\\"\"; Description struct { Name string \"json:\\\"Name\\\" storage:\\\"slot=0,offset=0\\\"\"; Identity string \"json:\\\"Identity\\\" storage:\\\"slot=1,offset=0\\\"\"; Website string \"json:\\\"Website\\\" storage:\\\"slot=2,offset=0\\\"\"; SecurityContact string \"json:\\\"SecurityContact\\\" storage:\\\"slot=3,offset=0\\\"\"; Details string \"json:\\\"Details\\\" storage:\\\"slot=4,offset=0\\\"\" } \"json:\\\"Description\\\" storage:\\\"slot=11,offset=0\\\"\"; CreationHeight *big.Int \"json:\\\"CreationHeight\\\" storage:\\\"slot=16,offset=0\\\"\" } \"json:\\\"Validator\\\" storage:\\\"slot=0,offset=0\\\"\"; Redelegations struct { Keys []*common.Address \"json:\\\"Keys\\\" storage:\\\"slot=0,offset=0\\\"\"; Map map[common.Address]*struct { Entry struct { DelegatorAddress common.Address \"json:\\\"DelegatorAddress\\\" storage:\\\"slot=0,offset=0\\\"\"; Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=1,offset=0\\\"\"; Reward *big.Int \"json:\\\"Reward\\\" storage:\\\"slot=2,offset=0\\\"\"; Undelegation struct { Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=0,offset=0\\\"\"; Epoch *big.Int \"json:\\\"Epoch\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"Undelegation\\\" storage:\\\"slot=3,offset=0\\\"\" } \"json:\\\"Entry\\\" storage:\\\"slot=0,offset=0\\\"\"; Index *big.Int \"json:\\\"Index\\\" storage:\\\"slot=5,offset=0\\\"\" } \"json:\\\"Map\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"Redelegations\\\" storage:\\\"slot=17,offset=0\\\"\"; Counters struct { NumBlocksToSign *big.Int \"json:\\\"NumBlocksToSign\\\" storage:\\\"slot=0,offset=0\\\"\"; NumBlocksSigned *big.Int \"json:\\\"NumBlocksSigned\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"Counters\\\" storage:\\\"slot=19,offset=0\\\"\"; BlockReward *big.Int \"json:\\\"BlockReward\\\" storage:\\\"slot=21,offset=0\\\"\"; TotalDelegation *big.Int \"json:\\\"TotalDelegation\\\" storage:\\\"slot=22,offset=0\\\"\"; TotalDelegationByOperator *big.Int \"json:\\\"TotalDelegationByOperator\\\" storage:\\\"slot=23,offset=0\\\"\" } \"json:\\\"Entry\\\" storage:\\\"slot=0,offset=0\\\"\"; Index *big.Int \"json:\\\"Index\\\" storage:\\\"slot=24,offset=0\\\"\" }"
type Map_11675a09c025 map[Address]*ValidatorWrapperMapEntry_

type Storage_Map_11675a09c025 struct {

	// &wrapper.tmplStruct{T:0xe, Name:"Map_11675a09c025", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0003885b0), (*wrapper.tmplField)(0xc000388820)}, Type:"map[Address]*ValidatorWrapperMapEntry_", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(ValidatorWrapperMapEntry_)113_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048c000), (*abi.Type)(0xc000156680)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           Map_11675a09c025
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Map_11675a09c025) Get(key Address) *Storage_ValidatorWrapperMapEntry_ {
	// Key:	&wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"ValidatorWrapperMapEntry_", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048c080), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	keyBytes := key.Hash().Bytes()
	hash := crypto.Keccak256Hash(append(keyBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()
	offset := 0
	if s.obj[key] == nil {

		instance := new(ValidatorWrapperMapEntry_)

		s.obj[key] = instance
	}

	return &Storage_ValidatorWrapperMapEntry_{
		obj:           s.obj[key],
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 800,
		dirty:         s.dirty,
	}
}

// Map_1e3b3171b974 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xe, Name:"Map_1e3b3171b974", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012c680), (*wrapper.tmplField)(0xc00012cea0)}, Type:"map[Address]*RedelegationMapEntry_", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(RedelegationMapEntry_)87_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fd780), (*abi.Type)(0xc0000fdc80)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "map[common.Address]*struct { Entry struct { DelegatorAddress common.Address \"json:\\\"DelegatorAddress\\\" storage:\\\"slot=0,offset=0\\\"\"; Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=1,offset=0\\\"\"; Reward *big.Int \"json:\\\"Reward\\\" storage:\\\"slot=2,offset=0\\\"\"; Undelegation struct { Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=0,offset=0\\\"\"; Epoch *big.Int \"json:\\\"Epoch\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"Undelegation\\\" storage:\\\"slot=3,offset=0\\\"\" } \"json:\\\"Entry\\\" storage:\\\"slot=0,offset=0\\\"\"; Index *big.Int \"json:\\\"Index\\\" storage:\\\"slot=5,offset=0\\\"\" }"
type Map_1e3b3171b974 map[Address]*RedelegationMapEntry_

type Storage_Map_1e3b3171b974 struct {

	// &wrapper.tmplStruct{T:0xe, Name:"Map_1e3b3171b974", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012c680), (*wrapper.tmplField)(0xc00012cea0)}, Type:"map[Address]*RedelegationMapEntry_", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(RedelegationMapEntry_)87_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fd780), (*abi.Type)(0xc0000fdc80)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           Map_1e3b3171b974
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Map_1e3b3171b974) Get(key Address) *Storage_RedelegationMapEntry_ {
	// Key:	&wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"RedelegationMapEntry_", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000fd800), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	keyBytes := key.Hash().Bytes()
	hash := crypto.Keccak256Hash(append(keyBytes, common.BigToHash(s.slot).Bytes()...))
	actual := hash.Big()
	offset := 0
	if s.obj[key] == nil {

		instance := new(RedelegationMapEntry_)

		s.obj[key] = instance
	}

	return &Storage_RedelegationMapEntry_{
		obj:           s.obj[key],
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 192,
		dirty:         s.dirty,
	}
}

// Map_3783ff483121 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xe, Name:"Map_3783ff483121", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000b85b0), (*wrapper.tmplField)(0xc0000b89c0)}, Type:"map[Address]*Bool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048c400), (*abi.Type)(0xc00048c580)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "map[common.Address]*bool"
type Map_3783ff483121 map[Address]*Bool

type Storage_Map_3783ff483121 struct {

	// &wrapper.tmplStruct{T:0xe, Name:"Map_3783ff483121", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000b85b0), (*wrapper.tmplField)(0xc0000b89c0)}, Type:"map[Address]*Bool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048c400), (*abi.Type)(0xc00048c580)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
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
	// Value: &wrapper.tmplField{Type:"Bool", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc00048c480), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
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
// &wrapper.tmplStruct{T:0xe, Name:"Map_624bcc0df495", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000388410), (*wrapper.tmplField)(0xc0003884e0)}, Type:"map[String]*Bool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_string_memory_ptr,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc000156880), (*abi.Type)(0xc000156a00)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "map[string]*bool"
type Map_624bcc0df495 map[String]*Bool

type Storage_Map_624bcc0df495 struct {

	// &wrapper.tmplStruct{T:0xe, Name:"Map_624bcc0df495", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000388410), (*wrapper.tmplField)(0xc0003884e0)}, Type:"map[String]*Bool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_string_memory_ptr,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc000156880), (*abi.Type)(0xc000156a00)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
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
	// Value: &wrapper.tmplField{Type:"Bool", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc000156980), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
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

// Slice_6bb0ade2b43c is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x4, Name:"Slice_6bb0ade2b43c", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000b8c30)}, Type:"[]*BLSPublicKey_", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000fce80), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_struct(BLSPublicKey_)32_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "[]*struct { Key [48]uint8 \"json:\\\"Key\\\" storage:\\\"slot=0,offset=0\\\"\" }"
type Slice_6bb0ade2b43c []*BLSPublicKey_

type Storage_Slice_6bb0ade2b43c struct {

	// &wrapper.tmplStruct{T:0x4, Name:"Slice_6bb0ade2b43c", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000b8c30)}, Type:"[]*BLSPublicKey_", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000fce80), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_struct(BLSPublicKey_)32_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           *Slice_6bb0ade2b43c
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Slice_6bb0ade2b43c) Length() int {
	rv := s.db.GetState(s.addr, common.BigToHash(s.slot))
	return int(rv.Big().Int64())
}

func (s *Storage_Slice_6bb0ade2b43c) Resize(length int) {
	// Value: &wrapper.tmplField{Type:"BLSPublicKey_", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000fcb80), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(big.NewInt(0).SetUint64(uint64(length))))

	_length := length
	_capacity := length + 50

	slice := make([]*BLSPublicKey_, _length, _capacity)

	copy(slice, *s.obj)
	*s.obj = slice
}

func (s *Storage_Slice_6bb0ade2b43c) Get(index int) *Storage_BLSPublicKey_ {
	// Value: &wrapper.tmplField{Type:"BLSPublicKey_", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000fcb80), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
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
// &wrapper.tmplStruct{T:0x4, Name:"Slice_760032d9d704", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000b8410)}, Type:"[]*Address", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000fd600), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_address)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "[]*common.Address"
type Slice_760032d9d704 []*Address

type Storage_Slice_760032d9d704 struct {

	// &wrapper.tmplStruct{T:0x4, Name:"Slice_760032d9d704", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000b8410)}, Type:"[]*Address", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000fd600), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_address)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
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
	// Value: &wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000fd500), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(big.NewInt(0).SetUint64(uint64(length))))

	_length := length
	_capacity := length + 50

	slice := make([]*Address, _length, _capacity)

	copy(slice, *s.obj)
	*s.obj = slice
}

func (s *Storage_Slice_760032d9d704) Get(index int) *Storage_Address {
	// Value: &wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000fd500), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
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

// Slice_ed073f3f1ff1 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x4, Name:"Slice_ed073f3f1ff1", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000388270)}, Type:"[]*Slot_", SolKind:abi.Type{Elem:(*abi.Type)(0xc000157000), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_struct(Slot_)128_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "[]*struct { EcdsaAddress common.Address \"json:\\\"EcdsaAddress\\\" storage:\\\"slot=0,offset=0\\\"\"; BLSPublicKey struct { Key [48]uint8 \"json:\\\"Key\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"BLSPublicKey\\\" storage:\\\"slot=1,offset=0\\\"\"; EffectiveStake struct { F *big.Int \"json:\\\"f\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"EffectiveStake\\\" storage:\\\"slot=3,offset=0\\\"\" }"
type Slice_ed073f3f1ff1 []*Slot_

type Storage_Slice_ed073f3f1ff1 struct {

	// &wrapper.tmplStruct{T:0x4, Name:"Slice_ed073f3f1ff1", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000388270)}, Type:"[]*Slot_", SolKind:abi.Type{Elem:(*abi.Type)(0xc000157000), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_struct(Slot_)128_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           *Slice_ed073f3f1ff1
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Slice_ed073f3f1ff1) Length() int {
	rv := s.db.GetState(s.addr, common.BigToHash(s.slot))
	return int(rv.Big().Int64())
}

func (s *Storage_Slice_ed073f3f1ff1) Resize(length int) {
	// Value: &wrapper.tmplField{Type:"Slot_", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc000156d80), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(big.NewInt(0).SetUint64(uint64(length))))

	_length := length
	_capacity := length + 50

	slice := make([]*Slot_, _length, _capacity)

	copy(slice, *s.obj)
	*s.obj = slice
}

func (s *Storage_Slice_ed073f3f1ff1) Get(index int) *Storage_Slot_ {
	// Value: &wrapper.tmplField{Type:"Slot_", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc000156d80), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	length := s.Length()
	if length <= index {
		s.Resize(index + 1)
	}

	base := crypto.Keccak256Hash(common.BigToHash(s.slot).Bytes()).Big()

	slotsPerItem := (128 + 31) / 32
	actual := big.NewInt(0).Add(base, big.NewInt(0).SetUint64(uint64(index*slotsPerItem)))
	offset := 0

	if (*s.obj)[index] == nil {

		instance := new(Slot_)

		(*s.obj)[index] = instance
	}

	return &Storage_Slot_{
		obj:           (*s.obj)[index],
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 128,
		dirty:         s.dirty,
	}
}

// AddressSet_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"AddressSet_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000b84e0), (*wrapper.tmplField)(0xc0000b8a90)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(AddressSet_)49_storage", TupleRawName:"AddressSet_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048c300), (*abi.Type)(0xc00048c380)}, TupleRawNames:[]string{"Keys", "Set"}, TupleType:(*reflect.rtype)(0xc000439080)}}
type AddressSet_ struct {
	Keys Slice_760032d9d704 `json:"Keys" storage:"slot=0,offset=0"`
	Set  Map_3783ff483121   `json:"Set" storage:"slot=1,offset=0"`
}

// BLSPublicKey_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"BLSPublicKey_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000b8270)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(BLSPublicKey_)32_storage", TupleRawName:"BLSPublicKey_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fcc00)}, TupleRawNames:[]string{"Key"}, TupleType:(*reflect.rtype)(0xc000438600)}}
type BLSPublicKey_ struct {
	Key Array_aaa956d410fd `json:"Key" storage:"slot=0,offset=0"`
}

// BLSPublicKeys_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"BLSPublicKeys_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000b8dd0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(BLSPublicKeys_)36_storage", TupleRawName:"BLSPublicKeys_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fcb00)}, TupleRawNames:[]string{"Keys"}, TupleType:(*reflect.rtype)(0xc000438660)}}
type BLSPublicKeys_ struct {
	Keys Slice_6bb0ade2b43c `json:"Keys" storage:"slot=0,offset=0"`
}

// CommissionRates_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"CommissionRates_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000b95f0), (*wrapper.tmplField)(0xc0000b96c0), (*wrapper.tmplField)(0xc0000b9790)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x60, stringKind:"t_struct(CommissionRates_)22_storage", TupleRawName:"CommissionRates_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fd080), (*abi.Type)(0xc0000fd200), (*abi.Type)(0xc0000fd280)}, TupleRawNames:[]string{"Rate", "MaxRate", "MaxChangeRate"}, TupleType:(*reflect.rtype)(0xc000438840)}}
type CommissionRates_ struct {
	Rate          Decimal `json:"Rate" storage:"slot=0,offset=0"`
	MaxRate       Decimal `json:"MaxRate" storage:"slot=1,offset=0"`
	MaxChangeRate Decimal `json:"MaxChangeRate" storage:"slot=2,offset=0"`
}

// Commission_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Commission_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000b9860), (*wrapper.tmplField)(0xc0000b9930)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x80, stringKind:"t_struct(Commission_)27_storage", TupleRawName:"Commission_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fd000), (*abi.Type)(0xc0000fd380)}, TupleRawNames:[]string{"CommissionRates", "UpdateHeight"}, TupleType:(*reflect.rtype)(0xc0004388a0)}}
type Commission_ struct {
	CommissionRates CommissionRates_ `json:"CommissionRates" storage:"slot=0,offset=0"`
	UpdateHeight    BigInt           `json:"UpdateHeight" storage:"slot=3,offset=0"`
}

// Committee_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Committee_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0003891e0), (*wrapper.tmplField)(0xc0003892b0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Committee_)137_storage", TupleRawName:"Committee_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000156c00), (*abi.Type)(0xc000156c80)}, TupleRawNames:[]string{"Epoch", "Slots"}, TupleType:(*reflect.rtype)(0xc000438000)}}
type Committee_ struct {
	Epoch BigInt `json:"Epoch" storage:"slot=0,offset=0"`
	Slots Slots_ `json:"Slots" storage:"slot=1,offset=0"`
}

// Counters_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Counters_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012d110), (*wrapper.tmplField)(0xc00012d1e0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Counters_)41_storage", TupleRawName:"Counters_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048d180), (*abi.Type)(0xc00048d200)}, TupleRawNames:[]string{"NumBlocksToSign", "NumBlocksSigned"}, TupleType:(*reflect.rtype)(0xc0004392c0)}}
type Counters_ struct {
	NumBlocksToSign BigInt `json:"NumBlocksToSign" storage:"slot=0,offset=0"`
	NumBlocksSigned BigInt `json:"NumBlocksSigned" storage:"slot=1,offset=0"`
}

// Description_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Description_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000b9ad0), (*wrapper.tmplField)(0xc0000b9e10), (*wrapper.tmplField)(0xc00012c0d0), (*wrapper.tmplField)(0xc00012c1a0), (*wrapper.tmplField)(0xc00012c270)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xa0, stringKind:"t_struct(Description_)15_storage", TupleRawName:"Description_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048ca80), (*abi.Type)(0xc00048cb80), (*abi.Type)(0xc00048cc00), (*abi.Type)(0xc00048cc80), (*abi.Type)(0xc00048cd00)}, TupleRawNames:[]string{"Name", "Identity", "Website", "SecurityContact", "Details"}, TupleType:(*reflect.rtype)(0xc0004391a0)}}
type Description_ struct {
	Name            String `json:"Name" storage:"slot=0,offset=0"`
	Identity        String `json:"Identity" storage:"slot=1,offset=0"`
	Website         String `json:"Website" storage:"slot=2,offset=0"`
	SecurityContact String `json:"SecurityContact" storage:"slot=3,offset=0"`
	Details         String `json:"Details" storage:"slot=4,offset=0"`
}

// Global_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Global_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000389450)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x0, stringKind:"t_struct(Global_t)", TupleRawName:"Global_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fde00)}, TupleRawNames:[]string{"ValidatorPool"}, TupleType:(*reflect.rtype)(0xc000438120)}}
type Global_t struct {
	ValidatorPool ValidatorPool_ `json:"ValidatorPool" storage:"slot=0,offset=0"`
}

// RedelegationMapEntry_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"RedelegationMapEntry_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012cd00), (*wrapper.tmplField)(0xc00012cdd0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xc0, stringKind:"t_struct(RedelegationMapEntry_)87_storage", TupleRawName:"RedelegationMapEntry_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fd880), (*abi.Type)(0xc0000fdb80)}, TupleRawNames:[]string{"Entry", "Index"}, TupleType:(*reflect.rtype)(0xc000438c00)}}
type RedelegationMapEntry_ struct {
	Entry Redelegation_ `json:"Entry" storage:"slot=0,offset=0"`
	Index BigInt        `json:"Index" storage:"slot=5,offset=0"`
}

// RedelegationMap_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"RedelegationMap_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012c5b0), (*wrapper.tmplField)(0xc00012cf70)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(RedelegationMap_)95_storage", TupleRawName:"RedelegationMap_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fd480), (*abi.Type)(0xc0000fd700)}, TupleRawNames:[]string{"Keys", "Map"}, TupleType:(*reflect.rtype)(0xc000438c60)}}
type RedelegationMap_ struct {
	Keys Slice_760032d9d704 `json:"Keys" storage:"slot=0,offset=0"`
	Map  Map_1e3b3171b974   `json:"Map" storage:"slot=1,offset=0"`
}

// Redelegation_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Redelegation_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012c820), (*wrapper.tmplField)(0xc00012c8f0), (*wrapper.tmplField)(0xc00012c9c0), (*wrapper.tmplField)(0xc00012cc30)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xa0, stringKind:"t_struct(Redelegation_)82_storage", TupleRawName:"Redelegation_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fd900), (*abi.Type)(0xc0000fd980), (*abi.Type)(0xc0000fda00), (*abi.Type)(0xc0000fda80)}, TupleRawNames:[]string{"DelegatorAddress", "Amount", "Reward", "Undelegation"}, TupleType:(*reflect.rtype)(0xc000438ba0)}}
type Redelegation_ struct {
	DelegatorAddress Address       `json:"DelegatorAddress" storage:"slot=0,offset=0"`
	Amount           BigInt        `json:"Amount" storage:"slot=1,offset=0"`
	Reward           BigInt        `json:"Reward" storage:"slot=2,offset=0"`
	Undelegation     Undelegation_ `json:"Undelegation" storage:"slot=3,offset=0"`
}

// Slot_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Slot_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012dee0), (*wrapper.tmplField)(0xc0003880d0), (*wrapper.tmplField)(0xc0003881a0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x80, stringKind:"t_struct(Slot_)128_storage", TupleRawName:"Slot_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000156e00), (*abi.Type)(0xc000156e80), (*abi.Type)(0xc000156f00)}, TupleRawNames:[]string{"EcdsaAddress", "BLSPublicKey", "EffectiveStake"}, TupleType:(*reflect.rtype)(0xc0004f8120)}}
type Slot_ struct {
	EcdsaAddress   Address       `json:"EcdsaAddress" storage:"slot=0,offset=0"`
	BLSPublicKey   BLSPublicKey_ `json:"BLSPublicKey" storage:"slot=1,offset=0"`
	EffectiveStake Decimal       `json:"EffectiveStake" storage:"slot=3,offset=0"`
}

// Slots_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Slots_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000388340)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(Slots_)132_storage", TupleRawName:"Slots_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000156d00)}, TupleRawNames:[]string{"Entrys"}, TupleType:(*reflect.rtype)(0xc0004f8180)}}
type Slots_ struct {
	Entrys Slice_ed073f3f1ff1 `json:"Entrys" storage:"slot=0,offset=0"`
}

// Undelegation_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Undelegation_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012ca90), (*wrapper.tmplField)(0xc00012cb60)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Undelegation_)73_storage", TupleRawName:"Undelegation_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fc900), (*abi.Type)(0xc0000fca00)}, TupleRawNames:[]string{"Amount", "Epoch"}, TupleType:(*reflect.rtype)(0xc0004383c0)}}
type Undelegation_ struct {
	Amount BigInt `json:"Amount" storage:"slot=0,offset=0"`
	Epoch  BigInt `json:"Epoch" storage:"slot=1,offset=0"`
}

// ValidatorPool_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"ValidatorPool_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000388f70), (*wrapper.tmplField)(0xc000389040), (*wrapper.tmplField)(0xc000389110), (*wrapper.tmplField)(0xc000389380)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xc0, stringKind:"t_struct(ValidatorPool_)150_storage", TupleRawName:"ValidatorPool_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fde80), (*abi.Type)(0xc000156800), (*abi.Type)(0xc000156b00), (*abi.Type)(0xc000156b80)}, TupleRawNames:[]string{"Validators", "SlotKeySet", "DescriptionIdentitySet", "Committee"}, TupleType:(*reflect.rtype)(0xc0004380c0)}}
type ValidatorPool_ struct {
	Validators             ValidatorWrapperMap_ `json:"Validators" storage:"slot=0,offset=0"`
	SlotKeySet             Map_624bcc0df495     `json:"SlotKeySet" storage:"slot=2,offset=0"`
	DescriptionIdentitySet Map_624bcc0df495     `json:"DescriptionIdentitySet" storage:"slot=3,offset=0"`
	Committee              Committee_           `json:"Committee" storage:"slot=4,offset=0"`
}

// ValidatorWrapperMapEntry_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"ValidatorWrapperMapEntry_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012da00), (*wrapper.tmplField)(0xc00012de10)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x320, stringKind:"t_struct(ValidatorWrapperMapEntry_)113_storage", TupleRawName:"ValidatorWrapperMapEntry_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048c100), (*abi.Type)(0xc000156580)}, TupleRawNames:[]string{"Entry", "Index"}, TupleType:(*reflect.rtype)(0xc00010be00)}}
type ValidatorWrapperMapEntry_ struct {
	Entry ValidatorWrapper_ `json:"Entry" storage:"slot=0,offset=0"`
	Index BigInt            `json:"Index" storage:"slot=24,offset=0"`
}

// ValidatorWrapperMap_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"ValidatorWrapperMap_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000388a90), (*wrapper.tmplField)(0xc000388ea0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(ValidatorWrapperMap_)121_storage", TupleRawName:"ValidatorWrapperMap_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fdf00), (*abi.Type)(0xc0000fdf80)}, TupleRawNames:[]string{"Keys", "Map"}, TupleType:(*reflect.rtype)(0xc00010be60)}}
type ValidatorWrapperMap_ struct {
	Keys Slice_760032d9d704 `json:"Keys" storage:"slot=0,offset=0"`
	Map  Map_11675a09c025   `json:"Map" storage:"slot=1,offset=0"`
}

// ValidatorWrapper_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"ValidatorWrapper_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012c4e0), (*wrapper.tmplField)(0xc00012d040), (*wrapper.tmplField)(0xc00012d380), (*wrapper.tmplField)(0xc00012d520), (*wrapper.tmplField)(0xc00012d5f0), (*wrapper.tmplField)(0xc00012d6c0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x300, stringKind:"t_struct(ValidatorWrapper_)108_storage", TupleRawName:"ValidatorWrapper_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048c180), (*abi.Type)(0xc00048d080), (*abi.Type)(0xc00048d100), (*abi.Type)(0xc00048d300), (*abi.Type)(0xc00048d380), (*abi.Type)(0xc00048d480)}, TupleRawNames:[]string{"Validator", "Redelegations", "Counters", "BlockReward", "TotalDelegation", "TotalDelegationByOperator"}, TupleType:(*reflect.rtype)(0xc00010bda0)}}
type ValidatorWrapper_ struct {
	Validator                 Validator_       `json:"Validator" storage:"slot=0,offset=0"`
	Redelegations             RedelegationMap_ `json:"Redelegations" storage:"slot=17,offset=0"`
	Counters                  Counters_        `json:"Counters" storage:"slot=19,offset=0"`
	BlockReward               BigInt           `json:"BlockReward" storage:"slot=21,offset=0"`
	TotalDelegation           BigInt           `json:"TotalDelegation" storage:"slot=22,offset=0"`
	TotalDelegationByOperator BigInt           `json:"TotalDelegationByOperator" storage:"slot=23,offset=0"`
}

// Validator_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Validator_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000b8340), (*wrapper.tmplField)(0xc0000b8b60), (*wrapper.tmplField)(0xc0000b92b0), (*wrapper.tmplField)(0xc0000b9380), (*wrapper.tmplField)(0xc0000b9450), (*wrapper.tmplField)(0xc0000b9520), (*wrapper.tmplField)(0xc0000b9a00), (*wrapper.tmplField)(0xc00012c340), (*wrapper.tmplField)(0xc00012c410)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x220, stringKind:"t_struct(Validator_)68_storage", TupleRawName:"Validator_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048c200), (*abi.Type)(0xc00048c280), (*abi.Type)(0xc00048c700), (*abi.Type)(0xc00048c780), (*abi.Type)(0xc00048c800), (*abi.Type)(0xc00048c900), (*abi.Type)(0xc00048c980), (*abi.Type)(0xc00048ca00), (*abi.Type)(0xc00048cf00)}, TupleRawNames:[]string{"ValidatorAddress", "OperatorAddresses", "SlotPubKeys", "LastEpochInCommittee", "MaxTotalDelegation", "Status", "Commission", "Description", "CreationHeight"}, TupleType:(*reflect.rtype)(0xc000439200)}}
type Validator_ struct {
	ValidatorAddress     Address        `json:"ValidatorAddress" storage:"slot=0,offset=0"`
	OperatorAddresses    AddressSet_    `json:"OperatorAddresses" storage:"slot=1,offset=0"`
	SlotPubKeys          BLSPublicKeys_ `json:"SlotPubKeys" storage:"slot=3,offset=0"`
	LastEpochInCommittee BigInt         `json:"LastEpochInCommittee" storage:"slot=4,offset=0"`
	MaxTotalDelegation   BigInt         `json:"MaxTotalDelegation" storage:"slot=5,offset=0"`
	Status               Uint8          `json:"Status" storage:"slot=6,offset=0"`
	Commission           Commission_    `json:"Commission" storage:"slot=7,offset=0"`
	Description          Description_   `json:"Description" storage:"slot=11,offset=0"`
	CreationHeight       BigInt         `json:"CreationHeight" storage:"slot=16,offset=0"`
}

type Storage_AddressSet_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"AddressSet_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000b84e0), (*wrapper.tmplField)(0xc0000b8a90)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(AddressSet_)49_storage", TupleRawName:"AddressSet_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048c300), (*abi.Type)(0xc00048c380)}, TupleRawNames:[]string{"Keys", "Set"}, TupleType:(*reflect.rtype)(0xc000439080)}}
	obj           *AddressSet_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_BLSPublicKey_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"BLSPublicKey_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000b8270)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(BLSPublicKey_)32_storage", TupleRawName:"BLSPublicKey_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fcc00)}, TupleRawNames:[]string{"Key"}, TupleType:(*reflect.rtype)(0xc000438600)}}
	obj           *BLSPublicKey_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_BLSPublicKeys_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"BLSPublicKeys_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000b8dd0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(BLSPublicKeys_)36_storage", TupleRawName:"BLSPublicKeys_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fcb00)}, TupleRawNames:[]string{"Keys"}, TupleType:(*reflect.rtype)(0xc000438660)}}
	obj           *BLSPublicKeys_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_CommissionRates_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"CommissionRates_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000b95f0), (*wrapper.tmplField)(0xc0000b96c0), (*wrapper.tmplField)(0xc0000b9790)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x60, stringKind:"t_struct(CommissionRates_)22_storage", TupleRawName:"CommissionRates_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fd080), (*abi.Type)(0xc0000fd200), (*abi.Type)(0xc0000fd280)}, TupleRawNames:[]string{"Rate", "MaxRate", "MaxChangeRate"}, TupleType:(*reflect.rtype)(0xc000438840)}}
	obj           *CommissionRates_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Commission_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Commission_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000b9860), (*wrapper.tmplField)(0xc0000b9930)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x80, stringKind:"t_struct(Commission_)27_storage", TupleRawName:"Commission_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fd000), (*abi.Type)(0xc0000fd380)}, TupleRawNames:[]string{"CommissionRates", "UpdateHeight"}, TupleType:(*reflect.rtype)(0xc0004388a0)}}
	obj           *Commission_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Committee_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Committee_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0003891e0), (*wrapper.tmplField)(0xc0003892b0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Committee_)137_storage", TupleRawName:"Committee_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000156c00), (*abi.Type)(0xc000156c80)}, TupleRawNames:[]string{"Epoch", "Slots"}, TupleType:(*reflect.rtype)(0xc000438000)}}
	obj           *Committee_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Counters_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Counters_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012d110), (*wrapper.tmplField)(0xc00012d1e0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Counters_)41_storage", TupleRawName:"Counters_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048d180), (*abi.Type)(0xc00048d200)}, TupleRawNames:[]string{"NumBlocksToSign", "NumBlocksSigned"}, TupleType:(*reflect.rtype)(0xc0004392c0)}}
	obj           *Counters_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Description_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Description_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000b9ad0), (*wrapper.tmplField)(0xc0000b9e10), (*wrapper.tmplField)(0xc00012c0d0), (*wrapper.tmplField)(0xc00012c1a0), (*wrapper.tmplField)(0xc00012c270)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xa0, stringKind:"t_struct(Description_)15_storage", TupleRawName:"Description_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048ca80), (*abi.Type)(0xc00048cb80), (*abi.Type)(0xc00048cc00), (*abi.Type)(0xc00048cc80), (*abi.Type)(0xc00048cd00)}, TupleRawNames:[]string{"Name", "Identity", "Website", "SecurityContact", "Details"}, TupleType:(*reflect.rtype)(0xc0004391a0)}}
	obj           *Description_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Global_t struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Global_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000389450)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x0, stringKind:"t_struct(Global_t)", TupleRawName:"Global_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fde00)}, TupleRawNames:[]string{"ValidatorPool"}, TupleType:(*reflect.rtype)(0xc000438120)}}
	obj           *Global_t
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_RedelegationMapEntry_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"RedelegationMapEntry_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012cd00), (*wrapper.tmplField)(0xc00012cdd0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xc0, stringKind:"t_struct(RedelegationMapEntry_)87_storage", TupleRawName:"RedelegationMapEntry_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fd880), (*abi.Type)(0xc0000fdb80)}, TupleRawNames:[]string{"Entry", "Index"}, TupleType:(*reflect.rtype)(0xc000438c00)}}
	obj           *RedelegationMapEntry_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_RedelegationMap_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"RedelegationMap_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012c5b0), (*wrapper.tmplField)(0xc00012cf70)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(RedelegationMap_)95_storage", TupleRawName:"RedelegationMap_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fd480), (*abi.Type)(0xc0000fd700)}, TupleRawNames:[]string{"Keys", "Map"}, TupleType:(*reflect.rtype)(0xc000438c60)}}
	obj           *RedelegationMap_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Redelegation_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Redelegation_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012c820), (*wrapper.tmplField)(0xc00012c8f0), (*wrapper.tmplField)(0xc00012c9c0), (*wrapper.tmplField)(0xc00012cc30)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xa0, stringKind:"t_struct(Redelegation_)82_storage", TupleRawName:"Redelegation_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fd900), (*abi.Type)(0xc0000fd980), (*abi.Type)(0xc0000fda00), (*abi.Type)(0xc0000fda80)}, TupleRawNames:[]string{"DelegatorAddress", "Amount", "Reward", "Undelegation"}, TupleType:(*reflect.rtype)(0xc000438ba0)}}
	obj           *Redelegation_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Slot_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Slot_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012dee0), (*wrapper.tmplField)(0xc0003880d0), (*wrapper.tmplField)(0xc0003881a0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x80, stringKind:"t_struct(Slot_)128_storage", TupleRawName:"Slot_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000156e00), (*abi.Type)(0xc000156e80), (*abi.Type)(0xc000156f00)}, TupleRawNames:[]string{"EcdsaAddress", "BLSPublicKey", "EffectiveStake"}, TupleType:(*reflect.rtype)(0xc0004f8120)}}
	obj           *Slot_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Slots_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Slots_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000388340)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(Slots_)132_storage", TupleRawName:"Slots_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000156d00)}, TupleRawNames:[]string{"Entrys"}, TupleType:(*reflect.rtype)(0xc0004f8180)}}
	obj           *Slots_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Undelegation_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Undelegation_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012ca90), (*wrapper.tmplField)(0xc00012cb60)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Undelegation_)73_storage", TupleRawName:"Undelegation_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fc900), (*abi.Type)(0xc0000fca00)}, TupleRawNames:[]string{"Amount", "Epoch"}, TupleType:(*reflect.rtype)(0xc0004383c0)}}
	obj           *Undelegation_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_ValidatorPool_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"ValidatorPool_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000388f70), (*wrapper.tmplField)(0xc000389040), (*wrapper.tmplField)(0xc000389110), (*wrapper.tmplField)(0xc000389380)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xc0, stringKind:"t_struct(ValidatorPool_)150_storage", TupleRawName:"ValidatorPool_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fde80), (*abi.Type)(0xc000156800), (*abi.Type)(0xc000156b00), (*abi.Type)(0xc000156b80)}, TupleRawNames:[]string{"Validators", "SlotKeySet", "DescriptionIdentitySet", "Committee"}, TupleType:(*reflect.rtype)(0xc0004380c0)}}
	obj           *ValidatorPool_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_ValidatorWrapperMapEntry_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"ValidatorWrapperMapEntry_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012da00), (*wrapper.tmplField)(0xc00012de10)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x320, stringKind:"t_struct(ValidatorWrapperMapEntry_)113_storage", TupleRawName:"ValidatorWrapperMapEntry_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048c100), (*abi.Type)(0xc000156580)}, TupleRawNames:[]string{"Entry", "Index"}, TupleType:(*reflect.rtype)(0xc00010be00)}}
	obj           *ValidatorWrapperMapEntry_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_ValidatorWrapperMap_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"ValidatorWrapperMap_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000388a90), (*wrapper.tmplField)(0xc000388ea0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(ValidatorWrapperMap_)121_storage", TupleRawName:"ValidatorWrapperMap_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fdf00), (*abi.Type)(0xc0000fdf80)}, TupleRawNames:[]string{"Keys", "Map"}, TupleType:(*reflect.rtype)(0xc00010be60)}}
	obj           *ValidatorWrapperMap_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_ValidatorWrapper_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"ValidatorWrapper_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00012c4e0), (*wrapper.tmplField)(0xc00012d040), (*wrapper.tmplField)(0xc00012d380), (*wrapper.tmplField)(0xc00012d520), (*wrapper.tmplField)(0xc00012d5f0), (*wrapper.tmplField)(0xc00012d6c0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x300, stringKind:"t_struct(ValidatorWrapper_)108_storage", TupleRawName:"ValidatorWrapper_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048c180), (*abi.Type)(0xc00048d080), (*abi.Type)(0xc00048d100), (*abi.Type)(0xc00048d300), (*abi.Type)(0xc00048d380), (*abi.Type)(0xc00048d480)}, TupleRawNames:[]string{"Validator", "Redelegations", "Counters", "BlockReward", "TotalDelegation", "TotalDelegationByOperator"}, TupleType:(*reflect.rtype)(0xc00010bda0)}}
	obj           *ValidatorWrapper_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Validator_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Validator_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000b8340), (*wrapper.tmplField)(0xc0000b8b60), (*wrapper.tmplField)(0xc0000b92b0), (*wrapper.tmplField)(0xc0000b9380), (*wrapper.tmplField)(0xc0000b9450), (*wrapper.tmplField)(0xc0000b9520), (*wrapper.tmplField)(0xc0000b9a00), (*wrapper.tmplField)(0xc00012c340), (*wrapper.tmplField)(0xc00012c410)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x220, stringKind:"t_struct(Validator_)68_storage", TupleRawName:"Validator_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048c200), (*abi.Type)(0xc00048c280), (*abi.Type)(0xc00048c700), (*abi.Type)(0xc00048c780), (*abi.Type)(0xc00048c800), (*abi.Type)(0xc00048c900), (*abi.Type)(0xc00048c980), (*abi.Type)(0xc00048ca00), (*abi.Type)(0xc00048cf00)}, TupleRawNames:[]string{"ValidatorAddress", "OperatorAddresses", "SlotPubKeys", "LastEpochInCommittee", "MaxTotalDelegation", "Status", "Commission", "Description", "CreationHeight"}, TupleType:(*reflect.rtype)(0xc000439200)}}
	obj           *Validator_
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

// Keys is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_AddressSet_) Keys() *Storage_Slice_760032d9d704 {
	// Field: &wrapper.tmplField{Type:"Slice_760032d9d704", ElemType:"Slice_760032d9d704", Name:"Keys", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000fd600), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_address)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Keys\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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

// Set is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_AddressSet_) Set() *Storage_Map_3783ff483121 {
	// Field: &wrapper.tmplField{Type:"Map_3783ff483121", ElemType:"Map_3783ff483121", Name:"Set", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048c400), (*abi.Type)(0xc00048c580)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Set\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Set == nil {

		instance := make(Map_3783ff483121)

		s.obj.Set = instance
	}

	return &Storage_Map_3783ff483121{
		obj:           s.obj.Set,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Key is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_BLSPublicKey_) Key() *Storage_Array_aaa956d410fd {
	// Field: &wrapper.tmplField{Type:"Array_aaa956d410fd", ElemType:"Array_aaa956d410fd", Name:"Key", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000fcc80), Size:48, T:0x5, numberOfBytes:0x40, stringKind:"t_array(t_bytes1)48_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Key\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
func (s *Storage_BLSPublicKeys_) Keys() *Storage_Slice_6bb0ade2b43c {
	// Field: &wrapper.tmplField{Type:"Slice_6bb0ade2b43c", ElemType:"Slice_6bb0ade2b43c", Name:"Keys", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000fce80), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_struct(BLSPublicKey_)32_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Keys\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Keys == nil {

		hash := s.db.GetState(s.addr, common.BigToHash(actual))
		instance := make(Slice_6bb0ade2b43c, hash.Big().Int64())

		s.obj.Keys = instance
	}

	return &Storage_Slice_6bb0ade2b43c{
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
func (s *Storage_CommissionRates_) Rate() *Storage_Decimal {
	// Field: &wrapper.tmplField{Type:"Decimal", ElemType:"Decimal", Name:"Rate", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fd100)}, TupleRawNames:[]string{"f"}, TupleType:(*reflect.rtype)(0xc0004387e0)}, Tag:"json:\"Rate\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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

// MaxRate is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_CommissionRates_) MaxRate() *Storage_Decimal {
	// Field: &wrapper.tmplField{Type:"Decimal", ElemType:"Decimal", Name:"MaxRate", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fd100)}, TupleRawNames:[]string{"f"}, TupleType:(*reflect.rtype)(0xc0004387e0)}, Tag:"json:\"MaxRate\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
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
func (s *Storage_CommissionRates_) MaxChangeRate() *Storage_Decimal {
	// Field: &wrapper.tmplField{Type:"Decimal", ElemType:"Decimal", Name:"MaxChangeRate", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fd100)}, TupleRawNames:[]string{"f"}, TupleType:(*reflect.rtype)(0xc0004387e0)}, Tag:"json:\"MaxChangeRate\" storage:\"slot=2,offset=0\"", Slot:2, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("2", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Decimal{
		obj:           &s.obj.MaxChangeRate,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// CommissionRates is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Commission_) CommissionRates() *Storage_CommissionRates_ {
	// Field: &wrapper.tmplField{Type:"CommissionRates_", ElemType:"CommissionRates_", Name:"CommissionRates", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x60, stringKind:"t_struct(CommissionRates_)22_storage", TupleRawName:"CommissionRates_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fd080), (*abi.Type)(0xc0000fd200), (*abi.Type)(0xc0000fd280)}, TupleRawNames:[]string{"Rate", "MaxRate", "MaxChangeRate"}, TupleType:(*reflect.rtype)(0xc000438840)}, Tag:"json:\"CommissionRates\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_CommissionRates_{
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
func (s *Storage_Commission_) UpdateHeight() *Storage_BigInt {
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

// Epoch is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Committee_) Epoch() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Epoch", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Epoch\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
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

// Slots is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Committee_) Slots() *Storage_Slots_ {
	// Field: &wrapper.tmplField{Type:"Slots_", ElemType:"Slots_", Name:"Slots", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(Slots_)132_storage", TupleRawName:"Slots_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000156d00)}, TupleRawNames:[]string{"Entrys"}, TupleType:(*reflect.rtype)(0xc0004f8180)}, Tag:"json:\"Slots\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Slots_{
		obj:           &s.obj.Slots,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// NumBlocksToSign is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Counters_) NumBlocksToSign() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"NumBlocksToSign", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"NumBlocksToSign\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.NumBlocksToSign == nil {

		instance := big.NewInt(0)

		s.obj.NumBlocksToSign = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.NumBlocksToSign,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// NumBlocksSigned is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Counters_) NumBlocksSigned() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"NumBlocksSigned", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"NumBlocksSigned\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.NumBlocksSigned == nil {

		instance := big.NewInt(0)

		s.obj.NumBlocksSigned = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.NumBlocksSigned,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
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

// ValidatorPool is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Global_t) ValidatorPool() *Storage_ValidatorPool_ {
	// Field: &wrapper.tmplField{Type:"ValidatorPool_", ElemType:"ValidatorPool_", Name:"ValidatorPool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xc0, stringKind:"t_struct(ValidatorPool_)150_storage", TupleRawName:"ValidatorPool_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fde80), (*abi.Type)(0xc000156800), (*abi.Type)(0xc000156b00), (*abi.Type)(0xc000156b80)}, TupleRawNames:[]string{"Validators", "SlotKeySet", "DescriptionIdentitySet", "Committee"}, TupleType:(*reflect.rtype)(0xc0004380c0)}, Tag:"json:\"ValidatorPool\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_ValidatorPool_{
		obj:           &s.obj.ValidatorPool,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 192,
		dirty:         s.dirty,
	}
}

// Entry is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_RedelegationMapEntry_) Entry() *Storage_Redelegation_ {
	// Field: &wrapper.tmplField{Type:"Redelegation_", ElemType:"Redelegation_", Name:"Entry", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xa0, stringKind:"t_struct(Redelegation_)82_storage", TupleRawName:"Redelegation_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fd900), (*abi.Type)(0xc0000fd980), (*abi.Type)(0xc0000fda00), (*abi.Type)(0xc0000fda80)}, TupleRawNames:[]string{"DelegatorAddress", "Amount", "Reward", "Undelegation"}, TupleType:(*reflect.rtype)(0xc000438ba0)}, Tag:"json:\"Entry\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Redelegation_{
		obj:           &s.obj.Entry,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 160,
		dirty:         s.dirty,
	}
}

// Index is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_RedelegationMapEntry_) Index() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Index", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Index\" storage:\"slot=5,offset=0\"", Slot:5, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("5", 10)
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
func (s *Storage_RedelegationMap_) Keys() *Storage_Slice_760032d9d704 {
	// Field: &wrapper.tmplField{Type:"Slice_760032d9d704", ElemType:"Slice_760032d9d704", Name:"Keys", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000fd600), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_address)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Keys\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
func (s *Storage_RedelegationMap_) Map() *Storage_Map_1e3b3171b974 {
	// Field: &wrapper.tmplField{Type:"Map_1e3b3171b974", ElemType:"Map_1e3b3171b974", Name:"Map", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(RedelegationMapEntry_)87_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fd780), (*abi.Type)(0xc0000fdc80)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Map\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Map == nil {

		instance := make(Map_1e3b3171b974)

		s.obj.Map = instance
	}

	return &Storage_Map_1e3b3171b974{
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
func (s *Storage_Redelegation_) DelegatorAddress() *Storage_Address {
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
func (s *Storage_Redelegation_) Amount() *Storage_BigInt {
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
func (s *Storage_Redelegation_) Reward() *Storage_BigInt {
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

// Undelegation is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Redelegation_) Undelegation() *Storage_Undelegation_ {
	// Field: &wrapper.tmplField{Type:"Undelegation_", ElemType:"Undelegation_", Name:"Undelegation", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Undelegation_)73_storage", TupleRawName:"Undelegation_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fc900), (*abi.Type)(0xc0000fca00)}, TupleRawNames:[]string{"Amount", "Epoch"}, TupleType:(*reflect.rtype)(0xc0004383c0)}, Tag:"json:\"Undelegation\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("3", 10)
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

// EcdsaAddress is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Slot_) EcdsaAddress() *Storage_Address {
	// Field: &wrapper.tmplField{Type:"Address", ElemType:"Address", Name:"EcdsaAddress", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"EcdsaAddress\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Address{
		obj:           &s.obj.EcdsaAddress,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 20,
		dirty:         s.dirty,
	}
}

// BLSPublicKey is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Slot_) BLSPublicKey() *Storage_BLSPublicKey_ {
	// Field: &wrapper.tmplField{Type:"BLSPublicKey_", ElemType:"BLSPublicKey_", Name:"BLSPublicKey", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(BLSPublicKey_)32_storage", TupleRawName:"BLSPublicKey_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fcc00)}, TupleRawNames:[]string{"Key"}, TupleType:(*reflect.rtype)(0xc000438600)}, Tag:"json:\"BLSPublicKey\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_BLSPublicKey_{
		obj:           &s.obj.BLSPublicKey,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 64,
		dirty:         s.dirty,
	}
}

// EffectiveStake is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Slot_) EffectiveStake() *Storage_Decimal {
	// Field: &wrapper.tmplField{Type:"Decimal", ElemType:"Decimal", Name:"EffectiveStake", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fd100)}, TupleRawNames:[]string{"f"}, TupleType:(*reflect.rtype)(0xc0004387e0)}, Tag:"json:\"EffectiveStake\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("3", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Decimal{
		obj:           &s.obj.EffectiveStake,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Entrys is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Slots_) Entrys() *Storage_Slice_ed073f3f1ff1 {
	// Field: &wrapper.tmplField{Type:"Slice_ed073f3f1ff1", ElemType:"Slice_ed073f3f1ff1", Name:"Entrys", SolKind:abi.Type{Elem:(*abi.Type)(0xc000157000), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_struct(Slot_)128_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Entrys\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Entrys == nil {

		hash := s.db.GetState(s.addr, common.BigToHash(actual))
		instance := make(Slice_ed073f3f1ff1, hash.Big().Int64())

		s.obj.Entrys = instance
	}

	return &Storage_Slice_ed073f3f1ff1{
		obj:           &s.obj.Entrys,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
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

// Validators is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_ValidatorPool_) Validators() *Storage_ValidatorWrapperMap_ {
	// Field: &wrapper.tmplField{Type:"ValidatorWrapperMap_", ElemType:"ValidatorWrapperMap_", Name:"Validators", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(ValidatorWrapperMap_)121_storage", TupleRawName:"ValidatorWrapperMap_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fdf00), (*abi.Type)(0xc0000fdf80)}, TupleRawNames:[]string{"Keys", "Map"}, TupleType:(*reflect.rtype)(0xc00010be60)}, Tag:"json:\"Validators\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_ValidatorWrapperMap_{
		obj:           &s.obj.Validators,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 64,
		dirty:         s.dirty,
	}
}

// SlotKeySet is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_ValidatorPool_) SlotKeySet() *Storage_Map_624bcc0df495 {
	// Field: &wrapper.tmplField{Type:"Map_624bcc0df495", ElemType:"Map_624bcc0df495", Name:"SlotKeySet", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_string_memory_ptr,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc000156880), (*abi.Type)(0xc000156a00)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"SlotKeySet\" storage:\"slot=2,offset=0\"", Slot:2, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("2", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.SlotKeySet == nil {

		instance := make(Map_624bcc0df495)

		s.obj.SlotKeySet = instance
	}

	return &Storage_Map_624bcc0df495{
		obj:           s.obj.SlotKeySet,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// DescriptionIdentitySet is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_ValidatorPool_) DescriptionIdentitySet() *Storage_Map_624bcc0df495 {
	// Field: &wrapper.tmplField{Type:"Map_624bcc0df495", ElemType:"Map_624bcc0df495", Name:"DescriptionIdentitySet", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_string_memory_ptr,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc000156880), (*abi.Type)(0xc000156a00)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"DescriptionIdentitySet\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("3", 10)
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

// Committee is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_ValidatorPool_) Committee() *Storage_Committee_ {
	// Field: &wrapper.tmplField{Type:"Committee_", ElemType:"Committee_", Name:"Committee", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Committee_)137_storage", TupleRawName:"Committee_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000156c00), (*abi.Type)(0xc000156c80)}, TupleRawNames:[]string{"Epoch", "Slots"}, TupleType:(*reflect.rtype)(0xc000438000)}, Tag:"json:\"Committee\" storage:\"slot=4,offset=0\"", Slot:4, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("4", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Committee_{
		obj:           &s.obj.Committee,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 64,
		dirty:         s.dirty,
	}
}

// Entry is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_ValidatorWrapperMapEntry_) Entry() *Storage_ValidatorWrapper_ {
	// Field: &wrapper.tmplField{Type:"ValidatorWrapper_", ElemType:"ValidatorWrapper_", Name:"Entry", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x300, stringKind:"t_struct(ValidatorWrapper_)108_storage", TupleRawName:"ValidatorWrapper_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048c180), (*abi.Type)(0xc00048d080), (*abi.Type)(0xc00048d100), (*abi.Type)(0xc00048d300), (*abi.Type)(0xc00048d380), (*abi.Type)(0xc00048d480)}, TupleRawNames:[]string{"Validator", "Redelegations", "Counters", "BlockReward", "TotalDelegation", "TotalDelegationByOperator"}, TupleType:(*reflect.rtype)(0xc00010bda0)}, Tag:"json:\"Entry\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_ValidatorWrapper_{
		obj:           &s.obj.Entry,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 768,
		dirty:         s.dirty,
	}
}

// Index is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_ValidatorWrapperMapEntry_) Index() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Index", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Index\" storage:\"slot=24,offset=0\"", Slot:24, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("24", 10)
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
func (s *Storage_ValidatorWrapperMap_) Keys() *Storage_Slice_760032d9d704 {
	// Field: &wrapper.tmplField{Type:"Slice_760032d9d704", ElemType:"Slice_760032d9d704", Name:"Keys", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000fd600), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_address)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Keys\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
func (s *Storage_ValidatorWrapperMap_) Map() *Storage_Map_11675a09c025 {
	// Field: &wrapper.tmplField{Type:"Map_11675a09c025", ElemType:"Map_11675a09c025", Name:"Map", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(ValidatorWrapperMapEntry_)113_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048c000), (*abi.Type)(0xc000156680)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Map\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Map == nil {

		instance := make(Map_11675a09c025)

		s.obj.Map = instance
	}

	return &Storage_Map_11675a09c025{
		obj:           s.obj.Map,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Validator is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_ValidatorWrapper_) Validator() *Storage_Validator_ {
	// Field: &wrapper.tmplField{Type:"Validator_", ElemType:"Validator_", Name:"Validator", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x220, stringKind:"t_struct(Validator_)68_storage", TupleRawName:"Validator_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048c200), (*abi.Type)(0xc00048c280), (*abi.Type)(0xc00048c700), (*abi.Type)(0xc00048c780), (*abi.Type)(0xc00048c800), (*abi.Type)(0xc00048c900), (*abi.Type)(0xc00048c980), (*abi.Type)(0xc00048ca00), (*abi.Type)(0xc00048cf00)}, TupleRawNames:[]string{"ValidatorAddress", "OperatorAddresses", "SlotPubKeys", "LastEpochInCommittee", "MaxTotalDelegation", "Status", "Commission", "Description", "CreationHeight"}, TupleType:(*reflect.rtype)(0xc000439200)}, Tag:"json:\"Validator\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Validator_{
		obj:           &s.obj.Validator,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 544,
		dirty:         s.dirty,
	}
}

// Redelegations is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_ValidatorWrapper_) Redelegations() *Storage_RedelegationMap_ {
	// Field: &wrapper.tmplField{Type:"RedelegationMap_", ElemType:"RedelegationMap_", Name:"Redelegations", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(RedelegationMap_)95_storage", TupleRawName:"RedelegationMap_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fd480), (*abi.Type)(0xc0000fd700)}, TupleRawNames:[]string{"Keys", "Map"}, TupleType:(*reflect.rtype)(0xc000438c60)}, Tag:"json:\"Redelegations\" storage:\"slot=17,offset=0\"", Slot:17, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("17", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_RedelegationMap_{
		obj:           &s.obj.Redelegations,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 64,
		dirty:         s.dirty,
	}
}

// Counters is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_ValidatorWrapper_) Counters() *Storage_Counters_ {
	// Field: &wrapper.tmplField{Type:"Counters_", ElemType:"Counters_", Name:"Counters", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Counters_)41_storage", TupleRawName:"Counters_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048d180), (*abi.Type)(0xc00048d200)}, TupleRawNames:[]string{"NumBlocksToSign", "NumBlocksSigned"}, TupleType:(*reflect.rtype)(0xc0004392c0)}, Tag:"json:\"Counters\" storage:\"slot=19,offset=0\"", Slot:19, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("19", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Counters_{
		obj:           &s.obj.Counters,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 64,
		dirty:         s.dirty,
	}
}

// BlockReward is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_ValidatorWrapper_) BlockReward() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"BlockReward", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"BlockReward\" storage:\"slot=21,offset=0\"", Slot:21, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("21", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.BlockReward == nil {

		instance := big.NewInt(0)

		s.obj.BlockReward = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.BlockReward,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// TotalDelegation is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_ValidatorWrapper_) TotalDelegation() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"TotalDelegation", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"TotalDelegation\" storage:\"slot=22,offset=0\"", Slot:22, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("22", 10)
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

// TotalDelegationByOperator is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_ValidatorWrapper_) TotalDelegationByOperator() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"TotalDelegationByOperator", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"TotalDelegationByOperator\" storage:\"slot=23,offset=0\"", Slot:23, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("23", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.TotalDelegationByOperator == nil {

		instance := big.NewInt(0)

		s.obj.TotalDelegationByOperator = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.TotalDelegationByOperator,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// ValidatorAddress is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Validator_) ValidatorAddress() *Storage_Address {
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

// OperatorAddresses is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Validator_) OperatorAddresses() *Storage_AddressSet_ {
	// Field: &wrapper.tmplField{Type:"AddressSet_", ElemType:"AddressSet_", Name:"OperatorAddresses", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(AddressSet_)49_storage", TupleRawName:"AddressSet_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048c300), (*abi.Type)(0xc00048c380)}, TupleRawNames:[]string{"Keys", "Set"}, TupleType:(*reflect.rtype)(0xc000439080)}, Tag:"json:\"OperatorAddresses\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_AddressSet_{
		obj:           &s.obj.OperatorAddresses,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 64,
		dirty:         s.dirty,
	}
}

// SlotPubKeys is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Validator_) SlotPubKeys() *Storage_BLSPublicKeys_ {
	// Field: &wrapper.tmplField{Type:"BLSPublicKeys_", ElemType:"BLSPublicKeys_", Name:"SlotPubKeys", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(BLSPublicKeys_)36_storage", TupleRawName:"BLSPublicKeys_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fcb00)}, TupleRawNames:[]string{"Keys"}, TupleType:(*reflect.rtype)(0xc000438660)}, Tag:"json:\"SlotPubKeys\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("3", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_BLSPublicKeys_{
		obj:           &s.obj.SlotPubKeys,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// LastEpochInCommittee is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Validator_) LastEpochInCommittee() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"LastEpochInCommittee", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"LastEpochInCommittee\" storage:\"slot=4,offset=0\"", Slot:4, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("4", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.LastEpochInCommittee == nil {

		instance := big.NewInt(0)

		s.obj.LastEpochInCommittee = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.LastEpochInCommittee,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// MaxTotalDelegation is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Validator_) MaxTotalDelegation() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"MaxTotalDelegation", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"MaxTotalDelegation\" storage:\"slot=5,offset=0\"", Slot:5, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("5", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.MaxTotalDelegation == nil {

		instance := big.NewInt(0)

		s.obj.MaxTotalDelegation = instance
	}

	return &Storage_BigInt{
		obj:           s.obj.MaxTotalDelegation,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Status is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Validator_) Status() *Storage_Uint8 {
	// Field: &wrapper.tmplField{Type:"Uint8", ElemType:"Uint8", Name:"Status", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:8, T:0x1, numberOfBytes:0x1, stringKind:"uint8", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Status\" storage:\"slot=6,offset=0\"", Slot:6, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("6", 10)
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

// Commission is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Validator_) Commission() *Storage_Commission_ {
	// Field: &wrapper.tmplField{Type:"Commission_", ElemType:"Commission_", Name:"Commission", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x80, stringKind:"t_struct(Commission_)27_storage", TupleRawName:"Commission_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fd000), (*abi.Type)(0xc0000fd380)}, TupleRawNames:[]string{"CommissionRates", "UpdateHeight"}, TupleType:(*reflect.rtype)(0xc0004388a0)}, Tag:"json:\"Commission\" storage:\"slot=7,offset=0\"", Slot:7, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("7", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Commission_{
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
func (s *Storage_Validator_) Description() *Storage_Description_ {
	// Field: &wrapper.tmplField{Type:"Description_", ElemType:"Description_", Name:"Description", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xa0, stringKind:"t_struct(Description_)15_storage", TupleRawName:"Description_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00048ca80), (*abi.Type)(0xc00048cb80), (*abi.Type)(0xc00048cc00), (*abi.Type)(0xc00048cc80), (*abi.Type)(0xc00048cd00)}, TupleRawNames:[]string{"Name", "Identity", "Website", "SecurityContact", "Details"}, TupleType:(*reflect.rtype)(0xc0004391a0)}, Tag:"json:\"Description\" storage:\"slot=11,offset=0\"", Slot:11, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("11", 10)
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
func (s *Storage_Validator_) CreationHeight() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"CreationHeight", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"CreationHeight\" storage:\"slot=16,offset=0\"", Slot:16, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("16", 10)
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
	// &wrapper.tmplStruct{T:0xf, Name:"Decimal", Fields:[]*wrapper.tmplField(nil), Type:"common.Dec", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xf, numberOfBytes:0x20, stringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fd100)}, TupleRawNames:[]string{"f"}, TupleType:(*reflect.rtype)(0xc0004387e0)}}
	s.SetValue(common.NewDecFromBigIntWithPrec(big.NewInt(0), common.Precision))
}

func (s *Storage_String) Clear() {
	// &wrapper.tmplStruct{T:0x3, Name:"String", Fields:[]*wrapper.tmplField(nil), Type:"string", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x3, numberOfBytes:0x20, stringKind:"t_string_memory_ptr", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	rv := s.Value()
	s.SetValue(string(bytes.Repeat([]byte{0x00}, len(rv))))
	s.SetValue("")
}

func (s *Storage_Uint8) Clear() {
	// &wrapper.tmplStruct{T:0x1, Name:"Uint8", Fields:[]*wrapper.tmplField(nil), Type:"uint8", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:8, T:0x1, numberOfBytes:0x1, stringKind:"uint8", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	s.SetValue(0)
}

func (s *Storage_AddressSet_) Save(obj *AddressSet_) {
	if obj.Keys != nil {
		s.Keys().Save(obj.Keys)
	}
	if obj.Set != nil {
		s.Set().Save(obj.Set)
	}
}

func (s *Storage_AddressSet_) Clear() {
	s.Keys().Clear()
	s.Set().Clear()
}

func (s *Storage_AddressSet_) load() *AddressSet_ {
	s.Keys().load()
	s.Set().load()
	return s.obj
}

func (s *Storage_AddressSet_) Load() (*AddressSet_, error) {
	src := s.load()
	des := AddressSet_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *AddressSet_) Copy() (*AddressSet_, error) {
	des := AddressSet_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
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

func (s *Storage_CommissionRates_) Save(obj *CommissionRates_) {
	s.Rate().SetValue(obj.Rate)
	s.MaxRate().SetValue(obj.MaxRate)
	s.MaxChangeRate().SetValue(obj.MaxChangeRate)
}

func (s *Storage_CommissionRates_) Clear() {
	s.Rate().Clear()
	s.MaxRate().Clear()
	s.MaxChangeRate().Clear()
}

func (s *Storage_CommissionRates_) load() *CommissionRates_ {
	s.Rate().Value()
	s.MaxRate().Value()
	s.MaxChangeRate().Value()
	return s.obj
}

func (s *Storage_CommissionRates_) Load() (*CommissionRates_, error) {
	src := s.load()
	des := CommissionRates_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *CommissionRates_) Copy() (*CommissionRates_, error) {
	des := CommissionRates_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_Commission_) Save(obj *Commission_) {
	s.CommissionRates().Save(&obj.CommissionRates)
	if obj.UpdateHeight != nil {
		s.UpdateHeight().SetValue(obj.UpdateHeight)
	}
}

func (s *Storage_Commission_) Clear() {
	s.CommissionRates().Clear()
	s.UpdateHeight().Clear()
}

func (s *Storage_Commission_) load() *Commission_ {
	s.CommissionRates().load()
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

func (s *Storage_Committee_) Save(obj *Committee_) {
	if obj.Epoch != nil {
		s.Epoch().SetValue(obj.Epoch)
	}
	s.Slots().Save(&obj.Slots)
}

func (s *Storage_Committee_) Clear() {
	s.Epoch().Clear()
	s.Slots().Clear()
}

func (s *Storage_Committee_) load() *Committee_ {
	s.Epoch().Value()
	s.Slots().load()
	return s.obj
}

func (s *Storage_Committee_) Load() (*Committee_, error) {
	src := s.load()
	des := Committee_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Committee_) Copy() (*Committee_, error) {
	des := Committee_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_Counters_) Save(obj *Counters_) {
	if obj.NumBlocksToSign != nil {
		s.NumBlocksToSign().SetValue(obj.NumBlocksToSign)
	}
	if obj.NumBlocksSigned != nil {
		s.NumBlocksSigned().SetValue(obj.NumBlocksSigned)
	}
}

func (s *Storage_Counters_) Clear() {
	s.NumBlocksToSign().Clear()
	s.NumBlocksSigned().Clear()
}

func (s *Storage_Counters_) load() *Counters_ {
	s.NumBlocksToSign().Value()
	s.NumBlocksSigned().Value()
	return s.obj
}

func (s *Storage_Counters_) Load() (*Counters_, error) {
	src := s.load()
	des := Counters_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Counters_) Copy() (*Counters_, error) {
	des := Counters_{}
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
	s.ValidatorPool().Save(&obj.ValidatorPool)
}

func (s *Storage_Global_t) Clear() {
	s.ValidatorPool().Clear()
}

func (s *Storage_Global_t) load() *Global_t {
	s.ValidatorPool().load()
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

func (s *Storage_RedelegationMapEntry_) Save(obj *RedelegationMapEntry_) {
	s.Entry().Save(&obj.Entry)
	if obj.Index != nil {
		s.Index().SetValue(obj.Index)
	}
}

func (s *Storage_RedelegationMapEntry_) Clear() {
	s.Entry().Clear()
	s.Index().Clear()
}

func (s *Storage_RedelegationMapEntry_) load() *RedelegationMapEntry_ {
	s.Entry().load()
	s.Index().Value()
	return s.obj
}

func (s *Storage_RedelegationMapEntry_) Load() (*RedelegationMapEntry_, error) {
	src := s.load()
	des := RedelegationMapEntry_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *RedelegationMapEntry_) Copy() (*RedelegationMapEntry_, error) {
	des := RedelegationMapEntry_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_RedelegationMap_) Save(obj *RedelegationMap_) {
	if obj.Keys != nil {
		s.Keys().Save(obj.Keys)
	}
	if obj.Map != nil {
		s.Map().Save(obj.Map)
	}
}

func (s *Storage_RedelegationMap_) Clear() {
	s.Keys().Clear()
	s.Map().Clear()
}

func (s *Storage_RedelegationMap_) load() *RedelegationMap_ {
	s.Keys().load()
	s.Map().load()
	return s.obj
}

func (s *Storage_RedelegationMap_) Load() (*RedelegationMap_, error) {
	src := s.load()
	des := RedelegationMap_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *RedelegationMap_) Copy() (*RedelegationMap_, error) {
	des := RedelegationMap_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_Redelegation_) Save(obj *Redelegation_) {
	s.DelegatorAddress().SetValue(obj.DelegatorAddress)
	if obj.Amount != nil {
		s.Amount().SetValue(obj.Amount)
	}
	if obj.Reward != nil {
		s.Reward().SetValue(obj.Reward)
	}
	s.Undelegation().Save(&obj.Undelegation)
}

func (s *Storage_Redelegation_) Clear() {
	s.DelegatorAddress().Clear()
	s.Amount().Clear()
	s.Reward().Clear()
	s.Undelegation().Clear()
}

func (s *Storage_Redelegation_) load() *Redelegation_ {
	s.DelegatorAddress().Value()
	s.Amount().Value()
	s.Reward().Value()
	s.Undelegation().load()
	return s.obj
}

func (s *Storage_Redelegation_) Load() (*Redelegation_, error) {
	src := s.load()
	des := Redelegation_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Redelegation_) Copy() (*Redelegation_, error) {
	des := Redelegation_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_Slot_) Save(obj *Slot_) {
	s.EcdsaAddress().SetValue(obj.EcdsaAddress)
	s.BLSPublicKey().Save(&obj.BLSPublicKey)
	s.EffectiveStake().SetValue(obj.EffectiveStake)
}

func (s *Storage_Slot_) Clear() {
	s.EcdsaAddress().Clear()
	s.BLSPublicKey().Clear()
	s.EffectiveStake().Clear()
}

func (s *Storage_Slot_) load() *Slot_ {
	s.EcdsaAddress().Value()
	s.BLSPublicKey().load()
	s.EffectiveStake().Value()
	return s.obj
}

func (s *Storage_Slot_) Load() (*Slot_, error) {
	src := s.load()
	des := Slot_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Slot_) Copy() (*Slot_, error) {
	des := Slot_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_Slots_) Save(obj *Slots_) {
	if obj.Entrys != nil {
		s.Entrys().Save(obj.Entrys)
	}
}

func (s *Storage_Slots_) Clear() {
	s.Entrys().Clear()
}

func (s *Storage_Slots_) load() *Slots_ {
	s.Entrys().load()
	return s.obj
}

func (s *Storage_Slots_) Load() (*Slots_, error) {
	src := s.load()
	des := Slots_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Slots_) Copy() (*Slots_, error) {
	des := Slots_{}
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

func (s *Storage_ValidatorPool_) Save(obj *ValidatorPool_) {
	s.Validators().Save(&obj.Validators)
	if obj.SlotKeySet != nil {
		s.SlotKeySet().Save(obj.SlotKeySet)
	}
	if obj.DescriptionIdentitySet != nil {
		s.DescriptionIdentitySet().Save(obj.DescriptionIdentitySet)
	}
	s.Committee().Save(&obj.Committee)
}

func (s *Storage_ValidatorPool_) Clear() {
	s.Validators().Clear()
	s.SlotKeySet().Clear()
	s.DescriptionIdentitySet().Clear()
	s.Committee().Clear()
}

func (s *Storage_ValidatorPool_) load() *ValidatorPool_ {
	s.Validators().load()
	s.SlotKeySet().load()
	s.DescriptionIdentitySet().load()
	s.Committee().load()
	return s.obj
}

func (s *Storage_ValidatorPool_) Load() (*ValidatorPool_, error) {
	src := s.load()
	des := ValidatorPool_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *ValidatorPool_) Copy() (*ValidatorPool_, error) {
	des := ValidatorPool_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_ValidatorWrapperMapEntry_) Save(obj *ValidatorWrapperMapEntry_) {
	s.Entry().Save(&obj.Entry)
	if obj.Index != nil {
		s.Index().SetValue(obj.Index)
	}
}

func (s *Storage_ValidatorWrapperMapEntry_) Clear() {
	s.Entry().Clear()
	s.Index().Clear()
}

func (s *Storage_ValidatorWrapperMapEntry_) load() *ValidatorWrapperMapEntry_ {
	s.Entry().load()
	s.Index().Value()
	return s.obj
}

func (s *Storage_ValidatorWrapperMapEntry_) Load() (*ValidatorWrapperMapEntry_, error) {
	src := s.load()
	des := ValidatorWrapperMapEntry_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *ValidatorWrapperMapEntry_) Copy() (*ValidatorWrapperMapEntry_, error) {
	des := ValidatorWrapperMapEntry_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_ValidatorWrapperMap_) Save(obj *ValidatorWrapperMap_) {
	if obj.Keys != nil {
		s.Keys().Save(obj.Keys)
	}
	if obj.Map != nil {
		s.Map().Save(obj.Map)
	}
}

func (s *Storage_ValidatorWrapperMap_) Clear() {
	s.Keys().Clear()
	s.Map().Clear()
}

func (s *Storage_ValidatorWrapperMap_) load() *ValidatorWrapperMap_ {
	s.Keys().load()
	s.Map().load()
	return s.obj
}

func (s *Storage_ValidatorWrapperMap_) Load() (*ValidatorWrapperMap_, error) {
	src := s.load()
	des := ValidatorWrapperMap_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *ValidatorWrapperMap_) Copy() (*ValidatorWrapperMap_, error) {
	des := ValidatorWrapperMap_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_ValidatorWrapper_) Save(obj *ValidatorWrapper_) {
	s.Validator().Save(&obj.Validator)
	s.Redelegations().Save(&obj.Redelegations)
	s.Counters().Save(&obj.Counters)
	if obj.BlockReward != nil {
		s.BlockReward().SetValue(obj.BlockReward)
	}
	if obj.TotalDelegation != nil {
		s.TotalDelegation().SetValue(obj.TotalDelegation)
	}
	if obj.TotalDelegationByOperator != nil {
		s.TotalDelegationByOperator().SetValue(obj.TotalDelegationByOperator)
	}
}

func (s *Storage_ValidatorWrapper_) Clear() {
	s.Validator().Clear()
	s.Redelegations().Clear()
	s.Counters().Clear()
	s.BlockReward().Clear()
	s.TotalDelegation().Clear()
	s.TotalDelegationByOperator().Clear()
}

func (s *Storage_ValidatorWrapper_) load() *ValidatorWrapper_ {
	s.Validator().load()
	s.Redelegations().load()
	s.Counters().load()
	s.BlockReward().Value()
	s.TotalDelegation().Value()
	s.TotalDelegationByOperator().Value()
	return s.obj
}

func (s *Storage_ValidatorWrapper_) Load() (*ValidatorWrapper_, error) {
	src := s.load()
	des := ValidatorWrapper_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *ValidatorWrapper_) Copy() (*ValidatorWrapper_, error) {
	des := ValidatorWrapper_{}
	if err := deepCopy(s, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Storage_Validator_) Save(obj *Validator_) {
	s.ValidatorAddress().SetValue(obj.ValidatorAddress)
	s.OperatorAddresses().Save(&obj.OperatorAddresses)
	s.SlotPubKeys().Save(&obj.SlotPubKeys)
	if obj.LastEpochInCommittee != nil {
		s.LastEpochInCommittee().SetValue(obj.LastEpochInCommittee)
	}
	if obj.MaxTotalDelegation != nil {
		s.MaxTotalDelegation().SetValue(obj.MaxTotalDelegation)
	}
	s.Status().SetValue(obj.Status)
	s.Commission().Save(&obj.Commission)
	s.Description().Save(&obj.Description)
	if obj.CreationHeight != nil {
		s.CreationHeight().SetValue(obj.CreationHeight)
	}
}

func (s *Storage_Validator_) Clear() {
	s.ValidatorAddress().Clear()
	s.OperatorAddresses().Clear()
	s.SlotPubKeys().Clear()
	s.LastEpochInCommittee().Clear()
	s.MaxTotalDelegation().Clear()
	s.Status().Clear()
	s.Commission().Clear()
	s.Description().Clear()
	s.CreationHeight().Clear()
}

func (s *Storage_Validator_) load() *Validator_ {
	s.ValidatorAddress().Value()
	s.OperatorAddresses().load()
	s.SlotPubKeys().load()
	s.LastEpochInCommittee().Value()
	s.MaxTotalDelegation().Value()
	s.Status().Value()
	s.Commission().load()
	s.Description().load()
	s.CreationHeight().Value()
	return s.obj
}

func (s *Storage_Validator_) Load() (*Validator_, error) {
	src := s.load()
	des := Validator_{}
	if err := deepCopy(src, &des); err != nil {
		return nil, err
	}
	return &des, nil
}

func (s *Validator_) Copy() (*Validator_, error) {
	des := Validator_{}
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

func (s *Storage_Map_11675a09c025) Save(obj Map_11675a09c025) {
	for k, v := range obj {
		s.Get(k).Save(v)
	}
}

func (s *Storage_Map_11675a09c025) Clear() {
	panic("not support to clear map")
}

func (s *Storage_Map_11675a09c025) load() Map_11675a09c025 {
	panic("not support to clear map")
}

func (s *Storage_Map_1e3b3171b974) Save(obj Map_1e3b3171b974) {
	for k, v := range obj {
		s.Get(k).Save(v)
	}
}

func (s *Storage_Map_1e3b3171b974) Clear() {
	panic("not support to clear map")
}

func (s *Storage_Map_1e3b3171b974) load() Map_1e3b3171b974 {
	panic("not support to clear map")
}

func (s *Storage_Map_3783ff483121) Save(obj Map_3783ff483121) {
	for k, v := range obj {
		s.Get(k).SetValue(*v)
	}
}

func (s *Storage_Map_3783ff483121) Clear() {
	panic("not support to clear map")
}

func (s *Storage_Map_3783ff483121) load() Map_3783ff483121 {
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

func (s *Storage_Slice_6bb0ade2b43c) Save(obj Slice_6bb0ade2b43c) {
	for i := 0; i < len(obj); i++ {
		s.Get(i).Save(obj[i])
	}
	s.Resize(len(obj))
}

func (s *Storage_Slice_6bb0ade2b43c) Clear() {
	for i := 0; i < s.Length(); i++ {
		s.Get(i).Clear()
	}
	s.Resize(0)
}

func (s *Storage_Slice_6bb0ade2b43c) load() Slice_6bb0ade2b43c {
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

func (s *Storage_Slice_ed073f3f1ff1) Save(obj Slice_ed073f3f1ff1) {
	for i := 0; i < len(obj); i++ {
		s.Get(i).Save(obj[i])
	}
	s.Resize(len(obj))
}

func (s *Storage_Slice_ed073f3f1ff1) Clear() {
	for i := 0; i < s.Length(); i++ {
		s.Get(i).Clear()
	}
	s.Resize(0)
}

func (s *Storage_Slice_ed073f3f1ff1) load() Slice_ed073f3f1ff1 {
	for i := 0; i < s.Length(); i++ {
		s.Get(i).load()
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
