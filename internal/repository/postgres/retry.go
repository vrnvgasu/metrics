package postgres

import (
	"github.com/vrnvgasu/metrics/pkg/retry"
)

func WithRetry(fn func() error, classifier *PostgresErrorClassifier) error {
	return retry.Retry(func() error {
		err := fn()
		if err == nil {
			return nil
		}

		// Классифицируем ошибку
		if classifier != nil && classifier.Classify(err) == NonRetriable {
			return err // Не-retryable ошибка, не повторяем
		}

		// Retryable ошибка
		return retry.NewRetryableError(err)
	})
}
