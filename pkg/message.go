package pkg

type msg struct {
	ReceiverId uint64 `json:"id"`
	SenderName string `json:"name"`
	SenderMail string `json:"mail"`
	Message    string `json:"msg"`
}

func (m *msg) process() bool {
	return senders[m.ReceiverId].Send(*m)
}

func getChannelContentMsg(c chan msg) []msg {
	msgs := make([]msg, 0, len(c))
	for len(c) > 0 {
		msgs = append(msgs, <-c)
	}
	return msgs
}
