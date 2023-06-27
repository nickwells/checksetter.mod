package checksetter_test

import (
	"fmt"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/checksetter.mod/v4/checksetter"
	"github.com/nickwells/param.mod/v5/param/paramset"
)

// ExampleSetter demonstrates how the Setter should be used with the param
// package
func ExampleSetter() {
	chkFuncs := []check.ValCk[string]{}
	ps := paramset.NewOrPanic()

	ps.Add("checks",
		&checksetter.Setter[string]{
			Value: &chkFuncs,
			Parser: checksetter.FindParserOrPanic[string](
				checksetter.StringCheckerName),
		},
		"help-text")

	ps.Parse([]string{"-checks", "OK, Length(GT(1))"})

	fmt.Printf("%d checks provided\n", len(chkFuncs))
	// Output:
	// 2 checks provided
}
