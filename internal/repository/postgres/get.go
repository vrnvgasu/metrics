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
	var (
		m   models.Metrics
		err error
	)
	q := `select * from metrics where type = $1 and id = $2 limit 1`

	if s.Tx != nil {
		err = s.Tx.QueryRowContext(ctx, q, mtype, id).Scan(&m.ID, &m.MType, &m.Delta, &m.Value, &m.Hash)
	} else {
		err = s.DB.QueryRowContext(ctx, q, mtype, id).Scan(&m.ID, &m.MType, &m.Delta, &m.Value, &m.Hash)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}

		return nil, fmt.Errorf("postgres.GetByTypeAndID QueryRowContext: %w", err)
	}

	return &m, nil
}
