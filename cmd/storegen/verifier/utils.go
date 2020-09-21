package verifier

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	slashSlash = []byte("//")
	moduleStr  = []byte("module")
)

// ModulePath returns the module path from the gomod file text.
// If it cannot find a module path, it returns an empty string.
// It is tolerant of unrelated problems in the go.mod file.
func ModulePath(mod []byte) string {
	for len(mod) > 0 {
		line := mod
		mod = nil
		if i := bytes.IndexByte(line, '\n'); i >= 0 {
			line, mod = line[:i], line[i+1:]
		}
		if i := bytes.Index(line, slashSlash); i >= 0 {
			line = line[:i]
		}
		line = bytes.TrimSpace(line)
		if !bytes.HasPrefix(line, moduleStr) {
			continue
		}
		line = line[len(moduleStr):]
		n := len(line)
		line = bytes.TrimSpace(line)
		if len(line) == n || len(line) == 0 {
			continue
		}

		if line[0] == '"' || line[0] == '`' {
			p, err := strconv.Unquote(string(line))
			if err != nil {
				return "" // malformed quoted string or multiline module path
			}
			return p
		}

		return string(line)
	}
	return "" // missing module path
}

func isFileExist(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	return false
}

func FullPackagePath(filename string) string {
	filename, _ = filepath.Abs(filename)
	if !isFileExist(filename) {
		return ""
	}

	parts := []string{}
	curdir := path.Dir(filename)

	for i := 0; i < 255; i++ {
		filename := path.Join(curdir, "go.mod")
		_, err := os.Stat(filename)
		if err != nil {
			base := path.Base(curdir)
			parts = append([]string{base}, parts...)
			curdir = path.Dir(curdir)
			if curdir == "/" {
				break
			}
			continue
		} else {
			data, _ := ioutil.ReadFile(filename)
			mod := ModulePath(data)
			log.Printf(":::: %v", mod)
			parts = append([]string{mod}, parts...)
			break
		}
	}
	return strings.Join(parts, "/")
}
