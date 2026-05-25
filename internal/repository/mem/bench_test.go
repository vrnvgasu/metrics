package mem

import (
	"context"
	"fmt"
	"testing"

	models "github.com/vrnvgasu/metrics/internal/model"
	"github.com/vrnvgasu/metrics/pkg/helper"
)

func newFilledStorage(n int) *Storage {
	s := NewMemStorage()
	ctx := context.Background()
	for i := range n {
		v := float64(i)
		_ = s.Save(ctx, &models.Metrics{
			ID:    fmt.Sprintf("gauge_%d", i),
			MType: models.Gauge,
			Value: helper.NewRefFloat64(v),
		})
	}
	return s
}

func BenchmarkStorageSave(b *testing.B) {
	s := newFilledStorage(30)
	ctx := context.Background()
	m := &models.Metrics{
		ID:    "Alloc",
		MType: models.Gauge,
		Value: helper.NewRefFloat64(1234567.0),
	}
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_ = s.Save(ctx, m)
	}
}

func BenchmarkStorageList(b *testing.B) {
	s := newFilledStorage(30)
	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_, _ = s.List(ctx)
	}
}

func BenchmarkStorageGetByTypeAndID(b *testing.B) {
	s := newFilledStorage(30)
	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_, _ = s.GetByTypeAndID(ctx, models.Gauge, "gauge_15")
	}
}
