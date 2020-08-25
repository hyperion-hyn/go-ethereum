// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package microstaking

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

type Decimal = common.Dec
type Storage_Decimal struct {
	// Builtin-Type

	// &wrapper.tmplStruct{T:0xf, Name:"Decimal", Fields:[]*wrapper.tmplField(nil), Type:"common.Dec", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xf, numberOfBytes:0x20, stringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032da00)}, TupleRawNames:[]string{"f"}, TupleType:(*reflect.rtype)(0xc00047c240)}}
	obj           *Decimal
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Decimal) Value() common.Dec {
	// &wrapper.tmplStruct{T:0xf, Name:"Decimal", Fields:[]*wrapper.tmplField(nil), Type:"common.Dec", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xf, numberOfBytes:0x20, stringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032da00)}, TupleRawNames:[]string{"f"}, TupleType:(*reflect.rtype)(0xc00047c240)}}
	hash := s.db.GetState(s.addr, common.BigToHash(s.slot))
	*s.obj = common.NewDecFromBigIntWithPrec(hash.Big(), common.Precision)
	return *s.obj
}

func (s *Storage_Decimal) SetValue(value common.Dec) {
	// &wrapper.tmplStruct{T:0xf, Name:"Decimal", Fields:[]*wrapper.tmplField(nil), Type:"common.Dec", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xf, numberOfBytes:0x20, stringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032da00)}, TupleRawNames:[]string{"f"}, TupleType:(*reflect.rtype)(0xc00047c240)}}
	hash := value.BigInt()
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(hash))
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
// &wrapper.tmplStruct{T:0x5, Name:"Array_aaa956d410fd", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0002e45b0)}, Type:"[48]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc00032d000), Size:48, T:0x5, numberOfBytes:0x40, stringKind:"t_array(t_bytes1)48_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "[48]uint8"
type Array_aaa956d410fd [48]Uint8

type Storage_Array_aaa956d410fd struct {

	// &wrapper.tmplStruct{T:0x5, Name:"Array_aaa956d410fd", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0002e45b0)}, Type:"[48]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc00032d000), Size:48, T:0x5, numberOfBytes:0x40, stringKind:"t_array(t_bytes1)48_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           *Array_aaa956d410fd
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Array_aaa956d410fd) Value() [48]Uint8 {
	// &wrapper.tmplStruct{T:0x5, Name:"Array_aaa956d410fd", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0002e45b0)}, Type:"[48]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc00032d000), Size:48, T:0x5, numberOfBytes:0x40, stringKind:"t_array(t_bytes1)48_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
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
	// &wrapper.tmplStruct{T:0x5, Name:"Array_aaa956d410fd", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0002e45b0)}, Type:"[48]Uint8", SolKind:abi.Type{Elem:(*abi.Type)(0xc00032d000), Size:48, T:0x5, numberOfBytes:0x40, stringKind:"t_array(t_bytes1)48_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
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

// Map_4dcba4f803c2 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xe, Name:"Map_4dcba4f803c2", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000badd0), (*wrapper.tmplField)(0xc0000bba00)}, Type:"map[Address]*Map3NodeWrapperMapEntry_", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(Map3NodeWrapperMapEntry_)107_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032d380), (*abi.Type)(0xc000481700)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "map[common.Address]*struct { Entry struct { Map3Node struct { NodeAddress common.Address \"json:\\\"NodeAddress\\\" storage:\\\"slot=0,offset=0\\\"\"; OperatorAddress common.Address \"json:\\\"OperatorAddress\\\" storage:\\\"slot=1,offset=0\\\"\"; PubKeys struct { Keys []*struct { Key [48]uint8 \"json:\\\"Key\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"Keys\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"PubKeys\\\" storage:\\\"slot=2,offset=0\\\"\"; Commission struct { F *big.Int \"json:\\\"f\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"Commission\\\" storage:\\\"slot=3,offset=0\\\"\"; Description struct { Name string \"json:\\\"Name\\\" storage:\\\"slot=0,offset=0\\\"\"; Identity string \"json:\\\"Identity\\\" storage:\\\"slot=1,offset=0\\\"\"; Website string \"json:\\\"Website\\\" storage:\\\"slot=2,offset=0\\\"\"; SecurityContact string \"json:\\\"SecurityContact\\\" storage:\\\"slot=3,offset=0\\\"\"; Details string \"json:\\\"Details\\\" storage:\\\"slot=4,offset=0\\\"\" } \"json:\\\"Description\\\" storage:\\\"slot=4,offset=0\\\"\"; CreationHeight *big.Int \"json:\\\"CreationHeight\\\" storage:\\\"slot=9,offset=0\\\"\"; Age struct { F *big.Int \"json:\\\"f\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"Age\\\" storage:\\\"slot=10,offset=0\\\"\"; Status uint8 \"json:\\\"Status\\\" storage:\\\"slot=11,offset=0\\\"\"; ActivationEpoch *big.Int \"json:\\\"ActivationEpoch\\\" storage:\\\"slot=12,offset=0\\\"\"; ReleaseEpoch struct { F *big.Int \"json:\\\"f\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"ReleaseEpoch\\\" storage:\\\"slot=13,offset=0\\\"\" } \"json:\\\"Map3Node\\\" storage:\\\"slot=0,offset=0\\\"\"; Microdelegations struct { Keys []*common.Address \"json:\\\"Keys\\\" storage:\\\"slot=0,offset=0\\\"\"; Map map[common.Address]*struct { Entry struct { DelegatorAddress common.Address \"json:\\\"DelegatorAddress\\\" storage:\\\"slot=0,offset=0\\\"\"; Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=1,offset=0\\\"\"; Reward *big.Int \"json:\\\"Reward\\\" storage:\\\"slot=2,offset=0\\\"\"; Undelegation struct { Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=0,offset=0\\\"\"; Epoch *big.Int \"json:\\\"Epoch\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"Undelegation\\\" storage:\\\"slot=3,offset=0\\\"\"; PendingDelegation struct { Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=0,offset=0\\\"\"; UnlockedEpoch struct { F *big.Int \"json:\\\"f\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"UnlockedEpoch\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"PendingDelegation\\\" storage:\\\"slot=5,offset=0\\\"\"; Renewal struct { IsRenew bool \"json:\\\"IsRenew\\\" storage:\\\"slot=0,offset=0\\\"\"; UpdateHeight *big.Int \"json:\\\"UpdateHeight\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"Renewal\\\" storage:\\\"slot=7,offset=0\\\"\" } \"json:\\\"Entry\\\" storage:\\\"slot=0,offset=0\\\"\"; Index *big.Int \"json:\\\"Index\\\" storage:\\\"slot=9,offset=0\\\"\" } \"json:\\\"Map\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"Microdelegations\\\" storage:\\\"slot=14,offset=0\\\"\"; RedelegationReference struct { ValidatorAddress common.Address \"json:\\\"ValidatorAddress\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"RedelegationReference\\\" storage:\\\"slot=16,offset=0\\\"\"; AccumulatedReward *big.Int \"json:\\\"AccumulatedReward\\\" storage:\\\"slot=17,offset=0\\\"\"; TotalDelegation *big.Int \"json:\\\"TotalDelegation\\\" storage:\\\"slot=18,offset=0\\\"\"; TotalPendingDelegation *big.Int \"json:\\\"TotalPendingDelegation\\\" storage:\\\"slot=19,offset=0\\\"\" } \"json:\\\"Entry\\\" storage:\\\"slot=0,offset=0\\\"\"; Index *big.Int \"json:\\\"Index\\\" storage:\\\"slot=20,offset=0\\\"\" }"
type Map_4dcba4f803c2 map[Address]*Map3NodeWrapperMapEntry_

type Storage_Map_4dcba4f803c2 struct {

	// &wrapper.tmplStruct{T:0xe, Name:"Map_4dcba4f803c2", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000badd0), (*wrapper.tmplField)(0xc0000bba00)}, Type:"map[Address]*Map3NodeWrapperMapEntry_", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(Map3NodeWrapperMapEntry_)107_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032d380), (*abi.Type)(0xc000481700)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           Map_4dcba4f803c2
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Map_4dcba4f803c2) Get(key Address) *Storage_Map3NodeWrapperMapEntry_ {
	// Key:	&wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"Map3NodeWrapperMapEntry_", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc00032d400), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
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
		numberOfBytes: 672,
		dirty:         s.dirty,
	}
}

