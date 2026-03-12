package mem

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/pkg/helper"
)

func TestGetByTypeAndID(t *testing.T) {
	t.Parallel()

	mMapDefault := models.MetricsMap{
		"1": {
			ID:    "1",
			MType: models.Gauge,
			Value: helper.NewRefFloat64(1.1),
		},
		"2": {
			ID:    "2",
			MType: models.Counter,
			Delta: helper.NewRefInt64(11),
		},
	}

	tests := []struct {
		name      string
		mMap      models.MetricsMap
		mtype     string
		id        string
		expectedM models.Metrics
		err       require.ErrorAssertionFunc
	}{
		{
			name:  "find gauge",
			mMap:  mMapDefault,
			mtype: models.Gauge,
			id:    "1",
			expectedM: models.Metrics{
				ID:    "1",
				MType: models.Gauge,
				Value: helper.NewRefFloat64(1.1),
			},
			err: require.NoError,
		},
		{
			name:  "find counter",
			mMap:  mMapDefault,
			mtype: models.Counter,
			id:    "2",
			expectedM: models.Metrics{
				ID:    "2",
				MType: models.Counter,
				Delta: helper.NewRefInt64(11),
			},
			err: require.NoError,
		},
		{
			name:      "type is not found",
			mMap:      mMapDefault,
			mtype:     "dummy",
			id:        "2",
			expectedM: models.Metrics{},
			err:       require.Error,
		},
		{
			name:      "type id not found",
			mMap:      mMapDefault,
			mtype:     models.Gauge,
			id:        "dummy",
			expectedM: models.Metrics{},
			err:       require.Error,
		},
		{
			name:      "storage is empty",
			mMap:      map[string]models.Metrics{},
			mtype:     models.Gauge,
			id:        "1",
			expectedM: models.Metrics{},
			err:       require.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &MemStorage{
				metrics: tt.mMap,
			}
			m, err := s.GetByTypeAndID(context.Background(), tt.mtype, tt.id)
			tt.err(t, err)
			if err == nil {
				require.Equal(t, tt.expectedM, *m)
			} else {
				require.Nil(t, m)
			}
		})
	}
}
