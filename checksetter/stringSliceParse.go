package checksetter

import (
	"fmt"
	"go/ast"
	"go/parser"

	"github.com/nickwells/check.mod/v2/check"
)

const (
	strSlcCFName = "check.StringSlice"
	strSlcCFDesc = "string slice check func"
)

type strSlcCFMaker func(*ast.CallExpr, string) (check.StringSlice, error)

var strSlcCFArgsToFunc map[string]strSlcCFMaker

func init() {
	strSlcCFArgsToFunc = map[string]strSlcCFMaker{
		"":                        makeStrSlcCF,
		intCFName:                 makeStrSlcCFIntCF,
		strCFName:                 makeStrSlcCFStrCF,
		strCFName + " ...":        makeStrSlcCFStrCFList,
		strCFName + ", string":    makeStrSlcCFStrCFStr,
		strSlcCFName + ", string": makeStrSlcCFStrSlcCFStr,
		strSlcCFName + " ...":     makeStrSlcCFStrSlcCFList,
	}
}

var strSlcCFNoParam = map[string]check.StringSlice{
	"NoDups": check.SliceHasNoDups[[]string, string],
}

var strSlcCFIntCF = map[string]func(check.ValCk[int]) check.StringSlice{
	"Length": check.SliceLength[[]string, string],
}

var strSlcCFStrCFList = map[string]func(...check.String) check.StringSlice{
	"SliceByPos": check.SliceByPos[[]string, string],
}

var strSlcCFStrCFStr = map[string]func(check.String, string) check.StringSlice{
	"SliceAny": check.SliceAny[[]string, string],
}

var strSlcCFStrCF = map[string]func(check.String) check.StringSlice{
	"SliceAll": check.SliceAll[[]string, string],
}

var strSlcCFStrSlcCFStr = map[string]func(check.StringSlice, string) check.StringSlice{
	"Not": check.Not[[]string],
}

var strSlcCFStrSlcCFList = map[string]func(...check.StringSlice) check.StringSlice{
	"And": check.And[[]string],
	"Or":  check.Or[[]string],
}

// makeStrSlcCF returns a StringSlice checker corresponding to the
// given name - this is for checkers that are not parameterised
func makeStrSlcCF(e *ast.CallExpr, fName string) (cf check.StringSlice, err error) {
	errIntro := fmt.Sprintf("can't make the %s func: %s():",
		strSlcCFName, fName)

	if e != nil {
		if err = checkArgCount(e, 0); err != nil {
			return nil, fmt.Errorf("%s %s", errIntro, err)
		}
	}

	if f, ok := strSlcCFNoParam[fName]; ok {
		return f, nil
	}
	return nil, fmt.Errorf("%s the name is not recognised", errIntro)
}

// makeStrSlcCFIntCF returns a StringSlice checker corresponding to the
// given name - this is for checkers that take a single integer-checker parameter
func makeStrSlcCFIntCF(e *ast.CallExpr, fName string) (cf check.StringSlice, err error) {
	errIntro := func() string {
		return fmt.Sprintf("can't make the %s func: %s(...):",
			strSlcCFName, fName)
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

	var icf check.ValCk[int]
	icf, err = getFuncIntCF(e.Args[0])
	if err != nil {
		return nil, fmt.Errorf("%s %s", errIntro(), err)
	}

	if f, ok := strSlcCFIntCF[fName]; ok {
		return f(icf), nil
	}

	return nil, fmt.Errorf("%s the name is not recognised", errIntro())
}

// makeStrSlcCFStrCF returns a StringSlice checker corresponding to the
// given name - this is for checkers that take a single String-checker parameter
func makeStrSlcCFStrCF(e *ast.CallExpr, fName string) (cf check.StringSlice, err error) {
	errIntro := func() string {
		return fmt.Sprintf("can't make the %s func: %s(...):",
			strSlcCFName, fName)
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

	var scf check.String
	scf, err = getFuncStrCF(e.Args[0])
	if err != nil {
		return nil, fmt.Errorf("%s %s", errIntro(), err)
	}

	if f, ok := strSlcCFStrCF[fName]; ok {
		return f(scf), nil
	}

	return nil, fmt.Errorf("%s the name is not recognised", errIntro())
}

// makeStrSlcCFStrCFStr returns a StringSlice checker corresponding to the
// given name - this is for checkers that take a string check parameter and a
// string
func makeStrSlcCFStrCFStr(e *ast.CallExpr, fName string) (cf check.StringSlice, err error) {
	var s string
	errIntro := func() string {
		return fmt.Sprintf("can't make the %s func: %s(%s, %s):",
			strSlcCFName, fName, strCFName, s)
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

	scf, err := getFuncStrCF(e.Args[0])
	if err != nil {
		return nil, fmt.Errorf("%s can't convert argument %d to %s: %s",
			errIntro(), 0, strCFName, err)
	}
	s, err = getString(e.Args[1])
	if err != nil {
		return nil, fmt.Errorf("%s %s", errIntro(), err)
	}

	if f, ok := strSlcCFStrCFStr[fName]; ok {
		return f(scf, s), nil
	}

	return nil, fmt.Errorf("%s the name is not recognised", errIntro())
}

// makeStrSlcCFStrCFList returns a StringSlice checker corresponding to
// the given name - this is for checkers that take a list of string
// check funcs
func makeStrSlcCFStrCFList(e *ast.CallExpr, fName string) (cf check.StringSlice, err error) {
	errIntro := "can't make the " + strSlcCFName +
		" func: " + fName + "(" + strSlcCFName + " ...):"
	defer func() {
		if r := recover(); r != nil {
			cf = nil
			err = fmt.Errorf("%s %v", errIntro, r)
		}
	}()

	fArgs := make([]check.String, 0, len(e.Args))
	for i, argExpr := range e.Args {
		scf, err := getFuncStrCF(argExpr)
		if err != nil {
			return nil, fmt.Errorf("%s can't convert argument %d to %s: %s",
				errIntro, i, strCFName, err)
		}
		fArgs = append(fArgs, scf)
	}

	if f, ok := strSlcCFStrCFList[fName]; ok {
		return f(fArgs...), nil
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

	sscf, err := getFuncStrSlcCF(e.Args[0])
	if err != nil {
		return nil, fmt.Errorf(
			"%s can't convert the first argument to %s: %s",
			errIntro(), strSlcCFName, err)
	}
	s, err = getString(e.Args[1])
	if err != nil {
		return nil, fmt.Errorf(
			"%s can't convert the second argument to a string: %s",
			errIntro(), err)
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
func getFuncStrSlcCF(elt ast.Expr) (cf check.StringSlice, err error) {
	defer func() {
		if r := recover(); r != nil {
			cf = nil
			err = fmt.Errorf("can't make the %s func: %v", strSlcCFName, r)
		}
	}()

	switch e := elt.(type) {
	case *ast.Ident:
		return makeStrSlcCF(nil, e.Name)
	case *ast.CallExpr:
		fd, err := getFuncDetails(e, strSlcCFName)
		if err != nil {
			return nil, err
		}

		maker, ok := strSlcCFArgsToFunc[fd.expectedArgs]
		if !ok {
			return nil, fmt.Errorf("%s has an unrecognised argument list: %s",
				fd.name, fd.expectedArgs)
		}

		return maker(e, fd.name)
	}

	return nil, fmt.Errorf("unexpected type: %T", elt)
}
