package swagger

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"
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

// Setup настраивает Swagger для gin роутера
func Setup(r *gin.Engine, cfg *Config) {
    if cfg == nil {
        cfg = DefaultConfig()
    }

    // Настройка информации о API
    docs.SwaggerInfo.Title = cfg.Title
    docs.SwaggerInfo.Description = cfg.Description
    docs.SwaggerInfo.Version = cfg.Version
    docs.SwaggerInfo.Host = cfg.Host
    docs.SwaggerInfo.BasePath = cfg.BasePath
    docs.SwaggerInfo.Schemes = cfg.Schemes

    // Добавление Swagger endpoint
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
