package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vrnvgasu/metrics/internal/repository"
)

type FindRequest struct {
	MType string `uri:"mtype" binding:"required"`
	Name  string `uri:"name" binding:"required"`
}

func (h *Handler) Find(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/plain")

	var req FindRequest
	if err := c.ShouldBindUri(&req); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)

		return
	}

	m, err := h.Storage.GetByTypeAndID(req.MType, req.Name)
	if err != nil {

		if errors.Is(err, repository.ErrNotFound) {
			http.Error(c.Writer, err.Error(), http.StatusNotFound)
		} else {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	if _, err = c.Writer.Write([]byte(m.ValueToString())); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)

		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}
