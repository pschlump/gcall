package main

/*
MIT License

Copyright (c) 2018 Philip Schlump

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"      //
	"github.com/ethereum/go-ethereum/accounts/abi/bind" //
	"github.com/ethereum/go-ethereum/common"            //
	"github.com/ethereum/go-ethereum/core/types"        //
	"github.com/ethereum/go-ethereum/crypto/sha3"       //
	"github.com/ethereum/go-ethereum/ethclient"         //
	"github.com/pschlump/GCall/bytecode"                //
	"github.com/pschlump/MiscLib"                       //
	"github.com/pschlump/ethrpc"                        // OLD: "github.com/onrik/ethrpc" - modified with new functions and functionality
	"github.com/pschlump/godebug"                       //
)

// -----------------------------------------------------------------------------------------------------
type CurrentWatchType struct {
	ContractName string
	EventName    string
}

var CurrentWatch []CurrentWatchType

// -----------------------------------------------------------------------------------------------------

type ABI_IO_Type struct {
	Name    string //
	Type    string //
	Indexed bool   //
}

type ABIType struct {
	Constant        bool          //
	IsConstructor   bool          //
	Name            string        // if name matches with contract-name, then constructor
	Inputs          []ABI_IO_Type //
	Outputs         []ABI_IO_Type //
	Payable         bool          //
	StateMutability string        // payable, nonpayable, view
	Anonymous       bool          //
	Type            string        // function, event
	MergeFrom       string
}

type ContractInfo struct {
	Name            string         // Name of this contract
	Address         string         // Address it is at - most recent version
	Version         string         // version string (semantic) pulled from source //$version$: v1.0.0
	LoadDateTime    string         // Pulled fomm block where contract loded
	FromAddress     string         // If different from global - overide owner - Account to unlock - owner of contracts
	KeyFile         string         // If different from global
	KeyFilePassword string         // If different from global
	ABI             []ABIType      // JSON parsed ABI
	RawABI          string         // Raw String version of ABI
	address         common.Address //
}

type GethInfoNetwork struct {
	GethURL_ws          string // ws://192.168.0.139:8546
	GethURL_http        string // http://192.168.0.139:8545
	FromAddress         string // Account to unlock - owner of contracts
	FromAddressPassword string //
	KeyFile             string //
	KeyFilePassword     string //
}

type GethInfo struct {
	GethURL_ws          string                     // ws://192.168.0.139:8546
	GethURL_http        string                     // http://192.168.0.139:8545
	ABIPath             []string                   //
	SRCPath             []string                   //
	FromAddress         string                     // Account to unlock - owner of contracts
	FromAddressPassword string                     //
	KeyFile             string                     //
	KeyFilePassword     string                     //
	ContractList        map[string]ContractInfo    // Map of Contracts
	ContractNames       []string                   // Orderd list of names - sorted
	DebugFlags          []string                   // List of debug flags - set by default
	UnlockSeconds       int                        //
	SolcCMD             string                     // xyzzy082
	AbigenCMD           string                     // xyzzy082
	RebuildBinary       []string                   //  xyzzy082 make or make abi, something like that
	BinariesInDir       string                     //  xyzzy082 What directory to look for .bin files
	ReadlineDir         string                     // xyzzy082
	AutoCurrentCheck    string                     //  xyzzy082 if "yes" then check all contrcts with address for current source
	PushAddressChanges  string                     //  xyzzy082 Sh Script To: ( cd address_dir ; git commit -m "Address Change" . )
	PullAddressChanges  string                     //  xyzzy082 Sh Script To: ( cd address_dir ; git pull )
	NetworkFlag         map[string]GethInfoNetwork //
	conn                *ethclient.Client          // geth connection
	rpc_client          *ethrpc.EthRPC             //
	CallOpts            *bind.CallOpts             // Call options to use throughout this session
	TransactOpts        *bind.TransactOpts         // Transaction auth options to use throughout this session
}

func listContracts() func(string) []string {
	return func(line string) []string {
		names := make([]string, 0)
		for key := range ContractAddressHash {
			names = append(names, key)
		}
		return names
	}
}

func listContractsMethods() func(string) []string {
	return func(line string) []string {
		names := make([]string, 0)
		for key, cl := range gCfg.ContractList {
			for _, aMethod := range cl.ABI {
				if aMethod.Name != "" && aMethod.Type == "function" {
					names = append(names, key+"."+aMethod.Name)
				}
			}
		}
		return names
	}
}

func listContractsEvents() func(string) []string {
	return func(line string) []string {
		names := make([]string, 0)
		for key, cl := range gCfg.ContractList {
			hasEvent := false
			for _, aMethod := range cl.ABI {
				if aMethod.Name != "" && aMethod.Type == "event" {
					hasEvent = true
					break
				}
			}
			if hasEvent {
				names = append(names, key)
				for _, aMethod := range cl.ABI {
					if aMethod.Name != "" && aMethod.Type == "event" {
						names = append(names, key+"."+aMethod.Name)
					}
				}
			}
		}
		return names
	}
}

// -------------------------------------------------------------------------------------------------
func (gCfg *GethInfo) SetTransactOpts() (err error) {
	// Account with "gas" to spend so that you can call the contract.
	file, err := os.Open(gCfg.KeyFile) // Read in file
	if err != nil {
		err = fmt.Errorf("Failed to open keyfile: %v, %s, %s", err, gCfg.KeyFile, godebug.LF())
		return
	}

	// Create a new trasaction call (will mine) for the contract.
	opts, err := bind.NewTransactor(bufio.NewReader(file), gCfg.KeyFilePassword)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sFailed to open keyfile: %v, %s, %s\n    This may be due to an invalid password.\n%s", MiscLib.ColorRed, err, gCfg.KeyFile, godebug.LF(), MiscLib.ColorReset)
		err = fmt.Errorf("Failed to read keyfile: %v, %s, %s", err, gCfg.KeyFile, godebug.LF())
		return
	}

	requestFrom := common.HexToAddress(gCfg.FromAddress)

	// function requestRelay(uint256 _payment, uint256 _blockReward, uint256 _seed)
	// payment needs to be sent as "msg.value" - as a part of the standard message
	opts.From = requestFrom

	opts.Value = big.NewInt(0) // Default payment - set to 0 wei.

	opts.GasPrice = big.NewInt(100000000000) // *big.Int // Gas price to use for the transaction execution (nil = gas price oracle)
	// opts.GasPrice = big.NewInt(0) // *big.Int // Gas price to use for the transaction execution (nil = gas price oracle)
	opts.GasLimit = 4712388 // uint64   // Gas limit to set for the transaction execution (0 = estimate)

	fmt.Printf("Value=%x GasLimit=%x\n", opts.Value, opts.GasLimit)

	gCfg.TransactOpts = opts
	return
}

// -------------------------------------------------------------------------------------------------
// GetContractAddress looks up the address for the contract by name.
func (gCfg *GethInfo) GetContractAddress(contractName string) (address common.Address, err error) {
	if ct, ok := gCfg.ContractList[contractName]; ok {
		return ct.address, nil
	}
	err = fmt.Errorf("Invalid contract name [%s]\n", contractName)
	return
}

// -------------------------------------------------------------------------------------------------
//	if !IsValidMethod(gCfg,methodName) {
// ContractList        map[string]ContractInfo // Map of Contracts
func (gCfg *GethInfo) IsValidMethodName(contractName, methodName string) (pos int, rv bool) {
	pos = -1
	cc := gCfg.ContractList[contractName]
	for ii, abi := range cc.ABI {
		if abi.Name == methodName {
			return ii, true
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------
// ABIMethodSet = gCfg.IsValidMethodNameSet(contractName, methodName, nParam)
func (gCfg *GethInfo) IsValidMethodNameSet(contractName, methodName string, need int) (pos []int, nfound int) {
	nfound = 0
	cc := gCfg.ContractList[contractName]
	for ii, abi := range cc.ABI {
		if abi.Name == methodName {
			if need == len(abi.Inputs) {
				nfound++
				pos = append(pos, ii)
			}
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------
// Bind2Contract binds a generic wrapper to an already deployed contract.
func Bind2Contract(ABI string, address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, *abi.ABI, error) {
	parsed, err := abi.JSON(strings.NewReader(ABI))
	if err != nil {
		return nil, nil, err
	}
	godebug.DbPf(gDebug["db12"], "Type of parsed = %T, value %s, %s\n", parsed, godebug.SVarI(parsed), godebug.LF())
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), &parsed, nil
}

// -------------------------------------------------------------------------------------------------
type ContractMgr struct {
	Contract *bind.BoundContract // Generic contract wrapper for the low level calls
	GCfg     *GethInfo
}

// -------------------------------------------------------------------------------------------------
// func (gCfg *GethInfo) IsValidMethodName(contractName, methodName string) (pos int, rv bool) {
func NewContractMgr(c *bind.BoundContract, gCfg *GethInfo) (rv *ContractMgr) {
	return &ContractMgr{
		Contract: c,
		GCfg:     gCfg,
	}
}

// -------------------------------------------------------------------------------------------------------
func (ctm *ContractMgr) CallContract(ABI ABIType, contractName, methodName string, params ...interface{}) (result interface{}, vv *types.Transaction, err error) {

	// fmt.Printf("CallContract params=%s\n", godebug.SVar(params))

	pos, found := ctm.GCfg.IsValidMethodName(contractName, methodName) // xyzzy - need to check for overloaded functions!
	if !found {
		err = fmt.Errorf("Error: %s.%s not found\n", contractName, methodName)
		return
	}
	godebug.DbPf(gDebug["db01"], "pos=%d, len of .ABI=%d\n", pos, len(ctm.GCfg.ContractList[contractName].ABI))
	// fmt.Printf("%sAT: %s After = %d%s\n", MiscLib.ColorYellow, godebug.LF(), len(ctm.GCfg.ContractList[contractName].ABI), MiscLib.ColorReset)
	abi := ctm.GCfg.ContractList[contractName].ABI[pos]
	isConst := abi.Constant
	godebug.DbPf(gDebug["db01"], "%sAT: %s, %s isConst=%v\n", MiscLib.ColorYellow, godebug.LF(), MiscLib.ColorReset, isConst)
	godebug.DbPf(gDebug["db01"], "isConst: %v\n", isConst)
	if len(abi.Inputs) != len(params) {
		fmt.Printf("Error: have %d arguments passed to function needing %d params, %s\n", len(params), len(abi.Inputs), godebug.LF())
		err = fmt.Errorf("Error: have %d arguments passed to function needing %d params\n", len(params), len(abi.Inputs))
		return
	}

	if !isConst {

		godebug.DbPf(gDebug["db01"], "%sAT: %s --- Doing a Transaction -- %s\n", MiscLib.ColorGreen, godebug.LF(), MiscLib.ColorReset)

		godebug.DbPf(gDebug["db1001"], "%sAT: %s --- Doing a Transaction -- call:%s with:%s %s\n", MiscLib.ColorGreen, godebug.LF(), methodName,
			godebug.SVar(params), MiscLib.ColorReset)

		vv, err = ctm.Transact(ctm.GCfg.TransactOpts, methodName, params...) // var vv *types.Transaction
		if err != nil {
			// fmt.Printf("...TransactOpts = %s\n", godebug.SVar(ctm.GCfg.TransactOpts))
			/*
				- Cause on call to function that takes an "address"
				   xyzzy000
				   Error on Contract call to CorpRegToken: abi: cannot use invalid as type array as argument, File: /Users/corwin/go/src/github.com/pschlump/GCall/gcall.go LineNo:1279
				   Error: abi: cannot use invalid as type array as argument on call to CorpRegToken.transfer params [6048de3601a6c4043deea717d95deac093763e6d 5000]
			*/
			fmt.Printf("Error on Contract call to %s: %s, %s\n", contractName, err, godebug.LF())
			return
		}

		godebug.DbPf(gDebug["db05"], "%sTransact: typeof(vv) = %T, vv = %s, %s%s\n", MiscLib.ColorGreen, vv, godebug.SVarI(vv), godebug.LF(), MiscLib.ColorReset)
		fmt.Printf("%sTx: %s%s\n", MiscLib.ColorGreen, godebug.SVarI(vv), MiscLib.ColorReset)
		// fmt.Printf("type %T\n", res)

		// xyzzy501

	} else {

		godebug.DbPf(gDebug["db04"], "%sAT: %s --- Doing Constant Call -- %s\n", MiscLib.ColorGreen, godebug.LF(), MiscLib.ColorReset)
		godebug.DbPf(gDebug["db09"], "ABI: %s, %s\n", godebug.SVarI(ABI), godebug.LF())

		// xyzzy - should this be "abi.Outputs[0].Type" ??? -- See above!  // xyzzy - This only handles the case of 1 return value - error if more than one.
		godebug.DbPf(gDebug["show-return-type"], "Type: %s <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< \n", ABI.Outputs[0].Type)
		switch ABI.Outputs[0].Type {
		case "string":
			result = new(string)
		case "uint256":
			result = new(*big.Int)
		case "int256":
			result = new(*big.Int)
		case "address":
			result = new(common.Address)
		case "int8":
			result = new(int8)
		case "int16":
			result = new(int16)
		case "int24", "int32":
			result = new(int32)
		case "int64", "int40", "int48", "int56":
			result = new(int64)
		case "int72", "int80", "int88", "int96", "int104", "int112", "int120", "int128":
			result = new(*big.Int)
		case "int":
			result = new(*big.Int)
		case "uint8":
			result = new(uint8)
		case "uint16":
			result = new(uint16)
		case "uint24", "uint32":
			result = new(uint32)
		case "uint64", "uint40", "uint48", "uint56":
			result = new(uint64)
		case "uint72", "uint80", "uint88", "uint96", "uint104", "uint112", "uint120", "uint128":
			result = new(*big.Int)
		case "uint":
			result = new(*big.Int)
		case "bool":
			result = new(bool)
		case "bytes32":
			// fmt.Printf("%s ********* New/Bad Type (Probably Fatal): %s AT: %s%s\n", MiscLib.ColorRed, ABI.Outputs[0].Type, godebug.LF(), MiscLib.ColorReset)
			result = new([32]byte)
		default:
			fmt.Printf("%sBad Type (Will Be Fatal): ---->>>>%s<<<<---- AT: %s%s\n", MiscLib.ColorRed, ABI.Outputs[0].Type, godebug.LF(), MiscLib.ColorReset)
			result = new(*interface{})
		}
		// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< xyzzyErr1
		// fmt.Printf("%x CallOpts\n", ctm.GCfg.CallOpts)
		err = ctm.Call(ctm.GCfg.CallOpts, result, methodName, params...)
		if err != nil {
			fmt.Printf("Error on Contract call to %s: %s, %s\n", methodName, err, godebug.LF())
			/*

					⇒  KeepGroup.getGroupIndex 0x23232323232
					Error on Contract call to getGroupIndex: abi: unmarshalling empty output, File: /Users/corwin/go/src/github.com/pschlump/GCall/gcall.go LineNo:1618
					Error: abi: unmarshalling empty output on call to KeepGroup.getGroupIndex params [0x23232323232]
					Error: abi: unmarshalling empty output on call to KeepGroup.getGroupIndex params [0x23232323232]

					----------------------------------------------------------------------------------------------------------------------

				   xyzzy000 -- cause by call to fucntion returning more than 1 argument -
				   	contracts/CorpRegToken.sol -
				   		function getCapTable(uint256 ii) public view returns ( address aa, uint256 nn_capTableData ) {
				   Error on Contract call to getCapTable: abi: cannot unmarshal common.Address into uint8, File: /Users/corwin/go/src/github.com/pschlump/GCall/gcall.go LineNo:1314
				   Error: abi: cannot unmarshal common.Address into uint8 on call to CorpRegToken.getCapTable params [0]
			*/
			return
		}

		godebug.DbPf(gDebug["db04"], "Call Returns: ->%s<-\n", godebug.SVar(result))
	}

	return
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (ctm *ContractMgr) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return ctm.Contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (ctm *ContractMgr) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return ctm.Contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (ctm *ContractMgr) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return ctm.Contract.Transact(opts, method, params...)
}

