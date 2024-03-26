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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/pschlump/GCall/jsonSyntaxErrorLib"
	"github.com/pschlump/dbgo"
	"github.com/pschlump/pw"
)

// -------------------------------------------------------------------------------------------------
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// -------------------------------------------------------------------------------------------------
func ExistsIsDir(name string) bool {
	fi, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	if fi.IsDir() {
		return true
	}
	return false
}

// -------------------------------------------------------------------------------------------------
// Get a list of filenames and directorys.
// -------------------------------------------------------------------------------------------------
func GetFilenames(dir string) (filenames, dirs []string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, nil
	}
	for _, fstat := range files {
		if !strings.HasPrefix(string(fstat.Name()), ".") {
			if fstat.IsDir() {
				dirs = append(dirs, fstat.Name())
			} else {
				filenames = append(filenames, fstat.Name())
			}
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------
func RmExt(filename string) string {
	var extension = filepath.Ext(filename)
	var name = filename[0 : len(filename)-len(extension)]
	return name
}

// RmExtIfHasExt will remove an extension from name if it exists.
// TODO: make ext an list of extensions and have it remove any that exists.
//
// name - example abc.xyz
// ext - example .xyz
//
// If extension is not on the end of name, then just return name.
func RmExtIfHasExt(name, ext string) (rv string) {
	rv = name
	if strings.HasSuffix(name, ext) {
		rv = name[0 : len(name)-len(ext)]
	}
	return
}

// -------------------------------------------------------------------------------------------------
var invalidMode = errors.New("Invalid Mode")

func Fopen(fn string, mode string) (file *os.File, err error) {
	file = nil
	if mode == "r" {
		file, err = os.Open(fn) // For read access.
	} else if mode == "w" {
		file, err = os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	} else if mode == "a" {
		file, err = os.OpenFile(fn, os.O_RDWR|os.O_APPEND, 0660)
		if err != nil {
			file, err = os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		}
	} else {
		err = invalidMode
	}
	return
}

// -------------------------------------------------------------------------------------------------
func ParseLineIntoWords(line string) []string {
	// rv := strings.Fields ( line )
	Pw := pw.NewParseWords()
	Pw.SetOptions("C", true, true)
	Pw.SetLine(line)
	rv := Pw.GetWords()
	return rv
}

// -------------------------------------------------------------------------------------------------
// This is to be used/implemented when we add
// 1. ability to chagne the prompt - using templates
// 2. ability to use templates in commands
func SetValue(name, val string) {
	// TODO
}

// ===============================================================================================================================================================================================
var isIntStringRe *regexp.Regexp
var isHexStringRe *regexp.Regexp
var trueValues map[string]bool
var boolValues map[string]bool

func init() {
	isIntStringRe = regexp.MustCompile("([+-])?[0-9][0-9]*")
	isHexStringRe = regexp.MustCompile("(0x)?[0-9a-fA-F][0-9a-fA-F]*")

	trueValues = make(map[string]bool)
	trueValues["t"] = true
	trueValues["T"] = true
	trueValues["yes"] = true
	trueValues["Yes"] = true
	trueValues["YES"] = true
	trueValues["1"] = true
	trueValues["true"] = true
	trueValues["True"] = true
	trueValues["TRUE"] = true
	trueValues["on"] = true
	trueValues["On"] = true
	trueValues["ON"] = true

	boolValues = make(map[string]bool)
	boolValues["t"] = true
	boolValues["T"] = true
	boolValues["yes"] = true
	boolValues["Yes"] = true
	boolValues["YES"] = true
	boolValues["1"] = true
	boolValues["true"] = true
	boolValues["True"] = true
	boolValues["TRUE"] = true
	boolValues["on"] = true
	boolValues["On"] = true
	boolValues["ON"] = true

	boolValues["f"] = true
	boolValues["F"] = true
	boolValues["no"] = true
	boolValues["No"] = true
	boolValues["NO"] = true
	boolValues["0"] = true
	boolValues["false"] = true
	boolValues["False"] = true
	boolValues["FALSE"] = true
	boolValues["off"] = true
	boolValues["Off"] = true
	boolValues["OFF"] = true
}

func IsIntString(s string) bool {
	return isIntStringRe.MatchString(s)
}

func ParseBool(s string) (b bool) {
	_, b = trueValues[s]
	return
	//if InArray(s, []string{"t", "T", "yes", "Yes", "YES", "1", "true", "True", "TRUE", "on", "On", "ON"}) {
	//	return true
	//}
	//return false
}

// -------------------------------------------------------------------------------------------------
func ConvToHexBigInt(s string) (rv *big.Int) {
	s = StripQuote(s)
	rv = big.NewInt(0)
	if strings.HasPrefix(s, "0x") {
		s = s[2:]
	}
	rv.SetString(s, 16)
	return
}

func ConvToDecBigInt(s string) (rv *big.Int) {
	s = StripQuote(s)
	rv = big.NewInt(0)
	rv.SetString(s, 10)
	return
}

func ConvToInt64(s string) (rv int64) {
	rv, _ = strconv.ParseInt(s, 10, 64)
	return
}

func ConvToUInt64(s string) (rv uint64) {
	t, _ := strconv.ParseInt(s, 10, 64)
	rv = uint64(t)
	return
}

func ConvToBool(s string) bool {
	return ParseBool(s)
}

func IsBool(s string) (ok bool) {
	_, ok = boolValues[s]
	return
}

func IsHexNumber(s string) (ok bool) {
	ok = isHexStringRe.MatchString(s)
	return
}

func IsNumber(s string) (ok bool) {
	ok = isIntStringRe.MatchString(s)
	return
}

func IsString(pp string) (rv bool) {
	return true
}

func HexOf(ss string, base int) (rv byte) { // still working on this
	t, err := strconv.ParseInt(ss, base, 64)
	if err != nil {
		fmt.Printf("Warning: HexOf: error with >%s< as input, %s\n", ss, err)
		return 0
	}
	rv = byte(t)
	return
}

func ConvNumberToByte32(pp string) (rv [32]byte) {
	// TBD xyzzy503
	pp = StripQuote(pp)
	base := 10
	if strings.HasPrefix(pp, "0x") {
		pp = pp[2:]
		base = 16
	}
	for ii := 0; ii < 32; ii++ {
		rv[ii] = 0
	}
	// xyzzy - if base == 16, then we do the hex thing, if == 10 then use a big.Int() -- TODO - not implemented yet.
	for ii := 0; ii < len(pp) && ii < 64; ii += 2 {
		if ii+2 <= len(pp) {
			rv[ii/2] = HexOf(pp[ii:ii+2], base)
		} else {
			rv[ii/2] = HexOf(pp[ii:ii+1]+"0", base)
		}
	}
	return
}

func ConvHexNumberToByte32(pp string) (rv [32]byte) {
	rv = ConvNumberToByte32(pp)
	return
}

func ConvStringToByte32(pp string) (rv [32]byte) {
	for ii := 0; ii < 32; ii++ {
		rv[ii] = 0
	}
	for ii := 0; ii < len(pp) && ii < 64; ii++ {
		rv[ii] = pp[ii]
	}
	return
}

// Close will Disconnect/Close connections to Geth
func (ctm *ContractMgr) Close() {
	// TBD
}

// -------------------------------------------------------------------------------------------------
func ReadConfig(fn string) (rv GethInfo) {

	// Create some defaults
	rv.UnlockSeconds = 600
	rv.ContractList = make(map[string]ContractInfo)
	rv.ContractList = make(map[string]ContractInfo)
	rv.ContractNames = make([]string, 0, 10)
	rv.DebugFlags = make([]string, 0, 10)

	data, err := ioutil.ReadFile(fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Must supply config file %s, error=%s\n", fn, err)
		os.Exit(1)
	}

	err = json.Unmarshal(data, &rv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing %s, error=%s\n", fn, err)
		PrintErrorJson(string(data), err)
		os.Exit(1)
	}

	if strings.HasPrefix(rv.FromAddressPassword, "$ENV$") {
		pw := os.Getenv(rv.FromAddressPassword[5:])
		if pw == "" {
			fmt.Printf("No password set for FromAddressPassword in %s\n", fn)
		}
		rv.FromAddressPassword = pw
	}

	if strings.HasPrefix(rv.KeyFilePassword, "$ENV$") {
		pw := os.Getenv(rv.KeyFilePassword[5:])
		if pw == "" {
			fmt.Printf("No password set for KeyFilePassword in %s\n", fn)
		}
		rv.KeyFilePassword = pw
	}

	return
}

// -------------------------------------------------------------------------------------------------
// name, abi := ReadABI(fn)
func ReadABI(fn string) (name string, ABI []ABIType, raw string) {

	dbgo.DbPf(gDebug["db19"], "fn=%s AT=%s\n", fn, dbgo.LF())

	// Infer name from fn
	s0 := strings.Split(filepath.Base(fn), "_") // maybee immediately preceding _sol_
	for i := 0; i < len(s0); i++ {
		if len(s0[i]) > 0 {
			s0 = s0[i:]
			break
		}
	}
	if len(s0) > 0 {
		name = s0[0]
	}
	dbgo.DbPf(gDebug["db19"], "fn=%s s0=%s name=->%s<-, AT=%s\n", fn, dbgo.SVar(s0), name, dbgo.LF())

	data, err := ioutil.ReadFile(fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Must supply ABI config file %s, error=%s\n", fn, err)
		os.Exit(1)
	}
	raw = string(data)

	err = json.Unmarshal(data, &ABI)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing ABI %s, error=%s\n", fn, err)
		PrintErrorJson(string(data), err)
		os.Exit(1)
	}
	return
}

// -------------------------------------------------------------------------------------------------
func ReadContractAddressHash(fn string) (ContractAddressHash map[string]AContractAddressType) {
	ContractAddressHash = make(map[string]AContractAddressType)
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Must supply ABI config file %s, error=%s\n", fn, err)
		os.Exit(1)
	}
	err = json.Unmarshal(data, &ContractAddressHash)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing ABI %s, error=%s\n", fn, err)
		PrintErrorJson(string(data), err)
		os.Exit(1)
	}
	return
}

