package checksetter

import (
	"fmt"
	"go/ast"
	"regexp"
	"strings"

	"github.com/nickwells/check.mod/v2/check"
)

const StringCheckerName = "string-checker"

var (
	strMakerArgs = []string{}
	strMaker     = MakerInfo[string]{
		Args: strMakerArgs,

		MF: func(e *ast.CallExpr, fName string) (
			cf check.ValCk[string], err error,
		) {
			funcs := map[string]check.ValCk[string]{
				"OK": check.ValOK[string],
			}

			maker, ok := funcs[fName]
			if !ok {
				return nil, fmt.Errorf("Unknown function: %q", fName)
			}

			defer func() {
				if err != nil {
					err = fmt.Errorf("%s(%s): %w",
						fName, strings.Join(strMakerArgs, ", "), err)
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
)

var (
	strMakerStrArgs = []string{"string"}
	strMakerStr     = MakerInfo[string]{
		Args: strMakerStrArgs,

		MF: func(e *ast.CallExpr, fName string) (
			cf check.ValCk[string], err error,
		) {
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
					err = fmt.Errorf("%s(%s): %w",
						fName, strings.Join(strMakerStrArgs, ", "), err)
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
)

var (
	strMakerIcheckerArgs = []string{IntCheckerName}
	strMakerIchecker     = MakerInfo[string]{
		Args: strMakerIcheckerArgs,

		MF: func(e *ast.CallExpr, fName string) (
			cf check.ValCk[string], err error,
		) {
			funcs := map[string]func(check.ValCk[int]) check.ValCk[string]{
				"Length": check.StringLength[string],
			}

			maker, ok := funcs[fName]
			if !ok {
				return nil, fmt.Errorf("Unknown function: %q", fName)
			}

			defer func() {
				if err != nil {
					err = fmt.Errorf("%s(%s): %w",
						fName, strings.Join(strMakerIcheckerArgs, ", "), err)
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
)

var (
	strMakerRegexpStrArgs = []string{"regexp", "string"}
	strMakerRegexpStr     = MakerInfo[string]{
		Args: strMakerRegexpStrArgs,

		MF: func(e *ast.CallExpr, fName string) (
			cf check.ValCk[string], err error,
		) {
			funcs := map[string]func(*regexp.Regexp, string) check.ValCk[string]{
				"MatchesPattern": check.StringMatchesPattern[string],
			}

			maker, ok := funcs[fName]
			if !ok {
				return nil, fmt.Errorf("Unknown function: %q", fName)
			}

			defer func() {
				if err != nil {
					err = fmt.Errorf("%s(%s): %w",
						fName, strings.Join(strMakerRegexpStrArgs, ", "), err)
				}
			}()
			defer func() {
				if r := recover(); r != nil {
					cf = nil
					err = fmt.Errorf("%v", r)
				}
			}()

			if err = checkArgCount(e, 2); err != nil { //nolint:mnd
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
)

var (
	strMakerStrcheckerStringArgs = []string{StringCheckerName, "string"}
	strMakerStrcheckerString     = MakerInfo[string]{
		Args: strMakerStrcheckerStringArgs,

		MF: func(e *ast.CallExpr, fName string) (
			cf check.ValCk[string], err error,
		) {
			funcs := map[string]func(check.ValCk[string], string) check.ValCk[string]{
				"Not": check.Not[string],
			}

			maker, ok := funcs[fName]
			if !ok {
				return nil, fmt.Errorf("Unknown function: %q", fName)
			}

			defer func() {
				if err != nil {
					err = fmt.Errorf("%s(%s): %w",
						fName, strings.Join(strMakerStrcheckerStringArgs, ", "), err)
				}
			}()
			defer func() {
				if r := recover(); r != nil {
					cf = nil
					err = fmt.Errorf("%v", r)
				}
			}()

			if err = checkArgCount(e, 2); err != nil { //nolint:mnd
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
)

var (
	strMakerMultiStrcheckerArgs = []string{"...", StringCheckerName}
	strMakerMultiStrchecker     = MakerInfo[string]{
		Args: strMakerMultiStrcheckerArgs,

		MF: func(e *ast.CallExpr, fName string) (
			cf check.ValCk[string], err error,
		) {
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
					err = fmt.Errorf("%s(%s): %w",
						fName, strings.Join(strMakerMultiStrcheckerArgs, ", "), err)
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
)
