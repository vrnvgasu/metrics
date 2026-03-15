package postgres

import (
	"context"
	"database/sql"
	"fmt"

	models "github.com/vrnvgasu/metrics/internal/model"
)

func (s *Storage) List(ctx context.Context) (models.MetricsList, error) {
	list := make(models.MetricsList, 0)
	q := `select * from metrics`

	var (
		rows *sql.Rows
		err  error
	)
	if s.Tx != nil {
		rows, err = s.Tx.QueryContext(ctx, q)
	} else {
		rows, err = s.DB.QueryContext(ctx, q)
	}

	if err != nil {
		return nil, fmt.Errorf("postgres.List QueryContext: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var m models.Metrics
		err = rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value, &m.Hash)
		if err != nil {
			return nil, fmt.Errorf("postgres.List Scan: %w", err)
		}

		list = append(list, m)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres.List Err: %w", err)
	}

	return list, nil
}
