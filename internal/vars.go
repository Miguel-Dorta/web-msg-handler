package internal

import (
	"github.com/Miguel-Dorta/logolang"
	"net/http"
	"time"
)

var (
	httpClient = http.Client{Timeout: 10 * time.Second}
	Log *logolang.Logger
	senders map[uint64]sender
	Version string
)
