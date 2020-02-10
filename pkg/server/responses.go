package server

import (
	"github.com/Miguel-Dorta/web-msg-handler/pkg"
	"net/http"
)

type httpResponse struct {
	status  int
	success bool
	msg     string
}

const statusUnknownError = 502

var (
	ErrNotFound = &httpResponse{
		success: false,
		status:  http.StatusNotFound,
		msg:     "not found",
	}
	ErrMethodNotAllowed = &httpResponse{
		success: false,
		status:  http.StatusMethodNotAllowed,
		msg:     "method not allowed",
	}
	ErrContentTypeNotAllowed = &httpResponse{
		success: false,
		status:  http.StatusBadRequest,
		msg:     pkg.MimeContentType + " not allowed",
	}
	ErrMalformedJSON = &httpResponse{
		success: false,
		status:  http.StatusBadRequest,
		msg:     "malformed JSON",
	}
	ErrInvalidMail = &httpResponse{
		success: false,
		status:  http.StatusBadRequest,
		msg:     "invalid email",
	}
	ErrRecaptchaVerificationFailed = &httpResponse{
		success: false,
		status:  http.StatusBadRequest,
		msg:     "reCAPTCHA verification failed",
	}
	ErrReadingBody = &httpResponse{
		success: false,
		status:  statusUnknownError,
		msg:     "unknown error reading request body",
	}
	ErrUnknown = &httpResponse{
		success: false,
		status:  statusUnknownError,
		msg:     "unknown error",
	}
	ErrInternalServerError = &httpResponse{
		success: false,
		status:  http.StatusInternalServerError,
		msg:     "internal server error",
	}
	ResponseOK = &httpResponse{
		success: true,
		status:  http.StatusOK,
		msg:     "",
	}
)
