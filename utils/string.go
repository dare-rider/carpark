package utils

import "strings"

func StringToBool(str string) bool {
	str = strings.ToLower(str)
	trueValues := []string{"true", "t", "yes", "y", "1"}
	return SliceContainsString(&trueValues, str)
}

// SliceContainsString returns true if a slice contains a string.
func SliceContainsString(sl *[]string, str string) bool {
	for _, n := range *sl {
		if str == n {
			return true
		}
	}
	return false
}
