package metric

import (
	"context"

	models "github.com/vrnvgasu/metrics/internal/model"
)

func (s *Service) AllMetrics(ctx context.Context) (models.MetricsList, error) {
	return s.storage.List(ctx)
}
