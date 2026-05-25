package errors

import (
	"fmt"
	"net/http"
)

// ServiceErrorType — строковый тип ошибки сервиса.
type ServiceErrorType string

const (
	ErrNotFound                ServiceErrorType = "not found"
	ErrInternal                ServiceErrorType = "internal server error"
	ErrUnprocessableEntity     ServiceErrorType = "unprocessable entity"
	ErrBadRequest              ServiceErrorType = "bad request"
	ErrServiceUnavailableError ServiceErrorType = "service unavailable"
)

// ServiceError — структурированная ошибка сервисного слоя с HTTP-кодом.
type ServiceError struct {
	Type        ServiceErrorType
	Message     string
	HTTPCode    int
	SourceError error
}

func (s *ServiceError) Error() string {
	if s.SourceError != nil {
		return fmt.Sprintf("[%s]: %s, error: %s", s.Type, s.Message, s.SourceError)
	}

	return fmt.Sprintf("[%s]: %s", s.Type, s.Message)
}

// NotFoundError возвращает ошибку 404.
func NotFoundError(err error) error {
	return &ServiceError{
		Type:        ErrNotFound,
		Message:     http.StatusText(http.StatusNotFound),
		HTTPCode:    http.StatusNotFound,
		SourceError: err,
	}
}

// UnprocessableEntity возвращает ошибку 422.
func UnprocessableEntity() error {
	return &ServiceError{
		Type:     ErrUnprocessableEntity,
		Message:  http.StatusText(http.StatusUnprocessableEntity),
		HTTPCode: http.StatusUnprocessableEntity,
	}
}

// BadRequestError возвращает ошибку 400.
func BadRequestError() error {
	return &ServiceError{
		Type:     ErrBadRequest,
		Message:  http.StatusText(http.StatusBadRequest),
		HTTPCode: http.StatusBadRequest,
	}
}

// InternalError возвращает ошибку 500.
func InternalError() error {
	return &ServiceError{
		Type:     ErrInternal,
		Message:  http.StatusText(http.StatusInternalServerError),
		HTTPCode: http.StatusInternalServerError,
	}
}