// Construct is a free data retrieval call binding the contract method.
func (ctm *ContractMgr) Construct(opts *bind.CallOpts, contractName string) (string, error) {
	ret0 := new(string)
	err := ctm.Contract.Call(opts, ret0, contractName) // constructors have same name as contrct itself
	return *ret0, err
}

// hash.Write([]byte("ReportGreetingEvent(string)"))
// cls is "event" or "function"
func (gCfg *GethInfo) GetCanonicalNameFor(contractName, itemName, cls string) (canonical string) {
	if it, ok := gCfg.ContractList[contractName]; ok {
		// ContractList        map[string]ContractInfo // Map of Contracts
		found := false
		s := itemName
	Loop:
		for _, anABI := range it.ABI {
			s = itemName
			if anABI.Name == itemName {
				if anABI.Type == cls {
					found = true
					s += "("
					com := ""
					for _, rr := range anABI.Inputs {
						s += com + rr.Type
						com = ","
					}
					s += ")"
					break Loop
				}
			}
		}
		if !found {
			// fmt.Printf("Invalid item to watch: [%s]\n", itemName)
			return "**error**"
		}
		return s
	}
	return "***errror***"
}

func (gCfg *GethInfo) GetHashForEvent(contractName, eventName string) (rv string) {
	hash := sha3.NewKeccak256()
	var buf []byte
	// hash.Write([]byte("ReportGreetingEvent(string)"))
	hash.Write([]byte(gCfg.GetCanonicalNameFor(contractName, eventName, "event")))
	buf = hash.Sum(buf)
	result := hex.EncodeToString(buf)
	rv = fmt.Sprintf("0x%s", result)
	return
}

