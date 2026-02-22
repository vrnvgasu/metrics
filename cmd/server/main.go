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

	"github.com/vrnvgasu/metrics/internal/handler"
	"github.com/vrnvgasu/metrics/internal/logger"
	"github.com/vrnvgasu/metrics/internal/repository"
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

	storage := repository.NewMemStorage()
	h := handler.NewHandler(storage)
	server := handler.NewServer(handler.NewRouter(h), cnf)

	serverErr := make(chan error)

	go func() {
		log.Println("starting server on: ", cnf.Address)
		if err := server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("shutting down server")
	case err := <-serverErr:
		return fmt.Errorf("server error: %w", err)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	return nil
}
