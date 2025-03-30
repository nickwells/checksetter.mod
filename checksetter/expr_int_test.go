package checksetter

import (
	"go/ast"
	"go/parser"
	"testing"

	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

// callExprTestVals returns various expressions for other tests
func callExprTestVals(t *testing.T) (callExprFFF, callExprFIF *ast.CallExpr) {
	t.Helper()

	var ok bool

	var err error

	var ce ast.Expr

	exprStr := `callFunc(OK, OK, OK)`

	ce, err = parser.ParseExpr(exprStr)
	if err != nil {
		t.Fatal("cannot parse the expression: ", exprStr, " error: ", err)
	}

	callExprFFF, ok = ce.(*ast.CallExpr)
	if !ok {
		t.Fatalf("the expression is not an ast.CallExpr: %T", callExprFFF)
	}

	exprStr = `callFunc(OK, 99, OK)`

	ce, err = parser.ParseExpr(exprStr)
	if err != nil {
		t.Fatal("cannot parse the expression: ", exprStr, " error: ", err)
	}

	callExprFIF, ok = ce.(*ast.CallExpr)
	if !ok {
		t.Fatalf("the expression is not an ast.CallExpr: %T", callExprFFF)
	}

	return
}

// exprTestVals returns various expressions for other tests
func exprTestVals(t *testing.T) (callExpr, litInt, litFloat, litStr ast.Expr) {
	t.Helper()

	exprStr := `callFunc(1, 1.5, "hello")`

	var err error

	callExpr, err = parser.ParseExpr(exprStr)
	if err != nil {
		t.Fatal("cannot parse the expression: ", exprStr, " error: ", err)
	}

	ce, ok := callExpr.(*ast.CallExpr)
	if !ok {
		t.Fatalf("the expression is not an ast.CallExpr: %T", callExpr)
	}

	litInt = ce.Args[0]
	litFloat = ce.Args[1]
	litStr = ce.Args[2]

	return
}

// exprTestBigVals returns various expressions representing large literal
// values for other tests
func exprTestBigVals(t *testing.T) (litInt, litFloat ast.Expr) {
	t.Helper()

	exprStr := `callFunc(999999999999999999999, 1e999999)`

	callExpr, err := parser.ParseExpr(exprStr)
	if err != nil {
		t.Fatal("cannot parse the expression: ", exprStr, " error: ", err)
	}

	ce, ok := callExpr.(*ast.CallExpr)
	if !ok {
		t.Fatalf("the expression is not an ast.CallExpr: %T", callExpr)
	}

	litInt = ce.Args[0]
	litFloat = ce.Args[1]

	return
}

func TestGetInt(t *testing.T) {
	callExpr, litInt, litFloat, litStr := exprTestVals(t)
	bigLitInt, _ := exprTestBigVals(t)

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		param       ast.Expr
		valExpected int64
	}{
		{
			ID:          testhelper.MkID("good"),
			param:       litInt,
			valExpected: 1,
		},
		{
			ID: testhelper.MkID(
				"bad - is a BasicLit but not an INT (FLOAT)"),
			param:  litFloat,
			ExpErr: testhelper.MkExpErr("isn't an INT, it's a FLOAT"),
		},
		{
			ID: testhelper.MkID(
				"bad - is a BasicLit but not an INT (STRING)"),
			param:  litStr,
			ExpErr: testhelper.MkExpErr("isn't an INT, it's a STRING"),
		},
		{
			ID:    testhelper.MkID("bad - not a BasicLit"),
			param: callExpr,
			ExpErr: testhelper.MkExpErr(
				"the expression isn't a BasicLit, it's a *ast.CallExpr"),
		},
		{
			ID:    testhelper.MkID("bad - int too big"),
			param: bigLitInt,
			ExpErr: testhelper.MkExpErr(
				`couldn't make an int from "999999999999999999999":` +
					` strconv.ParseInt: parsing "999999999999999999999":` +
					` value out of range`),
		},
	}

	for _, tc := range testCases {
		val, err := getInt64(tc.param)
		if testhelper.CheckExpErr(t, err, tc) &&
			err == nil {
			testhelper.DiffInt(t, tc.IDStr(), "int", val, tc.valExpected)
		}
	}
}

