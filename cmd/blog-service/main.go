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

	"gitlab.com/Nikolay-Yakunin/blog-service/config"
	"gitlab.com/Nikolay-Yakunin/blog-service/internal/auth"
	"gitlab.com/Nikolay-Yakunin/blog-service/internal/comments"
	"gitlab.com/Nikolay-Yakunin/blog-service/internal/posts"
	"gitlab.com/Nikolay-Yakunin/blog-service/internal/users"
	"gitlab.com/Nikolay-Yakunin/blog-service/pkg/auth/oauth"
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

	// Загружаем конфигурацию из файла и переменных окружения
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Подключаемся к базе данных
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}
	log.Printf("Connecting to database with DSN: %s", dsn)
	db, dbErr := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if dbErr != nil {
		log.Fatal("Failed to connect to database:", dbErr)
	}
	log.Println("Successfully connected to database")

	// Инициализируем соединение с БД (если пакет database это делает)
	// database.InitDB(db) // Закомментировано, т.к. db передается напрямую

	// Миграции - используем make migrate-up
	/* log.Println("Running database migrations...")
	if err := db.AutoMigrate(&users.User{}, &auth.RevokedToken{}, &posts.Post{}, &comments.Comment{}); err != nil { // Добавлены Post и Comment на всякий случай, но строка закомментирована
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migrations completed successfully") */

	// Инициализируем репозитории и сервисы
	userRepo := users.NewUserRepository(db)
	userService := users.NewUserService(userRepo)
	postRepo := posts.NewPostRepository(db)
	postService := posts.NewPostService(postRepo)
	commentRepo := comments.NewCommentRepository(db)
	commentService := comments.NewCommentService(commentRepo)

	// Инициализируем OAuth конфигурацию (возвращаем старый способ)
	oauthConfig := oauth.NewConfig()

	// Инициализируем черный список токенов
	auth.InitTokenBlacklist(db)

	// Настраиваем режим Gin
	if cfg.App.Name == "production" { // Сверяемся с полем в конфиге
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// Подключаем стандартные middleware
	r.Use(auth.LoggerMiddleware())   // Используем из internal/auth
	r.Use(auth.RecoveryMiddleware()) // Используем из internal/auth

	// Подключаем Swagger
	swagger.RegisterRoutes(r)

	// Базовые маршруты
	r.GET("/health", healthHandler)

	// Создаем группу API v1
	apiV1 := r.Group("/api/v1")

	// Регистрируем обработчики, передавая группу apiV1

	// Auth
	authHandler := auth.NewHandler(oauthConfig, userService)
	authHandler.RegisterRoutes(apiV1)

	// Users
	userHandler := users.NewHandler(userService)
	userHandler.RegisterRoutes(apiV1) // Используем единый метод

	// Posts
	postHandler := posts.NewHandler(postService, cfg) // Передаем cfg
	postHandler.Register(r)                           // Используем существующий метод Register(*gin.Engine)

	// Comments
	commentHandler := comments.NewHandler(commentService, cfg) // Передаем cfg
	commentHandler.Register(r)                                 // Используем существующий метод Register(*gin.Engine)

	// Запускаем сервер
	port := cfg.Server.Port
	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s in %s mode", port, cfg.App.Name)
	log.Printf("Swagger documentation available at http://localhost:%s/swagger/index.html", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
