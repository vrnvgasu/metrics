package audit

import (
	"context"
	"fmt"
	"time"
)

type Event struct {
	observers []Observer
}

func NewEvent() *Event {
	return &Event{
		observers: make([]Observer, 0),
	}
}

func (e *Event) Register(observer Observer) {
	e.observers = append(e.observers, observer)
}

func (e *Event) Notify(_ context.Context, metricIDList []string, ip string) error {
	message := EventMessage{
		Ts:        time.Now().Unix(),
		Metrics:   metricIDList,
		IpAddress: ip,
	}
	for _, observer := range e.observers {
		if err := observer.Observe(message); err != nil {
			return fmt.Errorf("audit event Observe: %w", err)
		}
	}

	return nil
}

func (e *Event) Close() error {
	for _, observer := range e.observers {
		if err := observer.Close(); err != nil {
			return fmt.Errorf("audit event Close: %w", err)
		}
	}

	return nil
}
