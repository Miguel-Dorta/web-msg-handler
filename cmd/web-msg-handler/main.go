package main

import (
	"flag"
	"fmt"
	"github.com/Miguel-Dorta/logolang"
	"github.com/Miguel-Dorta/web-msg-handler/internal"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/server"
	"os"
	"strconv"
)

var (
	configPath   string
	logFile      string
	port         int
	verboseLevel int
	version      bool
)

func init() {
	flag.StringVar(&configPath, "config", "config.json", "set config path")
	flag.StringVar(&logFile, "log-file", "", "set log file")
	flag.IntVar(&port, "port", 8080, "set port")
	flag.IntVar(&verboseLevel, "verbose", 3, "set verbose level. 0=no-log, 1=critical, 2=errors, 3=info, 4=debug")
	flag.BoolVar(&version, "version", false, "print version and exit")

	flag.Parse()
}

// checkFlags checks if the values assigned to the parsed flags are valid.
// It will also print version and exit if that flag is set to true.
func checkFlags() {
	if port < 0 || port > 65535 {
		_, _ = fmt.Fprintln(os.Stderr, "invalid port")
		os.Exit(1)
	}

	if verboseLevel < 0 || verboseLevel > 4 {
		_, _ = fmt.Fprintln(os.Stderr, "invalid verbose level")
		os.Exit(1)
	}

	if version {
		_, _ = fmt.Fprintln(os.Stdout, internal.Version)
		os.Exit(0)
	}
}

func main() {
	checkFlags()

	logger := logolang.NewLogger()
	if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "cannot open log file \"%s\": %s", logFile, err)
			os.Exit(1)
		}
		safeF := &logolang.SafeWriter{W: f}
		logger = logolang.NewLoggerWriters(safeF, safeF, safeF, safeF)
	}
	logger.Color = false
	logger.Level = verboseLevel

	server.Log = logger
	server.Run(configPath, strconv.Itoa(port))
}

