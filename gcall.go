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
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/chzyer/readline"                        //
	"github.com/ethereum/go-ethereum/accounts/abi/bind" //
	"github.com/ethereum/go-ethereum/common"            //
	"github.com/ethereum/go-ethereum/ethclient"         //
	"github.com/pschlump/GCall/args"                    //
	"github.com/pschlump/MiscLib"                       //
	"github.com/pschlump/ethrpcx"
	"github.com/pschlump/godebug" //
)

// OLD: "github.com/onrik/ethrpc" - modified with new functions -

var Cfg = flag.String("cfg", "cfg.json", "config file for this call")            // 0
var GethURL_ws = flag.String("ws", "", "Websocket address for geth")             // 1 -- needed for event watch
var GethURL_http = flag.String("http", "", "HTTP address for geth")              // 2 -- needed for JSON RPC calls
var Debug = flag.String("debug", "", "comma sep list of debugt flags to toggle") // 3
var Input = flag.String("input", "", "Input file")                               // 4 - Xyzzy -- not implemented yet --
var Output = flag.String("output", "", "Output file")                            // 5 - Xyzzy -- not implemented yet --
var Network = flag.String("network", "", "name of network connected to")         // 6

var gDebug map[string]bool
var in *os.File
var out *os.File
var log0 *os.File
var g_prompt0_x string = "--> "
var g_prompt string = "--> "
var g_if_cli bool = true
var line_no = 0
var m_line_no = 0
var contractAddrFile = "gcall.addr.json"

var CurrentWatchMap map[CurrentWatchType]bool

/*
{
	"Greeter": {
		"ContractAddress": "0x88b8dc1fa683c44356b0f644b57fd5ce6ca357e9"
	}
}
*/
type AContractAddressType struct {
	ContractAddress string
}

var ContractAddressHash map[string]AContractAddressType
var gCfg GethInfo

