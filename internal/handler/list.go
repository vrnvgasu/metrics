package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vrnvgasu/metrics/internal/handler/response"
)

// List обрабатывает GET / — возвращает HTML-страницу со списком всех метрик.
func (h *Handler) List(c *gin.Context) {
	list, err := h.MetricService.AllMetrics(c)
	if err != nil {
		response.ResponseError(c, err)

		return
	}
	items := make([]string, 0, len(list))
	for _, metrics := range list {
		items = append(items, fmt.Sprintf("%s: %s", metrics.ID, metrics.ValueToString()))
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"items": items,
	})
}
