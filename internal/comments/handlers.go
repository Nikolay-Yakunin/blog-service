package comments

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"gitlab.com/Nikolay-Yakunin/blog-service/config"
	"gitlab.com/Nikolay-Yakunin/blog-service/pkg/middleware"
)


// Handler обрабатывает HTTP-запросы для работы с комментариями
// Содержит сервисный слой для бизнес-логики и конфигурацию приложения
type Handler struct {
	service Service    // Сервис для работы с бизнес-логикой комментариев
	config  *config.Config // Конфигурация приложения, включая JWT настройки
}

// NewHandler создает новый обработчик HTTP-запросов для комментариев
// Внедряет зависимости: сервис и конфигурацию
func NewHandler(service Service, cfg *config.Config) *Handler {
	return &Handler{
		service: service,
		config:  cfg,
	}
}

// Register регистрирует все пути обработки HTTP-запросов
// Группирует все эндпоинты под /api/v1 и защищает их middleware аутентификации
func (h *Handler) Register(router *gin.Engine) {
	// Группируем все пути под /api/v1
	comments := router.Group("/api/v1")
	// Защищаем все эндпоинты проверкой JWT токена
	comments.Use(middleware.AuthMiddleware(h.config.JWT.SecretKey)) // Добавляем middleware
	{
		// GET /api/v1/posts/:postId/comments - получение всех комментариев поста
        // POST /api/v1/posts/:postId/comments - создание нового комментария
        // PUT /api/v1/comments/:id - обновление существующего комментария
        // DELETE /api/v1/comments/:id - удаление комментария
		// Регистрируем обработчики для различных эндпоинтов
		comments.GET("/posts/:postId/comments", h.GetPostComments)
		comments.POST("/posts/:postId/comments", h.CreateComment)
		comments.PUT("/comments/:id", h.UpdateComment)
		comments.DELETE("/comments/:id", h.DeleteComment)
	}
}

// GetPostComments возвращает все комментарии для конкретного поста
// Поддерживает древовидную структуру комментариев (с вложенными ответами)
// @Security JWT
// @Summary Получить комментарии поста
// @Description Получает все комментарии для указанного поста
// @Tags comments
// @Param postId path int true "ID поста"
// @Success 200 {array} Comment
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func (h *Handler) GetPostComments(c *gin.Context) {
	// 1. Извлекаем и валидируем ID поста из URL
	postID, err := strconv.ParseUint(c.Param("postId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(
			http.StatusBadRequest,
			"Invalid post ID",
			err.Error(),
		))
		return
	}

	// 2. Получаем комментарии через сервисный слой
    // Комментарии будут отсортированы по времени создания (сначала новые)
	comments, err := h.service.GetPostComments(uint(postID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(
			http.StatusInternalServerError,
			"Failed to fetch comments",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, comments)
}

// CreateComment создает новый комментарий для поста
// Поддерживает создание как корневых комментариев, так и ответов на другие комментарии
// @Security JWT
// @Summary Создать комментарий
// @Description Создает новый комментарий для указанного поста
// @Tags comments
// @Accept json
// @Produce json
// @Param postId path int true "ID поста"
// @Param comment body Comment true "Данные комментария"
// @Success 201 {object} Comment
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func (h *Handler) CreateComment(c *gin.Context) {
	// 1. Парсим данные комментария из тела запроса
	var comment Comment
	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(
			http.StatusBadRequest,
			"Invalid comment data",
			err.Error(),
		))
		return
	}

	// 2. Устанавливаем ID автора из JWT токена
    // Это предотвращает подделку авторства комментария
	// Получаем ID пользователя из токена
	userID := c.GetUint("userID")
	comment.AuthorID = userID

	// 3. Устанавливаем ID поста из URL параметра
	postID, err := strconv.ParseUint(c.Param("postId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(
			http.StatusBadRequest,
			"Invalid post ID",
			err.Error(),
		))
		return
	}
	comment.PostID = uint(postID)

	if err := h.service.CreateComment(&comment); err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(
			http.StatusInternalServerError,
			"Failed to create comment",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// UpdateComment обновляет существующий комментарий
// Проверяет права доступа: только автор или модератор может изменить комментарий
// @Security JWT
// @Summary Обновить комментарий
// @Description Обновляет существующий комментарий
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "ID комментария"
// @Param comment body Comment true "Обновленные данные"
// @Success 200 {object} Comment
// @Failure 400,403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func (h *Handler) UpdateComment(c *gin.Context) {
	// 1. Парсим обновленные данные комментария
	var comment Comment
	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(
			http.StatusBadRequest,
			"Invalid comment data",
			err.Error(),
		))
		return
	}

	// 2. Извлекаем и валидируем ID комментария из URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(
			http.StatusBadRequest,
			"Invalid comment ID",
			err.Error(),
		))
		return
	}
	comment.ID = uint(id)

	// 3. Получаем данные пользователя из JWT токена для проверки прав
	userID := c.GetUint("userID")
	userRole := c.GetString("userRole")

	// 4. Пытаемся обновить комментарий
    // Сервис проверит права доступа (авторство или роль модератора)
	if err := h.service.UpdateComment(&comment, userID, userRole); err != nil {
		// 5. Определяем правильный статус ошибки
		status := http.StatusInternalServerError
		message := "Failed to update comment"
		
		if err == ErrUnauthorized {
			status = http.StatusForbidden // 403 для ошибок доступа
			message = "Unauthorized to modify this comment"
		}
		c.JSON(status, NewErrorResponse(
			status,
			message,
			err.Error(),
		))
		return
	}

	// 6. Получаем обновленный комментарий для ответа
    // Это гарантирует, что клиент получит актуальные данные
	// Получаем обновленный комментарий из базы
	updatedComment, err := h.service.GetComment(comment.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(
			http.StatusInternalServerError,
			"Failed to fetch updated comment",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, updatedComment)
}

// DeleteComment удаляет комментарий
// Выполняет мягкое удаление, сохраняя комментарий в базе со статусом "deleted"
// Также рекурсивно помечает удаленными все ответы на этот комментарий
// @Security JWT
// @Summary Удалить комментарий
// @Description Удаляет существующий комментарий
// @Tags comments
// @Param id path int true "ID комментария"
// @Success 204 "No Content"
// @Failure 400,403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
func (h *Handler) DeleteComment(c *gin.Context) {
	// 1. Извлекаем и валидируем ID комментария
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(
			http.StatusBadRequest,
			"Invalid comment ID",
			err.Error(),
		))
		return
	}

	// 2. Получаем данные пользователя из JWT токена
	userID := c.GetUint("userID")
	userRole := c.GetString("userRole")

	// 3. Пытаемся удалить комментарий
    // Сервис проверит права доступа и выполнит мягкое удаление
	if err := h.service.DeleteComment(uint(id), userID, userRole); err != nil {
		status := http.StatusInternalServerError
		message := "Failed to delete comment"
		
		if err == ErrUnauthorized {
			status = http.StatusForbidden
			message = "Unauthorized to delete this comment"
		}
		c.JSON(status, NewErrorResponse(
			status,
			message,
			err.Error(),
		))
		return
	}

	// 4. Возвращаем 204 No Content при успешном удалении
	c.Status(http.StatusNoContent)
}
