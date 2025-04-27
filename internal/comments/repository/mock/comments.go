package mock

import (
	"github.com/stretchr/testify/mock"
	"gitlab.com/Nikolay-Yakunin/blog-service/internal/comments"
)

type CommentsRepositoryMock struct {
	mock.Mock
}

func (r *CommentsRepositoryMock) Create(comment *comments.Comment) error {
	args := r.Called(comment)
	return args.Error(0)
}

func (r *CommentsRepositoryMock) GetByID(id uint) (*comments.Comment, error) {
	args := r.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*comments.Comment), args.Error(1)
}

func (r *CommentsRepositoryMock) GetByPostID(postID uint) ([]comments.Comment, error) {
	args := r.Called(postID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]comments.Comment), args.Error(1)
}

func (r *CommentsRepositoryMock) Update(comment *comments.Comment) error {
	args := r.Called(comment)
	return args.Error(0)
}

func (r *CommentsRepositoryMock) Delete(id uint) error {
	args := r.Called(id)
	return args.Error(0)
}
