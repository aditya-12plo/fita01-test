package utils

import "github.com/gobeam/stringy"

func KebabStyle(parameter string, typeCase string) string {

	str := stringy.New(parameter)
	result := str.KebabCase()

	// Lorem-Ipsum-Dolor-Sit-Amet
	if typeCase == "lower" {
		return result.ToLower()
	} else {
		return result.ToUpper()
	}
}
