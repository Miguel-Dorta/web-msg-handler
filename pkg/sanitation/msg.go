package sanitation

import (
	"strings"
	"unicode"
)

func SanitizeMsg(msg string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) || unicode.IsControl(r) {
			return r
		}
		return -1
	}, msg)
}
