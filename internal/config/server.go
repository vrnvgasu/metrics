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
	CryptoKey       string `env:"CRYPTO_KEY"`
	ConfigFile      string `env:"CONFIG"`

	AuditFile string `env:"AUDIT_FILE"`
	AuditURL  string `env:"AUDIT_URL"`
}

func NewServerCnf() *ServerCnf {
	return &ServerCnf{
		Address:         "localhost:8080",
		LogLevel:        "info",
		StoreInterval:   300,
		FileStoragePath: "store.json",
		Restore:         true,
	}
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
