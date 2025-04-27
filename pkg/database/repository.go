package database

import "gorm.io/gorm"

// BaseRepository предоставляет базовую имплементацию для работы с БД
type BaseRepository struct {
	DB *gorm.DB
}

// NewBaseRepository создает новый экземпляр базового репозитория
func NewBaseRepository(db *gorm.DB) BaseRepository {
	return BaseRepository{DB: db}
}
