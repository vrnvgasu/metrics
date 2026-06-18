package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type ServerFileCnf struct {
	Address         string `json:"address"`
	Restore         *bool  `json:"restore"`
	StoreInterval   string `json:"store_interval"`
	FileStoragePath string `json:"store_file"`
	DatabaseDSN     string `json:"database_dsn"`
	CryptoKey       string `json:"crypto_key"`
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

func (cnf *ServerCnf) ApplyFile(f *ServerFileCnf) {
	if f.Address != "" {
		cnf.Address = f.Address
	}
	if f.Restore != nil {
		cnf.Restore = *f.Restore
	}
	if f.StoreInterval != "" {
		if secs, ok := parseDurationSec(f.StoreInterval); ok {
			cnf.StoreInterval = secs
		}
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
}

func (cnf *AgentCnf) ApplyFile(f *AgentFileCnf) {
	if f.Address != "" {
		cnf.Address = f.Address
	}
	if f.ReportInterval != "" {
		if secs, ok := parseDurationSec(f.ReportInterval); ok {
			cnf.ReportInterval = secs
		}
	}
	if f.PollInterval != "" {
		if secs, ok := parseDurationSec(f.PollInterval); ok {
			cnf.PollInterval = secs
		}
	}
	if f.CryptoKey != "" {
		cnf.CryptoKey = f.CryptoKey
	}
}

// FindConfigPath ищет путь к файлу конфигурации: сначала в аргументах (-c/--config),
// затем в переменной окружения CONFIG.
func FindConfigPath() string {
	args := os.Args[1:]
	for i, arg := range args {
		switch {
		case (arg == "-c" || arg == "--config") && i+1 < len(args):
			return args[i+1]
		case strings.HasPrefix(arg, "-c="):
			return strings.TrimPrefix(arg, "-c=")
		case strings.HasPrefix(arg, "--config="):
			return strings.TrimPrefix(arg, "--config=")
		}
	}
	return os.Getenv("CONFIG")
}

func parseDurationSec(s string) (int, bool) {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, false
	}

	return int(d.Seconds()), true
}
