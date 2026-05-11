package agent

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/vrnvgasu/metrics/internal/config"
	models "github.com/vrnvgasu/metrics/internal/model"
)

func Test_addPollCount(t *testing.T) {
	t.Parallel()

	a := NewAgent(nil, 100)
	require.Equal(t, len(a.Metrics), 0)

	a.addPollCount()
	require.Equal(t, len(a.Metrics), 1)

	var m models.Metrics
	select {
	case m = <-a.Metrics:
	default:
		t.Error("metrics channel closed")
	}

	var delta int64 = 1
	require.Equal(t, models.Metrics{
		ID:    pollCount,
		MType: models.Counter,
		Delta: &delta,
	}, m)
	require.Nil(t, m.Value)

	require.Equal(t, int64(1), a.pollCount.Load())
}

func Test_addRandomValue(t *testing.T) {
	t.Parallel()

	a := NewAgent(nil, 100)
	require.Equal(t, len(a.Metrics), 0)

	a.addRandomValue()
	require.Equal(t, len(a.Metrics), 1)

	var mActual models.Metrics
	select {
	case mActual = <-a.Metrics:
	default:
		t.Error("metrics channel closed")
	}

	require.Equal(t, randomValue, mActual.ID)
	require.Equal(t, models.Gauge, mActual.MType)
	require.NotNil(t, mActual.Value)
	require.Nil(t, mActual.Delta)
}

func Test_addStatsMetric(t *testing.T) {
	t.Parallel()

	var memStats runtime.MemStats
	a := NewAgent(nil, 100)
	require.Equal(t, 0, len(a.Metrics))

	a.addStatsMetric(memStats)
	require.Equal(t, len(gaugesMemStatNames), len(a.Metrics))

	actualNames := make(map[string]struct{}, len(a.Metrics))

loop:
	for {
		var m models.Metrics
		select {
		case m = <-a.Metrics:
		default:
			break loop
		}

		var existInList bool
		for _, name := range gaugesMemStatNames {
			if name == m.ID {
				existInList = true

				_, ok := actualNames[name]
				require.False(t, ok, fmt.Sprintf("name %s the name appeared twice", name))

				actualNames[name] = struct{}{}

				break
			}
		}

		require.True(t, existInList, fmt.Sprintf("metric.ID [%s] is not exist in gauges stat names", m.ID))
	}
}

func TestCollect(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
	defer cancel()

	a := NewAgent(nil, 100)
	require.Equal(t, 0, len(a.Metrics))

	err := a.Collect(ctx, &config.AgentCnf{PollInterval: 1})
	require.NoError(t, err)
	require.LessOrEqual(t, len(gaugesMemStatNames)+1+1, len(a.Metrics))
}

func BenchmarkAddStatsMetric(b *testing.B) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	a := NewAgent(nil, len(gaugesMemStatNames)+10)
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		a.addStatsMetric(memStats)
		// чтобы не блокировать следующую итерацию
		for len(a.Metrics) > 0 {
			<-a.Metrics
		}
	}
}

func BenchmarkAddRandomValue(b *testing.B) {
	a := NewAgent(nil, 10)
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		a.addRandomValue()
		<-a.Metrics
	}
}
