package api

type Response struct {
	Success bool   `json:"success"`
	Err     string `json:"error,omitempty"`
}
