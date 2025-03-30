package checksetter

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/nickwells/check.mod/v2/check"
)

// Int64CheckerName is the value to use to select the Parser to use when
// creating checkers for int64 values
const Int64CheckerName = "int64-checker"

var (
	i64MakerArgs = []string{}
	i64Maker     = MakerInfo[int64]{
		Args: i64MakerArgs,

		MF: func(e *ast.CallExpr, fName string) (
			cf check.ValCk[int64], err error,
		) {
			funcs := map[string]check.ValCk[int64]{
				"OK": check.ValOK[int64],
			}

			maker, ok := funcs[fName]
			if !ok {
				return nil, fmt.Errorf(errFmtUnknownFunc, fName)
			}

			defer func() {
				if err != nil {
					err = fmt.Errorf("%s(%s): %w",
						fName, strings.Join(i64MakerArgs, ", "), err)
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
	i64MakerI64Args = []string{"int64"}
	i64MakerI64     = MakerInfo[int64]{
		Args: i64MakerI64Args,

		MF: func(e *ast.CallExpr, fName string) (
			cf check.ValCk[int64], err error,
		) {
			funcs := map[string]func(int64) check.ValCk[int64]{
				"EQ":          check.ValEQ[int64],
				"GT":          check.ValGT[int64],
				"GE":          check.ValGE[int64],
				"LT":          check.ValLT[int64],
				"LE":          check.ValLE[int64],
				"Divides":     check.ValDivides[int64],
				"IsAMultiple": check.ValIsAMultiple[int64],
			}

			maker, ok := funcs[fName]
			if !ok {
				return nil, fmt.Errorf(errFmtUnknownFunc, fName)
			}

			defer func() {
				if err != nil {
					err = fmt.Errorf("%s(%s): %w",
						fName, strings.Join(i64MakerI64Args, ", "), err)
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

			i, err := getInt64(e.Args[0])
			if err != nil {
				return nil, err
			}

			return maker(i), nil
		},
	}
)

var (
	i64MakerI64I64Args = []string{"int64", "int64"}
	i64MakerI64I64     = MakerInfo[int64]{
		Args: i64MakerI64I64Args,

		MF: func(e *ast.CallExpr, fName string) (
			cf check.ValCk[int64], err error,
		) {
			funcs := map[string]func(int64, int64) check.ValCk[int64]{
				"Between": check.ValBetween[int64],
			}

			maker, ok := funcs[fName]
			if !ok {
				return nil, fmt.Errorf(errFmtUnknownFunc, fName)
			}

			defer func() {
				if err != nil {
					err = fmt.Errorf("%s(%s): %w",
						fName, strings.Join(i64MakerI64I64Args, ", "), err)
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

			i1, err := getInt64(e.Args[0])
			if err != nil {
				return nil, err
			}

			i2, err := getInt64(e.Args[1])
			if err != nil {
				return nil, err
			}

			return maker(i1, i2), nil
		},
	}
)

var (
	i64MakerI64checkerStringArgs = []string{Int64CheckerName, "string"}
	i64MakerI64checkerString     = MakerInfo[int64]{
		Args: i64MakerI64checkerStringArgs,

		MF: func(e *ast.CallExpr, fName string) (
			cf check.ValCk[int64], err error,
		) {
			funcs := map[string]func(check.ValCk[int64], string) check.ValCk[int64]{
				"Not": check.Not[int64],
			}

			maker, ok := funcs[fName]
			if !ok {
				return nil, fmt.Errorf(errFmtUnknownFunc, fName)
			}

			defer func() {
				if err != nil {
					err = fmt.Errorf("%s(%s): %w",
						fName, strings.Join(i64MakerI64checkerStringArgs, ", "), err)
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

			ckFunc, err := getCheckFunc[int64](e, 0, Int64CheckerName)
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
	i64MakerMultiI64checkerArgs = []string{"...", Int64CheckerName}
	i64MakerMultiI64checker     = MakerInfo[int64]{
		Args: i64MakerMultiI64checkerArgs,

		MF: func(e *ast.CallExpr, fName string) (
			cf check.ValCk[int64], err error,
		) {
			funcs := map[string]func(...check.ValCk[int64]) check.ValCk[int64]{
				"And": check.And[int64],
				"Or":  check.Or[int64],
			}

			maker, ok := funcs[fName]
			if !ok {
				return nil, fmt.Errorf(errFmtUnknownFunc, fName)
			}

			defer func() {
				if err != nil {
					err = fmt.Errorf("%s(%s): %w",
						fName,
						strings.Join(i64MakerMultiI64checkerArgs, ", "), err)
				}
			}()
			defer func() {
				if r := recover(); r != nil {
					cf = nil
					err = fmt.Errorf("%v", r)
				}
			}()

			checkFuncs, err := getCheckFuncs[int64](e, Int64CheckerName)
			if err != nil {
				return nil, err
			}

			return maker(checkFuncs...), nil
		},
	}
)
