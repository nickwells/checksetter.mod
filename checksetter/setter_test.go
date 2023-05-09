package checksetter_test

import (
	"fmt"
	"go/ast"
	"testing"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/checksetter.mod/v3/checksetter"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestSetter(t *testing.T) {
	type setterTestType struct{}
	setterTypeName := "checksetter.Setter[checksetter_test.setterTestType]"
	badParser, err := checksetter.MakeParser(
		setterTypeName+"-bad",
		map[string]checksetter.MakerInfo[setterTestType]{})
	if err != nil {
		t.Fatal("couldn't create the bad test Parser: " + err.Error())
	}
	goodParser, err := checksetter.MakeParser(
		setterTypeName+"-good",
		map[string]checksetter.MakerInfo[setterTestType]{
			"OK": {
				MF: func(_ *ast.CallExpr, _ string) (
					check.ValCk[setterTestType], error,
				) {
					return check.ValOK[setterTestType], nil
				},
			},
		},
	)
	if err != nil {
		t.Fatal("couldn't create the good test Parser: " + err.Error())
	}

	value := []check.ValCk[setterTestType]{}

	type param struct {
		val, expCrntVal string
		testhelper.ExpErr
	}
	testCases := []struct {
		testhelper.ID
		testhelper.ExpPanic
		s                checksetter.Setter[setterTestType]
		expCrntValInit   string
		vals             []param
		expAllowedValues string
	}{
		{
			ID: testhelper.MkID("bad value"),
			ExpPanic: testhelper.MkExpPanic(
				setterTypeName +
					" Check failed: the Value to be set is nil"),
		},
		{
			ID: testhelper.MkID("good value, nil Parser"),
			s:  checksetter.Setter[setterTestType]{Value: &value},
			ExpPanic: testhelper.MkExpPanic(
				"The Parser for " + setterTypeName + " has not been set"),
		},
		{
			ID: testhelper.MkID("good value, bad Parser"),
			s: checksetter.Setter[setterTestType]{
				Value:  &value,
				Parser: &badParser,
			},
			ExpPanic: testhelper.MkExpPanic(
				"The Parser for " +
					setterTypeName +
					" can't make any check-funcs"),
		},
		{
			ID: testhelper.MkID("good value, good Parser"),
			s: checksetter.Setter[setterTestType]{
				Value:  &value,
				Parser: &goodParser,
			},
			expCrntValInit: "no checks",
			vals: []param{
				{
					val:        `OK`,
					expCrntVal: `one check: "OK"`,
				},
				{
					val:        `OK, OK`,
					expCrntVal: `2 checks: "OK, OK"`,
				},
				{
					val:        `OK,, OK`,
					expCrntVal: `2 checks: "OK, OK"`,
					ExpErr: testhelper.MkExpErr(
						"expected operand, found ','"),
				},
			},
			expAllowedValues: "a list of" +
				" checksetter.Setter[checksetter_test.setterTestType]-good " +
				"functions separated by ','." +
				" Write the checks as if you were writing code." +
				" The functions recognised are:" +
				"\n\n" +
				"    checksetter.Setter[checksetter_test.setterTestType]-good" +
				" functions:" +
				"\n" +
				"        OK()",
		},
	}

	for _, tc := range testCases {
		panicked, panicVal := testhelper.PanicSafe(
			func() {
				tc.s.CheckSetter(tc.IDStr())
			})
		if testhelper.CheckExpPanic(t, panicked, panicVal, tc) {
			t.Log(tc.IDStr())
			t.Errorf("\t: unexpected panic\n")
		}
		if !panicked {
			testhelper.DiffString(t,
				tc.IDStr(), "Current Value",
				tc.s.CurrentValue(), tc.expCrntValInit)
			for _, p := range tc.vals {
				id := fmt.Sprintf("with val: %q", p.val)
				err := tc.s.SetWithVal("", p.val)
				testhelper.CheckExpErrWithID(t, id, err, p)
				testhelper.DiffString(t,
					tc.IDStr(), id,
					tc.s.CurrentValue(), p.expCrntVal)
			}
			testhelper.DiffString(t,
				tc.IDStr(), "Allowed Values",
				tc.s.AllowedValues(), tc.expAllowedValues)
		}
	}
}
