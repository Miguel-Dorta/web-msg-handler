package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// TODO tests

// Filename is the default filename of the web-msg-handler config
const Filename = "config.toml"

var (
	// Directory is the default setting directory path
	Directory = "/etc/opt/web-msg-handler"

	// ErrInvalidPort is returned when the config have a invalid port
	ErrInvalidPort = errors.New("invalid port: must be between 0 and 65535")
)

// Config represents the structure of the web-msg-handler config
type Config struct {
	Port       int    `toml:"port"`
	Verbose    int    `toml:"verbose"`
	PIDFile    string `toml:"pid_file"`
	LogOutFile string `toml:"log_output_file"`
	LogErrFile string `toml:"log_error_file"`
}

// Load will read the config from Directory and return a Config object
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
