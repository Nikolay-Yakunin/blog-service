package users

import (
	"errors"
	"gorm.io/gorm"

	"gitlab.com/Nikolay-Yakunin/blog-service/pkg/database"
)

// UserRepository реализует интерфейс Repository для работы с БД
type UserRepository struct {
	database.BaseRepository
}

func NewUserRepository(db *gorm.DB) Repository {
	return &UserRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

// Create создает нового пользователя в базе данных
func (r *UserRepository) Create(user *User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}
	return r.DB.Create(user).Error
}

// GetByID возвращает пользователя по его ID
// Возвращает nil, nil если пользователь не найден
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
