package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/internal/repository"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	metrics, err := h.Storage.GetByTypeAndID(body.MType, body.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, NewValueResponseFromMetric(metrics))
}
