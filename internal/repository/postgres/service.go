package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"

	serviceerrors "github.com/vrnvgasu/metrics/internal/service/errors"
)

type Storage struct {
	*sql.DB
	*sql.Tx
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) Start(ctx context.Context, dsn string) (err error) {
	classifier := NewPostgresErrorClassifier()

	err = serviceerrors.Retry(func() error {
		s.DB, err = sql.Open("pgx", dsn)
		if err != nil {
			if classifier.Classify(err) == NonRetriable {
				return err
			}

			return serviceerrors.NewRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("postgres service.Start Open: %w", err)
	}

	err = serviceerrors.Retry(func() error {
		err = s.DB.PingContext(ctx)
		if err != nil {
			if classifier.Classify(err) == NonRetriable {
				return err
			}

			return serviceerrors.NewRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("postgres service.Start Ping: %w", err)
	}

	err = serviceerrors.Retry(func() error {
		err = s.Migrate(ctx)
		if err != nil {
			if classifier.Classify(err) == NonRetriable {
				return err
			}

			return serviceerrors.NewRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("postgres service.Start Migrate: %w", err)
	}

	return nil
}

func (s *Storage) Stop() error {
	return s.DB.Close()
}
