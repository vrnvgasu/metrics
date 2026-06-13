// Package audit реализует сервис аудита событий (Observer/Publisher).
package audit

import (
	"context"
	"fmt"

	"github.com/vrnvgasu/metrics/internal/config"
)

//go:generate mockgen -destination=./mocks/mock.go . Observer
type Observer interface {
	Observe(context.Context, EventMessage) error
	Close() error
}

type EventMessage struct {
	Ts        int64    `json:"ts"`
	Metrics   []string `json:"metrics"`
	IpAddress string   `json:"ip_address"`
}

func NewAudit(cfg *config.ServerCnf) (*Event, error) {
	publisher := NewEvent()

	if cfg.AuditFile != "" {
		audit, err := NewFileAudit(cfg.AuditFile)
		if err != nil {
			return nil, fmt.Errorf("could not create audit file: %w", err)
		}

		publisher.Register(audit)
	}
	if cfg.AuditURL != "" {
		audit := NewURLAudit(cfg.AuditURL)
		publisher.Register(audit)
	}

	return publisher, nil
}
