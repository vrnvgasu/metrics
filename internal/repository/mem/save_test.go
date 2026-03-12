package mem

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/pkg/helper"
)

func TestMemStorageAdd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		metrics         map[string]models.Metrics
		m               models.Metrics
		expectedMetrics map[string]models.Metrics
		expectedErr     require.ErrorAssertionFunc
	}{
		{
			name:    "no metrics; add gauge metric",
			metrics: map[string]models.Metrics{},
			m:       models.Metrics{ID: "1", MType: models.Gauge},
			expectedMetrics: map[string]models.Metrics{
				"1": {ID: "1", MType: models.Gauge},
			},
			expectedErr: require.NoError,
		},
		{
			name:    "no metrics; add counter metric",
			metrics: map[string]models.Metrics{},
			m:       models.Metrics{ID: "1", MType: models.Counter},
			expectedMetrics: map[string]models.Metrics{
				"1": {ID: "1", MType: models.Counter},
			},
			expectedErr: require.NoError,
		},
		{
			name:            "no metrics; add wrong metric",
			metrics:         map[string]models.Metrics{},
			m:               models.Metrics{ID: "1", MType: "dummy"},
			expectedMetrics: map[string]models.Metrics{},
			expectedErr:     require.Error,
		},
		{
			name: "has metrics; rewrite gauge metric",
			metrics: map[string]models.Metrics{
				"1": {
					ID:    "1",
					MType: models.Gauge,
					Value: helper.NewRefFloat64(1),
				},
				"11": {
					ID:    "11",
					MType: models.Gauge,
					Value: helper.NewRefFloat64(11),
				},
			},
			m: models.Metrics{
				ID:    "1",
				MType: models.Gauge,
				Value: helper.NewRefFloat64(2),
			},
			expectedMetrics: map[string]models.Metrics{
				"1": {
					ID:    "1",
					MType: models.Gauge,
					Value: helper.NewRefFloat64(2),
				},
				"11": {
					ID:    "11",
					MType: models.Gauge,
					Value: helper.NewRefFloat64(11),
				},
			},
			expectedErr: require.NoError,
		},
		{
			name: "has metrics; rewrite counter metric",
			metrics: map[string]models.Metrics{
				"1": {
					ID:    "1",
					MType: models.Counter,
					Delta: helper.NewRefInt64(1),
				},
				"11": {
					ID:    "11",
					MType: models.Counter,
					Delta: helper.NewRefInt64(11),
				},
			},
			m: models.Metrics{
				ID:    "1",
				MType: models.Counter,
				Delta: helper.NewRefInt64(2),
			},
			expectedMetrics: map[string]models.Metrics{
				"1": {
					ID:    "1",
					MType: models.Counter,
					Delta: helper.NewRefInt64(3),
				},
				"11": {
					ID:    "11",
					MType: models.Counter,
					Delta: helper.NewRefInt64(11),
				},
			},
			expectedErr: require.NoError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &MemStorage{
				metrics: tt.metrics,
			}

			err := s.Add(context.Background(), &tt.m)
			tt.expectedErr(t, err)
			for k, v := range tt.expectedMetrics {
				m, ok := s.metrics[k]
				require.True(t, ok)
				require.Equal(t, v, m)
			}
		})
	}
}
