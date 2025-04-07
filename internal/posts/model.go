package posts

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
	// GetByTitle возвращает пост по его заголовку
	GetByTitle(title string) (*Post, error)
	// GetBySlug возвращает пост по его слагу
	GetBySlug(slug string) (*Post, error)
	// GetByAuthor возвращает посты автора
	GetByAuthor(authorID uint) ([]Post, error)
	// GetByTag возвращает посты по тегу
	GetByTag(tag string) ([]Post, error)
	// GetByID возвращает пост по его ID
	GetByID(id uint) (*Post, error)
	// GetByPublishedAt возвращает посты, опубликованные в указанный период
	GetByPublishedAt(from, to time.Time) ([]Post, error)
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
	// GetPostByTitle получает пост по его заголовку
	GetPostByTitle(title string) (*Post, error)
	// GetPostBySlug получает пост по его слагу
	GetPostBySlug(slug string) (*Post, error)
	// GetPostsByAuthor возвращает посты автора
	GetPostsByAuthor(authorID uint) ([]Post, error)
	// GetPostsByTag возвращает посты по тегу
	GetPostsByTag(tag string) ([]Post, error)
	// GetPostsByPublishedAt возвращает посты, опубликованные в указанный период
	GetPostsByPublishedAt(from, to time.Time) ([]Post, error)
	// UpdatePost обновляет существующий пост
	UpdatePost(post *Post) error
	// IncrementViewCount увеличивает счетчик просмотров поста
	IncrementViewCount(id uint) error
	// DeletePost удаляет пост
	DeletePost(id uint) error
	// ListPosts получает список постов с пагинацией
	ListPosts(offset, limit int) ([]Post, error)
}