func (gCfg *GethInfo) GetNameForTopic(aTopic string) (eventName string) {
	for contractName, it := range gCfg.ContractList {
		for _, anABI := range it.ABI {
			if anABI.Type == "event" {
				godebug.DbPf(gDebug["db18"], "anABI=%s, %s\n", godebug.SVarI(anABI), godebug.LF())
				eventName = anABI.Name
				hh := gCfg.GetHashForEvent(contractName, eventName)
				godebug.DbPf(gDebug["db18"], "checking hash contractName[%s] eventName[%s], hh=%s, AT:%s\n ", contractName, eventName, hh, godebug.LF())
				if hh == aTopic {
					return
				}
			}
		}
	}
	return "**error not found**"
}

// gCfg.CurrentContract(cmds[1], addr)
func (gCfg *GethInfo) CurrentContract(contractName string, address string) (rv bool) {
	if address == "" {
		return true
	}
	// 1. call the code to get back the compile script from ethereum/geth
	ethBin, err := gCfg.rpc_client.EthGetCode(address, "latest")
	if err != nil {
		fmt.Printf("Attempt to validate code at %s: %s failed - error connecting to geth: %s\n", contractName, address, err)
	}

	// 3. read the .bin file - from make
	binFn := fmt.Sprintf("%s/%s_sol_%s.bin", gCfg.BinariesInDir, contractName, contractName)
	binFn = fmt.Sprintf("%s/%s.bin-runtime", gCfg.BinariesInDir, contractName)
	solcCode, err := ioutil.ReadFile(binFn)
	if err != nil {
		fmt.Printf("uable to read %s - to validate binary in Ethereum, %s, %s\n", binFn, err, godebug.LF())
		return
	}

	// xyzzy - do exec to get!
	solcVersion := `solc, the solidity compiler commandline interface
Version: 0.4.21+commit.dfe3193c.Darwin.appleclang
`

	if len(ethBin) < 2 {
		fmt.Printf("Failed to get an ethBin address with atleast 2 chars ->%s<- for contract %s\n", ethBin, contractName)
		return false
	}

	if len(solcCode) < 2 {
		fmt.Printf("Failed to read binary code with atleast 2 chars ->%s<- for contract %s\n", solcCode, contractName)
		return false
	}

	// 4. check and see if same.  ethBin[2:] chops off the 0x at the beginning.
	ok := bytecode.VerifyCode(contractName, solcVersion, string(solcCode), ethBin[2:])
	if !ok {
		fmt.Printf("%sBinary for contrct %s did not match%s\n", MiscLib.ColorYellow, contractName, MiscLib.ColorReset)
		godebug.DbPf(gDebug["db21"], "ethBin: ->%s<-, bin ->%s<-\n", ethBin[2:], ethBin)
		godebug.DbPf(gDebug["db21"], "len: ethBin: %d, bin %d\n", len(ethBin[2:]), len(ethBin))
	}
	return ok
}

