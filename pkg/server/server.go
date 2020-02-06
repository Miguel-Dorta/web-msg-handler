package server

// Package server will manage all the HTTP request made to web-msg-handler.

import (
	"context"
	"encoding/json"
	"github.com/Miguel-Dorta/logolang"
	"github.com/Miguel-Dorta/web-msg-handler/api"
	"github.com/Miguel-Dorta/web-msg-handler/pkg"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/config"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/plugin"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/recaptcha"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/sanitation"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

var (
	log   *logolang.Logger
	sites map[string]*config.Site
)

func init() {
	log = logolang.NewLogger()
	log.Color = false
}

// Run will start a HTTP server in the port provided using the config file path provided.
// It ends when a termination or interrupt signal is received.
// It can end the program execution prematurely.
func Run() {
	var (
		conf *config.Config
		err error
	)
	conf, sites, err = config.LoadConfig()
	if err != nil {
		log.Criticalf("error loading config: %s", err)
		os.Exit(1)
	}

	http.HandleFunc("/", handle)
	srv := http.Server{Addr: ":" + strconv.Itoa(conf.Port)}

	serverClosed := make(chan bool)
	go func() {
		defer close(serverClosed)

		quit := make(chan os.Signal, 2)
		signal.Notify(quit, quitSignals...)
		<-quit // Block until quit signal is received

		log.Info("Shutting down")

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Criticalf("error while shutting down: %s", err)
			os.Exit(1)
		}
	}()

	log.Infof("Listening port %s", srv.Addr[1:])
	if err = srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Criticalf("Unexpected error which closed the server: %s", err)
		os.Exit(1)
	}
	<-serverClosed
}

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

	site, ok := sites[siteID]
	if !ok {
		log.Debugf("[Request %d] Site ID not found: %d", requestID, siteID)
		statusWriter(w, ErrNotFound)
		return
	}

	if method := r.Method; method != http.MethodPost {
		log.Debugf("[Request %d] Invalid method: %s", requestID, method)
		statusWriter(w, ErrMethodNotAllowed)
		return
	}

	if contentType := r.Header.Get(pkg.MimeContentType); contentType != pkg.MimeJSON {
		log.Debugf("[Request %d] Invalid content type: %s", requestID, contentType)
		statusWriter(w, ErrContentTypeNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("[Request %d] Error while reading body: %s", requestID, err)
		statusWriter(w, ErrReadingBody)
		return
	}

	var r2 api.Request
	if err = json.Unmarshal(body, &r2); err != nil {
		log.Debugf("[Request %d] Malformed JSON: %s", requestID, err)
		statusWriter(w, ErrMalformedJSON)
		return
	}

	if !sanitation.IsValidMail(r2.Mail) {
		log.Debugf("[Request %d] Invalid email", requestID)
		statusWriter(w, ErrInvalidMail)
		return
	}

	msgJS, err := plugin.MsgToJS(sanitation.SanitizeName(r2.Name), r2.Mail, sanitation.SanitizeMsg(r2.Msg))
	if err != nil {
		log.Errorf("[Request %d] Error parsing msg to JS: %s", requestID, err)
		statusWriter(w, ErrInternalServerError)
		return
	}

	if err = recaptcha.CheckRecaptcha(site.RecaptchaSecret, r2.Recaptcha); err != nil {
		log.Debugf("[Request %d] Recaptcha verification failed: %s", requestID, err)
		statusWriter(w, ErrRecaptchaVerificationFailed)
		return
	}

	if err = plugin.Exec(site.SenderType, site.ConfigJS, msgJS); err != nil {
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
	w.Header().Set(pkg.MimeContentType, pkg.MimeJSON)
	w.WriteHeader(resp.status)

	data, _ := json.Marshal(api.Response{
		Success: resp.success,
		Err:     resp.msg,
	})

	if _, err := w.Write(data); err != nil {
		log.Errorf("error writing response: %s", err)
	}
}
