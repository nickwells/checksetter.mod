package checksetter

import (
	"fmt"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/param.mod/v6/psetter"
)

// Setter satisfies the param.Setter interface. Important points of
// difference are that you need to provide both the Parser (use the
// checksetter.FindParserOrPanic func) and the Value when initialising the
// Setter. Also, you will need to pass the address of the Setter rather than
// the Setter itself, this is because the SetWithVal method takes a pointer
// receiver.
type Setter[T any] struct {
	psetter.ValueReqMandatory
	Parser *Parser[T]

	Value *[]check.ValCk[T]

	paramVal string
	valSet   bool
}

// SetWithVal (called when a value follows the parameter) splits the value
// into a slice of check.Int64's and sets the Value accordingly.
func (s *Setter[T]) SetWithVal(_ string, paramVal string) error {
	v, err := s.Parser.Parse(paramVal)
	if err != nil {
		return err
	}

	*s.Value = v
	s.paramVal = paramVal
	s.valSet = true

	return nil
}

// AllowedValues returns a description of the allowed values. It includes the
// separator to be used
func (s Setter[T]) AllowedValues() string {
	return AllowedValues(s.Parser.CheckerName(), s.Parser.MakerFuncs())
}

// CurrentValue returns the current setting of the parameter value
func (s Setter[T]) CurrentValue() string {
	val := ""

	switch len(*s.Value) {
	case 0:
		val = "no checks"
	case 1:
		val = "one check"
	default:
		val = fmt.Sprintf("%d checks", len(*s.Value))
	}

	if s.valSet {
		val += fmt.Sprintf(": %q", s.paramVal)
	}

	return val
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil or the Parser is nil.
func (s Setter[T]) CheckSetter(name string) {
	if s.Value == nil {
		var v T

		panic(
			psetter.NilValueMessage(name,
				fmt.Sprintf("checksetter.Setter[%T]", v)))
	}

	if s.Parser == nil {
		var v T

		panic(fmt.Sprintf(
			"The Parser for checksetter.Setter[%T] has not been set",
			v))
	}

	if len(s.Parser.Makers()) == 0 {
		var v T

		panic(fmt.Sprintf(
			"The Parser for checksetter.Setter[%T] can't make any check-funcs",
			v))
	}
}
