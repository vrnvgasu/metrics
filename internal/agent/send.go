package agent

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/vrnvgasu/metrics/internal/config"
	models "github.com/vrnvgasu/metrics/internal/model"
)

const (
	path = "update"
)

func (a *Agent) SendMetrics(ctx context.Context, cnf *config.AgentCnf) error {
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

			if err := a.SendMetric(*m, cnf.Address); err != nil {
				return fmt.Errorf("sending metric: %w", err)
			}
		}
		time.Sleep(time.Duration(cnf.ReportInterval) * time.Second)
	}
}

func (a *Agent) SendMetric(m models.Metrics, address string) error {
	url := fmt.Sprintf("http://%s/%s/%s/%s/%s", address, path, m.MType, m.ID, m.ValueToString())
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
