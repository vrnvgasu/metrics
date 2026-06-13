package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mockhandler "github.com/vrnvgasu/metrics/internal/handler/mocks"
	"go.uber.org/mock/gomock"

	"github.com/vrnvgasu/metrics/internal/config"
	"github.com/vrnvgasu/metrics/internal/repository/mem"
	"github.com/vrnvgasu/metrics/internal/service/audit"
	"github.com/vrnvgasu/metrics/internal/service/metric"
)

func TestPing(t *testing.T) {
	t.Parallel()

	publisher, err := audit.NewAudit(&config.ServerCnf{})
	require.NoError(t, err)

	s := metric.NewService(mem.NewMemStorage())

	controller := gomock.NewController(t)
	defer controller.Finish()

	tests := []struct {
		name           string
		healthService  func() HealthService
		expectedStatus int
	}{
		{
			name: "ping success",
			healthService: func() HealthService {
				mock := mockhandler.NewMockHealthService(controller)
				mock.EXPECT().CheckPing(gomock.Any()).Return(nil)

				return mock
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "ping error",

			healthService: func() HealthService {
				mock := mockhandler.NewMockHealthService(controller)
				mock.EXPECT().CheckPing(gomock.Any()).Return(errors.New("db unavailable"))

				return mock
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := NewHandler(s, tt.healthService(), publisher, &config.ServerCnf{})

			req := httptest.NewRequest(http.MethodGet, "/ping", http.NoBody)
			w := httptest.NewRecorder()

			NewRouter(h).ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
