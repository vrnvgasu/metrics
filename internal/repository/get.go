package repository

import (
	models "github.com/vrnvgasu/metrics/internal/model"
)

func (ms *MemStorage) GetByTypeAndID(mtype, id string) (models.Metrics, error) {
	m, ok := ms.metrics[id]
	if !ok || m.MType != mtype {
		return models.Metrics{}, ErrNotFound
	}

	return m, nil
}
