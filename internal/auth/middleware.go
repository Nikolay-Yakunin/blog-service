package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/Nikolay-Yakunin/blog-service/pkg/auth/jwt"
	"gorm.io/gorm"
)

// JWTMiddleware проверяет JWT-токен в заголовке Authorization и
// добавляет ID и роль пользователя в контекст запроса
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем заголовок Authorization
		authHeader := c.GetHeader("Authorization")
		var tokenStr string

		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// Пробуем взять токен из cookie
			cookie, err := c.Cookie("token")
			if err == nil {
				tokenStr = cookie
			}
		}

		if tokenStr == "" {
			fmt.Println("JWTMiddleware: токен не найден ни в заголовке, ни в cookie")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "требуется авторизация"})
			c.Abort()
			return
		}

		// Проверяем токен
		claims, err := jwt.ValidateToken(tokenStr)
		if err != nil {
			fmt.Println("JWTMiddleware: невалидный токен:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "недействительный токен"})
			c.Abort()
			return
		}

		// Проверяем, не отозван ли токен
		if blacklist != nil && blacklist.IsRevoked(claims.ID) {
			fmt.Println("JWTMiddleware: токен отозван (ID:", claims.ID, ")")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "токен отозван"})
			c.Abort()
			return
		}

		fmt.Printf("JWTMiddleware: user_id=%v (type %T), role=%v\n", claims.UserID, claims.UserID, claims.Role)

		// Добавляем информацию о пользователе в контекст
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}

// Глобальный blacklist, инициализируемый при запуске приложения
var blacklist *TokenBlacklist

// InitTokenBlacklist инициализирует blacklist для отозванных токенов
func InitTokenBlacklist(db *gorm.DB) {
	blacklist = &TokenBlacklist{db: db}
}

// Middleware для логирования запросов
func LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// Middleware для восстановления после паники
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.RecoveryWithWriter(gin.DefaultErrorWriter)
}
