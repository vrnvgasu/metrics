package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/vrnvgasu/metrics/internal/config"
	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/pkg/compress"
	"github.com/vrnvgasu/metrics/pkg/crypto"
	"github.com/vrnvgasu/metrics/pkg/hash"
)

const (
	path            = "updates"
	batchCount      = 100
	hashHeader      = "HashSHA256"
	encryptedHeader = "X-Encrypted"
	realIPHeader    = "X-Real-IP"
)

func (a *Agent) SendMetrics(ctx context.Context, cnf *config.AgentCnf) error {
	chMetrics := make(chan *models.Metrics, batchCount)

	go func() {
		defer close(chMetrics)
		for m := range a.Metrics {
			chMetrics <- &m
		}
	}()

	workersCount := cnf.RateLimit
	if workersCount < 1 {
		workersCount = 1
	}

	errGroup, ctx := errgroup.WithContext(ctx)
	for i := 0; i < workersCount; i++ {
		errGroup.Go(func() error {
			batch := make([]*models.Metrics, 0, batchCount)
			for m := range chMetrics {
				batch = append(batch, m)
				if len(batch) < batchCount {
					continue
				}

				if err := a.sendBatch(batch, cnf); err != nil {
					return fmt.Errorf("agent.SendMetrics SendBatch: %w", err)
				}

				batch = make([]*models.Metrics, 0, batchCount)

				select {
				case <-ctx.Done():
				case <-time.After(time.Duration(cnf.ReportInterval) * time.Second):
				}
			}

			if len(batch) > 0 {
				return a.sendBatch(batch, cnf)
			}

			return nil
		})
	}

	if err := errGroup.Wait(); err != nil {
		return fmt.Errorf("agent.SendMetrics Wait: %w", err)
	}

	return nil
}

func (a *Agent) sendBatch(m []*models.Metrics, cnf *config.AgentCnf) error {
	body, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("agent.SendBatch Marshal: %w", err)
	}

	cBody, err := compress.GzipCompress(body)
	if err != nil {
		return fmt.Errorf("agent.SendBatch GzipCompress: %w", err)
	}

	sendBody := cBody
	if a.publicKey != nil {
		sendBody, err = crypto.Encrypt(a.publicKey, cBody)
		if err != nil {
			return fmt.Errorf("agent.SendBatch Encrypt: %w", err)
		}
	}

	updateURL := fmt.Sprintf("http://%s/%s", cnf.Address, path)
	req, err := http.NewRequest(http.MethodPost, updateURL, bytes.NewBuffer(sendBody))
	if err != nil {
		return fmt.Errorf("agent.SendBatch NewRequest: %w", err)
	}

	if err = a.SetHeaderHashSHA256(cnf.Key, body, req); err != nil {
		return fmt.Errorf("agent.SendBatch SetHeaderHashSHA256: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	if a.publicKey != nil {
		req.Header.Set(encryptedHeader, "true")
	}
	if a.realIP != "" {
		req.Header.Set(realIPHeader, a.realIP)
	}

	resp, err := a.Client.Do(req)
	if err != nil {
		return fmt.Errorf("agent.SendBatch Do: %w", err)
	}

	_, err = io.Copy(io.Discard, resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return fmt.Errorf("agent.SendBatch Copy: %w", err)
	}

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