// -------------------------------------------------------------------------------------------------
func StripQuote(s string) string {
	if len(s) > 0 && s[0] == '"' { // only double quotes around prompt with blanks in it.
		s = s[1:]
		if len(s) > 0 && s[len(s)-1] == '"' {
			s = s[:len(s)-1]
		}
	} else if len(s) > 0 && s[0] == '\'' { // only double quotes around prompt with blanks in it.
		s = s[1:]
		if len(s) > 0 && s[len(s)-1] == '\'' {
			s = s[:len(s)-1]
		}
	}
	return s
}

func PrintErrorJson(js string, err error) (rv string) {
	rv = jsonSyntaxErrorLib.GenerateSyntaxError(js, err)
	fmt.Printf("%s\n", rv)
	return
}

// KeysFromMap returns an array of keys from a map.
func KeysFromMap(a interface{}) (keys []string) {
	xkeys := reflect.ValueOf(a).MapKeys()
	keys = make([]string, len(xkeys))
	for ii, vv := range xkeys {
		keys[ii] = vv.String()
	}
	return
}

// SetDebugFlags convers from --debug csv,csv -> gDebug
func SetDebugFlags() {
	if Debug != nil {
		df := strings.Split(*Debug, ",")
		for _, dd := range df {
			if _, have := gDebug[dd]; have {
				gDebug[dd] = !gDebug[dd]
			} else {
				gDebug[dd] = true
			}
		}
	}
}

