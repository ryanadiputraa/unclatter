package oauth

import (
	"github.com/ryanadiputraa/unclatter/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOauth interface {
	GetSignInURL() string
}

type googleOauth struct {
	config *config.GoogleOauth
}

func NewGoogleOauth(config *config.GoogleOauth) GoogleOauth {
	return &googleOauth{
		config: config,
	}
}

func (g *googleOauth) GetSignInURL() string {
	config := &oauth2.Config{
		ClientID:     g.config.ClientID,
		ClientSecret: g.config.ClientSecret,
		RedirectURL:  g.config.RedirectURL,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	}
	return config.AuthCodeURL(g.config.State, oauth2.SetAuthURLParam("prompt", "select_account"))
}
