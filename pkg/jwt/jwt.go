package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ryanadiputraa/unclatter/config"
)

const expiresAt = time.Hour * 24

type JWT struct {
	AccessToken string `json:"access_token"`
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type JWTTokens interface {
	GenereateJWTWithClaims(userID string) (*JWT, error)
}

type jwtTokens struct {
	config *config.JWT
}

func NewJWTTokens(config *config.JWT) JWTTokens {
	return &jwtTokens{
		config: config,
	}
}

func (j *jwtTokens) GenereateJWTWithClaims(userID string) (*JWT, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresAt)),
		},
	})
	accessToken, err := token.SignedString([]byte(j.config.Secret))
	if err != nil {
		return nil, err
	}

	return &JWT{AccessToken: accessToken}, nil
}
