package checksetter_test

import (
	"fmt"
	"go/ast"
	"testing"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/checksetter.mod/v4/checksetter"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestMakeParser(t *testing.T) {
	const (
		checkerNameTMP1 = "TestMakeParser-1"
		checkerNameTMP2 = "TestMakeParser-2"
	)

	type funcArgs struct {
		funcName string
		expArgs  []string
		testhelper.ExpErr
	}

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		checkerName string
		makers      map[string]checksetter.MakerInfo[int]
		expMakers   []string
		expArgs     []funcArgs
	}{
		{
			ID:          testhelper.MkID("good"),
			checkerName: checkerNameTMP1,
			makers:      map[string]checksetter.MakerInfo[int]{},
			expArgs: []funcArgs{
				{
					funcName: "Any",
					ExpErr:   testhelper.MkExpErr(`Unknown maker: "Any"`),
				},
			},
		},
		{
			ID: testhelper.MkID("bad - duplicate"),
			ExpErr: testhelper.MkExpErr(
				`A Parser for "` + checkerNameTMP1 + `" already exists`),
			checkerName: checkerNameTMP1,
			makers:      map[string]checksetter.MakerInfo[int]{},
		},
		{
			ID:          testhelper.MkID("good - with makers"),
			checkerName: checkerNameTMP2,
			makers: map[string]checksetter.MakerInfo[int]{
				"func1": {
					Args: []string{"int", "int"},
					MF: func(_ *ast.CallExpr, _ string) (
						check.ValCk[int], error,
					) {
						return nil, nil
					},
				},
				"func2": {
					Args: []string{"int", "string"},
					MF: func(_ *ast.CallExpr, _ string) (
						check.ValCk[int], error,
					) {
						return nil, nil
					},
				},
			},
			expMakers: []string{"func1", "func2"},
			expArgs: []funcArgs{
				{
					funcName: "Any",
					ExpErr:   testhelper.MkExpErr(`Unknown maker: "Any"`),
				},
				{
					funcName: "func1",
					expArgs:  []string{"int", "int"},
				},
				{
					funcName: "func2",
					expArgs:  []string{"int", "string"},
				},
			},
		},
	}

	for _, tc := range testCases {
		p, err := checksetter.MakeParser(tc.checkerName, tc.makers)
		if testhelper.CheckExpErr(t, err, tc) && err == nil {
			testhelper.DiffString(t, tc.IDStr(), "checker name",
				p.CheckerName(), tc.checkerName)
			testhelper.DiffInt(t, tc.IDStr(), "number of maker funcs",
				len(p.Makers()), len(tc.makers))
			testhelper.DiffStringSlice(t, tc.IDStr(), "maker names",
				p.Makers(), tc.expMakers)

			for _, ea := range tc.expArgs {
				args, err := p.Args(ea.funcName)
				if testhelper.CheckExpErrWithID(t,
					tc.IDStr()+" - args for "+ea.funcName,
					err, ea) && err == nil {
					testhelper.DiffStringSlice(t,
						tc.IDStr(), "args for "+ea.funcName,
						args, ea.expArgs)
				}
			}
		}
	}
}

func TestFindParser(t *testing.T) {
	checkerName := checksetter.IntCheckerName
	expType := "*checksetter.Parser[int]"

	badCheckerName := "nonesuch"

	p, err := checksetter.FindParser[int](checkerName)
	if err != nil {
		t.Errorf("unexpected error when retrieving %q as an int", checkerName)
	} else if actType := fmt.Sprintf("%T", p); actType != expType {
		t.Log("Bad Type")
		t.Log("\t: expected type: " + expType)
		t.Log("\t:   actual type: " + actType)
		t.Errorf("\t: unexpected result when retrieving %q as an int",
			checkerName)
	}

	badTypeMsg := fmt.Sprintf(
		"The Parser for %q is of the wrong type (%s)", checkerName, expType)
	badTypeID := fmt.Sprintf("retrieving %q as a string", checkerName)
	badTypeErr := struct {
		testhelper.ID
		testhelper.ExpErr
	}{
		ID:     testhelper.MkID(badTypeID),
		ExpErr: testhelper.MkExpErr(badTypeMsg),
	}
	badTypePanic := struct {
		testhelper.ID
		testhelper.ExpPanic
	}{
		ID:       testhelper.MkID(badTypeID),
		ExpPanic: testhelper.MkExpPanic(badTypeMsg),
	}

	_, err = checksetter.FindParser[string](checkerName)
	testhelper.CheckExpErr(t, err, badTypeErr)

	panicked, panicVal := testhelper.PanicSafe(func() {
		_ = checksetter.FindParserOrPanic[string](checkerName)
	})
	testhelper.CheckExpPanicError(t, panicked, panicVal, badTypePanic)

	noParserMsg := fmt.Sprintf(
		"There is no Parser registered for %q", badCheckerName)
	noParserID := fmt.Sprintf("retrieving %q", badCheckerName)
	noParserErr := struct {
		testhelper.ID
		testhelper.ExpErr
	}{
		ID:     testhelper.MkID(noParserID),
		ExpErr: testhelper.MkExpErr(noParserMsg),
	}
	noParserPanic := struct {
		testhelper.ID
		testhelper.ExpPanic
	}{
		ID:       testhelper.MkID(noParserID),
		ExpPanic: testhelper.MkExpPanic(noParserMsg),
	}

	_, err = checksetter.FindParser[string](badCheckerName)
	testhelper.CheckExpErr(t, err, noParserErr)

	panicked, panicVal = testhelper.PanicSafe(func() {
		_ = checksetter.FindParserOrPanic[string](badCheckerName)
	})
	testhelper.CheckExpPanicError(t, panicked, panicVal, noParserPanic)
}

