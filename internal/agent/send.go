package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/vrnvgasu/metrics/internal/config"
	"github.com/vrnvgasu/metrics/internal/logger"
	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/pkg/compress"
)

const (
	path       = "updates"
	batchCount = 100
)

type ServerNotAvailableError struct {
	Err error
	Msg string
}

type ResponseError struct {
	Err error
	Msg string
}

func (e *ServerNotAvailableError) Error() string {
	return fmt.Sprintf("%s: %v", e.Msg, e.Err)
}

func (e *ServerNotAvailableError) Unwrap() error {
	return e.Err
}

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

			if err := a.SendMetric(batch, cnf.Address); err != nil {
				var availableErr *ServerNotAvailableError
				if !errors.As(err, &availableErr) {
					return fmt.Errorf("sending metric: %w", err)
				}

				logger.Log.Errorf("failed to send metrics: %v", err)
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

	url := fmt.Sprintf("http://%s/%s", address, path)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(cBody))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	resp, err := a.Client.Do(req)
	if err != nil {
		return &ServerNotAvailableError{
			Err: err,
			Msg: "sending metric:",
		}
	}

	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}
	defer resp.Body.Close()

	return nil
}
