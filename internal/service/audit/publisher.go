package audit

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/vrnvgasu/metrics/internal/logger"
)

type Event struct {
	observers []Observer
	mu        sync.Mutex
	semaphore chan struct{}
}

func NewEvent() *Event {
	return &Event{
		observers: make([]Observer, 0),
		semaphore: make(chan struct{}, 10),
	}
}

func (e *Event) Register(observer Observer) {
	e.mu.Lock()
	e.observers = append(e.observers, observer)
	e.mu.Unlock()
}

func (e *Event) Notify(ctx context.Context, metricIDList []string, ip string) {
	message := EventMessage{
		Ts:        time.Now().Unix(),
		Metrics:   metricIDList,
		IpAddress: ip,
	}

	go e.dispatch(ctx, message)
}

func (e *Event) dispatch(ctx context.Context, message EventMessage) {
	e.mu.Lock()
	observers := make([]Observer, len(e.observers))
	copy(observers, e.observers)
	e.mu.Unlock()

	for _, observer := range observers {
		select {
		case <-ctx.Done():
			return
		case e.semaphore <- struct{}{}:
		}

		go func() {
			defer func() { <-e.semaphore }()
			if err := observer.Observe(ctx, message); err != nil {
				logger.Log.Errorf("audit event Observe: %s", err.Error())
			}
		}()
	}
}

func (e *Event) Close() error {
	for _, observer := range e.observers {
		if err := observer.Close(); err != nil {
			return fmt.Errorf("audit event Close: %w", err)
		}
	}

	return nil
}
