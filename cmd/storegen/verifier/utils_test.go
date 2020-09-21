package verifier

import (
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/log"
)

func TestFullPackagePath(t *testing.T) {
	t.Log("hahahah")
	pkg := FullPackagePath("../../../consensus/atlas/storage/consortium_storage.go")
	expected := "github.com/ethereum/go-ethereum/consensus/atlas/storage"
	if pkg != expected {
		t.Errorf("wrong full package path, exepcted: %v, got: %v", expected, pkg)
	}
}

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(log.LvlDebug), log.StreamHandler(os.Stdout, log.TerminalFormat(true))))
}
