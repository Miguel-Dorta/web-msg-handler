package config

import (
	"encoding/json"
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"path/filepath"
)

// TODO tests

var (
	Directory = "/etc/web-msg-handler"
	Filename = "config.toml"
)

type config struct {
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

func LoadConfig() (*config, map[string]*siteConfig, error) {
	data, err := ioutil.ReadFile(filepath.Join(Directory, Filename))
	if err != nil {
		return nil, nil, fmt.Errorf("error loading config: %w", err)
	}

	var c config
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

func loadSites(path string) (map[string]*siteConfig, error) {
	sites, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error listing sites directory: %w", err)
	}

	sitesMap := make(map[string]*siteConfig, len(sites))
	for _, site := range sites {
		if !site.Mode().IsRegular() {
			continue
		}

		sitePath := filepath.Join(path, site.Name())
		data, err := ioutil.ReadFile(sitePath)
		if err != nil {
			return nil, fmt.Errorf("error reading file \"%s\": %w", sitePath, err)
		}

		var sc siteConfig
		if err := toml.Unmarshal(data, &sc); err != nil {
			return nil, fmt.Errorf("error parsing site config from file \"%s\": %w", sitePath, err)
		}

		sitesMap[sc.ID] = &sc
	}

	return sitesMap, nil
}
