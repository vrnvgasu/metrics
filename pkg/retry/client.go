package retry

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// ClientConfig — конфигурация HTTP-клиента с поддержкой retry.
type ClientConfig struct {
	retryConfig *RetryConfig
}

// Client — HTTP-клиент с автоматическими повторными попытками при 5xx-ошибках.
type Client struct {
	*http.Client
	cnf *ClientConfig
}

// NewClient создает Client с переданным конфигом (nil — использует настройки по умолчанию).
func NewClient(cfg *ClientConfig) *Client {
	if cfg == nil {
		cfg = &ClientConfig{}
	}

	return &Client{
		Client: &http.Client{},
		cnf:    cfg,
	}
}

// Do выполняет HTTP-запрос с повторными попытками при временных ошибках и 5xx-ответах.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	var lastResp *http.Response
	err := RetryWithSettings(func() error {
		if bodyBytes != nil {
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		resp, err := c.Client.Do(req)
		if err != nil {
			return NewRetryableError(err)
		}

		lastResp = resp

		if resp.StatusCode < http.StatusInternalServerError {
			return nil
		}

		if c.shouldRetry(resp) {
			resp.Body.Close()

			return NewRetryableError(fmt.Errorf("status code: %d", resp.StatusCode))
		}

		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}, c.cnf.retryConfig)

	if err != nil {
		if lastResp != nil {
			lastResp.Body.Close()
		}

		return nil, fmt.Errorf("request failed: %v", err)
	}

	return lastResp, nil
}

func (c *Client) shouldRetry(resp *http.Response) bool {
	return resp.StatusCode >= http.StatusInternalServerError &&
		resp.StatusCode != http.StatusNotImplemented
}
