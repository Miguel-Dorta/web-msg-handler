package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/sender"
	"io/ioutil"
)

type config struct {
	Sites []site `json:"sites"`
}

type site struct {
	ID              uint64      `json:"id"`
	URL             string      `json:"url"`
	RecaptchaSecret string      `json:"recaptcha-secret"`
	Sender          interface{} `json:"sender"`
}

type mail struct {
	Mailto   string `json:"mailto"`
	Username string `json:"username"`
	Password string `json:"password"`
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
}

type telegram struct {
	ChatID   string `json:"chat-id"`
	BotToken string `json:"bot-token"`
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

	senders := make(map[uint64]sender.Sender, len(c.Sites))

	for _, s := range c.Sites {
		var parsedSender sender.Sender
		if mailSender, ok := s.Sender.(*mail); ok {
			parsedSender = &sender.Mail{
				Url:             s.URL,
				RecaptchaSecret: s.RecaptchaSecret,
				Mailto:          mailSender.Mailto,
				Username:        mailSender.Username,
				Password:        mailSender.Password,
				Hostname:        mailSender.Hostname,
				Port:            mailSender.Port,
			}
		} else if telegramSender, ok := s.Sender.(*telegram); ok {
			parsedSender = &sender.Telegram{
				Url:             s.URL,
				RecaptchaSecret: s.RecaptchaSecret,
				ChatId:          telegramSender.ChatID,
				BotToken:        telegramSender.BotToken,
			}
		} else {
			return nil, errors.New("error parsing config file: invalid sender")
		}

		if _, exists := senders[s.ID]; exists {
			return nil, fmt.Errorf("conflicting IDs in config file (ID: %d)", s.ID)
		}
		senders[s.ID] = parsedSender
	}

	return senders, nil
}