func (gCfg *GethInfo) CheckAllContractsAreCurrent() {
	godebug.DbPf(gDebug["db24"], "automatically check contracts are current at this point, %s\n", godebug.LF())
	for name, ca := range ContractAddressHash {
		// fmt.Printf("AT: %s contract %s, %s\n", godebug.LF(), name, godebug.SVar(ca))
		if ca.ContractAddress == "" { // check for missing address - indicates contract not loaded
			fmt.Printf("%s • %sContract: %s - no address for contract\n", MiscLib.ColorCyan, MiscLib.ColorReset, name)
		} else {
			cok := gCfg.CurrentContract(name, ca.ContractAddress)
			if cok {
				// ✓ U+2713  • U+2715 ✕
				fmt.Printf("%s ✓ %sContract: %s\n", MiscLib.ColorGreen, MiscLib.ColorReset, name)
			} else {
				fmt.Printf("%s ✕ %sContract Did Not Match: %s\n", MiscLib.ColorRed, MiscLib.ColorReset, name)
			}
		}
	}
}

// listOfEvents, err := GetListOfEventsFor(contractName)
func GetListOfEventsFor(contractName string) (evList []string, err error) {
	ABIx, ok := gCfg.ContractList[contractName]
	if !ok {
		fmt.Printf("Contract [%s] is not defined, defined contracts are: %s, %s\n", contractName, gCfg.ContractNames, godebug.LF())
		err = fmt.Errorf("Contract [%s] is not defined, defined contracts are: %s, %s\n", contractName, gCfg.ContractNames, godebug.LF())
		return
	}

	godebug.DbPf(gDebug["db011"], "contractName [%s] %s\n", contractName, godebug.LF())
	ABIraw := ABIx.RawABI

	contractAddress, err := gCfg.GetContractAddress(contractName)
	if err != nil {
		fmt.Printf("Contract missing address [%s]\n", contractName)
		err = fmt.Errorf("Contract missing address [%s]\n", contractName)
		return
	}

	Contract, parsedABI, err := Bind2Contract(ABIraw, contractAddress, gCfg.conn, gCfg.conn, gCfg.conn)
	if err != nil {
		fmt.Printf("Error on Bind2Contract: %s, %s\n", err, godebug.LF())
		err = fmt.Errorf("Error on Bind2Contract: %s, %s\n", err, godebug.LF())
		return
	}
	_ = Contract

	godebug.DbPf(gDebug["db011"], "AT: %s\n", godebug.LF())

	for event := range parsedABI.Events {
		evList = append(evList, event)
	}
	return
}

