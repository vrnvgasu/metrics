package main

import (
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/spf13/pflag"

	"github.com/vrnvgasu/metrics/internal/config"
)

func parseFlags() *config.ServerCnf {
	cnf := config.NewServerCnf()

	pflag.StringVarP(&cnf.ConfigFile, "config", "c", os.Getenv("CONFIG"), "path to config file")
	pflag.Parse()

	if cnf.ConfigFile != "" {
		fileCnf, err := config.LoadServerFileCnf(cnf.ConfigFile)
		if err != nil {
			log.Fatal("error loading config file: ", err)
		}
		cnf.ApplyFile(fileCnf)
	}

	pflag.StringVarP(&cnf.Address, "address", "a", cnf.Address, "address:port to listen on")
	pflag.StringVarP(&cnf.LogLevel, "loglevel", "l", cnf.LogLevel, "log level")
	pflag.IntVarP(&cnf.StoreInterval, "storeInterval", "i", cnf.StoreInterval, "store interval in seconds")
	pflag.StringVarP(&cnf.FileStoragePath, "fileStoragePath", "f", cnf.FileStoragePath, "file storage path")
	pflag.BoolVarP(&cnf.Restore, "restore", "r", cnf.Restore, "restore storage data")
	pflag.StringVarP(&cnf.DatabaseDSN, "database_dsn", "d", cnf.DatabaseDSN, "database DSN in format: host=localhost port=5432 user=metrics password=metrics dbname=metrics sslmode=disable")
	pflag.StringVarP(&cnf.Key, "key", "k", cnf.Key, "hash key")
	pflag.StringVar(&cnf.AuditFile, "audit-file", cnf.AuditFile, "path to audit file")
	pflag.StringVar(&cnf.AuditURL, "audit-url", cnf.AuditURL, "audit server URL")
	pflag.StringVar(&cnf.CryptoKey, "crypto-key", cnf.CryptoKey, "path to private key file for decryption")
	pflag.Parse()

	if err := env.Parse(cnf); err != nil {
		log.Fatal("error parsing server env: ", err)
	}

	return cnf
}
