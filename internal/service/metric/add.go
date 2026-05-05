package metric

import (
	"context"
	"errors"
	"fmt"

	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/internal/repository"
)

func (s *Service) CreateOrUpdate(ctx context.Context, list []*models.Metrics) error {
	err := s.storage.DoInTransaction(ctx, func(ctx context.Context) error {
		return s.createOrUpdate(ctx, list)
	})

	if err != nil {
		return fmt.Errorf("metric.CreateOrUpdate DoInTransaction: %w", err)
	}

	return nil
}

func (s *Service) createOrUpdate(ctx context.Context, list []*models.Metrics) error {
	for _, metrics := range list {
		switch metrics.MType {
		case models.Gauge:
			if err := s.addGauge(ctx, metrics); err != nil {
				return fmt.Errorf("metric.createOrUpdate addGauge: %w", err)
			}
		case models.Counter:
			if err := s.addCounter(ctx, metrics); err != nil {
				return fmt.Errorf("metric.createOrUpdate addCounter: %w", err)
			}
		default:
			return fmt.Errorf("metrics type %s not supported: %w", metrics.MType, repository.ErrNotSupport)
		}
	}

	return nil
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
