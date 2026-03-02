package service

import (
	"fmt"

	"github.com/vrnvgasu/metrics/internal/config"
	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/internal/service/handler"
)

type Storage interface {
	List() models.MetricsList
	Add(m models.Metrics) error
}

type Service struct {
	stop bool

	storage  Storage
	cnf      config.ServerCnf
	producer *handler.Producer
}

func NewService(storage Storage, cnf config.ServerCnf) (*Service, error) {
	p, err := handler.NewProducer(cnf.FileStoragePath)
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
