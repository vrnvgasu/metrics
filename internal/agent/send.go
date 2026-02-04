package agent

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	models "github.com/vrnvgasu/metrics/internal/model"
)

const (
	path = "http://localhost:8080/update"
)

func (a *Agent) SendMetrics(ctx context.Context, interval time.Duration) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		for {
			m := a.pollMetric()
			if m == nil {
				break
			}

			if err := a.SendMetric(*m); err != nil {
				return fmt.Errorf("sending metric: %w", err)
			}
		}
		time.Sleep(interval * time.Second)
	}
}

func (a *Agent) SendMetric(m models.Metrics) error {
	url := fmt.Sprintf("%s/%s/%s/%s", path, m.MType, m.ID, m.ValueToString())
	resp, err := a.Client.Post(url, "text/plain", http.NoBody)
	if err != nil {
		return fmt.Errorf("sending metric: %w", err)
	}

	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}
	defer resp.Body.Close()

	return nil
}
