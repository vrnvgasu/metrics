package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/internal/repository"
)

func (s *Storage) GetByTypeAndID(ctx context.Context, mtype, id string) (*models.Metrics, error) {
	var m models.Metrics
	q := `select metric_id, type, delta, value, hash from metrics where type = $1 and metric_id = $2 limit 1`
	err := s.getDB(ctx).QueryRowContext(ctx, q, mtype, id).Scan(&m.ID, &m.MType, &m.Delta, &m.Value, &m.Hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}

		return nil, fmt.Errorf("postgres.GetByTypeAndID QueryRowContext: %w", err)
	}

	return &m, nil
}
