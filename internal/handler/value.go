package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	models "github.com/vrnvgasu/metrics/internal/model"
	serviceerrors "github.com/vrnvgasu/metrics/internal/service/errors"
)

type ValueRequest struct {
	ID    string `json:"id" binding:"required"`
	MType string `json:"type" binding:"required"`
}

type ValueResponse struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func NewValueResponseFromMetric(m models.Metrics) ValueResponse {
	return ValueResponse{
		ID:    m.ID,
		MType: m.MType,
		Delta: m.Delta,
		Value: m.Value,
	}
}

func (h *Handler) Value(c *gin.Context) {
	var body ValueRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		h.responseError(c, serviceerrors.BadRequestError(err.Error()))

		return
	}

	metrics, err := h.MetricService.FindByTypeAndID(c, body.MType, body.ID)
	if err != nil {
		h.responseError(c, err)

		return
	}

	c.JSON(http.StatusOK, NewValueResponseFromMetric(*metrics))
}
