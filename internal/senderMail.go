package internal

type senderMail struct {
	Mailto string `json:"mailto"`
}

func (sm *senderMail) Send(name, mail, msg string) bool {
	return false
}
