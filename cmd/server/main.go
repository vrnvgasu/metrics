package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vrnvgasu/metrics/internal/config"
	"github.com/vrnvgasu/metrics/internal/handler"
	"github.com/vrnvgasu/metrics/internal/logger"
	"github.com/vrnvgasu/metrics/internal/service/audit"
	"github.com/vrnvgasu/metrics/internal/service/healthcheck"
	"github.com/vrnvgasu/metrics/internal/service/metric"
	"github.com/vrnvgasu/metrics/internal/service/store"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cnf := parseFlags()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	err := logger.Initialize(cnf.LogLevel)
	if err != nil {
		return fmt.Errorf("could not initialize logger: %w", err)
	}

	storage, err := initStorage(ctx, cnf)
	if err != nil {
		return fmt.Errorf("could not initialize storage: %w", err)
	}

	metricService := metric.NewService(storage)
	healthService := healthcheck.NewService(storage)

	publisher, err := audit.NewAudit(cnf)
	if err != nil {
		return fmt.Errorf("could not initialize audit: %w", err)
	}
	defer publisher.Close()

	h := handler.NewHandler(metricService, healthService, publisher, cnf)
	router := handler.NewServer(handler.NewRouter(h), cnf)

	storeService, err := store.NewService(storage, metricService, *cnf)
	if err != nil {
		return fmt.Errorf("could not create service: %w", err)
	}

	serverErr, err := start(ctx, cnf, router, storeService)
	if err != nil {
		return err
	}

	return wait(ctx, serverErr, router, storeService)
}

func start(
	ctx context.Context, cnf *config.ServerCnf, router *handler.Server, server *store.Service,
) (chan error, error) {
	serverErr := make(chan error)

	if err := server.Restore(ctx); err != nil {
		return nil, fmt.Errorf("could not restore server: %w", err)
	}

	go func() {
		logger.Log.Infof("starting router on: %s", cnf.Address)
		if err := router.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()
	go func() {
		logger.Log.Infof("starting store data with interval: %d to file: %s", cnf.StoreInterval, cnf.FileStoragePath)
		if err := server.StoreInterval(ctx); err != nil {
			serverErr <- err
		}
	}()

	return serverErr, nil
}

func wait(ctx context.Context, serverErr chan error, router *handler.Server, server *store.Service) error {
	select {
	case <-ctx.Done():
		logger.Log.Info("shutting down router")
	case err := <-serverErr:
		return fmt.Errorf("router error: %w", err)
	}

	if err := server.Stop(); err != nil {
		return fmt.Errorf("could not stop server: %w", err)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := router.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("router shutdown error: %w", err)
	}

	return nil
}
