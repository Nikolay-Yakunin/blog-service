package posts

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/Nikolay-Yakunin/blog-service/config"
	"gitlab.com/Nikolay-Yakunin/blog-service/pkg/middleware"
)

// Handler обрабатывает HTTP-запросы для работы с постами
type Handler struct {
	service Service
	config  *config.Config
}

// NewHandler создает новый обработчик HTTP-запросов для постов
func NewHandler(service Service, cfg *config.Config) *Handler {
	return &Handler{
		service: service,
		config:  cfg,
	}
}

// Register регистрирует все пути обработки HTTP-запросов
func (h *Handler) Register(router *gin.Engine) {
	posts := router.Group("/api/v1/posts")
	{
		// Публичные эндпоинты
		posts.GET("", h.ListPosts)
		posts.GET("/:id", h.GetPost)
		posts.GET("/slug/:slug", h.GetPostBySlug)

		// Защищенные эндпоинты
		authorized := posts.Use(middleware.AuthMiddleware())
		{
			authorized.POST("", h.CreatePost)
			authorized.PUT("/:id", h.UpdatePost)
			authorized.DELETE("/:id", h.DeletePost)
		}
	}
}

// ListPosts возвращает список постов с пагинацией
// @Summary Получить список постов
// @Tags posts
// @Param offset query int false "Смещение"
// @Param limit query int false "Количество записей"
// @Success 200 {array} Post
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/posts [get]
func (h *Handler) ListPosts(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	posts, err := h.service.ListPosts(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(
			http.StatusInternalServerError,
			"Failed to fetch posts",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, posts)
}

// GetPost возвращает пост по ID
// @Summary Получить пост по ID
// @Tags posts
// @Param id path int true "ID поста"
// @Success 200 {object} Post
// @Failure 404,500 {object} ErrorResponse
// @Router /api/v1/posts/{id} [get]
func (h *Handler) GetPost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(
			http.StatusBadRequest,
			"Invalid post ID",
			err.Error(),
		))
		return
	}

	post, err := h.service.GetPost(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		message := "Failed to fetch post"

		if err == ErrPostNotFound {
			status = http.StatusNotFound
			message = "Post not found"
		}

		c.JSON(status, NewErrorResponse(
			status,
			message,
			err.Error(),
		))
		return
	}

	// Увеличиваем счетчик просмотров
	go h.service.IncrementViewCount(uint(id))

	c.JSON(http.StatusOK, post)
}

// CreatePost создает новый пост
// @Security JWT
// @Summary Создать новый пост
// @Tags posts
// @Accept json
// @Produce json
// @Param post body Post true "Данные поста"
// @Success 201 {object} Post
// @Failure 400,401 {object} ErrorResponse
// @Router /api/v1/posts [post]
func (h *Handler) CreatePost(c *gin.Context) {
	var post Post
	if err := c.BindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(
			http.StatusBadRequest,
			"Invalid post data",
			err.Error(),
		))
		return
	}

	// Устанавливаем ID автора из JWT
	post.AuthorID = c.GetUint("userID")

	if err := h.service.CreatePost(&post); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(
			http.StatusBadRequest,
			"Failed to create post",
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, post)
}

// UpdatePost обновляет существующий пост
// @Security JWT
// @Summary Обновить пост
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "ID поста"
// @Param post body Post true "Данные поста"
// @Success 200 {object} Post
// @Failure 400,401,404 {object} ErrorResponse
// @Router /api/v1/posts/{id} [put]
func (h *Handler) UpdatePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(
			http.StatusBadRequest,
			"Invalid post ID",
			err.Error(),
		))
		return
	}

	var post Post
	if err := c.BindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(
			http.StatusBadRequest,
			"Invalid post data",
			err.Error(),
		))
		return
	}

	post.ID = uint(id)
	// Проверка прав (автор или админ)
	if !h.canModifyPost(c, post.AuthorID) {
		c.JSON(http.StatusForbidden, NewErrorResponse(
			http.StatusForbidden,
			"Unauthorized",
			ErrUnauthorized.Error(),
		))
		return
	}

	if err := h.service.UpdatePost(&post); err != nil {
		status := http.StatusInternalServerError
		message := "Failed to update post"

		if err == ErrPostNotFound {
			status = http.StatusNotFound
			message = "Post not found"
		}

		c.JSON(status, NewErrorResponse(
			status,
			message,
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, post)
}

// DeletePost удаляет пост
// @Security JWT
// @Summary Удалить пост
// @Tags posts
// @Param id path int true "ID поста"
// @Success 204 "No Content"
// @Failure 400,401,404 {object} ErrorResponse
// @Router /api/v1/posts/{id} [delete]
func (h *Handler) DeletePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(
			http.StatusBadRequest,
			"Invalid post ID",
			err.Error(),
		))
		return
	}

	post, err := h.service.GetPost(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		message := "Failed to fetch post"

		if err == ErrPostNotFound {
			status = http.StatusNotFound
			message = "Post not found"
		}

		c.JSON(status, NewErrorResponse(
			status,
			message,
			err.Error(),
		))
		return
	}

	// Проверка прав (автор или админ)
	if !h.canModifyPost(c, post.AuthorID) {
		c.JSON(http.StatusForbidden, NewErrorResponse(
			http.StatusForbidden,
			"Unauthorized",
			ErrUnauthorized.Error(),
		))
		return
	}

	if err := h.service.DeletePost(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(
			http.StatusInternalServerError,
			"Failed to delete post",
			err.Error(),
		))
		return
	}

	c.Status(http.StatusNoContent)
}

// GetPostByTitle возвращает пост по его заголовку
func (h *Handler) GetPostByTitle(c *gin.Context) {
	title := c.Param("title")
	post, err := h.service.GetPostByTitle(title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

// GetPostsByPublishedAt возвращает посты, опубликованные в указанный период
func (h *Handler) GetPostsByPublishedAt(c *gin.Context) {
	from, err := time.Parse(time.RFC3339, c.Query("from"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'from' date"})
		return
	}

	to, err := time.Parse(time.RFC3339, c.Query("to"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'to' date"})
		return
	}

	posts, err := h.service.GetPostsByPublishedAt(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// GetPostsByTag возвращает посты по тегу
func (h *Handler) GetPostsByTag(c *gin.Context) {
	tag := c.Param("tag")
	posts, err := h.service.GetPostsByTag(tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// GetPostBySlug возвращает пост по его slug
// @Summary Получить пост по slug
// @Tags posts
// @Param slug path string true "Slug поста"
// @Success 200 {object} Post
// @Failure 404,500 {object} ErrorResponse
// @Router /api/v1/posts/slug/{slug} [get]
func (h *Handler) GetPostBySlug(c *gin.Context) {
	slug := c.Param("slug")
	post, err := h.service.GetPostBySlug(slug)
	if err != nil {
		status := http.StatusInternalServerError
		if err == ErrPostNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// Увеличиваем счетчик просмотров
	go h.service.IncrementViewCount(post.ID)

	c.JSON(http.StatusOK, post)
}

// GetPostsByAuthor возвращает посты автора
func (h *Handler) GetPostsByAuthor(c *gin.Context) {
	userID := c.GetUint("userID")
	posts, err := h.service.GetPostsByAuthor(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// canModifyPost проверяет, может ли текущий пользователь изменять пост
func (h *Handler) canModifyPost(c *gin.Context, authorID uint) bool {
	userID := c.GetUint("userID")
	userRole := c.GetString("userRole")
	return userID == authorID || userRole == "admin"
}