// Map_624bcc0df495 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xe, Name:"Map_624bcc0df495", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000bbad0), (*wrapper.tmplField)(0xc0000bbc70)}, Type:"map[String]*Bool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_string_memory_ptr,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fa600), (*abi.Type)(0xc0000fa700)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "map[string]*bool"
type Map_624bcc0df495 map[String]*Bool

type Storage_Map_624bcc0df495 struct {

	// &wrapper.tmplStruct{T:0xe, Name:"Map_624bcc0df495", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000bbad0), (*wrapper.tmplField)(0xc0000bbc70)}, Type:"map[String]*Bool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_string_memory_ptr,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fa600), (*abi.Type)(0xc0000fa700)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
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
	// Value: &wrapper.tmplField{Type:"Bool", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000fa680), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
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

// Map_64b0ddf77b70 is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xe, Name:"Map_64b0ddf77b70", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000bbd40), (*wrapper.tmplField)(0xc0001221a0)}, Type:"map[Address]*DelegationIndexMapEntry_", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(DelegationIndexMapEntry_)125_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032c900), (*abi.Type)(0xc00032ce00)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "map[common.Address]*struct { Entry struct { NodeAddress common.Address \"json:\\\"NodeAddress\\\" storage:\\\"slot=0,offset=0\\\"\"; IsOperator bool \"json:\\\"IsOperator\\\" storage:\\\"slot=0,offset=20\\\"\" } \"json:\\\"Entry\\\" storage:\\\"slot=0,offset=0\\\"\"; Index *big.Int \"json:\\\"Index\\\" storage:\\\"slot=1,offset=0\\\"\" }"
type Map_64b0ddf77b70 map[Address]*DelegationIndexMapEntry_

type Storage_Map_64b0ddf77b70 struct {

	// &wrapper.tmplStruct{T:0xe, Name:"Map_64b0ddf77b70", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000bbd40), (*wrapper.tmplField)(0xc0001221a0)}, Type:"map[Address]*DelegationIndexMapEntry_", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(DelegationIndexMapEntry_)125_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032c900), (*abi.Type)(0xc00032ce00)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           Map_64b0ddf77b70
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Map_64b0ddf77b70) Get(key Address) *Storage_DelegationIndexMapEntry_ {
	// Key:	&wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"DelegationIndexMapEntry_", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc00032c980), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
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

// Map_67e3f4d129ef is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xe, Name:"Map_67e3f4d129ef", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0002e5d40), (*wrapper.tmplField)(0xc0000baa90)}, Type:"map[Address]*MicrodelegationMapEntry_", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(MicrodelegationMapEntry_)78_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc000480580), (*abi.Type)(0xc000481080)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "map[common.Address]*struct { Entry struct { DelegatorAddress common.Address \"json:\\\"DelegatorAddress\\\" storage:\\\"slot=0,offset=0\\\"\"; Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=1,offset=0\\\"\"; Reward *big.Int \"json:\\\"Reward\\\" storage:\\\"slot=2,offset=0\\\"\"; Undelegation struct { Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=0,offset=0\\\"\"; Epoch *big.Int \"json:\\\"Epoch\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"Undelegation\\\" storage:\\\"slot=3,offset=0\\\"\"; PendingDelegation struct { Amount *big.Int \"json:\\\"Amount\\\" storage:\\\"slot=0,offset=0\\\"\"; UnlockedEpoch struct { F *big.Int \"json:\\\"f\\\" storage:\\\"slot=0,offset=0\\\"\" } \"json:\\\"UnlockedEpoch\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"PendingDelegation\\\" storage:\\\"slot=5,offset=0\\\"\"; Renewal struct { IsRenew bool \"json:\\\"IsRenew\\\" storage:\\\"slot=0,offset=0\\\"\"; UpdateHeight *big.Int \"json:\\\"UpdateHeight\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"Renewal\\\" storage:\\\"slot=7,offset=0\\\"\" } \"json:\\\"Entry\\\" storage:\\\"slot=0,offset=0\\\"\"; Index *big.Int \"json:\\\"Index\\\" storage:\\\"slot=9,offset=0\\\"\" }"
type Map_67e3f4d129ef map[Address]*MicrodelegationMapEntry_

type Storage_Map_67e3f4d129ef struct {

	// &wrapper.tmplStruct{T:0xe, Name:"Map_67e3f4d129ef", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0002e5d40), (*wrapper.tmplField)(0xc0000baa90)}, Type:"map[Address]*MicrodelegationMapEntry_", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(MicrodelegationMapEntry_)78_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc000480580), (*abi.Type)(0xc000481080)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           Map_67e3f4d129ef
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Map_67e3f4d129ef) Get(key Address) *Storage_MicrodelegationMapEntry_ {
	// Key:	&wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"MicrodelegationMapEntry_", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc000480600), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
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

// Map_bc58594e503f is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0xe, Name:"Map_bc58594e503f", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000122270), (*wrapper.tmplField)(0xc0001225b0)}, Type:"map[Address]*DelegationIndexMap_", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(DelegationIndexMap_)133_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fa000), (*abi.Type)(0xc0000fa280)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "map[common.Address]*struct { Keys []*common.Address \"json:\\\"Keys\\\" storage:\\\"slot=0,offset=0\\\"\"; Map map[common.Address]*struct { Entry struct { NodeAddress common.Address \"json:\\\"NodeAddress\\\" storage:\\\"slot=0,offset=0\\\"\"; IsOperator bool \"json:\\\"IsOperator\\\" storage:\\\"slot=0,offset=20\\\"\" } \"json:\\\"Entry\\\" storage:\\\"slot=0,offset=0\\\"\"; Index *big.Int \"json:\\\"Index\\\" storage:\\\"slot=1,offset=0\\\"\" } \"json:\\\"Map\\\" storage:\\\"slot=1,offset=0\\\"\" }"
type Map_bc58594e503f map[Address]*DelegationIndexMap_

type Storage_Map_bc58594e503f struct {

	// &wrapper.tmplStruct{T:0xe, Name:"Map_bc58594e503f", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000122270), (*wrapper.tmplField)(0xc0001225b0)}, Type:"map[Address]*DelegationIndexMap_", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(DelegationIndexMap_)133_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fa000), (*abi.Type)(0xc0000fa280)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
	obj           Map_bc58594e503f
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

func (s *Storage_Map_bc58594e503f) Get(key Address) *Storage_DelegationIndexMap_ {
	// Key:	&wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	// Value: &wrapper.tmplField{Type:"DelegationIndexMap_", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc0000fa080), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
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

// Slice_2f1765e04dad is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x4, Name:"Slice_2f1765e04dad", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0002e4750)}, Type:"[]*BLSPublicKey_", SolKind:abi.Type{Elem:(*abi.Type)(0xc00032d800), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_struct(BLSPublicKey_)20_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "[]*struct { Key [48]uint8 \"json:\\\"Key\\\" storage:\\\"slot=0,offset=0\\\"\" }"
type Slice_2f1765e04dad []*BLSPublicKey_

type Storage_Slice_2f1765e04dad struct {

	// &wrapper.tmplStruct{T:0x4, Name:"Slice_2f1765e04dad", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0002e4750)}, Type:"[]*BLSPublicKey_", SolKind:abi.Type{Elem:(*abi.Type)(0xc00032d800), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_struct(BLSPublicKey_)20_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
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
	// Value: &wrapper.tmplField{Type:"BLSPublicKey_", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc00032d780), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(big.NewInt(0).SetUint64(uint64(length))))

	_length := length
	_capacity := length + 50

	slice := make([]*BLSPublicKey_, _length, _capacity)

	copy(slice, *s.obj)
	*s.obj = slice
}

func (s *Storage_Slice_2f1765e04dad) Get(index int) *Storage_BLSPublicKey_ {
	// Value: &wrapper.tmplField{Type:"BLSPublicKey_", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc00032d780), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
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
// &wrapper.tmplStruct{T:0x4, Name:"Slice_760032d9d704", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0002e5930)}, Type:"[]*Address", SolKind:abi.Type{Elem:(*abi.Type)(0xc00032c800), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_address)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
// "[]*common.Address"
type Slice_760032d9d704 []*Address

type Storage_Slice_760032d9d704 struct {

	// &wrapper.tmplStruct{T:0x4, Name:"Slice_760032d9d704", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0002e5930)}, Type:"[]*Address", SolKind:abi.Type{Elem:(*abi.Type)(0xc00032c800), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_address)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}}
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
	// Value: &wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc00032c700), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
	s.db.SetState(s.addr, common.BigToHash(s.slot), common.BigToHash(big.NewInt(0).SetUint64(uint64(length))))

	_length := length
	_capacity := length + 50

	slice := make([]*Address, _length, _capacity)

	copy(slice, *s.obj)
	*s.obj = slice
}

func (s *Storage_Slice_760032d9d704) Get(index int) *Storage_Address {
	// Value: &wrapper.tmplField{Type:"Address", ElemType:"", Name:"", SolKind:abi.Type{Elem:(*abi.Type)(0xc00032c700), Size:0, T:0xd, numberOfBytes:0x0, stringKind:"", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"", Slot:<nil>, Offset:0}
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
// &wrapper.tmplStruct{T:0x6, Name:"BLSPublicKey_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0002e4680)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(BLSPublicKey_)20_storage", TupleRawName:"BLSPublicKey_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032cf80)}, TupleRawNames:[]string{"Key"}, TupleType:(*reflect.rtype)(0xc0001cfe00)}}
type BLSPublicKey_ struct {
	Key Array_aaa956d410fd `json:"Key" storage:"slot=0,offset=0"`
}

// BLSPublicKeys_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"BLSPublicKeys_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0002e4820)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(BLSPublicKeys_)24_storage", TupleRawName:"BLSPublicKeys_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032d700)}, TupleRawNames:[]string{"Keys"}, TupleType:(*reflect.rtype)(0xc00047c180)}}
type BLSPublicKeys_ struct {
	Keys Slice_2f1765e04dad `json:"Keys" storage:"slot=0,offset=0"`
}

// DelegationIndexMapEntry_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"DelegationIndexMapEntry_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000122000), (*wrapper.tmplField)(0xc0001220d0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(DelegationIndexMapEntry_)125_storage", TupleRawName:"DelegationIndexMapEntry_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032ca00), (*abi.Type)(0xc00032cc80)}, TupleRawNames:[]string{"Entry", "Index"}, TupleType:(*reflect.rtype)(0xc0001cfc20)}}
type DelegationIndexMapEntry_ struct {
	Entry DelegationIndex_ `json:"Entry" storage:"slot=0,offset=0"`
	Index BigInt           `json:"Index" storage:"slot=1,offset=0"`
}

// DelegationIndexMap_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"DelegationIndexMap_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000122340), (*wrapper.tmplField)(0xc0001224e0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(DelegationIndexMap_)133_storage", TupleRawName:"DelegationIndexMap_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fa100), (*abi.Type)(0xc0000fa180)}, TupleRawNames:[]string{"Keys", "Map"}, TupleType:(*reflect.rtype)(0xc0003223c0)}}
type DelegationIndexMap_ struct {
	Keys Slice_760032d9d704 `json:"Keys" storage:"slot=0,offset=0"`
	Map  Map_64b0ddf77b70   `json:"Map" storage:"slot=1,offset=0"`
}

// DelegationIndex_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"DelegationIndex_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000bbe10), (*wrapper.tmplField)(0xc0000bbee0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(DelegationIndex_)120_storage", TupleRawName:"DelegationIndex_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032ca80), (*abi.Type)(0xc00032cb00)}, TupleRawNames:[]string{"NodeAddress", "IsOperator"}, TupleType:(*reflect.rtype)(0xc0001cfb60)}}
type DelegationIndex_ struct {
	NodeAddress Address `json:"NodeAddress" storage:"slot=0,offset=0"`
	IsOperator  Bool    `json:"IsOperator" storage:"slot=0,offset=20"`
}

// Description_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Description_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0002e4a90), (*wrapper.tmplField)(0xc0002e4b60), (*wrapper.tmplField)(0xc0002e4c30), (*wrapper.tmplField)(0xc0002e4d00), (*wrapper.tmplField)(0xc0002e4dd0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xa0, stringKind:"t_struct(Description_)15_storage", TupleRawName:"Description_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032db80), (*abi.Type)(0xc00032dc80), (*abi.Type)(0xc00032dd00), (*abi.Type)(0xc00032dd80), (*abi.Type)(0xc00032de00)}, TupleRawNames:[]string{"Name", "Identity", "Website", "SecurityContact", "Details"}, TupleType:(*reflect.rtype)(0xc00047c360)}}
type Description_ struct {
	Name            String `json:"Name" storage:"slot=0,offset=0"`
	Identity        String `json:"Identity" storage:"slot=1,offset=0"`
	Website         String `json:"Website" storage:"slot=2,offset=0"`
	SecurityContact String `json:"SecurityContact" storage:"slot=3,offset=0"`
	Details         String `json:"Details" storage:"slot=4,offset=0"`
}

// Global_t is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Global_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000122c30)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x0, stringKind:"t_struct(Global_t)", TupleRawName:"Global_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fa900)}, TupleRawNames:[]string{"Map3NodePool"}, TupleType:(*reflect.rtype)(0xc0003225a0)}}
type Global_t struct {
	Map3NodePool Map3NodePool_ `json:"Map3NodePool" storage:"slot=0,offset=0"`
}

// Map3NodePool_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Map3NodePool_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0001228f0), (*wrapper.tmplField)(0xc0001229c0), (*wrapper.tmplField)(0xc000122a90), (*wrapper.tmplField)(0xc000122b60)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xa0, stringKind:"t_struct(Map3NodePool_)148_storage", TupleRawName:"Map3NodePool_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032d200), (*abi.Type)(0xc00052a080), (*abi.Type)(0xc0000fa380), (*abi.Type)(0xc0000fa800)}, TupleRawNames:[]string{"Nodes", "DelegationIndexMapByDelegator", "SlotKeySet", "DescriptionIdentitySet"}, TupleType:(*reflect.rtype)(0xc0003224e0)}}
type Map3NodePool_ struct {
	Nodes                         Map3NodeWrapperMap_ `json:"Nodes" storage:"slot=0,offset=0"`
	DelegationIndexMapByDelegator Map_bc58594e503f    `json:"DelegationIndexMapByDelegator" storage:"slot=2,offset=0"`
	SlotKeySet                    Map_624bcc0df495    `json:"SlotKeySet" storage:"slot=3,offset=0"`
	DescriptionIdentitySet        Map_624bcc0df495    `json:"DescriptionIdentitySet" storage:"slot=4,offset=0"`
}

// Map3NodeWrapperMapEntry_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Map3NodeWrapperMapEntry_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000bb790), (*wrapper.tmplField)(0xc0000bb860)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x2a0, stringKind:"t_struct(Map3NodeWrapperMapEntry_)107_storage", TupleRawName:"Map3NodeWrapperMapEntry_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032d480), (*abi.Type)(0xc000481600)}, TupleRawNames:[]string{"Entry", "Index"}, TupleType:(*reflect.rtype)(0xc00047ca20)}}
type Map3NodeWrapperMapEntry_ struct {
	Entry Map3NodeWrapper_ `json:"Entry" storage:"slot=0,offset=0"`
	Index BigInt           `json:"Index" storage:"slot=20,offset=0"`
}

// Map3NodeWrapperMap_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Map3NodeWrapperMap_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000122750), (*wrapper.tmplField)(0xc000122820)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Map3NodeWrapperMap_)115_storage", TupleRawName:"Map3NodeWrapperMap_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032d280), (*abi.Type)(0xc00032d300)}, TupleRawNames:[]string{"Keys", "Map"}, TupleType:(*reflect.rtype)(0xc0001080c0)}}
type Map3NodeWrapperMap_ struct {
	Keys Slice_760032d9d704 `json:"Keys" storage:"slot=0,offset=0"`
	Map  Map_4dcba4f803c2   `json:"Map" storage:"slot=1,offset=0"`
}

// Map3NodeWrapper_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Map3NodeWrapper_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000bb1e0), (*wrapper.tmplField)(0xc0000bb2b0), (*wrapper.tmplField)(0xc0000bb450), (*wrapper.tmplField)(0xc0000bb520), (*wrapper.tmplField)(0xc0000bb5f0), (*wrapper.tmplField)(0xc0000bb6c0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x280, stringKind:"t_struct(Map3NodeWrapper_)102_storage", TupleRawName:"Map3NodeWrapper_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032d500), (*abi.Type)(0xc000480400), (*abi.Type)(0xc000481200), (*abi.Type)(0xc000481380), (*abi.Type)(0xc000481400), (*abi.Type)(0xc000481500)}, TupleRawNames:[]string{"Map3Node", "Microdelegations", "RedelegationReference", "AccumulatedReward", "TotalDelegation", "TotalPendingDelegation"}, TupleType:(*reflect.rtype)(0xc00047c9c0)}}
type Map3NodeWrapper_ struct {
	Map3Node               Map3Node_              `json:"Map3Node" storage:"slot=0,offset=0"`
	Microdelegations       MicrodelegationMap_    `json:"Microdelegations" storage:"slot=14,offset=0"`
	RedelegationReference  RedelegationReference_ `json:"RedelegationReference" storage:"slot=16,offset=0"`
	AccumulatedReward      BigInt                 `json:"AccumulatedReward" storage:"slot=17,offset=0"`
	TotalDelegation        BigInt                 `json:"TotalDelegation" storage:"slot=18,offset=0"`
	TotalPendingDelegation BigInt                 `json:"TotalPendingDelegation" storage:"slot=19,offset=0"`
}

// Map3Node_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Map3Node_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0002e4410), (*wrapper.tmplField)(0xc0002e44e0), (*wrapper.tmplField)(0xc0002e48f0), (*wrapper.tmplField)(0xc0002e49c0), (*wrapper.tmplField)(0xc0002e4ea0), (*wrapper.tmplField)(0xc0002e4f70), (*wrapper.tmplField)(0xc0002e5040), (*wrapper.tmplField)(0xc0002e5110), (*wrapper.tmplField)(0xc0002e51e0), (*wrapper.tmplField)(0xc0002e52b0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x1c0, stringKind:"t_struct(Map3Node_)45_storage", TupleRawName:"Map3Node_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032d580), (*abi.Type)(0xc00032d600), (*abi.Type)(0xc00032d680), (*abi.Type)(0xc00032d980), (*abi.Type)(0xc00032db00), (*abi.Type)(0xc000480080), (*abi.Type)(0xc000480100), (*abi.Type)(0xc000480180), (*abi.Type)(0xc000480200), (*abi.Type)(0xc000480300)}, TupleRawNames:[]string{"NodeAddress", "OperatorAddress", "PubKeys", "Commission", "Description", "CreationHeight", "Age", "Status", "ActivationEpoch", "ReleaseEpoch"}, TupleType:(*reflect.rtype)(0xc00047c3c0)}}
type Map3Node_ struct {
	NodeAddress     Address        `json:"NodeAddress" storage:"slot=0,offset=0"`
	OperatorAddress Address        `json:"OperatorAddress" storage:"slot=1,offset=0"`
	PubKeys         BLSPublicKeys_ `json:"PubKeys" storage:"slot=2,offset=0"`
	Commission      Decimal        `json:"Commission" storage:"slot=3,offset=0"`
	Description     Description_   `json:"Description" storage:"slot=4,offset=0"`
	CreationHeight  BigInt         `json:"CreationHeight" storage:"slot=9,offset=0"`
	Age             Decimal        `json:"Age" storage:"slot=10,offset=0"`
	Status          Uint8          `json:"Status" storage:"slot=11,offset=0"`
	ActivationEpoch BigInt         `json:"ActivationEpoch" storage:"slot=12,offset=0"`
	ReleaseEpoch    Decimal        `json:"ReleaseEpoch" storage:"slot=13,offset=0"`
}

// MicrodelegationMapEntry_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"MicrodelegationMapEntry_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000ba8f0), (*wrapper.tmplField)(0xc0000ba9c0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x140, stringKind:"t_struct(MicrodelegationMapEntry_)78_storage", TupleRawName:"MicrodelegationMapEntry_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000480680), (*abi.Type)(0xc000480f80)}, TupleRawNames:[]string{"Entry", "Index"}, TupleType:(*reflect.rtype)(0xc00047c840)}}
type MicrodelegationMapEntry_ struct {
	Entry Microdelegation_ `json:"Entry" storage:"slot=0,offset=0"`
	Index BigInt           `json:"Index" storage:"slot=9,offset=0"`
}

// MicrodelegationMap_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"MicrodelegationMap_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0002e5c70), (*wrapper.tmplField)(0xc0000bac30)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(MicrodelegationMap_)86_storage", TupleRawName:"MicrodelegationMap_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000480480), (*abi.Type)(0xc000480500)}, TupleRawNames:[]string{"Keys", "Map"}, TupleType:(*reflect.rtype)(0xc00047c8a0)}}
type MicrodelegationMap_ struct {
	Keys Slice_760032d9d704 `json:"Keys" storage:"slot=0,offset=0"`
	Map  Map_67e3f4d129ef   `json:"Map" storage:"slot=1,offset=0"`
}

// Microdelegation_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Microdelegation_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0002e5e10), (*wrapper.tmplField)(0xc0002e5ee0), (*wrapper.tmplField)(0xc0000ba000), (*wrapper.tmplField)(0xc0000ba0d0), (*wrapper.tmplField)(0xc0000ba340), (*wrapper.tmplField)(0xc0000ba5b0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x120, stringKind:"t_struct(Microdelegation_)68_storage", TupleRawName:"Microdelegation_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000480700), (*abi.Type)(0xc000480780), (*abi.Type)(0xc000480800), (*abi.Type)(0xc000480880), (*abi.Type)(0xc000480a80), (*abi.Type)(0xc000480d00)}, TupleRawNames:[]string{"DelegatorAddress", "Amount", "Reward", "Undelegation", "PendingDelegation", "Renewal"}, TupleType:(*reflect.rtype)(0xc00047c7e0)}}
type Microdelegation_ struct {
	DelegatorAddress  Address            `json:"DelegatorAddress" storage:"slot=0,offset=0"`
	Amount            BigInt             `json:"Amount" storage:"slot=1,offset=0"`
	Reward            BigInt             `json:"Reward" storage:"slot=2,offset=0"`
	Undelegation      Undelegation_      `json:"Undelegation" storage:"slot=3,offset=0"`
	PendingDelegation PendingDelegation_ `json:"PendingDelegation" storage:"slot=5,offset=0"`
	Renewal           Renewal_           `json:"Renewal" storage:"slot=7,offset=0"`
}

// PendingDelegation_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"PendingDelegation_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000ba1a0), (*wrapper.tmplField)(0xc0000ba270)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(PendingDelegation_)55_storage", TupleRawName:"PendingDelegation_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000480b00), (*abi.Type)(0xc000480b80)}, TupleRawNames:[]string{"Amount", "UnlockedEpoch"}, TupleType:(*reflect.rtype)(0xc00047c6c0)}}
type PendingDelegation_ struct {
	Amount        BigInt  `json:"Amount" storage:"slot=0,offset=0"`
	UnlockedEpoch Decimal `json:"UnlockedEpoch" storage:"slot=1,offset=0"`
}

// RedelegationReference_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"RedelegationReference_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000bb380)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(RedelegationReference_)89_storage", TupleRawName:"RedelegationReference_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000481280)}, TupleRawNames:[]string{"ValidatorAddress"}, TupleType:(*reflect.rtype)(0xc00047c960)}}
type RedelegationReference_ struct {
	ValidatorAddress Address `json:"ValidatorAddress" storage:"slot=0,offset=0"`
}

// Renewal_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Renewal_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000ba410), (*wrapper.tmplField)(0xc0000ba4e0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Renewal_)73_storage", TupleRawName:"Renewal_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000480d80), (*abi.Type)(0xc000480e00)}, TupleRawNames:[]string{"IsRenew", "UpdateHeight"}, TupleType:(*reflect.rtype)(0xc00047c780)}}
type Renewal_ struct {
	IsRenew      Bool   `json:"IsRenew" storage:"slot=0,offset=0"`
	UpdateHeight BigInt `json:"UpdateHeight" storage:"slot=1,offset=0"`
}

// Undelegation_ is an auto generated low-level Go binding around an user-defined struct.
// &wrapper.tmplStruct{T:0x6, Name:"Undelegation_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0002e5380), (*wrapper.tmplField)(0xc0002e5450)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Undelegation_)50_storage", TupleRawName:"Undelegation_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000480900), (*abi.Type)(0xc000480980)}, TupleRawNames:[]string{"Amount", "Epoch"}, TupleType:(*reflect.rtype)(0xc00047c600)}}
type Undelegation_ struct {
	Amount BigInt `json:"Amount" storage:"slot=0,offset=0"`
	Epoch  BigInt `json:"Epoch" storage:"slot=1,offset=0"`
}

type Storage_BLSPublicKey_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"BLSPublicKey_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0002e4680)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(BLSPublicKey_)20_storage", TupleRawName:"BLSPublicKey_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032cf80)}, TupleRawNames:[]string{"Key"}, TupleType:(*reflect.rtype)(0xc0001cfe00)}}
	obj           *BLSPublicKey_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_BLSPublicKeys_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"BLSPublicKeys_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0002e4820)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(BLSPublicKeys_)24_storage", TupleRawName:"BLSPublicKeys_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032d700)}, TupleRawNames:[]string{"Keys"}, TupleType:(*reflect.rtype)(0xc00047c180)}}
	obj           *BLSPublicKeys_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_DelegationIndexMapEntry_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"DelegationIndexMapEntry_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000122000), (*wrapper.tmplField)(0xc0001220d0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(DelegationIndexMapEntry_)125_storage", TupleRawName:"DelegationIndexMapEntry_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032ca00), (*abi.Type)(0xc00032cc80)}, TupleRawNames:[]string{"Entry", "Index"}, TupleType:(*reflect.rtype)(0xc0001cfc20)}}
	obj           *DelegationIndexMapEntry_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_DelegationIndexMap_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"DelegationIndexMap_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000122340), (*wrapper.tmplField)(0xc0001224e0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(DelegationIndexMap_)133_storage", TupleRawName:"DelegationIndexMap_", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fa100), (*abi.Type)(0xc0000fa180)}, TupleRawNames:[]string{"Keys", "Map"}, TupleType:(*reflect.rtype)(0xc0003223c0)}}
	obj           *DelegationIndexMap_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_DelegationIndex_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"DelegationIndex_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000bbe10), (*wrapper.tmplField)(0xc0000bbee0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(DelegationIndex_)120_storage", TupleRawName:"DelegationIndex_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032ca80), (*abi.Type)(0xc00032cb00)}, TupleRawNames:[]string{"NodeAddress", "IsOperator"}, TupleType:(*reflect.rtype)(0xc0001cfb60)}}
	obj           *DelegationIndex_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Description_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Description_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0002e4a90), (*wrapper.tmplField)(0xc0002e4b60), (*wrapper.tmplField)(0xc0002e4c30), (*wrapper.tmplField)(0xc0002e4d00), (*wrapper.tmplField)(0xc0002e4dd0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xa0, stringKind:"t_struct(Description_)15_storage", TupleRawName:"Description_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032db80), (*abi.Type)(0xc00032dc80), (*abi.Type)(0xc00032dd00), (*abi.Type)(0xc00032dd80), (*abi.Type)(0xc00032de00)}, TupleRawNames:[]string{"Name", "Identity", "Website", "SecurityContact", "Details"}, TupleType:(*reflect.rtype)(0xc00047c360)}}
	obj           *Description_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Global_t struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Global_t", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000122c30)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x0, stringKind:"t_struct(Global_t)", TupleRawName:"Global_t", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fa900)}, TupleRawNames:[]string{"Map3NodePool"}, TupleType:(*reflect.rtype)(0xc0003225a0)}}
	obj           *Global_t
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Map3NodePool_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Map3NodePool_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0001228f0), (*wrapper.tmplField)(0xc0001229c0), (*wrapper.tmplField)(0xc000122a90), (*wrapper.tmplField)(0xc000122b60)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xa0, stringKind:"t_struct(Map3NodePool_)148_storage", TupleRawName:"Map3NodePool_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032d200), (*abi.Type)(0xc00052a080), (*abi.Type)(0xc0000fa380), (*abi.Type)(0xc0000fa800)}, TupleRawNames:[]string{"Nodes", "DelegationIndexMapByDelegator", "SlotKeySet", "DescriptionIdentitySet"}, TupleType:(*reflect.rtype)(0xc0003224e0)}}
	obj           *Map3NodePool_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Map3NodeWrapperMapEntry_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Map3NodeWrapperMapEntry_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000bb790), (*wrapper.tmplField)(0xc0000bb860)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x2a0, stringKind:"t_struct(Map3NodeWrapperMapEntry_)107_storage", TupleRawName:"Map3NodeWrapperMapEntry_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032d480), (*abi.Type)(0xc000481600)}, TupleRawNames:[]string{"Entry", "Index"}, TupleType:(*reflect.rtype)(0xc00047ca20)}}
	obj           *Map3NodeWrapperMapEntry_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Map3NodeWrapperMap_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Map3NodeWrapperMap_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc000122750), (*wrapper.tmplField)(0xc000122820)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Map3NodeWrapperMap_)115_storage", TupleRawName:"Map3NodeWrapperMap_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032d280), (*abi.Type)(0xc00032d300)}, TupleRawNames:[]string{"Keys", "Map"}, TupleType:(*reflect.rtype)(0xc0001080c0)}}
	obj           *Map3NodeWrapperMap_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Map3NodeWrapper_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Map3NodeWrapper_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000bb1e0), (*wrapper.tmplField)(0xc0000bb2b0), (*wrapper.tmplField)(0xc0000bb450), (*wrapper.tmplField)(0xc0000bb520), (*wrapper.tmplField)(0xc0000bb5f0), (*wrapper.tmplField)(0xc0000bb6c0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x280, stringKind:"t_struct(Map3NodeWrapper_)102_storage", TupleRawName:"Map3NodeWrapper_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032d500), (*abi.Type)(0xc000480400), (*abi.Type)(0xc000481200), (*abi.Type)(0xc000481380), (*abi.Type)(0xc000481400), (*abi.Type)(0xc000481500)}, TupleRawNames:[]string{"Map3Node", "Microdelegations", "RedelegationReference", "AccumulatedReward", "TotalDelegation", "TotalPendingDelegation"}, TupleType:(*reflect.rtype)(0xc00047c9c0)}}
	obj           *Map3NodeWrapper_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Map3Node_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Map3Node_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0002e4410), (*wrapper.tmplField)(0xc0002e44e0), (*wrapper.tmplField)(0xc0002e48f0), (*wrapper.tmplField)(0xc0002e49c0), (*wrapper.tmplField)(0xc0002e4ea0), (*wrapper.tmplField)(0xc0002e4f70), (*wrapper.tmplField)(0xc0002e5040), (*wrapper.tmplField)(0xc0002e5110), (*wrapper.tmplField)(0xc0002e51e0), (*wrapper.tmplField)(0xc0002e52b0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x1c0, stringKind:"t_struct(Map3Node_)45_storage", TupleRawName:"Map3Node_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032d580), (*abi.Type)(0xc00032d600), (*abi.Type)(0xc00032d680), (*abi.Type)(0xc00032d980), (*abi.Type)(0xc00032db00), (*abi.Type)(0xc000480080), (*abi.Type)(0xc000480100), (*abi.Type)(0xc000480180), (*abi.Type)(0xc000480200), (*abi.Type)(0xc000480300)}, TupleRawNames:[]string{"NodeAddress", "OperatorAddress", "PubKeys", "Commission", "Description", "CreationHeight", "Age", "Status", "ActivationEpoch", "ReleaseEpoch"}, TupleType:(*reflect.rtype)(0xc00047c3c0)}}
	obj           *Map3Node_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_MicrodelegationMapEntry_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"MicrodelegationMapEntry_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000ba8f0), (*wrapper.tmplField)(0xc0000ba9c0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x140, stringKind:"t_struct(MicrodelegationMapEntry_)78_storage", TupleRawName:"MicrodelegationMapEntry_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000480680), (*abi.Type)(0xc000480f80)}, TupleRawNames:[]string{"Entry", "Index"}, TupleType:(*reflect.rtype)(0xc00047c840)}}
	obj           *MicrodelegationMapEntry_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_MicrodelegationMap_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"MicrodelegationMap_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0002e5c70), (*wrapper.tmplField)(0xc0000bac30)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(MicrodelegationMap_)86_storage", TupleRawName:"MicrodelegationMap_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000480480), (*abi.Type)(0xc000480500)}, TupleRawNames:[]string{"Keys", "Map"}, TupleType:(*reflect.rtype)(0xc00047c8a0)}}
	obj           *MicrodelegationMap_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Microdelegation_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Microdelegation_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0002e5e10), (*wrapper.tmplField)(0xc0002e5ee0), (*wrapper.tmplField)(0xc0000ba000), (*wrapper.tmplField)(0xc0000ba0d0), (*wrapper.tmplField)(0xc0000ba340), (*wrapper.tmplField)(0xc0000ba5b0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x120, stringKind:"t_struct(Microdelegation_)68_storage", TupleRawName:"Microdelegation_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000480700), (*abi.Type)(0xc000480780), (*abi.Type)(0xc000480800), (*abi.Type)(0xc000480880), (*abi.Type)(0xc000480a80), (*abi.Type)(0xc000480d00)}, TupleRawNames:[]string{"DelegatorAddress", "Amount", "Reward", "Undelegation", "PendingDelegation", "Renewal"}, TupleType:(*reflect.rtype)(0xc00047c7e0)}}
	obj           *Microdelegation_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_PendingDelegation_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"PendingDelegation_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000ba1a0), (*wrapper.tmplField)(0xc0000ba270)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(PendingDelegation_)55_storage", TupleRawName:"PendingDelegation_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000480b00), (*abi.Type)(0xc000480b80)}, TupleRawNames:[]string{"Amount", "UnlockedEpoch"}, TupleType:(*reflect.rtype)(0xc00047c6c0)}}
	obj           *PendingDelegation_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_RedelegationReference_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"RedelegationReference_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000bb380)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(RedelegationReference_)89_storage", TupleRawName:"RedelegationReference_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000481280)}, TupleRawNames:[]string{"ValidatorAddress"}, TupleType:(*reflect.rtype)(0xc00047c960)}}
	obj           *RedelegationReference_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Renewal_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Renewal_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0000ba410), (*wrapper.tmplField)(0xc0000ba4e0)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Renewal_)73_storage", TupleRawName:"Renewal_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000480d80), (*abi.Type)(0xc000480e00)}, TupleRawNames:[]string{"IsRenew", "UpdateHeight"}, TupleType:(*reflect.rtype)(0xc00047c780)}}
	obj           *Renewal_
	db            StateDB
	addr          common.Address
	slot          *big.Int
	offset        int
	numberOfBytes int
	dirty         StateValues
}

