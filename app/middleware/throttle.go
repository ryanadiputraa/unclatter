package middleware

import (
	"net/http"

	_http "github.com/ryanadiputraa/unclatter/pkg/http"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(rate.Limit(100), 200)

func ThrottleMiddleware(next http.Handler, rw _http.ResponseWriter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			rw.WriteErrMessage(w, http.StatusTooManyRequests, "too fast")
			return
		}
		next.ServeHTTP(w, r)
	}
}
