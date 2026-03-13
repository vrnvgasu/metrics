package metric

import (
	"context"
	"errors"
	"fmt"

	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/internal/repository"
)

func (s *Service) CreateOrUpdate(ctx context.Context, m *models.Metrics) error {
	switch m.MType {
	case models.Gauge:
		return s.addGauge(ctx, m)
	case models.Counter:
		return s.addCounter(ctx, m)
	default:
		return fmt.Errorf("metrics type %s not supported: %w", m.MType, repository.ErrNotSupport)
	}
}

func (s *Service) addGauge(ctx context.Context, m *models.Metrics) error {
	return s.storage.Save(ctx, m)
}

func (s *Service) addCounter(ctx context.Context, m *models.Metrics) error {
	oldM, err := s.storage.GetByTypeAndID(ctx, m.MType, m.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return s.storage.Save(ctx, m)
		}

		return fmt.Errorf("metric.addCounter GetByTypeAndID: %w", err)
	}

	oldDelta := *oldM.Delta
	oldDelta += *m.Delta
	oldM.Delta = &oldDelta

	return s.storage.Save(ctx, oldM)
}
