package checksetter

import (
	"fmt"
	"go/ast"
	"go/parser"
)

// getElts returns a slice of expressions. If s does not represent an array
// of ast.Expr's then a non-nil error is returned.
func getElts(s, desc string) ([]ast.Expr, error) {
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
