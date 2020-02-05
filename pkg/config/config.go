package config

import (
	"encoding/json"
	"fmt"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/plugin"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"path/filepath"
)

// TODO tests

var (
	Directory = "/etc/web-msg-handler"
	Filename = "config.toml"
)

type Config struct {
	Port       int    `toml:"port"`
	Verbose    int    `toml:"verbose"`
	SitesPath  string `toml:"sites"`
	PIDFile    string `toml:"pid_file"`
	LockFile   string `toml:"lock_file"`
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
	RecaptchaSecret, SenderType, ConfigJS string
}

func LoadConfig() (*Config, map[string]*Site, error) {
	data, err := ioutil.ReadFile(filepath.Join(Directory, Filename))
	if err != nil {
		return nil, nil, fmt.Errorf("error loading config: %w", err)
	}

	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, nil, fmt.Errorf("error parsing config: %w", err)
	}

	sitesPath := c.SitesPath
	if !filepath.IsAbs(sitesPath) {
		sitesPath = filepath.Join(Directory, sitesPath)
	}
	sites, err := loadSites(sitesPath)
	if err != nil {
		return nil, nil, err
	}
	return &c, sites, nil
}

func loadSites(path string) (map[string]*Site, error) {
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

		configJSON, err := plugin.ArgsToJS(sc.SenderConfig)
		if err != nil {
			return nil, fmt.Errorf("error generating config JSON for plugin %s: %w", s.Name(), err)
		}

		sitesMap[sc.ID] = &Site{
			RecaptchaSecret: sc.RecaptchaSecret,
			SenderType:      sc.SenderType,
			ConfigJS:        configJSON,
		}
	}

	return sitesMap, nil
}
