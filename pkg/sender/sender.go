package sender

// Sender represents any type that's able to verify a ReCaptcha response and Send a form that fits in api.Request,
// both by its own.
type Sender interface {
	CheckRecaptcha(resp string) error
	Send(name, mail, msg string) error
}
