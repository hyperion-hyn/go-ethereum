package wrapper

import (
	"bytes"
	"crypto"
	"fmt"
	"go/format"
	"reflect"
	"regexp"
	"strings"
	"text/template"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func Store(types []string, layouts []string, pkg string, lang bind.Lang) (string, error) {
	var (
		// structs is the map of all reclared structs shared by passed contracts.
		declarations = make(map[string]*tmplStruct)
		structs      = make(map[string]*tmplStruct)
		defines      = make(map[string]*tmplStruct)
		basics       = make(map[string]*tmplStruct)
		variables    = make(map[string]abi.Type)
	)

	for i := 0; i < len(types); i++ {
		// Parse the actual ABI to generate the binding for
		layout, err := JSON(strings.NewReader(layouts[i]))
		if err != nil {
			return "", err
		}

		for _, typ := range layout.Types {
			_, err := bindStructType[lang](typ, declarations)
			if err != nil {
				return "", err
			}
		}

		for name, typ := range declarations {
			switch typ.T {
			case abi.TupleTy:
				structs[name] = typ
			default:
				if len(typ.Fields) == 0 {
					basics[name] = typ
				} else {
					defines[name] = typ
				}
			}
		}

		for _, variable := range layout.Storage {
			variables[variable.Label] = variable.Type
		}
	}

	// Generate the contract template data content and render it
	data := &tmplData{
		Package:   pkg,
		Defines:   defines,
		Structs:   structs,
		Basics:    basics,
		Variables: variables,
	}
	buffer := new(bytes.Buffer)

	// enhance
	isBasicTypeWrapper := func (val interface{}) bool {
		return isBasicOrDefineType(val, basics)
	}
	isDefineTypeWrapper := func (val interface{}) bool {
		return isBasicOrDefineType(val, defines)
	}
	// enhance - end

	funcs := map[string]interface{}{
		"bindtype":             bindType[lang],
		"capitalise":           capitalise,
		"decapitalise":         decapitalise,
		"sha1":                 sha1,
		"isptr":                isptr,
		"isarray":              isarray,
		"isslice":              isslice,
		"ismap":                ismap,
		"isFixedSizeByteArray": isFixedSizeByteArray,
		"match":                match,
		"GetReflectType":       GetReflectType,

		// enhance
		"isBasicType":          isBasicTypeWrapper,
		"isDefineType":         isDefineTypeWrapper,
	}
	tmpl := template.Must(template.New("").Funcs(funcs).Parse(tmplSource[lang]))
	if err := tmpl.Execute(buffer, data); err != nil {
		return "", err
	}
	// For Go bindings pass the code through gofmt to clean it up
	if lang == bind.LangGo {
		code, err := format.Source(buffer.Bytes())
		if err != nil {
			return "", fmt.Errorf("%v\n%s", err, buffer)
		}
		return string(code), nil
	}
	// For all others just return as is for now
	return buffer.String(), nil
}

// methodNormalizer is a name transformer that modifies Solidity method names to
// conform to target language naming concentions.
var methodNormalizer = map[bind.Lang]func(string) string{
	bind.LangGo: abi.ToCamelCase,
}

// capitalise makes a camel-case string which starts with an upper case character.
func capitalise(input string) string {
	return abi.ToCamelCase(input)
}

// decapitalise makes a camel-case string which starts with a lower case character.
func decapitalise(input string) string {
	if len(input) == 0 {
		return input
	}

	goForm := abi.ToCamelCase(input)
	return strings.ToLower(goForm[:1]) + goForm[1:]
}

func sha1(input string) string {
	h := crypto.SHA1.New()
	h.Write([]byte(input))
	d := h.Sum(nil)
	return hexutil.Encode(d)
}

func isT(val interface{}, T byte) bool {
	switch v := val.(type) {
	case abi.Type:
		return v.T == T
	case *abi.Type:
		return v.T == T
	case tmplField:
		return v.SolKind.T == T
	case *tmplField:
		return v.SolKind.T == T
	case tmplStruct:
		return v.SolKind.T == T
	case *tmplStruct:
		return v.SolKind.T == T
	default:
		return false
	}
}

func isptr(val interface{}) bool {
	return isT(val, abi.PointerTy)
}

func isarray(val interface{}) bool {
	return isT(val, abi.ArrayTy)
}

func isslice(val interface{}) bool {
	return isT(val, abi.SliceTy)
}

func ismap(val interface{}) bool {
	return isT(val, abi.MappingTy)
}

func isFixedSizeByteArray(val interface{}) bool {
	switch v := val.(type) {
	case reflect.Type:
		if v.Kind() == reflect.Array {
			switch v.Elem().Kind() {
			case reflect.Uint8, reflect.Uint16, reflect.Uint64:
				return true
			case reflect.Int8, reflect.Int16, reflect.Int64:
				return true
			}
		}
		return false
	default:
		return false
	}
}

func match(value string, pattern string) bool {
	re := regexp.MustCompile(pattern)
	return re.MatchString(value)
}


// enhance
func isBasicOrDefineType(val interface{}, basicOrDefineTypes map[string]*tmplStruct) bool {
	switch v := val.(type) {
	case tmplField:
		_, ok := basicOrDefineTypes[v.Type]
		return ok
	case *tmplField:
		_, ok := basicOrDefineTypes[v.Type]
		return ok
	case tmplStruct:
		_, ok := basicOrDefineTypes[v.Type]
		return ok
	case *tmplStruct:
		_, ok := basicOrDefineTypes[v.Type]
		return ok
	default:
		return false
	}
}
