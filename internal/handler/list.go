package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) List(c *gin.Context) {
	mMap := h.Storage.Map()
	items := make([]string, 0, len(mMap))
	for _, m := range mMap {
		items = append(items, fmt.Sprintf("%s: %s", m.ID, m.ValueToString()))
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"items": items,
	})
}