type Storage_Undelegation_ struct {

	// &wrapper.tmplStruct{T:0x6, Name:"Undelegation_", Fields:[]*wrapper.tmplField{(*wrapper.tmplField)(0xc0002e5380), (*wrapper.tmplField)(0xc0002e5450)}, Type:"", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Undelegation_)50_storage", TupleRawName:"Undelegation_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000480900), (*abi.Type)(0xc000480980)}, TupleRawNames:[]string{"Amount", "Epoch"}, TupleType:(*reflect.rtype)(0xc00047c600)}}
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
	// Field: &wrapper.tmplField{Type:"Array_aaa956d410fd", ElemType:"Array_aaa956d410fd", Name:"Key", SolKind:abi.Type{Elem:(*abi.Type)(0xc00032d000), Size:48, T:0x5, numberOfBytes:0x40, stringKind:"t_array(t_bytes1)48_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Key\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Slice_2f1765e04dad", ElemType:"Slice_2f1765e04dad", Name:"Keys", SolKind:abi.Type{Elem:(*abi.Type)(0xc00032d800), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_struct(BLSPublicKey_)20_storage)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Keys\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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

// Entry is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_DelegationIndexMapEntry_) Entry() *Storage_DelegationIndex_ {
	// Field: &wrapper.tmplField{Type:"DelegationIndex_", ElemType:"DelegationIndex_", Name:"Entry", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(DelegationIndex_)120_storage", TupleRawName:"DelegationIndex_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032ca80), (*abi.Type)(0xc00032cb00)}, TupleRawNames:[]string{"NodeAddress", "IsOperator"}, TupleType:(*reflect.rtype)(0xc0001cfb60)}, Tag:"json:\"Entry\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Slice_760032d9d704", ElemType:"Slice_760032d9d704", Name:"Keys", SolKind:abi.Type{Elem:(*abi.Type)(0xc00032c800), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_address)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Keys\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
func (s *Storage_DelegationIndexMap_) Map() *Storage_Map_64b0ddf77b70 {
	// Field: &wrapper.tmplField{Type:"Map_64b0ddf77b70", ElemType:"Map_64b0ddf77b70", Name:"Map", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(DelegationIndexMapEntry_)125_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032c900), (*abi.Type)(0xc00032ce00)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Map\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Map == nil {

		instance := make(Map_64b0ddf77b70)

		s.obj.Map = instance
	}

	return &Storage_Map_64b0ddf77b70{
		obj:           s.obj.Map,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// NodeAddress is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_DelegationIndex_) NodeAddress() *Storage_Address {
	// Field: &wrapper.tmplField{Type:"Address", ElemType:"Address", Name:"NodeAddress", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"NodeAddress\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

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
	// Field: &wrapper.tmplField{Type:"Map3NodePool_", ElemType:"Map3NodePool_", Name:"Map3NodePool", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xa0, stringKind:"t_struct(Map3NodePool_)148_storage", TupleRawName:"Map3NodePool_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032d200), (*abi.Type)(0xc00052a080), (*abi.Type)(0xc0000fa380), (*abi.Type)(0xc0000fa800)}, TupleRawNames:[]string{"Nodes", "DelegationIndexMapByDelegator", "SlotKeySet", "DescriptionIdentitySet"}, TupleType:(*reflect.rtype)(0xc0003224e0)}, Tag:"json:\"Map3NodePool\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Map3NodeWrapperMap_", ElemType:"Map3NodeWrapperMap_", Name:"Nodes", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Map3NodeWrapperMap_)115_storage", TupleRawName:"Map3NodeWrapperMap_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032d280), (*abi.Type)(0xc00032d300)}, TupleRawNames:[]string{"Keys", "Map"}, TupleType:(*reflect.rtype)(0xc0001080c0)}, Tag:"json:\"Nodes\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
func (s *Storage_Map3NodePool_) DelegationIndexMapByDelegator() *Storage_Map_bc58594e503f {
	// Field: &wrapper.tmplField{Type:"Map_bc58594e503f", ElemType:"Map_bc58594e503f", Name:"DelegationIndexMapByDelegator", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(DelegationIndexMap_)133_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fa000), (*abi.Type)(0xc0000fa280)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"DelegationIndexMapByDelegator\" storage:\"slot=2,offset=0\"", Slot:2, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("2", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.DelegationIndexMapByDelegator == nil {

		instance := make(Map_bc58594e503f)

		s.obj.DelegationIndexMapByDelegator = instance
	}

	return &Storage_Map_bc58594e503f{
		obj:           s.obj.DelegationIndexMapByDelegator,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// SlotKeySet is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodePool_) SlotKeySet() *Storage_Map_624bcc0df495 {
	// Field: &wrapper.tmplField{Type:"Map_624bcc0df495", ElemType:"Map_624bcc0df495", Name:"SlotKeySet", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_string_memory_ptr,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fa600), (*abi.Type)(0xc0000fa700)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"SlotKeySet\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("3", 10)
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
func (s *Storage_Map3NodePool_) DescriptionIdentitySet() *Storage_Map_624bcc0df495 {
	// Field: &wrapper.tmplField{Type:"Map_624bcc0df495", ElemType:"Map_624bcc0df495", Name:"DescriptionIdentitySet", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_string_memory_ptr,t_bool)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc0000fa600), (*abi.Type)(0xc0000fa700)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"DescriptionIdentitySet\" storage:\"slot=4,offset=0\"", Slot:4, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Map3NodeWrapper_", ElemType:"Map3NodeWrapper_", Name:"Entry", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x280, stringKind:"t_struct(Map3NodeWrapper_)102_storage", TupleRawName:"Map3NodeWrapper_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032d500), (*abi.Type)(0xc000480400), (*abi.Type)(0xc000481200), (*abi.Type)(0xc000481380), (*abi.Type)(0xc000481400), (*abi.Type)(0xc000481500)}, TupleRawNames:[]string{"Map3Node", "Microdelegations", "RedelegationReference", "AccumulatedReward", "TotalDelegation", "TotalPendingDelegation"}, TupleType:(*reflect.rtype)(0xc00047c9c0)}, Tag:"json:\"Entry\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
		numberOfBytes: 640,
		dirty:         s.dirty,
	}
}

// Index is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeWrapperMapEntry_) Index() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"Index", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Index\" storage:\"slot=20,offset=0\"", Slot:20, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("20", 10)
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
	// Field: &wrapper.tmplField{Type:"Slice_760032d9d704", ElemType:"Slice_760032d9d704", Name:"Keys", SolKind:abi.Type{Elem:(*abi.Type)(0xc00032c800), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_address)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Keys\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
func (s *Storage_Map3NodeWrapperMap_) Map() *Storage_Map_4dcba4f803c2 {
	// Field: &wrapper.tmplField{Type:"Map_4dcba4f803c2", ElemType:"Map_4dcba4f803c2", Name:"Map", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(Map3NodeWrapperMapEntry_)107_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032d380), (*abi.Type)(0xc000481700)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Map\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Map == nil {

		instance := make(Map_4dcba4f803c2)

		s.obj.Map = instance
	}

	return &Storage_Map_4dcba4f803c2{
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
	// Field: &wrapper.tmplField{Type:"Map3Node_", ElemType:"Map3Node_", Name:"Map3Node", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x1c0, stringKind:"t_struct(Map3Node_)45_storage", TupleRawName:"Map3Node_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032d580), (*abi.Type)(0xc00032d600), (*abi.Type)(0xc00032d680), (*abi.Type)(0xc00032d980), (*abi.Type)(0xc00032db00), (*abi.Type)(0xc000480080), (*abi.Type)(0xc000480100), (*abi.Type)(0xc000480180), (*abi.Type)(0xc000480200), (*abi.Type)(0xc000480300)}, TupleRawNames:[]string{"NodeAddress", "OperatorAddress", "PubKeys", "Commission", "Description", "CreationHeight", "Age", "Status", "ActivationEpoch", "ReleaseEpoch"}, TupleType:(*reflect.rtype)(0xc00047c3c0)}, Tag:"json:\"Map3Node\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
		numberOfBytes: 448,
		dirty:         s.dirty,
	}
}

// Microdelegations is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeWrapper_) Microdelegations() *Storage_MicrodelegationMap_ {
	// Field: &wrapper.tmplField{Type:"MicrodelegationMap_", ElemType:"MicrodelegationMap_", Name:"Microdelegations", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(MicrodelegationMap_)86_storage", TupleRawName:"MicrodelegationMap_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000480480), (*abi.Type)(0xc000480500)}, TupleRawNames:[]string{"Keys", "Map"}, TupleType:(*reflect.rtype)(0xc00047c8a0)}, Tag:"json:\"Microdelegations\" storage:\"slot=14,offset=0\"", Slot:14, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("14", 10)
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

// RedelegationReference is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3NodeWrapper_) RedelegationReference() *Storage_RedelegationReference_ {
	// Field: &wrapper.tmplField{Type:"RedelegationReference_", ElemType:"RedelegationReference_", Name:"RedelegationReference", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(RedelegationReference_)89_storage", TupleRawName:"RedelegationReference_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000481280)}, TupleRawNames:[]string{"ValidatorAddress"}, TupleType:(*reflect.rtype)(0xc00047c960)}, Tag:"json:\"RedelegationReference\" storage:\"slot=16,offset=0\"", Slot:16, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("16", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_RedelegationReference_{
		obj:           &s.obj.RedelegationReference,
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
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"AccumulatedReward", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"AccumulatedReward\" storage:\"slot=17,offset=0\"", Slot:17, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("17", 10)
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
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"TotalDelegation", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"TotalDelegation\" storage:\"slot=18,offset=0\"", Slot:18, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("18", 10)
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
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"TotalPendingDelegation", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"TotalPendingDelegation\" storage:\"slot=19,offset=0\"", Slot:19, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("19", 10)
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

// NodeAddress is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3Node_) NodeAddress() *Storage_Address {
	// Field: &wrapper.tmplField{Type:"Address", ElemType:"Address", Name:"NodeAddress", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:20, T:0x7, numberOfBytes:0x14, stringKind:"t_address", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"NodeAddress\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

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

// PubKeys is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3Node_) PubKeys() *Storage_BLSPublicKeys_ {
	// Field: &wrapper.tmplField{Type:"BLSPublicKeys_", ElemType:"BLSPublicKeys_", Name:"PubKeys", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(BLSPublicKeys_)24_storage", TupleRawName:"BLSPublicKeys_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032d700)}, TupleRawNames:[]string{"Keys"}, TupleType:(*reflect.rtype)(0xc00047c180)}, Tag:"json:\"PubKeys\" storage:\"slot=2,offset=0\"", Slot:2, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("2", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_BLSPublicKeys_{
		obj:           &s.obj.PubKeys,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Commission is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3Node_) Commission() *Storage_Decimal {
	// Field: &wrapper.tmplField{Type:"Decimal", ElemType:"Decimal", Name:"Commission", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032da00)}, TupleRawNames:[]string{"f"}, TupleType:(*reflect.rtype)(0xc00047c240)}, Tag:"json:\"Commission\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("3", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Decimal{
		obj:           &s.obj.Commission,
		db:            s.db,
		addr:          s.addr,
		slot:          actual,
		offset:        offset,
		numberOfBytes: 32,
		dirty:         s.dirty,
	}
}

