package mem

import (
	"context"
)

func (s *Storage) DoInTransaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}
