package checksetter_test

import (
	"fmt"
	"testing"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/checksetter.mod/checksetter"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestChkString(t *testing.T) {
	testCases := []struct {
		name                string
		arg                 string
		errExpected         bool
		errContains         []string
		strToCheck          string
		stringErrorExpected bool
		stringErrContains   []string
	}{
		{
			name:       "ok - one check function (one int param), string ok",
			arg:        "LenEQ(3)",
			strToCheck: "two",
		},
		{
			name:       "ok - one check function (two int param), string ok",
			arg:        "LenBetween(2, 3)",
			strToCheck: "two",
		},
		{
			name:       "ok - one check function (check string prefix), string ok",
			arg:        `HasPrefix("rc")`,
			strToCheck: "rc001",
		},
		{
			name:       "ok - two check functions, string ok",
			arg:        `HasPrefix("rc"), LenEQ(5)`,
			strToCheck: "rc001",
		},
		{
			name:                "ok - three check functions, string bad",
			arg:                 `Or(And(HasPrefix("rc"), LenEQ(4)), LenEQ(3))`,
			strToCheck:          "rc002",
			stringErrorExpected: true,
			stringErrContains: []string{
				"the length of the value (5) must equal 4",
				" or the length of the value (5) must equal 3",
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
				"can't make the check.String func: ",
				"Impossible checks passed to StringLenBetween",
			},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)
		var checks []check.String
		var checker = checksetter.String{Value: &checks}
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
				if err = chk(tc.strToCheck); err != nil {
					break
				}
			}
			testhelper.CheckError(t, tcID+" (string checks)",
				err, tc.stringErrorExpected, tc.stringErrContains)

		}
	}
}
