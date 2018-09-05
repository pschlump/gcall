#!/bin/bash

# Unlock account

curl \
	-X POST \
	-H 'Content-Type: application/json;charset=UTF-8' \
	-H 'Accept: application/json, text/plain, /' \
	-H 'Cache-Control: no-cache' \
	--data '{"jsonrpc":"2.0","method":"personal_unlockAccount","params":[ "0xc2a56884538778bacd91aa5bf343bf882c5fb18b", "password", 15000],"id":67}' \
	http://192.168.0.158:8545
