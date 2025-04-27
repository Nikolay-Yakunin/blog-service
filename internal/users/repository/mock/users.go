package mock

import (
	"github.com/stretchr/testify/mock"
	"gitlab.com/Nikolay-Yakunin/blog-service/internal/users"
)

type UsersRepositoryMock struct {
	mock.Mock
}

// FindActive implements users.Repository.
func (r *UsersRepositoryMock) FindActive() ([]users.User, error) {
	panic("unimplemented")
}

// FindByRole implements users.Repository.
func (r *UsersRepositoryMock) FindByRole(role users.Role) ([]users.User, error) {
	panic("unimplemented")
}

// GetByEmail implements users.Repository.
func (r *UsersRepositoryMock) GetByEmail(email string) (*users.User, error) {
	panic("unimplemented")
}

// GetByProviderID implements users.Repository.
func (r *UsersRepositoryMock) GetByProviderID(provider users.Provider, providerID string) (*users.User, error) {
	panic("unimplemented")
}

func (r *UsersRepositoryMock) Create(user *users.User) error {
	args := r.Called(user)
	return args.Error(0)
}

func (r *UsersRepositoryMock) GetByID(id uint) (*users.User, error) {
	args := r.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.User), args.Error(1)
}

func (r *UsersRepositoryMock) GetByOAuth(provider, providerID string) (*users.User, error) {
	args := r.Called(provider, providerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.User), args.Error(1)
}

func (r *UsersRepositoryMock) Update(user *users.User) error {
	args := r.Called(user)
	return args.Error(0)
}

func (r *UsersRepositoryMock) Delete(id uint) error {
	args := r.Called(id)
	return args.Error(0)
}
