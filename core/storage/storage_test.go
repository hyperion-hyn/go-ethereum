package storage

// generate wrapper_test.go
// cat data/build/contracts/Storage.json | jq -c '.abi' | abigen --abi - --pkg storage --type StorageWrapper --out wrapper_test.go

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)


type Description struct {
	name string `storage:"slot0"`
	url  string `storage:"slot1"`
}

type Delegation struct {
	amount      int `storage:"slot0"`
	blockNumber int `storage:"slot1"`
}

type Validator struct {
	desc        Description  `storage:"slot0"`
	delegations []Delegation `storage:"slot1"`
}

type Donation struct {
	Name   string `storage:"slot0"`
	amount int    `storage:"slot1"`
}
type ValidatorList struct {
	Name       string              `storage:"slot0"`
	author     string              `storage:"slot1"`
	count      int                 `storage:"slot2"`
	Desc       Description         `storage:"slot3"`
	validators []Validator         `storage:"slot4"`
	donations  map[string]Donation `storage:"slot5"`
}

type GlobalVariables struct {
	Version       int           `storage:"slot0"`
	Name          string        `storage:"slot1"`
	ValidatorList ValidatorList `storage:"slot2"`
}

type BuiltContract struct {
	ABI interface{}				`json:"abi"`
	Bytecode string			`json:"bytecode"`
	DeployedBytecode string	`json:"deployedBytecode"`
}

var testKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")

