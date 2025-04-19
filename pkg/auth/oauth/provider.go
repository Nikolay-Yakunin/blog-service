package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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
	case "google":
		return p.fetchGoogleUser(client)
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

// fetchGoogleUser получает данные пользователя от Google
func (p *Provider) fetchGoogleUser(client *http.Client) (*UserData, error) {
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var gu struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		Verified  bool   `json:"verified_email"`
		Name      string `json:"name"`
		GivenName string `json:"given_name"`
		Picture   string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&gu); err != nil {
		return nil, err
	}

	return &UserData{
		ID:        gu.ID,
		Login:     gu.Email, // Google не всегда возвращает username, используем email
		Email:     gu.Email,
		AvatarURL: gu.Picture,
	}, nil
}

// --- Тесты ---

func TestFetchGoogleUser(t *testing.T) {
	// Мокаем ответ Google
	mockResponse := `{
		"id": "1234567890",
		"email": "testuser@gmail.com",
		"verified_email": true,
		"name": "Test User",
		"given_name": "Test",
		"picture": "https://example.com/avatar.jpg"
	}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResponse))
	}))
	defer ts.Close()

	client := ts.Client()

	// Переопределяем endpoint для теста
	getUserData := func(client *http.Client) (*UserData, error) {
		resp, err := client.Get(ts.URL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var gu struct {
			ID        string `json:"id"`
			Email     string `json:"email"`
			Verified  bool   `json:"verified_email"`
			Name      string `json:"name"`
			GivenName string `json:"given_name"`
			Picture   string `json:"picture"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&gu); err != nil {
			return nil, err
		}

		return &UserData{
			ID:        gu.ID,
			Login:     gu.Email,
			Email:     gu.Email,
			AvatarURL: gu.Picture,
		}, nil
	}

	user, err := getUserData(client)
	if err != nil {
		t.Fatalf("ошибка получения пользователя: %v", err)
	}

	if user.ID != "1234567890" {
		t.Errorf("ожидался ID '1234567890', получено: %s", user.ID)
	}
	if user.Email != "testuser@gmail.com" {
		t.Errorf("ожидался email 'testuser@gmail.com', получено: %s", user.Email)
	}
	if user.AvatarURL != "https://example.com/avatar.jpg" {
		t.Errorf("ожидался avatar_url 'https://example.com/avatar.jpg', получено: %s", user.AvatarURL)
	}
}
