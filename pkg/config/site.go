package config

import (
	"encoding/json"
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"path/filepath"
)

// Site is the object generated for each site when loading the config.
// It consists in a RecaptchaSecret, a SenderName (that will match the name of a plugin)
// and a ConfigJSON that will be generated from the settings.toml
type Site struct {
	RecaptchaSecret, WebUrl, SenderName, ConfigJSON string
}

// siteConfig is the internal type for unmarshalling the site configs.
type siteConfig struct {
	ID              string                 `toml:"id"`
	RecaptchaSecret string                 `toml:"recaptcha_secret"`
	SenderType      string                 `toml:"sender_type"`
	WebUrl          string                 `toml:"web_url"`
	SenderConfig    map[string]interface{} `toml:"sender"`
}

// SitesDirectory is the name of the subdirectory (of Directory) that contains the site configs.
const SitesDirectory = "sites"

// LoadSites will read the site configs and return a map where the key is the site ID and the value is the site itself.
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

		if _, exists := sitesMap[sc.ID]; exists {
			return nil, fmt.Errorf("site ID collition: %s", sc.ID)
		}

		configJSON, err := json.Marshal(sc.SenderConfig)
		if err != nil {
			return nil, fmt.Errorf("error generating config JSON for plugin %s: %w", s.Name(), err)
		}

		if sc.WebUrl == "" {
			sc.WebUrl = "*"
		}

		sitesMap[sc.ID] = &Site{
			RecaptchaSecret: sc.RecaptchaSecret,
			WebUrl:          sc.WebUrl,
			SenderName:      sc.SenderType,
			ConfigJSON:      string(configJSON),
		}
	}

	return sitesMap, nil
}
