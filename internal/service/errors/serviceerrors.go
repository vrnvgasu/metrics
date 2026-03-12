package errors

import (
	"fmt"
	"net/http"
)

type ServiceErrorType string

const (
	ErrNotFound                ServiceErrorType = "not found"
	ErrInternal                ServiceErrorType = "internal server error"
	ErrUnprocessableEntity     ServiceErrorType = "unprocessable entity"
	ErrBadRequest              ServiceErrorType = "bad request"
	ErrServiceUnavailableError ServiceErrorType = "service unavailable"
)

type ServiceError struct {
	Type        ServiceErrorType
	Title       string
	Message     string
	HTTPCode    int
	SourceError error
}

func (s *ServiceError) Error() string {
	if s.SourceError != nil {
		return fmt.Sprintf("[%s] - %s: %s, error: %s", s.Type, s.Title, s.Message, s.SourceError)
	}

	return fmt.Sprintf("[%s] - %s: %s", s.Type, s.Title, s.Message)
}

func NotFoundError(err error, message string) error {
	return &ServiceError{
		Type:        ErrNotFound,
		Title:       "Not found",
		HTTPCode:    http.StatusNotFound,
		Message:     message,
		SourceError: err,
	}
}

func UnprocessableEntity(message string) error {
	return &ServiceError{
		Type:     ErrUnprocessableEntity,
		Title:    "Unprocessable entity",
		Message:  message,
		HTTPCode: http.StatusUnprocessableEntity,
	}
}

func BadRequestError(message string) error {
	return &ServiceError{
		Type:     ErrBadRequest,
		Title:    "Bad request",
		Message:  message,
		HTTPCode: http.StatusBadRequest,
	}
}

func InternalError(message string) error {
	return &ServiceError{
		Type:     ErrInternal,
		Title:    "Internal Server Error",
		Message:  message,
		HTTPCode: http.StatusInternalServerError,
	}
}
