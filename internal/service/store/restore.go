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

	for _, m := range list {
		if err = s.storage.Add(ctx, &m); err != nil {
			return fmt.Errorf("server.Restore Add: %w", err)
		}
	}

	return nil
}
