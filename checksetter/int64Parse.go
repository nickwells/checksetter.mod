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

type int64CFMaker func(*ast.CallExpr, string) (check.Int64, error)

var int64CFArgsToFunc map[string]int64CFMaker

func init() {
	int64CFArgsToFunc = map[string]int64CFMaker{
		"int":                    makeInt64CFInt,
		"int, int":               makeInt64CFIntInt,
		int64CFName + ", string": makeInt64CFInt64CFStr,
		int64CFName + " ...":     makeInt64CFInt64CFList,
	}
}

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

// makeInt64CFInt returns an Int64 checker corresponding to the given name -
// this is for checkers that take a single integer parameter
func makeInt64CFInt(e *ast.CallExpr, fName string) (cf check.Int64, err error) {
	var i int64
	errIntro := func() string {
		return fmt.Sprintf("can't make the %s func: %s(%d):",
			int64CFName, fName, i)
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

	i, err = getArgAsInt(e, 0)
	if err != nil {
		return nil, fmt.Errorf("%s %s", errIntro(), err)
	}

	if f, ok := int64CFInt[fName]; ok {
		return f(i), nil
	}

	return nil, fmt.Errorf("%s the name is not recognised", errIntro())
}

// makeInt64CFIntInt returns an Int64 checker corresponding to the given name
// - this is for checkers that take two integer parameters
func makeInt64CFIntInt(e *ast.CallExpr, fName string) (cf check.Int64, err error) {
	var i, j int64
	errIntro := func() string {
		return fmt.Sprintf("can't make the %s func: %s(%d, %d):",
			int64CFName, fName, i, j)
	}
	defer func() {
		if r := recover(); r != nil {
			cf = nil
			err = fmt.Errorf("%s %v",
				errIntro(), r)
		}
	}()

	if err = checkArgCount(e, 2); err != nil {
		return nil, fmt.Errorf("%s %s", errIntro(), err)
	}

	i, err = getArgAsInt(e, 0)
	if err != nil {
		return nil, fmt.Errorf("%s %s", errIntro(), err)
	}
	j, err = getArgAsInt(e, 1)
	if err != nil {
		return nil, fmt.Errorf("%s %s", errIntro(), err)
	}

	if f, ok := int64CFIntInt[fName]; ok {
		return f(i, j), nil
	}

	return nil, fmt.Errorf("%s the name is not recognised", errIntro())
}

// makeInt64CFInt64CFStr returns an Int64 checker corresponding to the given
// name - this is for checkers that take an Int64 check func and a string
// parameter
func makeInt64CFInt64CFStr(e *ast.CallExpr, fName string) (cf check.Int64, err error) {
	var s string
	errIntro := func() string {
		return fmt.Sprintf("can't make the %s func: %s(%s, %s):",
			int64CFName, fName, int64CFName, s)
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
			errIntro(), int64CFName, err)
	}
	icf, err := getFuncInt64CF(argExpr)
	if err != nil {
		return nil, fmt.Errorf("%s can't convert argument %d to %s: %s",
			errIntro(), 0, int64CFName, err)
	}
	s, err = getArgAsString(e, 1)
	if err != nil {
		return nil, fmt.Errorf("%s %s", errIntro(), err)
	}

	if f, ok := int64CFInt64CFStr[fName]; ok {
		return f(icf, s), nil
	}

	return nil, fmt.Errorf("%s the name is not recognised", errIntro())
}

// makeInt64CFInt64CFList returns an Int64 checker corresponding to the given
// name - this is for checkers that take a list of int64 check funcs
func makeInt64CFInt64CFList(e *ast.CallExpr, fName string) (cf check.Int64, err error) {
	errIntro := "can't make the " + int64CFName +
		" func: " + fName + "(" + int64CFName + " ...):"
	defer func() {
		if r := recover(); r != nil {
			cf = nil
			err = fmt.Errorf("%s %v", errIntro, r)
		}
	}()

	fArgs := make([]check.Int64, 0, len(e.Args))
	for i, argExpr := range e.Args {
		scf, err := getFuncInt64CF(argExpr)
		if err != nil {
			return nil, fmt.Errorf("%s can't convert argument %d to %s: %s",
				errIntro, i, int64CFName, err)
		}
		fArgs = append(fArgs, scf)
	}

	if f, ok := int64CFInt64CFList[fName]; ok {
		return f(fArgs...), nil
	}

	return nil, fmt.Errorf("%s the name is not recognised", errIntro)
}

// int64CFParse returns a slice of int64 check functions and a nil error if
// the string is successfully parsed or nil and an error if the string
// couldn't be converted to a slice of check functions.
func int64CFParse(s string) ([]check.Int64, error) {
	expr, err := parser.ParseExpr("[]T{\n" + s + "}")
	if err != nil {
		return nil, err
	}

	v := make([]check.Int64, 0, 1)
	cl, ok := expr.(*ast.CompositeLit)
	if !ok {
		return nil, fmt.Errorf("unexpected type for the collection of %s: %T",
			int64CFDesc, expr)
	}
	_, ok = cl.Type.(*ast.ArrayType)
	if !ok {
		return nil, fmt.Errorf("unexpected type for the array of %s: %T",
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
func getFuncInt64CF(elt ast.Expr) (cf check.Int64, err error) {
	e, ok := elt.(*ast.CallExpr)
	if !ok {
		return nil, fmt.Errorf("unexpected type: %T", elt)
	}

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

	makeF, ok := int64CFArgsToFunc[fd.expectedArgs]
	if ok {
		return makeF(e, fd.name)
	} else {
		return nil, fmt.Errorf("%s has an unrecognised argument list: %s",
			fd.name, fd.expectedArgs)
	}
}
