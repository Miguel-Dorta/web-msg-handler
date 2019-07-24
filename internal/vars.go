package internal

import (
	"github.com/Miguel-Dorta/logolang"
)

var (
	Log *logolang.Logger
	senders map[uint64]sender
	Version string
)
