package service

import (
	"fmt"

	"github.com/vrnvgasu/metrics/internal/service/handler"
)

func (s *Service) Restore() error {
	if !s.cnf.Restore {
		return nil
	}

	consumer, err := handler.NewConsumer(s.cnf.FileStoragePath)
	if err != nil {
		return fmt.Errorf("server.Restore NewConsumer: %w", err)
	}

	list, err := consumer.ReadMetrics()
	if err != nil {
		return fmt.Errorf("server.Restore ReadMetrics: %w", err)
	}

	for _, m := range list {
		if err = s.storage.Add(m); err != nil {
			return fmt.Errorf("server.Restore Add: %w", err)
		}
	}

	return nil
}
