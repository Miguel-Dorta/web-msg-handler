package config

import (
	"encoding/json"
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"path/filepath"
)

type Site struct {
	RecaptchaSecret, SenderName, ConfigJS string
}

type siteConfig struct {
	ID              string                 `toml:"id"`
	RecaptchaSecret string                 `toml:"recaptcha_secret"`
	SenderType      string                 `toml:"sender_type"`
	SenderConfig    map[string]interface{} `toml:"sender"`
}

const SitesDirectory = "sites"

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
