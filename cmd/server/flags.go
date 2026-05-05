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
	pflag.IntVarP(&cnf.StoreInterval, "storeInterval", "i", 300, "store interval in seconds")
	pflag.StringVarP(&cnf.FileStoragePath, "fileStoragePath", "f", "store.json", "file storage path")
	pflag.BoolVarP(&cnf.Restore, "restore", "r", true, "restore storage data")
	pflag.StringVarP(&cnf.DatabaseDSN, "database_dsn", "d", "", "database DSN in format: host=localhost port=5432 user=metrics password=metrics dbname=metrics sslmode=disable")
	pflag.StringVarP(&cnf.Key, "key", "k", "", "hash key")
	pflag.Parse()

	if err := env.Parse(cnf); err != nil {
		log.Fatal("error parsing server env: ", err)
	}

	return cnf
}
