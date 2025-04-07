package posts

import (
	"errors"
	"time"

	"gitlab.com/Nikolay-Yakunin/blog-service/pkg/database"
	"gorm.io/gorm"
)

// PostRepository реализует интерфейс Repository и предоставляет методы
// для работы с постами в базе данных. Использует GORM как ORM и поддерживает
// все базовые CRUD операции.
type PostRepository struct {
	database.BaseRepository
}

// NewPostRepository создает новый экземпляр репозитория постов.
// Принимает инициализированное подключение к базе данных через GORM
// и возвращает указатель на PostRepository.
func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

// GetByTitle возвращает пост по его заголовку
func (r *PostRepository) GetByTitle(title string) (*Post, error) {
	var post Post
	if err := r.DB.Where("title = ?", title).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

// Возможно придумаю как реализовать поиск по токенам
// // GetByDescription возвращает пост по его описанию
// func (r *PostRepository) GetByDescription(description string) (*Post, error) {
// 	var post Post
// 	if err := r.DB.Where("description = ?", description).First(&post).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
// 	return &post, nil
// }

// GetBySlug возвращает пост по его слагу
func (r *PostRepository) GetBySlug(slug string) (*Post, error) {
	var post Post
	if err := r.DB.Where("slug = ?", slug).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

// GetByAuthor возвращает посты автора
func (r *PostRepository) GetByAuthor(authorID uint) ([]Post, error) {
	var posts []Post
	err := r.DB.Where("author_id = ?", authorID).
		Order("created_at DESC").
		Find(&posts).Error
	return posts, err
}

// GetByTag возвращает посты по тегу
func (r *PostRepository) GetByTag(tag string) ([]Post, error) {
	var posts []Post
	err := r.DB.Where("? = ANY(tags)", tag).
		Order("created_at DESC").
		Find(&posts).Error
	return posts, err
}

// GetByPublishedAt возвращает посты, опубликованные в указанный период
func (r *PostRepository) GetByPublishedAt(from, to time.Time) ([]Post, error) {
	var posts []Post
	err := r.DB.Where("published_at BETWEEN ? AND ?", from, to).
		Order("published_at DESC").
		Find(&posts).Error
	return posts, err
}

// Create создает новый пост в базе данных.
// Принимает указатель на структуру Post, которая должна содержать
// все необходимые поля. Возвращает error в случае неудачи.
// При успешном создании, пост получает ID и временные метки.
func (r *PostRepository) Create(post *Post) error {
	return r.DB.Create(post).Error
}

// GetByID возвращает пост по его идентификатору.
// Если пост не найден, возвращает (nil, nil).
// В случае других ошибок возвращает (nil, error).
func (r *PostRepository) GetByID(id uint) (*Post, error) {
	var post Post
	if err := r.DB.First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

// Update обновляет существующий пост в базе данных.
// Принимает указатель на структуру Post с обновленными данными.
// Пост должен иметь валидный ID. Возвращает error в случае неудачи.
// Автоматически обновляет временную метку updated_at.
func (r *PostRepository) Update(post *Post) error {
	return r.DB.Save(post).Error
}

// Delete выполняет мягкое удаление поста по его идентификатору.
// Запись остается в базе данных, но помечается как удаленная
// путем установки временной метки deleted_at.
// Возвращает error в случае неудачи.
func (r *PostRepository) Delete(id uint) error {
	return r.DB.Delete(&Post{}, id).Error
}

// List возвращает список постов с пагинацией.
// Принимает offset (смещение от начала) и limit (максимальное количество записей).
// Возвращает срез постов и error в случае неудачи.
// Посты сортируются по дате создания в обратном порядке (новые первые).
// Удаленные записи (soft deleted) не включаются в результат.
func (r *PostRepository) List(offset, limit int) ([]Post, error) {
	var posts []Post
	err := r.DB.Offset(offset).Limit(limit).Order("created_at DESC").Find(&posts).Error
	return posts, err
}
