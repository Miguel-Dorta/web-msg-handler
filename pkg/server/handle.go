package server

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Miguel-Dorta/web-msg-handler/api"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/mime"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/plugin"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/recaptcha"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/sanitation"
	"io/ioutil"
	"net/http"
	"time"
)

// handle is the function executed for each HTTP request received by web-msg-handler.
//
// It will:
//
// - Assign an ID to every request (corresponding to a timestamp of the EPOCH nanosecond when it was received
// for debugging and logging purposes.
//
// - Check if the Sender ID is correct
//
// - Check if the HTTP method used is POST
//
// - Check if the Content-Type header is the MIME JSON.
//
// - Check if the request body is valid.
//
// - Check if the email provided is valid.
//
// - Check if the request have passed the ReCaptcha verification.
//
// - Send the message
func handle(w http.ResponseWriter, r *http.Request) {
	// Request ID for logging purposes
	requestID := time.Now().UnixNano()
	log.Debugf("[Request %d] Received: %+v", requestID, r)

	siteID := r.URL.Path[1:]

	// Check if site exists
	site, ok := sites[siteID]
	if !ok {
		log.Debugf("[Request %d] Site ID not found: %d", requestID, siteID)
		statusWriter(w, ErrNotFound)
		return
	}

	// Check if method is valid
	if method := r.Method; method != http.MethodPost {
		log.Debugf("[Request %d] Invalid method: %s", requestID, method)
		statusWriter(w, ErrMethodNotAllowed)
		return
	}

	// Check if content-type is valid
	if contentType := r.Header.Get(mime.ContentType); contentType != mime.JSON {
		log.Debugf("[Request %d] Invalid content type: %s", requestID, contentType)
		statusWriter(w, ErrContentTypeNotAllowed)
		return
	}

	// Read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("[Request %d] Error while reading body: %s", requestID, err)
		statusWriter(w, ErrReadingBody)
		return
	}

	// Parse body
	var r2 api.Request
	if err = json.Unmarshal(body, &r2); err != nil {
		log.Debugf("[Request %d] Malformed JSON: %s", requestID, err)
		statusWriter(w, ErrMalformedJSON)
		return
	}

	// Check valid email
	if !sanitation.IsValidMail(r2.Mail) {
		log.Debugf("[Request %d] Invalid email", requestID)
		statusWriter(w, ErrInvalidMail)
		return
	}

	// Sanitate and serialize input
	msgJS, err := msgToJSON(sanitation.SanitizeName(r2.Name), r2.Mail, sanitation.SanitizeMsg(r2.Msg))
	if err != nil {
		log.Errorf("[Request %d] Error parsing msg to JSON: %s", requestID, err)
		statusWriter(w, ErrUnknown)
		return
	}

	// Check recaptcha
	if err = recaptcha.CheckRecaptcha(site.RecaptchaSecret, r2.Recaptcha); err != nil {
		log.Debugf("[Request %d] Recaptcha verification failed: %s", requestID, err)
		statusWriter(w, ErrRecaptchaVerificationFailed)
		return
	}

	// Exec plugin
	if err = plugin.Exec(site.SenderName, site.ConfigJSON, msgJS); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Errorf("[Request %d] Sender %s took too long", requestID, site.SenderName)
			statusWriter(w, ErrGatewayTimeout)
			return
		}

		log.Errorf("[Request %d] Sender failed: %s", requestID, err)
		statusWriter(w, ErrInternalServerError)
		return
	}

	statusWriter(w, ResponseOK)
	log.Debugf("[Request %d] Success", requestID)
}

// statusWriter will write a response to the http.ResponseWriter provided.
// That response will be sent with the status code provided,
// and its body will consists in a JSON represented by api.Response with the success status and error provided.
func statusWriter(w http.ResponseWriter, resp *httpResponse) {
	w.Header().Set(mime.ContentType, mime.JSON)
	w.WriteHeader(resp.status)

	data, _ := json.Marshal(api.Response{
		Success: resp.success,
		Err:     resp.msg,
	})

	if _, err := w.Write(data); err != nil {
		log.Errorf("error writing response: %s", err)
	}
}

// msgToJSON takes the name, mail and msg provided and creates the JSON that will be passed to the site
func msgToJSON(name, mail, msg string) (string, error) {
	data, err := json.Marshal(map[string]string{
		"name": name,
		"mail": mail,
		"msg": msg,
	})
	if err != nil {
		return "", err
	}
	return string(data), nil
}
