package checksetter_test

import (
	"testing"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/checksetter.mod/v3/checksetter"
)

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

func TestCurrentValueString(t *testing.T) {
	var checks []check.String
	setter := checksetter.String{Value: &checks}

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

func TestCurrentValueInt64(t *testing.T) {
	var checks []check.Int64
	setter := checksetter.Int64{Value: &checks}

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

func TestCurrentValueFloat64(t *testing.T) {
	var checks []check.Float64
	setter := checksetter.Float64{Value: &checks}

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
