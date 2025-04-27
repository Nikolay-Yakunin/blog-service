// Package users определяет модели данных и интерфейсы для работы с пользователями системы
//
// Основные компоненты:
//   - User: модель пользователя
//   - Role: уровни доступа пользователей
//   - Provider: поддерживаемые OAuth провайдеры
//   - Repository: интерфейс хранилища
//   - Service: интерфейс бизнес-логики
package users

import "time"

// Role определяет уровень доступа пользователя в системе
// Роли образуют иерархию от гостя до администратора
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

// Provider определяет поддерживаемые системой OAuth провайдеры
// Каждый провайдер требует соответствующей настройки в конфигурации приложения
type Provider string

const (
	ProviderGithub   Provider = "github"
	ProviderGoogle   Provider = "google"
	ProviderVk       Provider = "vk"       // TODO: Нужно будет разобарться как у них создать приложение
	ProviderGitlab   Provider = "gitlab"   // Возможно добавлю в будущем
	ProviderFacebook Provider = "facebook" // Возможно добавлю в будущем
)

// User представляет собой основную модель пользователя системы
//
// Валидация полей:
//   - Username: обязательное, уникальное
//   - Email: обязательное, уникальное, валидный email
//   - Provider: обязательное, одно из определенных значений
//   - ProviderID: обязательное, уникальное в рамках провайдера
//
// Индексы:
//   - PRIMARY KEY (id)
//   - UNIQUE INDEX (username)
//   - UNIQUE INDEX (email)
//   - UNIQUE INDEX (provider, provider_id)
//   - INDEX (deleted_at)
type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"size:255;not null;unique"` // Уникальное имя пользователя
	Email    string `json:"email" gorm:"size:255;not null;unique"`    // Уникальный email пользователя

	// OAuth данные для аутентификации
	Provider   Provider `json:"provider" gorm:"size:20;not null"`     // Тип провайдера OAuth
	ProviderID string   `json:"provider_id" gorm:"size:255;not null"` // ID пользователя у провайдера

	// Данные профиля
	Avatar string `json:"avatar" gorm:"size:255"` 					  // URL аватара пользователя
	Bio    string `json:"bio" gorm:"type:text"`   					  // Описание профиля

	// Системные настройки
	Role     Role `json:"role" gorm:"type:varchar(20);default:'user'"` // Роль пользователя
	IsActive bool `json:"is_active" gorm:"default:true"`               // Статус активности

	// Временные метки
	LastLogin *time.Time `json:"last_login,omitempty"`              // Время последнего входа
	CreatedAt time.Time  `json:"created_at"`                        // Время создания записи
	UpdatedAt time.Time  `json:"updated_at"`                        // Время последнего обновления
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"` // Время удаления (soft delete)
}

// Repository описывает методы для работы с хранилищем пользователей
// Реализации должны обеспечивать потокобезопасность операций
type Repository interface {
	// Create создает нового пользователя
	Create(user *User) error
	// GetByID возвращает пользователя по ID
	GetByID(id uint) (*User, error)
	// GetByEmail возвращает пользователя по email
	GetByEmail(email string) (*User, error)
	// GetByProviderID возвращает пользователя по ID провайдера
	GetByProviderID(provider Provider, providerID string) (*User, error)
	// FindByRole возвращает список пользователей с указанной ролью
	FindByRole(role Role) ([]User, error)
	// FindActive возвращает список активных пользователей
	FindActive() ([]User, error)
	// Update обновляет данные пользователя
	Update(user *User) error
	// Delete удаляет пользователя
	Delete(id uint) error
}

// Service описывает бизнес-логику работы с пользователями
// Реализации должны производить валидацию входных данных
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
