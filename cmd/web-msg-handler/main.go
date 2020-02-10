package main

import (
	"github.com/Miguel-Dorta/logolang"
	"os"
)

var log *logolang.Logger

func init() {
	log = logolang.NewLogger()
	log.Color = false
}

func main() {
	if err := cmdRoot.Execute(); err != nil {
		log.Critical(err.Error())
		os.Exit(1)
	}
}