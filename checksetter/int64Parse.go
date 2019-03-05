package checksetter

import (
	"fmt"
	"go/ast"
	"go/parser"

	"github.com/nickwells/check.mod/check"
)

const (
	int64CFName = "check.Int64"
	int64CFDesc = "int check func"
)

var int64CFInt = map[string]func(int64) check.Int64{
	"EQ":          check.Int64EQ,
	"GT":          check.Int64GT,
	"GE":          check.Int64GE,
	"LT":          check.Int64LT,
	"LE":          check.Int64LE,
	"Divides":     check.Int64Divides,
	"IsAMultiple": check.Int64IsAMultiple,
}

var int64CFIntInt = map[string]func(int64, int64) check.Int64{
	"Between": check.Int64Between,
}

var int64CFInt64CFStr = map[string]func(check.Int64, string) check.Int64{
	"Not": check.Int64Not,
}

var int64CFInt64CFList = map[string]func(...check.Int64) check.Int64{
	"And": check.Int64And,
	"Or":  check.Int64Or,
}

// makeInt64CFInt returns an Int64 checker corresponding to the
// given name - this is for checkers that take a single integer parameter
func makeInt64CFInt(name string, i int64) check.Int64 {
	if f, ok := int64CFInt[name]; ok {
		return f(i)
	}
	return nil
}

// makeInt64CFIntInt returns an Int64 checker corresponding to the
// given name - this is for checkers that take two integer parameters
func makeInt64CFIntInt(name string, i, j int64) check.Int64 {
	if f, ok := int64CFIntInt[name]; ok {
		return f(i, j)
	}
	return nil
}

// makeInt64CFInt64CFStr returns an Int64 checker corresponding to the given name
// - this is for checkers that take an Int64 check func and a string
// parameter
func makeInt64CFInt64CFStr(name string, cf check.Int64, s string) check.Int64 {
	if f, ok := int64CFInt64CFStr[name]; ok {
		return f(cf, s)
	}
	return nil
}

// makeInt64CFInt64CFList returns an Int64 checker corresponding to the given
// name - this is for checkers that take a list of int64 check funcs
func makeInt64CFInt64CFList(name string, cf ...check.Int64) check.Int64 {
	if f, ok := int64CFInt64CFList[name]; ok {
		return f(cf...)
	}
	return nil
}

// int64CFParse returns a slice of int64 check functions and a nil error
// if the string is successfully parsed or nil and an error if the string
// couldn't be converted to a slice of check functions.
func int64CFParse(s string) ([]check.Int64, error) {
	expr, err := parser.ParseExpr("[]T{\n" + s + "}")
	if err != nil {
		return nil, err
	}

	v := make([]check.Int64, 0, 1)
	cl, ok := expr.(*ast.CompositeLit)
	if !ok {
		return nil,
			fmt.Errorf("unexpected type for the collection of %s: %T",
				int64CFDesc, expr)
	}
	_, ok = cl.Type.(*ast.ArrayType)
	if !ok {
		return nil,
			fmt.Errorf("unexpected type for the array of %s: %T",
				int64CFDesc, cl.Type)
	}

	for _, elt := range cl.Elts {
		f, err := getFuncInt64CF(elt)
		if err != nil {
			return nil, fmt.Errorf("bad function: %s", err)
		}
		v = append(v, f)
	}

	return v, nil
}

// getFuncInt64CF will process the expression and return an Int64 checker or
// nil
func getFuncInt64CF(elt ast.Expr) (check.Int64, error) {
	switch e := elt.(type) {
	case *ast.CallExpr:
		return callInt64CFMaker(e)
	}
	return nil, fmt.Errorf("unexpected type: %T", elt)
}

// callInt64CFMaker calls the appropriate makeInt64CF... and returns the
// results
func callInt64CFMaker(e *ast.CallExpr) (cf check.Int64, err error) {
	defer func() {
		if r := recover(); r != nil {
			cf = nil
			err = fmt.Errorf("Cannot create the %s func: %v", int64CFName, r)
		}
	}()

	fd, err := getFuncDetails(e, int64CFName)
	if err != nil {
		return nil, err
	}

	var f check.Int64

	switch fd.expectedArgs {
	case "int":
		i, err := getArgAsInt(e, fd, 0)
		if err != nil {
			return nil, err
		}
		f = makeInt64CFInt(fd.name, i)
		if f == nil {
			return nil, fmt.Errorf("cannot create the %s: %s(%d)",
				int64CFDesc, fd.name, i)
		}
	case "int, int":
		i, err := getArgAsInt(e, fd, 0)
		if err != nil {
			return nil, err
		}
		j, err := getArgAsInt(e, fd, 1)
		if err != nil {
			return nil, err
		}
		f = makeInt64CFIntInt(fd.name, i, j)
		if f == nil {
			return nil, fmt.Errorf("cannot create the %s: %s(%d, %d)",
				int64CFDesc, fd.name, i, j)
		}
	case int64CFName + ", string":
		argExpr, err := getArg(e, fd, 0)
		if err != nil {
			return nil, fmt.Errorf(
				"couldn't create the %s argument for the %s: %s(...): %s ",
				int64CFDesc, int64CFDesc, fd.name, err)
		}
		cssf, err := getFuncInt64CF(argExpr)
		if err != nil {
			return nil, fmt.Errorf(
				"couldn't create the %s argument for the %s: %s(...): %s ",
				int64CFDesc, int64CFDesc, fd.name, err)
		}
		s, err := getArgAsString(e, fd, 1)
		if err != nil {
			return nil, err
		}
		f = makeInt64CFInt64CFStr(fd.name, cssf, s)
		if f == nil {
			return nil, fmt.Errorf("cannot create the %s: %s(...)",
				int64CFDesc, fd.name)
		}
	case int64CFName + " ...":
		scfArgs := make([]check.Int64, 0, len(e.Args))
		for i, argExpr := range e.Args {
			scf, err := getFuncInt64CF(argExpr)
			if err != nil {
				return nil, fmt.Errorf(
					"couldn't create the %s argument (%d) for the %s: %s(...): %s ",
					int64CFDesc, i, int64CFDesc, fd.name, err)
			}
			scfArgs = append(scfArgs, scf)
		}
		f = makeInt64CFInt64CFList(fd.name, scfArgs...)
		if f == nil {
			return nil, fmt.Errorf("cannot create the %s: %s(...)",
				int64CFDesc, fd.name)
		}
	default:
		return nil, fmt.Errorf("unexpected argument list: %s", fd.expectedArgs)
	}
	return f, nil
}
