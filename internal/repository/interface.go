// Package repository определяет интерфейсы хранилища метрик.
package repository

import (
	"context"

	models "github.com/vrnvgasu/metrics/internal/model"
)

// Transaction — маркерный интерфейс транзакции.
type Transaction interface {
}

// Storage — интерфейс хранилища метрик.
//
//go:generate mockgen -destination=./mocks/mock.go . Storage
type Storage interface {
	// List возвращает все метрики.
	List(context.Context) (models.MetricsList, error)
	// Save сохраняет или обновляет метрику.
	Save(ctx context.Context, m *models.Metrics) error
	// GetByTypeAndID возвращает метрику по типу и имени.
	GetByTypeAndID(ctx context.Context, mtype models.MetricType, id string) (*models.Metrics, error)
	// Ping проверяет доступность хранилища.
	Ping(context.Context) error
	// DoInTransaction выполняет fn внутри транзакции.
	DoInTransaction(ctx context.Context, fn func(context.Context) error) error
}
