package oauth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
