package repository

import (
	"fmt"

	models "github.com/vrnvgasu/metrics/internal/model"
)

type MemStorage struct {
	Metrics map[string]models.Metrics
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Metrics: make(map[string]models.Metrics),
	}
}

func (ms *MemStorage) Add(m models.Metrics) error {
	switch m.MType {
	case models.Gauge:
		ms.addGauge(m)

		return nil
	case models.Counter:
		ms.addCounter(m)

		return nil
	default:
		return fmt.Errorf("metrics type %s not supported", m.MType)
	}
}

func (ms *MemStorage) addGauge(m models.Metrics) {
	ms.Metrics[m.ID] = m
}

func (ms *MemStorage) addCounter(m models.Metrics) {
	oldM, ok := ms.Metrics[m.ID]
	if !ok {
		ms.Metrics[m.ID] = m

		return
	}

	oldDelta := *oldM.Delta
	oldDelta += *m.Delta
	oldM.Delta = &oldDelta

	ms.Metrics[m.ID] = oldM
}
