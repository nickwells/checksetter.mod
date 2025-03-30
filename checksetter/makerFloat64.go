package checksetter

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/nickwells/check.mod/v2/check"
)

// Float64CheckerName is the value to use to select the Parser to use when
// creating checkers for float64 values
const Float64CheckerName = "float64-checker"

var (
	f64MakerArgs = []string{}
	f64Maker     = MakerInfo[float64]{
		Args: f64MakerArgs,

		MF: func(e *ast.CallExpr, fName string) (
			cf check.ValCk[float64], err error,
		) {
			funcs := map[string]check.ValCk[float64]{
				"OK": check.ValOK[float64],
			}

			maker, ok := funcs[fName]
			if !ok {
				return nil, fmt.Errorf("Unknown function: %q", fName)
			}

			defer func() {
				if err != nil {
					err = fmt.Errorf("%s(%s): %w",
						fName, strings.Join(f64MakerArgs, ", "), err)
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
	f64MakerF64Args = []string{"float64"}
	f64MakerF64     = MakerInfo[float64]{
		Args: f64MakerF64Args,

		MF: func(e *ast.CallExpr, fName string) (
			cf check.ValCk[float64], err error,
		) {
			funcs := map[string]func(float64) check.ValCk[float64]{
				"GT": check.ValGT[float64],
				"GE": check.ValGE[float64],
				"LT": check.ValLT[float64],
				"LE": check.ValLE[float64],
			}

			maker, ok := funcs[fName]
			if !ok {
				return nil, fmt.Errorf("Unknown function: %q", fName)
			}

			defer func() {
				if err != nil {
					err = fmt.Errorf("%s(%s): %w",
						fName, strings.Join(f64MakerF64Args, ", "), err)
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

			f, err := getFloat64(e.Args[0])
			if err != nil {
				return nil, err
			}

			return maker(f), nil
		},
	}
)

var (
	f64MakerF64F64Args = []string{"float64", "float64"}
	f64MakerF64F64     = MakerInfo[float64]{
		Args: f64MakerF64F64Args,

		MF: func(e *ast.CallExpr, fName string) (
			cf check.ValCk[float64], err error,
		) {
			funcs := map[string]func(float64, float64) check.ValCk[float64]{
				"Between": check.ValBetween[float64],
			}

			maker, ok := funcs[fName]
			if !ok {
				return nil, fmt.Errorf("Unknown function: %q", fName)
			}

			defer func() {
				if err != nil {
					err = fmt.Errorf("%s(%s): %w",
						fName, strings.Join(f64MakerF64F64Args, ", "), err)
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

			f1, err := getFloat64(e.Args[0])
			if err != nil {
				return nil, err
			}

			f2, err := getFloat64(e.Args[1])
			if err != nil {
				return nil, err
			}

			return maker(f1, f2), nil
		},
	}
)

var (
	f64MakerF64checkerStringArgs = []string{Float64CheckerName, "string"}
	f64MakerF64checkerString     = MakerInfo[float64]{
		Args: f64MakerF64checkerStringArgs,

		MF: func(e *ast.CallExpr, fName string) (
			cf check.ValCk[float64], err error,
		) {
			funcs := map[string]func(check.ValCk[float64], string) check.ValCk[float64]{
				"Not": check.Not[float64],
			}

			maker, ok := funcs[fName]
			if !ok {
				return nil, fmt.Errorf("Unknown function: %q", fName)
			}

			defer func() {
				if err != nil {
					err = fmt.Errorf("%s(%s): %w",
						fName,
						strings.Join(f64MakerF64checkerStringArgs, ", "),
						err)
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

			ckFunc, err := getCheckFunc[float64](e, 0, Float64CheckerName)
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
	f64MakerMultiF64checkerArgs = []string{"...", Float64CheckerName}
	f64MakerMultiF64checker     = MakerInfo[float64]{
		Args: f64MakerMultiF64checkerArgs,

		MF: func(e *ast.CallExpr, fName string) (
			cf check.ValCk[float64], err error,
		) {
			funcs := map[string]func(...check.ValCk[float64]) check.ValCk[float64]{
				"And": check.And[float64],
				"Or":  check.Or[float64],
			}

			maker, ok := funcs[fName]
			if !ok {
				return nil, fmt.Errorf("Unknown function: %q", fName)
			}

			defer func() {
				if err != nil {
					err = fmt.Errorf("%s(%s): %w",
						fName, strings.Join(f64MakerMultiF64checkerArgs, ", "), err)
				}
			}()
			defer func() {
				if r := recover(); r != nil {
					cf = nil
					err = fmt.Errorf("%v", r)
				}
			}()

			checkFuncs, err := getCheckFuncs[float64](e, Float64CheckerName)
			if err != nil {
				return nil, err
			}

			return maker(checkFuncs...), nil
		},
	}
)
