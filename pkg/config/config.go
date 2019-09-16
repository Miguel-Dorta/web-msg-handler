package config

// Package config is the package that manages the functions related to the config file of web-msg-handler.

import (
	"encoding/json"
	"fmt"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/sender"
	"io/ioutil"
)

// config represents the structure of the config file of web-msg-handler.
// It consists in an array of sites, which contains an numeric ID, an URL a ReCaptcha Secret and a Sender.
// The sender is an object that contains a key "type" that indicates which kind of sender it is, and a dynamic "settings" object.
//
// The "settings" object will have the following keys depending of the sender type value.
//
// Mail type
// 		"mailto":   string
// 		"username": string
// 		"password": string
// 		"hostname": string
// 		"port":     string
//
// Telegram type
// 		"chatID":   string
// 		"botToken": string
type config struct {
	Sites []struct {
		ID              uint64 `json:"id"`
		URL             string `json:"url"`
		RecaptchaSecret string `json:"recaptchaSecret"`
		Sender          struct {
			Type     string            `json:"type"`
			Settings map[string]string `json:"settings"`
		} `json:"sender"`
	} `json:"sites"`
}

// LoadConfig will read the config from the path provided and return a map of sender.Sender with uint64 key.
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
		switch s.Sender.Type {
		case "mail":
			parsedSender = &sender.Mail{
				URL: s.URL,
				RecaptchaSecret: s.RecaptchaSecret,
				Mailto: s.Sender.Settings["mailto"],
				Username: s.Sender.Settings["username"],
				Password: s.Sender.Settings["password"],
				Hostname: s.Sender.Settings["hostname"],
				Port: s.Sender.Settings["port"],
			}
		case "telegram":
			parsedSender = &sender.Telegram{
				URL: s.URL,
				RecaptchaSecret: s.RecaptchaSecret,
				ChatId: s.Sender.Settings["chatID"],
				BotToken: s.Sender.Settings["botToken"],
			}
		default:
			return nil, fmt.Errorf("error parsing config file in site with ID %d: invalid sender type \"%s\"", s.ID, s.Sender.Type)
		}

		if _, exists := senders[s.ID]; exists {
			return nil, fmt.Errorf("conflicting IDs in config file (ID: %d)", s.ID)
		}
		senders[s.ID] = parsedSender
	}

	return senders, nil
}
