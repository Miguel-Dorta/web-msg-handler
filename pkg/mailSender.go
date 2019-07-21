package pkg

type mailSender struct {
	Mailto string `json:"mailto"`
}

func (ms *mailSender) Send(msg msg) bool {
	return false
}
