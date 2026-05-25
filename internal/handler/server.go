package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vrnvgasu/metrics/internal/config"
)

// Server — обертка над http.Server.
type Server struct {
	*http.Server
}

// NewServer создает сервер с адресом из конфига.
func NewServer(router *gin.Engine, cnf *config.ServerCnf) *Server {
	return &Server{
		Server: &http.Server{
			Addr:    cnf.Address,
			Handler: router,
		},
	}
}

// Run запускает HTTP-сервер.
func (s *Server) Run() error {
	if err := s.ListenAndServe(); err != nil {
		return fmt.Errorf("server Run: %w", err)
	}

	return nil
}
