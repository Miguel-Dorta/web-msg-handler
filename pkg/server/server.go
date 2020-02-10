package server

// Package server will manage all the HTTP request made to web-msg-handler.

import (
	"context"
	"github.com/Miguel-Dorta/logolang"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/config"
	"golang.org/x/sys/unix"
	"net/http"
	"os"
	"os/signal"
	"strconv"
)

var (
	log   *logolang.Logger
	sites map[string]*config.Site
)

// Run will start a HTTP server in the port provided using the config file path provided.
// It ends when a termination or interrupt signal is received.
// It can end the program execution prematurely.
func Run(port int, logger *logolang.Logger) {
	log = logger
	err := loadSites()
	if err != nil {
		log.Criticalf("error loading sites config: %s", err)
		os.Exit(1)
	}

	http.HandleFunc("/", handle)
	srv := http.Server{Addr: ":" + strconv.Itoa(port)}

	serverClosed := make(chan bool)
	go func() {
		var (
			quit = make(chan os.Signal, 2)
			reload = make(chan os.Signal, 1)
		)
		signal.Notify(quit, unix.SIGTERM, unix.SIGINT)
		signal.Notify(reload, unix.SIGUSR1)
		defer close(serverClosed)

		for {
			select {
			case <-reload:
				if err := loadSites(); err != nil {
					log.Errorf("error reloading sites config: %s", err)
					log.Info("preserving previous config")
				}
			case <-quit:
				log.Info("Shutting down")
				if err := srv.Shutdown(context.Background()); err != nil {
					log.Criticalf("error while shutting down: %s", err)
					os.Exit(1)
				}
				return
			}
		}
	}()

	log.Infof("Listening port %s", srv.Addr[1:])
	if err = srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Criticalf("Unexpected error which closed the server: %s", err)
		os.Exit(1)
	}
	<-serverClosed
}

// loadSites loads the site configs and sets it to the package variable "sites"
func loadSites() error {
	s, err := config.LoadSites()
	if err != nil {
		return err
	}
	sites = s
	return nil
}
