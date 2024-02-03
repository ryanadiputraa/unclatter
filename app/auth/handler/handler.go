package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ryanadiputraa/unclatter/app/auth"
	"github.com/ryanadiputraa/unclatter/app/user"
	"github.com/ryanadiputraa/unclatter/app/validation"
	"github.com/ryanadiputraa/unclatter/config"
	"github.com/ryanadiputraa/unclatter/pkg/jwt"
	"github.com/ryanadiputraa/unclatter/pkg/logger"
	"github.com/ryanadiputraa/unclatter/pkg/oauth"
)

type handler struct {
	config      *config.Config
	log         logger.Logger
	authService auth.AuthService
	userService user.UserService
	googleOauth oauth.GoogleOauth
	jwtTokens   jwt.JWTTokens
}

func NewHandler(
	r *echo.Group,
	config *config.Config,
	log logger.Logger,
	authService auth.AuthService,
	userService user.UserService,
	googleOauth oauth.GoogleOauth,
	jwtTokens jwt.JWTTokens,
) {
	h := &handler{
		config:      config,
		log:         log,
		authService: authService,
		userService: userService,
		googleOauth: googleOauth,
		jwtTokens:   jwtTokens,
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
			h.log.Error("auth handler: fail to exchange", err.Error())
			return h.redirectWithError(validation.ExchangeCodeFailed)(c)
		}

		user, err := h.userService.CreateUser(c.Request().Context(), user.NewUserArg{
			Email:     userInfo.Email,
			FirstName: userInfo.FirstName,
			LastName:  userInfo.LastName,
		})
		if err != nil {
			return h.redirectWithError(validation.ServerErr)(c)
		}

		_, err = h.authService.AddUserAuthProvider(c.Request().Context(), auth.NewAuthProviderArg{
			Provider:       "google",
			ProviderUserID: userInfo.ID,
			UserID:         user.ID,
		})
		if err != nil {
			return h.redirectWithError(validation.ServerErr)(c)
		}

		jwtTokens, err := h.jwtTokens.GenereateJWTWithClaims(user.ID)
		if err != nil {
			h.log.Error("auth handler: fail to generate jwt tokens", err.Error())
			return h.redirectWithError(validation.ServerErr)(c)
		}

		return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%v/auth?access_token=%v", h.config.FrontendURL, jwtTokens.AccessToken))
	}
}

func (h *handler) redirectWithError(err string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, h.config.FrontendURL+"/auth?err="+err)
	}
}
