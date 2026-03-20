package mem

import (
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
	}{
		{
			name:    "no metrics; add metric",
			metrics: map[string]models.Metrics{},
			m:       models.Metrics{ID: "1", MType: "someMetric"},
			expectedMetrics: map[string]models.Metrics{
				"1": {ID: "1", MType: "someMetric"},
			},
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
					Delta: helper.NewRefInt64(2),
				},
				"11": {
					ID:    "11",
					MType: models.Counter,
					Delta: helper.NewRefInt64(11),
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &Storage{
				metrics: tt.metrics,
			}

			err := s.Save(t.Context(), &tt.m)
			require.NoError(t, err)
			for k, v := range tt.expectedMetrics {
				m, ok := s.metrics[k]
				require.True(t, ok)
				require.Equal(t, v, m)
			}
		})
	}
}
