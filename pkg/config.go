package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type config struct {
	TgSites   []telegramSites `json:"telegram-sites"`
	MailSites []mailSites     `json:"mail-sites"`
}

type telegramSites struct {
	Id     uint64         `json:"id"`
	Sender telegramSender `json:"sender"`
}

type mailSites struct {
	Id     uint64     `json:"id"`
	Sender mailSender `json:"sender"`
}

func LoadConfig(path string) (map[uint64]sender, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file from path \"%s\": %s", path, err)
	}

	var c config
	if err = json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("error parsing config file from path \"%s\": %s", path, err)
	}

	senders := make(map[uint64]sender, len(c.MailSites) + len(c.TgSites))

	for _, tgs := range c.TgSites {
		senders[tgs.Id] = &tgs.Sender
	}
	for _, ms := range c.MailSites {
		senders[ms.Id] = &ms.Sender
	}

	return senders, nil
}
