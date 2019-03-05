package checksetter

import (
	"errors"
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
func makeStrSlcCF(name string) check.StringSlice {
	if f, ok := strSlcCFNoParam[name]; ok {
		return f
	}
	return nil
}

// makeStrSlcCFInt returns a StringSlice checker corresponding to the
// given name - this is for checkers that take a single integer parameter
func makeStrSlcCFInt(name string, i int64) check.StringSlice {
	if f, ok := strSlcCFInt[name]; ok {
		return f(int(i))
	}
	return nil
}

// makeStrSlcCFIntInt returns a StringSlice checker corresponding to the
// given name - this is for checkers that take two integer parameters
func makeStrSlcCFIntInt(name string, i, j int64) check.StringSlice {
	if f, ok := strSlcCFIntInt[name]; ok {
		return f(int(i), int(j))
	}
	return nil
}

// makeStrSlcCFStrCF returns a StringSlice checker corresponding to the
// given name - this is for checkers that take a string check parameter
func makeStrSlcCFStrCF(name string, csf check.String) check.StringSlice {
	if f, ok := strSlcCFStrCF[name]; ok {
		return f(csf)
	}
	return nil
}

// makeStrSlcCFStrSlcCFStr returns a StringSlice checker corresponding to the
// given name - this is for checkers that take a string slice check func and
// a string
func makeStrSlcCFStrSlcCFStr(name string, cf check.StringSlice, s string) check.StringSlice {
	if f, ok := strSlcCFStrSlcCFStr[name]; ok {
		return f(cf, s)
	}
	return nil
}

// makeStrSlcCFStrSlcCFList returns a StringSlice checker corresponding to
// the given name - this is for checkers that take a list of string slice
// check funcs
func makeStrSlcCFStrSlcCFList(name string, cf ...check.StringSlice) check.StringSlice {
	if f, ok := strSlcCFStrSlcCFList[name]; ok {
		return f(cf...)
	}
	return nil
}

// strSlcCFParse returns a slice of string slice check functions and a nil
// error if the string is successfully parsed or nil and an error if the
// string couldn't be converted to a slice of check functions.
func strSlcCFParse(s string) ([]check.StringSlice, error) {
	expr, err := parser.ParseExpr("[]T{\n" + s + "}")
	if err != nil {
		return nil, err
	}

	v := make([]check.StringSlice, 0, 1)
	cl, ok := expr.(*ast.CompositeLit)
	if !ok {
		return nil,
			fmt.Errorf("unexpected type for the collection of %s: %T",
				strSlcCFDesc, expr)
	}
	_, ok = cl.Type.(*ast.ArrayType)
	if !ok {
		return nil,
			fmt.Errorf("unexpected type for the array of %s: %T",
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
		f := makeStrSlcCF(e.Name)
		if f == nil {
			return nil,
				fmt.Errorf("unknown unparameterised %s: %s",
					strSlcCFDesc, e.Name)
		}
		return f, nil
	}
	return nil, fmt.Errorf("unexpected type: %T", elt)
}

// callStrSlcCFMaker calls the appropriate makeStrSlcCF... and returns the
// results
func callStrSlcCFMaker(e *ast.CallExpr) (cf check.StringSlice, err error) {
	defer func() {
		if r := recover(); r != nil {
			cf = nil
			err = fmt.Errorf("Cannot create the %s func: %v", strSlcCFName, r)
		}
	}()

	fd, err := getFuncDetails(e, strSlcCFName)
	if err != nil {
		return nil, err
	}

	var f func([]string) error

	switch fd.expectedArgs {
	case "":
		f = makeStrSlcCF(fd.name)
		if f == nil {
			return nil, errors.New(
				"cannot create the " + strSlcCFDesc + ": " + fd.name)
		}
	case "int":
		i, err := getArgAsInt(e, fd, 0)
		if err != nil {
			return nil, err
		}
		f = makeStrSlcCFInt(fd.name, i)
		if f == nil {
			return nil, fmt.Errorf("cannot create the %s: %s(%d)",
				strSlcCFDesc, fd.name, i)
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
		f = makeStrSlcCFIntInt(fd.name, i, j)
		if f == nil {
			return nil, fmt.Errorf("cannot create the %s: %s(%d, %d)",
				strSlcCFDesc, fd.name, i, j)
		}
	case strCFName:
		argExpr, err := getArg(e, fd, 0)
		if err != nil {
			return nil, fmt.Errorf(
				"couldn't create the %s argument for the %s: %s(...): %s ",
				strCFDesc, strSlcCFDesc, fd.name, err)
		}
		csf, err := getFuncStrCF(argExpr)
		if err != nil {
			return nil, fmt.Errorf(
				"couldn't create the %s argument for the %s: %s(...): %s ",
				strCFDesc, strSlcCFDesc, fd.name, err)
		}
		f = makeStrSlcCFStrCF(fd.name, csf)
		if f == nil {
			return nil, fmt.Errorf("cannot create the %s: %s(...)",
				strSlcCFDesc, fd.name)
		}
	case strSlcCFName + ", string":
		argExpr, err := getArg(e, fd, 0)
		if err != nil {
			return nil, fmt.Errorf(
				"couldn't create the %s argument for the %s: %s(...): %s ",
				strSlcCFDesc, strSlcCFDesc, fd.name, err)
		}
		cssf, err := getFuncStrSlcCF(argExpr)
		if err != nil {
			return nil, fmt.Errorf(
				"couldn't create the %s argument for the %s: %s(...): %s ",
				strSlcCFDesc, strSlcCFDesc, fd.name, err)
		}
		s, err := getArgAsString(e, fd, 1)
		if err != nil {
			return nil, err
		}
		f = makeStrSlcCFStrSlcCFStr(fd.name, cssf, s)
		if f == nil {
			return nil, fmt.Errorf("cannot create the %s: %s(...)",
				strSlcCFDesc, fd.name)
		}
	case strSlcCFName + " ...":
		sscfArgs := make([]check.StringSlice, 0, len(e.Args))
		for i, argExpr := range e.Args {
			sscf, err := getFuncStrSlcCF(argExpr)
			if err != nil {
				return nil, fmt.Errorf(
					"couldn't create the %s argument (%d) for the %s: %s(...): %s ",
					strSlcCFDesc, i, strSlcCFDesc, fd.name, err)
			}
			sscfArgs = append(sscfArgs, sscf)
		}
		f = makeStrSlcCFStrSlcCFList(fd.name, sscfArgs...)
		if f == nil {
			return nil, fmt.Errorf("cannot create the %s: %s(...)",
				strSlcCFDesc, fd.name)
		}
	default:
		return nil, fmt.Errorf("unexpected argument list: %s", fd.expectedArgs)
	}
	return f, nil
}
