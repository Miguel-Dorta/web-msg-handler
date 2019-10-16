// +build !windows

package server

import (
	"golang.org/x/sys/unix"
	"os"
)

var quitSignals = []os.Signal{unix.SIGTERM, unix.SIGINT}
