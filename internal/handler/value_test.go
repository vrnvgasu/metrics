package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/internal/repository"
	"github.com/vrnvgasu/metrics/pkg/helper"
)

func TestValue(t *testing.T) {
	t.Parallel()

	const path = "/value"

	repo := repository.NewMemStorage()
	err := repo.Add(models.Metrics{
		ID:    "1",
		MType: models.Gauge,
		Value: helper.NewRefFloat64(842315.916000),
	})
	require.NoError(t, err)
	err = repo.Add(models.Metrics{
		ID:    "2",
		MType: models.Counter,
		Delta: helper.NewRefInt64(11),
	})
	require.NoError(t, err)

	tests := []struct {
		name                string
		method              string
		body                any
		expectedStatus      int
		expectedContentType string
		expectedBody        any
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := NewHandler(repo)

			body, err := json.Marshal(tt.body)
			require.NoError(t, err)
			request := httptest.NewRequest(tt.method, path, bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			NewRouter(h).ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			assert.Contains(t, res.Header.Get("Content-Type"), tt.expectedContentType)

			if res.StatusCode == http.StatusOK {
				v := &ValueResponse{}
				err = json.NewDecoder(res.Body).Decode(v)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedBody, *v)
			}
		})
	}
}
