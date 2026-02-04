package handler

import models "github.com/vrnvgasu/metrics/internal/model"

type Storage interface {
	Add(m models.Metrics) error
}

type Handler struct {
	Storage Storage
}

func NewHandler(s Storage) *Handler {
	return &Handler{
		Storage: s,
	}
}
