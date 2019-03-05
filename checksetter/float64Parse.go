package checksetter

import (
	"fmt"
	"go/ast"
	"go/parser"

	"github.com/nickwells/check.mod/check"
)

const (
	float64CFName = "check.Float64"
	float64CFDesc = "float check func"
)

var float64CFFloat = map[string]func(float64) check.Float64{
	"GT": check.Float64GT,
	"GE": check.Float64GE,
	"LT": check.Float64LT,
	"LE": check.Float64LE,
}

var float64CFFloatFloat = map[string]func(float64, float64) check.Float64{
	"Between": check.Float64Between,
}

var float64CFFloat64CFStr = map[string]func(check.Float64, string) check.Float64{
	"Not": check.Float64Not,
}

var float64CFFloat64CFList = map[string]func(...check.Float64) check.Float64{
	"And": check.Float64And,
	"Or":  check.Float64Or,
}

// makeFloat64CFFloat returns a Float64 checker corresponding to the
// given name - this is for checkers that take a single float parameter
func makeFloat64CFFloat(name string, i float64) check.Float64 {
	if f, ok := float64CFFloat[name]; ok {
		return f(i)
	}
	return nil
}

// makeFloat64CFFloatFloat returns a Float64 checker corresponding to the
// given name - this is for checkers that take two float parameters
func makeFloat64CFFloatFloat(name string, i, j float64) check.Float64 {
	if f, ok := float64CFFloatFloat[name]; ok {
		return f(i, j)
	}
	return nil
}

// makeFloat64CFFloat64CFStr returns a Float64 checker corresponding to the
// given name - this is for checkers that take a Float64 check func and a
// string parameter
func makeFloat64CFFloat64CFStr(name string, cf check.Float64, s string) check.Float64 {
	if f, ok := float64CFFloat64CFStr[name]; ok {
		return f(cf, s)
	}
	return nil
}

// makeFloat64CFFloat64CFList returns a Float64 checker corresponding to the
// given name - this is for checkers that take a list of float check funcs
func makeFloat64CFFloat64CFList(name string, cf ...check.Float64) check.Float64 {
	if f, ok := float64CFFloat64CFList[name]; ok {
		return f(cf...)
	}
	return nil
}

// float64CFParse returns a slice of float64 check functions and a nil error
// if the string is successfully parsed or nil and an error if the string
// couldn't be converted to a slice of check functions.
func float64CFParse(s string) ([]check.Float64, error) {
	expr, err := parser.ParseExpr("[]T{\n" + s + "}")
	if err != nil {
		return nil, err
	}

	v := make([]check.Float64, 0, 1)
	cl, ok := expr.(*ast.CompositeLit)
	if !ok {
		return nil,
			fmt.Errorf("unexpected type for the collection of %s: %T",
				float64CFDesc, expr)
	}
	_, ok = cl.Type.(*ast.ArrayType)
	if !ok {
		return nil,
			fmt.Errorf("unexpected type for the array of %s: %T",
				float64CFDesc, cl.Type)
	}

	for _, elt := range cl.Elts {
		f, err := getFuncFloat64CF(elt)
		if err != nil {
			return nil, fmt.Errorf("bad function: %s", err)
		}
		v = append(v, f)
	}

	return v, nil
}

// getFuncFloat64CF will process the expression and return a Float64 checker or
// nil
func getFuncFloat64CF(elt ast.Expr) (check.Float64, error) {
	switch e := elt.(type) {
	case *ast.CallExpr:
		return callFloat64CFMaker(e)
	}
	return nil, fmt.Errorf("unexpected type: %T", elt)
}

// callFloat64CFMaker calls the appropriate makeFloat64CF... and returns the
// results
func callFloat64CFMaker(e *ast.CallExpr) (cf check.Float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			cf = nil
			err = fmt.Errorf("Cannot create the %s func: %v", float64CFName, r)
		}
	}()

	fd, err := getFuncDetails(e, float64CFName)
	if err != nil {
		return nil, err
	}

	var f check.Float64

	switch fd.expectedArgs {
	case "float":
		i, err := getArgAsFloat(e, fd, 0)
		if err != nil {
			return nil, err
		}
		f = makeFloat64CFFloat(fd.name, i)
		if f == nil {
			return nil, fmt.Errorf("cannot create the %s: %s(%f)",
				float64CFDesc, fd.name, i)
		}
	case "float, float":
		i, err := getArgAsFloat(e, fd, 0)
		if err != nil {
			return nil, err
		}
		j, err := getArgAsFloat(e, fd, 1)
		if err != nil {
			return nil, err
		}
		f = makeFloat64CFFloatFloat(fd.name, i, j)
		if f == nil {
			return nil, fmt.Errorf("cannot create the %s: %s(%f, %f)",
				float64CFDesc, fd.name, i, j)
		}
	case float64CFName + ", string":
		argExpr, err := getArg(e, fd, 0)
		if err != nil {
			return nil, fmt.Errorf(
				"couldn't create the %s argument for the %s: %s(...): %s ",
				float64CFDesc, float64CFDesc, fd.name, err)
		}
		cssf, err := getFuncFloat64CF(argExpr)
		if err != nil {
			return nil, fmt.Errorf(
				"couldn't create the %s argument for the %s: %s(...): %s ",
				float64CFDesc, float64CFDesc, fd.name, err)
		}
		s, err := getArgAsString(e, fd, 1)
		if err != nil {
			return nil, err
		}
		f = makeFloat64CFFloat64CFStr(fd.name, cssf, s)
		if f == nil {
			return nil, fmt.Errorf("cannot create the %s: %s(...)",
				float64CFDesc, fd.name)
		}
	case float64CFName + " ...":
		scfArgs := make([]check.Float64, 0, len(e.Args))
		for i, argExpr := range e.Args {
			scf, err := getFuncFloat64CF(argExpr)
			if err != nil {
				return nil, fmt.Errorf(
					"couldn't create the %s argument (%d) for the %s: %s(...): %s ",
					float64CFDesc, i, float64CFDesc, fd.name, err)
			}
			scfArgs = append(scfArgs, scf)
		}
		f = makeFloat64CFFloat64CFList(fd.name, scfArgs...)
		if f == nil {
			return nil, fmt.Errorf("cannot create the %s: %s(...)",
				float64CFDesc, fd.name)
		}
	default:
		return nil, fmt.Errorf("unexpected argument list: %s", fd.expectedArgs)
	}
	return f, nil
}
