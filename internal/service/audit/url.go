package audit

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type URLAudit struct {
	client *resty.Client
}

func NewURLAudit(url string) *URLAudit {
	return &URLAudit{
		client: resty.New().SetBaseURL(url),
	}
}

func (a *URLAudit) Observe(e EventMessage) error {
	resp, err := a.client.R().SetBody(e).Post("")
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
