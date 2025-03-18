package comments

import "time"

// Status определяет текущее состояние комментария
type Status string

const (
	// StatusActive - комментарий активен
	StatusActive Status = "active"
	// StatusDeleted - комментарий удален
	StatusDeleted Status = "deleted"
	// StatusHidden - комментарий скрыт модератором
	StatusHidden Status = "hidden"
)

// Comment представляет собой сущность комментария
type Comment struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Content  string `json:"content" gorm:"type:text;not null"`
	PostID   uint   `json:"post_id" gorm:"index"`
	AuthorID uint   `json:"author_id"`
	ParentID *uint  `json:"parent_id,omitempty" gorm:"index"` // Для древовидной структуры
	Status   Status `json:"status" gorm:"type:varchar(20);default:'active'"`

	// Древовидная структура
	Parent  *Comment  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Replies []Comment `json:"replies,omitempty" gorm:"foreignKey:ParentID"`

	// Метаданные
	Likes int `json:"likes" gorm:"default:0"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// Repository описывает методы для работы с хранилищем комментариев
type Repository interface {
	// Create создает новый комментарий
	Create(comment *Comment) error
	// GetByID возвращает комментарий по его ID
	GetByID(id uint) (*Comment, error)
	// GetByPostID возвращает все комментарии к посту
	GetByPostID(postID uint) ([]Comment, error)
	// Update обновляет существующий комментарий
	Update(comment *Comment) error
	// Delete удаляет комментарий
	Delete(id uint) error
}

// Service описывает бизнес-логику работы с комментариями
type Service interface {
	// CreateComment создает новый комментарий
	CreateComment(comment *Comment) error
	// GetComment получает комментарий по ID
	GetComment(id uint) (*Comment, error)
	// Обновляем сигнатуру метода, добавляя userID и userRole
	UpdateComment(comment *Comment, userID uint, userRole string) error
	// Обновляем сигнатуру метода, добавляя userID и userRole
	DeleteComment(id uint, userID uint, userRole string) error
	// GetPostComments получает все комментарии к посту
	GetPostComments(postID uint) ([]Comment, error)
}
