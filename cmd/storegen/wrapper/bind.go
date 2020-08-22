package wrapper

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"strconv"

	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// bindType is a set of type binders that convert Solidity types to some supported
// programming language types.
var bindType = map[bind.Lang]func(kind abi.Type, structs map[string]*tmplStruct) string{
	bind.LangGo: bindTypeGo,
}

func nameOfBasicTypeGo(kind *abi.Type) string {
	switch kind.T {
	case abi.AddressTy:
		return "Address"
	case abi.IntTy, abi.UintTy:
		parts := regexp.MustCompile(`(u)?int([0-9]*)`).FindStringSubmatch(kind.String())
		switch parts[2] {
		case "8", "16", "32", "64":
			name := fmt.Sprintf("%sint%s", parts[1], parts[2])
			return capitalise(name)
		}
		return "BigInt"
	case abi.FixedBytesTy:
		return fmt.Sprintf("Bytes%d", kind.Size)
	case abi.BytesTy:
		return "Bytes"
	case abi.FunctionTy:
		return "Function"
	case abi.StringTy:
		return "String"
	case abi.BoolTy:
		return "Bool"
	case abi.TupleTy:
		panic("basic type do not support tuple/struct")
	case abi.PointerTy:
		panic("basic type do not support pointer")
	case abi.MappingTy:
		panic("basic type do not support mapping")
	case abi.DecimalTy:
		return "Decimal"
	default:
		// string, bool types
		return capitalise(kind.String())
	}
}

func bindBasicTypeGo(kind abi.Type) (string, string) {
	typ := GetReflectType(&kind)
	return nameOfBasicTypeGo(&kind), typ.String()
}

// bindTypeGo converts solidity types to Go ones. Since there is no clear mapping
// from all Solidity types to Go ones (e.g. uint17), those that cannot be exactly
// mapped will use an upscaled type (e.g. BigDecimal).
func bindTypeGo(kind abi.Type, structs map[string]*tmplStruct) string {
	switch kind.T {
	case abi.TupleTy:
		return structs[kind.TupleRawName+kind.String()].Name
	case abi.ArrayTy:
		modifier := map[bool]string{true: "*", false: ""}[GetReflectType(kind.Elem).Kind() == reflect.Ptr && GetReflectType(kind.Elem) != bigT]
		return fmt.Sprintf("[%d]", kind.Size) + modifier + bindTypeGo(*kind.Elem, structs)
	case abi.SliceTy:
		modifier := map[bool]string{true: "*", false: ""}[GetReflectType(kind.Elem).Kind() == reflect.Ptr && GetReflectType(kind.Elem) != bigT]
		return "[]" + modifier + bindTypeGo(*kind.Elem, structs)
	default:
		name, _ := bindBasicTypeGo(kind)
		return name
	}
}

// bindStructType is a set of type binders that convert Solidity tuple types to some supported
// programming language struct definition.
var bindStructType = map[bind.Lang]func(kind abi.Type, structs map[string]*tmplStruct) (string, error){
	bind.LangGo: bindStructTypeGo,
}

