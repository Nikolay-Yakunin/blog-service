// Package database предоставляет функциональность для работы с базой данных
package database

import (
	"gorm.io/gorm"
)

// Глобальный экземпляр соединения с БД
var db *gorm.DB

// InitDB инициализирует соединение с базой данных
func InitDB(database *gorm.DB) {
	db = database
}

// GetDB возвращает глобальный экземпляр соединения с БД
func GetDB() *gorm.DB {
	return db
}
