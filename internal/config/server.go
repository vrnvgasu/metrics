package config

type ServerCnf struct {
	Address string `env:"ADDRESS"`
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
