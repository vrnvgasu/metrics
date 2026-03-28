package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Ping(c *gin.Context) {
	if err := h.HealthService.CheckPing(c); err != nil {
		h.responseError(c, err)

		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
