package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

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

	agentClient := agent.NewAgent(&http.Client{})
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
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