func TestGetFloat(t *testing.T) {
	callExpr, litInt, litFloat, litStr := exprTestVals(t)
	_, bigLitFloat := exprTestBigVals(t)

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		param       ast.Expr
		valExpected float64
	}{
		{
			ID:          testhelper.MkID("good - FLOAT"),
			param:       litFloat,
			valExpected: 1.5,
		},
		{
			ID:          testhelper.MkID("good - INT"),
			param:       litInt,
			valExpected: 1.0,
		},
		{
			ID:     testhelper.MkID("bad - is a BasicLit but not a FLOAT"),
			param:  litStr,
			ExpErr: testhelper.MkExpErr("isn't a FLOAT/INT, it's a STRING"),
		},
		{
			ID:    testhelper.MkID("bad - not a BasicLit"),
			param: callExpr,
			ExpErr: testhelper.MkExpErr(
				"the expression isn't a BasicLit, it's a *ast.CallExpr"),
		},
		{
			ID:    testhelper.MkID("bad - not a BasicLit"),
			param: bigLitFloat,
			ExpErr: testhelper.MkExpErr(
				`couldn't make a float from "1e999999":` +
					` strconv.ParseFloat: parsing "1e999999":` +
					` value out of range`),
		},
	}

	for _, tc := range testCases {
		val, err := getFloat64(tc.param)
		if testhelper.CheckExpErr(t, err, tc) &&
			err == nil {
			testhelper.DiffFloat(t, tc.IDStr(), "float",
				val, tc.valExpected, 0.000001)
		}
	}
}

func TestGetString(t *testing.T) {
	callExpr, litInt, litFloat, litStr := exprTestVals(t)

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		param       ast.Expr
		valExpected string
	}{
		{
			ID:          testhelper.MkID("good"),
			param:       litStr,
			valExpected: "hello",
		},
		{
			ID:     testhelper.MkID("bad - is a BasicLit but not a STRING"),
			param:  litInt,
			ExpErr: testhelper.MkExpErr("isn't a STRING, it's a INT"),
		},
		{
			ID:     testhelper.MkID("bad - is a BasicLit but not a STRING"),
			param:  litFloat,
			ExpErr: testhelper.MkExpErr("isn't a STRING, it's a FLOAT"),
		},
		{
			ID:    testhelper.MkID("bad - not a BasicLit"),
			param: callExpr,
			ExpErr: testhelper.MkExpErr(
				"the expression isn't a BasicLit, it's a *ast.CallExpr"),
		},
	}

	for _, tc := range testCases {
		val, err := getString(tc.param)
		if testhelper.CheckExpErr(t, err, tc) &&
			err == nil {
			testhelper.DiffString(t, tc.IDStr(), "string", val, tc.valExpected)
		}
	}
}

func TestGetElts(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		s      string
		expLen int
	}{
		{
			ID:     testhelper.MkID("good, empty"),
			s:      "",
			expLen: 0,
		},
		{
			ID:     testhelper.MkID("good, one entry"),
			s:      "Call",
			expLen: 1,
		},
		{
			ID:     testhelper.MkID("good, two entries"),
			s:      "Call, Call2(1)",
			expLen: 2,
		},
		{
			ID:     testhelper.MkID("bad, syntax error"),
			ExpErr: testhelper.MkExpErr("expected 'EOF', found '}'"),
			s:      "}",
			expLen: 0,
		},
	}

	for _, tc := range testCases {
		elts, err := getElts(tc.s, "test")
		if testhelper.CheckExpErr(t, err, tc) &&
			err == nil {
			testhelper.DiffInt(t, tc.IDStr(), "length", len(elts), tc.expLen)
		}
	}
}

func TestCheckArgCount(t *testing.T) {
	e, _, _, _ := exprTestVals(t)

	callExpr, ok := e.(*ast.CallExpr)
	if !ok {
		t.Fatalf("the expression is not an ast.CallExpr: %T", e)
	}

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		expArgs int
	}{
		{
			ID:      testhelper.MkID("good"),
			expArgs: 3,
		},
		{
			ID: testhelper.MkID("bad, expect too many"),
			ExpErr: testhelper.MkExpErr(
				"the call has 3 arguments, it should have 4"),
			expArgs: 4,
		},
		{
			ID: testhelper.MkID("bad, expect too few"),
			ExpErr: testhelper.MkExpErr(
				"the call has 3 arguments, it should have 2"),
			expArgs: 2,
		},
	}

	for _, tc := range testCases {
		err := checkArgCount(callExpr, tc.expArgs)
		testhelper.CheckExpErr(t, err, tc)
	}
}

