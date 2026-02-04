package agent

import (
	"container/list"
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
	Metrics   list.List
	pollCount atomic.Int64
	mu        sync.Mutex
}

func NewAgent(client Client) *Agent {
	return &Agent{
		Client:  client,
		Metrics: list.List{},
	}
}

func (a *Agent) pushMetric(m models.Metrics) {
	a.mu.Lock()
	a.Metrics.PushBack(m)
	a.mu.Unlock()
}

func (a *Agent) pollMetric() *models.Metrics {
	a.mu.Lock()
	defer a.mu.Unlock()

	el := a.Metrics.Front()
	if el == nil {
		return nil
	}

	a.Metrics.Remove(el)

	m := el.Value.(models.Metrics)

	return &m
}
