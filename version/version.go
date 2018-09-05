package version

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

// SemanticVersion returns true if the cmp operator is true between the versionString and
// the "to" string.  cmp operator is one of ">", ">=", "<", "<=", "==", "!=".
// Version strings are of the form ".* v[0-9]+.[0-9]+.[0-9]+[ .*]?"
// to strings are of the form [0-9]+.[0-9]+.[0-9]+
func SemanticVersion(versionString, cmp, to string) (rv bool) {

	vOff, err := ExtractVersion(versionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid version string: %s\n", err)
		return false
	}
	vCmp, err := ExtractVersion(" " + to)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid version string: %s\n", err)
		return false
	}

	switch cmp {
	case ">":
		return vOff > vCmp
	case ">=":
		return vOff >= vCmp
	case "<":
		return vOff < vCmp
	case "<=":
		return vOff <= vCmp
	case "==":
		return vOff == vCmp
	case "!=":
		return vOff != vCmp
	}

	return
}

// ExtractVersion picks out the semantic vesion from a version string and returns it as a 64 bit integer.
// It is assumed that versions are limited to 0..10000.
func ExtractVersion(versionString string) (rv int64, err error) {
	xx := versionStr.FindStringSubmatch(versionString)
	if len(xx) == 4 {
		p1, e1 := strconv.ParseInt(xx[1], 10, 32)
		if e1 != nil {
			err = fmt.Errorf("Error with part 1 of version, ->%s<- is not a number, %s, from:%s\n", xx[1], versionString, e1)
			return
		}
		p2, e1 := strconv.ParseInt(xx[2], 10, 32)
		if e1 != nil {
			err = fmt.Errorf("Error with part 2 of version, ->%s<- is not a number, %s, from:%s\n", xx[2], versionString, e1)
			return
		}
		p3, e1 := strconv.ParseInt(xx[3], 10, 32)
		if e1 != nil {
			err = fmt.Errorf("Error with part 3 of version, ->%s<- is not a number, %s, from:%s\n", xx[3], versionString, e1)
			return
		}

		rv = 100000000*p1 + 10000*p2 + p3
	} else {
		err = fmt.Errorf("Failed to match, %s, results = %q", versionString, xx)
	}
	return
}

var versionStr *regexp.Regexp

func init() {
	versionStr = regexp.MustCompile("(?m).*[ \t]*v?([0-9][0-9]*)\\.([0-9][0-9]*)\\.([0-9][0-9]*)")
}
