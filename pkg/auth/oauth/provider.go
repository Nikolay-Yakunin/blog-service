package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

// UserData представляет нормализованные данные пользователя от OAuth провайдеров
type UserData struct {
	ID        string
	Login     string
	Email     string
	AvatarURL string
	Provider  string
}

// Provider обрабатывает процесс аутентификации через OAuth
type Provider struct {
	config *Config
}

func NewProvider(config *Config) *Provider {
	return &Provider{config: config}
}

// GetUserData получает информацию о пользователе от OAuth провайдера
func (p *Provider) GetUserData(ctx context.Context, provider string, code string) (*UserData, error) {
	var config *oauth2.Config
	switch provider {
	case "github":
		config = p.config.Github
	case "google":
		config = p.config.Google
	case "vk":
		config = p.config.VK
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %w", err)
	}

	client := config.Client(ctx, token)
	userData, err := p.fetchUserData(client, provider)
	if err != nil {
		return nil, err
	}
	userData.Provider = provider

	return userData, nil
}

// fetchUserData получает данные пользователя в зависимости от провайдера
func (p *Provider) fetchUserData(client *http.Client, provider string) (*UserData, error) {
	switch provider {
	case "github":
		return p.fetchGithubUser(client)
	// case "google":	// Uncomment when implementing Google OAuth	TODO: Write fetchGoogleUser
	//     return p.fetchGoogleUser(client)
	// case "vk":	// Uncomment when implementing VK OAuth TODO: Write fetchVKUser
	//     return p.fetchVKUser(client)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// fetchGithubUser получает данные пользователя от GitHub
func (p *Provider) fetchGithubUser(client *http.Client) (*UserData, error) {
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var gh struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&gh); err != nil {
		return nil, err
	}

	return &UserData{
		ID:        fmt.Sprintf("%d", gh.ID),
		Login:     gh.Login,
		Email:     gh.Email,
		AvatarURL: gh.AvatarURL,
	}, nil
}

// Similar implementations for Google and VK...
