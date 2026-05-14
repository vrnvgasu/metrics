// Package retry предоставляет механизм повторных попыток для временных ошибок.
package retry

import (
	"errors"
	"time"
)

// RetryConfig — настройки стратегии повторных попыток.
type RetryConfig struct {
	MaxRetries         int
	StartRetryInterval time.Duration
	AddRetryPeriod     time.Duration
}

// DefaultConfig возвращает конфиг по умолчанию: 3 попытки, интервал 1с с шагом 2с.
func DefaultConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:         3,
		StartRetryInterval: 1 * time.Second,
		AddRetryPeriod:     2 * time.Second,
	}
}

// RetryableError — обертка над ошибкой, сигнализирующая о том, что операцию можно повторить.
type RetryableError struct {
	Err error
}

// NewRetryableError оборачивает err в RetryableError.
func NewRetryableError(err error) error {
	return &RetryableError{Err: err}
}

func (e *RetryableError) Error() string {
	return e.Err.Error()
}

func (e *RetryableError) Unwrap() error {
	return e.Err
}

// RetryWithSettings выполняет f, повторяя при RetryableError согласно cnf.
func RetryWithSettings(f func() error, cnf *RetryConfig) error {
	if cnf == nil {
		cnf = DefaultConfig()
	}

	var (
		retryInterval = cnf.StartRetryInterval
		attempt       = 0
	)
	for {
		err := f()
		if err == nil {
			return nil
		}

		var re *RetryableError
		if !errors.As(err, &re) {
			return err
		}
		if attempt >= cnf.MaxRetries {
			return err
		}

		time.Sleep(retryInterval)

		attempt++
		retryInterval += cnf.AddRetryPeriod
	}
}

// Retry выполняет f с настройками по умолчанию.
func Retry(f func() error) error {
	return RetryWithSettings(f, DefaultConfig())
}
