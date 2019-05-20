package checksetter_test

import (
	"flag"
	"strings"
	"testing"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/checksetter.mod/v2/checksetter"
	"github.com/nickwells/testhelper.mod/testhelper"
)

const (
	allowedValsDir    = "testdata"
	allowedValsSubDir = "allowedVals"
)

var aValGFC = testhelper.GoldenFileCfg{
	DirNames: []string{allowedValsDir, allowedValsSubDir},
	Sfx:      "txt",
}

var updateAVals = flag.Bool("upd-avals", false,
	"update the files holding the allowed values messages")

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

	testhelper.CheckAgainstGoldenFile(t, name, val,
		aValGFC.PathName(name), *updateAVals)
}

func TestAllowedValuesString(t *testing.T) {
	var checks []check.String
	setter := checksetter.String{Value: &checks}

	val := []byte(setter.AllowedValues())
	name := avalTNameAbbr(t)

	testhelper.CheckAgainstGoldenFile(t, name, val,
		aValGFC.PathName(name), *updateAVals)
}

func TestAllowedValuesInt64(t *testing.T) {
	var checks []check.Int64
	setter := checksetter.Int64{Value: &checks}

	val := []byte(setter.AllowedValues())
	name := avalTNameAbbr(t)

	testhelper.CheckAgainstGoldenFile(t, name, val,
		aValGFC.PathName(name), *updateAVals)
}

func TestAllowedValuesFloat64(t *testing.T) {
	var checks []check.Float64
	setter := checksetter.Float64{Value: &checks}

	val := []byte(setter.AllowedValues())
	name := avalTNameAbbr(t)

	testhelper.CheckAgainstGoldenFile(t, name, val,
		aValGFC.PathName(name), *updateAVals)
}
