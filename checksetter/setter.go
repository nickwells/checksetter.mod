package checksetter

import (
	"fmt"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/param.mod/v5/param/psetter"
)

type Setter[T any] struct {
	psetter.ValueReqMandatory
	Parser *Parser[T]

	Value *[]check.ValCk[T]
}

// SetWithVal (called when a value follows the parameter) splits the value
// into a slice of check.Int64's and sets the Value accordingly.
func (s Setter[T]) SetWithVal(_ string, paramVal string) error {
	v, err := s.Parser.Parse(paramVal)
	if err != nil {
		return err
	}
	*s.Value = v

	return nil
}

// AllowedValues returns a description of the allowed values. It includes the
// separator to be used
func (s Setter[T]) AllowedValues() string {
	return AllowedValues(s.Parser.CheckerName(), s.Parser.MakerFuncs())
}

// CurrentValue returns the current setting of the parameter value
func (s Setter[T]) CurrentValue() string {
	switch len(*s.Value) {
	case 0:
		return "no checks"
	case 1:
		return "one check"
	}
	return fmt.Sprintf("%d checks", len(*s.Value))
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
}
