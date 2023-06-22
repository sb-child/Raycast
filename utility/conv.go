package utility

import (
	"regexp"
	"strconv"
)

func ExtractNumber(s string) (int, error) {
	re := regexp.MustCompile(`\d+`)
	match := re.FindString(s)
	return strconv.Atoi(match)
}

func RemoveElement[T comparable](slice []T, x T) []T {
	for i, v := range slice {
		if v == x {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
