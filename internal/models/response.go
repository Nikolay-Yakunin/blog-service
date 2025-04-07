package models

// ErrorResponse представляет структуру ответа с ошибкой
// @Description Стандартная структура ответа с ошибкой
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

// SuccessResponse представляет структуру успешного ответа
// @Description Стандартная структура успешного ответа
type SuccessResponse struct {
    Code    int         `json:"code" example:"200" swagger:"description=HTTP код ответа"`
    Message string      `json:"message,omitempty" example:"Операция выполнена успешно" swagger:"description=Сообщение об успехе"`
    Data    interface{} `json:"data,omitempty" swagger:"description=Данные ответа"`
}

// NewSuccessResponse создает новый экземпляр SuccessResponse
func NewSuccessResponse(code int, message string, data interface{}) *SuccessResponse {
    return &SuccessResponse{
        Code:    code,
        Message: message,
        Data:    data,
    }
}
