package pkg

type telegramSender struct {
	ChatId   string `json:"chat-id"`
	BotToken string `json:"bot-token"`
}

func (tgSender *telegramSender) Send(msg msg) bool {
	return false
}
