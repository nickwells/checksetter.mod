package checksetter

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/nickwells/check.mod/v2/check"
)

const IntCheckerName = "int-checker"

var (
	iMakerArgs = []string{}
	iMaker     = MakerInfo[int]{
		Args: iMakerArgs,

		MF: func(e *ast.CallExpr, fName string) (
			cf check.ValCk[int], err error,
		) {
			funcs := map[string]check.ValCk[int]{
				"OK": check.ValOK[int],
			}

			maker, ok := funcs[fName]
			if !ok {
				return nil, fmt.Errorf("Unknown function: %q", fName)
			}

			defer func() {
				if err != nil {
					err = fmt.Errorf("%s(%s): %w",
						fName, strings.Join(iMakerArgs, ", "), err)
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
	iMakerIArgs = []string{"int"}
	iMakerI     = MakerInfo[int]{
		Args: iMakerIArgs,

		MF: func(e *ast.CallExpr, fName string) (
			cf check.ValCk[int], err error,
		) {
			funcs := map[string]func(int) check.ValCk[int]{
				"EQ":          check.ValEQ[int],
				"GT":          check.ValGT[int],
				"GE":          check.ValGE[int],
				"LT":          check.ValLT[int],
				"LE":          check.ValLE[int],
				"Divides":     check.ValDivides[int],
				"IsAMultiple": check.ValIsAMultiple[int],
			}

			maker, ok := funcs[fName]
			if !ok {
				return nil, fmt.Errorf("Unknown function: %q", fName)
			}

			defer func() {
				if err != nil {
					err = fmt.Errorf("%s(%s): %w",
						fName, strings.Join(iMakerIArgs, ", "), err)
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

			i, err := getInt(e.Args[0])
			if err != nil {
				return nil, err
			}

			return maker(i), nil
		},
	}
)

var (
	iMakerIIArgs = []string{"int", "int"}
	iMakerII     = MakerInfo[int]{
		Args: iMakerIIArgs,

		MF: func(e *ast.CallExpr, fName string) (
			cf check.ValCk[int], err error,
		) {
			funcs := map[string]func(int, int) check.ValCk[int]{
				"Between": check.ValBetween[int],
			}

			maker, ok := funcs[fName]
			if !ok {
				return nil, fmt.Errorf("Unknown function: %q", fName)
			}

			defer func() {
				if err != nil {
					err = fmt.Errorf("%s(%s): %w",
						fName, strings.Join(iMakerIIArgs, ", "), err)
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

			i1, err := getInt(e.Args[0])
			if err != nil {
				return nil, err
			}

			i2, err := getInt(e.Args[1])
			if err != nil {
				return nil, err
			}

			return maker(i1, i2), nil
		},
	}
)

var (
	iMakerIcheckerStringArgs = []string{IntCheckerName, "string"}
	iMakerIcheckerString     = MakerInfo[int]{
		Args: iMakerIcheckerStringArgs,

		MF: func(e *ast.CallExpr, fName string) (
			cf check.ValCk[int], err error,
		) {
			funcs := map[string]func(check.ValCk[int], string) check.ValCk[int]{
				"Not": check.Not[int],
			}

			maker, ok := funcs[fName]
			if !ok {
				return nil, fmt.Errorf("Unknown function: %q", fName)
			}

			defer func() {
				if err != nil {
					err = fmt.Errorf("%s(%s): %w",
						fName, strings.Join(iMakerIcheckerStringArgs, ", "), err)
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

			ckFunc, err := getCheckFunc[int](e, 0, IntCheckerName)
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
	iMakerMultiIcheckerArgs = []string{"...", IntCheckerName}
	iMakerMultiIchecker     = MakerInfo[int]{
		Args: iMakerMultiIcheckerArgs,

		MF: func(e *ast.CallExpr, fName string) (
			cf check.ValCk[int], err error,
		) {
			funcs := map[string]func(...check.ValCk[int]) check.ValCk[int]{
				"And": check.And[int],
				"Or":  check.Or[int],
			}

			maker, ok := funcs[fName]
			if !ok {
				return nil, fmt.Errorf("Unknown function: %q", fName)
			}

			defer func() {
				if err != nil {
					err = fmt.Errorf("%s(%s): %w",
						fName, strings.Join(iMakerMultiIcheckerArgs, ", "), err)
				}
			}()
			defer func() {
				if r := recover(); r != nil {
					cf = nil
					err = fmt.Errorf("%v", r)
				}
			}()

			checkFuncs, err := getCheckFuncs[int](e, IntCheckerName)
			if err != nil {
				return nil, err
			}

			return maker(checkFuncs...), nil
		},
	}
)
