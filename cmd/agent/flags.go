package main

import (
	"github.com/spf13/pflag"

	"github.com/vrnvgasu/metrics/internal/config"
)

func parseFlags() *config.AgentCnf {
	cnf := &config.AgentCnf{}

	pflag.StringVarP(&cnf.Address, "address", "a", "localhost:8080", "address:port to listen on")
	pflag.IntVarP(&cnf.ReportInterval, "reportInterval", "r", 10, "frequency of sending metrics to the server")
	pflag.IntVarP(&cnf.PollInterval, "pollInterval", "p", 2, "frequency of polling metrics from the package")

	pflag.Parse()

	return cnf
}
