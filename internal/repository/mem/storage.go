package mem

import (
	"sync"

	models "github.com/vrnvgasu/metrics/internal/model"
)

type Storage struct {
	mu      sync.RWMutex
	metrics models.MetricsMap
}

func NewMemStorage() *Storage {
	return &Storage{
		metrics: make(models.MetricsMap),
	}
}
