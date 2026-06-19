package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/vrnvgasu/metrics/internal/agent"
	"github.com/vrnvgasu/metrics/pkg/crypto"
	"github.com/vrnvgasu/metrics/pkg/retry"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func info() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}

func run() error {
	info()
	cnf := parseFlags()

	bufSize := cnf.ReportInterval / cnf.PollInterval * 100

	agentClient := agent.NewAgent(retry.NewClient(nil), bufSize)

	if cnf.CryptoKey != "" {
		pub, err := crypto.LoadPublicKey(cnf.CryptoKey)
		if err != nil {
			return fmt.Errorf("could not load public key: %w", err)
		}
		agentClient.SetPublicKey(pub)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
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
