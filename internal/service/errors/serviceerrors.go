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

func NotFoundError(err error) error {
	return &ServiceError{
		Type:        ErrNotFound,
		Message:     http.StatusText(http.StatusNotFound),
		HTTPCode:    http.StatusNotFound,
		SourceError: err,
	}
}

func UnprocessableEntity() error {
	return &ServiceError{
		Type:     ErrUnprocessableEntity,
		Message:  http.StatusText(http.StatusUnprocessableEntity),
		HTTPCode: http.StatusUnprocessableEntity,
	}
}

func BadRequestError() error {
	return &ServiceError{
		Type:     ErrBadRequest,
		Message:  http.StatusText(http.StatusBadRequest),
		HTTPCode: http.StatusBadRequest,
	}
}

func InternalError() error {
	return &ServiceError{
		Type:     ErrInternal,
		Message:  http.StatusText(http.StatusInternalServerError),
		HTTPCode: http.StatusInternalServerError,
	}
}
