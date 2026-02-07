package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	models "github.com/vrnvgasu/metrics/internal/model"
)

type UpdateRequest struct {
	MType string `uri:"mtype" binding:"required"`
	Name  string `uri:"name" binding:"required"`
	Value string `uri:"value" binding:"required"`
}

func (h *Handler) Update(c *gin.Context) {
	c.Writer.Header().Add("Content-Type", "text/plain")

	if !strings.Contains(c.Request.Header.Get("Content-Type"), "text/plain") {
		http.Error(c.Writer, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)

		return
	}

	var req UpdateRequest
	if err := c.ShouldBindUri(&req); err != nil {
		http.Error(c.Writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	metric, err := models.NewMetricsFromStrings(req.MType, req.Name, req.Value)
	if err != nil {
		http.Error(c.Writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	if err = h.Storage.Add(metric); err != nil {
		http.Error(c.Writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
