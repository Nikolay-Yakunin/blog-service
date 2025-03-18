package oauth

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/vk"
)

// Config contains OAuth configuration for all providers
type Config struct {
	Github *oauth2.Config
	Google *oauth2.Config
	VK     *oauth2.Config
}

// NewConfig creates OAuth configuration from environment variables
func NewConfig() *Config {
	return &Config{
		Github: &oauth2.Config{
			ClientID:     os.Getenv("OAUTH_GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("OAUTH_GITHUB_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("OAUTH_GITHUB_REDIRECT_URL"),
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
		},
		Google: &oauth2.Config{
			ClientID:     os.Getenv("OAUTH_GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("OAUTH_GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("OAUTH_GOOGLE_REDIRECT_URL"),
			Scopes:       []string{"profile", "email"},
			Endpoint:     google.Endpoint,
		},
		VK: &oauth2.Config{
			ClientID:     os.Getenv("OAUTH_VK_CLIENT_ID"),
			ClientSecret: os.Getenv("OAUTH_VK_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("OAUTH_VK_REDIRECT_URL"),
			Scopes:       []string{"email"},
			Endpoint:     vk.Endpoint,
		},
	}
}
