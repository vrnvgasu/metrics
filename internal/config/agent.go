package config

import (
	"fmt"
)

type AgentCnf struct {
	Address        string
	ReportInterval int
	PollInterval   int
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
