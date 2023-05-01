package checksetter

import (
	"fmt"
	"go/ast"

	"github.com/nickwells/check.mod/v2/check"
)

const StringSliceCheckerName = "string-slice-checker"

var strSlcMaker = MakerInfo[[]string]{
	Args: []string{},

	MF: func(e *ast.CallExpr, fName string) (cf check.ValCk[[]string], err error) {
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

var strSlcMakerIchecker = MakerInfo[[]string]{
	Args: []string{IntCheckerName},

	MF: func(e *ast.CallExpr, fName string) (cf check.ValCk[[]string], err error) {
		funcs := map[string]func(check.ValCk[int]) check.ValCk[[]string]{
			"Length": check.SliceLength[[]string],
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

var strSlcMakerStrSlccheckerString = MakerInfo[[]string]{
	Args: []string{StringSliceCheckerName, "string"},

	MF: func(e *ast.CallExpr, fName string) (cf check.ValCk[[]string], err error) {
		funcs := map[string]func(check.ValCk[[]string], string) check.ValCk[[]string]{
			"Not": check.Not[[]string],
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

var strSlcMakerStrcheckerString = MakerInfo[[]string]{
	Args: []string{StringCheckerName, "string"},

	MF: func(e *ast.CallExpr, fName string) (cf check.ValCk[[]string], err error) {
		funcs := map[string]func(check.ValCk[string], string) check.ValCk[[]string]{
			"SliceAny": check.SliceAny[[]string],
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

var strSlcMakerStrchecker = MakerInfo[[]string]{
	Args: []string{StringCheckerName},

	MF: func(e *ast.CallExpr, fName string) (cf check.ValCk[[]string], err error) {
		funcs := map[string]func(check.ValCk[string]) check.ValCk[[]string]{
			"SliceAll": check.SliceAll[[]string],
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

		ckFunc, err := getCheckFunc[string](e, 0, StringCheckerName)
		if err != nil {
			return nil, err
		}

		return maker(ckFunc), nil
	},
}

var strSlcMakerMultiStrchecker = MakerInfo[[]string]{
	Args: []string{"...", StringCheckerName},

	MF: func(e *ast.CallExpr, fName string) (cf check.ValCk[[]string], err error) {
		funcs := map[string]func(...check.ValCk[string]) check.ValCk[[]string]{
			"SliceByPos": check.SliceByPos[[]string],
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

var strSlcMakerMultiStrSlcchecker = MakerInfo[[]string]{
	Args: []string{"...", StringSliceCheckerName},

	MF: func(e *ast.CallExpr, fName string) (cf check.ValCk[[]string], err error) {
		funcs := map[string]func(...check.ValCk[[]string]) check.ValCk[[]string]{
			"And": check.And[[]string],
			"Or":  check.Or[[]string],
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

		checkFuncs, err := getCheckFuncs[[]string](e, StringSliceCheckerName)
		if err != nil {
			return nil, err
		}

		return maker(checkFuncs...), nil
	},
}
