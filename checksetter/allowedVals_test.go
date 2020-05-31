package checksetter_test

import (
	"strings"
	"testing"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/checksetter.mod/v3/checksetter"
	"github.com/nickwells/testhelper.mod/testhelper"
)

const (
	allowedValsDir    = "testdata"
	allowedValsSubDir = "allowedVals"
)

var aValGFC = testhelper.GoldenFileCfg{
	DirNames: []string{allowedValsDir, allowedValsSubDir},
	Sfx:      "txt",

	UpdFlagName: "upd-avals",
}

func init() {
	aValGFC.AddUpdateFlag()
}

// avalTNameAbbv returns the abbreviated test name - with a common prefix
// removed
func avalTNameAbbr(t *testing.T) string {
	t.Helper()
	return strings.TrimPrefix(t.Name(), "TestAllowedValues")
}
func TestAllowedValuesStringSlice(t *testing.T) {
	var checks []check.StringSlice
	setter := checksetter.StringSlice{Value: &checks}

	val := []byte(setter.AllowedValues())
	name := avalTNameAbbr(t)

	aValGFC.Check(t, name, name, val)
}

func TestAllowedValuesString(t *testing.T) {
	var checks []check.String
	setter := checksetter.String{Value: &checks}

	val := []byte(setter.AllowedValues())
	name := avalTNameAbbr(t)

	aValGFC.Check(t, name, name, val)
}

func TestAllowedValuesInt64(t *testing.T) {
	var checks []check.Int64
	setter := checksetter.Int64{Value: &checks}

	val := []byte(setter.AllowedValues())
	name := avalTNameAbbr(t)

	aValGFC.Check(t, name, name, val)
}

func TestAllowedValuesFloat64(t *testing.T) {
	var checks []check.Float64
	setter := checksetter.Float64{Value: &checks}

	val := []byte(setter.AllowedValues())
	name := avalTNameAbbr(t)

	aValGFC.Check(t, name, name, val)
}
