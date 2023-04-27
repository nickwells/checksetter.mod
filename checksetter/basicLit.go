package checksetter

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

// getInt64 converts the expression which is expected to be a BasicLit into the
// corresponding int64
func getInt64(e ast.Expr) (int64, error) {
	v, ok := e.(*ast.BasicLit)
	if !ok {
		return 0,
			fmt.Errorf("the expression should have been a literal not %T", e)
	}

	if v.Kind != token.INT {
		return 0,
			fmt.Errorf("the expression should have been a literal INT, was %s",
				v.Kind)
	}
	i, err := strconv.ParseInt(v.Value, 0, 64)
	if err != nil {
		return 0,
			fmt.Errorf("Couldn't convert '%s' into an int64: %s",
				v.Value, err)
	}
	return i, nil
}

// getInt converts the expression which is expected to be a BasicLit into the
// corresponding int
func getInt(e ast.Expr) (int, error) {
	v, ok := e.(*ast.BasicLit)
	if !ok {
		return 0,
			fmt.Errorf("the expression should have been a literal not %T", e)
	}

	if v.Kind != token.INT {
		return 0,
			fmt.Errorf("the expression should have been a literal INT, was %s",
				v.Kind)
	}
	i, err := strconv.ParseInt(v.Value, 0, 64)
	if err != nil {
		return 0,
			fmt.Errorf("Couldn't convert '%s' into an int: %s",
				v.Value, err)
	}
	return int(i), nil
}

// getFloat converts the expression which is expected to be a BasicLit into the
// corresponding float
func getFloat(e ast.Expr) (float64, error) {
	v, ok := e.(*ast.BasicLit)
	if !ok {
		return 0,
			fmt.Errorf("the expression should have been a literal not %T", e)
	}

	if v.Kind != token.FLOAT {
		return 0,
			fmt.Errorf(
				"the expression should have been a literal FLOAT, was %s",
				v.Kind)
	}
	f, err := strconv.ParseFloat(v.Value, 64)
	if err != nil {
		return 0,
			fmt.Errorf("Couldn't convert '%s' into a float: %s",
				v.Value, err)
	}
	return f, nil
}

// getString converts the expression which is expected to be a BasicLit into the
// corresponding string
func getString(e ast.Expr) (string, error) {
	v, ok := e.(*ast.BasicLit)
	if !ok {
		return "",
			fmt.Errorf("the expression should have been a literal not %T", e)
	}

	if v.Kind != token.STRING {
		return "",
			fmt.Errorf(
				"the expression should have been a literal STRING, was %s",
				v.Kind)
	}
	return strings.Trim(v.Value, `"`), nil
}
