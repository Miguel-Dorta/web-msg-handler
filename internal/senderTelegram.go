package internal

type senderTelegram struct {
	ChatId   string `json:"chat-id"`
	BotToken string `json:"bot-token"`
}

func (st *senderTelegram) Send(name, mail, msg string) bool {
	return false
}
