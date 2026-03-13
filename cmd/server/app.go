package main

import (
	"context"
	"fmt"

	"github.com/vrnvgasu/metrics/internal/config"
	"github.com/vrnvgasu/metrics/internal/repository"
	"github.com/vrnvgasu/metrics/internal/repository/mem"
	"github.com/vrnvgasu/metrics/internal/repository/postgres"
)

func initStorage(ctx context.Context, cnf *config.ServerCnf) (dbStorage repository.Storage, err error) {
	switch {
	case cnf.DatabaseDSN != "":
		provider := postgres.NewStorage()
		if err = provider.Start(ctx, cnf.DatabaseDSN); err != nil {
			return nil, fmt.Errorf("init storage postgres start: %w", err)
		}
		if err = provider.Migrate(ctx); err != nil {
			return nil, fmt.Errorf("init storage postgres migrate: %w", err)
		}

		return provider, nil
	default:
		return mem.NewMemStorage(), nil
	}
}
