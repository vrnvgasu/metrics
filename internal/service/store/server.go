package store

import (
	"context"
	"fmt"

	"github.com/vrnvgasu/metrics/internal/config"
	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/internal/service/store/provider"
)

type Storage interface {
	List(context.Context) models.MetricsList
	Add(context.Context, *models.Metrics) error
}

type Service struct {
	stop bool

	storage  Storage
	cnf      config.ServerCnf
	producer *provider.Producer
}

func NewService(storage Storage, cnf config.ServerCnf) (*Service, error) {
	p, err := provider.NewProducer(cnf.FileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("server.StoreInterval NewProducer: %w", err)
	}

	return &Service{
		storage:  storage,
		cnf:      cnf,
		producer: p,
	}, nil
}

func (s *Service) Stop() error {
	s.stop = true

	return s.producer.Close()
}
