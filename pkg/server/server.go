package server

// Package server will manage all the HTTP request made to web-msg-handler.

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Miguel-Dorta/logolang"
	"github.com/Miguel-Dorta/web-msg-handler/api"
	"github.com/Miguel-Dorta/web-msg-handler/pkg"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/config"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/sanitation"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/sender"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"
)

const statusUnknownError = 502

var (
	Log        *logolang.Logger
	sites      map[uint64]sender.Sender

	closing    = false
	requestsWG = &sync.WaitGroup{}
)

// Run will start a HTTP server in the port provided using the config file path provided.
// It ends when a termination or interrupt signal is received.
// It can end the program execution prematurely.
func Run(configFile, port string) {
	var err error
	sites, err = config.LoadConfig(configFile)
	if err != nil {
		Log.Criticalf("error loading config file: %s", err)
		os.Exit(1)
	}

	http.HandleFunc("/", handle)
	srv := http.Server{Addr: ":" + port}

	go func() {
		quit := make(chan os.Signal, 2)
		signal.Notify(quit, quitSignals...)
		<-quit // Block until quit signal is received

		Log.Info("Shutting down")

		closing = true
		requestsWG.Wait()

		if err := srv.Shutdown(context.Background()); err != nil {
			Log.Criticalf("error while shutting down: %s", err)
			os.Exit(1)
		}
	}()

	Log.Infof("Listening port %s", srv.Addr[1:])
	if err = srv.ListenAndServe(); err != http.ErrServerClosed {
		Log.Criticalf("Unexpected error which closed the server: %s", err)
		os.Exit(1)
	}
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
	Log.Debugf("[Request %d] Received: %+v", requestID, r)

	// Close request if the closing var is set to true
	if closing {
		Log.Debugf("[Request %d] Reject request. Closing server.", requestID)
		statusWriter(w, http.StatusServiceUnavailable, false, "closing server")
		return
	}
	requestsWG.Add(1) // This is after the closing check because if it's before it could never stop
	defer requestsWG.Done()

	url := r.URL
	id, err := strconv.ParseUint(url.Path[1:], 10, 64)
	if err != nil {
		Log.Debugf("[Request %d] Failed to parse ID: %s", requestID, url.Path[1:])
		statusWriter(w, http.StatusNotFound, false, fmt.Sprintf("path %s not found", url))
		return
	}

	s, senderExists := sites[id]
	if !senderExists {
		Log.Debugf("[Request %d] ID not found: %d", requestID, id)
		statusWriter(w, http.StatusNotFound, false, fmt.Sprintf("path %s not found", url))
		return
	}

	if method := r.Method; method != http.MethodPost {
		Log.Debugf("[Request %d] Invalid method: %s", requestID, method)
		statusWriter(w, http.StatusMethodNotAllowed, false, fmt.Sprintf("method %s not supported", method))
		return
	}

	if contentType := r.Header.Get(pkg.MimeContentType); contentType != pkg.MimeJSON {
		Log.Debugf("[Request %d] Invalid content type: %s", requestID, contentType)
		statusWriter(w, http.StatusBadRequest, false, fmt.Sprintf("content-type %s not supported", contentType))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Log.Errorf("[Request %d] Error while reading body: %s", requestID, err)
		statusWriter(w, statusUnknownError, false, fmt.Sprintf("unknown error while reading request body: %s", err.Error()))
		return
	}

	var r2 api.Request
	if err = json.Unmarshal(body, &r2); err != nil {
		Log.Debugf("[Request %d] Malformed JSON: %s", requestID, err)
		statusWriter(w, http.StatusBadRequest, false, "malformed JSON")
		return
	}

	if !sanitation.IsValidMail(r2.Mail) {
		Log.Debugf("[Request %d] Invalid email", requestID)
		statusWriter(w, http.StatusBadRequest, false, "invalid email")
		return
	}

	if err = s.CheckRecaptcha(r2.Recaptcha); err != nil {
		Log.Debugf("[Request %d] Recaptcha verification failed: %s", requestID, err)
		statusWriter(w, http.StatusBadRequest, false, "recaptcha verification failed")
		return
	}

	if err = s.Send(sanitation.SanitizeName(r2.Name), r2.Mail, sanitation.SanitizeMsg(r2.Msg)); err != nil {
		Log.Debugf("[Request %d] Sender failed: %s", requestID, err)
		statusWriter(w, http.StatusServiceUnavailable, false, "error sending message")
		return
	}

	statusWriter(w, http.StatusOK, true, "")
	Log.Debugf("[Request %d] Success", requestID)
}

// statusWriter will write a response to the http.ResponseWriter provided.
// That response will be sent with the status code provided,
// and its body will consists in a JSON represented by api.Response with the success status and error provided.
func statusWriter(w http.ResponseWriter, statusCode int, success bool, msg string) {
	w.Header().Set(pkg.MimeContentType, pkg.MimeJSON)
	w.WriteHeader(statusCode)

	data, _ := json.Marshal(api.Response{
		Success: success,
		Err:     msg,
	})

	if _, err := w.Write(data); err != nil {
		Log.Errorf("error writing response: %s", err)
	}
}
