package auth

import (
    "time"
    "gorm.io/gorm"
)

// RevokedToken хранит информацию об отозванных токенах
type RevokedToken struct {
    TokenID   string    `gorm:"primaryKey"`
    RevokedAt time.Time
    ExpiresAt time.Time
}

// TokenBlacklist для работы с отозванными токенами
type TokenBlacklist struct {
    db *gorm.DB
}

func (b *TokenBlacklist) IsRevoked(tokenID string) bool {
    var token RevokedToken
    result := b.db.First(&token, "token_id = ?", tokenID)
    return result.Error == nil
}

func (b *TokenBlacklist) RevokeToken(tokenID string, expiresAt time.Time) error {
    token := RevokedToken{
        TokenID:   tokenID,
        RevokedAt: time.Now(),
        ExpiresAt: expiresAt,
    }
    return b.db.Create(&token).Error
}

// Периодическая очистка устаревших записей
func (b *TokenBlacklist) CleanupExpired() error {
    return b.db.Delete(&RevokedToken{}, "expires_at < ?", time.Now()).Error
}
