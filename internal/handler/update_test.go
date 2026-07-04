package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/vrnvgasu/metrics/internal/config"
	mockhandler "github.com/vrnvgasu/metrics/internal/handler/mocks"
	"github.com/vrnvgasu/metrics/internal/repository/mem"
	"github.com/vrnvgasu/metrics/internal/service/audit"
	serviceerrors "github.com/vrnvgasu/metrics/internal/service/errors"
	"github.com/vrnvgasu/metrics/internal/service/metric"
)

func TestUpdate(t *testing.T) {
	t.Parallel()

	publisher, err := audit.NewAudit(&config.ServerCnf{})
	require.NoError(t, err)

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
			expectedStatus:      http.StatusNotFound,
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
			expectedStatus:      http.StatusNotFound,
			expectedContentType: "text/plain",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mem.NewMemStorage()
			s := metric.NewService(repo)
			h := NewHandler(s, nil, publisher, &config.ServerCnf{})

			request := httptest.NewRequest(tt.method, tt.path, http.NoBody)
			request.Header.Add("Content-Type", tt.contentType)
			w := httptest.NewRecorder()

			mustNewRouter(t, h).ServeHTTP(w, request)

			res := w.Result()
			res.Body.Close()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			assert.Contains(t, res.Header.Get("Content-Type"), tt.expectedContentType)
		})
	}
}

// TestUpdate_ServiceError покрывает 500-ветку Update при ошибке CreateOrUpdate.
func TestUpdate_ServiceError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	ms := mockhandler.NewMockMetricService(ctrl)
	ms.EXPECT().CreateOrUpdate(gomock.Any(), gomock.Any()).Return(serviceerrors.InternalError())

	h := NewHandler(ms, nil, newTestPublisher(t), &config.ServerCnf{})

	req := httptest.NewRequest(http.MethodPost, "/update/gauge/test/1", http.NoBody)
	w := httptest.NewRecorder()
	mustNewRouter(t, h).ServeHTTP(w, req)

	res := w.Result()
	res.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestUpdate_BindError(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	h := NewHandler(metric.NewService(mem.NewMemStorage()), nil, newTestPublisher(t), &config.ServerCnf{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", http.NoBody)

	h.Update(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
