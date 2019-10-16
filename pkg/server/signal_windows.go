// +build !unix

package server

import "os"

var quitSignals = []os.Signal{os.Interrupt}
