package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"gopkg.in/urfave/cli.v1"

	"github.com/ethereum/go-ethereum/cmd/storegen/verifier"
)

var (
	outputFlag = cli.StringFlag{
		Name:  "output",
		Usage: "output file",
	}
	arrayFlag = cli.IntFlag{
		Name:  "array",
		Usage: "array elements count",
		Value: 1,
	}
	mapFlag = cli.IntFlag{
		Name:  "map",
		Usage: "map elements count",
		Value: 1,
	}
)

var commandVerify = cli.Command{
	Name:      "verify",
	Usage:     "verify generated storage wrapper",
	ArgsUsage: "[ <wrapperfile> ]",
	Description: `
Verify generated wrapper file.

`,
	Flags: []cli.Flag{
		outputFlag,
		pkgFlag,
		arrayFlag,
		mapFlag,
	},
	Action: func(ctx *cli.Context) error {
		return verify(ctx)
	},
}

const tmplSource = `
// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package {{.Package}}

import (
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"time"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/rawdb"
    "github.com/ethereum/go-ethereum/core/state"

{{.Imports}}
)
type Storage = wrapper.Storage_Global_t
type UpdateFn func(storage *Storage)


var cases = []UpdateFn {
	{{.Cases}},
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func applyStatements(statements []UpdateFn) common.Hash {
	db := rawdb.NewMemoryDatabase()
	stateDB, _ := state.New(common.Hash{}, state.NewDatabase(db), nil)

	addr := common.Address{}
	var global wrapper.Global_t
	storage := wrapper.New(&global, stateDB, addr, big.NewInt(0))
	for _, f := range cases {
		f(storage)
	}
	hash := stateDB.IntermediateRoot(false)
	return hash
}

const COUNT = 100
func main() {
	last := common.Hash{}
	for i := 0; i < COUNT; i++ {
		hash := applyStatements(cases)
		if i == 0 {
			last = hash
		}

		if hash != last {
			log.Fatal("hash do not match, last: %v, %d: %v", last, i, hash)
		}
	}
	fmt.Printf("running %d times, hash is %s each time\n", COUNT, last.String())
}
`

type UpdateFn func(name string) error

var cases = []UpdateFn{
	func(name string) error { return errors.New("no problem") },
}

func verify(c *cli.Context) error {
	output := c.String(outputFlag.Name)
	pkg := c.String(pkgFlag.Name)
	arrayLimitation := c.Int(arrayFlag.Name)
	mapLimitation := c.Int(mapFlag.Name)

	var fout *os.File
	if output == "-" || output == "" {
		fout = os.Stdout
	} else {
		var err error
		fout, err = os.Create(output)
		if err != nil {
			log.Fatalf("failed to open file %s: %v", output, err)
		}
	}

	if pkg == "" {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		pkg = path.Base(dir)
	}

	filename := c.Args().First()

	fset := token.NewFileSet() // positions are relative to fset
	// Parse src but stop after processing the imports.
	f, err := parser.ParseFile(fset, filename, nil, parser.AllErrors)
	if err != nil {
		log.Fatalf("failed to parse file %s: %v", filename, err)
	}

	buffer := new(bytes.Buffer)
	global := f.Scope.Objects["Global_t"]
	Statements := make([]string, 0)
	visitor := verifier.Visitor{
		Path:       "storage",
		File:       f,
		Statements: &Statements,
		ArrayLimit: arrayLimitation,
		MapLimit:   mapLimitation,
	}
	ast.Walk(visitor, global.Decl.(*ast.TypeSpec))

	imports := fmt.Sprintf(`wrapper "%s"`, verifier.FullPackagePath(filename))

	Cases := []string{}
	for _, statement := range Statements {
		line := fmt.Sprintf("func(storage *Storage) { %v}", statement)
		Cases = append(Cases, line)
	}
	data := struct {
		Cases   string
		Package string
		Imports string
	}{
		Cases:   strings.Join(Cases, ",\n"),
		Package: pkg,
		Imports: imports,
	}

	funcs := map[string]interface{}{}
	buffer = new(bytes.Buffer)
	tmpl := template.Must(template.New("").Funcs(funcs).Parse(tmplSource))
	if err := tmpl.Execute(buffer, data); err != nil {
		log.Fatalf("failed to execute template:%v", err)
	}

	// For Go bindings pass the code through gofmt to clean it up
	code, err := format.Source(buffer.Bytes())
	if err != nil {
		fmt.Fprintf(fout, buffer.String())
		log.Fatalf("failed to format code: %v\n%v", err, buffer.String())
	}

	fmt.Fprint(fout, string(code))
	fout.Close()
	return nil
}