var (
	abiJSON string
	abiBin string
	deployedCode string
)

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(log.LvlDebug), log.StreamHandler(os.Stdout, log.TerminalFormat(true))))

	// Read entire file content, giving us little control but
	// making it very simple. No need to close the file.
	content, err := ioutil.ReadFile("./data/build/contracts/Storage.json")
	if err != nil {
		panic("failed to read ./data/build/contracts/Storage.json")
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

// Tests parseTag
func TestParseTag(t *testing.T) {
	var tests = []struct {
		tag      string // input
		expected int    // expected result
		err      error
	}{
		{"", 0, errors.New(fmt.Sprintf("invalid tag: "))},
		{"slot", 0, errors.New(fmt.Sprintf("invalid tag: slot"))},
		{"slot0", 0, nil},
		{"slot99", 99, nil},
		{"slot-10", 0, errors.New(fmt.Sprintf("invalid tag: slot-10"))},
	}

	for _, tt := range tests {
		actual, err := parseTag(tt.tag)
		if actual != tt.expected {
			if err.Error() != tt.err.Error() {
				t.Errorf("parseTag(%s): expected %d, actual %d, '%v+', '%v+'", tt.tag, tt.expected, actual, tt.err, err)
			}
		}
	}
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

func TestBlockchainViaPackParameters(t *testing.T) {
	addr, sim, bgCtx := setupBlockchain(t, abiJSON, abiBin)
	defer sim.Close()

	testAddr := crypto.PubkeyToAddress(testKey.PublicKey)

	parsed, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		t.Errorf("could not get code at test addr: %v", err)
	}

	input, err := parsed.Pack("Hello")
	if err != nil {
		t.Errorf("could pack receive function on contract: %v", err)
	}

	// make sure you can call the contract
	res, err := sim.CallContract(bgCtx, ethereum.CallMsg{
		From: testAddr,
		To:   &addr,
		Data: input,
	}, nil)
	if err != nil {
		t.Errorf("could not call receive method on contract: %v", err)
	}
	if len(res) == 0 {
		t.Errorf("result of contract call was empty: %v", res)
	}

	type Response struct {
		Res string
	}

	var response Response

	err = parsed.Unpack(&response, "Hello", res)
	if err != nil {
		t.Errorf("could not unpack response %v: %v", err, res)
	}

	if response.Res != "hello world" {
		t.Errorf("response from calling contract was expected to be 'hello world' instead received '%v'", res)
	}
}


func TestBlockchainViaBinding(t *testing.T) {
	addr, sim, _ := setupBlockchain(t, abiJSON, abiBin)
	defer sim.Close()

	wrapper, err := NewStorageWrapper(addr, sim)
	if err != nil {
		t.Errorf("could not new a StorageWrapper: %v", err)
	}

	rv, err := wrapper.Hello(nil)
	if err != nil || rv != "hello world" {
		t.Errorf("response from calling contract was expected to be 'hello world' instead received '%v'", rv)
	}
}

func testSetUnexportedField (t *testing.T, sim *backends.SimulatedBackend, addr common.Address, globalVariables *GlobalVariables) {
	state, _ := sim.Blockchain().State()
	storage := NewStorage(state, addr, 0, globalVariables, nil)

	name := storage.GetByName("ValidatorList").GetByName("Desc").GetByName("name")
	name.SetValue("TheWonderingEarth")

	if globalVariables.ValidatorList.Desc.name != "TheWonderingEarth" {
		t.Errorf("failed to set value via storage")
	}

	// NOTE: no Flush in this case.
}

func testWriteViaStorageAndReadFromContract(t *testing.T, sim *backends.SimulatedBackend, addr common.Address, wrapper *StorageWrapper, globalVariables *GlobalVariables) {
		state, err := sim.Blockchain().State()
		var target int = 0b101010
		storage := NewStorage(state, addr, 0, globalVariables, nil)
		version := storage.Get("Version")
		version.SetValue(target)
		if globalVariables.Version != target {
			t.Errorf("failed to set .Version, expect %v got %v", target, globalVariables.Version)
		}

		// write modifications to statedb
		storage.Flush()
		// flush state in a new block
		if err = sim.FlushStateInNewBlock(storage.StateDB()); err != nil {
			t.Errorf("failed to FlushStateInNewBlock, err: %v", err)
		}

		rv, err := wrapper.Version(nil)
		if err != nil || rv.Cmp(big.NewInt(int64(target))) != 0 {
			t.Errorf("response from Version() was expected to be %v instead received %v, err: %v", target, rv, err)
		}

		target = 7788
		state, err = sim.Blockchain().State()
		state.SetBalance(addr, big.NewInt(int64(target)))
		if err = sim.FlushStateInNewBlock(state); err != nil {
			t.Errorf("failed to FlushStateInNewBlock, err: %v", err)
		}

		rv, err = wrapper.Balance(nil)
		if err != nil  || rv.Cmp(big.NewInt(int64(target))) != 0 {
			t.Errorf("response from Version() was expected to be %v instead received %v, err: %v", target, rv, err)
		}
}


// Tests that storage manipulation
func TestStorageManipulation(t *testing.T) {
	addr, sim, _ := setupBlockchain(t, abiJSON, abiBin)
	defer sim.Close()

	log.Debug("Blockchain", "deployed", addr)

	// smartcontract binding wrapper
	wrapper, err := NewStorageWrapper(addr, sim)
	if err != nil {
		t.Errorf("could not new a StorageWrapper: %v", err)
	}

	var globalVariables GlobalVariables = GlobalVariables{
		ValidatorList: ValidatorList{
			Name:   "atlas",
			author: "hyperion",
			count:  11,
			Desc: Description{
				name: "hyperion",
				url:  "https://www.hyn.space",
			},
			validators: nil,
		},
	}

	if reflect.ValueOf(&globalVariables).Elem().FieldByName("ValidatorList").FieldByName("Name").String() != globalVariables.ValidatorList.Name {
		t.Errorf("FieldByName to retrive globalVariables.ValidatorList.Name got wrong value")
	}
	reflect.Indirect(reflect.ValueOf(&globalVariables)).FieldByName("ValidatorList").FieldByName("Name").SetString("ATLAS")
	if globalVariables.ValidatorList.Name != "ATLAS" {
		t.Errorf("failed to set globalVariables.ValidatorList.Name")
	}


	testSetUnexportedField(t, sim, addr, &globalVariables)
	testWriteViaStorageAndReadFromContract(t, sim, addr, wrapper, &globalVariables)

	state, err := sim.Blockchain().State()
	storage := NewStorage(state, addr, 0, &globalVariables, nil)
	
	log.Debug("TestStorageManipulation", "validatorList", globalVariables)
	// name := storage.GetByName("validators").GetByName("desc").GetByName("name")
	// name := storage.GetByName("validators").GetByName("name")
	{
		name := storage.GetByName("ValidatorList").GetByName("Name")
		log.Debug("result", "validatorList.Name", name.Value())
	}

	{
		name := storage.GetByName("ValidatorList").GetByName("author")
		log.Debug("result", "validatorList.author", name.Value())
	}

	{
		name := storage.GetByName("ValidatorList").GetByName("Desc").GetByName("name")
		log.Debug("result", "validatorList.Desc.name", name.Value())
	}

	{
		name := storage.GetByName("ValidatorList").GetByName("author")
		log.Debug("result", "validatorList.author", name.Value())
		name.SetValue("harmony")
		log.Debug("result", "validatorList.author", globalVariables.ValidatorList.author)
	}

	{
		name := storage.GetByName("ValidatorList").GetByName("count")
		log.Debug("result", "validatorList.count", name.Value())
		name.SetValue(22)
		log.Debug("result", "validatorList.count", globalVariables.ValidatorList.count)
	}

	{
		log.Debug("compare", "validatorList.validators == nil", globalVariables.ValidatorList.validators == nil)
		log.Debug("compare", "len(validatorList.validators)", len(globalVariables.ValidatorList.validators))
		validators := storage.GetByName("ValidatorList").GetByName("validators")
		log.Debug("result", "validatorList.validators", validators.Value())
		vv := validators.Value().([]Validator)
		t := Validator{
			desc: Description{
				name: "temp",
				url:  "http://www.hyn.space",
			},
			delegations: nil,
		}
		vv = append(vv, t)
		log.Debug("result", "validatorList.validators", vv)
		log.Debug("result", "validatorList.validators", globalVariables.ValidatorList.validators)
		validators.GetByIndex(1).SetValue(t)
		log.Debug("result", "validatorList.validators", globalVariables.ValidatorList.validators)

		validators.GetByIndex(2).GetByName("desc").GetByName("name").SetValue("haha")
		log.Debug("result", "validatorList.validators", globalVariables.ValidatorList.validators)
	}

	{
		globalVariables.ValidatorList.donations = make(map[string]Donation)
		globalVariables.ValidatorList.donations["what"] = Donation{
			Name:   "who-donation",
			amount: 8899,
		}
		donations := storage.GetByName("ValidatorList").GetByName("donations")
		log.Debug("result", "validatorList.donations", donations.Value())
		donations.GetByName("what").SetValue(Donation{
			Name:   "who-donation",
			amount: 7788,
		})
		val := donations.GetByName("what").Value().(Donation)
		val.Name = "6688"

		// donations.GetByName("what").GetByName("Name").SetValue("6688")
		// log.Debug("result", "validatorList.donations['what'].name", validatorList.donations["what"].Name)
		// m := validatorList.donations["what"]
		// m.Name = "abc"
		// log.Debug("result", "validatorList.donations['what'].name", validatorList.donations["what"].Name)
		// vv := donations.Value().(map[string]string)
		// t := Validator{
		//     desc:        Description{
		//         name: "temp",
		//         url: "http://www.hyn.space",
		//     },
		//     delegations: nil,
		// }
		// log.Debug("result", "validatorList.validators", vv)
		// log.Debug("result", "validatorList.validators", validatorList.validators)
		// donations.GetByIndex(1).SetValue(t)
		// log.Debug("result", "validatorList.validators", validatorList.validators)
		//
		// donations.GetByIndex(2).GetByName("desc").GetByName("name").SetValue("haha")
		// log.Debug("result", "validatorList.validators", validatorList.validators)
	}
}
