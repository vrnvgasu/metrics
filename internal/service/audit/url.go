package audit

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type URLAudit struct {
	client *resty.Client
}

func NewURLAudit(url string) *URLAudit {
	return &URLAudit{
		client: resty.New().
			SetBaseURL(url).
			SetRetryCount(3).
			SetRetryWaitTime(30 * time.Second).
			SetRetryMaxWaitTime(90 * time.Second),
	}
}

func (a *URLAudit) Observe(ctx context.Context, e EventMessage) error {
	resp, err := a.client.R().SetContext(ctx).SetBody(e).Post("")
	if err != nil {
		return fmt.Errorf("urlAudit.Sub Post: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("urlAudit.Sub Post: status %d", resp.StatusCode())
	}

	return nil
}

func (a *URLAudit) Close() error {
	return nil
}
