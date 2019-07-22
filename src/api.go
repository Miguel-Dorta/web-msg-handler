package src

type request struct {
	Name string `json:"name"`
	Mail string `json:"mail"`
	Msg  string `json:"msg"`
}

type response struct {
	Success bool   `json:"success"`
	Err     string `json:"error,omitempty"`
}
