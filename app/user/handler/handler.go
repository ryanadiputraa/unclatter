package handler

import (
	"net/http"

	"github.com/ryanadiputraa/unclatter/app/middleware"
	"github.com/ryanadiputraa/unclatter/app/user"
	"github.com/ryanadiputraa/unclatter/app/validation"
	_http "github.com/ryanadiputraa/unclatter/pkg/http"
)

type handler struct {
	rw          _http.ResponseWriter
	userService user.UserService
}

func NewUserHandler(web *http.ServeMux, rw _http.ResponseWriter, userService user.UserService, authMiddleware middleware.AuthMiddleware) {
	h := &handler{
		rw:          rw,
		userService: userService,
	}

	web.Handle("GET /api/users", authMiddleware.ParseJWTToken(h.getUserInfo()))
}

func (h *handler) getUserInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ac := r.Context().(*middleware.AuthContext)
		user, err := h.userService.GetUserInfo(ac.Context, ac.UserID)
		if err != nil {
			if vErr, ok := err.(*validation.Error); ok {
				h.rw.WriteErrMessage(w, validation.HttpErrMap[vErr.Err], vErr.Message)
				return
			}
			h.rw.WriteErrMessage(w, http.StatusInternalServerError, "internal server error")
			return
		}

		h.rw.WriteResponseData(w, http.StatusOK, user)
	}
}
