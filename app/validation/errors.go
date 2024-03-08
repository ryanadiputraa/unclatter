package validation

import (
	"errors"
	"net/http"
)

const (
	// Service error
	BadRequest   = "bad_request"
	Unauthorized = "unauthorized"
	Forbidden    = "forbidden"
	NotFound     = "not_found"
	ServerErr    = "server_err"

	// Oauth errror
	InvalidCallbackParam = "invalid_callback_param"
	ExchangeCodeFailed   = "exchange_code_failed"
)

var (
	HttpErrMap = map[string]int{
		BadRequest:   http.StatusBadRequest,
		Unauthorized: http.StatusUnauthorized,
		Forbidden:    http.StatusForbidden,
		NotFound:     http.StatusNotFound,
		ServerErr:    http.StatusInternalServerError,
	}
)

type Error struct {
	Err     string
	Message string
}

func NewError(errCode, msg string) *Error {
	return &Error{
		Err:     errCode,
		Message: msg,
	}
}

func (e Error) Error() string {
	err := errors.New(e.Message)
	return err.Error()
}
