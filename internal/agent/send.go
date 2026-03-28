package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/vrnvgasu/metrics/internal/config"
	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/pkg/compress"
	"github.com/vrnvgasu/metrics/pkg/hash"
)

const (
	path       = "updates"
	batchCount = 100
	hashHeader = "HashSHA256"
)

func (a *Agent) SendMetrics(ctx context.Context, cnf *config.AgentCnf) error {
	batch := make([]*models.Metrics, 0, batchCount)

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

			batch = append(batch, m)
			if len(batch) < batchCount {
				continue
			}

			if err := a.SendMetric(batch, cnf); err != nil {
				return fmt.Errorf("agent.SendMetrics SendMetric: %w", err)
			}
		}
		time.Sleep(time.Duration(cnf.ReportInterval) * time.Second)
	}
}

func (a *Agent) SendMetric(m []*models.Metrics, cnf *config.AgentCnf) error {
	body, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("agent.SendMetric Marshal: %w", err)
	}

	cBody, err := compress.GzipCompress(body)
	if err != nil {
		return fmt.Errorf("agent.SendMetric GzipCompress: %w", err)
	}

	updateURL := fmt.Sprintf("http://%s/%s", cnf.Address, path)
	req, err := http.NewRequest(http.MethodPost, updateURL, bytes.NewBuffer(cBody))
	if err != nil {
		return fmt.Errorf("agent.SendMetric NewRequest: %w", err)
	}

	if err = a.SetHeaderHashSHA256(cnf.Key, body, req); err != nil {
		return fmt.Errorf("agent.SendMetric SetHeaderHashSHA256: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	resp, err := a.Client.Do(req)
	if err != nil {
		return fmt.Errorf("agent.SendMetric Do: %w", err)
	}

	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return fmt.Errorf("agent.SendMetric Copy: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func (a *Agent) SetHeaderHashSHA256(secret string, body []byte, r *http.Request) (err error) {
	if secret == "" {
		return nil
	}

	h, err := hash.PrepareHeaderHashSHA256(secret, body)
	if err != nil {
		return fmt.Errorf("agent.SetHeaderHashSHA256 PrepareHeaderHashSHA256: %w", err)
	}

	r.Header.Set(hashHeader, h)

	return
}
