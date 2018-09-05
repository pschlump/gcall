// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Greeter

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

// GreeterABI is the input ABI used to generate the binding from.
const GreeterABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_amt\",\"type\":\"uint256\"}],\"name\":\"withdrawAmount\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"xgreeting\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"greeter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_greeting\",\"type\":\"uint256\"}],\"name\":\"setNGreeting\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"kill\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getNGreeting\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"aNum\",\"type\":\"uint256\"}],\"name\":\"test01\",\"outputs\":[{\"name\":\"id\",\"type\":\"uint256\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ngreeting\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_greeting\",\"type\":\"string\"}],\"name\":\"setGreeting\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"greet\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_greeting\",\"type\":\"string\"}],\"name\":\"greeter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"getGreeting\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"greeting\",\"type\":\"string\"}],\"name\":\"ReportGreetingEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"greeting\",\"type\":\"string\"}],\"name\":\"ReportGreetingChangedEvent\",\"type\":\"event\"}]"

// Greeter is an auto generated Go binding around an Ethereum contract.
type Greeter struct {
	GreeterCaller     // Read-only binding to the contract
	GreeterTransactor // Write-only binding to the contract
	GreeterFilterer   // Log filterer for contract events
}

// GreeterCaller is an auto generated read-only Go binding around an Ethereum contract.
type GreeterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GreeterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GreeterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GreeterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GreeterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GreeterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GreeterSession struct {
	Contract     *Greeter          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GreeterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GreeterCallerSession struct {
	Contract *GreeterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// GreeterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GreeterTransactorSession struct {
	Contract     *GreeterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// GreeterRaw is an auto generated low-level Go binding around an Ethereum contract.
type GreeterRaw struct {
	Contract *Greeter // Generic contract binding to access the raw methods on
}

// GreeterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GreeterCallerRaw struct {
	Contract *GreeterCaller // Generic read-only contract binding to access the raw methods on
}

// GreeterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GreeterTransactorRaw struct {
	Contract *GreeterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGreeter creates a new instance of Greeter, bound to a specific deployed contract.
func NewGreeter(address common.Address, backend bind.ContractBackend) (*Greeter, error) {
	contract, err := bindGreeter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Greeter{GreeterCaller: GreeterCaller{contract: contract}, GreeterTransactor: GreeterTransactor{contract: contract}, GreeterFilterer: GreeterFilterer{contract: contract}}, nil
}

// NewGreeterCaller creates a new read-only instance of Greeter, bound to a specific deployed contract.
func NewGreeterCaller(address common.Address, caller bind.ContractCaller) (*GreeterCaller, error) {
	contract, err := bindGreeter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GreeterCaller{contract: contract}, nil
}

// NewGreeterTransactor creates a new write-only instance of Greeter, bound to a specific deployed contract.
func NewGreeterTransactor(address common.Address, transactor bind.ContractTransactor) (*GreeterTransactor, error) {
	contract, err := bindGreeter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GreeterTransactor{contract: contract}, nil
}

// NewGreeterFilterer creates a new log filterer instance of Greeter, bound to a specific deployed contract.
func NewGreeterFilterer(address common.Address, filterer bind.ContractFilterer) (*GreeterFilterer, error) {
	contract, err := bindGreeter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GreeterFilterer{contract: contract}, nil
}

// bindGreeter binds a generic wrapper to an already deployed contract.
func bindGreeter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(GreeterABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Greeter *GreeterRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Greeter.Contract.GreeterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Greeter *GreeterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Greeter.Contract.GreeterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Greeter *GreeterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Greeter.Contract.GreeterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Greeter *GreeterCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Greeter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Greeter *GreeterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Greeter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Greeter *GreeterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Greeter.Contract.contract.Transact(opts, method, params...)
}

// GetNGreeting is a free data retrieval call binding the contract method 0x43ae9581.
//
// Solidity: function getNGreeting() constant returns(uint256)
func (_Greeter *GreeterCaller) GetNGreeting(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Greeter.contract.Call(opts, out, "getNGreeting")
	return *ret0, err
}

// GetNGreeting is a free data retrieval call binding the contract method 0x43ae9581.
//
// Solidity: function getNGreeting() constant returns(uint256)
func (_Greeter *GreeterSession) GetNGreeting() (*big.Int, error) {
	return _Greeter.Contract.GetNGreeting(&_Greeter.CallOpts)
}

// GetNGreeting is a free data retrieval call binding the contract method 0x43ae9581.
//
// Solidity: function getNGreeting() constant returns(uint256)
func (_Greeter *GreeterCallerSession) GetNGreeting() (*big.Int, error) {
	return _Greeter.Contract.GetNGreeting(&_Greeter.CallOpts)
}

// Greet is a free data retrieval call binding the contract method 0xcfae3217.
//
// Solidity: function greet() constant returns(string)
func (_Greeter *GreeterCaller) Greet(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Greeter.contract.Call(opts, out, "greet")
	return *ret0, err
}

// Greet is a free data retrieval call binding the contract method 0xcfae3217.
//
// Solidity: function greet() constant returns(string)
func (_Greeter *GreeterSession) Greet() (string, error) {
	return _Greeter.Contract.Greet(&_Greeter.CallOpts)
}

// Greet is a free data retrieval call binding the contract method 0xcfae3217.
//
// Solidity: function greet() constant returns(string)
func (_Greeter *GreeterCallerSession) Greet() (string, error) {
	return _Greeter.Contract.Greet(&_Greeter.CallOpts)
}

// Ngreeting is a free data retrieval call binding the contract method 0x8e554e98.
//
// Solidity: function ngreeting() constant returns(uint256)
func (_Greeter *GreeterCaller) Ngreeting(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Greeter.contract.Call(opts, out, "ngreeting")
	return *ret0, err
}

// Ngreeting is a free data retrieval call binding the contract method 0x8e554e98.
//
// Solidity: function ngreeting() constant returns(uint256)
func (_Greeter *GreeterSession) Ngreeting() (*big.Int, error) {
	return _Greeter.Contract.Ngreeting(&_Greeter.CallOpts)
}

// Ngreeting is a free data retrieval call binding the contract method 0x8e554e98.
//
// Solidity: function ngreeting() constant returns(uint256)
func (_Greeter *GreeterCallerSession) Ngreeting() (*big.Int, error) {
	return _Greeter.Contract.Ngreeting(&_Greeter.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Greeter *GreeterCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Greeter.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Greeter *GreeterSession) Owner() (common.Address, error) {
	return _Greeter.Contract.Owner(&_Greeter.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Greeter *GreeterCallerSession) Owner() (common.Address, error) {
	return _Greeter.Contract.Owner(&_Greeter.CallOpts)
}

// Xgreeting is a free data retrieval call binding the contract method 0x0753fea8.
//
// Solidity: function xgreeting() constant returns(string)
func (_Greeter *GreeterCaller) Xgreeting(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Greeter.contract.Call(opts, out, "xgreeting")
	return *ret0, err
}

// Xgreeting is a free data retrieval call binding the contract method 0x0753fea8.
//
// Solidity: function xgreeting() constant returns(string)
func (_Greeter *GreeterSession) Xgreeting() (string, error) {
	return _Greeter.Contract.Xgreeting(&_Greeter.CallOpts)
}

// Xgreeting is a free data retrieval call binding the contract method 0x0753fea8.
//
// Solidity: function xgreeting() constant returns(string)
func (_Greeter *GreeterCallerSession) Xgreeting() (string, error) {
	return _Greeter.Contract.Xgreeting(&_Greeter.CallOpts)
}

// GetGreeting is a paid mutator transaction binding the contract method 0xfe50cc72.
//
// Solidity: function getGreeting() returns(string)
func (_Greeter *GreeterTransactor) GetGreeting(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Greeter.contract.Transact(opts, "getGreeting")
}

// GetGreeting is a paid mutator transaction binding the contract method 0xfe50cc72.
//
// Solidity: function getGreeting() returns(string)
func (_Greeter *GreeterSession) GetGreeting() (*types.Transaction, error) {
	return _Greeter.Contract.GetGreeting(&_Greeter.TransactOpts)
}

// GetGreeting is a paid mutator transaction binding the contract method 0xfe50cc72.
//
// Solidity: function getGreeting() returns(string)
func (_Greeter *GreeterTransactorSession) GetGreeting() (*types.Transaction, error) {
	return _Greeter.Contract.GetGreeting(&_Greeter.TransactOpts)
}

// Greeter is a paid mutator transaction binding the contract method 0xfaf27bca.
//
// Solidity: function greeter(_greeting string) returns()
func (_Greeter *GreeterTransactor) Greeter(opts *bind.TransactOpts, _greeting string) (*types.Transaction, error) {
	return _Greeter.contract.Transact(opts, "greeter", _greeting)
}

// Greeter is a paid mutator transaction binding the contract method 0xfaf27bca.
//
// Solidity: function greeter(_greeting string) returns()
func (_Greeter *GreeterSession) Greeter(_greeting string) (*types.Transaction, error) {
	return _Greeter.Contract.Greeter(&_Greeter.TransactOpts, _greeting)
}

// Greeter is a paid mutator transaction binding the contract method 0xfaf27bca.
//
// Solidity: function greeter(_greeting string) returns()
func (_Greeter *GreeterTransactorSession) Greeter(_greeting string) (*types.Transaction, error) {
	return _Greeter.Contract.Greeter(&_Greeter.TransactOpts, _greeting)
}

// Kill is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: function kill() returns()
func (_Greeter *GreeterTransactor) Kill(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Greeter.contract.Transact(opts, "kill")
}

// Kill is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: function kill() returns()
func (_Greeter *GreeterSession) Kill() (*types.Transaction, error) {
	return _Greeter.Contract.Kill(&_Greeter.TransactOpts)
}

// Kill is a paid mutator transaction binding the contract method 0x41c0e1b5.
//
// Solidity: function kill() returns()
func (_Greeter *GreeterTransactorSession) Kill() (*types.Transaction, error) {
	return _Greeter.Contract.Kill(&_Greeter.TransactOpts)
}

// SetGreeting is a paid mutator transaction binding the contract method 0xa4136862.
//
// Solidity: function setGreeting(_greeting string) returns()
func (_Greeter *GreeterTransactor) SetGreeting(opts *bind.TransactOpts, _greeting string) (*types.Transaction, error) {
	return _Greeter.contract.Transact(opts, "setGreeting", _greeting)
}

// SetGreeting is a paid mutator transaction binding the contract method 0xa4136862.
//
// Solidity: function setGreeting(_greeting string) returns()
func (_Greeter *GreeterSession) SetGreeting(_greeting string) (*types.Transaction, error) {
	return _Greeter.Contract.SetGreeting(&_Greeter.TransactOpts, _greeting)
}

// SetGreeting is a paid mutator transaction binding the contract method 0xa4136862.
//
// Solidity: function setGreeting(_greeting string) returns()
func (_Greeter *GreeterTransactorSession) SetGreeting(_greeting string) (*types.Transaction, error) {
	return _Greeter.Contract.SetGreeting(&_Greeter.TransactOpts, _greeting)
}

// SetNGreeting is a paid mutator transaction binding the contract method 0x2df045ff.
//
// Solidity: function setNGreeting(_greeting uint256) returns()
func (_Greeter *GreeterTransactor) SetNGreeting(opts *bind.TransactOpts, _greeting *big.Int) (*types.Transaction, error) {
	return _Greeter.contract.Transact(opts, "setNGreeting", _greeting)
}

// SetNGreeting is a paid mutator transaction binding the contract method 0x2df045ff.
//
// Solidity: function setNGreeting(_greeting uint256) returns()
func (_Greeter *GreeterSession) SetNGreeting(_greeting *big.Int) (*types.Transaction, error) {
	return _Greeter.Contract.SetNGreeting(&_Greeter.TransactOpts, _greeting)
}

// SetNGreeting is a paid mutator transaction binding the contract method 0x2df045ff.
//
// Solidity: function setNGreeting(_greeting uint256) returns()
func (_Greeter *GreeterTransactorSession) SetNGreeting(_greeting *big.Int) (*types.Transaction, error) {
	return _Greeter.Contract.SetNGreeting(&_Greeter.TransactOpts, _greeting)
}

// Test01 is a paid mutator transaction binding the contract method 0x7cc22ea9.
//
// Solidity: function test01(aNum uint256) returns(id uint256)
func (_Greeter *GreeterTransactor) Test01(opts *bind.TransactOpts, aNum *big.Int) (*types.Transaction, error) {
	return _Greeter.contract.Transact(opts, "test01", aNum)
}

// Test01 is a paid mutator transaction binding the contract method 0x7cc22ea9.
//
// Solidity: function test01(aNum uint256) returns(id uint256)
func (_Greeter *GreeterSession) Test01(aNum *big.Int) (*types.Transaction, error) {
	return _Greeter.Contract.Test01(&_Greeter.TransactOpts, aNum)
}

// Test01 is a paid mutator transaction binding the contract method 0x7cc22ea9.
//
// Solidity: function test01(aNum uint256) returns(id uint256)
func (_Greeter *GreeterTransactorSession) Test01(aNum *big.Int) (*types.Transaction, error) {
	return _Greeter.Contract.Test01(&_Greeter.TransactOpts, aNum)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_Greeter *GreeterTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Greeter.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_Greeter *GreeterSession) Withdraw() (*types.Transaction, error) {
	return _Greeter.Contract.Withdraw(&_Greeter.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_Greeter *GreeterTransactorSession) Withdraw() (*types.Transaction, error) {
	return _Greeter.Contract.Withdraw(&_Greeter.TransactOpts)
}

// WithdrawAmount is a paid mutator transaction binding the contract method 0x0562b9f7.
//
// Solidity: function withdrawAmount(_amt uint256) returns()
func (_Greeter *GreeterTransactor) WithdrawAmount(opts *bind.TransactOpts, _amt *big.Int) (*types.Transaction, error) {
	return _Greeter.contract.Transact(opts, "withdrawAmount", _amt)
}

// WithdrawAmount is a paid mutator transaction binding the contract method 0x0562b9f7.
//
// Solidity: function withdrawAmount(_amt uint256) returns()
func (_Greeter *GreeterSession) WithdrawAmount(_amt *big.Int) (*types.Transaction, error) {
	return _Greeter.Contract.WithdrawAmount(&_Greeter.TransactOpts, _amt)
}

// WithdrawAmount is a paid mutator transaction binding the contract method 0x0562b9f7.
//
// Solidity: function withdrawAmount(_amt uint256) returns()
func (_Greeter *GreeterTransactorSession) WithdrawAmount(_amt *big.Int) (*types.Transaction, error) {
	return _Greeter.Contract.WithdrawAmount(&_Greeter.TransactOpts, _amt)
}

// GreeterReportGreetingChangedEventIterator is returned from FilterReportGreetingChangedEvent and is used to iterate over the raw logs and unpacked data for ReportGreetingChangedEvent events raised by the Greeter contract.
type GreeterReportGreetingChangedEventIterator struct {
	Event *GreeterReportGreetingChangedEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GreeterReportGreetingChangedEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GreeterReportGreetingChangedEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GreeterReportGreetingChangedEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GreeterReportGreetingChangedEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GreeterReportGreetingChangedEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GreeterReportGreetingChangedEvent represents a ReportGreetingChangedEvent event raised by the Greeter contract.
type GreeterReportGreetingChangedEvent struct {
	Greeting string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterReportGreetingChangedEvent is a free log retrieval operation binding the contract event 0xfc2d6907db7ce1823e4bae783aeda0363c138aba19e7a497ebc3dd80a0a3453a.
//
// Solidity: event ReportGreetingChangedEvent(greeting string)
func (_Greeter *GreeterFilterer) FilterReportGreetingChangedEvent(opts *bind.FilterOpts) (*GreeterReportGreetingChangedEventIterator, error) {

	logs, sub, err := _Greeter.contract.FilterLogs(opts, "ReportGreetingChangedEvent")
	if err != nil {
		return nil, err
	}
	return &GreeterReportGreetingChangedEventIterator{contract: _Greeter.contract, event: "ReportGreetingChangedEvent", logs: logs, sub: sub}, nil
}

// WatchReportGreetingChangedEvent is a free log subscription operation binding the contract event 0xfc2d6907db7ce1823e4bae783aeda0363c138aba19e7a497ebc3dd80a0a3453a.
//
// Solidity: event ReportGreetingChangedEvent(greeting string)
func (_Greeter *GreeterFilterer) WatchReportGreetingChangedEvent(opts *bind.WatchOpts, sink chan<- *GreeterReportGreetingChangedEvent) (event.Subscription, error) {

	logs, sub, err := _Greeter.contract.WatchLogs(opts, "ReportGreetingChangedEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GreeterReportGreetingChangedEvent)
				if err := _Greeter.contract.UnpackLog(event, "ReportGreetingChangedEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// GreeterReportGreetingEventIterator is returned from FilterReportGreetingEvent and is used to iterate over the raw logs and unpacked data for ReportGreetingEvent events raised by the Greeter contract.
type GreeterReportGreetingEventIterator struct {
	Event *GreeterReportGreetingEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GreeterReportGreetingEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GreeterReportGreetingEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GreeterReportGreetingEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GreeterReportGreetingEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GreeterReportGreetingEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GreeterReportGreetingEvent represents a ReportGreetingEvent event raised by the Greeter contract.
type GreeterReportGreetingEvent struct {
	Greeting string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterReportGreetingEvent is a free log retrieval operation binding the contract event 0xcd8a940e2f36da995f021b62a2ae60d67db78ea2004f4c3aec74c8458db73e4b.
//
// Solidity: event ReportGreetingEvent(greeting string)
func (_Greeter *GreeterFilterer) FilterReportGreetingEvent(opts *bind.FilterOpts) (*GreeterReportGreetingEventIterator, error) {

	logs, sub, err := _Greeter.contract.FilterLogs(opts, "ReportGreetingEvent")
	if err != nil {
		return nil, err
	}
	return &GreeterReportGreetingEventIterator{contract: _Greeter.contract, event: "ReportGreetingEvent", logs: logs, sub: sub}, nil
}

// WatchReportGreetingEvent is a free log subscription operation binding the contract event 0xcd8a940e2f36da995f021b62a2ae60d67db78ea2004f4c3aec74c8458db73e4b.
//
// Solidity: event ReportGreetingEvent(greeting string)
func (_Greeter *GreeterFilterer) WatchReportGreetingEvent(opts *bind.WatchOpts, sink chan<- *GreeterReportGreetingEvent) (event.Subscription, error) {

	logs, sub, err := _Greeter.contract.WatchLogs(opts, "ReportGreetingEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GreeterReportGreetingEvent)
				if err := _Greeter.contract.UnpackLog(event, "ReportGreetingEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
