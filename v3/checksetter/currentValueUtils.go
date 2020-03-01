package checksetter

import "fmt"

// currentValue returns a standard message for the current value depending on
// the number of checks in the list (passed as checkCount)
func currentValue(checkCount int) string {
	switch checkCount {
	case 0:
		return "no checks"
	case 1:
		return "one check"
	}

	return fmt.Sprintf("%d checks", checkCount)
}
