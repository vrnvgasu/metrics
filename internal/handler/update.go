package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

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
		h.responseError(c, serviceerrors.BadRequestError())

		return
	}

	if ok := h.validateHeaderHashSHA256(c); !ok {
		return
	}

	metric, err := models.NewMetricsFromStrings(req.MType, req.Name, req.Value)
	if err != nil {
		h.responseError(c, serviceerrors.BadRequestError())

		return
	}

	if err = h.MetricService.CreateOrUpdate(c, []*models.Metrics{&metric}); err != nil {
		h.responseError(c, err)

		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
