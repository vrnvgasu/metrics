// Package handler содержит HTTP-хендлеры сервера метрик.
package handler

import (
	"context"

	"github.com/vrnvgasu/metrics/internal/config"
	models "github.com/vrnvgasu/metrics/internal/model"
)

// MetricService — интерфейс сервиса метрик.
//
//go:generate mockgen -destination=./mocks/metric.go . MetricService
type MetricService interface {
	// CreateOrUpdate создает или обновляет список метрик.
	CreateOrUpdate(context.Context, models.MetricsList) error
	// FindByTypeAndID возвращает метрику по типу и имени.
	FindByTypeAndID(ctx context.Context, mtype models.MetricType, id string) (*models.Metrics, error)
	// AllMetrics возвращает все метрики.
	AllMetrics(context.Context) (models.MetricsList, error)
}

// HealthService — интерфейс проверки доступности хранилища.
//
//go:generate mockgen -destination=./mocks/health.go . HealthService
type HealthService interface {
	// CheckPing проверяет соединение с хранилищем.
	CheckPing(ctx context.Context) error
}

// Publisher — интерфейс публикации событий аудита.
type Publisher interface {
	// Notify отправляет уведомление об обновлении метрик.
	Notify(ctx context.Context, metricIDList []string, ip string)
}

// Handler содержит зависимости HTTP-хендлеров.
type Handler struct {
	MetricService MetricService
	HealthService HealthService

	publisher Publisher

	cfg *config.ServerCnf
}

// NewHandler создает Handler с переданными зависимостями.
func NewHandler(m MetricService, h HealthService, p Publisher, cfg *config.ServerCnf) *Handler {
	return &Handler{
		MetricService: m,
		HealthService: h,
		publisher:     p,
		cfg:           cfg,
	}
}
