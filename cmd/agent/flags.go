package main

import (
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/spf13/pflag"

	"github.com/vrnvgasu/metrics/internal/config"
)

type agentConfigValue struct {
	cnf *config.AgentCnf
}

func (v *agentConfigValue) String() string { return v.cnf.ConfigFile }
func (v *agentConfigValue) Type() string   { return "string" }
func (v *agentConfigValue) Set(path string) error {
	fileCnf, err := config.LoadAgentFileCnf(path)
	if err != nil {
		return err
	}
	if err = v.cnf.ApplyFile(fileCnf); err != nil {
		return err
	}
	v.cnf.ConfigFile = path
	return nil
}

func parseFlags() *config.AgentCnf {
	cnf := config.NewAgentCnf()

	cfgVal := &agentConfigValue{cnf: cnf}
	if path := os.Getenv("CONFIG"); path != "" {
		if err := cfgVal.Set(path); err != nil {
			log.Fatal("error loading config file: ", err)
		}
	}

	pflag.VarP(cfgVal, "config", "c", "path to config file")
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
