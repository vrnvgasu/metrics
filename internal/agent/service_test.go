package agent

import (
	"testing"

	"github.com/stretchr/testify/require"

	models "github.com/vrnvgasu/metrics/internal/model"
)

func TestAgent_pollMetric(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                    string
		metric                  []models.Metrics
		expectedMetric          *models.Metrics
		expectedRemainingLength int
	}{
		{
			name: "poll from 3 elements",
			metric: []models.Metrics{
				{
					ID:    "1",
					MType: "type1",
					Delta: nil,
					Value: func() *float64 {
						v := 1.0
						return &v
					}(),
				},
				{
					ID:    "2",
					MType: "type2",
					Delta: nil,
					Value: func() *float64 {
						v := 2.0
						return &v
					}(),
				},
				{
					ID:    "3",
					MType: "type3",
					Delta: nil,
					Value: func() *float64 {
						v := 3.0
						return &v
					}(),
				},
			},
			expectedMetric: &models.Metrics{
				ID:    "1",
				MType: "type1",
				Delta: nil,
				Value: func() *float64 {
					v := 1.0
					return &v
				}(),
			},
			expectedRemainingLength: 2,
		},
		{
			name: "poll from 1 elements",
			metric: []models.Metrics{
				{
					ID:    "1",
					MType: "type1",
					Delta: nil,
					Value: func() *float64 {
						v := 1.0
						return &v
					}(),
				},
			},
			expectedMetric: &models.Metrics{
				ID:    "1",
				MType: "type1",
				Delta: nil,
				Value: func() *float64 {
					v := 1.0
					return &v
				}(),
			},
			expectedRemainingLength: 0,
		},
		{
			name:                    "poll from 0 elements",
			metric:                  []models.Metrics{},
			expectedMetric:          nil,
			expectedRemainingLength: 0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			a := NewAgent(nil, 100)
			for _, m := range tt.metric {
				a.Metrics <- m
			}

			require.Equal(t, tt.expectedMetric, a.pollMetric())
			require.Equal(t, tt.expectedRemainingLength, len(a.Metrics))
		})
	}
}

func TestAgent_pushMetric(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		metric         []models.Metrics
		expectedLength int
	}{
		{
			name:           "add 1",
			metric:         []models.Metrics{{ID: "1"}},
			expectedLength: 1,
		},
		{
			name:           "add 2",
			metric:         []models.Metrics{{ID: "1"}, {ID: "2"}},
			expectedLength: 2,
		},
		{
			name:           "add 3",
			metric:         []models.Metrics{{ID: "1"}, {ID: "2"}, {ID: "3"}},
			expectedLength: 3,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			a := NewAgent(nil, 100)
			for _, m := range tt.metric {
				a.pushMetric(m)
			}

			require.Equal(t, len(a.Metrics), tt.expectedLength)
		})
	}
}
