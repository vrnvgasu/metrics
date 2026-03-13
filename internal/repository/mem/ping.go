package mem

import "context"

func (s *Storage) Ping(_ context.Context) error {
	return nil
}
