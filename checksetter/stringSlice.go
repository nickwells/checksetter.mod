package checksetter

import (
	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/param.mod/v5/param/psetter"
)

// StringSlice can be used to set a list of checkers for a slice of
// strings.
type StringSlice struct {
	psetter.ValueReqMandatory

	Value *[]check.StringSlice
}

// SetWithVal (called when a value follows the parameter) splits the value
// into a slice of check.StringSlice's and sets the Value accordingly. It
// will return an error if a check is breached.
func (s StringSlice) SetWithVal(_ string, paramVal string) error {
	v, err := stringSliceCFParse(paramVal)
	if err != nil {
		return err
	}
	*s.Value = v

	return nil
}

// AllowedValues returns a description of the allowed values. It includes the
// separator to be used
func (s StringSlice) AllowedValues() string {
	return allowedValues(strSlcCFName)
}

// CurrentValue returns the current setting of the parameter value
func (s StringSlice) CurrentValue() string {
	return currentValue(len(*s.Value))
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s StringSlice) CheckSetter(name string) {
	if s.Value == nil {
		panic(psetter.NilValueMessage(name, "checksetter.StringSlice"))
	}
}
