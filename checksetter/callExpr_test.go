package checksetter

import (
	"fmt"
	"go/ast"
	"go/parser"
	"testing"

	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestCheckArgCount(t *testing.T) {
	callExprs := []struct {
		s  string
		ce *ast.CallExpr
	}{
		{s: `callFunc()`},
		{s: `callFunc(1)`},
		{s: `callFunc(1, 2)`},
	}
	for i, ce := range callExprs {
		expr, err := parser.ParseExpr(ce.s)
		if err != nil {
			t.Fatal("cannot parse the expression: ", ce.s, " error: ", err)
		}

		callExpr, ok := expr.(*ast.CallExpr)
		if !ok {
			t.Fatalf("the expression is not a call: %T", expr)
		}
		callExprs[i].ce = callExpr
	}

	testCases := []struct {
		name             string
		ce               *ast.CallExpr
		n                int
		errExpected      bool
		errShouldContain []string
	}{
		{
			name: "good - 0 params",
			ce:   callExprs[0].ce,
			n:    0,
		},
		{
			name:        "bad - 0 params",
			ce:          callExprs[0].ce,
			n:           99,
			errExpected: true,
			errShouldContain: []string{
				"the call has 0 arguments, it should have 99",
			},
		},
		{
			name: "good - 1 params",
			ce:   callExprs[1].ce,
			n:    1,
		},
		{
			name:        "bad - 1 params",
			ce:          callExprs[1].ce,
			n:           99,
			errExpected: true,
			errShouldContain: []string{
				"the call has 1 arguments, it should have 99",
			},
		},
		{
			name: "good - 2 params",
			ce:   callExprs[2].ce,
			n:    2,
		},
		{
			name:        "bad - 2 params",
			ce:          callExprs[2].ce,
			n:           99,
			errExpected: true,
			errShouldContain: []string{
				"the call has 2 arguments, it should have 99",
			},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)
		err := checkArgCount(tc.ce, tc.n)
		testhelper.CheckError(t, tcID, err, tc.errExpected, tc.errShouldContain)
	}
}

func TestGetArg(t *testing.T) {
	callExprs := []struct {
		s  string
		ce *ast.CallExpr
	}{
		{s: `callFunc()`},
		{s: `callFunc(1)`},
		{s: `callFunc(1, 2)`},
	}
	for i, ce := range callExprs {
		expr, err := parser.ParseExpr(ce.s)
		if err != nil {
			t.Fatal("cannot parse the expression: ", ce.s, " error: ", err)
		}

		callExpr, ok := expr.(*ast.CallExpr)
		if !ok {
			t.Fatalf("the expression is not a call: %T", expr)
		}
		callExprs[i].ce = callExpr
	}

	testCases := []struct {
		name             string
		ce               *ast.CallExpr
		n                int
		errExpected      bool
		errShouldContain []string
	}{
		{
			name:        "bad - 0 params, get with idx: 0",
			ce:          callExprs[0].ce,
			n:           0,
			errExpected: true,
			errShouldContain: []string{
				"can't get argument 0, too few parameters",
			},
		},
		{
			name:        "bad - 0 params, get with idx: 1",
			ce:          callExprs[0].ce,
			n:           1,
			errExpected: true,
			errShouldContain: []string{
				"can't get argument 1, too few parameters",
			},
		},
		{
			name: "good - 1 params, get with idx: 0",
			ce:   callExprs[1].ce,
			n:    0,
		},
		{
			name:        "bad - 1 params, get with idx: 1",
			ce:          callExprs[1].ce,
			n:           1,
			errExpected: true,
			errShouldContain: []string{
				"can't get argument 1, too few parameters",
			},
		},
		{
			name: "good - 2 params, get with idx: 0",
			ce:   callExprs[2].ce,
			n:    0,
		},
		{
			name: "good - 2 params, get with idx: 1",
			ce:   callExprs[2].ce,
			n:    1,
		},
		{
			name:        "bad - 2 params, get with idx: 2",
			ce:          callExprs[2].ce,
			n:           2,
			errExpected: true,
			errShouldContain: []string{
				"can't get argument 2, too few parameters",
			},
		},
		{
			name:        "bad - 2 params, get with idx: 99",
			ce:          callExprs[2].ce,
			n:           99,
			errExpected: true,
			errShouldContain: []string{
				"can't get argument 99, too few parameters",
			},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)

		_, err := getArg(tc.ce, tc.n)
		testhelper.CheckError(t, tcID, err, tc.errExpected, tc.errShouldContain)
	}
}

func TestGetArgAsInt(t *testing.T) {
	callExprs := []struct {
		s  string
		ce *ast.CallExpr
	}{
		{s: `callFunc(1)`},
		{s: `callFunc(1, "hello")`},
	}
	for i, ce := range callExprs {
		expr, err := parser.ParseExpr(ce.s)
		if err != nil {
			t.Fatal("cannot parse the expression: ", ce.s, " error: ", err)
		}

		callExpr, ok := expr.(*ast.CallExpr)
		if !ok {
			t.Fatalf("the expression is not a call: %T", expr)
		}
		callExprs[i].ce = callExpr
	}

	testCases := []struct {
		name             string
		ce               *ast.CallExpr
		n                int
		expVal           int64
		errExpected      bool
		errShouldContain []string
	}{
		{
			name:   "good - 1 params, get with idx: 0",
			ce:     callExprs[0].ce,
			n:      0,
			expVal: 1,
		},
		{
			name:        "bad - 1 params, get with idx: 1",
			ce:          callExprs[0].ce,
			n:           1,
			errExpected: true,
			errShouldContain: []string{
				"can't get argument 1, too few parameters",
			},
		},
		{
			name:   "good - 2 params, get with idx: 0",
			ce:     callExprs[1].ce,
			n:      0,
			expVal: 1,
		},
		{
			name:        "bad - 2 params, get with idx: 1 - not an int",
			ce:          callExprs[1].ce,
			n:           1,
			errExpected: true,
			errShouldContain: []string{
				"can't convert argument 1 to an int:",
				" the expression should have been a literal INT, was STRING",
			},
		},
		{
			name:        "bad - 2 params, get with idx: 2",
			ce:          callExprs[1].ce,
			n:           2,
			errExpected: true,
			errShouldContain: []string{
				"can't get argument 2, too few parameters",
			},
		},
		{
			name:        "bad - 2 params, get with idx: 99",
			ce:          callExprs[1].ce,
			n:           99,
			errExpected: true,
			errShouldContain: []string{
				"can't get argument 99, too few parameters",
			},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)

		val, err := getArgAsInt(tc.ce, tc.n)
		if testhelper.CheckError(t, tcID, err,
			tc.errExpected, tc.errShouldContain) &&
			err == nil {
			if val != tc.expVal {
				t.Log(tcID)
				t.Logf("\t: expected: %d", tc.expVal)
				t.Logf("\t:      got: %d", val)
				t.Error("\t: unexpected value")
			}
		}
	}
}

func TestGetArgAsFloat(t *testing.T) {
	callExprs := []struct {
		s  string
		ce *ast.CallExpr
	}{
		{s: `callFunc(1.1)`},
		{s: `callFunc(1.1, "hello")`},
	}
	for i, ce := range callExprs {
		expr, err := parser.ParseExpr(ce.s)
		if err != nil {
			t.Fatal("cannot parse the expression: ", ce.s, " error: ", err)
		}

		callExpr, ok := expr.(*ast.CallExpr)
		if !ok {
			t.Fatalf("the expression is not a call: %T", expr)
		}
		callExprs[i].ce = callExpr
	}

	testCases := []struct {
		name             string
		ce               *ast.CallExpr
		n                int
		expVal           float64
		errExpected      bool
		errShouldContain []string
	}{
		{
			name:   "good - 1 params, get with idx: 0",
			ce:     callExprs[0].ce,
			n:      0,
			expVal: 1.1,
		},
		{
			name:        "bad - 1 params, get with idx: 1",
			ce:          callExprs[0].ce,
			n:           1,
			errExpected: true,
			errShouldContain: []string{
				"can't get argument 1, too few parameters",
			},
		},
		{
			name:   "good - 2 params, get with idx: 0",
			ce:     callExprs[1].ce,
			n:      0,
			expVal: 1.1,
		},
		{
			name:        "bad - 2 params, get with idx: 1 - not a float",
			ce:          callExprs[1].ce,
			n:           1,
			errExpected: true,
			errShouldContain: []string{
				"can't convert argument 1 to a float:",
				" the expression should have been a literal FLOAT, was STRING",
			},
		},
		{
			name:        "bad - 2 params, get with idx: 2",
			ce:          callExprs[1].ce,
			n:           2,
			errExpected: true,
			errShouldContain: []string{
				"can't get argument 2, too few parameters",
			},
		},
		{
			name:        "bad - 2 params, get with idx: 99",
			ce:          callExprs[1].ce,
			n:           99,
			errExpected: true,
			errShouldContain: []string{
				"can't get argument 99, too few parameters",
			},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)

		val, err := getArgAsFloat(tc.ce, tc.n)
		if testhelper.CheckError(t, tcID, err,
			tc.errExpected, tc.errShouldContain) &&
			err == nil {
			if val != tc.expVal {
				t.Log(tcID)
				t.Logf("\t: expected: %f", tc.expVal)
				t.Logf("\t:      got: %f", val)
				t.Error("\t: unexpected value")
			}
		}
	}
}
