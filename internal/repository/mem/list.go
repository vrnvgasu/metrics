package mem

import (
	"context"

	models "github.com/vrnvgasu/metrics/internal/model"
)

func (s *Storage) List(_ context.Context) (models.MetricsList, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]models.Metrics, 0, len(s.metrics))
	for _, m := range s.metrics {
		list = append(list, m)
	}

	return list, nil
}
