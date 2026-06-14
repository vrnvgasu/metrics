// Package config содержит конфигурационные структуры сервера и агента.
package config

import (
	"fmt"
)

// AgentCnf — конфигурация агента сбора метрик.
type AgentCnf struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	Key            string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
	CryptoKey      string `env:"CRYPTO_KEY"`
}

func (a *AgentCnf) String() string {
	return fmt.Sprintf("%v", *a)
}

func (a *AgentCnf) Set(value string) error {
	a.Address = value

	return nil
}

func (a *AgentCnf) Type() string {
	return "agent configuration"
}
