package sender

import (
	"encoding/json"
	"fmt"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/client"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/recaptcha"
	"html"
)

const (
	telegramBotApiUrl = "https://api.telegram.org/bot"
	sendMsgMethod     = "/sendMessage"
	parseModeHtml     = "HTML"
)

// Telegram represents a type that implements the Sender interface that send forms via Telegram Bot API
type Telegram struct {
	URL             string `json:"url"`
	RecaptchaSecret string `json:"recaptcha-secret"`
	ChatId          string `json:"chat-id"`
	BotToken        string `json:"bot-token"`
}

// requestJSON represents the request that Telegram.Send() will do to the Telegram Bot API
type requestJSON struct {
	ChatID                 string `json:"chat_id"`
	Text                   string `json:"text"`
	ParseMode              string `json:"parse_mode"`
	DisableWebImagePreview bool   `json:"disable_web_page_preview"`
}

// CheckRecaptcha will check if the ReCaptcha response provided have passed the ReCaptcha verification using its
// internal secret.
func (st *Telegram) CheckRecaptcha(resp string) error {
	return recaptcha.CheckRecaptcha(st.RecaptchaSecret, resp)
}

// Send will send the form provided via Telegram Bot API
func (st *Telegram) Send(name, mail, msg string) error {
	data, err := json.Marshal(requestJSON{
		ChatID:                 st.ChatId,
		Text:                   st.createMessage(name, mail, msg),
		ParseMode:              parseModeHtml,
		DisableWebImagePreview: true,
	})
	if err != nil {
		return fmt.Errorf("error parsing message JSON: %s", err)
	}

	resp, err := client.PostJSON(telegramBotApiUrl+st.BotToken+sendMsgMethod, data)
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

// createMessage will return an string containing a styled message from the form provided.
func (st *Telegram) createMessage(name, mail, msg string) string {
	return fmt.Sprintf(
		"Message from %s\n" +
			"\n" +
			"<b>Name</b>: %s\n" +
			"<b>Email</b>: %s\n" +
			"<b>Message</b>: %s",
		st.URL,
		html.EscapeString(name),
		html.EscapeString(mail),
		html.EscapeString(msg),
	)
}
