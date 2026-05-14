package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vrnvgasu/metrics/internal/handler/response"
	"github.com/vrnvgasu/metrics/internal/logger"
	models "github.com/vrnvgasu/metrics/internal/model"
	serviceerrors "github.com/vrnvgasu/metrics/internal/service/errors"
)

// UpdateJSONRequest — тело запроса для POST /update/ и POST /updates.
type UpdateJSONRequest struct {
	ID    string            `json:"id" binding:"required"`   // имя метрики
	MType models.MetricType `json:"type" binding:"required"` // тип: gauge или counter
	Delta *int64            `json:"delta,omitempty"`         // значение counter
	Value *float64          `json:"value,omitempty"`         // значение gauge
}

// ToMetrics конвертирует запрос в модель Metrics.
func (r *UpdateJSONRequest) ToMetrics() models.Metrics {
	return models.Metrics{
		ID:    r.ID,
		MType: r.MType,
		Delta: r.Delta,
		Value: r.Value,
	}
}

// UpdateJSON обрабатывает POST /update/ — обновляет одну метрику из JSON-тела.
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

// UpdateJSONList обрабатывает POST /updates — пакетное обновление метрик из JSON-массива.
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
