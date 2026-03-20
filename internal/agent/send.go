package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/vrnvgasu/metrics/internal/config"
	models "github.com/vrnvgasu/metrics/internal/model"
	serviceerrors "github.com/vrnvgasu/metrics/internal/service/errors"
	"github.com/vrnvgasu/metrics/pkg/compress"
)

const (
	path       = "updates"
	batchCount = 100
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

			err := serviceerrors.Retry(func() error {
				if err := a.SendMetric(batch, cnf.Address); err != nil {
					var urlErr *url.Error
					if errors.As(err, &urlErr) {
						return serviceerrors.NewRetryableError(err)
					}

					return err
				}

				return nil
			})
			if err != nil {
				return fmt.Errorf("sending metric: %w", err)
			}
		}
		time.Sleep(time.Duration(cnf.ReportInterval) * time.Second)
	}
}

func (a *Agent) SendMetric(m []*models.Metrics, address string) error {
	body, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("marshaling metric: %w", err)
	}

	cBody, err := compress.GzipCompress(body)
	if err != nil {
		return fmt.Errorf("compressing metric: %w", err)
	}

	updateURL := fmt.Sprintf("http://%s/%s", address, path)
	req, err := http.NewRequest(http.MethodPost, updateURL, bytes.NewBuffer(cBody))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	resp, err := a.Client.Do(req)
	if err != nil {
		return err
	}

	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}
	defer resp.Body.Close()

	return nil
}
