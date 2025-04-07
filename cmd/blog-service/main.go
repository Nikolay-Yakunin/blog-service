package main

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
)

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

	// Базовые маршруты
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"version": "1.0.0",
			"env":     environment,
		})
	})

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
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
