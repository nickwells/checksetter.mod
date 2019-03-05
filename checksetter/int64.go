package checksetter

import (
	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/param.mod/v2/param"
	"github.com/nickwells/param.mod/v2/param/psetter"
)

// Int64 can be used to set a list of checkers for an int64.
type Int64 struct {
	param.ValueReqMandatory

	Value *[]check.Int64
}

// SetWithVal (called when a value follows the parameter) splits the value
// into a slice of check.Int64's and sets the Value accordingly.
func (s Int64) SetWithVal(_ string, paramVal string) error {
	v, err := int64CFParse(paramVal)
	if err != nil {
		return err
	}
	*s.Value = v

	return nil
}

// AllowedValues returns a description of the allowed values. It includes the
// separator to be used
func (s Int64) AllowedValues() string {
	return allowedValues(int64CFName)
}

// CurrentValue returns the current setting of the parameter value
func (s Int64) CurrentValue() string {
	return currentValue(len(*s.Value))
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s Int64) CheckSetter(name string) {
	if s.Value == nil {
		panic(psetter.NilValueMessage(name, "checksetter.Int64"))
	}
}
