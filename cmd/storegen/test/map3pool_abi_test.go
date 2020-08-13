// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package test

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// Map3PoolWrapperABI is the input ABI used to generate the binding from.
const Map3PoolWrapperABI = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"Length\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"Version\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// Map3PoolWrapper is an auto generated Go binding around an Ethereum contract.
type Map3PoolWrapper struct {
	Map3PoolWrapperCaller     // Read-only binding to the contract
	Map3PoolWrapperTransactor // Write-only binding to the contract
	Map3PoolWrapperFilterer   // Log filterer for contract events
}

// Map3PoolWrapperCaller is an auto generated read-only Go binding around an Ethereum contract.
type Map3PoolWrapperCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Map3PoolWrapperTransactor is an auto generated write-only Go binding around an Ethereum contract.
type Map3PoolWrapperTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Map3PoolWrapperFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Map3PoolWrapperFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Map3PoolWrapperSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Map3PoolWrapperSession struct {
	Contract     *Map3PoolWrapper  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Map3PoolWrapperCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Map3PoolWrapperCallerSession struct {
	Contract *Map3PoolWrapperCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// Map3PoolWrapperTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Map3PoolWrapperTransactorSession struct {
	Contract     *Map3PoolWrapperTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// Map3PoolWrapperRaw is an auto generated low-level Go binding around an Ethereum contract.
type Map3PoolWrapperRaw struct {
	Contract *Map3PoolWrapper // Generic contract binding to access the raw methods on
}

// Map3PoolWrapperCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Map3PoolWrapperCallerRaw struct {
	Contract *Map3PoolWrapperCaller // Generic read-only contract binding to access the raw methods on
}

// Map3PoolWrapperTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Map3PoolWrapperTransactorRaw struct {
	Contract *Map3PoolWrapperTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMap3PoolWrapper creates a new instance of Map3PoolWrapper, bound to a specific deployed contract.
func NewMap3PoolWrapper(address common.Address, backend bind.ContractBackend) (*Map3PoolWrapper, error) {
	contract, err := bindMap3PoolWrapper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Map3PoolWrapper{Map3PoolWrapperCaller: Map3PoolWrapperCaller{contract: contract}, Map3PoolWrapperTransactor: Map3PoolWrapperTransactor{contract: contract}, Map3PoolWrapperFilterer: Map3PoolWrapperFilterer{contract: contract}}, nil
}

// NewMap3PoolWrapperCaller creates a new read-only instance of Map3PoolWrapper, bound to a specific deployed contract.
func NewMap3PoolWrapperCaller(address common.Address, caller bind.ContractCaller) (*Map3PoolWrapperCaller, error) {
	contract, err := bindMap3PoolWrapper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Map3PoolWrapperCaller{contract: contract}, nil
}

// NewMap3PoolWrapperTransactor creates a new write-only instance of Map3PoolWrapper, bound to a specific deployed contract.
func NewMap3PoolWrapperTransactor(address common.Address, transactor bind.ContractTransactor) (*Map3PoolWrapperTransactor, error) {
	contract, err := bindMap3PoolWrapper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Map3PoolWrapperTransactor{contract: contract}, nil
}

// NewMap3PoolWrapperFilterer creates a new log filterer instance of Map3PoolWrapper, bound to a specific deployed contract.
func NewMap3PoolWrapperFilterer(address common.Address, filterer bind.ContractFilterer) (*Map3PoolWrapperFilterer, error) {
	contract, err := bindMap3PoolWrapper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Map3PoolWrapperFilterer{contract: contract}, nil
}

// bindMap3PoolWrapper binds a generic wrapper to an already deployed contract.
func bindMap3PoolWrapper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Map3PoolWrapperABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Map3PoolWrapper *Map3PoolWrapperRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Map3PoolWrapper.Contract.Map3PoolWrapperCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Map3PoolWrapper *Map3PoolWrapperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Map3PoolWrapper.Contract.Map3PoolWrapperTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Map3PoolWrapper *Map3PoolWrapperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Map3PoolWrapper.Contract.Map3PoolWrapperTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Map3PoolWrapper *Map3PoolWrapperCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Map3PoolWrapper.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Map3PoolWrapper *Map3PoolWrapperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Map3PoolWrapper.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Map3PoolWrapper *Map3PoolWrapperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Map3PoolWrapper.Contract.contract.Transact(opts, method, params...)
}

// Length is a free data retrieval call binding the contract method 0x82172882.
//
// Solidity: function Length() constant returns(uint256)
func (_Map3PoolWrapper *Map3PoolWrapperCaller) Length(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Map3PoolWrapper.contract.Call(opts, out, "Length")
	return *ret0, err
}

// Length is a free data retrieval call binding the contract method 0x82172882.
//
// Solidity: function Length() constant returns(uint256)
func (_Map3PoolWrapper *Map3PoolWrapperSession) Length() (*big.Int, error) {
	return _Map3PoolWrapper.Contract.Length(&_Map3PoolWrapper.CallOpts)
}

// Length is a free data retrieval call binding the contract method 0x82172882.
//
// Solidity: function Length() constant returns(uint256)
func (_Map3PoolWrapper *Map3PoolWrapperCallerSession) Length() (*big.Int, error) {
	return _Map3PoolWrapper.Contract.Length(&_Map3PoolWrapper.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0xbb62860d.
//
// Solidity: function Version() constant returns(int256)
func (_Map3PoolWrapper *Map3PoolWrapperCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Map3PoolWrapper.contract.Call(opts, out, "Version")
	return *ret0, err
}

// Version is a free data retrieval call binding the contract method 0xbb62860d.
//
// Solidity: function Version() constant returns(int256)
func (_Map3PoolWrapper *Map3PoolWrapperSession) Version() (*big.Int, error) {
	return _Map3PoolWrapper.Contract.Version(&_Map3PoolWrapper.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0xbb62860d.
//
// Solidity: function Version() constant returns(int256)
func (_Map3PoolWrapper *Map3PoolWrapperCallerSession) Version() (*big.Int, error) {
	return _Map3PoolWrapper.Contract.Version(&_Map3PoolWrapper.CallOpts)
}
