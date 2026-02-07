package main

import (
	"fmt"
	"log"

	"github.com/vrnvgasu/metrics/internal/handler"
	"github.com/vrnvgasu/metrics/internal/repository"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	storage := repository.NewMemStorage()
	h := handler.NewHandler(storage)
	server := handler.NewServer(handler.NewRouter(h))
	if err := server.Run(); err != nil {
		return fmt.Errorf("could not start server: %w", err)
	}

	return nil
}
