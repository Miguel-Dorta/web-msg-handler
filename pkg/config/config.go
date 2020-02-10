package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// TODO tests

const Filename = "config.toml"

var (
	Directory = "/etc/web-msg-handler"
	ErrInvalidPort = errors.New("invalid port: must be between 0 and 65535")
)

type Config struct {
	Port       int    `toml:"port"`
	Verbose    int    `toml:"verbose"`
	PIDFile    string `toml:"pid_file"`
	LogOutFile string `toml:"log_output_file"`
	LogErrFile string `toml:"log_error_file"`
}

func Load() (*Config, error) {
	data, err := ioutil.ReadFile(filepath.Join(Directory, Filename))
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	if c.Port < 0 || c.Port > 65535 {
		return nil, ErrInvalidPort
	}

	return &c, nil
}
