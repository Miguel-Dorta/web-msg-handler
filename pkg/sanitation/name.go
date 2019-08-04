package sanitation

import (
	"strings"
	"unicode"
)

// SanitizeName will remove illegal characters from the "name" field of api.Request
func SanitizeName(name string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, name)
}
