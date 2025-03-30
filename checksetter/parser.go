package checksetter

import (
	"fmt"
	"go/ast"
	"maps"
	"slices"

	"github.com/nickwells/check.mod/v2/check"
)

// MakerFunc is the type of a function that converts a CallExpr into a check
// func. The string should contain the name of the function to be generated.
type MakerFunc[T any] func(*ast.CallExpr, string) (check.ValCk[T], error)

// MakerInfo holds details of the function used to generate a check func. It
// holds the function value and a list of the arguments that it should
// take. Note that the Args member is present purely for documentation
// purposes; all the functions generated should be created by makers taking
// the listed args.
type MakerInfo[T any] struct {
	Args []string
	MF   MakerFunc[T]
}

// Parser records the available maker functions for generating the named
// checker functions. Note that a given parser will generate checker
// functions all of the same type.
type Parser[T any] struct {
	checkerName string
	makers      map[string]MakerInfo[T]
}

// MakeParser creates a new parser and adds it to the Parser register. It
// will return an error if there is already an entry with the same name.
func MakeParser[T any](checkerName string, makers map[string]MakerInfo[T]) (
	Parser[T], error,
) {
	if _, exists := parserRegister[checkerName]; exists {
		return Parser[T]{},
			fmt.Errorf("A Parser for %q already exists", checkerName)
	}

	p := Parser[T]{
		checkerName: checkerName,
		makers:      makers,
	}

	parserRegister[checkerName] = &p

	return p, nil
}

// Makers returns a slice holding the names of the checker-makers. The names
// are in alphabetical order.
//
// This can be used to construct the Allowed Values message for a setter.
func (p Parser[T]) Makers() []string {
	return slices.Sorted(maps.Keys(p.makers))
}

// Args returns the args that the named maker function takes.
//
// This can be used to construct the Allowed Values message for a setter.
func (p Parser[T]) Args(makerName string) ([]string, error) {
	mi, ok := p.makers[makerName]
	if !ok {
		return []string{}, fmt.Errorf("Unknown maker: %q", makerName)
	}

	return mi.Args, nil
}

// MakerFuncs returns a map of all the functions that the parser recognises
// and the arguments they each expect.
//
// This can be used to construct the Allowed Values message for a setter.
func (p Parser[T]) MakerFuncs() map[string][]string {
	makerFuncs := make(map[string][]string)

	for k, mi := range p.makers {
		makerFuncs[k] = mi.Args
	}

	return makerFuncs
}

// CheckerName returns the name of the checkers that will be generated
func (p Parser[T]) CheckerName() string {
	return p.checkerName
}

// Parse will parse the given string and return a slice of check.ValCk
// functions of the appropriate type and an error. The error will be nil if
// the parsing was successful, otherwise an error describing the problem and
// a nil slice will be returned.
func (p Parser[T]) Parse(s string) ([]check.ValCk[T], error) {
	exprs, err := getElts(s, p.checkerName)
	if err != nil {
		return nil, err
	}

	ckFuncs := make([]check.ValCk[T], 0, len(exprs))

	for _, e := range exprs {
		f, err := p.ParseExpr(e)
		if err != nil {
			return nil,
				fmt.Errorf("Can't make %s function: %s",
					p.checkerName, err)
		}

		ckFuncs = append(ckFuncs, f)
	}

	return ckFuncs, nil
}

// runMaker finds the appropriate function makerName and calls it passing the
// CallExpr and the function name.
func (p Parser[T]) runMaker(e *ast.CallExpr, makerName string) (
	check.ValCk[T], error,
) {
	maker, ok := p.makers[makerName]
	if !ok {
		return nil, fmt.Errorf("%s is an unknown function", makerName)
	}

	return maker.MF(e, makerName)
}

// CallExprMaker finds the function name using the information given in the
// CallExpr and then calls the maker func to build the check.ValCk
// function. It returns whatever that returns with no further processing
func (p Parser[T]) CallExprMaker(e *ast.CallExpr) (check.ValCk[T], error) {
	maker, err := getFuncName(e)
	if err != nil {
		return nil, err
	}

	return p.runMaker(e, maker)
}

// IdentMaker finds the function name using the information given in the
// Ident and then calls the maker func to build the check.ValCk
// function. It returns whatever that returns with no further processing
func (p Parser[T]) IdentMaker(e *ast.Ident) (check.ValCk[T], error) {
	return p.runMaker(nil, e.Name)
}

// ParseExpr parses an individual element from the list of functions. You
// should only need to call this if you are writing your own checker parser.
func (p Parser[T]) ParseExpr(elt ast.Expr) (cf check.ValCk[T], err error) {
	defer func() {
		if r := recover(); r != nil {
			cf = nil
			err = fmt.Errorf("Can't create the %s func: %v", p.checkerName, r)
		}
	}()

	switch e := elt.(type) {
	case *ast.Ident:
		return p.IdentMaker(e)
	case *ast.CallExpr:
		return p.CallExprMaker(e)
	}

	return nil, fmt.Errorf("unexpected type: %T", elt)
}

// getFuncName returns the function name from the call expression.
func getFuncName(e *ast.CallExpr) (string, error) {
	fID, ok := e.Fun.(*ast.Ident)
	if !ok {
		return "", fmt.Errorf("Syntax error: unexpected call type: %T", e.Fun)
	}

	return fID.Name, nil
}
