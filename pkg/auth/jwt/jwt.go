package jwt

import (
	"os"
    "time"
    "github.com/golang-jwt/jwt/v5"
    "gitlab.com/Nikolay-Yakunin/blog-service/internal/users"
)

// Настройки JWT
const (
    TokenTTL = 24 * time.Hour
)

// GenerateToken создает новый JWT токен для пользователя
func GenerateToken(user *users.User) (string, error) {
    key := []byte(os.Getenv("JWT_SECRET_KEY"))

    claims := Claims{
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenTTL)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
        UserID: user.ID,
        Role:   user.Role,
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(key)
}
