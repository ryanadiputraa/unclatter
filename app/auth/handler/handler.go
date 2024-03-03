package handler

import (
	"fmt"
	"net/http"

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
	web *http.ServeMux,
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

	web.HandleFunc("GET /auth/signin/google", h.googleSignIn())
	web.HandleFunc("GET /auth/signin/google/callback", h.googleCallback())
}

func (h *handler) googleSignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := h.googleOauth.GetSignInURL()
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

func (h *handler) googleCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		code := query.Get("code")
		state := query.Get("state")

		if state != h.config.GoogleOauth.State || code == "" {
			h.log.Warn(fmt.Sprintf("auth handler: invalid callback; state=\"%v\" code=\"%v\"", state, code))
			h.redirectWithError(w, r, validation.InvalidCallbackParam)
			return
		}

		userInfo, err := h.googleOauth.ExchangeCodeWithUserInfo(r.Context(), code)
		if err != nil {
			h.log.Error("auth handler: fail to exchange", err.Error())
			h.redirectWithError(w, r, validation.ExchangeCodeFailed)
			return
		}

		user, err := h.userService.CreateUser(r.Context(), user.NewUserArg{
			Email:     userInfo.Email,
			FirstName: userInfo.FirstName,
			LastName:  userInfo.LastName,
		})
		if err != nil {
			h.redirectWithError(w, r, validation.ServerErr)
			return
		}

		_, err = h.authService.AddUserAuthProvider(r.Context(), auth.NewAuthProviderArg{
			Provider:       "google",
			ProviderUserID: userInfo.ID,
			UserID:         user.ID,
		})
		if err != nil {
			h.redirectWithError(w, r, validation.ServerErr)
			return
		}

		jwtTokens, err := h.jwtTokens.GenereateJWTWithClaims(user.ID)
		if err != nil {
			h.log.Error("auth handler: fail to generate jwt tokens", err.Error())
			h.redirectWithError(w, r, validation.ServerErr)
			return
		}

		http.Redirect(w, r, fmt.Sprintf(
			"%v/auth?access_token=%v&expires_at=%v", h.config.FrontendURL, jwtTokens.AccessToken, jwtTokens.ExpiresAt,
		), http.StatusSeeOther)
	}
}

func (h *handler) redirectWithError(w http.ResponseWriter, r *http.Request, err string) {
	http.Redirect(w, r, h.config.FrontendURL+"/auth?err="+err, http.StatusSeeOther)
}
