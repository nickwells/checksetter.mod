package checksetter

import (
	"fmt"
	"go/ast"
	"strings"
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

// getArg ...
func getArg(e *ast.CallExpr, fd funcDetails, idx int) (ast.Expr, error) {
	if idx >= len(e.Args) {
		return nil,
			fmt.Errorf(
				"couldn't get argument %d for the %s (%s), too few parameters",
				idx, fd.setName, fd.name)
	}
	return e.Args[idx], nil
}

// getArgAsInt gets an int64 arg from the list of arguments to the function
// returning an error if there are any problems
func getArgAsInt(e *ast.CallExpr, fd funcDetails, idx int) (int64, error) {
	argExpr, err := getArg(e, fd, idx)
	if err != nil {
		return 0, err
	}
	i, err := getInt(argExpr)
	if err != nil {
		return 0,
			fmt.Errorf(
				"couldn't get argument %d for the %s (%s) as an int: %s",
				idx, fd.setName, fd.name, err)
	}
	return i, nil
}

// getArgAsFloat gets a float64 arg from the list of arguments to the function
// returning an error if there are any problems
func getArgAsFloat(e *ast.CallExpr, fd funcDetails, idx int) (float64, error) {
	argExpr, err := getArg(e, fd, idx)
	if err != nil {
		return 0.0, err
	}
	f, err := getFloat(argExpr)
	if err != nil {
		return 0.0, fmt.Errorf(
			"couldn't get argument %d for the %s (%s) as a float: %s",
			idx, fd.setName, fd.name, err)
	}
	return f, nil
}

// getArgAsString gets a string arg from the list of arguments to the
// function returning an error if there are any problems
func getArgAsString(e *ast.CallExpr, fd funcDetails, idx int) (string, error) {
	argExpr, err := getArg(e, fd, idx)
	if err != nil {
		return "", err
	}
	s, err := getString(argExpr)
	if err != nil {
		return "",
			fmt.Errorf(
				"couldn't get argument %d for the %s (%s) as a string: %s",
				idx, fd.setName, fd.name, err)
	}
	return strings.Trim(s, "\""), nil
}
