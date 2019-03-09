package checksetter

import (
	"fmt"
	"testing"

	"github.com/nickwells/testhelper.mod/testhelper"
)

// parseTestInfo holds details of parse tests
type parseTestInfo struct {
	name             string
	s                string
	parser           string
	lenExpected      int
	errExpected      bool
	errShouldContain []string
}

func TestParse(t *testing.T) {
	testCases := []parseTestInfo{
		{
			name:        "good - 1 check",
			s:           `GT(5)`,
			parser:      "int64",
			lenExpected: 1,
		},
		{
			name:        "good - 2 checks",
			s:           `GT(5), LT(9)`,
			parser:      "int64",
			lenExpected: 2,
		},
		{
			name:        "good - 1 check",
			s:           `GT(5.0)`,
			parser:      "float64",
			lenExpected: 1,
		},
		{
			name:        "good - 2 checks",
			s:           `GT(5.0), LT(9.0)`,
			parser:      "float64",
			lenExpected: 2,
		},
		{
			name:        "good - 1 check",
			s:           `LenGT(5)`,
			parser:      "string",
			lenExpected: 1,
		},
		{
			name:        "good - 2 checks",
			s:           `LenGT(5), LenLT(9)`,
			parser:      "string",
			lenExpected: 2,
		},
		{
			name:        "good - 1 check",
			s:           `LenGT(5)`,
			parser:      "stringSlice",
			lenExpected: 1,
		},
		{
			name:        "good - 2 checks",
			s:           `LenGT(5), LenLT(9)`,
			parser:      "stringSlice",
			lenExpected: 2,
		},
		{
			name:             "bad - syntax: ,,",
			s:                `,,`,
			parser:           "all",
			errExpected:      true,
			errShouldContain: []string{"expected operand, found ','"},
		},
		{
			name:             "bad - syntax: }",
			s:                `}`,
			parser:           "all",
			errExpected:      true,
			errShouldContain: []string{"expected 'EOF', found '}'"},
		},
		{
			name:             "bad - syntax: {",
			s:                `{`,
			parser:           "all",
			errExpected:      true,
			errShouldContain: []string{"missing ',' before newline in composite literal"},
		},
		{
			name:             "bad - syntax: {}",
			s:                `{}`,
			parser:           "all",
			errExpected:      true,
			errShouldContain: []string{"bad function: unexpected type: *ast.CompositeLit"},
		},
	}

	for i, tc := range testCases {
		if tc.parser == "int64" || tc.parser == "all" {
			slc, err := int64CFParse(tc.s)
			checkParseResults(t, i, len(slc), err, tc)
		}
		if tc.parser == "float64" || tc.parser == "all" {
			slc, err := float64CFParse(tc.s)
			checkParseResults(t, i, len(slc), err, tc)
		}
		if tc.parser == "string" || tc.parser == "all" {
			slc, err := stringCFParse(tc.s)
			checkParseResults(t, i, len(slc), err, tc)
		}
		if tc.parser == "stringSlice" || tc.parser == "all" {
			slc, err := stringSliceCFParse(tc.s)
			checkParseResults(t, i, len(slc), err, tc)
		}
	}
}

// checkParseResults ...
func checkParseResults(t *testing.T, tNum int, slcLen int, err error, tc parseTestInfo) {
	tcID := fmt.Sprintf("test %d: %s", tNum, tc.name)
	if testhelper.CheckError(t, tcID, err,
		tc.errExpected, tc.errShouldContain) && err == nil {
		if slcLen != tc.lenExpected {
			t.Log(tcID)
			t.Logf("\t: expected: %d", tc.lenExpected)
			t.Logf("\t:      got: %d", slcLen)
			t.Errorf("\t: unexpected number of checks\n")
		}
	}
}
