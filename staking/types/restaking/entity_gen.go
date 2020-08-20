// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package restaking

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

	// &wrapper.tmplStruct{T:0x7, Name:"Address", Fields:[]*wrapper.tmplField(nil), Type:"common.Address", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x11, Type:(*reflect.rtype)(0x48ba1e0), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *Address
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Address) Value() common.Address {
	// &wrapper.tmplStruct{T:0x7, Name:"Address", Fields:[]*wrapper.tmplField(nil), Type:"common.Address", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x11, Type:(*reflect.rtype)(0x48ba1e0), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}

	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	data := hash.Bytes()[32-(s.offset+s.numberOfBytes) : 32-s.offset]

	*s.obj = common.BytesToAddress(data)
	return *s.obj
}

func (s *Storage_Address) SetValue(value common.Address) {
	// &wrapper.tmplStruct{T:0x7, Name:"Address", Fields:[]*wrapper.tmplField(nil), Type:"common.Address", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x11, Type:(*reflect.rtype)(0x48ba1e0), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
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

	// &wrapper.tmplStruct{T:0x1, Name:"BigInt", Fields:[]*wrapper.tmplField(nil), Type:"*big.Int", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x4907560), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           BigInt
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_BigInt) Value() *big.Int {
	// &wrapper.tmplStruct{T:0x1, Name:"BigInt", Fields:[]*wrapper.tmplField(nil), Type:"*big.Int", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x4907560), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}

	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	data := hash.Bytes()[32-(s.offset+s.numberOfBytes) : 32-s.offset]

	*s.obj = *(big.NewInt(0).SetBytes(data))
	retval := big.NewInt(0).SetBytes(data)
	return retval
}

func (s *Storage_BigInt) SetValue(value *big.Int) {
	// &wrapper.tmplStruct{T:0x1, Name:"BigInt", Fields:[]*wrapper.tmplField(nil), Type:"*big.Int", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x4907560), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
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

	// &wrapper.tmplStruct{T:0x2, Name:"Bool", Fields:[]*wrapper.tmplField(nil), Type:"bool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x1, Type:(*reflect.rtype)(0x47d5a40), Size:0, T:0x2, NumberOfBytes:0x1, StringKind:"t_bool", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *Bool
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Bool) Value() bool {
	// &wrapper.tmplStruct{T:0x2, Name:"Bool", Fields:[]*wrapper.tmplField(nil), Type:"bool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x1, Type:(*reflect.rtype)(0x47d5a40), Size:0, T:0x2, NumberOfBytes:0x1, StringKind:"t_bool", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}

	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	data := hash.Bytes()[32-(s.offset+s.numberOfBytes) : 32-s.offset]

	*s.obj = bool(big.NewInt(0).SetBytes(data).Int64() != 0)
	return *s.obj
}

func (s *Storage_Bool) SetValue(value bool) {
	// &wrapper.tmplStruct{T:0x2, Name:"Bool", Fields:[]*wrapper.tmplField(nil), Type:"bool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x1, Type:(*reflect.rtype)(0x47d5a40), Size:0, T:0x2, NumberOfBytes:0x1, StringKind:"t_bool", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
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

	// &wrapper.tmplStruct{T:0xb, Name:"Decimal", Fields:[]*wrapper.tmplField(nil), Type:"common.Dec", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0x48ff7a0), Size:0, T:0xb, NumberOfBytes:0x20, StringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc000592f30)}, TupleRawNames:[]string{"f"}}}
	obj           *Decimal
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Decimal) Value() common.Dec {
	// &wrapper.tmplStruct{T:0xb, Name:"Decimal", Fields:[]*wrapper.tmplField(nil), Type:"common.Dec", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0x48ff7a0), Size:0, T:0xb, NumberOfBytes:0x20, StringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc000592f30)}, TupleRawNames:[]string{"f"}}}
	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	*s.obj = common.NewDecFromBigIntWithPrec(hash.Big(), common.Precision)
	return *s.obj
}

func (s *Storage_Decimal) SetValue(value common.Dec) {
	// &wrapper.tmplStruct{T:0xb, Name:"Decimal", Fields:[]*wrapper.tmplField(nil), Type:"common.Dec", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0x48ff7a0), Size:0, T:0xb, NumberOfBytes:0x20, StringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc000592f30)}, TupleRawNames:[]string{"f"}}}
	hash := value.BigInt()
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(hash))
	*s.obj = value
}

type String = string
type Storage_String struct {
	// Builtin-Type

	// &wrapper.tmplStruct{T:0x3, Name:"String", Fields:[]*wrapper.tmplField(nil), Type:"string", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x18, Type:(*reflect.rtype)(0x47e0080), Size:0, T:0x3, NumberOfBytes:0x20, StringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *String
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_String) Value() string {
	// &wrapper.tmplStruct{T:0x3, Name:"String", Fields:[]*wrapper.tmplField(nil), Type:"string", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x18, Type:(*reflect.rtype)(0x47e0080), Size:0, T:0x3, NumberOfBytes:0x20, StringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	rv := GetStateAsBytes(s.db, s.addr, s.slot)
	*s.obj = string(rv)
	return *s.obj
}

func (s *Storage_String) SetValue(value string) {
	// &wrapper.tmplStruct{T:0x3, Name:"String", Fields:[]*wrapper.tmplField(nil), Type:"string", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x18, Type:(*reflect.rtype)(0x47e0080), Size:0, T:0x3, NumberOfBytes:0x20, StringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	SetStateAsBytes(s.db, s.addr, s.slot, []byte(value))
	*s.obj = value
}

type Uint8 = uint8
type Storage_Uint8 struct {
	// Builtin-Type

	// &wrapper.tmplStruct{T:0x1, Name:"Uint8", Fields:[]*wrapper.tmplField(nil), Type:"uint8", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x8, Type:(*reflect.rtype)(0x47e0300), Size:0, T:0x1, NumberOfBytes:0x1, StringKind:"uint8", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *Uint8
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Uint8) Value() uint8 {
	// &wrapper.tmplStruct{T:0x1, Name:"Uint8", Fields:[]*wrapper.tmplField(nil), Type:"uint8", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x8, Type:(*reflect.rtype)(0x47e0300), Size:0, T:0x1, NumberOfBytes:0x1, StringKind:"uint8", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}

	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	data := hash.Bytes()[32-(s.offset+s.numberOfBytes) : 32-s.offset]

	*s.obj = uint8(big.NewInt(0).SetBytes(data).Uint64())
	return *s.obj
}

