package sender

import (
	"fmt"
	"html"
	"net/smtp"
	"strings"
)

type Mail struct {
	Url             string `json:"url"`
	RecaptchaSecret string `json:"recaptcha-secret"`
	Mailto          string `json:"mailto"`
	Sender          struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Hostname string `json:"hostname"`
		Port     string `json:"port"`
	} `json:"sender"`
}

func lfToBr(str string) string {
	replacer := strings.NewReplacer("\r", "", "\n", "<br>")
	return replacer.Replace(str)
}

func (sm *Mail) createMessage(name, mail, msg string) []byte {
	return []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: Message from %s\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n" +
			"\r\n"+
			"<html><body>" +
			"<b>Name</b>: %s<br>"+
			"<b>Email</b>: %s<br>"+
			"<b>Message</b>: %s" +
			"</body></html>\r\n",
		sm.Sender.Username,
		sm.Mailto,
		sm.Url,
		html.EscapeString(name),
		html.EscapeString(mail),
		lfToBr(html.EscapeString(msg)),
	))
}

func (sm *Mail) Send(name, mail, msg string) error {
	err := smtp.SendMail(
		sm.Sender.Hostname+":"+sm.Sender.Port,
		smtp.PlainAuth("", sm.Sender.Username, sm.Sender.Password, sm.Sender.Hostname),
		sm.Sender.Username,
		[]string{sm.Mailto},
		sm.createMessage(name, mail, msg),
	)
	if err != nil {
		return fmt.Errorf("error sending mail: %s", err)
	}
	return nil
}

func (sm *Mail) CheckRecaptcha(resp string) error {
	return checkRecaptcha(sm.RecaptchaSecret, resp)
}

