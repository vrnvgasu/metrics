package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/vrnvgasu/metrics/internal/config"
	"github.com/vrnvgasu/metrics/internal/logger"
	models "github.com/vrnvgasu/metrics/internal/model"
)

const (
	path = "update"
)

type ServerNotAvailableError struct {
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

func (a *Agent) SendMetric(m models.Metrics, address string) error {
	body, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("marshaling metric: %w", err)
	}

	url := fmt.Sprintf("http://%s/%s", address, path)
	resp, err := a.Client.Post(url, "text/plain", bytes.NewBuffer(body))
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
