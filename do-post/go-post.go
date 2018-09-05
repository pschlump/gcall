package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
)

// SEE: https://stackoverflow.com/questions/24455147/how-do-i-send-a-json-string-in-a-post-request-in-go

/*
   IP="http://127.0.0.1:8545/

   curl  \
   -H "Content-Type: application/json" \
   -X POST \
   --data "{\"jsonrpc\":\"2.0\", \"method\":\"eth_getTransactionReceipt\",\"params\":[\"${Tx}\"],\"id\":1}" \
   ${IP}| tee ,receipt.out

   check-json-syntax -p ,receipt.out | tee receipt.out.json

{
	"jsonrpc":"2.0",
	"id":1,
	"result":{
		"blockHash":"0x5917d2cf36968d8bac114dd5d84d64ba8ff488a602644a4d110cbb5cb444d55f",
		"blockNumber":"0x68ec",
		"contractAddress":null,
		"cumulativeGasUsed":"0x5daa",
		"from":"0x023e291a99d21c944a871adcc44561a58f99bdbc",
		"gasUsed":"0x5daa",
		"logs":[],
		"logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"status":"0x0",
		"to":"0x568ec7fd7c2904f26cbc398d34d12f07d85839f5",
		"transactionHash":"0x736aceb665ae6aadef64368c05d1516f9a00af23d16178e9573ac70f63297651",
		"transactionIndex":"0x0"
	}
}

--------------------------------------------------------------------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------------------------------------------------------------------
--------------------------------------------------------------------------------------------------------------------------------------------------

success s=->
{
	"JsonRPC": "2.0",
	"Id": 1,
	"Result": {
		"BlockHash": "0xd17909170cd43e08ef41455a2f58e4daabfd5811c0f7c1a597fe36758951bf68",
		"BlockNumber": "0x7573",
		"ContractAddress": "",
		"CumulativeGasUsed": "0x167f2",
		"From": "0x6ffba2d0f4c8fd7961f516af43c55fe2d56f6044",
		"GasUsed": "0x167f2",
		"Logs": [
			{
				"address": "0x53578f1ec1130ffdc4400623f3109faa7dbc8d4c",
				"blockHash": "0xd17909170cd43e08ef41455a2f58e4daabfd5811c0f7c1a597fe36758951bf68",
				"blockNumber": "0x7573",
				"data": "0x000000000000000000000000000000000000000000000000000000000000004700000000000000000000000000000000000000000000000000000000000000c8000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000007573",
				"logIndex": "0x0",
				"removed": false,
				"topics": [
					"0x72e71cec9908703e388c4f242f74c20f1afb406dcfea214956e1ddce0e1e6ff3"
				],
				"transactionHash": "0x54f8c2d9fda7c2b94f78ca5ecd6d081c067934a04bac2e1071ac27725b1c2417",
				"transactionIndex": "0x0"
			}
		],
		"LogsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000200000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000080",
		"Status": "0x1",
		"To": "0x53578f1ec1130ffdc4400623f3109faa7dbc8d4c",
		"TransactionHash": "0x54f8c2d9fda7c2b94f78ca5ecd6d081c067934a04bac2e1071ac27725b1c2417",
		"TransactionIndex": "0x0",
		"StatusInt": 1,
		"HasLog": true
	}
}

*/

// Pulled from go-ethereum source and fixed
type Txdata struct {
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

type CfgType struct {
	URLToCall string
}

var Cfg = flag.String("cfg", "cfg.json", "config file for this call") // 0

func main() {

	flag.Parse() // Parse CLI arguments to this, --cfg <name>.json

	fns := flag.Args()
	if len(fns) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: do-post [--cfg cfg.json]\n")

		os.Exit(1)
	}

	gCfg := ReadConfig(*Cfg) // var gCfg GethInfo

	// URLToCall := "http://10.51.245.213:8545/"
	// txHash := "0x736aceb665ae6aadef64368c05d1516f9a00af23d16178e9573ac70f63297651"
	for _, txHash := range fns {
		s, status, err := FetchReceipt(gCfg.URLToCall, txHash)
		if err != nil {
			fmt.Printf("%serror getting tranaction: %s%s\n", MiscLib.ColorRed, err, MiscLib.ColorReset)
		} else if status == 0 {
			fmt.Printf("%sfailed s=->%s<-%s\n", MiscLib.ColorYellow, s, MiscLib.ColorReset)
		} else {
			fmt.Printf("%ssuccess s=->%s<-%s\n", MiscLib.ColorGreen, s, MiscLib.ColorReset)
		}
	}
}

func ReadConfig(fn string) (cfg CfgType) {
	buf, err := ioutil.ReadFile(fn)
	if err != nil {
		fmt.Printf("Error reading [%s] error: %s\n", fn, err)
		os.Exit(1)
	}

	err = json.Unmarshal(buf, &cfg)
	if err != nil {
		fmt.Printf("Error parsing [%s] error: %s\n", fn, err)
		os.Exit(1)
	}

	return
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

const db0101 = false
