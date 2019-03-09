package checksetter

import (
	"fmt"
	"go/ast"
	"go/parser"

	"github.com/nickwells/check.mod/check"
)

const (
	strSlcCFName = "check.StringSlice"
	strSlcCFDesc = "string slice check func"
)

var strSlcCFNoParam = map[string]check.StringSlice{
	"NoDups": check.StringSliceNoDups,
}
var strSlcCFInt = map[string]func(int) check.StringSlice{
	"LenEQ": check.StringSliceLenEQ,
	"LenGT": check.StringSliceLenGT,
	"LenLT": check.StringSliceLenLT,
}

var strSlcCFIntInt = map[string]func(int, int) check.StringSlice{
	"LenBetween": check.StringSliceLenBetween,
}

var strSlcCFStrCF = map[string]func(check.String) check.StringSlice{
	"String": check.StringSliceStringCheck,
}

var strSlcCFStrSlcCFStr = map[string]func(check.StringSlice, string) check.StringSlice{
	"Not": check.StringSliceNot,
}

var strSlcCFStrSlcCFList = map[string]func(...check.StringSlice) check.StringSlice{
	"And": check.StringSliceAnd,
	"Or":  check.StringSliceOr,
}

// makeStrSlcCF returns a StringSlice checker corresponding to the
// given name - this is for checkers that are not parameterised
func makeStrSlcCF(_ *ast.CallExpr, fName string) (cf check.StringSlice, err error) {
	if f, ok := strSlcCFNoParam[fName]; ok {
		return f, nil
	}
	return nil, fmt.Errorf(
		"can't make the %s func: %s: the name is not recognised",
		strSlcCFName, fName)
}