func init() {
	gDebug = make(map[string]bool)
	gDebug["ethrpcx.echo"] = false
	gDebug["dump.ABIMethod"] = false
	gDebug["db01"] = false
	gDebug["db02"] = false
	gDebug["db03"] = false
	gDebug["db04"] = false
	gDebug["db05"] = false
	gDebug["db06"] = false
	gDebug["db07"] = false
	gDebug["db08"] = false
	gDebug["db09"] = false
	gDebug["db10"] = false
	gDebug["db11"] = false
	gDebug["db12"] = false
	gDebug["show.ignored.event"] = true

	in = os.Stdin // to support the --input and --output flags
	out = os.Stdout
	g_if_cli = true

	CurrentWatchMap = make(map[CurrentWatchType]bool)
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem("mode",
		readline.PcItem("vi"),
		readline.PcItem("emacs"),
	),
	readline.PcItem("list",
		readline.PcItemDynamic(listContracts()),
	),
	readline.PcItem("current",
		readline.PcItemDynamic(listContracts()),
	),
	readline.PcItem("watch",
		readline.PcItemDynamic(listContractsEvents()),
	),
	readline.PcItem("proxyfor",
		readline.PcItemDynamic(listContracts()),
		readline.PcItemDynamic(listContracts()),
	),
	readline.PcItem("quit"),
	readline.PcItem("exit"),
	readline.PcItem("bye"),
	readline.PcItem("prompt"),
	readline.PcItem("help"),
	readline.PcItemDynamic(listContractsMethods()),
)

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func main() {

	// -------------------------------------------------------------------------------------------------------------

	flag.Parse() // Parse CLI arguments to this, --cfg <name>.json

	gCfg = ReadConfig(*Cfg) // var gCfg GethInfo

	// fmt.Printf("gCfg.SRCPath=%s\n", gCfg.SRCPath)

	if Cfg == nil {
		fmt.Printf("--cfg is a required parameter\n")
		os.Exit(1)
	}

	for _, dd := range gCfg.DebugFlags {
		gDebug[dd] = true
	}

	// var Network = flag.String("network", "", "name of network connected to")         // 6
	if Network != nil && *Network != "" {
		contractAddrFile = fmt.Sprintf("gcall.%s.addr.json", *Network)
		fmt.Printf("%sNetwork Set To: %s%s\n", MiscLib.ColorGreen, *Network, MiscLib.ColorReset)
	}

	if !Exists(*Cfg) {
		fmt.Printf("Missing %s - required configuration file\n", *Cfg)
		os.Exit(1)
	}

	// -------------------------------------------------------------------------------------------------------------------------------
	// xyzzy012-network may need to open gcall.log specific to network - if read it for addresses
	// -------------------------------------------------------------------------------------------------------------------------------

	log0, err := Fopen("gcall.log", "a")
	if err != nil {
		fmt.Printf("Unable to open log for append, %s, %s\n", "gcall.log", err)
	}
	if Network != nil && *Network != "" {
		if log0 != nil {
			fmt.Fprintf(log0, `{ "lt":"network", "network": %q }`+"\n\n", *Network)
		}
		if override, ok := gCfg.NetworkFlag[*Network]; ok {
			var IfNotEmpty = func(a, b string) string {
				if a != "" {
					return a
				}
				return b
			}
			gCfg.GethURL_ws = IfNotEmpty(override.GethURL_ws, gCfg.GethURL_ws)
			gCfg.GethURL_http = IfNotEmpty(override.GethURL_http, gCfg.GethURL_http)
			gCfg.FromAddress = IfNotEmpty(override.FromAddress, gCfg.FromAddress)
			gCfg.FromAddressPassword = IfNotEmpty(override.FromAddressPassword, gCfg.FromAddressPassword)
			gCfg.KeyFile = IfNotEmpty(override.KeyFile, gCfg.KeyFile)
			gCfg.KeyFilePassword = IfNotEmpty(override.KeyFilePassword, gCfg.KeyFilePassword)
		}
	}
	if *GethURL_http != "" || *GethURL_ws != "" {
		gCfg.GethURL_ws = *GethURL_ws
		gCfg.GethURL_http = *GethURL_http
	}

	// -------------------------------------------------------------------------------------------------------------------------------
	// Recompile any binaries, regenerate any ABI's from SRCPath -> ABIPath (xyzzy010)
	// -------------------------------------------------------------------------------------------------------------------------------

	// -------------------------------------------------------------------------------------------------------------------------------
	// may need to read in log and find addresses of contracts (migrations) before open for append.  (xyzzy012)
	// -------------------------------------------------------------------------------------------------------------------------------

	ContractAddressHash = ReadContractAddressHash(contractAddrFile)

	SetDebugFlags() // convers from --debug csv,csv -> gDebug

	// Read in Map of ABI's for potential calls
	dirty := false
	for _, ap := range gCfg.ABIPath {
		if !ExistsIsDir(ap) {
			fmt.Printf("%sError: [%s] appears to not be correct for path to .abi files - not found, should be a sirectory%s\n",
				MiscLib.ColorRed, ap, MiscLib.ColorReset)
		}
		fns, dirs := GetFilenames(ap)
		if len(dirs) > 0 {
			fmt.Printf("Warning: ignoring sub-directory(s) %s in %s\n", dirs, ap)
			fmt.Printf("      If you need to include the sub-directories then list them in the config file (%s).\n", *Cfg)
		}
		if len(fns) == 0 {
			fmt.Printf("%sError: [%s] appears to not be correct for path to .abi files - none found%s\n", MiscLib.ColorRed, ap, MiscLib.ColorReset)
		}
		godebug.DbPf(gDebug["db1000"], "%sAT: %s fns=%s %s\n", MiscLib.ColorYellow, godebug.LF(), fns, MiscLib.ColorReset)
		missing := make([]string, 0, 15)
		for _, fn := range fns {
			// ----------------------------------------------------------------------------------------------
			// filter fns for .abi at end!
			// ----------------------------------------------------------------------------------------------
			if extension := filepath.Ext(fn); extension != ".abi" {
				// fmt.Printf("Extension [%s]\n", extension) // cool, extension includes '.'
				continue
			}
			fn2 := ap + "/" + fn
			name, abi, raw := ReadABI(fn2)
			name = RmExtIfHasExt(name, ".abi")
			godebug.DbPf(gDebug["db1000"], "%sAT: %s name=[%s] %s\n", MiscLib.ColorYellow, godebug.LF(), name, MiscLib.ColorReset)

			if ct, ok := gCfg.ContractList[name]; ok {
				ct.ABI = abi
				godebug.DbPf(gDebug["db1000"], "%sAT: %s len ct.ABI=%d%s\n", MiscLib.ColorYellow, godebug.LF(), len(abi), MiscLib.ColorReset)
				ct.RawABI = raw
				gCfg.ContractList[name] = ct
			} else {
				dirty = true
				godebug.DbPf(gDebug["db1000"], "%sAT: %s len ct.ABI=%d%s\n", MiscLib.ColorCyan, godebug.LF(), len(abi), MiscLib.ColorReset)
				gCfg.ContractList[name] = ContractInfo{
					Name:    name,
					ABI:     abi,
					RawABI:  raw,
					Version: "v1.0.0", // see-below - pull from source
					// Address         string	 -- // -- See-Below - need to find this / report that we don't have this.
				}
				// fmt.Printf("%sAT: %s len ct.ABI=%d%s\n", MiscLib.ColorYellow, godebug.LF(), len(abi), MiscLib.ColorReset)
				gCfg.ContractNames = append(gCfg.ContractNames, name)
			}

			ai, ok := ContractAddressHash[name]
			if !ok {
				missing = append(missing, name) //	fmt.Printf("Info: Missing address for contract [%s] in %s\n", name, contractAddrFile)
			} else {
				ct := gCfg.ContractList[name]
				ct.Address = ai.ContractAddress
				ct.address = common.HexToAddress(ct.Address)
				gCfg.ContractList[name] = ct
			}
		}
		if len(missing) > 0 {
			fmt.Printf("Info: Missing address for contracts %+v in %s, this may be due to contract inheritance.\n", missing, contractAddrFile)
		}
	}

	// --------------- readline setup -------------------------------------------------------------------------------------------------
	// func NewEx(cfg *Config) (*Instance, error) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:              "\033[31m⇒ \033[0m ",
		HistoryFile:         gCfg.ReadlineDir + "/.gcall.readline.data", // "./.readline.data",
		AutoComplete:        completer,
		InterruptPrompt:     "^C",
		EOFPrompt:           "exit",
		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		fmt.Printf("Error setting up readline: %s\n", err)
		os.Exit(1)
	}
	defer rl.Close()

	// xyzzy012 - read the "log" and find out the address and migration info

	// 		xyzzy012 - find block# and last date-time that each contract was loaded
	// 		xyzzy012 - compare load date/time to source date-time to see if out of date
	// 		xyzzy012 - for any out-of-date contracts - compile and compare binary - to see if modified

	// 		xyzzy - Post process data - find constrctors

	if dirty {
		godebug.DbPf(gDebug["db06"], "Config is dirty -need- to update\n")
	}

	godebug.DbPf(gDebug["dumpCfg0"], "AT: %s gCfg=%s\n", godebug.LF(), godebug.SVarI(gCfg))

	// --------------------------------------------------------------------------------------------------------
	// Get connection to Geth
	conn, err := ethclient.Dial(gCfg.GethURL_ws)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v at address: %s", err, gCfg.GethURL_ws)
	}

	// xyzzy - do we need the other kind of connection to geth? rpcclient?

	gCfg.conn = conn

	client := ethrpcx.NewEthRPC(gCfg.GethURL_http)
	client.Db01 = gDebug["ethrpcx.echo"] // Echo message sent on JSON RPC	// xyzzy - settable from --> cli
	client.Db02 = gDebug["ethrpcx.echo"] // Echo responce recevied back
	client.LogWith = func(s string) {
		if log0 != nil {
			fmt.Fprintf(log0, `{ "lt":"callreturn", "data": %q }`+"\n\n", s)
		}
	}

	gCfg.rpc_client = client

	// --------------------------------------------------------------------------------------------------------
	version, err := gCfg.rpc_client.Web3ClientVersion()
	if err != nil {
		fmt.Printf("Error! AT: %s %serr=%s%s\n", godebug.LF(), MiscLib.ColorRed, err, MiscLib.ColorReset)
	}

	fmt.Printf("%sConnected To GETH Version: %s%s\n", MiscLib.ColorGreen, version, MiscLib.ColorReset)

	// --------------------------------------------------------------------------------------------------------
	if gCfg.UnlockSeconds < 600 {
		fmt.Printf("Warning: unlockSeconds is too small, minimum 60\n")
	}

	doUnlockAccount := func() {
		// Unlock specified account - keep track of time so will re-unlock it as necessary
		// Find go code to call for unlock of account
		// Xyzzy - Create a "go-routine" that will lock and call to unlock account for duration of processing [ set to 10 min for unlock ]
		x := gCfg.rpc_client.PersonalUnlockAccount(gCfg.FromAddress, gCfg.FromAddressPassword, gCfg.UnlockSeconds)
		if x == nil {
			godebug.DbPf(gDebug["db06"], "AT: %s %sSuccess on personal.unlockAccount()%s\n", godebug.LF(), MiscLib.ColorGreen, MiscLib.ColorReset)
			fmt.Printf("%sSuccess Account(%s) Unlocked%s\n", MiscLib.ColorGreen, gCfg.FromAddress, MiscLib.ColorReset)
		} else {
			fmt.Printf("Error! AT: %s %sx=%s%s\n", godebug.LF(), MiscLib.ColorRed, x, MiscLib.ColorReset)
		}
		godebug.DbPf(gDebug["db07"], "AT: %s %sSuccessfully unlocked account - ready to start main loop%s\n", godebug.LF(), MiscLib.ColorGreen, MiscLib.ColorReset)
	}

	err = gCfg.SetTransactOpts()
	if err != nil {
		fmt.Printf("Error: %s, %s\n", err, godebug.LF())
		os.Exit(1)
	}
	/********
	  type CallOpts struct {
	  	Pending bool           // Whether to operate on the pending state or the last known one
	  	From    common.Address // Optional the sender address, otherwise the first account is used

	  	Context context.Context // Network context to support cancellation and timeouts (nil = no timeout)
	  }
	  ********/
	requestFrom := common.HexToAddress(gCfg.FromAddress)
	gCfg.CallOpts = &bind.CallOpts{
		Pending: false,
		From:    requestFrom,
		Context: nil,
	}

	doUnlockAccount()

	// --------------------------------------------------------------------------------------------------------
	ticker := time.NewTicker(time.Duration(gCfg.UnlockSeconds-300) * time.Second)
	quitUnlock := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				doUnlockAccount()
			case <-quitUnlock:
				ticker.Stop()
				return
			}
		}
	}()

	// --------------------------------------------------------------------------------------------------------

	// -------------------------------------------------------------------------------------------------------------------------------
	// Validate loaded contracts at addresses
	// An "auto" flag that will check each contract v.s. source === AutoCurrentCheck
	// -------------------------------------------------------------------------------------------------------------------------------
	if ParseBool(gCfg.AutoCurrentCheck) {
		gCfg.CheckAllContractsAreCurrent()
	}

	// -------------------------------------------------------------------------------------------------------------------------------
	// read/Parse commands - Run the interpreter
	// -------------------------------------------------------------------------------------------------------------------------------

	g_if_cli = true

	SetValue("__input_file__", "--stdin--")
	SetValue("__line_no__", "0")
	rl.SetVimMode(true)
	line_no = 0
	m_line_no = 0
	g_prompt = g_prompt0_x
