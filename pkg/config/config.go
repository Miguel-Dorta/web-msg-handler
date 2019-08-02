package config

import (
	"encoding/json"
	"fmt"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/sender"
	"io/ioutil"
)

type config struct {
	SitesTg   []sitesTelegram `json:"telegram-sites"`
	SitesMail []sitesMail     `json:"mail-sites"`
}

func LoadConfig(path string) (map[uint64]sender.Sender, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file from path \"%s\": %s", path, err)
	}

	var c config
	if err = json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("error parsing config file from path \"%s\": %s", path, err)
	}

	senders := make(map[uint64]sender.Sender, len(c.SitesMail) + len(c.SitesTg))

	for _, stg := range c.SitesTg {
		if _, exists := senders[stg.Id]; exists {
			return nil, fmt.Errorf("conflicting IDs in config file (ID: %d)", stg.Id)
		}
		senders[stg.Id] = &stg.Site
	}

	for _, sm := range c.SitesMail {
		if _, exists := senders[sm.Id]; exists {
			return nil, fmt.Errorf("conflicting IDs in config file (ID: %d)", sm.Id)
		}
		senders[sm.Id] = &sm.Site
	}

	return senders, nil
}

