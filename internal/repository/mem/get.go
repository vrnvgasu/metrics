package mem

import (
	"context"

	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/internal/repository"
)

func (ms *MemStorage) GetByTypeAndID(_ context.Context, mtype, id string) (*models.Metrics, error) {
	m, ok := ms.metrics[id]
	if !ok || m.MType != mtype {
		return nil, repository.ErrNotFound
	}

	return &m, nil
}
