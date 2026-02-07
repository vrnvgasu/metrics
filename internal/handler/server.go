package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	*http.Server
}

func NewServer(router *gin.Engine) *Server {
	return &Server{
		Server: &http.Server{
			Addr:    ":8080",
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
