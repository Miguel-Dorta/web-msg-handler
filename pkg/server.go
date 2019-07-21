package pkg

import (
	"context"
	"encoding/json"
	"golang.org/x/sys/unix"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
)

type clientRequest struct {
	Name string `json:"name"`
	Mail string `json:"mail"`
	Msg string `json:"msg"`
}

func doStuff() {
	var err error
	senders, err = LoadConfig(configFile)
	if err != nil {
		// TODO
	}

	var wg sync.WaitGroup

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.URL.Path[1:], 10, 64)
		if err != nil {
			// TODO 404
		}

		if _, exists := senders[id]; !exists {
			// TODO 404
		}

		if r.Method != http.MethodPost {
			// TODO 405
		}

		if r.Header.Get("Content-Type") != "application/json" {
			// TODO 400
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			// TODO
		}

		var cr clientRequest
		if err = json.Unmarshal(body, &cr); err != nil {
			// TODO 400
		}

		wg.Add(1)
		msg{
			ReceiverId: id,
			SenderName: cr.Name,
			SenderMail: cr.Mail,
			Message: cr.Msg,
		}.process(&wg)
	})

	srv := http.Server{Addr: ":8080"}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// TODO
		}
	}()

	quit := make(chan os.Signal, 2)
	signal.Notify(quit, unix.SIGTERM, unix.SIGINT)
	<- quit // Block until quit signal is received
	// TODO log server closing
	err = srv.Shutdown(context.Background())
	if err != nil {
		// TODO
	}
	wg.Wait()
}
