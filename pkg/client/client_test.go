package client_test

import (
	"context"
	"fmt"
	"github.com/Miguel-Dorta/web-msg-handler/pkg"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/client"
	"golang.org/x/sys/unix"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"testing"
)

func TestPostJSON(t *testing.T) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		errs := make([]string, 0, 4)
		if method := r.Method; method != http.MethodPost {
			errs = append(errs, "Method is not POST")
		}

		if contentType := r.Header.Get(pkg.MimeContentType); contentType != pkg.MimeJSON {
			errs = append(errs, "Content-Type is not JSON")
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			errs = append(errs, "Cannot read body")
		}

		if string(body) != "{\"test\": \"hi\"}" {
			errs = append(errs, fmt.Sprintf("Unexpected body:\n-> Expected: {\"test\": \"hi\"}\n-> Found: %s", string(body)))
		}

		w.Header().Set(pkg.MimeContentType, pkg.MimeJSON)
		statusCode := 200
		resp := "{\"success\": true}"
		if len(errs) != 0 {
			statusCode = 400
			resp = fmt.Sprintf("{\"errs\": \"%v\"}", errs)
		}
		w.WriteHeader(statusCode)
		if _, err := w.Write([]byte(resp)); err != nil {
			t.Errorf("error writing response: %s", err)
		}
	})
	srv := http.Server{Addr: ":8080"}
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, unix.SIGTERM, unix.SIGINT)
	end := make(chan bool, 1)

	go func() {
		<-quit // Block until quit signal is received
		if err := srv.Shutdown(context.Background()); err != nil {
			t.Errorf("error while shutting down server: %s", err)
		}
		end <- true //Send end status
	}()

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			t.Errorf("Unexpected error which closed the server: %s", err)
		}
	}()

	resp, err := client.PostJSON("http://localhost:8080", []byte("{\"test\": \"hi\"}"))
	if err != nil {
		t.Errorf("Error when doing POST request: %s", err)
		return
	}

	if string(resp) != "{\"success\": true}" {
		t.Errorf("Returned error: %s", string(resp))
	}

	quit <- unix.SIGTERM
	<-end //Block until server ends
}
