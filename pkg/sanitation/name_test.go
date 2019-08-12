package sanitation_test

import (
	"github.com/Miguel-Dorta/web-msg-handler/pkg/sanitation"
	"testing"
)

func TestSanitizeName(t *testing.T) {
	result := sanitation.SanitizeName(`this will only accept printable text,
 that means, not control characters` + string(28))
	expected := "this will only accept printable text, that means, not control characters"
	if result != expected {
		t.Errorf("Invalid sanitation.\n" +
			"-> Expected output: \"%s\"\n" +
			"-> Found output: \"%s\"",
			expected, result)
	}
}
