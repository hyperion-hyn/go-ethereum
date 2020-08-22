package wrapper

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var (
	// typeRegex parses the abi sub types
	typeRegex  = regexp.MustCompile(`t_([_a-z]+)([0-9]+)?(\((.*)\)([_a-z0-9]*)?)?`)
	arrayRegex = regexp.MustCompile(`^.*\[([0-9]*)\].*$`)
)

var typeCache sync.Map // map[reflect.Type]structFields

func NewType(t string, definition Type, types map[string]Type) (typ abi.Type, err error) {
	if f, ok := typeCache.Load(t); ok {
		return f.(abi.Type), nil
	}

	typ, err = newType(t, definition, types)
	if err != nil {
		return typ, err
	}

	f, _ := typeCache.LoadOrStore(t, typ)
	return f.(abi.Type), err
}

func newType(t string, definition Type, types map[string]Type) (typ abi.Type, err error) {
	typ.SetStringKind(t)
	typ.SetNumberOfBytes(uint(definition.NumberOfBytes.Uint64()))

	// parse the type and size of the abi-type.
	matches := typeRegex.FindAllStringSubmatch(t, -1)
	if len(matches) == 0 {
		return abi.Type{}, fmt.Errorf("invalid type '%v'", t)
	}
	parsedType := matches[0]

	// varSize is the size of the variable
	var varSize int
	if len(parsedType[2]) > 0 {
		var err error
		varSize, err = strconv.Atoi(parsedType[2])
		if err != nil {
			return abi.Type{}, fmt.Errorf("abi: error parsing variable size: %v", err)
		}
	} else {
		if parsedType[0] == "uint" || parsedType[0] == "int" {
			// this should fail because it means that there's something wrong with
			// the abi type (the compiler should always format it to the size...always)
			return abi.Type{}, fmt.Errorf("unsupported arg type: %s", t)
		}
	}
	// varType is the parsed abi type
	switch varType := parsedType[1]; varType {
	case "int":
		typ.Size = varSize
		typ.T = abi.IntTy
	case "uint":
		typ.Size = varSize
		typ.T = abi.UintTy
	case "bool":
		typ.T = abi.BoolTy
	case "address":
		typ.Size = 20
		typ.T = abi.AddressTy
	case "string", "string_storage", "string_memory_ptr":
		typ.T = abi.StringTy
	case "bytes", "bytes_storage":
		if varSize == 0 && definition.Encoding != "bytes" {
			return abi.Type{}, fmt.Errorf("abi: invalid encoding for %v, expected %v but got %v", t, "bytes", definition.Encoding)
		}

		if varSize == 1 {
			typ.T = abi.UintTy
			typ.Size = 8
			typ.SetStringKind("uint8")
		} else if varSize == 0 {
			typ.T = abi.BytesTy
		} else {
			typ.T = abi.FixedBytesTy
			typ.Size = varSize
		}

	case "struct":
		var (
			fields     []reflect.StructField
			elems      []*abi.Type
			names      []string
			expression string // canonical parameter expression
		)
		expression += "("
		overloadedNames := make(map[string]string)

		for idx, c := range definition.Members {
			cType, err := NewType(c.Type, types[c.Type], types)
			if err != nil {
				return abi.Type{}, err
			}
			fieldName, err := overloadedArgName(c.Label, overloadedNames)
			if err != nil {
				return abi.Type{}, err
			}
			overloadedNames[fieldName] = fieldName
			fields = append(fields, reflect.StructField{
				Name: fieldName, // reflect.StructOf will panic for any exported field.
				Type: GetReflectType(&cType),
				Tag:  reflect.StructTag("json:\"" + c.Label + "\"" + " " + fmt.Sprintf(`storage:"slot=%s,offset=%d"`, c.Slot.String(), c.Offset)),
			})
			elems = append(elems, &cType)
			names = append(names, c.Label)
			expression += cType.StringKind()
			if idx != len(definition.Members)-1 {
				expression += ","
			}
		}
		expression += ")"

		typ.TupleType = reflect.StructOf(fields)
		typ.TupleElems = elems
		typ.TupleRawNames = names
		typ.T = abi.TupleTy
		typ.TupleRawName = parsedType[4]

	case "function":
		typ.T = abi.FunctionTy
		typ.Size = 24

	case "array":
		// https://solidity.readthedocs.io/en/v0.6.11/types.html
		// The type of an array of fixed size k and element type T is written as T[k], and an array of dynamic size as T[].
		//
		// For example, an array of 5 dynamic arrays of uint is written as uint[][5]. The notation is reversed compared
		// to some other languages. In Solidity, X[3] is always an array containing three elements of type X, even if X
		// is itself an array. This is not the case in other languages such as C.
		cType, err := NewType(definition.Base, types[definition.Base], types)
		if err != nil {
			return abi.Type{}, err
		}

		if definition.Encoding == "dynamic_array" {
			varSize = 0
		} else {
			matches := arrayRegex.FindAllStringSubmatch(definition.Label, -1)
			if len(matches) == 0 {
				return abi.Type{}, fmt.Errorf("invalid type '%v'", t)
			}
			parsedType := matches[0]
			varSize = 0
			if len(parsedType[1]) > 0 {
				varSize, err = strconv.Atoi(parsedType[1])
				if err != nil {
					return abi.Type{}, fmt.Errorf("abi: error parsing variable size: %v", err)
				}
			}
		}

		typ.Elem = &cType
		if varSize == 0 {
			typ.T = abi.SliceTy

			switch GetReflectType(&cType).Kind() {
			case reflect.Ptr, reflect.Slice, reflect.Map:
			default:
				embeddedType := &abi.Type{}
				embeddedType.T = abi.PointerTy
				embeddedType.Elem = &cType
				typ.Elem = embeddedType
			}
		} else {
			typ.T = abi.ArrayTy
			typ.Size = varSize
		}

	case "mapping":
		if definition.Encoding != "mapping" {
			return abi.Type{}, fmt.Errorf("abi: invalid encoding for %v, expected %v but got %v", t, "mapping", definition.Encoding)
		}

		var (
			elems []*abi.Type
		)
		typ.T = abi.MappingTy
		kType, err := NewType(definition.Key, types[definition.Key], types)
		if err != nil {
			return abi.Type{}, err
		}

		vType, err := NewType(definition.Value, types[definition.Value], types)
		if err != nil {
			return abi.Type{}, err
		}

		var embeddedType *abi.Type = &vType
		switch GetReflectType(&vType).Kind() {
		case reflect.Ptr, reflect.Slice, reflect.Map:
		default:
			embeddedType = &abi.Type{}
			embeddedType.T = abi.PointerTy
			embeddedType.Elem = &vType
		}
		elems = append(elems, &kType, embeddedType)
		typ.TupleElems = elems

	default:
		return abi.Type{}, fmt.Errorf("unsupported arg type: %s", t)
	}

	return
}