func (s *Storage_Uint8) SetValue(value uint8) {
	// &wrapper.tmplStruct{T:0x1, Name:"Uint8", Fields:[]*wrapper.tmplField(nil), Type:"uint8", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x8, Type:(*reflect.rtype)(0x47e0300), Size:0, T:0x1, NumberOfBytes:0x1, StringKind:"uint8", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	val := math.PaddedBigBytes(big.NewInt(0).SetUint64(uint64(value)), 32)

	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	data := hash.Bytes()
	copy(data[32-(s.offset+s.numberOfBytes):32-s.offset], val[len(val)-s.numberOfBytes:])
	hash.SetBytes(data)
	s.db.SetState(s.addr, common.BigToHash(s.slot), hash)

	*s.obj = value
}

// Array_aaa956d410fd is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x5, Name:"Array_aaa956d410fd", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017ca80)}, Type:"[48]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc0005927e0), Kind:0x11, Type:(*reflect.rtype)(0x47ebe80), Size:48, T:0x5, NumberOfBytes:0x40, StringKind:"t_array(t_bytes1)48_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
// "[48]uint8"
type Array_aaa956d410fd [48]Uint8

type Storage_Array_aaa956d410fd struct {

	// &wrapper.tmplStruct{T:0x5, Name:"Array_aaa956d410fd", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017ca80)}, Type:"[48]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc0005927e0), Kind:0x11, Type:(*reflect.rtype)(0x47ebe80), Size:48, T:0x5, NumberOfBytes:0x40, StringKind:"t_array(t_bytes1)48_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *Array_aaa956d410fd
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Array_aaa956d410fd) Value() [48]Uint8 {
	// &wrapper.tmplStruct{T:0x5, Name:"Array_aaa956d410fd", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017ca80)}, Type:"[48]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc0005927e0), Kind:0x11, Type:(*reflect.rtype)(0x47ebe80), Size:48, T:0x5, NumberOfBytes:0x40, StringKind:"t_array(t_bytes1)48_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
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
	// &wrapper.tmplStruct{T:0x5, Name:"Array_aaa956d410fd", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017ca80)}, Type:"[48]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc0005927e0), Kind:0x11, Type:(*reflect.rtype)(0x47ebe80), Size:48, T:0x5, NumberOfBytes:0x40, StringKind:"t_array(t_bytes1)48_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
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
// &wrapper.tmplStruct{T:0xd, Name:"Map_11675a09c025", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017c0e0), (*wrapper.tmplField)(0xc00017c380)}, Type:"map[Address]*ValidatorWrapperMapEntry_", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc000571800), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_address,t_struct(ValidatorWrapperMapEntry_)113_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc00058bc20), (*abi.Type)(0xc00058bd40)}, TupleRawNames:[]string(nil)}}
// "map[common.Address]*struct { Entry struct { Validator struct { ValidatorAddress common.Address \"json:\\\"ValidatorAddress\\\" storage:\\\"slot=0,offset=0\\\"\"; OperatorAddresses struct { Keys []*common.Address \"json:\\\"Keys\\\" storage:\\\"slot=0,offset=0\\\"\"; Set map[common.Address]*bool \"json:\\\"Set\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"OperatorAddresses\\\" storage:\\\"slot=1,offset=0\\\"\"; SlotPubKeys struct { Keys []*struct { Key [48]uint8 \"json:\\\"Key\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"Keys\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"SlotPubKeys\\\" storage:\\\"slot=3,offset=0\\\"\"; LastEpochInCommittee *big.Int \"json:\\\"LastEpochInCommittee\\\" storage:\\\"slot=4,offset=0\\\"\"; MaxTotalDelegation *big.Int \"json:\\\"MaxTotalDelegation\\\" storage:\\\"slot=5,offset=0\\\"\"; Status uint8 \"json:\\\"Status\\\" storage:\\\"slot=6,offset=0\\\"\"; Commission struct { CommissionRates struct { Rate struct { F *big.Int \"json:\\\"f\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"Rate\\\" storage:\\\"slot=0,offset=0\\\"\"; MaxRate struct { F *big.Int \"json:\\\"f\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"MaxRate\\\" storage:\\\"slot=1,offset=0\\\"\"; MaxChangeRate struct { F *big.Int \"json:\\\"f\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"MaxChangeRate\\\" storage:\\\"slot=2,offset=0\\\"\" } \"json:\\\"CommissionRates\\\" storage:\\\"slot=0,offset=0\\\"\"; UpdateHeight *big.Int \"json:\\\"UpdateHeight\\\" storage:\\\"slot=3,offset=0\\\"\" } \"json:\\\"Commission\\\" storage:\\\"slot=7,offset=0\\\"\"; Description struct { Name string \"json:\\\"Name\\\" storage:\\\"slot=0,offset=0\\\"\"; Identity string \"json:\\\"Identity\\\" storage:\\\"slot=1,offset=0\\\"\"; Website string \"json:\\\"Website\\\" storage:\\\"slot=2,offset=0\\\"\"; SecurityContact string \"json:\\\"SecurityContact\\\" storage:\\\"slot=3,offset=0\\\"\"; Details string \"json:\\\"Details\\\" storage:\\\"slot=4,offset=0\\\"\" } \"json:\\\"Description\\\" storage:\\\"slot=11,offset=0\\\"\"; CreationHeight *big.Int \"json:\\\"CreationHeight\\\" storage:\\\"slot=16,offset=0\\\"\" } \"json:\\\"Validator\\\" storage:\\\"slot=0,offset=0\\\"\"; Redelegations struct { Keys []*common.Address \"json:\\\"Keys\\\" storage:\\\"slot=0,offset=0\\\"\"; Map map[common.Address]*struct { Entry struct { DelegatorAddress common.Address \"json:\\\"DelegatorAddress\\\" storage:\\\"slot=0,offset=0\\\"\"; Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=1,offset=0\\\"\"; Reward *big.Int \"json:\\\"Reward\\\" storage:\\\"slot=2,offset=0\\\"\"; Undelegation struct { Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=0,offset=0\\\"\"; Epoch *big.Int \"json:\\\"Epoch\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"Undelegation\\\" storage:\\\"slot=3,offset=0\\\"\" } \"json:\\\"Entry\\\" storage:\\\"slot=0,offset=0\\\"\"; Index *big.Int \"json:\\\"Index\\\" storage:\\\"slot=5,offset=0\\\"\" } \"json:\\\"Map\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"Redelegations\\\" storage:\\\"slot=17,offset=0\\\"\"; Counters struct { NumBlocksToSign *big.Int \"json:\\\"NumBlocksToSign\\\" storage:\\\"slot=0,offset=0\\\"\"; NumBlocksSigned *big.Int \"json:\\\"NumBlocksSigned\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"Counters\\\" storage:\\\"slot=19,offset=0\\\"\"; BlockReward *big.Int \"json:\\\"BlockReward\\\" storage:\\\"slot=21,offset=0\\\"\"; TotalDelegation *big.Int \"json:\\\"TotalDelegation\\\" storage:\\\"slot=22,offset=0\\\"\"; TotalDelegationByOperator *big.Int \"json:\\\"TotalDelegationByOperator\\\" storage:\\\"slot=23,offset=0\\\"\" } \"json:\\\"Entry\\\" storage:\\\"slot=0,offset=0\\\"\"; Index *big.Int \"json:\\\"Index\\\" storage:\\\"slot=24,offset=0\\\"\" }"
type Map_11675a09c025 map[Address]*ValidatorWrapperMapEntry_

type Storage_Map_11675a09c025 struct {

	// &wrapper.tmplStruct{T:0xd, Name:"Map_11675a09c025", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017c0e0), (*wrapper.tmplField)(0xc00017c380)}, Type:"map[Address]*ValidatorWrapperMapEntry_", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc000571800), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_address,t_struct(ValidatorWrapperMapEntry_)113_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc00058bc20), (*abi.Type)(0xc00058bd40)}, TupleRawNames:[]string(nil)}}
	obj           Map_11675a09c025
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Map_11675a09c025) Get(key Address) *Storage_ValidatorWrapperMapEntry_ {
	// Key:	&wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x11, Type:(*reflect.rtype)(0x48ba1e0), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"ValidatorWrapperMapEntry_", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0xc0005655c0), Size:0, T:0x6, NumberOfBytes:0x320, StringKind:"t_struct(ValidatorWrapperMapEntry_)113_storage", TupleRawName:"ValidatorWrapperMapEntry_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00058bdd0), (*abi.Type)(0xc0005a4870)}, TupleRawNames:[]string{"Entry", "Index"}}, Tag:"", Slot:<nil>, Offset:0}
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
// &wrapper.tmplStruct{T:0xd, Name:"Map_1e3b3171b974", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000718e00), (*wrapper.tmplField)(0xc0007195e0)}, Type:"map[Address]*RedelegationMapEntry_", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc0005715c0), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_address,t_struct(RedelegationMapEntry_)87_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc000593a70), (*abi.Type)(0xc000593b00)}, TupleRawNames:[]string(nil)}}
// "map[common.Address]*struct { Entry struct { DelegatorAddress common.Address \"json:\\\"DelegatorAddress\\\" storage:\\\"slot=0,offset=0\\\"\"; Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=1,offset=0\\\"\"; Reward *big.Int \"json:\\\"Reward\\\" storage:\\\"slot=2,offset=0\\\"\"; Undelegation struct { Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=0,offset=0\\\"\"; Epoch *big.Int \"json:\\\"Epoch\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"Undelegation\\\" storage:\\\"slot=3,offset=0\\\"\" } \"json:\\\"Entry\\\" storage:\\\"slot=0,offset=0\\\"\"; Index *big.Int \"json:\\\"Index\\\" storage:\\\"slot=5,offset=0\\\"\" }"
type Map_1e3b3171b974 map[Address]*RedelegationMapEntry_

type Storage_Map_1e3b3171b974 struct {

	// &wrapper.tmplStruct{T:0xd, Name:"Map_1e3b3171b974", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000718e00), (*wrapper.tmplField)(0xc0007195e0)}, Type:"map[Address]*RedelegationMapEntry_", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc0005715c0), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_address,t_struct(RedelegationMapEntry_)87_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc000593a70), (*abi.Type)(0xc000593b00)}, TupleRawNames:[]string(nil)}}
	obj           Map_1e3b3171b974
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Map_1e3b3171b974) Get(key Address) *Storage_RedelegationMapEntry_ {
	// Key:	&wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x11, Type:(*reflect.rtype)(0x48ba1e0), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"RedelegationMapEntry_", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0xc000565500), Size:0, T:0x6, NumberOfBytes:0xc0, StringKind:"t_struct(RedelegationMapEntry_)87_storage", TupleRawName:"RedelegationMapEntry_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000593b90), (*abi.Type)(0xc0005a4120)}, TupleRawNames:[]string{"Entry", "Index"}}, Tag:"", Slot:<nil>, Offset:0}
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
// &wrapper.tmplStruct{T:0xd, Name:"Map_3783ff483121", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017d180), (*wrapper.tmplField)(0xc00017d260)}, Type:"map[Address]*Bool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc000570ae0), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_address,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0005922d0), (*abi.Type)(0xc000592360)}, TupleRawNames:[]string(nil)}}
// "map[common.Address]*bool"
type Map_3783ff483121 map[Address]*Bool

type Storage_Map_3783ff483121 struct {

	// &wrapper.tmplStruct{T:0xd, Name:"Map_3783ff483121", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017d180), (*wrapper.tmplField)(0xc00017d260)}, Type:"map[Address]*Bool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc000570ae0), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_address,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0005922d0), (*abi.Type)(0xc000592360)}, TupleRawNames:[]string(nil)}}
	obj           Map_3783ff483121
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Map_3783ff483121) Get(key Address) *Storage_Bool {
	// Key:	&wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x11, Type:(*reflect.rtype)(0x48ba1e0), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"Bool", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x47b4780), Size:0, T:0x2, NumberOfBytes:0x1, StringKind:"t_bool", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"", Slot:<nil>, Offset:0}
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
// &wrapper.tmplStruct{T:0xd, Name:"Map_624bcc0df495", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017db20), (*wrapper.tmplField)(0xc00017dc00)}, Type:"map[String]*Bool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc000571b60), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_string_memory_ptr,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0005a54d0), (*abi.Type)(0xc0005a55f0)}, TupleRawNames:[]string(nil)}}
// "map[string]*bool"
type Map_624bcc0df495 map[String]*Bool

type Storage_Map_624bcc0df495 struct {

	// &wrapper.tmplStruct{T:0xd, Name:"Map_624bcc0df495", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017db20), (*wrapper.tmplField)(0xc00017dc00)}, Type:"map[String]*Bool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc000571b60), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_string_memory_ptr,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0005a54d0), (*abi.Type)(0xc0005a55f0)}, TupleRawNames:[]string(nil)}}
	obj           Map_624bcc0df495
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Map_624bcc0df495) Get(key String) *Storage_Bool {
	// Key:	&wrapper.tmplField{Type:"String", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x18, Type:(*reflect.rtype)(0x47e0080), Size:0, T:0x3, NumberOfBytes:0x20, StringKind:"t_string_memory_ptr", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"Bool", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x47b4780), Size:0, T:0x2, NumberOfBytes:0x1, StringKind:"t_bool", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"", Slot:<nil>, Offset:0}
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
// &wrapper.tmplStruct{T:0x4, Name:"Slice_6bb0ade2b43c", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017d500)}, Type:"[]*BLSPublicKey_", SolKind:abi.Type{Elem:(*abi.Type)(0xc0005926c0), Kind:0x17, Type:(*reflect.rtype)(0xc000565340), Size:0, T:0x4, NumberOfBytes:0x20, StringKind:"t_array(t_struct(BLSPublicKey_)32_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
// "[]*struct { Key [48]uint8 \"json:\\\"Key\\\" storage:\\\"slot=0,offset=0\\\"\" }"
type Slice_6bb0ade2b43c []*BLSPublicKey_

type Storage_Slice_6bb0ade2b43c struct {

	// &wrapper.tmplStruct{T:0x4, Name:"Slice_6bb0ade2b43c", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017d500)}, Type:"[]*BLSPublicKey_", SolKind:abi.Type{Elem:(*abi.Type)(0xc0005926c0), Kind:0x17, Type:(*reflect.rtype)(0xc000565340), Size:0, T:0x4, NumberOfBytes:0x20, StringKind:"t_array(t_struct(BLSPublicKey_)32_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
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
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(big.NewInt(0).SetUint64(uint64(length))))

	slice := make([]*BLSPublicKey_, length, length+50)
	copy(slice, *s.obj)
	*s.obj = slice
}

func (s *Storage_Slice_6bb0ade2b43c) Get(index int) *Storage_BLSPublicKey_ {
	// Value: &wrapper.tmplField{Type:"BLSPublicKey_", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0xc000565300), Size:0, T:0x6, NumberOfBytes:0x40, StringKind:"t_struct(BLSPublicKey_)32_storage", TupleRawName:"BLSPublicKey_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000592750)}, TupleRawNames:[]string{"Key"}}, Tag:"", Slot:<nil>, Offset:0}
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
// &wrapper.tmplStruct{T:0x4, Name:"Slice_760032d9d704", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017cfc0)}, Type:"[]*Address", SolKind:abi.Type{Elem:(*abi.Type)(0xc000592120), Kind:0x17, Type:(*reflect.rtype)(0xc0005652c0), Size:0, T:0x4, NumberOfBytes:0x20, StringKind:"t_array(t_address)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
// "[]*common.Address"
type Slice_760032d9d704 []*Address

type Storage_Slice_760032d9d704 struct {

	// &wrapper.tmplStruct{T:0x4, Name:"Slice_760032d9d704", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017cfc0)}, Type:"[]*Address", SolKind:abi.Type{Elem:(*abi.Type)(0xc000592120), Kind:0x17, Type:(*reflect.rtype)(0xc0005652c0), Size:0, T:0x4, NumberOfBytes:0x20, StringKind:"t_array(t_address)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
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
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(big.NewInt(0).SetUint64(uint64(length))))

	slice := make([]*Address, length, length+50)
	copy(slice, *s.obj)
	*s.obj = slice
}

func (s *Storage_Slice_760032d9d704) Get(index int) *Storage_Address {
	// Value: &wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x48d3980), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"", Slot:<nil>, Offset:0}
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
// &wrapper.tmplStruct{T:0x4, Name:"Slice_ed073f3f1ff1", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017ce00)}, Type:"[]*Slot_", SolKind:abi.Type{Elem:(*abi.Type)(0xc0005a59e0), Kind:0x17, Type:(*reflect.rtype)(0xc0005656c0), Size:0, T:0x4, NumberOfBytes:0x20, StringKind:"t_array(t_struct(Slot_)128_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
// "[]*struct { EcdsaAddress common.Address \"json:\\\"EcdsaAddress\\\" storage:\\\"slot=0,offset=0\\\"\"; BLSPublicKey struct { Key [48]uint8 \"json:\\\"Key\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"BLSPublicKey\\\" storage:\\\"slot=1,offset=0\\\"\"; EffectiveStake struct { F *big.Int \"json:\\\"f\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"EffectiveStake\\\" storage:\\\"slot=3,offset=0\\\"\" }"
type Slice_ed073f3f1ff1 []*Slot_

type Storage_Slice_ed073f3f1ff1 struct {

	// &wrapper.tmplStruct{T:0x4, Name:"Slice_ed073f3f1ff1", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017ce00)}, Type:"[]*Slot_", SolKind:abi.Type{Elem:(*abi.Type)(0xc0005a59e0), Kind:0x17, Type:(*reflect.rtype)(0xc0005656c0), Size:0, T:0x4, NumberOfBytes:0x20, StringKind:"t_array(t_struct(Slot_)128_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
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
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(big.NewInt(0).SetUint64(uint64(length))))

	slice := make([]*Slot_, length, length+50)
	copy(slice, *s.obj)
	*s.obj = slice
}

func (s *Storage_Slice_ed073f3f1ff1) Get(index int) *Storage_Slot_ {
	// Value: &wrapper.tmplField{Type:"Slot_", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0xc000565680), Size:0, T:0x6, NumberOfBytes:0x80, StringKind:"t_struct(Slot_)128_storage", TupleRawName:"Slot_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0005a4ab0), (*abi.Type)(0xc0005a4b40), (*abi.Type)(0xc0005a4bd0)}, TupleRawNames:[]string{"EcdsaAddress", "BLSPublicKey", "EffectiveStake"}}, Tag:"", Slot:<nil>, Offset:0}
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
// &wrapper.tmplStruct{T:0x6, Name:"AddressSet_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017d0a0), (*wrapper.tmplField)(0xc00017d340)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type AddressSet_ struct {
	Keys Slice_760032d9d704 `json:"Keys" storage:"slot=0,offset=0"`
	Set  Map_3783ff483121   `json:"Set" storage:"slot=1,offset=0"`
}

// BLSPublicKey_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"BLSPublicKey_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017cb60)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type BLSPublicKey_ struct {
	Key Array_aaa956d410fd `json:"Key" storage:"slot=0,offset=0"`
}

// BLSPublicKeys_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"BLSPublicKeys_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017d5e0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type BLSPublicKeys_ struct {
	Keys Slice_6bb0ade2b43c `json:"Keys" storage:"slot=0,offset=0"`
}

// CommissionRates_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"CommissionRates_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0007180e0), (*wrapper.tmplField)(0xc0007181c0), (*wrapper.tmplField)(0xc0007182a0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type CommissionRates_ struct {
	Rate          Decimal `json:"Rate" storage:"slot=0,offset=0"`
	MaxRate       Decimal `json:"MaxRate" storage:"slot=1,offset=0"`
	MaxChangeRate Decimal `json:"MaxChangeRate" storage:"slot=2,offset=0"`
}

// Commission_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Commission_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000718380), (*wrapper.tmplField)(0xc000718460)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Commission_ struct {
	CommissionRates CommissionRates_ `json:"CommissionRates" storage:"slot=0,offset=0"`
	UpdateHeight    BigInt           `json:"UpdateHeight" storage:"slot=3,offset=0"`
}

