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

	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/internal/repository/mem"
	"github.com/vrnvgasu/metrics/internal/service/metric"
	"github.com/vrnvgasu/metrics/pkg/helper"
)

func TestUpdateJSON(t *testing.T) {
	t.Parallel()

	const path = "/update"

	tests := []struct {
		name                string
		method              string
		body                any
		expectedStatus      int
		expectedContentType string
	}{
		{
			name:   "success gauge",
			method: http.MethodPost,
			body: UpdateJSONRequest{
				ID:    "11",
				MType: models.Gauge,
				Value: helper.NewRefFloat64(1.1),
			},
			expectedStatus:      http.StatusOK,
			expectedContentType: "application/json",
		},
		{
			name:   "success counter",
			method: http.MethodPost,
			body: UpdateJSONRequest{
				ID:    "12",
				MType: models.Counter,
				Delta: helper.NewRefInt64(11),
			},
			expectedStatus:      http.StatusOK,
			expectedContentType: "application/json",
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mem.NewMemStorage()
			s := metric.NewService(repo)
			h := NewHandler(s, nil)

			body, err := json.Marshal(tt.body)
			require.NoError(t, err)
			request := httptest.NewRequest(tt.method, path, bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			NewRouter(h).ServeHTTP(w, request)

			res := w.Result()
			res.Body.Close()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			assert.Contains(t, res.Header.Get("Content-Type"), tt.expectedContentType)

			if res.StatusCode == http.StatusOK {
				req := tt.body.(UpdateJSONRequest)
				metric, err := repo.GetByTypeAndID(t.Context(), req.MType, req.ID)
				require.NoError(t, err)
				assert.NotEmpty(t, metric.ID)
				assert.Equal(t, req.ID, metric.ID)
				assert.Equal(t, req.MType, metric.MType)

				if metric.Delta != nil {
					assert.Equal(t, *req.Delta, *metric.Delta)
				}
				if metric.Value != nil {
					assert.Equal(t, *req.Value, *metric.Value)
				}
			}
		})
	}
}

func TestUpdateJSONGzipCompress(t *testing.T) {
	t.Parallel()

	var (
		path          = "/update"
		updateRequest = UpdateJSONRequest{
			ID:    "11",
			MType: models.Gauge,
			Value: helper.NewRefFloat64(1.1),
		}
		buf bytes.Buffer
	)

	repo := mem.NewMemStorage()
	s := metric.NewService(repo)
	h := NewHandler(s, nil)

	body, err := json.Marshal(updateRequest)
	require.NoError(t, err)

	gzipWriter := gzip.NewWriter(&buf)
	_, err = gzipWriter.Write(body)
	gzipWriter.Close()
	require.NoError(t, err)

	request := httptest.NewRequest(http.MethodPost, path, bytes.NewBuffer(buf.Bytes()))
	request.Header.Set("Content-Encoding", "gzip")
	w := httptest.NewRecorder()

	NewRouter(h).ServeHTTP(w, request)

	res := w.Result()
	res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	metric, err := repo.GetByTypeAndID(t.Context(), updateRequest.MType, updateRequest.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, metric.ID)
	assert.Equal(t, updateRequest.ID, metric.ID)
	assert.Equal(t, updateRequest.MType, metric.MType)
	assert.Equal(t, *updateRequest.Value, *metric.Value)
}
