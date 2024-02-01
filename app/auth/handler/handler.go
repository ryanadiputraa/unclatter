package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ryanadiputraa/unclatter/app/user"
	"github.com/ryanadiputraa/unclatter/pkg/oauth"
)

type handler struct {
	userService user.UserService
	googleOauth oauth.GoogleOauth
}

func NewHandler(r *echo.Group, userService user.UserService, googleOauth oauth.GoogleOauth) {
	h := &handler{
		userService: userService,
		googleOauth: googleOauth,
	}

	r.GET("/signin/google", h.GoogleSignIn())
}

func (h *handler) GoogleSignIn() echo.HandlerFunc {
	return func(c echo.Context) error {
		url := h.googleOauth.GetSignInURL()
		return c.Redirect(http.StatusTemporaryRedirect, url)
	}
}
