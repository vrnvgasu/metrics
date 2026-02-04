package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vrnvgasu/metrics/internal/repository"
)

func TestUpdate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		method              string
		path                string
		contentType         string
		expectedStatus      int
		expectedContentType string
	}{
		{
			name:                "success: gauge 0",
			method:              http.MethodPost,
			path:                "/update/gauge/test/0",
			contentType:         "text/plain",
			expectedStatus:      http.StatusOK,
			expectedContentType: "text/plain",
		},
		{
			name:                "success: gauge 0.1",
			method:              http.MethodPost,
			path:                "/update/gauge/test/0.1",
			contentType:         "text/plain",
			expectedStatus:      http.StatusOK,
			expectedContentType: "text/plain",
		},
		{
			name:                "success: gauge 1.1",
			method:              http.MethodPost,
			path:                "/update/gauge/test/1.1",
			contentType:         "text/plain",
			expectedStatus:      http.StatusOK,
			expectedContentType: "text/plain",
		},
		{
			name:                "success: gauge -1.1",
			method:              http.MethodPost,
			path:                "/update/gauge/test/-1.1",
			contentType:         "text/plain",
			expectedStatus:      http.StatusOK,
			expectedContentType: "text/plain",
		},
		{
			name:                "failed: gauge empty value",
			method:              http.MethodPost,
			path:                "/update/gauge/test",
			contentType:         "text/plain",
			expectedStatus:      http.StatusNotFound,
			expectedContentType: "text/plain",
		},
		{
			name:                "failed: gauge wrong value",
			method:              http.MethodPost,
			path:                "/update/gauge/test/dummy",
			contentType:         "text/plain",
			expectedStatus:      http.StatusBadRequest,
			expectedContentType: "text/plain",
		},
		{
			name:                "failed: gauge wrong method",
			method:              http.MethodPatch,
			path:                "/update/gauge/test/1",
			contentType:         "text/plain",
			expectedStatus:      http.StatusMethodNotAllowed,
			expectedContentType: "text/plain",
		},
		{
			name:                "failed: gauge wrong content type",
			method:              http.MethodPost,
			path:                "/update/gauge/test/1",
			contentType:         "application/json",
			expectedStatus:      http.StatusUnsupportedMediaType,
			expectedContentType: "text/plain",
		},

		{
			name:                "success: counter 0",
			method:              http.MethodPost,
			path:                "/update/counter/test/0",
			contentType:         "text/plain",
			expectedStatus:      http.StatusOK,
			expectedContentType: "text/plain",
		},
		{
			name:                "success: counter 1",
			method:              http.MethodPost,
			path:                "/update/counter/test/1",
			contentType:         "text/plain",
			expectedStatus:      http.StatusOK,
			expectedContentType: "text/plain",
		},
		{
			name:                "success: counter -1",
			method:              http.MethodPost,
			path:                "/update/counter/test/1",
			contentType:         "text/plain",
			expectedStatus:      http.StatusOK,
			expectedContentType: "text/plain",
		},
		{
			name:                "failed: counter empty value",
			method:              http.MethodPost,
			path:                "/update/counter/test",
			contentType:         "text/plain",
			expectedStatus:      http.StatusNotFound,
			expectedContentType: "text/plain",
		},
		{
			name:                "failed: counter wrong value",
			method:              http.MethodPost,
			path:                "/update/counter/test/dummy",
			contentType:         "text/plain",
			expectedStatus:      http.StatusBadRequest,
			expectedContentType: "text/plain",
		},
		{
			name:                "failed: counter incorrect value",
			method:              http.MethodPost,
			path:                "/update/counter/test/1.1",
			contentType:         "text/plain",
			expectedStatus:      http.StatusBadRequest,
			expectedContentType: "text/plain",
		},
		{
			name:                "failed: counter wrong method",
			method:              http.MethodPatch,
			path:                "/update/counter/test/1",
			contentType:         "text/plain",
			expectedStatus:      http.StatusMethodNotAllowed,
			expectedContentType: "text/plain",
		},
		{
			name:                "failed: counter wrong content type",
			method:              http.MethodPost,
			path:                "/update/counter/test/1",
			contentType:         "application/json",
			expectedStatus:      http.StatusUnsupportedMediaType,
			expectedContentType: "text/plain",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := NewHandler(repository.NewMemStorage())

			request := httptest.NewRequest(tt.method, tt.path, http.NoBody)
			request.Header.Add("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			h.Update(w, request)
			res := w.Result()
			res.Body.Close()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			assert.Contains(t, res.Header.Get("Content-Type"), tt.expectedContentType)
		})
	}
}