// Committee_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Committee_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000719dc0), (*wrapper.tmplField)(0xc00017c000)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Committee_ struct {
	Epoch BigInt `json:"Epoch" storage:"slot=0,offset=0"`
	Slots Slots_ `json:"Slots" storage:"slot=1,offset=0"`
}

// Counters_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Counters_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000719880), (*wrapper.tmplField)(0xc000719960)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Counters_ struct {
	NumBlocksToSign BigInt `json:"NumBlocksToSign" storage:"slot=0,offset=0"`
	NumBlocksSigned BigInt `json:"NumBlocksSigned" storage:"slot=1,offset=0"`
}

// Description_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Description_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000718620), (*wrapper.tmplField)(0xc000718700), (*wrapper.tmplField)(0xc0007187e0), (*wrapper.tmplField)(0xc0007188c0), (*wrapper.tmplField)(0xc0007189a0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Description_ struct {
	Name            String `json:"Name" storage:"slot=0,offset=0"`
	Identity        String `json:"Identity" storage:"slot=1,offset=0"`
	Website         String `json:"Website" storage:"slot=2,offset=0"`
	SecurityContact String `json:"SecurityContact" storage:"slot=3,offset=0"`
	Details         String `json:"Details" storage:"slot=4,offset=0"`
}

// Global_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Global_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000e2000)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Global_t struct {
	ValidatorPool ValidatorPool_ `json:"ValidatorPool" storage:"slot=0,offset=0"`
}

