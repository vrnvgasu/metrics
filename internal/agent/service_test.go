package agent

import (
	"testing"

	"github.com/stretchr/testify/require"

	models "github.com/vrnvgasu/metrics/internal/model"
)

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
