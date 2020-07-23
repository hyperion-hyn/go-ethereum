package wrapper

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
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
	typ.StringKind = t
	typ.NumberOfBytes = uint(definition.NumberOfBytes.Uint64())

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
		typ.Kind, typ.Type = reflectIntKindAndType(false, varSize)
		typ.Size = varSize
		typ.T = abi.IntTy
	case "uint":
		typ.Kind, typ.Type = reflectIntKindAndType(true, varSize)
		typ.Size = varSize
		typ.T = abi.UintTy
	case "bool":
		typ.Kind = reflect.Bool
		typ.T = abi.BoolTy
		typ.Type = reflect.TypeOf(bool(false))
	case "address":
		typ.Kind = reflect.Array
		typ.Type = addressT
		typ.Size = 20
		typ.T = abi.AddressTy
	case "string", "string_storage", "string_memory_ptr":
		typ.Kind = reflect.String
		typ.Type = stringT
		typ.T = abi.StringTy
	case "bytes", "bytes_storage":
		if varSize == 0 {
			if definition.Encoding != "bytes" {
				return abi.Type{}, fmt.Errorf("abi: invalid encoding for %v, expected %v but got %v", t, "bytes", definition.Encoding)
			}

			typ.T = abi.BytesTy
			typ.Kind = reflect.Slice
			typ.Type = reflect.SliceOf(reflect.TypeOf(byte(0)))
		} else if varSize == 1 {
			typ.T = abi.UintTy
			typ.Kind = reflect.Uint8
			typ.Type = uint8T
			typ.StringKind = "uint8"
		} else {
			typ.T = abi.FixedBytesTy
			typ.Kind = reflect.Array
			typ.Size = varSize
			typ.Type = reflect.ArrayOf(varSize, reflect.TypeOf(byte(0)))
		}

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
			typ.Kind = reflect.Slice
			switch cType.Type.Kind() {
			case reflect.Ptr, reflect.Slice, reflect.Map:
			default:
				cType.Type = reflect.PtrTo(cType.Type)
				cType.Kind = cType.Type.Kind()
			}
			typ.Type = reflect.SliceOf(reflect.TypeOf(cType.Type))
		} else {
			typ.T = abi.ArrayTy
			typ.Kind = reflect.Array
			typ.Size = varSize
			typ.Type = reflect.ArrayOf(varSize, reflect.TypeOf(cType.Type))
		}

	case "struct":
		var (
			fields []reflect.StructField
			elems  []*abi.Type
			names  []string
		)

		for _, c := range definition.Members {
			cType, err := NewType(c.Type, types[c.Type], types)
			if err != nil {
				return abi.Type{}, err
			}
			if abi.ToCamelCase(c.Label) == "" {
				return abi.Type{}, errors.New("abi: purely anonymous or underscored field is not supported")
			}
			fields = append(fields, reflect.StructField{
				Name: abi.ToCamelCase(c.Label), // reflect.StructOf will panic for any exported field.
				Type: cType.Type,
				Tag:  reflect.StructTag("json:\"" + c.Label + "\"" + " " + fmt.Sprintf(`storage:"slot=%s,offset=%d"`, c.Slot.String(), c.Offset)),
			})
			elems = append(elems, &cType)
			names = append(names, c.Label)
		}
		typ.Kind = reflect.Struct
		typ.Type = reflect.StructOf(fields)
		typ.TupleElems = elems
		typ.TupleRawNames = names
		typ.T = abi.TupleTy
		typ.TupleRawName = parsedType[4]

	case "mapping":
		if definition.Encoding != "mapping" {
			return abi.Type{}, fmt.Errorf("abi: invalid encoding for %v, expected %v but got %v", t, "mapping", definition.Encoding)
		}

		var (
			elems []*abi.Type
		)
		typ.Kind = reflect.Map
		typ.T = abi.MappingTy
		kType, err := NewType(definition.Key, types[definition.Key], types)
		if err != nil {
			return abi.Type{}, err
		}

		vType, err := NewType(definition.Value, types[definition.Value], types)
		if err != nil {
			return abi.Type{}, err
		}

		switch vType.Type.Kind() {
		case reflect.Ptr, reflect.Slice, reflect.Map:
		default:
			vType.Type = reflect.PtrTo(vType.Type)
			vType.Kind = vType.Type.Kind()
		}
		elems = append(elems, &kType, &vType)

		typ.Type = reflect.MapOf(kType.Type, vType.Type)
		typ.TupleElems = elems

	default:
		return abi.Type{}, fmt.Errorf("unsupported arg type: %s", t)
	}

	return
}

// reflectIntKind returns the reflect using the given size and
// unsignedness.
func reflectIntKindAndType(unsigned bool, size int) (reflect.Kind, reflect.Type) {
	switch size {
	case 8:
		if unsigned {
			return reflect.Uint8, uint8T
		}
		return reflect.Int8, int8T
	case 16:
		if unsigned {
			return reflect.Uint16, uint16T
		}
		return reflect.Int16, int16T
	case 32:
		if unsigned {
			return reflect.Uint32, uint32T
		}
		return reflect.Int32, int32T
	case 64:
		if unsigned {
			return reflect.Uint64, uint64T
		}
		return reflect.Int64, int64T
	}
	return reflect.Ptr, bigT
}
