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
func makeStrCFInt(name string, i int64) check.String {
	if f, ok := strCFInt[name]; ok {
		return f(int(i))
	}
	return nil
}

// makeStrCFIntInt returns a String checker corresponding to the
// given name - this is for checkers that take two integer parameters
func makeStrCFIntInt(name string, i, j int64) check.String {
	if f, ok := strCFIntInt[name]; ok {
		return f(int(i), int(j))
	}
	return nil
}

// makeStrCFStr returns a String checker corresponding to the
// given name - this is for checkers that take a single string parameter
func makeStrCFStr(name, s string) check.String {
	if f, ok := strCFStr[name]; ok {
		return f(s)
	}
	return nil
}

// makeStrCFREStr returns a String checker corresponding to the given name -
// this is for checkers that take a single string parameter
func makeStrCFREStr(name string, re *regexp.Regexp, s string) check.String {
	if f, ok := strCFREStr[name]; ok {
		return f(re, s)
	}
	return nil
}

// makeStrCFStrCFStr returns a String checker corresponding to the given name
// - this is for checkers that take a string check func and a string
// parameter
func makeStrCFStrCFStr(name string, cf check.String, s string) check.String {
	if f, ok := strCFStrCFStr[name]; ok {
		return f(cf, s)
	}
	return nil
}

// makeStrCFStrCFList returns a String checker corresponding to the given
// name - this is for checkers that take a list of string check funcs
func makeStrCFStrCFList(name string, cf ...check.String) check.String {
	if f, ok := strCFStrCFList[name]; ok {
		return f(cf...)
	}
	return nil
}

// strCFParse returns a slice of string slice check functions and a nil error
// if the string is successfully parsed or nil and an error if the string
// couldn't be converted to a slice of check functions.
func strCFParse(s string) ([]check.String, error) {
	expr, err := parser.ParseExpr("[]T{\n" + s + "}")
	if err != nil {
		return nil, err
	}

	v := make([]check.String, 0, 1)
	cl, ok := expr.(*ast.CompositeLit)
	if !ok {
		return nil,
			fmt.Errorf("unexpected type for the collection of %s: %T",
				strCFDesc, expr)
	}
	_, ok = cl.Type.(*ast.ArrayType)
	if !ok {
		return nil,
			fmt.Errorf("unexpected type for the array of %s: %T",
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
	defer func() {
		if r := recover(); r != nil {
			cf = nil
			err = fmt.Errorf("Cannot create the %s func: %v", strCFName, r)
		}
	}()

	fd, err := getFuncDetails(e, strCFName)
	if err != nil {
		return nil, err
	}

	var f func(string) error

	switch fd.expectedArgs {
	case "int":
		i, err := getArgAsInt(e, fd, 0)
		if err != nil {
			return nil, err
		}
		f = makeStrCFInt(fd.name, i)
		if f == nil {
			return nil, fmt.Errorf("cannot create the %s: %s(%d)",
				strCFDesc, fd.name, i)
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
		f = makeStrCFIntInt(fd.name, i, j)
		if f == nil {
			return nil, fmt.Errorf("cannot create the %s: %s(%d, %d)",
				strCFDesc, fd.name, i, j)
		}
	case "string":
		s, err := getArgAsString(e, fd, 0)
		if err != nil {
			return nil, err
		}
		f = makeStrCFStr(fd.name, s)
		if f == nil {
			return nil, fmt.Errorf("cannot create the %s: %s(...)",
				strCFDesc, fd.name)
		}
	case "regexp, string":
		reStr, err := getArgAsString(e, fd, 0)
		if err != nil {
			return nil, err
		}
		re, err := regexp.Compile(reStr)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot create the regexp parameter for the %s (%s): %s",
				strCFDesc, fd.name, err)
		}
		reDesc, err := getArgAsString(e, fd, 1)
		if err != nil {
			return nil, err
		}
		f = makeStrCFREStr(fd.name, re, reDesc)
		if f == nil {
			return nil, fmt.Errorf("cannot create the %s: %s(...)",
				strCFDesc, fd.name)
		}
	case strCFName + ", string":
		argExpr, err := getArg(e, fd, 0)
		if err != nil {
			return nil, fmt.Errorf(
				"couldn't create the %s argument for the %s: %s(...): %s ",
				strCFDesc, strCFDesc, fd.name, err)
		}
		cssf, err := getFuncStrCF(argExpr)
		if err != nil {
			return nil, fmt.Errorf(
				"couldn't create the %s argument for the %s: %s(...): %s ",
				strCFDesc, strCFDesc, fd.name, err)
		}
		s, err := getArgAsString(e, fd, 1)
		if err != nil {
			return nil, err
		}
		f = makeStrCFStrCFStr(fd.name, cssf, s)
		if f == nil {
			return nil, fmt.Errorf("cannot create the %s: %s(...)",
				strCFDesc, fd.name)
		}
	case strCFName + " ...":
		scfArgs := make([]check.String, 0, len(e.Args))
		for i, argExpr := range e.Args {
			scf, err := getFuncStrCF(argExpr)
			if err != nil {
				return nil, fmt.Errorf(
					"couldn't create the %s argument (%d) for the %s: %s(...): %s ",
					strCFDesc, i, strCFDesc, fd.name, err)
			}
			scfArgs = append(scfArgs, scf)
		}
		f = makeStrCFStrCFList(fd.name, scfArgs...)
		if f == nil {
			return nil, fmt.Errorf("cannot create the %s: %s(...)",
				strCFDesc, fd.name)
		}
	default:
		return nil, fmt.Errorf("unexpected argument list: %s", fd.expectedArgs)
	}
	return f, nil
}
