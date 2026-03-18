package errors

import (
	"errors"
	"time"
)

type RetryableError struct {
	Err error
}

func NewRetryableError(err error) error {
	return &RetryableError{Err: err}
}

func (e *RetryableError) Error() string {
	return e.Err.Error()
}

func (e *RetryableError) Unwrap() error {
	return e.Err
}

const (
	maxRetries         = 3
	startRetryInterval = 1 * time.Second
	addRetryPeriod     = 2 * time.Second
)

func Retry(f func() error) error {
	var (
		retryInterval = startRetryInterval
		attempt       = 0
	)
	for {
		err := f()
		if err == nil {
			return nil
		}

		var re *RetryableError
		if !errors.As(err, &re) {
			return re.Err
		}
		if attempt >= maxRetries {
			return err
		}

		time.Sleep(retryInterval)

		attempt++
		retryInterval += addRetryPeriod
	}
}
