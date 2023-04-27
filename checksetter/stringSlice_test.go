package checksetter_test

import (
	"testing"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/checksetter.mod/v3/checksetter"
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
			arg:       "Length(EQ(2))",
			ssToCheck: []string{"one", "two"},
		},
		{
			ID: testhelper.MkID(
				"ok - one check function (two int param), slice ok"),
			arg:       "Length(Between(2, 3))",
			ssToCheck: []string{"one", "two"},
		},
		{
			ID: testhelper.MkID(
				"ok - one check function (check string param), slice ok"),
			arg:       `SliceAll(HasPrefix("rc"))`,
			ssToCheck: []string{"rc001", "rc002"},
		},
		{
			ID:        testhelper.MkID("ok - three check functions, slice ok"),
			arg:       `SliceAny(HasPrefix("rc"), "must start with 'rc'"), Length(EQ(2)), NoDups`,
			ssToCheck: []string{"rc001", "rc002"},
		},
		{
			ID:        testhelper.MkID("ok - three check functions, slice bad"),
			arg:       `SliceAny(And(HasPrefix("rc"), Length(EQ(5))), "must be length 5, starting 'rc'"), Length(EQ(3)), NoDups`,
			ssToCheck: []string{"rc001", "rc002", "rc002"},
			sliceErr: testhelper.MkExpErr(
				"duplicate list entries: 1 and 2 are both: rc002"),
		},
		{
			ID: testhelper.MkID(
				"bad - one check function (two int param) - unknown func"),
			arg: "Length(BetweenXXX(2, 3))",
			ExpErr: testhelper.MkExpErr(
				"bad function",
				"can't make the check.StringSlice func: Length(...):",
				"unknown check.Int: BetweenXXX",
			),
		},
		{
			ID: testhelper.MkID(
				"bad - one check function (two int param) - invalid params"),
			arg: "Length(Between(5, 3))",
			ExpErr: testhelper.MkExpErr(
				"bad function: ",
				"can't make the check.StringSlice func: Length(...):",
				"Impossible checks passed to ValBetween",
				"the lower limit (5) must be less than the upper limit (3)",
			),
		},
		{
			ID:        testhelper.MkID("ok - one check function (check.String list)"),
			arg:       `SliceByPos(EQ("RC"), MatchesPattern("[1-9][0-9]*", "numeric"))`,
			ssToCheck: []string{"RC", "9"},
		},
	}

	for _, tc := range testCases {
		var checks []check.StringSlice
		checker := checksetter.StringSlice{Value: &checks}
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
