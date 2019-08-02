package sanitation

import (
	"strings"
	"unicode"
)

func SanitizeName(name string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, name)
}
