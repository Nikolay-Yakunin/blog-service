# Техническое задание для клиента blog-service

## Общая информация

**Blog-service** — это RESTful API для управления блогом с поддержкой OAuth-аутентификации, комментариев, поиска и ролевой модели пользователей. Документация доступна через Swagger.

---

## Архитектура и компоненты

- **API**: REST, JSON, JWT-авторизация, OAuth (GitHub, Google, VK)
- **Основные сущности**: Пользователь, Пост, Комментарий
- **Роли пользователей**: guest, user, verified, moderator, admin
- **Документация**: Swagger (OpenAPI 2.0)
- **База данных**: PostgreSQL (структура — см. миграции)
- **Права доступа**: через JWT и роли

---

## API (эндпоинты)

### Аутентификация и OAuth

- `GET /api/v1/auth/{provider}/login` — начало OAuth-авторизации
- `GET /api/v1/auth/{provider}/callback` — callback после авторизации

### Пользователи

- `GET    /api/v1/users/{id}` — получить пользователя
- `PUT    /api/v1/users/{id}` — обновить пользователя
- `POST   /api/v1/users/{id}/verify` — верифицировать пользователя
- `DELETE /api/v1/users/{id}` — удалить пользователя

### Посты

- `GET    /api/v1/posts` — список постов (параметры: offset, limit)
- `GET    /api/v1/posts/{id}` — получить пост по ID
- `GET    /api/v1/posts/slug/{slug}` — получить пост по slug
- `POST   /api/v1/posts` — создать пост
- `PUT    /api/v1/posts/{id}` — обновить пост
- `DELETE /api/v1/posts/{id}` — удалить пост

### Комментарии

- `GET    /api/v1/posts/{postId}/comments` — получить комментарии к посту
- `POST   /api/v1/posts/{postId}/comments` — добавить комментарий к посту
- `PUT    /api/v1/comments/{id}` — обновить комментарий
- `DELETE /api/v1/comments/{id}` — удалить комментарий

### Служебные

- `GET /health` — проверка статуса API

---

## Модели данных

### User

```go
{
  "id": uint,
  "username": string,
  "email": string,
  "provider": string, // github, google, vk и др.
  "provider_id": string,
  "avatar": string,
  "bio": string,
  "role": string, // guest, user, verified, moderator, admin
  "is_active": bool,
  "last_login": string,
  "created_at": string,
  "updated_at": string,
  "deleted_at": string
}
```

### Post

```go
{
  "id": uint,
  "title": string,
  "slug": string,
  "description": string,
  "raw_content": string,   // markdown
  "html_content": string,  // html
  "status": string,        // draft, published, archived
  "tags": [string],
  "view_count": int,
  "author_id": uint,
  "created_at": string,
  "updated_at": string,
  "published_at": string,
  "comments": [Comment] // опционально
}
```

### Comment

```go
{
  "id": uint,
  "content": string,
  "post_id": uint,
  "author_id": uint,
  "parent_id": uint, // для вложенных комментариев
  "status": string,  // active, deleted, hidden
  "likes": int,
  "created_at": string,
  "updated_at": string,
  "deleted_at": string,
  "replies": [Comment] // опционально, для дерева
}
```

---

## Авторизация и безопасность

- Все защищённые эндпоинты требуют JWT-токен в заголовке:  
  `Authorization: Bearer <token>`
- Для OAuth поддерживаются: GitHub (реализовано), Google и VK (в процессе)
- Роли пользователей влияют на доступность операций (например, только автор или админ может удалять пост)

---

## Особенности реализации

- **Посты**: поддержка markdown и html, статусы, теги, счетчик просмотров, soft delete, связи с комментариями.
- **Комментарии**: древовидная структура, soft delete, модерация, лайки.
- **Пользователи**: OAuth, soft delete, роли, профили.
- **Swagger**: подробная документация, примеры запросов/ответов.
- **Миграции**: структура таблиц users, posts, comments, revoked_tokens.
- **Конфигурация**: через YAML и переменные окружения.
- **Права**: через middleware и JWT.

---

## Примеры запросов

### Получить список постов

```http
GET /api/v1/posts?offset=0&limit=10
Authorization: Bearer <token>
```

### Создать пост

```http
POST /api/v1/posts
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Новый пост",
  "raw_content": "# Пример",
  "tags": ["go", "api"]
}
```

### Получить комментарии к посту

```http
GET /api/v1/posts/1/comments
Authorization: Bearer <token>
```

---

## Дополнительно

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **OAuth**: для каждого провайдера нужны переменные окружения (см. `docs/oauth.md`)
- **Планируется**: полнотекстовый поиск, метрики, интеграция с ElasticSearch, Prometheus, Grafana.

---

## Рекомендации для клиента

1. Реализовать работу с JWT (авторизация, обновление токена).
2. Поддерживать все статусы и поля моделей (особенно для постов и комментариев).
3. Обрабатывать ошибки и коды ответов согласно Swagger.
4. Для OAuth — предусмотреть редиректы и обработку callback.
5. Для пагинации использовать параметры `offset` и `limit`.
6. Для вложенных комментариев — поддерживать дерево.
7. Для ролей — предусмотреть UI/UX для разных прав.

---

Если нужны примеры конкретных запросов, схемы или дополнительные детали по какому-либо разделу — дай знать!
