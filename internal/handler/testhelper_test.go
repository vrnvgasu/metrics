package handler

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func mustNewRouter(t *testing.T, h *Handler) *gin.Engine {
	t.Helper()
	r, err := NewRouter(h)
	if err != nil {
		t.Fatalf("NewRouter: %v", err)
	}
	return r
}
