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

func root(_ *cobra.Command, _ []string) {
	c := loadConf()
	if err := si.Register(c.PIDFile); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	server.Run(c.Port)
}

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

func restart(_ *cobra.Command, _ []string) {
	stop(nil, nil)
	root(nil, nil)
}

func stop(_ *cobra.Command, _ []string) {
	c := loadConf()
	p, err := si.Find(c.PIDFile)
	if err != nil {
		log.Errorf("error finding other running instance of web-msg-handler: %s", err)
		os.Exit(1)
	}

	if p == nil {
		os.Exit(0)
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

func version(_ *cobra.Command, _ []string) {
	fmt.Println(internal.Version)
}

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

func getAlias(pidFile string) string {
	pidFile = strings.TrimRight(pidFile, ".pid")
	si.Dir = filepath.Dir(pidFile)
	return filepath.Base(pidFile)
}