var helpText = `GCall - a tool for accessing Ethereum (solidity) contracts

help 
	This command.  Prints out this text.

current 
	Check that the source is current with the compiled contracts on chain.
	This will check that all contracts that have an address are current.

current Contract
	Check that the source is current with the compiled contracts on chain.
	This checks that just the specified contract is current.

list 
	Show a list of contracts that can be called.   Contract must have an address
	to be called.

list ContractName	 
	Show the callable functions in that contract and what the parameters and return
	values are.

Name.Method ParametersOptional 
	Call the named method with the parameters.  If this
	is a non-transactional (constant) call then it may return data.  Note: variables 
	in the contract that are public show up as functions without any parameters that
	can be called to get the current value.
	
	You can put quotes around parameters if they are strings and need to contain blanks.
	
	Numbers are assumed to be in base 10, unless preceded with '0x'.  

	Booleans are one of yes|Yes|YES|no|No|NO|true|True|TRUE|false|False|False.

	bytes32 values are 0x000000000000 hex numbers

quit 
	Exit from QCall (Also exit, :q, :wq, \q, bye, logout, quit;, exit;, bye;)

prompt String
	Allows setting of a different prompt for the program.  Eventually this will
	have the ability to display the current connected network name (not implemented
	yet).  Probably some other useful information too.

	{{.__network__}}	display the network name
	{{.__line_no__}}	display the current input line number
	{{env "NAME"}}		Pull in exported environment variable

watch ContractName 
watch ContractName.EventName 
	Watch a contract or an event in a contract and report when the event
	occurs.

	watch --list		Show the currently watched contracts and events

proxyfor ProxyContractName ContractName
	Make ProxyContractName have all the ABI calls of ContrctName so it
	can be called as a proxy.

setValue ### - set the amount of "Wei" to send with the transaction.


`

// I bid you adieu!

/* vim: set noai ts=4 sw=4: */
