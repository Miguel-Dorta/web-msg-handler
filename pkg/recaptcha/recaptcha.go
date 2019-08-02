package recaptcha

import (
	"encoding/json"
	"errors"
	"fmt"
)

const recaptchaVerifyUrl = "https://www.google.com/recaptcha/api/siteverify"

type recaptchaRequest struct {
	Secret string `json:"secret"`
	Response string `json:"response"`
}

type recaptchaResponse struct {
	Success bool `json:"success"`
	Errors []string `json:"error-codes"`
}

func checkRecaptcha(secret, response string) error {
	data, err := json.Marshal(recaptchaRequest{
		Secret: secret,
		Response: response,
	})
	if err != nil {
		return fmt.Errorf("error parsing recaptcha request JSON: %s", err.Error())
	}

	rawResp, err := postJson(recaptchaVerifyUrl, data)
	if err != nil {
		return fmt.Errorf("error doing request for reCaptcha verification: %s", err.Error())
	}

	var resp recaptchaResponse
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
