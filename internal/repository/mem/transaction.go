package mem

import (
	"context"

	"github.com/vrnvgasu/metrics/internal/repository"
)

func (s *Storage) WithTx(ctx context.Context) (repository.Storage, error) {
	return s, nil
}

func (s *Storage) Commit() error {
	return nil
}

func (s *Storage) Rollback() error {
	return nil
}
