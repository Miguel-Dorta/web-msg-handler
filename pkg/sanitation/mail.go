package sanitation

import "regexp"

var regexMail = regexp.MustCompile("^[-0-9A-Za-z!#$%&'*+/=?^_`{|}.~]+@[-0-9A-Za-z_.~]+\\.[A-Za-z]+$")

// IsValidMail checks if the mail provided is valid.
func IsValidMail(mail string) bool {
	return regexMail.MatchString(mail)
}
