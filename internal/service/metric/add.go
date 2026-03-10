package metric

import (
	"context"
	"errors"

	models "github.com/vrnvgasu/metrics/internal/model"
	seriveerrors "github.com/vrnvgasu/metrics/internal/service/errors"
)

func (s *Service) CreateOrUpdate(ctx context.Context, m *models.Metrics) error {
	err := s.storage.Add(ctx, m)
	if errors.Is(err, errors.ErrUnsupported) {
		return seriveerrors.UnprocessableEntity(err.Error())
	}

	return err
}
