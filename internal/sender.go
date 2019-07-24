package internal

type sender interface {
	CheckRecaptcha(resp string) error
	Send(name, mail, msg string) error
}
