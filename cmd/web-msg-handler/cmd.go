package main

import (
	"fmt"
	"github.com/Miguel-Dorta/logolang"
	"github.com/Miguel-Dorta/si"
	"github.com/Miguel-Dorta/web-msg-handler/internal"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/config"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/server"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	// Cobra commands
	cmdRoot = &cobra.Command{
		Use: "web-msg-handler",
		Run: root,
	}
	cmdReload = &cobra.Command{
		Use: "reload",
		Short: "reload site configs",
		Run: reload,
	}
	cmdRestart = &cobra.Command{
		Use: "restart",
		Short: "restart web-msg-handler",
		Run: restart,
	}
	cmdStop = &cobra.Command{
		Use: "stop",
		Short: "stop web-msg-handler",
		Run: stop,
	}
	cmdVersion = &cobra.Command{
		Use: "version",
		Short: "print version and exit",
		Run: version,
	}

	installationPath string
)

func init() {
	cmdRoot.PersistentFlags().StringVarP(&installationPath, "installation-path", "p", config.Directory, "set installation path")
	cmdRoot.AddCommand(cmdReload, cmdRestart, cmdStop, cmdVersion)
}

// root will execute when no command is given.
// It starts the service if no other instance is running.
func root(_ *cobra.Command, _ []string) {
	c := loadConf()
	if err := si.Register(c.PIDFile); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	server.Run(c.Port, log)
}

// reload will execute when "reload" command is given.
// It will send a SIGUSR1 signal to a running process in order to reload its sites configs.
func reload(_ *cobra.Command, _ []string) {
	c := loadConf()
	p, err := si.Find(c.PIDFile)
	if err != nil {
		log.Errorf("error finding other running instance of web-msg-handler: %s", err)
		os.Exit(1)
	}

	if p == nil {
		log.Error("there are no running instance of web-msg-handler")
		os.Exit(1)
	}

	if err := p.Signal(unix.SIGUSR1); err != nil {
		log.Errorf("error sending reload signal to other instance of web-msg-handler: %s", err)
		os.Exit(1)
	}

	log.Debug("reload signal delivered. Check the log of web-msg-handler for detecting error")
}

// restart will execute when "restart" command is given.
// It will stop other instance and start this one.
func restart(_ *cobra.Command, _ []string) {
	stop(nil, nil)
	root(nil, nil)
}

// stop will execute when "stop" command is given.
// It will stop another instance of web-msg-handler if it's running.
func stop(_ *cobra.Command, _ []string) {
	c := loadConf()
	p, err := si.Find(c.PIDFile)
	if err != nil {
		log.Errorf("error finding other running instance of web-msg-handler: %s", err)
		os.Exit(1)
	}

	if p == nil {
		return
	}

	if err := p.Signal(os.Interrupt); err != nil {
		log.Errorf("error sending interrupt signal to the other instance of web-msg-handler: %s", err)
		os.Exit(1)
	}

	if _, err := p.Wait(); err != nil {
		log.Errorf("the other instance of web-msg-handler returned an error: %s", err)
		os.Exit(1)
	}
}

// version will execute when "version" command is given.
// It will print web-msg-handler's version and exit.
func version(_ *cobra.Command, _ []string) {
	fmt.Println(internal.Version)
}

// loadConf returns the config that exists in installationPath,
// apply it to this package logger and
// return the config with the PID alias in config.PIDFile
func loadConf() *config.Config {
	config.Directory = installationPath

	c, err := config.Load()
	if err != nil {
		log.Criticalf("error loading config: %s", err)
		os.Exit(1)
	}

	stdout := setWriter(c.LogOutFile, os.Stdout)
	stderr := setWriter(c.LogErrFile, os.Stderr)
	log = logolang.NewLoggerWriters(stdout, stdout, stderr, stderr)
	log.Color = false
	log.Level = c.Verbose

	c.PIDFile = getAlias(c.PIDFile)
	return c
}

// setWriter will return an io.Writer for the file of the specified path or the writer provided if path == ""
func setWriter(path string, w io.Writer) io.Writer {
	if path == "" {
		return &logolang.SafeWriter{W:w}
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Criticalf("error opening log file %s: %s", path, err)
		os.Exit(1)
	}
	return &logolang.SafeWriter{W: f}
}

// getAlias gets a path for a pith file, sets si.Dir to that path's parent dir,
// and return the alias for that file (filename without ".pid")
func getAlias(pidFile string) string {
	pidFile = strings.TrimRight(pidFile, ".pid")
	si.Dir = filepath.Dir(pidFile)
	return filepath.Base(pidFile)
}
