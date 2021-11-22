package util

import "strings"

// contains performs a case-insensitive search for the item in the input array
func Contains(arr []string, item string) bool {
	if len(arr) == 0 {
		return false
	}
	for _, str := range arr {
		if strings.EqualFold(str, item) {
			return true
		}
	}
	return false
}
