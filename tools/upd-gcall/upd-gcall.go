package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"

	"github.com/pschlump/GCall/jsonSyntaxErrorLib"
	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

// Running migration: 2_initial_migration.js
//   Deploying SimpleToken...
//  ... 0x62cdc5ac81af06524df81393d3f4aa5a2ec90ed089dcb6e0ea09a58f84020568
//  SimpleToken: 0xbfae441f574715049dc33b75cd5f5dd9dcd5c459

func main() {

	nameAddr := make(map[string]string)
	name := ""

	// fmt.Printf("%s%s%s\n", MiscLib.ColorCyan, "yep", MiscLib.ColorReset)
	st := 0
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		godebug.Printf(db2, "%sLine:%s st:%d %s\n", MiscLib.ColorCyan, line, st, MiscLib.ColorReset)
		if startMigration.MatchString(line) {
			st = 1
		} else if st == 1 && deploying.MatchString(line) {
			// grab name
			match := deploying.FindStringSubmatch(line)
			name = match[1]
			godebug.Printf(db1, "AT: %s name ->%s<-\n", godebug.LF(), name)
			st = 2
		} else if st == 1 && replacing.MatchString(line) {
			// grab name
			match := replacing.FindStringSubmatch(line)
			name = match[1]
			godebug.Printf(db1, "AT: %s name ->%s<-\n", godebug.LF(), name)
			st = 2
		} else if st == 2 {
			st = 3
		} else if st == 3 {
			// grab Addr
			reS := name + ": (0x.*)"
			grabAddr, err := regexp.Compile(reS)
			if err != nil {
				fmt.Printf("Invalid RE: %s error %s\n", reS, err)
				os.Exit(1)
			}
			match := grabAddr.FindStringSubmatch(line)
			addr := match[1]
			nameAddr[name] = addr
			godebug.Printf(db1, "AT: %s name ->%s%s%s<- addr ->%s<-\n", godebug.LF(), MiscLib.ColorGreen, name, MiscLib.ColorReset, addr)
			st = 1
		}

	}

	// TODO / xyzzy - udpate the gcall.addr.cfg
	contractAddrFile := "gcall.addr.cfg"
	ContractAddressHash := ReadContractAddressHash(contractAddrFile)
	start := time.Now()
	ts := start.Format(time.RFC3339)

	for name, addr := range nameAddr {
		ContractAddressHash[name] = AContractAddressType{
			ContractAddress: addr,
			LoadedAt:        ts,
		}
	}

	WriteContractAddressHash(contractAddrFile, ContractAddressHash)
}

var deploying *regexp.Regexp
var replacing *regexp.Regexp
var startMigration *regexp.Regexp

func init() {
	deploying = regexp.MustCompile("Deploying ([^.]*)")
	replacing = regexp.MustCompile("Replacing ([^.]*)")
	// Running migration: 2_deploy_contracts.js
	startMigration = regexp.MustCompile("Running mig.* 2_.*js")
}

type AContractAddressType struct {
	ContractAddress string
	LoadedAt        string
}

var ContractAddressHash map[string]AContractAddressType

func ReadContractAddressHash(fn string) (ContractAddressHash map[string]AContractAddressType) {
	ContractAddressHash = make(map[string]AContractAddressType)
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &ContractAddressHash)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing ABI %s, error=%s\n", fn, err)
		PrintErrorJson(string(data), err)
		os.Exit(1)
	}
	return
}

func WriteContractAddressHash(fn string, ContractAddressHash map[string]AContractAddressType) {
	ioutil.WriteFile(fn, []byte(godebug.SVarI(ContractAddressHash)+"\n"), 0644)
}

func PrintErrorJson(js string, err error) (rv string) {
	rv = jsonSyntaxErrorLib.GenerateSyntaxError(js, err)
	fmt.Printf("%s\n", rv)
	return
}

const db1 = true
const db2 = true
