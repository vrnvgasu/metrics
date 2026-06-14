package healthcheck

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vrnvgasu/metrics/internal/repository/mem"
)

func TestNewService(t *testing.T) {
	t.Parallel()

	storage := mem.NewMemStorage()
	s := NewService(storage)
	require.NotNil(t, s)
	assert.Equal(t, storage, s.DB)
}

func TestCheckPing(t *testing.T) {
	t.Parallel()

	storage := mem.NewMemStorage()
	s := NewService(storage)

	err := s.CheckPing(t.Context())
	assert.NoError(t, err)
}