// bindStructTypeGo converts a Solidity tuple type to a Go one and records the mapping
// in the given map.
// Notably, this function will resolve and record nested struct recursively.
func bindStructTypeGo(kind abi.Type, structs map[string]*tmplStruct) (string, error) {
	switch kind.T {
	case abi.TupleTy:
		// We compose raw struct name and canonical parameter expression
		// together here. The reason is before solidity v0.5.11, kind.TupleRawName
		// is empty, so we use canonical parameter expression to distinguish
		// different struct definition. From the consideration of backward
		// compatibility, we concat these two together so that if kind.TupleRawName
		// is not empty, it can have unique id.
		id := kind.TupleRawName
		if isBuiltinType(id) {
			return "", errors.New(fmt.Sprintf("%s is builtin type.", id))
		}

		if s, exist := structs[id]; exist {
			return s.Name, nil
		}

		if kind.TupleRawName == "Decimal" {
			kind.T = abi.DecimalTy
			name, typ := bindBasicTypeGo(kind)

			structs[name] = &tmplStruct{
				Name:    name,
				T:       kind.T,
				Type:    typ,
				SolKind: kind,
			}
			return name, nil
		}

		var fields []*tmplField
		var underlyingType reflect.Type
		underlyingType = GetReflectType(&kind)
		for i, elem := range kind.TupleElems {
			field, err := bindStructTypeGo(*elem, structs)
			if err != nil {
				return "", err
			}
			fieldType, _ := underlyingType.FieldByName(abi.ToCamelCase(kind.TupleRawNames[i]))
			slot, offset, _ := parseTag(string(fieldType.Tag.Get("storage")))
			fields = append(fields, &tmplField{
				Type:     field,
				Name:     capitalise(kind.TupleRawNames[i]),
				Tag:      string(fieldType.Tag),
				SolKind:  *elem,
				ElemType: field,
				Slot:     slot,
				Offset:   offset,
			})
		}
		name := kind.TupleRawName
		if name == "" {
			name = fmt.Sprintf("Struct%d", len(structs))
		}
		structs[id] = &tmplStruct{
			Name:    name,
			Fields:  fields,
			T:       kind.T,
			SolKind: kind,
		}
		return name, nil

	case abi.ArrayTy:
		name := fmt.Sprintf("Array_%s", typeHash(kind.String()))
		if isBuiltinType(name) {
			return "", errors.New(fmt.Sprintf("%s is builtin type.", name))
		}

		if s, exist := structs[name]; exist {
			return s.Name, nil
		}

		field, err := bindStructTypeGo(*kind.Elem, structs)
		if err != nil {
			return "", err
		}

		typ := fmt.Sprintf("[%d]%s", kind.Size, field)
		var fields []*tmplField
		fields = append(fields, &tmplField{
			Type:    field,
			Name:    "",
			SolKind: *kind.Elem,
		})

		structs[name] = &tmplStruct{
			Name:    name,
			T:       kind.T,
			Type:    typ,
			SolKind: kind,
			Fields:  fields,
		}
		return name, nil

	case abi.SliceTy:
		name := fmt.Sprintf("Slice_%s", typeHash(kind.String()))
		if isBuiltinType(name) {
			return "", errors.New(fmt.Sprintf("%s is builtin type.", name))
		}

		if s, exist := structs[name]; exist {
			return s.Name, nil
		}

		field, err := bindStructTypeGo(*kind.Elem, structs)
		if err != nil {
			return "", err
		}

		modifier := map[bool]string{true: "*", false: ""}[isT(kind.Elem, abi.PointerTy)]
		typ := fmt.Sprintf("[]%s%s", modifier, field)
		var fields []*tmplField
		fields = append(fields, &tmplField{
			Type:    field,
			Name:    "",
			SolKind: *kind.Elem,
		})

		structs[name] = &tmplStruct{
			Name:    name,
			T:       kind.T,
			Type:    typ,
			SolKind: kind,
			Fields:  fields,
		}
		return name, nil

	case abi.MappingTy:
		name := fmt.Sprintf("Map_%s", typeHash(kind.String()))
		if isBuiltinType(name) {
			return "", errors.New(fmt.Sprintf("%s is builtin type.", name))
		}

		if s, exist := structs[name]; exist {
			return s.Name, nil
		}

		var fields []*tmplField
		for i := 0; i < 2; i++ {
			field, err := bindStructTypeGo(*kind.TupleElems[i], structs)
			if err != nil {
				return "", err
			}
			fields = append(fields, &tmplField{
				Type:    field,
				Name:    "",
				SolKind: *kind.TupleElems[i],
			})
		}

		modifier := map[bool]string{true: "*", false: ""}[isT(kind.TupleElems[1], abi.PointerTy)]
		typ := fmt.Sprintf("map[%s]%s%s", fields[0].Type,
			modifier,
			fields[1].Type)
		structs[name] = &tmplStruct{
			Name:    name,
			T:       kind.T,
			Type:    typ,
			SolKind: kind,
			Fields:  fields,
		}
		return name, nil

	case abi.PointerTy:
		name, _ := bindStructTypeGo(*kind.Elem, structs)
		return name, nil

	default:
		name, typ := bindBasicTypeGo(kind)
		if s, exist := structs[name]; exist {
			return s.Name, nil
		}

		structs[name] = &tmplStruct{
			Name:    name,
			T:       kind.T,
			Type:    typ,
			SolKind: kind,
		}
		return name, nil
	}
}

func isBytes(name string) bool {
	re := regexp.MustCompile("(Bytes|Bytes([1-9]|[12][0-9]|3[0-2]))")
	return re.MatchString(name)
}

func isBuiltinType(name string) bool {
	switch name {
	case
		"Address",
		"Bool",
		"Int8", "Int16", "Int32", "Int64",
		"Uint8", "Uint16", "Uint32", "Uint64",
		"BigInt",
		"String",
		"Bytes":
		return true
	default:
		if isBytes(name) {
			return true
		}
		return false
	}
}

func parseTag(tag string) (slot *big.Int, offset int, err error) {
	re := regexp.MustCompile(`slot=([^,]+),offset=([0-9]+)`)
	matches := re.FindAllStringSubmatch(tag, -1)
	if len(matches) == 0 {
		return big.NewInt(0), 0, errors.New(fmt.Sprintf("invalid tag: %s", tag))
	}

	parsed := matches[0]
	slot = big.NewInt(0)
	slot, ok := slot.SetString(parsed[1], 10)
	if !ok {
		return big.NewInt(0), 0, errors.New(fmt.Sprintf("invalid tag: %s", tag))
	}

	offset, err = strconv.Atoi(parsed[2])
	if err != nil {
		return big.NewInt(0), 0, errors.New(fmt.Sprintf("invalid tag: %s", tag))
	}

	return slot, offset, nil
}

func typeHash(value string) string {
	digest := sha3.Sum256([]byte(value))
	return hexutil.Encode(digest[:6])[2:]
}
