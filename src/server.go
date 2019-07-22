package src

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
)

const (
	mimeContentType    = "Content-Type"
	mimeJson           = "application/json"
	statusUnknownError = 502
)

var (
	configFile = ""
	port = "8080"
	senders map[uint64]sender
)

func start() {
	if err := LoadConfig(configFile); err != nil {
		// TODO
	}

	http.HandleFunc("/", handle)

	srv := http.Server{Addr: ":" + port}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// TODO
		}
	}()

	quit := make(chan os.Signal, 2)
	signal.Notify(quit, unix.SIGTERM, unix.SIGINT)
	<-quit // Block until quit signal is received

	// TODO log server closing
	if err := srv.Shutdown(context.Background()); err != nil {
		// TODO
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.URL.Path[1:], 10, 64)
	if err != nil {
		statusWriter(w, http.StatusNotFound, false, fmt.Sprintf("path %s not found", r.URL))
		return
	}

	s, senderExists := senders[id]
	if !senderExists {
		statusWriter(w, http.StatusNotFound, false, fmt.Sprintf("path %s not found", r.URL))
		return
	}

	if method := r.Method; method != http.MethodPost {
		statusWriter(w, http.StatusMethodNotAllowed, false, fmt.Sprintf("method %s not supported", method))
		return
	}

	if contentType := r.Header.Get(mimeContentType); contentType != mimeJson {
		statusWriter(w, http.StatusBadRequest, false, fmt.Sprintf("content-type %s not supported", contentType))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		statusWriter(w, statusUnknownError, false, fmt.Sprintf("unknown error while reading request body: %s", err.Error()))
		return
	}

	var r2 request
	if err = json.Unmarshal(body, &r2); err != nil {
		statusWriter(w, http.StatusBadRequest, false, "malformed JSON")
		return
	}

	if !s.Send(r2.Name, r2.Mail, r2.Msg) {
		statusWriter(w, http.StatusServiceUnavailable, false, "error sending message")
		return
	}

	statusWriter(w, http.StatusOK, true, "")
}

func statusWriter(w http.ResponseWriter, statusCode int, success bool, msg string) {
	w.Header().Set(mimeContentType, mimeJson)
	w.WriteHeader(statusCode)

	data, _ := json.Marshal(response{
		Success: success,
		Err:     msg,
	})

	if _, err := w.Write(data); err != nil {
		// Log
		return fmt.Errorf("error writing response: %s", err.Error())
	}
}
