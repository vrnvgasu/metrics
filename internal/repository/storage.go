package repository

import models "github.com/vrnvgasu/metrics/internal/model"

var Storage *MemStorage

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
	case models.Counter:
		ms.addCounter(m)
	}

	return nil
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

	oldValue := *oldM.Value
	oldValue += *m.Value
	oldM.Value = &oldValue
}
