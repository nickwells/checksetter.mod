package checksetter

import "testing"

// reportUnknownFuncErr checks that the error is not nil and has the right
// value and reports an error if not
func reportUnknownFuncErr(t *testing.T, err error, cName, fName string) {
	t.Helper()

	expErrMsg := `Unknown function: "nonesuch"`

	if err == nil || err.Error() != expErrMsg {
		t.Logf("Checker: %q, function: %q", cName, fName)

		if err == nil {
			t.Log("\t: No error was reported")
		} else {
			t.Logf("\t: error returned: %q", err)
			t.Logf("\t: error expected: %q", expErrMsg)
		}

		t.Error("\t: name should have been flagged as unknown\n")
	}
}

// initAllParsersCheckedRegister ...
func initAllParsersCheckedRegister() map[string]bool {
	pa := ParsersAvailable()
	checked := map[string]bool{}

	for _, cn := range pa {
		checked[cn] = false
	}

	return checked
}

// confirmAllParsersChecked ...
func confirmAllParsersChecked(t *testing.T, checked map[string]bool) {
	t.Helper()

	for cn, ok := range checked {
		if !ok {
			t.Fatalf("%q has no tests for unknown functions", cn)
		}
	}
}

// getParserRegisterEntry[T any] gets the named parser and confirms that the
// type is correct. It is a fatal error if not.
func getParserRegisterEntry[T any](t *testing.T, cName string) *Parser[T] {
	t.Helper()

	p, ok := parserRegister[cName].(*Parser[T])
	if !ok {
		t.Fatal("bad ParserRegister entry: ", cName)
	}

	return p
}

func TestMakerUnknownFunc(t *testing.T) {
	checked := initAllParsersCheckedRegister()

	{
		p := getParserRegisterEntry[float64](t, Float64CheckerName)

		makers := p.Makers()
		for _, fName := range makers {
			mi := p.makers[fName]
			_, err := mi.MF(nil, "nonesuch")
			reportUnknownFuncErr(t, err, p.checkerName, fName)
		}

		checked[p.checkerName] = true
	}
	{
		p := getParserRegisterEntry[int](t, IntCheckerName)

		makers := p.Makers()
		for _, fName := range makers {
			mi := p.makers[fName]
			_, err := mi.MF(nil, "nonesuch")
			reportUnknownFuncErr(t, err, p.checkerName, fName)
		}

		checked[p.checkerName] = true
	}
	{
		p := getParserRegisterEntry[int64](t, Int64CheckerName)

		makers := p.Makers()
		for _, fName := range makers {
			mi := p.makers[fName]
			_, err := mi.MF(nil, "nonesuch")
			reportUnknownFuncErr(t, err, p.checkerName, fName)
		}

		checked[p.checkerName] = true
	}
	{
		p := getParserRegisterEntry[string](t, StringCheckerName)

		makers := p.Makers()
		for _, fName := range makers {
			mi := p.makers[fName]
			_, err := mi.MF(nil, "nonesuch")
			reportUnknownFuncErr(t, err, p.checkerName, fName)
		}

		checked[p.checkerName] = true
	}
	{
		p := getParserRegisterEntry[[]string](t, StringSliceCheckerName)

		makers := p.Makers()
		for _, fName := range makers {
			mi := p.makers[fName]
			_, err := mi.MF(nil, "nonesuch")
			reportUnknownFuncErr(t, err, p.checkerName, fName)
		}

		checked[p.checkerName] = true
	}

	confirmAllParsersChecked(t, checked)
}
