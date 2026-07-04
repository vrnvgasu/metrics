package handler

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vrnvgasu/metrics/internal/config"
)

func TestNewServer(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	cnf := &config.ServerCnf{Address: "localhost:9999"}

	s := NewServer(r, cnf)
	require.NotNil(t, s)
	assert.Equal(t, "localhost:9999", s.Addr)
}

func TestServer_Run_Error(t *testing.T) {
	t.Parallel()

	// невалидный адрес → ListenAndServe вернет ошибку сразу
	srv := NewServer(gin.New(), &config.ServerCnf{Address: "localhost:-1"})
	err := srv.Run()
	require.Error(t, err)
}
