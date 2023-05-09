package checksetter

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/nickwells/check.mod/v2/check"
)

const StringSliceCheckerName = "string-slice-checker"

var (
	strSlcMakerArgs = []string{}
	strSlcMaker     = MakerInfo[[]string]{
		Args: strSlcMakerArgs,

		MF: func(e *ast.CallExpr, fName string) (
			cf check.ValCk[[]string], err error,
		) {
			funcs := map[string]check.ValCk[[]string]{
				"OK":     check.ValOK[[]string],
				"NoDups": check.SliceHasNoDups[[]string],
			}

			maker, ok := funcs[fName]
			if !ok {
				return nil, fmt.Errorf("Unknown function: %q", fName)
			}

			defer func() {
				if err != nil {
					err = fmt.Errorf("%s(%s): %w",
						fName, strings.Join(strSlcMakerArgs, ", "), err)
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
	strSlcMakerIcheckerArgs = []string{IntCheckerName}
	strSlcMakerIchecker     = MakerInfo[[]string]{
		Args: strSlcMakerIcheckerArgs,

		MF: func(e *ast.CallExpr, fName string) (
			cf check.ValCk[[]string], err error,
		) {
			funcs := map[string]func(check.ValCk[int]) check.ValCk[[]string]{
				"Length": check.SliceLength[[]string],
			}

			maker, ok := funcs[fName]
			if !ok {
				return nil, fmt.Errorf("Unknown function: %q", fName)
			}

			defer func() {
				if err != nil {
					err = fmt.Errorf("%s(%s): %w",
						fName,
						strings.Join(strSlcMakerIcheckerArgs, ", "),
						err)
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

var strSlcMakerStrSlccheckerStringArgs = []string{
	StringSliceCheckerName,
	"string",
}

var strSlcMakerStrSlccheckerString = MakerInfo[[]string]{
	Args: strSlcMakerStrSlccheckerStringArgs,

	MF: func(e *ast.CallExpr, fName string) (
		cf check.ValCk[[]string], err error,
	) {
		funcs := map[string]func(
			check.ValCk[[]string], string,
		) check.ValCk[[]string]{
			"Not": check.Not[[]string],
		}

		maker, ok := funcs[fName]
		if !ok {
			return nil, fmt.Errorf("Unknown function: %q", fName)
		}

		defer func() {
			if err != nil {
				err = fmt.Errorf("%s(%s): %w",
					fName,
					strings.Join(strSlcMakerStrSlccheckerStringArgs, ", "),
					err)
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

		ckFunc, err := getCheckFunc[[]string](e, 0, StringSliceCheckerName)
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

var (
	strSlcMakerStrcheckerStringArgs = []string{StringCheckerName, "string"}
	strSlcMakerStrcheckerString     = MakerInfo[[]string]{
		Args: strSlcMakerStrcheckerStringArgs,

		MF: func(e *ast.CallExpr, fName string) (
			cf check.ValCk[[]string], err error,
		) {
			funcs := map[string]func(
				check.ValCk[string], string,
			) check.ValCk[[]string]{
				"SliceAny": check.SliceAny[[]string],
			}

			maker, ok := funcs[fName]
			if !ok {
				return nil, fmt.Errorf("Unknown function: %q", fName)
			}

			defer func() {
				if err != nil {
					err = fmt.Errorf("%s(%s): %w",
						fName,
						strings.Join(strSlcMakerStrcheckerStringArgs,
							", "), err)
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
)

var (
	strSlcMakerStrcheckerArgs = []string{StringCheckerName}
	strSlcMakerStrchecker     = MakerInfo[[]string]{
		Args: strSlcMakerStrcheckerArgs,

		MF: func(e *ast.CallExpr, fName string) (
			cf check.ValCk[[]string], err error,
		) {
			funcs := map[string]func(check.ValCk[string]) check.ValCk[[]string]{
				"SliceAll": check.SliceAll[[]string],
			}

			maker, ok := funcs[fName]
			if !ok {
				return nil, fmt.Errorf("Unknown function: %q", fName)
			}

			defer func() {
				if err != nil {
					err = fmt.Errorf("%s(%s): %w",
						fName, strings.Join(strSlcMakerStrcheckerArgs, ", "), err)
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

			ckFunc, err := getCheckFunc[string](e, 0, StringCheckerName)
			if err != nil {
				return nil, err
			}

			return maker(ckFunc), nil
		},
	}
)

var (
	strSlcMakerMultiStrcheckerArgs = []string{"...", StringCheckerName}
	strSlcMakerMultiStrchecker     = MakerInfo[[]string]{
		Args: strSlcMakerMultiStrcheckerArgs,

		MF: func(e *ast.CallExpr, fName string) (
			cf check.ValCk[[]string], err error,
		) {
			funcs := map[string]func(...check.ValCk[string]) check.ValCk[[]string]{
				"SliceByPos": check.SliceByPos[[]string],
			}

			maker, ok := funcs[fName]
			if !ok {
				return nil, fmt.Errorf("Unknown function: %q", fName)
			}

			defer func() {
				if err != nil {
					err = fmt.Errorf("%s(%s): %w",
						fName,
						strings.Join(strSlcMakerMultiStrcheckerArgs, ", "),
						err)
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

var (
	strSlcMakerMultiStrSlccheckerArgs = []string{"...", StringSliceCheckerName}
	strSlcMakerMultiStrSlcchecker     = MakerInfo[[]string]{
		Args: strSlcMakerMultiStrSlccheckerArgs,

		MF: func(e *ast.CallExpr, fName string) (
			cf check.ValCk[[]string], err error,
		) {
			funcs := map[string]func(
				...check.ValCk[[]string]) check.ValCk[[]string]{
				"And": check.And[[]string],
				"Or":  check.Or[[]string],
			}

			maker, ok := funcs[fName]
			if !ok {
				return nil, fmt.Errorf("Unknown function: %q", fName)
			}

			defer func() {
				if err != nil {
					err = fmt.Errorf("%s(%s): %w",
						fName,
						strings.Join(strSlcMakerMultiStrSlccheckerArgs, ", "),
						err)
				}
			}()
			defer func() {
				if r := recover(); r != nil {
					cf = nil
					err = fmt.Errorf("%v", r)
				}
			}()

			checkFuncs, err := getCheckFuncs[[]string](e, StringSliceCheckerName)
			if err != nil {
				return nil, err
			}

			return maker(checkFuncs...), nil
		},
	}
)