// RedelegationMapEntry_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"RedelegationMapEntry_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000719420), (*wrapper.tmplField)(0xc000719500)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type RedelegationMapEntry_ struct {
	Entry Redelegation_ `json:"Entry" storage:"slot=0,offset=0"`
	Index BigInt        `json:"Index" storage:"slot=5,offset=0"`
}

// RedelegationMap_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"RedelegationMap_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000718d20), (*wrapper.tmplField)(0xc0007196c0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type RedelegationMap_ struct {
	Keys Slice_760032d9d704 `json:"Keys" storage:"slot=0,offset=0"`
	Map  Map_1e3b3171b974   `json:"Map" storage:"slot=1,offset=0"`
}

// Redelegation_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Redelegation_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000718ee0), (*wrapper.tmplField)(0xc000718fc0), (*wrapper.tmplField)(0xc0007190a0), (*wrapper.tmplField)(0xc000719340)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Redelegation_ struct {
	DelegatorAddress Address       `json:"DelegatorAddress" storage:"slot=0,offset=0"`
	Amount           BigInt        `json:"Amount" storage:"slot=1,offset=0"`
	Reward           BigInt        `json:"Reward" storage:"slot=2,offset=0"`
	Undelegation     Undelegation_ `json:"Undelegation" storage:"slot=3,offset=0"`
}

