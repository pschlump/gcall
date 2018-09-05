package args

import (
	"fmt"
	"testing"

	"github.com/pschlump/godebug"
)

func Test_GzipServer(t *testing.T) {
	aa, err := GetArgs(ArgConfigType{
		Config: map[string]AnArgType{
			"network": AnArgType{
				Name:    "--network",
				Abrev:   "-n",
				Default: "testnet",
			},
		},
	}, 1, []string{"ignoreMe", "--network", "testnet", "list", "--opt1", "bob"})
	if err != nil {
		aa.Usage()
	} else {
		fmt.Printf("results: %s\n", godebug.SVarI(aa))
	}
	// t.Fatalf("Test %d: Could not create HTTP request: %v", ii, err)
	exp := `{
	"Set": {
		"network": "testnet"
	},
	"Remainder": [
		"list",
		"--opt1",
		"bob"
	]
}`
	if got := godebug.SVarI(aa); got != exp {
		t.Fatalf("Test: did not get expected ->%s< got ->%s<-\n", exp, got)
	}

}
