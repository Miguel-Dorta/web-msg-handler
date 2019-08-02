package config

import (
	"github.com/Miguel-Dorta/web-msg-handler/pkg/sender"
)

type sitesMail struct {
	Id   uint64      `json:"id"`
	Site sender.Mail `json:"site"`
}
