package store

import (
	"context"
	"fmt"

	"github.com/vrnvgasu/metrics/internal/service/store/provider"
)

func (s *Service) Restore(ctx context.Context) error {
	if !s.cnf.Restore {
		return nil
	}

	consumer, err := provider.NewConsumer(s.cnf.FileStoragePath)
	if err != nil {
		return fmt.Errorf("server.Restore NewConsumer: %w", err)
	}

	list, err := consumer.ReadMetrics()
	if err != nil {
		return fmt.Errorf("server.Restore ReadMetrics: %w", err)
	}

	if err = s.metricService.CreateOrUpdate(ctx, list); err != nil {
		return fmt.Errorf("server.Restore CreateOrUpdate: %w", err)
	}

	return nil
}
