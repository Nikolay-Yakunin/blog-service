package users

import "time"

// Role определяет уровень доступа пользователя
type Role string

const (
    // RoleGuest - неавторизованный пользователь
    RoleGuest Role = "guest"
    // RoleUser - авторизованный но не верифицированный пользователь
    RoleUser Role = "user"
    // RoleVerified - верифицированный пользователь
    RoleVerified Role = "verified"
    // RoleModerator - модератор
    RoleModerator Role = "moderator"
    // RoleAdmin - администратор
    RoleAdmin Role = "admin"
)

// Provider определяет тип OAuth провайдера
type Provider string

const (
    ProviderGithub Provider = "github"
    ProviderGoogle Provider = "google"
    ProviderVk Provider = "vk"  // TODO: Нужно будет разобарться как у них создать приложение
    ProviderGitlab Provider = "gitlab"  // Возможно добавлю в будущем
    ProviderFacebook Provider = "facebook"  // Возможно добавлю в будущем
)

// User представляет собой модель пользователя
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Username  string    `json:"username" gorm:"size:255;not null;unique"`
    Email     string    `json:"email" gorm:"size:255;not null;unique"`
    
    // OAuth данные
    Provider  Provider  `json:"provider" gorm:"size:20;not null"`
    ProviderID string  `json:"provider_id" gorm:"size:255;not null"`
    
    // Профиль
    Avatar    string    `json:"avatar" gorm:"size:255"`
    Bio       string    `json:"bio" gorm:"type:text"`
    
    // Уровень доступа
    Role      Role      `json:"role" gorm:"type:varchar(20);default:'user'"`
    IsActive  bool      `json:"is_active" gorm:"default:true"`
    
    // Метаданные
    LastLogin *time.Time `json:"last_login,omitempty"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// Repository описывает методы для работы с хранилищем пользователей
type Repository interface {
    // Create создает нового пользователя
    Create(user *User) error
    // GetByID возвращает пользователя по ID
    GetByID(id uint) (*User, error)
    // GetByEmail возвращает пользователя по email
    GetByEmail(email string) (*User, error)
    // GetByProviderID возвращает пользователя по ID провайдера
    GetByProviderID(provider Provider, providerID string) (*User, error)
    // Update обновляет данные пользователя
    Update(user *User) error
    // Delete удаляет пользователя
    Delete(id uint) error
}

// Service описывает бизнес-логику работы с пользователями
type Service interface {
    // Register регистрирует нового пользователя через OAuth
    Register(provider Provider, providerData map[string]interface{}) (*User, error)
    // GetUser получает пользователя по ID
    GetUser(id uint) (*User, error)
    // UpdateUser обновляет данные пользователя
    UpdateUser(user *User) error
    // VerifyUser повышает уровень доступа пользователя до верифицированного
    VerifyUser(id uint) error
    // DeactivateUser деактивирует пользователя
    DeactivateUser(id uint) error
    // UpdateLastLogin обновляет время последнего входа
    UpdateLastLogin(id uint) error
}
