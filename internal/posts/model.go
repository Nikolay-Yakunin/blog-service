package posts

import (
	"time"

	"gitlab.com/Nikolay-Yakunin/blog-service/internal/comments"
)

// Status определяет текущее состояние поста
// @Description Статус поста
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
// @Description Пост в блоге
type Post struct {
	ID          uint   `json:"id" gorm:"primaryKey" example:"1"`
	Title       string `json:"title" gorm:"size:255;not null" example:"Как настроить Swagger в Go"`
	Slug        string `json:"slug" gorm:"uniqueIndex;size:255" example:"how-to-setup-swagger-in-go"`
	Description string `json:"description" gorm:"size:500" example:"Подробное руководство по настройке документации API с помощью Swagger в Go-приложениях"`

	// Контент
	RawContent  string `json:"raw_content" gorm:"type:text" example:"# Заголовок\n\nМаркдаун контент поста..."`        // Оригинальный Markdown
	HTMLContent string `json:"html_content" gorm:"type:text" example:"<h1>Заголовок</h1><p>HTML контент поста...</p>"` // Отрендеренный HTML

	// Метаданные
	Status    Status   `json:"status" gorm:"type:varchar(20);default:'draft'" example:"published" enums:"draft,published,archived"`
	Tags      []string `json:"tags" gorm:"type:text[]" example:"golang,swagger,api"`
	ViewCount int64    `json:"view_count" gorm:"default:0" example:"42"`

	// Временные метки
	CreatedAt   time.Time  `json:"created_at" example:"2025-01-01T00:00:00Z"`
	UpdatedAt   time.Time  `json:"updated_at" example:"2025-01-02T00:00:00Z"`
	PublishedAt *time.Time `json:"published_at" example:"2025-01-03T12:00:00Z"`

	// Связи
	AuthorID uint `json:"author_id" example:"5"`
	// Используем тип comments.CommentRef вместо comments.Comment для избежания рекурсии
	// swaggerignore: true
	Comments []comments.Comment `json:"comments,omitempty" gorm:"foreignKey:PostID" swaggerignore:"true"`
}

// PostResponse используется для ответа API с упрощенной структурой комментариев
// @Description Ответ API с постом
type PostResponse struct {
	ID          uint      `json:"id" example:"1"`
	Title       string    `json:"title" example:"Как настроить Swagger в Go"`
	Slug        string    `json:"slug" example:"how-to-setup-swagger-in-go"`
	Description string    `json:"description" example:"Подробное руководство по настройке документации API с помощью Swagger в Go-приложениях"`
	RawContent  string    `json:"raw_content" example:"# Заголовок\n\nМаркдаун контент поста..."`
	HTMLContent string    `json:"html_content" example:"<h1>Заголовок</h1><p>HTML контент поста...</p>"`
	Status      Status    `json:"status" example:"published"`
	Tags        []string  `json:"tags" example:"golang,swagger,api"`
	ViewCount   int64     `json:"view_count" example:"42"`
	CreatedAt   time.Time `json:"created_at" example:"2025-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2025-01-02T00:00:00Z"`
	AuthorID    uint      `json:"author_id" example:"5"`
	// Список ID комментариев
	CommentIDs []uint `json:"comment_ids,omitempty" example:"1,2,3"`
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
