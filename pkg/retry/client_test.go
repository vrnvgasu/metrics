package retry

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	c := NewClient(nil)
	require.NotNil(t, c)

	c2 := NewClient(&ClientConfig{})
	require.NotNil(t, c2)
}

func TestClient_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{"200 ok", http.StatusOK, false},
		{"404 not found", http.StatusNotFound, false},
		{"501 not implemented", http.StatusNotImplemented, true},
	}

	cfg := &ClientConfig{retryConfig: &RetryConfig{MaxRetries: 1, StartRetryInterval: 0, AddRetryPeriod: 0}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer srv.Close()

			c := NewClient(cfg)
			req, err := http.NewRequest(http.MethodGet, srv.URL, http.NoBody)
			require.NoError(t, err)

			resp, err := c.Do(req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				resp.Body.Close()
				assert.Equal(t, tt.statusCode, resp.StatusCode)
			}
		})
	}
}

func TestClient_Do_WithBody(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	cfg := &ClientConfig{retryConfig: &RetryConfig{MaxRetries: 0}}
	c := NewClient(cfg)

	req, err := http.NewRequest(http.MethodPost, srv.URL, strings.NewReader("body"))
	require.NoError(t, err)

	resp, err := c.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestClient_ShouldRetry(t *testing.T) {
	t.Parallel()

	c := NewClient(nil)

	assert.True(t, c.shouldRetry(&http.Response{StatusCode: http.StatusInternalServerError}))
	assert.True(t, c.shouldRetry(&http.Response{StatusCode: http.StatusServiceUnavailable}))
	assert.False(t, c.shouldRetry(&http.Response{StatusCode: http.StatusNotImplemented}))
	assert.False(t, c.shouldRetry(&http.Response{StatusCode: http.StatusOK}))
}
