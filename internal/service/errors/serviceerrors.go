package errors

import (
	"fmt"
	"net/http"
)

type ServiceErrorType string

const (
	ErrNotFound            ServiceErrorType = "not found"
	ErrInternal            ServiceErrorType = "internal server error"
	ErrUnprocessableEntity ServiceErrorType = "unprocessable entity"
	ErrBadRequest          ServiceErrorType = "bad request"
)

type ServiceError struct {
	Type        ServiceErrorType
	Title       string
	Message     string
	HttpCode    int
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
		HttpCode:    http.StatusNotFound,
		Message:     message,
		SourceError: err,
	}
}

func UnprocessableEntity(message string) error {
	return &ServiceError{
		Type:     ErrUnprocessableEntity,
		Title:    "Unprocessable entity",
		Message:  message,
		HttpCode: http.StatusUnprocessableEntity,
	}
}

func BadRequestError(message string) error { // todo: rename to BadRequest
	return &ServiceError{
		Type:     ErrBadRequest,
		Title:    "Bad request",
		Message:  message,
		HttpCode: http.StatusBadRequest,
	}
}
