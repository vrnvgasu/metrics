package metric

import (
	"context"
	"errors"
	"fmt"

	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/internal/repository"
)

func (s *Service) CreateOrUpdate(ctx context.Context, list []*models.Metrics) (err error) {
	txStorage, err := s.storage.WithTx(ctx)
	if err != nil {
		return fmt.Errorf("metric.CreateOrUpdateList WithTx: %w", err)
	}

Loop:
	for _, m := range list {
		switch m.MType {
		case models.Gauge:
			if err = s.addGauge(ctx, txStorage, m); err != nil {
				break Loop
			}
		case models.Counter:
			if err = s.addCounter(ctx, txStorage, m); err != nil {
				break Loop
			}
		default:
			err = fmt.Errorf("metrics type %s not supported: %w", m.MType, repository.ErrNotSupport)
			break Loop
		}
	}

	if err != nil {
		txStorage.Rollback()

		return err
	}

	return txStorage.Commit()
}

func (s *Service) addGauge(ctx context.Context, storage repository.Storage, m *models.Metrics) error {
	return storage.Save(ctx, m)
}

func (s *Service) addCounter(ctx context.Context, storage repository.Storage, m *models.Metrics) error {
	oldM, err := storage.GetByTypeAndID(ctx, m.MType, m.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return storage.Save(ctx, m)
		}

		return fmt.Errorf("metric.addCounter GetByTypeAndID: %w", err)
	}

	oldDelta := *oldM.Delta
	oldDelta += *m.Delta
	oldM.Delta = &oldDelta

	return storage.Save(ctx, oldM)
}