func TestGetCheckFunc(t *testing.T) {
	exprStr := `callFunc(1, OK)`

	ce, err := parser.ParseExpr(exprStr)
	if err != nil {
		t.Fatal("cannot parse the expression: ", exprStr, " error: ", err)
	}

	callExpr, ok := ce.(*ast.CallExpr)
	if !ok {
		t.Fatalf("the expression is not an ast.CallExpr: %T", callExpr)
	}

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		idx int
	}{
		{
			ID:  testhelper.MkID("good"),
			idx: 1,
		},
		{
			ID: testhelper.MkID("bad, not a checker"),
			ExpErr: testhelper.MkExpErr(
				"can't convert argument 0 to ",
				"unexpected type: *ast.BasicLit",
			),
			idx: 0,
		},
		{
			ID: testhelper.MkID("bad, idx too big"),
			ExpErr: testhelper.MkExpErr(
				"index (2) is too large, there are only 2 arguments"),
			idx: 2,
		},
		{
			ID:     testhelper.MkID("bad, idx < 0"),
			ExpErr: testhelper.MkExpErr("index (-1) must be >= 0"),
			idx:    -1,
		},
	}

	for _, tc := range testCases {
		_, err = getCheckFunc[float64](callExpr, tc.idx, Float64CheckerName)
		testhelper.CheckExpErr(t, err, tc)
		_, err = getCheckFunc[int64](callExpr, tc.idx, Int64CheckerName)
		testhelper.CheckExpErr(t, err, tc)
		_, err = getCheckFunc[int](callExpr, tc.idx, IntCheckerName)
		testhelper.CheckExpErr(t, err, tc)
		_, err = getCheckFunc[string](callExpr, tc.idx, StringCheckerName)
		testhelper.CheckExpErr(t, err, tc)
		_, err = getCheckFunc[[]string](callExpr, tc.idx, StringSliceCheckerName)
		testhelper.CheckExpErr(t, err, tc)
	}

	expErr := testhelper.MkExpErr(
		`there is no Parser registered for "nonesuch"`)
	_, err = getCheckFunc[int](callExpr, 1, "nonesuch")
	testhelper.CheckExpErrWithID(t, "nonesuch", err, expErr)
}

func TestGetCheckFuncs(t *testing.T) {
	callExprFFF, callExprFIF := callExprTestVals(t)

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		expLen int
	}{
		{
			ID: testhelper.MkID("good"),
			ExpErr: testhelper.MkExpErr(
				"can't convert argument 1 to ",
				"unexpected type: *ast.BasicLit"),
			expLen: 3,
		},
	}

	for _, tc := range testCases {
		{
			valCks, err := getCheckFuncs[float64](
				callExprFFF, Float64CheckerName)
			if err != nil {
				t.Log(tc.IDStr())
				t.Errorf("\tunexpected error: %s", err)

				continue
			}

			testhelper.DiffInt(t, tc.IDStr(), "length of checker slice",
				len(valCks), tc.expLen)

			_, err = getCheckFuncs[float64](callExprFIF, Float64CheckerName)
			testhelper.CheckExpErr(t, err, tc)
		}
		{
			valCks, err := getCheckFuncs[int64](
				callExprFFF, Int64CheckerName)
			if err != nil {
				t.Log(tc.IDStr())
				t.Errorf("\tunexpected error: %s", err)

				continue
			}

			testhelper.DiffInt(t, tc.IDStr(), "length of checker slice",
				len(valCks), tc.expLen)

			_, err = getCheckFuncs[int64](callExprFIF, Int64CheckerName)
			testhelper.CheckExpErr(t, err, tc)
		}
		{
			valCks, err := getCheckFuncs[int](
				callExprFFF, IntCheckerName)
			if err != nil {
				t.Log(tc.IDStr())
				t.Errorf("\tunexpected error: %s", err)

				continue
			}

			testhelper.DiffInt(t, tc.IDStr(), "length of checker slice",
				len(valCks), tc.expLen)

			_, err = getCheckFuncs[int](callExprFIF, IntCheckerName)
			testhelper.CheckExpErr(t, err, tc)
		}
		{
			valCks, err := getCheckFuncs[string](
				callExprFFF, StringCheckerName)
			if err != nil {
				t.Log(tc.IDStr())
				t.Errorf("\tunexpected error: %s", err)

				continue
			}

			testhelper.DiffInt(t, tc.IDStr(), "length of checker slice",
				len(valCks), tc.expLen)

			_, err = getCheckFuncs[string](callExprFIF, StringCheckerName)
			testhelper.CheckExpErr(t, err, tc)
		}
		{
			valCks, err := getCheckFuncs[[]string](
				callExprFFF, StringSliceCheckerName)
			if err != nil {
				t.Log(tc.IDStr())
				t.Errorf("\tunexpected error: %s", err)

				continue
			}

			testhelper.DiffInt(t, tc.IDStr(), "length of checker slice",
				len(valCks), tc.expLen)

			_, err = getCheckFuncs[[]string](callExprFIF, StringSliceCheckerName)
			testhelper.CheckExpErr(t, err, tc)
		}
	}

	expErr := testhelper.MkExpErr(
		`there is no Parser registered for "nonesuch"`)
	_, err := getCheckFuncs[int](nil, "nonesuch")

	testhelper.CheckExpErrWithID(t, "nonesuch", err, expErr)
}
