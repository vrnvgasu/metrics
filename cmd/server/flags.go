package main

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/spf13/pflag"

	"github.com/vrnvgasu/metrics/internal/config"
)

func parseFlags() *config.ServerCnf {
	cnf := &config.ServerCnf{}

	pflag.StringVarP(&cnf.Address, "address", "a", "localhost:8080", "address:port to listen on")
	pflag.StringVarP(&cnf.LogLevel, "loglevel", "l", "info", "log level")
	pflag.Parse()

	if err := env.Parse(cnf); err != nil {
		log.Fatal("error parsing server env: ", err)
	}

	return cnf
}
