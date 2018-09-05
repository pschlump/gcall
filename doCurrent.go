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
	"fmt"
	"os/exec"

	"github.com/chzyer/readline" //
	"github.com/pschlump/GCall/args"
	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

func doCurrent(cmds []string, rl *readline.Instance) {

	// ---- Process arguments to the list command -----------------------------------------------------
	aa, err := args.GetArgs(&args.ArgConfigType{
		Config: map[string]args.AnArgType{},
	}, 1, cmds)
	if err != nil {
		aa.Usage("current")
	} else {
		godebug.Printf(gDebug["db19"], "results: %s\n", godebug.SVarI(aa))
	}

	// Run make to rebuild any binaries that we need. -------------------------------------------------
	makeArgs := gCfg.RebuildBinary
	if len(makeArgs) > 1 {
		exCmd := exec.Command(makeArgs[0], makeArgs[1:]...)
		err := exCmd.Run()
		if err != nil {
			fmt.Printf("Build of binaries did not work, %s\n", err)
		}
	} else {
		exCmd := exec.Command(makeArgs[0])
		err := exCmd.Run()
		if err != nil {
			fmt.Printf("Build of binaries did not work, %s\n", err)
		}
	}

	if len(aa.Remainder) > 0 {
		for ii, anArg := range aa.Remainder {
			godebug.Printf(gDebug["db24"], "current contract[%s] #%d, %s\n", anArg, ii, godebug.LF())
			if addr, ok := ContractAddressHash[anArg]; ok {
				if addr.ContractAddress != "" {
					cok := gCfg.CurrentContract(anArg, addr.ContractAddress)
					if cok {
						fmt.Printf("%s ✓ %sContract: %s\n", MiscLib.ColorGreen, MiscLib.ColorReset, anArg)
					} else {
						fmt.Printf("%s ✕ %sContract Did Not Match: %s\n", MiscLib.ColorRed, MiscLib.ColorReset, anArg)
					}
					continue
				}
				fmt.Printf("Do not have an address for %s specified in %s\n", anArg, "gcall.addr.json") // xyzzy - get actual file name
			}
		}
	} else {
		gCfg.CheckAllContractsAreCurrent()
	}

}
