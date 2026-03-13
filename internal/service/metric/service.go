package metric

import (
	"github.com/vrnvgasu/metrics/internal/repository"
)

type Service struct {
	storage repository.Storage
}

func NewService(storage repository.Storage) *Service {
	return &Service{
		storage: storage,
	}
}
