package http

import (
	"encoding/json"
	"net/http"

	"github.com/ryanadiputraa/unclatter/app/pagination"
)

type responseWriter struct{}

type ResponseWriter interface {
	WriteResponseData(w http.ResponseWriter, code int, data any)
	WriteResponseDataWithPagination(w http.ResponseWriter, code int, data any, meta pagination.Meta)
	WriteErrMessage(w http.ResponseWriter, code int, message string)
	WriteErrDetails(w http.ResponseWriter, code int, message string, errMap map[string]string)
}

type ResponseData struct {
	Data any `json:"data"`
}

type ResponseWithPaginationMeta struct {
	Data any             `json:"data"`
	Meta pagination.Meta `json:"meta"`
}

type ErrMessage struct {
	Message string `json:"message"`
}

type ErrDetails struct {
	Message string            `json:"message"`
	Error   map[string]string `json:"error"`
}

func NewResponseWriter() ResponseWriter {
	return &responseWriter{}
}

func (rw *responseWriter) setJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func (rw *responseWriter) WriteResponseData(w http.ResponseWriter, code int, data any) {
	rw.setJSONHeader(w)
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&ResponseData{
		Data: data,
	})
}

func (rw *responseWriter) WriteResponseDataWithPagination(w http.ResponseWriter, code int, data any, meta pagination.Meta) {
	rw.setJSONHeader(w)
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&ResponseWithPaginationMeta{
		Data: data,
		Meta: meta,
	})
}

func (rw *responseWriter) WriteErrMessage(w http.ResponseWriter, code int, message string) {
	rw.setJSONHeader(w)
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&ErrMessage{
		Message: message,
	})
}

func (rw *responseWriter) WriteErrDetails(w http.ResponseWriter, code int, message string, errMap map[string]string) {
	rw.setJSONHeader(w)
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&ErrDetails{
		Message: message,
		Error:   errMap,
	})
}
