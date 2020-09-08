package storage

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/hyperion-hyn/bls/ffi/go/bls"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

type BuiltContract struct {
	ABI              interface{} `json:"abi"`
	Bytecode         string      `json:"bytecode"`
	DeployedBytecode string      `json:"deployedBytecode"`
}

var testKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")

var (
	abiJSON      string
	abiBin       string
	deployedCode string
)

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(log.LvlDebug), log.StreamHandler(os.Stdout, log.TerminalFormat(true))))

	filename := "../build/contracts/consortium.json"
	// Read entire file content, giving us little control but
	// making it very simple. No need to close the file.
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("failed to read " + filename)
	}

	var builtContract BuiltContract
	err = json.Unmarshal(content, &builtContract)
	if err != nil {
		log.Debug("init", "Unmarshal failed", err, "content", content)
	}

	abi, _ := json.Marshal(builtContract.ABI)
	log.Debug("init", "abi", string(abi))
	log.Debug("init", "bytecode", builtContract.Bytecode)
	log.Debug("init", "deployedBytecode", builtContract.DeployedBytecode)

	abiJSON = string(abi)
	abiBin = builtContract.Bytecode
	deployedCode = builtContract.DeployedBytecode
}

func setupBlockchain(t *testing.T, abiJSON string, abiBin string) (common.Address, *backends.SimulatedBackend, context.Context) {
	testAddr := crypto.PubkeyToAddress(testKey.PublicKey)
	sim := backends.NewSimulatedBackend(
		core.GenesisAlloc{
			testAddr: {Balance: big.NewInt(10000000000)},
		},
		100000000,
	)

	bgCtx := context.Background()

	parsed, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		t.Errorf("could not get code at test addr: %v", err)
	}
	contractAuth := bind.NewKeyedTransactor(testKey)
	contractAuth.GasLimit = 10000000
	addr, _, _, err := bind.DeployContract(contractAuth, parsed, common.FromHex(abiBin), sim)
	log.Debug("setup", "deployed", addr, "test", testAddr)
	if err != nil {
		t.Errorf("could not deploy contract: %v", err)
	}

	sim.Commit()

	return addr, sim, bgCtx
}

func testReadViaStorageAndWriteFromContract(t *testing.T, sim *backends.SimulatedBackend, addr common.Address) {
	stateDB, err := sim.Blockchain().State()
	if err != nil {
		t.Errorf("could not get a new mutable state based on the current HEAD block")
	}

	var global Global_t
	storage := New(&global, stateDB, addr, big.NewInt(0))

	{
		// .Version
		wrapper, err := NewConsortium(addr, sim)
		if err != nil {
			t.Errorf("could not new a StorageWrapper: %v", err)
		}

		versionContract, err := wrapper.Version(nil)
		if err != nil {
			t.Errorf("failed to call function Version, %v", err)
		}

		expected := int32(22)
		if versionContract != expected {
			t.Errorf("Version expected %v instead %v", expected, versionContract)
		}
	}

	{
		// .committee
		type Validator struct {
			acctPrivKey   *ecdsa.PrivateKey
			acctPubKey    ecdsa.PublicKey
			signerPrivKey *bls.SecretKey
			signerPubKey  *bls.PublicKey
		}

		const numberOfCandidates = 7
		var candidates []Validator = make([]Validator, numberOfCandidates)
		membersStorage := storage.Committee().Members()
		for i := 0; i < numberOfCandidates; i++ {
			acctPrivKey, _ := crypto.GenerateKey()
			signerPrivKey, _ := crypto.GenerateBLSKey()

			candidates[i].acctPrivKey = acctPrivKey
			candidates[i].acctPubKey = acctPrivKey.PublicKey

			candidates[i].signerPrivKey = signerPrivKey
			candidates[i].signerPubKey = signerPrivKey.GetPublicKey()
		}

		for i := 0; i < numberOfCandidates; i++ {
			currentCandidate := candidates[i]

			// acctPrivKeyHex := hex.EncodeToString(crypto.FromECDSA(currentCandidate.acctPrivKey))
			// acctPubKeyHex := crypto.PubkeyToAddress(currentCandidate.acctPubKey).String()
			// signerPrivKeyHex := hex.EncodeToString(currentCandidate.signerPrivKey.Serialize())
			// signerPubKeyHex := hex.EncodeToString(currentCandidate.signerPubKey.Serialize())
			//
			// t.Logf("%d: acct   (priv): %s", i, acctPrivKeyHex)
			// t.Logf("%d: acct   (addr): %s", i, acctPubKeyHex)
			// t.Logf("%d: signer (priv): %s", i, signerPrivKeyHex)
			// t.Logf("%d: signer (pub) : %s", i, signerPubKeyHex)

			currentStorage := membersStorage.Get(i)
			currentStorage.Coinbase().SetValue(crypto.PubkeyToAddress(currentCandidate.acctPubKey))
			currentStorage.PublicKey().SetValue(currentCandidate.signerPubKey.Serialize())

			acctPubKey := currentStorage.Coinbase().Value()
			if acctPubKey != crypto.PubkeyToAddress(currentCandidate.acctPubKey) {
				t.Errorf("%d expected %s got %s", i, crypto.PubkeyToAddress(currentCandidate.acctPubKey), acctPubKey)
			}

			signerPubKey := currentStorage.PublicKey().Value()
			if bytes.Compare(signerPubKey, currentCandidate.signerPubKey.Serialize()) != 0 {
				t.Errorf("%d expected %s got %s", i, currentCandidate.signerPubKey.Serialize(), signerPubKey)
			}
		}
	}

}

// Tests that storage manipulation
func TestStorageManipulation(t *testing.T) {
	addr, sim, _ := setupBlockchain(t, abiJSON, abiBin)
	defer sim.Close()

	testReadViaStorageAndWriteFromContract(t, sim, addr)
}
