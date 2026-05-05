package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vrnvgasu/metrics/internal/handler/response"
	models "github.com/vrnvgasu/metrics/internal/model"
	serviceerrors "github.com/vrnvgasu/metrics/internal/service/errors"
)

type FindRequest struct {
	MType string `uri:"mtype" binding:"required"`
	Name  string `uri:"name" binding:"required"`
}

func (h *Handler) Find(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/plain")

	var req FindRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.ResponseError(c, serviceerrors.BadRequestError())

		return
	}

	metrics, err := h.MetricService.FindByTypeAndID(c, models.MetricType(req.MType), req.Name)
	if err != nil {
		response.ResponseError(c, err)

		return
	}

	if _, err = c.Writer.Write([]byte(metrics.ValueToString())); err != nil {
		response.ResponseError(c, err)

		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}
