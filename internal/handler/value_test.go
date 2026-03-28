package handler

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vrnvgasu/metrics/internal/config"
	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/internal/repository/mem"
	"github.com/vrnvgasu/metrics/internal/service/metric"
	"github.com/vrnvgasu/metrics/pkg/helper"
)

func TestValue(t *testing.T) {
	t.Parallel()

	const path = "/value"

	ctx := t.Context()

	repo := mem.NewMemStorage()
	err := repo.Save(ctx, &models.Metrics{
		ID:    "1",
		MType: models.Gauge,
		Value: helper.NewRefFloat64(842315.916000),
	})
	require.NoError(t, err)
	err = repo.Save(ctx, &models.Metrics{
		ID:    "2",
		MType: models.Counter,
		Delta: helper.NewRefInt64(11),
	})
	require.NoError(t, err)

	s := metric.NewService(repo)

	tests := []struct {
		name                string
		method              string
		body                any
		expectedStatus      int
		expectedContentType string
		expectedBody        any
		gzipResponse        bool
	}{
		{
			name:   "success gauge",
			method: http.MethodPost,
			body: ValueRequest{
				ID:    "1",
				MType: models.Gauge,
			},
			expectedStatus:      http.StatusOK,
			expectedContentType: "application/json",
			expectedBody: ValueResponse{
				ID:    "1",
				MType: models.Gauge,
				Value: helper.NewRefFloat64(842315.916000),
			},
		},
		{
			name:   "success counter",
			method: http.MethodPost,
			body: ValueRequest{
				ID:    "2",
				MType: models.Counter,
			},
			expectedStatus:      http.StatusOK,
			expectedContentType: "application/json",
			expectedBody: ValueResponse{
				ID:    "2",
				MType: models.Counter,
				Delta: helper.NewRefInt64(11),
			},
		},
		{
			name:   "wrong request",
			method: http.MethodPost,
			body: struct {
				SomeID string `json:"someID"`
			}{SomeID: "1"},
			expectedStatus:      http.StatusBadRequest,
			expectedContentType: "application/json",
		},
		{
			name:                "dummy request",
			method:              http.MethodPost,
			body:                "dummy",
			expectedStatus:      http.StatusBadRequest,
			expectedContentType: "application/json",
		},
		{
			name:                "null request",
			method:              http.MethodPost,
			body:                nil,
			expectedStatus:      http.StatusBadRequest,
			expectedContentType: "application/json",
		},
		{
			name:   "wrong method",
			method: http.MethodPut,
			body: UpdateJSONRequest{
				ID:    "1",
				MType: models.Gauge,
				Value: helper.NewRefFloat64(1.1),
			},
			expectedStatus:      http.StatusNotFound,
			expectedContentType: "text/plain",
		},
		{
			name:   "success gzip",
			method: http.MethodPost,
			body: ValueRequest{
				ID:    "1",
				MType: models.Gauge,
			},
			expectedStatus:      http.StatusOK,
			expectedContentType: "application/json",
			expectedBody: ValueResponse{
				ID:    "1",
				MType: models.Gauge,
				Value: helper.NewRefFloat64(842315.916000),
			},
			gzipResponse: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := NewHandler(s, nil, &config.ServerCnf{})

			body, err := json.Marshal(tt.body)
			require.NoError(t, err)
			request := httptest.NewRequest(tt.method, path, bytes.NewBuffer(body))
			if tt.gzipResponse {
				request.Header.Set("Accept-Encoding", "gzip")
			}
			w := httptest.NewRecorder()

			NewRouter(h).ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			assert.Contains(t, res.Header.Get("Content-Type"), tt.expectedContentType)

			if res.StatusCode == http.StatusOK {
				reader := res.Body
				if tt.gzipResponse {
					reader, err = gzip.NewReader(reader)
					require.NoError(t, err)
				}

				v := &ValueResponse{}
				err = json.NewDecoder(reader).Decode(v)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedBody, *v)
			}
		})
	}
}
