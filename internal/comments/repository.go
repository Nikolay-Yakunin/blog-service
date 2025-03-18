package comments

import (
	"gorm.io/gorm"

	"gitlab.com/Nikolay-Yakunin/blog-service/pkg/database"
)

// CommentRepository реализует интерфейс Repository для работы с БД
type CommentRepository struct {
	database.BaseRepository
}

// NewCommentRepository создает новый экземпляр репозитория комментариев
func NewCommentRepository(db *gorm.DB) Repository {
	return &CommentRepository{
		BaseRepository: database.NewBaseRepository(db),
	}
}

// Create сохраняет новый комментарий в базу данных
func (r *CommentRepository) Create(comment *Comment) error {
	return r.DB.Create(comment).Error
}

// GetByID получает комментарий по ID вместе с вложенными ответами
// Использует GORM Preload для загрузки связанных данных
func (r *CommentRepository) GetByID(id uint) (*Comment, error) {
	var comment Comment
	if err := r.DB.Preload("Replies").First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetByPostID получает все корневые комментарии для поста с полной загрузкой вложенности
func (r *CommentRepository) GetByPostID(postID uint) ([]Comment, error) {
	var comments []Comment

	// Выбираем только корневые комментарии (parent_id IS NULL)
	err := r.DB.Where("post_id = ? AND parent_id IS NULL", postID).
		Preload("Replies.Replies.Replies"). // Загружаем до 3 уровней вложенности
		Order("created_at DESC").
		Find(&comments).Error

	if err != nil {
		return nil, err // Added error check
	}

	return comments, nil
}

// Update обновляет существующий комментарий
func (r *CommentRepository) Update(comment *Comment) error {
	return r.DB.Save(comment).Error
}

// Delete выполняет мягкое удаление комментария и всех его ответов
func (r *CommentRepository) Delete(id uint) error {
	// Начинаем транзакцию
	tx := r.DB.Begin()

	// Удаляем основной комментарий
	if err := tx.Model(&Comment{}).
		Where("id = ?", id).
		Update("status", StatusDeleted).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Рекурсивно удаляем все ответы
	if err := r.recursiveDeleteReplies(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// recursiveDeleteReplies рекурсивно удаляет все ответы на комментарий
func (r *CommentRepository) recursiveDeleteReplies(tx *gorm.DB, parentID uint) error {
	// Находим все прямые ответы
	var replies []Comment
	if err := tx.Where("parent_id = ?", parentID).Find(&replies).Error; err != nil {
		return err
	}

	// Для каждого ответа
	for _, reply := range replies {
		// Удаляем его
		if err := tx.Model(&Comment{}).
			Where("id = ?", reply.ID).
			Update("status", StatusDeleted).Error; err != nil {
			return err
		}

		// Рекурсивно удаляем его ответы
		if err := r.recursiveDeleteReplies(tx, reply.ID); err != nil {
			return err
		}
	}

	return nil
}