// Slot_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Slot_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017c9a0), (*wrapper.tmplField)(0xc00017cc40), (*wrapper.tmplField)(0xc00017cd20)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Slot_ struct {
	EcdsaAddress   Address       `json:"EcdsaAddress" storage:"slot=0,offset=0"`
	BLSPublicKey   BLSPublicKey_ `json:"BLSPublicKey" storage:"slot=1,offset=0"`
	EffectiveStake Decimal       `json:"EffectiveStake" storage:"slot=3,offset=0"`
}

// Slots_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Slots_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000719ea0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Slots_ struct {
	Entrys Slice_ed073f3f1ff1 `json:"Entrys" storage:"slot=0,offset=0"`
}

// Undelegation_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Undelegation_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000719180), (*wrapper.tmplField)(0xc000719260)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type Undelegation_ struct {
	Amount BigInt `json:"Amount" storage:"slot=0,offset=0"`
	Epoch  BigInt `json:"Epoch" storage:"slot=1,offset=0"`
}

// ValidatorPool_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"ValidatorPool_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017da40), (*wrapper.tmplField)(0xc00017dce0), (*wrapper.tmplField)(0xc00017ddc0), (*wrapper.tmplField)(0xc00017dea0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type ValidatorPool_ struct {
	Validators             ValidatorWrapperMap_ `json:"Validators" storage:"slot=0,offset=0"`
	SlotKeySet             Map_624bcc0df495     `json:"SlotKeySet" storage:"slot=2,offset=0"`
	DescriptionIdentitySet Map_624bcc0df495     `json:"DescriptionIdentitySet" storage:"slot=3,offset=0"`
	Committee              Committee_           `json:"Committee" storage:"slot=4,offset=0"`
}

// ValidatorWrapperMapEntry_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"ValidatorWrapperMapEntry_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017c1c0), (*wrapper.tmplField)(0xc00017c2a0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type ValidatorWrapperMapEntry_ struct {
	Entry ValidatorWrapper_ `json:"Entry" storage:"slot=0,offset=0"`
	Index BigInt            `json:"Index" storage:"slot=24,offset=0"`
}

// ValidatorWrapperMap_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"ValidatorWrapperMap_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017c700), (*wrapper.tmplField)(0xc00017d960)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type ValidatorWrapperMap_ struct {
	Keys Slice_760032d9d704 `json:"Keys" storage:"slot=0,offset=0"`
	Map  Map_11675a09c025   `json:"Map" storage:"slot=1,offset=0"`
}

// ValidatorWrapper_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"ValidatorWrapper_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000718c40), (*wrapper.tmplField)(0xc0007197a0), (*wrapper.tmplField)(0xc000719a40), (*wrapper.tmplField)(0xc000719b20), (*wrapper.tmplField)(0xc000719c00), (*wrapper.tmplField)(0xc000719ce0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
type ValidatorWrapper_ struct {
	Validator                 Validator_       `json:"Validator" storage:"slot=0,offset=0"`
	Redelegations             RedelegationMap_ `json:"Redelegations" storage:"slot=17,offset=0"`
	Counters                  Counters_        `json:"Counters" storage:"slot=19,offset=0"`
	BlockReward               BigInt           `json:"BlockReward" storage:"slot=21,offset=0"`
	TotalDelegation           BigInt           `json:"TotalDelegation" storage:"slot=22,offset=0"`
	TotalDelegationByOperator BigInt           `json:"TotalDelegationByOperator" storage:"slot=23,offset=0"`
}

