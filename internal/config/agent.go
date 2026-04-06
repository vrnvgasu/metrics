package config

import (
	"fmt"
)

type AgentCnf struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	Key            string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
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
