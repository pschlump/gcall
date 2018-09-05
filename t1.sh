#!/bin/bash

# Get the highest nonce for this account

curl \
	-X POST \
	-H 'Content-Type: application/json;charset=UTF-8' \
	-H 'Accept: application/json, text/plain, /' \
	-H 'Cache-Control: no-cache' \
	--data '{"jsonrpc":"2.0","method":"eth_getTransactionCount","params":[ "0xc2a56884538778bacd91aa5bf343bf882c5fb18b", "latest"],"id":67}' \
	http://192.168.0.158:8545

