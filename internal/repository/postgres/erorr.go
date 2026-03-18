package postgres

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type PGErrorClassification int

const (
	// NonRetriable - операцию не следует повторять
	NonRetriable PGErrorClassification = iota

	// Retriable - операцию можно повторить
	Retriable
)

type PostgresErrorClassifier struct{}

func NewPostgresErrorClassifier() *PostgresErrorClassifier {
	return &PostgresErrorClassifier{}
}

func (c *PostgresErrorClassifier) Classify(err error) PGErrorClassification {
	if err == nil {
		return NonRetriable
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return ClassifyPgError(pgErr)
	}

	return NonRetriable
}

func ClassifyPgError(pgErr *pgconn.PgError) PGErrorClassification {
	switch pgErr.Code {
	// Класс 08 - Ошибки соединения
	case pgerrcode.ConnectionException,
		pgerrcode.ConnectionDoesNotExist,
		pgerrcode.ConnectionFailure:
		return Retriable
	default:
		return NonRetriable
	}
}
