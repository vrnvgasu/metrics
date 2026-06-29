package middleware

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func testTrustedSubnetRouter(subnet *net.IPNet) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(TrustedSubnet(subnet))
	r.POST("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	return r
}

func mustCIDR(t *testing.T, cidr string) *net.IPNet {
	t.Helper()
	_, n, err := net.ParseCIDR(cidr)
	if err != nil {
		t.Fatalf("parse cidr %q: %v", cidr, err)
	}
	return n
}

func TestTrustedSubnet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		subnet   *net.IPNet
		realIP   string
		wantCode int
	}{
		{
			name:     "nil subnet allows any",
			subnet:   nil,
			realIP:   "",
			wantCode: http.StatusOK,
		},
		{
			name:     "ip in subnet",
			subnet:   mustCIDR(t, "192.168.0.0/24"),
			realIP:   "192.168.0.42",
			wantCode: http.StatusOK,
		},
		{
			name:     "ip out of subnet",
			subnet:   mustCIDR(t, "192.168.0.0/24"),
			realIP:   "10.0.0.1",
			wantCode: http.StatusForbidden,
		},
		{
			name:     "missing header",
			subnet:   mustCIDR(t, "192.168.0.0/24"),
			realIP:   "",
			wantCode: http.StatusForbidden,
		},
		{
			name:     "invalid ip",
			subnet:   mustCIDR(t, "192.168.0.0/24"),
			realIP:   "not-an-ip",
			wantCode: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := testTrustedSubnetRouter(tt.subnet)
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			if tt.realIP != "" {
				req.Header.Set(realIPHeader, tt.realIP)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}
