package auth

import "github.com/golang-jwt/jwt/v5"

type User struct {
    ID       uint   `json:"id" gorm:"primaryKey"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Role     string `json:"role"`
}

type Claims struct {
    jwt.RegisteredClaims
    UserID uint   `json:"user_id"`
    Role   string `json:"role"`
}