package agent

import (
	"io"
	"net/http"
	"sync"
	"sync/atomic"

	models "github.com/vrnvgasu/metrics/internal/model"
)

type Client interface {
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

type Agent struct {
	Client    Client
	Metrics   chan models.Metrics
	pollCount atomic.Int64
	mu        sync.Mutex
}

func NewAgent(client Client, bufSize int) *Agent {
	return &Agent{
		Client:  client,
		Metrics: make(chan models.Metrics, bufSize),
	}
}

func (a *Agent) pushMetric(m models.Metrics) {
	a.Metrics <- m
}

func (a *Agent) pollMetric() *models.Metrics {
	select {
	case m := <-a.Metrics:
		return &m
	default:
		return nil
	}
}
