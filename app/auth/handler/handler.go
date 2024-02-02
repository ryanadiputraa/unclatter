package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ryanadiputraa/unclatter/app/user"
	"github.com/ryanadiputraa/unclatter/app/validation"
	"github.com/ryanadiputraa/unclatter/config"
	"github.com/ryanadiputraa/unclatter/pkg/logger"
	"github.com/ryanadiputraa/unclatter/pkg/oauth"
)

type handler struct {
	config      *config.Config
	log         logger.Logger
	userService user.UserService
	googleOauth oauth.GoogleOauth
}

func NewHandler(r *echo.Group, config *config.Config, log logger.Logger, userService user.UserService, googleOauth oauth.GoogleOauth) {
	h := &handler{
		config:      config,
		log:         log,
		userService: userService,
		googleOauth: googleOauth,
	}

	r.GET("/signin/google", h.GoogleSignIn())
	r.GET("/signin/google/callback", h.GoogleCallback())
}

func (h *handler) GoogleSignIn() echo.HandlerFunc {
	return func(c echo.Context) error {
		url := h.googleOauth.GetSignInURL()
		return c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func (h *handler) GoogleCallback() echo.HandlerFunc {
	return func(c echo.Context) error {
		code := c.QueryParam("code")
		state := c.QueryParam("state")

		if state != h.config.GoogleOauth.State || code == "" {
			h.log.Warn(fmt.Sprintf("auth handler: invalid callback; state=\"%v\" code=\"%v\"", state, code))
			return h.redirectWithError(validation.InvalidCallbackParam)(c)
		}

		userInfo, err := h.googleOauth.ExchangeCodeWithUserInfo(c.Request().Context(), code)
		if err != nil {
			h.log.Error("auth handler: fail to exchange; " + err.Error())
			return h.redirectWithError(validation.ExchangeCodeFailed)(c)
		}

		_, err = h.userService.CreateUser(c.Request().Context(), user.CreateUserArg{
			Email:     userInfo.Email,
			FirstName: userInfo.FirstName,
			LastName:  userInfo.LastName,
		})
		if err != nil {
			h.log.Error("auth handler: fail to register user; " + err.Error())
			return h.redirectWithError(validation.ServerErr)(c)
		}

		// TODO: save auth provider

		// TODO: redirect with jwt tokens

		return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%v/auth/token=<token>", h.config.FrontendURL))
	}
}

func (h *handler) redirectWithError(err string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, h.config.FrontendURL+"/auth?err="+err)
	}
}
