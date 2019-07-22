package main

import (
	"flag"
	"fmt"
	"github.com/Miguel-Dorta/logolang"
	"github.com/Miguel-Dorta/web-msg-handler/internal"
	"os"
	"strconv"
)

var (
	configPath string
	logFile string
	port int
	verbose int
	version bool
)

func init() {
	flag.StringVar(&configPath, "config", "config.json", "set config path")
	flag.StringVar(&logFile, "log-file", "", "set log file")
	flag.IntVar(&port, "port", 8080, "set port")
	flag.IntVar(&verbose, "verbose", 2, "set verbose level. 0=no-log, 1=critical, 2=errors, 3=info, 4 debug")
	flag.BoolVar(&version, "version", false, "print version and exit")

	flag.Parse()
}

func checkFlags() {
	if port < 0 || port > 65535 {
		_, _ = fmt.Fprintln(os.Stderr, "invalid port")
		os.Exit(1)
	}

	if verbose < 0 || verbose > 4 {
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

	if logFile == "" {
		internal.Log = logolang.NewLogger()
	} else {
		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "cannot open log file \"%s\": %s", logFile, err)
			os.Exit(1)
		}
		safeF := &logolang.SafeWriter{W: f}
		internal.Log = logolang.NewLoggerWriters(safeF, safeF, safeF, safeF)
	}
	internal.Log.Color = false
	internal.Log.Level = verbose

	internal.Run(configPath, strconv.Itoa(port))
}

