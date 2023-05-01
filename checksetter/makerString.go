package checksetter

import (
	"fmt"
	"go/ast"
	"regexp"

	"github.com/nickwells/check.mod/v2/check"
)

const StringCheckerName = "string-checker"

var strMaker = MakerInfo[string]{
	Args: []string{},

	MF: func(e *ast.CallExpr, fName string) (cf check.ValCk[string], err error) {
		funcs := map[string]check.ValCk[string]{
			"OK": check.ValOK[string],
		}

		maker, ok := funcs[fName]
		if !ok {
			return nil, fmt.Errorf("Unknown function: %q", fName)
		}

		defer func() {
			if err != nil {
				err = fmt.Errorf("%s(...): %w", fName, err)
			}
		}()
		defer func() {
			if r := recover(); r != nil {
				cf = nil
				err = fmt.Errorf("%v", r)
			}
		}()

		if e != nil {
			if err = checkArgCount(e, 0); err != nil {
				return nil, err
			}
		}

		return maker, nil
	},
}

var strMakerStr = MakerInfo[string]{
	Args: []string{"string"},

	MF: func(e *ast.CallExpr, fName string) (cf check.ValCk[string], err error) {
		funcs := map[string]func(string) check.ValCk[string]{
			"EQ":        check.ValEQ[string],
			"GT":        check.ValGT[string],
			"GE":        check.ValGE[string],
			"LT":        check.ValLT[string],
			"LE":        check.ValLE[string],
			"HasPrefix": check.StringHasPrefix[string],
			"HasSuffix": check.StringHasSuffix[string],
		}

		maker, ok := funcs[fName]
		if !ok {
			return nil, fmt.Errorf("Unknown function: %q", fName)
		}

		defer func() {
			if err != nil {
				err = fmt.Errorf("%s(...): %w", fName, err)
			}
		}()
		defer func() {
			if r := recover(); r != nil {
				cf = nil
				err = fmt.Errorf("%v", r)
			}
		}()

		if err = checkArgCount(e, 1); err != nil {
			return nil, err
		}

		str, err := getString(e.Args[0])
		if err != nil {
			return nil, err
		}

		return maker(str), nil
	},
}

var strMakerIchecker = MakerInfo[string]{
	Args: []string{IntCheckerName},

	MF: func(e *ast.CallExpr, fName string) (cf check.ValCk[string], err error) {
		funcs := map[string]func(check.ValCk[int]) check.ValCk[string]{
			"Length": check.StringLength[string],
		}

		maker, ok := funcs[fName]
		if !ok {
			return nil, fmt.Errorf("Unknown function: %q", fName)
		}

		defer func() {
			if err != nil {
				err = fmt.Errorf("%s(...): %w", fName, err)
			}
		}()
		defer func() {
			if r := recover(); r != nil {
				cf = nil
				err = fmt.Errorf("%v", r)
			}
		}()

		if err = checkArgCount(e, 1); err != nil {
			return nil, err
		}

		ckFunc, err := getCheckFunc[int](e, 0, IntCheckerName)
		if err != nil {
			return nil, err
		}

		return maker(ckFunc), nil
	},
}

var strMakerRegexpStr = MakerInfo[string]{
	Args: []string{"regexp", "string"},

	MF: func(e *ast.CallExpr, fName string) (cf check.ValCk[string], err error) {
		funcs := map[string]func(*regexp.Regexp, string) check.ValCk[string]{
			"MatchesPattern": check.StringMatchesPattern[string],
		}

		maker, ok := funcs[fName]
		if !ok {
			return nil, fmt.Errorf("Unknown function: %q", fName)
		}

		defer func() {
			if err != nil {
				err = fmt.Errorf("%s(...): %w", fName, err)
			}
		}()
		defer func() {
			if r := recover(); r != nil {
				cf = nil
				err = fmt.Errorf("%v", r)
			}
		}()

		if err = checkArgCount(e, 2); err != nil {
			return nil, err
		}

		reStr, err := getString(e.Args[0])
		if err != nil {
			return nil, err
		}

		re, err := regexp.Compile(reStr)
		if err != nil {
			return nil, fmt.Errorf("the regexp doesn't compile: %s", err)
		}

		reDesc, err := getString(e.Args[1])
		if err != nil {
			return nil, err
		}

		return maker(re, reDesc), nil
	},
}

var strMakerStrcheckerString = MakerInfo[string]{
	Args: []string{StringCheckerName, "string"},

	MF: func(e *ast.CallExpr, fName string) (cf check.ValCk[string], err error) {
		funcs := map[string]func(check.ValCk[string], string) check.ValCk[string]{
			"Not": check.Not[string],
		}

		maker, ok := funcs[fName]
		if !ok {
			return nil, fmt.Errorf("Unknown function: %q", fName)
		}

		defer func() {
			if err != nil {
				err = fmt.Errorf("%s(...): %w", fName, err)
			}
		}()
		defer func() {
			if r := recover(); r != nil {
				cf = nil
				err = fmt.Errorf("%v", r)
			}
		}()

		if err = checkArgCount(e, 2); err != nil {
			return nil, err
		}

		ckFunc, err := getCheckFunc[string](e, 0, StringCheckerName)
		if err != nil {
			return nil, err
		}

		s, err := getString(e.Args[1])
		if err != nil {
			return nil, err
		}

		return maker(ckFunc, s), nil
	},
}

var strMakerMultiStrchecker = MakerInfo[string]{
	Args: []string{"...", StringCheckerName},

	MF: func(e *ast.CallExpr, fName string) (cf check.ValCk[string], err error) {
		funcs := map[string]func(...check.ValCk[string]) check.ValCk[string]{
			"And": check.And[string],
			"Or":  check.Or[string],
		}

		maker, ok := funcs[fName]
		if !ok {
			return nil, fmt.Errorf("Unknown function: %q", fName)
		}

		defer func() {
			if err != nil {
				err = fmt.Errorf("%s(...): %w", fName, err)
			}
		}()
		defer func() {
			if r := recover(); r != nil {
				cf = nil
				err = fmt.Errorf("%v", r)
			}
		}()

		checkFuncs, err := getCheckFuncs[string](e, StringCheckerName)
		if err != nil {
			return nil, err
		}

		return maker(checkFuncs...), nil
	},
}
