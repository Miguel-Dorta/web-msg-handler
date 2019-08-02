package api

type Request struct {
	Name      string `json:"name"`
	Mail      string `json:"mail"`
	Msg       string `json:"msg"`
	Recaptcha string `json:"g-recaptcha-response"`
}
