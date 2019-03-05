package checksetter_test

import (
	"fmt"
	"testing"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/checksetter.mod/checksetter"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestChkStringSlice(t *testing.T) {
	testCases := []struct {
		name               string
		arg                string
		errExpected        bool
		errContains        []string
		ssToCheck          []string
		sliceErrorExpected bool
		sliceErrContains   []string
	}{
		{
			name:      "ok - one check function (no params), slice ok",
			arg:       "NoDups",
			ssToCheck: []string{"one", "two"},
		},
		{
			name:      "ok - one check function (one int param), slice ok",
			arg:       "LenEQ(2)",
			ssToCheck: []string{"one", "two"},
		},
		{
			name:      "ok - one check function (two int param), slice ok",
			arg:       "LenBetween(2, 3)",
			ssToCheck: []string{"one", "two"},
		},
		{
			name:      "ok - one check function (check string param), slice ok",
			arg:       `String(HasPrefix("rc"))`,
			ssToCheck: []string{"rc001", "rc002"},
		},
		{
			name:      "ok - three check functions, slice ok",
			arg:       `String(HasPrefix("rc")), LenEQ(2), NoDups`,
			ssToCheck: []string{"rc001", "rc002"},
		},
		{
			name:               "ok - three check functions, slice bad",
			arg:                `String(And(HasPrefix("rc"), LenEQ(5))), LenEQ(3), NoDups`,
			ssToCheck:          []string{"rc001", "rc002", "rc002"},
			sliceErrorExpected: true,
			sliceErrContains: []string{
				"list entries: 1 and 2 are duplicates, both are: rc002",
			},
		},
		{
			name:        "bad - one check function (two int param) - unknown func",
			arg:         "LenBetweenXXX(2, 3)",
			errExpected: true,
			errContains: []string{
				"bad function",
				"LenBetweenXXX",
			},
		},
		{
			name:        "bad - one check function (two int param) - invalid params",
			arg:         "LenBetween(5, 3)",
			errExpected: true,
			errContains: []string{
				"bad function: ",
				"Cannot create the check.StringSlice func: ",
				"Impossible checks passed to StringSliceLenBetween",
			},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)
		var checks []check.StringSlice
		var checker = checksetter.StringSlice{Value: &checks}
		err := checker.SetWithVal("dummy", tc.arg)
		if testhelper.CheckError(t, tcID, err,
			tc.errExpected, tc.errContains) &&
			err == nil {
			for i, chk := range checks {
				if chk == nil {
					t.Logf("%s:\n", tcID)
					t.Errorf("\t: nil check found at check slice element %d", i)
					continue
				}
				if err = chk(tc.ssToCheck); err != nil {
					break
				}
			}
			testhelper.CheckError(t, tcID+" (slice checks)",
				err, tc.sliceErrorExpected, tc.sliceErrContains)

		}
	}
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

func TestCurrentValueStringSlice(t *testing.T) {
	var checks []check.StringSlice
	setter := checksetter.StringSlice{Value: &checks}

	for i, expVal := range []string{"no checks", "one check", "2 checks"} {
		val := setter.CurrentValue()
		if expVal != val {
			t.Logf("after %d additions\n", i)
			t.Logf("\tcurrent value should be '%s'\n", expVal)
			t.Errorf("\t                but was '%s'\n", val)
		}
		checks = append(checks, nil)
	}
}
