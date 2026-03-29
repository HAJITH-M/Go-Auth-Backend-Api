package authconfig

import (
	"go-auth-backend-api/internal/config/env"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GoogleOAuthConfig() *oauth2.Config {

	return &oauth2.Config{
		ClientID:     env.AppEnv.CLIENT_ID,
		ClientSecret: env.AppEnv.CLIENT_SECRET,
		RedirectURL:  env.AppEnv.REDIRECT_URL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
