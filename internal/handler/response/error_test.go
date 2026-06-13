package response

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	serviceerrors "github.com/vrnvgasu/metrics/internal/service/errors"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestResponseError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{
			name:           "service error not found",
			err:            serviceerrors.NotFoundError(nil),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "service error with source",
			err:            serviceerrors.NotFoundError(errors.New("original")),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "unhandled error",
			err:            errors.New("unexpected"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodGet, "/", http.NoBody)

			ResponseError(c, tt.err)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