func TestParseInt(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		expr        string
		passingVals map[int][]int
		failingVals map[int][]int
		expLen      int
	}{
		{
			ID: testhelper.MkID("bad: no-such name"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" nonesuch is an unknown function"),
			expr: "nonesuch",
		},
		{
			ID: testhelper.MkID("bad: not a named function"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" Syntax error: unexpected call type: *ast.FuncLit"),
			expr: "func(){}()",
		},
		{
			ID:     testhelper.MkID("bad: multi-ck, empty list entry at start"),
			ExpErr: testhelper.MkExpErr("expected operand, found ','"),
			expr:   ",OK",
		},
		{
			ID:     testhelper.MkID("bad: multi-ck, empty list entry - middle"),
			ExpErr: testhelper.MkExpErr("expected operand, found ','"),
			expr:   "OK,,OK",
		},
		{
			ID:          testhelper.MkID("no-params: good: just name"),
			expr:        "OK",
			passingVals: map[int][]int{0: {1, -1, 0, 9999, -9999}},
			expLen:      1,
		},
		{
			ID:          testhelper.MkID("no-params: good: just name call"),
			expr:        "OK()",
			passingVals: map[int][]int{0: {1, -1, 0, 9999, -9999}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("no-params: bad: with params"),
			ExpErr: testhelper.MkExpErr(
				"Can't make int-checker function: OK()",
				"the call has 1 arguments, it should have 0"),
			expr:   "OK(42)",
			expLen: 1,
		},
		{
			ID:   testhelper.MkID("no-params: good: multi"),
			expr: "OK, OK()",
			passingVals: map[int][]int{
				0: {1, -1, 0, 9999, -9999},
				1: {42},
			},
			expLen: 2,
		},
		{
			ID:          testhelper.MkID("1 int param: good: EQ"),
			expr:        "EQ(1)",
			passingVals: map[int][]int{0: {1}},
			failingVals: map[int][]int{0: {0, 2, -1}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("1 int param: bad: EQ, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" EQ(int):" +
				" \"`hello`\" isn't an INT, it's a STRING"),
			expr: "EQ(`hello`)",
		},
		{
			ID: testhelper.MkID("1 int param: bad: EQ, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" EQ(int):" +
				" the call has 0 arguments, it should have 1"),
			expr: "EQ()",
		},
		{
			ID: testhelper.MkID("1 int param: bad: EQ, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" EQ(int):" +
				" the call has 2 arguments, it should have 1"),
			expr: "EQ(1, 2)",
		},
		{
			ID:          testhelper.MkID("1 int param: good: GT"),
			expr:        "GT(1)",
			passingVals: map[int][]int{0: {2, 99}},
			failingVals: map[int][]int{0: {-1, 0, 1}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("1 int param: bad: GT, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" GT(int):" +
				" \"`hello`\" isn't an INT, it's a STRING"),
			expr: "GT(`hello`)",
		},
		{
			ID: testhelper.MkID("1 int param: bad: GT, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" GT(int):" +
				" the call has 0 arguments, it should have 1"),
			expr: "GT()",
		},
		{
			ID: testhelper.MkID("1 int param: bad: GT, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" GT(int):" +
				" the call has 2 arguments, it should have 1"),
			expr: "GT(1, 2)",
		},
		{
			ID:          testhelper.MkID("1 int param: good: GE"),
			expr:        "GE(1)",
			passingVals: map[int][]int{0: {1, 2, 99}},
			failingVals: map[int][]int{0: {-99, -1, 0}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("1 int param: bad: GE, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" GE(int):" +
				" \"`hello`\" isn't an INT, it's a STRING"),
			expr: "GE(`hello`)",
		},
		{
			ID: testhelper.MkID("1 int param: bad: GE, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" GE(int):" +
				" the call has 0 arguments, it should have 1"),
			expr: "GE()",
		},
		{
			ID: testhelper.MkID("1 int param: bad: GE, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" GE(int):" +
				" the call has 2 arguments, it should have 1"),
			expr: "GE(1, 2)",
		},
		{
			ID:          testhelper.MkID("1 int param: good: LT"),
			expr:        "LT(1)",
			passingVals: map[int][]int{0: {-99, -1, 0}},
			failingVals: map[int][]int{0: {1, 2, 99}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("1 int param: bad: LT, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" LT(int):" +
				" \"`hello`\" isn't an INT, it's a STRING"),
			expr: "LT(`hello`)",
		},
		{
			ID: testhelper.MkID("1 int param: bad: LT, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" LT(int):" +
				" the call has 0 arguments, it should have 1"),
			expr: "LT()",
		},
		{
			ID: testhelper.MkID("1 int param: bad: LT, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" LT(int):" +
				" the call has 2 arguments, it should have 1"),
			expr: "LT(1, 2)",
		},
		{
			ID:          testhelper.MkID("1 int param: good: LE"),
			expr:        "LE(1)",
			passingVals: map[int][]int{0: {1, 0, -1, -99}},
			failingVals: map[int][]int{0: {2, 99}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("1 int param: bad: LE, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" LE(int):" +
				" \"`hello`\" isn't an INT, it's a STRING"),
			expr: "LE(`hello`)",
		},
		{
			ID: testhelper.MkID("1 int param: bad: LE, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" LE(int):" +
				" the call has 0 arguments, it should have 1"),
			expr: "LE()",
		},
		{
			ID: testhelper.MkID("1 int param: bad: LE, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" LE(int):" +
				" the call has 2 arguments, it should have 1"),
			expr: "LE(1, 2)",
		},
		{
			ID:          testhelper.MkID("1 int param: good: Divides"),
			expr:        "Divides(60)",
			passingVals: map[int][]int{0: {1, 2, 3, 4, 5, 6, 12, -2}},
			failingVals: map[int][]int{0: {7, 8, 9}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("1 int param: bad: Divides, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" Divides(int):" +
				" \"`hello`\" isn't an INT, it's a STRING"),
			expr: "Divides(`hello`)",
		},
		{
			ID: testhelper.MkID("1 int param: bad: Divides, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" Divides(int):" +
				" the call has 0 arguments, it should have 1"),
			expr: "Divides()",
		},
		{
			ID: testhelper.MkID("1 int param: bad: Divides, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" Divides(int):" +
				" the call has 2 arguments, it should have 1"),
			expr: "Divides(1, 2)",
		},
		{
			ID:          testhelper.MkID("1 int param: good: IsAMultiple"),
			expr:        "IsAMultiple(10)",
			passingVals: map[int][]int{0: {10, 20, 30, -10}},
			failingVals: map[int][]int{0: {11, 19, 15}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("1 int param: bad: IsAMultiple, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" IsAMultiple(int):" +
				" \"`hello`\" isn't an INT, it's a STRING"),
			expr: "IsAMultiple(`hello`)",
		},
		{
			ID: testhelper.MkID("1 int param: bad: IsAMultiple, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" IsAMultiple(int):" +
				" the call has 0 arguments, it should have 1"),
			expr: "IsAMultiple()",
		},
		{
			ID: testhelper.MkID("1 int param: bad: IsAMultiple, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" IsAMultiple(int):" +
				" the call has 2 arguments, it should have 1"),
			expr: "IsAMultiple(1, 2)",
		},
		{
			ID:   testhelper.MkID("1 int param: good: multi-checks"),
			expr: "IsAMultiple(10), Divides(60), GT(10)",
			passingVals: map[int][]int{
				0: {20, 30, 70},
				1: {20, 30, 15},
				2: {20, 30, 12},
			},
			failingVals: map[int][]int{
				0: {2, 3, 5, 12},
				1: {7, 70},
				2: {2, 3, 4, 5},
			},
			expLen: 3,
		},
		{
			ID:          testhelper.MkID("2 int param: good: Between"),
			expr:        "Between(10, 12)",
			passingVals: map[int][]int{0: {10, 11, 12}},
			failingVals: map[int][]int{0: {9, 13}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("2 int param: bad: Between, bad args (1st)"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" Between(int, int):" +
				" \"`hello`\" isn't an INT, it's a STRING"),
			expr: "Between(`hello`, 1)",
		},
		{
			ID: testhelper.MkID("2 int param: bad: Between, bad args (2nd)"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" Between(int, int):" +
				" \"`hello`\" isn't an INT, it's a STRING"),
			expr: "Between(1, `hello`)",
		},
		{
			ID: testhelper.MkID("2 int param: bad: Between, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" Between(int, int):" +
				" Impossible checks passed to ValBetween:" +
				" the lower limit (12) must be less than the upper limit (10)"),
			expr: "Between(12, 10)",
		},
		{
			ID: testhelper.MkID("2 int param: bad: Between, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" Between(int, int):" +
				" the call has 1 arguments, it should have 2"),
			expr: "Between(9)",
		},
		{
			ID: testhelper.MkID("2 int param: bad: Between, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" Between(int, int):" +
				" the call has 3 arguments, it should have 2"),
			expr: "Between(9, 10, 11)",
		},
		{
			ID:          testhelper.MkID("int-ckr, str param: good: Not"),
			expr:        "Not(EQ(10), `not 10`)",
			passingVals: map[int][]int{0: {9, 11}},
			failingVals: map[int][]int{0: {10}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("int-ckr, str param: bad: Not, bad args (1st)"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" Not(int-checker, string):" +
				" can't convert argument 0 to int-checker:" +
				" unexpected type: *ast.BasicLit"),
			expr: "Not(1, `hello`)",
		},
		{
			ID: testhelper.MkID("int-ckr, str param: bad: Not, bad args (2nd)"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" Not(int-checker, string):" +
				" \"1\" isn't a STRING, it's a INT"),
			expr: "Not(EQ(1), 1)",
		},
		{
			ID: testhelper.MkID("int-ckr, str param: bad: Not, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" Not(int-checker, string):" +
				" the call has 3 arguments, it should have 2"),
			expr: "Not(EQ(1), `Hello`, `World`)",
		},
		{
			ID: testhelper.MkID("int-ckr, str param: bad: Not, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" Not(int-checker, string):" +
				" the call has 1 arguments, it should have 2"),
			expr: "Not(EQ(1))",
		},
		{
			ID:          testhelper.MkID("...int-ckr param: good: And: 1 CF"),
			expr:        "And(EQ(10))",
			passingVals: map[int][]int{0: {10}},
			failingVals: map[int][]int{0: {9, 11}},
			expLen:      1,
		},
		{
			ID:          testhelper.MkID("...int-ckr param: good: And: 3 CF"),
			expr:        "And(EQ(10), OK, LT(11))",
			passingVals: map[int][]int{0: {10}},
			failingVals: map[int][]int{0: {9, 11}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("...int-ckr param: bad: And: not CF"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" And(..., int-checker):" +
				" can't convert argument 1 to int-checker:" +
				" nonesuch is an unknown function"),
			expr: "And(EQ(10), nonesuch, LT(11))",
		},
		{
			ID:          testhelper.MkID("...int-ckr param: good: Or: 1 CF"),
			expr:        "Or(EQ(10))",
			passingVals: map[int][]int{0: {10}},
			failingVals: map[int][]int{0: {9, 11}},
			expLen:      1,
		},
		{
			ID:          testhelper.MkID("...int-ckr param: good: Or: 3 CF"),
			expr:        "Or(EQ(200), LT(11))",
			passingVals: map[int][]int{0: {10, 9, 8, 200}},
			failingVals: map[int][]int{0: {199, 11}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("...int-ckr param: bad: Or: not CF"),
			ExpErr: testhelper.MkExpErr("Can't make int-checker function:" +
				" Or(..., int-checker):" +
				" can't convert argument 1 to int-checker:" +
				" nonesuch is an unknown function"),
			expr: "Or(EQ(10), nonesuch, LT(11))",
		},
	}

	parser := checksetter.FindParserOrPanic[int](
		checksetter.IntCheckerName)
	for _, tc := range testCases {
		vcs, err := parser.Parse(tc.expr)
		if testhelper.CheckExpErr(t, err, tc) &&
			err == nil &&
			testhelper.DiffInt(t, tc.IDStr(), "number of ValCk funcs",
				len(vcs), tc.expLen) {
			for vcIdx, vc := range vcs {
				for _, pVal := range tc.passingVals[vcIdx] {
					if err = vc(pVal); err != nil {
						t.Log(tc.IDStr())
						t.Logf(
							"\t: error when checking %v with ValCk: %d",
							pVal, vcIdx)
						t.Error("\t: Bad check")
					}
				}

				for _, fVal := range tc.failingVals[vcIdx] {
					if err = vc(fVal); err == nil {
						t.Log(tc.IDStr())
						t.Logf(
							"\t: missing error when checking %v with ValCk: %d",
							fVal, vcIdx)
						t.Error("\t: Bad check")
					}
				}
			}
		}
	}
}

func TestParseInt64(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		expr        string
		passingVals map[int][]int64
		failingVals map[int][]int64
		expLen      int
	}{
		{
			ID: testhelper.MkID("bad: no-such name"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" nonesuch is an unknown function"),
			expr: "nonesuch",
		},
		{
			ID: testhelper.MkID("bad: not a named function"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" Syntax error: unexpected call type: *ast.FuncLit"),
			expr: "func(){}()",
		},
		{
			ID:     testhelper.MkID("bad: multi-ck, empty list entry at start"),
			ExpErr: testhelper.MkExpErr("expected operand, found ','"),
			expr:   ",OK",
		},
		{
			ID:     testhelper.MkID("bad: multi-ck, empty list entry - middle"),
			ExpErr: testhelper.MkExpErr("expected operand, found ','"),
			expr:   "OK,,OK",
		},
		{
			ID:          testhelper.MkID("no-params: good: just name"),
			expr:        "OK",
			passingVals: map[int][]int64{0: {1, -1, 0, 9999, -9999}},
			expLen:      1,
		},
		{
			ID:          testhelper.MkID("no-params: good: just name call"),
			expr:        "OK()",
			passingVals: map[int][]int64{0: {1, -1, 0, 9999, -9999}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("no-params: bad: with params"),
			ExpErr: testhelper.MkExpErr(
				"Can't make int64-checker function: OK()",
				"the call has 1 arguments, it should have 0"),
			expr:   "OK(42)",
			expLen: 1,
		},
		{
			ID:   testhelper.MkID("no-params: good: multi"),
			expr: "OK, OK()",
			passingVals: map[int][]int64{
				0: {1, -1, 0, 9999, -9999},
				1: {42},
			},
			expLen: 2,
		},
		{
			ID:          testhelper.MkID("1 int64 param: good: EQ"),
			expr:        "EQ(1)",
			passingVals: map[int][]int64{0: {1}},
			failingVals: map[int][]int64{0: {0, 2, -1}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("1 int64 param: bad: EQ, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" EQ(int64):" +
				" \"`hello`\" isn't an INT, it's a STRING"),
			expr: "EQ(`hello`)",
		},
		{
			ID: testhelper.MkID("1 int64 param: bad: EQ, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" EQ(int64):" +
				" the call has 0 arguments, it should have 1"),
			expr: "EQ()",
		},
		{
			ID: testhelper.MkID("1 int64 param: bad: EQ, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" EQ(int64):" +
				" the call has 2 arguments, it should have 1"),
			expr: "EQ(1, 2)",
		},
		{
			ID:          testhelper.MkID("1 int64 param: good: GT"),
			expr:        "GT(1)",
			passingVals: map[int][]int64{0: {2, 99}},
			failingVals: map[int][]int64{0: {-1, 0, 1}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("1 int64 param: bad: GT, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" GT(int64):" +
				" \"`hello`\" isn't an INT, it's a STRING"),
			expr: "GT(`hello`)",
		},
		{
			ID: testhelper.MkID("1 int64 param: bad: GT, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" GT(int64):" +
				" the call has 0 arguments, it should have 1"),
			expr: "GT()",
		},
		{
			ID: testhelper.MkID("1 int64 param: bad: GT, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" GT(int64):" +
				" the call has 2 arguments, it should have 1"),
			expr: "GT(1, 2)",
		},
		{
			ID:          testhelper.MkID("1 int64 param: good: GE"),
			expr:        "GE(1)",
			passingVals: map[int][]int64{0: {1, 2, 99}},
			failingVals: map[int][]int64{0: {-99, -1, 0}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("1 int64 param: bad: GE, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" GE(int64):" +
				" \"`hello`\" isn't an INT, it's a STRING"),
			expr: "GE(`hello`)",
		},
		{
			ID: testhelper.MkID("1 int64 param: bad: GE, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" GE(int64):" +
				" the call has 0 arguments, it should have 1"),
			expr: "GE()",
		},
		{
			ID: testhelper.MkID("1 int64 param: bad: GE, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" GE(int64):" +
				" the call has 2 arguments, it should have 1"),
			expr: "GE(1, 2)",
		},
		{
			ID:          testhelper.MkID("1 int64 param: good: LT"),
			expr:        "LT(1)",
			passingVals: map[int][]int64{0: {-99, -1, 0}},
			failingVals: map[int][]int64{0: {1, 2, 99}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("1 int64 param: bad: LT, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" LT(int64):" +
				" \"`hello`\" isn't an INT, it's a STRING"),
			expr: "LT(`hello`)",
		},
		{
			ID: testhelper.MkID("1 int64 param: bad: LT, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" LT(int64):" +
				" the call has 0 arguments, it should have 1"),
			expr: "LT()",
		},
		{
			ID: testhelper.MkID("1 int64 param: bad: LT, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" LT(int64):" +
				" the call has 2 arguments, it should have 1"),
			expr: "LT(1, 2)",
		},
		{
			ID:          testhelper.MkID("1 int64 param: good: LE"),
			expr:        "LE(1)",
			passingVals: map[int][]int64{0: {1, 0, -1, -99}},
			failingVals: map[int][]int64{0: {2, 99}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("1 int64 param: bad: LE, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" LE(int64):" +
				" \"`hello`\" isn't an INT, it's a STRING"),
			expr: "LE(`hello`)",
		},
		{
			ID: testhelper.MkID("1 int64 param: bad: LE, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" LE(int64):" +
				" the call has 0 arguments, it should have 1"),
			expr: "LE()",
		},
		{
			ID: testhelper.MkID("1 int64 param: bad: LE, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" LE(int64):" +
				" the call has 2 arguments, it should have 1"),
			expr: "LE(1, 2)",
		},
		{
			ID:          testhelper.MkID("1 int64 param: good: Divides"),
			expr:        "Divides(60)",
			passingVals: map[int][]int64{0: {1, 2, 3, 4, 5, 6, 12, -2}},
			failingVals: map[int][]int64{0: {7, 8, 9}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("1 int64 param: bad: Divides, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" Divides(int64):" +
				" \"`hello`\" isn't an INT, it's a STRING"),
			expr: "Divides(`hello`)",
		},
		{
			ID: testhelper.MkID("1 int64 param: bad: Divides, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" Divides(int64):" +
				" the call has 0 arguments, it should have 1"),
			expr: "Divides()",
		},
		{
			ID: testhelper.MkID("1 int64 param: bad: Divides, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" Divides(int64):" +
				" the call has 2 arguments, it should have 1"),
			expr: "Divides(1, 2)",
		},
		{
			ID:          testhelper.MkID("1 int64 param: good: IsAMultiple"),
			expr:        "IsAMultiple(10)",
			passingVals: map[int][]int64{0: {10, 20, 30, -10}},
			failingVals: map[int][]int64{0: {11, 19, 15}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("1 int64 param: bad: IsAMultiple, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" IsAMultiple(int64):" +
				" \"`hello`\" isn't an INT, it's a STRING"),
			expr: "IsAMultiple(`hello`)",
		},
		{
			ID: testhelper.MkID(
				"1 int64 param: bad: IsAMultiple, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" IsAMultiple(int64):" +
				" the call has 0 arguments, it should have 1"),
			expr: "IsAMultiple()",
		},
		{
			ID: testhelper.MkID(
				"1 int64 param: bad: IsAMultiple, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" IsAMultiple(int64):" +
				" the call has 2 arguments, it should have 1"),
			expr: "IsAMultiple(1, 2)",
		},
		{
			ID:   testhelper.MkID("1 int64 param: good: multi-checks"),
			expr: "IsAMultiple(10), Divides(60), GT(10)",
			passingVals: map[int][]int64{
				0: {20, 30, 70},
				1: {20, 30, 15},
				2: {20, 30, 12},
			},
			failingVals: map[int][]int64{
				0: {2, 3, 5, 12},
				1: {7, 70},
				2: {2, 3, 4, 5},
			},
			expLen: 3,
		},
		{
			ID:          testhelper.MkID("2 int64 param: good: Between"),
			expr:        "Between(10, 12)",
			passingVals: map[int][]int64{0: {10, 11, 12}},
			failingVals: map[int][]int64{0: {9, 13}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("2 int64 param: bad: Between, bad args (1st)"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" Between(int64, int64):" +
				" \"`hello`\" isn't an INT, it's a STRING"),
			expr: "Between(`hello`, 1)",
		},
		{
			ID: testhelper.MkID("2 int64 param: bad: Between, bad args (2nd)"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" Between(int64, int64):" +
				" \"`hello`\" isn't an INT, it's a STRING"),
			expr: "Between(1, `hello`)",
		},
		{
			ID: testhelper.MkID("2 int64 param: bad: Between, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" Between(int64, int64):" +
				" Impossible checks passed to ValBetween:" +
				" the lower limit (12) must be less than the upper limit (10)"),
			expr: "Between(12, 10)",
		},
		{
			ID: testhelper.MkID("2 int64 param: bad: Between, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" Between(int64, int64):" +
				" the call has 1 arguments, it should have 2"),
			expr: "Between(9)",
		},
		{
			ID: testhelper.MkID("2 int64 param: bad: Between, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" Between(int64, int64):" +
				" the call has 3 arguments, it should have 2"),
			expr: "Between(9, 10, 11)",
		},
		{
			ID:          testhelper.MkID("int64-ckr, str param: good: Not"),
			expr:        "Not(EQ(10), `not 10`)",
			passingVals: map[int][]int64{0: {9, 11}},
			failingVals: map[int][]int64{0: {10}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID(
				"int64-ckr, str param: bad: Not, bad args (1st)"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" Not(int64-checker, string):" +
				" can't convert argument 0 to int64-checker:" +
				" unexpected type: *ast.BasicLit"),
			expr: "Not(1, `hello`)",
		},
		{
			ID: testhelper.MkID(
				"int64-ckr, str param: bad: Not, bad args (2nd)"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" Not(int64-checker, string):" +
				" \"1\" isn't a STRING, it's a INT"),
			expr: "Not(EQ(1), 1)",
		},
		{
			ID: testhelper.MkID(
				"int64-ckr, str param: bad: Not, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" Not(int64-checker, string):" +
				" the call has 3 arguments, it should have 2"),
			expr: "Not(EQ(1), `Hello`, `World`)",
		},
		{
			ID: testhelper.MkID("int64-ckr, str param: bad: Not, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" Not(int64-checker, string):" +
				" the call has 1 arguments, it should have 2"),
			expr: "Not(EQ(1))",
		},
		{
			ID:          testhelper.MkID("...int64-ckr param: good: And: 1 CF"),
			expr:        "And(EQ(10))",
			passingVals: map[int][]int64{0: {10}},
			failingVals: map[int][]int64{0: {9, 11}},
			expLen:      1,
		},
		{
			ID:          testhelper.MkID("...int64-ckr param: good: And: 3 CF"),
			expr:        "And(EQ(10), OK, LT(11))",
			passingVals: map[int][]int64{0: {10}},
			failingVals: map[int][]int64{0: {9, 11}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("...int64-ckr param: bad: And: not CF"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" And(..., int64-checker):" +
				" can't convert argument 1 to int64-checker:" +
				" nonesuch is an unknown function"),
			expr: "And(EQ(10), nonesuch, LT(11))",
		},
		{
			ID:          testhelper.MkID("...int64-ckr param: good: Or: 1 CF"),
			expr:        "Or(EQ(10))",
			passingVals: map[int][]int64{0: {10}},
			failingVals: map[int][]int64{0: {9, 11}},
			expLen:      1,
		},
		{
			ID:          testhelper.MkID("...int64-ckr param: good: Or: 3 CF"),
			expr:        "Or(EQ(200), LT(11))",
			passingVals: map[int][]int64{0: {10, 9, 8, 200}},
			failingVals: map[int][]int64{0: {199, 11}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("...int64-ckr param: bad: Or: not CF"),
			ExpErr: testhelper.MkExpErr("Can't make int64-checker function:" +
				" Or(..., int64-checker):" +
				" can't convert argument 1 to int64-checker:" +
				" nonesuch is an unknown function"),
			expr: "Or(EQ(10), nonesuch, LT(11))",
		},
	}

	parser := checksetter.FindParserOrPanic[int64](
		checksetter.Int64CheckerName)
	for _, tc := range testCases {
		vcs, err := parser.Parse(tc.expr)
		if testhelper.CheckExpErr(t, err, tc) &&
			err == nil &&
			testhelper.DiffInt(t, tc.IDStr(), "number of ValCk funcs",
				len(vcs), tc.expLen) {
			for vcIdx, vc := range vcs {
				for _, pVal := range tc.passingVals[vcIdx] {
					if err = vc(pVal); err != nil {
						t.Log(tc.IDStr())
						t.Logf(
							"\t: error when checking %v with ValCk: %d",
							pVal, vcIdx)
						t.Error("\t: Bad check")
					}
				}

				for _, fVal := range tc.failingVals[vcIdx] {
					if err = vc(fVal); err == nil {
						t.Log(tc.IDStr())
						t.Logf(
							"\t: missing error when checking %v with ValCk: %d",
							fVal, vcIdx)
						t.Error("\t: Bad check")
					}
				}
			}
		}
	}
}

func TestParseFloat64(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		expr        string
		passingVals map[int][]float64
		failingVals map[int][]float64
		expLen      int
	}{
		{
			ID: testhelper.MkID("bad: no-such name"),
			ExpErr: testhelper.MkExpErr("Can't make float64-checker function:" +
				" nonesuch is an unknown function"),
			expr: "nonesuch",
		},
		{
			ID: testhelper.MkID("bad: not a named function"),
			ExpErr: testhelper.MkExpErr("Can't make float64-checker function:" +
				" Syntax error: unexpected call type: *ast.FuncLit"),
			expr: "func(){}()",
		},
		{
			ID:     testhelper.MkID("bad: multi-ck, empty list entry at start"),
			ExpErr: testhelper.MkExpErr("expected operand, found ','"),
			expr:   ",OK",
		},
		{
			ID:     testhelper.MkID("bad: multi-ck, empty list entry - middle"),
			ExpErr: testhelper.MkExpErr("expected operand, found ','"),
			expr:   "OK,,OK",
		},
		{
			ID:          testhelper.MkID("no-params: good: just name"),
			expr:        "OK",
			passingVals: map[int][]float64{0: {1, -1, 0, 9999, -9999}},
			expLen:      1,
		},
		{
			ID:          testhelper.MkID("no-params: good: just name call"),
			expr:        "OK()",
			passingVals: map[int][]float64{0: {1, -1, 0, 9999, -9999}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("no-params: bad: with params"),
			ExpErr: testhelper.MkExpErr(
				"Can't make float64-checker function: OK()",
				"the call has 1 arguments, it should have 0"),
			expr:   "OK(42)",
			expLen: 1,
		},
		{
			ID:   testhelper.MkID("no-params: good: multi"),
			expr: "OK, OK()",
			passingVals: map[int][]float64{
				0: {1, -1, 0, 9999, -9999},
				1: {42},
			},
			expLen: 2,
		},
		{
			ID:          testhelper.MkID("1 float64 param: good: GT"),
			expr:        "GT(1)",
			passingVals: map[int][]float64{0: {2, 99}},
			failingVals: map[int][]float64{0: {-1, 0, 1}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("1 float64 param: bad: GT, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make float64-checker function:" +
				" GT(float64):" +
				" \"`hello`\" isn't a FLOAT/INT, it's a STRING"),
			expr: "GT(`hello`)",
		},
		{
			ID: testhelper.MkID("1 float64 param: bad: GT, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make float64-checker function:" +
				" GT(float64):" +
				" the call has 0 arguments, it should have 1"),
			expr: "GT()",
		},
		{
			ID: testhelper.MkID("1 float64 param: bad: GT, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make float64-checker function:" +
				" GT(float64):" +
				" the call has 2 arguments, it should have 1"),
			expr: "GT(1, 2)",
		},
		{
			ID:          testhelper.MkID("1 float64 param: good: GE"),
			expr:        "GE(1)",
			passingVals: map[int][]float64{0: {1, 2, 99}},
			failingVals: map[int][]float64{0: {-99, -1, 0}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("1 float64 param: bad: GE, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make float64-checker function:" +
				" GE(float64):" +
				" \"`hello`\" isn't a FLOAT/INT, it's a STRING"),
			expr: "GE(`hello`)",
		},
		{
			ID: testhelper.MkID("1 float64 param: bad: GE, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make float64-checker function:" +
				" GE(float64):" +
				" the call has 0 arguments, it should have 1"),
			expr: "GE()",
		},
		{
			ID: testhelper.MkID("1 float64 param: bad: GE, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make float64-checker function:" +
				" GE(float64):" +
				" the call has 2 arguments, it should have 1"),
			expr: "GE(1, 2)",
		},
		{
			ID:          testhelper.MkID("1 float64 param: good: LT"),
			expr:        "LT(1)",
			passingVals: map[int][]float64{0: {-99, -1, 0}},
			failingVals: map[int][]float64{0: {1, 2, 99}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("1 float64 param: bad: LT, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make float64-checker function:" +
				" LT(float64):" +
				" \"`hello`\" isn't a FLOAT/INT, it's a STRING"),
			expr: "LT(`hello`)",
		},
		{
			ID: testhelper.MkID("1 float64 param: bad: LT, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make float64-checker function:" +
				" LT(float64):" +
				" the call has 0 arguments, it should have 1"),
			expr: "LT()",
		},
		{
			ID: testhelper.MkID("1 float64 param: bad: LT, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make float64-checker function:" +
				" LT(float64):" +
				" the call has 2 arguments, it should have 1"),
			expr: "LT(1, 2)",
		},
		{
			ID:          testhelper.MkID("1 float64 param: good: LE"),
			expr:        "LE(1)",
			passingVals: map[int][]float64{0: {1, 0, -1, -99}},
			failingVals: map[int][]float64{0: {2, 99}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("1 float64 param: bad: LE, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make float64-checker function:" +
				" LE(float64):" +
				" \"`hello`\" isn't a FLOAT/INT, it's a STRING"),
			expr: "LE(`hello`)",
		},
		{
			ID: testhelper.MkID("1 float64 param: bad: LE, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make float64-checker function:" +
				" LE(float64):" +
				" the call has 0 arguments, it should have 1"),
			expr: "LE()",
		},
		{
			ID: testhelper.MkID("1 float64 param: bad: LE, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make float64-checker function:" +
				" LE(float64):" +
				" the call has 2 arguments, it should have 1"),
			expr: "LE(1, 2)",
		},
		{
			ID:   testhelper.MkID("1 float64 param: good: multi-checks"),
			expr: "LT(60), GT(10)",
			passingVals: map[int][]float64{
				0: {20, 30, 40},
				1: {20, 30, 11},
			},
			failingVals: map[int][]float64{
				0: {60, 61},
				1: {7, -70},
			},
			expLen: 2,
		},
		{
			ID:          testhelper.MkID("2 float64 param: good: Between"),
			expr:        "Between(10, 12)",
			passingVals: map[int][]float64{0: {10, 11, 12}},
			failingVals: map[int][]float64{0: {9, 13}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID(
				"2 float64 param: bad: Between, bad args (1st)"),
			ExpErr: testhelper.MkExpErr("Can't make float64-checker function:" +
				" Between(float64, float64):" +
				" \"`hello`\" isn't a FLOAT/INT, it's a STRING"),
			expr: "Between(`hello`, 1)",
		},
		{
			ID: testhelper.MkID(
				"2 float64 param: bad: Between, bad args (2nd)"),
			ExpErr: testhelper.MkExpErr("Can't make float64-checker function:" +
				" Between(float64, float64):" +
				" \"`hello`\" isn't a FLOAT/INT, it's a STRING"),
			expr: "Between(1, `hello`)",
		},
		{
			ID: testhelper.MkID("2 float64 param: bad: Between, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make float64-checker function:" +
				" Between(float64, float64):" +
				" Impossible checks passed to ValBetween:" +
				" the lower limit (12) must be less than the upper limit (10)"),
			expr: "Between(12, 10)",
		},
		{
			ID: testhelper.MkID("2 float64 param: bad: Between, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make float64-checker function:" +
				" Between(float64, float64):" +
				" the call has 1 arguments, it should have 2"),
			expr: "Between(9)",
		},
		{
			ID: testhelper.MkID("2 float64 param: bad: Between, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make float64-checker function:" +
				" Between(float64, float64):" +
				" the call has 3 arguments, it should have 2"),
			expr: "Between(9, 10, 11)",
		},
		{
			ID:          testhelper.MkID("float64-ckr, str param: good: Not"),
			expr:        "Not(GE(10), `not GE(10)`)",
			passingVals: map[int][]float64{0: {10, 9, 8}},
			failingVals: map[int][]float64{0: {11, 12}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID(
				"float64-ckr, str param: bad: Not, bad args (1st)"),
			ExpErr: testhelper.MkExpErr("Can't make float64-checker function:" +
				" Not(float64-checker, string):" +
				" can't convert argument 0 to float64-checker:" +
				" unexpected type: *ast.BasicLit"),
			expr: "Not(1, `hello`)",
		},
		{
			ID: testhelper.MkID(
				"float64-ckr, str param: bad: Not, bad args (2nd)"),
			ExpErr: testhelper.MkExpErr("Can't make float64-checker function:" +
				" Not(float64-checker, string):" +
				" \"1\" isn't a STRING, it's a INT"),
			expr: "Not(GE(1), 1)",
		},
		{
			ID: testhelper.MkID(
				"float64-ckr, str param: bad: Not, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make float64-checker function:" +
				" Not(float64-checker, string):" +
				" the call has 3 arguments, it should have 2"),
			expr: "Not(GE(1), `Hello`, `World`)",
		},
		{
			ID: testhelper.MkID(
				"float64-ckr, str param: bad: Not, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make float64-checker function:" +
				" Not(float64-checker, string):" +
				" the call has 1 arguments, it should have 2"),
			expr: "Not(GE(1))",
		},
		{
			ID: testhelper.MkID(
				"...float64-ckr param: good: And: 1 CF"),
			expr:        "And(GE(10))",
			passingVals: map[int][]float64{0: {10}},
			failingVals: map[int][]float64{0: {9, 11}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID(
				"...float64-ckr param: good: And: 3 CF"),
			expr:        "And(GE(10), OK, LT(11))",
			passingVals: map[int][]float64{0: {10}},
			failingVals: map[int][]float64{0: {9, 11, 12}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("...float64-ckr param: bad: And: not CF"),
			ExpErr: testhelper.MkExpErr("Can't make float64-checker function:" +
				" And(..., float64-checker):" +
				" can't convert argument 1 to float64-checker:" +
				" nonesuch is an unknown function"),
			expr: "And(GE(10), nonesuch, LT(11))",
		},
		{
			ID: testhelper.MkID(
				"...float64-ckr param: good: Or: 1 CF"),
			expr:        "Or(LE(10))",
			passingVals: map[int][]float64{0: {10}},
			failingVals: map[int][]float64{0: {12, 11}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID(
				"...float64-ckr param: good: Or: 2 CF"),
			expr:        "Or(GE(200), LT(11))",
			passingVals: map[int][]float64{0: {200, 201, 10.9}},
			failingVals: map[int][]float64{0: {199, 11}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("...float64-ckr param: bad: Or: not CF"),
			ExpErr: testhelper.MkExpErr("Can't make float64-checker function:" +
				" Or(..., float64-checker):" +
				" can't convert argument 1 to float64-checker:" +
				" nonesuch is an unknown function"),
			expr: "Or(GE(10), nonesuch, LT(11))",
		},
	}

	parser := checksetter.FindParserOrPanic[float64](
		checksetter.Float64CheckerName)
	for _, tc := range testCases {
		vcs, err := parser.Parse(tc.expr)
		if testhelper.CheckExpErr(t, err, tc) &&
			err == nil &&
			testhelper.DiffInt(t, tc.IDStr(), "number of ValCk funcs",
				len(vcs), tc.expLen) {
			for vcIdx, vc := range vcs {
				for _, pVal := range tc.passingVals[vcIdx] {
					if err = vc(pVal); err != nil {
						t.Log(tc.IDStr())
						t.Logf(
							"\t: error when checking %v with ValCk: %d",
							pVal, vcIdx)
						t.Error("\t: Bad check")
					}
				}

				for _, fVal := range tc.failingVals[vcIdx] {
					if err = vc(fVal); err == nil {
						t.Log(tc.IDStr())
						t.Logf(
							"\t: missing error when checking %v with ValCk: %d",
							fVal, vcIdx)
						t.Error("\t: Bad check")
					}
				}
			}
		}
	}
}

func TestParseString(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		expr        string
		passingVals map[int][]string
		failingVals map[int][]string
		expLen      int
	}{
		{
			ID: testhelper.MkID("bad: no-such name"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" nonesuch is an unknown function"),
			expr: "nonesuch",
		},
		{
			ID: testhelper.MkID("bad: not a named function"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" Syntax error: unexpected call type: *ast.FuncLit"),
			expr: "func(){}()",
		},
		{
			ID:     testhelper.MkID("bad: multi-ck, empty list entry at start"),
			ExpErr: testhelper.MkExpErr("expected operand, found ','"),
			expr:   ",OK",
		},
		{
			ID:     testhelper.MkID("bad: multi-ck, empty list entry - middle"),
			ExpErr: testhelper.MkExpErr("expected operand, found ','"),
			expr:   "OK,,OK",
		},
		{
			ID:          testhelper.MkID("no-params: good: just name"),
			expr:        "OK",
			passingVals: map[int][]string{0: {"", "Hello"}},
			expLen:      1,
		},
		{
			ID:          testhelper.MkID("no-params: good: just name call"),
			expr:        "OK()",
			passingVals: map[int][]string{0: {"", "Hello"}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("no-params: bad: with params"),
			ExpErr: testhelper.MkExpErr(
				"Can't make string-checker function: OK()",
				"the call has 1 arguments, it should have 0"),
			expr:   "OK(42)",
			expLen: 1,
		},
		{
			ID:   testhelper.MkID("no-params: good: multi"),
			expr: "OK, OK()",
			passingVals: map[int][]string{
				0: {"", "Hello"},
				1: {"World"},
			},
			expLen: 2,
		},
		{
			ID:          testhelper.MkID("1 string param: good: EQ"),
			expr:        "EQ(`A`)",
			passingVals: map[int][]string{0: {"A"}},
			failingVals: map[int][]string{0: {"B", "C"}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("1 string param: bad: EQ, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" EQ(string):" +
				" \"42\" isn't a STRING, it's a INT"),
			expr: "EQ(42)",
		},
		{
			ID: testhelper.MkID("1 string param: bad: EQ, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" EQ(string):" +
				" the call has 0 arguments, it should have 1"),
			expr: "EQ()",
		},
		{
			ID: testhelper.MkID("1 string param: bad: EQ, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" EQ(string):" +
				" the call has 2 arguments, it should have 1"),
			expr: "EQ(`Hello`, `World`)",
		},
		{
			ID:          testhelper.MkID("1 string param: good: GT"),
			expr:        "GT(`D`)",
			passingVals: map[int][]string{0: {`E`, `F`}},
			failingVals: map[int][]string{0: {`A`, `B`}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("1 string param: bad: GT, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" GT(string):" +
				" \"42\" isn't a STRING, it's a INT"),
			expr: "GT(42)",
		},
		{
			ID: testhelper.MkID("1 string param: bad: GT, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" GT(string):" +
				" the call has 0 arguments, it should have 1"),
			expr: "GT()",
		},
		{
			ID: testhelper.MkID("1 string param: bad: GT, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" GT(string):" +
				" the call has 2 arguments, it should have 1"),
			expr: "GT(`A`, `B`)",
		},
		{
			ID:          testhelper.MkID("1 string param: good: GE"),
			expr:        "GE(`D`)",
			passingVals: map[int][]string{0: {`D`, `E`}},
			failingVals: map[int][]string{0: {`A`, `B`}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("1 string param: bad: GE, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" GE(string):" +
				" \"42\" isn't a STRING, it's a INT"),
			expr: "GE(42)",
		},
		{
			ID: testhelper.MkID("1 string param: bad: GE, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" GE(string):" +
				" the call has 0 arguments, it should have 1"),
			expr: "GE()",
		},
		{
			ID: testhelper.MkID("1 string param: bad: GE, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" GE(string):" +
				" the call has 2 arguments, it should have 1"),
			expr: "GE(1, 2)",
		},
		{
			ID:          testhelper.MkID("1 string param: good: LT"),
			expr:        "LT(`D`)",
			passingVals: map[int][]string{0: {`A`, `B`}},
			failingVals: map[int][]string{0: {`D`, `E`}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("1 string param: bad: LT, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" LT(string):" +
				" \"42\" isn't a STRING, it's a INT"),
			expr: "LT(42)",
		},
		{
			ID: testhelper.MkID("1 string param: bad: LT, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" LT(string):" +
				" the call has 0 arguments, it should have 1"),
			expr: "LT()",
		},
		{
			ID: testhelper.MkID("1 string param: bad: LT, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" LT(string):" +
				" the call has 2 arguments, it should have 1"),
			expr: "LT(1, 2)",
		},
		{
			ID:          testhelper.MkID("1 string param: good: LE"),
			expr:        "LE(`D`)",
			passingVals: map[int][]string{0: {`A`, `D`}},
			failingVals: map[int][]string{0: {`E`, `F`}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("1 string param: bad: LE, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" LE(string):" +
				" \"42\" isn't a STRING, it's a INT"),
			expr: "LE(42)",
		},
		{
			ID: testhelper.MkID("1 string param: bad: LE, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" LE(string):" +
				" the call has 0 arguments, it should have 1"),
			expr: "LE()",
		},
		{
			ID: testhelper.MkID("1 string param: bad: LE, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" LE(string):" +
				" the call has 2 arguments, it should have 1"),
			expr: "LE(1, 2)",
		},
		{
			ID:          testhelper.MkID("1 string param: good: HasPrefix"),
			expr:        "HasPrefix(`D`)",
			passingVals: map[int][]string{0: {`Delay`, `Death`}},
			failingVals: map[int][]string{0: {`Elephant`, `Fatuous`}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("1 string param: bad: HasPrefix, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" HasPrefix(string):" +
				" \"42\" isn't a STRING, it's a INT"),
			expr: "HasPrefix(42)",
		},
		{
			ID: testhelper.MkID("1 string param: bad: HasPrefix, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" HasPrefix(string):" +
				" the call has 0 arguments, it should have 1"),
			expr: "HasPrefix()",
		},
		{
			ID: testhelper.MkID(
				"1 string param: bad: HasPrefix, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" HasPrefix(string):" +
				" the call has 2 arguments, it should have 1"),
			expr: "HasPrefix(1, 2)",
		},
		{
			ID:          testhelper.MkID("1 string param: good: HasSuffix"),
			expr:        "HasSuffix(`ed`)",
			passingVals: map[int][]string{0: {`Red`, `Subborned`}},
			failingVals: map[int][]string{0: {`Elephant`, `Fatuous`}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("1 string param: bad: HasSuffix, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" HasSuffix(string):" +
				" \"42\" isn't a STRING, it's a INT"),
			expr: "HasSuffix(42)",
		},
		{
			ID: testhelper.MkID("1 string param: bad: HasSuffix, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" HasSuffix(string):" +
				" the call has 0 arguments, it should have 1"),
			expr: "HasSuffix()",
		},
		{
			ID: testhelper.MkID(
				"1 string param: bad: HasSuffix, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" HasSuffix(string):" +
				" the call has 2 arguments, it should have 1"),
			expr: "HasSuffix(1, 2)",
		},
		{
			ID:          testhelper.MkID("1 int-checker param: good: Length"),
			expr:        "Length(GT(3))",
			passingVals: map[int][]string{0: {`Reed`, `Subborned`}},
			failingVals: map[int][]string{0: {`Red`, `It`}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("1 int-checker param: bad: Length, bad args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" Length(int-checker):" +
				" can't convert argument 0 to int-checker:" +
				" unexpected type: *ast.BasicLit"),
			expr: "Length(42)",
		},
		{
			ID: testhelper.MkID(
				"1 int-checker param: bad: Length, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" Length(int-checker):" +
				" the call has 0 arguments, it should have 1"),
			expr: "Length()",
		},
		{
			ID: testhelper.MkID(
				"1 int-checker param: bad: Length, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" Length(int-checker):" +
				" the call has 2 arguments, it should have 1"),
			expr: "Length(GE(1), Divides(40))",
		},
		{
			ID: testhelper.MkID(
				"1 regexp, 1 string param: good: MatchesPattern"),
			expr:        "MatchesPattern(`a.*z`, `name`)",
			passingVals: map[int][]string{0: {`a to z`, `a-z`}},
			failingVals: map[int][]string{0: {`Red`, `It`}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID(
				"1 regexp, 1 string param: bad: MatchesPattern, bad args 1st"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" MatchesPattern(regexp, string):" +
				" \"42\" isn't a STRING, it's a INT"),
			expr: "MatchesPattern(42, `name`)",
		},
		{
			ID: testhelper.MkID(
				"1 regexp, 1 string param: bad: MatchesPattern, bad regexp"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" MatchesPattern(regexp, string):" +
				" the regexp doesn't compile:" +
				" error parsing regexp: missing closing ]: `[Wworld``"),
			expr: "MatchesPattern(`Hello, [Wworld`, `name`)",
		},
		{
			ID: testhelper.MkID(
				"1 regexp, 1 string param: bad: MatchesPattern, bad args 2nd"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" MatchesPattern(regexp, string):" +
				" \"42\" isn't a STRING, it's a INT"),
			expr: "MatchesPattern(`name`, 42)",
		},
		{
			ID: testhelper.MkID(
				"1 regexp, 1 string param: bad: MatchesPattern, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" MatchesPattern(regexp, string):" +
				" the call has 0 arguments, it should have 2"),
			expr: "MatchesPattern()",
		},
		{
			ID: testhelper.MkID(
				"1 regexp, 1 string param: bad: MatchesPattern, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" MatchesPattern(regexp, string):" +
				" the call has 1 arguments, it should have 2"),
			expr: "MatchesPattern(`.*`)",
		},
		{
			ID: testhelper.MkID(
				"1 regexp, 1 string param: bad: MatchesPattern, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" MatchesPattern(regexp, string):" +
				" the call has 3 arguments, it should have 2"),
			expr: "MatchesPattern(`.*`, `any`, ` value`)",
		},
		{
			ID:   testhelper.MkID("1 string param: good: multi-checks"),
			expr: "LT(`D`), GT(`B`)",
			passingVals: map[int][]string{
				0: {`A`, `C`},
				1: {`C`, `D`},
			},
			failingVals: map[int][]string{
				0: {`D`, `E`},
				1: {`A`, `B`},
			},
			expLen: 2,
		},
		{
			ID:          testhelper.MkID("string-ckr, str param: good: Not"),
			expr:        "Not(GE(`D`), `not GE(\"D\")`)",
			passingVals: map[int][]string{0: {`D`, `E`}},
			failingVals: map[int][]string{0: {`A`, `C`}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID(
				"string-ckr, str param: bad: Not, bad args (1st)"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" Not(string-checker, string):" +
				" can't convert argument 0 to string-checker:" +
				" unexpected type: *ast.BasicLit"),
			expr: "Not(1, `hello`)",
		},
		{
			ID: testhelper.MkID(
				"string-ckr, str param: bad: Not, bad args (2nd)"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" Not(string-checker, string):" +
				" \"42\" isn't a STRING, it's a INT"),
			expr: "Not(GE(`A`), 42)",
		},
		{
			ID: testhelper.MkID(
				"string-ckr, str param: bad: Not, too many args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" Not(string-checker, string):" +
				" the call has 3 arguments, it should have 2"),
			expr: "Not(GE(`A`), `Hello`, `World`)",
		},
		{
			ID: testhelper.MkID(
				"string-ckr, str param: bad: Not, too few args"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" Not(string-checker, string):" +
				" the call has 1 arguments, it should have 2"),
			expr: "Not(GE(`A`))",
		},
		{
			ID: testhelper.MkID(
				"...string-ckr param: good: And: 1 CF"),
			expr:        "And(GE(`D`))",
			passingVals: map[int][]string{0: {`D`, `E`}},
			failingVals: map[int][]string{0: {`A`}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID(
				"...string-ckr param: good: And: 3 CF"),
			expr:        "And(GE(`D`), OK, LT(`F`))",
			passingVals: map[int][]string{0: {`D`, `E`}},
			failingVals: map[int][]string{0: {`A`, `F`}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("...string-ckr param: bad: And: not CF"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" And(..., string-checker):" +
				" can't convert argument 1 to string-checker:" +
				" nonesuch is an unknown function"),
			expr: "And(GE(`A`), nonesuch, LT(`B`))",
		},
		{
			ID:          testhelper.MkID("...string-ckr param: good: Or: 1 CF"),
			expr:        "Or(LE(`D`))",
			passingVals: map[int][]string{0: {`A`, `D`}},
			failingVals: map[int][]string{0: {`E`, `Z`}},
			expLen:      1,
		},
		{
			ID:          testhelper.MkID("...string-ckr param: good: Or: 2 CF"),
			expr:        "Or(GE(`D`), LT(`F`))",
			passingVals: map[int][]string{0: {`E`, `A`}},
			failingVals: map[int][]string{0: {}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("...string-ckr param: bad: Or: not CF"),
			ExpErr: testhelper.MkExpErr("Can't make string-checker function:" +
				" Or(..., string-checker):" +
				" can't convert argument 1 to string-checker:" +
				" nonesuch is an unknown function"),
			expr: "Or(GE(`A`), nonesuch, LT(`B`))",
		},
	}

	parser := checksetter.FindParserOrPanic[string](
		checksetter.StringCheckerName)
	for _, tc := range testCases {
		vcs, err := parser.Parse(tc.expr)
		if testhelper.CheckExpErr(t, err, tc) &&
			err == nil &&
			testhelper.DiffInt(t, tc.IDStr(), "number of ValCk funcs",
				len(vcs), tc.expLen) {
			for vcIdx, vc := range vcs {
				for _, pVal := range tc.passingVals[vcIdx] {
					if err = vc(pVal); err != nil {
						t.Log(tc.IDStr())
						t.Logf(
							"\t: error when checking %v with ValCk: %d",
							pVal, vcIdx)
						t.Error("\t: Bad check")
					}
				}

				for _, fVal := range tc.failingVals[vcIdx] {
					if err = vc(fVal); err == nil {
						t.Log(tc.IDStr())
						t.Logf(
							"\t: missing error when checking %v with ValCk: %d",
							fVal, vcIdx)
						t.Error("\t: Bad check")
					}
				}
			}
		}
	}
}

func TestParseStringSlice(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		expr        string
		passingVals map[int][][]string
		failingVals map[int][][]string
		expLen      int
	}{
		{
			ID: testhelper.MkID("bad: no-such name"),
			ExpErr: testhelper.MkExpErr(
				"Can't make string-slice-checker function:" +
					" nonesuch is an unknown function"),
			expr: "nonesuch",
		},
		{
			ID: testhelper.MkID("bad: not a named function"),
			ExpErr: testhelper.MkExpErr(
				"Can't make string-slice-checker function:" +
					" Syntax error: unexpected call type: *ast.FuncLit"),
			expr: "func(){}()",
		},
		{
			ID:     testhelper.MkID("bad: multi-ck, empty list entry at start"),
			ExpErr: testhelper.MkExpErr("expected operand, found ','"),
			expr:   ",OK",
		},
		{
			ID:     testhelper.MkID("bad: multi-ck, empty list entry - middle"),
			ExpErr: testhelper.MkExpErr("expected operand, found ','"),
			expr:   "OK,,OK",
		},
		{
			ID:          testhelper.MkID("no-params: good: just name"),
			expr:        "OK",
			passingVals: map[int][][]string{0: {{"hi", "world"}, {"hi"}, {}}},
			expLen:      1,
		},
		{
			ID:          testhelper.MkID("no-params: good: just name call"),
			expr:        "OK()",
			passingVals: map[int][][]string{0: {{"hi", "world"}, {"hi"}, {}}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("no-params: bad: with params"),
			ExpErr: testhelper.MkExpErr(
				"Can't make string-slice-checker function: OK()",
				"the call has 1 arguments, it should have 0"),
			expr:   "OK(42)",
			expLen: 1,
		},
		{
			ID:   testhelper.MkID("no-params: good: multi"),
			expr: "OK, OK()",
			passingVals: map[int][][]string{
				0: {{"hi", "world"}, {"hi"}, {}},
				1: {{"hi", "world"}, {"hi"}, {}},
			},
			expLen: 2,
		},
		{
			ID:          testhelper.MkID("no-params: good: just name"),
			expr:        "NoDups",
			passingVals: map[int][][]string{0: {{"hi", "world"}, {"hi"}, {}}},
			failingVals: map[int][][]string{
				0: {
					{"hi", "hi"},
					{"hi", "world", "hi"},
				},
			},
			expLen: 1,
		},
		{
			ID:          testhelper.MkID("no-params: good: just name call"),
			expr:        "NoDups()",
			passingVals: map[int][][]string{0: {{"hi", "world"}, {"hi"}, {}}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("no-params: bad: with params"),
			ExpErr: testhelper.MkExpErr(
				"Can't make string-slice-checker function: NoDups()",
				"the call has 1 arguments, it should have 0"),
			expr:   "NoDups(42)",
			expLen: 1,
		},
		{
			ID:   testhelper.MkID("no-params: good: multi"),
			expr: "NoDups, NoDups()",
			passingVals: map[int][][]string{
				0: {{"hi", "world"}, {"hi"}, {}},
				1: {{"hi", "world"}, {"hi"}, {}},
			},
			failingVals: map[int][][]string{
				0: {{"hi", "hi"}, {"hi", "world", "hi"}},
				1: {{"hi", "hi"}, {"hi", "world", "hi"}},
			},
			expLen: 2,
		},
		{
			ID:          testhelper.MkID("1 int-checker param: good: Length"),
			expr:        "Length(GT(3))",
			passingVals: map[int][][]string{0: {{`A`, `B`, `C`, `D`}}},
			failingVals: map[int][][]string{0: {{`Red`, `It`}}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("1 int-checker param: bad: Length, bad args"),
			ExpErr: testhelper.MkExpErr(
				"Can't make string-slice-checker function:" +
					" Length(int-checker):" +
					" can't convert argument 0 to int-checker:" +
					" unexpected type: *ast.BasicLit"),
			expr: "Length(42)",
		},
		{
			ID: testhelper.MkID(
				"1 int-checker param: bad: Length, too few args"),
			ExpErr: testhelper.MkExpErr(
				"Can't make string-slice-checker function:" +
					" Length(int-checker):" +
					" the call has 0 arguments, it should have 1"),
			expr: "Length()",
		},
		{
			ID: testhelper.MkID(
				"1 int-checker param: bad: Length, too many args"),
			ExpErr: testhelper.MkExpErr(
				"Can't make string-slice-checker function:" +
					" Length(int-checker):" +
					" the call has 2 arguments, it should have 1"),
			expr: "Length(GE(1), Divides(40))",
		},
		{
			ID:          testhelper.MkID("str-slc-ckr, str param: good: Not"),
			expr:        "Not(Length(GE(1)), `not GE(10)`)",
			passingVals: map[int][][]string{0: {{}}},
			failingVals: map[int][][]string{0: {{`a`}, {`a`, `b`}}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID(
				"str-slc-ckr, str param: bad: Not, bad args (1st)"),
			ExpErr: testhelper.MkExpErr(
				"Can't make string-slice-checker function:" +
					" Not(string-slice-checker, string):" +
					" can't convert argument 0 to string-slice-checker:" +
					" unexpected type: *ast.BasicLit"),
			expr: "Not(1, `hello`)",
		},
		{
			ID: testhelper.MkID(
				"str-slc-ckr, str param: bad: Not, bad args (2nd)"),
			ExpErr: testhelper.MkExpErr(
				"Can't make string-slice-checker function:" +
					" Not(string-slice-checker, string):" +
					" \"1\" isn't a STRING, it's a INT"),
			expr: "Not(Length(GE(1)), 1)",
		},
		{
			ID: testhelper.MkID(
				"str-slc-ckr, str param: bad: Not, too many args"),
			ExpErr: testhelper.MkExpErr(
				"Can't make string-slice-checker function:" +
					" Not(string-slice-checker, string):" +
					" the call has 3 arguments, it should have 2"),
			expr: "Not(GE(1), `Hello`, `World`)",
		},
		{
			ID: testhelper.MkID(
				"str-slc-ckr, str param: bad: Not, too few args"),
			ExpErr: testhelper.MkExpErr(
				"Can't make string-slice-checker function:" +
					" Not(string-slice-checker, string):" +
					" the call has 1 arguments, it should have 2"),
			expr: "Not(GE(1))",
		},
		{
			ID: testhelper.MkID(
				"string-checker, string param: good: SliceAny"),
			expr: "SliceAny(Length(GT(3)), `good-checker`)",
			passingVals: map[int][][]string{
				0: {
					{`Absolute`, `B`, `C`, `D`},
					{``, `Double`},
				},
			},
			failingVals: map[int][][]string{0: {{`Red`, `It`}}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID(
				"string-checker, string param: bad: SliceAny, bad args (1st)"),
			ExpErr: testhelper.MkExpErr(
				"Can't make string-slice-checker function:" +
					" SliceAny(string-checker, string):" +
					" can't convert argument 0 to string-checker:" +
					" unexpected type: *ast.BasicLit"),
			expr: "SliceAny(42, `hello`)",
		},
		{
			ID: testhelper.MkID(
				"string-checker, string param: bad: SliceAny, bad args (2nd)"),
			ExpErr: testhelper.MkExpErr(
				"Can't make string-slice-checker function:" +
					" SliceAny(string-checker, string):" +
					" \"42\" isn't a STRING, it's a INT"),
			expr: "SliceAny(OK, 42)",
		},
		{
			ID: testhelper.MkID(
				"string-checker, string param: bad: SliceAny, too few args"),
			ExpErr: testhelper.MkExpErr(
				"Can't make string-slice-checker function:" +
					" SliceAny(string-checker, string):" +
					" the call has 0 arguments, it should have 2"),
			expr: "SliceAny()",
		},
		{
			ID: testhelper.MkID(
				"string-checker, string param: bad: SliceAny, too few args"),
			ExpErr: testhelper.MkExpErr(
				"Can't make string-slice-checker function:" +
					" SliceAny(string-checker, string):" +
					" the call has 1 arguments, it should have 2"),
			expr: "SliceAny(Length(GT(3)))",
		},
		{
			ID: testhelper.MkID(
				"string-checker, string param: bad: SliceAny, too many args"),
			ExpErr: testhelper.MkExpErr(
				"Can't make string-slice-checker function:" +
					" SliceAny(string-checker, string):" +
					" the call has 3 arguments, it should have 2"),
			expr: "SliceAny(Length(GE(1)), `str-slc-ckr`, 42)",
		},
		{
			ID:   testhelper.MkID("string-checker param: good: SliceAll"),
			expr: "SliceAll(Length(GT(3)))",
			passingVals: map[int][][]string{
				0: {
					{`Absolute`, `Batch`, `Cold`, `Dogs`},
					{`1234`, `Double`},
				},
			},
			failingVals: map[int][][]string{
				0: {
					{`Red`, `Redden`},
					{`Blue`, `Blu`},
				},
			},
			expLen: 1,
		},
		{
			ID: testhelper.MkID(
				"string-checker param: bad: SliceAll, bad args"),
			ExpErr: testhelper.MkExpErr(
				"Can't make string-slice-checker function:" +
					" SliceAll(string-checker):" +
					" can't convert argument 0 to string-checker:" +
					" unexpected type: *ast.BasicLit"),
			expr: "SliceAll(42)",
		},
		{
			ID: testhelper.MkID(
				"string-checker param: bad: SliceAll, too few args"),
			ExpErr: testhelper.MkExpErr(
				"Can't make string-slice-checker function:" +
					" SliceAll(string-checker):" +
					" the call has 0 arguments, it should have 1"),
			expr: "SliceAll()",
		},
		{
			ID: testhelper.MkID(
				"string-checker param: bad: SliceAll, too many args"),
			ExpErr: testhelper.MkExpErr(
				"Can't make string-slice-checker function:" +
					" SliceAll(string-checker):" +
					" the call has 2 arguments, it should have 1"),
			expr: "SliceAll(Length(GE(1)), 42)",
		},

		{
			ID:   testhelper.MkID("string-checker param: good: SliceByPos"),
			expr: "SliceByPos(HasPrefix(`RC`), OK, Length(GT(3)))",
			passingVals: map[int][][]string{
				0: {
					{`RC42`, `anything`, `Cold`, `Dogs`},
					{`RC1234`, ``, `four`},
				},
			},
			failingVals: map[int][][]string{
				0: {
					{`RC`, ``, `Red`},
					{`rc`, `Blue`, `Blue`},
					{`R`, `Clue`, `Blue`},
				},
			},
			expLen: 1,
		},
		{
			ID: testhelper.MkID(
				"string-checker param: bad: SliceByPos, bad args"),
			ExpErr: testhelper.MkExpErr(
				"Can't make string-slice-checker function:" +
					" SliceByPos(..., string-checker):" +
					" can't convert argument 0 to string-checker:" +
					" unexpected type: *ast.BasicLit"),
			expr: "SliceByPos(42)",
		},
		{
			ID:   testhelper.MkID("...str-slc-ckr param: good: And: 2 CF"),
			expr: "And(NoDups, Length(LE(3)))",
			passingVals: map[int][][]string{
				0: {
					{`a`}, {`a`, `b`, `c`},
				},
			},
			failingVals: map[int][][]string{
				0: {
					{`a`, `a`, `b`},
					{`a`, `b`, `c`, `d`},
				},
			},
			expLen: 1,
		},
		{
			ID: testhelper.MkID("...str-slc-ckr param: bad: And: not CF"),
			ExpErr: testhelper.MkExpErr(
				"Can't make string-slice-checker function:" +
					" And(..., string-slice-checker):" +
					" can't convert argument 1 to string-slice-checker:" +
					" nonesuch is an unknown function"),
			expr: "And(OK, nonesuch, NoDups)",
		},
		{
			ID:   testhelper.MkID("...str-slc-ckr param: good: Or: 1 CF"),
			expr: "Or(NoDups)",
			passingVals: map[int][][]string{
				0: {
					{`a`},
				},
			},
			failingVals: map[int][][]string{
				0: {
					{`a`, `a`},
				},
			},
			expLen: 1,
		},
		{
			ID:   testhelper.MkID("...str-slc-ckr param: good: Or: 2 CF"),
			expr: "Or(NoDups, Length(GT(2)))",
			passingVals: map[int][][]string{
				0: {
					{`a`},
					{`a`, `a`, `b`},
				},
			},
			failingVals: map[int][][]string{0: {{`a`, `a`}}},
			expLen:      1,
		},
		{
			ID: testhelper.MkID("...str-slc-ckr param: bad: Or: not CF"),
			ExpErr: testhelper.MkExpErr(
				"Can't make string-slice-checker function:" +
					" Or(..., string-slice-checker):" +
					" can't convert argument 1 to string-slice-checker:" +
					" nonesuch is an unknown function"),
			expr: "Or(Length(GE(10)), nonesuch, OK)",
		},
	}

	parser := checksetter.FindParserOrPanic[[]string](
		checksetter.StringSliceCheckerName)
	for _, tc := range testCases {
		vcs, err := parser.Parse(tc.expr)
		if testhelper.CheckExpErr(t, err, tc) &&
			err == nil &&
			testhelper.DiffInt(t, tc.IDStr(), "number of ValCk funcs",
				len(vcs), tc.expLen) {
			for vcIdx, vc := range vcs {
				for _, pVal := range tc.passingVals[vcIdx] {
					if err = vc(pVal); err != nil {
						t.Log(tc.IDStr())
						t.Logf(
							"\t: error when checking %v with ValCk: %d",
							pVal, vcIdx)
						t.Error("\t: Bad check")
					}
				}

				for _, fVal := range tc.failingVals[vcIdx] {
					if err = vc(fVal); err == nil {
						t.Log(tc.IDStr())
						t.Logf(
							"\t: missing error when checking %v with ValCk: %d",
							fVal, vcIdx)
						t.Error("\t: Bad check")
					}
				}
			}
		}
	}
}
