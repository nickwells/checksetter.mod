package checksetter

import (
	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/param.mod/v5/param/psetter"
)

// Float64 can be used to set a list of checkers for a float64.
type Float64 struct {
	psetter.ValueReqMandatory

	Value *[]check.Float64
}

// SetWithVal (called when a value follows the parameter) splits the value
// into a slice of check.Float64's and sets the Value accordingly.
func (s Float64) SetWithVal(_ string, paramVal string) error {
	v, err := float64CFParse(paramVal)
	if err != nil {
		return err
	}
	*s.Value = v

	return nil
}

// AllowedValues returns a description of the allowed values. It includes the
// separator to be used
func (s Float64) AllowedValues() string {
	return allowedValues(float64CFName)
}

// CurrentValue returns the current setting of the parameter value
func (s Float64) CurrentValue() string {
	return currentValue(len(*s.Value))
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s Float64) CheckSetter(name string) {
	if s.Value == nil {
		panic(psetter.NilValueMessage(name, "checksetter.Float64"))
	}
}
