package jwt

import (
	"fmt"
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
	ParseJWTClaims(accessToken string) (*Claims, error)
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

func (j *jwtTokens) ParseJWTClaims(accessToken string) (*Claims, error) {
	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method: %v", t.Header["alg"])
		}
		return []byte(j.config.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}
