package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type ServerFileCnf struct {
	Address         string `json:"address"`
	Restore         *bool  `json:"restore"`
	StoreInterval   string `json:"store_interval"`
	FileStoragePath string `json:"store_file"`
	DatabaseDSN     string `json:"database_dsn"`
	CryptoKey       string `json:"crypto_key"`
	TrustedSubnet   string `json:"trusted_subnet"`
}

type AgentFileCnf struct {
	Address        string `json:"address"`
	ReportInterval string `json:"report_interval"`
	PollInterval   string `json:"poll_interval"`
	CryptoKey      string `json:"crypto_key"`
}

func LoadServerFileCnf(path string) (*ServerFileCnf, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read server config file: %w", err)
	}

	cnf := &ServerFileCnf{}
	if err = json.Unmarshal(data, cnf); err != nil {
		return nil, fmt.Errorf("failed to parse server config file: %w", err)
	}

	return cnf, nil
}

func LoadAgentFileCnf(path string) (*AgentFileCnf, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read agent config file: %w", err)
	}

	cnf := &AgentFileCnf{}
	if err = json.Unmarshal(data, cnf); err != nil {
		return nil, fmt.Errorf("failed to parse agent config file: %w", err)
	}

	return cnf, nil
}

func (cnf *ServerCnf) ApplyFile(f *ServerFileCnf) error {
	if f.Address != "" {
		cnf.Address = f.Address
	}
	if f.Restore != nil {
		cnf.Restore = *f.Restore
	}
	if f.StoreInterval != "" {
		secs, err := parseDurationSec(f.StoreInterval)
		if err != nil {
			return fmt.Errorf("store_interval: %w", err)
		}
		cnf.StoreInterval = secs
	}
	if f.FileStoragePath != "" {
		cnf.FileStoragePath = f.FileStoragePath
	}
	if f.DatabaseDSN != "" {
		cnf.DatabaseDSN = f.DatabaseDSN
	}
	if f.CryptoKey != "" {
		cnf.CryptoKey = f.CryptoKey
	}
	if f.TrustedSubnet != "" {
		cnf.TrustedSubnet = f.TrustedSubnet
	}

	return nil
}

func (cnf *AgentCnf) ApplyFile(f *AgentFileCnf) error {
	if f.Address != "" {
		cnf.Address = f.Address
	}
	if f.ReportInterval != "" {
		secs, err := parseDurationSec(f.ReportInterval)
		if err != nil {
			return fmt.Errorf("report_interval: %w", err)
		}
		cnf.ReportInterval = secs
	}
	if f.PollInterval != "" {
		secs, err := parseDurationSec(f.PollInterval)
		if err != nil {
			return fmt.Errorf("poll_interval: %w", err)
		}
		cnf.PollInterval = secs
	}
	if f.CryptoKey != "" {
		cnf.CryptoKey = f.CryptoKey
	}

	return nil
}

func parseDurationSec(s string) (int, error) {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("invalid duration %q: %w", s, err)
	}

	return int(d.Seconds()), nil
}
