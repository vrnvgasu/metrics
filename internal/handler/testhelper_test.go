package handler

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/vrnvgasu/metrics/internal/config"
	"github.com/vrnvgasu/metrics/internal/service/audit"
)

func mustNewRouter(t *testing.T, h *Handler) *gin.Engine {
	t.Helper()
	r, err := NewRouter(h)
	if err != nil {
		t.Fatalf("NewRouter: %v", err)
	}
	return r
}

func newTestPublisher(t *testing.T) *audit.Event {
	t.Helper()
	p, err := audit.NewAudit(&config.ServerCnf{})
	require.NoError(t, err)
	return p
}
