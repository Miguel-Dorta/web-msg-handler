package recaptcha

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/client"
)

const recaptchaVerifyUrl = "https://www.google.com/recaptcha/api/siteverify"

type request struct {
	Secret string `json:"secret"`
	Response string `json:"response"`
}

type response struct {
	Success bool `json:"success"`
	Errors []string `json:"error-codes"`
}

func CheckRecaptcha(secret, userResponse string) error {
	data, err := json.Marshal(request{
		Secret:   secret,
		Response: userResponse,
	})
	if err != nil {
		return fmt.Errorf("error parsing recaptcha request JSON: %s", err.Error())
	}

	rawResp, err := client.PostJSON(recaptchaVerifyUrl, data)
	if err != nil {
		return fmt.Errorf("error doing request for reCaptcha verification: %s", err.Error())
	}

	var resp response
	if err = json.Unmarshal(rawResp, &resp); err != nil {
		return fmt.Errorf("error parsing reCaptcha server response: %s", err.Error())
	}

	if !resp.Success {
		errStr := "recaptcha verification failed"
		if len(resp.Errors) != 0 {
			errStr += ":"
			for _, e := range resp.Errors {
				errStr += " \"" + e + "\""
			}
		}
		return errors.New(errStr)
	}

	return nil
}