package pkg

import (
	"os"
	"os/signal"
	"sync"
)

func doStuff() {
	var err error
	senders, err = LoadConfig(configFile)
	if err != nil {
		// TODO
	}

	newRequests := make(chan msg, maxNewRequest)
	// Listen for http requests (and send it)

	var wg sync.WaitGroup
	wg.Add(1)
	go start(newRequests, &wg)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<- quit // Block until interrupt signal is received
	// Log closing
	close(newRequests)
	wg.Wait()
}
