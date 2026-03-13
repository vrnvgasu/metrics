package postgres

import (
	"context"
	"fmt"

	models "github.com/vrnvgasu/metrics/internal/model"
)

func (s *Storage) Save(ctx context.Context, m *models.Metrics) error {
	q := `insert into metrics (id, type, delta, value, hash) values ($1, $2, $3, $4, $5)
		on conflict (id, type) do update set 
					delta = EXCLUDED.delta,
					value = EXCLUDED.value,
					hash = EXCLUDED.hash`
	_, err := s.DB.ExecContext(ctx, q, m.ID, m.MType, m.Delta, m.Value, m.Hash)
	if err != nil {
		return fmt.Errorf("postgres.Save ExecContext: %w", err)
	}

	return nil
}
