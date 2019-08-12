package sanitation_test

import (
	"github.com/Miguel-Dorta/web-msg-handler/pkg/sanitation"
	"testing"
)

func TestIsValidMail(t *testing.T) {
	mails := []struct{
		addr string
		isValid bool
	}{
		{"me@me.me", true},
		{"eMaIl1234._test@domain.tld", true},
		{"me@me@me.com", false},
		{"hello@world.", false},
		{"@domain.com", false},
		{"domain.com", false},
		{"-123TEXTtext!#$%&'*+/=?^_`{|}~.@-123TEXTtext_.~.TEXTtext", true},
		{"\\@text.com", false},
		{"subdomain@test.subdomain1.subdomain2.com", true},
		{"slash/test@test.com", true},
		{"notld@test", false},
		{"invalidChars@?*'.123", false},
	}

	for _, mail := range mails {
		result := sanitation.IsValidMail(mail.addr)
		if result != mail.isValid {
			t.Errorf("Incorrect validation:\n" +
				"-> Addr: \"%s\"\n" +
				"-> Expected: %v\n" +
				"-> Found: %v",
				mail.addr, mail.isValid, result)
		}
	}
}
