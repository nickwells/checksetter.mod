package checksetter_test

import (
	"testing"

	"github.com/nickwells/checksetter.mod/v4/checksetter"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

const (
	testDataDir = "testdata"
	avalSubDir  = "allowedVals"
)

var gfcAval = testhelper.GoldenFileCfg{
	DirNames:               []string{testDataDir, avalSubDir},
	Sfx:                    "txt",
	UpdFlagName:            "upd-aval-files",
	KeepBadResultsFlagName: "keep-bad-results-aval",
}

func init() {
	gfcAval.AddUpdateFlag()
	gfcAval.AddKeepBadResultsFlag()
}

func TestAllowedValues(t *testing.T) {
	var (
		intParser         = checksetter.FindParserOrPanic[int](checksetter.IntCheckerName)
		int64Parser       = checksetter.FindParserOrPanic[int64](checksetter.Int64CheckerName)
		float64Parser     = checksetter.FindParserOrPanic[float64](checksetter.Float64CheckerName)
		stringParser      = checksetter.FindParserOrPanic[string](checksetter.StringCheckerName)
		stringSliceParser = checksetter.FindParserOrPanic[[]string](checksetter.StringSliceCheckerName)
	)

	testCases := []struct {
		testhelper.ID
		name       string
		makerFuncs map[string][]string
	}{
		{
			ID:         testhelper.MkID("no-funcs"),
			name:       "nonesuch",
			makerFuncs: map[string][]string{},
		},
		{
			ID:         testhelper.MkID(checksetter.IntCheckerName),
			name:       checksetter.IntCheckerName,
			makerFuncs: intParser.MakerFuncs(),
		},
		{
			ID:         testhelper.MkID(checksetter.Int64CheckerName),
			name:       checksetter.Int64CheckerName,
			makerFuncs: int64Parser.MakerFuncs(),
		},
		{
			ID:         testhelper.MkID(checksetter.Float64CheckerName),
			name:       checksetter.Float64CheckerName,
			makerFuncs: float64Parser.MakerFuncs(),
		},
		{
			ID:         testhelper.MkID(checksetter.StringCheckerName),
			name:       checksetter.StringCheckerName,
			makerFuncs: stringParser.MakerFuncs(),
		},
		{
			ID:         testhelper.MkID(checksetter.StringSliceCheckerName),
			name:       checksetter.StringSliceCheckerName,
			makerFuncs: stringSliceParser.MakerFuncs(),
		},
	}

	for _, tc := range testCases {
		av := checksetter.AllowedValues(tc.name, tc.makerFuncs)
		gfcAval.Check(t, "Allowed Values", tc.Name, []byte(av))
	}
}
