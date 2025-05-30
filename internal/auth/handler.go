package auth

import (
	"net/http"
	"strings"
	"time"
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.com/Nikolay-Yakunin/blog-service/internal/users"
	"gitlab.com/Nikolay-Yakunin/blog-service/pkg/auth/jwt"
	"gitlab.com/Nikolay-Yakunin/blog-service/pkg/auth/oauth"
)

type Handler struct {
	oauth  *oauth.Provider
	users  users.Service
	config *oauth.Config
}

func NewHandler(oauthConfig *oauth.Config, userService users.Service) *Handler {
	return &Handler{
		oauth:  oauth.NewProvider(oauthConfig),
		users:  userService,
		config: oauthConfig,
	}
}

// Login инициирует процесс OAuth аутентификации
func (h *Handler) Login(c *gin.Context) {
	provider := c.Param("provider")
	var authURL string

	switch provider {
	case "github":
		authURL = h.config.Github.AuthCodeURL("state")
	case "google":
		authURL = h.config.Google.AuthCodeURL("state")
	case "vk":
		authURL = h.config.VK.AuthCodeURL("state")
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported provider"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// Убрал редирект ПОЧЕМУ НЕ ПУШИТСЯ 
// Callback обрабатывает ответ от OAuth провайдера
func (h *Handler) Callback(c *gin.Context) {
    provider := c.Param("provider")
    code := c.Query("code")

    userData, err := h.oauth.GetUserData(c.Request.Context(), provider, code)
    if err != nil {
        c.String(http.StatusInternalServerError, "OAuth error: %s", err.Error())
        return
    }

    user, err := h.users.Register(users.Provider(provider), map[string]interface{}{
        "id":         userData.ID,
        "login":      userData.Login,
        "email":      userData.Email,
        "avatar_url": userData.AvatarURL,
    })
    if err != nil {
        c.String(http.StatusInternalServerError, "User error: %s", err.Error())
        return
    }

    tokenUser := &jwt.TokenUser{
        ID:   user.ID,
        Role: user.Role,
    }
    token, err := jwt.GenerateToken(tokenUser)
    if err != nil {
        c.String(http.StatusInternalServerError, "Token error")
        return
    }

    // Формируем Set-Cookie вручную!
    cookie := fmt.Sprintf(
        "token=%s; Path=/; Max-Age=%d; HttpOnly; Secure; SameSite=None",
        token, 60*60*24,
    )
    c.Writer.Header().Add("Set-Cookie", cookie)

    // CORS
    c.Writer.Header().Set("Access-Control-Allow-Origin", "https://nikolay-yakunin.github.io")
    c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

    // Возвращаем простую HTML-страницу для popup
    c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
        <html>
          <body>
            <script>
              window.close();
            </script>
            <p>Авторизация успешна, окно можно закрыть</p>
          </body>
        </html>
    `))
}

// Logout обрабатывает выход пользователя из системы
func (h *Handler) Logout(c *gin.Context) {
	// Получаем заголовок Authorization
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "токен не предоставлен"})
		return
	}

	// Извлекаем токен из заголовка
	tokenStr := authHeader
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
	}

	// Проверяем токен
	claims, err := jwt.ValidateToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "недействительный токен"})
		return
	}

	// Добавляем токен в черный список
	if blacklist != nil {
		// Получаем время истечения токена или устанавливаем текущее время + TTL
		expiresAt := time.Now().Add(jwt.TokenTTL)
		err = blacklist.RevokeToken(claims.ID, expiresAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при выходе из системы"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "успешно вышли из системы"})
}

// RegisterRoutes регистрирует маршруты обработчика в Gin-роутере
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		// OAuth маршруты
		auth.GET("/login/:provider", h.Login)
		auth.GET("/callback/:provider", h.Callback)

		// Маршрут выхода (требует аутентификации)
		auth.POST("/logout", h.Logout)
	}
}
