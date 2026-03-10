package repository

import (
	"context"

	models "github.com/vrnvgasu/metrics/internal/model"
)

func (ms *MemStorage) List(_ context.Context) models.MetricsList {
	list := make([]models.Metrics, 0, len(ms.metrics))
	for _, m := range ms.metrics {
		list = append(list, m)
	}

	return list
}
