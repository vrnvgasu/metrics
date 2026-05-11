package mem

import (
	"context"

	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/internal/repository"
)

func (s *Storage) GetByTypeAndID(_ context.Context, mtype models.MetricType, id string) (*models.Metrics, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.metrics[id]
	if !ok || m.MType != mtype {
		return nil, repository.ErrNotFound
	}

	return &m, nil
}
