package mem

import (
	models "github.com/vrnvgasu/metrics/internal/model"
)

type MemStorage struct {
	metrics models.MetricsMap
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: make(models.MetricsMap),
	}
}
