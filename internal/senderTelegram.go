package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	telegramBotApiUrl = "https://api.telegram.org/bot"
	sendMsgMethod     = "/sendMessage"
	parseModeMarkdown = "Markdown"
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
		"Message from %s\n\n*Name*: %s\n*Email*: %s\n*Message*: %s",
		st.Url, name, mail, msg,
	)
}

func (st *senderTelegram) Send(name, mail, msg string) error {
	data, err := json.Marshal(messageSend{
		ChatId:                 st.ChatId,
		Text:                   st.createMessage(name, mail, msg),
		ParseMode:              parseModeMarkdown,
		DisableWebImagePreview: true,
	})
	if err != nil {
		return fmt.Errorf("error parsing message JSON: %s", err)
	}

	resp, err := httpClient.Post(telegramBotApiUrl+st.BotToken+sendMsgMethod, mimeJson, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed http request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %s", err)
	}

	if err = resp.Body.Close(); err != nil {
		return fmt.Errorf("error closing response body: %s", err)
	}

	var bodyJson map[string]interface{}
	if err = json.Unmarshal(body, &bodyJson); err != nil {
		return fmt.Errorf("error parsing response JSON: %s", err)
	}

	if !bodyJson["ok"].(bool) {
		return fmt.Errorf("request failed: %s", body)
	}

	return nil
}
