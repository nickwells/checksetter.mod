package checksetter_test

import (
	"testing"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/checksetter.mod/checksetter"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestChkStringSlice(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		arg       string
		ssToCheck []string
		sliceErr  testhelper.ExpErr
	}{
		{
			ID: testhelper.MkID(
				"ok - one check function (no params), slice ok"),
			arg:       "NoDups",
			ssToCheck: []string{"one", "two"},
		},
		{
			ID: testhelper.MkID(
				"ok - one check function (no params), slice ok"),
			arg:       "NoDups()",
			ssToCheck: []string{"one", "two"},
		},
		{
			ID: testhelper.MkID(
				"ok - one check function (no params), param given"),
			arg:       "NoDups(1)",
			ssToCheck: []string{"one", "two"},
			ExpErr: testhelper.MkExpErr(
				"bad function",
				"can't make the check.StringSlice func:",
				"the call has 1 arguments, it should have 0",
			),
		},
		{
			ID: testhelper.MkID(
				"ok - one check function (one int param), slice ok"),
			arg:       "LenEQ(2)",
			ssToCheck: []string{"one", "two"},
		},
		{
			ID: testhelper.MkID(
				"ok - one check function (two int param), slice ok"),
			arg:       "LenBetween(2, 3)",
			ssToCheck: []string{"one", "two"},
		},
		{
			ID: testhelper.MkID(
				"ok - one check function (check string param), slice ok"),
			arg:       `String(HasPrefix("rc"))`,
			ssToCheck: []string{"rc001", "rc002"},
		},
		{
			ID:        testhelper.MkID("ok - three check functions, slice ok"),
			arg:       `String(HasPrefix("rc")), LenEQ(2), NoDups`,
			ssToCheck: []string{"rc001", "rc002"},
		},
		{
			ID:        testhelper.MkID("ok - three check functions, slice bad"),
			arg:       `String(And(HasPrefix("rc"), LenEQ(5))), LenEQ(3), NoDups`,
			ssToCheck: []string{"rc001", "rc002", "rc002"},
			sliceErr: testhelper.MkExpErr(
				"list entries: 1 and 2 are duplicates, both are: 'rc002'"),
		},
		{
			ID: testhelper.MkID(
				"bad - one check function (two int param) - unknown func"),
			arg: "LenBetweenXXX(2, 3)",
			ExpErr: testhelper.MkExpErr(
				"bad function",
				"LenBetweenXXX",
			),
		},
		{
			ID: testhelper.MkID(
				"bad - one check function (two int param) - invalid params"),
			arg: "LenBetween(5, 3)",
			ExpErr: testhelper.MkExpErr(
				"bad function: ",
				"can't make the check.StringSlice func: ",
				"Impossible checks passed to StringSliceLenBetween",
			),
		},
		{
			ID: testhelper.MkID(
				"ok - one check function (check.String list)"),
			arg:       `StringCheckByPos(Equals("RC"), MatchesPattern("[1-9][0-9]*", "numeric"))`,
			ssToCheck: []string{"RC", "9"},
		},
	}

	for _, tc := range testCases {
		var checks []check.StringSlice
		var checker = checksetter.StringSlice{Value: &checks}
		err := checker.SetWithVal("dummy", tc.arg)
		if testhelper.CheckExpErr(t, err, tc) &&
			err == nil {
			tcID := tc.IDStr()
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
			testhelper.CheckExpErrWithID(t, "SLICE: "+tcID, err, tc.sliceErr)
		}
	}
}
