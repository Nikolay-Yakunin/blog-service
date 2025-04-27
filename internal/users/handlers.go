// Package users предоставляет HTTP обработчики для управления пользователями.
// Включает в себя эндпоинты для получения информации о пользователях, обновления профилей,
// а также административные функции управления статусами и ролями пользователей.
package users

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler предоставляет HTTP-обработчики для управления пользователями
type Handler struct {
	userService Service
}

// NewHandler создает новый экземпляр обработчика пользователей
func NewHandler(service Service) *Handler {
	return &Handler{
		userService: service,
	}
}

// RegisterRoutes регистрирует маршруты обработчика в Gin-роутере
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		// Публичные маршруты
		users.GET("/:id", h.GetUser)

		// Требуют аутентификации
		authenticated := users.Group("/")
		authenticated.Use(AuthRequiredMiddleware())
		{
			authenticated.GET("/me", h.GetCurrentUser)
			authenticated.PUT("/me", h.UpdateCurrentUser)
		}

		// Только для администраторов
		admin := users.Group("/admin")
		admin.Use(AdminRequiredMiddleware())
		{
			admin.GET("", h.ListUsers)
			admin.PUT("/:id/verify", h.VerifyUser)
			admin.PUT("/:id/role", h.UpdateUserRole)
			admin.PUT("/:id/deactivate", h.DeactivateUser)
		}
	}
}

// GetUser возвращает информацию о пользователе по ID
func (h *Handler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат ID"})
		return
	}

	user, err := h.userService.GetUser(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "пользователь не найден"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetCurrentUser возвращает информацию о текущем авторизованном пользователе
func (h *Handler) GetCurrentUser(c *gin.Context) {
	// Получаем ID текущего пользователя из контекста (установлен middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	id, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения информации о пользователе"})
		return
	}

	user, err := h.userService.GetUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "пользователь не найден"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUserProfile обновляет данные профиля текущего пользователя
type UpdateUserProfile struct {
	Username string `json:"username"`
	Bio      string `json:"bio"`
}

// UpdateCurrentUser обновляет информацию о текущем пользователе
func (h *Handler) UpdateCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	id, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения информации о пользователе"})
		return
	}

	var profileData UpdateUserProfile
	if err := c.ShouldBindJSON(&profileData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных"})
		return
	}

	user, err := h.userService.GetUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "пользователь не найден"})
		return
	}

	// Проверяем, не занято ли имя пользователя
	if user.Username != profileData.Username {
		repo, ok := h.userService.(*UserService)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "внутренняя ошибка сервера"})
			return
		}

		existing, err := repo.repo.(*UserRepository).GetByUsername(profileData.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if existing != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "имя пользователя уже занято"})
			return
		}
	}

	// Обновляем данные пользователя
	user.Username = profileData.Username
	user.Bio = profileData.Bio

	if err := h.userService.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ListUsers возвращает список пользователей (для администраторов)
func (h *Handler) ListUsers(c *gin.Context) {
	role := c.Query("role")
	activeOnly := c.Query("active") == "true"

	var users []User
	var err error

	if role != "" {
		// Используем GetByRole, который должен быть добавлен в Service
		users, err = h.ListUsersByRole(Role(role))
	} else if activeOnly {
		// Используем GetActiveUsers, который должен быть добавлен в Service
		users, err = h.ListActiveUsers()
	} else {
		// Здесь можно добавить дополнительную логику для получения всех пользователей
		// или реализовать пагинацию
		c.JSON(http.StatusBadRequest, gin.H{"error": "необходимо указать фильтр (role или active)"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// ListUsersByRole получает список пользователей с указанной ролью
func (h *Handler) ListUsersByRole(role Role) ([]User, error) {
	userService, ok := h.userService.(*UserService)
	if !ok {
		return nil, errors.New("unable to cast to UserService")
	}

	return userService.repo.FindByRole(role)
}

// ListActiveUsers получает список активных пользователей
func (h *Handler) ListActiveUsers() ([]User, error) {
	userService, ok := h.userService.(*UserService)
	if !ok {
		return nil, errors.New("unable to cast to UserService")
	}

	return userService.repo.FindActive()
}

// VerifyUser изменяет роль пользователя на верифицированную (для администраторов)
func (h *Handler) VerifyUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат ID"})
		return
	}

	if err := h.userService.VerifyUser(uint(id)); err != nil {
		if errors.Is(err, errors.New("user not found")) {
			c.JSON(http.StatusNotFound, gin.H{"error": "пользователь не найден"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "пользователь верифицирован"})
}

// UpdateRoleRequest содержит данные для обновления роли пользователя
type UpdateRoleRequest struct {
	Role Role `json:"role" binding:"required"`
}

// UpdateUserRole обновляет роль пользователя (для администраторов)
func (h *Handler) UpdateUserRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат ID"})
		return
	}

	var roleData UpdateRoleRequest
	if err := c.ShouldBindJSON(&roleData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных"})
		return
	}

	// Чот тип валидации
	validRoles := map[Role]bool{
		RoleUser:      true,
		RoleVerified:  true,
		RoleModerator: true,
		RoleAdmin:     true,
	}

	if !validRoles[roleData.Role] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимая роль"})
		return
	}

	user, err := h.userService.GetUser(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "пользователь не найден"})
		return
	}

	user.Role = roleData.Role
	if err := h.userService.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "роль обновлена", "user": user})
}

// DeactivateUser деактивирует учетную запись пользователя (для администраторов)
func (h *Handler) DeactivateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат ID"})
		return
	}

	if err := h.userService.DeactivateUser(uint(id)); err != nil {
		if errors.Is(err, errors.New("user not found")) {
			c.JSON(http.StatusNotFound, gin.H{"error": "пользователь не найден"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "пользователь деактивирован"})
}

// Middleware для проверки наличия аутентификации
func AuthRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "требуется авторизация"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Middleware для проверки наличия прав администратора
func AdminRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "требуется авторизация"})
			c.Abort()
			return
		}

		userRole, ok := role.(Role)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка получения роли пользователя"})
			c.Abort()
			return
		}

		if userRole != RoleAdmin && userRole != RoleModerator {
			c.JSON(http.StatusForbidden, gin.H{"error": "недостаточно прав"})
			c.Abort()
			return
		}

		c.Next()
	}
}
