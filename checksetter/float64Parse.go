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

type float64CFMaker func(*ast.CallExpr, string) (check.Float64, error)

var float64CFArgsToFunc map[string]float64CFMaker

func init() {
	float64CFArgsToFunc = map[string]float64CFMaker{
		"float":                    makeFloat64CFFloat,
		"float, float":             makeFloat64CFFloatFloat,
		float64CFName + ", string": makeFloat64CFFloat64CFStr,
		float64CFName + " ...":     makeFloat64CFFloat64CFList,
	}
}

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
func makeFloat64CFFloat(e *ast.CallExpr, fName string) (cf check.Float64, err error) {
	var v float64
	errIntro := func() string {
		return fmt.Sprintf("can't make the %s func: %s(%f):",
			float64CFName, fName, v)
	}
	defer func() {
		if r := recover(); r != nil {
			cf = nil
			err = fmt.Errorf("%s %v", errIntro(), r)
		}
	}()

	if err = checkArgCount(e, 1); err != nil {
		return nil, fmt.Errorf("%s %s", errIntro(), err)
	}

	v, err = getArgAsFloat(e, 0)
	if err != nil {
		return nil, fmt.Errorf("%s %s", errIntro(), err)
	}

	if f, ok := float64CFFloat[fName]; ok {
		return f(v), nil
	}

	return nil, fmt.Errorf("%s the name is not recognised", errIntro())
}

// makeFloat64CFFloatFloat returns a Float64 checker corresponding to the
// given name - this is for checkers that take two float parameters
func makeFloat64CFFloatFloat(e *ast.CallExpr, fName string) (cf check.Float64, err error) {
	var v, w float64
	errIntro := func() string {
		return fmt.Sprintf("can't make the %s func: %s(%f, %f):",
			float64CFName, fName, v, w)
	}
	defer func() {
		if r := recover(); r != nil {
			cf = nil
			err = fmt.Errorf("%s %v", errIntro(), r)
		}
	}()

	if err = checkArgCount(e, 2); err != nil {
		return nil, fmt.Errorf("%s %s", errIntro(), err)
	}

	v, err = getArgAsFloat(e, 0)
	if err != nil {
		return nil, fmt.Errorf("%s %s", errIntro(), err)
	}

	w, err = getArgAsFloat(e, 0)
	if err != nil {
		return nil, fmt.Errorf("%s %s", errIntro(), err)
	}

	if f, ok := float64CFFloatFloat[fName]; ok {
		return f(v, w), nil
	}

	return nil, fmt.Errorf("%s the name is not recognised", errIntro())
}

// makeFloat64CFFloat64CFStr returns a Float64 checker corresponding to the
// given name - this is for checkers that take a Float64 check func and a
// string parameter
func makeFloat64CFFloat64CFStr(e *ast.CallExpr, fName string) (cf check.Float64, err error) {
	var s string
	errIntro := func() string {
		return fmt.Sprintf("can't make the %s func: %s(%s, %s):",
			float64CFName, fName, float64CFName, s)
	}
	defer func() {
		if r := recover(); r != nil {
			cf = nil
			err = fmt.Errorf("%s %v", errIntro(), r)
		}
	}()

	if err = checkArgCount(e, 2); err != nil {
		return nil, fmt.Errorf("%s %s", errIntro(), err)
	}

	argExpr, err := getArg(e, 0)
	if err != nil {
		return nil, fmt.Errorf("%s can't get the %s argument: %s",
			errIntro(), float64CFName, err)
	}
	fcf, err := getFuncFloat64CF(argExpr)
	if err != nil {
		return nil, fmt.Errorf("%s can't convert argument %d to %s: %s",
			errIntro(), 0, float64CFName, err)
	}
	s, err = getArgAsString(e, 1)
	if err != nil {
		return nil, fmt.Errorf("%s %s", errIntro(), err)
	}

	if f, ok := float64CFFloat64CFStr[fName]; ok {
		return f(fcf, s), nil
	}

	return nil, fmt.Errorf("%s the name is not recognised", errIntro())
}

// makeFloat64CFFloat64CFList returns a Float64 checker corresponding to the
// given name - this is for checkers that take a list of float check funcs
func makeFloat64CFFloat64CFList(e *ast.CallExpr, fName string) (cf check.Float64, err error) {
	errIntro := func() string {
		return fmt.Sprintf("can't make the %s func: %s(%s ...):",
			float64CFName, fName, float64CFName)
	}
	defer func() {
		if r := recover(); r != nil {
			cf = nil
			err = fmt.Errorf("%s %v", errIntro(), r)
		}
	}()

	fArgs := make([]check.Float64, 0, len(e.Args))
	for i, argExpr := range e.Args {
		fcf, err := getFuncFloat64CF(argExpr)
		if err != nil {
			return nil, fmt.Errorf("%s can't convert argument %d to %s: %s",
				errIntro(), i, float64CFName, err)
		}
		fArgs = append(fArgs, fcf)
	}

	if f, ok := float64CFFloat64CFList[fName]; ok {
		return f(fArgs...), nil
	}

	return nil, fmt.Errorf("%s the name is not recognised", errIntro())
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
		return nil, fmt.Errorf("unexpected type for the collection of %s: %T",
			float64CFDesc, expr)
	}
	_, ok = cl.Type.(*ast.ArrayType)
	if !ok {
		return nil, fmt.Errorf("unexpected type for the array of %s: %T",
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
func getFuncFloat64CF(elt ast.Expr) (cf check.Float64, err error) {
	e, ok := elt.(*ast.CallExpr)
	if !ok {
		return nil, fmt.Errorf("unexpected type: %T", elt)
	}

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

	switch fd.expectedArgs {
	case "float":
		return makeFloat64CFFloat(e, fd.name)
	case "float, float":
		return makeFloat64CFFloatFloat(e, fd.name)
	case float64CFName + ", string":
		return makeFloat64CFFloat64CFStr(e, fd.name)
	case float64CFName + " ...":
		return makeFloat64CFFloat64CFList(e, fd.name)
	default:
		return nil, fmt.Errorf("%s has an unexpected argument list: %s",
			fd.name, fd.expectedArgs)
	}
}
