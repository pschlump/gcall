package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/pschlump/MiscLib"
)

func main() {
	// fmt.Printf("%s%s%s\n", MiscLib.ColorCyan, "yep", MiscLib.ColorReset)
	st := 0
	nErr := 0
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		// fmt.Printf("%sLine:%s%s\n", MiscLib.ColorCyan, line, MiscLib.ColorReset)
		if hasWarning.MatchString(line) {
			// fmt.Printf("%sFound warning ->%s<-%s\n", MiscLib.ColorRed, line, MiscLib.ColorReset)
			st = 3
		}
		if hasError.MatchString(line) {
			// consider outputting in red!
			nErr++
		}
		if st == 0 {
			if nErr > 0 {
				fmt.Printf("%s%s%s\n", MiscLib.ColorRed, line, MiscLib.ColorReset)
			} else {
				fmt.Printf("%s\n", line) // or do something else with line
			}
		} else {
			st--
		}
	}
	if nErr > 0 {
		os.Exit(1)
	}
	os.Exit(0)
}

var hasWarning *regexp.Regexp
var hasError *regexp.Regexp

func init() {
	hasWarning = regexp.MustCompile(" Warning:")
	hasError = regexp.MustCompile(" Error:")
}
