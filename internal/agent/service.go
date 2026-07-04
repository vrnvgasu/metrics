// Package agent реализует агент сбора и отправки метрик на сервер.
package agent

import (
	"crypto/rsa"
	"io"
	"net/http"
	"sync/atomic"

	models "github.com/vrnvgasu/metrics/internal/model"
	pb "github.com/vrnvgasu/metrics/internal/proto"
)

// Client — интерфейс HTTP-клиента агента.
//
//go:generate mockgen -destination=./mocks/mock.go . Client
type Client interface {
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
	Do(req *http.Request) (resp *http.Response, err error)
}

// Agent собирает и отправляет метрики на сервер.
type Agent struct {
	Client     Client
	Metrics    chan models.Metrics
	publicKey  *rsa.PublicKey
	realIP     string
	grpcClient pb.MetricsClient
	pollCount  atomic.Int64
}

// NewAgent создает агента с буфером канала метрик размером bufSize.
func NewAgent(client Client, bufSize int) *Agent {
	return &Agent{
		Client:  client,
		Metrics: make(chan models.Metrics, bufSize),
	}
}

// SetPublicKey устанавливает RSA публичный ключ для шифрования запросов.
func (a *Agent) SetPublicKey(pub *rsa.PublicKey) {
	a.publicKey = pub
}

// SetRealIP устанавливает IP-адрес хоста агента для заголовка X-Real-IP.
func (a *Agent) SetRealIP(ip string) {
	a.realIP = ip
}

// SetGRPCClient включает отправку метрик через gRPC-клиент вместо HTTP.
func (a *Agent) SetGRPCClient(client pb.MetricsClient) {
	a.grpcClient = client
}

func (a *Agent) pushMetric(m models.Metrics) {
	a.Metrics <- m
}
