package service

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vrnvgasu/metrics/internal/config"
	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/internal/repository"
	"github.com/vrnvgasu/metrics/pkg/helper"
)

func TestStoreRestore(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		fileBody        string
		restore         bool
		expectedMetrics models.MetricsList
	}{
		{
			name:     "restore list",
			fileBody: `[{"id":"1","type":"gauge","value":842315.916}]`,
			restore:  true,
			expectedMetrics: models.MetricsList{
				{
					ID:    "1",
					MType: models.Gauge,
					Value: helper.NewRefFloat64(842315.916),
				},
			},
		},
		{
			name:            "restore empty",
			fileBody:        `[]`,
			restore:         true,
			expectedMetrics: models.MetricsList{},
		},
		{
			name:            "don't restore list",
			fileBody:        `[{"id":"1","type":"gauge","value":842315.916}]`,
			restore:         false,
			expectedMetrics: models.MetricsList{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fileName := tt.name + "_testRestore.json"

			file, err := os.Create(fileName)
			require.NoError(t, err)
			defer file.Close()

			_, err = file.Write([]byte(tt.fileBody))
			require.NoError(t, err)

			repo := repository.NewMemStorage()
			service, err := NewService(repo, config.ServerCnf{
				Restore:         tt.restore,
				FileStoragePath: fileName,
			})
			require.NoError(t, err)

			err = service.Restore()
			require.NoError(t, err)

			require.Equal(t, tt.expectedMetrics, repo.List())

			defer func() {
				service.Stop()
				os.Remove(fileName)
			}()
		})
	}
}
