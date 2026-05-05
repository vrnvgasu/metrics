package handler

import (
	"context"

	"github.com/vrnvgasu/metrics/internal/config"
	models "github.com/vrnvgasu/metrics/internal/model"
)

type MetricService interface {
	CreateOrUpdate(context.Context, []*models.Metrics) error
	FindByTypeAndID(ctx context.Context, mtype models.MetricType, id string) (*models.Metrics, error)
	AllMetrics(context.Context) (models.MetricsList, error)
}

type HealthService interface {
	CheckPing(ctx context.Context) error
}

type Handler struct {
	MetricService MetricService
	HealthService HealthService

	cfg *config.ServerCnf
}

func NewHandler(m MetricService, h HealthService, cfg *config.ServerCnf) *Handler {
	return &Handler{
		MetricService: m,
		HealthService: h,
		cfg:           cfg,
	}
}
