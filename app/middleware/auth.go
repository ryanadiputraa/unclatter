package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ryanadiputraa/unclatter/config"
	_http "github.com/ryanadiputraa/unclatter/pkg/http"
	"github.com/ryanadiputraa/unclatter/pkg/jwt"
	"github.com/ryanadiputraa/unclatter/pkg/logger"
)

type AuthMiddleware struct {
	log       logger.Logger
	config    *config.JWT
	rw        _http.ResponseWriter
	jwtTokens jwt.JWTTokens
}

type AuthContext struct {
	UserID string
	context.Context
}

func NewAuthMiddleware(log logger.Logger, config *config.JWT, rw _http.ResponseWriter, jwtTokens jwt.JWTTokens) *AuthMiddleware {
	return &AuthMiddleware{
		log:       log,
		config:    config,
		rw:        rw,
		jwtTokens: jwtTokens,
	}
}

func (m *AuthMiddleware) ParseJWTToken(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if len(authorization) == 0 {
			m.rw.WriteErrMessage(w, http.StatusUnauthorized, "missing Authorization header")
			return
		}

		tokens := strings.Split(authorization, " ")
		if len(tokens) < 2 || tokens[0] != "Bearer" {
			m.rw.WriteErrMessage(w, http.StatusForbidden, "invalid Authorization header, expecting Bearer Token")
			return
		}

		claims, err := m.jwtTokens.ParseJWTClaims(tokens[1])
		if err != nil {
			m.rw.WriteErrMessage(w, http.StatusForbidden, err.Error())
			return
		}

		ac := &AuthContext{
			UserID:  claims.UserID,
			Context: r.Context(),
		}
		rc := r.WithContext(ac)

		next(w, rc)
	})
}
