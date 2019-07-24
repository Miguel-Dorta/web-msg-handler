package internal

import (
	"fmt"
	"net/smtp"
)

type senderMail struct {
	Url    string `json:"url"`
	Mailto string `json:"mailto"`
	Sender struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Hostname string `json:"hostname"`
		Port     string `json:"port"`
	} `json:"sender"`
}

func (sm *senderMail) createMessage(name, mail, msg string) []byte {
	return []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: Message from %s\r\n"+
			"\r\n"+
			"Name: %s\r\n"+
			"Email: %s\r\n"+
			"Message: %s\r\n",
		sm.Sender.Username,
		sm.Mailto,
		sm.Url,
		name, mail, msg,
	))
}

func (sm *senderMail) Send(name, mail, msg string) error {
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
