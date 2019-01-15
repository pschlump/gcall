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
	"context"
	"fmt"
	"log"

	ethereum "github.com/ethereum/go-ethereum"   //
	"github.com/ethereum/go-ethereum/common"     //
	"github.com/ethereum/go-ethereum/core/types" //
	"github.com/pschlump/MiscLib"                //
	"github.com/pschlump/godebug"                //
)

// --------------------------------------------------------------------------------------------------------
//
// doWatch watches a contract for events.
//
// Input:
// 		contractName, eventName		-- Watch to watch - if eventName == "" then watch all events on contract
//		gCfg 						-- Config
//			gCfg.ContractList[ "contractName" ]
//			gCfg.GetContractAddress( "contractName" )
//			gCfg.conn
//			gCfg.GetNameForTopic(log.Topics[0].String())
//
//
// Uses:
// 		Bind2Contract(...)
// 		ReturnTypeConverter(marshalledValues)
//		TypeOfSlice(marshalledValues)				((debug only))
//
// TODO:
// 		xyzzyW000 - a "quit" chanel for exiting event watches!
// 		xyzzyW001 - a delete-all-watch that clears all the watched stuff.
//
// --------------------------------------------------------------------------------------------------------

func doWatch(contractName, eventName string) (err error) {

	if ABIx, ok := gCfg.ContractList[contractName]; ok { // check that it exists
		godebug.DbPf(gDebug["db11"], "contractName [%s] eventName [%s], %s\n", contractName, eventName, godebug.LF())
		godebug.DbPf(gDebug["db11"], "Found contract [before overload check], %s, %s\n", contractName, godebug.LF())
		ABIraw := ABIx.RawABI

		contractAddress, err := gCfg.GetContractAddress(contractName)
		if err != nil {
			fmt.Printf("Usage: watch ContractName || watch ContractName.EventName - invalid name for contract [%s] specifed.\n", contractName)
			return err
		}

		/* Contract - parse into the go-eth format */
		_, parsedABI, err := Bind2Contract(ABIraw, contractAddress, gCfg.conn, gCfg.conn, gCfg.conn) // keep just the parsedABI
		if err != nil {
			fmt.Printf("Error on Bind2Contract: %s, %s\n", err, godebug.LF())
			return err
		}

		query := ethereum.FilterQuery{
			Addresses: []common.Address{contractAddress},
		}

		godebug.DbPf(gDebug["db11"], "AT: %s\n", godebug.LF())

		var ch = make(chan types.Log)
		ctx := context.Background()

		godebug.DbPf(gDebug["db11"], "AT: %s\n", godebug.LF()) // last working line with truffle, "Subscribe: notifications not supported"

		sub, err := gCfg.conn.SubscribeFilterLogs(ctx, query, ch)
		if err != nil {
			log.Println("Subscribe:", err) // xyzzy  - fix
			return err
		}

		godebug.DbPf(gDebug["db11"], "AT: %s\n", godebug.LF())

		// list out the current watched events! -- capture current events in list
		if watching, ok := CurrentWatchMap[CurrentWatchType{ContractName: contractName, EventName: eventName}]; !ok || !watching {
			CurrentWatchMap[CurrentWatchType{ContractName: contractName, EventName: eventName}] = true
			CurrentWatch = append(CurrentWatch, CurrentWatchType{ContractName: contractName, EventName: eventName})
		} else {
			fmt.Printf("Already watching %s.%s\n", contractName, eventName)
			return err
		}

		// xyzzyW000 - a "quit" chanel for exiting event watches!
		// xyzzyW001 - a delete-all-watch that clears all the watched stuff.

		go func() {
			for {
				godebug.DbPf(gDebug["db11"], "%sWaiting for event at 'select' - AT: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
				select {
				case log := <-ch:
					if len(log.Topics) > 0 {
						name := gCfg.GetNameForTopic(log.Topics[0].String())
						godebug.DbPf(gDebug["db18"], "name [%s] eventName [%s], %s\n", name, eventName, godebug.LF())
						if eventName == "" || name == eventName {
							fmt.Printf("%sCaught Event Log:%s, %s%s\n", MiscLib.ColorGreen, godebug.LF(), godebug.SVarI(log), MiscLib.ColorReset)
							godebug.DbPf(gDebug["db15"], "%sAT:%s name ->%s<-%s\n", MiscLib.ColorYellow, godebug.LF(), name, MiscLib.ColorReset)

							if event, ok := parsedABI.Events[name]; ok {
								godebug.DbPf(gDebug["db15"], "%sAT: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
								arguments := event.Inputs                                 // get the inputs to the event - these will determine the unpack.
								marshalledValues, err := arguments.UnpackValues(log.Data) // marshalledValues is an array of interface{}
								if err != nil {
									fmt.Printf("Error on unmarshalling event data: %s eventName:%s\n", err, name)
								} else {
									// 1. Output of watch "bytes32" data - display better as a hex string
									// 0xBBbbBB... for 32 bytes instead of an array of byte.
									typeModified := ReturnTypeConverter(marshalledValues)
									fmt.Printf("%sEvent Data: %s%s\n", MiscLib.ColorGreen, godebug.SVarI(typeModified), MiscLib.ColorReset)
									godebug.DbPf(gDebug["db15ev"], "%sAT: %s %T %s\n", MiscLib.ColorCyan, godebug.LF(), marshalledValues, MiscLib.ColorReset)
									godebug.DbPf(gDebug["db15ev"], "%sAT: %s %T %s\n", MiscLib.ColorCyan, godebug.LF(), marshalledValues[0], MiscLib.ColorReset)
									if gDebug["db15ev"] {
										TypeOfSlice(marshalledValues)
									}
								}
							} else {
								fmt.Printf("Error failed to lookup event [%s] in ABI\n", name)
							}
						} else {
							godebug.DbPf(gDebug["show.ignored.event"], "%s%s.%s - event ignored; not watched%s\n", MiscLib.ColorYellow, contractName, name, MiscLib.ColorReset)
						}
					}
				case err := <-sub.Err():
					fmt.Printf("AT: %s, error=%s\n", godebug.LF(), err)
					return
				}
				godebug.DbPf(gDebug["db11"], "AT: %s\n", godebug.LF())
			}
		}()

		return nil

	}
	return fmt.Errorf("contrct %s did not exist - no ABI or incorrect contract name", contractName)

}
