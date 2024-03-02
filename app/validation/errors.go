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
	Timeout      = "bad_gateway"

	// Oauth errror
	InvalidCallbackParam = "invalid_callback_param"
	ExchangeCodeFailed   = "exchange_code_failed"
)

var (
	HttpErrMap = map[string]int{
		"bad_request":  http.StatusBadRequest,
		"unauthorized": http.StatusUnauthorized,
		"forbidden":    http.StatusForbidden,
		"not_found":    http.StatusNotFound,
		"server_err":   http.StatusInternalServerError,
		"bad_gateway":  http.StatusBadGateway,
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
