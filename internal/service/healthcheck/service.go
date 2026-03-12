package healthcheck

import "context"

type DBStorage interface {
	Ping(context.Context) error
}

type Service struct {
	DB DBStorage
}

func NewService(db DBStorage) *Service {
	return &Service{
		DB: db,
	}
}
