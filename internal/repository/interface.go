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
	GetByTypeAndID(ctx context.Context, mtype, id string) (*models.Metrics, error)
	Ping(context.Context) error

	WithTx(ctx context.Context) (Storage, error)
	Commit() error
	Rollback() error
}
