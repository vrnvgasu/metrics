package repository

import (
	"context"

	models "github.com/vrnvgasu/metrics/internal/model"
)

func (ms *MemStorage) GetByTypeAndID(_ context.Context, mtype, id string) (*models.Metrics, error) {
	m, ok := ms.metrics[id]
	if !ok || m.MType != mtype {
		return nil, ErrNotFound
	}

	return &m, nil
}
