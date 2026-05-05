package repository

import (
	"context"

	models "github.com/vrnvgasu/metrics/internal/model"
)

type Transaction interface {
}

//go:generate mockgen -destination=./mocks/mock.go . Storage
type Storage interface {
	List(context.Context) (models.MetricsList, error)
	Save(ctx context.Context, m *models.Metrics) error
	GetByTypeAndID(ctx context.Context, mtype models.MetricType, id string) (*models.Metrics, error)
	Ping(context.Context) error

	DoInTransaction(ctx context.Context, fn func(context.Context) error) error
}
