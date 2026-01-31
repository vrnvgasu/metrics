package main

import (
	"log"
	"net/http"

	"github.com/vrnvgasu/metrics/internal/handler/server"
	"github.com/vrnvgasu/metrics/internal/repository"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	repository.Storage = repository.NewMemStorage()

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", server.Update)

	return http.ListenAndServe(":8080", mux)
}