Loop:
	for { // loop until ReadLine returns nil (signalling EOF) or until a quit command is entered.
		SetValue("__line_no__", fmt.Sprintf("%d", line_no))

		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				println()
				break Loop // exit loop when EOF(^D) is entered
			} else {
				continue Loop
			}
		} else if err == io.EOF {
			println()
			break Loop // exit loop when EOF(^D) is entered
		}

		// xyzzy72 - line := ExecuteATemplate (*result, g_data)

		cmds := ParseLineIntoWords(line)
		if len(cmds) <= 0 {
			continue
		}

		if len(cmds[0]) > 0 && cmds[0][0:1] == "#" { // process comments '#'.*
			continue
		}

		switch cmds[0] {

		// -------------------------------------------------------------------------------------------------------
		// proxyfor ProxyName ActualContract
		case "forproxy", "proxyfor":
			doForProxy(cmds, rl)

		// -------------------------------------------------------------------------------------------------------
		// mode						Print out current mode
		// mode vi					Set to vi mode
		// mode emacs				Set to emacs mode
		case "mode":
			doMode(cmds, rl)

		// -------------------------------------------------------------------------------------------------------
		// current 				-- check all scripts again
		// current Contract		-- check just this one script
		case "current":
			doCurrent(cmds, rl)

		// -------------------------------------------------------------------------------------------------------
		// list
		// list Contract
		// list --watched
		case "list":
			var getParams = func(in []ABI_IO_Type, out []ABI_IO_Type) (rv string) {
				rv = "("
				com := " "
				for _, pp := range in {
					rv += com + pp.Type + " " + pp.Name
					com = ", "
				}
				if len(in) <= 0 {
					rv += ")"
				} else {
					rv += " )"
				}
				if len(out) > 0 {
					rv += " returns ( "
					com = ""
					for _, pp := range out {
						rv += com + pp.Type
						com = ","
					}
					rv += " )"
				}
				return
			}

			// ---- Process arguments to the list command -----------------------------------------------------
			aa, err := args.GetArgs(&args.ArgConfigType{
				Config: map[string]args.AnArgType{
					"watched": args.AnArgType{
						Name:    "--watched",
						Abrev:   "-w",
						NoValue: true,
					},
				},
			}, 1, cmds)
			if err != nil {
				aa.Usage("list")
			} else {
				godebug.DbPf(gDebug["db19"], "results: %s\n", godebug.SVarI(aa))
			}

			// ---- Implementation  --------------------------------------------------------------------------
			if err != nil {
			} else if val, ok := aa.Set["watched"]; ok && val == "1" {
				fmt.Printf("%-30s %-30s\n", "Contract Name", "Event Name")
				fmt.Printf("%-30s %-30s\n", "-------------------", "-----------------")
				for _, ww := range CurrentWatch {
					if ww.EventName == "" {
						fmt.Printf("%-30s %-30s\n", ww.ContractName, "--all events--")
					} else {
						fmt.Printf("%-30s %-30s\n", ww.ContractName, ww.EventName)
					}
				}
			} else if len(aa.Remainder) > 0 {
				for _, contractName := range aa.Remainder {
					// contractName := aa.Remainder[0]
					ABIx, ok := gCfg.ContractList[contractName]
					if !ok {
						fmt.Printf("Contract [%s] is not defined, defined contracts are: %s, %s\n", contractName, gCfg.ContractNames, godebug.LF())
						continue Loop
					}
					ABI := ABIx.ABI
					godebug.DbPf(gDebug["db18"], "ABI=%s\n", godebug.SVar(ABI))
					fmt.Printf("%-30s %-5s %-2s %-42s\n", "Method Name", "Const", "$", "Params")
					fmt.Printf("%-30s %-5s %-2s %-42s\n", "------------------------------", "-----", "--", "------------------------------------------")

					// ----------------------------- sort list -------------------------------------------------------------------------------------------------------------------
					nameMap := make(map[string]int)
					for ii, aAbi := range ABI { // unsorted version
						if aAbi.Name != "" {
							nameMap[aAbi.Name] = ii
						}
					}
					keys := KeysFromMap(nameMap)
					sort.Strings(keys)

					//for _, aAbi := range ABI { // unsorted version
					//	if aAbi.Name != "" {
					for _, key := range keys {
						at := nameMap[key]
						aAbi := ABI[at]
						pay := " "
						PAY := "X"

						// fmt.Printf("xyzzy099 %s\n", godebug.SVarI(aAbi))

						// "type": "constructor"
						// "type": "fallback"

						if aAbi.Constant {
							pay = fmt.Sprintf("%sconst%s", MiscLib.ColorGreen, MiscLib.ColorReset)
							PAY = " "
						} else if aAbi.Type == "function" {
							pay = fmt.Sprintf("%sTx%s", MiscLib.ColorCyan, MiscLib.ColorReset)
							// xyzzy ---------------------------------------------------------------------------------------------------------
							PAY = " "
							if aAbi.Payable {
								PAY = "$E"
							}
							// xyzzy ---------------------------------------------------------------------------------------------------------
						} else if aAbi.Type == "event" {
							pay = fmt.Sprintf("%sevent%s", MiscLib.ColorYellow, MiscLib.ColorReset)
							PAY = " "
							// } else if aAbi.Type == "constructor" {
							// 	pay = fmt.Sprintf("%scons.%s", MiscLib.ColorYellow, MiscLib.ColorReset)
							// } else if aAbi.Type == "fallback" {
							// 	pay = fmt.Sprintf("%sfallb.%s", MiscLib.ColorYellow, MiscLib.ColorReset)
						} else {
							pay = fmt.Sprintf("%s-?-%s", MiscLib.ColorRed, MiscLib.ColorReset)
							PAY = " "
						}
						fmt.Printf("%-30s %-17s %-2s %-42s\n", aAbi.Name, pay, PAY, getParams(aAbi.Inputs, aAbi.Outputs)) // note 11 chars added for color set/reset
						//	}
					}
					for _, aAbi := range ABI { // unsorted version
						if aAbi.Name == "" {
							pay := " "
							PAY := " "

							// "type": "constructor"
							// "type": "fallback"

							if aAbi.Type == "constructor" {
								pay = fmt.Sprintf("%sconstructor%s", MiscLib.ColorCyan, MiscLib.ColorReset)
								fmt.Printf("%-30s %-17s    %-42s\n", pay, "", getParams(aAbi.Inputs, aAbi.Outputs)) // note 11 chars added for color set/reset
							} else if aAbi.Type == "fallback" {
								pay = fmt.Sprintf("%sfallback(proxy)%s", MiscLib.ColorCyan, MiscLib.ColorReset)
								PAY = " "
								if aAbi.Payable {
									PAY = "$E"
								}
								fmt.Printf("%-30s %-17s %-2s %-42s\n", pay, "", PAY, getParams(aAbi.Inputs, aAbi.Outputs)) // note 11 chars...
							} else {
								pay = fmt.Sprintf("%s-?-%s", MiscLib.ColorRed, MiscLib.ColorReset)
								fmt.Printf("%-30s %-17s    %-42s\n", aAbi.Name, pay, getParams(aAbi.Inputs, aAbi.Outputs)) // note 11 chars added for color set/reset
							}
						}
					}
				}
			} else {
				if db82 {
					fmt.Printf("List - else case, at %s\n", godebug.LF())
					fmt.Printf("gCfg.ContractNames = %s, should not be empty\n", godebug.SVarI(gCfg.ContractNames))
					fmt.Printf("ContractAddressHash = %s, should not be empty, should match ContractNames\n", godebug.SVarI(ContractAddressHash))
				}
				// tmp := KeysFromMap(gCfg.ContractNames)
				sort.Strings(gCfg.ContractNames)
				fmt.Printf("%4s %-30s %-42s\n", "No.", "Contract Name", "Address")
				fmt.Printf("%4s %-30s %-42s\n", "---", "------------------------------", "------------------------------------------")
				for ii, contractName := range gCfg.ContractNames {
					addr := ContractAddressHash[contractName].ContractAddress
					fmt.Printf(" %03d %-30s %-42s\n", ii+1, contractName, addr)
				}
			}

		// prompt "NewPromptString"
		case "prompt":
			doPrompt(cmds, rl)

		// echo a b c
		case "echo":
			doEcho(cmds, rl)

		// setValue ####
		case "setValue":
			doSetValue(cmds, rl)

		// setValue ####
		case "getValue":
			doGetValue(cmds, rl)

		// quit
		case "quit", "exit", ":q", ":wq", "\\q", "bye", "logout", "quit;", "exit;", "bye;":
			break Loop

		// help
		case "help":
			doHelp(cmds, rl)

		// watch --list		- list currently watched events.
		// watch --verbose 	- to print out entire transaction		(TBD)
		// watch --delete  	- To remove watch on an event			(TBD)
		case "watch":
			// /Users/corwin/go/src/www.2c-why.com/Corp-Reg/MidGeth/EthEventWatch		xyzzyEthWatch
			contractName := ""
			eventName := ""
			if len(cmds) > 1 { // xyzzy len(cmds) == 2?
				contractName = cmds[1]
			} else {
				fmt.Printf("Usage: watch ContractName || watch ContractName.EventName - No name was specified\n")
				continue Loop
			}
			if contractName == "--list" {
				// var CurrentWatch []CurrentWatchType
				// CurrentWatch = append ( CurrentWatch, &CurrentWatchType{ ContractName : contractName, EventName : eventName, })
				fmt.Printf("%-30s %-30s\n", "Contract Name", "Event Name")
				fmt.Printf("%-30s %-30s\n", "------------------------", "------------------------")
				for _, cw := range CurrentWatch {
					fmt.Printf("%-30s %-30s\n", cw.ContractName, cw.EventName)
				}
				continue Loop
			}
			if strings.Index(contractName, ".") == -1 && len(cmds) > 2 {
				contractName, eventName = cmds[1], cmds[2]
				fmt.Printf("new Processing [%s].[%s]\n", contractName, eventName)
			} else if p2 := strings.Split(contractName, "."); len(p2) == 2 {
				contractName, eventName = p2[0], p2[1]
			} // xyzzy else if len(p2) > 2 ??

			// xyzzy - --color red|green|cyan|yellow - what collor to print the event in.
			// xyzzy - check to see if a watch on this already exists - if so - then skip adding!
			// xyzzy - log that we are setting up a watch

			if len(eventName) == 0 {
				listOfEvents, err := GetListOfEventsFor(contractName)
				if err != nil {
					fmt.Printf("Did not get a list of events: %s\n", err)
				}
				fmt.Printf("watch contractName: [%s] for all events %s\n", contractName, godebug.SVar(listOfEvents))
			} else {
				fmt.Printf("watch contractName: [%s] for event named: [%s]\n", contractName, eventName)
			}

			doWatch(contractName, eventName)

		default:
			// /Users/corwin/go/src/www.2c-why.com/Corp-Reg/MidGeth/EthContractCall	xyzzyEthCall
			godebug.DbPf(gDebug["db02"], "cmd: %s was not recognized, %s\n", cmds[0], godebug.LF())
			godebug.DbPf(gDebug["db32"], "cmds[...]: %s, %s\n", cmds, godebug.LF())
			done := false

			// 		xyzzy - may want to check for "./" or / at leading of cmd - indicates not a contract call.

			// 		xyzzy - run a contract Name.Func Plist...
			// Greeter.getGreeting p1, p2 ...  -- Note bob.sh would split into "bob" and "sh" - but not be found as a contract.
			// if you need to run a script that has the same name as a contract then ./name.sh will do the trick, or use full path.
			p2 := strings.Split(cmds[0], ".")
			if len(p2) == 2 {
				contractName, methodName := p2[0], p2[1]
				godebug.DbPf(gDebug["dump.contractInfo"], "contractName [%s] methodName [%s], %s\n", contractName, methodName, godebug.LF())
				godebug.DbPf(gDebug["dump.contractInfo"], "%sAT: %s After = %d%s\n", MiscLib.ColorYellow, godebug.LF(), len(gCfg.ContractList[contractName].ABI), MiscLib.ColorReset)
				if ABIx, ok := gCfg.ContractList[contractName]; ok { // check that it exists	// xyzzy900
					godebug.DbPf(gDebug["db01"], "contractName [%s] methodName [%s], %s\n", contractName, methodName, godebug.LF())
					godebug.DbPf(gDebug["db01"], "Found contract [before overload check], %s, %s\n", contractName, godebug.LF())
					done = true // If looks like a contract but mis-match of parameters - then call it an error, don't look for a script.
					ABIraw := ABIx.RawABI

					address, err := gCfg.GetContractAddress(contractName)
					if err != nil {
						fmt.Printf("Error: %s, %s\n", err, godebug.LF())
						continue Loop
					}

					Contract, parsedABI, err := Bind2Contract(ABIraw, address, gCfg.conn, gCfg.conn, gCfg.conn)
					if err != nil {
						fmt.Printf("Error on Bind2Contract: %s, %s\n", err, godebug.LF())
						continue Loop
					}
					_ = parsedABI

					ctm := NewContractMgr(Contract, &gCfg)

					// fmt.Printf("%sAT: %s After = %d%s\n", MiscLib.ColorYellow, godebug.LF(), len(ctm.GCfg.ContractList[contractName].ABI), MiscLib.ColorReset)

					ABIMethod := -1 // gCfg.ContractList[contractName].ABI[ABIMethod]
					{
						var ok bool
						if ABIMethod, ok = gCfg.IsValidMethodName(contractName, methodName); !ok {
							fmt.Printf("Invalid method name [%s] not defined in contract [%s], %s\n", methodName, contractName, godebug.LF())
							continue Loop
						}
					}

					godebug.DbPf(gDebug["dump.ABIMethod"], "ABIMethod = %d\n", ABIMethod)

					// this is be the "CallContract" section
					param := cmds[1:]
					nParam := len(param)
					godebug.DbPf(gDebug["db05"], "%sParam: %s, nParam=%d %s, %s\n", MiscLib.ColorCyan, param, nParam, MiscLib.ColorReset, godebug.LF())

					ABIMethodSet, nfound := gCfg.IsValidMethodNameSet(contractName, methodName, nParam)
					if nfound == 0 {
						fmt.Printf("Incorrect number of parameters passed, no methods have %d parameters, you passed %s, %s\n", nParam, param, godebug.LF())
						// Find the set of name/match - then - report on the set. (xyzzy001) -- better error reporting
						// xyzzy gCfg.PrintMatchingFunctions ( ABIMethodSet, contractName, methodName )
						continue Loop
					}

					// --------------------------------------------------------------------------------
					// Cooerce to each fo the items in ABIMethodSet
					// --------------------------------------------------------------------------------

					// First!
					// "type": "bytes1[]"

					// Then add....
					// "type": "int"
					// "type": "uint"
					// "type": "int8"		8,16,24,32,...64
					// "type": "uint8"		8,16,24,32,...64
					// "type": "address"
					// "type": "bytes1[]"
					// "type": "bytes1"		1..32

					// "type": "enum"		It's an int [ range check? ]

					// -------------------------------------------------------------------------------- --------------------------------------------------------------------------------
					matchOverload := true
					iParam := make([]interface{}, nParam, nParam+1)
					usedItemNo := ABIMethod
					// fmt.Printf("%sAT: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
					if nParam > 0 {
						// fmt.Printf("%sAT: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
						matchOverload = false

					Outer:
						for _, itemNo := range ABIMethodSet {
							matchOverload = true
							// fmt.Printf("%sAT: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
							usedItemNo = itemNo
							godebug.DbPf(gDebug["db10"], "%sAT: %s usedItemNo=%d %s\n", MiscLib.ColorCyan, godebug.LF(), usedItemNo, MiscLib.ColorReset)

						TypeConv:
							for iP, aParam := range param {

								// -------------------------------------------------------------------------------------------------------------------
								// -------------------------------------------------------------------------------------------------------------------
								// -------------------------------------------------------------------------------------------------------------------
								// xyzzy500 // bytes32 type v.s. uint256 v.s. address - what to do? -- Need conversion to correct type for pass to contract
								// -------------------------------------------------------------------------------- ----------------------------------
								// -------------------------------------------------------------------------------- ----------------------------------
								// -------------------------------------------------------------------------------- ----------------------------------

								// fmt.Printf("%sAT: %s iP=%d aParam=->%s<- %s\n", MiscLib.ColorCyan, godebug.LF(), iP, aParam, MiscLib.ColorReset)
								switch gCfg.ContractList[contractName].ABI[itemNo].Inputs[iP].Type {
								case "bool":
									// fmt.Printf("%sAT: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
									if IsBool(aParam) {
										// fmt.Printf("%sAT: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
										iParam[iP] = ConvToBool(param[iP])
									} else {
										// fmt.Printf("%sAT: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
										matchOverload = false
										break TypeConv
									}
								case "string":
									// fmt.Printf("%sAT: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
									iParam[iP] = StripQuote(param[iP])
								case "uint256":
									// fmt.Printf("%sAT: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
									if IsNumber(aParam) {
										// fmt.Printf("%sAT: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
										iParam[iP] = ConvToDecBigInt(param[iP])
									} else if IsHexNumber(aParam) {
										// fmt.Printf("%sAT: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
										iParam[iP] = ConvToHexBigInt(param[iP])
									} else {
										// fmt.Printf("%sAT: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
										matchOverload = false
										break TypeConv
									}

								case "int256":
									// fmt.Printf("%sAT: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
									if IsNumber(aParam) {
										// fmt.Printf("%sAT: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
										iParam[iP] = ConvToDecBigInt(param[iP])
									} else if IsHexNumber(aParam) {
										// fmt.Printf("%sAT: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
										iParam[iP] = ConvToHexBigInt(param[iP])
									} else {
										// fmt.Printf("%sAT: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
										matchOverload = false
										break TypeConv
									}

								case "address":
									// fmt.Printf("%sAT: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
									s := StripQuote(param[iP])
									if len(s) > 2 && s[0:2] != "0x" {
										s = "0x" + s
									}
									godebug.DbPf(gDebug["db31"], "%s Converting ->%s<- string into address param #[%d], %s%s\n", MiscLib.ColorCyan, s, iP, godebug.LF(), MiscLib.ColorReset)
									a := common.HexToAddress(s) // a := common.Address(s)
									iParam[iP] = a

									// xyzzy000 - may bemissing other types!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

								case "bytes32":
									if IsNumber(aParam) {
										iParam[iP] = ConvNumberToByte32(param[iP])
									} else if IsHexNumber(aParam) {
										iParam[iP] = ConvHexNumberToByte32(param[iP])
									} else if IsString(aParam) {
										s := StripQuote(param[iP])
										iParam[iP] = ConvStringToByte32(s)
									} else {
										// fmt.Printf("%sAT: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
										matchOverload = false
										break TypeConv
									}

								default:
									fmt.Printf("%sBad Type: %s AT: %s%s\n", MiscLib.ColorCyan, gCfg.ContractList[contractName].ABI[itemNo].Inputs[iP].Type, godebug.LF(), MiscLib.ColorReset)
									iParam[iP] = param[iP]
								}
							}
							// fmt.Printf("%sAT: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
							if matchOverload {
								godebug.DbPf(gDebug["db10"], "%sAT: %s -- successful match and coerce of parameters -- %s\n", MiscLib.ColorGreen, godebug.LF(), MiscLib.ColorReset)
								break Outer // Success!
							}
						}
					}

					if !matchOverload {
						fmt.Printf("AT: %s - failed to match overloded function paramters\n", godebug.LF())
						// xyzzy gCfg.PrintMatchingFunctions ( ABIMethodSet, contractName, methodName )
					} else {

						godebug.DbPf(gDebug["db02"], "%sAT: %s -- .CallContract() now -- %s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)

						res, Tx, err := ctm.CallContract(gCfg.ContractList[contractName].ABI[usedItemNo], contractName, methodName, iParam...)

						if err != nil {
							fmt.Printf("Error: %s on call to %s.%s params %s\n", err, contractName, methodName, param)
							if log0 != nil {
								fmt.Printf("Error: %s on call to %s.%s params %s\n", err, contractName, methodName, param)
							}
						} else {
							fmt.Printf("Tranaction Is: %s\n", godebug.SVarI(res))

							// create a log entry for later -- so we can see when contract was run.
							if log0 != nil {
								if err != nil {
									fmt.Fprintf(log0, `{ "lt":"call", "callTo":"%s.%s", "params": %s, "err": %s }`+"\n\n",
										contractName, methodName, godebug.SVar(param), godebug.SVar(err))
								} else if Tx != nil {
									fmt.Fprintf(log0, `{ "lt":"call", "callTo":"%s.%s", "params": %s, "Tx": %s }`+"\n\n",
										contractName, methodName, godebug.SVar(param), godebug.SVar(Tx))
								} else {
									fmt.Fprintf(log0, `{ "lt":"call", "callTo":"%s.%s", "params": %s }`+"\n\n",
										contractName, methodName, godebug.SVar(param))
								}
							}

							// -------------------------------------------------------------------------------------------------------------------
							// xyzzy501 // Take "tx" and make call to get receipt - as a go-process - backgorund
							// -------------------------------------------------------------------------------- ----------------------------------

							/*
								   IP="http://127.0.0.1:8545/

								   curl  \
								   -H "Content-Type: application/json" \
								   -X POST \
								   --data "{\"jsonrpc\":\"2.0\", \"method\":\"eth_getTransactionReceipt\",\"params\":[\"${Tx}\"],\"id\":1}" \
								   ${IP}| tee ,receipt.out

								   check-json-syntax -p ,receipt.out | tee receipt.out.json


									values := map[string]string{"username": username, "password": password}

									jsonValue, err := json.Marshal(values)

									resp, err := http.Post(authAuthenticatorUrl, "application/json", bytes.NewBuffer(jsonValue))
							*/
							if err != nil {
								fmt.Fprintf(os.Stderr, `{ "lt":"call/err", "callTo":"%s.%s", "params":%s, "err":%q, "at":%q }`+"\n",
									contractName, methodName, godebug.SVar(param), err, godebug.LF())
							} else if Tx != nil {
								fmt.Fprintf(os.Stderr, `{ "lt":"call/tx", "callTo":"%s.%s", "params": %s, "Tx":%q }`+"\n",
									contractName, methodName, godebug.SVar(param), godebug.SVar(Tx))
							} else {
								if len(param) == 0 {
									fmt.Fprintf(os.Stderr, `{ "lt":"call/view", "callTo":"%s.%s" }`+"\n", contractName, methodName)
								} else {
									fmt.Fprintf(os.Stderr, `{ "lt":"call/view", "callTo":"%s.%s", "params": %s }`+"\n",
										contractName, methodName, godebug.SVar(param))
								}
							}

							var TxdataUnpack Txdata
							err = json.Unmarshal([]byte(godebug.SVarI(Tx)), &TxdataUnpack)
							if err != nil {
								fmt.Printf("Error unpacking data: %s\n", err)
							}
							TxHash := TxdataUnpack.Hash

							fmt.Printf("%sTxHash: %s%s\n", MiscLib.ColorCyan, TxHash, MiscLib.ColorReset)

							if Tx != nil {
								go func(URLToCall, txHash string) {
									var s string
									var status int
									var err error
									for i := 0; i < 100; i++ {
										s, status, err = FetchReceipt(URLToCall, txHash)
										if err == nil {
											break
										}
										// fmt.Printf("Before Sleep %d\n", i)
										fmt.Printf(".")
										time.Sleep(5 * time.Second)
										//fmt.Printf("After Sleep\n")
										fmt.Printf(".")
									}
									if err != nil {
										fmt.Printf("%serror getting tranaction: %s%s\n", MiscLib.ColorRed, err, MiscLib.ColorReset)
									} else if status == 0 {
										fmt.Printf("%sfailed s=->%s<-%s\n", MiscLib.ColorYellow, s, MiscLib.ColorReset)
									} else {
										fmt.Printf("%sTranaction succeded success Receipt:%s%s\n", MiscLib.ColorGreen, s, MiscLib.ColorReset)
									}
								}(gCfg.GethURL_http, TxHash)
							}

						}

					}

				} // xyzzy900
			}

			if !done {
				// 		TODO xyzzy002 - run a file ./fn, or in "path"
				fmt.Printf("TODO/not implemented yet: Search for script to run - using path\n")
				if len(cmds) > 0 && Exists(cmds[0]) {
					fmt.Printf("%sTODO/Command exists - should run it! AT: %s%s\n", MiscLib.ColorGreen, godebug.LF(), MiscLib.ColorReset)
					done = true // bold assumption!
				}
			}

			if !done {
				fmt.Printf("cmd: %s was not recognized, %s\n", cmds[0], godebug.LF())
			}

		} // switch

	} // main for loop
}

const db0101 = false
const db82 = false

// I bid you adieu!

/* vim: set noai ts=4 sw=4: */
