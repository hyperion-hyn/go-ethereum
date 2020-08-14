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
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ConsortiumABI is the input ABI used to generate the binding from.
const ConsortiumABI = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// Consortium is an auto generated Go binding around an Ethereum contract.
type Consortium struct {
	ConsortiumCaller     // Read-only binding to the contract
	ConsortiumTransactor // Write-only binding to the contract
	ConsortiumFilterer   // Log filterer for contract events
}

// ConsortiumCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConsortiumCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConsortiumTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConsortiumTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConsortiumFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConsortiumFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConsortiumSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConsortiumSession struct {
	Contract     *Consortium       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ConsortiumCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConsortiumCallerSession struct {
	Contract *ConsortiumCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// ConsortiumTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConsortiumTransactorSession struct {
	Contract     *ConsortiumTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ConsortiumRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConsortiumRaw struct {
	Contract *Consortium // Generic contract binding to access the raw methods on
}

// ConsortiumCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConsortiumCallerRaw struct {
	Contract *ConsortiumCaller // Generic read-only contract binding to access the raw methods on
}

// ConsortiumTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConsortiumTransactorRaw struct {
	Contract *ConsortiumTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConsortium creates a new instance of Consortium, bound to a specific deployed contract.
func NewConsortium(address common.Address, backend bind.ContractBackend) (*Consortium, error) {
	contract, err := bindConsortium(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Consortium{ConsortiumCaller: ConsortiumCaller{contract: contract}, ConsortiumTransactor: ConsortiumTransactor{contract: contract}, ConsortiumFilterer: ConsortiumFilterer{contract: contract}}, nil
}

// NewConsortiumCaller creates a new read-only instance of Consortium, bound to a specific deployed contract.
func NewConsortiumCaller(address common.Address, caller bind.ContractCaller) (*ConsortiumCaller, error) {
	contract, err := bindConsortium(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConsortiumCaller{contract: contract}, nil
}

// NewConsortiumTransactor creates a new write-only instance of Consortium, bound to a specific deployed contract.
func NewConsortiumTransactor(address common.Address, transactor bind.ContractTransactor) (*ConsortiumTransactor, error) {
	contract, err := bindConsortium(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConsortiumTransactor{contract: contract}, nil
}

// NewConsortiumFilterer creates a new log filterer instance of Consortium, bound to a specific deployed contract.
func NewConsortiumFilterer(address common.Address, filterer bind.ContractFilterer) (*ConsortiumFilterer, error) {
	contract, err := bindConsortium(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConsortiumFilterer{contract: contract}, nil
}

// bindConsortium binds a generic wrapper to an already deployed contract.
func bindConsortium(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ConsortiumABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Consortium *ConsortiumRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Consortium.Contract.ConsortiumCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Consortium *ConsortiumRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Consortium.Contract.ConsortiumTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Consortium *ConsortiumRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Consortium.Contract.ConsortiumTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Consortium *ConsortiumCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Consortium.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Consortium *ConsortiumTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Consortium.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Consortium *ConsortiumTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Consortium.Contract.contract.Transact(opts, method, params...)
}
