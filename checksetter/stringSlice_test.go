package checksetter_test

import (
	"fmt"
	"testing"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/checksetter.mod/checksetter"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestChkStringSlice(t *testing.T) {
	testCases := []struct {
		name               string
		arg                string
		errExpected        bool
		errContains        []string
		ssToCheck          []string
		sliceErrorExpected bool
		sliceErrContains   []string
	}{
		{
			name:      "ok - one check function (no params), slice ok",
			arg:       "NoDups",
			ssToCheck: []string{"one", "two"},
		},
		{
			name:      "ok - one check function (one int param), slice ok",
			arg:       "LenEQ(2)",
			ssToCheck: []string{"one", "two"},
		},
		{
			name:      "ok - one check function (two int param), slice ok",
			arg:       "LenBetween(2, 3)",
			ssToCheck: []string{"one", "two"},
		},
		{
			name:      "ok - one check function (check string param), slice ok",
			arg:       `String(HasPrefix("rc"))`,
			ssToCheck: []string{"rc001", "rc002"},
		},
		{
			name:      "ok - three check functions, slice ok",
			arg:       `String(HasPrefix("rc")), LenEQ(2), NoDups`,
			ssToCheck: []string{"rc001", "rc002"},
		},
		{
			name:               "ok - three check functions, slice bad",
			arg:                `String(And(HasPrefix("rc"), LenEQ(5))), LenEQ(3), NoDups`,
			ssToCheck:          []string{"rc001", "rc002", "rc002"},
			sliceErrorExpected: true,
			sliceErrContains: []string{
				"list entries: 1 and 2 are duplicates, both are: rc002",
			},
		},
		{
			name:        "bad - one check function (two int param) - unknown func",
			arg:         "LenBetweenXXX(2, 3)",
			errExpected: true,
			errContains: []string{
				"bad function",
				"LenBetweenXXX",
			},
		},
		{
			name:        "bad - one check function (two int param) - invalid params",
			arg:         "LenBetween(5, 3)",
			errExpected: true,
			errContains: []string{
				"bad function: ",
				"can't make the check.StringSlice func: ",
				"Impossible checks passed to StringSliceLenBetween",
			},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)
		var checks []check.StringSlice
		var checker = checksetter.StringSlice{Value: &checks}
		err := checker.SetWithVal("dummy", tc.arg)
		if testhelper.CheckError(t, tcID, err,
			tc.errExpected, tc.errContains) &&
			err == nil {
			for i, chk := range checks {
				if chk == nil {
					t.Logf("%s:\n", tcID)
					t.Errorf("\t: nil check found at check slice element %d", i)
					continue
				}
				if err = chk(tc.ssToCheck); err != nil {
					break
				}
			}
			testhelper.CheckError(t, tcID+" (slice checks)",
				err, tc.sliceErrorExpected, tc.sliceErrContains)

		}
	}
}
