package checksetter

import (
	"go/ast"
	"go/parser"
	"testing"

	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestGetInt(t *testing.T) {
	exprStr := `callFunc(1, "hello")`
	expr, err := parser.ParseExpr(exprStr)
	if err != nil {
		t.Fatal("cannot parse the expression: ", exprStr, " error: ", err)
	}

	callExpr, ok := expr.(*ast.CallExpr)
	if !ok {
		t.Fatalf("the expression is not a call: %T", expr)
	}
	exprArgInt := callExpr.Args[0]
	exprArgNotInt := callExpr.Args[1]

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		param       ast.Expr
		valExpected int64
	}{
		{
			ID:          testhelper.MkID("good"),
			param:       exprArgInt,
			valExpected: 1,
		},
		{
			ID:    testhelper.MkID("bad - is a BasicLit but not an INT"),
			param: exprArgNotInt,
			ExpErr: testhelper.MkExpErr(
				"the expression should have been a literal INT, was STRING"),
		},
		{
			ID:    testhelper.MkID("bad - not a BasicLit"),
			param: expr,
			ExpErr: testhelper.MkExpErr(
				"the expression should have been a literal not *ast.CallExpr"),
		},
	}

	for _, tc := range testCases {
		val, err := getInt64(tc.param)
		if testhelper.CheckExpErr(t, err, tc) &&
			err == nil {
			testhelper.DiffInt64(t, tc.IDStr(), "int", val, tc.valExpected)
		}
	}
}

func TestGetFloat(t *testing.T) {
	exprStr := `callFunc(1.1, "hello")`
	expr, err := parser.ParseExpr(exprStr)
	if err != nil {
		t.Fatal("cannot parse the expression: ", exprStr, " error: ", err)
	}

	callExpr, ok := expr.(*ast.CallExpr)
	if !ok {
		t.Fatalf("the expression is not a call: %T", expr)
	}
	exprArgFloat := callExpr.Args[0]
	exprArgNotFloat := callExpr.Args[1]

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		param       ast.Expr
		valExpected float64
	}{
		{
			ID:          testhelper.MkID("good"),
			param:       exprArgFloat,
			valExpected: 1.1,
		},
		{
			ID:    testhelper.MkID("bad - is a BasicLit but not a FLOAT"),
			param: exprArgNotFloat,
			ExpErr: testhelper.MkExpErr(
				"the expression should have been a literal FLOAT, was STRING"),
		},
		{
			ID:    testhelper.MkID("bad - not a BasicLit"),
			param: expr,
			ExpErr: testhelper.MkExpErr(
				"the expression should have been a literal not *ast.CallExpr"),
		},
	}

	for _, tc := range testCases {
		val, err := getFloat(tc.param)
		if testhelper.CheckExpErr(t, err, tc) &&
			err == nil {
			testhelper.DiffFloat64(t, tc.IDStr(), "float",
				val, tc.valExpected, 0.000001)
		}
	}
}

func TestGetString(t *testing.T) {
	exprStr := `callFunc(1.1, "hello")`
	expr, err := parser.ParseExpr(exprStr)
	if err != nil {
		t.Fatal("cannot parse the expression: ", exprStr, " error: ", err)
	}

	callExpr, ok := expr.(*ast.CallExpr)
	if !ok {
		t.Fatalf("the expression is not a call: %T", expr)
	}
	exprArgString := callExpr.Args[1]
	exprArgNotString := callExpr.Args[0]

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		param       ast.Expr
		valExpected string
	}{
		{
			ID:          testhelper.MkID("good"),
			param:       exprArgString,
			valExpected: "hello",
		},
		{
			ID:    testhelper.MkID("bad - is a BasicLit but not a STRING"),
			param: exprArgNotString,
			ExpErr: testhelper.MkExpErr(
				"the expression should have been a literal STRING, was FLOAT"),
		},
		{
			ID:    testhelper.MkID("bad - not a BasicLit"),
			param: expr,
			ExpErr: testhelper.MkExpErr(
				"the expression should have been a literal not *ast.CallExpr"),
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
