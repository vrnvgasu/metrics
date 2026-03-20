package metric

import (
	"context"
	"errors"
	"fmt"

	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/internal/repository"
	serviceerrors "github.com/vrnvgasu/metrics/internal/service/errors"
)

func (s *Service) FindByTypeAndID(ctx context.Context, mtype, id string) (*models.Metrics, error) {
	m, err := s.storage.GetByTypeAndID(ctx, mtype, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, serviceerrors.NotFoundError(fmt.Errorf("mtype not found: %w", err))
	}

	return m, err
}
