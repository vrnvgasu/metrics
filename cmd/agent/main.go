package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/vrnvgasu/metrics/internal/agent"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cnf := parseFlags()

	bufSize := cnf.ReportInterval / cnf.PollInterval * 100

	agentClient := agent.NewAgent(&http.Client{}, bufSize)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	errGroup, runtimeCtx := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		return agentClient.Collect(runtimeCtx, cnf)
	})
	errGroup.Go(func() error {
		return agentClient.SendMetrics(runtimeCtx, cnf)
	})

	if err := errGroup.Wait(); err != nil {
		return err
	}

	return nil
}
