

#MIT License
#
#Copyright (c) 2018 Philip Schlump
#
#Permission is hereby granted, free of charge, to any person obtaining a copy
#of this software and associated documentation files (the "Software"), to deal
#in the Software without restriction, including without limitation the rights
#to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
#copies of the Software, and to permit persons to whom the Software is
#furnished to do so, subject to the following conditions:
#
#The above copyright notice and this permission notice shall be included in all
#copies or substantial portions of the Software.
#
#THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
#IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
#FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
#AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
#LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
#OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
#SOFTWARE.


all: ./contracts/Greeter_sol_Greeter.abi ./contracts/mortal_sol_mortal.abi
	go build
	cp gcall ~/bin

install:
	cp GCall ~/bin

./contracts/mortal_sol_mortal.abi: ./contracts/mortal.sol
	( cd contracts ; solc --abi -o . mortal.sol )
	( cd contracts ; /Users/corwin/bin/check-json-syntax -p <mortal_sol_mortal.abi >,tmp ; mv ,tmp mortal_sol_mortal.abi )
	( cd contracts ; solc --bin -o . mortal.sol )
	( cd contracts ; solc --bin-runtime -o . Greeter.sol )
	( cd contracts ; abigen --abi mortal_sol_mortal.abi --pkg mortal --out mortal.go )
	mkdir -p ./lib/mortal
	cp ./contracts/mortal.go ./lib/mortal
	( cd ./lib/mortal ; go build )

./contracts/Greeter_sol_Greeter.abi: ./contracts/Greeter.sol
	( cd contracts ; solc --abi -o . Greeter.sol ; mv Greeter.abi Greeter_sol_Greeter.abi )
	( cd contracts ; /Users/corwin/bin/check-json-syntax -p <Greeter_sol_Greeter.abi >,tmp ; mv ,tmp Greeter_sol_Greeter.abi )
	( cd contracts ; solc --bin -o . Greeter.sol ; mv Greeter.bin Greeter_sol_Greeter.bin )
	( cd contracts ; solc --bin-runtime -o . Greeter.sol )
	( cd contracts ; abigen --abi Greeter_sol_Greeter.abi --pkg Greeter --out Greeter.go )
	mkdir -p ./lib/Greeter
	cp ./contracts/Greeter.go ./lib/Greeter
	cp ./contracts/Greeter_sol_Greeter.abi ./abi
	( cd ./lib/Greeter ; go build )

x1:
	cp ./contracts/mortal_sol_mortal.abi ./abi
	cp ./contracts/Greeter_sol_Greeter.abi ./abi

#testnet_migrate:
#	truffle migrate --reset	--network testnet

local_migrate:
	truffle migrate --reset	--network local


# note
# solc_when_soljs:
# 	( cd contracts ; solc --abi Greeter.sol )
