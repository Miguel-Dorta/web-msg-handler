package src

type sender interface {
	Send(name, mail, msg string) bool
}
