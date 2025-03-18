// Package users предоставляет реализацию хранилища пользователей в PostgreSQL
package users

import (
	"errors"
	"gorm.io/gorm"

	"gitlab.com/Nikolay-Yakunin/blog-service/pkg/database"
)

// UserRepository реализует интерфейс Repository для работы с PostgreSQL
// Предоставляет CRUD операции для сущности User
type UserRepository struct {
	database.BaseRepository
}

// NewUserRepository создает новый экземпляр репозитория пользователей
//
// Пример использования:
//
//	db := database.GetConnection()
//	repo := users.NewUserRepository(db)
func NewUserRepository(db *gorm.DB) Repository {
	return &UserRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

// Create создает нового пользователя в базе данных
//
// Возвращает error:
//   - если user == nil
//   - если произошла ошибка при создании записи в БД
//   - если нарушены ограничения уникальности
func (r *UserRepository) Create(user *User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}
	return r.DB.Create(user).Error
}

// GetByID возвращает пользователя по его ID
//
// Возвращает:
//   - (nil, nil) если пользователь не найден
//   - (nil, error) если произошла ошибка БД
//   - (*User, nil) если пользователь успешно найден
func (r *UserRepository) GetByID(id uint) (*User, error) {
	if id == 0 {
		return nil, errors.New("invalid user ID")
	}
	
	var user User
	if err := r.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*User, error) {
	var user User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByProviderID(provider Provider, providerID string) (*User, error) {
	var user User
	if err := r.DB.Where("provider = ? AND provider_id = ?", provider, providerID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *User) error {
	return r.DB.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.DB.Delete(&User{}, id).Error
}
