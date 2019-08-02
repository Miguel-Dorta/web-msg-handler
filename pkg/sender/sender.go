package sender

type Sender interface {
	CheckRecaptcha(resp string) error
	Send(name, mail, msg string) error
}