// Description is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3Node_) Description() *Storage_Description_ {
	// Field: &wrapper.tmplField{Type:"Description_", ElemType:"Description_", Name:"Description", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0xa0, stringKind:"t_struct(Description_)15_storage", TupleRawName:"Description_", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032db80), (*abi.Type)(0xc00032dc80), (*abi.Type)(0xc00032dd00), (*abi.Type)(0xc00032dd80), (*abi.Type)(0xc00032de00)}, TupleRawNames:[]string{"Name", "Identity", "Website", "SecurityContact", "Details"}, TupleType:(*reflect.rtype)(0xc00047c360)}, Tag:"json:\"Description\" storage:\"slot=4,offset=0\"", Slot:4, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("4", 10)
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
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"CreationHeight", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"CreationHeight\" storage:\"slot=9,offset=0\"", Slot:9, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("9", 10)
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
	// Field: &wrapper.tmplField{Type:"Decimal", ElemType:"Decimal", Name:"Age", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032da00)}, TupleRawNames:[]string{"f"}, TupleType:(*reflect.rtype)(0xc00047c240)}, Tag:"json:\"Age\" storage:\"slot=10,offset=0\"", Slot:10, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("10", 10)
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
	// Field: &wrapper.tmplField{Type:"Uint8", ElemType:"Uint8", Name:"Status", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:8, T:0x1, numberOfBytes:0x1, stringKind:"uint8", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Status\" storage:\"slot=11,offset=0\"", Slot:11, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("11", 10)
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

