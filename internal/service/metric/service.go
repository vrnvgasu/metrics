// Package metric реализует бизнес-логику работы с метриками.
package metric

import (
	"github.com/vrnvgasu/metrics/internal/repository"
)

// Service реализует бизнес-логику работы с метриками.
type Service struct {
	storage repository.Storage
}

// NewService создает Service с переданным хранилищем.
func NewService(storage repository.Storage) *Service {
	return &Service{
		storage: storage,
	}
}
