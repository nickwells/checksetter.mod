package checksetter

import (
	"fmt"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/param.mod/v2/param"
)

// String can be used to set a list of checkers for a string.
type String struct {
	param.ValueReqMandatory

	Value *[]check.String
}

// SetWithVal (called when a value follows the parameter) splits the value
// into a slice of check.String's and sets the Value accordingly.
func (s String) SetWithVal(_ string, paramVal string) error {
	v, err := strCFParse(paramVal)
	if err != nil {
		return err
	}
	*s.Value = v

	return nil
}

// AllowedValues returns a description of the allowed values. It includes the
// separator to be used
func (s String) AllowedValues() string {
	rval := "a list of " + strCFName + " separated by ','.\n"
	rval += `
Write the checks as if you were writing code.

The functions recognised are:` + allowedVals(strCFName)

	return rval
}

// CurrentValue returns the current setting of the parameter value
func (s String) CurrentValue() string {
	switch len(*s.Value) {
	case 0:
		return "no checks"
	case 1:
		return "one check"
	}

	return fmt.Sprintf("%d checks", len(*s.Value))
}

// CheckSetter panics if the setter has not been properly created - if the
// Value is nil.
func (s String) CheckSetter(name string) {
	if s.Value == nil {
		panic(name +
			": String Check failed:" +
			" the Value to be set is nil")
	}
}
