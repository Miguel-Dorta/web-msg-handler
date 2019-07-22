package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/sys/unix"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

const (
	mimeContentType    = "Content-Type"
	mimeJson           = "application/json"
	statusUnknownError = 502
)

func Run(configFile, port string) {
	if err := loadConfig(configFile); err != nil {
		Log.Criticalf("error loading config file: %s", err)
		os.Exit(1)
	}

	http.HandleFunc("/", handle)
	srv := http.Server{Addr: ":" + port}

	go func() {
		quit := make(chan os.Signal, 2)
		signal.Notify(quit, unix.SIGTERM, unix.SIGINT)
		<-quit // Block until quit signal is received

		Log.Info("Shutting down")
		if err := srv.Shutdown(context.Background()); err != nil {
			Log.Criticalf("error while shutting down: %s", err)
			os.Exit(1)
		}
	}()

	Log.Infof("Listening port %s", srv.Addr[1:])
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		Log.Criticalf("Unexpected error which closed the server: %s", err)
		os.Exit(1)
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	// Request ID for logging purposes
	requestId := time.Now().UnixNano()
	Log.Debugf("[Request %d] Received: %+v", requestId, r)

	url := r.URL
	id, err := strconv.ParseUint(url.Path[1:], 10, 64)
	if err != nil {
		Log.Debugf("[Request %d] Failed to parse ID: %s", requestId, url.Path[1:])
		statusWriter(w, http.StatusNotFound, false, fmt.Sprintf("path %s not found", url))
		return
	}

	s, senderExists := senders[id]
	if !senderExists {
		Log.Debugf("[Request %d] ID not found: %d", requestId, id)
		statusWriter(w, http.StatusNotFound, false, fmt.Sprintf("path %s not found", url))
		return
	}

	if method := r.Method; method != http.MethodPost {
		Log.Debugf("[Request %d] Invalid method: %s", requestId, method)
		statusWriter(w, http.StatusMethodNotAllowed, false, fmt.Sprintf("method %s not supported", method))
		return
	}

	if contentType := r.Header.Get(mimeContentType); contentType != mimeJson {
		Log.Debugf("[Request %d] Invalid content type: %s", requestId, contentType)
		statusWriter(w, http.StatusBadRequest, false, fmt.Sprintf("content-type %s not supported", contentType))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Log.Errorf("[Request %d] Error while reading body: %s", requestId, err)
		statusWriter(w, statusUnknownError, false, fmt.Sprintf("unknown error while reading request body: %s", err.Error()))
		return
	}

	var r2 request
	if err = json.Unmarshal(body, &r2); err != nil {
		Log.Debugf("[Request %d] Malformed JSON: %s", requestId, err)
		statusWriter(w, http.StatusBadRequest, false, "malformed JSON")
		return
	}

	if !s.Send(r2.Name, r2.Mail, r2.Msg) {
		Log.Debugf("[Request %d] Sender failed", requestId)
		statusWriter(w, http.StatusServiceUnavailable, false, "error sending message")
		return
	}

	statusWriter(w, http.StatusOK, true, "")
	Log.Debugf("[Request %d] Success", requestId)
}

func statusWriter(w http.ResponseWriter, statusCode int, success bool, msg string) {
	w.Header().Set(mimeContentType, mimeJson)
	w.WriteHeader(statusCode)

	data, _ := json.Marshal(response{
		Success: success,
		Err:     msg,
	})

	if _, err := w.Write(data); err != nil {
		Log.Errorf("error writing response: %s", err)
	}
}
