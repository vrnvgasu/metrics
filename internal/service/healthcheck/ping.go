package healthcheck

import (
	"context"

	seriveerrors "github.com/vrnvgasu/metrics/internal/service/errors"
)

func (s *Service) CheckPing(ctx context.Context) error {
	if err := s.DB.Ping(ctx); err != nil {
		return seriveerrors.InternalError("healthcheck: service unavailable")
	}

	return nil
}
