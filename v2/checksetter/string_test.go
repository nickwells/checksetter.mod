package checksetter_test

import (
	"testing"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/checksetter.mod/checksetter"
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
				"ok - one chk func (one int param), string ok"),
			arg: "LenEQ(3)",
			s:   "two",
		},
		{
			ID: testhelper.MkID(
				"ok - one chk func (two int param), string ok"),
			arg: "LenBetween(2, 3)",
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
			arg: `HasPrefix("rc"), LenEQ(5)`,
			s:   "rc001",
		},
		{
			ID:  testhelper.MkID("ok - two chk funcs, string bad"),
			arg: `LenEQ(5), Or(And(HasPrefix("rc"), LenEQ(4)), LenEQ(3))`,
			s:   "rc002",
			strExpErr: testhelper.MkExpErr(
				"the length of the value (5) must equal 4",
				" or the length of the value (5) must equal 3"),
		},
		{
			ID: testhelper.MkID(
				"bad - one chk func (two int param) - unknown func"),
			arg:    "LenBetweenXXX(2, 3)",
			ExpErr: testhelper.MkExpErr("bad function", "LenBetweenXXX"),
		},
		{
			ID: testhelper.MkID(
				"bad - one chk func (two int param) - invalid params"),
			arg: "LenBetween(5, 3)",
			ExpErr: testhelper.MkExpErr(
				"bad function: ",
				"can't make the check.String func: ",
				"Impossible checks passed to StringLenBetween",
			),
		},
		{
			ID: testhelper.MkID(
				"bad - one chk func (two int param) - too many params"),
			arg: "LenBetween(1, 2, 3)",
			ExpErr: testhelper.MkExpErr(
				"bad function: ",
				"can't make the check.String func: ",
				"the call has 3 arguments, it should have 2",
			),
		},
		{
			ID: testhelper.MkID(
				"bad - one chk func (two int param) - too few params"),
			arg: "LenBetween(1)",
			ExpErr: testhelper.MkExpErr(
				"bad function: ",
				"can't make the check.String func: ",
				"the call has 1 arguments, it should have 2",
			),
		},
	}

	for _, tc := range testCases {
		var checks []check.String
		var checker = checksetter.String{Value: &checks}
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
