package mem

import (
	models "github.com/vrnvgasu/metrics/internal/model"
)

type Storage struct {
	metrics models.MetricsMap
}

func NewMemStorage() *Storage {
	return &Storage{
		metrics: make(models.MetricsMap),
	}
}
