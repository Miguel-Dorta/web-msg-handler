package recaptcha
// Package recaptcha is the package that manages the function related to the Google's ReCaptcha verification.

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/client"
	"net/url"
	"strings"
)

const recaptchaVerifyUrl = "https://www.google.com/recaptcha/api/siteverify"

// response represents the response that Google's ReCaptcha servers returns telling if the request sent+
// passes the ReCaptcha verification.
type response struct {
	Success bool `json:"success"`
	Errors []string `json:"error-codes"`
}

// CheckRecaptcha checks if the response provided have passed the ReCaptcha verification with the secret provided.
func CheckRecaptcha(secret, userResponse string) error {
	// When empty secret, omit verification
	if secret == "" {
		return nil
	}

	data := url.Values{}
	data.Set("secret", secret)
	data.Set("response", userResponse)

	rawResp, err := client.PostForm(recaptchaVerifyUrl, data)
	if err != nil {
		return fmt.Errorf("error doing request for reCaptcha verification: %s", err.Error())
	}

	var resp response
	if err = json.Unmarshal(rawResp, &resp); err != nil {
		return fmt.Errorf("error parsing reCaptcha server response: %s", err.Error())
	}

	if !resp.Success {
		sb := new(strings.Builder)
		sb.WriteString("recaptcha verification failed: ")
		if len(resp.Errors) == 0 {
			sb.WriteString("unknown reason")
		} else {
			sb.WriteString("reasons:")
			for _, e := range resp.Errors {
				sb.WriteString(fmt.Sprintf("\n -> \"%s\"", e))
			}
		}
		return errors.New(sb.String())
	}

	return nil
}
