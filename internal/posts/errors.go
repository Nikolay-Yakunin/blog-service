package posts

import "errors"

// Определение ошибок для работы с постами
var (
	// ErrEmptyTitle возвращается при попытке создать/обновить пост с пустым заголовком
	ErrEmptyTitle = errors.New("заголовок поста не может быть пустым")

	// ErrEmptyContent возвращается при попытке создать/обновить пост с пустым содержимым
	ErrEmptyContent = errors.New("содержимое поста не может быть пустым")

	// ErrPostNotFound возвращается, когда пост не найден в базе данных
	ErrPostNotFound = errors.New("пост не найден")

	// ErrInvalidStatus возвращается при попытке установить недопустимый статус поста
	ErrInvalidStatus = errors.New("недопустимый статус поста")

	// ErrUnauthorized возвращается при попытке выполнить операцию без необходимых прав
	ErrUnauthorized = errors.New("недостаточно прав для выполнения операции")
)

// ErrorResponse представляет структуру ответа с ошибкой
type ErrorResponse struct {
	Code    int    `json:"code" example:"400" swagger:"description=HTTP код ошибки"`
	Message string `json:"message" example:"Неверный формат данных" swagger:"description=Описание ошибки"`
	Details string `json:"details,omitempty" example:"Поле 'title' не может быть пустым" swagger:"description=Дополнительные детали ошибки"`
}

// NewErrorResponse создает новый экземпляр ErrorResponse
func NewErrorResponse(code int, message string, details string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	}
}
