package main

import (
	"log"
	"net/http"

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

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", h.Update)

	return http.ListenAndServe(":8080", mux)
}
