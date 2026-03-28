package mem

import (
	"context"

	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/internal/repository"
)

func (s *Storage) GetByTypeAndID(_ context.Context, mtype, id string) (*models.Metrics, error) {
	m, ok := s.metrics[id]
	if !ok || m.MType != mtype {
		return nil, repository.ErrNotFound
	}

	return &m, nil
}
