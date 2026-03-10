package store

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/vrnvgasu/metrics/internal/config"
	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/internal/repository"
	"github.com/vrnvgasu/metrics/pkg/helper"
)

func TestStoreInterval(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name string
		repo func() (*repository.MemStorage, error)
	}{
		{
			name: "save list",
			repo: func() (*repository.MemStorage, error) {
				repo := repository.NewMemStorage()
				err := repo.Add(ctx, &models.Metrics{
					ID:    "1",
					MType: models.Gauge,
					Value: helper.NewRefFloat64(842315.916000),
				})
				if err != nil {
					return nil, err
				}

				return repo, nil
			},
		},
		{
			name: "save empty",
			repo: func() (*repository.MemStorage, error) {
				return repository.NewMemStorage(), nil
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fileName := tt.name + "_testStoreInterval.json"

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			repo, err := tt.repo()
			require.NoError(t, err)

			service, err := NewService(repo, config.ServerCnf{
				StoreInterval:   1,
				FileStoragePath: fileName,
			})
			require.NoError(t, err)

			service.StoreInterval(ctx)

			file, err := os.Open(fileName)
			require.NoError(t, err)
			defer file.Close()

			bytes, err := io.ReadAll(file)
			require.NoError(t, err)

			expected, err := json.Marshal(repo.List(ctx))
			require.NoError(t, err)
			require.JSONEq(t, string(expected), string(bytes))

			defer func() {
				service.Stop()
				os.Remove(fileName)
			}()
		})
	}
}
