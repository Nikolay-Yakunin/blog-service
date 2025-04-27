package service

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"gitlab.com/Nikolay-Yakunin/blog-service/internal/comments"
	mockRepo "gitlab.com/Nikolay-Yakunin/blog-service/internal/comments/repository/mock"
)

func TestCommentsService_CreateComment(t *testing.T) {
	repo := new(mockRepo.CommentsRepositoryMock)
	service := comments.NewCommentService(repo)

	tests := []struct {
		name    string
		comment *comments.Comment
		mockErr error
		wantErr bool
	}{
		{
			name: "Success create",
			comment: &comments.Comment{
				AuthorID: 1,
				PostID:   1,
				Content:  "Test comment",
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "Empty content",
			comment: &comments.Comment{
				AuthorID: 1,
				PostID:  1,
				Content: "",
			},
			mockErr: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.On("Create", tt.comment).Return(tt.mockErr).Maybe()
			err := service.CreateComment(tt.comment)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestCommentsService_GetThreadByPostID(t *testing.T) {
	repo := new(mockRepo.CommentsRepositoryMock)
	service := comments.NewCommentService(repo)

	mockComments := []comments.Comment{
		{
			ID:       1,
			AuthorID: 1,
			PostID:   1,
			Content:  "Parent comment",
		},
		{
			ID:       2,
			AuthorID: 2,
			PostID:   1,
			Content:  "Child comment",
			ParentID: new(uint),
		},
	}
	*mockComments[1].ParentID = 1

	tests := []struct {
		name         string
		postID      uint
		mockComments []comments.Comment
		mockErr     error
		wantErr     bool
	}{
		{
			name:         "Success get thread",
			postID:      1,
			mockComments: mockComments,
			mockErr:     nil,
			wantErr:     false,
		},
		{
			name:         "Post not found",
			postID:      999,
			mockComments: nil,
			mockErr:     comments.ErrPostNotFound,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.On("GetByPostID", tt.postID).
				Return(tt.mockComments, tt.mockErr)

			comments, err := service.GetPostComments(tt.postID)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.mockComments, comments)
		})
	}
}

