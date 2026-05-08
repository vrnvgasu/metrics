package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vrnvgasu/metrics/internal/handler/response"
	"github.com/vrnvgasu/metrics/internal/logger"
	models "github.com/vrnvgasu/metrics/internal/model"
	serviceerrors "github.com/vrnvgasu/metrics/internal/service/errors"
)

type UpdateJSONRequest struct {
	ID    string            `json:"id" binding:"required"`   // имя метрики
	MType models.MetricType `json:"type" binding:"required"` // параметр, принимающий значение gauge или counter
	Delta *int64            `json:"delta,omitempty"`         // значение метрики в случае передачи counter
	Value *float64          `json:"value,omitempty"`         // значение метрики в случае передачи gauge
}

func (r *UpdateJSONRequest) ToMetrics() models.Metrics {
	return models.Metrics{
		ID:    r.ID,
		MType: r.MType,
		Delta: r.Delta,
		Value: r.Value,
	}
}

func (h *Handler) UpdateJSON(c *gin.Context) {
	var body UpdateJSONRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		response.ResponseError(c, serviceerrors.BadRequestError())

		return
	}

	if err := h.MetricService.CreateOrUpdate(c, []models.Metrics{body.ToMetrics()}); err != nil {
		response.ResponseError(c, err)

		return
	}

	if err := h.publisher.Notify(c, []string{body.ID}, c.ClientIP()); err != nil {
		logger.Log.Errorf("failed to notify audit: %s", err.Error())
	}

	c.JSON(http.StatusOK, http.NoBody)
}

func (h *Handler) UpdateJSONList(c *gin.Context) {
	var body []UpdateJSONRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		response.ResponseError(c, serviceerrors.BadRequestError())

		return
	}

	metricList := make(models.MetricsList, 0, len(body))
	for _, v := range body {
		metricList = append(metricList, v.ToMetrics())
	}
	if err := h.MetricService.CreateOrUpdate(c, metricList); err != nil {
		response.ResponseError(c, err)

		return
	}

	if err := h.publisher.Notify(c, metricList.IDList(), c.ClientIP()); err != nil {
		logger.Log.Errorf("failed to notify audit: %s", err.Error())
	}

	c.JSON(http.StatusOK, http.NoBody)
}