// GetReflectType returns the reflection type of the ABI type.
func GetReflectType(t *abi.Type) reflect.Type {
	switch t.T {
	case abi.IntTy:
		return abi.ReflectIntType(false, t.Size)
	case abi.UintTy:
		return abi.ReflectIntType(true, t.Size)
	case abi.BoolTy:
		return reflect.TypeOf(false)
	case abi.StringTy:
		return reflect.TypeOf("")
	case abi.SliceTy:
		return reflect.SliceOf(GetReflectType(t.Elem))
	case abi.ArrayTy:
		return reflect.ArrayOf(t.Size, GetReflectType(t.Elem))
	case abi.TupleTy:
		return t.TupleType
	case abi.AddressTy:
		return reflect.TypeOf(common.Address{})
	case abi.FixedBytesTy:
		return reflect.ArrayOf(t.Size, reflect.TypeOf(byte(0)))
	case abi.BytesTy:
		return reflect.SliceOf(reflect.TypeOf(byte(0)))
	case abi.HashTy:
		// hashtype currently not used
		return reflect.ArrayOf(32, reflect.TypeOf(byte(0)))
	case abi.FixedPointTy:
		// fixedpoint type currently not used
		return reflect.ArrayOf(32, reflect.TypeOf(byte(0)))
	case abi.FunctionTy:
		return reflect.ArrayOf(24, reflect.TypeOf(byte(0)))
	case abi.PointerTy:
		return reflect.PtrTo(GetReflectType(t.Elem))
	case abi.MappingTy:
		return reflect.MapOf(GetReflectType(t.TupleElems[0]), GetReflectType(t.TupleElems[1]))
	case abi.DecimalTy:
		return reflect.TypeOf(common.Dec{})
	default:
		panic("Invalid type")
	}
}

func overloadedArgName(rawName string, names map[string]string) (string, error) {
	fieldName := abi.ToCamelCase(rawName)
	if fieldName == "" {
		return "", errors.New("abi: purely anonymous or underscored field is not supported")
	}
	// Handle overloaded fieldNames
	_, ok := names[fieldName]
	for idx := 0; ok; idx++ {
		fieldName = fmt.Sprintf("%s%d", abi.ToCamelCase(rawName), idx)
		_, ok = names[fieldName]
	}
	return fieldName, nil
}
