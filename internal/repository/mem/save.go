package mem

import (
	"context"

	models "github.com/vrnvgasu/metrics/internal/model"
)

func (s *Storage) Save(_ context.Context, m *models.Metrics) error {
	s.metrics[m.ID] = *m

	return nil
}
