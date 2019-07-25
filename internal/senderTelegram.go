package internal

import (
	"encoding/json"
	"fmt"
	"html"
)

const (
	telegramBotApiUrl = "https://api.telegram.org/bot"
	sendMsgMethod     = "/sendMessage"
	parseModeHtml     = "HTML"
)

type senderTelegram struct {
	Url             string `json:"url"`
	RecaptchaSecret string `json:"recaptcha-secret"`
	ChatId          string `json:"chat-id"`
	BotToken        string `json:"bot-token"`
}

type messageSend struct {
	ChatId                 string `json:"chat_id"`
	Text                   string `json:"text"`
	ParseMode              string `json:"parse_mode"`
	DisableWebImagePreview bool   `json:"disable_web_page_preview"`
}

func (st *senderTelegram) createMessage(name, mail, msg string) string {
	return fmt.Sprintf(
		"Message from %s\n" +
			"\n" +
			"<b>Name</b>: %s\n" +
			"<b>Email</b>: %s\n" +
			"<b>Message</b>: %s",
		st.Url,
		html.EscapeString(name),
		html.EscapeString(mail),
		html.EscapeString(msg),
	)
}

func (st *senderTelegram) Send(name, mail, msg string) error {
	data, err := json.Marshal(messageSend{
		ChatId:                 st.ChatId,
		Text:                   st.createMessage(name, mail, msg),
		ParseMode:              parseModeHtml,
		DisableWebImagePreview: true,
	})
	if err != nil {
		return fmt.Errorf("error parsing message JSON: %s", err)
	}

	resp, err := postJson(telegramBotApiUrl+st.BotToken+sendMsgMethod, data)
	if err != nil {
		return fmt.Errorf("error doing request to Telegram servers: %s", err.Error())
	}

	var respJson map[string]interface{}
	if err = json.Unmarshal(resp, &respJson); err != nil {
		return fmt.Errorf("error parsing response JSON: %s", err)
	}

	if !respJson["ok"].(bool) {
		return fmt.Errorf("request failed: %s", resp)
	}

	return nil
}

func (st *senderTelegram) CheckRecaptcha(resp string) error {
	return checkRecaptcha(st.RecaptchaSecret, resp)
}
