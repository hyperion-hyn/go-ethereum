package verifier

import (
	"encoding/hex"
	"fmt"
	"go/ast"
	"math/big"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

var seed = rand.New(rand.NewSource(99))

const (
	ARRAY_LIMITATION = 1
	MAP_LIMITATION   = 1
)

type Visitor struct {
	Path       string
	File       *ast.File
	Statements *[]string
}

func NewVisitor(path string, visitor Visitor) Visitor {
	retval := Visitor{
		Path:       path,
		File:       visitor.File,
		Statements: visitor.Statements,
	}
	return retval
}

func IsBasicType(typeName string) bool {
	switch typeName {
	case "Address":
		return true
	case "Int8", "Int16", "Int32", "Int64", "Int256",
		"Uint8", "Uint16", "Uint32", "Uint64", "Uint256":
		return true
	case "BigInt":
		return true
	case "Bytes",
		"Bytes1", "Bytes2", "Bytes3", "Bytes4", "Bytes5", "Bytes6", "Bytes7", "Bytes8", "Bytes9", "Bytes10",
		"Bytes11", "Bytes12", "Bytes13", "Bytes14", "Bytes15", "Bytes16", "Bytes17", "Bytes18", "Bytes19", "Bytes110",
		"Bytes21", "Bytes22", "Bytes23", "Bytes24", "Bytes25", "Bytes26", "Bytes27", "Bytes28", "Bytes29", "Bytes210",
		"Bytes31", "Bytes32":
		return true
	case "String":
		return true
	case "Bool":
		return true
	case "Decimal":
		return true
	default:
		return false
	}
}
func getType(expr ast.Expr, File *ast.File) (typeName string, typ *ast.Object, isBasicType bool) {
	switch t := expr.(type) {
	case *ast.StarExpr:
		typeName = t.X.(*ast.Ident).Name
	case *ast.Ident:
		typeName = t.Name
	case *ast.SelectorExpr:
		typeName = fmt.Sprintf("%s.%s", t.X.(*ast.Ident).Name, t.Sel.Name)
	default:
		panic("unhandled")
	}

	typ, _ = File.Scope.Objects[typeName]
	isBasicType = IsBasicType(typeName)
	return typeName, typ, isBasicType
}

func getRandomBytes(n int) []byte {
	retval := make([]byte, n)
	for i := 0; i < n; i++ {
		retval[i] = byte(rand.Intn(255) & 255)
	}
	return retval
}

func BytesToString(data []byte, fixed bool) string {
	items := make([]string, len(data))
	for i, v := range data {
		items[i] = fmt.Sprintf("0x%x", v)
	}

	length := ""
	if fixed {
		length = fmt.Sprintf("%d", len(data))
	}
	return fmt.Sprintf("[%s]byte{%s}", length, strings.Join(items, ","))
}

func getRandomValue(typeName string) string {
	switch typeName {
	case "Address":
		privateKey, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		return fmt.Sprintf(`common.HexToAddress("%s")`, addr.String())
	case "Int8", "Int16", "Int32", "Int64", "Int256",
		"Uint8", "Uint16", "Uint32", "Uint64", "Uint256":
		re := regexp.MustCompile("(?:Int|Uint)([0-9])*")
		matches := re.FindAllStringSubmatch(typeName, -1)
		m := matches[0]
		length, _ := strconv.Atoi(m[1])

		privateKey, _ := crypto.GenerateKey()
		data := crypto.FromECDSA(privateKey)
		val := big.NewInt(0).And(big.NewInt(0).SetBytes(data), big.NewInt(0).Sub(big.NewInt(0).Lsh(big.NewInt(1), uint(length)), big.NewInt(1)))
		return fmt.Sprintf("%v", val)
	case "BigInt":
		privateKey, _ := crypto.GenerateKey()
		return fmt.Sprintf(`func() *big.Int { v, _ := big.NewInt(0).SetString("%s", 16); return v}()`, hex.EncodeToString(crypto.FromECDSA(privateKey)))
	case "Bytes":
		privateKey, _ := crypto.GenerateKey()
		data := crypto.FromECDSA(privateKey)
		return BytesToString(data, false)
	case "Bytes1", "Bytes2", "Bytes3", "Bytes4", "Bytes5", "Bytes6", "Bytes7", "Bytes8", "Bytes9", "Bytes10",
		"Bytes11", "Bytes12", "Bytes13", "Bytes14", "Bytes15", "Bytes16", "Bytes17", "Bytes18", "Bytes19", "Bytes20",
		"Bytes21", "Bytes22", "Bytes23", "Bytes24", "Bytes25", "Bytes26", "Bytes27", "Bytes28", "Bytes29", "Bytes30",
		"Bytes31", "Bytes32":
		re := regexp.MustCompile("Bytes([0-9])*")
		matches := re.FindAllStringSubmatch(typeName, -1)
		m := matches[0]
		length, _ := strconv.Atoi(m[1])
		return BytesToString(getRandomBytes(length), true)
	case "String":
		privateKey, _ := crypto.GenerateKey()
		return fmt.Sprintf(`"%s"`, hex.EncodeToString(crypto.FromECDSA(privateKey)))
	case "Bool":
		n := rand.Intn(2)
		var flag bool
		if n == 0 {
			flag = false
		} else {
			flag = true
		}
		return fmt.Sprintf("%v", flag)
	case "Decimal":
		val := big.NewInt(0)
		val = val.Rand(seed, big.NewInt(0).Lsh(big.NewInt(1), 48))
		return fmt.Sprintf(`func() common.Dec { v, _ := big.NewInt(0).SetString("%s", 16); d := common.NewDecFromBigIntWithPrec(v, 18); return d}()`, val.String())
	default:
		return "Unknown"
	}
}

func (v Visitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.TypeSpec:
		isBasicType := IsBasicType(n.Name.Name)
		if !isBasicType {
			return v
		}
		statement := fmt.Sprintf("%s.SetValue(%v)", v.Path, getRandomValue(n.Name.Name))
		*v.Statements = append(*v.Statements, statement)
		return nil

	case *ast.Field:
		name := n.Names[0].Name
		retval := NewVisitor(fmt.Sprintf("%s.%s()", v.Path, name), v)
		typeName := n.Type.(*ast.Ident).Name
		typ, _ := v.File.Scope.Objects[typeName]
		ast.Walk(retval, typ.Decl.(*ast.TypeSpec))
		return v

	case *ast.ArrayType:
		var length int
		var isFixedSize bool
		if n.Len == nil {
			// slice, limit to 65535, anyone can increase this limitation
			length = ARRAY_LIMITATION
			isFixedSize = false
		} else {
			// array
			length, _ = strconv.Atoi(n.Len.(*ast.BasicLit).Value)
			isFixedSize = true
		}

		typeName, typ, _ := getType(n.Elt, v.File)

		if isFixedSize && typeName == "Uint8" {
			// fixed-size bytes
			val := BytesToString(getRandomBytes(length), true)
			statement := fmt.Sprintf("%s.SetValue(%s)", v.Path, val)
			*v.Statements = append(*v.Statements, statement)
		} else {
			for i := 0; i < length; i++ {
				retval := NewVisitor(fmt.Sprintf("%s.Get(%d)", v.Path, i), v)
				ast.Walk(retval, typ.Decl.(*ast.TypeSpec))
			}
			return v
		}

		return v

	case *ast.MapType:
		keyTypeName, _, _ := getType(n.Key, v.File)
		for i := 0; i < MAP_LIMITATION; i++ {
			keyValue := getRandomValue(keyTypeName)
			visitor := NewVisitor(fmt.Sprintf("%s.Get(%s)", v.Path, keyValue), v)
			_, typ, _ := getType(n.Value, v.File)
			ast.Walk(visitor, typ.Decl.(*ast.TypeSpec))
		}
		return v

	default:
		return v
	}
}
