package utils

import "strconv"

func StringToUint(element string) (uint, bool) {
	var uintElement uint

	element64, err := strconv.ParseUint(element, 10, 0)
	if err != nil {
		return uintElement, false
	}

	return uint(element64), true
}
