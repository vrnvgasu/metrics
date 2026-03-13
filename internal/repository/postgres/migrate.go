package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/pressly/goose/v3"

	"github.com/vrnvgasu/metrics/migrations"
)

const (
	migrateTimeout = 30 * time.Second
)

func (s *Storage) Migrate(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, migrateTimeout)
	defer cancel()

	goose.SetBaseFS(migrations.MigrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("postgres.Migrate SetDialect: %w", err)
	}

	if err := goose.UpContext(ctx, s.DB, "."); err != nil {
		return fmt.Errorf("postgres.Migrate Up: %w", err)
	}

	return nil
}
