package auth

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "gitlab.com/Nikolay-Yakunin/blog-service/pkg/auth/oauth"
    "gitlab.com/Nikolay-Yakunin/blog-service/internal/users"
    "gitlab.com/Nikolay-Yakunin/blog-service/pkg/auth/jwt"
)

type Handler struct {
    oauth   *oauth.Provider
    users   users.Service
    config  *oauth.Config
}

func NewHandler(oauthConfig *oauth.Config, userService users.Service) *Handler {
    return &Handler{
        oauth:   oauth.NewProvider(oauthConfig),
        users:   userService,
        config:  oauthConfig,
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

// Callback обрабатывает ответ от OAuth провайдера
func (h *Handler) Callback(c *gin.Context) {
    provider := c.Param("provider")
    code := c.Query("code")
    
    userData, err := h.oauth.GetUserData(c.Request.Context(), provider, code)
    if (err != nil) {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    user, err := h.users.Register(users.Provider(provider), map[string]interface{}{
        "id":         userData.ID,
        "login":      userData.Login,
        "email":      userData.Email,
        "avatar_url": userData.AvatarURL,
    })
    if (err != nil) {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Генерируем JWT токен
    token, err := jwt.GenerateToken(user)
    if (err != nil) {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка генерации токена"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "token": token,
        "user":  user,
    })
}
