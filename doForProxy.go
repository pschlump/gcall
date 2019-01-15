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
	"fmt"

	"github.com/chzyer/readline" //
	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

// proxyfor ProxyName ActualContract
func doForProxy(cmds []string, rl *readline.Instance) {
	fmt.Printf("forproxy - setup a proxy copntract, %s\n", cmds)
	var err error

	// 1. Take contract ProxyName - verify that it has a
	//	{
	//		"payable": true,
	//		"stateMutability": "payable",
	//		"type": "fallback"
	//	},
	// in it's ABI - if so then
	//
	// 2. Take the ActualContract and pull out all the calls in it that are type "function"
	// 3. Merge these into the ABI for the proxy
	//		{ "mergedfrom": "contract_name" }
	//	added to each ABI row that is copied in.

	if len(cmds) != 3 {
		fmt.Printf("Usage: proxyfor ProxyName ActualContracName\n")
		return
	}

	fmt.Printf("Proxy: %s Actual: %s\n", cmds[1], cmds[2])
	ProxyName := cmds[1]
	ActualName := cmds[2]

	if addr, ok := ContractAddressHash[ProxyName]; !ok || addr.ContractAddress == "" {
		fmt.Printf("Contract %s not defined, must have an address\n", ProxyName)
		return
	}
	if addr, ok := ContractAddressHash[ActualName]; !ok || addr.ContractAddress == "" {
		fmt.Printf("Contract %s not defined, must have an address\n", ActualName)
		return
	}

	ABIx, ok := gCfg.ContractList[ProxyName]
	if !ok {
		fmt.Printf("Contract [%s] is not defined, ABI not defined\n", ProxyName)
		return
	}

	ABI := ABIx.ABI
	ProxyABIraw := ABIx.RawABI

	isProxy := false
	godebug.DbPf(gDebug["db38"], "ABI=%s\n", godebug.SVar(ABI))
	for _, aAbi := range ABI {
		if aAbi.Type == "fallback" {
			isProxy = true
			break
		}
	}

	godebug.DbPf(gDebug["db38"], "isProxy = %v, %s\n", isProxy, godebug.LF())

	if isProxy {
		fmt.Printf("%s is a proxy contract\n", ProxyName)
	} else {
		fmt.Printf("Error: %s is %s*NOT*%s a proxy contract\n", ProxyName, MiscLib.ColorRed, MiscLib.ColorReset)
		return
	}

	ABIto, ok := gCfg.ContractList[ActualName]
	if !ok {
		fmt.Printf("Contract [%s] is not defined, ABI not defined\n", ActualName)
		return
	}
	ActualABIraw := ABIto.RawABI

	// Keep a lookup table of the names in the original proxy contract - can not override any named items.
	// Warn below if you are there is a name overlap.
	names := make(map[string]int)
	for _, aAbi := range ABIx.ABI { // unsorted version
		if aAbi.Name != "" && aAbi.Type == "function" {
			names[aAbi.Name] = 1
		}
	}

	MergeData := make([]ABIType, 0, len(ABIto.ABI))
	for _, aAbi := range ABIto.ABI { // unsorted version
		if aAbi.Type == "function" {
			aAbi.MergeFrom = ActualName
			MergeData = append(MergeData, aAbi)
		}
		if aAbi.Type == "event" {
			aAbi.MergeFrom = ActualName
			MergeData = append(MergeData, aAbi)
		}
	}

	godebug.DbPf(gDebug["db38"], "%s are the found functions, %s\n", godebug.SVarI(MergeData), godebug.LF())

	godebug.DbPf(gDebug["db38"], "Before: %s\n", godebug.SVarI(gCfg.ContractList[ProxyName]))
	ABIx.ABI = append(ABIx.ABI, MergeData...)
	gCfg.ContractList[ProxyName] = ABIx
	godebug.DbPf(gDebug["db38"], "After: %s\n", godebug.SVarI(gCfg.ContractList[ProxyName]))

	// We need to mangle Raw?

	proxyTmp := make([]map[string]interface{}, 0, 25)
	actualTmp := make([]map[string]interface{}, 0, 25)
	err = json.Unmarshal([]byte(ProxyABIraw), &proxyTmp)
	if err != nil {
		fmt.Printf("Unable to parse >%s<, error %s\n", ProxyABIraw, err)
	}
	err = json.Unmarshal([]byte(ActualABIraw), &actualTmp)
	if err != nil {
		fmt.Printf("Unable to parse >%s<, error %s\n", ActualABIraw, err)
	}

	// fmt.Printf("actualTmp=%s, %s\n", godebug.SVarI(actualTmp), godebug.LF())

	for _, dat := range actualTmp {
		newName := ""
		if vv, ok := dat["name"]; ok {
			if w, ok := vv.(string); ok {
				newName = w
			}
		}
		if found, ok := names[newName]; ok && found == 1 {
			fmt.Printf("%sWarning: %s.%s is hiden/orverwriden by proxy contract %s.%s - calls to that will be handled by proxy.%s\n", MiscLib.ColorYellow, ActualName, newName, ProxyName, newName, MiscLib.ColorReset)
		} else {
			if vv, ok := dat["type"]; ok {
				if w, ok := vv.(string); ok {
					if w == "function" {
						proxyTmp = append(proxyTmp, dat)
					}
					if w == "event" {
						proxyTmp = append(proxyTmp, dat)
					}
				}
			}
		}
		names[newName] = 2
	}

	godebug.DbPf(gDebug["db39"], "RawABI Before: %s\n", godebug.SVarI(gCfg.ContractList[ProxyName]))
	godebug.DbPf(gDebug["db39"], "RawABI Before: -->%s<--\n", gCfg.ContractList[ProxyName].RawABI)

	ABIx.RawABI = godebug.SVarI(proxyTmp)
	gCfg.ContractList[ProxyName] = ABIx

	godebug.DbPf(gDebug["db39"], "RawABI After: %s\n", godebug.SVarI(gCfg.ContractList[ProxyName]))
	godebug.DbPf(gDebug["db39"], "RawABI After: -->%s<--\n", gCfg.ContractList[ProxyName].RawABI)

}
