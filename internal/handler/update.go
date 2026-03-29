package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vrnvgasu/metrics/internal/handler/response"
	models "github.com/vrnvgasu/metrics/internal/model"
	serviceerrors "github.com/vrnvgasu/metrics/internal/service/errors"
)

type UpdateRequest struct {
	MType string `uri:"mtype" binding:"required"`
	Name  string `uri:"name" binding:"required"`
	Value string `uri:"value" binding:"required"`
}

func (h *Handler) Update(c *gin.Context) {
	c.Writer.Header().Add("Content-Type", "text/plain")

	var req UpdateRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.ResponseError(c, serviceerrors.BadRequestError())

		return
	}

	metric, err := models.NewMetricsFromStrings(req.MType, req.Name, req.Value)
	if err != nil {
		response.ResponseError(c, serviceerrors.BadRequestError())

		return
	}

	if err = h.MetricService.CreateOrUpdate(c, []*models.Metrics{&metric}); err != nil {
		response.ResponseError(c, err)

		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
