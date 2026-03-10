package repository

import (
	"context"
	"fmt"

	models "github.com/vrnvgasu/metrics/internal/model"
)

func (ms *MemStorage) Add(_ context.Context, m *models.Metrics) error {
	switch m.MType {
	case models.Gauge:
		ms.addGauge(*m)

		return nil
	case models.Counter:
		ms.addCounter(*m)

		return nil
	default:
		return fmt.Errorf("metrics type %s not supported: %w", m.MType, ErrNotSupport)
	}
}

func (ms *MemStorage) addGauge(m models.Metrics) {
	ms.metrics[m.ID] = m
}

func (ms *MemStorage) addCounter(m models.Metrics) {
	oldM, ok := ms.metrics[m.ID]
	if !ok {
		ms.metrics[m.ID] = m

		return
	}

	oldDelta := *oldM.Delta
	oldDelta += *m.Delta
	oldM.Delta = &oldDelta

	ms.metrics[m.ID] = oldM
}
