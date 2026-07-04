package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vrnvgasu/metrics/internal/config"
)

func TestInitStorage_Mem(t *testing.T) {
	t.Parallel()

	// без DatabaseDSN — должно вернуться in-memory хранилище
	storage, err := initStorage(context.Background(), &config.ServerCnf{})
	require.NoError(t, err)
	require.NotNil(t, storage)
}
