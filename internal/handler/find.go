package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	seriveerrors "github.com/vrnvgasu/metrics/internal/service/errors"
)

type FindRequest struct {
	MType string `uri:"mtype" binding:"required"`
	Name  string `uri:"name" binding:"required"`
}

func (h *Handler) Find(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/plain")

	var req FindRequest
	if err := c.ShouldBindUri(&req); err != nil {
		h.responseError(c, seriveerrors.BadRequestError(err.Error()))

		return
	}

	m, err := h.MetricService.FindByTypeAndID(c, req.MType, req.Name)
	if err != nil {
		h.responseError(c, err)

		return
	}

	if _, err = c.Writer.Write([]byte(m.ValueToString())); err != nil {
		h.responseError(c, err)

		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}
