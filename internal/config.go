package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type config struct {
	SitesTg   []sitesTelegram `json:"telegram-sites"`
	SitesMail []sitesMail     `json:"mail-sites"`
}

type sitesTelegram struct {
	Id     uint64         `json:"id"`
	Sender senderTelegram `json:"sender"`
}

type sitesMail struct {
	Id     uint64     `json:"id"`
	Sender senderMail `json:"sender"`
}

func loadConfig(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading config file from path \"%s\": %s", path, err)
	}

	var c config
	if err = json.Unmarshal(data, &c); err != nil {
		return fmt.Errorf("error parsing config file from path \"%s\": %s", path, err)
	}

	senders = make(map[uint64]sender, len(c.SitesMail) + len(c.SitesTg))

	for _, stg := range c.SitesTg {
		if _, exists := senders[stg.Id]; exists {
			return fmt.Errorf("conflicting IDs in config file (ID: %d)", stg.Id)
		}
		senders[stg.Id] = &stg.Sender
	}

	for _, sm := range c.SitesMail {
		if _, exists := senders[sm.Id]; exists {
			return fmt.Errorf("conflicting IDs in config file (ID: %d)", sm.Id)
		}
		senders[sm.Id] = &sm.Sender
	}

	return nil
}