// Pulled from go-ethereum source and fixed
type TxdataConverted struct {
	AccountNonce uint64          `json:"nonce"`
	Price        *big.Int        `json:"gasPrice"`
	GasLimit     uint64          `json:"gas"`
	Recipient    *common.Address `json:"to"`
	Amount       *big.Int        `json:"value"`
	Payload      []byte          `json:"input"`
	V            *big.Int        `json:"v"`
	R            *big.Int        `json:"r"`
	S            *big.Int        `json:"s"`
	Hash         *common.Hash    `json:"hash"`
}

type Txdata struct {
	AccountNonce string `json:"nonce"`
	Price        string `json:"gasPrice"`
	GasLimit     string `json:"gas"`
	Recipient    string `json:"to"`
	Amount       string `json:"value"`
	Payload      string `json:"input"`
	V            string `json:"v"`
	R            string `json:"r"`
	S            string `json:"s"`
	Hash         string `json:"hash"`
}

type GetTranactionReceiptBodyType struct {
	BlockHash         string        // hex decode
	BlockNumber       string        // hex decode
	ContractAddress   string        // hex decode, only set when loading a contract
	CumulativeGasUsed string        //
	From              string        //
	GasUsed           string        //
	Logs              []interface{} //
	LogsBloom         string        // hex decoe - long
	Status            string        // hex decode, 0 indicates failure, 1 indicates success
	To                string        // address
	TransactionHash   string        // Hash
	TransactionIndex  string        // hex decode
	StatusInt         int           //
	HasLog            bool          //
}

