package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gitlab.com/Nikolay-Yakunin/blog-service/internal/users"
)

// Настройки JWT
const (
	TokenTTL = 24 * time.Hour
)

// TokenUser представляет минимальные данные пользователя для генерации токена
type TokenUser struct {
	ID   uint
	Role users.Role
}

// Генерирует случайный идентификатор для токена
func generateTokenID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateToken создает новый JWT токен для пользователя
func GenerateToken(user *TokenUser) (string, error) {
	key := []byte(os.Getenv("JWT_SECRET_KEY"))
	if len(key) == 0 {
		return "", errors.New("JWT_SECRET_KEY не настроен")
	}

	// Генерируем уникальный ID для токена
	tokenID, err := generateTokenID()
	if err != nil {
		return "", err
	}

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        tokenID,
		},
		ID:     tokenID,
		UserID: user.ID,
		Role:   user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}

// ValidateToken проверяет JWT токен и возвращает содержащиеся в нем данные
func ValidateToken(tokenString string) (*Claims, error) {
	key := []byte(os.Getenv("JWT_SECRET_KEY"))
	if len(key) == 0 {
		return nil, errors.New("JWT_SECRET_KEY не настроен")
	}

	// Создаем функцию для валидации ключа
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что используется правильный алгоритм
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неверный метод подписи")
		}
		return key, nil
	}

	// Создаем новые claims
	claims := &Claims{}

	// Парсим токен
	token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("недействительный токен")
	}

	return claims, nil
}
