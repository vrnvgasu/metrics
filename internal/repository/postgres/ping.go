package postgres

import "context"

func (s *Storage) Ping(ctx context.Context) error {
	return s.DB.PingContext(ctx)
}
