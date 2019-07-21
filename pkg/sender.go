package pkg

type sender interface {
	Send(msg msg) bool
}
