package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/keep-network/keep-core/go/interface/lib/FakeStake"
	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

//	"github.com/miguelmota/go-web3-example/fakeStake"

func main() {
	// client, err := ethclient.Dial("ws://rinkeby.infura.io/ws")
	client, err := ethclient.Dial("ws://192.168.0.157:8546")

	if err != nil {
		log.Fatal(err)
	}

	// fakeStakeAddress := "921837fa45906264e3a69e17eff42e3867e889eb"
	fakeStakeAddress := "4fea0149986c13cd36a877dcbdf1b88713c4b04f" // single param version of FakeStaker
	fakeStakeAddress = "8e01ad08008696efe34fb3215ab5b82c13e89f86"  // 2 param version from Antonio
	fakeStakeAddress = "4fea0149986c13cd36a877dcbdf1b88713c4b04f"  // single param version of FakeStaker

	priv := "f5f04d04718e6fd795577fbff5ca1322998c54de8bccb5630493e3d5635d3e12"

	key, err := crypto.HexToECDSA(priv)

	contractAddress := common.HexToAddress(fakeStakeAddress)
	fakeStakeClient, err := FakeStake.NewFakeStake(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("AT: %s\n", godebug.LF())

	auth := bind.NewKeyedTransactor(key)

	fmt.Printf("AT: %s\n", godebug.LF())
	// not sure why I have to set this when using testrpc
	// var nonce int64 = 2323232323891
	// auth.Nonce = big.NewInt(nonce)
	auth.Value = big.NewInt(4712388)         // uint64   // Gas limit to set for the transaction execution (0 = estimate)
	auth.GasPrice = big.NewInt(100000000000) // *big.Int // Gas price to use for the transaction execution (nil = gas price oracle)
	auth.GasLimit = 4712388                  // uint64   // Gas limit to set for the transaction execution (0 = estimate)
	fmt.Printf("AT: %s\n", godebug.LF())

	// tx, err := fakeStakeClient.Stake(auth, big.NewInt(22))
	tx, err := fakeStakeClient.Stake(auth, "22")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("AT: %s\n", godebug.LF())

	fmt.Printf("Pending TX: 0x%x\n", tx.Hash())

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	fmt.Printf("AT: %s\n", godebug.LF())

	var ch = make(chan types.Log)
	ctx := context.Background()

	fmt.Printf("AT: %s\n", godebug.LF()) // last working line with truffle, "Subscribe: notifications not supported"

	sub, err := client.SubscribeFilterLogs(ctx, query, ch)
	if err != nil {
		log.Println("Subscribe:", err)
		return
	}

	fmt.Printf("AT: %s\n", godebug.LF())
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			fmt.Printf("%sWaiting for event at 'select' - AT: %s%s\n", MiscLib.ColorCyan, godebug.LF(), MiscLib.ColorReset)
			select {
			case err := <-sub.Err():
				fmt.Printf("AT: %s\n", godebug.LF())
				log.Fatal(err)
			case log := <-ch:
				fmt.Printf("%sCaught Event Log:%s, %s%s\n", MiscLib.ColorGreen, godebug.LF(), godebug.SVarI(log), MiscLib.ColorReset)
			}
			fmt.Printf("AT: %s\n", godebug.LF())
		}
	}()

	fmt.Printf("Wating for completion of go-routine - AT: %s\n", godebug.LF())
	wg.Wait()
	fmt.Printf("That's all folks... AT: %s\n", godebug.LF())

}
