// Package swagger предоставляет утилиты для работы с Swagger документацией API
package swagger

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Config содержит настройки для Swagger
type Config struct {
	Title       string
	Description string
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() *Config {
	return &Config{
		Title:       "API Documentation",
		Description: "API documentation",
		Version:     "1.0",
		Host:        "localhost:8080",
		BasePath:    "/api/v1",
		Schemes:     []string{"http", "https"},
	}
}

// RegisterRoutes регистрирует маршруты для Swagger UI
func RegisterRoutes(r *gin.Engine) {
	// Документация Swagger будет доступна по URL /swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
