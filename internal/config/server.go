package config

// ServerCnf — конфигурация сервера метрик.
type ServerCnf struct {
	Address         string `env:"ADDRESS"`
	LogLevel        string `env:"LOG_LEVEL"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	Key             string `env:"KEY"`

	AuditFile string `env:"AUDIT_FILE"`
	AuditURL  string `env:"AUDIT_URL"`
}

func (s *ServerCnf) String() string {
	return "Address: " + s.Address
}

func (s *ServerCnf) Set(value string) error {
	s.Address = value

	return nil
}

func (s *ServerCnf) Type() string {
	return "server configuration"
}
