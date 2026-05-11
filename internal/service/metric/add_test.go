package metric

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/internal/repository"
	"github.com/vrnvgasu/metrics/internal/repository/mem"
	mockrepository "github.com/vrnvgasu/metrics/internal/repository/mocks"
	"github.com/vrnvgasu/metrics/pkg/helper"
)

func Test_createOrUpdate(t *testing.T) {
	t.Parallel()

	controller := gomock.NewController(t)
	defer controller.Finish()

	tests := []struct {
		name        string
		m           models.Metrics
		storage     func() repository.Storage
		expectedErr require.ErrorAssertionFunc
	}{
		{
			name: "no metrics; add gauge metric",
			m:    models.Metrics{ID: "1", MType: models.Gauge},
			storage: func() repository.Storage {
				storageMock := mockrepository.NewMockStorage(controller)
				storageMock.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

				return storageMock
			},
			expectedErr: require.NoError,
		},
		{
			name: "no metrics; add counter metric",
			m:    models.Metrics{ID: "1", MType: models.Counter},
			storage: func() repository.Storage {
				storageMock := mockrepository.NewMockStorage(controller)
				storageMock.EXPECT().GetByTypeAndID(gomock.Any(), models.Counter, "1").Return(nil, repository.ErrNotFound)
				storageMock.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

				return storageMock
			},
			expectedErr: require.NoError,
		},
		{
			name: "no metrics; add wrong metric",
			m:    models.Metrics{ID: "1", MType: "dummy"},
			storage: func() repository.Storage {
				storageMock := mockrepository.NewMockStorage(controller)

				return storageMock
			},
			expectedErr: require.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &Service{
				storage: tt.storage(),
			}

			err := s.createOrUpdate(t.Context(), []models.Metrics{tt.m})
			tt.expectedErr(t, err)
		})
	}
}

// benchPayload — реалистичный батч из 30 метрик (27 gauge + PollCount counter + 2 gopsutil)
var benchBatch = func() models.MetricsList {
	list := make(models.MetricsList, 0, 30)
	names := []string{
		"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys",
		"HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased",
		"HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys",
		"MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC",
		"NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys",
		"Sys", "TotalAlloc", "TotalMemory", "FreeMemory",
	}
	for _, name := range names {
		v := 1234567.0
		list = append(list, models.Metrics{ID: name, MType: models.Gauge, Value: &v})
	}
	delta := int64(42)
	list = append(list, models.Metrics{ID: "PollCount", MType: models.Counter, Delta: &delta})
	return list
}()

func BenchmarkCreateOrUpdate(b *testing.B) {
	s := NewService(mem.NewMemStorage())
	ctx := context.Background()
	// прогреть: добавить счётчик, чтобы покрыть путь addCounter с существующей записью
	_ = s.CreateOrUpdate(ctx, models.MetricsList{
		{ID: "PollCount", MType: models.Counter, Delta: helper.NewRefInt64(1)},
	})
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_ = s.CreateOrUpdate(ctx, benchBatch)
	}
}
