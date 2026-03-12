package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	*sql.DB
}

func NewService() *Storage {
	return &Storage{}
}

func (s *Storage) Start(ctx context.Context, dsn string) (err error) {
	s.DB, err = sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("postgres service.Start Open: %w", err)
	}

	if err = s.DB.PingContext(ctx); err != nil {
		return fmt.Errorf("postgres service.Start Ping: %w", err)
	}

	return nil
}

func (s *Storage) Stop() error {
	return s.DB.Close()
}
