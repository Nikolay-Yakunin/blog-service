package comments

import "errors"

var (
	ErrCommentNotFound = errors.New("comment not found")
	ErrPostNotFound    = errors.New("post not found")
	ErrUnauthorized    = errors.New("unauthorized to modify this comment")
	ErrEmptyContent    = errors.New("comment content cannot be empty")
)

// ErrorResponse представляет структуру ответа с ошибкой
type ErrorResponse struct {
	Code   int    `json:"code" example:"400" swagger:"description=HTTP код ошибки"`
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