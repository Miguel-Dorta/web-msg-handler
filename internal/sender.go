package internal

type sender interface {
	Send(name, mail, msg string) error
}
