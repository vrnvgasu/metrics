package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vrnvgasu/metrics/internal/handler/response"
)

func (h *Handler) Ping(c *gin.Context) {
	if err := h.HealthService.CheckPing(c); err != nil {
		response.ResponseError(c, err)

		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