// ActivationEpoch is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Map3Node_) ActivationEpoch() *Storage_BigInt {
	// Field: &wrapper.tmplField{Type:"BigInt", ElemType:"BigInt", Name:"ActivationEpoch", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:256, T:0x1, numberOfBytes:0x20, stringKind:"t_uint256", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"ActivationEpoch\" storage:\"slot=12,offset=0\"", Slot:12, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("12", 10)
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
	// Field: &wrapper.tmplField{Type:"Decimal", ElemType:"Decimal", Name:"ReleaseEpoch", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032da00)}, TupleRawNames:[]string{"f"}, TupleType:(*reflect.rtype)(0xc00047c240)}, Tag:"json:\"ReleaseEpoch\" storage:\"slot=13,offset=0\"", Slot:13, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("13", 10)
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
	// Field: &wrapper.tmplField{Type:"Microdelegation_", ElemType:"Microdelegation_", Name:"Entry", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x120, stringKind:"t_struct(Microdelegation_)68_storage", TupleRawName:"Microdelegation_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000480700), (*abi.Type)(0xc000480780), (*abi.Type)(0xc000480800), (*abi.Type)(0xc000480880), (*abi.Type)(0xc000480a80), (*abi.Type)(0xc000480d00)}, TupleRawNames:[]string{"DelegatorAddress", "Amount", "Reward", "Undelegation", "PendingDelegation", "Renewal"}, TupleType:(*reflect.rtype)(0xc00047c7e0)}, Tag:"json:\"Entry\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Slice_760032d9d704", ElemType:"Slice_760032d9d704", Name:"Keys", SolKind:abi.Type{Elem:(*abi.Type)(0xc00032c800), Size:0, T:0x4, numberOfBytes:0x20, stringKind:"t_array(t_address)dyn_storage", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Keys\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
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
func (s *Storage_MicrodelegationMap_) Map() *Storage_Map_67e3f4d129ef {
	// Field: &wrapper.tmplField{Type:"Map_67e3f4d129ef", ElemType:"Map_67e3f4d129ef", Name:"Map", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0xe, numberOfBytes:0x20, stringKind:"t_mapping(t_address,t_struct(MicrodelegationMapEntry_)78_storage)", TupleRawName:"", TupleElems:[]*abi.Type{(*abi.Type)(0xc000480580), (*abi.Type)(0xc000481080)}, TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"Map\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("1", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	if s.obj.Map == nil {

		instance := make(Map_67e3f4d129ef)

		s.obj.Map = instance
	}

	return &Storage_Map_67e3f4d129ef{
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

// Undelegation is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Microdelegation_) Undelegation() *Storage_Undelegation_ {
	// Field: &wrapper.tmplField{Type:"Undelegation_", ElemType:"Undelegation_", Name:"Undelegation", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Undelegation_)50_storage", TupleRawName:"Undelegation_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000480900), (*abi.Type)(0xc000480980)}, TupleRawNames:[]string{"Amount", "Epoch"}, TupleType:(*reflect.rtype)(0xc00047c600)}, Tag:"json:\"Undelegation\" storage:\"slot=3,offset=0\"", Slot:3, Offset:0}
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

// PendingDelegation is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Microdelegation_) PendingDelegation() *Storage_PendingDelegation_ {
	// Field: &wrapper.tmplField{Type:"PendingDelegation_", ElemType:"PendingDelegation_", Name:"PendingDelegation", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(PendingDelegation_)55_storage", TupleRawName:"PendingDelegation_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000480b00), (*abi.Type)(0xc000480b80)}, TupleRawNames:[]string{"Amount", "UnlockedEpoch"}, TupleType:(*reflect.rtype)(0xc00047c6c0)}, Tag:"json:\"PendingDelegation\" storage:\"slot=5,offset=0\"", Slot:5, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("5", 10)
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

// Renewal is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Microdelegation_) Renewal() *Storage_Renewal_ {
	// Field: &wrapper.tmplField{Type:"Renewal_", ElemType:"Renewal_", Name:"Renewal", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x40, stringKind:"t_struct(Renewal_)73_storage", TupleRawName:"Renewal_", TupleElems:[]*abi.Type{(*abi.Type)(0xc000480d80), (*abi.Type)(0xc000480e00)}, TupleRawNames:[]string{"IsRenew", "UpdateHeight"}, TupleType:(*reflect.rtype)(0xc00047c780)}, Tag:"json:\"Renewal\" storage:\"slot=7,offset=0\"", Slot:7, Offset:0}
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
	// Field: &wrapper.tmplField{Type:"Decimal", ElemType:"Decimal", Name:"UnlockedEpoch", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x6, numberOfBytes:0x20, stringKind:"t_struct(Decimal)4_storage", TupleRawName:"Decimal", TupleElems:[]*abi.Type{(*abi.Type)(0xc00032da00)}, TupleRawNames:[]string{"f"}, TupleType:(*reflect.rtype)(0xc00047c240)}, Tag:"json:\"UnlockedEpoch\" storage:\"slot=1,offset=0\"", Slot:1, Offset:0}
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

// ValidatorAddress is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_RedelegationReference_) ValidatorAddress() *Storage_Address {
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

// IsRenew is an auto generated low-level Go binding around an user-defined struct.
func (s *Storage_Renewal_) IsRenew() *Storage_Bool {
	// Field: &wrapper.tmplField{Type:"Bool", ElemType:"Bool", Name:"IsRenew", SolKind:abi.Type{Elem:(*abi.Type)(nil), Size:0, T:0x2, numberOfBytes:0x1, stringKind:"t_bool", TupleRawName:"", TupleElems:[]*abi.Type(nil), TupleRawNames:[]string(nil), TupleType:reflect.Type(nil)}, Tag:"json:\"IsRenew\" storage:\"slot=0,offset=0\"", Slot:0, Offset:0}
	var slot *big.Int
	slot, _ = big.NewInt(0).SetString("0", 10)
	offset := 0

	var actual *big.Int = big.NewInt(0).Add(s.slot, slot)
	return &Storage_Bool{
		obj:           &s.obj.IsRenew,
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

