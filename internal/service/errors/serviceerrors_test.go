package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceError_Error(t *testing.T) {
	t.Parallel()

	t.Run("with source error", func(t *testing.T) {
		t.Parallel()
		err := &ServiceError{
			Type:        ErrNotFound,
			Message:     "not found",
			SourceError: errors.New("original"),
		}
		require.Contains(t, err.Error(), "not found")
		require.Contains(t, err.Error(), "original")
	})

	t.Run("without source error", func(t *testing.T) {
		t.Parallel()
		err := &ServiceError{
			Type:    ErrBadRequest,
			Message: "bad request",
		}
		require.Contains(t, err.Error(), "bad request")
	})
}

func TestNotFoundError(t *testing.T) {
	t.Parallel()
	err := NotFoundError(errors.New("src"))
	var se *ServiceError
	require.True(t, errors.As(err, &se))
	assert.Equal(t, ErrNotFound, se.Type)
	assert.Equal(t, 404, se.HTTPCode)
}

func TestUnprocessableEntity(t *testing.T) {
	t.Parallel()
	err := UnprocessableEntity()
	var se *ServiceError
	require.True(t, errors.As(err, &se))
	assert.Equal(t, ErrUnprocessableEntity, se.Type)
	assert.Equal(t, 422, se.HTTPCode)
}

func TestBadRequestError(t *testing.T) {
	t.Parallel()
	err := BadRequestError()
	var se *ServiceError
	require.True(t, errors.As(err, &se))
	assert.Equal(t, ErrBadRequest, se.Type)
	assert.Equal(t, 400, se.HTTPCode)
}

func TestInternalError(t *testing.T) {
	t.Parallel()
	err := InternalError()
	var se *ServiceError
	require.True(t, errors.As(err, &se))
	assert.Equal(t, ErrInternal, se.Type)
	assert.Equal(t, 500, se.HTTPCode)
}
