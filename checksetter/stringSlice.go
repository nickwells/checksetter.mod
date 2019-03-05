package checksetter

import (
	"errors"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/param.mod/v2/param"
	"github.com/nickwells/param.mod/v2/param/psetter"
)

// StringSlice can be used to set a list of checkers for a slice of
// strings.
type StringSlice struct {
	Value *[]check.StringSlice
}

// ValueReq returns param.Mandatory indicating that some value must follow
// the parameter
func (s StringSlice) ValueReq() param.ValueReq { return param.Mandatory }

// Set (called when there is no following value) returns an error
func (s StringSlice) Set(_ string) error {
	return errors.New("no value given (it should be followed by '=...')")
}

// SetWithVal (called when a value follows the parameter) splits the value
// into a slice of check.StringSlice's and sets the Value accordingly. It
// will return an error if a check is breached.
func (s StringSlice) SetWithVal(_ string, paramVal string) error {
	v, err := strSlcCFParse(paramVal)
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
