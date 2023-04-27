package checksetter_test

import (
	"testing"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/checksetter.mod/v3/checksetter"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestChkString(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		arg string
		testhelper.ExpErr
		s         string
		strExpErr testhelper.ExpErr
	}{
		{
			ID:  testhelper.MkID("ok - one chk func (no params), string ok"),
			arg: "OK",
			s:   "two",
		},
		{
			ID:  testhelper.MkID("ok - one chk func (no params), string ok"),
			arg: "OK()",
			s:   "two",
		},
		{
			ID:  testhelper.MkID("bad - one chk func (no params), param given"),
			arg: "OK(1)",
			s:   "two",
			ExpErr: testhelper.MkExpErr(
				"bad function: ",
				"can't make the check.String func: ",
				"the call has 1 arguments, it should have 0",
			),
		},
		{
			ID: testhelper.MkID(
				"ok - one chk func (one intCF param), string ok"),
			arg: "Length(EQ(3))",
			s:   "two",
		},
		{
			ID: testhelper.MkID(
				"ok - one chk func (one intCF param), string ok"),
			arg: "Length(Between(2, 3))",
			s:   "two",
		},
		{
			ID: testhelper.MkID(
				"ok - one chk func (check string prefix), string ok"),
			arg: `HasPrefix("rc")`,
			s:   "rc001",
		},
		{
			ID:  testhelper.MkID("ok - two chk funcs, string ok"),
			arg: `HasPrefix("rc"), Length(EQ(5))`,
			s:   "rc001",
		},
		{
			ID:  testhelper.MkID("ok - two chk funcs, string bad"),
			arg: `Length(EQ(5)), Or(And(HasPrefix("rc"), Length(EQ(4))), Length(EQ(3)))`,
			s:   "rc002",
			strExpErr: testhelper.MkExpErr(
				"the length of the string (5) is incorrect",
				"the value (5) must equal 4",
				"the value (5) must equal 3"),
		},
		{
			ID: testhelper.MkID(
				"bad - one chk func (two int param) - unknown func"),
			arg: "Length(BetweenXXX(2, 3))",
			ExpErr: testhelper.MkExpErr(
				"bad function",
				"can't make the check.String func: Length(...):",
				"unknown check.Int: BetweenXXX",
			),
		},
		{
			ID: testhelper.MkID(
				"bad - one chk func (two int param) - invalid params"),
			arg: "Length(Between(5, 3))",
			ExpErr: testhelper.MkExpErr(
				"bad function: ",
				"can't make the check.String func: ",
				"Impossible checks passed to ValBetween",
			),
		},
		{
			ID: testhelper.MkID(
				"bad - one chk func (two int param) - too many params"),
			arg: "Length(Between(1, 2, 3))",
			ExpErr: testhelper.MkExpErr(
				"bad function: ",
				"can't make the check.String func: ",
				"the call has 3 arguments, it should have 2",
			),
		},
		{
			ID: testhelper.MkID(
				"bad - one chk func (two int param) - too few params"),
			arg: "Length(Between(1))",
			ExpErr: testhelper.MkExpErr(
				"bad function: ",
				"can't make the check.String func: ",
				"the call has 1 arguments, it should have 2",
			),
		},
	}

	for _, tc := range testCases {
		var checks []check.String
		checker := checksetter.String{Value: &checks}
		err := checker.SetWithVal("dummy", tc.arg)

		if testhelper.CheckExpErr(t, err, tc) &&
			err == nil {
			for i, chk := range checks {
				if chk == nil {
					t.Logf("%s:\n", tc.IDStr())
					t.Errorf("\t: nil check found at check slice element %d", i)
					continue
				}
				if err = chk(tc.s); err != nil {
					break
				}
			}
			testhelper.CheckExpErrWithID(t, tc.IDStr()+" (string checks)",
				err, tc.strExpErr)
		}
	}
}
