package sanitation_test

import (
	"github.com/Miguel-Dorta/web-msg-handler/pkg/sanitation"
	"testing"
)

func TestSanitizeMsg(t *testing.T) {
	expected := `this is a literal text
it has control and printable characters
except the characters following #99:`
	result := sanitation.SanitizeMsg(expected + string(28))
	if result != expected {
		t.Errorf("Invalid sanitation.\n" +
			"-> Expected output: \"%s\"\n" +
			"-> Found output: \"%s\"",
			expected, result)
	}
}
