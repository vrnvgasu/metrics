package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/vrnvgasu/metrics/internal/config"
	mockhandler "github.com/vrnvgasu/metrics/internal/handler/mocks"
	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/internal/repository/mem"
	"github.com/vrnvgasu/metrics/internal/service/audit"
	serviceerrors "github.com/vrnvgasu/metrics/internal/service/errors"
	"github.com/vrnvgasu/metrics/internal/service/metric"
	"github.com/vrnvgasu/metrics/pkg/helper"
)

func TestList(t *testing.T) {
	ctx := t.Context()

	repo := mem.NewMemStorage()
	err := repo.Save(ctx, &models.Metrics{
		ID:    "1",
		MType: models.Gauge,
		Value: helper.NewRefFloat64(11.1),
	})
	require.NoError(t, err)
	err = repo.Save(ctx, &models.Metrics{
		ID:    "2",
		MType: models.Counter,
		Delta: helper.NewRefInt64(11),
	})
	require.NoError(t, err)

	s := metric.NewService(repo)

	publisher, err := audit.NewAudit(&config.ServerCnf{})
	require.NoError(t, err)

	tests := []struct {
		name                string
		method              string
		path                string
		expectedStatus      int
		expectedContentType string
	}{
		{
			name:                "success",
			method:              http.MethodGet,
			path:                "/",
			expectedStatus:      http.StatusOK,
			expectedContentType: "text/html",
		},
		{
			name:                "Bad method",
			method:              http.MethodPost,
			path:                "/dummy",
			expectedStatus:      http.StatusNotFound,
			expectedContentType: "text/plain",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(s, nil, publisher, &config.ServerCnf{})

			request := httptest.NewRequest(tt.method, tt.path, http.NoBody)
			w := httptest.NewRecorder()

			mustNewRouter(t, h).ServeHTTP(w, request)

			res := w.Result()
			res.Body.Close()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			assert.Contains(t, res.Header.Get("Content-Type"), tt.expectedContentType)
		})
	}
}

func TestList_ServiceError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	ms := mockhandler.NewMockMetricService(ctrl)
	ms.EXPECT().AllMetrics(gomock.Any()).Return(nil, serviceerrors.InternalError())

	h := NewHandler(ms, nil, newTestPublisher(t), &config.ServerCnf{})

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	w := httptest.NewRecorder()
	mustNewRouter(t, h).ServeHTTP(w, req)

	res := w.Result()
	res.Body.Close()
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}
