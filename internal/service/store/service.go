package store

import (
	"context"
	"fmt"

	"github.com/vrnvgasu/metrics/internal/config"
	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/internal/repository"
	"github.com/vrnvgasu/metrics/internal/service/store/provider"
)

type MetricService interface {
	CreateOrUpdate(context.Context, []*models.Metrics) error
}

type Service struct {
	stop bool

	storage       repository.Storage
	cnf           config.ServerCnf
	producer      *provider.Producer
	metricService MetricService
}

func NewService(
	storage repository.Storage, ms MetricService, cnf config.ServerCnf,
) (*Service, error) {
	p, err := provider.NewProducer(cnf.FileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("server.StoreInterval NewProducer: %w", err)
	}

	return &Service{
		storage:       storage,
		cnf:           cnf,
		producer:      p,
		metricService: ms,
	}, nil
}

func (s *Service) Stop() error {
	s.stop = true

	return s.producer.Close()
}
