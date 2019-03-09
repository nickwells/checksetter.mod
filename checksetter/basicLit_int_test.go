package checksetter

import (
	"fmt"
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
		name             string
		param            ast.Expr
		valExpected      int64
		errExpected      bool
		errShouldContain []string
	}{
		{
			name:        "good",
			param:       exprArgInt,
			valExpected: 1,
		},
		{
			name:        "bad - is a BasicLit but not an INT",
			param:       exprArgNotInt,
			errExpected: true,
			errShouldContain: []string{
				"the expression should have been a literal INT, was STRING",
			},
		},
		{
			name:        "bad - not a BasicLit",
			param:       expr,
			errExpected: true,
			errShouldContain: []string{
				"the expression should have been a literal not *ast.CallExpr",
			},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)
		val, err := getInt(tc.param)
		if testhelper.CheckError(t, tcID,
			err, tc.errExpected, tc.errShouldContain) &&
			err == nil {
			if val != tc.valExpected {
				t.Log(tcID)
				t.Logf("\t: expected: %d\n", tc.valExpected)
				t.Logf("\t:      got: %d\n", val)
				t.Errorf("\t: value unexpected\n")
			}
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
		name             string
		param            ast.Expr
		valExpected      float64
		errExpected      bool
		errShouldContain []string
	}{
		{
			name:        "good",
			param:       exprArgFloat,
			valExpected: 1.1,
		},
		{
			name:        "bad - is a BasicLit but not a FLOAT",
			param:       exprArgNotFloat,
			errExpected: true,
			errShouldContain: []string{
				"the expression should have been a literal FLOAT, was STRING",
			},
		},
		{
			name:        "bad - not a BasicLit",
			param:       expr,
			errExpected: true,
			errShouldContain: []string{
				"the expression should have been a literal not *ast.CallExpr",
			},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)
		val, err := getFloat(tc.param)
		if testhelper.CheckError(t, tcID,
			err, tc.errExpected, tc.errShouldContain) &&
			err == nil {
			if val != tc.valExpected {
				t.Log(tcID)
				t.Logf("\t: expected: %f\n", tc.valExpected)
				t.Logf("\t:      got: %f\n", val)
				t.Errorf("\t: value unexpected\n")
			}
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
		name             string
		param            ast.Expr
		valExpected      string
		errExpected      bool
		errShouldContain []string
	}{
		{
			name:        "good",
			param:       exprArgString,
			valExpected: `"hello"`,
		},
		{
			name:        "bad - is a BasicLit but not a STRING",
			param:       exprArgNotString,
			errExpected: true,
			errShouldContain: []string{
				"the expression should have been a literal STRING, was FLOAT",
			},
		},
		{
			name:        "bad - not a BasicLit",
			param:       expr,
			errExpected: true,
			errShouldContain: []string{
				"the expression should have been a literal not *ast.CallExpr",
			},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)
		val, err := getString(tc.param)
		if testhelper.CheckError(t, tcID,
			err, tc.errExpected, tc.errShouldContain) &&
			err == nil {
			if val != tc.valExpected {
				t.Log(tcID)
				t.Logf("\t: expected: %s\n", tc.valExpected)
				t.Logf("\t:      got: %s\n", val)
				t.Errorf("\t: value unexpected\n")
			}
		}
	}
}
