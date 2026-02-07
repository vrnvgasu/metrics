package repository

import models "github.com/vrnvgasu/metrics/internal/model"

func (ms *MemStorage) Map() models.MetricsMap {
	return ms.metrics
}
