package checksetter_test

import (
	"testing"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/checksetter.mod/v3/checksetter"
	"github.com/nickwells/testhelper.mod/testhelper"
)

// panicSafeCheckSetterString calls the CheckSetter method on the
// checksetter and returns values showing whether the call panicked and if so
// what error it found
func panicSafeCheckSetterString(cs checksetter.String) (panicked bool, panicVal interface{}) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicVal = r
		}
	}()
	cs.CheckSetter("dummy")
	return false, nil
}

func TestCheckSetterString(t *testing.T) {
	var checks []check.String

	badSetter := checksetter.String{}
	panicked, panicVal := panicSafeCheckSetterString(badSetter)
	testhelper.PanicCheckString(t, "bad setter - Value not set",
		panicked, true, panicVal,
		[]string{
			"String Check failed:",
			"the Value to be set is nil",
		})
	goodSetter := checksetter.String{Value: &checks}
	panicked, panicVal = panicSafeCheckSetterString(goodSetter)
	testhelper.PanicCheckString(t, "good setter",
		panicked, false, panicVal, []string{})
}

// panicSafeCheckSetterStringSlice calls the CheckSetter method on the
// checksetter and returns values showing whether the call panicked and if so
// what error it found
func panicSafeCheckSetterStringSlice(cs checksetter.StringSlice) (panicked bool, panicVal interface{}) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicVal = r
		}
	}()
	cs.CheckSetter("dummy")
	return false, nil
}

func TestCheckSetterStringSlice(t *testing.T) {
	var checks []check.StringSlice

	badSetter := checksetter.StringSlice{}
	panicked, panicVal := panicSafeCheckSetterStringSlice(badSetter)
	testhelper.PanicCheckString(t, "bad setter - Value not set",
		panicked, true, panicVal,
		[]string{
			"StringSlice Check failed:",
			"the Value to be set is nil",
		})
	goodSetter := checksetter.StringSlice{Value: &checks}
	panicked, panicVal = panicSafeCheckSetterStringSlice(goodSetter)
	testhelper.PanicCheckString(t, "good setter",
		panicked, false, panicVal, []string{})
}

// panicSafeCheckSetterInt64 calls the CheckSetter method on the
// checksetter and returns values showing whether the call panicked and if so
// what error it found
func panicSafeCheckSetterInt64(cs checksetter.Int64) (panicked bool, panicVal interface{}) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicVal = r
		}
	}()
	cs.CheckSetter("dummy")
	return false, nil
}

func TestCheckSetterInt64(t *testing.T) {
	var checks []check.Int64

	badSetter := checksetter.Int64{}
	panicked, panicVal := panicSafeCheckSetterInt64(badSetter)
	testhelper.PanicCheckString(t, "bad setter - Value not set",
		panicked, true, panicVal,
		[]string{
			"Int64 Check failed:",
			"the Value to be set is nil",
		})
	goodSetter := checksetter.Int64{Value: &checks}
	panicked, panicVal = panicSafeCheckSetterInt64(goodSetter)
	testhelper.PanicCheckString(t, "good setter",
		panicked, false, panicVal, []string{})
}

// panicSafeCheckSetterFloat64 calls the CheckSetter method on the
// checksetter and returns values showing whether the call panicked and if so
// what error it found
func panicSafeCheckSetterFloat64(cs checksetter.Float64) (panicked bool, panicVal interface{}) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicVal = r
		}
	}()
	cs.CheckSetter("dummy")
	return false, nil
}

func TestCheckSetterFloat64(t *testing.T) {
	var checks []check.Float64

	badSetter := checksetter.Float64{}
	panicked, panicVal := panicSafeCheckSetterFloat64(badSetter)
	testhelper.PanicCheckString(t, "bad setter - Value not set",
		panicked, true, panicVal,
		[]string{
			"Float64 Check failed:",
			"the Value to be set is nil",
		})
	goodSetter := checksetter.Float64{Value: &checks}
	panicked, panicVal = panicSafeCheckSetterFloat64(goodSetter)
	testhelper.PanicCheckString(t, "good setter",
		panicked, false, panicVal, []string{})
}

func TestIsASetter(t *testing.T) {
	setterF := checksetter.Float64{}
	_ = setterF.AllowedValues()
	setterI := checksetter.Int64{}
	_ = setterI.AllowedValues()
	setterS := checksetter.String{}
	_ = setterS.AllowedValues()
	setterSS := checksetter.StringSlice{}
	_ = setterSS.AllowedValues()
}
