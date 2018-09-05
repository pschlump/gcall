package version

import "testing"

func Test_ExtractVersion01(t *testing.T) {

	tests := []struct {
		vString   string
		expResult int64
		expError  bool
	}{
		{
			vString:   "Version: 1.2.3",
			expResult: (10000*10000)*1 + 10000*2 + 3,
			expError:  false,
		},
		{
			vString:   "Version: v1.2.3",
			expResult: (10000*10000)*1 + 10000*2 + 3,
			expError:  false,
		},
		{
			vString:   " 2.3.4",
			expResult: (10000*10000)*2 + 10000*3 + 4,
			expError:  false,
		},
		{
			vString:   " v2.3.4",
			expResult: (10000*10000)*2 + 10000*3 + 4,
			expError:  false,
		},
		{
			vString:   "Version: 0.4.7",
			expResult: (10000*10000)*0 + 10000*4 + 7,
			expError:  false,
		},
		{
			vString:   "Version: v0.4.7",
			expResult: (10000*10000)*0 + 10000*4 + 7,
			expError:  false,
		},
		{
			vString:   " 2.3",
			expResult: (10000*10000)*2 + 10000*3 + 4,
			expError:  true,
		},
	}

	for ii, test := range tests {
		rv, err := ExtractVersion(test.vString)
		if err != nil {
			if test.expError == false {
				t.Errorf("Error %2d, Invalid error : %s\n", ii, err)
			}
		} else {
			if rv != test.expResult {
				t.Errorf("Error %2d, Invalid result : expected %d got %d \n", ii, test.expResult, rv)
			}
		}
	}

}

// func SemanticVersion(versionString, cmp, to string) (rv bool) {
func Test_SemanticVersion01(t *testing.T) {

	tests := []struct {
		vString   string
		CmpTo     string
		CmpOp     string
		expResult bool
	}{
		{
			vString:   "My Compiler Version: v1.2.3 - Jan 2, 2018",
			CmpOp:     "==",
			CmpTo:     "1.2.3",
			expResult: true,
		},
		{
			vString:   "My Compiler Version: v1.2.3 - Jan 2, 2018",
			CmpOp:     ">=",
			CmpTo:     "1.0.4",
			expResult: true,
		},
		{
			vString:   "Version: 0.4.21+commit.dfe3193c.Darwin.appleclang",
			CmpOp:     ">=",
			CmpTo:     "0.4.7",
			expResult: true,
		},
	}
	//		expResult: (10000*10000)*2 + 10000*3 + 4,

	for ii, test := range tests {
		rv := SemanticVersion(test.vString, test.CmpOp, test.CmpTo)
		if rv != test.expResult {
			t.Errorf("Error %2d, Invalid result : expected %v got %v \n", ii, test.expResult, rv)
		}
	}

}

/* vim: set noai ts=4 sw=4: */
