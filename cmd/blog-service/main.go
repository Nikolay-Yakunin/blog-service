package main

// @title Blog Service API
// @version 1.0
// @description API для управления блогом с поддержкой OAuth, комментариев и поиска
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT токен в формате Bearer {token}

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gitlab.com/Nikolay-Yakunin/blog-service/internal/auth"
	"gitlab.com/Nikolay-Yakunin/blog-service/internal/users"
	"gitlab.com/Nikolay-Yakunin/blog-service/pkg/auth/oauth"
	"gitlab.com/Nikolay-Yakunin/blog-service/pkg/database"
	"gitlab.com/Nikolay-Yakunin/blog-service/pkg/swagger"

	// Импорт swagger документации
	_ "gitlab.com/Nikolay-Yakunin/blog-service/docs"
)

// @Summary Check API health
// @Description Get status of the API
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"version": "1.0.0",
		"env":     os.Getenv("APP_ENV"),
	})
}

func main() {
	// Определяем, в каком окружении запускаемся
	environment := os.Getenv("APP_ENV")

	// Загружаем соответствующий .env файл
	if environment == "docker" || environment == "" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, using environment variables")
		}
	} else if environment == "local" {
		if err := godotenv.Load(".env.local"); err != nil {
			log.Println("No .env.local file found, trying default .env")
			if err := godotenv.Load(); err != nil {
				log.Println("No .env file found, using environment variables")
			}
		}
	}

	// Подключаемся к базе данных
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	log.Printf("Connecting to database with DSN: %s", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Successfully connected to database")

	// Инициализируем соединение с БД
	database.InitDB(db)

	// Миграция моделей
	log.Println("Running database migrations...")
	if err := db.AutoMigrate(&users.User{}, &auth.RevokedToken{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migrations completed successfully")

	// Инициализируем репозитории и сервисы
	userRepo := users.NewUserRepository(db)
	userService := users.NewUserService(userRepo)

	// Инициализируем OAuth конфигурацию
	oauthConfig := oauth.NewConfig()

	// Инициализируем черный список токенов
	auth.InitTokenBlacklist(db)

	// Настраиваем режим Gin в зависимости от окружения
	if environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Создаем роутер
	r := gin.Default()

	// Подключаем middleware
	r.Use(auth.LoggerMiddleware())
	r.Use(auth.RecoveryMiddleware())
	r.Use(auth.JWTMiddleware())

	// Подключаем Swagger
	swagger.RegisterRoutes(r)

	// Базовые маршруты
	r.GET("/health", healthHandler)

	// Создаем группу API маршрутов
	api := r.Group("/api/v1")

	// Регистрируем обработчики аутентификации
	authHandler := auth.NewHandler(oauthConfig, userService)
	authHandler.RegisterRoutes(api)

	// Регистрируем обработчики пользователей
	userHandler := users.NewHandler(userService)
	userHandler.RegisterRoutes(api)

	// Запускаем сервер
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s in %s mode", port, environment)
	log.Printf("Swagger documentation available at http://localhost:%s/swagger/index.html", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
