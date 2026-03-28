package main

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/spf13/pflag"

	"github.com/vrnvgasu/metrics/internal/config"
)

func parseFlags() *config.AgentCnf {
	cnf := &config.AgentCnf{}

	pflag.StringVarP(&cnf.Address, "address", "a", "localhost:8080", "address:port to listen on")
	pflag.IntVarP(&cnf.ReportInterval, "reportInterval", "r", 10, "frequency of sending metrics to the server")
	pflag.IntVarP(&cnf.PollInterval, "pollInterval", "p", 2, "frequency of polling metrics from the package")
	pflag.StringVarP(&cnf.Key, "key", "k", "", "hash key")
	pflag.Parse()

	if err := env.Parse(cnf); err != nil {
		log.Fatal("error parsing agent env: ", err)
	}

	return cnf
}
