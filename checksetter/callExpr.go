package checksetter

import (
	"fmt"
	"go/ast"
)

// funcDetails records details about a function grouping them into a single
// place for convenience
type funcDetails struct {
	name         string
	setName      string
	expectedArgs string
}

// getFuncDetails will return the function name and expected arguments and
// return any errors
func getFuncDetails(e *ast.CallExpr, funcSet string) (funcDetails, error) {
	fd := funcDetails{setName: funcSet}

	nameToArgs, ok := mapOfNameToArgMaps[funcSet]
	if !ok {
		return fd, fmt.Errorf("unknown function set: %s", funcSet)
	}

	fID, ok := e.Fun.(*ast.Ident)
	if !ok {
		return fd, fmt.Errorf("%s: syntax error: unexpected call type: %T",
			funcSet, e.Fun)
	}
	fd.name = fID.Name

	expArgs, ok := nameToArgs[fd.name]
	if !ok {
		return fd, fmt.Errorf("unknown %s: %s", funcSet, fd.name)
	}
	fd.expectedArgs = expArgs

	return fd, nil
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

// getArg retrieves the idx'th argument from the call expression
func getArg(e *ast.CallExpr, idx int) (ast.Expr, error) {
	if idx >= len(e.Args) {
		return nil, fmt.Errorf("can't get argument %d, too few parameters", idx)
	}
	return e.Args[idx], nil
}

// getArgAsInt64 gets an int64 arg from the list of arguments to the function
// returning an error if there are any problems
func getArgAsInt64(e *ast.CallExpr, idx int) (int64, error) {
	argExpr, err := getArg(e, idx)
	if err != nil {
		return 0, err
	}
	i, err := getInt64(argExpr)
	if err != nil {
		return 0, fmt.Errorf("can't convert argument %d to an int64: %s",
			idx, err)
	}
	return i, nil
}

// getArgAsInt gets an int arg from the list of arguments to the function
// returning an error if there are any problems
func getArgAsInt(e *ast.CallExpr, idx int) (int, error) {
	argExpr, err := getArg(e, idx)
	if err != nil {
		return 0, err
	}
	i, err := getInt(argExpr)
	if err != nil {
		return 0, fmt.Errorf("can't convert argument %d to an int: %s",
			idx, err)
	}
	return i, nil
}

// getArgAsFloat gets a float64 arg from the list of arguments to the function
// returning an error if there are any problems
func getArgAsFloat(e *ast.CallExpr, idx int) (float64, error) {
	argExpr, err := getArg(e, idx)
	if err != nil {
		return 0.0, err
	}
	f, err := getFloat(argExpr)
	if err != nil {
		return 0.0, fmt.Errorf("can't convert argument %d to a float: %s",
			idx, err)
	}
	return f, nil
}
