package checksetter

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"

	"github.com/nickwells/check.mod/v2/check"
)

// getInt64 converts the expression which is expected to be a BasicLit into the
// corresponding int64
func getInt64(e ast.Expr) (int64, error) {
	v, ok := e.(*ast.BasicLit)
	if !ok {
		return 0, fmt.Errorf("the expression isn't a BasicLit, it's a %T", e)
	}

	if v.Kind != token.INT {
		return 0, fmt.Errorf("%q isn't an INT, it's a %s", v.Value, v.Kind)
	}
	i, err := strconv.ParseInt(v.Value, 0, 64)
	if err != nil {
		return 0, fmt.Errorf("Couldn't make an int from %q: %s", v.Value, err)
	}
	return i, nil
}

// getInt converts the expression which is expected to be a BasicLit into the
// corresponding int
func getInt(e ast.Expr) (int, error) {
	i, err := getInt64(e)
	return int(i), err
}

// getFloat64 converts the expression which is expected to be a BasicLit into
// the corresponding float64
func getFloat64(e ast.Expr) (float64, error) {
	v, ok := e.(*ast.BasicLit)
	if !ok {
		return 0, fmt.Errorf("the expression isn't a BasicLit, it's a %T", e)
	}

	if v.Kind != token.FLOAT && v.Kind != token.INT {
		return 0, fmt.Errorf("%q isn't a FLOAT/INT, it's a %s", v.Value, v.Kind)
	}
	f, err := strconv.ParseFloat(v.Value, 64)
	if err != nil {
		return 0, fmt.Errorf("Couldn't make a float %q: %s", v.Value, err)
	}
	return f, nil
}

// getString converts the expression which is expected to be a BasicLit into the
// corresponding string
func getString(e ast.Expr) (string, error) {
	v, ok := e.(*ast.BasicLit)
	if !ok {
		return "", fmt.Errorf("the expression isn't a BasicLit, it's a %T", e)
	}

	if v.Kind != token.STRING {
		return "", fmt.Errorf("%q isn't a STRING, it's a %s", v.Value, v.Kind)
	}
	return strings.Trim(v.Value, `"`), nil
}

// getElts returns a slice of expressions. If s does not represent an array
// of ast.Expr's then a non-nil error is returned.
func getElts(s, desc string) (e []ast.Expr, err error) {
	defer func() {
		if r := recover(); r != nil {
			e, err = nil, fmt.Errorf("%s: unexpected parse error: %q", desc, r)
		}
	}()

	expr, err := parser.ParseExpr("[]T{\n" + s + "}")
	if err != nil {
		return nil, err
	}

	cl, ok := expr.(*ast.CompositeLit)
	if !ok {
		return nil, fmt.Errorf("unexpected type for the collection of %s: %T",
			desc, expr)
	}
	_, ok = cl.Type.(*ast.ArrayType)
	if !ok {
		return nil, fmt.Errorf("unexpected type for the array of %s: %T",
			desc, cl.Type)
	}

	return cl.Elts, nil
}

// checkArgCount will return an error if the number of arguments in the
// CallExpr is not equal to the given value, nil otherwise
func checkArgCount(e *ast.CallExpr, n int) error {
	if len(e.Args) != n {
		return fmt.Errorf("the call has %d arguments, it should have %d",
			len(e.Args), n)
	}
	return nil
}

// getCheckFuncs[T any] returns a slice of check-funcs from the CallExpr
func getCheckFuncs[T any](e *ast.CallExpr, checkerName string) (
	[]check.ValCk[T], error,
) {
	parser, err := FindParser[T](checkerName)
	if err != nil {
		return nil, err
	}

	checkFuncs := make([]check.ValCk[T], 0, len(e.Args))
	for i, expr := range e.Args {
		cf, err := parser.ParseExpr(expr)
		if err != nil {
			return nil,
				fmt.Errorf("can't convert argument %d to %s: %s",
					i, checkerName, err)
		}
		checkFuncs = append(checkFuncs, cf)
	}

	return checkFuncs, nil
}

// getCheckFunc[T any] returns a single check-func from the CallExpr argument
// at index idx
func getCheckFunc[T any](e *ast.CallExpr, idx int, checkerName string) (
	check.ValCk[T], error,
) {
	if idx < 0 {
		return nil, fmt.Errorf("Index (%d) must be >= 0", idx)
	}
	if idx >= len(e.Args) {
		return nil,
			fmt.Errorf("Index (%d) is too large, there are only %d arguments",
				idx, len(e.Args))
	}
	parser, err := FindParser[T](checkerName)
	if err != nil {
		return nil, err
	}

	cf, err := parser.ParseExpr(e.Args[idx])
	if err != nil {
		return nil,
			fmt.Errorf("can't convert argument %d to %s: %s",
				idx, checkerName, err)
	}

	return cf, nil
}