// makeStrSlcCFInt returns a StringSlice checker corresponding to the
// given name - this is for checkers that take a single integer parameter
func makeStrSlcCFInt(e *ast.CallExpr, fName string) (cf check.StringSlice, err error) {
	var i int64
	errIntro := func() string {
		return fmt.Sprintf("can't make the %s func: %s(%d):",
			strSlcCFName, fName, i)
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

	if f, ok := strSlcCFInt[fName]; ok {
		return f(int(i)), nil
	}

	return nil, fmt.Errorf("%s the name is not recognised", errIntro())
}

// makeStrSlcCFIntInt returns a StringSlice checker corresponding to the
// given name - this is for checkers that take two integer parameters
func makeStrSlcCFIntInt(e *ast.CallExpr, fName string) (cf check.StringSlice, err error) {
	var i, j int64
	errIntro := func() string {
		return fmt.Sprintf("can't make the %s func: %s(%d, %d):",
			strSlcCFName, fName, i, j)
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

	i, err = getArgAsInt(e, 0)
	if err != nil {
		return nil, fmt.Errorf("%s %s", errIntro(), err)
	}
	j, err = getArgAsInt(e, 1)
	if err != nil {
		return nil, fmt.Errorf("%s %s", errIntro(), err)
	}

	if f, ok := strSlcCFIntInt[fName]; ok {
		return f(int(i), int(j)), nil
	}

	return nil, fmt.Errorf("%s the name is not recognised", errIntro())
}

// makeStrSlcCFStrCF returns a StringSlice checker corresponding to the
// given name - this is for checkers that take a string check parameter
func makeStrSlcCFStrCF(e *ast.CallExpr, fName string) (cf check.StringSlice, err error) {
	errIntro := "can't make the " + strSlcCFName +
		" func: " + fName + "(" + strCFName + "):"
	defer func() {
		if r := recover(); r != nil {
			cf = nil
			err = fmt.Errorf("%s %v", errIntro, r)
		}
	}()

	if err = checkArgCount(e, 1); err != nil {
		return nil, fmt.Errorf("%s %s", errIntro, err)
	}

	argExpr, err := getArg(e, 0)
	if err != nil {
		return nil, fmt.Errorf("%s can't get the %s argument: %s",
			errIntro, strCFName, err)
	}
	scf, err := getFuncStrCF(argExpr)
	if err != nil {
		return nil, fmt.Errorf("%s can't convert argument %d to %s: %s",
			errIntro, 0, strCFName, err)
	}

	if f, ok := strSlcCFStrCF[fName]; ok {
		return f(scf), nil
	}

	return nil, fmt.Errorf("%s the name is not recognised", errIntro)
}

// makeStrSlcCFStrSlcCFStr returns a StringSlice checker corresponding to the
// given name - this is for checkers that take a string slice check func and
// a string
func makeStrSlcCFStrSlcCFStr(e *ast.CallExpr, fName string) (cf check.StringSlice, err error) {
	var s string
	errIntro := func() string {
		return fmt.Sprintf("can't make the %s func: %s(%s, %s):",
			strSlcCFName, fName, strSlcCFName, s)
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
			errIntro(), strSlcCFName, err)
	}
	sscf, err := getFuncStrSlcCF(argExpr)
	if err != nil {
		return nil, fmt.Errorf("%s can't convert argument %d to %s: %s",
			errIntro(), 0, strSlcCFName, err)
	}
	s, err = getArgAsString(e, 1)
	if err != nil {
		return nil, fmt.Errorf("%s %s", errIntro(), err)
	}

	if f, ok := strSlcCFStrSlcCFStr[fName]; ok {
		return f(sscf, s), nil
	}

	return nil, fmt.Errorf("%s the name is not recognised", errIntro())
}

// makeStrSlcCFStrSlcCFList returns a StringSlice checker corresponding to
// the given name - this is for checkers that take a list of string slice
// check funcs
func makeStrSlcCFStrSlcCFList(e *ast.CallExpr, fName string) (cf check.StringSlice, err error) {
	errIntro := "can't make the " + strSlcCFName +
		" func: " + fName + "(" + strSlcCFName + " ...):"
	defer func() {
		if r := recover(); r != nil {
			cf = nil
			err = fmt.Errorf("%s %v", errIntro, r)
		}
	}()

	fArgs := make([]check.StringSlice, 0, len(e.Args))
	for i, argExpr := range e.Args {
		sscf, err := getFuncStrSlcCF(argExpr)
		if err != nil {
			return nil, fmt.Errorf("%s can't convert argument %d to %s: %s",
				errIntro, i, strSlcCFName, err)
		}
		fArgs = append(fArgs, sscf)
	}

	if f, ok := strSlcCFStrSlcCFList[fName]; ok {
		return f(fArgs...), nil
	}

	return nil, fmt.Errorf("%s the name is not recognised", errIntro)
}

// stringSliceCFParse returns a slice of string slice check functions and a nil
// error if the string is successfully parsed or nil and an error if the
// string couldn't be converted to a slice of check functions.
func stringSliceCFParse(s string) ([]check.StringSlice, error) {
	expr, err := parser.ParseExpr("[]T{\n" + s + "}")
	if err != nil {
		return nil, err
	}

	v := make([]check.StringSlice, 0, 1)
	cl, ok := expr.(*ast.CompositeLit)
	if !ok {
		return nil, fmt.Errorf("unexpected type for the collection of %s: %T",
			strSlcCFDesc, expr)
	}
	_, ok = cl.Type.(*ast.ArrayType)
	if !ok {
		return nil, fmt.Errorf("unexpected type for the array of %s: %T",
			strSlcCFDesc, cl.Type)
	}

	for _, elt := range cl.Elts {
		f, err := getFuncStrSlcCF(elt)
		if err != nil {
			return nil, fmt.Errorf("bad function: %s", err)
		}
		v = append(v, f)
	}

	return v, nil
}

// getFuncStrSlcCF will process the expression and return a string slice
// checker or nil
func getFuncStrSlcCF(elt ast.Expr) (check.StringSlice, error) {
	switch e := elt.(type) {
	case *ast.CallExpr:
		return callStrSlcCFMaker(e)
	case *ast.Ident:
		return makeStrSlcCF((*ast.CallExpr)(nil), e.Name)
	}
	return nil, fmt.Errorf("unexpected type: %T", elt)
}

// callStrSlcCFMaker calls the appropriate makeStrSlcCF... and returns the
// results
func callStrSlcCFMaker(e *ast.CallExpr) (cf check.StringSlice, err error) {
	defer func() {
		if r := recover(); r != nil {
			cf = nil
			err = fmt.Errorf("can't make the %s func: %v", strSlcCFName, r)
		}
	}()

	fd, err := getFuncDetails(e, strSlcCFName)
	if err != nil {
		return nil, err
	}

	switch fd.expectedArgs {
	case "":
		return makeStrSlcCF(e, fd.name)
	case "int":
		return makeStrSlcCFInt(e, fd.name)
	case "int, int":
		return makeStrSlcCFIntInt(e, fd.name)
	case strCFName:
		return makeStrSlcCFStrCF(e, fd.name)
	case strSlcCFName + ", string":
		return makeStrSlcCFStrSlcCFStr(e, fd.name)
	case strSlcCFName + " ...":
		return makeStrSlcCFStrSlcCFList(e, fd.name)
	default:
		return nil, fmt.Errorf("%s has an unexpected argument list: %s",
			fd.name, fd.expectedArgs)
	}
}
