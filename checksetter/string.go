package checksetter

import (
	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/param.mod/v5/param/psetter"
)

// String can be used to set a list of checkers for a string.
type String struct {
	psetter.ValueReqMandatory

	Value *[]check.String
}

// SetWithVal (called when a value follows the parameter) splits the value
// into a slice of check.String's and sets the Value accordingly.
func (s String) SetWithVal(_ string, paramVal string) error {
	v, err := stringCFParse(paramVal)
	if err != nil {
		return err
	}
	*s.Value = v

	return nil
}

// AllowedValues returns a description of the allowed values. It includes the
// separator to be used
func (s String) AllowedValues() string {
	return allowedValues(strCFName)
}

// CurrentValue returns the current setting of the parameter value
func (s String) CurrentValue() string {
	return currentValue(len(*s.Value))
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s String) CheckSetter(name string) {
	if s.Value == nil {
		panic(psetter.NilValueMessage(name, "checksetter.String"))
	}
}
