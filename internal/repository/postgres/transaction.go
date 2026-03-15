package postgres

import (
	"context"
	"fmt"

	"github.com/vrnvgasu/metrics/internal/repository"
)

func (s *Storage) WithTx(_ context.Context) (repository.Storage, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("postgres.BeginTx Begin: %w", err)
	}

	return &Storage{
		DB: s.DB,
		Tx: tx,
	}, nil
}

func (s *Storage) Commit() error {
	return s.Tx.Commit()
}

func (s *Storage) Rollback() error {
	return s.Tx.Rollback()
}
