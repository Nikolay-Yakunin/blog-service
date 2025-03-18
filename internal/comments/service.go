package comments

import (
	"errors"
	"fmt"
)

// CommentSvc реализует бизнес-логику работы с комментариями
type CommentSvc struct {
	repo Repository
}

var (
	ErrCommentNotFound = errors.New("comment not found")
	ErrUnauthorized    = errors.New("unauthorized to modify this comment")
	ErrEmptyContent    = errors.New("comment content cannot be empty")
)

// NewCommentService создает новый экземпляр сервиса комментариев
func NewCommentService(repo Repository) Service {
	return &CommentSvc{repo: repo}
}

// CreateComment создает новый комментарий
// Проверяет наличие контента перед созданием
func (s *CommentSvc) CreateComment(comment *Comment) error {
	// Базовая валидация
	if comment.Content == "" {
		return ErrEmptyContent
	}
	return s.repo.Create(comment)
}

// GetComment получает комментарий по ID
func (s *CommentSvc) GetComment(id uint) (*Comment, error) {
	return s.repo.GetByID(id)
}

// UpdateComment обновляет существующий комментарий
func (s *CommentSvc) UpdateComment(comment *Comment, userID uint, userRole string) error {
	// Базовая валидация
	if comment.Content == "" {
		return ErrEmptyContent
	}

	existing, err := s.repo.GetByID(comment.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch comment: %w", err)
	}

	// Проверяем права на редактирование
	if !s.canModifyComment(existing.AuthorID, userID, userRole) {
		return ErrUnauthorized
	}

	existing.Content = comment.Content
	return s.repo.Update(existing)
}

// DeleteComment удаляет комментарий
func (s *CommentSvc) DeleteComment(id uint, userID uint, userRole string) error {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to fetch comment: %w", err)
	}

	if !s.canModifyComment(existing.AuthorID, userID, userRole) {
		return ErrUnauthorized
	}

	return s.repo.Delete(id)
}

// canModifyComment проверяет права на модификацию комментария
func (s *CommentSvc) canModifyComment(authorID, userID uint, userRole string) bool {
	// Автор может модифицировать свой комментарий
	if authorID == userID {
		return true
	}

	// Администратор или модератор может модифицировать любой комментарий
	if userRole == "admin" || userRole == "moderator" {
		return true
	}

	return false
}

// GetPostComments получает все комментарии для конкретного поста
func (s *CommentSvc) GetPostComments(postID uint) ([]Comment, error) {
	return s.repo.GetByPostID(postID)
}
