package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"path/filepath"
)

// TODO tests

const (
	SitesDirectory = "sites"
	Filename = "config.toml"
)

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

type siteConfig struct {
	ID              string                 `toml:"id"`
	RecaptchaSecret string                 `toml:"recaptcha_secret"`
	SenderType      string                 `toml:"sender_type"`
	SenderConfig    map[string]interface{} `toml:"sender"`
}

type Site struct {
	RecaptchaSecret, SenderName, ConfigJS string
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

func LoadSites() (map[string]*Site, error) {
	path := filepath.Join(Directory, SitesDirectory)

	sites, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error listing sites directory: %w", err)
	}

	sitesMap := make(map[string]*Site, len(sites))
	for _, s := range sites {
		if !s.Mode().IsRegular() {
			continue
		}

		sitePath := filepath.Join(path, s.Name())
		data, err := ioutil.ReadFile(sitePath)
		if err != nil {
			return nil, fmt.Errorf("error reading file \"%s\": %w", sitePath, err)
		}

		var sc siteConfig
		if err := toml.Unmarshal(data, &sc); err != nil {
			return nil, fmt.Errorf("error parsing site config from file \"%s\": %w", sitePath, err)
		}

		configJSON, err := json.Marshal(sc.SenderConfig)
		if err != nil {
			return nil, fmt.Errorf("error generating config JSON for plugin %s: %w", s.Name(), err)
		}

		sitesMap[sc.ID] = &Site{
			RecaptchaSecret: sc.RecaptchaSecret,
			SenderName:      sc.SenderType,
			ConfigJS:        string(configJSON),
		}
	}

	return sitesMap, nil
}
