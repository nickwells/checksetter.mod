package checksetter

import (
	"testing"

	"github.com/nickwells/testhelper.mod/testhelper"
)

// parseTestInfo holds details of parse tests
type parseTestInfo struct {
	testhelper.ID
	testhelper.ExpErr
	s           string
	parser      string
	lenExpected int
}

func TestParse(t *testing.T) {
	testCases := []parseTestInfo{
		{
			ID:          testhelper.MkID("good - 1 check"),
			s:           `GT(5)`,
			parser:      "int64",
			lenExpected: 1,
		},
		{
			ID:          testhelper.MkID("good - 2 checks"),
			s:           `GT(5), LT(9)`,
			parser:      "int64",
			lenExpected: 2,
		},
		{
			ID:          testhelper.MkID("good - 1 check"),
			s:           `GT(5.0)`,
			parser:      "float64",
			lenExpected: 1,
		},
		{
			ID:          testhelper.MkID("good - 2 checks"),
			s:           `GT(5.0), LT(9.0)`,
			parser:      "float64",
			lenExpected: 2,
		},
		{
			ID:          testhelper.MkID("good - 1 check"),
			s:           `Length(GT(5))`,
			parser:      "string",
			lenExpected: 1,
		},
		{
			ID:          testhelper.MkID("good - 2 checks"),
			s:           `Length(GT(5)), Length(LT(9))`,
			parser:      "string",
			lenExpected: 2,
		},
		{
			ID:          testhelper.MkID("good - 1 check"),
			s:           `Length(GT(5))`,
			parser:      "stringSlice",
			lenExpected: 1,
		},
		{
			ID:          testhelper.MkID("good - 2 checks"),
			s:           `Length(GT(5)), Length(LT(9))`,
			parser:      "stringSlice",
			lenExpected: 2,
		},
		{
			ID:     testhelper.MkID("bad - syntax: ,,"),
			s:      `,,`,
			parser: "all",
			ExpErr: testhelper.MkExpErr("expected operand, found ','"),
		},
		{
			ID:     testhelper.MkID("bad - syntax: }"),
			s:      `}`,
			parser: "all",
			ExpErr: testhelper.MkExpErr("expected 'EOF', found '}'"),
		},
		{
			ID:     testhelper.MkID("bad - syntax: {"),
			s:      `{`,
			parser: "all",
			ExpErr: testhelper.MkExpErr(
				"missing ',' before newline in composite literal"),
		},
		{
			ID:     testhelper.MkID("bad - syntax: {}"),
			s:      `{}`,
			parser: "all",
			ExpErr: testhelper.MkExpErr(
				"bad function: unexpected type: *ast.CompositeLit"),
		},
	}

	for _, tc := range testCases {
		if tc.parser == "int64" || tc.parser == "all" {
			slc, err := int64CFParse(tc.s)
			checkParseResults(t, len(slc), err, tc)
		}
		if tc.parser == "float64" || tc.parser == "all" {
			slc, err := float64CFParse(tc.s)
			checkParseResults(t, len(slc), err, tc)
		}
		if tc.parser == "string" || tc.parser == "all" {
			slc, err := stringCFParse(tc.s)
			checkParseResults(t, len(slc), err, tc)
		}
		if tc.parser == "stringSlice" || tc.parser == "all" {
			slc, err := stringSliceCFParse(tc.s)
			checkParseResults(t, len(slc), err, tc)
		}
	}
}

// checkParseResults checks that the results of parsing are as expected
func checkParseResults(t *testing.T, slcLen int, err error, tc parseTestInfo) {
	t.Helper()

	if testhelper.CheckExpErr(t, err, tc) && err == nil {
		testhelper.DiffInt(t, tc.IDStr(), "# of checks", slcLen, tc.lenExpected)
	}
}
