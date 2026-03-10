package handler

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/internal/repository"
	"github.com/vrnvgasu/metrics/internal/service/metric"
	"github.com/vrnvgasu/metrics/pkg/helper"
)

func TestFind(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	repo := repository.NewMemStorage()
	err := repo.Add(ctx, &models.Metrics{
		ID:    "1",
		MType: models.Gauge,
		Value: helper.NewRefFloat64(842315.916000),
	})
	require.NoError(t, err)
	err = repo.Add(ctx, &models.Metrics{
		ID:    "2",
		MType: models.Counter,
		Delta: helper.NewRefInt64(11),
	})
	require.NoError(t, err)

	s := metric.NewService(repo)

	tests := []struct {
		name                string
		method              string
		path                string
		contentType         string
		expectedStatus      int
		expectedContentType string
		expectedValue       string
	}{
		{
			name:                "Find gauge",
			method:              http.MethodGet,
			path:                "/value/gauge/1",
			contentType:         "text/plain",
			expectedStatus:      http.StatusOK,
			expectedContentType: "text/plain",
			expectedValue:       "842315.916",
		},
		{
			name:                "Find counter",
			method:              http.MethodGet,
			path:                "/value/counter/2",
			contentType:         "text/plain",
			expectedStatus:      http.StatusOK,
			expectedContentType: "text/plain",
			expectedValue:       "11",
		},
		{
			name:                "Not find by ID",
			method:              http.MethodGet,
			path:                "/value/gauge/some",
			contentType:         "text/plain",
			expectedStatus:      http.StatusNotFound,
			expectedContentType: "text/plain",
			expectedValue:       "",
		},
		{
			name:                "Not find by type",
			method:              http.MethodGet,
			path:                "/value/some/1",
			contentType:         "text/plain",
			expectedStatus:      http.StatusNotFound,
			expectedContentType: "text/plain",
			expectedValue:       "",
		},
		{
			name:                "Bad method",
			method:              http.MethodPost,
			path:                "/value/gauge/1",
			contentType:         "text/plain",
			expectedStatus:      http.StatusNotFound,
			expectedContentType: "text/plain",
			expectedValue:       "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.NoError(t, err)
			h := NewHandler(s)

			request := httptest.NewRequest(tt.method, tt.path, http.NoBody)
			request.Header.Add("Content-Type", tt.contentType)
			w := httptest.NewRecorder()

			NewRouter(h).ServeHTTP(w, request)

			res := w.Result()
			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			assert.Contains(t, res.Header.Get("Content-Type"), tt.expectedContentType)
			if tt.expectedValue != "" {
				assert.Equal(t, tt.expectedValue, string(body))
			}
		})
	}
}
