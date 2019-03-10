package checksetter

import (
	"fmt"
	"go/ast"
	"go/parser"
	"regexp"

	"github.com/nickwells/check.mod/check"
)

const (
	strCFName = "check.String"
	strCFDesc = "string check func"
)

var strCFInt = map[string]func(int) check.String{
	"LenEQ": check.StringLenEQ,
	"LenGT": check.StringLenGT,
	"LenLT": check.StringLenLT,
}

var strCFIntInt = map[string]func(int, int) check.String{
	"LenBetween": check.StringLenBetween,
}

var strCFStr = map[string]func(string) check.String{
	"Equals":    check.StringEquals,
	"HasPrefix": check.StringHasPrefix,
	"HasSuffix": check.StringHasSuffix,
}

var strCFREStr = map[string]func(*regexp.Regexp, string) check.String{
	"MatchesPattern": check.StringMatchesPattern,
}

var strCFStrCFStr = map[string]func(check.String, string) check.String{
	"Not": check.StringNot,
}

var strCFStrCFList = map[string]func(...check.String) check.String{
	"And": check.StringAnd,
	"Or":  check.StringOr,
}

// makeStrCFInt returns a String checker corresponding to the
// given name - this is for checkers that take a single integer parameter
func makeStrCFInt(e *ast.CallExpr, fName string) (cf check.String, err error) {
	var i int64
	errIntro := func() string {
		return fmt.Sprintf("can't make the %s func: %s(%d):",
			strCFName, fName, i)
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

	if f, ok := strCFInt[fName]; ok {
		return f(int(i)), nil
	}

	return nil, fmt.Errorf("%s the name is not recognised", errIntro())
}

// makeStrCFIntInt returns a String checker corresponding to the
// given name - this is for checkers that take two integer parameters
func makeStrCFIntInt(e *ast.CallExpr, fName string) (cf check.String, err error) {
	var i, j int64
	errIntro := func() string {
		return fmt.Sprintf("can't make the %s func: %s(%d, %d):",
			strCFName, fName, i, j)
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

	if f, ok := strCFIntInt[fName]; ok {
		return f(int(i), int(j)), nil
	}

	return nil, fmt.Errorf("%s the name is not recognised", errIntro())
}

// makeStrCFStr returns a String checker corresponding to the
// given name - this is for checkers that take a single string parameter
func makeStrCFStr(e *ast.CallExpr, fName string) (cf check.String, err error) {
	var s string
	errIntro := func() string {
		return fmt.Sprintf("can't make the %s func: %s(%s):",
			strCFName, fName, s)
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

	s, err = getArgAsString(e, 0)
	if err != nil {
		return nil, fmt.Errorf("%s %s", errIntro(), err)
	}

	if f, ok := strCFStr[fName]; ok {
		return f(s), nil
	}

	return nil, fmt.Errorf("%s the name is not recognised", errIntro())
}

// makeStrCFREStr returns a String checker corresponding to the given name -
// this is for checkers that take a regular expression and a single string
// parameter
func makeStrCFREStr(e *ast.CallExpr, fName string) (cf check.String, err error) {
	var reStr, reDesc string
	errIntro := func() string {
		return fmt.Sprintf("can't make the %s func: %s(%s, %s):",
			strCFName, fName, reStr, reDesc)
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

	reStr, err = getArgAsString(e, 0)
	if err != nil {
		return nil, fmt.Errorf("%s can't get the regexp: %s", errIntro(), err)
	}

	re, err := regexp.Compile(reStr)
	if err != nil {
		return nil, fmt.Errorf("%s the regexp doesn't compile: %s",
			errIntro(), err)
	}

	reDesc, err = getArgAsString(e, 1)
	if err != nil {
		return nil, fmt.Errorf("%s can't get the regexp description: %s",
			errIntro(), err)
	}

	if f, ok := strCFREStr[fName]; ok {
		return f(re, reDesc), nil
	}

	return nil, fmt.Errorf("%s the name is not recognised", errIntro())
}

// makeStrCFStrCFStr returns a String checker corresponding to the given name
// - this is for checkers that take a string check func and a string
// parameter
func makeStrCFStrCFStr(e *ast.CallExpr, fName string) (cf check.String, err error) {
	var s string
	errIntro := func() string {
		return fmt.Sprintf("can't make the %s func: %s(%s, %s):",
			strCFName, fName, strCFName, s)
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
			errIntro(), strCFName, err)
	}
	scf, err := getFuncStrCF(argExpr)
	if err != nil {
		return nil, fmt.Errorf("%s can't convert argument %d to %s: %s",
			errIntro(), 0, strCFName, err)
	}
	s, err = getArgAsString(e, 1)
	if err != nil {
		return nil, fmt.Errorf("%s %s", errIntro(), err)
	}

	if f, ok := strCFStrCFStr[fName]; ok {
		return f(scf, s), nil
	}

	return nil, fmt.Errorf("%s the name is not recognised", errIntro())
}

// makeStrCFStrCFList returns a String checker corresponding to the given
// name - this is for checkers that take a list of string check funcs
func makeStrCFStrCFList(e *ast.CallExpr, fName string) (cf check.String, err error) {
	errIntro := "can't make the " + strCFName +
		" func: " + fName + "(" + strCFName + " ...):"
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

	if f, ok := strCFStrCFList[fName]; ok {
		return f(fArgs...), nil
	}

	return nil, fmt.Errorf("%s the name is not recognised", errIntro)
}

// strCFParse returns a slice of string check functions and a nil error if
// the string is successfully parsed or nil and an error if the string
// couldn't be converted to a slice of check functions.
func stringCFParse(s string) ([]check.String, error) {
	expr, err := parser.ParseExpr("[]T{\n" + s + "}")
	if err != nil {
		return nil, err
	}

	v := make([]check.String, 0, 1)
	cl, ok := expr.(*ast.CompositeLit)
	if !ok {
		return nil, fmt.Errorf("unexpected type for the collection of %s: %T",
			strCFDesc, expr)
	}
	_, ok = cl.Type.(*ast.ArrayType)
	if !ok {
		return nil, fmt.Errorf("unexpected type for the array of %s: %T",
			strCFDesc, cl.Type)
	}

	for _, elt := range cl.Elts {
		f, err := getFuncStrCF(elt)
		if err != nil {
			return nil, fmt.Errorf("bad function: %s", err)
		}
		v = append(v, f)
	}

	return v, nil
}

// getFuncStrCF will process the expression and return a string checker or
// nil
func getFuncStrCF(elt ast.Expr) (check.String, error) {
	switch e := elt.(type) {
	case *ast.CallExpr:
		return callStrCFMaker(e)
	}
	return nil, fmt.Errorf("unexpected type: %T", elt)
}

// callStrCFMaker calls the appropriate makeStrCF... and returns the
// results
func callStrCFMaker(e *ast.CallExpr) (cf check.String, err error) {
	fd, err := getFuncDetails(e, strCFName)
	if err != nil {
		return nil, err
	}

	switch fd.expectedArgs {
	case "int":
		return makeStrCFInt(e, fd.name)
	case "int, int":
		return makeStrCFIntInt(e, fd.name)
	case "string":
		return makeStrCFStr(e, fd.name)
	case "regexp, string":
		return makeStrCFREStr(e, fd.name)
	case strCFName + ", string":
		return makeStrCFStrCFStr(e, fd.name)
	case strCFName + " ...":
		return makeStrCFStrCFList(e, fd.name)
	default:
		return nil, fmt.Errorf("%s has an unexpected argument list: %s",
			fd.name, fd.expectedArgs)
	}
}
