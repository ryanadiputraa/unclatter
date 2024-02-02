package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ryanadiputraa/unclatter/app/validation"
	"github.com/ryanadiputraa/unclatter/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	googleAPIURL = "https://www.googleapis.com"
)

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Picture   string `json:"picture"`
	Locale    string `json:"locale"`
}

type GoogleOauth interface {
	GetSignInURL() string
	ExchangeCodeWithUserInfo(ctx context.Context, code string) (*User, error)
}

type googleOauth struct {
	config       *config.GoogleOauth
	oauth2Config *oauth2.Config
}

func NewGoogleOauth(config *config.GoogleOauth) GoogleOauth {
	return &googleOauth{
		config: config,
		oauth2Config: &oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			RedirectURL:  config.RedirectURL,
			Endpoint:     google.Endpoint,
			Scopes: []string{
				fmt.Sprintf("%v/auth/userinfo.email", googleAPIURL), fmt.Sprintf("%v/auth/userinfo.profile", googleAPIURL),
			},
		},
	}
}

func (g *googleOauth) GetSignInURL() string {
	return g.oauth2Config.AuthCodeURL(g.config.State, oauth2.SetAuthURLParam("prompt", "select_account"))
}

func (g *googleOauth) ExchangeCodeWithUserInfo(ctx context.Context, code string) (*User, error) {
	token, err := g.oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(fmt.Sprintf("%v/oauth2/v2/userinfo?access_token=%v", googleAPIURL, token.AccessToken))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(validation.BadRequest)
	}

	var user *User
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return user, nil
}
