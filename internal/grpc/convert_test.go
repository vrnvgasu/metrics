package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	models "github.com/vrnvgasu/metrics/internal/model"
	pb "github.com/vrnvgasu/metrics/internal/proto"
)

func refInt64(v int64) *int64       { return &v }
func refFloat64(v float64) *float64 { return &v }

func TestToModels(t *testing.T) {
	t.Parallel()

	in := []*pb.Metric{
		pb.Metric_builder{Id: "PollCount", Type: pb.Metric_COUNTER, Delta: 5}.Build(),
		pb.Metric_builder{Id: "Alloc", Type: pb.Metric_GAUGE, Value: 3.14}.Build(),
	}

	got := ToModels(in)

	require.Len(t, got, 2)

	assert.Equal(t, "PollCount", got[0].ID)
	assert.Equal(t, models.Counter, got[0].MType)
	require.NotNil(t, got[0].Delta)
	assert.Equal(t, int64(5), *got[0].Delta)
	assert.Nil(t, got[0].Value)

	assert.Equal(t, "Alloc", got[1].ID)
	assert.Equal(t, models.Gauge, got[1].MType)
	require.NotNil(t, got[1].Value)
	assert.InDelta(t, 3.14, *got[1].Value, 1e-9)
	assert.Nil(t, got[1].Delta)
}

func TestToProto(t *testing.T) {
	t.Parallel()

	in := []*models.Metrics{
		{ID: "PollCount", MType: models.Counter, Delta: refInt64(7)},
		{ID: "Alloc", MType: models.Gauge, Value: refFloat64(2.5)},
	}

	got := ToProto(in)

	require.Len(t, got, 2)

	assert.Equal(t, "PollCount", got[0].GetId())
	assert.Equal(t, pb.Metric_COUNTER, got[0].GetType())
	assert.Equal(t, int64(7), got[0].GetDelta())

	assert.Equal(t, "Alloc", got[1].GetId())
	assert.Equal(t, pb.Metric_GAUGE, got[1].GetType())
	assert.InDelta(t, 2.5, got[1].GetValue(), 1e-9)
}

// TestToModels_ToProto_RoundTrip проверяет, что конвертация туда-обратно сохраняет данные.
func TestToModels_ToProto_RoundTrip(t *testing.T) {
	t.Parallel()

	src := []*models.Metrics{
		{ID: "PollCount", MType: models.Counter, Delta: refInt64(42)},
		{ID: "Alloc", MType: models.Gauge, Value: refFloat64(1.5)},
	}

	got := ToModels(ToProto(src))

	require.Len(t, got, 2)
	assert.Equal(t, int64(42), *got[0].Delta)
	assert.InDelta(t, 1.5, *got[1].Value, 1e-9)
}
