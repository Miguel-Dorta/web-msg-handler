package sender

import (
	"fmt"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/recaptcha"
	"html"
	"net/smtp"
	"strings"
)

// Mail represents a type that implements the Sender interface that send forms via SMTP
type Mail struct {
	URL             string `json:"url"`
	RecaptchaSecret string `json:"recaptcha-secret"`
	Mailto          string `json:"mailto"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	Hostname        string `json:"hostname"`
	Port            string `json:"port"`
}

// CheckRecaptcha will check if the ReCaptcha response provided have passed the ReCaptcha verification using its
// internal secret.
func (sm *Mail) CheckRecaptcha(resp string) error {
	return recaptcha.CheckRecaptcha(sm.RecaptchaSecret, resp)
}

// Send will send the form provided via SMTP
func (sm *Mail) Send(name, mail, msg string) error {
	err := smtp.SendMail(
		sm.Hostname+":"+sm.Port,
		smtp.PlainAuth("", sm.Username, sm.Password, sm.Hostname),
		sm.Username,
		[]string{sm.Mailto},
		sm.createMessage(name, mail, msg),
	)
	if err != nil {
		return fmt.Errorf("error sending mail: %s", err)
	}
	return nil
}

// createMessage will return a byte slice containing a styled message from the form provided.
func (sm *Mail) createMessage(name, mail, msg string) []byte {
	return []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: Message from %s\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"\r\n"+
			"<html><body>"+
			"<b>Name</b>: %s<br>"+
			"<b>Email</b>: %s<br>"+
			"<b>Message</b>: %s"+
			"</body></html>\r\n",
		sm.Username,
		sm.Mailto,
		sm.URL,
		html.EscapeString(name),
		html.EscapeString(mail),
		lfToBr(html.EscapeString(msg)),
	))
}

// lfToBr will replace the End-of-Line characters ("\r\n" and "\n") for the HTML tag "<br>" (without quotes).
// <br> is the HTML tag that represents line breaks.
func lfToBr(str string) string {
	replacer := strings.NewReplacer("\r", "", "\n", "<br>")
	return replacer.Replace(str)
}