type GetTranactionReceiptType struct {
	JsonRPC string
	Id      int
	Result  GetTranactionReceiptBodyType
}

var n_id = 1

func FetchReceipt(URLToCall, txHash string) (rv string, status int, err error) {

	values := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getTransactionReceipt",
		"id":      n_id,
		"params":  []string{txHash},
	}
	n_id++

	jsonValue, err := json.Marshal(values)
	if err != nil {
		// xyzzy messge
		return
	}

	// URLToCall := "http://10.51.245.213:8545/"

	resp, err := http.Post(URLToCall, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		if db0101 {
			fmt.Printf("Error: %s\n", err)
		}
		return
	}

	if db0101 {
		fmt.Printf("resp: %s err: %s\n", godebug.SVarI(resp), err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if db0101 {
			fmt.Printf("Failed to read body: %s\n", err)
		}
		return
	}
	respstatus := resp.StatusCode
	if respstatus == 200 {
		if db0101 {
			fmt.Printf("Status 200/success - Body is: ->%s<-\n", string(body))
		}
	}

	var bodyDecode GetTranactionReceiptType

	err = json.Unmarshal(body, &bodyDecode)
	if err != nil {
		if db0101 {
			fmt.Printf("Error: %s failed to decode body\n", err)
		}
		return
	}

	bodyDecode.Result.HasLog = len(bodyDecode.Result.Logs) > 0
	StatusInt, err := strconv.ParseInt(bodyDecode.Result.Status, 0, 64)
	if err != nil {
		if db0101 {
			fmt.Printf("Error: %s failed to parse status\n", err)
		}
		return
	}

	bodyDecode.Result.StatusInt = int(StatusInt)
	status = bodyDecode.Result.StatusInt
	rv = fmt.Sprintf("\n%s\n", godebug.SVarI(bodyDecode))

	return
}

