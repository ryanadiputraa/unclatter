package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/ryanadiputraa/unclatter/config"
	"github.com/ryanadiputraa/unclatter/pkg/jwt"
	"github.com/ryanadiputraa/unclatter/pkg/logger"
)

type AuthMiddleware struct {
	log       logger.Logger
	config    *config.JWT
	jwtTokens jwt.JWTTokens
}

type RequestContext struct {
	UserID string
	echo.Context
}

func NewAuthMiddleware(log logger.Logger, config *config.JWT, jwtTokens jwt.JWTTokens) *AuthMiddleware {
	return &AuthMiddleware{
		log:       log,
		config:    config,
		jwtTokens: jwtTokens,
	}
}

func (m *AuthMiddleware) ParseJWTToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header
		authorization := header["Authorization"]
		if len(authorization) == 0 {
			return c.JSON(http.StatusForbidden, map[string]any{
				"message": "missing Authorization header",
			})
		}

		tokens := strings.Split(authorization[0], " ")
		if len(tokens) < 2 || tokens[0] != "Bearer" {
			return c.JSON(http.StatusForbidden, map[string]any{
				"message": "invalid Authorization header, expecting Bearer Token",
			})
		}

		claims, err := m.jwtTokens.ParseJWTClaims(tokens[1])
		if err != nil {
			return c.JSON(http.StatusForbidden, map[string]any{
				"message": err.Error(),
			})
		}

		cc := &RequestContext{
			UserID:  claims.UserID,
			Context: c,
		}

		return next(cc)
	}
}
