package pkg

import "sync"

func start(newRequests chan msg, wg *sync.WaitGroup) {
	pending, err := LoadPending("/home")
	if err != nil {
		//TODO
	}

	mainLoop:
	for {
		processAllRequest(newRequests, pending)

		select {
		case r, ok := <- newRequests:
			if !ok {
				break mainLoop
			}
			if success := r.process(); !success {
				pending <- r
			}
		case r := <- pending:
			if success := r.process(); !success {
				pending <- r
			}
		}
	}

	if err = SavePending("", pending); err != nil {
		// TODO
	}

	wg.Done()
}

func processAllRequest(newRequests, pending chan msg) {
	for len(newRequests) > 0 {
		r := <- newRequests
		if success := r.process(); !success {
			pending <- r
		}
	}
}
