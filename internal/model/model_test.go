package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vrnvgasu/metrics/pkg/helper"
)

func TestNewMetricsFromStrings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		mType       string
		id          string
		value       string
		expected    Metrics
		expectedErr require.ErrorAssertionFunc
	}{
		{
			name:  "success: gauge",
			mType: string(Gauge),
			id:    "1",
			value: "0.1",
			expected: Metrics{
				ID:    "1",
				MType: "gauge",
				Value: helper.NewRefFloat64(0.1),
			},
			expectedErr: require.NoError,
		},
		{
			name:        "failed: gauge wrong value",
			mType:       string(Gauge),
			id:          "1",
			value:       "dummy",
			expected:    Metrics{},
			expectedErr: require.Error,
		},

		{
			name:  "success: counter",
			mType: string(Counter),
			id:    "1",
			value: "1",
			expected: Metrics{
				ID:    "1",
				MType: "counter",
				Delta: helper.NewRefInt64(1),
			},
			expectedErr: require.NoError,
		},
		{
			name:        "failed: counter incorrect value",
			mType:       string(Counter),
			id:          "1",
			value:       "0.1",
			expected:    Metrics{},
			expectedErr: require.Error,
		},
		{
			name:        "failed: counter wrong value",
			mType:       string(Counter),
			id:          "1",
			value:       "dummy",
			expected:    Metrics{},
			expectedErr: require.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m, err := NewMetricsFromStrings(tt.mType, tt.id, tt.value)
			tt.expectedErr(t, err)
			assert.Equal(t, tt.expected, m)
		})
	}
}
