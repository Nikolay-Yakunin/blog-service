# API Спецификация

## Swagger Аннотации

Пример аннотации для эндпоинта:
```go
// @Summary Получение поста
// @Description Получение поста по ID
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "ID поста"
// @Success 200 {object} models.Post
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/posts/{id} [get]
```

## Основная спецификация API

### Аутентификация
```
GET /api/v1/auth/{provider}/login
GET /api/v1/auth/{provider}/callback
```

### Пользователи
```
GET    /api/v1/users/{id}
PUT    /api/v1/users/{id}
POST   /api/v1/users/{id}/verify
DELETE /api/v1/users/{id}
```

### Посты
```
GET    /api/v1/posts
GET    /api/v1/posts/{id}
POST   /api/v1/posts
PUT    /api/v1/posts/{id}
DELETE /api/v1/posts/{id}
```

### Комментарии
```
GET    /api/v1/posts/{postId}/comments
POST   /api/v1/posts/{postId}/comments
PUT    /api/v1/comments/{id}
DELETE /api/v1/comments/{id}
```
