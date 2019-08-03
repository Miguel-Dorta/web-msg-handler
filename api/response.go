package api

// Response represents the content of the request that web-msg-handler will reply.
// It's always a JSON with a boolean "success" field that indicates if the request was accepted successfully
// and an "error" field that indicates the error found in the case of a failed request.
type Response struct {
	Success bool   `json:"success"`
	Err     string `json:"error,omitempty"`
}
