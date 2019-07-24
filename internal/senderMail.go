package internal

type senderMail struct {
	Url    string `json:"url"`
	Mailto string `json:"mailto"`
}

func (sm *senderMail) Send(name, mail, msg string) error {
	return nil
}
