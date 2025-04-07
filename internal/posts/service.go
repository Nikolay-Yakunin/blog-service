package posts

import (
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

// PostService реализует бизнес-логику работы с постами
type PostService struct {
	repo Repository
}

// NewPostService создает новый экземпляр сервиса постов
func NewPostService(repo Repository) *PostService {
	return &PostService{
		repo: repo,
	}
}

// CreatePost создает новый пост
func (s *PostService) CreatePost(post *Post) error {
	// Валидация
	if err := s.validatePost(post); err != nil {
		return err
	}

	// Генерация слага из заголовка
	post.Slug = slug.Make(post.Title)

	// Рендеринг HTML из Markdown
	post.HTMLContent = s.renderHTML(post.RawContent)

	// Установка начальных значений
	post.Status = StatusDraft
	post.ViewCount = 0
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	return s.repo.Create(post)
}

// GetPost получает пост по ID
func (s *PostService) GetPost(id uint) (*Post, error) {
	post, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, ErrPostNotFound
	}
	return post, nil
}

// UpdatePost обновляет существующий пост
func (s *PostService) UpdatePost(post *Post) error {
	// Валидация
	if err := s.validatePost(post); err != nil {
		return err
	}

	// Проверяем существование поста
	existing, err := s.repo.GetByID(post.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrPostNotFound
	}

	// Обновляем HTML контент если изменился Markdown
	if post.RawContent != existing.RawContent {
		post.HTMLContent = s.renderHTML(post.RawContent)
	}

	// Обновляем слаг если изменился заголовок
	if post.Title != existing.Title {
		post.Slug = slug.Make(post.Title)
	}

	// Если пост публикуется впервые
	if post.Status == StatusPublished && existing.Status != StatusPublished {
		now := time.Now()
		post.PublishedAt = &now
	}

	post.UpdatedAt = time.Now()
	return s.repo.Update(post)
}

// DeletePost удаляет пост
func (s *PostService) DeletePost(id uint) error {
	// Проверяем существование поста
	if _, err := s.GetPost(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

// ListPosts получает список постов с пагинацией
func (s *PostService) ListPosts(offset, limit int) ([]Post, error) {
	return s.repo.List(offset, limit)
}

// GetPostByTitle получает пост по его заголовку
func (s *PostService) GetPostByTitle(title string) (*Post, error) {
	post, err := s.repo.GetByTitle(title)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, ErrPostNotFound
	}
	return post, nil
}

// GetPostBySlug получает пост по его слагу
func (s *PostService) GetPostBySlug(slug string) (*Post, error) {
	post, err := s.repo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, ErrPostNotFound
	}
	return post, nil
}

// GetPostsByAuthor возвращает посты автора
func (s *PostService) GetPostsByAuthor(authorID uint) ([]Post, error) {
	post, err := s.repo.GetByAuthor(authorID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, ErrPostNotFound
	}
	return post, nil
}

// GetPostsByTag возвращает посты по тегу
func (s *PostService) GetPostsByTag(tag string) ([]Post, error) {
	posts, err := s.repo.GetByTag(tag)
	if err != nil {
		return nil, err
	}
	if posts == nil {
		return nil, ErrPostNotFound
	}
	return posts, nil
}

// GetPostsByPublishedAt возвращает посты, опубликованные в указанный период
func (s *PostService) GetPostsByPublishedAt(from, to time.Time) ([]Post, error) {
	posts, err := s.repo.GetByPublishedAt(from, to)
	if err != nil {
		return nil, err
	}
	if posts == nil {
		return nil, ErrPostNotFound
	}
	return posts, nil
}

// IncrementViewCount увеличивает счетчик просмотров поста
func (s *PostService) IncrementViewCount(id uint) error {
	post, err := s.GetPost(id)
	if err != nil {
		return err
	}
	post.ViewCount++
	return s.repo.Update(post)
}

// validatePost проверяет корректность данных поста
func (s *PostService) validatePost(post *Post) error {
	if strings.TrimSpace(post.Title) == "" {
		return ErrEmptyTitle
	}
	if strings.TrimSpace(post.RawContent) == "" {
		return ErrEmptyContent
	}
	if post.Status != "" && 
	   post.Status != StatusDraft && 
	   post.Status != StatusPublished && 
	   post.Status != StatusArchived {
		return ErrInvalidStatus
	}
	return nil
}

// renderHTML конвертирует Markdown в HTML с санитизацией
func (s *PostService) renderHTML(markdown string) string {
	// Конвертируем Markdown в HTML
	unsafe := blackfriday.Run([]byte(markdown))
	
	// Санитизируем HTML
	p := bluemonday.UGCPolicy()
	return string(p.SanitizeBytes(unsafe))
}
