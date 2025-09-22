package services

import "fmt"

// AppError - это кастомная структура ошибки для нашего приложения.
// Она позволяет сервисному слою определять HTTP-статус, который должен быть возвращен.
type AppError struct {
	StatusCode int
	Message    string
	err        error // Внутренняя (исходная) ошибка
}

// Error реализует стандартный интерфейс error.
func (e *AppError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.err)
	}
	return e.Message
}

// Unwrap позволяет использовать errors.Is и errors.As для исходной ошибки.
func (e *AppError) Unwrap() error {
	return e.err
}

// Конструкторы для стандартных типов ошибок

func NewNotFoundError(message string, err error) error {
	return &AppError{StatusCode: 404, Message: message, err: err}
}

func NewBadRequestError(message string, err error) error {
	return &AppError{StatusCode: 400, Message: message, err: err}
}

func NewUnauthorizedError(message string, err error) error {
	return &AppError{StatusCode: 401, Message: message, err: err}
}

func NewForbiddenError(message string, err error) error {
	return &AppError{StatusCode: 403, Message: message, err: err}
}

func NewConflictError(message string, err error) error {
	return &AppError{StatusCode: 409, Message: message, err: err}
}

func NewInternalServerError(message string, err error) error {
	return &AppError{StatusCode: 500, Message: message, err: err}
}
