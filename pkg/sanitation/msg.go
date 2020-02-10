package sanitation

import (
	"strings"
	"unicode"
)

// SanitizeMsg will remove illegal characters from the "msg" field of api.Request
func SanitizeMsg(msg string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) || r == '\n' {
			return r
		}
		return -1
	}, msg)
}
