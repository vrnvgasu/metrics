package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/vrnvgasu/metrics/internal/config"
	"github.com/vrnvgasu/metrics/pkg/hash"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func hashRouter(cnf *config.ServerCnf) http.Handler {
	r := gin.New()
	r.Use(Hash(cnf))
	r.POST("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	return r
}

func TestHash(t *testing.T) {
	t.Parallel()

	body := []byte(`{"id":"Alloc","type":"gauge","value":1.0}`)

	validSig, err := hash.PrepareHeaderHashSHA256("secret", body)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		cnf            *config.ServerCnf
		body           []byte
		signature      string
		expectedStatus int
	}{
		{
			name:           "empty key; skip validation",
			cnf:            &config.ServerCnf{},
			body:           body,
			signature:      "any",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "valid signature",
			cnf:            &config.ServerCnf{Key: "secret"},
			body:           body,
			signature:      validSig,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid signature",
			cnf:            &config.ServerCnf{Key: "secret"},
			body:           body,
			signature:      "wrong",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing signature header",
			cnf:            &config.ServerCnf{Key: "secret"},
			body:           body,
			signature:      "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(tt.body))
			if tt.signature != "" {
				req.Header.Set("HashSHA256", tt.signature)
			}
			w := httptest.NewRecorder()

			hashRouter(tt.cnf).ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
