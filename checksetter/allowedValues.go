package checksetter

import (
	"sort"
	"strings"
)

// getCheckFuncDesc returns a string describing the function and its
// arguments
func getCheckFuncDesc(fName string, args []string) string {
	return fName + "(" + strings.Join(args, ", ") + ")"
}

// getOrderedNames returns the function names in alphabetical order
func getOrderedNames(nameToArgs map[string][]string) []string {
	names := make([]string, 0, len(nameToArgs))
	for k := range nameToArgs {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// allowedValFuncs will return a string showing all the allowed values for the
// given family of check functions. It will also show the allowed values for
// any referenced families of check functions.
func allowedValFuncs(checkerName string, makerFuncs map[string][]string) string {
	type toShowDetails struct {
		shown      bool
		makerFuncs map[string][]string
	}
	toShow := map[string]toShowDetails{
		checkerName: {makerFuncs: makerFuncs},
	}

	allowedVals := make([]string, 0)
	hasNew := true
	for hasNew {
		hasNew = false
		for k, v := range toShow {
			if !v.shown {
				v.shown = true
				toShow[k] = v

				names := getOrderedNames(v.makerFuncs)

				funcSet := make([]string, 0, len(names)+1)
				funcSet = append(funcSet, "the allowed "+k+" functions are:")

				for _, fn := range names {
					funcSet = append(funcSet,
						"    "+getCheckFuncDesc(fn, v.makerFuncs[fn]))
					for _, arg := range v.makerFuncs[fn] {
						if _, ok := toShow[arg]; !ok {
							if p, ok := parserRegister[arg]; ok {
								hasNew = true
								toShow[arg] = toShowDetails{
									makerFuncs: p.MakerFuncs(),
								}
							} else {
								toShow[arg] = toShowDetails{
									shown: true,
								}
							}
						}
					}
				}
				allowedVals = append(allowedVals, strings.Join(funcSet, "\n"))
			}
		}
	}

	return strings.Join(allowedVals, "\n\n")
}

// AllowedValues returns a string descibing the allowed values for the given
// class of Check functions
func AllowedValues(checkerName string, makerFuncs map[string][]string) string {
	rval := "a list of " + checkerName + " functions separated by ','.\n"
	rval += `
Write the checks as if you were writing code.

The functions recognised are:
` + allowedValFuncs(checkerName, makerFuncs)

	return rval
}
