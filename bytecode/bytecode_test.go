package bytecode

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func Test_Bytecode01(t *testing.T) {

	tests := []struct {
		solcVersion  string
		binRuntimeFn string
		ethGetCodeFn string
		shouldMatch  bool
		debugFlag    bool
	}{
		{
			solcVersion: `solc, the solidity compiler commandline interface
Version: 0.4.21+commit.dfe3193c.Darwin.appleclang
`,
			binRuntimeFn: "../testdata/Greeter.bin-runtime",
			ethGetCodeFn: "../testdata/0x84aef122b06582b68d3e57c11ac4ed75aef01aeb",
			shouldMatch:  true,
			debugFlag:    true,
		},
	}

	type Result struct {
		Result string
	}
	var getResult Result

	for ii, test := range tests {
		db1 = test.debugFlag
		solcCode, err := ioutil.ReadFile(test.binRuntimeFn)
		if err != nil {
			fmt.Printf("Error -unable to read file %s, %s\n", test.binRuntimeFn, err)
			t.Errorf("Error %2d, IO Error on test\n", ii)
			continue
		}
		eth, err := ioutil.ReadFile(test.ethGetCodeFn)
		if err != nil {
			fmt.Printf("Error -unable to read file %s, %s\n", test.ethGetCodeFn, err)
			t.Errorf("Error %2d, IO Error on test\n", ii)
			continue
		}
		err = json.Unmarshal(eth, &getResult)
		if err != nil {
			fmt.Printf("Error -unable to read file %s, %s\n", test.binRuntimeFn, err)
			t.Errorf("Error %2d, Unmarshal Error on test\n", ii)
			continue
		}
		ok := VerifyCode(test.solcVersion, string(solcCode), getResult.Result[2:])
		if ok != test.shouldMatch {
			t.Errorf("Error %2d, Invalid test\n", ii)
		}
	}

}

/* vim: set noai ts=4 sw=4: */
