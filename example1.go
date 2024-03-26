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
	"log"
	"math/big"

	"github.com/pschlump/dbgo"
	"github.com/pschlump/ethrpcx"
)

// -------------------------------------------------------------------------------------------------
// From: https://github.com/onrik/ethrpc
// An example!
func SendFunds() {
	client := ethrpcx.NewEthRPC("http://127.0.0.1:8545")

	version, err := client.Web3ClientVersion()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Version: %v\n", version)

	// Send 1 eth
	txid, err := client.EthSendTransaction(ethrpcx.T{
		From:  "0x6247cf0412c6462da2a51d05139e2a3c6c630f0a",
		To:    "0xcfa202c4268749fbb5136f2b68f7402984ed444b",
		Value: big.NewInt(1000000000000000000),
	}) // TODO: err
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Txid = %s\n", dbgo.SVarI(txid))
}
