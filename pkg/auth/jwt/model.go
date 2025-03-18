package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"gitlab.com/Nikolay-Yakunin/blog-service/internal/users"
)

// Claims определяет данные, хранимые в JWT токене
type Claims struct {
    jwt.RegisteredClaims
    UserID uint        `json:"user_id"`
    Role   users.Role  `json:"role"`
}