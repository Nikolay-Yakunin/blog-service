// Package users предоставляет сервис для управления пользователями.
// Включает в себя функционал регистрации через OAuth, управления ролями и статусами пользователей.
package users

import (
	"errors"
	"time"
	"strings"
	"fmt"
)

// UserService реализует бизнес-логику работы с пользователями
type UserService struct {
    repo Repository
}

// NewUserService создает новый экземпляр сервиса пользователей
func NewUserService(repo Repository) *UserService {
    return &UserService{
        repo: repo,
    }
}

// ValidateEmail проверяет корректность email адреса
func (s *UserService) validateEmail(email string) error {
    // Простая проверка на наличие @ и домена
    if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
        return errors.New("invalid email format")
    }
    // Проверка на уникальность
    existing, err := s.repo.GetByEmail(email)
    if err != nil {
        return err
    }
    if existing != nil {
        return errors.New("email already exists")
    }
    return nil
}

// Register регистрирует нового пользователя на основе данных OAuth провайдера
// Если пользователь уже существует, возвращает существующего пользователя
//
// Параметры:
//   - provider: тип OAuth провайдера (github, google, etc.)
//   - providerData: карта с данными пользователя от провайдера
//
// Возвращает:
//   - *User: созданный или существующий пользователь
//   - error: ошибка при создании/получении пользователя
//
// Пример:
//
//	data := map[string]interface{}{
//	    "id": "12345",
//	    "login": "username",
//	    "email": "user@example.com",
//	    "avatar_url": "https://example.com/avatar.jpg",
//	}
//	user, err := service.Register(users.ProviderGithub, data)
func (s *UserService) Register(provider Provider, providerData map[string]interface{}) (*User, error) {
	// Проверяем наличие необходимых полей в providerData
	providerID, ok := providerData["id"].(string)
	if !ok {
		return nil, errors.New("missing or invalid provider ID")
	}

	login, ok := providerData["login"].(string)
	if !ok {
		return nil, errors.New("missing or invalid login")
	}

	email, ok := providerData["email"].(string)
	if !ok {
		return nil, errors.New("missing or invalid email")
	}

	avatarURL, ok := providerData["avatar_url"].(string)
	if !ok {
		avatarURL = "" // Используем пустую строку, если аватар отсутствует
	}

	// Добавляем валидацию email
    if err := s.validateEmail(email); err != nil {
        return nil, fmt.Errorf("email validation failed: %w", err)
    }

	// Ищем существующего пользователя
	existingUser, err := s.repo.GetByProviderID(provider, providerID)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return existingUser, nil
	}

	// Создаем нового пользователя
	now := time.Now()
	user := &User{
		Username:   login,
		Email:      email,
		Provider:   provider,
		ProviderID: providerID,
		Avatar:     avatarURL,
		Role:       RoleUser,
		IsActive:   true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser возвращает пользователя по его ID
//
// Возвращает nil, error если пользователь не найден
func (s *UserService) GetUser(id uint) (*User, error) {
	return s.repo.GetByID(id)
}

// UpdateUser обновляет данные существующего пользователя
//
// Возвращает error если пользователь не найден или произошла ошибка обновления
func (s *UserService) UpdateUser(user *User) error {
	return s.repo.Update(user)
}

// VerifyUser изменяет роль пользователя на RoleVerified
//
// Возвращает error если пользователь не найден или произошла ошибка обновления
func (s *UserService) VerifyUser(id uint) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	user.Role = RoleVerified
	return s.repo.Update(user)
}

// DeactivateUser отключает пользователя, устанавливая IsActive в false
//
// Возвращает error если пользователь не найден или произошла ошибка обновления
func (s *UserService) DeactivateUser(id uint) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	user.IsActive = false
	return s.repo.Update(user)
}

// UpdateLastLogin обновляет время последнего входа пользователя
//
// Возвращает error если пользователь не найден или произошла ошибка обновления
func (s *UserService) UpdateLastLogin(id uint) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	now := time.Now()
	user.LastLogin = &now
	return s.repo.Update(user)
}

// GetUsersByRole возвращает список пользователей с указанной ролью
func (s *UserService) GetUsersByRole(role Role) ([]User, error) {
    return s.repo.FindByRole(role)
}

// GetActiveUsers возвращает список активных пользователей
func (s *UserService) GetActiveUsers() ([]User, error) {
    return s.repo.FindActive()
}
