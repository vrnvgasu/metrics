package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) List(c *gin.Context) {
	list, err := h.MetricService.AllMetrics(c)
	if err != nil {
		h.responseError(c, err)

		return
	}
	items := make([]string, 0, len(list))
	for _, m := range list {
		items = append(items, fmt.Sprintf("%s: %s", m.ID, m.ValueToString()))
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"items": items,
	})
}
