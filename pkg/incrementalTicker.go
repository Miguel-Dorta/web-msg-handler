package pkg

import "time"

type incrementalTicker struct {
	C    chan int
	quit chan bool
}

func makeIncrementalTicker(d time.Duration) *incrementalTicker {
	it := incrementalTicker{
		C:    make(chan int, 1),
		quit: make(chan bool, 1),
	}
	go it.run(d)
	return &it
}

func (it *incrementalTicker) stop() {
	it.quit <- true
}

func (it *incrementalTicker) run(d time.Duration) {
	multiplier := 1
	for i := 0; ; i++ {
		select {
		case <-it.quit:
			return
		default:
			it.C <- i
			time.Sleep(d * time.Duration(multiplier))
			multiplier *= 2
		}
	}
}