// Validator_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Validator_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017cee0), (*wrapper.tmplField)(0xc00017d420), (*wrapper.tmplField)(0xc00017d6c0), (*wrapper.tmplField)(0xc00017d7a0), (*wrapper.tmplField)(0xc00017d880), (*wrapper.tmplField)(0xc000718000), (*wrapper.tmplField)(0xc000718540), (*wrapper.tmplField)(0xc000718a80), (*wrapper.tmplField)(0xc000718b60)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
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

	// &wrapper.tmplStruct{T:0x6, Name:"AddressSet_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017d0a0), (*wrapper.tmplField)(0xc00017d340)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *AddressSet_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_BLSPublicKey_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"BLSPublicKey_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017cb60)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *BLSPublicKey_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_BLSPublicKeys_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"BLSPublicKeys_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017d5e0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *BLSPublicKeys_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_CommissionRates_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"CommissionRates_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0007180e0), (*wrapper.tmplField)(0xc0007181c0), (*wrapper.tmplField)(0xc0007182a0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *CommissionRates_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Commission_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Commission_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000718380), (*wrapper.tmplField)(0xc000718460)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *Commission_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Committee_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Committee_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000719dc0), (*wrapper.tmplField)(0xc00017c000)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *Committee_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Counters_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Counters_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000719880), (*wrapper.tmplField)(0xc000719960)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *Counters_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Description_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Description_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000718620), (*wrapper.tmplField)(0xc000718700), (*wrapper.tmplField)(0xc0007187e0), (*wrapper.tmplField)(0xc0007188c0), (*wrapper.tmplField)(0xc0007189a0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *Description_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Global_t struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Global_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000e2000)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *Global_t
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_RedelegationMapEntry_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"RedelegationMapEntry_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000719420), (*wrapper.tmplField)(0xc000719500)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *RedelegationMapEntry_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_RedelegationMap_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"RedelegationMap_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000718d20), (*wrapper.tmplField)(0xc0007196c0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *RedelegationMap_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Redelegation_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Redelegation_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000718ee0), (*wrapper.tmplField)(0xc000718fc0), (*wrapper.tmplField)(0xc0007190a0), (*wrapper.tmplField)(0xc000719340)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *Redelegation_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Slot_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Slot_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017c9a0), (*wrapper.tmplField)(0xc00017cc40), (*wrapper.tmplField)(0xc00017cd20)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *Slot_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Slots_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Slots_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000719ea0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *Slots_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Undelegation_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Undelegation_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000719180), (*wrapper.tmplField)(0xc000719260)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *Undelegation_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_ValidatorPool_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"ValidatorPool_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017da40), (*wrapper.tmplField)(0xc00017dce0), (*wrapper.tmplField)(0xc00017ddc0), (*wrapper.tmplField)(0xc00017dea0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *ValidatorPool_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_ValidatorWrapperMapEntry_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"ValidatorWrapperMapEntry_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017c1c0), (*wrapper.tmplField)(0xc00017c2a0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *ValidatorWrapperMapEntry_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_ValidatorWrapperMap_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"ValidatorWrapperMap_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017c700), (*wrapper.tmplField)(0xc00017d960)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *ValidatorWrapperMap_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_ValidatorWrapper_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"ValidatorWrapper_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000718c40), (*wrapper.tmplField)(0xc0007197a0), (*wrapper.tmplField)(0xc000719a40), (*wrapper.tmplField)(0xc000719b20), (*wrapper.tmplField)(0xc000719c00), (*wrapper.tmplField)(0xc000719ce0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
	obj           *ValidatorWrapper_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Validator_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Validator_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc00017cee0), (*wrapper.tmplField)(0xc00017d420), (*wrapper.tmplField)(0xc00017d6c0), (*wrapper.tmplField)(0xc00017d7a0), (*wrapper.tmplField)(0xc00017d880), (*wrapper.tmplField)(0xc000718000), (*wrapper.tmplField)(0xc000718540), (*wrapper.tmplField)(0xc000718a80), (*wrapper.tmplField)(0xc000718b60)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x0, Type:reflect.Type(nil), Size:0, T:0x0, NumberOfBytes:0x0, StringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}}
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
	// Field: &wrapper.tmplField{Type:"Slice_760032d9d704", ElemType:"Slice_760032d9d704", Name:"Keys", SolKind:abi.Type{Elem:(*abi.Type)(0xc000592120), Kind:0x17, Type:(*reflect.rtype)(0xc0005652c0), Size:0, T:0x4, NumberOfBytes:0x20, StringKind:"t_array(t_address)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Keys\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Map_3783ff483121", ElemType:"Map_3783ff483121", Name:"Set", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc000570ae0), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_address,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0005922d0), (*abi.Type)(0xc000592360)}, TupleRawNames:[]string(nil)}, Tag:"json:\"Set\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Array_aaa956d410fd", ElemType:"Array_aaa956d410fd", Name:"Key", SolKind:abi.Type{Elem:(*abi.Type)(0xc0005927e0), Kind:0x11, Type:(*reflect.rtype)(0x47ebe80), Size:48, T:0x5, NumberOfBytes:0x40, StringKind:"t_array(t_bytes1)48_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Key\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Slice_6bb0ade2b43c", ElemType:"Slice_6bb0ade2b43c", Name:"Keys", SolKind:abi.Type{Elem:(*abi.Type)(0xc0005926c0), Kind:0x17, Type:(*reflect.rtype)(0xc000565340), Size:0, T:0x4, NumberOfBytes:0x20, StringKind:"t_array(t_struct(BLSPublicKey_)32_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Keys\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Decimal", ElemType:"Decimal", Name:"Rate", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000570fc0), Size:0, T:0x6, NumberOfBytes:0x20, StringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc000592f30)}, TupleRawNames:[]string{"f"}}, Tag:"json:\"Rate\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Decimal", ElemType:"Decimal", Name:"MaxRate", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000570fc0), Size:0, T:0x6, NumberOfBytes:0x20, StringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc000592f30)}, TupleRawNames:[]string{"f"}}, Tag:"json:\"MaxRate\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Decimal", ElemType:"Decimal", Name:"MaxChangeRate", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000570fc0), Size:0, T:0x6, NumberOfBytes:0x20, StringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc000592f30)}, TupleRawNames:[]string{"f"}}, Tag:"json:\"MaxChangeRate\" storage:\"slot=2,offset=0\"", Slot:2, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"CommissionRates_", ElemType:"CommissionRates_", Name:"CommissionRates", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000571020), Size:0, T:0x6, NumberOfBytes:0x60, StringKind:"t_struct(CommissionRates_)22_storage", TupleRawName:"CommissionRates_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000592ea0), (*abi.Type)(0xc000593050), (*abi.Type)(0xc0005930e0)}, TupleRawNames:[]string{"Rate", "MaxRate", "MaxChangeRate"}}, Tag:"json:\"CommissionRates\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"UpdateHeight", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x4907560), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"UpdateHeight\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Epoch", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x4907560), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Epoch\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Slots_", ElemType:"Slots_", Name:"Slots", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000571ce0), Size:0, T:0x6, NumberOfBytes:0x20, StringKind:"t_struct(Slots_)132_storage", TupleRawName:"Slots_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0005a5950)}, TupleRawNames:[]string{"Entrys"}}, Tag:"json:\"Slots\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"NumBlocksToSign", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x4907560), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"NumBlocksToSign\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"NumBlocksSigned", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x4907560), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"NumBlocksSigned\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"String", ElemType:"String", Name:"Name", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x18, Type:(*reflect.rtype)(0x47e0080), Size:0, T:0x3, NumberOfBytes:0x20, StringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Name\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"String", ElemType:"String", Name:"Identity", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x18, Type:(*reflect.rtype)(0x47e0080), Size:0, T:0x3, NumberOfBytes:0x20, StringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Identity\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"String", ElemType:"String", Name:"Website", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x18, Type:(*reflect.rtype)(0x47e0080), Size:0, T:0x3, NumberOfBytes:0x20, StringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Website\" storage:\"slot=2,offset=0\"", Slot:2, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"String", ElemType:"String", Name:"SecurityContact", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x18, Type:(*reflect.rtype)(0x47e0080), Size:0, T:0x3, NumberOfBytes:0x20, StringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"SecurityContact\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"String", ElemType:"String", Name:"Details", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x18, Type:(*reflect.rtype)(0x47e0080), Size:0, T:0x3, NumberOfBytes:0x20, StringKind:"t_string_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Details\" storage:\"slot=4,offset=0\"", Slot:4, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"ValidatorPool_", ElemType:"ValidatorPool_", Name:"ValidatorPool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000571e00), Size:0, T:0x6, NumberOfBytes:0xc0, StringKind:"t_struct(ValidatorPool_)150_storage", TupleRawName:"ValidatorPool_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0005a5200), (*abi.Type)(0xc0005a5440), (*abi.Type)(0xc0005a5710), (*abi.Type)(0xc0005a57a0)}, TupleRawNames:[]string{"Validators", "SlotKeySet", "DescriptionIdentitySet", "Committee"}}, Tag:"json:\"ValidatorPool\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Redelegation_", ElemType:"Redelegation_", Name:"Entry", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000571500), Size:0, T:0x6, NumberOfBytes:0xa0, StringKind:"t_struct(Redelegation_)82_storage", TupleRawName:"Redelegation_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000593c20), (*abi.Type)(0xc000593cb0), (*abi.Type)(0xc000593d40), (*abi.Type)(0xc000593dd0)}, TupleRawNames:[]string{"DelegatorAddress", "Amount", "Reward", "Undelegation"}}, Tag:"json:\"Entry\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Index", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x4907560), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Index\" storage:\"slot=5,offset=0\"", Slot:5, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Slice_760032d9d704", ElemType:"Slice_760032d9d704", Name:"Keys", SolKind:abi.Type{Elem:(*abi.Type)(0xc000592120), Kind:0x17, Type:(*reflect.rtype)(0xc0005652c0), Size:0, T:0x4, NumberOfBytes:0x20, StringKind:"t_array(t_address)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Keys\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Map_1e3b3171b974", ElemType:"Map_1e3b3171b974", Name:"Map", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc0005715c0), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_address,t_struct(RedelegationMapEntry_)87_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc000593a70), (*abi.Type)(0xc000593b00)}, TupleRawNames:[]string(nil)}, Tag:"json:\"Map\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Address", ElemType:"Address", Name:"DelegatorAddress", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x11, Type:(*reflect.rtype)(0x48ba1e0), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"DelegatorAddress\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Amount", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x4907560), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Amount\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Reward", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x4907560), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Reward\" storage:\"slot=2,offset=0\"", Slot:2, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Undelegation_", ElemType:"Undelegation_", Name:"Undelegation", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000571440), Size:0, T:0x6, NumberOfBytes:0x40, StringKind:"t_struct(Undelegation_)73_storage", TupleRawName:"Undelegation_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000593e60), (*abi.Type)(0xc000593ef0)}, TupleRawNames:[]string{"Amount", "Epoch"}}, Tag:"json:\"Undelegation\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Address", ElemType:"Address", Name:"EcdsaAddress", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x11, Type:(*reflect.rtype)(0x48ba1e0), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"EcdsaAddress\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"BLSPublicKey_", ElemType:"BLSPublicKey_", Name:"BLSPublicKey", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000570d80), Size:0, T:0x6, NumberOfBytes:0x40, StringKind:"t_struct(BLSPublicKey_)32_storage", TupleRawName:"BLSPublicKey_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000592750)}, TupleRawNames:[]string{"Key"}}, Tag:"json:\"BLSPublicKey\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Decimal", ElemType:"Decimal", Name:"EffectiveStake", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000570fc0), Size:0, T:0x6, NumberOfBytes:0x20, StringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc000592f30)}, TupleRawNames:[]string{"f"}}, Tag:"json:\"EffectiveStake\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Slice_ed073f3f1ff1", ElemType:"Slice_ed073f3f1ff1", Name:"Entrys", SolKind:abi.Type{Elem:(*abi.Type)(0xc0005a59e0), Kind:0x17, Type:(*reflect.rtype)(0xc0005656c0), Size:0, T:0x4, NumberOfBytes:0x20, StringKind:"t_array(t_struct(Slot_)128_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Entrys\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Amount", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x4907560), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Amount\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Epoch", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x4907560), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Epoch\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"ValidatorWrapperMap_", ElemType:"ValidatorWrapperMap_", Name:"Validators", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000571a40), Size:0, T:0x6, NumberOfBytes:0x40, StringKind:"t_struct(ValidatorWrapperMap_)121_storage", TupleRawName:"ValidatorWrapperMap_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0005a5290), (*abi.Type)(0xc0005a5320)}, TupleRawNames:[]string{"Keys", "Map"}}, Tag:"json:\"Validators\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Map_624bcc0df495", ElemType:"Map_624bcc0df495", Name:"SlotKeySet", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc000571b60), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_string_memory_ptr,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0005a54d0), (*abi.Type)(0xc0005a55f0)}, TupleRawNames:[]string(nil)}, Tag:"json:\"SlotKeySet\" storage:\"slot=2,offset=0\"", Slot:2, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Map_624bcc0df495", ElemType:"Map_624bcc0df495", Name:"DescriptionIdentitySet", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc000571b60), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_string_memory_ptr,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0005a54d0), (*abi.Type)(0xc0005a55f0)}, TupleRawNames:[]string(nil)}, Tag:"json:\"DescriptionIdentitySet\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Committee_", ElemType:"Committee_", Name:"Committee", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000571d40), Size:0, T:0x6, NumberOfBytes:0x40, StringKind:"t_struct(Committee_)137_storage", TupleRawName:"Committee_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0005a5830), (*abi.Type)(0xc0005a58c0)}, TupleRawNames:[]string{"Epoch", "Slots"}}, Tag:"json:\"Committee\" storage:\"slot=4,offset=0\"", Slot:4, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"ValidatorWrapper_", ElemType:"ValidatorWrapper_", Name:"Entry", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000571740), Size:0, T:0x6, NumberOfBytes:0x300, StringKind:"t_struct(ValidatorWrapper_)108_storage", TupleRawName:"ValidatorWrapper_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00058be60), (*abi.Type)(0xc0005938c0), (*abi.Type)(0xc0005a4360), (*abi.Type)(0xc0005a45a0), (*abi.Type)(0xc0005a4630), (*abi.Type)(0xc0005a46c0)}, TupleRawNames:[]string{"Validator", "Redelegations", "Counters", "BlockReward", "TotalDelegation", "TotalDelegationByOperator"}}, Tag:"json:\"Entry\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Index", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x4907560), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Index\" storage:\"slot=24,offset=0\"", Slot:24, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Slice_760032d9d704", ElemType:"Slice_760032d9d704", Name:"Keys", SolKind:abi.Type{Elem:(*abi.Type)(0xc000592120), Kind:0x17, Type:(*reflect.rtype)(0xc0005652c0), Size:0, T:0x4, NumberOfBytes:0x20, StringKind:"t_array(t_address)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Keys\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Map_11675a09c025", ElemType:"Map_11675a09c025", Name:"Map", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x15, Type:(*reflect.rtype)(0xc000571800), Size:0, T:0xd, NumberOfBytes:0x20, StringKind:"t_mapping(t_address,t_struct(ValidatorWrapperMapEntry_)113_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc00058bc20), (*abi.Type)(0xc00058bd40)}, TupleRawNames:[]string(nil)}, Tag:"json:\"Map\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Validator_", ElemType:"Validator_", Name:"Validator", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000571200), Size:0, T:0x6, NumberOfBytes:0x220, StringKind:"t_struct(Validator_)68_storage", TupleRawName:"Validator_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00058bef0), (*abi.Type)(0xc000592000), (*abi.Type)(0xc0005925a0), (*abi.Type)(0xc000592b40), (*abi.Type)(0xc000592c60), (*abi.Type)(0xc000592cf0), (*abi.Type)(0xc000592d80), (*abi.Type)(0xc000593320), (*abi.Type)(0xc0005937a0)}, TupleRawNames:[]string{"ValidatorAddress", "OperatorAddresses", "SlotPubKeys", "LastEpochInCommittee", "MaxTotalDelegation", "Status", "Commission", "Description", "CreationHeight"}}, Tag:"json:\"Validator\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"RedelegationMap_", ElemType:"RedelegationMap_", Name:"Redelegations", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000571620), Size:0, T:0x6, NumberOfBytes:0x40, StringKind:"t_struct(RedelegationMap_)95_storage", TupleRawName:"RedelegationMap_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000593950), (*abi.Type)(0xc0005939e0)}, TupleRawNames:[]string{"Keys", "Map"}}, Tag:"json:\"Redelegations\" storage:\"slot=17,offset=0\"", Slot:17, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Counters_", ElemType:"Counters_", Name:"Counters", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc0005716e0), Size:0, T:0x6, NumberOfBytes:0x40, StringKind:"t_struct(Counters_)41_storage", TupleRawName:"Counters_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0005a43f0), (*abi.Type)(0xc0005a4480)}, TupleRawNames:[]string{"NumBlocksToSign", "NumBlocksSigned"}}, Tag:"json:\"Counters\" storage:\"slot=19,offset=0\"", Slot:19, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"BlockReward", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x4907560), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"BlockReward\" storage:\"slot=21,offset=0\"", Slot:21, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"TotalDelegation", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x4907560), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"TotalDelegation\" storage:\"slot=22,offset=0\"", Slot:22, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"TotalDelegationByOperator", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x4907560), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"TotalDelegationByOperator\" storage:\"slot=23,offset=0\"", Slot:23, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Address", ElemType:"Address", Name:"ValidatorAddress", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x11, Type:(*reflect.rtype)(0x48ba1e0), Size:20, T:0x7, NumberOfBytes:0x14, StringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"ValidatorAddress\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"AddressSet_", ElemType:"AddressSet_", Name:"OperatorAddresses", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000570b40), Size:0, T:0x6, NumberOfBytes:0x40, StringKind:"t_struct(AddressSet_)49_storage", TupleRawName:"AddressSet_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000592090), (*abi.Type)(0xc000592240)}, TupleRawNames:[]string{"Keys", "Set"}}, Tag:"json:\"OperatorAddresses\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"BLSPublicKeys_", ElemType:"BLSPublicKeys_", Name:"SlotPubKeys", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000570de0), Size:0, T:0x6, NumberOfBytes:0x20, StringKind:"t_struct(BLSPublicKeys_)36_storage", TupleRawName:"BLSPublicKeys_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000592630)}, TupleRawNames:[]string{"Keys"}}, Tag:"json:\"SlotPubKeys\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"LastEpochInCommittee", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x4907560), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"LastEpochInCommittee\" storage:\"slot=4,offset=0\"", Slot:4, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"MaxTotalDelegation", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x4907560), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"MaxTotalDelegation\" storage:\"slot=5,offset=0\"", Slot:5, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Uint8", ElemType:"Uint8", Name:"Status", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x8, Type:(*reflect.rtype)(0x47e0300), Size:0, T:0x1, NumberOfBytes:0x1, StringKind:"uint8", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"Status\" storage:\"slot=6,offset=0\"", Slot:6, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Commission_", ElemType:"Commission_", Name:"Commission", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc000571080), Size:0, T:0x6, NumberOfBytes:0x80, StringKind:"t_struct(Commission_)27_storage", TupleRawName:"Commission_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000592e10), (*abi.Type)(0xc000593200)}, TupleRawNames:[]string{"CommissionRates", "UpdateHeight"}}, Tag:"json:\"Commission\" storage:\"slot=7,offset=0\"", Slot:7, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Description_", ElemType:"Description_", Name:"Description", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x19, Type:(*reflect.rtype)(0xc0005711a0), Size:0, T:0x6, NumberOfBytes:0xa0, StringKind:"t_struct(Description_)15_storage", TupleRawName:"Description_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0005933b0), (*abi.Type)(0xc0005934d0), (*abi.Type)(0xc000593560), (*abi.Type)(0xc0005935f0), (*abi.Type)(0xc000593680)}, TupleRawNames:[]string{"Name", "Identity", "Website", "SecurityContact", "Details"}}, Tag:"json:\"Description\" storage:\"slot=11,offset=0\"", Slot:11, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"CreationHeight", SolKind:abi.Type{Elem:(*abi.Type)(nil), Kind:0x16, Type:(*reflect.rtype)(0x4907560), Size:256, T:0x1, NumberOfBytes:0x20, StringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil)}, Tag:"json:\"CreationHeight\" storage:\"slot=16,offset=0\"", Slot:16, Offset:0}
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
