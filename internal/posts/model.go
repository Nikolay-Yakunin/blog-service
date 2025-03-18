package post

import (
	"time"

	"gitlab.com/Nikolay-Yakunin/blog-service/internal/comments"
)

// Status определяет текущее состояние поста
type Status string

const (
	// StatusDraft - пост в черновике
	StatusDraft Status = "draft"
	// StatusPublished - пост опубликован
	StatusPublished Status = "published"
	// StatusArchived - пост в архиве
	StatusArchived Status = "archived"
)

// Post представляет собой основную сущность блога
type Post struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Title       string `json:"title" gorm:"size:255;not null"`
	Slug        string `json:"slug" gorm:"uniqueIndex;size:255"`
	Description string `json:"description" gorm:"size:500"`

	// Контент
	RawContent  string `json:"raw_content" gorm:"type:text"`  // Оригинальный Markdown
	HTMLContent string `json:"html_content" gorm:"type:text"` // Отрендеренный HTML

	// Метаданные
	Status    Status   `json:"status" gorm:"type:varchar(20);default:'draft'"`
	Tags      []string `json:"tags" gorm:"type:text[]"`
	ViewCount int64    `json:"view_count" gorm:"default:0"`

	// Временные метки
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at"`

	// Связи
	AuthorID uint               `json:"author_id"`
	Comments []comments.Comment `json:"comments,omitempty" gorm:"foreignKey:PostID"`
}

// Repository описывает методы для работы с хранилищем постов
type Repository interface {
	// Create создает новый пост
	Create(post *Post) error
	// GetByID возвращает пост по его ID
	GetByID(id uint) (*Post, error)
	// Update обновляет существующий пост
	Update(post *Post) error
	// Delete удаляет пост
	Delete(id uint) error
	// List возвращает список постов с пагинацией
	List(offset, limit int) ([]Post, error)
}

// Service описывает бизнес-логику работы с постами
type Service interface {
	// CreatePost создает новый пост
	CreatePost(post *Post) error
	// GetPost получает пост по ID
	GetPost(id uint) (*Post, error)
	// UpdatePost обновляет существующий пост
	UpdatePost(post *Post) error
	// DeletePost удаляет пост
	DeletePost(id uint) error
	// ListPosts получает список постов с пагинацией
	ListPosts(offset, limit int) ([]Post, error)
}
