package config

import (
	"github.com/Miguel-Dorta/web-msg-handler/pkg/sender"
)

type sitesTelegram struct {
	Id   uint64          `json:"id"`
	Site sender.Telegram `json:"site"`
}
