package pkg

import (
	"golang.org/x/sys/unix"
	"os"
	"os/signal"
	"sync"
	"time"
)

type msg struct {
	ReceiverId uint64 `json:"id"`
	SenderName string `json:"name"`
	SenderMail string `json:"mail"`
	Message    string `json:"msg"`
}

func (m *msg) process(wg *sync.WaitGroup) {
	quit := make(chan os.Signal)
	signal.Notify(quit, unix.SIGINT, unix.SIGTERM)
	it := makeIncrementalTicker(time.Second * 10)

	for exitLoop := false; !exitLoop; {
		select {
		case try := <-it.C:
			if try > maxTries {
				exitLoop = true
				break
			}

			if senders[m.ReceiverId].Send(*m) {
				exitLoop = true
			}
		case <-quit:
			exitLoop = true
			// TODO DUMP
		}
	}
	wg.Done()
}
