package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vrnvgasu/metrics/internal/config"
	grpcserver "github.com/vrnvgasu/metrics/internal/grpc"
	"github.com/vrnvgasu/metrics/internal/handler"
	"github.com/vrnvgasu/metrics/internal/logger"
	"github.com/vrnvgasu/metrics/internal/service/audit"
	"github.com/vrnvgasu/metrics/internal/service/healthcheck"
	"github.com/vrnvgasu/metrics/internal/service/metric"
	"github.com/vrnvgasu/metrics/internal/service/store"
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
	logger.Log.Infof("Build version: %s", buildVersion)
	logger.Log.Infof("Build date: %s", buildDate)
	logger.Log.Infof("Build commit: %s", buildCommit)
}

func run() error {
	cnf := parseFlags()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	err := logger.Initialize(cnf.LogLevel)
	if err != nil {
		return fmt.Errorf("could not initialize logger: %w", err)
	}

	info()

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
	r, err := handler.NewRouter(h)
	if err != nil {
		return fmt.Errorf("could not create router: %w", err)
	}
	server := handler.NewServer(r, cnf)

	storeService, err := store.NewService(storage, metricService, *cnf)
	if err != nil {
		return fmt.Errorf("could not create service: %w", err)
	}

	var grpcSrv *grpcserver.Server
	if cnf.GRPCAddress != "" {
		grpcSrv, err = grpcserver.NewServer(metricService, publisher, cnf)
		if err != nil {
			return fmt.Errorf("could not create gRPC server: %w", err)
		}
	}

	serverErr, err := start(ctx, cnf, server, grpcSrv, storeService)
	if err != nil {
		return err
	}

	return wait(ctx, serverErr, server, grpcSrv, storeService)
}

func start(
	ctx context.Context,
	cnf *config.ServerCnf,
	server *handler.Server,
	grpcSrv *grpcserver.Server,
	storeService *store.Service,
) (chan error, error) {
	serverErr := make(chan error, 3)

	if err := storeService.Restore(ctx); err != nil {
		return nil, fmt.Errorf("could not restore server: %w", err)
	}

	go func() {
		logger.Log.Info("starting pprof server on: localhost:6060")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			logger.Log.Warnf("pprof server: %v", err)
		}
	}()
	go func() {
		logger.Log.Infof("starting server on: %s", cnf.Address)
		if err := server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()
	if grpcSrv != nil {
		go func() {
			logger.Log.Infof("starting gRPC server on: %s", cnf.GRPCAddress)
			if err := grpcSrv.Run(); err != nil {
				serverErr <- err
			}
		}()
	}
	go func() {
		logger.Log.Infof("starting store data with interval: %d to file: %s", cnf.StoreInterval, cnf.FileStoragePath)
		if err := storeService.StoreInterval(ctx); err != nil {
			serverErr <- err
		}
	}()

	return serverErr, nil
}

func wait(
	ctx context.Context,
	serverErr chan error,
	server *handler.Server,
	grpcSrv *grpcserver.Server,
	storeService *store.Service,
) error {
	select {
	case <-ctx.Done():
		logger.Log.Info("shutting down server")
	case err := <-serverErr:
		return fmt.Errorf("server error: %w", err)
	}

	if grpcSrv != nil {
		grpcSrv.Stop()
	}

	if err := storeService.Stop(); err != nil {
		return fmt.Errorf("could not stop server: %w", err)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	return nil
}
