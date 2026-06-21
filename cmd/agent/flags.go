package main

import (
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/spf13/pflag"

	"github.com/vrnvgasu/metrics/internal/config"
)

func parseFlags() *config.AgentCnf {
	cnf := config.NewAgentCnf()

	pflag.StringVarP(&cnf.ConfigFile, "config", "c", os.Getenv("CONFIG"), "path to config file")
	pflag.Parse()

	if cnf.ConfigFile != "" {
		fileCnf, err := config.LoadAgentFileCnf(cnf.ConfigFile)
		if err != nil {
			log.Fatal("error loading config file: ", err)
		}
		cnf.ApplyFile(fileCnf)
	}

	pflag.StringVarP(&cnf.Address, "address", "a", cnf.Address, "address:port to listen on")
	pflag.IntVarP(&cnf.ReportInterval, "reportInterval", "r", cnf.ReportInterval, "frequency of sending metrics to the server")
	pflag.IntVarP(&cnf.PollInterval, "pollInterval", "p", cnf.PollInterval, "frequency of polling metrics from the package")
	pflag.StringVarP(&cnf.Key, "key", "k", cnf.Key, "hash key")
	pflag.IntVarP(&cnf.RateLimit, "rateLimit", "l", cnf.RateLimit, "rate limit")
	pflag.StringVar(&cnf.CryptoKey, "crypto-key", cnf.CryptoKey, "path to public key file for encryption")
	pflag.Parse()

	if err := env.Parse(cnf); err != nil {
		log.Fatal("error parsing agent env: ", err)
	}

	return cnf
}
