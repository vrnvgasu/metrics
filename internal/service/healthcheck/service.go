// Package healthcheck реализует проверку доступности хранилища.
package healthcheck

import (
	"github.com/vrnvgasu/metrics/internal/repository"
)

type Service struct {
	DB repository.Storage
}

func NewService(db repository.Storage) *Service {
	return &Service{
		DB: db,
	}
}
