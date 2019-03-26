package checksetter

import (
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
		testhelper.ID
		testhelper.ExpErr
		ce *ast.CallExpr
		n  int
	}{
		{
			ID: testhelper.MkID("good - 0 params"),
			ce: callExprs[0].ce,
			n:  0,
		},
		{
			ID: testhelper.MkID("bad - 0 params"),
			ce: callExprs[0].ce,
			n:  99,
			ExpErr: testhelper.MkExpErr(
				"the call has 0 arguments, it should have 99"),
		},
		{
			ID: testhelper.MkID("good - 1 params"),
			ce: callExprs[1].ce,
			n:  1,
		},
		{
			ID: testhelper.MkID("bad - 1 params"),
			ce: callExprs[1].ce,
			n:  99,
			ExpErr: testhelper.MkExpErr(
				"the call has 1 arguments, it should have 99"),
		},
		{
			ID: testhelper.MkID("good - 2 params"),
			ce: callExprs[2].ce,
			n:  2,
		},
		{
			ID: testhelper.MkID("bad - 2 params"),
			ce: callExprs[2].ce,
			n:  99,
			ExpErr: testhelper.MkExpErr(
				"the call has 2 arguments, it should have 99"),
		},
	}

	for _, tc := range testCases {
		err := checkArgCount(tc.ce, tc.n)
		testhelper.CheckExpErr(t, err, tc)
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
		testhelper.ID
		testhelper.ExpErr
		ce *ast.CallExpr
		n  int
	}{
		{
			ID: testhelper.MkID("bad - 0 params, get with idx: 0"),
			ce: callExprs[0].ce,
			n:  0,
			ExpErr: testhelper.MkExpErr(
				"can't get argument 0, too few parameters"),
		},
		{
			ID: testhelper.MkID("bad - 0 params, get with idx: 1"),
			ce: callExprs[0].ce,
			n:  1,
			ExpErr: testhelper.MkExpErr(
				"can't get argument 1, too few parameters"),
		},
		{
			ID: testhelper.MkID("good - 1 params, get with idx: 0"),
			ce: callExprs[1].ce,
			n:  0,
		},
		{
			ID: testhelper.MkID("bad - 1 params, get with idx: 1"),
			ce: callExprs[1].ce,
			n:  1,
			ExpErr: testhelper.MkExpErr(
				"can't get argument 1, too few parameters"),
		},
		{
			ID: testhelper.MkID("good - 2 params, get with idx: 0"),
			ce: callExprs[2].ce,
			n:  0,
		},
		{
			ID: testhelper.MkID("good - 2 params, get with idx: 1"),
			ce: callExprs[2].ce,
			n:  1,
		},
		{
			ID: testhelper.MkID("bad - 2 params, get with idx: 2"),
			ce: callExprs[2].ce,
			n:  2,
			ExpErr: testhelper.MkExpErr(
				"can't get argument 2, too few parameters"),
		},
		{
			ID: testhelper.MkID("bad - 2 params, get with idx: 99"),
			ce: callExprs[2].ce,
			n:  99,
			ExpErr: testhelper.MkExpErr(
				"can't get argument 99, too few parameters"),
		},
	}

	for _, tc := range testCases {
		_, err := getArg(tc.ce, tc.n)
		testhelper.CheckExpErr(t, err, tc)
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
		testhelper.ID
		testhelper.ExpErr
		ce     *ast.CallExpr
		n      int
		expVal int64
	}{
		{
			ID:     testhelper.MkID("good - 1 params, get with idx: 0"),
			ce:     callExprs[0].ce,
			n:      0,
			expVal: 1,
		},
		{
			ID: testhelper.MkID("bad - 1 params, get with idx: 1"),
			ce: callExprs[0].ce,
			n:  1,
			ExpErr: testhelper.MkExpErr(
				"can't get argument 1, too few parameters"),
		},
		{
			ID:     testhelper.MkID("good - 2 params, get with idx: 0"),
			ce:     callExprs[1].ce,
			n:      0,
			expVal: 1,
		},
		{
			ID: testhelper.MkID("bad - 2 params, get with idx: 1 - not an int"),
			ce: callExprs[1].ce,
			n:  1,
			ExpErr: testhelper.MkExpErr(
				"can't convert argument 1 to an int:",
				" the expression should have been a literal INT, was STRING"),
		},
		{
			ID: testhelper.MkID("bad - 2 params, get with idx: 2"),
			ce: callExprs[1].ce,
			n:  2,
			ExpErr: testhelper.MkExpErr(
				"can't get argument 2, too few parameters"),
		},
		{
			ID: testhelper.MkID("bad - 2 params, get with idx: 99"),
			ce: callExprs[1].ce,
			n:  99,
			ExpErr: testhelper.MkExpErr(
				"can't get argument 99, too few parameters"),
		},
	}

	for _, tc := range testCases {
		val, err := getArgAsInt(tc.ce, tc.n)
		if testhelper.CheckExpErr(t, err, tc) &&
			err == nil {
			if val != tc.expVal {
				t.Log(tc.IDStr())
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
		testhelper.ID
		testhelper.ExpErr
		ce     *ast.CallExpr
		n      int
		expVal float64
	}{
		{
			ID:     testhelper.MkID("good - 1 params, get with idx: 0"),
			ce:     callExprs[0].ce,
			n:      0,
			expVal: 1.1,
		},
		{
			ID: testhelper.MkID("bad - 1 params, get with idx: 1"),
			ce: callExprs[0].ce,
			n:  1,
			ExpErr: testhelper.MkExpErr(
				"can't get argument 1, too few parameters"),
		},
		{
			ID:     testhelper.MkID("good - 2 params, get with idx: 0"),
			ce:     callExprs[1].ce,
			n:      0,
			expVal: 1.1,
		},
		{
			ID: testhelper.MkID("bad - 2 params, get with idx: 1 - not a float"),
			ce: callExprs[1].ce,
			n:  1,
			ExpErr: testhelper.MkExpErr(
				"can't convert argument 1 to a float:",
				" the expression should have been a literal FLOAT, was STRING"),
		},
		{
			ID: testhelper.MkID("bad - 2 params, get with idx: 2"),
			ce: callExprs[1].ce,
			n:  2,
			ExpErr: testhelper.MkExpErr(
				"can't get argument 2, too few parameters"),
		},
		{
			ID: testhelper.MkID("bad - 2 params, get with idx: 99"),
			ce: callExprs[1].ce,
			n:  99,
			ExpErr: testhelper.MkExpErr(
				"can't get argument 99, too few parameters"),
		},
	}

	for _, tc := range testCases {
		val, err := getArgAsFloat(tc.ce, tc.n)
		if testhelper.CheckExpErr(t, err, tc) &&
			err == nil {
			if val != tc.expVal {
				t.Log(tc.IDStr())
				t.Logf("\t: expected: %f", tc.expVal)
				t.Logf("\t:      got: %f", val)
				t.Error("\t: unexpected value")
			}
		}
	}
}
