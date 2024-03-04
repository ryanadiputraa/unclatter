package http

import (
	"encoding/json"
	"net/http"
)

type responseWriter struct{}

type ResponseWriter interface {
	WriteResponseData(w http.ResponseWriter, code int, data any)
	WriteErrMessage(w http.ResponseWriter, code int, message string)
}

type ResponseData struct {
	Data any `json:"data"`
}

type ErrMessage struct {
	Message string `json:"message"`
}

func NewResponseWriter() ResponseWriter {
	return &responseWriter{}
}

func (rw *responseWriter) WriteResponseData(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&ResponseData{
		Data: data,
	})
}

func (rw *responseWriter) WriteErrMessage(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&ErrMessage{
		Message: message,
	})
}
