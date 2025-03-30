package checksetter

import (
	"fmt"
	"maps"
	"slices"
)

func init() {
	_, err := MakeParser(
		Float64CheckerName,
		map[string]MakerInfo[float64]{
			"OK":      f64Maker,
			"GT":      f64MakerF64,
			"GE":      f64MakerF64,
			"LT":      f64MakerF64,
			"LE":      f64MakerF64,
			"Between": f64MakerF64F64,
			"Not":     f64MakerF64checkerString,
			"And":     f64MakerMultiF64checker,
			"Or":      f64MakerMultiF64checker,
		})
	if err != nil {
		panic(err)
	}

	_, err = MakeParser(
		IntCheckerName,
		map[string]MakerInfo[int]{
			"OK":          iMaker,
			"EQ":          iMakerI,
			"GT":          iMakerI,
			"GE":          iMakerI,
			"LT":          iMakerI,
			"LE":          iMakerI,
			"Divides":     iMakerI,
			"IsAMultiple": iMakerI,
			"Between":     iMakerII,
			"Not":         iMakerIcheckerString,
			"And":         iMakerMultiIchecker,
			"Or":          iMakerMultiIchecker,
		})
	if err != nil {
		panic(err)
	}

	_, err = MakeParser(
		Int64CheckerName,
		map[string]MakerInfo[int64]{
			"OK":          i64Maker,
			"EQ":          i64MakerI64,
			"GT":          i64MakerI64,
			"GE":          i64MakerI64,
			"LT":          i64MakerI64,
			"LE":          i64MakerI64,
			"Divides":     i64MakerI64,
			"IsAMultiple": i64MakerI64,
			"Between":     i64MakerI64I64,
			"Not":         i64MakerI64checkerString,
			"And":         i64MakerMultiI64checker,
			"Or":          i64MakerMultiI64checker,
		})
	if err != nil {
		panic(err)
	}

	_, err = MakeParser(
		StringCheckerName,
		map[string]MakerInfo[string]{
			"OK":             strMaker,
			"EQ":             strMakerStr,
			"GT":             strMakerStr,
			"GE":             strMakerStr,
			"LT":             strMakerStr,
			"LE":             strMakerStr,
			"HasPrefix":      strMakerStr,
			"HasSuffix":      strMakerStr,
			"Length":         strMakerIchecker,
			"MatchesPattern": strMakerRegexpStr,
			"Not":            strMakerStrcheckerString,
			"And":            strMakerMultiStrchecker,
			"Or":             strMakerMultiStrchecker,
		})
	if err != nil {
		panic(err)
	}

	_, err = MakeParser(
		StringSliceCheckerName,
		map[string]MakerInfo[[]string]{
			"OK":         strSlcMaker,
			"NoDups":     strSlcMaker,
			"Length":     strSlcMakerIchecker,
			"Not":        strSlcMakerStrSlccheckerString,
			"SliceAny":   strSlcMakerStrcheckerString,
			"SliceAll":   strSlcMakerStrchecker,
			"SliceByPos": strSlcMakerMultiStrchecker,
			"And":        strSlcMakerMultiStrSlcchecker,
			"Or":         strSlcMakerMultiStrSlcchecker,
		})
	if err != nil {
		panic(err)
	}
}

// anyParser is a minimal interface that can be satisfied by any Parser
// because it does not take or return any type-specific values
type anyParser interface {
	CheckerName() string
	Makers() []string
	Args(string) ([]string, error)
	MakerFuncs() map[string][]string
}

// parserRegister records the parsers that we have created. Note that it
// doesn't hold Parsers but instead holds 'anyParser' values, this is because
// each Parser is of a different type . When you retrieve the parser you
// should check that it really is of the type that you expect.
var parserRegister = map[string]anyParser{}

// FindParser finds a pre-registered parser with the given checker name. It
// will return nil if there is no such Parser already registered. Note that
// the return type is 'any'; it is the caller's responsibility to check that
// it is of the type required.
func FindParser[T any](checkerName string) (*Parser[T], error) {
	anyParser, ok := parserRegister[checkerName]
	if !ok {
		return nil,
			fmt.Errorf("there is no Parser registered for %q", checkerName)
	}

	parser, ok := anyParser.(*Parser[T])
	if !ok {
		return nil,
			fmt.Errorf("the Parser for %q is of the wrong type (%T)",
				checkerName, anyParser)
	}

	return parser, nil
}

// FindParserOrPanic finds a pre-registered parser with the given checker
// name or else it panics
func FindParserOrPanic[T any](checkerName string) *Parser[T] {
	parser, err := FindParser[T](checkerName)
	if err != nil {
		panic(err)
	}

	return parser
}

// ParsersAvailable returns a sorted list of all the available parsers
func ParsersAvailable() []string {
	return slices.Sorted(maps.Keys(parserRegister))
}
