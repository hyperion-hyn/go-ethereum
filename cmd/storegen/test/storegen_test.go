package test

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"testing"

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

	filename := "./data/build/contracts/Map3Pool.json"
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
	state, err := sim.Blockchain().State()
	if err != nil {
		t.Errorf("could not get a new mutable state based on the current HEAD block")
	}

	var global Global_t
	storage := New(&global, state, addr, big.NewInt(0))

	{
		// .Version
		wrapper, err := NewMap3PoolWrapper(addr, sim)
		if err != nil {
			t.Errorf("could not new a StorageWrapper: %v", err)
		}

		versionContract, err := wrapper.Version(nil)
		if err != nil {
			t.Errorf("failed to call function Version, %v", err)
		}

		versionStorage := storage.Version().Value()
		log.Debug("TEST", "versionStorage", versionStorage)
		if versionStorage.Cmp(versionContract) != 0 {
			t.Errorf("response from calling contract was expected to be %v instead received %v", versionContract, versionStorage)
		}
	}

	{
		// .Name
		nameStorage := storage.Name().Value()
		expected := "Hyperion"
		if nameStorage != expected {
			t.Errorf("response from calling contract was expected to be '%v' %d instead received '%v' %d", expected, len(expected), nameStorage, len(nameStorage))
		}

		if nameStorage != global.Name {
			t.Errorf(" field expected to be %v instead received %v", global.Name, nameStorage)
		}
	}

	{
		// .Node.NodeAddress
		nodeAddressStorage := storage.Node().NodeAddress().Value()
		expected := common.HexToAddress("0xA07306b4d845BD243Da172aeE557893172ccd04a")
		if nodeAddressStorage != expected {
			t.Errorf("response from calling contract was expected to be %v instead received %v", expected, nodeAddressStorage)
		}
	}

	{
		// .Node.Commission.CommissionRates.Rate
		rateStorage := storage.Node().Commission().CommissionRates().Rate().Value()
		expected := big.NewInt(0).Mul(big.NewInt(0x33), big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil))
		if rateStorage.Cmp(expected) != 0 {
			t.Errorf("response from calling contract was expected to be %v instead received %v", expected, rateStorage)
		}

		if global.Node.Commission.CommissionRates.Rate.Cmp(expected) != 0 {
			t.Errorf("field expected to be %v instead %v", expected, global.Node.Commission.CommissionRates.Rate)
		}
	}

	{
		{
			// .Node.Commission.CommissionRates.MaxRate
			expected, _ := common.NewDecFromStr("5.11")
			rateStorage := storage.Node().Commission().CommissionRates().MaxRate().Value()
			if !rateStorage.Equal(expected) {
				t.Errorf("response from calling contract was expected to be %v instead received %v", expected, rateStorage)
			}

			if !global.Node.Commission.CommissionRates.MaxRate.Equal(expected) {
				t.Errorf("field expected to be %v instead %v", expected, global.Node.Commission.CommissionRates.MaxRate)
			}
		}

		{
			// .Node.Commission.CommissionRates.MaxRate
			expected, _ := common.NewDecFromStr("7.7")
			storage.Node().Commission().CommissionRates().MaxRate().SetValue(expected)

			rateStorage := storage.Node().Commission().CommissionRates().MaxRate().Value()
			if !rateStorage.Equal(expected) {
				t.Errorf("response from calling contract was expected to be %v instead received %v", expected, rateStorage)
			}

			if !global.Node.Commission.CommissionRates.MaxRate.Equal(expected) {
				t.Errorf("field expected to be %v instead %v", expected, global.Node.Commission.CommissionRates.MaxRate)
			}
		}
	}

	{
		// .Node.Description.Serial
		serialStorage := storage.Node().Description().Serial().Value()
		expected, _ := hex.DecodeString("123456789a")
		if bytes.Compare(serialStorage[:], expected) != 0 {
			t.Errorf("response from calling contract was expected to be %v instead received %v", expected, serialStorage)
		}

		expected, _ = hex.DecodeString("deadbeef00")
		copy(serialStorage[:], expected)
		storage.Node().Description().Serial().SetValue(serialStorage)

		serialStorage = storage.Node().Description().Serial().Value()
		if bytes.Compare(serialStorage[:], expected) != 0 {
			t.Errorf("response from calling contract was expected to be %v instead received %v", expected, serialStorage)
		}
	}

	{
		// .Node.Description.Symbol
		var expected [9]byte
		for i := 0; i < len(expected); i++ {
			expected[i] = byte(i & 0xff)
		}

		signatureStorage := storage.Node().Description().Symbol().Value()
		if bytes.Compare(signatureStorage[:], expected[:]) != 0 {
			t.Errorf("response from calling contract was expected to be %v instead received %v", expected, signatureStorage)
		}

		for i := 0; i < len(expected); i++ {
			expected[i] = byte(255 - i)
		}
		storage.Node().Description().Symbol().SetValue(expected)

		signatureStorage = storage.Node().Description().Symbol().Value()
		if bytes.Compare(signatureStorage[:], expected[:]) != 0 {
			t.Errorf("response from calling contract was expected to be %v instead received %v", expected, signatureStorage)
		}
	}

	{
		// .Node.Description.Signature
		var expected [300]byte
		for i := 0; i < len(expected); i++ {
			expected[i] = byte(i & 0xff)
		}

		signatureStorage := storage.Node().Description().Signature().Value()
		if bytes.Compare(signatureStorage[:], expected[:]) != 0 {
			t.Errorf("response from calling contract was expected to be %v instead received %v", expected, signatureStorage)
		}

		for i := 0; i < len(expected); i++ {
			expected[i] = byte(255 - i)
		}
		storage.Node().Description().Signature().SetValue(expected)

		signatureStorage = storage.Node().Description().Signature().Value()
		if bytes.Compare(signatureStorage[:], expected[:]) != 0 {
			t.Errorf("response from calling contract was expected to be %v instead received %v", expected, signatureStorage)
		}
	}

	{
		// .Node.Description.Feature
		var expected [4]uint32
		expected[0] = 0xdeadbeef
		expected[1] = 0xbeeddeed
		expected[2] = 0xfacecafe
		expected[3] = 0xfeedc0de

		for i := 0; i < len(expected); i++ {
			featureStorage := storage.Node().Description().Feature().Get(uint64(i)).Value()
			if featureStorage != expected[i] {
				t.Errorf("response from calling contract was expected to be %x instead received %x", expected[i], featureStorage)
			}
		}

		expected[0], expected[3] = expected[3], expected[0]
		for i := 0; i < len(expected); i++ {
			storage.Node().Description().Feature().Get(uint64(i)).SetValue(expected[i])
		}

		for i := 0; i < len(expected); i++ {
			featureStorage := storage.Node().Description().Feature().Get(uint64(i)).Value()
			if featureStorage != expected[i] {
				t.Errorf("response from calling contract was expected to be %x instead received %x", expected[i], featureStorage)
			}
		}
	}

	{
		// .Node.Description.Version[0]
		versionStorage := storage.Node().Description().Version().Get(0).Value()
		expected := big.NewInt(0xbeef)
		if versionStorage.Cmp(expected) != 0 {
			t.Errorf("response from calling contract was expected to be %v instead received %v", expected, versionStorage)
		}

		if versionStorage.Cmp(global.Node.Description.Version[0]) != 0 {
			t.Errorf(" field expected to be %v instead received %v", global.Node.Description.Version[0], versionStorage)
		}
	}

	{
		// .Node.Description.Version[1]
		versionStorage := storage.Node().Description().Version().Get(1).Value()
		expected := big.NewInt(0xdead)
		if versionStorage.Cmp(expected) != 0 {
			t.Errorf("response from calling contract was expected to be %v instead received %v", expected, versionStorage)
		}

		if versionStorage.Cmp(global.Node.Description.Version[1]) != 0 {
			t.Errorf(" field expected to be %v instead received %v", global.Node.Description.Version[1], versionStorage)
		}
	}

	{
		// .Node.Description.Name
		nameStorage := storage.Node().Description().Name().Value()
		expected := "Hyperion - 海伯利安"
		if nameStorage != expected {
			t.Errorf("response from calling contract was expected to be %v instead received %v", expected, nameStorage)
		}

		if nameStorage != global.Node.Description.Name {
			t.Errorf(" field expected to be %v instead received %v", global.Node.Description.Name, nameStorage)
		}
	}

	{
		// .Node.Description.Name
		// Set/Get
		expected := "Map3"
		storage.Node().Description().Name().SetValue(expected)
		if expected != global.Node.Description.Name {
			t.Errorf(" field expected to be %v instead received %v", expected, global.Node.Description.Name)
		}

		nameStorage := storage.Node().Description().Name().Value()
		if nameStorage != expected {
			t.Errorf("response from calling contract was expected to be %v instead received %v", expected, nameStorage)
		}

		if nameStorage != global.Node.Description.Name {
			t.Errorf(" field expected to be %v instead received %v", global.Node.Description.Name, nameStorage)
		}
	}

	{
		// .Node.Description.Details
		nameStorage := storage.Node().Description().Details().Value()
		expected := "Hyperion, a decentralized map platform, aims to achieve the “One Map” vision - to provide an unified view of global map data and service, and to make it universally accessible just like a public utility for 10B people.\n海伯利安是去中心化的地图生态。"
		if nameStorage != expected {
			t.Errorf("response from calling contract was expected to be %v instead received %v", expected, nameStorage)
		}

		if nameStorage != global.Node.Description.Details {
			t.Errorf(" field expected to be %v instead received %v", global.Node.Description.Details, nameStorage)
		}
	}

	{
		// .Node.Description.Details
		// Set/Get
		expected := "Map3 is a decentralized map service network to safeguard Digital Location Autonomy.\nMap3是去中心化地图服务网络。"
		storage.Node().Description().Details().SetValue(expected)
		if expected != global.Node.Description.Details {
			t.Errorf(" field expected to be %v instead received %v", expected, global.Node.Description.Details)
		}

		nameStorage := storage.Node().Description().Details().Value()
		if nameStorage != expected {
			t.Errorf("response from calling contract was expected to be %v instead received %v", expected, nameStorage)
		}

		if nameStorage != global.Node.Description.Details {
			t.Errorf(" field expected to be %v instead received %v", global.Node.Description.Details, nameStorage)
		}
	}

	{
		// .Node.NodeKeys
		nodeKeysStorage := storage.Node().NodeKeys().Value()
		expected := []byte("MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDQBkQd2vUJtyNa2MBw4i8S0N9kQAAHwWdr1D5CPWgv/9GsGVCAUmLZhLV6E5JcrsL3fcKpak+oO+X3chffgOANVolvwqPUJif1ciimoMiEOU7+auLhTpRohX44phoCJ7J9C1nklTx1L6YHDrnMpvlAuRf0V6HM5Ro0L56LUMwZmwIDAQAB")
		if bytes.Compare(nodeKeysStorage, expected) != 0 {
			t.Errorf("response from calling contract was expected to be %v instead received %v", expected, nodeKeysStorage)
		}

		if bytes.Compare(nodeKeysStorage, global.Node.NodeKeys) != 0 {
			t.Errorf(" field expected to be %v instead received %v", global.Node.NodeKeys, nodeKeysStorage)
		}
	}

	{
		// .Node.NodeKeys
		// Set/Get
		expected := []byte("MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCRRM4LWDW6x/8rHP0yte94a2LG17+6niq0uBq8h5AnwB5v6N0tHOoAA5nz18EkD4Lvp/NyUPCaAmWZyFQ3eHR5iv4bUItt5PJWbFGXSMWOxZyeoZjylK+V8fpbgjHq9a4JlMLzWtGJ/6f5/49uVXaUsfSiDL+zJawrdAjiM5/xyQIDAQAB")
		storage.Node().NodeKeys().SetValue(expected)
		if bytes.Compare(expected, global.Node.NodeKeys) != 0 {
			t.Errorf(" field expected to be %v instead received %v", expected, global.Node.NodeKeys)
		}

		nodeKeysStorage := storage.Node().NodeKeys().Value()
		if bytes.Compare(nodeKeysStorage, expected) != 0 {
			t.Errorf("response from calling contract was expected to be %v instead received %v", expected, nodeKeysStorage)
		}

		if bytes.Compare(nodeKeysStorage, global.Node.NodeKeys) != 0 {
			t.Errorf(" field expected to be %v instead received %v", global.Node.NodeKeys, nodeKeysStorage)
		}
	}

	{
		// .Node[].TotalDelegation
		//
		// pool.Nodes[0xA07306b4d845BD243Da172aeE557893172ccd04a].TotalDelegation = 0xdeadbeef;

		// Set/Get
		expected := 0xdeadbeef
		addr1 := common.HexToAddress("A07306b4d845BD243Da172aeE557893172ccd04a")
		amountStorage := storage.Pool().Nodes().Get(addr1).TotalDelegation().Value()
		if amountStorage.Cmp(big.NewInt(0).SetUint64(uint64(expected))) != 0 {
			t.Errorf("response from calling contract was expected to be %v instead received %v", expected, amountStorage)
		}

		if amountStorage.Cmp(global.Pool.Nodes[addr1].TotalDelegation) != 0 {
			t.Errorf(" field expected to be %v instead received %v", global.Pool.Nodes[addr1].TotalDelegation, amountStorage)
		}
	}

	{
		// .Node[].Microdelegations[].PendingDelegationsfixed2dimension[][].Amount
		//
		// pool.Nodes[0xA07306b4d845BD243Da172aeE557893172ccd04a].Microdelegations[0x3CB0B0B6D52885760A5404eb0A593B979c88BcEF].PendingDelegationsfixed2dimension[2][1].Amount = 0xbeef;
		// pool.Nodes[0xA07306b4d845BD243Da172aeE557893172ccd04a].Microdelegations[0x3CB0B0B6D52885760A5404eb0A593B979c88BcEF].PendingDelegationsfixed2dimension[0][0].Amount = 0xdead;

		// Set/Get
		{
			expected := 0xdead
			addr1 := common.HexToAddress("A07306b4d845BD243Da172aeE557893172ccd04a")
			addr2 := common.HexToAddress("3CB0B0B6D52885760A5404eb0A593B979c88BcEF")
			amountStorage := storage.Pool().Nodes().Get(addr1).Microdelegations().Get(addr2).PendingDelegationsfixed2dimension().Get(0).Get(0).Amount().Value()
			if amountStorage.Cmp(big.NewInt(0).SetUint64(uint64(expected))) != 0 {
				t.Errorf("response from calling contract was expected to be %v instead received %v", expected, amountStorage)
			}

			if amountStorage.Cmp(global.Pool.Nodes[addr1].Microdelegations[addr2].PendingDelegationsfixed2dimension[0][0].Amount) != 0 {
				t.Errorf(" field expected to be %v instead received %v", global.Pool.Nodes[addr1].Microdelegations[addr2].PendingDelegationsfixed2dimension[0][0].Amount, amountStorage)
			}
		}

		{
			expected := 0xbeef
			addr1 := common.HexToAddress("A07306b4d845BD243Da172aeE557893172ccd04a")
			addr2 := common.HexToAddress("3CB0B0B6D52885760A5404eb0A593B979c88BcEF")
			amountStorage := storage.Pool().Nodes().Get(addr1).Microdelegations().Get(addr2).PendingDelegationsfixed2dimension().Get(2).Get(1).Amount().Value()
			if amountStorage.Cmp(big.NewInt(0).SetUint64(uint64(expected))) != 0 {
				t.Errorf("response from calling contract was expected to be %v instead received %v", expected, amountStorage)
			}

			if amountStorage.Cmp(global.Pool.Nodes[addr1].Microdelegations[addr2].PendingDelegationsfixed2dimension[2][1].Amount) != 0 {
				t.Errorf(" field expected to be %v instead received %v", global.Pool.Nodes[addr1].Microdelegations[addr2].PendingDelegationsfixed2dimension[2][1].Amount, amountStorage)
			}
		}
	}

	{
		// .Version
		wrapper, err := NewMap3PoolWrapper(addr, sim)
		if err != nil {
			t.Errorf("could not new a StorageWrapper: %v", err)
		}

		lengthContract, err := wrapper.Length(nil)
		if err != nil {
			t.Errorf("failed to call function Length, %v", err)
		}

		expected := big.NewInt(10)
		if expected.Cmp(lengthContract) != 0 {
			t.Errorf("response from calling contract was expected to be %v instead received %v", lengthContract, expected)
		}
	}

	{
		// .Node[].Microdelegations[].PendingDelegations[].Amount
		//
		// pool.Nodes[0xA07306b4d845BD243Da172aeE557893172ccd04a].Microdelegations[0x3CB0B0B6D52885760A5404eb0A593B979c88BcEF].PendingDelegations[5].Amount = 0xdeaf;

		// Set/Get
		{
			addr1 := common.HexToAddress("A07306b4d845BD243Da172aeE557893172ccd04a")
			addr2 := common.HexToAddress("3CB0B0B6D52885760A5404eb0A593B979c88BcEF")

			// length
			{
				expected := int64(10)
				lengthStorage := storage.Pool().Nodes().Get(addr1).Microdelegations().Get(addr2).PendingDelegations().Length()

				if big.NewInt(expected).Cmp(lengthStorage) != 0 {
					t.Errorf(" length expected to be %v instead received %v", expected, lengthStorage.Uint64())
				}
			}

			// expand
			{
				storage.Pool().Nodes().Get(addr1).Microdelegations().Get(addr2).PendingDelegations().Get(20).Amount().Value()

				expected := int64(21)
				lengthStorage := storage.Pool().Nodes().Get(addr1).Microdelegations().Get(addr2).PendingDelegations().Length()
				if big.NewInt(expected).Cmp(lengthStorage) != 0 {
					t.Errorf(" length expected to be %v instead received %v", expected, lengthStorage.Uint64())
				}

				if int(expected) != len(global.Pool.Nodes[addr1].Microdelegations[addr2].PendingDelegations) {
					t.Errorf(" length expected to be %v instead %v", expected, len(global.Pool.Nodes[addr1].Microdelegations[addr2].PendingDelegations))
				}
			}

			// shrink
			{
				expected := int64(15)
				storage.Pool().Nodes().Get(addr1).Microdelegations().Get(addr2).PendingDelegations().Resize(15)
				lengthStorage := storage.Pool().Nodes().Get(addr1).Microdelegations().Get(addr2).PendingDelegations().Length()

				if big.NewInt(expected).Cmp(lengthStorage) != 0 {
					t.Errorf(" length expected to be %v instead received %v", expected, lengthStorage.Uint64())
				}

				if int(expected) != len(global.Pool.Nodes[addr1].Microdelegations[addr2].PendingDelegations) {
					t.Errorf(" length expected to be %v instead %v", expected, len(global.Pool.Nodes[addr1].Microdelegations[addr2].PendingDelegations))
				}
			}

			{
				expected := 0x7788

				amountStorage := storage.Pool().Nodes().Get(addr1).Microdelegations().Get(addr2).PendingDelegations().Get(5).Amount().Value()
				if amountStorage.Cmp(big.NewInt(0).SetUint64(uint64(expected))) != 0 {
					t.Errorf("response from calling contract was expected to be %v instead received %v", expected, amountStorage)
				}

				if amountStorage.Cmp(global.Pool.Nodes[addr1].Microdelegations[addr2].PendingDelegations[5].Amount) != 0 {
					t.Errorf(" field expected to be %v instead received %v", global.Pool.Nodes[addr1].Microdelegations[addr2].PendingDelegations[5].Amount, amountStorage)
				}
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
