package users

import (
	"errors"
	"time"
)

type userService struct {
	repo Repository
}

func NewUserService(repo Repository) Service {
	return &userService{repo: repo}
}

// Register регистрирует или возвращает существующего пользователя на основе OAuth данных
func (s *userService) Register(provider Provider, providerData map[string]interface{}) (*User, error) {
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

func (s *userService) GetUser(id uint) (*User, error) {
	return s.repo.GetByID(id)
}

func (s *userService) UpdateUser(user *User) error {
	return s.repo.Update(user)
}

func (s *userService) VerifyUser(id uint) error {
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

func (s *userService) DeactivateUser(id uint) error {
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

func (s *userService) UpdateLastLogin(id uint) error {
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
