package store

import (
	"context"
	"fmt"
	"time"
)

func (s *Service) StoreInterval(ctx context.Context) error {
	for !s.stop {
		time.Sleep(time.Duration(s.cnf.StoreInterval) * time.Second)

		select {
		case <-ctx.Done():
			return nil
		default:
			if err := s.flushOnFile(ctx); err != nil {
				return fmt.Errorf("s.StoreInterval flushOnFile: %w", err)
			}
		}
	}

	return nil
}

func (s *Service) flushOnFile(ctx context.Context) error {
	list := s.storage.List(ctx)
	if err := s.producer.WriteMetrics(&list); err != nil {
		return fmt.Errorf("s.flush WriteMetrics: %w", err)
	}

	return nil
}
