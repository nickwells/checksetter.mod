package checksetter

import (
	"sort"
	"strings"
)

var mapOfNameToArgMaps = map[string]map[string]string{
	strSlcCFName: {
		"NoDups":     "",
		"LenEQ":      "int",
		"LenGT":      "int",
		"LenLT":      "int",
		"LenBetween": "int, int",
		"String":     strCFName,
		"Contains":   strCFName + ", string",
		"Not":        strSlcCFName + ", string",
		"And":        strSlcCFName + " ...",
		"Or":         strSlcCFName + " ...",
	},
	strCFName: {
		"LenEQ":          "int",
		"LenGT":          "int",
		"LenLT":          "int",
		"LenBetween":     "int, int",
		"Equals":         "string",
		"HasPrefix":      "string",
		"HasSuffix":      "string",
		"MatchesPattern": "regexp, string",
		"Not":            strCFName + ", string",
		"And":            strCFName + " ...",
		"Or":             strCFName + " ...",
	},
	int64CFName: {
		"EQ":          "int",
		"GT":          "int",
		"GE":          "int",
		"LT":          "int",
		"LE":          "int",
		"Divides":     "int",
		"IsAMultiple": "int",
		"Between":     "int, int",
		"Not":         int64CFName + ", string",
		"And":         int64CFName + " ...",
		"Or":          int64CFName + " ...",
	},
	float64CFName: {
		"GT":      "float",
		"GE":      "float",
		"LT":      "float",
		"LE":      "float",
		"Between": "float, float",
		"Not":     float64CFName + ", string",
		"And":     float64CFName + " ...",
		"Or":      float64CFName + " ...",
	},
}

// trimArg removes any white space and '...' from the argument string and
// returns what's left. This gives you a string that can be tested to see if
// it's the name of a set of allowed values
func trimArg(a string) string {
	a = strings.TrimLeft(a, " ")
	a = strings.TrimRight(a, " ")
	a = strings.TrimSuffix(a, "...")
	a = strings.TrimRight(a, " ")
	return a
}

// hasNewArgType will add any of the arguments to the function which are
// themselves families of check functions to the map of work remaining
// (provided they haven't already been reported). It returns true if any new
// work has been added.
func hasNewArgType(args string, shown map[string]bool) bool {
	hasNew := false
	for _, a := range strings.Split(args, ",") {
		a = trimArg(a)
		if _, ok := mapOfNameToArgMaps[a]; ok {
			if !shown[a] {
				shown[a] = false
				hasNew = true
			}
		}
	}
	return hasNew
}

// allowedValFuncs will return a string showing all the allowed values for the
// given family of check functions. It will also show the allowed values for
// any referenced families of check functions.
func allowedValFuncs(s string) string {
	shown := map[string]bool{
		s: false,
	}
	hasNew := true
	rval := ""
	for hasNew {
		hasNew = false
		for k, isShown := range shown {
			if !isShown {
				shown[k] = true
				nameToArgs := mapOfNameToArgMaps[k]
				rval += "\nfor " + k + " allowed values are:\n"
				for _, fn := range getOrderedNames(nameToArgs) {
					rval += "    " + fn
					if args := nameToArgs[fn]; args != "" {
						rval += "(" + args + ")"
						if hasNewArgType(args, shown) {
							hasNew = true
						}
					}
					rval += "\n"
				}
				rval += "\n"
			}
		}
	}

	return rval
}

// getOrderedNames returns the function names in alphabetical order
func getOrderedNames(nameToArgs map[string]string) []string {
	names := make([]string, 0, len(nameToArgs))
	for k := range nameToArgs {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// allowedValues returns a string descibing the allowed values for the given
// class of Check functions
func allowedValues(cfName string) string {
	rval := "a list of " + cfName + " functions separated by ','.\n"
	rval += `
Write the checks as if you were writing code.

The functions recognised are:` + allowedValFuncs(cfName)

	return rval
}
