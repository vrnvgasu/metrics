package metric

import (
	"context"

	models "github.com/vrnvgasu/metrics/internal/model"
)

type Storage interface {
	List(context.Context) models.MetricsList
	Add(context.Context, *models.Metrics) error
	GetByTypeAndID(ctx context.Context, mtype, id string) (*models.Metrics, error)
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}
