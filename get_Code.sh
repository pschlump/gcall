#!/bin/bash

# Unlock account

curl \
	-X POST \
	-H 'Content-Type: application/json;charset=UTF-8' \
	-H 'Accept: application/json, text/plain, /' \
	-H 'Cache-Control: no-cache' \
	--data '{"jsonrpc":"2.0","method":"eth_getCode","params":[ "0x84aef122b06582b68d3e57c11ac4ed75aef01aeb" , "latest"],"id":65}' \
	http://192.168.0.158:8545

