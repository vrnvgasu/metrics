package handler

import models "github.com/vrnvgasu/metrics/internal/model"

type Storage interface {
	Add(m models.Metrics) error
	GetByTypeAndID(mtype, id string) (models.Metrics, error)
	Map() models.MetricsMap
}

type Handler struct {
	Storage Storage
}

func NewHandler(s Storage) *Handler {
	return &Handler{
		Storage: s,
	}
}
