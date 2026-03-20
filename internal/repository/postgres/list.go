package postgres

import (
	"context"
	"fmt"

	models "github.com/vrnvgasu/metrics/internal/model"
)

func (s *Storage) List(ctx context.Context) (models.MetricsList, error) {
	list := make(models.MetricsList, 0)
	q := `select metric_id, type, delta, value, hash from metrics`

	rows, err := s.getDB(ctx).QueryContext(ctx, q)
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
