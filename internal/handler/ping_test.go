package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vrnvgasu/metrics/internal/config"
	"github.com/vrnvgasu/metrics/internal/repository/mem"
	"github.com/vrnvgasu/metrics/internal/service/audit"
	"github.com/vrnvgasu/metrics/internal/service/healthcheck"
	"github.com/vrnvgasu/metrics/internal/service/metric"
)

type mockHealthService struct {
	err error
}

func (m *mockHealthService) CheckPing(_ context.Context) error {
	return m.err
}

func TestPing(t *testing.T) {
	t.Parallel()

	publisher, err := audit.NewAudit(&config.ServerCnf{})
	require.NoError(t, err)

	s := metric.NewService(mem.NewMemStorage())

	tests := []struct {
		name           string
		healthService  HealthService
		expectedStatus int
	}{
		{
			name:           "ping success",
			healthService:  healthcheck.NewService(mem.NewMemStorage()),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "ping error",
			healthService:  &mockHealthService{err: errors.New("db unavailable")},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := NewHandler(s, tt.healthService, publisher, &config.ServerCnf{})

			req := httptest.NewRequest(http.MethodGet, "/ping", http.NoBody)
			w := httptest.NewRecorder()

			NewRouter(h).ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
