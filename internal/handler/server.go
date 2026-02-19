package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vrnvgasu/metrics/internal/config"
)

type Server struct {
	*http.Server
}

func NewServer(router *gin.Engine, cnf *config.ServerCnf) *Server {
	return &Server{
		Server: &http.Server{
			Addr:    cnf.Address,
			Handler: router,
		},
	}
}

func (s *Server) Run() error {
	if err := s.ListenAndServe(); err != nil {
		return fmt.Errorf("server Run: %w", err)
	}

	return nil
}
