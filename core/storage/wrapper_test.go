// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package storage

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
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// StorageWrapperABI is the input ABI used to generate the binding from.
const StorageWrapperABI = "[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"constant\":true,\"inputs\":[],\"name\":\"Hello\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"res\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"Version\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"v\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"Name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"Balance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// StorageWrapper is an auto generated Go binding around an Ethereum contract.
type StorageWrapper struct {
	StorageWrapperCaller     // Read-only binding to the contract
	StorageWrapperTransactor // Write-only binding to the contract
	StorageWrapperFilterer   // Log filterer for contract events
}

// StorageWrapperCaller is an auto generated read-only Go binding around an Ethereum contract.
type StorageWrapperCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageWrapperTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StorageWrapperTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageWrapperFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StorageWrapperFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageWrapperSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StorageWrapperSession struct {
	Contract     *StorageWrapper   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StorageWrapperCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StorageWrapperCallerSession struct {
	Contract *StorageWrapperCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// StorageWrapperTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StorageWrapperTransactorSession struct {
	Contract     *StorageWrapperTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// StorageWrapperRaw is an auto generated low-level Go binding around an Ethereum contract.
type StorageWrapperRaw struct {
	Contract *StorageWrapper // Generic contract binding to access the raw methods on
}

// StorageWrapperCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StorageWrapperCallerRaw struct {
	Contract *StorageWrapperCaller // Generic read-only contract binding to access the raw methods on
}

// StorageWrapperTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StorageWrapperTransactorRaw struct {
	Contract *StorageWrapperTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStorageWrapper creates a new instance of StorageWrapper, bound to a specific deployed contract.
func NewStorageWrapper(address common.Address, backend bind.ContractBackend) (*StorageWrapper, error) {
	contract, err := bindStorageWrapper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StorageWrapper{StorageWrapperCaller: StorageWrapperCaller{contract: contract}, StorageWrapperTransactor: StorageWrapperTransactor{contract: contract}, StorageWrapperFilterer: StorageWrapperFilterer{contract: contract}}, nil
}

// NewStorageWrapperCaller creates a new read-only instance of StorageWrapper, bound to a specific deployed contract.
func NewStorageWrapperCaller(address common.Address, caller bind.ContractCaller) (*StorageWrapperCaller, error) {
	contract, err := bindStorageWrapper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StorageWrapperCaller{contract: contract}, nil
}

// NewStorageWrapperTransactor creates a new write-only instance of StorageWrapper, bound to a specific deployed contract.
func NewStorageWrapperTransactor(address common.Address, transactor bind.ContractTransactor) (*StorageWrapperTransactor, error) {
	contract, err := bindStorageWrapper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StorageWrapperTransactor{contract: contract}, nil
}

// NewStorageWrapperFilterer creates a new log filterer instance of StorageWrapper, bound to a specific deployed contract.
func NewStorageWrapperFilterer(address common.Address, filterer bind.ContractFilterer) (*StorageWrapperFilterer, error) {
	contract, err := bindStorageWrapper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StorageWrapperFilterer{contract: contract}, nil
}

// bindStorageWrapper binds a generic wrapper to an already deployed contract.
func bindStorageWrapper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StorageWrapperABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StorageWrapper *StorageWrapperRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _StorageWrapper.Contract.StorageWrapperCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StorageWrapper *StorageWrapperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StorageWrapper.Contract.StorageWrapperTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StorageWrapper *StorageWrapperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StorageWrapper.Contract.StorageWrapperTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StorageWrapper *StorageWrapperCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _StorageWrapper.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StorageWrapper *StorageWrapperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StorageWrapper.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StorageWrapper *StorageWrapperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StorageWrapper.Contract.contract.Transact(opts, method, params...)
}

// Balance is a free data retrieval call binding the contract method 0x0ef67887.
//
// Solidity: function Balance() view returns(uint256)
func (_StorageWrapper *StorageWrapperCaller) Balance(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _StorageWrapper.contract.Call(opts, out, "Balance")
	return *ret0, err
}

// Balance is a free data retrieval call binding the contract method 0x0ef67887.
//
// Solidity: function Balance() view returns(uint256)
func (_StorageWrapper *StorageWrapperSession) Balance() (*big.Int, error) {
	return _StorageWrapper.Contract.Balance(&_StorageWrapper.CallOpts)
}

// Balance is a free data retrieval call binding the contract method 0x0ef67887.
//
// Solidity: function Balance() view returns(uint256)
func (_StorageWrapper *StorageWrapperCallerSession) Balance() (*big.Int, error) {
	return _StorageWrapper.Contract.Balance(&_StorageWrapper.CallOpts)
}

// Hello is a free data retrieval call binding the contract method 0xbcdfe0d5.
//
// Solidity: function Hello() pure returns(string res)
func (_StorageWrapper *StorageWrapperCaller) Hello(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _StorageWrapper.contract.Call(opts, out, "Hello")
	return *ret0, err
}

// Hello is a free data retrieval call binding the contract method 0xbcdfe0d5.
//
// Solidity: function Hello() pure returns(string res)
func (_StorageWrapper *StorageWrapperSession) Hello() (string, error) {
	return _StorageWrapper.Contract.Hello(&_StorageWrapper.CallOpts)
}

// Hello is a free data retrieval call binding the contract method 0xbcdfe0d5.
//
// Solidity: function Hello() pure returns(string res)
func (_StorageWrapper *StorageWrapperCallerSession) Hello() (string, error) {
	return _StorageWrapper.Contract.Hello(&_StorageWrapper.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x8052474d.
//
// Solidity: function Name() view returns(string)
func (_StorageWrapper *StorageWrapperCaller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _StorageWrapper.contract.Call(opts, out, "Name")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x8052474d.
//
// Solidity: function Name() view returns(string)
func (_StorageWrapper *StorageWrapperSession) Name() (string, error) {
	return _StorageWrapper.Contract.Name(&_StorageWrapper.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x8052474d.
//
// Solidity: function Name() view returns(string)
func (_StorageWrapper *StorageWrapperCallerSession) Name() (string, error) {
	return _StorageWrapper.Contract.Name(&_StorageWrapper.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0xbb62860d.
//
// Solidity: function Version() view returns(int256 v)
func (_StorageWrapper *StorageWrapperCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _StorageWrapper.contract.Call(opts, out, "Version")
	return *ret0, err
}

// Version is a free data retrieval call binding the contract method 0xbb62860d.
//
// Solidity: function Version() view returns(int256 v)
func (_StorageWrapper *StorageWrapperSession) Version() (*big.Int, error) {
	return _StorageWrapper.Contract.Version(&_StorageWrapper.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0xbb62860d.
//
// Solidity: function Version() view returns(int256 v)
func (_StorageWrapper *StorageWrapperCallerSession) Version() (*big.Int, error) {
	return _StorageWrapper.Contract.Version(&_StorageWrapper.CallOpts)
}