// TypeOfSlice print out slice types.  Used in debuging.
func TypeOfSlice(t interface{}) {
	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(t)

		for i := 0; i < s.Len(); i++ {
			fmt.Printf("i=%d: type=%T\n", i, s.Index(i))
		}
	}
}

// ---------------------------------------------------------------------------------------------------------------------

// Section Note:
// 		1. Output of watch "bytes32" data - display better as a hex string
// 			0xBBbbBB... for 32 bytes instead of an array of byte.
//		2. There are going to be other types ( address? uint256 etc) that may need converstion.

// ReturnTypeConverter will Convert return type to have correct datay types so that JSON marshal/unmarshal
// will display it correclty.
func ReturnTypeConverter(rt []interface{}) (rv []interface{}) {
	for ii := 0; ii < len(rt); ii++ {
		t := rt[ii]
		tT := fmt.Sprintf("%T", t)
		if tT == "[32]uint8" {
			uu, ok := t.([32]uint8)
			if !ok {
				panic("Should have conveted")
			}
			var ft EthBytes32
			for jj := 0; jj < 32; jj++ {
				ft[jj] = uu[jj]
			}
			rv = append(rv, ft)
		} else {
			rv = append(rv, t)
		}
		/*
			switch reflect.TypeOf(t).Kind() {
			case reflect.Slice:
				s := reflect.ValueOf(t)

				for i := 0; i < s.Len(); i++ {
					fmt.Printf("i=%d: type=%T", i, s.Index(i))
				}
			default:
				rv = append(rv, rt[ii])
			}
		*/
	}
	return
}

// EthBytes32 is setup to meat the interface{} specification for JSON.
type EthBytes32 [32]uint8

// MarshalJSON takes a named type of [32]uint8 === bytes32 from the ethereum
// return and convers it into a single hex string.
func (ww EthBytes32) MarshalJSON() ([]byte, error) {
	fmt.Printf("In the MarshalJSON for .GCcall/[32]uint8, %s\n", godebug.LF())
	// return []byte(fmt.Sprintf("\"%x\"", ww)), nil
	return []byte(`"0x` + hex.EncodeToString(ww[:]) + `"`), nil
}

// UnmarshalJSON convers a EthBytes32 ([32]uint8 === bytes32) into a usable return
// value.   This is really to meat the interface{} specification for JSON.
func (ww *EthBytes32) UnmarshalJSON(b []byte) error {
	fmt.Printf("In the UnmarshalJSON for .GCcall/[32]uint8, %s\n", godebug.LF())
	// First, deserialize everything into a local data type that matches with data.
	var objMap string
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	xx, err := hex.DecodeString(objMap)
	if err != nil {
		return err
	}

	var wwSlice []uint8

	wwSlice = append(wwSlice, xx...)

	// *ww = wwSlice
	for i := 0; i < 32; i++ {
		if i < len(wwSlice) {
			(*ww)[i] = wwSlice[i]
		} else {
			(*ww)[i] = 0
		}
	}

	fmt.Printf("got ->%x<-, %s\n", ww, godebug.LF())
	return nil
}
